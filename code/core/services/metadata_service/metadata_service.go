package metadata_service

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/code/core"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type MetadataService struct {
	sharedMDServiceSpace *MetadataServiceSpace
	appMDServiceSpace    *MetadataServiceSpace
}

/* Package init */
func Initialize() {
	instance = &MetadataService{}
	instance.sharedMDServiceSpace = newMetadataServiceSpace(metadata_typedefs.METADATA_SPACE_SHARED)
	instance.appMDServiceSpace    = newMetadataServiceSpace(metadata_typedefs.METADATA_SPACE_APP)
}

var instance *MetadataService
func Instance() *MetadataService {
	if instance == nil {
		logger.LogError("Metadata Service instance is null")
	}
	return instance
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
	    	continue
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

