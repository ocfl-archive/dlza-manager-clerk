package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	"net/http"
)

type StatusController struct {
	ClientClerkHandlerService pb.ClerkHandlerServiceClient
}

func NewStatusController(clientClerkHandlerService pb.ClerkHandlerServiceClient) *StatusController {
	return &StatusController{ClientClerkHandlerService: clientClerkHandlerService}
}

// CheckStatus godoc
// @Summary		Check status
// @Description	Checking status of upload
// @Security 	ApiKeyAuth
// @ID 			check-status
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/status/{id} [get]
func (s *StatusController) CheckStatus(ctx *gin.Context) {

	id := ctx.Param("id")
	status, err := s.ClientClerkHandlerService.CheckStatus(ctx, &pb.Id{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, models.ArchivingStatus{Id: status.Id, Status: status.Status, LastChanged: status.LastChanged})
}

// AlterStatus godoc
// @Summary		Alter status
// @Description	Altering status of upload
// @Security 	ApiKeyAuth
// @ID 			alter-status
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/status [patch]
func (s *StatusController) AlterStatus(ctx *gin.Context) {
	statusObject := pb.StatusObject{}
	err := ctx.ShouldBindJSON(&statusObject)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	_, err = s.ClientClerkHandlerService.AlterStatus(ctx, &statusObject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// CreateStatus godoc
// @Summary		Create status
// @Description	Creating status of upload
// @Security 	ApiKeyAuth
// @ID 			create-status
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/status [post]
func (s *StatusController) CreateStatus(ctx *gin.Context) {
	statusObject := pb.StatusObject{}
	err := ctx.ShouldBindJSON(&statusObject)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	id, err := s.ClientClerkHandlerService.CreateStatus(ctx, &statusObject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, models.ArchivingStatus{Id: id.Id})
}