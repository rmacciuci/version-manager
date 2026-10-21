package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rmacciuci/version-manager/internal/semver"
)

// gitIn runs git in dir and fails the test on error.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}

	return strings.TrimSpace(string(out))
}

// TestRun creates a remote with existing tags, runs the tool in a clone and
// checks the pushed tag, the VERSION file and the GitHub Actions outputs.
func TestRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	remote := filepath.Join(t.TempDir(), "remote.git")
	gitIn(t, t.TempDir(), "init", "--bare", "--quiet", "-b", "main", remote)

	seed := t.TempDir()
	gitIn(t, seed, "init", "--quiet", "-b", "main")
	gitIn(t, seed, "-c", "user.email=t@t", "-c", "user.name=t", "-c", "commit.gpgsign=false", "commit", "--quiet", "--allow-empty", "-m", "init")
	for _, tag := range []string{"v0.9.0", "api/v1.9.0", "api/v1.10.0", "web/v2.0.0"} {
		gitIn(t, seed, "-c", "tag.gpgsign=false", "tag", tag)
	}
	gitIn(t, seed, "remote", "add", "origin", remote)
	gitIn(t, seed, "push", "--quiet", "origin", "main", "--tags")

	work := t.TempDir()
	gitIn(t, work, "clone", "--quiet", remote, ".")
	gitIn(t, work, "config", "tag.gpgsign", "false")
	t.Chdir(work)

	output := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", output)

	if err := run([]string{"--minor=true", "--prefix=api/"}); err != nil {
		t.Fatal(err)
	}

	if got := gitIn(t, work, "ls-remote", "--tags", "origin", "api/v1.11.0"); got == "" {
		t.Fatal("api/v1.11.0 was not pushed")
	}

	version, err := os.ReadFile("VERSION")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(version)) != "api/v1.11.0" {
		t.Fatalf("VERSION = %q", version)
	}

	outputs, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(outputs) != "tag=api/v1.11.0\nversion=v1.11.0\n" {
		t.Fatalf("outputs = %q", outputs)
	}
}

func TestRunRequiresABump(t *testing.T) {
	// run touches the global git config, so never let a test use the real one.
	t.Setenv("HOME", t.TempDir())

	if err := run([]string{"--prefix=api/"}); !errors.Is(err, semver.ErrNoBump) {
		t.Fatalf("want ErrNoBump, got %v", err)
	}
}
