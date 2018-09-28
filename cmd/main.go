package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/jcsw/go-api-learn/pkg/application"
)

var env string

func main() {
	flag.StringVar(&env, "env", "prod", "app environment")
	flag.Parse()

	app := application.App{}
	app.Initialize(env)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		app.Stop()
	}()

	app.Start()
}
