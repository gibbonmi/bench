package worktree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/responsebound"
)

// overBound writes one response line more than the bound allows, so the owner starts its
// spill before the caller writes anything else.
func overBound(t *testing.T, stream io.Writer) {
	t.Helper()
	for i := 1; i <= 11; i++ {
		if _, err := fmt.Fprintf(stream, "line %02d\n", i); err != nil {
			t.Fatal(err)
		}
	}
}

// namedSpill answers the path that the spill line of one finished response names.
func namedSpill(t *testing.T, response string) string {
	t.Helper()
	for _, line := range strings.SplitAfter(response, "\n") {
		if _, path, found := strings.Cut(strings.TrimSuffix(line, "}\n"), ",path="); found && strings.HasPrefix(line, "spilled{") {
			return path
		}
	}
	t.Fatalf("response = %q, want a spill line", response)
	return ""
}

// spillAt sends one over-bound response from a process whose repository root is root, and
// answers its spill path.
func spillAt(t *testing.T, root string) string {
	t.Helper()
	var sink bytes.Buffer
	owner := responsebound.New(&sink, &sink, func() string { return root })
	overBound(t, owner.Stdout())
	owner.Finish()
	return namedSpill(t, sink.String())
}

// BO16: the retirement of an assignment removes its spill directory, and a spill in the
// primary scope stays.
func TestRetirementDropsResponseSpills(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "spill-release")
	bindEnv(t, homeEnv, home)
	retired := spillAt(t, creation.Path)
	primary := spillAt(t, root)
	requireTest(t, filepath.Dir(retired) != filepath.Dir(primary), "assignment spill %s shares the primary scope of %s", retired, primary)
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "landed-spill-release", creation.Path}, &stdout, &stderr)
	requireTest(t, code == 0, "release = (%d, %q, %q)", code, stdout.String(), stderr.String())
	_, err := os.Stat(filepath.Dir(retired))
	requireTest(t, os.IsNotExist(err), "the release kept the assignment spill directory: %v", err)
	_, err = os.Stat(primary)
	requireTest(t, err == nil, "the release removed the primary spill: %v", err)
}

// BO72: a `release` that runs in the assignment it retires keeps its over-bound spill.
// The spill starts before the retirement, so a spill in the assignment's own scope would
// leave the line naming a removed file.
func TestRetiringVerbSpillsToPrimary(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "spill-own-release")
	bindEnv(t, homeEnv, home)
	args := []string{"--request", "landed-spill-own-release", creation.Path}
	var sink bytes.Buffer
	owner := responsebound.New(&sink, &sink, func() string { return creation.Path }, append([]string{"worktree", "release"}, args...)...)
	overBound(t, owner.Stdout())
	code := ReleaseCommand(root, home, args, owner.Stdout(), owner.Stderr())
	owner.Finish()
	requireTest(t, code == 0, "release = (%d, %q)", code, sink.String())
	path := namedSpill(t, sink.String())
	requireTest(t, !strings.Contains(path, creation.Assignment.ID), "spill %s lies in the retired assignment's scope", path)
	data, err := os.ReadFile(path)
	requireTest(t, err == nil, "the release removed its own spill: %v", err)
	requireTest(t, strings.HasPrefix(string(data), "line 01\n") && strings.Count(string(data), "\n") > 11, "spill = %q, want the complete release response", data)
}
