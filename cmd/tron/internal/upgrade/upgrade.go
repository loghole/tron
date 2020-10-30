package upgrade

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

var ErrVersionNotFound = errors.New("not found")

const (
	repositoryURL = "https://github.com/loghole/tron.git"
	releasesURL   = "https://api.github.com/repos/loghole/tron/releases"
	versionLdflag = `-X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=%s'`

	minTronVersion = "v0.3.0"
	printReleases  = 10

	cmdGit = "git"
	cmdGo  = "go"
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
	u.printer.Println(color.Reset, "Available versions:")

	for idx, r := range u.releases {
		u.printer.Printf(color.Reset, "\t%s, published at: %s\n", color.CyanString(r.TagName), r.PublishedAt)

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
		u.printer.Printf(color.Reset, "You already use version %s\n", color.CyanString(rel.TagName))

		return nil
	}

	if err := u.downloadAndInstall(rel); err != nil {
		return err
	}

	u.printer.Println(color.FgGreen, "Success")

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

	u.printer.Printf(color.Reset, "Download tron %s\n", color.CyanString(rel.TagName))

	if err := u.download(rel, dir); err != nil {
		return simplerr.Wrap(err, "failed to download")
	}

	u.printer.Println(color.Reset, "Install tron...")

	if err := u.install(rel, dir); err != nil {
		return simplerr.Wrap(err, "failed to install")
	}

	return nil
}

func (u *Upgrade) download(rel *release, dir string) error {
	cmd := exec.Command(cmdGit, "clone", repositoryURL)
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	cmd = exec.Command(cmdGit, "checkout", rel.TagName) // nolint:gosec //all good
	cmd.Dir = filepath.Join(dir, "tron")

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	return nil
}

func (u *Upgrade) install(rel *release, dir string) error {
	path, err := binaryPath()
	if err != nil {
		return simplerr.Wrapf(err, "failed to get binary path %v", err)
	}

	args := []string{`build`, `-o`, path, `-ldflags`, fmt.Sprintf(versionLdflag, rel.TagName), `main.go`}

	cmd := exec.Command(cmdGo, args...) // nolint:gosec //all good
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

	constraint, err := semver.NewConstraint(">= " + minTronVersion)
	if err != nil {
		return nil, simplerr.Wrap(err, "build semver constraint failed")
	}

	for _, rel := range dest {
		if !constraint.Check(rel.version()) {
			continue
		}

		result = append(result, rel)
	}

	return result, nil
}

func binaryPath() (string, error) {
	binaryPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(binaryPath)
}
