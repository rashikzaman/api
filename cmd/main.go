package main

import (
	"rashikzaman/api/application"
	"rashikzaman/api/config"
	"rashikzaman/api/db"
	"rashikzaman/api/http"
	"rashikzaman/api/log"
)

func main() {
	app := application.Application{}
	logger := log.NewLogger()

	config, err := config.InitConfig("../.env")
	if err != nil {
		logger.Fatal(err, err.Error())
	}

	db, err := db.InitDB(config.GetDBConfig())
	if err != nil {
		logger.Fatal(err, err.Error())
	}

	if err != nil {
		logger.Fatal(err, "Failed to create session store")

		return
	}

	//app.Config = config
	app.DB = db
	app.Config = config

	http.RunHTTPServer(app)
}
