package download

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

const (
	repositoryURL = "https://github.com/loghole/tron.git"
	releasesURL   = "https://api.github.com/repos/loghole/tron/releases"
	versionLdflag = `-X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=%s'`

	tronMinVersion = "v0.4.0"
	printReleases  = 8

	cmdGit = "git"
	cmdGo  = "go"
)

type Tron struct {
	printer  stdout.Printer
	releases []*models.Release
}

func NewTron(printer stdout.Printer, stable bool) (*Tron, error) {
	list, err := releasesList()
	if err != nil {
		return nil, simplerr.Wrap(err, "get releases list failed")
	}

	list, err = filterReleases(list, stable)
	if err != nil {
		return nil, simplerr.Wrap(err, "filter releases list failed")
	}

	return &Tron{printer: printer, releases: list}, nil
}

func (t *Tron) ListVersions() error {
	t.printer.Println(color.Reset, "Available versions:")

	for idx, r := range t.releases {
		if version.CliVersion == r.TagName {
			t.printer.Printf(color.FgGreen, "\t%s, published at: %s, already installed\n", r.TagName, r.PublishedAt)
		} else {
			t.printer.Printf(color.Reset, "\t%s, published at: %s\n", color.CyanString(r.TagName), r.PublishedAt)
		}

		if idx+1 >= printReleases {
			return nil
		}
	}

	return nil
}

func (t *Tron) Upgrade(tag string) error {
	rel, err := t.findReleaseByTag(tag)
	if err != nil {
		return simplerr.Wrapf(err, "version '%s'", tag)
	}

	if version.CliVersion == rel.TagName {
		t.printer.Printf(color.Reset, "You already use version %s\n", color.CyanString(rel.TagName))

		return nil
	}

	if err := t.downloadAndInstall(rel); err != nil {
		return err
	}

	return nil
}

func (t *Tron) findReleaseByTag(tag string) (*models.Release, error) {
	if len(t.releases) == 0 {
		t.printer.Println(color.FgCyan, "No versions available ¯\\_(ツ)_/¯")
		os.Exit(1)
	}

	if tag == "latest" {
		return t.releases[0], nil
	}

	for _, rel := range t.releases {
		if tag == rel.TagName {
			return rel, nil
		}
	}

	return nil, ErrVersionNotFound
}

func (t *Tron) downloadAndInstall(rel *models.Release) error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return simplerr.Wrap(err, "failed to get cache dir")
	}

	dir = filepath.Join(dir, "tron")

	t.printer.Printf(color.Reset, "Download tron %s\n", color.CyanString(rel.TagName))

	if err := t.download(dir); err != nil {
		return simplerr.Wrap(err, "failed to download")
	}

	if err := t.checkout(rel, dir); err != nil {
		return simplerr.Wrap(err, "failed to checkout")
	}

	t.printer.Println(color.Reset, "Install tron...")

	if err := t.install(rel, dir); err != nil {
		return simplerr.Wrap(err, "failed to install")
	}

	return nil
}

func (t *Tron) checkout(rel *models.Release, dir string) error {
	cmd := exec.Command(cmdGit, "fetch", "--all", "--tags")
	cmd.Dir = filepath.Join(dir, "tron")

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

func (t *Tron) download(dir string) error {
	if _, err := os.Stat(filepath.Join(dir, "tron", ".git")); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return simplerr.Wrapf(err, "can't create dir '%s'", dir)
	}

	cmd := exec.Command(cmdGit, "clone", repositoryURL)
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	return nil
}

func (t *Tron) install(rel *models.Release, dir string) error {
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

func filterReleases(list []*models.Release, stable bool) ([]*models.Release, error) {
	result := make([]*models.Release, 0, len(list))

	minVersion, err := semver.NewVersion(tronMinVersion)
	if err != nil {
		return nil, simplerr.Wrap(err, "parse min semver version failed")
	}

	for _, rel := range list {
		if minVersion.Compare(rel.Version()) >= 0 {
			continue
		}

		if stable && rel.Version().Prerelease() != "" {
			continue
		}

		result = append(result, rel)
	}

	return result, nil
}

func releasesList() ([]*models.Release, error) {
	resp, err := http.Get(releasesURL) // nolint:gosec,bodyclose,noctx //body is closed
	if err != nil {
		return nil, simplerr.Wrap(err, "request failed")
	}

	defer helpers.Close(resp.Body)

	result := make([]*models.Release, 0)

	if err := jsoniter.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, simplerr.Wrap(err, "unmarshal failed")
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
