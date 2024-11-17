package executor

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"log"
	"regexp"
	"xdznOJ/sandbox/container/config"
)

// Compile is for static language to be compiled
func Compile(ctx context.Context, cli *client.Client, image string,
	language string, path string) (int, string) {
	eval := config.EvalLang()
	if config.LangStringToCode(language) == 0 {
		log.Fatal("language not supported")
		return -1, ""
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", eval[config.LangStringToCode(language)]},
	},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: path,
					Target: "/tests/data",
				},
			},
		}, nil, nil, "")

	if err != nil {
		return 2, err.Error()
	}

	out, err := runContainer(ctx, cli, resp.ID)
	if err != nil {
		return 2, err.Error()
	}

	r, _ := regexp.Compile("error")
	if r.MatchString(out) {
		return 0, out
	}

	return 1, ""
}
