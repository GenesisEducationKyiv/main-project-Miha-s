package logger

import (
	"btc-test-task/internal/common/configuration/config"
	"errors"
	"os"
)

type FileLoggerWriter struct {
	file *os.File
}

func NewFileLoggerWriter(conf *config.Config) (*FileLoggerWriter, error) {
	if len(conf.LogFile) == 0 {
		return nil, errors.New("empty log file path")
	}

	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Log.Errorf("Failed to log to file %v, error %v", conf.LogFile, err)
		return nil, err
	}

	return &FileLoggerWriter{
		file: file,
	}, nil
}

func (logger *FileLoggerWriter) Write(p []byte) (n int, err error) {
	return logger.file.Write(p)
}

func (logger *FileLoggerWriter) Close() {
	logger.file.Close()
}
