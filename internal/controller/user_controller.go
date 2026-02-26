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
	"github.com/lib/pq"
)

type UserController interface {
	GetUserV1(ctx *gin.Context)
	CreateUserV1(ctx *gin.Context)
	CreateAdminV1(ctx *gin.Context)
	GenerateUserTokenV1(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{userService: service}
}

func (c *userController) GetUserV1(ctx *gin.Context) {
	payload := ctx.MustGet(middleware.AuthorizationPayload).(*util.Payload)

	response, err := c.userService.GetUserV1(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *userController) CreateUserV1(ctx *gin.Context) {
	var request model.PostUserV1Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
		return
	}

	payload := ctx.MustGet(middleware.AuthorizationPayload).(*util.Payload)
	hasAccess := c.userService.HasAccess(
		ctx,
		payload,
		[]database.RootUserRole{
			database.RootUserRoleAdmin,
			database.RootUserRoleCoach,
		},
	)
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return
	}

	response, err := c.userService.CreateUserV1(ctx, request)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == database.PgErrUniqueViolation {
			ctx.JSON(http.StatusConflict, service.PrepareErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) CreateAdminV1(ctx *gin.Context) {
	var request model.PostUserV1Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
		return
	}

	response, err := c.userService.CreateAdminV1(ctx, request)
	if err != nil {
		if errors.Is(err, service.ErrAdminCreationIsNotAllowed) {
			ctx.JSON(http.StatusNotFound, service.PrepareEmptyResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) GenerateUserTokenV1(ctx *gin.Context) {
	var request model.PostUserTokenV1Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
		return
	}

	response, err := c.userService.GenerateUserTokenV1(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, service.PrepareErrorResponse(errors.New("invalid credentials")))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
