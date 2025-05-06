package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ocfl-archive/dlza-manager-clerk/auth"
	"github.com/ocfl-archive/dlza-manager-clerk/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(key string, controllers ...controller.Controller) *gin.Engine {
	router := gin.Default()

	//Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRouter := router.Group("/api")
	baseRouter.Use(auth.JwtAuthMiddleware(key))

	for _, cntr := range controllers {
		subRouter := baseRouter.Group(cntr.Path())
		cntr.InitRoutes(subRouter)
	}

	return router
}
