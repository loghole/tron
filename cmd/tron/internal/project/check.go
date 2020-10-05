package project

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
)

const (
	GoMinVersion     = "1.13.15"
	GitMinVersion    = "2.3.4"
	ProtocMinVersion = "3.3.0"
)

type Checker struct {
	printer stdout.Printer
}

func NewChecker(printer stdout.Printer) *Checker {
	return &Checker{
		printer: printer,
	}
}

func (c *Checker) CheckRequirements() (failed bool) {
	c.printer.VerbosePrint(color.FgMagenta, "Check system requirements:\n")

	checks := []struct {
		name string
		fn   func() error
	}{
		{name: "git", fn: c.checkGitVersion},
		{name: "golang", fn: c.checkGoVersion},
		{name: "protoc", fn: c.checkProtocVersion},
	}

	for _, check := range checks {
		c.printer.VerbosePrintf(color.FgBlack, "\t%s version: ", check.name)

		if err := check.fn(); err != nil {
			c.printer.VerbosePrintf(color.FgRed, "FAIL: %v\n", err)
			failed = true

			continue
		}

		c.printer.VerbosePrint(color.FgGreen, "OK\n")
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
		return fmt.Errorf("should be >= %s current version is %s", minReqVersion, v.String())
	}

	return nil
}

func (c *Checker) extractVersion(s string) (string, error) {
	matches := models.VersionRegexp.FindStringSubmatch(s)
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("string is not semantic version")
}
