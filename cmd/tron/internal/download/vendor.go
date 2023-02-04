package download

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type VendorPB struct {
	project *models.Project
	printer stdout.Printer

	replacer map[string]string
	exists   map[string]struct{}
	imports  []string
}

func NewVendor(project *models.Project, printer stdout.Printer) *VendorPB {
	return &VendorPB{
		project: project,
		printer: printer,

		replacer: map[string]string{
			"google/type":     "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/api":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/rpc":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/protobuf": "https://raw.githubusercontent.com/google/protobuf/master/src/",
		},
		exists:  make(map[string]struct{}),
		imports: make([]string, 0),
	}
}

func (v *VendorPB) Download() error {
	v.printer.VerbosePrintln(color.FgMagenta, "Vendor proto imports")

	if err := v.copyProjectFiles(); err != nil {
		return fmt.Errorf("copy project files: %w", err)
	}

	if err := v.scanFiles(); err != nil {
		return fmt.Errorf("scan files: %w", err)
	}

	if err := v.downloadFiles(); err != nil {
		return fmt.Errorf("download files: %w", err)
	}

	v.printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func (v *VendorPB) copyProjectFiles() error {
	for _, file := range v.project.ProtoFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read file '%s': %w", file, err)
		}

		path := filepath.Join(
			v.project.AbsPath,
			models.ProjectPathVendorPB,
			v.project.Name,
			strings.TrimPrefix(file, v.project.AbsPath),
		)

		if err := helpers.WriteToFile(path, data); err != nil {
			return fmt.Errorf("write file '%s': %w", file, err)
		}
	}

	return nil
}

func (v *VendorPB) scanFiles() error {
	for _, file := range v.project.ProtoFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read file '%s': %w", file, err)
		}

		if err := v.findImports(data); err != nil {
			return fmt.Errorf("find imports in file '%s': %w", file, err)
		}
	}

	return nil
}

func (v *VendorPB) downloadFiles() (err error) {
	var val string

	for len(v.imports) > 0 {
		val, v.imports = v.imports[0], v.imports[1:]

		if strings.HasPrefix(val, v.project.Name) {
			continue
		}

		v.printer.VerbosePrintf(color.Reset, "\tvendor '%s': ", color.YellowString(val))

		if err := v.curlProto(val); err != nil {
			v.printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		v.printer.VerbosePrintln(color.FgGreen, "OK")
	}

	return nil
}

func (v *VendorPB) curlProto(name string) error {
	link, ok := v.importLink(name)
	if !ok {
		return fmt.Errorf("'%s': %w", name, ErrBadImport)
	}

	resp, err := http.Get(link) //nolint:gosec,bodyclose,noctx //body is closed
	if err != nil {
		return err
	}

	defer helpers.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("link: %s code: %d: %w", link, resp.StatusCode, ErrReqFailed)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := v.findImports(data); err != nil {
		return err
	}

	return helpers.WriteToFile(filepath.Join(v.project.AbsPath, models.ProjectPathVendorPB, name), data)
}

func (v *VendorPB) findImports(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := models.ImportRegexp.FindStringSubmatch(scanner.Text())
		if len(m) == 0 {
			continue
		}

		if _, ok := v.exists[m[1]]; ok {
			continue
		}

		v.exists[m[1]], v.imports = struct{}{}, append(v.imports, m[1])
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (v *VendorPB) importLink(link string) (string, bool) {
	const (
		github      = "github.com"
		githubRaw   = "https://raw.githubusercontent.com"
		importParts = 3
	)

	parts := strings.Split(link, "/")

	for repoPartsCount := len(parts); repoPartsCount >= 1; repoPartsCount-- {
		repo := strings.Join(parts[:repoPartsCount], "/")

		replacer, ok := v.replacer[repo]
		if !ok {
			continue
		}

		return replacer + link, true
	}

	if len(parts) > importParts && strings.EqualFold(parts[0], github) {
		parts[0] = githubRaw

		part1 := strings.Join(parts[:importParts], "/")
		part2 := strings.Join(parts[importParts:], "/")

		return fmt.Sprintf("%s/master/%s", part1, part2), true
	}

	return "", false
}
