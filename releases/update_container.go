package releases

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

func cloneRepos() (cloned bool) {
	cloned = true
	yosysRepo := "https://github.com/The-OpenROAD-Project/yosys"
	cmd := exec.Command("git", "clone", yosysRepo, "releases/build/yosys")

	path, err := os.Getwd()
	if err != nil {
		log.Printf("%s", err.Error())
		cloned = false
		return
	}
	cmd.Dir = path

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
		cloned = false
		return
	}

	// TODO: clone openroad

	return
}
func deleteRepos() (deleted bool) {
	deleted = true

	// TODO: delete yosys and openroad

	return
}

func updateDockerImage() (updated bool) {
	updated = true

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		updated = false
		return
	}
	cli.NegotiateAPIVersion(ctx)

	// clone latest @master of repos
	cloned := cloneRepos()
	if !cloned {
		updated = false
		return
	}
	defer deleteRepos()

	tar := new(archivex.TarFile)
	tar.Create("/tmp/openroad.tar")
	tar.AddAll("releases/build", true)
	tar.Close()

	dockerBuildContext, err := os.Open("/tmp/openroad.tar")
	if err != nil {
		updated = false
		return
	}
	defer dockerBuildContext.Close()

	dockerBuildOptions := types.ImageBuildOptions{
		Dockerfile:     "build/Dockerfile",
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{"edaac/openroad"},
	}

	buildResponse, err := cli.ImageBuild(ctx, dockerBuildContext, dockerBuildOptions)
	if err != nil {
		updated = false
		fmt.Printf("%s", err.Error())
		return
	}
	defer buildResponse.Body.Close()

	// writeToLog(buildResponse.Body)

	return
}

func writeToLog(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Println(string(n))
	}
	return nil
}
