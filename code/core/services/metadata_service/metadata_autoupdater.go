package metadata_service

import (
	"context"
	"strconv"
	"time"

	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/v2/utils/logger"
)

const kRedisKeyPrefixAppMetadataLastUpdatedTimestamp = "metadata_last_up:"
const kRedisKeySharedMetadataLastUpdatedTimestamp = "metadata_last_up:shared"

var lastUpdatedAppTimestamp int64
var lastUpdatedSharedTimestamp int64

func MarkMetadataAsUpdated(space metadata_typedefs.MetadataSpace, ctx context.Context) error {
	var key string
	if space == metadata_typedefs.METADATA_SPACE_SHARED {
		key = kRedisKeySharedMetadataLastUpdatedTimestamp
	} else {
		key = kRedisKeyPrefixAppMetadataLastUpdatedTimestamp + config.GetAppName()
	}

	timestamp := time.Now().Unix()

	err := redis_adaptor.Write(key, strconv.FormatInt(timestamp, 10), 7*24*time.Hour, ctx)
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

func RefreshMetadata() {
	defer ReleaseInstanceRW()
	_ = InstanceRW()

	RefreshLastUpdatedTimestamps()
}

func RefreshLastUpdatedTimestamps() {
	lastUpdatedSharedTimestamp = time.Now().Unix()
	lastUpdatedAppTimestamp = time.Now().Unix()
}

func CheckIfMetadataUpToDate(space metadata_typedefs.MetadataSpace, ctx context.Context) bool {
	var key string
	var lastUpdatedTimestamp int64
	if space == metadata_typedefs.METADATA_SPACE_SHARED {
		key = kRedisKeySharedMetadataLastUpdatedTimestamp
		lastUpdatedTimestamp = lastUpdatedSharedTimestamp
	} else {
		key = kRedisKeyPrefixAppMetadataLastUpdatedTimestamp + config.GetAppName()
		lastUpdatedTimestamp = lastUpdatedAppTimestamp
	}

	metadataEditedTimestampString, ok := redis_adaptor.Read(key, ctx)
	if !ok {
		// No last updated timestamp. Create the timestamp as now, and Force refresh to be safe
		err := MarkMetadataAsUpdated(space, ctx)
		if err != nil {
			logger.LogError("error trying to force mark metadata as updated" +
				"|space=" + space.String() +
				"|error=" + err.Error())
		}
		return false
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

func startAutoUpdater(space metadata_typedefs.MetadataSpace, ctx context.Context) {
	RefreshLastUpdatedTimestamps()

	ticker := time.NewTicker(time.Second * time.Duration(config.GetEnvironmentConfiguration().MetadataAutoUpdaterPollSeconds))
	checkAndAutoUpdateMetadata(space, ctx)
	for range ticker.C {
		checkAndAutoUpdateMetadata(space, ctx)
	}
}

func checkAndAutoUpdateMetadata(space metadata_typedefs.MetadataSpace, ctx context.Context) {
	if !CheckIfMetadataUpToDate(space, ctx) {
		logger.LogInfo("refresh triggered to re-fetch stale metadata" +
			"|space=" + space.String())
		RefreshMetadata()
	}
}
