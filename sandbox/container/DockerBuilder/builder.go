package DockerBuilder

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"xdznOJ/sandbox/container/config"
)

// getContext package the './docker' into tar file
// so the dockerfile can be executed in the image
func getContext(path string) io.Reader {
	ctx, _ := archive.TarWithOptions(path, &archive.TarOptions{})
	return ctx
}

// newImage create new container for the code to be
// compiled and run.
func newImage(cli *client.Client, imageName string) error {
	ctx := context.Background()
	if cli == nil {
		cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	}

	_, basePath, _, _ := runtime.Caller(0)
	dockerCtx := getContext(filepath.Dir(basePath) + "/docker")

	opt := types.ImageBuildOptions{
		Tags:       []string{imageName},
		CPUSetCPUs: config.DefaultCPUSetCPUs,
		Memory:     config.DefaultMemory,
		ShmSize:    config.DefaultSharedMemory,
		Dockerfile: "Dockerfile",
	}

	resp, err := cli.ImageBuild(ctx, dockerCtx, opt)
	if err != nil {
		return errors.New(fmt.Sprintf("build docker image failed , err : %v", err))
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("read docker image output failed , err : %v", err))
	}

	return nil
}
