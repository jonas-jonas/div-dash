package mocks

import (
	"context"
	"div-dash/internal/db"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) HasLastSuccessfulJobExpired(ctx context.Context, name string, duration time.Duration) (bool, error) {
	args := m.Called(ctx, name, duration)
	return args.Bool(0), args.Error(1)
}

func (m *MockJobService) StartJob(ctx context.Context, name string) (db.StartJobRow, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(db.StartJobRow), args.Error(1)

}
func (m *MockJobService) FinishJob(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockJobService) FailJob(ctx context.Context, id int32, message string) error {
	args := m.Called(ctx, id, message)
	return args.Error(0)
}
