package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	DatabaseUrl string `json:"database_url" split_words:"true" required:"true"`
	Port        int    `json:"port" default:"8080"`
	JwtSecret   string `json:"jwt_secret" split_words:"true" required:"true"`
	AllowSignup bool   `json:"allow_signup" split_words:"true" default:"false"`
}

func LoadConfig(filename string) (*Config, error) {
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		if err := loadEnvironment(filename); err != nil {
			return nil, err
		}
	}
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}
