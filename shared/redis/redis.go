package redis

import (
	"soliant-mock-api/shared/logger"

	goredis "github.com/innovafour/redis"
)

type CacheRepository goredis.CacheRepository

func NewClient(prefix string) (goredis.CacheRepository, error) {
	var redisClient goredis.CacheRepository
	var err error

	createRedisDto := goredis.CreateNewRedisDTO{
		Host:   "localhost:6379",
		Prefix: prefix,
	}

	redisClient, err = goredis.NewClient(createRedisDto)

	if err != nil {
		logger.Error("failed to create redis client", "error", err)

		return nil, err
	}

	return redisClient, nil
}
