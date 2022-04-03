package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	SMTP     SmtpConfiguration
	Token    TokenConfiguration
	IEX      IEXConfiguration
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type SmtpConfiguration struct {
	Server   string
	Port     int
	Password string
}

type TokenConfiguration struct {
	Key        string
	TokenValid uint32
}

type IEXConfiguration struct {
	Token   string
	BaseUrl string
}

var configuration Configuration

var configOnce sync.Once

func ReadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func unmarshalConfiguration() {
	err := viper.Unmarshal(&configuration)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal configuration object: %s", err))
	}

	log.Println("Read config file")
}

func Config() *Configuration {
	configOnce.Do(unmarshalConfiguration)
	return &configuration
}
