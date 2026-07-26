package preprelease

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/git"
)

// TestPrepReleaseWritesEvidence is the story 5 acceptance: ship green is exit 0 with the
// named index and the full artifact set. The artifact count is read from the plan the
// build step reads rather than written down here, so a route that exits 0 without
// building anything cannot pass.
func TestPrepReleaseWritesEvidence(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	probe := r.prepRelease(nil)
	probe.RequireExit(0)

	if !r.Exists(evidenceIndexPath) {
		t.Fatalf("ship green left no %s\nstderr:\n%s", evidenceIndexPath, probe.Stderr)
	}
	if got, want := countFiles(t, filepath.Join(r.Root, filepath.FromSlash(evidenceArtifacts))), planTargetCount(t, r.Root); got != want {
		t.Fatalf("artifact set = %d files, want one per planned target (%d)", got, want)
	}
	requireContains(t, "preflight argv", r.ReadFile(preflightArgvFile), "--mode\nverify")
}

// TestPrepReleaseRunsStressTags is the story 5 cross-compile row. crossCompileMatrix is
// a no-op without the tag, so a forgotten -tags stress runs a check that silently
// returns nil; the conformance entry point records which way it was compiled.
func TestPrepReleaseRunsStressTags(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	r.prepRelease(nil).RequireExit(0)

	record := r.ReadFile(conformanceFile)
	requireContains(t, "conformance invocation", record, "stress=on")
	requireContains(t, "conformance invocation", record, "tier=ship")
}

// TestPrepReleaseRequiresDevGreen is the story 6 acceptance. All four rejected cache
// states are enumerated, because a file-existence check would accept a red or
// tree-mismatched record and claim ship green over a tree the dev tier never passed.
func TestPrepReleaseRequiresDevGreen(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)

	for _, test := range []struct {
		name  string
		spoil func(*testing.T, shipRepo)
		cause string
	}{
		{
			name:  "absent",
			spoil: func(t *testing.T, r shipRepo) { removeVerdict(t, r) },
			cause: "absent",
		},
		{
			name: "red",
			spoil: func(t *testing.T, r shipRepo) {
				r.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
				r.CommitAll("red gate")
				r.RunEnv(r.runEnv(nil), "bash", r.benchScript(), "gate").RequireExit(1)
			},
			cause: "recorded red",
		},
		{
			name: "bound to a different tree",
			spoil: func(t *testing.T, r shipRepo) {
				r.WriteFile("drift.txt", "the tree moved after the verdict was recorded\n")
				r.CommitAll("drift")
			},
			cause: "working tree changed",
		},
		{
			name:  "stale",
			spoil: func(t *testing.T, r shipRepo) { ageVerdict(t, r) },
			cause: "verdict expired",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := newShipRepo(t)
			test.spoil(t, r)

			probe := r.prepRelease(nil)
			if probe.ExitCode == 0 {
				t.Fatalf("prep-release accepted a %s verdict\nstderr:\n%s", test.name, probe.Stderr)
			}
			requireContains(t, "refusal", probe.Stderr, test.cause)
			requireContains(t, "refusal", probe.Stderr, "bench gate")
			if r.Exists(evidenceIndexPath) {
				t.Fatalf("a refused run still produced %s", evidenceIndexPath)
			}
		})
	}
}

// TestPrepReleaseAllSurfaces is the story 8 acceptance: a route added to bin/bench.sh
// alone would pass a single-surface test. The three shipped paths are the kit CLI a
// maintainer types, the by-path CLI a linked repo carries, and the compiled binary,
// which is what every route_binary case and every hook ultimately execs. The refusal is
// the observation because it needs no toolchain: identical bytes from all three is what
// one implementation looks like.
func TestPrepReleaseAllSurfaces(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	kit := contract.SubjectRoot(t)
	f := contract.NewFixture(t)
	f.BenchEnv(map[string]string{"BENCH_KIT": kit}, "link").RequireExit(0)

	var first string
	for _, surface := range []struct {
		name string
		argv []string
	}{
		{"kit CLI", []string{"bash", filepath.Join(kit, "bin", "bench.sh"), "prep-release"}},
		{"linked by-path CLI", []string{"bash", filepath.Join(f.Root, ".bench", "bin", "bench.sh"), "prep-release"}},
		{"compiled binary", []string{filepath.Join(kit, "dist", "bench"), "prep-release"}},
	} {
		probe := f.Run(surface.argv[0], surface.argv[1:]...)
		if probe.ExitCode == 0 {
			t.Fatalf("%s: prep-release exited 0 with no dev-green verdict", surface.name)
		}
		requireContains(t, surface.name+" refusal", probe.Stderr, "no current dev-green verdict")
		if first == "" {
			first = probe.Stderr
			continue
		}
		if probe.Stderr != first {
			t.Fatalf("%s reached a different implementation:\n%q\nversus\n%q", surface.name, probe.Stderr, first)
		}
	}
}

// TestPrepReleaseMissingTool is the edge row for a required tool absent from PATH. A
// thin route that just shells out surfaces an opaque `exec: not found` from inside a
// script instead of naming what the maintainer has to install.
func TestPrepReleaseMissingTool(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)

	for _, tool := range []string{"go", "node"} {
		t.Run(tool, func(t *testing.T) {
			t.Parallel()
			r := newShipRepo(t)
			binary := filepath.Join(contract.SubjectRoot(t), "dist", "bench")

			probe := r.RunEnv(r.runEnv(map[string]string{"PATH": pathWithout(t, tool)}), binary, "prep-release")
			if probe.ExitCode == 0 {
				t.Fatalf("prep-release ran with %s off PATH", tool)
			}
			requireContains(t, "missing-tool diagnostic", probe.Stderr, "required tool is missing or not executable: "+tool)
		})
	}
}

