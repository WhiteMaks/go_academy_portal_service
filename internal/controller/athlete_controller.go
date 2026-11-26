package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AthleteController interface {
	CreateAthleteV1(ctx *gin.Context)
	GetAthletesV1(ctx *gin.Context)
	GetAthleteV1(ctx *gin.Context)
	UpdateAthleteV1(ctx *gin.Context)
}

type athleteController struct {
	athleteService service.AthleteService
	userService    service.UserService
}

func NewAthleteController(athleteService service.AthleteService, userService service.UserService) AthleteController {
	return &athleteController{
		athleteService: athleteService,
		userService:    userService,
	}
}

func (c *athleteController) CreateAthleteV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var request model.PostAthleteV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapAthleteCreateBindError(err)))
		return
	}

	response, err := c.athleteService.CreateAthleteV1(ctx.Request.Context(), request)
	if err != nil {
		writeAthleteServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *athleteController) GetAthletesV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var query model.GetAthletesV1Query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapAthleteListBindError(err)))
		return
	}

	response, err := c.athleteService.GetAthletesV1(ctx.Request.Context(), query)
	if err != nil {
		writeAthleteServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *athleteController) GetAthleteV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parseAthleteID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid athlete id")))
		return
	}

	response, err := c.athleteService.GetAthleteV1(ctx.Request.Context(), id)
	if err != nil {
		writeAthleteServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *athleteController) UpdateAthleteV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parseAthleteID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid athlete id")))
		return
	}

	var request model.PutAthleteV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapAthleteCreateBindError(err)))
		return
	}

	response, err := c.athleteService.UpdateAthleteV1(ctx.Request.Context(), id, request)
	if err != nil {
		writeAthleteServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *athleteController) requireAdmin(ctx *gin.Context) bool {
	payload := ctx.MustGet(middleware.AuthorizationPayload).(*util.Payload)
	hasAccess := c.userService.HasAccess(
		ctx,
		payload,
		[]database.RootUserRole{database.RootUserRoleAdmin},
	)
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return false
	}
	return true
}

func mapAthleteCreateBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrAthleteValidation
	}
	return err
}

func mapAthleteListBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrAthleteListQuery
	}
	return err
}

func parseAthleteID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid athlete id")
	}
	return id, nil
}

func writeAthleteServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrAthleteListQuery):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrAthleteValidation),
		errors.Is(err, service.ErrAthleteParentPartial):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrAthleteNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("athlete not found")))
	case errors.Is(err, service.ErrAthletePhoneConflict):
		ctx.JSON(http.StatusConflict, service.PrepareErrorResponse(service.ErrAthletePhoneConflict))
	default:
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
	}
}
