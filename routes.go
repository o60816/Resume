// routes.go

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func initializeRoutes() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.Static("/image", "./image")
	router.Static("/profile", "./profile")

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	router.GET("/login", showLoginPage)
	router.POST("/login", login)
	router.GET("/logout", logout)
	router.GET("/edit/person/:userId", showEditPage)
	router.GET("/person/:userId", showUserPage)

	router.GET("/user/update/:userId", showUpdateUserPage)
	router.POST("/user/update/:userId", userHandler)

	router.GET("/work/add", showAddWorkPage)
	router.GET("/work/update/:workId", showUpdateWorkPage)
	router.POST("/work/:workId", workHandler)
	router.DELETE("/work/:workId", workHandler)

	router.GET("/project/add/:workId", showAddProjectPage)
	router.GET("/project/update/:projectId", showUpdateProjectPage)
	router.POST("/project/:projectId", projectHandler)
	router.DELETE("/project/:projectId", projectHandler)

	if err := router.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
