package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
}

// InitializeStore creates and returns a new instance of StorageService.
// You may need to adjust the implementation to match your actual initialization logic.
func InitializeStore() *StorageService {
	return &StorageService{
		// Initialize fields as needed, e.g. redisClient: ...
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
