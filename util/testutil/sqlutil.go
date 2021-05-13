package testutil

import (
	"database/sql/driver"
	"time"
)

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyTransactionId struct{}

func (a AnyTransactionId) Match(v driver.Value) bool {
	_, ok := v.(string)
	//TODO
	return ok
}

type AnyAccountId struct{}

func (a AnyAccountId) Match(v driver.Value) bool {
	_, ok := v.(string)
	//TODO
	return ok
}
