package metadata_fetchers

import (
	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/v2/utils/logger"
)

func NewMetadataFetcher(metadataSpace metadata_typedefs.MetadataSpace) metadata_typedefs.IMetadataFetcher {
	var fetcher metadata_typedefs.IMetadataFetcher

	if metadataSpace == metadata_typedefs.METADATA_SPACE_SHARED {
		if config.GetEnvironmentConfiguration().SharedMetadataSourceURL == "local" {
			fetcher = NewMetadataFetcherFilesystem(config.GetSharedMetadataFilesPath())
		} else {
			fetcher = NewMetadataFetcherS3(config.GetEnvironmentConfiguration().SharedMetadataSourceURL, config.GetEnvironmentConfiguration().AdminToolConfig.SharedMetadataS3BucketName)
		}
	}

	if metadataSpace == metadata_typedefs.METADATA_SPACE_APP {
		if config.GetEnvironmentConfiguration().AppMetadataSourceURL == "local" {
			fetcher = NewMetadataFetcherFilesystem(config.GetAppMetadataFilesPath())
		} else {
			fetcher = NewMetadataFetcherS3(config.GetEnvironmentConfiguration().AppMetadataSourceURL, config.GetEnvironmentConfiguration().AdminToolConfig.AppMetadataS3BucketName)
		}
	}

	if fetcher == nil {
		logger.LogFatal("Unable to create metadata fetcher|metadataspace=" + string(metadataSpace))
		return nil
	}

	return fetcher
}
