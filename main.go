package main

import (
	_ "embed"

	console "github.com/leancodebox/dbhelper/cmd"
	"github.com/leancodebox/dbhelper/util/app"
)

//go:embed config.example.toml
var envExample string

func main() {
	app.InitStart()
	app.DefaultConfig(envExample)

	console.Execute()

}
