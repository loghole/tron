package helpers

import (
	"io"
	"os"
	"path/filepath"
)

func WriteToFile(path string, data []byte) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer Close(file)

	if _, err := file.Write(data); err != nil {
		return err
	}

	return err
}

func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}
