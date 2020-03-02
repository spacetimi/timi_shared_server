package encryption_utils

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
)

func HashAndSaltPassword(password string) (string, error) {
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
