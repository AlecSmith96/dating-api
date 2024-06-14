package adapters

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseConnectionString string `yaml:"database-connection-string"`
	JwtExpiryMillis          int    `yaml:"jwt-expiry-millis"`
	JwtSecretKey             string `yaml:"jwt-secret-key"`
}

func NewConfig() (*Config, error) {
	var conf Config
	err := cleanenv.ReadConfig("dev-config.yaml", &conf)
	if err != nil {
		err = cleanenv.ReadEnv(&conf)
		if err != nil {
			return nil, err
		}
	}

	return &conf, nil
}
