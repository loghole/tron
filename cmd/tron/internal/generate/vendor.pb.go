package generate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	github      = "github.com"
	githubRaw   = "https://raw.githubusercontent.com"
	importParts = 3
)

var (
	ErrBadImport = errors.New("bad import")
	ErrReqFailed = errors.New("request failed")
)

type vendorPB struct {
	printer  stdout.Printer
	project  *project.Project
	replacer map[string]string
	exists   map[string]struct{}
	imports  []string
}

func VendorPB(p *project.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Start vendor pb")

	generator := &vendorPB{
		printer: printer,
		project: p,
		replacer: map[string]string{
			"google/type":     "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/api":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/rpc":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/protobuf": "https://raw.githubusercontent.com/google/protobuf/master/src/",
		},
		exists:  make(map[string]struct{}),
		imports: make([]string, 0),
	}

	return generator.run()
}

func (v *vendorPB) run() (err error) {
	if err := v.scanProtos(v.project.Protos); err != nil {
		return err
	}

	for val := ""; len(v.imports) > 0; {
		val, v.imports = v.imports[0], v.imports[1:]

		v.printer.VerbosePrintf(color.Reset, "\tvendor %s: ", val)

		switch {
		case strings.HasPrefix(val, v.project.Module):
			err = v.copyProto(val)
		default:
			err = v.curlProto(val)
		}

		if err != nil {
			v.printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		v.printer.VerbosePrintln(color.FgGreen, "OK")
	}

	return nil
}

func (v *vendorPB) scanProtos(protos []*models.Proto) error {
	for _, proto := range protos {
		for _, val := range proto.Imports {
			if _, ok := v.exists[val]; ok {
				continue
			}

			v.exists[val], v.imports = struct{}{}, append(v.imports, val)
		}
	}

	return nil
}

func (v *vendorPB) copyProto(name string) error {
	filename := strings.TrimPrefix(strings.TrimPrefix(name, v.project.Module), "/")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return helpers.WriteToFile(path.Join(models.ProjectPathVendorPB, name), data)
}

func (v *vendorPB) curlProto(name string) error {
	link, ok := v.importLink(name)
	if !ok {
		return ErrBadImport
	}

	resp, err := http.Get(link) // nolint:gosec,bodyclose,noctx //body is closed
	if err != nil {
		return err
	}

	defer helpers.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return simplerr.Wrapf(ErrReqFailed, "link: %s code: %d", link, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := v.findImports(bytes.NewReader(data)); err != nil {
		return err
	}

	return helpers.WriteToFile(path.Join(models.ProjectPathVendorPB, name), data)
}

func (v *vendorPB) importLink(s string) (string, bool) {
	parts := strings.Split(s, "/")

	for repoPartsCount := len(parts); repoPartsCount >= 1; repoPartsCount-- {
		repo := strings.Join(parts[:repoPartsCount], "/")

		replacer, ok := v.replacer[repo]
		if !ok {
			continue
		}

		return strings.Join([]string{replacer, s}, ""), true
	}

	if len(parts) > importParts && strings.EqualFold(parts[0], github) {
		parts[0] = githubRaw

		part1 := strings.Join(parts[:importParts], "/")
		part2 := strings.Join(parts[importParts:], "/")

		return fmt.Sprintf("%s/master/%s", part1, part2), true
	}

	return "", false
}

func (v *vendorPB) findImports(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		m := models.ImportRegexp.FindStringSubmatch(scanner.Text())
		if len(m) == 0 {
			continue
		}

		if _, ok := v.exists[m[0]]; ok {
			continue
		}

		v.exists[m[1]], v.imports = struct{}{}, append(v.imports, m[1])
	}

	return nil
}
