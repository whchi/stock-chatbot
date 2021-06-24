package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whchi/stock-chatbot/routes/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/line/callback", api.LineEventHandler)
	r.GET("/", func(c *gin.Context){
		c.String(http.StatusOK, "Hello World!")
	})
	return r
}
