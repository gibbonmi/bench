package worktree

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// installSpec is the deliverable that an install fixture binds, and installDeclared is the
// header line that declares one persisted ledger field in it.
const (
	installSpec     = "specs/install/spec.md"
	installDeclared = "Persisted fields: intent-ledger:commitment.fixture\n"
)

// installShape is one install repository: whether it declares Bench build inputs, the
// declaration line of the default-branch spec, and the deliverable that the assignment
// binds. An empty deliverable binds nothing.
type installShape struct {
	inputs      bool
	declaration string
	deliverable string
}

// declaredInstall is the shape that each install refusal admits.
var declaredInstall = installShape{inputs: true, declaration: installDeclared, deliverable: installSpec}

// newInstallAssignment builds a repository whose default branch holds a Go main package,
// a wrapper, and installSpec, and one owned assignment that starts at that commit.
func newInstallAssignment(t *testing.T, request string, shape installShape) ownedAssignment {
	t.Helper()
	root := newWorktreeRepo(t)
	writeGoMainFixture(t, root)
	mustMkdirAll(t, wrapperDir(root), 0o755)
	mustWrite(t, wrapperPath(root), []byte("#!/bin/sh\nexit 0\n"), 0o755)
	mustMkdirAll(t, filepath.Dir(filepath.Join(root, installSpec)), 0o755)
	mustWrite(t, filepath.Join(root, installSpec), []byte("# install\n\nStatus: staged\n"+shape.declaration), 0o644)
	gitRun(t, root, "add", "go.mod", "cmd", "scripts", "bin", installSpec)
	gitRun(t, root, withIdentity("commit", "-qm", "install sources")...)
	if shape.inputs {
		commitLandingBuildInputs(t, root)
	}
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, request, request)
	if shape.deliverable != "" {
		commitmenttest.EditCommitment(t, root, func(state *intent.CommitmentState) {
			state.Bindings = append(state.Bindings, intent.DeliveryBinding{Assignment: creation.Assignment.ID, Request: creation.Assignment.Request,
				Milestone: "fixture", Outcome: "install", Deliverable: shape.deliverable, Identity: "fixture"})
		})
	}
	return ownedAssignment{repoHome: repoHome{root, home}, creation: creation}
}

// doctorCall is one doctorChild call, with the primary executable's bytes when it ran.
type doctorCall struct {
	dir       string
	argv, env []string
	published string
}

// doctorRecorder is the doctorChild join a row drives instead of the repair. It answers exit.
type doctorRecorder struct {
	calls []doctorCall
	exit  int
}

func (r *doctorRecorder) join(dir string, argv, env []string, _ io.Writer) int {
	published, _ := os.ReadFile(freshness.PublishedExecutable(dir))
	r.calls = append(r.calls, doctorCall{dir: dir, argv: argv, env: env, published: string(published)})
	return r.exit
}

// install runs the install form of the build verb on f's assignment, with the build and
// the repair recorded, and with running as the executable that this process runs. Each
// edit then changes one more seam.
func install(t *testing.T, f ownedAssignment, build *buildRecorder, doctor *doctorRecorder, running string, edits ...func(*joins)) verbResult {
	t.Helper()
	j := buildJoins(build)
	j.resolveRunningBinary = func() (string, error) { return running, nil }
	j.doctorChild = doctor.join
	for _, edit := range edits {
		edit(&j)
	}
	return runVerb(t, verbBuild, f.callWith(j, "--install", f.creation.Assignment.Label))
}

// requireInstalled fails unless the primary executable of f holds body, its digest equals
// the target's executable digest, and the primary broker manifest binds it by that digest.
func requireInstalled(t *testing.T, f ownedAssignment, body string) string {
	t.Helper()
	primary, built := freshness.PublishedExecutable(f.root), freshness.PublishedExecutable(f.creation.Path)
	requireTest(t, fixtureFileText(t, primary) == body, "primary executable = %q, want the candidate %q", fixtureFileText(t, primary), body)
	digest, err := brokermanifest.Digest(primary)
	mustNoError(t, err)
	builtDigest, err := brokermanifest.Digest(built)
	mustNoError(t, err)
	requireTest(t, digest == builtDigest, "primary digest %s, worktree digest %s; want one digest", digest, builtDigest)
	fields, err := brokermanifest.Read(manifestPath(f.root))
	mustNoError(t, err)
	requireTest(t, fields["path"] == primary && fields["sha256"] == digest, "primary manifest = %v, want path %s and digest %s", fields, primary, digest)
	return digest
}

