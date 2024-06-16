package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type UpdateVersionRequest struct {
	Image string `json:"image" required:"true"`
	Tag   string `json:"tag" required:"true"`
}

func main() {
	config, err := LoadConfig(".env")
	if err != nil {
		println("Failed to load configuration, ensure the required environment variables are set.")
		return
	}
	sshKeys, err := LoadSSHKeys(config)
	if err != nil {
		println("Failed to load SSH keys. Please ensure all Github accounts have SSH keys.")
		return
	}

	conn, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		println(fmt.Sprintf("Unable to connect to database: %v", err))
		return
	}

	queries := New(conn)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		versions, err := queries.GetVersions(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(200, versions)
	})

	r.GET("/:service", func(c *gin.Context) {
		service := c.Param("service")
		versions, err := queries.GetVersionsForService(c.Request.Context(), strings.ToLower(service))
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(200, versions)
	})

	r.POST("/:service/update-version", func(c *gin.Context) {
		signature := c.GetHeader("Signature")
		if signature == "" {
			c.JSON(400, gin.H{"error": "missing Signature header"})
			return
		}

		service := c.Param("service")

		var body UpdateVersionRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		isValidSig := false
		var account string
		for _, key := range sshKeys {
			if err := VerifySignature(key.Key, []byte(fmt.Sprintf("%s:%s", body.Image, body.Tag)), []byte(signature)); err == nil {
				isValidSig = true
				account = key.Username
				break
			}
		}

		if !isValidSig {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		svc, err := queries.CreateNewVersion(c.Request.Context(), CreateNewVersionParams{
			ServiceID: strings.ToLower(service),
			Image:     body.Image,
			Tag:       body.Tag,
			CreatedBy: strings.ToLower(account),
		})

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(200, gin.H{
			"id": svc.ID,
		})
	})

	r.Run(config.ListenAddress)
}
