package runner

import (
	"github.com/gin-gonic/gin"
)

// InitializeRoutes initializes the router with handlers
func InitializeRoutes(router *gin.Engine) {
	router.GET("/status", statusResponseHandler)
	router.POST("/start", startResponseHander)
}
