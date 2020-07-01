package target

import (
	"log"
)

type Device struct {
	name          string
	devicePath    string
	pciDevicePath string
	enable        bool
	path 		  string
}

func (d *Device) Enable() bool {
	return d.enable
}

func (d *Device) SetDeviceSysFsEnable(enable bool) {
	//fmt.Println(d.path)
	if enable {
		if err := setAttr(d.path + "/enable", ByteOne); err != nil {
			log.Println(err)
			return
		}
	} else {
		if err := setAttr(d.path + "/enable", ByteZero); err != nil {
			log.Println(err)
			return
		}
	}

	//if bytes.Equal(enable, ByteOne) {
	d.SetEnable(enable)
		//return
	//}
	//
	//d.SetEnable(false)
}

func (d *Device) SetEnable(enable bool) {
	d.enable = enable
}

func (d *Device) IsDisabled() bool {
	if d.Enable() {
		return false
	}
	return true
}