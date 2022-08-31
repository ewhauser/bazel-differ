package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteTargetsFile(targets map[string]bool, output string) {
	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()

	datawriter := bufio.NewWriter(file)
	defer datawriter.Flush()

	for k := range targets {
		_, _ = datawriter.WriteString(k + "\n")
	}
}

func ReadHashFile(filename string) (map[string]string, error) {
	x := map[string]string{}
	startingContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil
	}
	err = json.Unmarshal(startingContent, &x)
	if err != nil {
		return nil, fmt.Errorf("error serializing hashes for file: %s", filename)
	}
	return x, nil
}

func WriteHashFile(filename string, data interface{}) (string, error) {
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return "", err
	}

	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filename, out.Bytes(), 0644)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
