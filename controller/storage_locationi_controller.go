package controller

import (
	"context"
	_ "github.com/ocfl-archive/dlza-manager-clerk/controller/docs"
	_ "github.com/ocfl-archive/dlza-manager-clerk/models"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StorageLocationController struct {
	ClientClerkHandler pbHandler.ClerkHandlerServiceClient
}

func (s *StorageLocationController) InitRoutes(storageLocationRouter *gin.RouterGroup) {
	storageLocationRouter.GET("/:id", s.GetStorageLocationsByTenantId)
	storageLocationRouter.POST("", s.SaveStorageLocation)
	storageLocationRouter.DELETE("/:id", s.DeleteStorageLocationById)
}

func (s *StorageLocationController) Path() string {
	return "/storage-location"
}

func NewStorageLocationController(clientClerkHandler pbHandler.ClerkHandlerServiceClient) Controller {
	return &StorageLocationController{ClientClerkHandler: clientClerkHandler}
}

// SaveStorageLocation godoc
// @Summary		Create storageLocation
// @Description	Add a new storageLocation
// @Security 	 ApiKeyAuth
// @ID create-storageLocation
// @Param		storageLocation's body models.StorageLocation true "Create storageLocation"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/storage-location [post]
func (s *StorageLocationController) SaveStorageLocation(ctx *gin.Context) {
	storageLocation := pb.StorageLocation{}
	err := ctx.ShouldBindJSON(&storageLocation)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err = s.ClientClerkHandler.SaveStorageLocation(cont, &storageLocation)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, storageLocation.Alias)
}

// DeleteStorageLocationById godoc
// @Summary		Delete storageLocation
// @Description	Delete a storageLocation
// @Security 	 ApiKeyAuth
// @ID 			delete-storageLocation
// @Param		id path string true "storage-location ID"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/storage-location/{id} [delete]
func (s *StorageLocationController) DeleteStorageLocationById(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")

	_, err := s.ClientClerkHandler.DeleteStorageLocationById(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// GetStorageLocationsByTenantId godoc
// @Summary		Find all storageLocations for tenant ID
// @Description	Finding all storageLocations for tenant ID
// @Security 	 ApiKeyAuth
// @ID 			find-all-storageLocations-for-tenant-id
// @Param		id path string true "tenant ID"
// @Produce		json
// @Success		200 {object} []models.StorageLocation
// @Failure 	400
// @Router		/storage-location/{id} [get]
func (s *StorageLocationController) GetStorageLocationsByTenantId(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")
	storageLocations, err := s.ClientClerkHandler.GetStorageLocationsByTenantId(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, storageLocations.StorageLocations)
}
