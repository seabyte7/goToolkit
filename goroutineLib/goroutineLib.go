package goroutineLib

import (
	"fmt"
	"go.uber.org/zap"
)

func RunGoroutine(goroutineName string, goroutineFunc func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logLib.Zap().Error("goroutine:%v panic, err:%v",
					zap.String("name", goroutineName),
					zap.String("error", err.(error).Error()))

				msg := fmt.Sprintf("goroutine:%v panic, err:%v", goroutineName, err.(error).Error())
				pushLib.PushTextMessageToDefault(msg)
			}
		}()

		goroutineFunc()
	}()
}
