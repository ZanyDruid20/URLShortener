package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ZanyDruid20/urlshortener/handler"
	"github.com/ZanyDruid20/urlshortener/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
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

	// Connect to MySQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "shortuser:shortpass@tcp(localhost:13306)/shortener"
	}
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}
	defer db.Close()
	fmt.Println("MySQL connected")

	// Pass db to store
	store.SetDB(db)

	r := gin.Default()

	// Add CORS middleware before your routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	err = r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
