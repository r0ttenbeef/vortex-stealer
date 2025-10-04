//go:build windows

package envtype

import (
	"errors"
	"os"
	"path/filepath"
	"vortex/hutil"
	"vortex/opsysinfo"
	"strings"

	"github.com/digitalocean/go-smbios/smbios"
	"golang.org/x/sys/windows/registry"
)

func vmwareRegKV() bool {

	regHKLM := map[string]string{
		"SYSTEM\\ControlSet001\\Control\\SystemInformation":                                  "SystemManufacturer",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0": "Identifier",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 1\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0": "Identifier",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 2\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0": "Identifier",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 3\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0": "Identifier",
	}

	regHKCU := "Software\\VMware, Inc."

	if _, err := registry.OpenKey(registry.CURRENT_USER, regHKCU, registry.QUERY_VALUE); err != registry.ErrNotExist {
		return true
	}

	for k, v := range regHKLM {
		reg, err := registry.OpenKey(registry.LOCAL_MACHINE, k, registry.QUERY_VALUE)
		defer reg.Close()
		if err == nil {
			val, _, _ := reg.GetStringValue(v)
			if strings.Contains(strings.ToLower(val), "vmware") {
				return true
			}
		}
	}

	return false
}

func vmwareDriveFiles() bool {

	drvFile := []string{
		"vmnet.sys", "vmmouse.sys", "vmusb.sys", "vm3dmp.sys", "vmci.sys",
		"vmhgfs.sys", "vmmemctl.sys", "vmx86.sys", "vmrawdsk.sys", "vmusbmouse.sys",
		"vmkdb.sys", "vmnetuserif.sys", "vmnetadapter.sys",
	}

	if _, err := os.Stat(filepath.Join(os.Getenv("ProgramW6432"), "VMware")); !errors.Is(err, os.ErrNotExist) {
		return true
	}

	for i := range drvFile {
		if _, err := os.Stat(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "drivers", drvFile[i])); !errors.Is(err, os.ErrNotExist) {
			return true
		}
	}

	return false

}

func vmwareMacAddr() bool {

	vmwareMacIds := []string{"00:05:69", "00:0c:29", "00:1C:14", "00:50:56"}

	macAddr := opsysinfo.LocalNetworkInfo()

	for i := range vmwareMacIds {
		if strings.HasPrefix(macAddr.MacAddress, vmwareMacIds[i]) {
			return true
		}
	}
	return false
}

func vmwareProcs() bool {

	vmwareProcList := []string{"vmtoolsd", "vmwaretray", "vmwareuser", "vgauthservice", "vmacthlp"}

	for i := range vmwareProcList {
		procId, _ := hutil.GetProcessIdByName(vmwareProcList[i])
		if procId != 0 {
			return true
		}
	}
	return false
}

func vmwareSMBIOS() bool {
	rc, _, err := smbios.Stream()
	if err != nil {
		return false
	}
	defer rc.Close()

	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		return false
	}

	for i := range ss {
		if strings.Contains(strings.ToLower(strings.Join(ss[i].Strings, " ")), "vmware") {
			return true
		}
	}

	return false
}

func detectVMWware() bool {

	if vmwareRegKV() || vmwareMacAddr() || vmwareDriveFiles() || vmwareProcs() || vmwareSMBIOS() {
		return true
	}

	return false
}
