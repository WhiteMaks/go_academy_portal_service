package repository

import (
	"context"
	"database/sql"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type UserRepository interface {
	CreateUserV1(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error)
	GetUserByUsernameV1(ctx context.Context, username string) (database.RootUser, error)
	IsAdminCreationAllowedV1(ctx context.Context) (sql.NullBool, error)
}

type userRepository struct {
	store database.Store
}

func NewUserRepository(store database.Store) UserRepository {
	return &userRepository{store: store}
}

func (r *userRepository) CreateUserV1(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
	return r.store.CreateUserV1(ctx, params)
}

func (r *userRepository) GetUserByUsernameV1(ctx context.Context, username string) (database.RootUser, error) {
	return r.store.GetUserByUsernameV1(ctx, username)
}

func (r *userRepository) IsAdminCreationAllowedV1(ctx context.Context) (sql.NullBool, error) {
	return r.store.IsAdminCreationAllowedV1(ctx)
}
