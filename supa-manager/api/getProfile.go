package api

import (
	"github.com/gin-gonic/gin"
	"supamanager.io/supa-manager/utils"
)

type ProfileReturn struct {
	Id               int32         `json:"id"`
	Auth0Id          *string       `json:"auth0_id"`
	PrimaryEmail     string        `json:"primary_email"`
	Username         *string       `json:"username"`
	FirstName        *string       `json:"first_name"`
	LastName         *string       `json:"last_name"`
	Mobile           *string       `json:"mobile"`
	IsAlphaUser      bool          `json:"is_alpha_user"`
	GotrueId         string        `json:"gotrue_id"`
	FreeProjectLimit int           `json:"free_project_limit"`
	DisabledFeatures []interface{} `json:"disabled_features"`
}

func (a *Api) getProfile(c *gin.Context) {
	account, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(200, ProfileReturn{
		Id:               account.ID,
		Auth0Id:          nil,
		PrimaryEmail:     account.Email,
		FirstName:        utils.PgTextToPointer(account.FirstName),
		LastName:         utils.PgTextToPointer(account.LastName),
		Mobile:           nil,
		IsAlphaUser:      true,
		GotrueId:         account.GotrueID,
		FreeProjectLimit: 9999,
		DisabledFeatures: []interface{}{},
	})
}
