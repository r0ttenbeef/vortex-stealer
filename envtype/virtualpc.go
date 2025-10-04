//go:build windows

package envtype

import (
	"vortex/hutil"

	"golang.org/x/sys/windows/registry"
)

func vpcRegKV() bool {

	if _, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Virtual Machine\\Guest\\Parameters", registry.QUERY_VALUE); err != registry.ErrNotExist {
		return true
	}

	return false
}

func vpcProcs() bool {

	vpcProcList := []string{"VMSrvc", "VMUSrvc"}

	for i := range vpcProcList {
		procId, _ := hutil.GetProcessIdByName(vpcProcList[i])
		if procId != 0 {
			return true
		}
	}

	return false
}

func detectVpc() bool {

	if vpcRegKV() || vpcProcs() {
		return true
	}

	return false
}
