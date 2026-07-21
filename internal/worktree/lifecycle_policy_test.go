package worktree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestReleaseDeadLeaseRemovesAndCompacts(t *testing.T) {
	root, creation := newOwnedAssignment(t, "dead-lease-release")
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(deadPidLine(t)), 0o600)

	var stdout bytes.Buffer
	code := ReleaseCommand(root, []string{"--request", "landed-dead-lease-release", creation.Path}, &stdout, io.Discard)
	requireTest(t, code == 0, "dead-lease release exit=%d stdout=%q", code, stdout.String())
	_, statErr := os.Stat(creation.Path)
	requireTest(t, os.IsNotExist(statErr), "dead-lease worktree remains: %v", statErr)
	_, err = assignmentByID(root, creation.Assignment.ID)
	requireTest(t, err != nil, "dead-lease assignment was not compacted")
	ledger, err := intent.Read(root)
	requireTest(t, err == nil && len(ledger.CleanupReceipts) == 2, "dead-lease receipts=%#v error=%v", ledger.CleanupReceipts, err)
}

func TestReleaseDeclaredBuildOutputRemoves(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	repositoryRoot := gitOutput(t, ".", "rev-parse", "--show-toplevel")
	declaration, err := os.ReadFile(filepath.Join(repositoryRoot, ".bench", "build-outputs.json"))
	mustNoError(t, err)
	mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), declaration, 0o644)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "landed-declared-output", "declared output")
	mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
	mustWrite(t, filepath.Join(creation.Path, "dist", "bench"), []byte("binary\n"), 0o755)

	var stdout bytes.Buffer
	code := ReleaseCommand(root, []string{"--request", "landed-declared-output", creation.Path}, &stdout, io.Discard)
	requireTest(t, code == 0, "declared-output release exit=%d stdout=%q", code, stdout.String())
	_, statErr := os.Stat(creation.Path)
	requireTest(t, os.IsNotExist(statErr), "declared-output worktree remains: %v", statErr)
}

func TestReleaseBuildOutputContainmentRetainsUnknownResidue(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(*testing.T, string)
	}{
		{"undeclared", func(t *testing.T, target string) {
			mustWrite(t, filepath.Join(target, "secret.txt"), []byte("secret\n"), 0o600)
		}},
		{"over-limit", func(t *testing.T, target string) {
			mustMkdirAll(t, filepath.Join(target, "dist"), 0o755)
			for i := 0; i <= ignoredEntryLimit; i++ {
				mustWrite(t, filepath.Join(target, "dist", fmt.Sprintf("output-%04d", i)), []byte("x"), 0o600)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			gitRun(t, root, "branch", "-M", "main")
			mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\nsecret.txt\n"), 0o644)
			gitRun(t, root, "add", ".gitignore")
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore residue")
			mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
			mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\"dist/\"]}\n"), 0o644)
			t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
			request := "landed-containment-" + tc.name
			creation := mustCreate(t, root, request, tc.name)
			tc.build(t, creation.Path)

			var stderr bytes.Buffer
			code := ReleaseCommand(root, []string{"--request", request, creation.Path}, io.Discard, &stderr)
			requireTest(t, code == 1 && strings.Contains(stderr.String(), "worktree retained (ignored)"),
				"%s release exit=%d stderr=%q", tc.name, code, stderr.String())
			_, statErr := os.Stat(creation.Path)
			requireTest(t, statErr == nil, "%s worktree removed: %v", tc.name, statErr)
		})
	}
}

