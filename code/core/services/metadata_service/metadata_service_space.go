package metadata_service

import (
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_fetchers"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
)

type MetadataServiceSpace struct {
    mdSpace metadata_typedefs.MetadataSpace
    mdVersionList *metadata_typedefs.MetadataVersionList
    mdFetcher metadata_typedefs.IMetadataFetcher

    mdCache map[string]*metadata_typedefs.MetadataCache			// version => MetadataCache for version
    mdManifests map[string]*metadata_typedefs.MetadataManifest 	// version => MetadataManifest for version
}

func newMetadataServiceSpace(metadataSpace metadata_typedefs.MetadataSpace) *MetadataServiceSpace {
    var err error

    msa := &MetadataServiceSpace{mdSpace: metadataSpace}

    msa.mdFetcher = metadata_fetchers.NewMetadataFetcher(metadataSpace)

    msa.mdVersionList, err = msa.mdFetcher.GetMetadataVersionList()
    if err != nil {
        logger.LogFatal("failed to load metadata version list|metadata_space=" + metadataSpace.String() +
                        "|error=" + err.Error())
    }

    // Load the manifest for every version specified in the version-list
    msa.mdManifests = make(map[string]*metadata_typedefs.MetadataManifest)
    for _, version := range msa.mdVersionList.Versions {
        manifest, err := msa.mdFetcher.GetMetadataManifestForVersion(version)
        if err != nil {
            logger.LogFatal("failed to load metadata manifest|metadata_space=" + metadataSpace.String() +
                            "|version=" + version +
                            "|error=" + err.Error())
        }
        msa.mdManifests[version] = manifest
    }

    // For versions marked as current in the version-list, load the metadata
    // and cache it in memory.
    // This will let us respond rapidly to metadata requests from current versions of the app
    msa.mdCache = make(map[string]*metadata_typedefs.MetadataCache)
    for _, currentVersion := range msa.mdVersionList.CurrentVersions {
        metadataCacheForVersion := metadata_typedefs.MetadataCache{}
        metadataCacheForVersion.Cache = make(map[string]string)

        msa.mdCache[currentVersion] = &metadataCacheForVersion

        manifestForVersion, ok := msa.mdManifests[currentVersion]
        if !ok {
            logger.LogFatal("failed to find manifest for a version marked as current" +
                            "|metadata_space=" + metadataSpace.String() +
                            "|version=" + currentVersion)
        }
        for _, manifestItem := range manifestForVersion.MetadataManifestItems {
            json, err := msa.mdFetcher.GetMetadataJsonByKey(manifestItem.MetadataKey, currentVersion)
            if err != nil {
                logger.LogFatal("failed to preload metadata|metadata_space=" + metadataSpace.String() +
                                "|version=" + currentVersion +
                                "|key=" + manifestItem.MetadataKey +
                                "|error=" + err.Error())
            }
            metadataCacheForVersion.Cache[manifestItem.MetadataKey] = json
        }
    }

    return msa
}

func (msa *MetadataServiceSpace) isMetadataHashUpToDate(key string, hash string, version *core.AppVersion) (bool, error) {
    if msa.mdVersionList.IsVersionValid(version) == false {
        return false, errors.New("invalid version")
    }

    manifest, ok := msa.mdManifests[version.String()]
    if !ok {
    	return false, errors.New("could not find manifest")
    }

    manifestItem := manifest.GetManifestItem(key)
    if manifestItem == nil {
        return false, errors.New("could not find manifest item")
    }

    return manifestItem.Hash == hash, nil
}

func (msa *MetadataServiceSpace) getMetadataJsonForItem(itemPtr metadata_typedefs.IMetadataItem, version *core.AppVersion) (string, error) {
    if msa.mdVersionList.IsVersionValid(version) == false {
        return "", errors.New("invalid version")
    }

    // If version is tagged under current versions, look inside metadata cache
    if msa.mdVersionList.IsVersionCurrent(version) {
        cachedMetadataForVersion, ok := msa.mdCache[version.String()]
        if ok {
            cachedMetadata, ok := cachedMetadataForVersion.Cache[itemPtr.GetKey()]
            if ok {
                return cachedMetadata, nil

            } else {
                logger.LogWarning("failed to find cached metadata item|metadata_space=" + msa.mdSpace.String() +
                                  "|version=" + version.String() +
                                  "|key=" + itemPtr.GetKey())
            }
        } else {
            logger.LogWarning("failed to find cached metadata|metadata_space=" + msa.mdSpace.String() +
                              "|version=" + version.String())
        }
    }

    // Else, load from fetcher

    metadataJson, err := msa.mdFetcher.GetMetadataJsonByKey(itemPtr.GetKey(), version.String())
    if err != nil {
        return "", errors.New("failed to find metadata from fetcher")
    }
    return metadataJson, nil
}

