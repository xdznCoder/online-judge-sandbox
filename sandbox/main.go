package main

import (
	"github.com/docker/docker/client"
	"log"
	"os"
	"path/filepath"
	"xdznOJ/sandbox/container/Build"
	"xdznOJ/sandbox/runner"
)

var cli *client.Client

func init() {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := Build.NewImage(cli, "runner:latest")

	if err != nil {
		log.Fatal(err)
	}
}

func Evaluate(code string, lang string, file string, in []string, out []string,
	count int, time int, mem int64) runner.CodeResult {

	data := runner.CodeData{
		ID:         "1",
		Language:   lang,
		Code:       code,
		Filename:   file,
		InputData:  in,
		OutputData: out,
		TestCount:  count,
		TimeLimit:  time,
		MemLimit:   mem,
	}
	return data.Run(cli)
}

func test(code string, lang string, filename string) {
	in := []string{"2", "4", "5"}
	out := []string{"4", "16", "25"}
	timeLimit := 2                       // Time in seconds
	memLimit := int64(500 * 1024 * 1024) // Memory in bytes

	res := Evaluate(code, lang, filename, in, out,
		len(in), timeLimit, memLimit)
	log.Print(res)
}

func main() {
	path, _ := os.Getwd()
	cfile, err := os.ReadFile(filepath.Join(path, "sandbox/test/test.c"))
	if err != nil {
		panic(err)
	}
	code := string(cfile)
	test(code, "c", "test.c")
}
