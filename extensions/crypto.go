package extensions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"

	config "GitlabTgBot/configuration"
)

func Encrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(config.GetConfigInstance().EncryptionKey))
	if err != nil {
		return nil, err
	}

	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(text []byte) string {
	block, err := aes.NewCipher([]byte(config.GetConfigInstance().EncryptionKey))
	if err != nil {
		log.Panic(err)
	}

	if len(text) < aes.BlockSize {
		log.Panic(err)
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		log.Panic(err)
	}
	return string(data)
}
