package identity_service

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
    "github.com/spacetimi/timi_shared_server/utils/logger"
)

var kCountersCollectionName string
var kCountersPrimaryKey string
var kCountersValueKey string
var kCountersPrimaryKeyValue string

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

func CreateNewUserID() (int64, error) {

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

