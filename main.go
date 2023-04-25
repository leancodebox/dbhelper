package main

import (
	_ "embed"

	console "github.com/leancodebox/dbhelper/cmd"
	"github.com/leancodebox/dbhelper/util/app"
)

//go:embed .env.example
var envExample string

func main() {
	app.InitStart()
	app.EnvExample(envExample)

	console.Execute()

}
