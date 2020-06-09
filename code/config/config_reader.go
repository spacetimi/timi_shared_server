package config

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"os"
)

func ReadConfigFile(filePath string, configObject interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("error opening config file: " + err.Error())
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.LogError("error closing config file" +
							"|file path=" + filePath +
							"|error=" + err.Error())
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configObject)
	if err != nil {
		return errors.New("error decoding configuration file: " + err.Error())
	}
	return nil
}