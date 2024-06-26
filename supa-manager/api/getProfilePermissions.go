package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"supamanager.io/supa-manager/permisions"
)

func (a *Api) getProfilePermissions(c *gin.Context) {
	acc, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgIds, err := a.queries.GetOrganizationIdsForAccountId(c, acc.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, permisions.ConstructPermissions(orgIds))
}
