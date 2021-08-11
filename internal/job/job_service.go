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
	nowFunc func() time.Time
}

func New(queries *db.Queries) *JobService {
	return &JobService{
		queries: queries,
		nowFunc: time.Now,
	}
}

func (j *JobService) hasLastSuccessfulJobExpired(ctx context.Context, name string, duration time.Duration) (bool, error) {
	lastJob, err := j.queries.GetLastJobByName(ctx, name)
	if err != nil {
		return true, err
	}
	if lastJob.HadError {
		return true, errors.New(lastJob.ErrorMessage.String)
	}
	var lastJobTimestamp time.Time
	if lastJob.Finished.Valid {
		lastJobTimestamp = time.Unix(lastJob.Finished.Int64, 0)
	} else {
		lastJobTimestamp = time.Unix(lastJob.Started, 0)
	}

	return lastJobTimestamp.Add(duration).Before(j.nowFunc()), nil
}

func (j *JobService) startJob(ctx context.Context, name string) (db.StartJobRow, error) {
	job, err := j.queries.StartJob(ctx, db.StartJobParams{
		Name:    name,
		Started: j.nowFunc().Unix(),
	})
	if err != nil {
		log.Printf("Could not start job '%s': %s", name, err.Error())
		return job, err
	}
	log.Printf("Starting job '%s' with id %d...", name, job.ID)
	return job, err
}

func (j *JobService) finishJob(ctx context.Context, id int32) error {
	job, err := j.queries.FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.nowFunc().Unix(),
			Valid: true,
		},
	})

	if err != nil {
		log.Printf("Could not finish job #%d: %s", id, err.Error())
		return err
	}

	log.Printf("Job %s#%d succeded in '%d'ms", job.Name, job.ID, (job.Finished.Int64 - job.Started))
	return nil
}

func (j *JobService) failJob(ctx context.Context, id int32, message string) error {
	job, err := j.queries.FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.nowFunc().Unix(),
			Valid: true,
		},
		ErrorMessage: sql.NullString{
			String: message,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Could not fail job #%d: %s", id, err.Error())
		return err
	}

	log.Printf("Job %s#%d failed in %ds", job.Name, job.ID, (job.Finished.Int64 - job.Started))
	return nil
}

type JobRunner func(context.Context) error
type JobDefinition struct {
	Key      string
	Validity time.Duration
}

func (j *JobService) RunJob(job JobDefinition, fn JobRunner) {

	ctx := context.Background()
	expired, err := j.hasLastSuccessfulJobExpired(ctx, job.Key, time.Duration(job.Validity))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Could not check for expiration of job %s: %e", job.Key, err)
		return
	}
	if !expired {
		log.Printf("Last execution of job %s is not expired, skipping...", job.Key)
		return
	}

	startedJob, err := j.startJob(ctx, job.Key)
	if err != nil {
		log.Printf("Could not start job %s: %e", job.Key, err)
		return
	}

	err = fn(ctx)

	if err != nil {
		j.failJob(ctx, startedJob.ID, err.Error())
		return
	}

	j.finishJob(ctx, startedJob.ID)
}
