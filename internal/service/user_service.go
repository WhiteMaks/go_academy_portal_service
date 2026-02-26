package service

import (
	"context"
	"errors"
	"slices"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/util"
)

var (
	ErrAdminCreationIsNotAllowed = errors.New("admin creation is not allowed")
)

type UserService interface {
	GetUserV1(ctx context.Context, username string) (model.GetUserV1Response, error)
	CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
	CreateAdminV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
	GenerateUserTokenV1(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error)
	HasAccess(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool
}

type userService struct {
	userRepository repository.UserRepository
	tokenMaker     util.TokenMaker
}

func NewUserService(userRepository repository.UserRepository, tokenMaker util.TokenMaker) UserService {
	return &userService{
		userRepository: userRepository,
		tokenMaker:     tokenMaker,
	}
}

func (s *userService) GetUserV1(ctx context.Context, username string) (model.GetUserV1Response, error) {
	userRecord, err := s.userRepository.GetUserByUsernameV1(ctx, username)
	if err != nil {
		return model.GetUserV1Response{}, err
	}

	response := model.GetUserV1Response{
		ID:       userRecord.ID,
		Username: userRecord.Username,
		Role:     string(userRecord.Role),
	}

	return response, nil
}

func (s *userService) CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	userRecordParams, err := PrepareCreateUserV1RecordParams(request, database.RootUserRoleAthlete, true)
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

func (s *userService) CreateAdminV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	isAdminCreationAllowed, err := s.userRepository.IsAdminCreationAllowedV1(ctx)
	if err != nil {
		return model.PostUserV1Response{}, err
	}

	if !isAdminCreationAllowed.Valid || !isAdminCreationAllowed.Bool {
		return model.PostUserV1Response{}, ErrAdminCreationIsNotAllowed
	}

	userRecordParams, err := PrepareCreateUserV1RecordParams(request, database.RootUserRoleAdmin, true)
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

func (s *userService) GenerateUserTokenV1(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error) {
	userRecord, err := s.userRepository.GetUserByUsernameV1(ctx, request.Username)
	if err != nil {
		return model.PostUserTokenV1Response{}, err
	}

	err = util.CheckPassword(request.Password, userRecord.Password)
	if err != nil {
		return model.PostUserTokenV1Response{}, err
	}

	token, err := s.tokenMaker.CreateToken(userRecord.Username, string(userRecord.Role))
	if err != nil {
		return model.PostUserTokenV1Response{}, err
	}

	response := model.PostUserTokenV1Response{
		Token: token,
	}

	return response, nil
}

func (s *userService) HasAccess(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
	userRecord, err := s.userRepository.GetUserByUsernameV1(ctx, payload.Username)
	if err != nil {
		return false
	}

	if !slices.Contains(roles, userRecord.Role) {
		return false
	}

	return true
}
