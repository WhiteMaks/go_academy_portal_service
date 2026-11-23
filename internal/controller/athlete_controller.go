package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
)

type AthleteController interface {
	GetAthletesV1(ctx *gin.Context)
	CreateAthleteV1(ctx *gin.Context)
}

type athleteController struct {
	athleteService service.AthleteService
	userService    service.UserService
}

func NewAthleteController(athleteService service.AthleteService, userService service.UserService) AthleteController {
	return &athleteController{athleteService: athleteService, userService: userService}
}

func (c *athleteController) GetAthletesV1(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return
	}
	response, err := c.athleteService.GetAthletesV1(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *athleteController) CreateAthleteV1(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return
	}
	var request model.PostAthleteV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid request body")))
		return
	}
	if !isUTCMidnight(request.BirthDate) {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid request body")))
		return
	}
	response, err := c.athleteService.CreateAthleteV1(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func isUTCMidnight(value time.Time) bool {
	_, offset := value.Zone()
	return !value.IsZero() && value.Year() >= 1 && offset == 0 &&
		value.Hour() == 0 && value.Minute() == 0 && value.Second() == 0 && value.Nanosecond() == 0
}

func (c *athleteController) isAdmin(ctx *gin.Context) bool {
	payload := ctx.MustGet(middleware.AuthorizationPayload).(*util.Payload)
	return c.userService.HasAccess(ctx, payload, []database.RootUserRole{database.RootUserRoleAdmin})
}
