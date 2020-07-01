package target

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/olekukonko/tablewriter"
)

const (
	SysFsNvmeGlob  = "/dev/nvme*n?"
	PortsPath      = "/sys/kernel/config/nvmet/ports/"
	SubsystemsPath = "/sys/kernel/config/nvmet/subsystems/"
)

const (
	ibUmadLoadCmd    = "modprobe ib_umad"
	mlx5IBLoadCmd    = "modprobe mlx5_ib"
	nvmetLoadCmd     = "modprobe nvmet"
	nvmetRDMALoadCmd = "modprobe nvmet_rdma"
)

type Target struct {
	subsystems []Subsystem
	ports      []Port
	disks      []Disk
}

func (t *Target) Subsystems() []Subsystem {
	return t.subsystems
}

func (t *Target) Ports() []Port {
	return t.ports
}

func (t *Target) PortsPointer() *[]Port {
	return &t.ports
}

func (t *Target) Disks() []Disk {
	return t.disks
}

func (t *Target) DisksPointer() *[]Disk {
	return &t.disks
}

func New() Target {
	return Target{
		[]Subsystem{},
		[]Port{},
		[]Disk{},
	}
}

func (t *Target) getDisks() {

	paths, err := filepath.Glob(SysFsNvmeGlob)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, path := range paths {
		name := filepath.Base(path)
		fmt.Println(name)

		pciDevicePath, err := os.Readlink(fmt.Sprintf("/sys/block/%s/device/device", name))
		if err != nil {
			log.Fatal(err)
		}
		pciDevicePath = filepath.Base(pciDevicePath)
		//fmt.Println(pciDevicePath)

		disk := Disk{name, path, pciDevicePath}
		t.disks = append(t.disks, disk)
	}
}

func (t *Target) GetPort(name string) {
	paths, err := filepath.Glob(PortsPath + name)
	if err != nil {
		log.Fatal(err)
		return
	}
	var addrTraddr string
	var addrTrsvcid string
	for _, path := range paths {
		paths, err := filepath.Glob(path + "/*")
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, path := range paths {

			switch filepath.Base(path) {
			case "addr_traddr":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				addrTraddr = string(file)
			case "addr_trsvcid":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				addrTrsvcid = string(file)
			default:
				continue
			}
		}
		t.ports = append(t.ports, NewPort(addrTraddr, addrTrsvcid, path, true))
	}
}

func (t *Target) getPorts() {

	paths, err := filepath.Glob(PortsPath + "/*")
	if err != nil {
		log.Fatal(err)
		return
	}
	var addrTraddr string
	var addrTrsvcid string
	for _, path := range paths {
		//fmt.Printf("path %s \n", path)

		paths, err := filepath.Glob(path + "/*")
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, path := range paths {
			//fmt.Printf("path %s \n", path)

			switch filepath.Base(path) {
			case "addr_traddr":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				addrTraddr = string(file)
			case "addr_trsvcid":
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return
				}
				addrTrsvcid = string(file)
			default:
				continue
			}

			//path = filepath.Dir(path) // check if probably not
		}
		t.ports = append(t.ports, NewPort(addrTraddr, addrTrsvcid, path, true))
	}
}

func (t *Target) getSubsystemsPorts() {
	for idxPort := range t.ports { // might need to pass t.ports instead of two loops move loop to get_subsystem_ports_read_soft_link
		for i := range t.subsystems {
			t.subsystems[i].get_subsystem_ports_read_soft_link(&t.ports[idxPort])
		}
	}
}

func (t *Target) getSubsystems() {
	paths, err := filepath.Glob(SubsystemsPath + "/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		//for _, path := range []string{"1", "2", "3"} {
		name := filepath.Base(path)
		subsystem := newSubsystem(name, path)
		t.subsystems = append(t.subsystems, subsystem)
	}
}

