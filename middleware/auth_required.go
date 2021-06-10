package middleware

import (
	"github.com/gin-gonic/gin"
	"server-app/utils"
	"net/http"
	"strings"
)

func AuthRequired(c *gin.Context) {
	auth := strings.Split(c.GetHeader("Authorization"), " ")

	if len(auth) == 2 && auth[0] == "Bearer" {
		username, err := utils.ParseToken(utils.AccessSecret, auth[1])
		if err == nil && username != "" {
			// handle
			if true {
				// c.Set("user", user)
				c.Next()
				return
			}
		}
	}

	c.Abort()
	c.String(http.StatusUnauthorized, "unauthorized")
}
