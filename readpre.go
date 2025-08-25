package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

func readUntilProxies(filePath string) ([]byte, error) {
	var buffer bytes.Buffer

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "proxies:") {
			break
		}
		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
