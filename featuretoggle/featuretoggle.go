package featuretoggle

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/delivery-much/dm-go/logger"
	"github.com/go-redis/redis"
)

var (
	// the redis client connection
	client redisClient
	// represents all of the service feature toggles (key-value pairs) saved in memory
	localMemory map[string]string
	// the name of the service currently being used
	serviceName string
)

// Mock mocks the feature toggle library, will use the keys provided as a param when acessing the feature toggles.
// Used for testing
func Mock(keys map[string]string) {
	localMemory = keys
}

// Reset resets the mock library to its empty state.
func Reset() {
	client = nil
	localMemory = map[string]string{}
	serviceName = ""
}

// Init inits the feature toggle library
func Init(c Config) error {
	cl, err := getRedisClient(c.Host, c.Port, c.DB)
	if err != nil {
		return err
	}
	client = cl
	serviceName = c.ServiceName

	// subscribe to the feature toggle channel and wait for changes
	channelPattern := fmt.Sprintf("__keyspace@%d__:*", c.DB)
	sub := cl.subscribe(channelPattern)
	if sub == nil {
		return fmt.Errorf("Failed to subscribe to feature toggle channel")
	}
	go waitForUpdates(sub)

	err = buildCache()
	if err != nil {
		return err
	}

	logger.NoCTX().Infof("Redis feature toggle started for service %s", c.ServiceName)
	return nil
}

// waitForUpdates receives a redis message subscriber,
// and waits forever for any updates on the redis cache for the specified service.
// When an update is received, rebuilds de cache.
func waitForUpdates(sub *redis.PubSub) {
	ch := sub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				logger.NoCTX().Info("Received a message via the feature toggle redis subscriber, but it was empty")
				return
			}

			separatedChannelName := strings.Split(msg.Channel, ":")
			if len(separatedChannelName) < 2 {
				logger.NoCTX().Infof(
					"Failed to process the feature toggle update message, the channel name was in a unexpected format (%s)",
					msg.Channel,
				)
				return
			}

			channelID := separatedChannelName[1]
			if msg.Payload == "hset" && channelID == serviceName {
				err := buildCache()
				if err != nil {
					logger.NoCTX().Infof("Failed to rebuild feature toggle redis with message: %s", err.Error())
					return
				}
				logger.NoCTX().Infof("Redis feature toggle rebuilt for %s", serviceName)
			}
		}
	}
}

// buildCache gets all the feature toggles for the specified service,
// and saves it to the local memory, so that the toggles can be accessed faster.
func buildCache() error {
	toggles, err := client.hgetall(serviceName)
	if err != nil {
		return fmt.Errorf("Failed to get toggles for service %s: %s", serviceName, err.Error())
	}

	localMemory = toggles
	return nil
}

// IsEnabled checks if given feature key is enabled in redis DB.
//
// returns the default value if:
//
// - the library was not initiated;
//
// - the key was not found;
//
// - the key value is empty;
//
// - the key value is not a boolean.
func IsEnabled(key string, defaultVal bool) (b bool) {
	if localMemory == nil {
		logger.NoCTX().Infof("IsEnabled for key %s, the library was not initiated", key)
		return defaultVal
	}

	val, ok := localMemory[key]
	if !ok || strings.TrimSpace(val) == "" {
		logger.NoCTX().Infof("IsEnabled for key %s, the value was not found or empty", key)
		return defaultVal
	}

	typeKey := fmt.Sprintf("%s.type", key)
	t, ok := localMemory[typeKey]
	if !ok || strings.TrimSpace(t) == "" {
		logger.NoCTX().Infof("IsEnabled for key %s, the value type was not found or empty", key)
		return defaultVal
	}

	b, err := strconv.ParseBool(val)
	if err != nil || t != "boolean" {
		logger.NoCTX().Infof("IsEnabled for key %s, the value was not a boolean", key)
		return defaultVal
	}

	return
}

// GetString returns the string value for the given key.
//
// returns the default value if:
//
// - the library was not initiated;
//
// - the key was not found;
//
// - the key value is empty.
func GetString(key string, defaultVal string) string {
	if localMemory == nil {
		logger.NoCTX().Infof("GetString for key %s, the library was not initiated", key)
		return defaultVal
	}

	val, ok := localMemory[key]
	if !ok || strings.TrimSpace(val) == "" {
		logger.NoCTX().Infof("GetString for key %s, the value was not found or empty", key)
		return defaultVal
	}

	typeKey := fmt.Sprintf("%s.type", key)
	t, ok := localMemory[typeKey]
	if !ok || strings.TrimSpace(t) == "" {
		logger.NoCTX().Infof("GetString for key %s, the value type was not found or empty", key)
		return defaultVal
	}

	if t != "string" {
		logger.NoCTX().Infof("GetString for key %s, the value was not a string", key)
		return defaultVal
	}

	return val
}

