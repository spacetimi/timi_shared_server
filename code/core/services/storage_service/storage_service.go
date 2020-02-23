package storage_service

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "reflect"
)

func GetBlobByPrimaryKeys(outBlobPtr storage_typedefs.IBlob,
                          ctx context.Context) error {

    if outBlobPtr == nil {
        return errors.New("blob ptr is nil")
    }

    primaryKeyValues, err := getPrimaryKeyValuesFromBlob(outBlobPtr)
    if err != nil {
        return errors.New("error getting primary key values from blob: " + err.Error())
    }

    // Check in redis first
    redisKey := getRedisKey(outBlobPtr.GetBlobName(), primaryKeyValues)
    redisValue, redisOk := redis_adaptor.Read(redisKey)
    if redisOk {
        err := json.Unmarshal([]byte(redisValue), outBlobPtr)
        if err == nil {
            // Successfully read the blob from redis
            return nil
        }
        logger.LogError("error deserializing blob from redis" +
                        "|blob name=" + outBlobPtr.GetBlobName() +
                        "|primary key values=" + fmt.Sprintf("%#v", primaryKeyValues) +
                        "|error=" + err.Error())
        // Fall-through
    }

    // Read the blob from DB
    dbSpace, err := getDBSpaceFromStorageSpace(outBlobPtr.GetStorageSpace())
    if err != nil {
        return errors.New("error resolving db space: " + err.Error())
    }

    collectionName := outBlobPtr.GetBlobName()
    primaryKeys := outBlobPtr.GetPrimaryKeys()

    err = mongo_adaptor.GetDataItemByPrimaryKeys(dbSpace, collectionName, primaryKeys, primaryKeyValues, outBlobPtr, ctx)
    if err != nil {
        return errors.New("error getting blob from db: " + err.Error())
    }

    // Write the blob to redis for faster reads next time
    err = writeBlobToRedis(redisKey, outBlobPtr)
    if err != nil {
        logger.LogError("error saving blob to redis" +
                        "|blob name=" + outBlobPtr.GetBlobName() +
                        "|primary key values=" + fmt.Sprintf("%#v", primaryKeyValues) +
                        "|error=" + err.Error())
        // Fall-through
    }

    return nil
}

func SetBlob(blobPtr storage_typedefs.IBlob, ctx context.Context) error {

    if blobPtr == nil {
        return errors.New("blob ptr is nil")
    }

    dbSpace, err := getDBSpaceFromStorageSpace(blobPtr.GetStorageSpace())
    if err != nil {
        return errors.New("error resolving db space: " + err.Error())
    }

    collectionName := blobPtr.GetBlobName()
    primaryKeys := blobPtr.GetPrimaryKeys()

    err = mongo_adaptor.WriteDataItemByPrimaryKeys(dbSpace, collectionName, primaryKeys, blobPtr, ctx)
    if err != nil {
        return errors.New("error writing blob to db: " + err.Error())
    }

    // Also write the blob to redis
    primaryKeyValues, err := getPrimaryKeyValuesFromBlob(blobPtr)
    if err != nil {
        logger.LogError("error getting primary key values while trying to save blob to redis" +
                        "|blob name=" + blobPtr.GetBlobName() +
                        "|error=" + err.Error())
        // Fall-through
    } else {
        redisKey := getRedisKey(blobPtr.GetBlobName(), primaryKeyValues)
        err = writeBlobToRedis(redisKey, blobPtr)
        if err != nil {
            logger.LogError("error saving blob to redis" +
                            "|blob name=" + blobPtr.GetBlobName() +
                            "|primary key values=" + fmt.Sprintf("%#v", primaryKeyValues) +
                            "|error=" + err.Error())
            // Fall-through
        }
    }

    return nil
}

/***** Private ******************************************************************/

func getDBSpaceFromStorageSpace(storageSpace storage_typedefs.StorageSpace) (mongo_adaptor.DBSpace, error) {
    switch storageSpace {
        case storage_typedefs.STORAGE_SPACE_SHARED: return mongo_adaptor.SHARED_DB, nil
        case storage_typedefs.STORAGE_SPACE_APP: return mongo_adaptor.APP_DB, nil
    }
    return -1, errors.New("invalid storage space")
}

func getRedisKey(blobName string, primaryKeyValues []interface{}) string {
    redisKey := blobName
    for _, value := range primaryKeyValues {
        var valueAsString string
        switch value.(type) {
        case int, int32, int64:
            valueAsString = fmt.Sprintf("%d", value)
        case bool:
            valueAsString = fmt.Sprintf("%t", value)
        case string:
            valueAsString = fmt.Sprintf("%s", value)
        case float32, float64:
            valueAsString = fmt.Sprintf("%f", value)
        default:
            logger.LogWarning("unsupported primary key type while forming redis key" +
                              "|blob name=" + blobName +
                              "|type=" + reflect.TypeOf(value).Name())
            // Fall back to binary so that it doesn't break
            valueAsString = fmt.Sprintf("%x", value)
        }
        redisKey = redisKey + ":" + valueAsString
    }
    return redisKey
}

func writeBlobToRedis(redisKey string, blobPtr storage_typedefs.IBlob) error {
    bytes, err := json.Marshal(blobPtr)
    if err != nil {
        return errors.New("error serializing blob: " + err.Error())
    }

    err = redis_adaptor.Write(redisKey, string(bytes))
    if err != nil {
        return errors.New("error writing blob to redis: " + err.Error())
    }

    return nil
}

func getPrimaryKeyValuesFromBlob(blobPtr storage_typedefs.IBlob) ([]interface{}, error) {
    var primaryKeyValues []interface{}
    primaryKeys := blobPtr.GetPrimaryKeys()


    p := reflect.ValueOf(blobPtr)
    if p.Kind() != reflect.Ptr {
        return nil, errors.New("not a struct ptr")
    }

    v := reflect.Indirect(p)
    if v.Kind() != reflect.Struct {
        return nil, errors.New("not a struct")
    }

    for _, primaryKey := range primaryKeys {
        value := v.FieldByName(primaryKey)
        if value.CanInterface() {
            primaryKeyValues = append(primaryKeyValues, value.Interface())
        } else {
            return nil, errors.New("inaccessible primary key field: " + primaryKey)
        }
    }

    return primaryKeyValues, nil
}
