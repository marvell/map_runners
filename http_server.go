package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) RunHTTPServer(ctx context.Context) error {
	router := gin.Default()

	router.Static("/static", "./static")
	router.GET("/:user_id", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.GET("/strava/login", app.StravaLoginHandler())
	router.GET("/strava/callback", app.StravaCallbackHandler())

	router.GET("/", app.IndexHandler())
	router.GET("/:user_id/data", app.UserHandler())

	log.Printf("http-server is running on %s", app.cfg.HTTPServerAddr)
	return router.Run(app.cfg.HTTPServerAddr)
}

func printError(c *gin.Context, err error) {
	log.Printf("ERROR %s", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"type":    fmt.Sprintf("%t", err),
			"message": err.Error(),
		},
	})
}
