package identity_service

import (
    "context"
    "errors"
    "github.com/dgrijalva/jwt-go"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/utils/encryption_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "strconv"
    "time"
)

var kCountersCollectionName string
var kCountersPrimaryKey string
var kCountersValueKey string
var kCountersPrimaryKeyValue string

// TODO: Get the following from somewhere else (preferably app scoped)
const kUserSessionExpirationTimeHours = 8
var kJwtSecretKey = []byte("8245EbB91cEa4808Ef064e1c53d0J8L6A8")

func Initialize() {
    kCountersCollectionName = "counters"
    kCountersPrimaryKey = "counter_name"
    kCountersValueKey = "counter_value"
    kCountersPrimaryKeyValue = config.GetAppName() + "::userid"

    // Use a delta of 0 to make sure the required counters-table document is set up
    _, err := mongo_adaptor.AtomicIncrement(mongo_adaptor.SHARED_DB,
                                            kCountersCollectionName,
                                            kCountersPrimaryKey,
                                            kCountersPrimaryKeyValue,
                                            kCountersValueKey,
                                            int64(0),
                                            context.Background())
    if err != nil {
        logger.LogError("error making sure counters table is set up: " + err.Error())
    }
}

func GetUserBlobById(userId int64, ctx context.Context) (*UserBlob, error) {
    user, err := loadUserBlobByUserId(userId, ctx)
    if err != nil {
        return nil, errors.New("error loading user blob: " + err.Error())
    }

    return user, nil
}

func CreateNewUserByUserNameAndPassword(userName string, password string, ctx context.Context) (*UserBlob, error) {
    uidm, err := loadUserNameToIdMappingByUserName(userName, ctx)
    if err == nil {
        return nil, errors.New("username \"" + userName + "\" already exists")
    }

    newUserId, err := createNewUserID()
    if err != nil {
        return nil, errors.New("error creating new user id: " + err.Error())
    }

    passwordHash, err := encryption_utils.HashAndSaltPassword(password)
    if err != nil {
        return nil, errors.New("error creating hash of password: " + err.Error())
    }

    newUserBlob := newUserBlob(newUserId)
    newUserBlob.CreatedTime = time.Now().Unix()
    newUserBlob.LastLoginTime = time.Now().Unix()
    newUserBlob.UserName = userName
    err = storage_service.SetBlob(newUserBlob, ctx)
    if err != nil {
        return nil, errors.New("error saving new user blob: " + err.Error())
    }

    uidm = newUserNameToIdMapping(userName)
    uidm.PasswordHash = passwordHash
    uidm.UserId = newUserId
    err = storage_service.SetBlob(uidm, ctx)
    if err != nil {
        return nil, errors.New("error saving new user name to id mapping: " + err.Error())
    }

    return newUserBlob, nil
}

func CheckAndGetUserBlobFromUserLoginCredentials(userName string, password string, ctx context.Context) (*UserBlob, error) {
    uidm, err := loadUserNameToIdMappingByUserName(userName, ctx)
    if err != nil {
        return nil, errors.New("error loading user name to id mapping: " + err.Error())
    }

    passwordOk := encryption_utils.VerifyPasswordWithHash(password, uidm.PasswordHash)
    if !passwordOk {
        return nil, errors.New("password doesn't match")
    }

    user, err := loadUserBlobByUserId(uidm.UserId, ctx)
    if err != nil {
        return nil, errors.New("error loading user blob for id: " + strconv.FormatInt(uidm.UserId, 10))
    }

    return user, nil
}

func CreateUserLoginToken(user *UserBlob) (string, error) {
    expiration := time.Now().Add(kUserSessionExpirationTimeHours * time.Hour)
    claims := &UserJWTClaims{
        UserId: user.UserId,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt:expiration.Unix(),
        },
    }

    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    jwtTokenString, err := jwtToken.SignedString(kJwtSecretKey)

    if err != nil {
        return "", errors.New("error creating user jwt token string: " + err.Error())
    }
    return jwtTokenString, nil
}

func CheckAndGetUserBlobFromUserLoginToken(jwtTokenString string, ctx context.Context) (*UserBlob, error) {
    claims := &UserJWTClaims{}

    jwtToken, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return kJwtSecretKey, nil
    })
    if err != nil {
        return nil, errors.New("error parsing user jwt token string: " + err.Error())
    }
    if !jwtToken.Valid {
        return nil, errors.New("invalid user jwt token")
    }

    user, err := loadUserBlobByUserId(claims.UserId, ctx)
    if err != nil {
        return nil, errors.New("error loading user from claim: " + err.Error() +
                               "|user id: " + strconv.FormatInt(claims.UserId, 10))
    }

    return user, nil
}

func createNewUserID() (int64, error) {

    // TODO: Don't use context: background
    newUserIdInterface, err := mongo_adaptor.AtomicIncrement(mongo_adaptor.SHARED_DB,
                                                             kCountersCollectionName,
                                                             kCountersPrimaryKey,
                                                             kCountersPrimaryKeyValue,
                                                             kCountersValueKey, 1,
                                                             context.Background())
    if err != nil {
        return -1, errors.New("error creating new UserId: " + err.Error())
    }

    newUserId, ok := newUserIdInterface.(int64)
    if !ok {
        return -1, errors.New("failed type assertion when creating new UserId")
    }

    return newUserId, nil
}

