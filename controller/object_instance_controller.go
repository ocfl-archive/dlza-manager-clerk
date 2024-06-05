package controller

import (
	"github.com/gin-gonic/gin"
	pb "gitlab.switch.ch/ub-unibas/dlza/dlza-manager/dlzamanagerproto"
	pbHandler "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-handler/handlerproto"
	"net/http"
)

type ObjectInstanceController struct {
	ClientClerkHandlerService pbHandler.ClerkHandlerServiceClient
}

func (o *ObjectInstanceController) InitRoutes(StorageInfoRouter *gin.RouterGroup) {
	StorageInfoRouter.GET("/:name", o.ObjectInstanceWithNameExists)
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
