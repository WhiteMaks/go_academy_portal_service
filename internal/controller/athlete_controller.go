package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AthleteController interface {
	GetAthletesV1(ctx *gin.Context)
	CreateAthleteV1(ctx *gin.Context)
	GetAthleteByIDV1(ctx *gin.Context)
	UpdateAthleteV1(ctx *gin.Context)
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

func (c *athleteController) GetAthleteByIDV1(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return
	}
	id, ok := parseAthleteID(ctx.Param("id"))
	if !ok {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid athlete id")))
		return
	}
	response, err := c.athleteService.GetAthleteByIDV1(ctx, id)
	if err != nil {
		c.writeAthleteServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *athleteController) UpdateAthleteV1(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, service.PrepareForbiddenErrorResponse())
		return
	}
	id, ok := parseAthleteID(ctx.Param("id"))
	if !ok {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid athlete id")))
		return
	}
	var request model.PutAthleteV1Request
	if err := decodeStrictJSON(ctx, &request); err != nil || validateAthleteUpdateRequest(request) != nil {
		ctx.JSON(http.StatusBadRequest, service.PrepareErrorResponse(errors.New("invalid request body")))
		return
	}
	response, err := c.athleteService.UpdateAthleteV1(ctx, id, request)
	if err != nil {
		c.writeAthleteServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func decodeStrictJSON(ctx *gin.Context, destination any) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return binding.Validator.ValidateStruct(destination)
}

func validateAthleteUpdateRequest(request model.PutAthleteV1Request) error {
	if !isUTCMidnight(request.BirthDate) {
		return errors.New("invalid request body")
	}
	for _, name := range []string{request.FirstName, request.SecondName} {
		if !isValidAthleteName(name, true) {
			return errors.New("invalid request body")
		}
	}
	for _, name := range []*string{request.MiddleName, request.ParentFirstName, request.ParentSecondName, request.ParentMiddleName} {
		if name != nil && !isValidAthleteName(*name, false) {
			return errors.New("invalid request body")
		}
	}
	return nil
}

func isValidAthleteName(value string, required bool) bool {
	trimmed := strings.TrimSpace(value)
	if required && trimmed == "" {
		return false
	}
	return utf8.RuneCountInString(trimmed) <= 100
}

func isUTCMidnight(value time.Time) bool {
	_, offset := value.Zone()
	return !value.IsZero() && value.Year() >= 1 && offset == 0 &&
		value.Hour() == 0 && value.Minute() == 0 && value.Second() == 0 && value.Nanosecond() == 0
}

func parseAthleteID(value string) (int64, bool) {
	if value == "" || value[0] < '1' || value[0] > '9' {
		return 0, false
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return 0, false
		}
	}
	id, err := strconv.ParseInt(value, 10, 64)
	return id, err == nil && id > 0
}

func (c *athleteController) writeAthleteServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrAthleteNotFound):
		ctx.JSON(http.StatusNotFound, service.PrepareErrorResponse(errors.New("athlete not found")))
	default:
		ctx.JSON(http.StatusInternalServerError, service.PrepareErrorResponse(errors.New("internal server error")))
	}
}

func (c *athleteController) isAdmin(ctx *gin.Context) bool {
	payload := ctx.MustGet(middleware.AuthorizationPayload).(*util.Payload)
	return c.userService.HasAccess(ctx, payload, []database.RootUserRole{database.RootUserRoleAdmin})
}
