package metadata_service

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "strconv"
    "time"
)

const kRedisKeyPrefixAppMetadataLastUpdatedTimestamp      = "metadata_last_up:"
const kRedisKeySharedMetadataLastUpdatedTimestamp   = "metadata_last_up:shared"

var lastUpdatedAppTimestamp int64
var lastUpdatedSharedTimestamp int64

func MarkMetadataAsUpdated(space metadata_typedefs.MetadataSpace) error {
    var key string
    if space == metadata_typedefs.METADATA_SPACE_SHARED {
        key = kRedisKeySharedMetadataLastUpdatedTimestamp
    } else {
        key = kRedisKeyPrefixAppMetadataLastUpdatedTimestamp + config.GetAppName()
    }

    timestamp := time.Now().Unix()

    err := redis_adaptor.Write(key, strconv.FormatInt(timestamp, 10))
    if err != nil {
        logger.LogError("error marking metadata as updated|" +
                        "|redis key=" + key +
                        "|timestamp=" + strconv.FormatInt(timestamp, 10) +
                        "|error=" + err.Error())
        return err
    }

    logger.LogInfo("marked metadata as updated" +
                   "|metadata space=" + space.String() +
                   "|timestamp=" + strconv.FormatInt(timestamp, 10))
    return nil
}

func RefreshLastUpdatedTimestamps() {
    lastUpdatedSharedTimestamp = time.Now().Unix()
    lastUpdatedAppTimestamp    = time.Now().Unix()
}

func startAutoUpdater(space metadata_typedefs.MetadataSpace) {
    RefreshLastUpdatedTimestamps()

    ticker := time.NewTicker(time.Second * time.Duration(config.GetEnvironmentConfiguration().MetadataAutoUpdaterPollSeconds))
    for range ticker.C {
        if !checkIfMetadataUpToDate(space) {
            logger.LogInfo("refresh triggered to re-fetch stale metadata" +
                           "|space=" + space.String())
            refreshMetadata()
        }
    }
}

func refreshMetadata() {
    defer ReleaseInstanceRW()
    _ = InstanceRW()

    RefreshLastUpdatedTimestamps()
}

func checkIfMetadataUpToDate(space metadata_typedefs.MetadataSpace) bool {
    var key string
    var lastUpdatedTimestamp int64
    if space == metadata_typedefs.METADATA_SPACE_SHARED {
        key = kRedisKeySharedMetadataLastUpdatedTimestamp
        lastUpdatedTimestamp = lastUpdatedSharedTimestamp
    } else {
        key = kRedisKeyPrefixAppMetadataLastUpdatedTimestamp + config.GetAppName()
        lastUpdatedTimestamp = lastUpdatedAppTimestamp
    }

    metadataEditedTimestampString, ok := redis_adaptor.Read(key)
    if !ok {
        // No last updated timestamp. We must be up to date
        return true
    }

    metadataEditedTimestamp, err := strconv.ParseInt(metadataEditedTimestampString, 10, 64)
    if err != nil {
        logger.LogError("error converting metadata last-updated-timestamp to string" +
                        "|metadata space=" + space.String() +
                        "|metadata edited timestamp=" + metadataEditedTimestampString +
                        "|last updated timestamp=" + strconv.FormatInt(lastUpdatedTimestamp, 10) +
                        "|error=" + err.Error())
        return false
    }

    if metadataEditedTimestamp > lastUpdatedTimestamp {
        return false
    }
    return true
}