// TestPrepReleaseHostilePath is the edge row for a repo root carrying a space and a glob
// character. The command hands that root to four shell scripts, where an unquoted
// expansion is this domain's first-named failure.
func TestPrepReleaseHostilePath(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	root := filepath.Join(t.TempDir(), "ship repo [v1]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("create hostile-path root: %v", err)
	}
	r := newShipRepoAt(t, root)

	probe := r.prepRelease(nil)
	probe.RequireExit(0)

	if !r.Exists(evidenceIndexPath) {
		t.Fatalf("hostile-path root produced no %s\nstderr:\n%s", evidenceIndexPath, probe.Stderr)
	}
	requireContains(t, "artifacts argv", r.ReadFile(artifactsArgvFile), resolved(t, root))
}

// TestPrepReleaseInterrupt is the edge row for SIGINT partway through. Evidence lands by
// atomic directory exchange, so an interrupt taken between staging and promotion must
// leave no index at all — a route that bypassed the exchange with a direct write would
// leave a partial one.
func TestPrepReleaseInterrupt(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	code := interruptDuringPreflight(t, r)
	if code == 0 {
		t.Fatal("interrupted run reported ship green")
	}
	if r.Exists(evidenceIndexPath) {
		t.Fatalf("an interrupted run left %s behind", evidenceIndexPath)
	}
}

// TestPrepReleaseIdempotent is the edge row for a second run on an unchanged tree: a
// route that appended to dist/ rather than promoting into it fails on the second run.
func TestPrepReleaseIdempotent(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	r.prepRelease(nil).RequireExit(0)
	first := countFiles(t, filepath.Join(r.Root, filepath.FromSlash(evidenceArtifacts)))
	firstIndex := r.ReadFile(evidenceIndexPath)

	r.prepRelease(nil).RequireExit(0)

	if second := countFiles(t, filepath.Join(r.Root, filepath.FromSlash(evidenceArtifacts))); second != first {
		t.Fatalf("second run produced %d artifacts, want the same %d", second, first)
	}
	if got := r.ReadFile(evidenceIndexPath); got != firstIndex {
		t.Fatalf("second run produced different evidence:\n%q\nversus\n%q", got, firstIndex)
	}
}

// TestPrepReleaseDeepCwd is the edge row for invocation below the repo root: every
// script argument is derived from the root, so a command that took the working
// directory for the root would hand three scripts a subdirectory.
func TestPrepReleaseDeepCwd(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	probe := contract.RunAt(t, r.Fixture, filepath.Join(r.Root, "scripts"), r.runEnv(nil), "bash", r.benchScript(), "prep-release")
	probe.RequireExit(0)

	if !r.Exists(evidenceIndexPath) {
		t.Fatalf("a deep-cwd run produced no %s\nstderr:\n%s", evidenceIndexPath, probe.Stderr)
	}
	if got, want := strings.SplitN(r.ReadFile(artifactsArgvFile), "\n", 2)[0], resolved(t, r.Root); got != want {
		t.Fatalf("build step received source root %q, want the repo root %q", got, want)
	}
}

// TestPrepReleaseDoesNotFalsifyItsOwnPrecheck is the edge row for the self-falsifying
// write: the verdict the run read is bound to a tree hash, so everything the run writes
// has to stay outside the subject. A clean tree afterwards is the observation, and a
// second run reusing the same verdict is the consequence.
func TestPrepReleaseDoesNotFalsifyItsOwnPrecheck(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	r := newShipRepo(t)

	r.prepRelease(nil).RequireExit(0)

	if status := r.Git("status", "--porcelain=v1", "--untracked-files=all"); strings.TrimSpace(status.Stdout) != "" {
		t.Fatalf("prep-release changed the subject it graded:\n%s", status.Stdout)
	}
	r.prepRelease(nil).RequireExit(0)
}

func removeVerdict(t *testing.T, r shipRepo) {
	t.Helper()
	if err := os.Remove(verdictPath(r)); err != nil && !os.IsNotExist(err) {
		t.Fatalf("remove verdict cache: %v", err)
	}
}

var recordedAtRE = regexp.MustCompile(`"recorded_at":"[^"]*"`)

// ageVerdict rewrites the recorded timestamp past the gate's freshness window. The
// record is edited rather than rebuilt so every other field stays exactly what the gate
// wrote, leaving expiry as the only thing the command can be reacting to.
func ageVerdict(t *testing.T, r shipRepo) {
	t.Helper()
	path := verdictPath(r)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read verdict cache: %v", err)
	}
	aged := time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339)
	stamped := recordedAtRE.ReplaceAllString(string(data), `"recorded_at":"`+aged+`"`)
	if stamped == string(data) {
		t.Fatalf("verdict cache carries no recorded_at to age:\n%s", data)
	}
	if err := os.WriteFile(path, []byte(stamped), 0o600); err != nil {
		t.Fatalf("write aged verdict cache: %v", err)
	}
}

func verdictPath(r shipRepo) string {
	return filepath.Join(r.Root, ".git", git.GateCacheFile)
}

// resolved is the path form the command's own root resolution produces, so an assertion
// against an argv compares like with like on hosts where the temp tree is symlinked.
func resolved(t *testing.T, path string) string {
	t.Helper()
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("resolve %s: %v", path, err)
	}
	return real
}
