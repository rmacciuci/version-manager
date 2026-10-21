// Command version-manager creates and pushes the next semver tag of the
// repository in the current directory. It runs as a GitHub Action.
//
//	version-manager --minor=true [--patch=false] [--major=false] [--prefix=api/]
//
// It writes the new tag (with prefix) to ./VERSION and, when running in
// GitHub Actions, sets the outputs "tag" (e.g. api/v1.2.3) and "version"
// (v1.2.3, usable as a Docker image tag).
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rmacciuci/version-manager/internal/semver"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "version-manager:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("version-manager", flag.ContinueOnError)
	var bump semver.Bump
	fs.BoolVar(&bump.Major, "major", false, "increase the major version")
	fs.BoolVar(&bump.Minor, "minor", false, "increase the minor version")
	fs.BoolVar(&bump.Patch, "patch", false, "increase the patch version")
	prefix := fs.String("prefix", "", `tag prefix for monorepos, e.g. "api/"`)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !bump.Major && !bump.Minor && !bump.Patch {
		return semver.ErrNoBump
	}

	fmt.Println("1. Setting up git")
	if err := git("config", "--global", "user.email", "actions@github.com"); err != nil {
		return err
	}
	if err := git("config", "--global", "user.name", "Version Manager"); err != nil {
		return err
	}
	// The action container runs as root on a workspace owned by the runner
	// user; without this, git refuses to work ("dubious ownership").
	if wd, err := os.Getwd(); err == nil {
		if err := git("config", "--global", "--add", "safe.directory", wd); err != nil {
			return err
		}
	}

	fmt.Println("2. Updating repository tags")
	if err := git("fetch", "origin", "--tags", "--quiet"); err != nil {
		return err
	}

	out, err := exec.Command("git", "tag", "--list").Output()
	if err != nil {
		return fmt.Errorf("git tag --list: %w", err)
	}

	latest, found := semver.Latest(strings.Split(string(out), "\n"), *prefix)
	if found {
		fmt.Printf("Last tag: %s%s\n", *prefix, latest)
	} else {
		fmt.Printf("No %svX.Y.Z tag yet, starting from v0.0.0\n", *prefix)
	}

	next, err := semver.Next(latest, bump)
	if err != nil {
		return err
	}

	tag := *prefix + next.String()
	fmt.Println("New tag:", tag)

	if err := git("pull", "--quiet"); err != nil {
		return err
	}

	fmt.Println("3. Creating and pushing the tag")
	if err := git("tag", tag); err != nil {
		return err
	}
	if err := git("push", "origin", tag); err != nil {
		return err
	}

	if err := os.WriteFile("VERSION", []byte(tag+"\n"), 0o644); err != nil {
		return fmt.Errorf("write VERSION: %w", err)
	}

	return writeOutputs(os.Getenv("GITHUB_OUTPUT"), tag, next.String())
}

// writeOutputs appends the step outputs when running in GitHub Actions.
func writeOutputs(path, tag, version string) error {
	if path == "" {
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open GITHUB_OUTPUT: %w", err)
	}

	_, werr := fmt.Fprintf(f, "tag=%s\nversion=%s\n", tag, version)

	return errors.Join(werr, f.Close())
}

func git(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}

	return nil
}
