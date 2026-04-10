package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
)

type ObjectController struct {
	ClientClerkHandlerService pbHandler.ClerkHandlerServiceClient
}

func (o *ObjectController) InitRoutes(StorageInfoRouter *gin.RouterGroup) {
	StorageInfoRouter.GET("/:checksum", o.GetObjectsByChecksum)
	StorageInfoRouter.GET("/resulting-quality/:id", o.GetResultingQualityForObject)
	StorageInfoRouter.GET("/needed-quality/:id", o.GetNeededQualityForObject)
	StorageInfoRouter.GET("/signature/:signature", o.GetObjectBySignature)
	StorageInfoRouter.POST("/create", o.CreateObjectAndInstance)
}

func (o *ObjectController) Path() string {
	return "/object"
}

func NewObjectController(clientClerkHandlerService pbHandler.ClerkHandlerServiceClient) Controller {
	return &ObjectController{ClientClerkHandlerService: clientClerkHandlerService}
}

// GetObjectsByChecksum godoc
// @Summary		Getting objects by checksum
// @Description	Getting objects by checksum
// @Security 	ApiKeyAuth
// @ID 			objects-by-checksum
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object/{checksum} [get]
func (o *ObjectController) GetObjectsByChecksum(ctx *gin.Context) {

	checksum := ctx.Param("checksum")
	objects, err := o.ClientClerkHandlerService.GetObjectsByChecksum(ctx, &pb.Id{Id: checksum})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objects)
}

// GetObjectBySignature godoc
// @Summary		Getting objects by signature
// @Description	Getting objects by signature
// @Security 	ApiKeyAuth
// @ID 			objects-by-signature
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object/signature/{signature} [get]
func (o *ObjectController) GetObjectBySignature(ctx *gin.Context) {

	signature := ctx.Param("signature")
	objects, err := o.ClientClerkHandlerService.GetObjectBySignature(ctx, &pb.Id{Id: signature})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objects)
}

// GetResultingQualityForObject godoc
// @Summary		Getting resulting quality by object id
// @Description	Getting resulting quality by object id
// @Security 	ApiKeyAuth
// @ID 			resulting-quality-by-object-id
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object/resulting-quality/{id} [get]
func (o *ObjectController) GetResultingQualityForObject(ctx *gin.Context) {

	checksum := ctx.Param("id")
	quality, err := o.ClientClerkHandlerService.GetResultingQualityForObject(ctx, &pb.Id{Id: checksum})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, quality)
}

// GetNeededQualityForObject godoc
// @Summary		Getting needed quality by object id
// @Description	Getting needed quality by object id
// @Security 	ApiKeyAuth
// @ID 			resulting-needed-by-object-id
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object/needed-quality/{id} [get]
func (o *ObjectController) GetNeededQualityForObject(ctx *gin.Context) {

	checksum := ctx.Param("id")
	quality, err := o.ClientClerkHandlerService.GetNeededQualityForObject(ctx, &pb.Id{Id: checksum})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, quality)
}

// CreateObjectAndInstance godoc
// @Summary		Create object with instance
// @Description	Creating object with instance
// @Security 	ApiKeyAuth
// @ID 			create-object-with-instance
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object/create [post]
func (o *ObjectController) CreateObjectAndInstance(ctx *gin.Context) {
	//statusId field is used to transfer partition id
	object := pb.ObjectAndFile{}
	err := ctx.ShouldBindJSON(&object)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	_, err = o.ClientClerkHandlerService.CreateObjectAndInstance(ctx, &object)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
