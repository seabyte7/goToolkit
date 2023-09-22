package fileLib

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GetWorkPath GetCurrPath Gets the working path of the current execution file
func GetWorkPath() string {
	exePath, err := os.Executable()
	if err != nil {
		// log it
		return ""
	}

	return filepath.Dir(exePath)
}

// IsFileExist Check whether the file exists
func IsFileExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return false
	}

	return true
}

// IsDirectoryExist Check whether the directory exists
func IsDirectoryExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// GetAbsoluteFileList Gets all files with the specified suffix in the specified directory
// suffix: suffix name, without.
func GetAbsoluteFileList(rootPath, suffix string) ([]string, error) {
	if !IsDirectoryExist(rootPath) {
		return nil, os.ErrNotExist
	}

	filePathList := make([]string, 0, 32)
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if len(suffix) > 0 && !strings.HasSuffix(path, suffix) {
			return nil
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		filePathList = append(filePathList, absolutePath)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePathList, nil
}

// GetFileList Gets a list of file names for the specified path
func GetFileList(rootPath, prefix, suffix string) ([]string, error) {
	if !IsDirectoryExist(rootPath) {
		return nil, os.ErrNotExist
	}

	fileList := make([]string, 0, 32)
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if len(prefix) > 0 && !strings.HasPrefix(path, prefix) {
			return nil
		}

		if len(suffix) > 0 && !strings.HasSuffix(path, suffix) {
			return nil
		}

		filename := filepath.Base(path)

		fileList = append(fileList, filename)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
