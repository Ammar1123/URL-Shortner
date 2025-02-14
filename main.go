// main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables for the clients and context
var (
	redisClient   *redis.Client
	mongoClient   *mongo.Client
	analyticsChan chan AnalyticsEvent
	ctx           = context.Background()
)

// initRedis initializes the Redis client.
func initRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

// initMongo initializes the MongoDB client.
func initMongo() {
	var err error
	mongoAddr := os.Getenv("MONGO_URI")
	if mongoAddr == "" {
		mongoAddr = "mongodb://localhost:27017"
	}
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoAddr))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
}

func main() {
	initRedis()
	initMongo()

	// Initialize the analytics channel and start the worker goroutine.
	analyticsChan = make(chan AnalyticsEvent, 100) // Buffer of 100 events.
	go analyticsWorker()

	// Set up the Gin router.
	router := gin.Default()

	router.Use(rateLimitMiddleware)
	router.POST("/shorten", shortenURLHandler)
	router.GET("/:shortID", redirectHandler)
	router.GET("/analytics", analyticsHandler)

	// Start the API server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
