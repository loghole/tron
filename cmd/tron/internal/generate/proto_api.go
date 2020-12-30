package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"
	"golang.org/x/tools/imports"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func ProtoAPI(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate proto api")

	for _, proto := range project.Protos {
		printer.VerbosePrintf(color.Reset, "\tgenerate from file '%s': ", color.YellowString(proto.RelPath))

		if err := generateProtoc(project, proto); err != nil {
			printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		if err := generateTronOptions(project, proto); err != nil {
			printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		if err := generateTransport(project, proto); err != nil {
			printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		printer.VerbosePrintln(color.FgGreen, "OK")
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func generateProtoc(project *models.Project, proto *models.Proto) error {
	if err := helpers.Mkdir(proto.PkgTypesFile()); err != nil {
		return err
	}

	var (
		output = filepath.Join(proto.PkgDir())
		pkgMap = strings.Join(project.ProtoPkgMap, ",")
	)

	args := []string{
		"--plugin=protoc-gen-go=" + filepath.Join(project.AbsPath, "bin", "protoc-gen-go"),
		"--plugin=protoc-gen-go-grpc=" + filepath.Join(project.AbsPath, "bin", "protoc-gen-go-grpc"),
		"--plugin=protoc-gen-grpc-gateway=" + filepath.Join(project.AbsPath, "bin", "protoc-gen-grpc-gateway"),
		"--plugin=protoc-gen-openapiv2=" + filepath.Join(project.AbsPath, "bin", "protoc-gen-openapiv2"),
		fmt.Sprintf("-I%s:%s", filepath.Dir(proto.RelPath), models.ProjectPathVendorPB),
		fmt.Sprintf("--go_out=%s:%s", pkgMap, output),
		fmt.Sprintf("--go-grpc_out=%s:%s", pkgMap, output),
		fmt.Sprintf("--grpc-gateway_out=%s:%s", pkgMap, output),
		fmt.Sprintf("--openapiv2_out=%s:%s", pkgMap, output),
		"--grpc-gateway_opt", "logtostderr=true",
		"--openapiv2_opt", "fqn_for_openapi_name=true",
		proto.RelPath,
	}

	if err := execProtoc(project.AbsPath, args); err != nil {
		return simplerr.Wrap(err, "generate protoc-gen-go")
	}

	return nil
}

func generateTronOptions(project *models.Project, proto *models.Proto) error {
	file, err := os.OpenFile(filepath.Join(project.AbsPath, proto.PkgTypesFile()), os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("can't open file '%s': %w", proto.PkgTypesFile(), err)
	}

	defer helpers.Close(file)

	result, err := scanAndWriteTronOptions(file)
	if err != nil {
		return fmt.Errorf("apply tron options: %w", err)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("truncate file '%s': %w", proto.PkgTypesFile(), err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("seek file '%s': %w", proto.PkgTypesFile(), err)
	}

	if _, err := file.Write(result); err != nil {
		return fmt.Errorf("write data to file '%s': %w", proto.PkgTypesFile(), err)
	}

	return nil
}

func generateTransport(project *models.Project, proto *models.Proto) error {
	defer os.Remove(proto.PkgSwaggerFile())

	if !proto.WithImpl() {
		return nil
	}

	data := templates.TransportData{
		Proto:   proto,
		Version: project.Version,
		Swagger: `"{}"`,
	}

	if swaggerData, err := ioutil.ReadFile(proto.PkgSwaggerFile()); err == nil {
		data.Swagger = "`" + string(swaggerData) + "`"
	}

	transport, err := helpers.ExecTemplate(templates.Transport, data)
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	if err := helpers.WriteToFile(proto.PkgTronFile(), []byte(transport)); err != nil {
		return fmt.Errorf("write file '%s': %w", proto.PkgTronFile(), err)
	}

	return nil
}

func scanAndWriteTronOptions(r io.Reader) ([]byte, error) {
	buf := make([]string, 0)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), models.TronOptionsTag) {
			buf = append(buf, scanner.Text())

			continue
		}

		opts, ok := models.ParseTronOptions(scanner.Text())
		if !ok {
			continue
		}

		buf = append(buf, opts.Apply(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner err: %w", err)
	}

	result := strings.Join(buf, "\n")

	formattedBytes, err := imports.Process("", []byte(result), nil)
	if err != nil {
		return nil, fmt.Errorf("imports process '%s': %w", result, err)
	}

	return formattedBytes, nil
}

func execProtoc(wd string, args []string) error {
	stderr := bytes.NewBuffer(nil)

	cmd := exec.Command("protoc", args...)
	cmd.Dir = wd
	cmd.Stderr = stderr

	o, err := cmd.Output()
	if err != nil {
		return simplerr.Wrapf(err, "stderr: %s", stderr.String())
	}

	if len(o) > 0 {
		return simplerr.Wrapf(ErrProtocUnexpectedOutput, string(o))
	}

	return nil
}
