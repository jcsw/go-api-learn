package infra

import (
	"log"
	"os"
)

var logger = GetConfiguredLogger()

// GetConfiguredLogger - Return a configured logger
func GetConfiguredLogger() *log.Logger {
	return log.New(os.Stdout, "go-api-learn ", log.LstdFlags)
}

func LogInfo(v ...interface{}) {
	logger.Println("INFO", (v))
}
