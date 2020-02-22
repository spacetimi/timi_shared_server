package storage_service

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
)

func GetBlobByPrimaryKeys(primaryKeyValues []interface{},
                          outBlobPtr storage_typedefs.IBlob,
                          ctx context.Context) error {

    if outBlobPtr == nil {
        return errors.New("blob ptr is nil")
    }

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
