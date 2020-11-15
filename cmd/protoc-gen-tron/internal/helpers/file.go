package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteToFile(path string, data []byte) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("can't mkdir: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't create file: %w", err)
	}

	defer Close(file)

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("can't write to file: %w", err)
	}

	return err
}

func Close(closer io.Closer) {
	_ = closer.Close()
}
