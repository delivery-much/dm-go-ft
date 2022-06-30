package infra

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient is a wrapper interface
type RedisClient interface {
	Get(string) (string, error)
	Set(string, string, time.Duration) error
	Del(string) error
	Count(string, time.Duration) (int64, error)
}

// RedisDB client of Redis
type RedisDB struct {
	*redis.Client
}

// GetRedisClient get connection with Redis
func GetRedisClient(host string, port string, db int) (redisClient *RedisDB, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       db,
	})

	redisClient = &RedisDB{client}

	_, err = client.Ping().Result()
	if err != nil {
		err = fmt.Errorf("Failed to connect to redis feature toggle with message: %s", err.Error())
	}
	return
}

// Get get a result from a key
func (c *RedisDB) Get(key string) (string, error) {
	return c.Client.Get(key).Result()
}
