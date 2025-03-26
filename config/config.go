package config

import (
	"os"
)

func GetHTTPPort() string {
	return os.Getenv("HTTP_PORT")
}

func GetDBConfig() string {
	//return config.envConfig.DBConfig
	return os.Getenv("DB_CONFIG")
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetClerkSecretKey() string {
	return os.Getenv("CLERK_SIGNING_SECRET_KEY")
}

func GetClerkSigningSecretKey() string {
	return os.Getenv("DB_CONFIG")
}
