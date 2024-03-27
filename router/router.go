package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/auth"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller"
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
