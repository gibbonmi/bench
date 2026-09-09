package probe

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/usage"
)

// PB12: the copy exists on disk, at its declared mode, while the run is in flight.
func TestProbePreservesTheSubjectBeforeTheRun(t *testing.T) {
	f := newFixture(t)
	listing := filepath.Join(scratchDir(t), "preserved")
	body := filepath.Join(filepath.Dir(listing), "preserved-body")
	// The copy is read while the run holds it, because a proven restore removes it before
	// the verb answers and the test would then read an absent file.
	hook := "find \"$BENCH_HOME/probe\" -type f -printf '%m %p\\n' > " + sanitize.ShellQuote(listing) + "\n" +
		"find \"$BENCH_HOME/probe\" -type f -exec cat {} + > " + sanitize.ShellQuote(body) + "\n"
	installStubGo(t, f, cannedFailure("caught"), 1, hook)
	if _, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	observed, err := os.ReadFile(listing)
	if err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(string(observed))
	prefix := "600 " + filepath.Join(f.home, "probe") + string(filepath.Separator)
	if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, "clamp.go") {
		t.Fatalf("preserved listing = %q, want one 0600 clamp.go below %s", line, prefix)
	}
	copied, err := os.ReadFile(body)
	if err != nil || string(copied) != clampSource {
		t.Fatalf("preserved copy = (%q, %v), want the start bytes", copied, err)
	}
}

// PB14: a restore that cannot write keeps the copy, names it, and exits 2, so a caller
// whose tree still holds the mutation has the file it needs to put back by hand.
func TestProbeReportsARestoreFailure(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "chmod 0500 .\n")
	t.Cleanup(func() { _ = os.Chmod(f.root, 0o755) })
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 2, probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no"))
	kept := preservedCopy(t, f)
	if !strings.Contains(out, "preserved[1]{path,reason}:\n  "+kept) {
		t.Fatalf("stdout = %q, want the preserved row naming %s", out, kept)
	}
	saved, err := os.ReadFile(kept)
	if err != nil || string(saved) != clampSource {
		t.Fatalf("preserved copy = (%q, %v), want the start bytes", saved, err)
	}
}

// PB43: the read-back is compared with the start bytes, not with the copy, so a copy that
// changed under the run cannot certify its own damage as a restore.
func TestProbeReportsAReadBackMismatch(t *testing.T) {
	f := newFixture(t)
	hook := "find \"$BENCH_HOME/probe\" -type f -exec truncate -s 0 {} +\n"
	installStubGo(t, f, cannedFailure("caught"), 1, hook)
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 2, probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no"))
	want := "the restored bytes differ from the bytes read at the start"
	if !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want the preserved row's reason %q", out, want)
	}
}

