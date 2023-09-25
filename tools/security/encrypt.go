package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"knaq-wallet/tools/logo"
	"strings"
)

type Encryption struct {
	SecretKey []byte
}

func (e Encryption) Encrypt(stringToEncrypt string) (encryptedString string, err error) {
	if stringToEncrypt == "" {
		return
	}
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(e.SecretKey)
	if err != nil {
		return encryptedString, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return encryptedString, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return encryptedString, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return fmt.Sprintf("%x", ciphertext), err
}

func (e Encryption) Decrypt(encryptedString string) (plainText string) {
	if encryptedString == "" {
		return ""
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return encryptedString
	}
	block, err := aes.NewCipher(e.SecretKey)
	if err != nil {
		return encryptedString
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return encryptedString
	}
	nonceSize := aesGCM.NonceSize()
	if len(enc) >= 0 && len(enc) >= nonceSize {
		var plainTextByte []byte
		nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
		plainTextByte, err = aesGCM.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return encryptedString
		}
		plainText = fmt.Sprintf("%s", plainTextByte)
	} else {
		logo.Printf("slice bounds out of range [%d:%d]\n", nonceSize, len(enc))
		return encryptedString
	}

	return
}

func (e Encryption) DecryptAes128Ecb(data string) (decryptedString string) {
	if data == "" {
		return ""
	}
	if !strings.HasSuffix(data, "=") {
		return data
	}
	base64Decrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	c, _ := aes.NewCipher(e.SecretKey[:32])
	decrypted := make([]byte, len(base64Decrypted))
	size := 16

	if len(base64Decrypted) >= 0 && len(base64Decrypted) >= size {
		for bs, be := 0, size; bs < len(base64Decrypted); bs, be = bs+size, be+size {
			c.Decrypt(decrypted[bs:be], base64Decrypted[bs:be])
		}

		decrypted = pkcs5UnPadding(decrypted)

		decryptedString = string(decrypted)
	} else {
		logo.Printf("slice bounds out of range [%d:%d]\n", size, len(base64Decrypted))
		return data
	}

	return
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
