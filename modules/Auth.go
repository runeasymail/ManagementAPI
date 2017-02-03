package modules

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/runeasymail/ManagementAPI/helpers"
	"net/http"
	"time"
)

func HandlerAuth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	log := logging.MustGetLogger("mail")

	log.Debug(helpers.Config.Auth)

	if username == "" && password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"result": false, "msg": "Username/password is required"})
		return
	}
	if helpers.Config.Auth.Username != username || helpers.Config.Auth.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"result": false, "msg": "Username/password is not correct"})
		return
	}

	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	claims := make(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(helpers.Config.Auth.SecretKey))

	c.JSON(200, gin.H{"result": true, "token": tokenString})
}
