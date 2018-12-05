package subsystems

import "testing"

func TestFindMountPoint(t *testing.T) {
	t.Logf("cpu subsystem mount point %v\n", findMountPoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", findMountPoint("cpuset"))
	t.Logf("memory subsystem mount point %v\n", findMountPoint("memory"))
}