func (t *Target) GetSubsystem(name string) {
	var subsystem Subsystem
	paths, err := filepath.Glob(SubsystemsPath + name)
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		name := filepath.Base(path)
		fmt.Println(path, name)
		subsystem = newSubsystem(name, path)
		for idxPort := range t.ports { // might need to pass t.ports instead of two loops move loop to get_subsystem_ports_read_soft_link
			subsystem.scan(&t.ports[idxPort])

			//for i := range t.subsystems {
			//	t.subsystems[i].get_subsystem_ports_read_soft_link(&t.ports[idxPort])
			//}
		}

		t.subsystems = append(t.subsystems, subsystem)
	}
	//for idxPort := range t.ports { // might need to pass t.ports instead of two loops move loop to get_subsystem_ports_read_soft_link
	//	for i := range t.subsystems {
	//		t.subsystems[i].get_subsystem_ports_read_soft_link(&t.ports[idxPort])
	//	}
	//}

}

func (t *Target) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subsystem", "Offload", "Namespace_ID", "Device_Path", "PCI_ADDR", "Device_Enable", "IP", "Port"})
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(t.getReport())
	table.Render()
}

func (t *Target) getReport() [][]string {
	data := &[][]string{}
	for _, subsystem := range t.subsystems {
		//fmt.Printf("subsystem.getReport() %s\n", subsystem.getReport())
		*data = append(*data, subsystem.getReport())
	}

	//fmt.Printf("Target getReport %s\n", data)

	return *data
}

func NewUUID() string {
	uuidByte := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuidByte)
	if n != len(uuidByte) || err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// variant bits; see section 4.1.1
	uuidByte[8] = uuidByte[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuidByte[6] = uuidByte[6]&^0xf0 | 0x40

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", uuidByte[0:4], uuidByte[4:6], uuidByte[6:8], uuidByte[8:10], uuidByte[10:])
	return uuid
}

func deletePath(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err)
	}
}

func DeleteForce() {

	paths, err := filepath.Glob(PortsPath + "/*")
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, path := range paths {
		deletePath(path)
	}

	paths1, err := filepath.Glob(SubsystemsPath + "/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths1 {
		deletePath(path)
	}
}

func (t *Target) addSubsystem(subsystem ...Subsystem) {
	t.subsystems = append(t.subsystems, subsystem...)
}
func (t *Target) deleteAll() {
	//t.devices = append(t.devices, t.devices[1:]...)
	//t.devices = append(t.devices, t.devices[1:]...)
	//t.devices = []Device{}
	t.subsystems = []Subsystem{}
	t.ports = []Port{}
	t.disks = []Disk{}
	//[]Subsystem{},
	//	[]Port{},
	//	[]Device{},
}

func (t *Target) getSubsystemsNameSpaces() {
	for i := range t.subsystems {
		t.subsystems[i].getNamespaces()
	}
}

func (t *Target) getSubsystemsDevices() {
	for i := range t.subsystems {
		t.subsystems[i].getDevices()
	}
}

func (t *Target) SetSubsystemsDeviceSysFsEnable(enable bool) {
	for i := range t.subsystems {
		t.subsystems[i].SetDeviceSysFsEnable(enable)
	}
}

func (t *Target) SetSubsystemsPortSysFsEnable(enable bool) {
	for i := range t.subsystems {
		t.subsystems[i].SetDeviceSysFsEnable(enable)
	}
}

func (t *Target) destroySubsystems() {
	for i := range t.subsystems {
		t.subsystems[i].destroyPorts()
		t.subsystems[i].destroyDisks()
		t.subsystems[i].destroyNameSpace()
		t.subsystems[i].destroySubsystem()
	}

	for i := len(t.subsystems) - 1; i >= 0; i-- {
		if len(t.subsystems[i].namespaces) == 0 {
			if len(t.subsystems[i].ports) == 0 {
				t.subsystems = append(t.subsystems[:i], t.subsystems[i+1:]...)
			}
		}
	}
}

func (t *Target) scanSubsystems() {
	t.getSubsystems()
	t.getSubsystemsNameSpaces()
	t.getSubsystemsDevices()
	t.getSubsystemsPorts()
}

func (t *Target) Scan() {
	t.getPorts()
	t.getDisks()
	t.scanSubsystems()
}

func (t *Target) UnlinkSubsystemsPorts() {
	for i := range t.subsystems {
		t.subsystems[i].UnlinkSubsystemsPorts()
	}
}

func (t *Target) LinkSubsystemsPorts() {
	for i := range t.subsystems {
		t.subsystems[i].LinkSubsystemPorts()
	}
}

