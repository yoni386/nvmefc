package target

import (
	"io/ioutil"
	"log"
	"bytes"
	"os"
	"path/filepath"
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"github.com/kataras/golog"
)

const (
	procPath = "/proc/cmdline"
)

const KHEX uint64 = 0x400

const (
	B  uint64 = 0x80000000
	KB uint64 = 0x200000
	MB uint64 = 0x800
	GB uint64 = 0x2
)

var (
	byteTrue = []byte{49, 10}
	byteZero = []byte{48}
	ByteZero = []byte{48}
	byteOne  = []byte{49}
	ByteOne  = []byte{49}
)

type Subsystem struct {
	name       string
	path       string
	namespaces []Namespace
	ports      []*Port
	offload     bool
}

func (s *Subsystem) Ports() []*Port {
	return s.ports
}

func newSubsystem(name string, path string) Subsystem {
	//ports := make([]*Port, 1, 1)
	//return Subsystem{name: name, path: path, ports: ports}
	offload := checkOffload(path)
	return Subsystem{name: name, path: path, offload: offload}
}

func Kernel() uint64 {
	file, err := ioutil.ReadFile(procPath)
	if err != nil {
		//log.Println(err)
		//golog.SetTimeFormat("03/01/2006 15:04")

		golog.Info("Kernel() This is an error ", err)

		//return -1
	}

	var hex uint64

	re := regexp.MustCompile("[0-9]+")
	procCmdline := bytes.Split(file, []byte(" "))

	for i := range procCmdline {
		if bytes.Contains(procCmdline[i], []byte("mem")) {
			s := string(procCmdline[i])
			memVal, err := strconv.ParseUint(re.FindString(s), 10, 0) // check if base 10 it right might be better to have 16
			if err != nil {
				//fmt.Println(err)
				//golog.SetTimeFormat("03/01/2006 15:04")

				golog.Info("This is an error 1111", err)
				//return 0
			}

			key := regexp.MustCompile("[A-Z]+").FindString(s)
			switch key {
			case "B":
				hex = memVal + B
			case "K":
				hex = (memVal + KB) * KHEX
			case "M":
				hex = (memVal + MB) * KHEX * KHEX
			case "G":
				hex = (memVal + GB) * KHEX * KHEX * KHEX

			default:
				fmt.Println("key is not known", key)
			}
		}
	}

	return hex
}

func (s *Subsystem) Namespaces() []Namespace {
	return s.namespaces
}

func checkOffload(path string) bool {
	var offload bool
	file, err := ioutil.ReadFile(path + "/attr_offload")
	if err != nil {
		log.Println(err)
	}
	if bytes.Equal(file, byteTrue) {
		offload = true
	}
	return offload
}

func setAttr(path string, byteValue []byte) error {
	if err := ioutil.WriteFile(path, byteValue, 0777); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func CreateSubsystem(name string, offload bool) error {
	if err := os.Mkdir(SubsystemsPath + name, 0777); err != nil {
		//log.Println(err)
		return err
	}
	if offload {
		//setAttr(SubsystemsPath + name + "/attr_offload", byteOne)
		if err := setAttr(SubsystemsPath + name + "/attr_offload", byteOne); err != nil {
			//log.Println(err)
			return err
		}
	}
	//setAttr(SubsystemsPath + name + "/attr_allow_any_host", byteOne)

	if err := setAttr(SubsystemsPath + name + "/attr_allow_any_host", byteOne); err != nil {
		//log.Println(err)
		return err
	}
	return nil
}

func SetSubsystemOffload(name string, enable bool) error {
	if enable {
		if err := setAttr(SubsystemsPath + name + "/attr_offload", byteOne); err != nil {
			return err
		}
		return nil
	} else {
		if err := setAttr(SubsystemsPath + name + "/attr_offload", ByteZero); err != nil {
			return err
		}

		return nil
	}
}

func (s *Subsystem) getNamespaces() {
	paths, err := filepath.Glob(fmt.Sprintf("%s/%s/*", s.path, "namespaces"))
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		name := filepath.Base(path)
		s.namespaces = append(s.namespaces, Namespace{name: name, path: path})
	}
}

func (s *Subsystem) scan(port *Port) {
	s.getNamespaces()
	s.getDevices()
	s.get_subsystem_ports_read_soft_link(port)
}

func (s *Subsystem) Offload() bool {
	return s.offload
}

func (s *Subsystem) IsOffload() string {
	if s.Offload() {
		return "true"
	}
	return "false"
}

