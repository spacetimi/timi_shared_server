package encryption_utils

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/base64"
    "encoding/hex"
    "errors"
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

