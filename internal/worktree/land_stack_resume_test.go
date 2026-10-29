// Stack resume tests: a stack whose release stopped names its resume, and the resume
// proves each published member before it finishes the stack.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// plantMemberResidue leaves an ignored file that no declaration covers in the worktree of
// member i of f, so the release of that member refuses, as in the unknown-ignored landing
// journey. The ignore rule sits in the common exclude file, so no member commits it. It
// answers the repair that removes the residue.
func plantMemberResidue(t *testing.T, f stackFixture, i int) (repair func()) {
	t.Helper()
	common := gitOutput(t, f.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	mustMkdirAll(t, filepath.Join(common, "info"), 0o755)
	mustWrite(t, filepath.Join(common, "info", "exclude"), []byte("private/\n"), 0o644)
	residue := filepath.Join(f.members[i].creation.Path, "private")
	mustMkdirAll(t, residue, 0o755)
	mustWrite(t, filepath.Join(residue, "output"), []byte("residue\n"), 0o600)
	return func() { mustRemove(t, residue) }
}

// stackResumeArgs is the stack resume of published over file with overrides as
// `--source-tip` values.
func stackResumeArgs(published, file string, overrides ...string) []string {
	args := []string{"--resume", published, "--stack", file}
	for _, o := range overrides {
		args = append(args, "--source-tip", o)
	}
	return args
}

// interruptedStack runs the stack of f with overrides while member stopped holds residue,
// requires the exit-3 record, and answers the stack file, the published commit, and the
// residue repair.
func interruptedStack(t *testing.T, f stackFixture, stopped int, overrides ...string) (file, published string, repair func()) {
	t.Helper()
	repair = plantMemberResidue(t, f, stopped)
	file, args := f.stack(t, overrides...)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 3 || !strings.Contains(r.stdout, ",worktree=incomplete:release}") {
		t.Fatalf("stack = (%d, %q, %q), want exit 3 at the release step", r.exit, r.stdout, r.stderr)
	}
	return file, gitOutput(t, f.root, "rev-parse", "main"), repair
}

// stackResumeSnapshot reads the refs, the worktrees, the intent ledger, and the destination
// checkout state at root.
func stackResumeSnapshot(t *testing.T, root string) string {
	t.Helper()
	ledger, err := os.ReadFile(filepath.Join(root, ".git", intent.Filename))
	mustNoError(t, err)
	return strings.Join([]string{
		gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)"),
		gitOutput(t, root, "worktree", "list", "--porcelain"),
		string(ledger),
		gitOutput(t, root, "status", "--porcelain=v1", "--untracked-files=all", "--ignored"),
	}, "\n--\n")
}

// TestLandStackIncompleteNamesTheResume is LS78. A stack called with one override whose
// release of member 2 fails exits 3, and its route names the stack resume of the published
// merge with that override.
func TestLandStackIncompleteNamesTheResume(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-incomplete", 2, nil)
	override := f.members[0].request + "=" + f.members[0].tip
	repair := plantMemberResidue(t, f, 1)
	defer repair()
	file, args := f.stack(t, override)
	r := runVerb(t, verbLand, f.call(args...))
	top := gitOutput(t, f.root, "rev-parse", "main")
	resume := "bench worktree land --resume '" + top + "' --stack '" + file + "' --source-tip '" + f.members[0].request + "'='" + f.members[0].tip + "'"
	next, printed := landingFaceNext(r.stdout, "landing stack stopped after publication")
	if want := landingRoute(faceLandIncomplete, resume, nil); r.exit != 3 || !printed || next != want {
		t.Fatalf("stack = (%d, %q, %q), want exit 3 and next %q", r.exit, r.stdout, r.stderr, want)
	}
}

