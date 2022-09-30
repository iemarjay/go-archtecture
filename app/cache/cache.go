package cache

import (
	"archtecture/app/env"
	"github.com/go-redis/redis"
	"time"
)

type Config struct {
	RedisUrl      string
	RedisPassword string
}

func NewConfig(env *env.Env) *Config {
	return &Config{
		RedisUrl:      env.RedisUrl,
		RedisPassword: env.RedisPassword,
	}
}

type Cache struct {
	redisClient *redis.Client
}

func NewCacheWithRedis(rc *redis.Client) *Cache {
	return &Cache{redisClient: rc}
}

func NewCacheWithRedisFromConfig(config *Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
	})

	return NewCacheWithRedis(client)
}

func (c *Cache) Set(key string, value string, expiry time.Duration) error {
	err := c.redisClient.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *Cache) Forget(key string) error {
	_, err := c.redisClient.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetOrDefault(key string, fail string) string {
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		return fail
	}

	return value
}