func (s *Subsystem) get_subsystem_ports_read_soft_link(port *Port) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/subsystems/*", port.path))
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		pathLink, err := os.Readlink(path)
		if err != nil {
			log.Fatal(err)
		}

		subName := filepath.Base(pathLink)
		if subName == s.name {
			//s.addPort(&Port{})
			s.ports = append(s.ports, port)
		}
	}
}

func (s *Subsystem) getIpReport() string {
	var ip string
	for _, port := range s.ports {
		ip += port.AddrTraddr()
	}
	return ip
}

func (s *Subsystem) addPort(port ...*Port) {
	s.ports = append(s.ports, port...)
}

func (s *Subsystem) getPortReport() string {
	var port string
	for _, ip := range s.ports {
		port += ip.addrTrsvcid
	}
	return port
}

func (s *Subsystem) getNamespacesReport() string {
	var namespaceName string

	for _, namespace := range s.namespaces {
		namespaceName += namespace.name + "\n"
	}
	return namespaceName
}

func (s *Subsystem) getNamespacesDevicePathReport() string {
	var devicePath string

	for _, namespace := range s.namespaces {
		devicePath += namespace.devicePath
	}
	return devicePath
}

func (s *Subsystem) getNamespacesPciDevicePathReport() string {
	var pciDevicePath string

	for _, namespace := range s.namespaces {
		pciDevicePath += namespace.pciDevicePath + "\n"
	}
	return pciDevicePath
}

func (s *Subsystem) getNamespacesDeviceEnableReport() string {
	var enable string

	for _, namespace := range s.namespaces {
		if namespace.Enable() {
			enable += "true\n"
			return enable
		}
		enable += "false\n"
	}

	return enable
}

func (s *Subsystem) getReport() []string {

	return []string{s.name, s.IsOffload(), s.getNamespacesReport(), s.getNamespacesDevicePathReport(), s.getNamespacesPciDevicePathReport(), s.getNamespacesDeviceEnableReport(), s.getIpReport(), s.getPortReport()}
	//return []string{s.name, s.path, s.getIpReport(), s.getPortReport(),}
}

func getPciDevicePath(name string) string {
	//name = strings.TrimSuffix(name, "\n")
	pciDevicePath, err := os.Readlink(fmt.Sprintf("/sys/block/%s/device/device", name))
	if err != nil {
		log.Fatal(err)
	}
	pciDevicePath = filepath.Base(pciDevicePath)
	return pciDevicePath
}


func (s *Subsystem) getDevices() {
	var devicePath string
	var pciDevicePath string
	var enable bool
	//var offload bool

	for i := range s.namespaces {
		paths, err := filepath.Glob(s.namespaces[i].path + "/*")
		if err != nil {
			log.Println(err)
		}
		for _, path := range paths {
			switch filepath.Base(path) {
			case "device_path":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				devicePath = string(file)
				pciDevicePath = getPciDevicePath(filepath.Base(strings.TrimSuffix(devicePath, "\n")))
			case "enable":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				//if string(file) == "1" {
				if bytes.Equal(file, byteTrue) {
					enable = true
				}

			default:
				continue
			}
			//fmt.Println(filepath.Dir(path))

			if len(devicePath) > 0 {
				s.namespaces[i].Device.devicePath = devicePath
				s.namespaces[i].Device.pciDevicePath = pciDevicePath
				s.namespaces[i].Device.enable = enable
				s.namespaces[i].Device.path = filepath.Dir(path)
			}
		}

	}
}

func (s *Subsystem) SetDeviceSysFsEnable(enable bool) {
	for i := range s.Namespaces() {
		s.Namespaces()[i].SetDeviceSysFsEnable(enable)
	}
}

// separate those two funcs to delete and update slice
func (s *Subsystem) destroyPorts() {
	for i := range s.ports {
		//if s.ports[i].IsDisabled() {
			fmt.Println("destroy", i, s.ports[i].AddrTrsvcid())
			s.ports[i].destroyPort()
		//}
	}

	//if err := os.Remove(fmt.Sprintf("%s/subsystems/%s", s.ports[i].path, s.name)); err != nil {
	//	log.Fatal(err)
		//	return
		//}
		//s.ports[i].SetState(false) // s.ports[i].Disable()??


	for i := len(s.ports) - 1; i >= 0; i-- {
		if s.ports[i].IsDisabled() {
			s.ports[i] = nil
			s.ports = append(s.ports[:i], s.ports[i+1:]...)
		}
	}
}

func (s *Subsystem) LinkSubsystemPorts() {
	for i := range s.ports {
		//if s.ports[i].IsDisabled() {
		//fmt.Println("link", i, s.name, s.ports[i].AddrTrsvcid())
		portID := strings.Trim(s.ports[i].AddrTrsvcid(), "\n")


		if err := LinkPort(s.name, portID); err != nil {
			log.Println(err)
			//return err
			return
		}

		if s.ports[i].IsDisabled() {
			s.ports[i].SetState(true)
		}

		//}
	}

	//if err := os.Remove(fmt.Sprintf("%s/subsystems/%s", s.ports[i].path, s.name)); err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//s.ports[i].SetState(false) // s.ports[i].Disable()??


	//for i := len(s.ports) - 1; i >= 0; i-- {
	//	if s.ports[i].IsDisabled() {
	//		s.ports[i] = nil
	//		s.ports = append(s.ports[:i], s.ports[i+1:]...)
	//	}
	//}
}


func (s *Subsystem) UnlinkSubsystemsPorts() {
	for i := range s.ports {
		//if s.ports[i].IsDisabled() {
		//fmt.Println("UnlinkSubsystemsPorts", i, s.ports[i].AddrTrsvcid())
		s.ports[i].UnLinkPort(s.name)
		//}
	}

	//if err := os.Remove(fmt.Sprintf("%s/subsystems/%s", s.ports[i].path, s.name)); err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//s.ports[i].SetState(false) // s.ports[i].Disable()??


}

func (s *Subsystem) destroyDisks() {
	for i := range s.namespaces {
		paths, err := filepath.Glob(s.namespaces[i].path + "/*")
		if err != nil {
			log.Println(err)
		}
		for _, path := range paths {
			//fmt.Printf("path: %s\n", path)
			if filepath.Base(path) == "enable" {
				if err := ioutil.WriteFile(path, byteZero, 0777); err != nil {
					log.Fatal(err)
				}
				s.namespaces[i].enable = false
			}
		}
	}
}

func (s *Subsystem) destroyNameSpace() {

	for i := range s.namespaces {
		if err := os.Remove(s.namespaces[i].path); err != nil {
			log.Fatal(err)
		}
	}

	for i := len(s.namespaces) - 1; i >= 0; i-- {
		if s.namespaces[i].IsDisabled() {
			//s.namespaces[i] = nil
			s.namespaces = append(s.namespaces[:i], s.namespaces[i+1:]...)
		}
	}

}
func (s *Subsystem) destroySubsystem() {
	if err := os.Remove(s.path); err != nil {
		log.Fatal(err)
	}
}