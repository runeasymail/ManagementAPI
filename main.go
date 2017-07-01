package main

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/runeasymail/ManagementAPI/helpers"
	"github.com/runeasymail/ManagementAPI/middlewares"
	"github.com/runeasymail/ManagementAPI/modules"
)

var log = logging.MustGetLogger("mail")

var git_sha string

func main() {

	logging.SetLevel(logging.DEBUG, "")
	var format = logging.MustStringFormatter(`%{color} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logging.SetFormatter(format)

	// read all config
	helpers.ConfigInit()

	// connect mysql
	helpers.InitMysql()

	// default gin engine
	r := gin.Default()

	// Add proper headers for CORS
	r.Use(middlewares.CORSMiddleware())

	// Auth
	r.POST("/auth", modules.HandlerAuth)

	//
	r.GET("/ping", func(c *gin.Context) {

		c.String(200, "pong")
	})

	authorized := r.Group("")
	authorized.Use(middlewares.AuthMiddleware())
	{
		// Domains
		authorized.GET("/domains", modules.HandlerGetAllDomains)
		authorized.POST("/domains", modules.HandlerAddNewDomain)
		authorized.DELETE("/domain", modules.HandlerRemoveDomain)

		// dkim
		authorized.POST("/dkim/add-new-domain", modules.HandlerNewDkimDomain)

		// Users
		authorized.POST("/users/:domain_id", modules.HandleUserAdd)
		authorized.POST("/user-change-password", modules.HandlerUserPasswordChange)
		authorized.GET("/users/:domain_id", modules.HandlerUserLists)
		authorized.DELETE("/users", modules.HandleUserDelete)

		// ssl
		authorized.POST("/ssl/get-left-days", modules.CheckSSLisValidHandler)
		authorized.POST("/ssl/upload", modules.UploadMySSLHandler)

		// let's encrypt
		authorized.POST("/ssl/letsencrypt/install", modules.LetsEncryptHandler)
	}

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

	log.Info("Version ", git_sha)
	log.Info("Starting app on port", helpers.Config.App.Port)
	r.Run(":" + helpers.Config.App.Port)
}
