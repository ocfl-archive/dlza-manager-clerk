package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
)

type Service struct {
	ServiceName    string         `yaml:"service_name" toml:"ServiceName"`
	Host           string         `yaml:"host" toml:"Host"`
	Port           int            `yaml:"port" toml:"Port"`
	DatabaseConfig DatabaseConfig `yaml:"database_config" toml:"DatabaseConfig"`
}

type Config struct {
	GraphQLConfig models.GraphQLConfig `yaml:"graphql_config" toml:"GraphQLConfig"`
	Handler       Service              `yaml:"handler" toml:"Handler"`
	Ingester      Service              `yaml:"ingester" toml:"Ingester"`
	Clerk         Service              `yaml:"clerk" toml:"Clerk"`
}

// GetConfig creates a new config from a given environment
func GetConfig(configFile string, fileType string) (config Config, err error) {
	if configFile == "" {
		fmt.Println("here 0")
		defaultConfig := "config.yml"
		if fileType == "toml" {
			defaultConfig = "config.toml"
			fmt.Println("here")
		}
		err = configor.Load(&config, defaultConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = configor.Load(&config, configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}
