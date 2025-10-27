package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UpdateEnvValue(filename, key, value string) error {
	// Open the .env file for reading
	file, err := os.OpenFile(filename, os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			lines = append(lines, line)
			continue
		}

		// Split key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
