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

type OrgBillingSubscription struct {
	NanoEnabled        bool `json:"nano_enabled"`
	BillingViaPartner  bool `json:"billing_via_partner"`
	CurrentPeriodEnd   int  `json:"current_period_end"`
	CurrentPeriodStart int  `json:"current_period_start"`
	NextInvoiceAt      int  `json:"next_invoice_at"`
	CustomerBalance    int  `json:"customer_balance"`
	Plan               struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"plan"`
	UsageBillingEnabled bool          `json:"usage_billing_enabled"`
	Addons              []interface{} `json:"addons"`
	ProjectAddons       []interface{} `json:"project_addons"`
	PaymentMethodType   string        `json:"payment_method_type"`
}

func (a *Api) platformOrganizationBillingSubscription(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(200, OrgBillingSubscription{
		NanoEnabled:        false,
		BillingViaPartner:  false,
		CurrentPeriodEnd:   2147385600, // Jan 18, 2038
		CurrentPeriodStart: 0,          // Start at 1st Jan 1970
		NextInvoiceAt:      2147385600, // Jan 18, 2038
		CustomerBalance:    999999,     // Balling 💰
		Plan: struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Id:   "enterprise",
			Name: "Enterprise",
		},
		UsageBillingEnabled: false,
		Addons:              []interface{}{},
		ProjectAddons:       []interface{}{},
		PaymentMethodType:   "none",
	})
}
