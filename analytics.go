package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AnalyticsEvent struct {
	ShortID   string    `bson:"short_id"`
	Timestamp time.Time `bson:"timestamp"`
	IP        string    `bson:"ip"`
}

// analyticsWorker reads events from analyticsChan and writes them to MongoDB.
func analyticsWorker() {
	collection := mongoClient.Database("urlshortener").Collection("analytics")
	for event := range analyticsChan {
		_, err := collection.InsertOne(context.Background(), bson.M{
			"short_id":  event.ShortID,
			"timestamp": event.Timestamp,
			"ip":        event.IP,
		})
		if err != nil {
			log.Printf("Failed to log analytics event: %v", err)
		}
	}
}
