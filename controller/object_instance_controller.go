package controller

import (
	"github.com/gin-gonic/gin"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"net/http"
)

type ObjectInstanceController struct {
	ClientClerkHandlerService pbHandler.ClerkHandlerServiceClient
}

func (o *ObjectInstanceController) InitRoutes(StorageInfoRouter *gin.RouterGroup) {
	StorageInfoRouter.GET("/:name", o.ObjectInstanceWithNameExists)
	StorageInfoRouter.GET("/alias-and-location/:alias/:locations-name", o.GetObjectInstancesBySignatureAndLocationsPathName)
}

func (o *ObjectInstanceController) Path() string {
	return "/object-instance"
}

func NewObjectInstanceController(clientClerkHandlerService pbHandler.ClerkHandlerServiceClient) Controller {
	return &ObjectInstanceController{ClientClerkHandlerService: clientClerkHandlerService}
}

// ObjectInstanceWithNameExists godoc
// @Summary		Getting object instance by name
// @Description	Getting object instance by name
// @Security 	ApiKeyAuth
// @ID 			object-instance-by-name
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object-instance/{name} [get]
func (o *ObjectInstanceController) ObjectInstanceWithNameExists(ctx *gin.Context) {

	name := ctx.Param("name")
	objectInstances, err := o.ClientClerkHandlerService.GetObjectInstancesByName(ctx, &pb.Id{Id: name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objectInstances)
}

// GetObjectInstancesBySignatureAndLocationsPathName godoc
// @Summary		Getting object instance by alias of object and storage locations name
// @Description	Getting object instance by alias of object and storage locations name
// @Security 	ApiKeyAuth
// @ID 			object-instance-by-alias-and-locations-name
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object-instance/alias-and-location/{alias}/{locations-name} [get]
func (o *ObjectInstanceController) GetObjectInstancesBySignatureAndLocationsPathName(ctx *gin.Context) {
	alias := ctx.Param("alias")
	locationsName := ctx.Param("locations-name")
	objectInstance, err := o.ClientClerkHandlerService.GetObjectInstancesBySignatureAndLocationsPathName(ctx, &pb.AliasAndLocationsName{Alias: alias, LocationsName: locationsName})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objectInstance)
}
