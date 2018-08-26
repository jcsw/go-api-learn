package main

import (
	"flag"

	"github.com/jcsw/go-api-learn/pkg/application"
)

var env string

func main() {
	flag.StringVar(&env, "env", "prod", "app environment")
	flag.Parse()

	app := application.App{}
	app.Initialize(env)
	app.Run()
}
