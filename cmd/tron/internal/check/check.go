package check

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/lissteron/simplerr"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	GoMinVersion     = "1.13.15"
	GitMinVersion    = "2.3.4"
	ProtocMinVersion = "3.3.0"
)

var (
	ErrNotSemanticVersion = errors.New("string is not semantic version")
	ErrCheckVersionFailed = errors.New("check failed")
)

type Checker struct {
	printer stdout.Printer
}

func NewChecker(printer stdout.Printer) *Checker {
	return &Checker{
		printer: printer,
	}
}

func (c *Checker) CheckAllRequirements() (failed bool) {
	return c.checkRequirements(map[string]func() error{
		"git":    c.checkGitVersion,
		"golang": c.checkGoVersion,
		"protoc": c.checkProtocVersion,
	})
}

func (c *Checker) CheckInitRequirements() (failed bool) {
	return c.checkRequirements(map[string]func() error{
		"git":    c.checkGitVersion,
		"golang": c.checkGoVersion,
	})
}

func (c *Checker) CheckGolang() (failed bool) {
	return c.checkRequirements(map[string]func() error{
		"golang": c.checkGoVersion,
	})
}

func (c *Checker) CheckProtoc() (failed bool) {
	return c.checkRequirements(map[string]func() error{
		"protoc": c.checkProtocVersion,
	})
}

func (c *Checker) checkRequirements(checks map[string]func() error) (failed bool) {
	c.printer.Print(color.FgMagenta, "Check system requirements:\n")

	for name, check := range checks {
		c.printer.Printf(color.Reset, "\t%s version: ", name)

		if err := check(); err != nil {
			c.printer.Printf(color.FgRed, "FAIL: %v\n", err)

			failed = true

			continue
		}

		c.printer.Print(color.FgGreen, "OK\n")
	}

	return !failed
}

func (c *Checker) checkGoVersion() error {
	return c.checkVersion(exec.Command("go", "version"), GoMinVersion)
}

func (c *Checker) checkProtocVersion() error {
	return c.checkVersion(exec.Command("protoc", "--version"), ProtocMinVersion)
}

func (c *Checker) checkGitVersion() error {
	return c.checkVersion(exec.Command("git", "version"), GitMinVersion)
}

func (c *Checker) checkVersion(cmd *exec.Cmd, minReqVersion string) error {
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	version, err := c.extractVersion(string(output))
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		fmt.Print("semver err", err)
	}

	constr, _ := semver.NewConstraint(">= " + minReqVersion)

	if !constr.Check(v) {
		return simplerr.Wrapf(ErrCheckVersionFailed, "should be >= %s current version is %s", minReqVersion, v.String())
	}

	return nil
}

func (c *Checker) extractVersion(s string) (string, error) {
	matches := models.VersionRegexp.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", ErrNotSemanticVersion
}
