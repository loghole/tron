package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"go/format"
	"os"
	"regexp"

	"golang.org/x/tools/imports"
)

var ErrModuleNotFound = errors.New("module not found")

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
		if m := rgxp.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan failed: %w", err)
	}

	return "", ErrModuleNotFound
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
