package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	InitRoutes(router *gin.RouterGroup)
	Path() string
}
