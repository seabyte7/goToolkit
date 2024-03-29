package securityLib

import (
	"crypto/md5"
	"fmt"
	"github.com/seabyte7/goToolkit/logLib"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

func MD5SumBytes(data []byte, upper bool) string {
	sumResult := md5.Sum(data)
	if upper {
		return fmt.Sprintf("%X", sumResult)
	}

	return fmt.Sprintf("%x", sumResult)
}

func MD5SumString(data string, upper bool) string {
	sumResult := md5.Sum([]byte(data))
	if upper {
		return fmt.Sprintf("%X", sumResult)
	}

	return fmt.Sprintf("%x", sumResult)
}

func MD5SumFile(filePath string, upper bool) string {
	hashObj := md5.New()
	fp, err := os.Open(filePath)
	if err != nil {
		logLib.Zap().Error("MD5SumFile open file failed.",
			zap.String("filePath", filePath),
			zap.String("error", err.Error()))
		return ""
	}
	defer fp.Close()

	if _, err := io.Copy(hashObj, fp); err != nil {
		logLib.Zap().Error("MD5SumFile io.Copy file to hashObj occur error.",
			zap.String("filePath", filePath),
			zap.String("error", err.Error()))
		return ""
	}

	sumResult := hashObj.Sum(nil)
	if upper {
		return fmt.Sprintf("%X", sumResult)
	}

	return fmt.Sprintf("%x", sumResult)
}

func MD5CompareString(data string, sourceMd5 string) bool {
	strSum := MD5SumString(data, true)
	sourceMd5 = strings.ToUpper(sourceMd5)
	if strSum == sourceMd5 {
		return true
	}

	return false
}

func MD5CompareBytes(data []byte, sourceMd5 string) bool {
	strSum := MD5SumBytes(data, true)
	sourceMd5 = strings.ToUpper(sourceMd5)
	if strSum == sourceMd5 {
		return true
	}

	return false
}

func MD5CompareFile(filePath string, sourceMd5 string) bool {
	strSum := MD5SumFile(filePath, true)
	sourceMd5 = strings.ToUpper(sourceMd5)
	if strSum == sourceMd5 {
		return true
	}

	return false
}

func MD5CompareFile2(srcFilePath string, destFilePath string) bool {
	srcSum := MD5SumFile(srcFilePath, true)
	destSum := MD5SumFile(destFilePath, true)
	if srcSum == destSum {
		return true
	}

	return false
}
