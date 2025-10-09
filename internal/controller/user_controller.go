package controller

import (
	"errors"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type UserController interface {
	CreateUserV1(ctx *gin.Context)
	GenerateUserTokenV1(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{service: service}
}

func (c *userController) CreateUserV1(ctx *gin.Context) {
	var request model.PostUserV1Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
		return
	}

	response, err := c.service.CreateUserV1(ctx, request)
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

func (c *userController) GenerateUserTokenV1(ctx *gin.Context) {
	var request model.PostUserTokenV1Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(err))
		return
	}

	response, err := c.service.GenerateUserTokenV1(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, service.PrepareErrorResponse(errors.New("invalid credentials")))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
