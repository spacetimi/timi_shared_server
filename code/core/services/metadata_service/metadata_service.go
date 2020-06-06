package metadata_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/code/core"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"sync"
)

type MetadataService struct {
	sharedMDServiceSpace *MetadataServiceSpace
	appMDServiceSpace    *MetadataServiceSpace
}

func Initialize() {
    // This is during Initialization. No need to take mutex lock
	instance = createInstance()
	go startAutoUpdater(metadata_typedefs.METADATA_SPACE_SHARED, context.Background())
	go startAutoUpdater(metadata_typedefs.METADATA_SPACE_APP, context.Background())
}

func createInstance() *MetadataService {
	newInstance := &MetadataService{}
	newInstance.sharedMDServiceSpace = newMetadataServiceSpace(metadata_typedefs.METADATA_SPACE_SHARED)
	newInstance.appMDServiceSpace    = newMetadataServiceSpace(metadata_typedefs.METADATA_SPACE_APP)

	return newInstance
}

var instance *MetadataService
var mutexForInstance sync.RWMutex

func Instance() *MetadataService {
    mutexForInstance.RLock()
	defer mutexForInstance.RUnlock()

	if instance == nil {
		logger.LogError("Metadata Service instance is null")
	}
	return instance
}

/**
 * Intended if you want to update the metadata
 * Only meant to be called from the admin tool / scripts
 *
 * This will create a new copy of instance so that any other requests
 * that are currently working with the old copy of instance
 * will go through using the old copy
 * Any requests that try to get a copy of instance while it is being modified here will have to wait
 * till the RW-instance is released
 */
func InstanceRW() *MetadataService {
	mutexForInstance.Lock()
	logger.LogInfo("Taking write lock on metadata service instance")

	instance = createInstance()
	return instance
}
/**
 * MUST be called after taking an InstanceRW
 * Failing to call this will lead to deadlocks
 */
func ReleaseInstanceRW() {
	mutexForInstance.Unlock()
	logger.LogInfo("Released write lock on metadata service instance")
}


