package infra

import (
	"log"
	"os"
)

var logger = GetLogger()

// GetLogger Return log
func GetLogger() *log.Logger {
	return log.New(os.Stdout, "go-api-learn ", log.LstdFlags)
}
