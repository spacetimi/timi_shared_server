package config

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/utils/file_utils"
)

func ReadConfigFile(filePath string, configObject interface{}) error {
	err := file_utils.ReadJsonFileIntoJsonObject(filePath, configObject)
	if err != nil {
		return errors.New("error reading config file: " + err.Error())
	}

	return nil
}
