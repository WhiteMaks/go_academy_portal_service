package service

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
)

type MicroserviceService interface {
	GetMicroserviceV1(ctx context.Context) (model.GetMicroserviceV1Response, error)
}

type microserviceService struct {
	userRepository repository.UserRepository
}

func NewMicroserviceService(userRepository repository.UserRepository) MicroserviceService {
	return &microserviceService{
		userRepository: userRepository,
	}
}

func (m microserviceService) GetMicroserviceV1(ctx context.Context) (model.GetMicroserviceV1Response, error) {
	isAdminCreationAllowed, err := m.userRepository.IsAdminCreationAllowedV1(ctx)
	if err != nil {
		return model.GetMicroserviceV1Response{}, err
	}

	response := model.GetMicroserviceV1Response{
		ReadyToUse: !isAdminCreationAllowed.Bool,
	}

	return response, nil
}
