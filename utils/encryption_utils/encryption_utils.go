package encryption_utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"github.com/spacetimi/timi_shared_server/v2/utils/string_utils"
	"golang.org/x/crypto/bcrypt"
)

////////////////////////////////////////////////////////////////////////////////

func HashAndSaltPassword(password string) (string, error) {
	// TODO: Get cost from config instead of bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error generating hash of password: " + err.Error())
	}

	return string(hash), nil
}

func VerifyPasswordWithHash(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

////////////////////////////////////////////////////////////////////////////////

func Generate_md5_hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))

	return hex.EncodeToString(hasher.Sum(nil))
}

func Generate_sha_hash(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

////////////////////////////////////////////////////////////////////////////////

func EncryptUsingAES(data string, key string) (string, error) {

	keyHash := Generate_md5_hash(key)
	block, err := aes.NewCipher([]byte(keyHash))
	if err != nil {
		return "", errors.New("error creating aes cipher for encryption: " + err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("error creating gcm for cipher for encryption: " + err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("error preparing nonce from rand for encryption: " + err.Error())
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(data), nil)
	return string_utils.EncodeBytesAsBase64String(cipherText), nil
}

func DecryptUsingAES(data string, key string) (string, error) {

	dataAsBytes, err := string_utils.DecodeBase64StringToBytes(data)
	if err != nil {
		return "", errors.New("error decoding secret from base 64: " + err.Error())
	}

	keyHashBytes := []byte(Generate_md5_hash(key))
	block, err := aes.NewCipher(keyHashBytes)
	if err != nil {
		return "", errors.New("error creating aes cipher for decryption: " + err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("error creating gcm for cipher for decryption: " + err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := dataAsBytes[:nonceSize], dataAsBytes[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", errors.New("error decrypting to plainText for decryption: " + err.Error())
	}

	return string(plainText), nil
}
