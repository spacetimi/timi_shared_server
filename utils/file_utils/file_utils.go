package file_utils

import (
	"bytes"
	"errors"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"net/http"
	"os"
	"strconv"
)

func DoesFileOrDirectoryExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
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
