package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare string
	CpuSet string
}

type Subsystem interface {
	Name() string
	Set(cgrouppath string,res *ResourceConfig) error
	Apply(cgrouppath string,pid int) error
	Remove(cgrouppath string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubsystem{},
		&MemorySubsystem{},
		&CpuSubsystem{},
	}
)