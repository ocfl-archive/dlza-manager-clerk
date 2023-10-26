package controller

import (
	"context"
	_ "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller/docs"
	_ "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CollectionController struct {
	ClientClerkHandler pb.ClerkHandlerServiceClient
}

func NewCollectionController(clientClerkIngestHandler pb.ClerkHandlerServiceClient) *CollectionController {
	return &CollectionController{ClientClerkHandler: clientClerkIngestHandler}
}

// CreateCollection godoc
// @Summary		Create collection
// @Description	Add a new collection
// @Security 	 ApiKeyAuth
// @ID 			create-collection
// @Param		collection's body models.Collection true "Create collection"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/collection [post]
func (col *CollectionController) CreateCollection(ctx *gin.Context) {
	collection := &pb.Collection{}
	err := ctx.ShouldBindJSON(&collection)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err = col.ClientClerkHandler.CreateCollection(cont, collection)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// UpdateCollection godoc
// @Summary		Update collection
// @Description	Update a collection
// @Security 	 ApiKeyAuth
// @ID 			update-collection
// @Param		collection's body models.Collection true "Update collection"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/collection [patch]
func (col *CollectionController) UpdateCollection(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	collection := pb.Collection{}
	err := ctx.ShouldBindJSON(&collection)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	_, err = col.ClientClerkHandler.UpdateCollection(cont, &collection)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// DeleteCollectionById godoc
// @Summary		Delete collection
// @Description	Delete a collection
// @Security 	 ApiKeyAuth
// @ID 			delete-collection
// @Param		id path string true "collection ID"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/collection/{id} [delete]
func (col *CollectionController) DeleteCollectionById(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")

	_, err := col.ClientClerkHandler.DeleteCollectionById(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// GetCollectionsByTenantId godoc
// @Summary		Find collections by tenant id
// @Description	Finding collections by tenant id
// @Security 	 ApiKeyAuth
// @ID 			find-collections-tenant-id
// @Param		id path string true "tenant ID"
// @Produce		json
// @Success		200 {object} []models.Collection
// @Failure 	400
// @Router		/collection/{id} [get]
func (col *CollectionController) GetCollectionsByTenantId(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")
	collections, err := col.ClientClerkHandler.GetCollectionsByTenantId(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, collections.Collections)
}
