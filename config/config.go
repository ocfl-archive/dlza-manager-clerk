package config

import (
	"github.com/jinzhu/configor"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	"log"
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
}

// GetConfig creates a new config from a given environment
func GetConfig(configFile string) (config Config, err error) {
	if configFile == "" {
		configFile = "config.yml"
	}
	err = configor.Load(&config, configFile)
	if err != nil {
		log.Fatal(err)
	}
	return
}
