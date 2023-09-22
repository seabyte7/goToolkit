package goToolkit

import (
	"go.uber.org/zap/zapcore"
	"goToolkit/logLib"
)

// Start goToolKit库的初始化
// 主要做一些依赖的初始化，目前主要是日志
func Start(filePath string, logLevel zapcore.Level) {
	logLib.StartLog(filePath, logLevel)
}

// Stop goToolKit库的关闭
func Stop() {
	logLib.StopLog()
}
