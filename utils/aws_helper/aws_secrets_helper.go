package aws_helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

/**
Reads secrets from aws-secrets-manager that are simple json key-value-pairs
*/
func ReadJsonSecret(secretName string, secretSubkey string) (string, error) {
	secretString, err := ReadSecret(secretName)
	if err != nil {
		return "", errors.New("error reading secret string: " + err.Error())
	}

	var secretAsMap map[string]string
	err = json.Unmarshal([]byte(secretString), &secretAsMap)
	if err != nil {
		return "", errors.New("error deserializing simple secret: " + err.Error())
	}

	secretValue, ok := secretAsMap[secretSubkey]
	if !ok {
		return "", errors.New("no secret found for subkey")
	}

	return secretValue, nil
}

func ReadSecret(secretName string) (string, error) {
	session, err := GetNewDefaultSession()
	if err != nil {
		return "", errors.New("error getting aws-session")
	}

	//Create a Secrets Manager client
	svc := secretsmanager.New(session)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", errors.New("error getting secret value from aws: " + err.Error())
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these methods will work:
	if result.SecretString != nil {
		secretString := *result.SecretString
		return secretString, nil

	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return "", errors.New("error decoding read secret: " + err.Error())
		}
		decodedBinarySecret := string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret, nil
	}
}
