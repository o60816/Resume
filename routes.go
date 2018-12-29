// routes.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func initializeRoutes() {
	router = gin.Default()
	router.Static("/image", "./image")

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/zh")
	})
	router.GET("/en", showMainPage)
	router.GET("/zh", showMainPage)
	router.GET("/edit/en", editPage)
	router.GET("/edit/zh", editPage)
	router.GET("/work/add", showAddWorkPage)
	router.GET("/work/update/:workId", showUpdateWorkPage)
	router.POST("/work/:workId", workHandler)
	router.DELETE("/work/:workId", workHandler)
	router.GET("/project/add/:workId", showAddProjectPage)
	router.GET("/project/update/:projectId", showUpdateProjectPage)
	router.POST("/project/:projectId", projectHandler)
	router.DELETE("/project/:projectId", projectHandler)

	router.Run(":80")
}