func (t *Target) SetSubsystemOffload(enable bool) {
	for i := range t.subsystems {
		SetSubsystemOffload(t.subsystems[i].name, enable)
	}
}

func (t *Target) DestroyPorts() {
	for i := range t.ports {
		if err := os.Remove(t.ports[i].path); err != nil {
			log.Println(err)
			return
		}
		t.ports[i].SetState(false)
	}

	for i := len(t.ports) - 1; i >= 0; i-- {
		if t.ports[i].State() {
			//t.ports[i] = nil
			t.ports = append(t.ports[:i], t.ports[i+1:]...)
		}
	}

}

func (t *Target) DestroyDisks() {
	//for i := range t.disks {
	//	if err := os.Remove(t.disks[i].path); err != nil {
	//		log.Fatal(err)
	//	}
	//}

	for i := len(t.disks) - 1; i >= 0; i-- {
		t.disks = append(t.disks[:i], t.disks[i+1:]...)
	}
}
func (t *Target) Destroy() {
	t.destroySubsystems()
	t.DestroyPorts()
	t.DestroyDisks()
}
func (t *Target) CreateNewSubsystem(name string, offload bool, ports *[]Port, disk Disk) error {
	//createSubsystem(name, offload)

	if err := CreateSubsystem(name, offload); err != nil {
		log.Println(err)
		golog.Error("This is an error", err)

		return err
	}

	//if err := CreateLinkPort(name, "172.17.100.8", "301"); err != nil {
	//	log.Println(err)
	//	return err
	//	return
	//}
	//t.

	//	for i := range *ports {
	//		//createPort(name, (*ports)[i].addrTraddr, (*ports)[i].addrTrsvcid)
	////CreateSub()
	//
	//		//if err := CreatePort(name, (*ports)[i].addrTraddr, (*ports)[i].addrTrsvcid); err != nil {
	//		//	log.Println(err)
	//		//	return err
	//		//}
	//
	//		t.GetPort(name)
	//
	//
	//	}
	return nil
}

func GetDisks() []string {
	paths, _ := filepath.Glob(SysFsNvmeGlob)

	return paths
}

func reloadModules(offloadMemSize, offloadBufferSize, numP2pQueues int, offloadMemStart uint64) {
	UnloadModules()
	if len(GetDisks()) > 0 {
		time.Sleep(3000 * time.Millisecond)
	}
	//time.Sleep(3000 * time.Millisecond)
	LoadModules(offloadMemSize, offloadBufferSize, numP2pQueues, offloadMemStart)
	if len(GetDisks()) < 1 {
		time.Sleep(3000 * time.Millisecond)
	}
	//time.Sleep(3000 * time.Millisecond)
}

func UnloadModules() {
	commands := []string{
		"systemctl stop systemd-udevd.service",
		"systemctl stop systemd-udevd-kernel.socket",
		"systemctl stop systemd-udevd-control.socket",
		"modprobe -r nvme_rdma nvmet_rdma nvmet nvme nvme_core",
	}

	for i := range commands {
		command := strings.Fields(commands[i])
		runCommand(command)
	}
}

func LoadModules(offloadMemSize, offloadBufferSize, numP2pQueues int, offloadMemStart uint64) {
	//cmd0 := "modprobe ib_umad"
	//cmd1 := "modprobe mlx5_ib"

	nvmeLoadCmd := fmt.Sprintf("modprobe nvme num_p2p_queues=%d", numP2pQueues)
	nvmetRDMALoadCmd := fmt.Sprintf("modprobe nvmet_rdma offload_mem_start=%#x offload_mem_size=%d offload_buffer_size=%d", offloadMemStart, offloadMemSize, offloadBufferSize)

	commands := []string{ibUmadLoadCmd, mlx5IBLoadCmd, nvmeLoadCmd, nvmetLoadCmd, nvmetRDMALoadCmd}

	fmt.Printf("\ncommands %s\n", commands)
	for i := range commands {
		command := strings.Fields(commands[i])
		runCommand(command)
	}
}

