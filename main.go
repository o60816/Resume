// main.go

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

type Project struct {
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

// var projectList = []Project{
// 	Project{ProjectName: "手機遠端查看並控制攝影機", Tech: "QT  framework、 c++、javascript、ffmpeg、darwin  streaming server、rtsp、onvif、muti-thread"},
// 	Project{ProjectName: "手機遠端即時查看設備與生產資訊", Tech: "OPC  Server,web  service、soap、muti-thread"},
// 	Project{ProjectName: "企業協作平台", Tech: "Soap、muti-thread、傳送文字、語音、圖片"},
// }

// var projectList2 = []Project{
// 	Project{ProjectName: "移植測試framework從Perl到Python", Tech: "Perl、Python"},
// 	Project{ProjectName: "完成多個測試程式的撰寫", Tech: "Python、Selenium、架設server"},
// 	Project{ProjectName: "Fake browser", Tech: "TLS、HTTP、HTTP2、SOCKS、c++、windows api、閱讀chrome與firefox source code、windbg"},
// }

// var workList = []work{
// 	work{Period: "Oct 2013-Apr 2016", Logo: "http://www.ccichain.net/images/up_images/2014420114640.bmp", Company: "海南金海漿紙業有限公司", Position: "軟體開發工程師", Content: "開發android、ios兩大移動平台的公司入口應用程式並將桌面平台的系統移植、整合到其中", ProjectList: projectList},
// 	work{Period: "May 2016~Present", Logo: "https://www.trendmicro.com/content/dam/trendmicro/global/en/global/logo/trend-micro-mobile.png", Company: "趨勢科技", Position: "軟體測試工程師", Content: "撰寫測試計畫、測試程式、測試工具、維護測試framework、架設測試伺服器、分析與解決客戶問題", ProjectList: projectList2},
// }

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
			rows2, err := db.Query("SELECT projectName, tech FROM project WHERE workId=?", work.Id)
			defer rows2.Close()
			if err != nil {
				log.Fatalln(err)
			}

			for rows2.Next() {
				var project Project
				rows2.Scan(&project.ProjectName, &project.Tech)
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

	router.Run(":8080")
}
