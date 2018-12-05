package subsystems

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"path"
)

func findMountPoint(subsystem string) string{
	fs,err := os.Open("/proc/self/mountinfo")
	if err != nil{
		return ""
	}
	defer fs.Close()

	scan := bufio.NewScanner(fs)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		txt := scan.Text()
		feild := strings.Split(txt," ")
		for _,opt := range strings.Split(feild[len(feild)-1],",") {
			if opt == subsystem {
				return feild[4]
			}
		}
	}

	if err := scan.Err();err != nil {
		return ""
	}

	return ""
}

func GetCgroupPath (subsystem string,cgroupPath string,autoCreate bool) (string, error){
	cgroupRoot := findMountPoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}