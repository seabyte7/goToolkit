package fileLib

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 获得当前的执行环境的位置
func GetCurrPath() string {
	exePath, err := os.Executable()
	if err != nil {
		// log it
		return ""
	}

	return filepath.Dir(exePath)
}

// 文件是否存在
func IsFileExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return false
	}

	return true
}

// 目录是否存在
func IsDirectoryExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// 获得指定目录下的指定后缀的所有文件
// suffix: 后缀名字，不带.
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

// 获得指定文件名字列表
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
