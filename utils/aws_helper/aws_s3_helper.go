package aws_helper

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   reader,
		ACL:    &publicReadACL,
	})

	if err != nil {
		return errors.New("error uploading to s3: " + err.Error())
	}

	return nil
}
