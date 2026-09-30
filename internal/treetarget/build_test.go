package treetarget

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
	"github.com/gibbonmi/bench/internal/worktree"
)

// The refusal lines, the repair command, and the child facts here are authored apart from
// Run: the spec fixes both refusal lines, the quoting of the label, the published build as
// the child, and the child environment. A started child, a changed word, an unquoted label,
// or an inherited run binary reds these rows. The build is a marker script that
// freshness.Publish seals, so each row grades the real seal.

// kitTarget is a label target whose worktree holds a kit tree.
type kitTarget struct {
	targetRepo
	// tree is the physical root of the labeled worktree.
	tree string
}

func newKitTarget(t *testing.T, label string) kitTarget {
	t.Helper()
	r := newTargetRepo(t, t.TempDir())
	tree := physicalPath(t, r.create(t, "tree-target-kit", label).Worktree)
	kittest.WriteTree(t, tree)
	return kitTarget{targetRepo: r, tree: tree}
}

// publish publishes a marker script as the worktree build, and answers the build and its
// marker.
func (k kitTarget) publish(t *testing.T) (build, marker string) {
	t.Helper()
	script, marker := markerScript(t, "exit 0")
	build = freshness.PublishedExecutable(k.tree)
	if err := freshness.Publish(k.tree, script, build, filepath.Dir(build), "tree-target"); err != nil {
		t.Fatal(err)
	}
	return build, marker
}

// appendByte appends one byte to the file at path.
func appendByte(t *testing.T, path string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

// TT43, TT44, TT45, TT46, TT47, TT48: a kit worktree target runs its own worktree build only
// after freshness.Verify accepts that build. A missing, stale, or changed build refuses at
// exit 1 before any child starts, and the refusal names the repair by label. The primary
// checkout never runs a worktree build.
func TestRunKitWorktreeBuild(t *testing.T) {
	t.Run("TT43 TT48 a current build", func(t *testing.T) {
		k := newKitTarget(t, "alpha")
		_, buildMarker := k.publish(t)
		wrapper, wrapperMarker := markerScript(t, "exit 0")
		running, runningMarker := markerScript(t, "exit 0")
		if stdout, stderr, code := k.run(wrapper, running, "status", "alpha"); code != 0 || stdout != "" || stderr != "" {
			t.Fatalf("status --in alpha = (%d, %q, %q), want the build's exit 0 and no parent output", code, stdout, stderr)
		}
		requireNoChild(t, wrapperMarker)
		requireNoChild(t, runningMarker)
		child := readMarker(t, buildMarker)
		if child.argv != "argv status" || child.dir != k.tree {
			t.Fatalf("build = %q in %q, want argv status in %q", child.argv, child.dir, k.tree)
		}
		for _, name := range []string{runbinary.Env, "BENCH_KIT"} {
			if value, ok := child.env[name]; ok {
				t.Errorf("child %s = %q, want it unset", name, value)
			}
		}
		if got, want := child.env[worktree.WrapperEnv], filepath.Join(k.tree, kittest.Wrapper); got != want {
			t.Errorf("child %s = %q, want the wrapper of the target %q", worktree.WrapperEnv, got, want)
		}
	})
	refusal := func(reason, label string) string {
		return "bench status --in: " + reason + "\nnext=bench worktree build " + label + "\n"
	}
	missing, stale := "worktree build is missing", "worktree build does not match the tree"
	for _, row := range []struct {
		name, label string
		// change edits the published build or its tree, and a row without it publishes no build.
		change func(t *testing.T, k kitTarget, build string)
		stderr string
	}{
		{name: "TT44 no build", label: "alpha", stderr: refusal(missing, "alpha")},
		{name: "TT45 a changed build input", label: "alpha", stderr: refusal(stale, "alpha"),
			change: func(t *testing.T, k kitTarget, _ string) { kittest.EditBuildInput(t, k.tree) }},
		{name: "TT46 one appended byte", label: "alpha", stderr: refusal(stale, "alpha"),
			change: func(t *testing.T, _ kitTarget, build string) { appendByte(t, build) }},
		{name: "TT47 a label that needs quoting", label: "my alpha", stderr: refusal(missing, "'my alpha'")},
	} {
		t.Run(row.name, func(t *testing.T) {
			k := newKitTarget(t, row.label)
			wrapper, wrapperMarker := markerScript(t, "exit 0")
			running, runningMarker := markerScript(t, "exit 0")
			markers := []string{wrapperMarker, runningMarker}
			if row.change != nil {
				build, buildMarker := k.publish(t)
				row.change(t, k, build)
				markers = append(markers, buildMarker)
			}
			if stdout, stderr, code := k.run(wrapper, running, "status", row.label); code != 1 || stdout != "" || stderr != row.stderr {
				t.Fatalf("status --in %q = (%d, %q, %q), want (1, \"\", %q)", row.label, code, stdout, stderr, row.stderr)
			}
			for _, marker := range markers {
				requireNoChild(t, marker)
			}
		})
	}
	t.Run("a kit primary checkout runs the wrapper", func(t *testing.T) {
		r := newTargetRepo(t, t.TempDir())
		kittest.WriteTree(t, r.root)
		wrapper, wrapperMarker := markerScript(t, "exit 0")
		if stdout, stderr, code := r.run(wrapper, "", "status", "primary"); code != 0 {
			t.Fatalf("status --in primary = (%d, %q, %q), want the wrapper's exit 0", code, stdout, stderr)
		}
		if child := readMarker(t, wrapperMarker); child.argv != "argv status" {
			t.Fatalf("wrapper argv = %q, want argv status", child.argv)
		}
	})
}
