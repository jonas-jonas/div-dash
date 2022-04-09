package job

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"time"

	"go.uber.org/zap"
)

type JobService struct {
	queries *db.Queries
	nowFunc func() time.Time
	logger  *zap.SugaredLogger
}

func NewJobService(queries *db.Queries, logger *zap.Logger) *JobService {
	return &JobService{
		queries: queries,
		nowFunc: time.Now,
		logger:  logger.Sugar(),
	}
}

func (j *JobService) isJobExpired(job *db.GetLastJobByNameRow, duration time.Duration) bool {
	var timestamp time.Time
	if job.Finished.Valid {
		timestamp = time.Unix(job.Finished.Int64, 0)
	} else {
		timestamp = time.Unix(job.Started, 0)
	}
	return job.HadError || timestamp.Add(duration).Before(j.nowFunc())
}

func (j *JobService) startJob(ctx context.Context, name string) (db.StartJobRow, error) {
	job, err := j.queries.StartJob(ctx, db.StartJobParams{
		Name:    name,
		Started: j.nowFunc().Unix(),
	})
	if err != nil {
		j.logger.Warnf("Could not start job '%s': %s", name, err.Error())
		return job, err
	}
	j.logger.Infof("Starting job '%s' with id %d...", name, job.ID)
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
		j.logger.Warnf("Could not finish job #%d: %s", id, err.Error())
		return err
	}

	j.logger.Infof("Job %s#%d succeded in '%d'ms", job.Name, job.ID, (job.Finished.Int64 - job.Started))
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
		j.logger.Warnf("Could not fail job #%d: %s", id, err.Error())
		return err
	}

	j.logger.Warnf("Job %s#%d failed in %ds", job.Name, job.ID, (job.Finished.Int64 - job.Started))
	return nil
}

type JobRunner func(context.Context) error
type JobDefinition struct {
	Key      string
	Validity time.Duration
}

func (j *JobService) RunJob(job JobDefinition, fn JobRunner) {

	ctx := context.Background()
	lastJob, err := j.queries.GetLastJobByName(ctx, job.Key)
	if err != nil && err != sql.ErrNoRows {
		j.logger.Warnf("Could not check for expiration of job %s: %e", job.Key, err)
		return
	}
	expired := j.isJobExpired(&lastJob, job.Validity)
	if !lastJob.HadError {
		if !expired {
			j.logger.Debugf("Last execution of job %s is not expired, skipping...", job.Key)
			return
		}
	} else {
		j.logger.Debugf("Last execution of job %s had error '%s' retrying", job.Key, lastJob.ErrorMessage.String)
	}

	startedJob, err := j.startJob(ctx, job.Key)
	if err != nil {
		j.logger.Warnf("Could not start job %s: %e", job.Key, err)
		return
	}

	err = fn(ctx)

	if err != nil {
		j.logger.Warnf("Job %s#%d failed with error '%s'", job.Key, startedJob.ID, err.Error())
		j.failJob(ctx, startedJob.ID, err.Error())
		return
	}

	j.finishJob(ctx, startedJob.ID)
}
