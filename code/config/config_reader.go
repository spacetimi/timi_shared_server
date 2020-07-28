package config

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/v2/utils/file_utils"
)

type IConfig interface {
	OnConfigLoaded()
}

func ReadConfigFile(filePath string, configObject IConfig) error {
	err := file_utils.ReadJsonFileIntoJsonObject(filePath, configObject)
	if err != nil {
		return errors.New("error reading config file: " + err.Error())
	}

	configObject.OnConfigLoaded()

	return nil
}
