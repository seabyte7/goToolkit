package logLib

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// The loglib package is not dependent on any other package within the goToolKit,
// and goToolKit should initialize logLib before it is used.

var (
	loggerInst *zap.Logger
)

func StartLog(filePath string, logLevel zapcore.Level) {
	// create standard file encoder
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("open log file:%v error:%v", filePath, err.Error()))
	}
	consoleOutput := zapcore.AddSync(consoleFile)

	// create json encoder
	jsonEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	jsonFilePath := filePath + "_json"
	jsonFile, err := os.OpenFile(jsonFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("open json log file:%v error:%v", jsonFilePath, err.Error()))
	}
	jsonOutput := zapcore.AddSync(jsonFile)

	// create core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleOutput, logLevel),
		zapcore.NewCore(jsonEncoder, jsonOutput, logLevel),
	)

	// create zap log
	logger := zap.New(core)
	if err != nil {
		panic(fmt.Errorf("new zaplog  failed:%v", err.Error()))
	}

	loggerInst = logger
}

func StopLog() {
	if loggerInst != nil {
		loggerInst.Sync()
	}
}

func Zap() *zap.Logger {
	return loggerInst
}

func Sugar() *zap.SugaredLogger {
	return loggerInst.Sugar()
}
