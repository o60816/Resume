// main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

type Project struct {
	Id          int    `json:"id"`
	WorkId      int    `json:"workid"`
	ProjectName string `json:"ProjectName"`
	Tech        string `json:"tech"`
}

type Work struct {
	Id          int       `json:"id"`
	Period      string    `json:"period"`
	Logo        string    `json:"logo"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	Content     string    `json:"content"`
	ProjectList []Project `json:"ProjectList"`
}

func main() {
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/Resume?parseTime=true")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/image", "./image")

	router.GET("/", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, period, logo, company, position, content FROM work")
		defer rows.Close()

		if err != nil {
			log.Fatalln(err)
		}

		workList := make([]Work, 0)

		for rows.Next() {
			var work Work
			rows.Scan(&work.Id, &work.Period, &work.Logo, &work.Company, &work.Position, &work.Content)
			projectList := make([]Project, 0)
			rows2, err := db.Query("SELECT id, workId, projectName, tech FROM project WHERE workId=?", work.Id)
			defer rows2.Close()
			if err != nil {
				log.Fatalln(err)
			}

			for rows2.Next() {
				var project Project
				rows2.Scan(&project.Id, &project.WorkId, &project.ProjectName, &project.Tech)
				projectList = append(projectList, project)
			}

			if err = rows2.Err(); err != nil {
				log.Fatalln(err)
			}
			work.ProjectList = projectList
			workList = append(workList, work)
		}

		if err = rows.Err(); err != nil {
			log.Fatalln(err)
		}

		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"works": workList,
			},
		)
	})

	router.GET("/work/create", func(c *gin.Context) {
		var work Work
		c.HTML(
			http.StatusOK,
			"addWork.html",
			gin.H{
				"method":  "POST",
				"work":    work,
				"btnName": "新增",
			},
		)
	})

	router.GET("/work/update/:workId", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, period, logo, company, position, content FROM work WHERE id=?", c.Param("workId"))

		if err != nil {
			log.Fatalln(err)
		}

		rows.Next()
		var work Work
		rows.Scan(&work.Id, &work.Period, &work.Logo, &work.Company, &work.Position, &work.Content)

		c.HTML(
			http.StatusOK,
			"addWork.html",
			gin.H{
				"method":  "PATCH",
				"work":    work,
				"btnName": "更新",
			},
		)
	})

	var query string
	router.POST("/work/:workId", func(c *gin.Context) {
		if c.PostForm("_method") == "POST" {
			query = fmt.Sprintf("INSERT INTO work(period, logo, company, position, content) VALUES('%s','%s','%s','%s','%s')", c.PostForm("period"), c.PostForm("logo"), c.PostForm("company"), c.PostForm("position"), c.PostForm("content"))
		} else {
			query = fmt.Sprintf("UPDATE work SET period='%s',logo='%s',company='%s',position='%s',content='%s' WHERE id='%s'", c.PostForm("period"), c.PostForm("logo"), c.PostForm("company"), c.PostForm("position"), c.PostForm("content"), c.Param("workId"))
		}

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.DELETE("/work/:workId", func(c *gin.Context) {
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Err": err})
			return
		}
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p) // re-throw panic after Rollback
			}
		}()
		if _, err = tx.Exec("DELETE FROM work WHERE id=?", c.Param("workId")); err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"Err": err})
			return
		}
		if _, err = tx.Exec("DELETE FROM project WHERE workId=?", c.Param("workId")); err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"Err": err})
			return
		}

		err = tx.Commit()
		c.JSON(http.StatusOK, gin.H{"status": "刪除成功"})
	})

	router.GET("/project/create/:workId", func(c *gin.Context) {
		var project Project
		project.WorkId, err = strconv.Atoi(c.Param("workId"))
		c.HTML(
			http.StatusOK,
			"addProject.html",
			gin.H{
				"method":  "POST",
				"project": project,
				"btnName": "新增",
			},
		)
	})

	router.GET("/project/update/:projectId", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, workId, projectName, tech FROM project WHERE id=?", c.Param("projectId"))
		if err != nil {
			log.Fatalln(err)
		}

		rows.Next()
		var project Project
		rows.Scan(&project.Id, &project.WorkId, &project.ProjectName, &project.Tech)

		c.HTML(
			http.StatusOK,
			"addProject.html",
			gin.H{
				"method":  "PATCH",
				"project": project,
				"btnName": "更新",
			},
		)
	})

	router.POST("/project/:projectId", func(c *gin.Context) {
		if c.PostForm("_method") == "POST" {
			query = fmt.Sprintf("INSERT INTO project(workId, projectName, tech) VALUES('%s','%s','%s')", c.PostForm("workId"), c.PostForm("projectName"), c.PostForm("tech"))
		} else {
			query = fmt.Sprintf("UPDATE project SET workId='%s',projectName='%s',tech='%s' WHERE id='%s'", c.PostForm("workId"), c.PostForm("projectName"), c.PostForm("tech"), c.Param("projectId"))
		}

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.DELETE("/project/:projectId", func(c *gin.Context) {

		if _, err := db.Exec("Delete from project WHERE id=?", c.Param("projectId")); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusOK, gin.H{"Err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "刪除成功"})
	})

	router.Run(":8080")
}
