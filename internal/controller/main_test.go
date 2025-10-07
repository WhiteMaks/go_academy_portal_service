package controller

import (
	"context"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

type mockUserService struct {
	createUserV1Func func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
}

func (m *mockUserService) CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	return m.createUserV1Func(ctx, request)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
