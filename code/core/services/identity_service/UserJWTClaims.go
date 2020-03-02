package identity_service

import "github.com/dgrijalva/jwt-go"

type UserJWTClaims struct {
    UserId int64
    jwt.StandardClaims
}
