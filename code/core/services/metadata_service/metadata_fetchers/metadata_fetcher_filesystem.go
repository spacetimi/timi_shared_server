package metadata_fetchers

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"io/ioutil"
	"os"
	"path"
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

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherFilesystem) SetMetadataJsonByKey(key string, metadataJson string, version string) error {
	filePath := mf.path + "/" + version  + "/" + key + ".json"
	err := ioutil.WriteFile(filePath, []byte(metadataJson), 0644)
	if err != nil {
		return errors.New("error writing file|error=" + err.Error())
	}
    return nil
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
func (mf *MetadataFetcherFilesystem) SetMetadataManifestForVersion(manifest *metadata_typedefs.MetadataManifest, version string) error {
	var err error
	var manifestJson []byte
	if config.GetEnvironmentConfiguration().AppEnvironment == config.PRODUCTION {
		manifestJson, err = json.Marshal(manifest)
    } else {
		manifestJson, err = json.MarshalIndent(manifest, "", "    ")
	}
	if err != nil {
		return errors.New("error serializing new manifest|error=" + err.Error())
	}

	filePath := mf.path + "/" + version + "/" + "MetadataManifest.json"
	_, err = os.Stat(filePath)
	if err != nil {
	    err = os.MkdirAll(path.Dir(filePath), 0755)
	    if err != nil {
	    	return errors.New("error creating path: " + path.Dir(filePath))
		}
	}

	err =  ioutil.WriteFile(filePath, []byte(manifestJson), 0644)
	if err != nil {
		return errors.New("error saving new manifest file|error=" + err.Error())
	}
    return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherFilesystem) SetMetadataVersionList(mvl *metadata_typedefs.MetadataVersionList) error {
    var err error
    var bytes []byte
	if config.GetEnvironmentConfiguration().AppEnvironment == config.PRODUCTION {
    	bytes, err = json.Marshal(mvl)
	} else {
		bytes, err = json.MarshalIndent(mvl, "", "    ")
	}
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

