package logging

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file:   file,
		logger: log.New(file, "", 0), // No default flags, we'll handle timestamps manually
	}, nil
}

func (l *Logger) logWithTimestamp(prefix, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("%s [%s] %s\n", prefix, timestamp, message)
}

func (l *Logger) Info(message string) {
	l.logWithTimestamp("INFO", message)
}

func (l *Logger) Error(message string) {
	l.logWithTimestamp("ERROR", message)
}

func (l *Logger) Debug(message string) {
	l.logWithTimestamp("DEBUG", message)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
