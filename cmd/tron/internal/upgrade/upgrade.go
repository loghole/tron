package upgrade

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

var ErrVersionNotFound = errors.New("not found")

const (
	releasesURL   = "https://api.github.com/repos/loghole/tron/releases"
	cmdGit        = "git"
	cmdMake       = "make"
	printReleases = 10
)

type Upgrade struct {
	printer  stdout.Printer
	releases []*release
}

func New(printer stdout.Printer) (*Upgrade, error) {
	list, err := releasesList()
	if err != nil {
		return nil, simplerr.Wrap(err, "get releases list failed")
	}

	return &Upgrade{printer: printer, releases: list}, nil
}

func (u *Upgrade) ListVersions() error {
	u.printer.Println(color.FgBlack, "Available versions")

	for idx, r := range u.releases {
		u.printer.Printf(color.FgBlack, "\t%s, published at: %s\n", color.CyanString(r.TagName), r.PublishedAt)

		if idx > printReleases {
			return nil
		}
	}

	return nil
}

func (u *Upgrade) Upgrade(tag string) error {
	rel, err := u.findReleaseByTag(tag)
	if err != nil {
		return simplerr.Wrapf(err, "version '%s'", tag)
	}

	if version.CliVersion == rel.TagName {
		u.printer.Printf(color.FgBlack, "You already use version %s\n", color.CyanString(rel.TagName))

		return nil
	}

	if err := u.downloadAndInstall(rel); err != nil {
		return err
	}

	u.printer.Println(color.FgGreen, "Success\n")

	return nil
}

func (u *Upgrade) findReleaseByTag(tag string) (*release, error) {
	if len(u.releases) == 0 {
		u.printer.Println(color.FgCyan, "No versions available ¯\\_(ツ)_/¯")
		os.Exit(1)
	}

	if tag == "latest" {
		return u.releases[0], nil
	}

	for _, rel := range u.releases {
		if tag == rel.TagName {
			return rel, nil
		}
	}

	return nil, ErrVersionNotFound
}

func (u *Upgrade) downloadAndInstall(rel *release) error {
	dir, err := ioutil.TempDir("", "tron-build")
	if err != nil {
		return simplerr.Wrap(err, "failed to create temp dir")
	}

	defer os.RemoveAll(dir)

	u.printer.Printf(color.FgBlack, "Download %s\n", color.CyanString(rel.TagName))

	cmd := exec.Command(cmdGit, "clone", "https://github.com/loghole/tron.git")
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	cmd = exec.Command(cmdGit, "checkout", rel.TagName) // nolint:gosec //all good
	cmd.Dir = filepath.Join(dir, "tron")

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	u.printer.Println(color.FgBlack, "Build...\n")

	cmd = exec.Command(cmdMake, "build")
	cmd.Dir = filepath.Join(dir, "tron", "cmd", "tron")

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	return nil
}

func releasesList() ([]*release, error) {
	resp, err := http.Get(releasesURL) // nolint:gosec,bodyclose,noctx //body is closed
	if err != nil {
		return nil, simplerr.Wrap(err, "request failed")
	}

	defer helpers.Close(resp.Body)

	dest := make([]*release, 0)

	if err := jsoniter.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return nil, simplerr.Wrap(err, "unmarshal failed")
	}

	result := make([]*release, 0, len(dest))

	for _, rel := range dest {
		if strings.Contains(rel.TagName, "cmd/tron/") {
			continue
		}

		result = append(result, rel)
	}

	return result, nil
}
