package triecache

import (
	"time"
)

type Cache interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (interface{}, error)
	Delete(string) error
	Keys(string) ([]string, error)
	Expire(string, time.Duration) error
	TTL(string) (int64, error)
	GetInt64(string) (int64, error)
	GetFloat64(string) (float64, error)
	Incr(string, time.Duration) (int64, error)
	IncrBy(string, int64, time.Duration) (int64, error)
	Decr(string, time.Duration) (int64, error)
	DecrBy(string, int64, time.Duration) (int64, error)
}
