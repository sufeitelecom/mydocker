package network

import (
	"path"
	"os"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

type IPAM struct {
	Subnetlocalpath string
	Subnets *map[string]string
}

const ipamDefaultLocalPath  = "/home/sufei/busybox/root/mydocker/network/ipam/subnet.json"

var ipAllocator = &IPAM{
	Subnetlocalpath:ipamDefaultLocalPath,
}

func (ipam *IPAM) dump() error {
	ipampath , _ := path.Split(ipamDefaultLocalPath)
	if _,err := os.Stat(ipampath);err != nil{
		if os.IsNotExist(err){
			os.MkdirAll(ipampath,0644)
		}else {
			return err
		}
	}
	subnetfile,err := os.OpenFile(ipamDefaultLocalPath,os.O_CREATE | os.O_TRUNC | os.O_WRONLY,0644)
	defer subnetfile.Close()

	if err != nil{
		return err
	}

	ipamjson,err :=json.Marshal(ipam.Subnets)
	if err != nil{
		return err
	}

	_,err = subnetfile.Write(ipamjson)
	if err != nil{
		return err
	}
	return nil
}

func (ipam *IPAM) load() error  {
	if _,err := os.Stat(ipamDefaultLocalPath);err != nil{
		if os.IsNotExist(err){
			return nil
		}else {
			return err
		}
	}

	subnetconfigfile,err:= os.Open(ipamDefaultLocalPath)
	if err != nil{
		return err
	}

	subnetjson := make([]byte,2000)
	n,err := subnetconfigfile.Read(subnetjson)
	if err != nil{
		return err
	}

	err = json.Unmarshal(subnetjson[:n],ipam.Subnets)
	if err != nil{
		log.Errorf("Error dump allocation info, %v",err)
		return err
	}
	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP,err error)  {
	ipam.Subnets = &map[string]string{}

	err = ipam.load()
	if err != nil{
		log.Errorf("Error dump allocation info, %v", err)
	}

	_,subnet,_ = net.ParseCIDR(subnet.String())

	one,size := subnet.Mask.Size()

	if _,exist := (*ipam.Subnets)[subnet.String()]; !exist{
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0",1<<uint8(size-one))
	}

	for c := range (*ipam.Subnets)[subnet.String()] {
		if (*ipam.Subnets)[subnet.String()][c] == '0'{
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)

			ip = subnet.IP
			//ip是一个uint数组，如192.168.1.0中192在ip[0]处，依次类推
			for t:=uint(4);t>0;t-=1{
				[]byte(ip)[4-t] += uint8(c >> ((t-1)*8))
			}
			ip[3] +=1
			break
		}
	}
	ipam.dump()
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}

	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3]-=1
	for t := uint(4); t > 0; t-=1 {
		c += int(releaseIP[t-1] - subnet.IP[t-1]) << ((4-t) * 8)
	}

	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.dump()
	return nil
}