package job

import (
	"bytes"
	"context"
	"database/sql"
	"div-dash/internal/db"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHasLastSuccessfulJobExpired(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(200, 0) },
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"}).
		AddRow(1, "job-name", 0, 100, "", false)

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*50)
	assert.True(t, expired)
	assert.Empty(t, err)
}

func TestHasLastSuccessfulJobExpiredNotExpiredReturnFalse(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(0, 120) },
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"}).
		AddRow(1, "job-name", 0, 100, "", false)

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*50)
	assert.False(t, expired)
	assert.Empty(t, err)
}

func TestHasLastSuccessfulJobExpiredNoEntriesReturnsTrue(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"})

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*50)
	assert.True(t, expired)
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestHasLastSuccessfulJobExpiredHadErrorReturnsTrue(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"}).
		AddRow(1, "job-name", 0, 100, "test-error", true)

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*50)
	assert.True(t, expired)
	assert.Equal(t, err.Error(), "test-error")
}

func TestHasLastSuccessfulJobExpiredNoFinishedTimestampReturnsTrue(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"}).
		AddRow(1, "job-name", 0, sql.NullInt64{Valid: false}, "", false)

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*50)
	assert.True(t, expired)
	assert.Nil(t, err)
}

func TestHasLastSuccessfulJobExpiredNoFinishedTimestampNoExpiredReturnsFalse(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	rows := sqlmock.NewRows([]string{"id", "name", "started", "finished", "error_message", "had_error"}).
		AddRow(1, "job-name", 0, sql.NullInt64{Valid: false}, "", false)

	mock.ExpectQuery("^-- name: GetLastJobByName :one .*$").WithArgs("job-name").WillReturnRows(rows)

	ctx := context.Background()

	expired, err := jobService.HasLastSuccessfulJobExpired(ctx, "job-name", time.Second*100)
	assert.False(t, expired)
	assert.Nil(t, err)
}

func TestStartJob(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(0, 0) },
	}

	rows := sqlmock.NewRows([]string{"id", "started"}).
		AddRow(1, 0)

	mock.ExpectQuery("^-- name: StartJob :one .*$").WithArgs("test-job", 0).WillReturnRows(rows)

	ctx := context.Background()

	job, err := jobService.StartJob(ctx, "test-job")
	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.Equal(t, "Starting job 'test-job' with id 1...\n", str.String())
}

func TestStartJobDbError(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(0, 0) },
	}

	mock.ExpectQuery("^-- name: StartJob :one .*$").WithArgs("test-job", 0).WillReturnError(errors.New("test-error"))

	ctx := context.Background()

	_, err := jobService.StartJob(ctx, "test-job")
	assert.Equal(t, err.Error(), "test-error")
	assert.Equal(t, "Could not start job 'test-job': test-error\n", str.String())
}

func TestFinishJob(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	rows := sqlmock.NewRows([]string{"name", "id", "started", "finished"}).
		AddRow("test-job", 1, 0, 100)

	mock.ExpectQuery("^-- name: FinishJob :one .*$").WithArgs(100, nil, 1).WillReturnRows(rows)

	ctx := context.Background()

	err := jobService.FinishJob(ctx, 1)
	assert.Nil(t, err)
}

func TestFinishJobDbError(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	mock.ExpectQuery("^-- name: FinishJob :one .*$").WithArgs(100, nil, 1).WillReturnError(errors.New("test-error"))

	ctx := context.Background()

	err := jobService.FinishJob(ctx, 1)
	assert.Equal(t, "test-error", err.Error())
}

func TestFailJob(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	rows := sqlmock.NewRows([]string{"name", "id", "started", "finished"}).
		AddRow("test-job", 1, 0, 100)
	mock.ExpectQuery("^-- name: FinishJob :one .*$").WithArgs(100, "test-error", 1).WillReturnRows(rows)

	ctx := context.Background()

	err := jobService.FailJob(ctx, 1, "test-error")
	assert.Nil(t, err)
}

func TestFailJobDbError(t *testing.T) {
	var str bytes.Buffer
	sdb, mock, _ := sqlmock.New()
	jobService := JobService{
		queries: db.New(sdb),
		logger:  log.New(&str, "", 0),
		nowFunc: func() time.Time { return time.Unix(100, 0) },
	}

	mock.ExpectQuery("^-- name: FinishJob :one .*$").WithArgs(100, "test-error", 1).WillReturnError(errors.New("test-dberror"))

	ctx := context.Background()

	err := jobService.FailJob(ctx, 1, "test-error")
	assert.Equal(t, "test-dberror", err.Error())
}

func TestNew(t *testing.T) {
	queries := &db.Queries{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := New(queries, logger)

	assert.Equal(t, queries, service.queries)
	assert.Equal(t, logger, service.logger)
	assert.NotNil(t, service.nowFunc)
}
