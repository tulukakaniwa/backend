package runner

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func statusResponseHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Up!",
	})
}

func startResponseHander(c *gin.Context) {
	// TODO: implement the flow runner
	c.JSON(http.StatusOK, gin.H{
		"status": "Not implemented yet",
	})
}
