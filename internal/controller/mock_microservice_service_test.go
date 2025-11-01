package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
)

type mockMicroserviceService struct {
	getMicroserviceV1Func func(ctx context.Context) (model.GetMicroserviceV1Response, error)
}

func (m *mockMicroserviceService) GetMicroserviceV1(ctx context.Context) (model.GetMicroserviceV1Response, error) {
	return m.getMicroserviceV1Func(ctx)
}
