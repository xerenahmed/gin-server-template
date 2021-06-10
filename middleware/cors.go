package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "*")

	if os.Getenv("MODE") == "dev" {
		c.Header("Access-Control-Allow-Origin", "*")
	} else {
		c.Header("Access-Control-Allow-Origin", "*.redmc.me")
	}

	c.Next()
}
