package scripting_utilities

import (
	"os"
	"path/filepath"
)

func CheckIfExecutableIsUpToDate(executableBuildTime int64, sourceDirectoryPath string) bool {

	/**
	 * Get all files in source directory (recursive) and check if any of them has mtime > build time
	 */

	fileList := GetFilesInDirectoryRecursive(sourceDirectoryPath, false)
	for _, file := range fileList {
		fileInfo, err := os.Stat(file)
		if err != nil {
			panic(err)
		}
		if fileInfo.ModTime().Unix() > executableBuildTime {
			return false
		}
	}

	return true
}

func GetFilesInDirectoryRecursive(directoryPath string, includeSubDirectories bool) []string {
	directoryInfo, err := os.Stat(directoryPath)
	if err != nil {
		panic(err)
	}
	if !directoryInfo.IsDir() {
		panic("Not a directory: " + directoryPath)
	}

	fileList := []string{}
	err = filepath.Walk(directoryPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() || includeSubDirectories {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return fileList
}

type ScriptError struct {
	Reason string
}
func (e ScriptError) Error() string {
	return e.Reason
}
