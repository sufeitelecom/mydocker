package subsystems

import (
	"io/ioutil"
	"path"
	"fmt"
	"strconv"
	"os"
)

type MemorySubsystem struct {
}

var _ Subsystem  = &MemorySubsystem{}


func (m *MemorySubsystem)Name() string {
	return "memory"
}

func (m *MemorySubsystem)Set(cgrouppath string,res *ResourceConfig) error{
	if systempath , err := GetCgroupPath(m.Name(),cgrouppath,true);err ==nil {
		if res.MemoryLimit != ""{
			if err := ioutil.WriteFile(path.Join(systempath,"memory.limit_in_bytes"),[]byte(res.MemoryLimit),0644);err != nil{
				return fmt.Errorf("set cgroup memory fail %v",err)
			}
		}
		return nil
	}else {
		return nil
	}
}

func (m *MemorySubsystem)Apply(cgrouppath string,pid int) error{
	if SubsyscgroupPath ,err := GetCgroupPath(m.Name(),cgrouppath,false);err == nil{
		if err := ioutil.WriteFile(path.Join(SubsyscgroupPath,"tasks"),[]byte(strconv.Itoa(pid)),0644);err != nil{
			return fmt.Errorf("set cgroup proc fail %v",err)
		}
		return nil
	}else {
		return fmt.Errorf("get cgroup %s error:%v",cgrouppath,err)
	}
}

func (m *MemorySubsystem)Remove(cgrouppath string) error  {
	if SubsyscgroupPath,err := GetCgroupPath(m.Name(),cgrouppath,false);err == nil{
		return os.Remove(SubsyscgroupPath)
	}else {
		return err
	}
}