package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel                 string `mapstructure:"LOG_LEVEL"`
	ServerBindAddr           string `mapstructure:"SERVER_BIND_ADDR"`
	DbHost                   string `mapstructure:"DB_HOST" validate:"required"`
	DbPort                   int    `mapstructure:"DB_PORT" validate:"required"`
	DbName                   string `mapstructure:"DB_NAME" validate:"required"`
	DbUser                   string `mapstructure:"DB_USER" validate:"required"`
	DbPassword               string `mapstructure:"DB_PASSWORD" validate:"required"`
	EthereumHost             string `mapstructure:"ETHEREUM_HOST" validate:"required"`
	EthereumBlockInitFrom    int    `mapstructure:"ETHEREUM_BLOCK_INIT_FROM" validate:"required"`
	EthereumBlockWorkerCount int    `mapstructure:"ETHEREUM_BLOCK_WORKER_COUNT"`
}

func NewConfig() (config *Config) {
	// set default value
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_BIND_ADDR", ":8080")
	viper.SetDefault("ETHEREUM_BLOKC_WORKER_COUNT", 100)

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

func (c *Config) GetDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		c.DbUser,
		c.DbPassword,
		c.DbHost,
		c.DbPort,
		c.DbName,
	)
}
