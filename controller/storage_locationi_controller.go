package controller

import (
	"context"
	_ "github.com/ocfl-archive/dlza-manager-clerk/controller/docs"
	_ "github.com/ocfl-archive/dlza-manager-clerk/models"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StorageLocationController struct {
	ClientClerkHandler pbHandler.ClerkHandlerServiceClient
}

func (s *StorageLocationController) InitRoutes(storageLocationRouter *gin.RouterGroup) {
	storageLocationRouter.GET("/collection/:alias/:size/:signature/:head", s.GetStorageLocationsStatusForCollectionAlias)
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

// GetStorageLocationsStatusForCollectionAlias godoc
// @Summary		Get storageLocations status for collection alias
// @Description	It will give a message if there is no possibility to archive data into some storage partition, otherwise it would deliver an empty string
// @Security 	ApiKeyAuth
// @ID 			get-storage-locations-status-for-collection-alias
// @Param		alias string true "collection alias"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/storage-location/collection/{alias}/{size}/{signature}/{head} [get]
func (s *StorageLocationController) GetStorageLocationsStatusForCollectionAlias(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	alias := ctx.Param("alias")
	size := ctx.Param("size")
	signature := ctx.Param("signature")
	head := ctx.Param("head")
	sizeInt64, _ := strconv.ParseInt(size, 10, 64)
	objectPb := &pb.Object{Signature: signature, Head: head}
	status, err := s.ClientClerkHandler.GetStorageLocationsStatusForCollectionAlias(cont, &pb.SizeAndId{Id: alias, Size: sizeInt64, Object: objectPb})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, status)
}
