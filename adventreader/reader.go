package adventreader

import (
	"os"
	"runtime"
	"strings"
)

func FromFile(suffix string) []byte {
	_, file, _, _ := runtime.Caller(1)
	subpath := strings.Split(file, ".")[0]
	data, err := os.ReadFile(subpath + suffix + ".txt")
	if err != nil {
		panic(err)
	}
	return data
}
