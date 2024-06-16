package api

import (
	"github.com/gin-gonic/gin"
	"github.com/harryet/supa-manager/database"
	"net/http"
	"strings"
)

type Organization struct {
	Slug              string        `json:"slug"`
	Name              string        `json:"name"`
	StripeCustomerId  string        `json:"stripe_customer_id"`
	SubscriptionId    string        `json:"subscription_id"`
	BillingEmail      interface{}   `json:"billing_email"`
	IsOwner           bool          `json:"is_owner"`
	OptInTags         []interface{} `json:"opt_in_tags"`
	Id                int32         `json:"id"`
	RestrictionData   interface{}   `json:"restriction_data"`
	RestrictionStatus interface{}   `json:"restriction_status"`
}

func (a *Api) getOrganizations(c *gin.Context) {
	account, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgs, err := a.queries.GetOrganizationsForAccountId(c, account.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	supaOrgs := []Organization{}
	for _, org := range orgs {
		supaOrgs = append(supaOrgs, Organization{
			Slug:              org.Slug,
			Name:              org.Name,
			StripeCustomerId:  "",
			SubscriptionId:    "",
			BillingEmail:      "billing@supamanager.io",
			IsOwner:           strings.ToLower(org.MemberRole) == "owner",
			OptInTags:         []interface{}{},
			Id:                org.ID,
			RestrictionData:   nil,
			RestrictionStatus: nil,
		})
	}

	c.JSON(200, supaOrgs)
}

type CreateOrgParams struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Size string `json:"size"`
	Tier string `json:"tier"`
}

// TODO: use tx
func (a *Api) platformCreateOrganization(c *gin.Context) {
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

func (a *Api) platformReachedFreeProjectLimit(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, []interface{}{})
}
