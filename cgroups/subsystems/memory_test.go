package subsystems

import (
	"testing"
	"os"
	"path"
)

func TestMemorySubsystem(t *testing.T) {
	mem := MemorySubsystem{}
	res := ResourceConfig{
		MemoryLimit:"1G",
	}
    testpath := "testmemory"

	t.Logf("start test memorysubsystem")

    if err := mem.Set(testpath,&res);err != nil{
    	t.Fatalf("Set fail:%v ",err)
	}

	stat, _ := os.Stat(path.Join(findMountPoint("memory"), testpath))
	t.Logf("cgroup stats: %+v", stat)

	if err := mem.Apply(testpath,os.Getpid());err != nil{
		t.Fatalf("Apply fail:%v",err)
	}

	if err := mem.Apply("",os.Getpid());err != nil{
		t.Fatalf("Move back root fail:%v",err)
	}

	if err := mem.Remove(testpath);err != nil{
		t.Fatalf("Remove fail:%v ",err)
	}
	t.Logf("finish test memorysubsystem")
}
