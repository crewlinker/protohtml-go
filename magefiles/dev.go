// Package main provides repo automation using mage.
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// init performs some sanity checks before running anything.
func init() {
	mustBeInRootIfNotTest()
}

// Dev groups commands for local development.
type Dev mg.Namespace

// Lint the codebase through static code analysis.
func (Dev) Lint() error {
	if err := sh.Run("golangci-lint", "run"); err != nil {
		return fmt.Errorf("failed to run golang-ci: %w", err)
	}

	if err := sh.Run("buf", "lint"); err != nil {
		return fmt.Errorf("failed to run buf lint: %w", err)
	}

	return nil
}

// Test tests all the code using Gingo, with an empty label filter.
func (Dev) Test() error {
	return (Dev{}).TestSome("")
}

// TestSome tests the whole repo using Ginkgo test runner with label filters applied.
func (Dev) TestSome(labelFilter string) error {
	if err := sh.Run(
		"go", "run", "-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "--fail-on-pending", "--race", "--trace",
		"--junit-report=test-report.xml",
		"--label-filter", labelFilter,
		"./...",
	); err != nil {
		return fmt.Errorf("failed to run ginkgo: %w", err)
	}

	return nil
}

// Generate generates code across the repository.
func (Dev) Generate() error {
	// generate rpc
	if err := sh.Run("buf", "generate"); err != nil {
		return fmt.Errorf("failed to generate protobuf: %w", err)
	}

	// generate templ code
	if err := sh.Run("templ", "generate"); err != nil {
		return fmt.Errorf("failed to generate templ: %w", err)
	}

	return nil
}

// error when wrong version format is used.
var errVersionFormat = fmt.Errorf("version must be in format vX,Y,Z")

// Release tags a new version and pushes it.
func (Dev) Release(version string) error {
	if !regexp.MustCompile(`^v([0-9]+).([0-9]+).([0-9]+)$`).Match([]byte(version)) {
		return errVersionFormat
	}

	if err := sh.Run("git", "tag", version); err != nil {
		return fmt.Errorf("failed to tag version: %w", err)
	}

	if err := sh.Run("git", "push", "origin", version); err != nil {
		return fmt.Errorf("failed to push version tag: %w", err)
	}

	return nil
}

// mustBeInRootIfNotTest checks that the command is run in the project root.
func mustBeInRootIfNotTest() {
	if _, err := os.ReadFile("go.mod"); err != nil && !strings.Contains(strings.Join(os.Args, ""), "-test.") {
		panic("must be in project root, couldn't stat go.mod file: " + err.Error())
	}
}
