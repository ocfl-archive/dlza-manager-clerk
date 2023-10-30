package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
)

type HostPort struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Config struct {
	GraphQLConfig models.GraphQLConfig `yaml:"graphql_config"`
	Handler       HostPort             `yaml:"handler"`
	Ingester      HostPort             `yaml:"ingester"`
	Clerk         HostPort             `yaml:"clerk"`
}

// GetConfig creates a new config from a given environment
func GetConfig() (config Config, err error) {
	err = configor.Load(&config, "config.yml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config", config)
	return
}
