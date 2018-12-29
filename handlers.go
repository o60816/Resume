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

func initUsedTable(language string) {
	if language == "en" {
		tblWork = "work_en"
		tblProject = "project_en"
	} else {
		tblWork = "work"
		tblProject = "project"
	}
}

func showMainPage(c *gin.Context) {
	language = c.Request.URL.String()[1:]
	router.LoadHTMLGlob(fmt.Sprintf("templates/%s/*", language))

	initUsedTable(language)
	workList, err := models.GetAllWork(tblWork)

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}

	for i := range workList {
		projectList, err := models.GetProjectByWorkId(tblProject, workList[i].Id)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusOK, gin.H{"Err": err})
		}
		workList[i].ProjectList = projectList
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"works": workList,
		},
	)
}

func editPage(c *gin.Context) {
	language = c.Request.URL.String()[6:]
	router.LoadHTMLGlob(fmt.Sprintf("templates/%s/*", language))

	initUsedTable(language)
	workList, err := models.GetAllWork(tblWork)

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}

	for i := range workList {
		projectList, err := models.GetProjectByWorkId(tblProject, workList[i].Id)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusOK, gin.H{"Err": err})
		}
		workList[i].ProjectList = projectList
	}

	c.HTML(
		http.StatusOK,
		"edit.html",
		gin.H{
			"works": workList,
		},
	)
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
		c.JSON(http.StatusOK, gin.H{"Err": err})
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
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}

	if method != "DELETE" {
		method = c.PostForm("_method")
		work.Period = c.PostForm("period")
		work.Logo = c.PostForm("logo")
		work.Company = c.PostForm("company")
		work.Position = c.PostForm("position")
		work.Content = c.PostForm("content")
	}

	if _, err = models.EditWork(tblWork, tblProject, method, work); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}
	if method == "DELETE" {
		c.JSON(http.StatusOK, gin.H{"status": "delete successfully"})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/edit/"+language)
	}
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
		c.JSON(http.StatusOK, gin.H{"Err": err})
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
		c.JSON(http.StatusOK, gin.H{"Err": err})
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
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}
	if method != "DELETE" {
		method = c.PostForm("_method")
		project.WorkId, err = strconv.Atoi(c.PostForm("workId"))
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusOK, gin.H{"Err": err})
		}
		project.ProjectName = c.PostForm("projectName")
		project.Tech = c.PostForm("tech")
	}
	if _, err := models.EditProject(tblProject, method, project); err != nil {
		log.Panic(err)
		c.JSON(http.StatusOK, gin.H{"Err": err})
	}

	if method == "DELETE" {
		c.JSON(http.StatusOK, gin.H{"status": "delete successfully"})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/edit/"+language)
	}
}
