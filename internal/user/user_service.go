package user

import (
	"context"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/util/security"
)

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

type CreateUserParams struct {
	Email    string
	Password string
}

func (u *UserService) CreateUser(params CreateUserParams) (db.User, error) {
	ctx := context.Background()

	passwordHash, err := security.HashPassword(params.Password)

	if err != nil {
		return db.User{}, err
	}

	user := db.CreateUserParams{
		Email:        params.Email,
		PasswordHash: passwordHash,
	}

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
