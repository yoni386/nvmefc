package target

type Disk struct {
	name          string
	path          string
	pciDevicePath string
}

func (d *Disk) PciDevicePath() string {
	return d.pciDevicePath
}

func (d *Disk) Path() string {
	return d.path
}

func (d *Disk) Name() string {
	return d.name
}
