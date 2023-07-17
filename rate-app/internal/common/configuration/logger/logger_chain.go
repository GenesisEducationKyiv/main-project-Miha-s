package logger

type LoggerWriter interface {
	Write(p []byte) (n int, err error)
	Close()
}

type LoggerWriterChain struct {
	current LoggerWriter
	next    LoggerWriter
}

func NewLoggerWriterChain(current LoggerWriter) *LoggerWriterChain {
	return &LoggerWriterChain{
		current: current,
	}
}

func (logger *LoggerWriterChain) SetNext(next LoggerWriter) {
	logger.next = next
}

func (logger *LoggerWriterChain) Write(p []byte) (n int, err error) {
	if logger.next != nil {
		count, err := logger.next.Write(p)
		if err != nil {
			return count, err
		}
	}

	return logger.current.Write(p)
}

func (logger *LoggerWriterChain) Close() {
	logger.next.Close()
	logger.current.Close()
}
