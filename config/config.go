package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
)

type Config struct {
	Server   models.Server        `yaml:"server"`
	Logger   models.LoggingConfig `yaml:"logger"`
	Handler  string               `yaml:"handler"`
	Ingester string               `yaml:"ingester"`
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
