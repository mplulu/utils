package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/p2play/p2play_backend/derror"
	"github.com/p2play/p2play_backend/log"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EncryptWithRandKey(data []byte) (key []byte, nonce []byte, result []byte, err error) {
	randKey := RandHexSeq(32)
	semiKey, err := hex.DecodeString(randKey)
	if err != nil {
		log.LogSerious("err encrypt 0 %v %v", err, randKey)
		return nil, nil, nil, err
	}
	fullKey := fmt.Sprintf("%s%s", randKey, kSalt)
	key, err = hex.DecodeString(fullKey)
	if err != nil {
		log.LogSerious("err encrypt 1 %v %v", err, randKey)
		return nil, nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.LogSerious("err encrypt 2 %v", err)
		return nil, nil, nil, err
	}
	nonce = make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.LogSerious("err encrypt 3 %v", err)
		return nil, nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.LogSerious("err encrypt 4 %v", err)
		return nil, nil, nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return semiKey, nonce, ciphertext, nil
}

func DecryptWithRandKey(key []byte, nonce []byte, ciphertext []byte) (result []byte, err error) {
	if len(key) == 0 || len(nonce) == 0 || len(ciphertext) == 0 {
		return nil, derror.NewErrorSingle("err:invalid_token", "reason", "cipher_empty")
	}
	semiKeyStr := hex.EncodeToString(key)
	fullKey := fmt.Sprintf("%s%s", semiKeyStr, kSalt)

	key, err = hex.DecodeString(fullKey)
	if err != nil {
		log.LogSerious("err decrypt 0 %v", err)
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.LogSerious("err decrypt 1 %v", err)
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.LogSerious("err decrypt 2 %v", err)
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return nil, derror.NewErrorSingle("err:invalid_token", "reason", "cipher_message_authentication_failed")
		}
		log.LogSerious("err decrypt 3 %v", err)
		return nil, err
	}
	return plaintext, nil
}

func DecodeBase64Str(value string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.LogSerious("err decode base 64 1 %v %s", err, value)
		return nil
	}
	return decoded
}

func HashPassword(password string) (hash string) {
	bytePassword := []byte(password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func CompareHashedPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetHMAC(value, secretKey string) string {
	hexKey := []byte(secretKey)
	hexValue := []byte(value)
	mac := hmac.New(sha256.New, hexKey)
	mac.Write(hexValue)
	resultHMAC := mac.Sum(nil)
	// fmt.Println("ee", value, "sec", secretKey, "hex", hex.EncodeToString(resultHMAC), "expect", expectedHMACValue)
	return hex.EncodeToString(resultHMAC)
}

func CheckHMAC(value, expectedHMACValue, secretKey string) bool {
	hexKey := []byte(secretKey)
	hexValue := []byte(value)
	mac := hmac.New(sha256.New, hexKey)
	mac.Write(hexValue)
	resultHMAC := mac.Sum(nil)
	// fmt.Println("ee", value, "sec", secretKey, "hex", hex.EncodeToString(resultHMAC), "expect", expectedHMACValue)
	if hex.EncodeToString(resultHMAC) == expectedHMACValue {
		return true
	}
	return false
}