// PB41: the restore runs before the render, so a subject the encoder cannot carry still
// leaves a clean tree behind its refusal.
func TestProbeRestoresBeforeARenderRefusal(t *testing.T) {
	f := newFixture(t)
	name := "data\aalert.txt"
	subject := filepath.Join(f.root, name)
	writeFixtureFile(t, subject, "alpha\n", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, name, "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	requireRow(t, out, code, 1, "error: unrepresentable TOON cell — ")
	restored, err := os.ReadFile(subject)
	if err != nil || string(restored) != "alpha\n" {
		t.Fatalf("subject = (%q, %v), want the start bytes", restored, err)
	}
	requireHomeEmpty(t, f)
}

// PB47: a subject the encoder cannot carry and a restore that cannot write answer at once.
// The caller reads the render refusal and still learns where the copy is, because the
// restore-failed verdict overrides every other answer.
func TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail(t *testing.T) {
	f := newFixture(t)
	// The BEL byte sits in the directory component, so the row's subject cell refuses while
	// the copy's own path, which carries only the base name, still renders.
	held := filepath.Join(f.root, "data\aalert")
	operand := filepath.Join("data\aalert", "alert.txt")
	writeFixtureFile(t, filepath.Join(held, "alert.txt"), "alpha\n", 0o644)
	installStubGo(t, f, cannedPass, 0, "chmod 0500 "+sanitize.ShellQuote(held)+"\n")
	t.Cleanup(func() { _ = os.Chmod(held, 0o755) })
	out, code := runProbe(t, operand, "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	requireRow(t, out, code, 2, "error: unrepresentable TOON cell — ")
	kept := preservedCopy(t, f)
	if !strings.Contains(out, "preserved[1]{path,reason}:\n  "+kept) {
		t.Fatalf("stdout = %q, want the preserved row naming %s", out, kept)
	}
	saved, err := os.ReadFile(kept)
	if err != nil || string(saved) != "alpha\n" {
		t.Fatalf("preserved copy = (%q, %v), want the start bytes", saved, err)
	}
}

// PB17: an executable subject stays executable across the mutation and the restore.
func TestProbeKeepsTheSubjectMode(t *testing.T) {
	f := newFixture(t)
	script := filepath.Join(f.root, "run.sh")
	writeFixtureFile(t, script, "#!/bin/sh\necho alpha\n", 0o755)
	installStubGo(t, f, cannedPass, 0, "")
	if _, code := runProbe(t, "run.sh", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$"); code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	info, err := os.Lstat(script)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("subject mode = %v, want 0755", info.Mode().Perm())
	}
}

// PB39: a subject whose name carries a space renders through the encoder, and the
// expectation derives through the same call.
func TestProbeRendersASubjectWithASpace(t *testing.T) {
	f := newFixture(t)
	writeFixtureFile(t, filepath.Join(f.root, "my pkg", "clamp.go"), "package spaced\n\nconst alpha = 1\n", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, "my pkg/clamp.go", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	requireRow(t, out, code, 1, probeRow(t, "silent", "my pkg/clamp.go", "swap", "passed", 0, "yes"))
}

// PB40: a subject whose last line has no newline restores byte-exact.
func TestProbeRestoresAFileWithoutTrailingNewline(t *testing.T) {
	f := newFixture(t)
	subject := filepath.Join(f.root, "notrailing.txt")
	writeFixtureFile(t, subject, "alpha", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, "notrailing.txt", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	requireRow(t, out, code, 1, probeRow(t, "silent", "notrailing.txt", "swap", "passed", 0, "yes"))
	got, err := os.ReadFile(subject)
	if err != nil || string(got) != "alpha" {
		t.Fatalf("restored bytes = (%q, %v), want %q", got, err, "alpha")
	}
}

// PB46: the verb writes only the preserved copy and the subject, so a completed probe
// adds no entry to the administration directory, the Bench home, or the tree.
func TestProbeWritesNoRecord(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	before := listRecordSurfaces(t, f)
	if _, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if after := listRecordSurfaces(t, f); before != after {
		t.Fatalf("record surfaces changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// listRecordSurfaces answers a stable listing of the three places a record could land.
func listRecordSurfaces(t *testing.T, f *fixture) string {
	t.Helper()
	var lines []string
	for _, dir := range []string{f.home, filepath.Join(f.root, ".git"), f.root} {
		walked := dir
		if err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if walked == f.root && entry.IsDir() && entry.Name() == ".git" {
				return filepath.SkipDir
			}
			lines = append(lines, path)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	return strings.Join(lines, "\n")
}

// PB7: the check form builds the request `bench test` builds for the same check name, so
// the two verbs cannot disagree about what a named check selects.
func TestProbeSelectsTheCheckForm(t *testing.T) {
	f := newFixture(t)
	parsed, line, code := usage.Parse(grammar, []string{"clamp.go", "--omit", "n < 0", "--check", "line-routing"})
	if line != "" {
		t.Fatalf("parse = (%q, %d), want an accepted selection", line, code)
	}
	if got := strings.Join(selectionArgs(parsed), " "); got != "--check line-routing" {
		t.Fatalf("selection args = %q, want %q", got, "--check line-routing")
	}
	got, line, _ := testreport.Prepare(f.root, selectionArgs(parsed))
	if line != "" {
		t.Fatalf("prepare refused the probe's check form: %q", line)
	}
	want, line, _ := testreport.Prepare(f.root, []string{"--check", "line-routing"})
	if line != "" {
		t.Fatalf("prepare refused the reference check form: %q", line)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("probe request = %#v, want %#v", got, want)
	}
}

// DG1: the check form names the check and the root conformance run pattern, so a caller
// reads which tests the probe selected rather than inferring them from the check name.
func TestProbeNamesTheCheckSelection(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, "clamp.go", "--omit", "return 0", "--check", "line-routing")
	want := selectionRow(t, "check", "line-routing", "^TestRootConformance$", "passed", 1)
	if code != 0 || !strings.Contains(out, want) {
		t.Fatalf("stdout = (%q, %d), want the check selection row %q and 0", out, code, want)
	}
}

// DG2: the package form names its expression and the exact run pattern it passed, and a
// selection that passed none names `all` rather than an empty cell.
func TestProbeNamesThePackageSelection(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, probeArgs()...)
	want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes") + selectionRow(t, "package", "./", "^TestClampNegative$", "passed", 1)
	requireRow(t, out, code, 0, want)
	bare, _ := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./")
	if row := selectionRow(t, "package", "./", "all", "passed", 1); !strings.Contains(bare, row) {
		t.Fatalf("stdout = %q, want the pattern-free selection row %q", bare, row)
	}
}

// PB38: the verdict mapping is the one place an outcome becomes a word and an exit.
func TestVerdictExitCodes(t *testing.T) {
	for _, tc := range []struct {
		kind    testreport.OutcomeKind
		verdict string
		code    int
	}{
		{testreport.OutcomeFailed, "bit", 0},
		{testreport.OutcomePassed, "silent", 1},
		{testreport.OutcomeBuildFailed, "invalid", 1},
		{testreport.OutcomeNoTestRun, "invalid", 1},
		{testreport.OutcomeRefused, "invalid", 1},
		{testreport.OutcomeInterrupted, "invalid", 1},
	} {
		t.Run(string(tc.kind), func(t *testing.T) {
			verdict, code := verdictFor(tc.kind)
			if verdict != tc.verdict || code != tc.code {
				t.Fatalf("verdictFor(%q) = (%q, %d), want (%q, %d)", tc.kind, verdict, code, tc.verdict, tc.code)
			}
		})
		t.Run(string(tc.kind)+" over a failed restore", func(t *testing.T) {
			preserved := preservation{file: filepath.Join(t.TempDir(), "clamp.go")}
			out, code := render(subject{display: "clamp.go"}, "swap", testreport.Outcome{Kind: tc.kind}, testreport.Request{}, "", preserved, false, "denied")
			if code != 2 {
				t.Fatalf("render(%q) over a failed restore exit = %d, want 2\n%s", tc.kind, code, out)
			}
			want := probeRow(t, "restore-failed", "clamp.go", "swap", string(tc.kind), 0, "no")
			if !strings.HasPrefix(out, want) {
				t.Fatalf("stdout = %q, want %q first", out, want)
			}
		})
	}
}
