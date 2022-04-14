package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type (
	ConfigProvider interface {
		Config() *Configuration
	}

	Configuration struct {
		Server   ServerConfiguration
		Database DatabaseConfiguration
		SMTP     SmtpConfiguration
		Token    TokenConfiguration
		IEX      IEXConfiguration
	}

	ServerConfiguration struct {
		Port int
	}

	DatabaseConfiguration struct {
		Host     string
		Port     int
		Database string
		Username string
		Password string
	}

	SmtpConfiguration struct {
		Server   string
		Port     int
		Password string
	}

	TokenConfiguration struct {
		Key        string
		TokenValid uint32
	}

	IEXConfiguration struct {
		Token   string
		BaseUrl string
	}
)

func NewConfig() *Configuration {
	return unmarshalConfiguration()
}

func unmarshalConfiguration() *Configuration {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // look for config in the working directory
	var configuration Configuration
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal configuration object: %s", err))
	}

	log.Println("Read config file")
	return &configuration
}
