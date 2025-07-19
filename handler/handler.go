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

// Handler function to create a short URL
func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest // Declare a variable to hold the incoming request data

	// Bind the incoming JSON to the UrlCreationRequest struct and check for errors
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a 400 error if binding fails
		return // Exit the function if there is an error
	}

	// Generate a short URL using the long URL and user ID
	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)

	// Save the mapping of the short URL to the original long URL and user ID
	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	host := "http://localhost:9808/" // Define the host URL for the shortener service

	// Respond with a JSON message containing the short URL
	c.JSON(200, gin.H{
		"message":   "short url created successfully", // Success message
		"short_url": host + shortUrl,                  // The full short URL
	})
}

// Handler function to redirect from a short URL to the original long URL
func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")           // Get the short URL parameter from the request path
	initialUrl := store.RetrieveInitialUrl(shortUrl) // Retrieve the original long URL from the store
	c.Redirect(302, initialUrl)               // Redirect the client to the original long URL with a 302 status
}