func TestBuildOutputDeclarationFailsClosed(t *testing.T) {
	absent := ""
	empty := "{\"schema\":1,\"paths\":[]}\n"
	noNewline := "{\"schema\":1,\"paths\":[\"build output/\"]}"
	malformed := "{"
	badSchema := "{\"schema\":2,\"paths\":[\"dist/\"]}\n"
	traversal := "{\"schema\":1,\"paths\":[\"../dist/\"]}\n"
	absolute := "{\"schema\":1,\"paths\":[\"/dist/\"]}\n"
	control := "{\"schema\":1,\"paths\":[\"dist/\\u0007\"]}\n"
	unknown := "{\"schema\":1,\"paths\":[\"dist/\"],\"extra\":true}\n"
	duplicate := "{\"schema\":1,\"paths\":[\"dist/\"],\"paths\":[\"dist/\"]}\n"
	oversized := "{\"schema\":1,\"paths\":[\"dist/" + strings.Repeat("x", 17*1024) + "\"]}\n"
	for _, tc := range []struct {
		name, declaration, ignored, residual, wantReason string
		wantCode                                         int
	}{
		{"absent", absent, "dist/", "dist/output", "ignored", 1},
		{"empty", empty, "dist/", "dist/output", "ignored", 1},
		{"valid-no-final-newline", noNewline, "build output/", "build output/result", "", 0},
		{"malformed-json", malformed, "dist/", "dist/output", "malformed", 1},
		{"bad-schema", badSchema, "dist/", "dist/output", "malformed", 1},
		{"traversal", traversal, "dist/", "dist/output", "malformed", 1},
		{"absolute", absolute, "dist/", "dist/output", "malformed", 1},
		{"control-byte", control, "dist/", "dist/output", "malformed", 1},
		{"unknown-field", unknown, "dist/", "dist/output", "malformed", 1},
		{"duplicate-field", duplicate, "dist/", "dist/output", "malformed", 1},
		{"oversized", oversized, "dist/", "dist/output", "malformed", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			gitRun(t, root, "branch", "-M", "main")
			mustWrite(t, filepath.Join(root, ".gitignore"), []byte(tc.ignored+"\n"), 0o644)
			gitRun(t, root, "add", ".gitignore")
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
			if tc.declaration != "" {
				mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
				mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte(tc.declaration), 0o644)
			}
			t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
			request := "landed-build-output-" + tc.name
			creation := mustCreate(t, root, request, tc.name)
			residual := filepath.Join(creation.Path, filepath.FromSlash(tc.residual))
			mustMkdirAll(t, filepath.Dir(residual), 0o755)
			mustWrite(t, residual, []byte("output\n"), 0o600)

			var stderr bytes.Buffer
			code := ReleaseCommand(root, []string{"--request", request, creation.Path}, io.Discard, &stderr)
			requireTest(t, code == tc.wantCode, "%s release exit=%d stderr=%q", tc.name, code, stderr.String())
			if tc.wantReason != "" {
				requireTest(t, strings.Contains(stderr.String(), "worktree retained ("+tc.wantReason+")"),
					"%s reason missing: %q", tc.name, stderr.String())
			}
			_, statErr := os.Stat(creation.Path)
			requireTest(t, (statErr == nil) == (tc.wantCode != 0), "%s tree existence error=%v", tc.name, statErr)
		})
	}
}

func TestResumeReconcilesDeadLeaseAndPreservesSafetyBranches(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))

	dead := mustCreate(t, root, "landed-resume-dead", "dead owner")
	markPending(t, root, dead.Assignment)
	deadLease, err := LeaseFile(dead.Path)
	mustNoError(t, err)
	mustWrite(t, deadLease, []byte(deadPidLine(t)), 0o600)

	live := mustCreate(t, root, "landed-resume-live", "live owner")
	markPending(t, root, live.Assignment)
	liveLease, err := LeaseFile(live.Path)
	mustNoError(t, err)
	mustWrite(t, liveLease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)

	unlanded := mustCreate(t, root, "landed-resume-unlanded", "preserved work")
	commitInWorktree(t, unlanded.Path, "unique.txt", "preserve\n", "unique work")
	markPending(t, root, unlanded.Assignment)

	chdir(t, root)
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	requireTest(t, code == 0 && stderr.String() == "", "resume exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	want := "bench resume: removed 1, recovered 0; retained live-lease=1 unmerged=1; pruned branches 0; reconciled 0; failed 0; open assignments 2\n"
	requireTest(t, stdout.String() == want, "resume summary=%q want=%q", stdout.String(), want)
	_, deadErr := os.Stat(dead.Path)
	_, liveErr := os.Stat(live.Path)
	_, unlandedErr := os.Stat(unlanded.Path)
	requireTest(t, os.IsNotExist(deadErr) && liveErr == nil && unlandedErr == nil,
		"resume tree states dead=%v live=%v unlanded=%v", deadErr, liveErr, unlandedErr)
}

