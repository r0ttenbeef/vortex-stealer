package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strings"
)

func Decrypt(encrypted_text string) string {
	cipherKey := []byte{
		0xD5, 0x03, 0xE6, 0xE0, 0x63, 0x52, 0x7C,
		0xB6, 0x24, 0xFE, 0x03, 0x63, 0xFF, 0xF9,
		0xB3, 0xBD, 0x20, 0x94, 0x1C, 0xAF, 0x70,
		0x84, 0x92, 0xB6, 0x90, 0x5F, 0x66, 0x43,
		0x4D, 0xCA, 0x72, 0x77}

	IV := []byte{
		0x27, 0xD8, 0xB1, 0xF0, 0xDB, 0x3C, 0xAB,
		0x3E, 0x20, 0x20, 0x21, 0x56, 0xBA, 0x1B,
		0x37, 0x13}

	decText, _ := base64.StdEncoding.DecodeString(encrypted_text)

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		panic(err)
	}

	decMode := cipher.NewCBCDecrypter(block, IV)
	textDec := make([]byte, len(decText))
	decMode.CryptBlocks(textDec, decText)

	decryptedStr := strings.TrimRight(string(textDec), "\x01\x02\x03\x04\x05\x06\x07\x08")

	return string(decryptedStr)
}

func B64Util(s string, method int) string {
	switch method {
	case 0:
		return "ERB1-7C" + base64.StdEncoding.EncodeToString([]byte(s))

	case 1:
		decodedStr, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, "ERB1-7C"))
		return string(decodedStr)

	default:
		return "you have to choose between encode and decode"
	}
}
