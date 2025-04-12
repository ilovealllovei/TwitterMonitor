package utils

import (
	"fmt"
	"log"
	"runtime"
)

// LogError 自定义日志函数，支持和 log.Fatalf 一样的入参格式
func LogError(format string, v ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(1)
	err := fmt.Errorf(format, v...)
	if ok {
		log.Printf("[ERROR] %s:%d - %v", file, line, err)
	} else {
		log.Printf("[ERROR] Unable to get caller information - %v", err)
	}
}
