package service

import "github.com/WhiteMaks/go_academy_portal_service/internal/model"

func PrepareErrorResponse(err error) model.ErrorResponse {
	return model.ErrorResponse{
		Message: err.Error(),
	}
}
