package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger

	infoFileLogger  *log.Logger
	errorFileLogger *log.Logger
)

var (
	InfoColor  = color.New(color.FgGreen).SprintFunc()
	ErrorColor = color.New(color.FgRed).SprintFunc()
)

var logDir = "log_files"

var ansiColorRegex = regexp.MustCompile("\x1b\\[[0-9;]*m")

// 清除ANSI颜色代码
func removeAnsiColors(args []interface{}) []interface{} {
	result := make([]interface{}, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			result[i] = ansiColorRegex.ReplaceAllString(v, "")
		default:
			result[i] = v
		}
	}
	return result
}

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

	// 控制台
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	// 文件
	infoFileLogger = log.New(appFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorFileLogger = log.New(appFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	// debug
	debugMultiWriter := io.MultiWriter(debugFile)
	debugLogger = log.New(debugMultiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func Info(v ...interface{}) {
	infoLogger.Println(append([]interface{}{InfoColor("[INFO]")}, v...)...)

	cleanedArgs := removeAnsiColors(v)
	infoFileLogger.Println(append([]interface{}{"[INFO]"}, cleanedArgs...)...)
}

func Error(v ...interface{}) {
	errorLogger.Println(append([]interface{}{ErrorColor("[ERROR]")}, v...)...)

	cleanedArgs := removeAnsiColors(v)
	errorFileLogger.Println(append([]interface{}{"[ERROR]"}, cleanedArgs...)...)
}

func Debug(v ...interface{}) {
	cleanedArgs := removeAnsiColors(v)
	debugLogger.Println(append([]interface{}{"[DEBUG]"}, cleanedArgs...)...)
}
