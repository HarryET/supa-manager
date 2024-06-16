package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectStatus(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	project, err := a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": project.Status})
}
