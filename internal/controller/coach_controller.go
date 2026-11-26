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

type CoachController interface {
	CreateCoachV1(ctx *gin.Context)
	GetCoachesV1(ctx *gin.Context)
	GetCoachV1(ctx *gin.Context)
	UpdateCoachV1(ctx *gin.Context)
}

type coachController struct {
	coachService service.CoachService
	userService  service.UserService
}

func NewCoachController(coachService service.CoachService, userService service.UserService) CoachController {
	return &coachController{
		coachService: coachService,
		userService:  userService,
	}
}

func (c *coachController) CreateCoachV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var request model.PostCoachV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapCoachCreateBindError(err)))
		return
	}

	response, err := c.coachService.CreateCoachV1(ctx.Request.Context(), request)
	if err != nil {
		writeCoachServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *coachController) GetCoachesV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var query model.GetCoachesV1Query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapCoachListBindError(err)))
		return
	}

	response, err := c.coachService.GetCoachesV1(ctx.Request.Context(), query)
	if err != nil {
		writeCoachServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *coachController) GetCoachV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parseCoachID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid coach id")))
		return
	}

	response, err := c.coachService.GetCoachV1(ctx.Request.Context(), id)
	if err != nil {
		writeCoachServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *coachController) UpdateCoachV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parseCoachID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid coach id")))
		return
	}

	var request model.PutCoachV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapCoachCreateBindError(err)))
		return
	}

	response, err := c.coachService.UpdateCoachV1(ctx.Request.Context(), id, request)
	if err != nil {
		writeCoachServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *coachController) requireAdmin(ctx *gin.Context) bool {
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

func mapCoachCreateBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrCoachValidation
	}
	return err
}

func mapCoachListBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrCoachListQuery
	}
	return err
}

func parseCoachID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid coach id")
	}
	return id, nil
}

func writeCoachServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrCoachListQuery):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrCoachValidation):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrCoachNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("coach not found")))
	case errors.Is(err, service.ErrCoachPhoneConflict):
		ctx.JSON(http.StatusConflict, service.PrepareErrorResponse(service.ErrCoachPhoneConflict))
	default:
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
	}
}
