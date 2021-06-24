package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/whchi/stock-chatbot/routes/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/line/callback", api.LineEventHandler)

	return r
}
