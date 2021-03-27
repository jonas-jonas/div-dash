package user

import (
	"context"
	"div-dash/internal/config"
	"div-dash/internal/db"
)

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

func (u *UserService) CreateUser(user db.CreateUserParams) (db.User, error) {
	ctx := context.Background()

	return config.Queries().CreateUser(ctx, user)
}

func (u *UserService) ListUsers() ([]db.User, error) {
	ctx := context.Background()

	return config.Queries().ListUsers(ctx)
}

func (u *UserService) DeleteUser(userId int64) error {
	ctx := context.Background()

	return config.Queries().DeleteUser(ctx, userId)
}

func (u *UserService) FindByEmail(email string) (db.User, error) {
	ctx := context.Background()

	return config.Queries().FindByEmail(ctx, email)
}

func (u *UserService) ExistsByEmail(email string) (bool, error) {
	ctx := context.Background()

	count, err := config.Queries().CountByEmail(ctx, email)

	return count > 0, err
}

func (u *UserService) FindById(userId int64) (db.User, error) {

	ctx := context.Background()

	return config.Queries().GetUser(ctx, userId)
}
