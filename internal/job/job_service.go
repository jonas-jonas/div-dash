package job

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"errors"
	"log"
	"time"
)

type JobService struct {
	queries *db.Queries
	logger  *log.Logger
	nowFunc func() time.Time
}

func New(queries *db.Queries, logger *log.Logger) *JobService {
	return &JobService{
		queries: queries,
		logger:  logger,
		nowFunc: time.Now,
	}
}

func (j *JobService) HasLastSuccessfulJobExpired(ctx context.Context, name string, duration time.Duration) (bool, error) {
	lastJob, err := j.queries.GetLastJobByName(ctx, name)
	if err != nil {
		return true, err
	}
	if lastJob.HadError {
		return true, errors.New(lastJob.ErrorMessage.String)
	}
	var lastJobTimestamp time.Time
	if lastJob.Finished.Valid {
		lastJobTimestamp = time.Unix(0, lastJob.Finished.Int64)
	} else {
		lastJobTimestamp = time.Unix(0, lastJob.Started)
	}

	return lastJobTimestamp.Add(duration).Before(j.nowFunc()), nil
}

func (j *JobService) StartJob(ctx context.Context, name string) (db.StartJobRow, error) {
	job, err := j.queries.StartJob(ctx, db.StartJobParams{
		Name:    name,
		Started: j.nowFunc().UnixNano(),
	})
	if err != nil {
		j.logger.Printf("Could not start job '%s': %s", name, err.Error())
		return job, err
	}
	j.logger.Printf("Starting job '%s' with id %d...", name, job.ID)
	return job, err
}

func (j *JobService) FinishJob(ctx context.Context, id int32) error {
	job, err := j.queries.FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.nowFunc().UnixNano(),
			Valid: true,
		},
	})

	if err != nil {
		j.logger.Printf("Could not finish job #%d: %s", id, err.Error())
		return err
	}

	j.logger.Printf("Job %s#%d succeded in '%d'ms", job.Name, job.ID, (job.Finished.Int64-job.Started)/1000)
	return nil
}

func (j *JobService) FailJob(ctx context.Context, id int32, message string) error {
	job, err := j.queries.FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.nowFunc().UnixNano(),
			Valid: true,
		},
		ErrorMessage: sql.NullString{
			String: message,
			Valid:  true,
		},
	})
	if err != nil {
		j.logger.Printf("Could not fail job #%d: %s", id, err.Error())
		return err
	}

	j.logger.Printf("Job %s#%d failed in %dms", job.Name, job.ID, (job.Finished.Int64-job.Started)/1000)
	return nil
}
