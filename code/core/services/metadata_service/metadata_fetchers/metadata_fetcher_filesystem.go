package metadata_fetchers

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"io/ioutil"
)

type MetadataFetcherFilesystem struct { // Implements IMetadataFetcher
	path string
}

func NewMetadataFetcherFilesystem(path string) metadata_typedefs.IMetadataFetcher {
	mf := MetadataFetcherFilesystem{path:path}
	return &mf
}

/********** Begin IMetadataFetcher implementation **********/
func (mf *MetadataFetcherFilesystem) GetMetadataJsonByKey(key string, version string) (string, error) {
	filePath := mf.path + "/" + version  + "/" + key + ".json"
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_FILE +
			            "|path=" + filePath +
			            "|error=" + err.Error())
		return "", errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_FILE)
	}
	return string(bytes), nil
}

func (mf *MetadataFetcherFilesystem) GetMetadataVersionList() (*metadata_typedefs.MetadataVersionList, error) {
	mvl := metadata_typedefs.MetadataVersionList{}

	filePath := mf.path + "/" + "MetadataVersionList.json"
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_VERSIONS_LIST +
						"|file_path=" + filePath +
			            "|error=" + err.Error())
		return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_VERSIONS_LIST)
	}

	err = json.Unmarshal(bytes, &mvl)
	if err != nil {
		logger.LogError(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_VERSIONS_LIST +
						"|file_path=" + filePath +
			            "|error=" + err.Error())
		return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_VERSIONS_LIST)
	}

	mvl.Initialize()
	return &mvl, nil
}

func (mf *MetadataFetcherFilesystem) GetMetadataManifestForVersion(version string) (*metadata_typedefs.MetadataManifest, error) {
	manifest := metadata_typedefs.MetadataManifest{}

	filePath := mf.path + "/" + version + "/" + "MetadataManifest.json"
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_MANIFEST +
						"|file_path=" + filePath +
			            "|error=" + err.Error())
		return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_MANIFEST)
	}

	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		logger.LogError(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_MANIFEST +
						"|file_path=" + filePath +
			            "|error=" + err.Error())
		return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_MANIFEST)
	}

	manifest.Initialize()
	return &manifest, nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherFilesystem) SetMetadataVersionList(mvl *metadata_typedefs.MetadataVersionList) error {
    bytes, err := json.Marshal(mvl)
    if err != nil {
    	return errors.New("error serializing metadata version list to json|error=" + err.Error())
	}

	filePath := mf.path + "/" + "MetadataVersionList.json"
	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return errors.New("error writing metadata version list file|error=" + err.Error())
	}

    return nil
}
/********** End IMetadataFetcher implementation **********/

