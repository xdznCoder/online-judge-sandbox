package executor

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func Execute(ctx context.Context, cli *client.Client, image string, language string,
	path string, count int, timeLimit int, memLimit int64) error {
	var exec string
	switch language {
	case "c":
		exec = fmt.Sprintf("/tests/evaluate %v %v %v", count, timeLimit, memLimit)
	case "cpp":
		exec = fmt.Sprintf("/tests/evaluate %v %v %v", count, timeLimit, memLimit)
	case "java":
		program := "java -cp /tests/data Main"
		exec = fmt.Sprintf("/tests/evaluate %v %v %v %v", count, timeLimit, memLimit, program)
	case "python":
		program := "python3 /tests/data/a.py 2>&1"
		exec = fmt.Sprintf("/tests/evaluate %v %v %v %v", count, timeLimit, memLimit, program)
	}

	// Creating namespace for container mode failed: Invalid argument.
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", exec},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/tests/data",
			},
			{
				Type:   mount.TypeBind,
				Source: "/sys/fs/cgroup",
				Target: "/sys/fs/cgroup",
			},
		},
		Privileged: true,
		CapDrop:    []string{"all"},
	}, nil, nil, "")

	if err != nil {
		return err
	}

	_, err = runContainer(ctx, cli, resp.ID)
	return err
}
