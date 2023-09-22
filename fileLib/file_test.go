package fileLib

import (
	"testing"
)

func TestIsFileExist(t *testing.T) {
	filePath := "./file.go"
	if !IsFileExist(filePath) {
		t.Errorf("TestIsFileExist file:%v isn't exist.", filePath)
	}
}

func TestIsDirectoryExist(t *testing.T) {
	path := "../fileLib"
	if !IsDirectoryExist(path) {
		t.Errorf("TestIsDirectoryExist file:%v isn't exist.", path)
	}
}

func TestGetAbsoluteFilelist(t *testing.T) {
	path := "../fileLib"
	filePathList, err := GetAbsoluteFileList(path, "go")
	if err != nil {
		t.Errorf("TestAbsoluteFilelist path:%v err:%v", path, err.Error())
	}

	if len(filePathList) != 2 {
		t.Errorf("TestAbsoluteFilelist path:%v file count:%v error", path, len(filePathList))
	}
}

func TestGetFileList(t *testing.T) {
	path := "../fileLib"
	fileList, err := GetFileList(path, "", "")
	if err != nil {
		t.Errorf("TestGetFileList path:%v err:%v", path, err.Error())
	}

	if len(fileList) != 2 {
		t.Errorf("TestGetFileList path:%v file count:%v error", path, len(fileList))
	}
}
