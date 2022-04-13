package job

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/logging"
	"div-dash/internal/timex"
	"time"
)

type (
	JobServiceProvider interface {
		JobService() *JobService
	}

	jobServiceDependencies interface {
		db.QueriesProvider
		logging.LoggerProvider
		timex.TimeHolderProvider
	}
)

type JobService struct {
	jobServiceDependencies
}

func NewJobService(j jobServiceDependencies) *JobService {
	return &JobService{
		jobServiceDependencies: j,
	}
}

func (j *JobService) isJobExpired(job *db.GetLastJobByNameRow, duration time.Duration) bool {
	var timestamp time.Time
	if job.Finished.Valid {
		timestamp = time.Unix(job.Finished.Int64, 0)
	} else {
		timestamp = time.Unix(job.Started, 0)
	}
	return job.HadError || timestamp.Add(duration).Before(j.TimeHolder().GetTime())
}

func (j *JobService) startJob(ctx context.Context, name string) (db.StartJobRow, error) {
	job, err := j.Queries().StartJob(ctx, db.StartJobParams{
		Name:    name,
		Started: j.TimeHolder().GetTime().Unix(),
	})
	if err != nil {
		j.Logger().Warnf("Could not start job '%s': %s", name, err.Error())
		return job, err
	}
	j.Logger().Infof("Starting job '%s' with id %d...", name, job.ID)
	return job, err
}

func (j *JobService) finishJob(ctx context.Context, id int32) error {
	job, err := j.Queries().FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.TimeHolder().GetTime().Unix(),
			Valid: true,
		},
	})

	if err != nil {
		j.Logger().Warnf("Could not finish job #%d: %s", id, err.Error())
		return err
	}

	j.Logger().Infof("Job %s#%d succeded in '%d'ms", job.Name, job.ID, (job.Finished.Int64 - job.Started))
	return nil
}

func (j *JobService) failJob(ctx context.Context, id int32, message string) error {
	job, err := j.Queries().FinishJob(ctx, db.FinishJobParams{
		ID: id,
		Finished: sql.NullInt64{
			Int64: j.TimeHolder().GetTime().Unix(),
			Valid: true,
		},
		ErrorMessage: sql.NullString{
			String: message,
			Valid:  true,
		},
	})
	if err != nil {
		j.Logger().Warnf("Could not fail job #%d: %s", id, err.Error())
		return err
	}

	j.Logger().Warnf("Job %s#%d failed in %ds", job.Name, job.ID, (job.Finished.Int64 - job.Started))
	return nil
}

type JobRunner func(context.Context) error
type JobDefinition struct {
	Key      string
	Validity time.Duration
}

func (j *JobService) RunJob(job JobDefinition, fn JobRunner) {

	ctx := context.Background()
	lastJob, err := j.Queries().GetLastJobByName(ctx, job.Key)
	if err != nil && err != sql.ErrNoRows {
		j.Logger().Warnf("Could not check for expiration of job %s: %e", job.Key, err)
		return
	}
	expired := j.isJobExpired(&lastJob, job.Validity)
	if !lastJob.HadError {
		if !expired {
			j.Logger().Debugf("Last execution of job %s is not expired, skipping...", job.Key)
			return
		}
	} else {
		j.Logger().Debugf("Last execution of job %s had error '%s' retrying", job.Key, lastJob.ErrorMessage.String)
	}

	startedJob, err := j.startJob(ctx, job.Key)
	if err != nil {
		j.Logger().Warnf("Could not start job %s: %e", job.Key, err)
		return
	}

	err = fn(ctx)

	if err != nil {
		j.Logger().Warnf("Job %s#%d failed with error '%s'", job.Key, startedJob.ID, err.Error())
		j.failJob(ctx, startedJob.ID, err.Error())
		return
	}

	j.finishJob(ctx, startedJob.ID)
}
