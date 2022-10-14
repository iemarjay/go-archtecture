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

type Redis struct {
	client *redis.Client
}

func NewCacheWithRedis(rc *redis.Client) *Redis {
	return &Redis{client: rc}
}

func NewCacheWithRedisFromConfig(config *Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
	})

	return NewCacheWithRedis(client)
}

func (c *Redis) Set(key string, value string, expiry time.Duration) error {
	err := c.client.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) Get(key string) (string, error) {
	value, err := c.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *Redis) Forget(key string) error {
	_, err := c.client.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) GetOrDefault(key string, fail string) string {
	value, err := c.client.Get(key).Result()
	if err != nil {
		return fail
	}

	return value
}
