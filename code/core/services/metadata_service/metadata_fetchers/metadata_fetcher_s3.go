package metadata_fetchers

import (
	"encoding/json"
	"errors"

	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/aws_helper"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type MetadataFetcherS3 struct { // Implements IMetadataFetcher
	url               string
	adminS3BucketName string // Used for admin tool functions
}

func NewMetadataFetcherS3(url string, adminS3BucketName string) metadata_typedefs.IMetadataFetcher {
	mf := MetadataFetcherS3{
		url:               url,
		adminS3BucketName: adminS3BucketName,
	}
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

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherS3) SetMetadataJsonByKey(key string, metadataJson string, version string) error {

	awsSession, err := aws_helper.GetNewDefaultSession()
	if err != nil {
		return errors.New("error creating aws session: " + err.Error())
	}

	err = aws_helper.UploadToS3(awsSession, []byte(metadataJson), mf.adminS3BucketName, "metadata/"+version+"/"+key+".json")
	if err != nil {
		return errors.New("error uploading metadata item: " + err.Error())
	}

	return nil
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

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherS3) SetMetadataManifestForVersion(manifest *metadata_typedefs.MetadataManifest, version string) error {
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

	awsSession, err := aws_helper.GetNewDefaultSession()
	if err != nil {
		return errors.New("error creating aws session: " + err.Error())
	}

	err = aws_helper.UploadToS3(awsSession, manifestJson, mf.adminS3BucketName, "metadata/"+version+"/MetadataManifest.json")
	if err != nil {
		return errors.New("error uploading metadata manifest: " + err.Error())
	}

	return nil
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mf *MetadataFetcherS3) SetMetadataVersionList(mvl *metadata_typedefs.MetadataVersionList) error {

	awsSession, err := aws_helper.GetNewDefaultSession()
	if err != nil {
		return errors.New("error creating aws session: " + err.Error())
	}

	var bytes []byte

	if config.GetEnvironmentConfiguration().AppEnvironment == config.PRODUCTION {
		bytes, err = json.Marshal(mvl)
	} else {
		bytes, err = json.MarshalIndent(mvl, "", "    ")
	}
	if err != nil {
		return errors.New("error serializing metadata version list: " + err.Error())
	}

	err = aws_helper.UploadToS3(awsSession, bytes, mf.adminS3BucketName, "metadata/MetadataVersionList.json")
	if err != nil {
		return errors.New("error uploading metadata version list: " + err.Error())
	}

	return nil
}

/********** End IMetadataFetcher implementation **********/
