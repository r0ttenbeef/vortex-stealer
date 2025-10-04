//go:build windows

package envtype

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/digitalocean/go-smbios/smbios"
	"golang.org/x/sys/windows/registry"
)

func kvmRegKV() bool {

	regHKLMsubKeys := []string{
		"SYSTEM\\ControlSet001\\Services\\vioscsi",
		"SYSTEM\\ControlSet001\\Services\\viostor",
		"SYSTEM\\ControlSet001\\Services\\VirtIO-FS Service",
		"SYSTEM\\ControlSet001\\Services\\VirtioSerial",
		"SYSTEM\\ControlSet001\\Services\\BALLOON",
		"SYSTEM\\ControlSet001\\Services\\BalloonService",
		"SYSTEM\\ControlSet001\\Services\\netkvm",
	}

	regHKLM := [][]string{
		{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", "Identifier"},
		{"HARDWARE\\Description\\System", "SystemBiosVersion"},
	}

	for i := range regHKLMsubKeys {
		if _, err := registry.OpenKey(registry.LOCAL_MACHINE, regHKLMsubKeys[i], registry.QUERY_VALUE); err != registry.ErrNotExist {
			return true
		}
	}

	for _, v := range regHKLM {
		reg, err := registry.OpenKey(registry.LOCAL_MACHINE, v[0], registry.QUERY_VALUE)
		defer reg.Close()
		if err == nil {
			val, _, _ := reg.GetStringValue(v[1])
			if strings.Contains(strings.ToLower(val), "qemu") {
				return true
			}
		}
	}

	return false
}

func kvmDriveFiles() bool {

	drvFiles := []string{
		"balloon.sys", "netkvm.sys", "pvpanic.sys", "viofs.sys",
		"viogpudo.sys", "vioinput.sys", "viorng.sys",
		"vioscsi.sys", "vioser.sys", "viostor.sys",
	}

	drvDirs := []string{"Virtio-Win", "qemu-ga", "Qemu-ga", "Spice Agent", "SPICE Guest Tools"}

	for i := range drvDirs {
		if _, err := os.Stat(filepath.Join(os.Getenv("ProgramW6432"), drvDirs[i])); !errors.Is(err, os.ErrNotExist) {
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

func kvmSMBIOS() bool {

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
		if strings.Contains(strings.ToLower(strings.Join(ss[i].Strings, " ")), "qemu") {
			return true
		}
	}

	return false
}

func detectKVM() bool {

	if kvmRegKV() || kvmDriveFiles() || kvmSMBIOS() {
		return true
	}

	return false
}
