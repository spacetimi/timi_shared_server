package aws_helper

import (
    "errors"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "strings"
)

func GetNewDefaultSession() (*session.Session, error) {

    awsSession := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    _, err := awsSession.Config.Credentials.Get()
    if err != nil {
        return nil, errors.New("error getting credentials for aws session: " + err.Error())
    }

    return awsSession, nil
}

func UploadToS3(awsSession *session.Session, bytes []byte, bucketName string, key string) error {
    reader := strings.NewReader(string(bytes))

    // TODO: Avi: Do not use public-read (not even on test and staging)
    publicReadACL := "public-read"

    uploader := s3manager.NewUploader(awsSession)
    _, err := uploader.Upload(&s3manager.UploadInput{
        Bucket:aws.String(bucketName),
        Key:aws.String(key),
        Body: reader,
        ACL: &publicReadACL,
    })

    if err != nil {
        return errors.New("error uploading to s3: " + err.Error())
    }

    return nil
}
