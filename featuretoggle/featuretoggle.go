package featuretoggle

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/delivery-much/dm-go/logger"
)

var client redisClient

// Init inits the feature toggle library
func Init(c Config) error {
	cl, err := getRedisClient(c.Host, c.Port, c.DB)
	if err != nil {
		return err
	}

	client = cl
	logger.Infof("Redis feature toggle started for service %s", c.ServiceName)
	return nil
}

// IsEnabled checks if given feature key is enabled in redis DB.
// returns the default value if:
// - the client was not instantiated;
// - the key was not found;
// - the key value is empty;
// - the key value is not a boolean.
func IsEnabled(key string, defaultVal bool) (b bool) {
	if client == nil {
		logger.Infof("IsEnabled for key %s, client is empty", key)
		return defaultVal
	}

	val, err := client.get(key)
	if err != nil {
		logger.Infof("IsEnabled for key %s, the client returned an error: %v", key, err)
		return defaultVal
	}
	if strings.TrimSpace(val) == "" {
		logger.Infof("IsEnabled for key %s, the value was not found or empty", key)
		return defaultVal
	}

	b, err = strconv.ParseBool(val)
	if err != nil {
		logger.Infof("IsEnabled for key %s, the value was not a boolean", key)
		return defaultVal
	}

	return
}

// GetString returns the string value for the given key.
// returns the default value if:
// - the client was not instantiated;
// - the key was not found;
// - the key value is empty.
func GetString(key string, defaultVal string) string {
	if client == nil {
		logger.Infof("GetString for key %s, client is empty", key)
		return defaultVal
	}

	val, err := client.get(key)
	if err != nil {
		logger.Infof("GetString for key %s, the client returned an error: %v", key, err)
		return defaultVal
	}
	if strings.TrimSpace(val) == "" {
		logger.Infof("GetString for key %s, the value was not found or empty", key)
		return defaultVal
	}

	return val
}

// GetNumber returns the number value for the given key.
// returns the default value if:
// - the client was not instantiated;
// - the key was not found;
// - the key value is empty;
// - the key value is not a boolean.
func GetNumber(key string, defaultVal float64) float64 {
	if client == nil {
		logger.Infof("GetNumber for key %s, client is empty", key)
		return defaultVal
	}

	val, err := client.get(key)
	if err != nil {
		logger.Infof("GetNumber for key %s, the client returned an error: %v", key, err)
		return defaultVal
	}
	if strings.TrimSpace(val) == "" {
		logger.Infof("GetNumber for key %s, the value was not found or empty", key)
		return defaultVal
	}

	n, err := strconv.ParseFloat(val, 10)
	if err != nil {
		logger.Infof("GetNumber for key %s, the value is not a valid number", key)
		return defaultVal
	}

	return n
}

// IsEnabledByPercent checks the redis key value for a percentage number (between 0 and 100),
// calculates an random number (also between 0 and 100), and returns true or false depending whether
// the calculated number is within the found percentage.
// returns false if:
// - the client was not instantiated;
// - the key was not found;
// - the key value is empty;
// - the key value is not a percentage (number between 0 and 100);
// - the random number greater than the found percentage.
func IsEnabledByPercent(key string) bool {
	if client == nil {
		logger.Infof("IsEnabledByPercent for key %s, client is empty", key)
		return false
	}

	val, err := client.get(key)
	if err != nil {
		logger.Infof("IsEnabledByPercent for key %s, the client returned an error: %v", key, err)
		return false
	}
	if strings.TrimSpace(val) == "" {
		logger.Infof("IsEnabledByPercent for key %s, the value was not found or empty", key)
		return false
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		logger.Infof("IsEnabledByPercent for key %s, the value is not a number", key)
		return false
	}

	if n > 100 || n < 0 {
		logger.Infof("IsEnabledByPercent for key %s, the value is not in percentage format", key)
		return false
	}

	r := rand.Intn(100)
	return r <= n
}
