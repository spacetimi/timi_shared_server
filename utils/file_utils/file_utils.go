package file_utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func DoesFileOrDirectoryExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetFileNameFromPath(filePath string) string {
	return filepath.Base(filePath)
}

func GetFileNamesFromPaths(filePaths []string) []string {
	var fileNames []string
	for _, filePath := range filePaths {
		fileNames = append(fileNames, GetFileNameFromPath(filePath))
	}

	return fileNames
}

func GetFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func GetFileModTimeUnix(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return -1, err
	}

	return fileInfo.ModTime().Unix(), nil
}

func ReadFileBytesFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		logger.LogError("failed fetching file from url|url=" + url +
			"|error=" + err.Error())
		return nil, errors.New("failed reading file")
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.LogError("failed to close http request|url=" + url +
				"|error=" + err.Error())
		}
	}()

	buffer := new(bytes.Buffer)
	n, err := buffer.ReadFrom(response.Body)
	if err != nil {
		logger.LogError("failed reading file contents from http response|url=" + url +
			"|error=" + err.Error())
		return nil, errors.New("failed reading file")
	}

	if n != response.ContentLength {
		logger.LogError("read file content length not equal to expected|url=" + url +
			"|expected_length=" + strconv.FormatInt(response.ContentLength, 10) +
			"|actual_length=" + strconv.FormatInt(n, 10))
		return nil, errors.New("failed reading file")
	}

	content := buffer.Bytes()
	return content, nil
}

func ReadFileFromURL(url string) (string, error) {
	fileBytes, err := ReadFileBytesFromURL(url)
	if err != nil {
		return "", err
	}

	return string(fileBytes), nil
}

func CopyFile(sourcePath string, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.New("error opening source file: " + err.Error())
	}
	defer func() {
		err := sourceFile.Close()
		if err != nil {
			logger.LogWarning("error closing source file" + "|file path=" + sourcePath + "|error=" + err.Error())
		}
	}()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return errors.New("error creating destination file: " + err.Error())
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return errors.New("error copying source to destination: " + err.Error())
	}

	err = destinationFile.Close()
	if err != nil {
		return errors.New("error closing destination file after copy: " + err.Error())
	}

	return nil
}

func ReadJsonFileIntoJsonObject(filePath string, jsonObject interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("error opening json file: " + err.Error())
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.LogError("error closing json file" +
				"|file path=" + filePath +
				"|error=" + err.Error())
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonObject)
	if err != nil {
		return errors.New("error decoding json file: " + err.Error())
	}
	return nil

}

func GetFileNamesInDirectory(path string) ([]string, error) {
	if !DoesFileOrDirectoryExist(path) {
		return nil, errors.New("no such directory")
	}

	directory, err := os.Open(path)
	if err != nil {
		return nil, errors.New("error opening directory: " + err.Error())
	}

	fileNames, err := directory.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}

func GetFilePathsInDirectoryMatchingPattern(path string, pattern string) ([]string, error) {
	if !DoesFileOrDirectoryExist(path) {
		return nil, errors.New("no such directory")
	}

	fileNames, err := filepath.Glob(path + "/" + pattern)
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}
