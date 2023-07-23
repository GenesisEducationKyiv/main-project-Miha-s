package logger

import "os"

type ConsoleLoggerWriter struct {
}

func NewConsoleLoggerWriter() *ConsoleLoggerWriter {
	return &ConsoleLoggerWriter{}
}

func (logger *ConsoleLoggerWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (logger *ConsoleLoggerWriter) Close() {}
