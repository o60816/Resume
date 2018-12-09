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
	Logo        string    `json:"Logo"`
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

var projectList2 = []project{
	project{ProjectName: "移植測試framework從Perl到Python", Tech: "Perl、Python"},
	project{ProjectName: "完成多個測試程式的撰寫", Tech: "Python、Selenium、架設server"},
	project{ProjectName: "Fake browser", Tech: "熟悉網路傳輸協定(TLS,HTTP,HTTP2,SOCKS)、c++、windows api、閱讀chrome與firefox source code、windbg。"},
}

var workList = []work{
	work{Period: "Oct 2013-Apr 2016", Logo: "http://www.ccichain.net/images/up_images/2014420114640.bmp", Company: "海南金海漿紙業有限公司", Position: "軟體開發工程師", Content: "開發android、ios兩大移動平台的公司入口應用程式並將桌面平台的系統移植、整合到其中。", ProjectList: projectList},
	work{Period: "May 2016~Present", Logo: "https://www.trendmicro.com/content/dam/trendmicro/global/en/global/logo/trend-micro-mobile.png", Company: "趨勢科技", Position: "軟體測試工程師", Content: "撰寫測試計畫、測試程式、測試工具，維護測試framework，架設測試伺服器，分析與解決客戶問題。", ProjectList: projectList2},
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
