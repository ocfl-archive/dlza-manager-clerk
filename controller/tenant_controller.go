package controller

import (
	"context"
	pb "gitlab.switch.ch/ub-unibas/dlza/dlza-manager/dlzamanagerproto"
	_ "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller/docs"
	_ "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pbHandler "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-handler/handlerproto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTenantController(clientClerkHandler pbHandler.ClerkHandlerServiceClient) Controller {
	return &TenantController{ClientClerkHandler: clientClerkHandler}
}

type TenantController struct {
	ClientClerkHandler pbHandler.ClerkHandlerServiceClient
}

func (t *TenantController) Path() string {
	return "/tenant"
}

func (t *TenantController) InitRoutes(tenantRouter *gin.RouterGroup) {

	tenantRouter.GET("", t.FindAllTenants)
	tenantRouter.GET("/:id", t.FindTenantById)
	tenantRouter.POST("", t.SaveTenant)
	tenantRouter.PATCH("", t.UpdateTenant)
	tenantRouter.DELETE("/:id", t.DeleteTenant)
}

// SaveTenant godoc
// @Summary		Create tenant
// @Description	Add a new tenant
// @Security 	 ApiKeyAuth
// @ID create-tenant
// @Param		tenant's body models.Tenant true "Create tenant"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/tenant [post]
func (t *TenantController) SaveTenant(ctx *gin.Context) {
	tenant := pb.Tenant{}
	err := ctx.ShouldBindJSON(&tenant)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	_, err = t.ClientClerkHandler.SaveTenant(cont, &tenant)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, tenant.Alias)
}

// UpdateTenant godoc
// @Summary		Update tenant
// @Description	Update a tenant
// @Security 	 ApiKeyAuth
// @ID update-tenant
// @Param		tenant's body models.Tenant true "Update tenant"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/tenant [patch]
func (t *TenantController) UpdateTenant(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	tenant := pb.Tenant{}
	err := ctx.ShouldBindJSON(&tenant)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "request failed"})
		return
	}
	_, err = t.ClientClerkHandler.UpdateTenant(cont, &tenant)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// DeleteTenant godoc
// @Summary		Delete tenant
// @Description	Delete a tenant
// @Security 	 ApiKeyAuth
// @ID delete-tenant
// @Param		id path string true "tenant ID"
// @Produce		json
// @Success		200
// @Failure 	400
// @Router		/tenant/{id} [delete]
func (t *TenantController) DeleteTenant(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")

	_, err := t.ClientClerkHandler.DeleteTenant(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// FindTenantById godoc
// @Summary		Find tenant by id
// @Description	Finding a tenant by id
// @Security 	 ApiKeyAuth
// @ID 			find-tenant-by-id
// @Param		id path string true "tenant ID"
// @Produce		json
// @Success		200 {object} models.Tenant
// @Failure 	400
// @Router		/tenant/{id} [get]
func (t *TenantController) FindTenantById(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	id := ctx.Param("id")
	tenant, err := t.ClientClerkHandler.FindTenantById(cont, &pb.Id{Id: id})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, tenant)
}

// FindAllTenants godoc
// @Summary		Find all tenants
// @Description	Finding all tenants
// @Security 	 ApiKeyAuth
// @ID 			find-all-tenants
// @Produce		json
// @Success		200 {object} []models.Tenant
// @Failure 	400
// @Router		/tenant [get]
func (t *TenantController) FindAllTenants(ctx *gin.Context) {
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	tenants, err := t.ClientClerkHandler.FindAllTenants(cont, &pb.NoParam{})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, tenants.Tenants)
}
