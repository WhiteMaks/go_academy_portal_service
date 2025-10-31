package service

import (
	"errors"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/util"
)

func PrepareEmptyResponse() model.EmptyResponse {
	return model.EmptyResponse{}
}

func PrepareForbiddenErrorResponse() model.ErrorResponse {
	return PrepareErrorResponse(errors.New("forbidden"))
}

func PrepareErrorResponse(err error) model.ErrorResponse {
	return model.ErrorResponse{
		Message: err.Error(),
	}
}

func PrepareCreateUserV1RecordParams(request model.PostUserV1Request, role database.RootUserRole, isActive bool) (database.CreateUserV1Params, error) {
	hashPassword, err := util.HashPassword(request.Password)
	if err != nil {
		return database.CreateUserV1Params{}, err
	}

	userRecordParams := database.CreateUserV1Params{
		PUsername: request.Username,
		PPassword: hashPassword,
		PRole:     role,
		PIsActive: isActive,
	}

	return userRecordParams, nil
}
