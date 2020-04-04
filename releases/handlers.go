package releases

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type gitHubPayload struct {
	Ref string
}

func githubResponseHandler(c *gin.Context) {
	var payload gitHubPayload
	if err := c.BindJSON(&payload); err != nil {
		c.Data(http.StatusBadRequest, "text/html",
			[]byte("Couldn't parse the payload"))
		return
	}

	// Update the Docker image only if the trigger was from `master`
	if payload.Ref == "refs/heads/master" {
		go updateDockerImage()
	}

	c.Data(http.StatusOK, "text/html",
		[]byte("OpenROAD Cloud triggered to update."))
}
