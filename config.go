package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConfigActions interface {
	Setup()
}

type MyConfig struct {
	PORT     string
	DATABASE struct {
		TYPE     string
		USERNAME string
		PASSWORD string
		HOSTNAME string
		PORT     string
		SCHEMA   string
	}
	JWT struct {
		SECRET           string
		SESSION_DURATION int
	}
	RABBITMQ struct {
		USERNAME string
		PASSWORD string
		HOST     string
		PORT     string
	}
}

func (config *MyConfig) Setup() {
	viper.SetConfigName("config") //config filename without the .JSON or .YAML extension
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	os.Setenv("PORT", config.PORT)
}
