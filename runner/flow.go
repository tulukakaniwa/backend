package runner

import (
	"context"
	"log"
	"os"

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

func runFlow(flowDir string, bucketName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
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
	}, nil, nil, bucketName)
	if err != nil {
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	waitChan, errChan := cli.ContainerWait(ctx, bucketName, container.WaitConditionNextExit)

	select {
	case e := <-errChan:
		log.Printf("Flow %s falied: %s", bucketName, e)
	case <-waitChan:
		log.Printf("Flow %s finished!", bucketName)

		// Remove container
		cli.ContainerRemove(ctx, bucketName, types.ContainerRemoveOptions{})

		// Upload flow directory
		compressedFlow := compressFlow(flowDir)
		err := uploadFlowDir(bucketName, compressedFlow)
		if err != nil {
			log.Printf("Flow %s failed to upload: %s", bucketName, err.Error())
		} else {
			log.Printf("Flow %s uploaded!", bucketName)
		}
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
