package config

import (
	"github.com/spacetimi/server/timi_shared/utils/go_vars_helper"
	"github.com/spacetimi/server/timi_shared/utils/logger"
	"os"
)

var _appName string
var _appEnvironmentString string

var _environmentConfiguration *EnvironmentConfiguration

// Package init
func init() {
	_appName = os.Getenv("app_name")
	if _appName == "" {
		panic("App Name not set")
	}

	_appEnvironmentString = os.Getenv("app_environment")
	if _appEnvironmentString == "" {
		panic("App Environment not set")
	}

	_environmentConfiguration = readEnvironmentConfiguration(GetAppConfigFilesPath(), _appEnvironmentString)
}

func GetEnvironmentConfiguration() *EnvironmentConfiguration {
	if _environmentConfiguration == nil {
		logger.LogFatal("Shared Configuration not initialized")
		return nil
	}
	return _environmentConfiguration
}

func GetAppName() string {
	return _appName
}

func GetAppConfigFilesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/server/" + GetAppName() + "/config"
}

func GetAppMetadataFilesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/server/" + GetAppName() + "/metadata"
}

func GetSharedMetadataFilesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/server/timi_shared/metadata"
}


