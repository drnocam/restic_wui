package fileop

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"encoding/base64"
	"os"
)

const file_name = "test.json"

const secret_word string = "agsdfgs4345dsf*/"

var key_bytes []byte = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		panic(err)
	}
	return data
}
func Encrypt(text []byte, secret_word string) (string, error) {
	block, err := aes.NewCipher([]byte(secret_word))
	if err != nil {
		return "", err
	}
	plainText := text
	cfb := cipher.NewCFBEncrypter(block, key_bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text []byte, secret_word string) (string, error) {
	block, err := aes.NewCipher([]byte(secret_word))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, key_bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func WriteSettings(data []byte) error {

	enc, er := Encrypt(data, secret_word)
	if er == nil {

		err := os.WriteFile(file_name, []byte(enc), 0644)
		return err
	}
	return nil
}

func ReadSettings() (string, error) {

	if data, err := os.ReadFile(file_name); err == nil {

		if datax, errx := Decrypt(data, secret_word); errx == nil {
			return datax, nil
		}
		return "", errors.New("File couldnt decrypt.")
	}
	return "", errors.New("File couldnt read.")
}
