package timex

import "time"

type (
	TimeHolderProvider interface {
		TimeHolder() TimeHolder
	}

	TimeHolder interface {
		GetTime() time.Time
	}

	timeHolderImpl struct {
	}
)

func NewTimeProvider() TimeHolder {
	return &timeHolderImpl{}
}

func (t *timeHolderImpl) GetTime() time.Time {
	return time.Now()
}
