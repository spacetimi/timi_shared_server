package metadata_fetchers

import (
    "encoding/json"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/file_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
)

type MetadataFetcherS3 struct { // Implements IMetadataFetcher
    url string
}

func NewMetadataFetcherS3(url string) metadata_typedefs.IMetadataFetcher {
    mf := MetadataFetcherS3{url:url}
    return &mf
}

/********** Begin IMetadataFetcher implementation **********/
func (mf *MetadataFetcherS3) GetMetadataJsonByKey(key string, version string) (string, error) {
    fileUrl := mf.url + "/" + version + "/" + key + ".json"
    fileContents, err := file_utils.ReadFileFromURL(fileUrl)
    if err != nil {
        logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_FILE +
                        "|url=" + fileUrl +
                        "|error=" + err.Error())
        return "", errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_FILE)
    }
    return fileContents, nil
}

func (mf *MetadataFetcherS3) GetMetadataVersionList() (*metadata_typedefs.MetadataVersionList, error) {
    mvl := metadata_typedefs.MetadataVersionList{}

    fileUrl := mf.url + "/" + "MetadataVersionList.json"
    fileContents, err := file_utils.ReadFileBytesFromURL(fileUrl)
    if err != nil {
        logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_VERSIONS_LIST +
                        "|url=" + fileUrl +
                        "|error=" + err.Error())
        return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_VERSIONS_LIST)
    }

    err = json.Unmarshal(fileContents, &mvl)
    if err != nil {
        logger.LogError(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_VERSIONS_LIST +
                        "|url=" + fileUrl +
                        "|error=" + err.Error())
        return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_VERSIONS_LIST)
    }

    mvl.Initialize()
    return &mvl, nil
}

func (mf *MetadataFetcherS3) GetMetadataManifestForVersion(version string) (*metadata_typedefs.MetadataManifest, error) {
    manifest := metadata_typedefs.MetadataManifest{}

    fileUrl := mf.url + "/" + version + "/" + "MetadataManifest.json"
    fileContents, err := file_utils.ReadFileBytesFromURL(fileUrl)
    if err != nil {
        logger.LogError(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_MANIFEST +
                        "|url=" + fileUrl +
                        "|error=" + err.Error())
        return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_READ_METADATA_MANIFEST)
    }

    err = json.Unmarshal(fileContents, &manifest)
    if err != nil {
        logger.LogError(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_MANIFEST +
                        "|url=" + fileUrl +
                        "|error=" + err.Error())
        return nil, errors.New(metadata_typedefs.ERROR_FAILED_TO_DESERIALIZE_METADATA_MANIFEST)
    }

    manifest.Initialize()
    return &manifest, nil
}
/********** End IMetadataFetcher implementation **********/

