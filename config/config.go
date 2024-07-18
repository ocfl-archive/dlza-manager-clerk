package config

import (
	"github.com/jinzhu/configor"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	"log"
	"os"
)

type Service struct {
	ServiceName string `yaml:"service_name" toml:"ServiceName"`
	Host        string `yaml:"host" toml:"Host"`
	Port        int    `yaml:"port" toml:"Port"`
}

type Config struct {
	GraphQLConfig  models.GraphQLConfig `yaml:"graphql_config" toml:"GraphQLConfig"`
	Handler        Service              `yaml:"handler" toml:"Handler"`
	StorageHandler Service              `yaml:"storage-handler" toml:"StorageHandler"`
	Clerk          Service              `yaml:"clerk" toml:"Clerk"`
	Jwt            string               `yaml:"jwt-key" toml:"JwtKey"`
}

// GetConfig creates a new config from a given environment
func GetConfig(configFile string) Config {
	conf := Config{}
	if configFile == "" {
		configFile = "config.yml"
	}
	err := configor.Load(&conf, configFile)
	if err != nil {
		log.Fatal(err)
	}
	if conf.Jwt == "" {
		conf.Jwt = os.Getenv("JWT_KEY")
	}
	return conf
}
