// main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/image", "./image")

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})

	router.Run(":8080")
}
