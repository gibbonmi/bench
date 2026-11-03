package treetarget

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
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
	build = freshness.PublishedExecutable(k.tree)
	return build, k.publishAt(t, build)
}

// publishAt publishes a marker script at build with a seal of the target tree, and answers
// its marker.
func (k kitTarget) publishAt(t *testing.T, build string) (marker string) {
	t.Helper()
	script, marker := markerScript(t, "exit 0")
	if err := os.MkdirAll(filepath.Dir(build), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(k.tree, script, build, filepath.Dir(build), "tree-target"); err != nil {
		t.Fatal(err)
	}
	return marker
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

// The two build refusals, and the stderr of a status call that refuses with reason for label.
const buildMissing, buildStale = "worktree build is missing", "worktree build does not match the tree"

func buildRefusal(reason, label string) string {
	return "bench status --in: " + reason + "\nnext=bench worktree build " + label + "\n"
}

// TT43, TT44, TT45, TT46, TT47, TT48: these rows hold no other matching candidate, so a kit
// worktree target runs its own worktree build only after freshness.Verify accepts that
// build. A missing, stale, or changed build refuses at exit 1 before any child starts, and
// the refusal names the repair by label. The primary checkout never runs a worktree build.
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
		if got, want := child.env[env.WrapperEnv], filepath.Join(k.tree, kittest.Wrapper); got != want {
			t.Errorf("child %s = %q, want the wrapper of the target %q", env.WrapperEnv, got, want)
		}
		if got := child.env[benchhome.Env]; got != k.home {
			t.Errorf("child %s = %q, want the resolved home %q", benchhome.Env, got, k.home)
		}
	})
	for _, row := range []struct {
		name, label string
		// change edits the published build or its tree, and a row without it publishes no build.
		change func(t *testing.T, k kitTarget, build string)
		stderr string
	}{
		{name: "TT44 no build", label: "alpha", stderr: buildRefusal(buildMissing, "'alpha'")},
		{name: "TT45 a changed build input", label: "alpha", stderr: buildRefusal(buildStale, "'alpha'"),
			change: func(t *testing.T, k kitTarget, _ string) { kittest.EditBuildInput(t, k.tree) }},
		{name: "TT46 one appended byte", label: "alpha", stderr: buildRefusal(buildStale, "'alpha'"),
			change: func(t *testing.T, _ kitTarget, build string) { appendByte(t, build) }},
		{name: "TT47 a label that needs quoting", label: "my alpha", stderr: buildRefusal(buildMissing, "'my alpha'")},
		// A line with no control rune prints as is, so the pasted repair names this label.
		{name: "TT47 a label with a backslash", label: `a\b`, stderr: buildRefusal(buildMissing, `'a\b'`)},
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

// candidateSeal is the state of one candidate executable of a kit worktree target.
type candidateSeal int

const (
	// absent is no file, or for the running executable a file with no seal.
	absent candidateSeal = iota
	// matching is a file whose seal matches the target tree.
	matching
	// changed is a sealed file with one appended byte, so its seal does not match.
	changed
)

// BC61, BC62, BC63, BC65: a kit worktree target runs the first candidate whose seal matches
// the target tree, in the order of its own build, the primary checkout's build, and the
// running executable. When no candidate matches, the missing refusal stays, and no
// candidate and no wrapper starts.
func TestRunReusesAMatchingBuild(t *testing.T) {
	const own, primary, running = 0, 1, 2
	for _, row := range []struct {
		name  string
		seals [3]candidateSeal
		// runs is the candidate that runs, or -1 when the call refuses.
		runs int
		// edit edits a build input of the target after the seals, so each intact seal holds
		// other sources.
		edit bool
	}{
		{name: "BC61 the primary build before the running executable", seals: [3]candidateSeal{absent, matching, matching}, runs: primary},
		{name: "BC61 the primary build after a changed own build", seals: [3]candidateSeal{changed, matching, absent}, runs: primary},
		{name: "BC62 the running executable after a changed primary build", seals: [3]candidateSeal{absent, changed, matching}, runs: running},
		{name: "BC62 the running executable after a changed own build", seals: [3]candidateSeal{changed, absent, matching}, runs: running},
		{name: "BC63 the own build before the primary build", seals: [3]candidateSeal{matching, matching, matching}, runs: own},
		{name: "BC65 no matching candidate", seals: [3]candidateSeal{absent, changed, absent}, runs: -1},
		{name: "BC65 an intact primary build sealed for other sources", seals: [3]candidateSeal{absent, matching, absent}, runs: -1, edit: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			k := newKitTarget(t, "alpha")
			wrapper, wrapperMarker := markerScript(t, "exit 0")
			paths := [3]string{freshness.PublishedExecutable(k.tree), freshness.PublishedExecutable(k.root), filepath.Join(t.TempDir(), "bench")}
			var markers [3]string
			for i, seal := range row.seals {
				switch {
				case seal != absent:
					markers[i] = k.publishAt(t, paths[i])
					if seal == changed {
						appendByte(t, paths[i])
					}
				case i == running:
					paths[i], markers[i] = markerScript(t, "exit 0")
				}
			}
			if row.edit {
				kittest.EditBuildInput(t, k.tree)
			}
			stdout, stderr, code := k.run(wrapper, paths[running], "status", "alpha")
			requireNoChild(t, wrapperMarker)
			if row.runs < 0 {
				if want := buildRefusal(buildMissing, "'alpha'"); code != 1 || stdout != "" || stderr != want {
					t.Fatalf("status --in alpha = (%d, %q, %q), want (1, \"\", %q)", code, stdout, stderr, want)
				}
			} else if code != 0 || stdout != "" || stderr != "" {
				t.Fatalf("status --in alpha = (%d, %q, %q), want the child's exit 0 and no parent output", code, stdout, stderr)
			}
			for i, marker := range markers {
				switch {
				case i == row.runs:
					if child := readMarker(t, marker); child.argv != "argv status" || child.dir != k.tree {
						t.Fatalf("candidate %d = %q in %q, want argv status in %q", i, child.argv, child.dir, k.tree)
					}
				case marker != "":
					requireNoChild(t, marker)
				}
			}
		})
	}
}
