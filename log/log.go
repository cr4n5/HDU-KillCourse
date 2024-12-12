package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

var logDir = "log_files"

func init() {
	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("无法创建日志目录: %v", err)
	}

	appFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开app.log文件: %v", err)
	}

	debugFile, err := os.OpenFile(filepath.Join(logDir, "debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开debug.log文件: %v", err)
	}

	appMultiWriter := io.MultiWriter(appFile, os.Stdout)
	debugMultiWriter := io.MultiWriter(debugFile)

	infoLogger = log.New(appMultiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLogger = log.New(appMultiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	debugLogger = log.New(debugMultiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func Info(v ...interface{}) {
	infoLogger.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func Error(v ...interface{}) {
	errorLogger.Println(append([]interface{}{"[ERROR]"}, v...)...)
}

func Debug(v ...interface{}) {
	debugLogger.Println(append([]interface{}{"[DEBUG]"}, v...)...)
}
