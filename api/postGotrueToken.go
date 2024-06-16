package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type GotrueToken struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	GotrueMetaSecurity struct {
		CaptchaToken string `json:"captcha_token"`
	} `json:"gotrue_meta_security"`
}

func (a *Api) postGotrueToken(c *gin.Context) {
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
