package store

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
}

// InitializeStore creates and returns a new instance of StorageService.
// You may need to adjust the implementation to match your actual initialization logic.
func InitializeStore() *StorageService {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Initialize MySQL DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "shortuser:shortpass@tcp(localhost:13306)/shortener"
	}
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err)
	}

	return &StorageService{
		redisClient: redisClient,
		db:          db,
	}
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	// Persist data mapping
	SaveUrlMapping(shortURL, initialLink, userUUId)

	// Retrieve initial URL
	retrievedUrl := RetrieveInitialUrl(shortURL)

	assert.Equal(t, initialLink, retrievedUrl)
}

func TestMySQLUrlMapping(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "shortuser:shortpass@tcp(localhost:13306)/shortener"
	}
	db, err := sql.Open("mysql", dbURL)
	assert.NoError(t, err)
	defer db.Close()

	SetDB(db)

	shortUrl := "testShort"
	originalUrl := "https://example.com"

	// Save mapping
	err = SaveUrlMappingMySQL(shortUrl, originalUrl)
	assert.NoError(t, err)

	// Retrieve mapping
	retrievedUrl, err := RetrieveInitialUrlMySQL(shortUrl)
	assert.NoError(t, err)
	assert.Equal(t, originalUrl, retrievedUrl)
}

func TestMySQLWrongConnection(t *testing.T) {
	// Intentionally use wrong password
	dbURL := "shortuser:wrongpass@tcp(localhost:13306)/shortener"
	db, err := sql.Open("mysql", dbURL)
	assert.NoError(t, err) // sql.Open does not connect immediately

	// Now actually try to ping
	err = db.Ping()
	assert.Error(t, err, "Expected error due to wrong password")
	db.Close()
}
