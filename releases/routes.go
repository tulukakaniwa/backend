package releases

import (
	"github.com/gin-gonic/gin"
)

// InitializeRoutes initializes the router with handlers
func InitializeRoutes(router *gin.Engine) {
	v1 := router.Group("/releases/v1")
	{
		v1.POST("/update", githubResponseHandler)
	}
}
