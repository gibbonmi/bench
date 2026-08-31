package coverage

import (
	"bytes"
	"context"
	"fmt"
	"go/build"
	"go/build/constraint"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate"
)

var packageLoadTimeout = bounds.PackageLoadTimeout

// executionScope is one Go test phase and the package directories the Go loader expands
// its operands to. The phase and its packages stay paired here, so no citation can
// borrow one phase's package scope and another phase's tag set.
type executionScope struct {
	phase    gate.TestExecution
	dirs     []string
	err      error
	resolved bool
}

// executionScopes carries one scope for each phase of a census.
func executionScopes(census []gate.TestExecution) []*executionScope {
	scopes := make([]*executionScope, 0, len(census))
	for _, phase := range census {
		scopes = append(scopes, &executionScope{phase: phase})
	}
	return scopes
}

// packages expands the phase's operands on first need and keeps the answer, so one
// coverage map runs the loader once for each phase rather than once for each cited file.
// The expansion waits for a citation that needs it: a row the grammar refuses before it
// reads the tree, such as a cited FIFO, must not start a loader that walks that tree.
func (s *executionScope) packages() ([]string, error) {
	if !s.resolved {
		s.dirs, s.err = selectedPackageDirs(s.phase)
		s.resolved = true
	}
	return s.dirs, s.err
}

// checkExecution accepts a citation only when one complete test phase both selects its
// package and builds its file. An empty census is inapplicable because the gate has no
// Go test phase for the tree.
//
// A phase the loader could not expand rejects the citation. An unknown package scope
// authorizes nothing, so the row reports rather than falls back to the tag answer alone.
func checkExecution(rn int, rel, path string, content []byte, scopes []*executionScope) []string {
	if len(scopes) == 0 {
		return nil
	}
	line, err := buildLine(content)
	if err != nil {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', whose //go:build expression does not parse: %s", rn, rel, err)}
	}
	dir, name := filepath.Split(path)
	selected := false
	for _, scope := range scopes {
		dirs, err := scope.packages()
		if err != nil {
			return []string{fmt.Sprintf("coverage map row %d cites '%s', could not expand packages for executed phase %q: %s", rn, rel, scope.phase.Name, err)}
		}
		if !containsDirectory(dirs, dir) {
			continue
		}
		selected = true
		ctxt := contextFor(scope.phase)
		if matched, err := ctxt.MatchFile(dir, name); err == nil && matched {
			return nil
		}
	}
	if !selected {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', which no executed test phase selects", rn, rel)}
	}
	return []string{fmt.Sprintf("coverage map row %d cites '%s', which no executed tag set builds (%s)", rn, rel, constraintOf(line))}
}

// selectedPackageDirs delegates package-pattern expansion to the installed Go loader.
// The invocation retains the phase directory, tags, environment, and Go -C setting so its
// package selection has the same terms as the phase that can execute the cited test.
//
// Only stdout is read. `go list` writes its matched-no-packages warning to stderr, and a
// combined stream would let that sentence pose as a package directory. The loader is
// asked for the strict answer, without -e, so a package operand it cannot resolve exits
// nonzero and becomes the caller's expansion failure. A pattern that legitimately
// matches nothing exits zero with no line, which selects nothing rather than fails.
func selectedPackageDirs(execution gate.TestExecution) ([]string, error) {
	args := make([]string, 0, len(execution.Packages)+6)
	if execution.GoC != "" {
		args = append(args, "-C", execution.GoC)
	}
	args = append(args, "list", "-f", "{{.Dir}}")
	if len(execution.Tags) != 0 {
		args = append(args, "-tags="+strings.Join(execution.Tags, ","))
	}
	packages := execution.Packages
	if len(packages) == 0 {
		packages = []string{"."}
	}
	args = append(args, packages...)
	cmd := exec.Command("go", args...)
	cmd.Dir = execution.Dir
	if len(execution.Env) != 0 {
		// The phase's own overrides reach the loader whole and last, so the loader answers
		// for the platform and module mode the phase runs under rather than the host's.
		cmd.Env = append(os.Environ(), execution.Env...)
	}
	var stdout bytes.Buffer
	// The run waits for the child, so an interrupted loader reports its terminal result
	// here rather than leaving a partial directory list behind. It also kills the whole
	// process group at the bound, because a grandchild that inherited the output pipe
	// would otherwise hold the read open long after the loader itself is gone.
	result := bounds.RunOutput(context.Background(), packageLoadTimeout, cmd, &stdout)
	if result.Status == bounds.ProcessTimeout {
		return nil, fmt.Errorf("go list: timed out after %s", packageLoadTimeout)
	}
	if result.Err != nil {
		if detail := strings.TrimSpace(string(result.Output)); detail != "" {
			return nil, fmt.Errorf("go list: %w: %s", result.Err, detail)
		}
		return nil, fmt.Errorf("go list: %w", result.Err)
	}
	return packageDirs(stdout.Bytes()), nil
}

// packageDirs reads one package directory for each line the loader prints. The split is
// by line rather than by field, because a package directory can hold a space.
func packageDirs(stdout []byte) []string {
	var dirs []string
	for _, line := range strings.Split(string(stdout), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			dirs = append(dirs, line)
		}
	}
	return dirs
}

func containsDirectory(dirs []string, want string) bool {
	want = filepath.Clean(want)
	for _, dir := range dirs {
		if filepath.Clean(dir) == want {
			return true
		}
	}
	return false
}

// contextFor is the build context one executed phase compiles in. It starts from the
// toolchain's own default, so the release tags and toolchain tags retain their real
// values, and the host GOOS and GOARCH stand where the phase declares nothing else.
// MatchFile then applies the file constraint and suffix.
//
// Three environment variables name the platform a file is selected for, so the phase's
// own settings for them override the host's. Every other variable the phase declares
// reaches the loader child instead, where Go itself applies it.
func contextFor(execution gate.TestExecution) build.Context {
	ctxt := build.Default
	ctxt.BuildTags = append(append([]string(nil), build.Default.BuildTags...), execution.Tags...)
	for _, entry := range execution.Env {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		switch key {
		case "GOOS":
			ctxt.GOOS = value
		case "GOARCH":
			ctxt.GOARCH = value
		case "CGO_ENABLED":
			ctxt.CgoEnabled = value == "1"
		}
	}
	return ctxt
}

// buildLine returns the file's //go:build line, and an error when that line holds an
// expression the toolchain cannot parse. A file with no such line returns the empty
// string, which is the always-satisfied case. The scan stops at the package clause.
func buildLine(content []byte) (string, error) {
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "package ") {
			return "", nil
		}
		if !constraint.IsGoBuild(line) {
			continue
		}
		if _, err := constraint.Parse(line); err != nil {
			return "", err
		}
		return line, nil
	}
	return "", nil
}

// constraintOf names the constraint that kept a selected file out of every phase. A
// file with no build line was refused by its filename suffix.
func constraintOf(line string) string {
	if line == "" {
		return "the filename's GOOS or GOARCH suffix"
	}
	return line
}
