package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"server-app/server/manager/db"
	"server-app/utils"
	"time"
)

func Login(db *db.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if len(username) < 4 || len(password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid username or password",
			})
			return
		}

		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// validations

		if true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid password",
			})
			return
		}

		token, err := utils.GenerateAccessToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "can't create token: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
