package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"supamanager.io/supa-manager/database"
)

type CreateOrgParams struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Size string `json:"size"`
	Tier string `json:"tier"`
}

// TODO: use tx
func (a *Api) postPlatformOrganizations(c *gin.Context) {
	account, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var params CreateOrgParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	org, err := a.queries.CreateOrganization(c.Request.Context(), params.Name)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	_, err = a.queries.CreateOrganizationMembership(c.Request.Context(), database.CreateOrganizationMembershipParams{
		OrganizationID: org.ID,
		AccountID:      account.ID,
		Role:           "OWNER",
	})

	c.AbortWithStatus(http.StatusCreated)
}
