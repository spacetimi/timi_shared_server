package admin

import (
    "errors"
    "github.com/dgrijalva/jwt-go"
    "time"
)

var hardcoded_credentials = map[string]string {
    "admin": "spacetimi1!",
}

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

type JWTClaims struct {
    Username string `json:"username"`
    IsAdminUser bool
    jwt.StandardClaims
}

func tryLoginWithAdminCredentials(request *AdminLoginRequest) (*AdminLoginResponse, error) {

    // TODO: Avi: Bypass some checks in local environment

    // TODO: Avi: Don't check against hardcoded hardcoded_credentials, use some other directory service instead
    reqdPassword, ok := hardcoded_credentials[request.username]
    if !ok {
        return nil, errors.New("unknown user")
    }

    if reqdPassword != request.password {
        return nil, errors.New("wrong password")
    }

    expiration := time.Now().Add(sessionExpirationTimeHours * time.Hour)

    claims := &JWTClaims{
        Username: request.username,
        IsAdminUser: true,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt:expiration.Unix(),
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

func checkAdminLoginClaim(jwtTokenString string) (bool, string, error) {
    claims := &JWTClaims{}

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
