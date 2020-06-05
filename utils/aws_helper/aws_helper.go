package aws_helper

import (
    "errors"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "strings"
)

func GetNewDefaultSession(optionalProfile string) (*session.Session, error) {

    var awsSession *session.Session
    if len(optionalProfile) == 0 {
        awsSession = session.Must(session.NewSessionWithOptions(session.Options{
            SharedConfigState: session.SharedConfigEnable,
        }))
    } else {
        awsSession = session.Must(session.NewSessionWithOptions(session.Options{
            Profile:optionalProfile,
            SharedConfigState: session.SharedConfigEnable,
        }))
    }

    _, err := awsSession.Config.Credentials.Get()
    if err != nil {
        return nil, errors.New("error getting credentials for aws session: " + err.Error())
    }

    return awsSession, nil
}

func UploadToS3(awsSession *session.Session, bytes []byte, bucketName string, key string) error {
    // krisa
    sess := session.Must(session.NewSessionWithOptions(session.Options{
       SharedConfigState: session.SharedConfigEnable,
    }))
    creds := stscreds.NewCredentials(sess, "arn:aws:iam::780204259180:role/Engineering-CrossAccountAccess")
    svc := s3.New(sess, &aws.Config{Credentials: creds})
    uploader := s3manager.NewUploaderWithClient(svc)

    reader := strings.NewReader(string(bytes))

    // TODO: Avi: Do not use public-read (not even on test and staging)
    publicReadACL := "public-read"

    // krisa
    // uploader := s3manager.NewUploader(awsSession)
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
