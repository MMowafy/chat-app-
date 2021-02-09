package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"instabug-task/api/application"
	"time"
)

type RedisService struct {
	client redis.Cmdable
}

func NewRedisService() *RedisService {
	client, _ := application.GetRedisConnectionByName("redis")
	return &RedisService{
		client,
	}
}

func (redisService *RedisService) checkRedisConnection() error {
	if redisService.client == nil {
		return errors.New("Failed to find openn redis connection")
	}
	return nil
}

func (redisService *RedisService) Set(key string, value interface{}, duration time.Duration) error {
	err := redisService.checkRedisConnection()

	if err != nil {
		return err
	}

	jsonValue, marshalError := json.Marshal(value)

	if marshalError != nil {
		return fmt.Errorf(
			"Error while marshaling key %s value %+v into json ERROR %s", key, value, marshalError.Error(),
		)
	}

	return redisService.client.Set(context.Background(), key, jsonValue, duration).Err()
}

func (redisService *RedisService) Get(key string) ([]byte, error) {
	err := redisService.checkRedisConnection()
	if err != nil {
		application.GetLogger().Errorf("can not fined opened redis connection")
		return nil, err
	}
	data, err := redisService.client.Get(context.Background(), key).Bytes()
	if err != nil {
		application.GetLogger().Errorf("can not find app key %s with error %s", key, err.Error())
		return nil, err
	}
	return data, nil
}

func (redisService *RedisService) Incr(key string) (int64, error) {
	err := redisService.checkRedisConnection()

	if err != nil {
		application.GetLogger().Errorf("Failed to increment chat in redis for key %s", key)
		return 0,err
	}

	return redisService.client.Incr(context.Background(), key).Result()
}
