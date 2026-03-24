package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
)

type ObjectInstanceController struct {
	ClientClerkHandlerService pbHandler.ClerkHandlerServiceClient
}

func (o *ObjectInstanceController) InitRoutes(StorageInfoRouter *gin.RouterGroup) {
	StorageInfoRouter.GET("/:name", o.ObjectInstanceWithNameExists)
	StorageInfoRouter.GET("/signature-and-location/:signature/:locations-name", o.GetObjectInstancesBySignatureAndLocationsPathName)
	StorageInfoRouter.GET("/raw-check/:object-id", o.GetRawObjectInstanceByObjectId)
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
// @Summary		Getting object instance by signature of object and storage locations name
// @Description	Getting object instance by signature of object and storage locations name
// @Security 	ApiKeyAuth
// @ID 			object-instance-by-signature-and-locations-name
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/object-instance/signature-and-location/{signature}/{locations-name} [get]
func (o *ObjectInstanceController) GetObjectInstancesBySignatureAndLocationsPathName(ctx *gin.Context) {
	signature := ctx.Param("signature")
	locationsName := ctx.Param("locations-name")
	objectInstance, err := o.ClientClerkHandlerService.GetObjectInstancesBySignatureAndLocationsPathName(ctx, &pb.AliasAndLocationsName{Alias: signature, LocationsName: locationsName})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objectInstance)
}

// GetRawObjectInstanceByObjectId godoc
// @Summary	Getting object instance by object ID if it has raw status
// @Description	Getting object instance by object ID if it has raw status
// @Security	ApiKeyAuth
// @ID	check-raw-object-instance-by-object-id
// @Produce	json
// @Success	200
// @Failure	400
// @Router	/object-instance/raw-check/{object-id} [get]
func (o *ObjectInstanceController) GetRawObjectInstanceByObjectId(ctx *gin.Context) {
	objectId := ctx.Param("object-id")
	objectInstance, err := o.ClientClerkHandlerService.CheckRawObjectInstanceByObjectId(ctx, &pb.Id{Id: objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, objectInstance)
}
