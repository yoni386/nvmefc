package target

import (
	"fmt"
	"os"
	"log"
)

type Namespace struct {
	name string
	path string
	Device
}


func CreateNameSpace(subsystemName string, id int, devicePath string, deviceNguid string, enable []byte) error {

	path := fmt.Sprintf("/sys/kernel/config/nvmet/subsystems/%s/namespaces/%d", subsystemName, id)

	if err := os.Mkdir(path, 0777); err != nil {
		log.Println(err)
		return err
	}

	//setAttr(path + "/device_path", []byte(devicePath))
	//setAttr(path + "/device_nguid", []byte(deviceNguid))
	//setAttr(path + "/enable", []byte(enable))
	//
	if err := setAttr(path + "/device_path", []byte(devicePath)); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(path + "/device_nguid", []byte(deviceNguid)); err != nil {
		//log.Println(err)
		return err
	}

	if err := setAttr(path + "/enable", []byte(enable)); err != nil {
		//log.Println(err)
		return err
	}

	return  nil
}