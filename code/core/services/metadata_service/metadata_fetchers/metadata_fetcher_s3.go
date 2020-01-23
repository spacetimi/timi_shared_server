package metadata_fetchers

import (
    "encoding/json"
    "errors"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/file_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "strings"
)

type MetadataFetcherS3 struct { // Implements IMetadataFetcher
    url string
    adminS3BucketName string   // Used for admin tool functions
}

func NewMetadataFetcherS3(url string, adminS3BucketName string) metadata_typedefs.IMetadataFetcher {
    mf := MetadataFetcherS3{
        url:url,
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
func (mf *MetadataFetcherS3) SetMetadataVersionList(mvl *metadata_typedefs.MetadataVersionList) error {

    awsSession := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    _, err := awsSession.Config.Credentials.Get()
    if err != nil {
        return errors.New("error getting credentials for aws session: " + err.Error())
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
    reader := strings.NewReader(string(bytes))

    // TODO: Avi: Do not use public-read (not even on test and staging)
    publicReadACL := "public-read"

    uploader := s3manager.NewUploader(awsSession)
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket:aws.String(mf.adminS3BucketName),
        Key:aws.String("metadata/MetadataVersionList.json"),
        Body: reader,
        ACL: &publicReadACL,
    })

    if err != nil {
        return errors.New("error uploading metadata version list: " + err.Error())
    }

    return nil
}
/********** End IMetadataFetcher implementation **********/

