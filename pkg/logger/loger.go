package logger

import (
	"log"
	"os"
	"runtime"
	"time"
)

// Logger is a reusable logger with file output
type Logger struct {
	LogFile *os.File
}

// NewLogger initializes a logger writing to a file
func Create_Logger() (*Logger, error) {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file) // Set global logging output to the file
	return &Logger{LogFile: file}, nil
}

// Close closes the log file
func (l *Logger) Close() {
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}

func LogWithDetails(message error) {
	// Get the current caller details (the calling function)
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		log.Println("Failed to get caller information")
	}
	// Get the function name from the program counter (pc)
	funcName := runtime.FuncForPC(pc).Name()
	// Log the message with the function name and line number
	log.Printf("%s [Function: %s] [Line: %d] %s", time.Now().Format(time.RFC3339), funcName, line, message)
}
