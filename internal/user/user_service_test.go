package user

import (
	"div-dash/internal/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "email@email.de", "password")

	mock.ExpectQuery("^-- name: CreateUser :one .*$").WillReturnRows(rows)

	userService := New()
	createUser := CreateUserParams{
		Email:    "email@email.de",
		Password: "password",
	}

	user, err := userService.CreateUser(createUser)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "email@email.de", user.Email)
	assert.Equal(t, "password", user.PasswordHash)
}

func TestListUsers(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()
	userService := New()
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "email@email.de", "password").
		AddRow(2, "email@email.de", "password").
		AddRow(3, "email@email.de", "password")

	mock.ExpectQuery("^-- name: ListUsers :many .*$").WillReturnRows(rows)

	users, err := userService.ListUsers()

	assert.Nil(t, err)
	assert.Len(t, users, 3)
}

func TestDeleteUser(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()
	userService := New()
	// rows := sqlmock.NewRows([]string{"id", "email", "password"}).
	// 	AddRow(1, "email@email.de", "password")
	mock.ExpectExec("^-- name: DeleteUser :exec .*$").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := userService.DeleteUser(int64(1))

	assert.Nil(t, err)
}

func TestFindByEmail(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()
	userService := New()

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "test@email.email", "password")
	mock.ExpectQuery("-- name: FindByEmail :one").WillReturnRows(rows)

	user, err := userService.FindByEmail("test@email.email")
	assert.Nil(t, err)
	assert.NotNil(t, user)
}

func TestExistsByEmail(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	config.SetDB(sdb)
	defer sdb.Close()
	userService := New()

	mock.ExpectQuery("-- name: CountByEmail :one").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := userService.ExistsByEmail("test@email.email")
	assert.Nil(t, err)
	assert.True(t, exists)
}
