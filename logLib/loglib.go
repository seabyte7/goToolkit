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

// Initialize a default log
// Most of the time it is used when the test code reports an error.
func newDefaultZapLog() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console", // 使用 console 编码器，它输出 plain text 格式
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"goToolKit_debug.log"}, // 输出到指定文件
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Zap() *zap.Logger {
	if loggerInst == nil {
		loggerInst = newDefaultZapLog()
	}

	return loggerInst
}

func Sugar() *zap.SugaredLogger {
	if loggerInst == nil {
		loggerInst = newDefaultZapLog()
	}

	return loggerInst.Sugar()
}
