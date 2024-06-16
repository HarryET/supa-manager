package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getConfigCatConfiguration(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"p": gin.H{
			"u": "https://cdn-global.configcat.com",
			"r": 0,
		},
		"f": gin.H{},
	})
}
