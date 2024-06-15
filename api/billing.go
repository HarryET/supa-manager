package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) platformOverdueInvoices(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, []interface{}{})
}

func (a *Api) platformSetupIntent(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret": "seti_1PSxxxxxPojXS6xxxxxxxxxx_secret_QIf9hLhCMsAtxxxxxxxxF9ReapCGOQP",
	})
}
