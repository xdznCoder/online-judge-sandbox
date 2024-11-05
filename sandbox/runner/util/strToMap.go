package util

import "strings"

func StrToMap(data string) map[string]string {
	mp := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSuffix(data, "\n"), "\n") {
		idx := strings.Index(line, "=")
		mp[line[0:idx]] = line[idx+1:]
	}
	return mp
}
