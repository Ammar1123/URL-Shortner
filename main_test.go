// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// init function to initialize dependencies before running tests.
func init() {
	initRedis()
	initMongo()
	// Initialize the analytics channel to avoid nil channel usage.
	analyticsChan = make(chan AnalyticsEvent, 100)
}

func TestShortenURLHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/shorten", shortenURLHandler)

	// Prepare a sample request.
	reqBody, _ := json.Marshal(gin.H{"url": "http://www.example.com"})
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Validate response.
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "short_url")
}

func TestRedirectHandlerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortID", redirectHandler)

	// Use a non-existent shortID.
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Validate that we get a 404.
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
