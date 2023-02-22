package main

import (
	_ "embed"

	console "github.com/purerun/dbhelper/cmd"
	"github.com/purerun/dbhelper/util/app"
)

//go:embed .env.example
var envExample string

func main() {
	app.InitStart()
	app.EnvExample(envExample)

	console.Execute()

}
