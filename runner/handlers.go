package runner

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/distribution/uuid"

	"github.com/gin-gonic/gin"
)

var flowDir = "/tmp/"

func loadEnv() {
	if value, exists := os.LookupEnv("FLOW_DIR"); exists {
		flowDir = value
	}
}

func statusResponseHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Up!",
	})
}

func startResponseHander(c *gin.Context) {
	designFile, err := c.FormFile("design")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("design upload failed: %s", err.Error()),
		})
		return
	}
	sdcFile, err := c.FormFile("sdc")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("sdc upload failed: %s", err.Error()),
		})
		return
	}
	designName := c.PostForm("design_name")
	if designName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing design_name"),
		})
		return
	}
	dieArea := c.PostForm("die_area")
	if dieArea == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing die_area"),
		})
		return
	}
	coreArea := c.PostForm("core_area")
	if coreArea == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("missing core_area"),
		})
		return
	}

	// Passed validation. Proceed ..
	loadEnv()

	// Create a directory for this run
	runDir := uuid.Generate()
	_ = os.MkdirAll(path.Join(flowDir, runDir.String()), os.ModePerm)

	// Save files to this directory
	filename := path.Join(flowDir, runDir.String(), filepath.Base(designFile.Filename))
	if err := c.SaveUploadedFile(designFile, filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't process design file: %s", err.Error()),
		})
		return
	}
	filename = path.Join(flowDir, runDir.String(), filepath.Base(sdcFile.Filename))
	if err = c.SaveUploadedFile(sdcFile, filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't process sdc file: %s", err.Error()),
		})
		return
	}

	// Create flow parameters
	filename = path.Join(flowDir, runDir.String(), "config.mk")
	params := flowParams{
		DesignName:   designName,
		VerilogFiles: path.Join("/cloud", designFile.Filename),
		SdcFile:      path.Join("/cloud", sdcFile.Filename),
		DieArea:      dieArea,
		CoreArea:     coreArea,
	}
	err = createConfigFile(filename, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't create flow config file: %s", err.Error()),
		})
		return
	}

	// Compress flow directory
	compressedFlow := compressFlow(path.Join(flowDir, runDir.String()))

	// Upload flow to bucket
	_, urlString, err := createBucket(runDir.String(), compressedFlow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("couldn't reserve storage for the flow: %s", err.Error()),
		})
		return
	}

	// Run flow
	go runFlow(path.Join(flowDir, runDir.String()), runDir.String())

	c.JSON(http.StatusOK, gin.H{
		"result_url": urlString,
		"message":    fmt.Sprintf("OpenROAD flow has started!"),
	})
}
