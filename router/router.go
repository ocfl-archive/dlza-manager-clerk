package router

import (
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/auth"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(tenantController *controller.TenantController,
	storageLocationController *controller.StorageLocationController, collectionController *controller.CollectionController,
	storagePartitionController *controller.StoragePartitionController, statusController *controller.StatusController) *gin.Engine {
	router := gin.Default()

	//Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRouter := router.Group("/api")
	baseRouter.Use(auth.JwtAuthMiddleware())

	tenantRouter := baseRouter.Group("/tenant")
	tenantRouter.GET("", tenantController.FindAllTenants)
	tenantRouter.GET("/:id", tenantController.FindTenantById)
	tenantRouter.POST("", tenantController.SaveTenant)
	tenantRouter.PATCH("", tenantController.UpdateTenant)
	tenantRouter.DELETE("/:id", tenantController.DeleteTenant)

	storageLocationRouter := baseRouter.Group("/storage-location")
	storageLocationRouter.GET("/:id", storageLocationController.GetStorageLocationsByTenantId)
	storageLocationRouter.POST("", storageLocationController.SaveStorageLocation)
	storageLocationRouter.DELETE("/:id", storageLocationController.DeleteStorageLocationById)

	storagePartitionRouter := baseRouter.Group("/storage-partition")
	storagePartitionRouter.POST("", storagePartitionController.CreateStoragePartition)

	collectionRouter := baseRouter.Group("/collection")
	collectionRouter.GET("/:id", collectionController.GetCollectionsByTenantId)
	collectionRouter.POST("", collectionController.CreateCollection)
	collectionRouter.PATCH("", collectionController.UpdateCollection)
	collectionRouter.DELETE("/:id", collectionController.DeleteCollectionById)

	statusRouter := baseRouter.Group("/status")
	statusRouter.GET("/:id", statusController.CheckStatus)
	statusRouter.POST("", statusController.CreateStatus)
	statusRouter.PATCH("", statusController.AlterStatus)

	return router
}