// plantCommitmentField writes a commitment field that this schema does not know into the
// ledger at root, as a build that adds a ledger field leaves the shared ledger.
func plantCommitmentField(t *testing.T, root string) {
	t.Helper()
	commitmenttest.PlantCommitmentField(t, root, "fixture")
	_, strict := intent.Read(root)
	requireTest(t, strict != nil, "the strict read accepts the planted field, so the row grades nothing")
}

// TestBuildInstall is CB13, CB14, CB16 to CB21, CB54 to CB56, and CB63. The install form
// builds the target, publishes its sealed bytes as the primary executable, and then runs
// the doctor repair. Each refusal comes before the build and writes no file.
func TestBuildInstall(t *testing.T) {
	t.Parallel()
	const candidate = "candidate-executable"
	t.Run("help names the install form", func(t *testing.T) {
		t.Parallel()
		r := runVerb(t, verbBuild, verbCall{root: t.TempDir(), home: t.TempDir(), args: []string{"--help"}})
		requireTest(t, r.exit == 0 && strings.Contains(r.stdout, "usage: "+usage.WorktreeBuild+"\n") && strings.Contains(r.stdout, usage.WorktreeBuildInstall), "help = (%d, %q), want both build forms", r.exit, r.stdout)
	})
	t.Run("publishes the sealed candidate and then repairs", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-publishes", declaredInstall)
		commitInWorktree(t, f.creation.Path, "cmd/bench/main.go", "package main\n\nfunc main() { println(\"candidate\") }\n", "candidate source")
		doctor := &doctorRecorder{}
		primary := freshness.PublishedExecutable(f.root)
		r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, doctor, primary)
		requireTest(t, r.exit == 0, "install = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
		digest := requireInstalled(t, f, candidate)
		requireTest(t, freshness.Verify(f.creation.Path, primary) == nil, "the primary seal does not verify against the target tree: %v", freshness.Verify(f.creation.Path, primary))
		requireTest(t, freshness.Verify(f.root, primary) != nil, "the primary seal verifies against the primary sources, want the target tree's sources")
		table, err := toon.Table(installTable, []string{"worktree", "sha256", "manifest"}, [][]string{{f.creation.Assignment.ID, digest, manifestPath(f.root)}})
		mustNoError(t, err)
		requireTest(t, r.stdout == table+"next[1]:\n  bench record assignment 'install' --help\n", "install printed %q, want the table %q and the dispatch of the bound spec", r.stdout, table)
		requireTest(t, len(doctor.calls) == 1, "doctor calls = %d, want 1", len(doctor.calls))
		call := doctor.calls[0]
		wantArgv := []string{wrapperPath(f.root), "doctor", "--fix"}
		requireTest(t, call.dir == f.root && slices.Equal(call.argv, wantArgv) && call.published == candidate, "doctor call = %+v, want %v in %s after the publication", call, wantArgv, f.root)
		var selected []string
		for _, entry := range call.env {
			requireTest(t, !strings.HasPrefix(entry, "BENCH_KIT=") && !strings.HasPrefix(entry, "BENCH_WRAPPER="), "doctor environment holds %q", entry)
			if value, ok := strings.CutPrefix(entry, runbinary.Env+"="); ok {
				selected = append(selected, value)
			}
		}
		requireTest(t, slices.Equal(selected, []string{primary}), "doctor selections = %q, want [%s]", selected, primary)
	})
	t.Run("the repair environment drops the routing and selects the primary executable", func(t *testing.T) {
		t.Parallel()
		base := []string{"BENCH_KIT=/kit", "BENCH_WRAPPER=/kit/bin/bench.sh", runbinary.Env + "=/kit/dist/bench", "KEEP=1"}
		got := installChildEnv(base, "/primary/dist/bench")
		requireTest(t, slices.Equal(got, []string{"KEEP=1", runbinary.Env + "=/primary/dist/bench"}), "repair environment = %q", got)
	})
	t.Run("a second install replaces both copies", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-replaces", declaredInstall)
		primary := freshness.PublishedExecutable(f.root)
		first := install(t, f, &buildRecorder{body: "first-candidate", version: "1.2.3"}, &doctorRecorder{}, primary)
		requireTest(t, first.exit == 0, "first install = (%d, %q)", first.exit, first.stderr)
		firstDigest := requireInstalled(t, f, "first-candidate")
		commitInWorktree(t, f.creation.Path, "cmd/bench/main.go", "package main\n\nfunc main() { println(\"second\") }\n", "second source")
		second := install(t, f, &buildRecorder{body: "second-candidate", version: "1.2.3"}, &doctorRecorder{}, primary)
		requireTest(t, second.exit == 0, "second install = (%d, %q)", second.exit, second.stderr)
		requireTest(t, requireInstalled(t, f, "second-candidate") != firstDigest, "the second install kept the first digest")
		requireTest(t, freshness.Verify(f.creation.Path, primary) == nil, "the second seal does not verify against the new commit: %v", freshness.Verify(f.creation.Path, primary))
	})
	t.Run("refuses before any write", func(t *testing.T) {
		t.Parallel()
		undeclared := installShape{inputs: true, deliverable: installSpec}
		for _, row := range []struct {
			name     string
			shape    installShape
			prepare  func(*testing.T, ownedAssignment)
			worktree bool
			want     func(ownedAssignment) string
		}{
			{name: "no build inputs", shape: installShape{declaration: installDeclared, deliverable: installSpec}, want: func(ownedAssignment) string { return "declares no Bench build inputs" }},
			{name: "a worktree executable runs", shape: declaredInstall, worktree: true, want: func(f ownedAssignment) string {
				return "the running executable " + freshness.PublishedExecutable(f.creation.Path) + " is not the primary executable"
			}},
			{name: "no bound assignment", shape: installShape{inputs: true, declaration: installDeclared}, want: func(ownedAssignment) string { return "has no bound deliverable" }},
			{name: "an undeclared default-branch spec", shape: undeclared, want: func(ownedAssignment) string { return "declares no persisted field" }},
			{name: "a declaration in the target tree alone", shape: undeclared, prepare: func(t *testing.T, f ownedAssignment) {
				commitInWorktree(t, f.creation.Path, installSpec, "# install\n\nStatus: staged\n"+installDeclared, "declare in the candidate")
			}, want: func(ownedAssignment) string { return "declares no persisted field" }},
			{name: "no default-branch copy", shape: installShape{inputs: true, deliverable: "specs/absent/spec.md"}, want: func(ownedAssignment) string { return "default-branch copy of specs/absent/spec.md" }},
			{name: "an invalid declaration", shape: installShape{inputs: true, declaration: "Persisted fields: gate-evidence:x\n", deliverable: installSpec}, want: func(ownedAssignment) string { return spec.ErrPersistedFieldsInvalid.Error() }},
			{name: "a retained assignment whose tree is gone", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) {
				mustNoError(t, os.RemoveAll(f.creation.Path))
			}, want: func(ownedAssignment) string { return refusalroute.Sentence(faceResetTreeMissing) }},
		} {
			t.Run(row.name, func(t *testing.T) {
				t.Parallel()
				f := newInstallAssignment(t, "install-refuses-"+strings.ReplaceAll(row.name, " ", "-"), row.shape)
				if row.prepare != nil {
					row.prepare(t, f)
				}
				running := freshness.PublishedExecutable(f.root)
				if row.worktree {
					running = freshness.PublishedExecutable(f.creation.Path)
				}
				build, doctor := &buildRecorder{body: candidate, version: "1.2.3"}, &doctorRecorder{}
				r := install(t, f, build, doctor, running)
				requireTest(t, r.exit == 1 && strings.Contains(r.stderr, row.want(f)), "install = (%d, %q), want exit 1 naming %q", r.exit, r.stderr, row.want(f))
				requireTest(t, len(build.calls) == 0 && len(doctor.calls) == 0, "build calls %d, doctor calls %d; want none", len(build.calls), len(doctor.calls))
				for _, path := range []string{filepath.Dir(freshness.PublishedExecutable(f.root)), manifestPath(f.root)} {
					_, err := os.Lstat(path)
					requireTest(t, os.IsNotExist(err), "refused install left %s: %v", path, err)
				}
			})
		}
	})
	t.Run("finds the binding past an unknown commitment field", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-tolerant", declaredInstall)
		plantCommitmentField(t, f.root)
		commitmenttest.RequireLedgerKept(t, f.root, func() {
			r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, &doctorRecorder{}, freshness.PublishedExecutable(f.root))
			requireTest(t, r.exit == 0, "install = (%d, %q), want exit 0", r.exit, r.stderr)
		})
		requireInstalled(t, f, candidate)
	})
	t.Run("a failed publication keeps the earlier set", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-publication-fails", declaredInstall)
		primary := freshness.PublishedExecutable(f.root)
		mustMkdirAll(t, filepath.Dir(primary), 0o755)
		mustWrite(t, primary, []byte("earlier-executable"), 0o755)
		mustWrite(t, primary+".seal", []byte("earlier-seal"), 0o644)
		// The manifest directory denies writes, so the owner fails at the manifest, after it
		// replaced the executable.
		mustNoError(t, os.Chmod(wrapperDir(f.root), 0o555))
		t.Cleanup(func() { _ = os.Chmod(wrapperDir(f.root), 0o755) })
		if probe, err := os.CreateTemp(wrapperDir(f.root), "probe-*"); err == nil {
			_ = probe.Close()
			_ = os.Remove(probe.Name())
			capability.Capability(t, capability.Privilege, "a 0o555 directory still takes writes from this user, so the publication cannot fail")
		}
		doctor := &doctorRecorder{}
		r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, doctor, primary)
		requireTest(t, r.exit == 1 && strings.Contains(r.stderr, "publication refused"), "install = (%d, %q), want the publication refusal", r.exit, r.stderr)
		requireTest(t, fixtureFileText(t, primary) == "earlier-executable" && fixtureFileText(t, primary+".seal") == "earlier-seal", "the failed publication left %q and seal %q", fixtureFileText(t, primary), fixtureFileText(t, primary+".seal"))
		_, err := os.Lstat(manifestPath(f.root))
		requireTest(t, os.IsNotExist(err) && len(doctor.calls) == 0, "manifest stat %v, doctor calls %d; want no manifest and no repair", err, len(doctor.calls))
		entries, err := os.ReadDir(filepath.Dir(primary))
		mustNoError(t, err)
		for _, entry := range entries {
			requireTest(t, entry.Name() == "bench" || entry.Name() == "bench.seal", "the failed publication left %s", entry.Name())
		}
	})
	t.Run("a failed publication keeps an earlier manifest", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-publication-keeps", declaredInstall)
		primary := freshness.PublishedExecutable(f.root)
		mustMkdirAll(t, filepath.Dir(primary), 0o755)
		earlier := map[string]string{primary: "earlier-executable", primary + ".seal": "earlier-seal", manifestPath(f.root): "earlier-manifest"}
		for path, body := range earlier {
			mustWrite(t, path, []byte(body), 0o644)
		}
		doctor := &doctorRecorder{}
		published := false
		// The seam answers as the owner answers a failure after its executable write: the
		// owner has already restored the earlier set, so only the install can change it.
		failAfterWrite := func(j *joins) {
			j.publishCandidate = func(string, string, string, string, string) error {
				published = true
				return errors.New("seal promotion failed after the executable write")
			}
		}
		r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, doctor, primary, failAfterWrite)
		requireTest(t, published && r.exit == 1 && strings.Contains(r.stderr, "publication refused"), "install = (%d, %q), want the publication refusal", r.exit, r.stderr)
		for path, body := range earlier {
			requireTest(t, fixtureFileText(t, path) == body, "the failed publication left %s = %q, want %q", path, fixtureFileText(t, path), body)
		}
		requireTest(t, len(doctor.calls) == 0, "doctor calls = %d, want no repair", len(doctor.calls))
	})
	t.Run("accepts a symlinked alias of the primary executable", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-alias", declaredInstall)
		primary := freshness.PublishedExecutable(f.root)
		mustMkdirAll(t, filepath.Dir(primary), 0o755)
		mustWrite(t, primary, []byte("earlier-executable"), 0o755)
		alias := filepath.Join(t.TempDir(), "bench")
		mustNoError(t, os.Symlink(primary, alias))
		r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, &doctorRecorder{}, alias)
		requireTest(t, r.exit == 0, "install under the alias %s = (%d, %q), want exit 0", alias, r.exit, r.stderr)
		requireInstalled(t, f, candidate)
	})
	t.Run("a failed repair keeps the candidate published", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "install-repair-fails", declaredInstall)
		doctor := &doctorRecorder{exit: 3}
		r := install(t, f, &buildRecorder{body: candidate, version: "1.2.3"}, doctor, freshness.PublishedExecutable(f.root))
		requireTest(t, r.exit != 0 && strings.Contains(r.stderr, "the doctor repair of the primary checkout exited 3"), "install = (%d, %q), want the repair refusal", r.exit, r.stderr)
		requireTest(t, len(doctor.calls) == 1, "doctor calls = %d, want 1", len(doctor.calls))
		requireInstalled(t, f, candidate)
	})
}
