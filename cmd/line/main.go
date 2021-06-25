package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whchi/stock-chatbot/pkg/setting"
	"github.com/whchi/stock-chatbot/routes"
)

func init() {
	setting.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routes.InitRouter()
	endPoint := "0.0.0.0:" + setting.ServerSetting.PORT
	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}
	server.ListenAndServe()
}
