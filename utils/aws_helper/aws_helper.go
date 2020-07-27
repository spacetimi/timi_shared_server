package aws_helper

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func GetNewDefaultSession() (*session.Session, error) {

	optionalProfile := getOptionalAwsProfile()

	var awsSession *session.Session
	if len(optionalProfile) == 0 {
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           optionalProfile,
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	_, err := awsSession.Config.Credentials.Get()
	if err != nil {
		return nil, errors.New("error getting credentials for aws session: " + err.Error())
	}

	return awsSession, nil
}

func getOptionalAwsProfile() string {
	return os.Getenv("aws_profile")
}