// TestLandStackResumeProvesEachMember is LS79. A published commit whose second parent is
// not the tip of a listed member refuses before any step, at the top merge and at a walked
// merge below it. The walked merge holds a commit with member 1's tree on the review base in
// place of member 1's tip, so only the second-parent proof refuses it. A released member
// whose terminal receipt names another source tip also refuses before any step.
func TestLandStackResumeProvesEachMember(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, detail string
		// forge breaks one proof of the published stack top and answers the commit to
		// resume.
		forge func(t *testing.T, f stackFixture, top string) string
	}{
		{"top", errStackPublishedCount.Error(), func(t *testing.T, f stackFixture, top string) string {
			forged := gitOutput(t, f.root, "commit-tree", top+"^{tree}", "-p", top+"^1", "-p", f.base, "-m", "forged top")
			gitRun(t, f.root, "update-ref", "refs/heads/main", forged)
			return forged
		}},
		{"walked", "stack member 1: published commit does not authenticate the reviewed source parent", func(t *testing.T, f stackFixture, top string) string {
			stranger := gitOutput(t, f.root, "commit-tree", f.members[0].tip+"^{tree}", "-p", f.base, "-m", "not a member")
			lower := gitOutput(t, f.root, "commit-tree", top+"^1^{tree}", "-p", f.base, "-p", stranger, "-m", "forged member 1")
			forged := gitOutput(t, f.root, "commit-tree", top+"^{tree}", "-p", lower, "-p", f.members[1].tip, "-m", "forged member 2")
			gitRun(t, f.root, "update-ref", "refs/heads/main", forged)
			return forged
		}},
		{"receipt", "stack member 1: ", func(t *testing.T, f stackFixture, top string) string {
			m := f.members[0]
			repo, _, err := cleanupIdentity(nil, f.root, m.creation.Path)
			mustNoError(t, err)
			receipt, found, err := intent.CleanupReceiptFor(f.root, repo, releaseOperation, m.creation.Path, intent.RequestDigest(m.request))
			if err != nil || !found {
				t.Fatalf("terminal receipt of member 1: found=%t error=%v", found, err)
			}
			receipt.BranchOID = f.base
			mustNoError(t, intent.PutCleanupReceipt(f.root, receipt))
			return top
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := stackLandingFixture(t, "stack-resume-proof-"+row.name, 2, nil)
			file, top, repair := interruptedStack(t, f, 1)
			// Without the residue, a release before the proofs would write.
			repair()
			forged := row.forge(t, f, top)
			// The green marker goes back to the review base, which each destination reaches.
			// So the destination proofs hold, the member proofs decide, and a marker step
			// before them would write.
			gitRun(t, f.root, "update-ref", "refs/bench/green/main", f.base)
			before := stackResumeSnapshot(t, f.root)
			r := runVerb(t, verbLand, f.call(stackResumeArgs(forged, file)...))
			if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail="+row.detail) || strings.Contains(r.stdout, "stack{") {
				t.Fatalf("resume of a forged merge = (%d, %q, %q), want exit 1 and the refusal %q", r.exit, r.stdout, r.stderr, row.detail)
			}
			if after := stackResumeSnapshot(t, f.root); after != before {
				t.Fatalf("resume of a forged merge changed state:\n%s\nwant\n%s", after, before)
			}
		})
	}
}

