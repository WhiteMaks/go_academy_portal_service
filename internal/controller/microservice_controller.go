package controller

import (
	"net/http"

	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/gin-gonic/gin"
)

type MicroserviceController interface {
	GetMicroserviceV1(ctx *gin.Context)
}

type microserviceController struct {
	microserviceService service.MicroserviceService
}

func NewMicroserviceController(service service.MicroserviceService) MicroserviceController {
	return &microserviceController{microserviceService: service}
}

func (m microserviceController) GetMicroserviceV1(ctx *gin.Context) {
	response, err := m.microserviceService.GetMicroserviceV1(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
