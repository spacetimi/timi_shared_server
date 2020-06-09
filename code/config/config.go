package config

import (
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"os"
)

var _appName string
var _appDirPath string
var _sharedDirPath string
var _appEnvironmentString string

var _environmentConfiguration *EnvironmentConfiguration

// Package init
func init() {
	_appName = os.Getenv("app_name")
	if _appName == "" {
		logger.LogFatal("app name not set")
	}

	_appDirPath = os.Getenv("app_dir_path")
	if _appDirPath == "" {
		logger.LogFatal("app dir path not set")
	}

	_sharedDirPath = os.Getenv("shared_dir_path")
	if _sharedDirPath == "" {
		logger.LogFatal("shared dir path not set")
	}

	_appEnvironmentString = os.Getenv("app_environment")
	if _appEnvironmentString == "" {
		logger.LogFatal("app environment not set")
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

