package aws_helper

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3(bytes []byte, bucketName string, key string) error {

	session, err := GetNewDefaultSession()
	if err != nil {
		return errors.New("error getting aws-session")
	}
	svc := s3.New(session)
	uploader := s3manager.NewUploaderWithClient(svc)

	reader := strings.NewReader(string(bytes))

	// TODO: Avi: Do not use public-read (not even on test and staging)
	publicReadACL := "public-read"

	_, err = uploader.Upload(&s3manager.UploadInput{
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
