package api

import (
	"github.com/gin-gonic/gin"
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
			BillingEmail:      org.BillingEmail,
			IsOwner:           org.MemberRole == "owner",
			OptInTags:         []interface{}{},
			Id:                org.ID,
			RestrictionData:   nil,
			RestrictionStatus: nil,
		})
	}

	c.JSON(200, gin.H{"data": supaOrgs})
}
