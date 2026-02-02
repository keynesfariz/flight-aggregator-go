package libs

import (
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func GetCacheClientInstance() *redis.Client {
	clientDB, _ := strconv.Atoi(GetEnv("REDIS_DATABASE", "0"))

	addr := GetEnv("REDIS_HOST", "localhost") + ":" + GetEnv("REDIS_PORT", "6379")
	maxRetries, _ := strconv.Atoi(GetEnv("REDIS_MAX_RETRIES", "3"))

	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:       addr,                         // use default redis address
			Password:   GetEnv("REDIS_PASSWORD", ""), // no password set
			DB:         clientDB,                     // no password set
			MaxRetries: maxRetries,
		})
	})

	return client
}
