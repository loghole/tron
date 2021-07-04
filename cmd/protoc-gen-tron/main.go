package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"regexp"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/generator"
)

var ErrModuleNotFound = errors.New("project module does not exists")

func main() {
	opt := protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}

	module, err := parseGoMod()
	if err != nil {
		log.Fatal(err)
	}

	gen := generator.NewGenerator(module)

	opt.Run(gen.Generate)
}

func parseGoMod() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	reg := regexp.MustCompile(`^module (.+)$`)

	for scanner.Scan() {
		if m := reg.FindStringSubmatch(scanner.Text()); len(m) > 1 {
			return m[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", ErrModuleNotFound
}
