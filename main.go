package main

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/runeasymail/ManagementAPI/helpers"
	"github.com/runeasymail/ManagementAPI/modules"
)

var log = logging.MustGetLogger("mail")

func main() {

	logging.SetLevel(logging.DEBUG, "")
	var format = logging.MustStringFormatter(`%{color} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logging.SetFormatter(format)

	// read all config
	helpers.ConfigInit()

	// connect mysql
	helpers.InitMysql()

	r := gin.Default()

	// Domains
	r.GET("/domains", modules.HandlerGetAllDomains)
	r.POST("/domains", modules.HandlerAddNewDomain)

	// Users
	r.POST("/users/:domain_id", modules.HandleUserAdd)
	r.GET("/users/:domain_id", modules.HandlerUserLists)

	// No Route
	r.NoRoute(func(c *gin.Context) {

		c.JSON(404, gin.H{
			"code":    404,
			"host":    c.Request.Host,
			"method":  c.Request.Method,
			"url":     c.Request.RequestURI,
			"message": "API method is not exist",
		})

	})

	log.Info("Starting app on port", helpers.Config.App.Port)
	r.Run("127.0.0.1:" + helpers.Config.App.Port)
}
