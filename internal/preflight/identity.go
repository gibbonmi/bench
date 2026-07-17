package preflight

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var exactTag = regexp.MustCompile(`^refs/tags/(v[0-9]+\.[0-9]+\.[0-9]+)$`)

func (r *runner) checkIdentity(ctx context.Context) error {
	ref := os.Getenv("BENCH_PREFLIGHT_REF")
	if ref == "" {
		ref = os.Getenv("GITHUB_REF")
	}
	m := exactTag.FindStringSubmatch(ref)
	if m == nil {
		return commandFailure{Failure{Kind: "identity", Message: "publish requires exact GITHUB_REF refs/tags/vMAJOR.MINOR.PATCH"}}
	}
	tag := m[1]
	version := strings.TrimPrefix(tag, "v")
	pkg, err := readPackageVersion(r.root)
	if err != nil {
		return commandFailure{Failure{Kind: "input", Message: "package.json version is unreadable"}}
	}
	toolchain, err := readToolchain(r.root)
	if err != nil {
		return commandFailure{Failure{Kind: "identity", Message: err.Error()}}
	}
	commit := ""
	if r.identity.SourceCommit != nil {
		commit = *r.identity.SourceCommit
	}
	tagCommit, err := gitOutput(r.root, "rev-list", "-n", "1", tag)
	if err != nil || tagCommit != commit {
		return commandFailure{Failure{Kind: "identity", Message: "tag does not resolve exactly to HEAD"}}
	}
	if pkg != version || r.binaryVersion != version {
		return commandFailure{Failure{Kind: "identity", Message: "tag, package version, and binary version must agree"}}
	}
	goVersion, err := exec.CommandContext(ctx, "go", "env", "GOVERSION").Output()
	if err != nil || strings.TrimSpace(string(goVersion)) != "go"+toolchain {
		return commandFailure{Failure{Kind: "identity", Message: "actual Go version must equal go.mod toolchain patch"}}
	}
	r.identity.Tag = &tag
	r.identity.PackageVersion = &pkg
	r.identity.BinaryVersion = &r.binaryVersion
	r.identity.Toolchain = &toolchain
	return nil
}

func (r *runner) checkAncestry(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "-C", r.root, "merge-base", "--is-ancestor", "HEAD", "origin/main")
	if err := cmd.Run(); err != nil {
		return commandFailure{Failure{Kind: "identity", Message: "tagged HEAD ancestry to origin/main could not be proven"}}
	}
	return nil
}

func (r *runner) checkChangelog() error {
	if r.identity.Tag == nil {
		return commandFailure{Failure{Kind: "identity", Message: "release identity must be green before changelog"}}
	}
	data, err := readRegular(filepath.Join(r.root, "CHANGELOG.md"))
	if err != nil {
		return commandFailure{Failure{Kind: "input", Message: "CHANGELOG.md is unreadable"}}
	}
	version := strings.TrimPrefix(*r.identity.Tag, "v")
	re := regexp.MustCompile(`(?m)^## \[` + regexp.QuoteMeta(version) + `\] - ([0-9]{4}-[0-9]{2}-[0-9]{2})$`)
	matches := re.FindAllSubmatch(data, -1)
	if len(matches) != 1 {
		return commandFailure{Failure{Kind: "identity", Message: "CHANGELOG.md must contain exactly one matching release heading"}}
	}
	if _, err := time.Parse("2006-01-02", string(matches[0][1])); err != nil {
		return commandFailure{Failure{Kind: "identity", Message: "CHANGELOG.md release heading date is invalid"}}
	}
	heading := string(matches[0][0])
	releaseAt := bytes.Index(data, matches[0][0])
	unreleasedHeading := []byte("## [Unreleased]")
	unreleasedAt := bytes.Index(data, unreleasedHeading)
	if bytes.Count(data, unreleasedHeading) != 1 || unreleasedAt >= releaseAt {
		return commandFailure{Failure{Kind: "identity", Message: "CHANGELOG.md must contain one Unreleased heading before the release"}}
	}
	between := string(data[unreleasedAt+len(unreleasedHeading) : releaseAt])
	for _, line := range strings.Split(between, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			return commandFailure{Failure{Kind: "identity", Message: "CHANGELOG.md has stranded content under Unreleased"}}
		}
	}
	r.identity.ChangelogHeading = &heading
	return nil
}

func readPackageVersion(root string) (string, error) {
	data, err := readRegular(filepath.Join(root, "package.json"))
	if err != nil {
		return "", err
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(data, &pkg) != nil || pkg.Version == "" {
		return "", errors.New("invalid package version")
	}
	return pkg.Version, nil
}
func readToolchain(root string) (string, error) {
	data, err := readRegular(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(?m)^toolchain go([0-9]+\.[0-9]+\.[0-9]+)$`)
	m := re.FindSubmatch(data)
	if m == nil {
		return "", errors.New("go.mod requires an exact patch toolchain directive")
	}
	return string(m[1]), nil
}
func gitOutput(root string, args ...string) (string, error) {
	argv := append([]string{"-C", root}, args...)
	out, err := exec.Command("git", argv...).Output()
	return strings.TrimSpace(string(out)), err
}

func readRegular(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", filepath.Base(path))
	}
	return os.ReadFile(path)
}
