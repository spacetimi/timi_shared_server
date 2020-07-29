package admin

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/utils/aws_helper"
)

const kCredentialsSecretName = "admin_login_creds"

const kCredentialsForLocalhostUsername = "admin"
const kCredentialsForLocalhostPasswd = "spacetimi1!"

const sessionExpirationTimeHours = 8

var kJWT_SECRET_KEY = []byte("b5128245ebb91cea4808ef064e1c53d0")

type AdminLoginRequest struct {
	username string
	password string
}

type AdminLoginResponse struct {
	jwtTokenString string
	expirationTime time.Time
}

type AdminUserJWTClaims struct {
	Username    string `json:"username"`
	IsAdminUser bool
	jwt.StandardClaims
}

func tryLoginWithAdminCredentials(request *AdminLoginRequest) (*AdminLoginResponse, error) {

	passwd, err := getAdminPasswordForUsername(request.username)
	if err != nil {
		return nil, errors.New("error validating credentials: " + err.Error())
	}

	if passwd != request.password {
		return nil, errors.New("wrong password")
	}

	expiration := time.Now().Add(sessionExpirationTimeHours * time.Hour)

	claims := &AdminUserJWTClaims{
		Username:    request.username,
		IsAdminUser: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the jwtClaims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	jwtTokenString, err := jwtToken.SignedString(kJWT_SECRET_KEY)
	if err != nil {
		return nil, errors.New("error creating jwt token string: " + err.Error())
	}

	response := AdminLoginResponse{
		jwtTokenString: jwtTokenString,
		expirationTime: expiration,
	}

	return &response, nil
}

func getAdminPasswordForUsername(username string) (string, error) {
	if config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL {
		if username != kCredentialsForLocalhostUsername {
			return "", errors.New("not an admin user")
		}
		return kCredentialsForLocalhostPasswd, nil
	}

	passwd, err := aws_helper.ReadJsonSecret(kCredentialsSecretName, username)
	if err != nil {
		return "", errors.New("probably not an admin user. error fetching credentials-passwd from aws secret: " + err.Error())
	}

	return passwd, nil
}

func checkAdminLoginClaim(jwtTokenString string) (bool, string, error) {
	claims := &AdminUserJWTClaims{}

	token, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return kJWT_SECRET_KEY, nil
	})
	if err != nil {
		return false, "", errors.New("error parsing jwt token string: " + err.Error())
	}
	if !token.Valid {
		return false, "", errors.New("invalid jwt token")
	}
	if !claims.IsAdminUser {
		return false, "", errors.New("claim not admin user")
	}

	return true, claims.Username, nil
}
