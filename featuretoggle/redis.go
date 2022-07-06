package featuretoggle

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type redisClient interface {
	get(key string) (string, error)
}

// RedisDB represents the Redis database client
type redisDB struct {
	*redis.Client
}

// getRedisClient starts the connection with Redis
func getRedisClient(host string, port string, db int) (rc redisClient, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       db,
	})

	rc = &redisDB{client}

	_, err = client.Ping().Result()
	if err != nil {
		err = fmt.Errorf("Failed to connect to redis feature toggle with message: %s", err.Error())
	}
	return
}

// Get get a result from a key
func (c *redisDB) get(key string) (string, error) {
	return c.Client.Get(key).Result()
}

// redisDBMock represents the Redis database client mock
type redisDBMock struct {
	throwErr  bool
	getResult string
}

// setGetResult sets the string that the mock should return when calling get
func (c *redisDBMock) setGetResult(res string) {
	c.getResult = res
}

func (c *redisDBMock) get(key string) (string, error) {
	if c.throwErr {
		return "", errors.New("Default error")
	}

	return c.getResult, nil
}
