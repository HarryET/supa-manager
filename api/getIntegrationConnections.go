package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getIntegrationConnections(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}})
}
