package cache

import "time"

type Cache interface {
	Set(string, string, time.Duration) error
	Get(string) (string, error)
	Forget(string) error
}
