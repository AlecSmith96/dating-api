package adapters

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseConnectionString string `yaml:"database-connection-string" env:"DATABASE_CONNECTION_STRING" env-required:"true"`
	JwtExpiryMillis          int    `yaml:"jwt-expiry-millis" env:"JWT_EXPIRY_MILLIS" env-required:"true"`
	JwtSecretKey             string `yaml:"jwt-secret-key" env:"JWT_SECRET_KEY" env-required:"true"`
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
