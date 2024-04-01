package main

import (
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"

	"stove-gitlab.sginfra.net/backend/template/app"
)

var Version = "None"

func main() {

	app, err := app.NewApp(Version)
	if err != nil {
		log.Error(err)
		return
	}

	if err := app.Start(); err != nil {
		log.Error(err)
		return
	}
}
