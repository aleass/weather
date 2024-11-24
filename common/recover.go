package common

import (
	"fmt"
	"runtime"
	"time"
)

// 记录错误信息和调用
func RecoverWithStackTrace(funcName func(duration time.Duration), duration time.Duration) {
	if r := recover(); r != nil {
		Logger.Error(fmt.Sprintf("退出了,发现错误recover : %v", r))
		_, file, line, ok := runtime.Caller(1)
		if ok {
			Logger.Error(fmt.Sprintf("Panic occurred at %s:%d", file, line))
		}

		// Print stack trace
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]
		Logger.Error(fmt.Sprintf("Stack trace:\n%s", stack))

		//至少睡眠30s
		if duration < 30*time.Second {
			duration = 30 * time.Second
		}

		time.Sleep(duration)

		if funcName == nil {
			return
		}
		funcName(duration)
	}
}
