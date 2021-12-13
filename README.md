[![Build Status](https://github.com/vsychov/go-rating-stars/actions/workflows/ci.yml/badge.svg)](https://github.com/vsychov/go-rating-stars/actions)
[![codecov](https://codecov.io/gh/vsychov/go-rating-stars/branch/master/graph/badge.svg?token=7V853A3LYA)](https://codecov.io/gh/vsychov/go-rating-stars)
[![Go Reference](https://pkg.go.dev/badge/github.com/vsychov/go-rating-stars.svg)](https://pkg.go.dev/github.com/vsychov/go-rating-stars)
[![Go Report Card](https://goreportcard.com/badge/github.com/vsychov/go-rating-stars)](https://goreportcard.com/report/github.com/vsychov/go-rating-stars)
[![Docker Repository on Quay](https://quay.io/repository/vsychov/go-rating-stars/status "Docker Repository on Quay")](https://quay.io/repository/vsychov/go-rating-stars)
---
#Load env
execute bash command
```bash
set -o allexport; source .env; set +o allexport
```

#Health check
```bash
curl -iL -XGET -H "X-Health-Check: 1" http://localhost:8080
```


#Storages
GoVote can use Postgress or Redis for store data.

#Generate mocks
```bash
mockery --output tests/mocks --recursive --name=StorageInterface
```

#Env variables
PGSQL_ADDR=pgsql
PGSQL_USER=gorm
PGSQL_PASSWORD=gorm
PGSQL_DBNAME=gorm
PGSQL_PORT=9920
PGSQL_SSLMODE=disable
PGSQL_TIMEZONE=Europe/London

REDIS_ADDR=
REDIS_DB=
REDIS_PASSWORD=

STORAGE_TYPE=pgsql or redis (default - pgsql), redis not fully implemented yet

CLIENT_IP_HEADER="X-Real-Ip" #gin TrustedPlatform