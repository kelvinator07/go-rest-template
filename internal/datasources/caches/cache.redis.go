package caches

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Del(key string) error
}

type redisCache struct {
	host     string
	db       int
	password string
	expires  time.Duration
	client   *redis.Client
}

func NewRedisCache(host string, db int, password string, expires time.Duration) (RedisCache, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	redisCache := &redisCache{
		host:     host,
		db:       db,
		password: password,
		expires:  expires,
		client:   redisClient,
	}

	return redisCache, nil
}

func (cache *redisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cache.client.Set(cache.client.Context(), key, json, cache.expires*time.Minute).Err()
}

func (cache *redisCache) Get(key string) (email string, err error) {
	val, err := cache.client.Get(cache.client.Context(), key).Result()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(val), &email)
	return email, err
}

func (cache *redisCache) Del(key string) error {
	return cache.client.Del(cache.client.Context(), key).Err()
}
