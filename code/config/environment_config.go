package config

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/spacetimi/timi_shared_server/v2/utils/logger"
)

type AppEnvironment_t int

const (
	LOCAL = iota
	TEST
	STAGING
	PRODUCTION
)

func (appEnvironment AppEnvironment_t) String() string {
	switch appEnvironment {
	case LOCAL:
		return "Local"
	case TEST:
		return "Test"
	case STAGING:
		return "Staging"
	case PRODUCTION:
		return "Production"
	}
	return ""
}

type EnvironmentConfiguration struct {
	AppEnvironment   AppEnvironment_t
	Port             int
	ApiServerBaseURL string

	// MetaData config
	SharedMetadataSourceURL        string
	AppMetadataSourceURL           string
	MetadataAutoUpdaterPollSeconds int

	// Admin tool config
	AdminToolConfig AdminToolConfiguration
}

type AdminToolConfiguration struct {
	SharedMetadataS3BucketName string
	AppMetadataS3BucketName    string
}

func readEnvironmentConfiguration(pathToConfigFiles string, appEnvString string) *EnvironmentConfiguration {

	switch appEnvString {
	case "Local": // OK
	case "Test": // OK
	case "Staging": // OK
	case "Production": // OK
	default:
		panic("Invalid app environment: " + appEnvString)
	}

	environmentConfigFilePath := pathToConfigFiles + "/environment_config." + strings.ToLower(appEnvString) + ".json"
	environmentConfigFile, err := os.Open(environmentConfigFilePath)
	if err != nil {
		logger.LogFatal("cannot open configuration file|file path=" + environmentConfigFilePath)
		return nil
	}
	defer func() {
		err := environmentConfigFile.Close()
		if err != nil {
			logger.LogError("error closing config file" +
				"|file path=" + environmentConfigFilePath +
				"|error=" + err.Error())
		}
	}()

	var environmentConfiguration *EnvironmentConfiguration
	decoder := json.NewDecoder(environmentConfigFile)
	err = decoder.Decode(&environmentConfiguration)
	if err != nil {
		logger.LogFatal("error decoding configuration file" +
			"|file path=" + environmentConfigFilePath +
			"|error=" + err.Error())
		return nil
	}

	return environmentConfiguration
}
