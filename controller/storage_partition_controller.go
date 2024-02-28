package controller

import (
	"context"
	pb "gitlab.switch.ch/ub-unibas/dlza/dlza-manager/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StoragePartitionController struct {
	ClientClerkStorageHandlerService pb.ClerkStorageHandlerServiceClient
}

func (s *StoragePartitionController) InitRoutes(storagePartitionRouter *gin.RouterGroup) {
	storagePartitionRouter.POST("", s.CreateStoragePartition)
}

func (s *StoragePartitionController) Path() string {
	return "/storage-partition"
}

func NewStoragePartitionController(clientClerkStorageHandlerService pb.ClerkStorageHandlerServiceClient) Controller {
	return &StoragePartitionController{ClientClerkStorageHandlerService: clientClerkStorageHandlerService}
}

// CreateStoragePartition godoc
// @Summary		Create storagePartition
// @Description	Add a new storagePartition
// @Security 	ApiKeyAuth
// @ID 			create-storagePartition
// @Param		tenant's body models.StoragePartition true "Create storagePartition"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/storage-partition [post]
func (s *StoragePartitionController) CreateStoragePartition(ctx *gin.Context) {
	storagePartition := pb.StoragePartition{}
	err := ctx.ShouldBindJSON(&storagePartition)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err = s.ClientClerkStorageHandlerService.CreateStoragePartition(cont, &storagePartition)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}
