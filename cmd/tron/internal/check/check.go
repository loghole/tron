package check

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/download"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/version"
)

const (
	GoMinVersion     = "1.17.0"
	GitMinVersion    = "2.3.4"
	ProtocMinVersion = "3.3.0"
)

var ErrCheckVersionFailed = errors.New("check failed")

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

func (c *Checker) CheckTron() {
	if err := c.checkTronIsLatest(); err != nil {
		c.printer.VerbosePrintf(color.Reset, "check tron version failed: %v\n", err)
	}
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

func (c *Checker) checkVersion(cmd *exec.Cmd, min string) error {
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	appVersion, err := models.ExtractVersion(string(output))
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(appVersion)
	if err != nil {
		fmt.Print("semver err", err)
	}

	constr, _ := semver.NewConstraint(">= " + min)

	if !constr.Check(v) {
		return fmt.Errorf("should be >= %s current version is %s: %w", min, v.String(), ErrCheckVersionFailed)
	}

	return nil
}

func (c *Checker) checkTronIsLatest() error {
	config := models.NewConfig()

	if err := config.Read(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if config.LastVersionCheck.After(time.Now().Truncate(time.Hour)) {
		return nil
	}

	config.LastVersionCheck = time.Now()

	if err := config.Write(); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	latest, err := download.LatestTronVersion()
	if err != nil {
		return fmt.Errorf("get latest tron version: %w", err)
	}

	currentVersion, err := semver.NewVersion(version.CliVersion)
	if err != nil {
		return fmt.Errorf("parse current version: %w", err)
	}

	latestVersion, err := semver.NewVersion(latest)
	if err != nil {
		return fmt.Errorf("parse latest version: %w", err)
	}

	if currentVersion.Compare(latestVersion) >= 0 {
		return nil
	}

	c.printer.Println(color.FgBlue, "New version of Tron is available!")
	c.printer.Printf(color.Reset, "\tcurrent: %s\n", color.YellowString(version.CliVersion))
	c.printer.Printf(color.Reset, "\tlatest: %s\n\n", color.GreenString(latest))

	return nil
}
