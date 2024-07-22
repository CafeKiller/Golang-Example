package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

// getLogFilePath 获取日志保存路径
func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

// getLogFileFullPath 获取完整的日志文件路径
func getLogFileFullPath() string {
	// 获取日志保存路径
	prefixPath := getLogFilePath()
	// 获取日志文件名
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// openLogFile 打开日志文件
func openLogFile(filePath string) *os.File {

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		// 如果不存在，则创建文件夹
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		// 如果没有权限，则报错
		log.Fatalf("Fail to OpenFile : %v", err)
	}

	// 打开文件，追加模式，如果文件不存在则创建
	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// 如果没有权限，则报错
		log.Fatalf("Fail to OpenFile : %v", err)
	}

	return handle
}

// mkDir 创建文件
func mkDir(filePath string) {

	// 获取当前工作目录
	dir, _ := os.Getwd()

	// 创建文件夹
	err := os.MkdirAll(dir+"/"+filePath, os.ModePerm)
	if err != nil {
		// 如果没有权限，则报错
		panic(err)
	}

}
