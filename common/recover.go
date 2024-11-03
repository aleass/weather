package common

import (
	"fmt"
	"runtime"
	"time"
)

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
		if duration == 0 {
			panic("exit:0")
		}
		time.Sleep(time.Minute * 5)
		funcName(duration)
	}
}
