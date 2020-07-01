package target

import (
	"os"
	"log"
)

type Port struct {
	addrTraddr  string
	addrTrsvcid string // need to trim new line
	path        string
	state       bool
}

func NewPort(addr_traddr string, addr_trsvcid string, path string, state bool) Port {
	return Port{addr_traddr, addr_trsvcid, path, state}
}

func LinkPort(subsystemName string, port string) error {
	src := SubsystemsPath + subsystemName
	dst := PortsPath + port + "/subsystems/" + subsystemName

	//fmt.Println("src, dst", src, dst)


	if err := os.Symlink(src, dst); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (p *Port) UnLinkPort(subsystemName string) error {
	//src := SubsystemsPath + subsystemName
	path := p.Path() + "/subsystems/" + subsystemName
	//dst := PortsPath + port + "/subsystems/" + subsystemName

	//if err := os.Symlink(src, dst); err != nil {
	//	log.Println(err)
	//	return err
	//}

	//log.Println(path)

	//if err := os.Remove(fmt.Sprintf("%s/subsystems/%s", s.ports[i].path, s.name)); err != nil {
	//	log.Fatal(err)
	//	return
	//}


	if err := os.Remove(path); err != nil {
		log.Println(err)
		return err
	}
	p.SetState(false) // s.ports[i].Disable()??
	//t.ports[i].SetState(false)

	return nil

}

func CreateLinkPort(subsystemName string, ip string, port string) error {
	if err := CreatePort(ip, port); err != nil {
		//log.Println(err)
		return err
	}

	if err := LinkPort(subsystemName, port); err != nil {
		//log.Println(err)
		return err
	}
	return nil
}

func CreatePort(ip string, port string) error {

	if err := os.Mkdir(PortsPath + port, 0777); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(PortsPath + port + "/addr_trsvcid", []byte(port)); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(PortsPath + port + "/addr_trtype", []byte("rdma")); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(PortsPath + port + "/addr_adrfam", []byte("ipv4")); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(PortsPath + port + "/addr_traddr", []byte(ip)); err != nil {
		//log.Println(err)
		return err
	}

	//setAttr(PortsPath + port + "/addr_trsvcid", []byte(port))
	//setAttr(PortsPath + port + "/addr_trtype", []byte("rdma"))
	//setAttr(PortsPath + port + "/addr_adrfam", []byte("ipv4"))
	//setAttr(PortsPath + port + "/addr_traddr", []byte(ip))

	//src := SubsystemsPath + subsystemName
	//dst := PortsPath + port + "/subsystems/" + subsystemName
	//
	//if err := os.Symlink(src, dst); err != nil {
	//	log.Println(err)
	//	return err
	//}

	return nil
}


func (p *Port) State() bool {
	return p.state
}

func (p *Port) SetState(state bool) {
	p.state = state
}

func (p *Port) Path() string {
	return p.path
}

func (p *Port) SetPath(path string) {
	p.path = path
}

func (p *Port) AddrTrsvcid() string {
	//return strings.Trim(p.AddrTrsvcid(), "\n")
	return p.addrTrsvcid
}

func (p *Port) SetAddrTrsvcid(addrTrsvcid string) {
	p.addrTrsvcid = addrTrsvcid
}

func (p *Port) AddrTraddr() string {
	return p.addrTraddr
}

func (p *Port) SetAddrTraddr(addrTraddr string) {
	p.addrTraddr = addrTraddr
}
func (p *Port) destroyPort() {
	log.Println(p.path)

	if err := os.Remove(p.path); err != nil {
		log.Fatal(err)
		//log.Println(err)
		//return err
	}
	p.SetState(false)
}

func (p *Port) IsDisabled() bool {
	if p.state {
		return false
	}
	return true
}