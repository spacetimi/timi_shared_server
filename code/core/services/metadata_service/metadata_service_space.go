package metadata_service

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_fetchers"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
)

type MetadataServiceSpace struct {
    mdSpace metadata_typedefs.MetadataSpace
    mdVersionList *metadata_typedefs.MetadataVersionList
    mdFetcher metadata_typedefs.IMetadataFetcher

    /* Cache of metadata for versions marked as current */
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

    // TODO: Avi: Rename current versions to cached versions
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

func (msa *MetadataServiceSpace) getMetadataManifestForVersion(version *core.AppVersion) (*metadata_typedefs.MetadataManifest, error) {
    if msa.mdVersionList.IsVersionValid(version) == false {
        return nil, errors.New("invalid version")
    }

    manifest, ok := msa.mdManifests[version.String()]
    if !ok {
        return nil, errors.New("could not find manifest for version")
    }

    return manifest, nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (msa *MetadataServiceSpace) setMetadataManifestForVersion(manifest *metadata_typedefs.MetadataManifest, version *core.AppVersion) error {
    if msa.mdVersionList.IsVersionValid(version) == false {
        return errors.New("invalid version")
    }

    msa.mdManifests[version.String()] = manifest

    err :=  msa.mdFetcher.SetMetadataManifestForVersion(manifest, version.String())
    if err != nil {
        return errors.New("error saving new manifest|error=" + err.Error())
    }

    return nil
}

func (msa *MetadataServiceSpace) isMetadataHashUpToDate(key string, hash string, version *core.AppVersion) (bool, error) {
    manifest, err := msa.getMetadataManifestForVersion(version)
    if err != nil {
        return false, err
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

/**
 * Only meant to be called from the admin tool / scripts
 */
func (msa *MetadataServiceSpace) setMetadataJsonForItem(itemPtr metadata_typedefs.IMetadataItem, version *core.AppVersion) error {
    if msa.mdVersionList.IsVersionValid(version) == false {
        return errors.New("invalid version")
    }

    var err error
    var metadataJsonBytes []byte
    if config.GetEnvironmentConfiguration().AppEnvironment == config.PRODUCTION {
        metadataJsonBytes, err = json.Marshal(itemPtr)
    } else {
        metadataJsonBytes, err = json.MarshalIndent(itemPtr, "", "    ")
    }
    if err != nil {
        return errors.New("error deserializing metadata|error=" + err.Error())
    }
    metadataJson := string(metadataJsonBytes)

    // If version is marked as current, also update the cache
    if msa.mdVersionList.IsVersionCurrent(version) {

        cachedMetadataForVersion, ok := msa.mdCache[version.String()]
        if !ok {
            logger.LogError("failed to find cached metadata while writing" +
                            "|metadata space=" + itemPtr.GetMetadataSpace().String() +
                            "|metadata item key=" + itemPtr.GetKey() +
                            "|version=" + version.String())
            return errors.New("failed to find cached metadata for version (this should theoretically never happen)")
        }
        cachedMetadataForVersion.Cache[itemPtr.GetKey()] = metadataJson
    }

    // Save the metadata item
    err = msa.mdFetcher.SetMetadataJsonByKey(itemPtr.GetKey(), metadataJson, version.String())
    if err != nil {
        return errors.New("error saving metadata json|error=" + err.Error())
    }

    // Generate hash of the file's contents
    hashBytes := md5.Sum([]byte(metadataJson))
    hash := hex.EncodeToString(hashBytes[:])

    // Update the metadata manifest for this version
    manifest, err := msa.getMetadataManifestForVersion(version)
    if err != nil {
        return errors.New("error getting metadata manifest for version|error=" + err.Error())
    }
    manifest.SetManifestItem(itemPtr.GetKey(), hash)
    err = msa.setMetadataManifestForVersion(manifest, version)
    if err != nil {
        return errors.New("error updating metadata manifest for version|error=" + err.Error())
    }

    return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (msa *MetadataServiceSpace) setCurrentVersions(newCurrentVersions []*core.AppVersion) error {
    if len(newCurrentVersions) == 0 {
        return errors.New("new current versions list cannot be empty")
    }

    msa.mdVersionList.CurrentVersions = nil
    for _, version := range newCurrentVersions {
        if version == nil {
            return errors.New("version cannot be null")
        }
        if !msa.mdVersionList.IsVersionValid(version) {
            return errors.New("invalid version: " + version.String())
        }
        msa.mdVersionList.CurrentVersions = append(msa.mdVersionList.CurrentVersions, version.String())
    }

    err := msa.mdFetcher.SetMetadataVersionList(msa.mdVersionList)
    if err != nil {
        return errors.New("could not save current metadata versions| error=" + err.Error())
    }

    return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (msa *MetadataServiceSpace) createNewVersion(version *core.AppVersion, markAsCurrent bool) error {
    _, ok := msa.mdManifests[version.String()]
    if ok {
        return errors.New("duplicate version")
    }

    err := msa.mdVersionList.CreateNewVersion(version, markAsCurrent)
    if err != nil {
        return errors.New("error adding new version to metadata version list: " + err.Error())
    }

    err = msa.mdFetcher.SetMetadataVersionList(msa.mdVersionList)
    if err != nil {
        return errors.New("error saving updated metadata version list: " + err.Error())
    }

    // Create empty manifest for new version
    newManifest := &metadata_typedefs.MetadataManifest{}
    newManifest.Initialize()
    msa.mdManifests[version.String()] = newManifest

    err = msa.mdFetcher.SetMetadataManifestForVersion(newManifest, version.String())
    if err != nil {
        return errors.New("error creating metadata manifest for new version: " + err.Error())
    }

    return nil
}
