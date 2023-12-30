package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	// api service
	ServerBindAddr string `mapstructure:"SERVER_BIND_ADDR"`
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
	return &c
}
