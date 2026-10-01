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
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
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

// spillAt sends one over-bound response from a process whose repository root is root, and
// answers its spill path below home.
func spillAt(t *testing.T, home, root string) string {
	t.Helper()
	var sink bytes.Buffer
	owner := responsebound.New(home, &sink, &sink, func() string { return root }, false)
	overBound(t, owner.Stdout())
	owner.Finish()
	return responseboundtest.Path(t, sink.String())
}

// BO16: the retirement of an assignment removes its spill directory, and a spill in the
// primary scope stays.
func TestRetirementDropsResponseSpills(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "spill-release")
	retired := spillAt(t, f.home, f.creation.Path)
	primary := spillAt(t, f.home, f.root)
	requireTest(t, filepath.Dir(retired) != filepath.Dir(primary), "assignment spill %s shares the primary scope of %s", retired, primary)
	release := runVerb(t, verbRelease, f.call("--request", "landed-spill-release", f.creation.Path))
	requireTest(t, release.exit == 0, "release = (%d, %q, %q)", release.exit, release.stdout, release.stderr)
	_, err := os.Stat(filepath.Dir(retired))
	requireTest(t, os.IsNotExist(err), "the release kept the assignment spill directory: %v", err)
	_, err = os.Stat(primary)
	requireTest(t, err == nil, "the release removed the primary spill: %v", err)
}

// BO72: a `release` that runs in the assignment it retires keeps its over-bound spill.
// The dispatcher marks the owner of a retiring leaf, as this test does. The spill starts
// before the retirement, so a spill in the assignment's own scope would leave the line
// naming a removed file. The release response reaches the owner after the retirement.
func TestRetiringVerbSpillsToPrimary(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "spill-own-release")
	var sink bytes.Buffer
	owner := responsebound.New(f.home, &sink, &sink, func() string { return f.creation.Path }, true)
	overBound(t, owner.Stdout())
	release := runVerb(t, verbRelease, f.call("--request", "landed-spill-own-release", f.creation.Path))
	_, stdoutErr := io.WriteString(owner.Stdout(), release.stdout)
	_, stderrErr := io.WriteString(owner.Stderr(), release.stderr)
	owner.Finish()
	requireTest(t, release.exit == 0 && stdoutErr == nil && stderrErr == nil, "release = (%d, %q), writes %v, %v", release.exit, sink.String(), stdoutErr, stderrErr)
	path := responseboundtest.Path(t, sink.String())
	requireTest(t, !strings.Contains(path, f.creation.Assignment.ID), "spill %s lies in the retired assignment's scope", path)
	data, err := os.ReadFile(path)
	requireTest(t, err == nil, "the release removed its own spill: %v", err)
	requireTest(t, strings.HasPrefix(string(data), "line 01\n") && strings.Count(string(data), "\n") > 11, "spill = %q, want the complete release response", data)
}
