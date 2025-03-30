package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/dotenv-org/godotenvvault"
)

type envConfig struct {
	HTTPPort              string `env:"HTTP_PORT"`
	DBConfig              string `env:"DB_CONFIG"`
	TestDBConfig          string `env:"TEST_DB_CONFIG"`
	Environment           string `env:"ENVIRONMENT"`
	SessionSecret         string `env:"SESSION_SECRET"`
	RedisHost             string `env:"REDIS_HOST"`
	RedisPassword         string `env:"REDIS_PASSWORD"`
	ClerkSecretKey        string `env:"CLERK_SECRET_KEY"`
	ClerkSigningSecretKey string `env:"CLERK_SIGNING_SECRET_KEY"`
	AWSAccessKeyID        string `env:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey    string `env:"AWS_SECRET_ACCESS_KEY"`
	TwilioAccountSID      string `env:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken       string `env:"TWILIO_AUTH_TOKEN"`
	TwilioPhoneNumber     string `env:"TWILIO_PHONE_NUMBER"`
}

type Config struct {
	envConfig envConfig
}

func InitConfig(filepath string) (Config, error) {
	config := Config{}
	envCfg := envConfig{}

	// load from filepath
	err := godotenvvault.Load(filepath)
	if err != nil {
		// fallback to os.Getenv if filepath fails, this is for dokku
		envCfg = envConfig{
			HTTPPort:              os.Getenv("HTTP_PORT"),
			DBConfig:              os.Getenv("DB_CONFIG"),
			TestDBConfig:          os.Getenv("TEST_DB_CONFIG"),
			Environment:           os.Getenv("ENVIRONMENT"),
			SessionSecret:         os.Getenv("SESSION_SECRET"),
			RedisHost:             os.Getenv("REDIS_HOST"),
			RedisPassword:         os.Getenv("REDIS_PASSWORD"),
			ClerkSecretKey:        os.Getenv("CLERK_SECRET_KEY"),
			ClerkSigningSecretKey: os.Getenv("CLERK_SIGNING_SECRET_KEY"),
			AWSAccessKeyID:        os.Getenv("AWS_ACCESS_KEY_ID"),
			AWSSecretAccessKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			TwilioAccountSID:      os.Getenv("TWILIO_ACCOUNT_SID"),
			TwilioAuthToken:       os.Getenv("TWILIO_AUTH_TOKEN"),
			TwilioPhoneNumber:     os.Getenv("TWILIO_PHONE_NUMBER"),
		}
	} else {
		// If filepath loading succeeds, parse env vars
		if err := env.Parse(&envCfg); err != nil {
			return config, err
		}
	}

	config.envConfig = envCfg
	return config, nil
}

// Getter methods remain the same as in the previous implementation
func (config Config) GetHTTPPort() string {
	return config.envConfig.HTTPPort
}

func (config Config) GetDBConfig() string {
	return config.envConfig.DBConfig
}

func (config Config) GetTestDBConfig() string {
	return config.envConfig.TestDBConfig
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

func (config Config) GetClerkSigningSecretKey() string {
	return config.envConfig.ClerkSigningSecretKey
}

func (config Config) GetAWSAccessKey() string {
	return config.envConfig.AWSAccessKeyID
}

func (config Config) GetAWSSecretAccessKey() string {
	return config.envConfig.AWSSecretAccessKey
}

func (config Config) GetTwilioAccountSID() string {
	return config.envConfig.TwilioAccountSID
}

func (config Config) GetTwilioAuthToken() string {
	return config.envConfig.TwilioAuthToken
}

func (config Config) GetTwilioPhoneNumber() string {
	return config.envConfig.TwilioPhoneNumber
}
