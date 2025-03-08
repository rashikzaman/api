package main

import (
	"rashikzaman/api/application"
	"rashikzaman/api/config"
	"rashikzaman/api/constants"
	"rashikzaman/api/db"
	"rashikzaman/api/http"
	"rashikzaman/api/log"

	"github.com/gin-contrib/sessions/redis"
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

	redisStore, err := redis.NewStore(
		constants.SessionRedisSize,
		"tcp",
		config.GetRedisHost(),
		config.GetRedisPassword(),
		[]byte(app.Config.GetSessionSecret()),
	)

	if err != nil {
		logger.Fatal(err, "Failed to create session store")

		return
	}

	app.Config = config
	app.DB = db

	http.RunHTTPServer(app, redisStore)
}
