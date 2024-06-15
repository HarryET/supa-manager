package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harryet/supa-manager/database"
	"github.com/jackc/pgx/v5"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type PlatformSignupBody struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	HcaptchaToken string `json:"hcaptchaToken"`
	RedirectTo    string `json:"redirectTo"`
}

func (a *Api) platformSignup(c *gin.Context) {
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

type GotrueToken struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	GotrueMetaSecurity struct {
		CaptchaToken string `json:"captcha_token"`
	} `json:"gotrue_meta_security"`
}

func (a *Api) gotrueToken(c *gin.Context) {
	var body GotrueToken
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	account, err := a.queries.GetAccountByEmail(c.Request.Context(), body.Email)
	if err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	if verified, err := argon2.VerifyEncoded([]byte(body.Password), []byte(account.PasswordHash)); err != nil || !verified {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "supamanager.io",
		Subject:   account.GotrueID,
		Audience:  []string{"supamanager.io"},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(0, 0, 1)},
		NotBefore: &jwt.NumericDate{Time: time.Now()},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// TODO: use a real secret
	signedJwt, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": signedJwt,
		"token_type":   "Bearer",
		"expires_in":   86400,
		// TODO: support refresh tokens
		"refresh_token": signedJwt,
		// TODO: look at GoTrue code for more fields
		"user": gin.H{
			"id":    account.ID,
			"email": account.Email,
			"app_metadata": gin.H{
				"provider": "email",
			},
		},
	})
}
