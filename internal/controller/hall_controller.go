package controller

import (
	"errors"
	"net/http"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HallController interface {
	CreateHallV1(ctx *gin.Context)
	GetBranchHallsV1(ctx *gin.Context)
	GetHallV1(ctx *gin.Context)
	UpdateHallV1(ctx *gin.Context)
}

type hallController struct {
	hallService service.HallService
	userService service.UserService
}

func NewHallController(hallService service.HallService, userService service.UserService) HallController {
	return &hallController{
		hallService: hallService,
		userService: userService,
	}
}

func (c *hallController) CreateHallV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var request model.PostHallV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapHallBindError(err)))
		return
	}

	response, err := c.hallService.CreateHallV1(ctx.Request.Context(), request)
	if err != nil {
		writeHallServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *hallController) GetBranchHallsV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	branchID, err := parsePositiveID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid branch id")))
		return
	}

	var query model.GetBranchHallsV1Query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapHallListBindError(err)))
		return
	}

	response, err := c.hallService.GetBranchHallsV1(ctx.Request.Context(), branchID, query)
	if err != nil {
		writeHallServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *hallController) GetHallV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parsePositiveID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid hall id")))
		return
	}

	response, err := c.hallService.GetHallV1(ctx.Request.Context(), id)
	if err != nil {
		writeHallServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *hallController) UpdateHallV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parsePositiveID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid hall id")))
		return
	}

	var request model.PutHallV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapHallBindError(err)))
		return
	}

	response, err := c.hallService.UpdateHallV1(ctx.Request.Context(), id, request)
	if err != nil {
		writeHallServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *hallController) requireAdmin(ctx *gin.Context) bool {
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

func mapHallBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrHallValidation
	}
	return err
}

func mapHallListBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrHallListQuery
	}
	return err
}

func writeHallServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrHallListQuery):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrHallValidation):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrHallCannotActivateInactiveBranch):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrBranchNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("branch not found")))
	case errors.Is(err, service.ErrHallNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("hall not found")))
	case errors.Is(err, service.ErrHallNameConflict):
		ctx.JSON(http.StatusConflict, service.PrepareErrorResponse(service.ErrHallNameConflict))
	default:
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
	}
}
