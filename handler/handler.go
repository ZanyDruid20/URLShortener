package handler

import (
	"net/http" // Import the net/http package for HTTP status codes

	"github.com/ZanyDruid20/urlshortener/shortener" // Import the shortener package for generating short URLs
	"github.com/ZanyDruid20/urlshortener/store"     // Import the store package for saving and retrieving URLs
	"github.com/gin-gonic/gin"                      // Import the Gin web framework
)

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"` // Field for the original long URL, required in JSON
	UserId  string `json:"user_id" binding:"required"`  // Field for the user ID, required in JSON
}

// Handler function to create a short URL using PostgreSQL
func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest // Declare a variable to hold the incoming request data

	// Bind the incoming JSON to the UrlCreationRequest struct and check for errors
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a 400 error if binding fails
		return                                                     // Exit the function if there is an error
	}

	// Generate a short URL using the long URL and user ID
	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)

	// Save mapping in PostgreSQL
	if err := store.SaveUrlMappingMySQL(shortUrl, creationRequest.LongUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL mapping"})
		return
	}

	host := "http://localhost:9808/" // Define the host URL for the shortener service

	// Respond with a JSON message containing the short URL
	c.JSON(200, gin.H{
		"message":   "short url created successfully", // Success message
		"short_url": host + shortUrl,                  // The full short URL
	})
}

// Handler function to redirect using PostgreSQL
func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")                  // Get the short URL parameter from the request path
	initialUrl := store.RetrieveInitialUrl(shortUrl) // Retrieve the original long URL from the store
	if initialUrl == "" {
		// Fallback to PostgreSQL if not found in Redis
		var err error
		initialUrl, err = store.RetrieveInitialUrlMySQL(shortUrl)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"}) // Return a 404 error if the short URL is not found
			return                                                             // Exit the function if there is an error
		}
		// Optionally cache in Redis for future requests
		store.SaveUrlMapping(shortUrl, initialUrl, "")
	}

	c.Redirect(302, initialUrl) // Redirect the client to the original long URL with a 302 status
}
