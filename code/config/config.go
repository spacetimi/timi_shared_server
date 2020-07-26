package config

import (
	"os"

	"github.com/spacetimi/timi_shared_server/utils/go_vars_helper"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

var _appName string
var _appDirPath string
var _sharedDirPath string
var _appEnvironmentString string

var _environmentConfiguration *EnvironmentConfiguration

func Initialize(appName string) {
	_appName = appName

	_appDirPath = os.Getenv("app_dir_path")
	if _appDirPath == "" {
		// Fallback to default
		_appDirPath = go_vars_helper.GOPATH + "/src/github.com/spacetimi/" + _appName
	}

	_sharedDirPath = os.Getenv("shared_dir_path")
	if _sharedDirPath == "" {
		// Fallback to default
		_sharedDirPath = go_vars_helper.GOPATH + "/src/github.com/spacetimi/timi_shared_server"
	}

	_appEnvironmentString = os.Getenv("app_environment")
	if _appEnvironmentString == "" {
		// Fallback to default
		_appEnvironmentString = "Local"
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
	return GetAppDirPath() + "/tmp/metadata"
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

func GetSharedDirPath() string {
	return _sharedDirPath
}

func GetSharedMetadataFilesPath() string {
	return GetSharedDirPath() + "/tmp/metadata"
}

func GetSharedResourcesPath() string {
	return GetSharedDirPath() + "/resources"
}

func GetSharedTemplateFilesPath() string {
	return GetSharedResourcesPath() + "/templates"
}

func GetSharedImageFilesPath() string {
	return GetSharedResourcesPath() + "/images"
}
