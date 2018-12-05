package cgroups

import "github.com/sufeitelecom/mydocker/cgroups/subsystems"

type CgroupManager struct {
	cgrouppath string
	res *subsystems.ResourceConfig
}

func NewCgroupManager(cgrouppath string) *CgroupManager {
	return &CgroupManager{
		cgrouppath:cgrouppath,
	}
}

func (c *CgroupManager)Set(res *subsystems.ResourceConfig) error  {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Set(c.cgrouppath, res)
	}
	return nil
}

func (c *CgroupManager)Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.cgrouppath,pid)
	}
	return nil
}

func (c *CgroupManager)Destory() error  {
	for  _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Remove(c.cgrouppath)
	}
	return nil
}