package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/dotenv-org/godotenvvault"
)

type envConfig struct {
	HTTPPort       string `env:"HTTP_PORT"`
	DBConfig       string `env:"DB_CONFIG"`
	Environment    string `env:"ENVIRONMENT"`
	SessionSecret  string `env:"SESSION_SECRET"`
	RedisHost      string `env:"REDIS_HOST"`
	RedisPassword  string `env:"REDIS_PASSWORD"`
	ClerkSecretKey string `env:"CLERK_SECRET_KEY"`
}

type Config struct {
	envConfig envConfig
}

func InitConfig(filepath string) (Config, error) {
	config := Config{}
	envConfig := envConfig{}

	err := godotenvvault.Load(filepath)
	if err != nil {
		return config, err
	}

	if err := env.Parse(&envConfig); err != nil {
		return config, err
	}

	config.envConfig = envConfig

	return config, nil
}

func (config Config) GetHTTPPort() string {
	return config.envConfig.HTTPPort
}

func (config Config) GetDBConfig() string {
	return config.envConfig.DBConfig
}

func (config Config) GetEnvironment() string {
	return config.envConfig.Environment
}

func (config Config) GetSessionSecret() string {
	return config.envConfig.SessionSecret
}

func (config Config) GetRedisHost() string {
	return config.envConfig.RedisHost
}

func (config Config) GetRedisPassword() string {
	return config.envConfig.RedisPassword
}

func (config Config) GetClerkSecretKey() string {
	return config.envConfig.ClerkSecretKey
}
