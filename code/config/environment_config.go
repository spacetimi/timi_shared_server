package config

import (
	"encoding/json"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"os"
	"strings"
)

type AppEnvironment_t int
const (
	LOCAL = iota
	TEST
	STAGING
	PRODUCTION
)

type EnvironmentConfiguration struct {
	AppEnvironment AppEnvironment_t
	Port int

	// MongoDB config
	SharedMongoURL string
	SharedDatabaseName string
	AppMongoURL string
	AppDatabaseName string

	// Redis config
	SharedRedisURL string
	SharedRedisPasswd string
	AppRedisURL string
	AppRedisPasswd string

	// MetaData config
	SharedMetadataSourceURL string
	AppMetadataSourceURL string
	MetadataAutoUpdaterPollSeconds int

	// Admin tool config
	AdminToolConfig AdminToolConfiguration
}

type AdminToolConfiguration struct {
	SharedMetadataS3BucketName string
	AppMetadataS3BucketName string
}

func readEnvironmentConfiguration(pathToConfigFiles string, appEnvString string) *EnvironmentConfiguration {

	switch appEnvString {
	case "Local": 		// OK
	case "Test": 		// OK
	case "Staging": 	// OK
	case "Production": 	// OK
	default:
		panic("Invalid app environment: " + appEnvString)
	}

	environmentConfigFilePath := pathToConfigFiles + "/environment_config." + strings.ToLower(appEnvString) + ".json"
	environmentConfigFile, err := os.Open(environmentConfigFilePath)
	if err != nil {
		logger.LogFatal("Cannot open configuration file at: " + environmentConfigFilePath)
		return nil
	}
	defer environmentConfigFile.Close()

	var environmentConfiguration *EnvironmentConfiguration
	decoder := json.NewDecoder(environmentConfigFile)
	err = decoder.Decode(&environmentConfiguration)
	if err != nil {
		logger.LogFatal("Error decoding configuration file at: " + environmentConfigFilePath + ". Error: " + err.Error())
		return nil
	}

	return environmentConfiguration
}
