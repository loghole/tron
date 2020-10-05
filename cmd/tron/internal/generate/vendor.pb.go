package generate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/project"
)

const (
	vendorDir = "vendor.pb"
)

type VendorPB struct {
	project   *project.Project
	importRgx *regexp.Regexp
	replacer  map[string]string
	exists    map[string]struct{}
	imports   []string
}

func NewVendorPB(pr *project.Project) *VendorPB {
	return &VendorPB{
		project:   pr,
		importRgx: regexp.MustCompile(`^import "(.*?)";$`),
		replacer: map[string]string{
			"google/type":     "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/api":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/rpc":      "https://raw.githubusercontent.com/googleapis/googleapis/master/",
			"google/protobuf": "https://raw.githubusercontent.com/google/protobuf/master/src/",
		},
		exists:  make(map[string]struct{}, 0),
		imports: make([]string, 0),
	}
}

func (v *VendorPB) Download() error {
	if err := v.scanProtos(v.project.Protos); err != nil {
		return err
	}

	if err := v.vendoring(); err != nil {
		return err
	}

	return nil
}

func (v *VendorPB) scanProtos(protos []*models.Proto) error {
	for _, proto := range protos {
		if err := v.scanFile(proto.Path); err != nil {
			return err
		}
	}

	return nil
}

func (v *VendorPB) scanFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer helpers.Close(file)

	return v.findImports(file)
}

func (v *VendorPB) findImports(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		m := v.importRgx.FindStringSubmatch(scanner.Text())
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

func (v *VendorPB) vendoring() (err error) {
	fmt.Println("Start vendoring:")

	for val := ""; len(v.imports) > 0; {
		val, v.imports = v.imports[0], v.imports[1:]

		fmt.Printf("    vendor %s: ", val)

		switch {
		case strings.HasPrefix(val, v.project.Module):
			err = v.copyProto(val)
		default:
			err = v.curlProto(val)
		}

		if err != nil {
			color.Red("FAIL: %v", err)

			return err
		}

		color.Green("OK")
	}

	return nil
}

func (v *VendorPB) copyProto(name string) error {
	filename := strings.TrimPrefix(strings.TrimPrefix(name, v.project.Module), "/")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return helpers.WriteToFile(path.Join(vendorDir, name), data)
}

func (v *VendorPB) curlProto(name string) error {
	link, ok := v.importLink(name)
	if !ok {
		return errors.New("bad import")
	}

	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer helpers.Close(resp.Body)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := v.findImports(bytes.NewReader(data)); err != nil {
		return err
	}

	return helpers.WriteToFile(path.Join(vendorDir, name), data)
}

func (v *VendorPB) importLink(s string) (string, bool) {
	parts := strings.Split(s, "/")

	for repoPartsCount := len(parts); repoPartsCount >= 1; repoPartsCount-- {
		repo := strings.Join(parts[:repoPartsCount], "/")

		replacer, ok := v.replacer[repo]
		if !ok {
			continue
		}

		return strings.Join([]string{replacer, s}, ""), true
	}

	return "", false
}

func (v *VendorPB) checkExists(s string) (bool, error) {
	_, err := os.Stat(path.Join(vendorDir, s))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
