package metadata_fetchers

import (
	"github.com/spacetimi/server/timi_shared/code/config"
	"github.com/spacetimi/server/timi_shared/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/server/timi_shared/utils/logger"
)

func NewMetadataFetcher(metadataSpace metadata_typedefs.MetadataSpace) metadata_typedefs.IMetadataFetcher {
	var fetcher metadata_typedefs.IMetadataFetcher

	if metadataSpace == metadata_typedefs.METADATA_SPACE_SHARED {
		if config.GetEnvironmentConfiguration().SharedMetadataSourceURL == "local" {
			fetcher = NewMetadataFetcherFilesystem(config.GetSharedMetadataFilesPath())
		} else {
			fetcher = NewMetadataFetcherS3(config.GetEnvironmentConfiguration().SharedMetadataSourceURL)
		}
	}

	if metadataSpace == metadata_typedefs.METADATA_SPACE_APP {
		if config.GetEnvironmentConfiguration().AppMetadataSourceURL == "local" {
			fetcher = NewMetadataFetcherFilesystem(config.GetAppMetadataFilesPath())
		} else {
			fetcher = NewMetadataFetcherS3(config.GetEnvironmentConfiguration().AppMetadataSourceURL)
		}
	}

	if fetcher == nil {
		logger.LogFatal("Unable to create metadata fetcher|metadataspace=" + string(metadataSpace))
		return nil
	}

	return fetcher
}
