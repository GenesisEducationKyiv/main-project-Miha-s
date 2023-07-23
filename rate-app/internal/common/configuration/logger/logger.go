package logger

import (
	"btc-test-task/internal/common/configuration/config"
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Log          *logrus.Logger
	loggingLevel = map[string]logrus.Level{
		"TRACE": logrus.TraceLevel,
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"FATAL": logrus.FatalLevel,
		"PANIC": logrus.PanicLevel,
	}
)

type CustomFormatter struct{}

func (mf *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	b.WriteString(fmt.Sprintf("%s: %s %s:%d %s\n",
		entry.Level, entry.Time.Format("2006-01-02 15:04:05,6"), fileName, entry.Caller.Line, entry.Message))
	return b.Bytes(), nil
}

func Init(conf *config.Config, writer LoggerWriter) error {
	Log = logrus.New()
	Log.SetFormatter(&CustomFormatter{})
	Log.SetReportCaller(true)

	Log.Level = getLogLevel(conf.LogLevel)
	Log.Out = writer

	return nil
}

func getLogLevel(logLevel string) logrus.Level {
	level, ok := loggingLevel[logLevel]
	if !ok {
		return logrus.TraceLevel
	} else {
		return level
	}
}
