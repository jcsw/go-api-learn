package logger

import (
	"log"
	"os"
)

var logger = configureLogger()

func configureLogger() *log.Logger {
	return log.New(os.Stdout, "go-api-learn ", log.LstdFlags)
}

func Info(v ...interface{}) {
	logger.Println("INFO", (v))
}

func Fatal(v ...interface{}) {
	logger.Fatalln("FATAL", (v))
}
