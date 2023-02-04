package download

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

type Deps struct {
	project *models.Project
	printer stdout.Printer
}

func NewDeps(project *models.Project, printer stdout.Printer) *Deps {
	return &Deps{
		project: project,
		printer: printer,
	}
}

func (p *Deps) InstallProtoPlugins() error {
	p.printer.Println(color.FgMagenta, "Install protobuf plugins")

	dir, err := tronCacheDir()
	if err != nil {
		return simplerr.Wrap(err, "get tron cache dir")
	}

	for _, dep := range models.ProtobufDeps() {
		if p.exists(dep.Out) && p.isActual(dep) {
			p.printer.VerbosePrintf(color.Reset, "\tplugin '%s' already exists\n", color.YellowString(dep.Out))

			continue
		}

		p.printer.VerbosePrintf(color.Reset, "\tinstall plugin '%s': ", color.YellowString(dep.Out))

		if err := p.install(dir, dep); err != nil {
			p.printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return simplerr.Wrapf(err, "download '%s'", dep.Out)
		}

		p.printer.VerbosePrintln(color.FgGreen, "OK")
	}

	p.printer.Println(color.FgBlue, "\tSuccess")

	return nil
}

func (p *Deps) install(dir string, dep *models.Dep) error {
	if err := p.download(dir, dep.Git); err != nil {
		return simplerr.Wrapf(err, "download '%s'", dep.Out)
	}

	if err := p.checkout(dir, dep.Git, dep.Version); err != nil {
		return simplerr.Wrapf(err, "checkout '%s' to version '%s'", dep.Out, dep.Version)
	}

	if err := p.build(dir, dep.Main, dep.Out, dep.LdFlag); err != nil {
		return simplerr.Wrapf(err, "build '%s'", dep.Out)
	}

	return nil
}

func (p *Deps) exists(target string) bool {
	_, err := os.Stat(filepath.Join(p.project.AbsPath, target))

	return !os.IsNotExist(err)
}

func (p *Deps) download(dir, repo string) error {
	stat, err := os.Stat(filepath.Join(dir, strings.TrimSuffix(filepath.Base(repo), filepath.Ext(repo)), ".git"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Skip download if git dir exists.
	if !os.IsNotExist(err) && stat.IsDir() {
		return nil
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return simplerr.Wrapf(err, "can't create dir '%s'", dir)
	}

	cmd := exec.Command(cmdGit, "clone", repo)
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	return nil
}

func (p *Deps) checkout(dir, repo, version string) error {
	cmd := exec.Command(cmdGit, "checkout", version) //nolint:gosec // need git version to checkout.
	cmd.Dir = filepath.Join(dir, strings.TrimSuffix(filepath.Base(repo), filepath.Ext(repo)))

	if err := cmd.Run(); err != nil {
		cmd := exec.Command(cmdGit, "fetch", "--all", "--tags")
		cmd.Dir = filepath.Join(dir, strings.TrimSuffix(filepath.Base(repo), filepath.Ext(repo)))

		if err := cmd.Run(); err != nil {
			return simplerr.Wrapf(err, "failed to run %s", cmd.String())
		}

		cmd = exec.Command(cmdGit, "checkout", version) //nolint:gosec // need git version to checkout.
		cmd.Dir = filepath.Join(dir, strings.TrimSuffix(filepath.Base(repo), filepath.Ext(repo)))

		if err := cmd.Run(); err != nil {
			return simplerr.Wrapf(err, "failed to run %s", cmd.String())
		}
	}

	return nil
}

func (p *Deps) build(dir, main, output, ldflags string) error {
	args := []string{
		"build",
		"-o", filepath.Join(p.project.AbsPath, output),
	}

	if ldflags != "" {
		args = append(args, "-ldflags", "-X "+ldflags)
	}

	//nolint:gosec // need output and ldflags.
	cmd := exec.Command(cmdGo, args...)
	cmd.Dir = filepath.Join(dir, main)

	if err := cmd.Run(); err != nil {
		return simplerr.Wrapf(err, "failed to run %s", cmd.String())
	}

	return nil
}

func (p *Deps) isActual(dep *models.Dep) bool {
	//nolint:gosec // need binary path.
	cmd := exec.Command(filepath.Join(p.project.AbsPath, dep.Out), "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	isActual, err := dep.IsActual(string(output))
	if err != nil {
		return false
	}

	return isActual
}
