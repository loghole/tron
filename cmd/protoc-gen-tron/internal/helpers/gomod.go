package helpers

import (
	"bufio"
	"fmt"
	"go/format"
	"os"
	"regexp"

	"golang.org/x/tools/imports"
)

func ModuleFromGoMod() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}

	rgxp := regexp.MustCompile(`^module (.+)$`)

	defer Close(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}

		if m := rgxp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	return "", fmt.Errorf("module not found")
}

func FormatGoFile(data string) ([]byte, error) {
	formattedBytes, err := format.Source([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("failed to format process: %w", err)
	}

	formattedBytes, err = imports.Process("", formattedBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to imports process: %w", err)
	}

	return formattedBytes, nil
}
