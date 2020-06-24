package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

type flowParams struct {
	DesignName   string
	VerilogFiles string
	SdcFile      string
	DieArea      string
	CoreArea     string
}

var flowDir = "/tmp/"

func loadEnv() {
	if value, exists := os.LookupEnv("FLOW_DIR"); exists {
		flowDir = value
	}
}

func compressFlow(runDir string) string {
	compressedFlow := runDir + ".tar"
	tar := new(archivex.TarFile)
	tar.Create(compressedFlow)
	tar.AddAll(runDir, true)
	tar.Close()
	return compressedFlow
}

func cloneRepo(repoURL string, cloneDir string) (cloned bool) {
	cloned = true

	cmd := exec.Command("git", "clone", repoURL, cloneDir)

	err := cmd.Run()

	if err != nil {
		log.Printf("Flow %s failed to clone from git repository %s: %s", cloneDir, repoURL, err.Error())
		cloned = false
		return
	}

	return
}

func createConfigFile(filename string, params flowParams) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString("export DESIGN_NAME = " + params.DesignName + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export PLATFORM = nangate45\n\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export VERILOG_FILES = " + params.VerilogFiles + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export SDC_FILE = " + params.SdcFile + "\n\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export DIE_AREA = " + params.DieArea + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export CORE_AREA = " + params.CoreArea + "\n"); err != nil {
		return err
	}
	f.Sync()

	return nil
}

func runFlow(flowDir string, flowUUID string, callbackURL string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Flow %s failed to connect to Docker engine: %s", flowUUID, err.Error())
		return
	}
	cli.NegotiateAPIVersion(ctx)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "openroadcloud/flow",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: flowDir,
				Target: "/cloud",
			},
		},
	}, nil, nil, flowUUID)
	if err != nil {
		log.Printf("Flow %s container failed to create: %s", flowUUID, err.Error())
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Flow %s container failed to start: %s", flowUUID, err.Error())
		return
	}

	waitChan, errChan := cli.ContainerWait(ctx, flowUUID, container.WaitConditionNextExit)

	select {
	case e := <-errChan:
		errorMessage := fmt.Sprintf("Flow %s falied: %s", flowUUID, e)
		log.Println(errorMessage)
		go notifyFlowFail(flowUUID, errorMessage)

	case <-waitChan:
		log.Printf("Flow %s finished!", flowUUID)

		// Remove container
		cli.ContainerRemove(ctx, flowUUID, types.ContainerRemoveOptions{})

		// Upload flow directory
		compressedFlow := compressFlow(path.Join(flowDir, "openroad-flow"))
		_, urlString, err := createBucket(flowUUID, compressedFlow)
		if err != nil {
			errorMessage := fmt.Sprintf("Flow %s failed to upload: %s", flowUUID, err.Error())
			log.Println(errorMessage)
			go notifyFlowFail(flowUUID, errorMessage)
		} else {
			log.Printf("Flow %s uploaded at %s!", flowUUID, urlString)
			go notifyFlowSuccess(flowUUID, urlString)
		}
	}

	// Remove flow directory
	os.RemoveAll(flowDir)
}

func startFlow(flowUUID string, repoURL string, callbackURL string) {
	go notifyFlowStart(flowUUID)

	loadEnv()

	// Create a directory for this run
	_ = os.MkdirAll(path.Join(flowDir, flowUUID), os.ModePerm)

	// clone repo
	cloned := cloneRepo(repoURL, path.Join(flowDir, flowUUID, "repo"))
	if !cloned {
		go notifyFlowFail(flowUUID, fmt.Sprintf("Cannot clone the repository from %s", repoURL))
		return
	}

	// Read openroad.yml
	conf, err := readFlowConf(path.Join(flowDir, flowUUID, "repo", "openroad.yml"))
	if err != nil {
		errorMessage := fmt.Sprintf("Couldn't read openroad.yml for flow %s: %s", flowUUID, err.Error())
		log.Println(errorMessage)
		go notifyFlowFail(flowUUID, errorMessage)

	}
	conf.DesignFiles = designFilesMap(conf.DesignFiles, func(f string) string {
		return "/cloud/repo/" + f
	})

	// Create flow parameters
	filename := path.Join(flowDir, flowUUID, "config.mk")
	params := flowParams{
		DesignName:   conf.DesignName,
		VerilogFiles: strings.Join(conf.DesignFiles, " "),
		SdcFile:      path.Join("/cloud/repo", conf.SdcFile),
		DieArea:      conf.DieArea,
		CoreArea:     conf.CoreArea,
	}
	err = createConfigFile(filename, params)
	if err != nil {
		errorMessage := fmt.Sprintf("Couldn't create flow config file for flow %s: %s", flowUUID, err.Error())
		log.Println(errorMessage)
		go notifyFlowFail(flowUUID, errorMessage)
		return
	}

	runFlow(path.Join(flowDir, flowUUID), flowUUID, callbackURL)
}

func designFilesMap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
