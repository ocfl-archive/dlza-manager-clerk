package controller

import (
	"github.com/gin-gonic/gin"
	pb "gitlab.switch.ch/ub-unibas/dlza/dlza-manager/dlzamanagerproto"
	pbHandler "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-handler/handlerproto"
	"net/http"
)

type ObjectController struct {
	ClientClerkHandlerService pbHandler.ClerkHandlerServiceClient
}

func (o *ObjectController) InitRoutes(StorageInfoRouter *gin.RouterGroup) {
	StorageInfoRouter.GET("/:checksum", o.GetObjectsByChecksum)
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