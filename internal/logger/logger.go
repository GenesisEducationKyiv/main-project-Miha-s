package logger

import (
	"log"
	"os"
	"btc-test-task/internal/config"
)

var (
	warningLogger *log.Logger
	infoLogger *log.Logger
	errorLogger *log.Logger
)

func LogInfo(str string) {
	infoLogger.Println(str)
}

func LogWarn(str string) {
	warningLogger.Println(str)
}

func LogErrorStr(str string) {
	errorLogger.Println(str)
}

func LogError(err error) {
	errorLogger.Println(err.Error())
}

func Init(conf *config.Config) {
	warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
