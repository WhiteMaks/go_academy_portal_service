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

type BranchController interface {
	CreateBranchV1(ctx *gin.Context)
	GetBranchesV1(ctx *gin.Context)
	GetBranchV1(ctx *gin.Context)
	UpdateBranchV1(ctx *gin.Context)
}

type branchController struct {
	branchService service.BranchService
	userService   service.UserService
}

func NewBranchController(branchService service.BranchService, userService service.UserService) BranchController {
	return &branchController{
		branchService: branchService,
		userService:   userService,
	}
}

func (c *branchController) CreateBranchV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var request model.PostBranchV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapBranchBindError(err)))
		return
	}

	response, err := c.branchService.CreateBranchV1(ctx.Request.Context(), request)
	if err != nil {
		writeBranchServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *branchController) GetBranchesV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	var query model.GetBranchesV1Query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapBranchListBindError(err)))
		return
	}

	response, err := c.branchService.GetBranchesV1(ctx.Request.Context(), query)
	if err != nil {
		writeBranchServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *branchController) GetBranchV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parsePositiveID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid branch id")))
		return
	}

	response, err := c.branchService.GetBranchV1(ctx.Request.Context(), id)
	if err != nil {
		writeBranchServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *branchController) UpdateBranchV1(ctx *gin.Context) {
	if !c.requireAdmin(ctx) {
		return
	}

	id, err := parsePositiveID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid branch id")))
		return
	}

	var request model.PutBranchV1Request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(mapBranchBindError(err)))
		return
	}

	response, err := c.branchService.UpdateBranchV1(ctx.Request.Context(), id, request)
	if err != nil {
		writeBranchServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *branchController) requireAdmin(ctx *gin.Context) bool {
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

func mapBranchBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrBranchValidation
	}
	return err
}

func mapBranchListBindError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return service.ErrBranchListQuery
	}
	return err
}

func writeBranchServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrBranchListQuery):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrBranchValidation):
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
	case errors.Is(err, service.ErrBranchNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("branch not found")))
	case errors.Is(err, service.ErrBranchNameConflict):
		ctx.JSON(http.StatusConflict, service.PrepareErrorResponse(service.ErrBranchNameConflict))
	default:
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
	}
}

func parsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
