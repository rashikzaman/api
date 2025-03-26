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

	db, err := db.InitDB(config.GetDBConfig())
	if err != nil {
		logger.Fatal(err, err.Error())
	}

	if err != nil {
		logger.Fatal(err, "Failed to create session store")

		return
	}

	app.DB = db

	http.RunHTTPServer(app)
}
