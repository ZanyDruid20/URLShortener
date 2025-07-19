package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZanyDruid20/urlshortener/handler"
	"github.com/ZanyDruid20/urlshortener/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func init() {
	godotenv.Load()

}

func main() {
	// Load .env file (optional, don't panic if missing)
	godotenv.Load()

	// Connect to Redis
	redisAddr := "localhost:6379"
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		redisAddr = addr
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Println("Redis connected:", pong)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hey Go URL Shortener !",
		})
	})

	// Test Redis set/get endpoint
	r.GET("/redis-test", func(c *gin.Context) {
		err := rdb.Set(ctx, "testkey", "testvalue", 0).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		val, err := rdb.Get(ctx, "testkey").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"testkey": val})
	})
	r.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	// Note that store initialization happens here
	store.InitializeStore()

	err = r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
