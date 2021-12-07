package adventreader

import (
    "bufio"
    "bytes"
	"os"
	"runtime"
	"strings"
)

func FromFile(filename string, offset ...int) []byte {
    skip := 1
    if len(offset) > 0 {
        skip = offset[0]
    }
	_, file, _, _ := runtime.Caller(skip)
    pathTokens := strings.Split(file, "/")
    dir := strings.Join(pathTokens[:len(pathTokens)-1], "/")
	data, err := os.ReadFile(dir + "/" + filename)
	if err != nil {
		panic(err)
	}
	return data
}

func LinesFromFile(filename string) []string {
    input := bytes.NewBuffer(FromFile(filename, 2))
    scanner := bufio.NewScanner(input)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines
}


