package config

import (
	"github.com/spacetimi/timi_shared_server/utils/go_vars_helper"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"os"
)

var _appName string
var _appDirPath string
var _appEnvironmentString string

var _environmentConfiguration *EnvironmentConfiguration

// Package init
func init() {
	_appName = os.Getenv("app_name")
	if _appName == "" {
		panic("App Name not set")
	}

	_appDirPath = os.Getenv("app_dir_path")
	if _appDirPath == "" {
		panic("App Dir Path not set")
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

func GetAppDirPath() string {
	return _appDirPath
}

func GetAppConfigFilesPath() string {
	return GetAppDirPath() + "/config"
}

func GetAppMetadataFilesPath() string {
	return GetAppDirPath() + "/metadata"
}

func GetAppResourcesPath() string {
	return GetAppDirPath() + "/resources"
}

func GetAppTemplateFilesPath() string {
    return GetAppResourcesPath() + "/templates"
}

func GetAppImageFilesPath() string {
	return GetAppResourcesPath() + "/images"
}

func GetSharedMetadataFilesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/timi_shared_server/metadata"
}

func GetSharedResourcesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/timi_shared_server/resources"
}

func GetSharedTemplateFilesPath() string {
    return GetSharedResourcesPath() + "/templates"
}

func GetSharedImageFilesPath() string {
	return go_vars_helper.GOPATH + "/src/github.com/spacetimi/timi_shared_server/resources/images"
}

