package config

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/v2/utils/aws_helper"
)

type PasswordConfig struct {
	Source string
	K1     string
	K2     string
}

func (pc *PasswordConfig) GetPassword() (string, error) {

	switch pc.Source {

	case "direct":
		return pc.K1, nil

	case "AWS_SECRETS":
		if len(pc.K1) == 0 {
			return "", errors.New("secret name not set")
		}
		if len(pc.K2) == 0 {
			return "", errors.New("secret subkey not set")
		}
		password, err := aws_helper.ReadJsonSecret(pc.K1, pc.K2)
		if err != nil || len(password) == 0 {
			return "", errors.New("error retrieving password from aws: " + err.Error())
		}
		return password, nil

	default:
		return "", errors.New("invalid source: " + pc.Source)
	}
}
