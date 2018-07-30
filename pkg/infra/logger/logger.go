package logger

import (
	"log"
	"os"
)

var logger = configureLogger()

func configureLogger() *log.Logger {
	return log.New(os.Stdout, "go-api-learn ", log.LstdFlags)
}

// Debug - Logging in level DEBUG
func Debug(v ...interface{}) {
	logger.Println("DEBUG", (v))
}

// Info - Logging in level INFO
func Info(v ...interface{}) {
	logger.Println("INFO", (v))
}

// Warn - Logging in level WARN
func Warn(v ...interface{}) {
	logger.Println("WARN", (v))
}

// Error - Logging in level ERROR
func Error(v ...interface{}) {
	logger.Println("ERROR", (v))
}

// Fatal - Logging in level FATAL
func Fatal(v ...interface{}) {
	logger.Fatalln("FATAL", (v))
}
