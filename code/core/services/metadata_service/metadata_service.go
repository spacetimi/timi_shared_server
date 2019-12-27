package metadata_service

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/server/timi_shared/code/core"
	"github.com/spacetimi/server/timi_shared/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/server/timi_shared/utils/logger"
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

func (ms *MetadataService) IsMetadataHashUpToDate(key string, hash string, space metadata_typedefs.MetadataSpace, version *core.AppVersion) (bool, error) {
	var msa *MetadataServiceSpace

	switch space {
	case metadata_typedefs.METADATA_SPACE_SHARED: msa = ms.sharedMDServiceSpace
	case metadata_typedefs.METADATA_SPACE_APP: msa = ms.appMDServiceSpace
	}

	if msa == nil {
		logger.LogError("unknown metadata space|metadata_space=" + space.String())
		return false, errors.New("unknown metadata space")
	}

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

	var msa *MetadataServiceSpace

	switch itemPtr.GetMetadataSpace() {
	case metadata_typedefs.METADATA_SPACE_SHARED: msa = ms.sharedMDServiceSpace
	case metadata_typedefs.METADATA_SPACE_APP: msa = ms.appMDServiceSpace
	}

	if msa == nil {
		logger.LogError("unknown metadata space|metadata_space=" + itemPtr.GetMetadataSpace().String())
		return errors.New("unknown metadata space")
	}

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

