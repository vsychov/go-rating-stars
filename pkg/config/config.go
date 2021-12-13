package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config is application configuration
type Config struct {
	PgsqlAddr     string `envconfig:"PGSQL_ADDR"`
	PgsqlUser     string `envconfig:"PGSQL_USER"`
	PgsqlPassword string `envconfig:"PGSQL_PASSWORD"`
	PgsqlDBName   string `envconfig:"PGSQL_DBNAME"`
	PgsqlPort     uint16 `envconfig:"PGSQL_PORT" default:"5432"`
	PgsqlSSLMode  string `envconfig:"PGSQL_SSLMODE" default:"disable"`
	PgsqlTimezone string `envconfig:"PGSQL_TIMEZONE" default:"Europe/London"`

	RedisAddr     string `envconfig:"REDIS_ADDR"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	StorageType    string `envconfig:"STORAGE_TYPE" default:"pgsql"`
	ClientIpHeader string `envconfig:"CLIENT_IP_HEADER" default:"X-Real-IP"`
}

// CreateFromEnv create new config instance from env variables
func CreateFromEnv() (conf Config, err error) {
	err = envconfig.Process("", &conf)
	if err != nil {
		return
	}

	return
}
