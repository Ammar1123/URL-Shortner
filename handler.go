package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	nanoid "github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

// shortenRequest defines the expected JSON payload for URL shortening.
type shortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func analyticsHandler(c *gin.Context) {
	collection := mongoClient.Database("urlshortener").Collection("analytics")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}
	defer cursor.Close(ctx)

	var events []AnalyticsEvent
	if err := cursor.All(ctx, &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse analytics data"})
		return
	}

	c.JSON(http.StatusOK, events)

}

// shortenURLHandler handles the POST /shorten endpoint.
func shortenURLHandler(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	generator, err := nanoid.Custom("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
		return
	}
	shortID := generator()

	if err := redisClient.Set(ctx, shortID, req.URL, 30*24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL mapping"})
		return
	}

	// Construct the short URL (adjust the domain as necessary).
	shortURL := "http://urlShortner.com/" + shortID
	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

// redirectHandler handles the GET /:shortID endpoint.
func redirectHandler(c *gin.Context) {
	shortID := c.Param("shortID")

	longURL, err := redisClient.Get(ctx, shortID).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL mapping", "details": err.Error()})
		}
		return
	}

	// Offload analytics logging and redirect
	analyticsEvent := AnalyticsEvent{
		ShortID:   shortID,
		Timestamp: time.Now(),
		IP:        c.ClientIP(),
	}
	select {
	case analyticsChan <- analyticsEvent:
	default:
	}
	c.Redirect(http.StatusFound, longURL)
}
