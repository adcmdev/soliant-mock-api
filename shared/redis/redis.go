package redis

import (
	"os"
	"soliant-mock-api/shared/logger"

	goredis "github.com/innovafour/redis"
)

type CacheRepository goredis.CacheRepository

func NewClient(prefix string) (goredis.CacheRepository, error) {
	var redisClient goredis.CacheRepository
	var err error

	host := "localhost:6379"

	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
	}

	createRedisDto := goredis.CreateNewRedisDTO{
		Host:   host,
		Prefix: prefix,
	}

	redisClient, err = goredis.NewClient(createRedisDto)

	if err != nil {
		logger.Error("failed to create redis client", "error", err)

		return nil, err
	}

	return redisClient, nil
}
