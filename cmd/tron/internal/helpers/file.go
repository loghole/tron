package helpers

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func WriteToFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" {
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

func Mkdir(path string) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Printf("close failed: %v", err)
	}
}
