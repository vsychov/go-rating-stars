package storage

import (
	"fmt"
	"github.com/vsychov/go-rating-stars/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// PostgresStorage pgsql instance
type PostgresStorage struct {
	client *gorm.DB
}

// CreatePostgres create new instance of PostgresStorage
func CreatePostgres(config config.Config) (storage StorageInterface, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.PgsqlAddr,
		config.PgsqlUser,
		config.PgsqlPassword,
		config.PgsqlDBName,
		config.PgsqlPort,
		config.PgsqlSSLMode,
		config.PgsqlTimezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&KeyValue{})
	if err != nil {
		return
	}

	return &PostgresStorage{
		client: db,
	}, nil
}

func (store *PostgresStorage) SetNX(key string, value float64, ttl time.Duration) (err error) {
	keyValue := KeyValue{
		ID:    key,
		Value: value,
		TTL:   time.Now().Unix() + int64(ttl.Seconds()),
	}

	dbResult := store.client.Create(&keyValue)
	err = dbResult.Error

	return
}

func (store *PostgresStorage) Get(key string) (value float64, err error) {
	keyValue := KeyValue{}
	result := store.client.First(&keyValue, "id = ?", key)

	return keyValue.Value, result.Error
}

func (store *PostgresStorage) Incr(key string, value float64) (result float64, err error) {
	keyValue := KeyValue{
		ID:    key,
		Value: value,
		TTL:   0,
	}

	dbResult := store.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"value": gorm.Expr("key_values.value + ?", keyValue.Value)}),
	}).Create(&keyValue)

	if dbResult.Error != nil {
		err = dbResult.Error
		return
	}

	return store.Get(key)
}

func (store *PostgresStorage) Cron() error {
	result := store.client.Where("ttl > 0 and ttl <= extract(epoch from now())").Delete(KeyValue{})

	return result.Error
}
