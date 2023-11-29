package featuretoggle

import (
	"fmt"

	"github.com/delivery-much/dm-go/logger"
	"github.com/go-redis/redis"
)

type redisClient interface {
	subscribe(pattern string) (subs *redis.PubSub)
	hgetall(namespace string) (map[string]string, error)
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

	configRes := client.ConfigSet("notify-keyspace-events", "KEA")
	if configRes == nil || configRes.Err() != nil {
		logger.NoCTX().Errorf(
			"Failed to configure feature toggle redis client to notify changes: %s. The library will be initiated anyway",
			configRes.Err().Error(),
		)
	}

	_, err = client.Ping().Result()
	if err != nil {
		err = fmt.Errorf("Failed to connect to redis feature toggle with message: %s", err.Error())
	}

	rc = &redisDB{client}
	return
}

// subscribes to a given channel pattern, to receive messages from.
// returns the subscriber
func (db *redisDB) subscribe(pattern string) (subs *redis.PubSub) {
	return db.Client.PSubscribe(pattern)
}

// hgetall gets all the redis keys and values for a given namespace
func (db *redisDB) hgetall(namespace string) (m map[string]string, err error) {
	resp := db.Client.HGetAll(namespace)
	if resp == nil || resp.Err() != nil {
		err = fmt.Errorf("Failed to get the redis pairs for namespace %s: %s", namespace, resp.Err().Error())
	}

	m = resp.Val()
	return
}
