package config

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServerBindAddr string `mapstructure:"SERVER_BIND_ADDR" validate:"required"`
	EthereumHost   string `mapstructure:"ETHEREUM_HOST" validate:"required"`
}

func NewConfig() (config *Config) {
	// set default value
	viper.SetDefault("SERVER_BIND_ADDR", ":8080")

	// bind env variable
	tags := GetMapStructureTag(Config{})
	for _, tag := range tags {
		viper.BindEnv(tag)
	}

	// load env file
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warning(err)
		} else {
			log.Fatal(err)
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("NewConfig failed: %v\n", err)
	}

	// validation
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		log.Fatalf("NewConfig failed: %v\n", err)
	}
	return &c
}