func LoadModulesppc64le(offloadMemSize, offloadBufferSize, numP2pQueues int, offloadMemStart string) {

	mem20 := "echo 0 > /sys/devices/system/memory/memory20/online"
	sleep2 := "sleep 2"
	mem21 := "echo 0 > /sys/devices/system/memory/memory21/online"

	nvmeLoadCmd := fmt.Sprintf("modprobe nvme num_p2p_queues=%d", numP2pQueues)
	nvmetRDMALoadCmd := fmt.Sprintf("modprobe nvmet_rdma offload_mem_start=%s offload_mem_size=%d offload_buffer_size=%d", offloadMemStart, offloadMemSize, offloadBufferSize)
	commands := []string{mem20, sleep2, mem21, nvmeLoadCmd, nvmetLoadCmd, nvmetRDMALoadCmd}

	fmt.Printf("\ncommands %s\n", commands)
	for i := range commands {
		command := strings.Fields(commands[i])
		runCommand(command)
	}
}

// LoadModulesDynamic nvmet module load
func LoadModulesDynamic(numP2pQueues int, offloadBufferSize int) {
	// maybe unload before load check if loaded first do stop ?
	var nvmeLoadCmd string

	if offloadBufferSize > 0 {
		nvmeLoadCmd = fmt.Sprintf("modprobe nvme num_p2p_queues=%d offloadBufferSize=%d", numP2pQueues, offloadBufferSize)
	} else {
		nvmeLoadCmd = fmt.Sprintf("modprobe nvme num_p2p_queues=%d", numP2pQueues)
	}

	// nvmeLoadCmd := fmt.Sprintf("modprobe nvme num_p2p_queues=%d", numP2pQueues)
	//cmd3 := "modprobe nvmet"
	//cmd4 := "modprobe nvmet_rdma"

	//cmd0 := "modprobe ib_umad"
	//cmd1 := "modprobe mlx5_ib"
	//cmd2 := fmt.Sprintf("modprobe nvme num_p2p_queues=%d", numP2pQueues)
	//cmd3 := "modprobe nvmet"
	//cmd4 := "modprobe nvmet_rdma"

	commands := []string{ibUmadLoadCmd, mlx5IBLoadCmd, nvmeLoadCmd, nvmetLoadCmd, nvmetRDMALoadCmd}

	fmt.Printf("\ncommands %s\n", commands)
	for i := range commands {
		command := strings.Fields(commands[i])
		runCommand(command)
	}
}

func runCommand(command []string) {
	//command := strings.Fields(command)
	fmt.Printf("\ncommand %s\n", command)
	cmd := exec.Command(command[0], command[1:]...)
	fmt.Printf("cmd.Args are %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cmd.Stdout is %s\n", cmd.Stdout)
	fmt.Printf("cmd.Process.Pid is %d\n", cmd.Process.Pid)
}

func CreateSub(name string, diskPath string, offload bool, ip string, id string) {
	//h, err := os.Hostname()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//var i int
	//for _, disk := range *disks {
	//target.disks[i].path
	//name := "sub_" + h + target.disks[i].name

	//name := fmt.Sprintf("sub_%s_%s", h, disk.name)
	//createSubsystem(s, true)

	//id, _ := strconv.Atoi(i)
	//i++
	//id := fmt.Sprintf("%d", i)
	////idA := fmt.Sprintf("300%s\n", id)
	//idA := fmt.Sprintf("300%s", id)
	//idB := fmt.Sprintf("400%s", id)
	////idB := fmt.Sprintf("400%s\n", id)
	//
	//fmt.Printf("s: %s idA: %s \n", name, idA)
	//fmt.Printf("s: %s idB: %s \n", name, idB)
	//
	//createSubsystem(name, offload)
	//createPort(name, Ips[0], idA)
	//createPort(name, Ips[1], idB)
	//createPort(name, "172.17.100.8\n", idA)
	//createPort(name, "172.17.200.8\n", idB)
	//createNameSpace(name, id, disk.path, newUUID(), byteOne)

	//}

	//createSubsystem(name, offload)
	//createPort(name, ip, id)
	//createNameSpace(name, id, diskPath, newUUID(), byteOne)
	//CreatePort(name, ips[1], idB)
}

//func createSubsystem(name string, offload bool) {
//func createSubsystem(name string, offload bool) {
//	if err := os.Mkdir(SubsystemsPath + name, 0777); err != nil {
//		log.Fatal(err)
//	}
//	if offload {
//		setAttr(SubsystemsPath + name + "/attr_offload", byteOne)
//	}
//	setAttr(SubsystemsPath + name + "/attr_allow_any_host", byteOne)
//}

