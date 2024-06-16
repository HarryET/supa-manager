package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DNSUpdateType string

const (
	DNSUpdateTypeCreate DNSUpdateType = "CREATE"
	DNSUpdateTypeDelete DNSUpdateType = "DELETE"
)

type RequestBody struct {
	Type       DNSUpdateType `json:"type"`
	Hostname   string        `json:"hostname"`
	ProjectRef string        `json:"project_ref"`
}

func main() {
	config, err := LoadConfig(".env")
	if err != nil {
		println("Failed to load configuration, ensure the required environment variables are set.")
		return
	}

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		tokenHeader := c.GetHeader("x-api-key")
		if tokenHeader != config.Token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var body RequestBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// TODO: DNS Work Here

		c.AbortWithStatus(http.StatusAccepted)
	})

	r.Run(config.ListenAddress)
}
