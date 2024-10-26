package config

var (
	DefaultLangType = 4

	DefaultCPUSetCPUs   = "1"
	DefaultMemory       = int64(512 * 1024 * 1024) // 512 MB
	DefaultSharedMemory = int64(64)                // 64 B
)

const (
	c = iota + 1
	cpp
	java
	python
)

func LangStringToCode(code string) int {
	switch code {
	case "c":
		return c
	case "cpp":
		return cpp
	case "java":
		return java
	default:
		return 0
	}
}

func EvalLang() []string {
	Eval := make([]string, DefaultLangType)

	Eval[c] = "gcc /tests/data/a.c -o /tests/data/a.out 2>&1"
	Eval[cpp] = "g++ -w -O2 /tests/data/a.cpp -o /tests/data/a.out 2>&1"
	Eval[java] = "javac -d /tests/data/ -cp /test/data/ /tests/data/Main.java 2>&1"

	return Eval
}
