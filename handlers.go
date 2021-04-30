package main

import (
	"Resume/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var tblWork string
var tblProject string
var language string
var currentUserId int

func initUsedTable(language string) {
	if language == "en" {
		tblWork = "work_en"
		tblProject = "project_en"
	} else {
		tblWork = "work"
		tblProject = "project"
	}
}

func showUserPage(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	router.LoadHTMLGlob("templates/login.html")
	user, err := models.GetUserById(userId)
	if err != nil || 0 == user.Id {
		log.Panic()
		return
	}

	router.LoadHTMLGlob(fmt.Sprintf("templates/%s/*", "zh"))

	initUsedTable(language)

	workList, err := models.GetAllWork(tblWork, userId)

	if err != nil {
		log.Panic(err)
		return
	}

	for i := range workList {
		projectList, err := models.GetProjectByWorkId(tblProject, workList[i].Id)
		if err != nil {
			log.Panic(err)
			return
		}
		workList[i].ProjectList = projectList
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"user":  user,
			"works": workList,
		},
	)
}

func showLoginPage(c *gin.Context) {
	if currentUserId != 0 {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/edit/person/%d", currentUserId))
		return
	}

	router.LoadHTMLGlob("templates/login.html")
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{},
	)
}

func login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	user, result := models.Login(account, password)
	if result != true {
		router.LoadHTMLGlob("templates/login.html")
		c.HTML(
			http.StatusOK,
			"login.html",
			gin.H{},
		)
		return
	}
	currentUserId = user.Id
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/edit/person/%d", user.Id))
}

func logout(c *gin.Context) {
	currentUserId = 0
	router.LoadHTMLGlob("templates/login.html")
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{},
	)
}

func showEditPage(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	if currentUserId != userId {
		currentUserId = 0
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/login"))
		return
	}
	user, err := models.GetUserById(userId)
	if err != nil {
		log.Panic(err)
		return
	}

	router.LoadHTMLGlob(fmt.Sprintf("templates/%s/*", "zh"))
	initUsedTable("zh")

	workList, err := models.GetAllWork(tblWork, userId)

	if err != nil {
		log.Panic(err)
		return
	}

	for i := range workList {
		projectList, err := models.GetProjectByWorkId(tblProject, workList[i].Id)
		if err != nil {
			log.Panic(err)
			return
		}
		workList[i].ProjectList = projectList
	}

	c.HTML(
		http.StatusOK,
		"edit.html",
		gin.H{
			"user":  user,
			"works": workList,
		},
	)
}

func showUpdateUserPage(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	user, err := models.GetUserById(userId)
	if err != nil {
		log.Panic(err)
		return
	}

	c.HTML(
		http.StatusOK,
		"editUser.html",
		gin.H{
			"user": user,
		},
	)
}

func userHandler(c *gin.Context) {
	var user models.User
	user.Id, _ = strconv.Atoi(c.Param("userId"))
	user.Name = c.PostForm("name")
	user.Title = c.PostForm("title")
	user.Information = c.PostForm("information")
	user.Introduction = c.PostForm("introduction")
	user.Motto = c.PostForm("motto")
	user.Photo = c.PostForm("photo")
	user.Github = c.PostForm("github")
	user.Linkedin = c.PostForm("linkedin")
	user.Facebook = c.PostForm("facebook")
	user.Email = c.PostForm("email")
	err := models.UpdateUser(user)
	if err != nil {
		log.Panic(err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/edit/person/%d", user.Id))
}

func showAddWorkPage(c *gin.Context) {
	var btnName string
	if strings.Contains(language, "en") {
		btnName = "Add"
	} else {
		btnName = "新增"
	}
	var work models.Work
	c.HTML(
		http.StatusOK,
		"addWork.html",
		gin.H{
			"method":  "POST",
			"work":    work,
			"btnName": btnName,
		},
	)
}

func showUpdateWorkPage(c *gin.Context) {
	work, err := models.GetWorkByID(tblWork, c.Param("workId"))
	if err != nil {
		log.Panic(err)
		return
	}

	var btnName string
	if strings.Contains(language, "en") {
		btnName = "Update"
	} else {
		btnName = "更新"
	}

	c.HTML(
		http.StatusOK,
		"addWork.html",
		gin.H{
			"method":  "PATCH",
			"work":    work,
			"btnName": btnName,
		},
	)
}

func workHandler(c *gin.Context) {
	method := c.Request.Method

	var err error
	var work models.Work
	work.Id, err = strconv.Atoi(c.Param("workId"))
	if err != nil {
		log.Panic(err)
		return
	}

	if method != "DELETE" {
		method = c.PostForm("_method")
		work.Period = c.PostForm("period")
		work.Logo = c.PostForm("logo")
		work.Company = c.PostForm("company")
		work.Position = c.PostForm("position")
		work.Content = c.PostForm("content")
		work.UserId = currentUserId
	}

	if _, err = models.EditWork(tblWork, tblProject, method, work); err != nil {
		log.Panic(err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/edit/person/%d", currentUserId))
}

func showAddProjectPage(c *gin.Context) {
	var btnName string
	if strings.Contains(language, "en") {
		btnName = "Add"
	} else {
		btnName = "新增"
	}

	var project models.Project
	var err error
	project.WorkId, err = strconv.Atoi(c.Param("workId"))
	if err != nil {
		log.Panic(err)
		return
	}
	c.HTML(
		http.StatusOK,
		"addProject.html",
		gin.H{
			"method":  "POST",
			"project": project,
			"btnName": btnName,
		},
	)
}

func showUpdateProjectPage(c *gin.Context) {
	var btnName string
	if strings.Contains(language, "en") {
		btnName = "Update"
	} else {
		btnName = "更新"
	}

	project, err := models.GetProjectById(tblProject, c.Param("projectId"))
	if err != nil {
		log.Panic(err)
		return
	}
	c.HTML(
		http.StatusOK,
		"addProject.html",
		gin.H{
			"method":  "PATCH",
			"project": project,
			"btnName": btnName,
		},
	)
}

func projectHandler(c *gin.Context) {
	method := c.Request.Method
	var err error
	var project models.Project
	project.Id, err = strconv.Atoi(c.Param("projectId"))
	if err != nil {
		log.Panic(err)
		return
	}
	if method != "DELETE" {
		method = c.PostForm("_method")
		project.WorkId, err = strconv.Atoi(c.PostForm("workId"))
		if err != nil {
			log.Panic(err)
			return
		}
		project.ProjectName = c.PostForm("projectName")
		project.Tech = c.PostForm("tech")
	}
	if _, err := models.EditProject(tblProject, method, project); err != nil {
		log.Panic(err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/edit/person/%d", currentUserId))
}
