//go:build windows

package opsysinfo

import (
	"fmt"
	"net"
	"os/user"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/registry"
)

type DeviceInfo struct {
	ScreenResolution     string
	Username             string
	CpuName              string
	Hostname             string `json:"hostname"`
	RAM                  uint64 `json:"ram"`
	VirtualizationSystem string `json:"virtualizationSystem"`
	Procs                int32  `json:"procs"`
	CPUArch              string `json:"cpuArch"`
	HostID               string `json:"hostid"`
}

type WindowsInfo struct {
	RegisteredOwner string
	ProductName     string
	DisplayVersion  string
}

type DiskStatus struct {
	DiskPartitions []string
	DiskTotalSpace uint64 `json:"total"`
	DiskFreeSpace  uint64 `json:"free"`
	DiskUsedSpace  uint64 `json:"used"`
}

type NicInfo struct {
	MacAddress   string
	LocalAddress string
}

// Get partitions and hard disk states
func HardDriveInfo() (diskstat DiskStatus) {
	partitions, _ := disk.Partitions(false)

	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		diskstat.DiskFreeSpace += usage.Free
		diskstat.DiskTotalSpace += usage.Total
		diskstat.DiskUsedSpace += usage.Used
		diskstat.DiskPartitions = append(diskstat.DiskPartitions, partition.Device)
	}

	return
}

// Windows information
func OsVersion() (winInfo WindowsInfo) {
	winInfoRegKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.QUERY_VALUE)
	defer winInfoRegKey.Close()

	winInfo.RegisteredOwner, _, _ = winInfoRegKey.GetStringValue("RegisteredOwner")
	winInfo.ProductName, _, _ = winInfoRegKey.GetStringValue("ProductName")
	winInfo.DisplayVersion, _, _ = winInfoRegKey.GetStringValue("DisplayVersion")

	return
}

// Machine specs of the device
func MachineSpecs() DeviceInfo {
	var (
		hostStat, _    = host.Info()
		memStat, _     = mem.VirtualMemory()
		cpuStat, _     = cpu.Info()
		userName, _    = user.Current()
		onlyUserString = strings.Split(userName.Username, "\\")
	)

	var deviceInfo = DeviceInfo{
		Hostname:             hostStat.Hostname,
		Username:             onlyUserString[len(onlyUserString)-1],
		RAM:                  memStat.Total / (1024 * 1024),
		VirtualizationSystem: hostStat.VirtualizationSystem,
		Procs:                cpuStat[0].Cores,
		CPUArch:              hostStat.KernelArch,
		HostID:               hostStat.HostID,
	}

	return deviceInfo
}

// Local network IP Address and Mac Address
func LocalNetworkInfo() (nicInfo NicInfo) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagBroadcast == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !strings.HasPrefix(ipnet.IP.String(), "169") && ipnet.IP.String() != "" && ipnet.IP.To4() != nil {
				netInterface, err := net.InterfaceByName(iface.Name)
				if err != nil {
					return
				}
				nicInfo.LocalAddress = ipnet.IP.String()
				nicInfo.MacAddress = netInterface.HardwareAddr.String()
				break
			}
		}
	}

	return
}

// Screen size resolution
func ScreenResolution() (DeviceInfo, int, int) {
	const (
		SM_CXSCREEN = 0
		SM_CYSCREEN = 1
	)

	var (
		data             DeviceInfo
		user32           = syscall.NewLazyDLL("User32.dll")
		getSystemMetrics = user32.NewProc("GetSystemMetrics")
	)

	cxScreen, _, _ := getSystemMetrics.Call(uintptr(SM_CXSCREEN))
	cyScreen, _, _ := getSystemMetrics.Call(uintptr(SM_CYSCREEN))

	data.ScreenResolution = fmt.Sprintf("%dx%d", cxScreen, cyScreen)

	return data, int(cxScreen), int(cyScreen)
}
