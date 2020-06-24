package runner

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func statusResponseHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Up!",
	})
}

func startResponseHander(c *gin.Context) {
	flowUUID := c.PostForm("flow_uuid")
	if flowUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing flow_uuid"),
		})
		return
	}
	flowRepoURL := c.PostForm("flow_repo_url")
	if flowRepoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing flow_repo_url"),
		})
		return
	}
	// callbackURL defined in notifications.go
	callbackURL = c.PostForm("callback_url")
	if callbackURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing callback_url"),
		})
		return
	}

	go startFlow(flowUUID, flowRepoURL, callbackURL)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("OpenROAD flow has started!"),
	})
}
