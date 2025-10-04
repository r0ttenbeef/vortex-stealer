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

func vboxRegKV() bool {

	regHKLM := [][]string{
		{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
		{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 1\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
		{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 2\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
		{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 3\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
		{"HARDWARE\\Description\\System", "SystemBiosVersion"},
		{"HARDWARE\\Description\\System", "VideoBiosVersion"},
	}

	regSubKeys := []string{
		"HARDWARE\\ACPI\\DSDT\\VBOX__",
		"HARDWARE\\ACPI\\FADT\\VBOX__",
		"HARDWARE\\ACPI\\RSDT\\VBOX__",
		"SOFTWARE\\Oracle\\VirtualBox Guest Additions",
		"SYSTEM\\ControlSet001\\Services\\VBoxGuest",
		"SYSTEM\\ControlSet001\\Services\\VBoxMouse",
		"SYSTEM\\ControlSet001\\Services\\VBoxService",
		"SYSTEM\\ControlSet001\\Services\\VBoxSF",
		"SYSTEM\\ControlSet001\\Services\\VBoxVideo",
	}

	for _, v := range regHKLM {
		reg, err := registry.OpenKey(registry.LOCAL_MACHINE, v[0], registry.QUERY_VALUE)
		defer reg.Close()
		if err == nil {
			val, _, _ := reg.GetStringValue(v[1])
			if strings.Contains(strings.ToLower(val), "vbox") || strings.Contains(strings.ToLower(val), "virtualbox") {
				return true
			}
		}
	}

	for i := range regSubKeys {
		if _, err := registry.OpenKey(registry.LOCAL_MACHINE, regSubKeys[i], registry.QUERY_VALUE); err != registry.ErrNotExist {
			return true
		}
	}

	return false
}

func vboxDriveFiles() bool {

	dllFiles := []string{
		"vboxdisp.dll", "vboxhook.dll", "vboxmrxnp.dll", "vboxogl.dll",
		"vboxoglarrayspu.dll", "vboxoglcrutil.dll", "vboxoglerrorspu.dll",
		"vboxoglfeedbackspu.dll", "vboxoglpackspu.dll", "vboxtray.exe",
		"vboxoglpassthroughspu.dll", "vboxservice.exe", "VBoxControl.exe",
	}

	drvFiles := []string{"VBoxMouse.sys", "VBoxGuest.sys", "VBoxSF.sys", "VBoxVideo.sys"}

	for i := range dllFiles {
		if _, err := os.Stat(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", dllFiles[i])); !errors.Is(err, os.ErrNotExist) {
			return true
		}
	}

	for i := range drvFiles {
		if _, err := os.Stat(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "drivers", drvFiles[i])); !errors.Is(err, os.ErrNotExist) {
			return true
		}
	}

	return false
}

func vboxMacAddr() bool {

	vboxMacId := "08:00:27"
	macAddr := opsysinfo.LocalNetworkInfo()

	if strings.HasPrefix(macAddr.MacAddress, vboxMacId) {
		return true
	}

	return false
}

func vboxProcs() bool {

	vboxProcList := []string{"vboxtray", "vboxservice"}

	for i := range vboxProcList {
		procId, _ := hutil.GetProcessIdByName(vboxProcList[i])
		if procId != 0 {
			return true
		}
	}

	return false
}

func vboxSMBIOS() bool {
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
		if strings.Contains(strings.ToLower(strings.Join(ss[i].Strings, " ")), "vbox") ||
			strings.Contains(strings.ToLower(strings.Join(ss[i].Strings, " ")), "virtualbox") ||
			strings.Contains(strings.ToLower(strings.Join(ss[i].Strings, " ")), "oracle corporation") {
			return true
		}
	}

	return false
}

func detectVBox() bool {

	if vboxRegKV() || vboxMacAddr() || vboxDriveFiles() || vboxProcs() || vboxSMBIOS() {
		return true
	}

	return false
}
