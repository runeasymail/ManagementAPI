package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Auth-token, X-Requested-With, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("options")
			c.JSON(200, "{}")
			return
		}
		// c.Next()
	}
}