func (ms *MetadataService) GetCurrentVersions(space metadata_typedefs.MetadataSpace) []string {
    msa := ms.getMetadataServiceSpace(space)

    return msa.mdVersionList.CurrentVersions
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (ms *MetadataService) SetCurrentVersions(newCurrentVersionStrings []string, space metadata_typedefs.MetadataSpace) error {
	if len(newCurrentVersionStrings) == 0 {
		return nil
	}

	var newCurrentVersions []*core.AppVersion
	for _, versionString := range newCurrentVersionStrings {
	    version, err := core.GetAppVersionFromString(versionString)
	    if err != nil {
	    	logger.LogError("Error parsing app version" +
							"|metadata space=" + space.String() +
	    					"|version string=" + versionString +
	    					"|error=" + err.Error())
	    	return errors.New("error parsing versions: " + err.Error())
		}
		newCurrentVersions = append(newCurrentVersions, version)
	}

	if len(newCurrentVersions) == 0 {
		return errors.New("error parsing versions")
	}

	msa := ms.getMetadataServiceSpace(space)

	err := msa.setCurrentVersions(newCurrentVersions)
	if err != nil {
		logger.LogError("error updating current versions" +
						"|metadata space=" + space.String() +
						"|error=" + err.Error())
		return errors.New("couldn't update current versions: " + err.Error())
	}

    return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (ms *MetadataService) CreateNewVersion(newVersion *core.AppVersion, space metadata_typedefs.MetadataSpace, markAsCurrent bool) error {
	validVersion, _ := ms.IsVersionValid(newVersion.String(), space)
	if validVersion {
		return errors.New("duplicate version. version already exists")
	}

	msa := ms.getMetadataServiceSpace(space)
	err := msa.createNewVersion(newVersion, markAsCurrent)
	if err != nil {
		logger.LogWarning("error creating new metadata version" +
						  "|metadata space=" + space.String() +
						  "|new version=" + newVersion.String() +
						  "|error=" + err.Error())
		return err
	}

	return nil
}

func (ms *MetadataService) GetAllVersions(space metadata_typedefs.MetadataSpace) []string {
	msa := ms.getMetadataServiceSpace(space)

	return msa.mdVersionList.Versions
}

func (ms *MetadataService) GetLatestDefinedVersion(space metadata_typedefs.MetadataSpace) (*core.AppVersion, error) {
	msa := ms.getMetadataServiceSpace(space)

	return msa.mdVersionList.GetLatestVersionDefined()
}

func (ms *MetadataService) IsVersionValid(versionString string, space metadata_typedefs.MetadataSpace) (bool, error) {
	msa := ms.getMetadataServiceSpace(space)

	version, err := core.GetAppVersionFromString(versionString)
	if err != nil {
		return false, errors.New("error parsing version string: " + err.Error())
	}
	valid := msa.mdVersionList.IsVersionValid(version)
	if !valid {
		return false, errors.New("no such version")
	}

	return true, nil
}

func (ms *MetadataService) GetMetadataManifestItemsInVersion(versionString string, space metadata_typedefs.MetadataSpace) ([]*metadata_typedefs.MetadataManifestItem, error) {
	msa := ms.getMetadataServiceSpace(space)

	version, err := core.GetAppVersionFromString(versionString)
	if err != nil {
		return nil, errors.New("error parsing version string: " + err.Error())
	}
	valid := msa.mdVersionList.IsVersionValid(version)
	if !valid {
		return nil, errors.New("no such version")
	}

	manifest, err := msa.getMetadataManifestForVersion(version)
	if err != nil {
		return nil, errors.New("error loading manifest: " + err.Error())
	}

	return manifest.MetadataManifestItems, nil
}

func (ms *MetadataService) GetMetadataManifestItemInVersion(metadataItemKey string, version *core.AppVersion, space metadata_typedefs.MetadataSpace) (*metadata_typedefs.MetadataManifestItem, error) {
	msa := ms.getMetadataServiceSpace(space)

	valid := msa.mdVersionList.IsVersionValid(version)
	if !valid {
		return nil, errors.New("no such version")
	}

	manifest, err := msa.getMetadataManifestForVersion(version)
	if err != nil {
		return nil, errors.New("error loading manifest: " + err.Error())
	}

	manifestItem := manifest.GetManifestItem(metadataItemKey)
	if manifestItem == nil {
		return nil, errors.New("error finding manifest item")
	}
	return manifestItem, nil
}

func (ms *MetadataService) IsMetadataHashUpToDate(key string, hash string, space metadata_typedefs.MetadataSpace, version *core.AppVersion) (bool, error) {
	msa := ms.getMetadataServiceSpace(space)

	result, err := msa.isMetadataHashUpToDate(key, hash, version)
	if err != nil {
		logger.LogError("failed to find if metadata up to data|space=" + space.String() +
			"|version=" + version.String() +
			"|key=" + key +
			"|hash=" + hash +
			"|error=" + err.Error())
		return false, errors.New("failed to check hash up to date for metadata")
	}

	return result, nil
}

func (ms *MetadataService) GetMetadataItem(itemPtr metadata_typedefs.IMetadataItem, version *core.AppVersion) error {
	if itemPtr == nil {
		logger.LogError("itemPtr is null")
		return errors.New("itemPtr is null")
	}

	msa := ms.getMetadataServiceSpace(itemPtr.GetMetadataSpace())

	var metadataJson string
	var err error
	metadataJson, err = msa.getMetadataJsonForItem(itemPtr, version)

	if err != nil {
		logger.LogError("Could not find metadata|metadata_space=" + itemPtr.GetMetadataSpace().String() +
						"|metadata_key=" + itemPtr.GetKey() +
						"|version=" + version.String() +
						"|error=" + err.Error())
		return errors.New("failed to find metadata")
	}

	if metadataJson == "" {
		logger.LogError("Could not find metadata|metadata_space=" + string(itemPtr.GetMetadataSpace()) +
			            "|metadata_key=" + itemPtr.GetKey() +
						"|version=" + version.String())
		return errors.New("failed to find metadata")
	}

	err = json.Unmarshal([]byte(metadataJson), itemPtr)
	if err != nil {
		logger.LogError("Error deserializing metadata json|metadata_space=" + string(itemPtr.GetMetadataSpace()) +
			            "|metadata_key=" + itemPtr.GetKey() +
						"|version=" + version.String() +
			            "|error=" + err.Error())
		return errors.New("failed deserializing metadata")
	}

	return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (ms *MetadataService) GetMetadataItemRawContent(metadataItemKey string, version *core.AppVersion, space metadata_typedefs.MetadataSpace) (string, error) {
	msa := ms.getMetadataServiceSpace(space)

	metadataItemJson, err := msa.mdFetcher.GetMetadataJsonByKey(metadataItemKey, version.String())
	if err != nil {
		return "", errors.New("error fetching metadata item raw content: " + err.Error())
	}

	return metadataItemJson, nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (ms *MetadataService) SetMetadataItem(itemPtr metadata_typedefs.IMetadataItem, version *core.AppVersion) error {
	if itemPtr == nil {
		logger.LogError("itemPtr is null")
		return errors.New("itemPtr is null")
	}

	msa := ms.getMetadataServiceSpace(itemPtr.GetMetadataSpace())
	err := msa.setMetadataJsonForItem(itemPtr, version)

	if err != nil {
		logger.LogError("error saving metadata item" +
						"|version=" + version.String() +
						"|metadata key=" + itemPtr.GetKey() +
						"|error=" + err.Error())
		return errors.New("error saving metadata item: " + err.Error())
	}

	return nil
}

func (ms *MetadataService) getMetadataServiceSpace(space metadata_typedefs.MetadataSpace) *MetadataServiceSpace {
	var msa *MetadataServiceSpace

	switch space {
	case metadata_typedefs.METADATA_SPACE_SHARED: msa = ms.sharedMDServiceSpace
	case metadata_typedefs.METADATA_SPACE_APP: msa = ms.appMDServiceSpace
	default:
		logger.LogError("undefined metadata space|space=" + space.String())
	}

	return msa
}

