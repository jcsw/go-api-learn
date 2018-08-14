package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/jcsw/go-api-learn/pkg/application"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

var env string

func main() {
	logger.Info("Server is starting...")

	flag.StringVar(&env, "env", "prod", "app environment")
	flag.Parse()

	app := application.App{}
	app.Initialize(env)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		app.Stop()
		close(done)
	}()

	app.Run()

	<-done
	logger.Info("Server stopped")
}
