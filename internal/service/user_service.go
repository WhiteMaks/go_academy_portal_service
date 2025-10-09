package service

import (
	"context"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
)

type UserService interface {
	CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	userRecordParams, err := PrepareCreateUserV1RecordParams(request)
	if err != nil {
		return model.PostUserV1Response{}, err
	}

	userRecord, err := s.userRepository.CreateUserV1(ctx, userRecordParams)
	if err != nil {
		return model.PostUserV1Response{}, err
	}

	response := model.PostUserV1Response{
		ID: userRecord.ID,
	}

	return response, nil
}
