package identity_service

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/dgrijalva/jwt-go"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/utils/encryption_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "os"
    "strconv"
    "time"
)

var kCountersCollectionName string
var kCountersPrimaryKey string
var kCountersValueKey string
var kCountersPrimaryKeyValue string

const kConfigFileName = "identity_service_config.json"
type Config_t struct {
    UserSessionExpiryHours int
    JWTSecretKey string
}
var Config *Config_t

func Initialize() {

    // Read config file
    configFilePath := config.GetAppConfigFilesPath() + "/" + kConfigFileName
    configFile, err := os.Open(configFilePath)
    if err != nil {
        logger.LogFatal("cannot open configuration file|file path=" + configFilePath)
        return
    }
    defer func() {
        err := configFile.Close()
        if err != nil {
            logger.LogError("error closing config file" +
                            "|file path=" + configFilePath +
                            "|error=" + err.Error())
        }
    }()

    decoder := json.NewDecoder(configFile)
    err = decoder.Decode(&Config)
    if err != nil {
        logger.LogFatal("error decoding configuration file" +
                        "|file path=" + configFilePath +
                        "|error=" + err.Error())
        return
    }

    kCountersCollectionName = "counters"
    kCountersPrimaryKey = "counter_name"
    kCountersValueKey = "counter_value"
    kCountersPrimaryKeyValue = config.GetAppName() + "::userid"

    // Use a delta of 0 to make sure the required counters-table document is set up
    _, err = mongo_adaptor.AtomicIncrement(mongo_adaptor.SHARED_DB,
                                            kCountersCollectionName,
                                            kCountersPrimaryKey,
                                            kCountersPrimaryKeyValue,
                                            kCountersValueKey,
                                            int64(0),
                                            context.Background())
    if err != nil {
        logger.LogFatal("error making sure counters table is set up|error=" + err.Error())
    }
}

////////////////////////////////////////////////////////////////////////////////
// Public API:

func GetUserBlobById(userId int64, ctx context.Context) (*UserBlob, error) {
    user, err := loadUserBlobByUserId(userId, ctx)
    if err != nil {
        return nil, errors.New("error loading user blob: " + err.Error())
    }

    return user, nil
}

func CreateNewUser(userName string, password string, userEmailAddress string, ctx context.Context) (*UserBlob, error) {
    uidm, err := loadUserNameToIdMappingByUserName(userName, ctx)
    if err == nil {
        return nil, errors.New("username \"" + userName + "\" already exists")
    }

    if userEmailAddress != "" {
        _, err := loadUserEmailToIdMappingByUserEmail(userEmailAddress, ctx)
        if err == nil {
            return nil, errors.New("email address \"" + userEmailAddress + "\" already in use")
        }
    }

    newUserId, err := createNewUserID()
    if err != nil {
        return nil, errors.New("error creating new user id: " + err.Error())
    }

    passwordHash, err := encryption_utils.HashAndSaltPassword(password)
    if err != nil {
        return nil, errors.New("error creating hash of password: " + err.Error())
    }

    if userEmailAddress != "" {
        ueidm := newUserEmailToIdMapping(userEmailAddress)
        ueidm.UserId = newUserId
        err = storage_service.SetBlob(ueidm, ctx)
        if err != nil {
            return nil, errors.New("error saving new user email to id mapping: " + err.Error())
        }
    }

    newUserBlob := newUserBlob(newUserId)
    newUserBlob.CreatedTime = time.Now().Unix()
    newUserBlob.LastLoginTime = time.Now().Unix()
    newUserBlob.UserName = userName
    newUserBlob.UserEmailAddress = userEmailAddress
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

func UpdateUserPassword(user *UserBlob, password string, ctx context.Context) error {

    passwordHash, err := encryption_utils.HashAndSaltPassword(password)
    if err != nil {
        return errors.New("error creating hash of password: " + err.Error())
    }

    uidm, err := loadUserNameToIdMappingByUserName(user.UserName, ctx)
    if err != nil {
        logger.LogError("error loading user name to id mapping" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|user name=" + user.UserName +
                        "|error=" + err.Error())
        return errors.New("error getting user name to id mapping: " + err.Error())
    }

    uidm.PasswordHash = passwordHash
    err = storage_service.SetBlob(uidm, ctx)
    if err != nil {
        logger.LogError("error saving user name to id mapping" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|user name=" + user.UserName +
                        "|error=" + err.Error())
        return errors.New("error saving new user name to id mapping: " + err.Error())
    }

    return nil
}

func SetUserEmailAddressVerified(userId int64, ctx context.Context) error {
    user, err := GetUserBlobById(userId, ctx)
    if err != nil {
        return err
    }

    if user.UserEmailAddressVerified {
        return nil
    }

    user.UserEmailAddressVerified = true
    err = storage_service.SetBlob(user, ctx)
    if err != nil {
        return errors.New("error saving user blob: " + err.Error())
    }

    return nil
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

func CheckAndGetUserBlobFromUserEmailAddress(userEmailAddress string, ctx context.Context) (*UserBlob, error) {
    ueidm, err := loadUserEmailToIdMappingByUserEmail(userEmailAddress, ctx)
    if err != nil {
        return nil, errors.New("error loading user email address to id mapping: " + err.Error())
    }

    user, err := loadUserBlobByUserId(ueidm.UserId, ctx)
    if err != nil {
        return nil, errors.New("error loading user blob for id: " + strconv.FormatInt(ueidm.UserId, 10))
    }

    return user, nil
}

func CreateUserLoginToken(user *UserBlob) (string, error) {
    expiration := time.Now().Add(time.Duration(Config.UserSessionExpiryHours) * time.Hour)
    claims := &UserJWTClaims{
        UserId: user.UserId,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt:expiration.Unix(),
        },
    }

    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    jwtTokenString, err := jwtToken.SignedString([]byte(Config.JWTSecretKey))

    if err != nil {
        return "", errors.New("error creating user jwt token string: " + err.Error())
    }
    return jwtTokenString, nil
}

func CheckAndGetUserBlobFromUserLoginToken(jwtTokenString string, ctx context.Context) (*UserBlob, error) {
    claims := &UserJWTClaims{}

    jwtToken, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(Config.JWTSecretKey), nil
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

