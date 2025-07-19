package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// we're defining the struct around redis
type StorageService struct {
	redisClient *redis.Client
}

// we're declaring top level for the storage and the context
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

func init() {
	initializeStore()
}

const cacheDuration = 6 * time.Hour

// we're initializing the store services and we are returning a store pointer as a result
func initializeStore() *StorageService {
	redisAddr := "localhost:6379"
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		redisAddr = addr
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}
	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

/*
we are saving the mapping of the url between the original url and then
the generated url
*/
func SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, cacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}

}

/*
After this we are going to retrieve the initial url once the generated one has been provided
and users will be calling the shortlink in the url , therefore wh have to retrieve the original url and then
redirect
*/
func RetrieveInitialUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
	return result

}
