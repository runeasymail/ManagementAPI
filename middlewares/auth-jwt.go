package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/runeasymail/ManagementAPI/helpers"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		myToken := c.Request.Header.Get("Auth-token")

		token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(helpers.Config.Auth.SecretKey), nil
		})

		if err != nil {
			c.String(200, "Token is not correct")
			c.Abort()
			return
		}

		if !token.Valid {
			c.String(200, string(err.Error()))
			c.Abort()
			return

		}

		c.Set("username", token.Claims.(jwt.MapClaims)["username"])

		c.Next()
	}
}