func TestRepositoryDeclaresDistBuildOutput(t *testing.T) {
	root := gitOutput(t, ".", "rev-parse", "--show-toplevel")
	paths, _, err := loadBuildOutputs(root)
	mustNoError(t, err)
	found := false
	for _, path := range paths {
		found = found || path == "dist/"
	}
	requireTest(t, found, "repository build outputs=%#v, want dist/", paths)
}

func TestReleaseLiveLeaseRetains(t *testing.T) {
	root, creation := newOwnedAssignment(t, "live-lease-release")
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)

	var stderr bytes.Buffer
	code := ReleaseCommand(root, []string{"--request", "landed-live-lease-release", creation.Path}, io.Discard, &stderr)
	requireTest(t, code == 1, "live-lease release exit=%d stderr=%q", code, stderr.String())
	requireTest(t, strings.Contains(stderr.String(), "worktree retained (live-lease)"), "live-lease reason missing: %q", stderr.String())
	_, statErr := os.Stat(creation.Path)
	requireTest(t, statErr == nil, "live-lease worktree removed: %v", statErr)
}

func TestReleaseMalformedLeaseRetainsAsUncertain(t *testing.T) {
	assertReleaseLeaseRetainedAsUncertain(t, "malformed-lease-release", func(lease string) {
		mustWrite(t, lease, []byte("garbage\n"), 0o600)
	})
}

func TestReleaseNumericLeaseWithBadTimestampAndExtraFieldsRetainsAsUncertain(t *testing.T) {
	assertReleaseLeaseRetainedAsUncertain(t, "numeric-bad-timestamp-extra-lease-release", func(lease string) {
		mustWrite(t, lease, []byte(fmt.Sprintf("%d not-a-time trailing\n", os.Getpid())), 0o600)
	})
}

func TestReleasePartialNumericLeaseRetainsAsUncertain(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
	}{
		{"pid-only", fmt.Sprintf("%d", os.Getpid())},
		{"missing-final-newline", fmt.Sprintf("%d 2026-07-15T00:00:00Z", os.Getpid())},
		{"truncated-timestamp", fmt.Sprintf("%d 2026-07-15T00:00:\n", os.Getpid())},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertReleaseLeaseRetainedAsUncertain(t, "partial-"+tc.name+"-lease-release", func(lease string) {
				mustWrite(t, lease, []byte(tc.content), 0o600)
			})
		})
	}
}

func TestReleaseSymlinkLeaseRetainsAsUncertain(t *testing.T) {
	assertReleaseLeaseRetainedAsUncertain(t, "symlink-lease-release", func(lease string) {
		target := filepath.Join(t.TempDir(), "lease-target")
		mustWrite(t, target, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)
		mustNoError(t, os.Symlink(target, lease))
	})
}

func TestReleaseDirectoryLeaseRetainsAsUncertain(t *testing.T) {
	assertReleaseLeaseRetainedAsUncertain(t, "directory-lease-release", func(lease string) {
		mustMkdirAll(t, lease, 0o700)
	})
}

func assertReleaseLeaseRetainedAsUncertain(t *testing.T, request string, makeLease func(string)) {
	t.Helper()
	root, creation := newOwnedAssignment(t, request)
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	makeLease(lease)

	var stderr bytes.Buffer
	code := ReleaseCommand(root, []string{"--request", "landed-" + request, creation.Path}, io.Discard, &stderr)
	requireTest(t, code == 1, "malformed-lease release exit=%d stderr=%q", code, stderr.String())
	requireTest(t, strings.Contains(stderr.String(), "worktree retained (uncertain)") && strings.Contains(stderr.String(), "lease state is unknown"),
		"malformed-lease reason missing: %q", stderr.String())
	_, statErr := os.Stat(creation.Path)
	requireTest(t, statErr == nil, "malformed-lease worktree removed: %v", statErr)
}