// TestLandStackResumeFinishes is LS80 and LS81. After the repair, the stack resume releases
// the active published member and exits 0. A second resume exits 0 and changes no ref, no
// ledger record, and no checkout. The resume walks from the published commit when the
// destination is past it, and it leaves the destination where it is. An override that
// differs from the file tip of its member is the effective tip of that member.
func TestLandStackResumeFinishes(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		// overrides prepares the members of f and answers the overrides of the stack and of
		// its resume.
		overrides func(t *testing.T, f stackFixture) []string
		// before runs after the repair and before the resume of top over file.
		before func(t *testing.T, f stackFixture, file, top string)
	}{
		{"published", func(t *testing.T, f stackFixture) []string {
			return []string{f.members[1].request + "=" + f.members[1].tip}
		}, nil},
		{"ahead", nil, func(t *testing.T, f stackFixture, _, _ string) {
			gitRun(t, f.root, withIdentity("commit", "--allow-empty", "-qm", "after the stack")...)
		}},
		// Member 2 commits after the file tip, and the override names that commit. So the
		// published second parent is the override, and the resume without it refuses before
		// any write.
		{"override", func(t *testing.T, f stackFixture) []string {
			path := f.members[1].creation.Path
			commitInWorktree(t, path, stackMemberFile(2), "member bytes after the file\n", "member source after the file")
			return []string{f.members[1].request + "=" + gitOutput(t, path, "rev-parse", "HEAD")}
		}, func(t *testing.T, f stackFixture, file, top string) {
			before := stackResumeSnapshot(t, f.root)
			r := runVerb(t, verbLand, f.call(stackResumeArgs(top, file)...))
			if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail="+errStackPublishedCount.Error()) {
				t.Fatalf("resume at the file tip = (%d, %q, %q), want exit 1 and the refusal %q", r.exit, r.stdout, r.stderr, errStackPublishedCount)
			}
			if after := stackResumeSnapshot(t, f.root); after != before {
				t.Fatalf("resume at the file tip changed state:\n%s\nwant\n%s", after, before)
			}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := stackLandingFixture(t, "stack-resume-finish-"+row.name, 2, nil)
			var overrides []string
			if row.overrides != nil {
				overrides = row.overrides(t, f)
			}
			file, top, repair := interruptedStack(t, f, 1, overrides...)
			repair()
			if row.before != nil {
				row.before(t, f, file, top)
			}
			destination := gitOutput(t, f.root, "rev-parse", "main")
			args := stackResumeArgs(top, file, overrides...)
			r := runVerb(t, verbLand, f.call(args...))
			want := "stack{destination_base=" + f.base + ",published_commit=" + top + ",landed=2,worktree=released}\n"
			if r.exit != 0 || !strings.HasSuffix(r.stdout, want) {
				t.Fatalf("resume = (%d, %q, %q), want exit 0 and the last line %q", r.exit, r.stdout, r.stderr, want)
			}
			if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
				t.Fatalf("resume moved the destination to %s, want %s", got, destination)
			}
			for _, m := range f.members {
				if assignmentActive(t, f.root, m.creation.Assignment.ID) {
					t.Fatalf("member %s is active after the resume", m.request)
				}
			}
			before := stackResumeSnapshot(t, f.root)
			again := runVerb(t, verbLand, f.call(args...))
			if again.exit != 0 || !strings.HasSuffix(again.stdout, want) {
				t.Fatalf("second resume = (%d, %q, %q), want exit 0 and the last line %q", again.exit, again.stdout, again.stderr, want)
			}
			if after := stackResumeSnapshot(t, f.root); after != before {
				t.Fatalf("second resume changed state:\n%s\nwant\n%s", after, before)
			}
		})
	}
}

// TestLandStackResumeKeepsUnlandedMembers is LS82. A stack of four with member 3 red
// publishes members 1 and 2, and the release of member 2 stops, so the stack exits 3. The
// resume releases member 2 and leaves members 3 and 4 active.
func TestLandStackResumeKeepsUnlandedMembers(t *testing.T) {
	t.Parallel()
	f := redStackFixture(t, "stack-resume-keep", 4, 2)
	file, top, repair := interruptedStack(t, f, 1)
	if second := gitOutput(t, f.root, "rev-parse", top+"^2"); second != f.members[1].tip {
		t.Fatalf("published second parent = %s, want the tip of member 2 %s", second, f.members[1].tip)
	}
	repair()
	r := runVerb(t, verbLand, f.call(stackResumeArgs(top, file)...))
	if r.exit != 0 || !strings.HasSuffix(r.stdout, ",landed=2,worktree=released}\n") {
		t.Fatalf("resume = (%d, %q, %q), want exit 0 with two landed members", r.exit, r.stdout, r.stderr)
	}
	for i, m := range f.members {
		if active := assignmentActive(t, f.root, m.creation.Assignment.ID); active != (i >= 2) {
			t.Fatalf("member %d active = %v after the resume, want active only for members 3 and 4", i+1, active)
		}
	}
}
