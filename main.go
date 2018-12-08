// main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

type project struct {
	ProjectName string `json:"ProjectName"`
	Tech        string `json:"tech"`
}

type work struct {
	Period      string    `json:"period"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	Content     string    `json:"content"`
	ProjectList []project `json:"ProjectList"`
}

var projectList = []project{
	project{ProjectName: "手機遠端查看並控制攝影機", Tech: "QT  framework、 c++、javascript、ffmpeg、darwin  streaming server、rtsp、onvif、muti-thread"},
	project{ProjectName: "手機遠端即時查看設備與生產資訊", Tech: "OPC  Server,web  service、soap、muti-thread"},
	project{ProjectName: "企業協作平台", Tech: "Soap、muti-thread、傳送文字、語音、圖片"},
}

var workList = []work{
	work{Period: "2013/10~2016/4", Company: "海南金海漿紙業有限公司", Position: "軟體開發工程師", Content: "工作內容:\n開發android、ios兩大移動平台的公司入口應用程式並將桌面平台的系統移植、整合到其中。", ProjectList: projectList},
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/image", "./image")

	router.GET("/", func(c *gin.Context) {
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
