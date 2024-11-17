package runner

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"os"
	"path/filepath"
	"xdznOJ/sandbox/container/executor"
)

// toFile transforms the information in data to file in container
func (data *CodeData) toFile() error {
	_, err := os.Stat(data.Path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(data.Path, 0755)
		if errDir != nil {
			return err
		}
	}

	// create the code file
	file, err := os.Create(filepath.Join(data.Path, data.Filename))
	if err != nil {
		return err
	}

	// write the code in file
	_, err = file.WriteString(data.Code)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	for i := 1; i <= data.TestCount; i++ {
		In, err := os.Create(filepath.Join(data.Path, fmt.Sprintf("in%v.txt", i)))
		if err != nil {
			return err
		}
		Out, err := os.Create(filepath.Join(data.Path, fmt.Sprintf("out%v.txt", i)))
		if err != nil {
			return err
		}

		_, err = In.WriteString(data.InputData[i-1])
		if err != nil {
			return err
		}
		_, err = Out.WriteString(data.OutputData[i-1])
		if err != nil {
			return err
		}

		_ = In.Close()
		_ = Out.Close()
	}
	return nil
}

// Run is the main function of this sandbox
// with the given code and limit, Run can start a docker container,
// execute the code and get the result.
func (data *CodeData) Run(cli *client.Client) CodeResult {

	// Create new cli
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx := context.Background()

	if err != nil {
		log.Println(err)
		return CodeResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	cwd, err := os.Getwd()
	data.Path = filepath.Join(cwd, "temp", data.ID)
	data.Image = "runner:latest"

	// load the information in data to file
	err = data.toFile()
	if err != nil {
		log.Println(err)
		return CodeResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	status, msg := executor.Compile(ctx, cli, data.Image, data.Language, data.Path)
	if status != 1 {
		log.Println(msg)
		if status == 0 {
			return CodeResult{
				ID:      data.ID,
				Status:  "COMPILATION ERROR",
				Message: msg,
			}
		}

		return CodeResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	err = executor.Execute(ctx, cli,
		data.Image,
		data.Language,
		data.Path,
		data.TestCount,
		data.TimeLimit,
		data.MemLimit)
	if err != nil {
		log.Println(err)
		return CodeResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	res := CodeResult{
		ID:     data.ID,
		Time:   make([]float64, data.TestCount),
		Memory: make([]float64, data.TestCount),
		Result: make([]string, data.TestCount),
		Error:  make([]string, data.TestCount),
		Status: "OK",
	}

	err = parseOutput(data.Path, data, &res)
	if err != nil {
		log.Println(err)
		return CodeResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	err = os.RemoveAll(data.Path)
	if err != nil {
		log.Println(err)
	}
	return res
}
