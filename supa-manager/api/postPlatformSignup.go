package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"supamanager.io/supa-manager/database"
)

type PlatformSignupBody struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	HcaptchaToken string `json:"hcaptchaToken"`
	RedirectTo    string `json:"redirectTo"`
}

func (a *Api) postPlatformSignup(c *gin.Context) {
	var body PlatformSignupBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if !a.config.AllowSignup {
		c.JSON(403, gin.H{"error": "Signup is disabled"})
		return
	}

	_, err := a.queries.GetAccountByEmail(c.Request.Context(), body.Email)
	if err == nil {
		c.JSON(409, gin.H{"error": "Email already in use"})
		return
	}

	if err != pgx.ErrNoRows {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	hash, err := a.argon.HashEncoded([]byte(body.Password))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	_, err = a.queries.CreateAccount(c.Request.Context(), database.CreateAccountParams{
		Email:        body.Email,
		PasswordHash: string(hash),
		Username:     "idekman",
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"status": "CREATED"})
}
