package controller

import (
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService service.OrderService
}

func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{OrderService: orderService}
}

// CopyFiles godoc
// @Summary		Copy files
// @Description	Copying all files from request
// @Security 	ApiKeyAuth
// @ID 			copy-files
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/order [post]
func (o *OrderController) CopyFiles(ctx *gin.Context) {
	incomingOrder := models.IncomingOrder{}
	err := ctx.ShouldBindJSON(&incomingOrder)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	_, err = o.OrderService.CopyFiles(incomingOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "request failed"})
		return
	}
	ctx.JSON(http.StatusOK, "Ok")
}
