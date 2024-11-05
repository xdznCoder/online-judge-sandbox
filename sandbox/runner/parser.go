package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"xdznOJ/sandbox/runner/util"
)

// parseOutput parses the output data and
// show the different judge results caused by the code
func parseOutput(path string, data *CodeData, res *CodeResult) error {
	for i := 1; i <= data.TestCount; i++ {
		data, err := os.ReadFile(filepath.Join(path, fmt.Sprintf("diff%v.txt", i)))
		diff := strings.TrimSpace(string(data))
		if err != nil {
			return err
		}

		// parse the running result from string to map
		data, err = os.ReadFile(filepath.Join(path, fmt.Sprintf("stats%d.txt", i)))
		stats := string(data)
		stMap := util.StrToMap(stats)
		returnValue, _ := strconv.Atoi(stMap["returnvalue"])
		termination := strings.TrimSpace(stMap["terminationreason"])
		time, _ := strconv.ParseFloat(strings.TrimSuffix(stMap["cputime"], "s"), 64)
		mem, _ := strconv.ParseFloat(strings.TrimSuffix(stMap["memory"], "B"), 64)

		res.Time[i-1] = time
		res.Memory[i-1] = mem

		if returnValue != 0 {
			switch returnValue {
			case 9, 15:
				{
					if termination == "cputime" {
						res.Result[i-1] = "TIME LIMIT EXCEEDED"
					} else if termination == "memory" {
						res.Result[i-1] = "MEMORY LIMIT EXCEEDED"
					} else {
						res.Result[i-1] = "ILLEGAL INSTRUCTIONS"
					}
				}
			default:
				res.Result[i-1] = "RUNTIME ERROR"
			}
		} else {
			if diff == "" {
				res.Result[i-1] = "ACCEPT"
			} else {
				res.Result[i-1] = "WRONG ANSWER"
			}
		}
	}
	return nil
}
