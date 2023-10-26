//go:generate swag init --parseDependency  --parseInternal -g .\swagg_config.go

package controller

import (
	_ "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller/docs"
)

// @title DLZA-archive API
// @version 1.0
// @description API for DLZA-archive

// @securityDefinitions.apikey	ApiKeyAuth
// @in 							header
// @name 						Authorization
// @description					Bearer Authentication with JWT

// @host localhost:8081
// @BasePath /api
