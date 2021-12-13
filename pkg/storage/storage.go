package storage

import (
	"go.uber.org/fx"
	"time"
)

type StorageInterface interface {
	// Cron method for execute storage-specific cron-jobs, e.g. cleanup ttl
	Cron() error
	// Get return record from storage, or error
	Get(key string) (float64, error)
	// SetNX create record if not exists
	SetNX(key string, value float64, ttl time.Duration) error
	// Incr atomic increment or decrement
	Incr(key string, value float64) (float64, error) //float64 precision is enough for votes
}

// Cron entrypoint
func Cron(storage StorageInterface, shutdowner fx.Shutdowner) error {
	err := storage.Cron()
	if err != nil {
		return err
	}

	return shutdowner.Shutdown()
}