// GetNumber returns the number value for the given key.
//
// returns the default value if:
//
// - the library was not initiated;
//
// - the key was not found;
//
// - the key value is empty;
//
// - the key value is not a boolean.
func GetNumber(key string, defaultVal float64) float64 {
	if localMemory == nil {
		logger.NoCTX().Infof("GetNumber for key %s, the library was not initiated", key)
		return defaultVal
	}

	val, ok := localMemory[key]
	if !ok || strings.TrimSpace(val) == "" {
		logger.NoCTX().Infof("GetNumber for key %s, the value was not found or empty", key)
		return defaultVal
	}

	typeKey := fmt.Sprintf("%s.type", key)
	t, ok := localMemory[typeKey]
	if !ok || strings.TrimSpace(t) == "" {
		logger.NoCTX().Infof("GetNumber for key %s, the value type was not found or empty", key)
		return defaultVal
	}

	n, err := strconv.ParseFloat(val, 64)
	if err != nil || t != "number" {
		logger.NoCTX().Infof("GetNumber for key %s, the value is not a valid number", key)
		return defaultVal
	}

	return n
}

// IsEnabledByPercent checks the redis key value for a percentage number (between 0 and 100),
// calculates an random number (also between 0 and 100), and returns true or false depending whether
// the calculated number is within the found percentage.
//
// returns false if:
//
// - the library was not initiated;
//
// - the key was not found;
//
// - the key value is empty;
//
// - the key value is not a percentage (number between 0 and 100);
//
// - the random number greater than the found percentage.
func IsEnabledByPercent(key string) bool {
	if localMemory == nil {
		logger.NoCTX().Infof("IsEnabledByPercent for key %s, the library was not initiated", key)
		return false
	}

	val, ok := localMemory[key]
	if !ok || strings.TrimSpace(val) == "" {
		logger.NoCTX().Infof("IsEnabledByPercent for key %s, the value was not found or empty", key)
		return false
	}

	typeKey := fmt.Sprintf("%s.type", key)
	t, ok := localMemory[typeKey]
	if !ok || strings.TrimSpace(t) == "" {
		logger.NoCTX().Infof("IsEnabledByPercent for key %s, the value type was not found or empty", key)
		return false
	}

	n, err := strconv.Atoi(val)
	if err != nil || t != "number" {
		logger.NoCTX().Infof("IsEnabledByPercent for key %s, the value is not a number", key)
		return false
	}

	if n > 100 || n < 0 {
		logger.NoCTX().Infof("IsEnabledByPercent for key %s, the value is not in percentage format", key)
		return false
	}

	r := rand.Intn(100)
	return r <= n
}

/*
Get gets a feature toggle by the key, but then parses that feature toggle value into the provided type (T) using json decoding.

If the provided type (T) is string, the raw value from the feature toggle will be returned.

returns the default value if:

- the library was not initiated;

- the key was not found;

- the key value is empty.

- the value stored in the key could not be parsed into the provided type (T)
*/
func Get[T any](key string, defaultVal T) (res T) {
	if localMemory == nil {
		logger.NoCTX().Infow("[Feature Toggle] The library was not initiated",
			"key", key,
			"method", "Get",
		)
		return defaultVal
	}

	val, ok := localMemory[key]
	if !ok || strings.TrimSpace(val) == "" {
		logger.NoCTX().Infow("[Feature Toggle] The value was not found",
			"key", key,
			"method", "Get",
		)
		return defaultVal
	}

	resPointer := new(T)
	if reflect.TypeOf(*resPointer).Kind() == reflect.String {
		res = any(val).(T)
		return
	}

	err := json.Unmarshal([]byte(val), resPointer)
	if err != nil {
		title := fmt.Sprintf("[Feature Toggle] Failed to parse the remote config value to a %T value", res)
		logger.NoCTX().Infow(title,
			"key", key,
			"method", "Get",
			"error", err.Error(),
		)

		return defaultVal
	}

	res = *resPointer
	return
}
