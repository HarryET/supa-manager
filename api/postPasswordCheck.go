package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustelem/zxcvbn"
	"net/http"
)

type PasswordCheckBody struct {
	Password string `json:"password"`
}

func (a *Api) postPasswordCheck(c *gin.Context) {
	var body PasswordCheckBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	result := zxcvbn.PasswordStrength(body.Password, nil)

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"score": result.Score,
		},
	})
}
