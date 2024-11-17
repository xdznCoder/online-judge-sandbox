package executor

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
)

func runContainer(ctx context.Context, cli *client.Client, ID string) (string, error) {
	if err := cli.ContainerStart(ctx, ID, container.StartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	// capture container logs
	out, err := cli.ContainerLogs(ctx, ID, container.LogsOptions{ShowStdout: true})
	defer func(out io.ReadCloser) {
		err := out.Close()
		if err != nil {
			log.Println(err)
		}
	}(out)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(out)
	if err != nil {
		return "", err
	}
	res := buffer.String()
	fmt.Println(res)

	// close the container
	err = cli.ContainerRemove(ctx, ID, container.RemoveOptions{})
	if err != nil {
		return "", err
	}
	return res, nil
}
