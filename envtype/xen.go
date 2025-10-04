//go:build windows

package envtype

import (
	"vortex/hutil"
	"vortex/opsysinfo"
	"strings"
)

func xenProcs() bool {

	procId, _ := hutil.GetProcessIdByName("xenservice")
	if procId != 0 {
		return true
	}

	return false
}

func xenMacAddr() bool {

	macAddr := opsysinfo.LocalNetworkInfo()

	if strings.HasPrefix(macAddr.MacAddress, "00:16:3E") {
		return true
	}

	return false
}

func detectXen() bool {

	if xenProcs() || xenMacAddr() {
		return true
	}

	return false
}
