package securityLib

import (
	"crypto/sha256"
	"fmt"
	"github.com/seabyte7/goToolkit/logLib"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

func SHA256SumBytes(data []byte, upper bool) string {
	sumResult := sha256.Sum256(data)
	if upper {
		return fmt.Sprintf("%X", sumResult)
	}

	return fmt.Sprintf("%x", sumResult)
}

func SHA256SumString(data string, upper bool) string {
	sumResult := sha256.Sum256([]byte(data))
	if upper {
		return fmt.Sprintf("%X", sumResult)
	}

	return fmt.Sprintf("%x", sumResult)
}

func SHA256SumFile(filePath string, upper bool) string {
	hashObj := sha256.New()
	fp, err := os.Open(filePath)
	if err != nil {
		logLib.Zap().Error("SHA256SumFile open file failed.",
			zap.String("filePath", filePath),
			zap.String("error", err.Error()))
		return ""
	}
	defer fp.Close()

	if _, err := io.Copy(hashObj, fp); err != nil {
		logLib.Zap().Error("SHA256SumFile io.Copy file to hashObj occur error.",
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

func SHA256CompareString(data string, sourceSHA256 string) bool {
	strSum := SHA256SumString(data, true)
	sourceSHA256 = strings.ToUpper(sourceSHA256)
	if sourceSHA256 == strSum {
		return true
	}

	return false
}

func SHA256CompareBytes(data []byte, sourceSHA256 string) bool {
	strSum := SHA256SumBytes(data, true)
	sourceSHA256 = strings.ToUpper(sourceSHA256)
	if sourceSHA256 == strSum {
		return true
	}

	return false
}

func SHA256CompareFile(filePath string, sourceSHA256 string) bool {
	strSum := SHA256SumFile(filePath, true)
	sourceSHA256 = strings.ToUpper(sourceSHA256)
	if strSum == sourceSHA256 {
		return true
	}

	return false
}

func SHA256CompareFile2(srcFilePath string, destFilePath string) bool {
	srcSum := SHA256SumFile(srcFilePath, true)
	destSum := SHA256SumFile(destFilePath, true)
	if srcSum == destSum {
		return true
	}

	return false
}
