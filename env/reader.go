package env

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func readOS() map[string]string {
	var kv = make(map[string]string)
	for _, keyValue := range os.Environ() {
		s := strings.SplitN(keyValue, "=", 2)
		kv[s[0]] = s[1]
	}

	return kv
}

func readFile(fileName string) (map[string]string, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var kv = make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		i := strings.SplitN(line, "=", 2)
		if len(i) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		key := strings.TrimSpace(i[0])
		kv[key] = strings.TrimSpace(i[1])
	}

	fmt.Printf("Values: \n %#v\n", kv)

	return kv, nil
}