//func ttt() {
//
//	target := New()
//	//t1.getPorts()
//	target.getPorts()
//
//	//fmt.Printf("ports %s \n", t1.ports)
//	////t1.ports = append(t1.ports, NewPort("1", "2", "1", true))
//	////t1.ports = append(t1.ports, NewPort("10", "20", "2", true))
//	////t1.ports = append(t1.ports, NewPort("30", "40", "3", true))
//	////
//	//
//	target.GetDisks()
//	//
//	target.getSubsystems()
//
//	//
//	target.getSubsystemsPorts()
//	target.getSubsystemsNameSpaces()
//
//	target.getSubsystemsDevices()
//
//	//fmt.Printf("ns %s \n", target.subsystems[0].namespaces)
//
//	//fmt.Printf("ports %s\n", target.subsystems[0].ports)
//
//	//target.destroyTarget()
//	//fmt.Printf("target port: %s\n", target.ports)
//	//fmt.Printf("target disks: %s\n", target.disks)
//	//fmt.Printf("port: %s\n", target.subsystems[5].ports)
//	//fmt.Printf("port: %p\n", &target.subsystems[5].ports[1])
//	//fmt.Printf("ports: %s\n", len(target.ports))
//	//fmt.Printf("ports: %s\n", len(target.subsystems[5].ports))
//	//fmt.Printf("ports: %s\n", len(target.ports))
//
//	target.PrintTable()
//}

func doIps(ips []net.Interface) []net.Interface {
	var newIps []net.Interface

	//fmt.Printf("ips len %v\n", len(ips))

	for i := range ips {
		fmt.Printf("ips len %v\n", len(ips))
		fmt.Printf("newIps len %v\n", len(newIps))
		if len(newIps) < 1 {
			newIps = append(newIps, ips[i])
			continue
		}
		//fmt.Print("newIps[i].HardwareAddr %s\n", newIps[i].HardwareAddr.String())
		//fmt.Printf("ips[i].HardwareAddr %s\n", ips[i].HardwareAddr.String())
		//fmt.Printf("newIps[i].HardwareAddr %s\n", newIps[i].HardwareAddr.String())
		if len(newIps) < len(ips) {
			if newIps[0].HardwareAddr.String() != ips[i].HardwareAddr.String() {
				newIps = append(newIps, ips[i])
			}
		}
		//if newIps[i].HardwareAddr.String() != ips[i].HardwareAddr.String() {
		//	newIps = append(newIps, ips[i])
		//}

		//if len(ips) > len(newIps) {
		//	continue
	}
	//}

	return newIps
}

//func GetIps() []string {
//	var ips []net.Interface
//	ifaces := getIps()
//	for i :=  range ifaces {
//		ip, err := ifaces[i].Addrs()
//		if err != nil {
//			//return
//			//return "", err
//		}
//		ips = append(ips, ip)
//	}
//}

func GetIps() []string {
	var ips []string

	ifaces, err := net.Interfaces()
	if err != nil {
		//return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			//return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			if strings.Contains(ip.String(), "10.144") {
				continue
			}

			if strings.Contains(ip.String(), "10.130") {
				continue
			}

			if strings.Contains(ip.String(), "10.12.151") {
				continue
			}

			ips = append(ips, ip.String())
			//fmt.Println(ip)
			//fmt.Println(iface.Name)
			//fmt.Println(iface.HardwareAddr)
			//fmt.Println(iface)

			//return ip.String(), nil
		}
	}
	//fmt.Println(ips)

	return ips
}
func getInterface() []net.Interface {
	var netIfaces []net.Interface
	ifaces, err := net.Interfaces()
	if err != nil {
		//return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			//return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			if strings.Contains(ip.String(), "10.144") {
				continue
			}

			if strings.Contains(ip.String(), "10.130") {
				continue
			}

			netIfaces = append(netIfaces, iface)
			//fmt.Println(ip)
			//fmt.Println(iface.Name)
			//fmt.Println(iface.HardwareAddr)
			//fmt.Println(iface)

			//return ip.String(), nil
		}
	}
	fmt.Println(netIfaces)

	return netIfaces
	//return "", errors.New("are you connected to the network?")
}

func GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return name
}
