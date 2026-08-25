package spec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/usage"
)

// writeSpec writes content to <dir>/specs/<slug>/spec.md and returns the path.
func writeSpec(t *testing.T, dir, slug, content string) string {
	t.Helper()
	specsDir := filepath.Join(dir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(specsDir, slug, "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeFolderSpec(t *testing.T, dir, slug, content string) string {
	t.Helper()
	path := filepath.Join(dir, "specs", slug, "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestImplementedPreservesEveryOtherByte pins the surviving flip source: Implemented
// rewrites the one line-start `Status: staged` and nothing else, whatever separator,
// trailing whitespace, or final newline the file carries.
func TestImplementedPreservesEveryOtherByte(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"space separator, trailing newline", "# spec\n\nStatus: staged\n\n## body\n", "# spec\n\nStatus: implemented\n\n## body\n"},
		{"tab separator preserved", "Status:\tstaged\nbody\n", "Status:\timplemented\nbody\n"},
		{"no trailing newline preserved", "# spec\nStatus: staged", "# spec\nStatus: implemented"},
		{"trailing whitespace on status line preserved", "Status:  staged  \n", "Status:  implemented  \n"},
		{"other staged word untouched", "Status: staged\nThe work is staged elsewhere.\n", "Status: implemented\nThe work is staged elsewhere.\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Implemented([]byte(tc.in))
			if err != nil {
				t.Fatalf("Implemented: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("content = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestImplementedDerivesExactBytesAndRefusesMalformedStatus(t *testing.T) {
	got, err := Implemented([]byte("x\nStatus: staged  \nlast"))
	if err != nil || string(got) != "x\nStatus: implemented  \nlast" {
		t.Fatalf("Implemented = %q, %v", got, err)
	}
	for _, content := range []string{"Status: implemented\n", "Status: staged\nStatus: staged\n", "status: staged\n"} {
		if _, err := Implemented([]byte(content)); err == nil {
			t.Fatalf("Implemented(%q) succeeded", content)
		}
	}
}

// TestSpecImplementedIsAnUnknownSubcommand pins FA3: the retired subcommand exits 2
// through the unknown-subcommand branch, and the named spec keeps every byte.
func TestSpecImplementedIsAnUnknownSubcommand(t *testing.T) {
	dir := t.TempDir()
	path := writeSpec(t, dir, "x", "# spec\n\nStatus: staged\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out, code := Command([]string{"implemented", "x"})
	if code != 2 {
		t.Fatalf("code = %d, want 2; out = %q", code, out)
	}
	if !strings.Contains(out, "usage: bench spec (unknown argument: implemented)") {
		t.Errorf("out = %q, want the unknown-subcommand usage line", out)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Errorf("spec content changed: %q -> %q", before, after)
	}
}

// TestSpecUsageNamesTheSurvivingSubcommands pins FA3's bare-argv half: the usage line
// offers retire and history only.
func TestSpecUsageNamesTheSurvivingSubcommands(t *testing.T) {
	out, code := Command(nil)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(out, "expected a subcommand: retire, history") {
		t.Errorf("out = %q, want the retire/history usage line", out)
	}
}

func TestSpecHelpPrintsThePublicSubcommandInventory(t *testing.T) {
	const want = "usage: bench spec <subcommand>\n\nsubcommands:\n  retire <slug>\n  history <slug>\n"
	for _, arg := range []string{"--help", "-h"} {
		out, code := Command([]string{arg})
		if code != 0 {
			t.Fatalf("Command(%q) code = %d, want 0; out = %q", arg, code, out)
		}
		if out != want {
			t.Errorf("Command(%q) out = %q, want %q", arg, out, want)
		}
	}
}

func TestResolveConvention(t *testing.T) {
	dir := t.TempDir()
	slugPath := writeSpec(t, dir, "mine", "Status: staged\n")

	// a bare slug resolves to <base>/specs/<slug>/spec.md
	_, resolved, _, ok, err := Resolve(dir, "mine")
	if err != nil || !ok {
		t.Fatalf("bare slug: ok=%v err=%v", ok, err)
	}
	if resolved != slugPath {
		t.Errorf("bare slug resolved = %q, want %q", resolved, slugPath)
	}

	// a path argument resolves as-given (cwd-relative), not via the fallback
	abs := slugPath
	_, resolved, _, ok, err = Resolve(dir, abs)
	if err != nil || !ok || resolved != abs {
		t.Errorf("path arg: resolved=%q ok=%v err=%v, want %q", resolved, ok, err, abs)
	}

	// a slug containing '/' gets no fallback
	_, _, tried, ok, _ := Resolve(dir, "sub/mine")
	if ok {
		t.Error("slug with '/': expected no resolution")
	}
	if len(tried) != 1 || tried[0] != "sub/mine" {
		t.Errorf("slug with '/': tried = %v, want only the as-given form", tried)
	}
}

func TestResolveBaseAnchorsFallbackFromAnyCwd(t *testing.T) {
	root := t.TempDir()
	writeSpec(t, root, "anchored", "Status: staged\n")
	// Resolve from a different cwd: the fallback must still hit <root>/specs.
	_, resolved, _, ok, err := Resolve(root, "anchored")
	if err != nil || !ok {
		t.Fatalf("anchored resolve: ok=%v err=%v", ok, err)
	}
	if resolved != filepath.Join(root, "specs", "anchored", "spec.md") {
		t.Errorf("resolved = %q, want root-anchored path", resolved)
	}
}

func TestResolveFolderFormFromDeeperCWD(t *testing.T) {
	root := t.TempDir()
	want := writeFolderSpec(t, root, "folder", "Status: staged\n")
	deep := filepath.Join(root, "deep", "cwd")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(deep)

	_, resolved, _, ok, err := Resolve(root, "folder")
	if err != nil || !ok {
		t.Fatalf("folder slug: ok=%v err=%v", ok, err)
	}
	if resolved != want {
		t.Errorf("resolved = %q, want %q", resolved, want)
	}
}

func TestFactsIncludesFolderSpecsAndMalformedEvidence(t *testing.T) {
	root := t.TempDir()
	writeSpec(t, root, "flat", "Status: staged\n")
	writeFolderSpec(t, root, "good", "Status: staged\nRoadmap: FT1\n")
	writeFolderSpec(t, root, "bad", "# no metadata\n")

	facts, err := Facts(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []Fact{
		{Slug: "bad", Path: "specs/bad/spec.md"},
		{Slug: "flat", Path: "specs/flat/spec.md", Status: "staged"},
		{Slug: "good", Path: "specs/good/spec.md", Status: "staged", RoadmapID: "FT1"},
	}
	if !reflect.DeepEqual(facts, want) {
		t.Fatalf("Facts = %#v, want %#v", facts, want)
	}
}

func TestFactsEnumeratesBracketedRootLiterally(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo[brackets]")
	writeFolderSpec(t, root, "folder", "Status: staged\n")

	facts, err := Facts(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []Fact{{Slug: "folder", Path: "specs/folder/spec.md", Status: "staged"}}
	if !reflect.DeepEqual(facts, want) {
		t.Fatalf("Facts = %#v, want %#v", facts, want)
	}
}

func TestFactsEnumeratesFolderCandidates(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, root string)
		want  []Fact
	}{
		{
			name: "skips non-directory entries and folders without spec files",
			setup: func(t *testing.T, root string) {
				writeFolderSpec(t, root, "zulu", "Status: staged\n")
				writeFolderSpec(t, root, "alpha", "Status: staged\n")
				if err := os.MkdirAll(filepath.Join(root, "specs", "missing"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "specs", "note.md"), []byte("not a folder spec\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: []Fact{
				{Slug: "alpha", Path: "specs/alpha/spec.md", Status: "staged"},
				{Slug: "zulu", Path: "specs/zulu/spec.md", Status: "staged"},
			},
		},
		{
			name: "follows a directory symlink",
			setup: func(t *testing.T, root string) {
				target := filepath.Join(root, "target")
				if err := os.MkdirAll(target, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(target, "spec.md"), []byte("Status: staged\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(filepath.Join(root, "specs"), 0o755); err != nil {
					t.Fatal(err)
				}
				link := filepath.Join(root, "specs", "linked")
				if err := os.Symlink(target, link); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
				}
			},
			want: []Fact{{Slug: "linked", Path: "specs/linked/spec.md", Status: "staged"}},
		},
		{
			name: "keeps a .md-suffixed directory name in the derived path",
			setup: func(t *testing.T, root string) {
				writeFolderSpec(t, root, "x.md", "Status: staged\n")
			},
			want: []Fact{{Slug: "x.md", Path: "specs/x.md/spec.md", Status: "staged"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.setup(t, root)

			facts, err := Facts(root)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(facts, tc.want) {
				t.Fatalf("Facts = %#v, want %#v", facts, tc.want)
			}
		})
	}
}

func TestLiveSpecPrimitives(t *testing.T) {
	if got, want := LiveSpecPath("ticket-19"), "specs/ticket-19/spec.md"; got != want {
		t.Errorf("LiveSpecPath = %q, want %q", got, want)
	}
	content := []byte("specs/one/spec.md specs/two_2/spec.md specs/one/spec.md\n```\nspecs/hidden/spec.md\n```\nspecs/<slug>/spec.md\n")
	if got, want := LiveSpecSlugs(content), []string{"one", "two_2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("LiveSpecSlugs = %q, want %q", got, want)
	}
}

func TestLiveSpecPathNormalizesExplicitPath(t *testing.T) {
	if got, want := LiveSpecPath("./specs/ticket-19/spec.md"), "specs/ticket-19/spec.md"; got != want {
		t.Errorf("LiveSpecPath = %q, want %q", got, want)
	}
	if got, want := LiveSpecSlug("./specs/ticket-19/spec.md"), "ticket-19"; got != want {
		t.Errorf("LiveSpecSlug = %q, want %q", got, want)
	}
}

func TestFactsRetainsHostileCandidates(t *testing.T) {
	root := t.TempDir()
	writeFolderSpec(t, root, "malformed", "Status: sta\xffged\n")
	writeFolderSpec(t, root, "oversized", string(bytes.Repeat([]byte("x"), int(bounds.ControlRecordLimit+1))))

	dangling := filepath.Join(root, "specs", "dangling", "spec.md")
	if err := os.MkdirAll(filepath.Dir(dangling), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing.md", dangling); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
	}

	fifo := filepath.Join(root, "specs", "fifo", "spec.md")
	if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}

	done := make(chan struct {
		facts []Fact
		err   error
	}, 1)
	go func() {
		facts, err := Facts(root)
		done <- struct {
			facts []Fact
			err   error
		}{facts, err}
	}()
	select {
	case got := <-done:
		if got.err != nil {
			t.Fatalf("Facts: %v", got.err)
		}
		want := []Fact{
			{Slug: "dangling", Path: "specs/dangling/spec.md"},
			{Slug: "fifo", Path: "specs/fifo/spec.md"},
			{Slug: "malformed", Path: "specs/malformed/spec.md"},
			{Slug: "oversized", Path: "specs/oversized/spec.md"},
		}
		if !reflect.DeepEqual(got.facts, want) {
			t.Fatalf("Facts = %#v, want %#v", got.facts, want)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("Facts blocked on a FIFO candidate, so it read before classifying the path")
	}
}

func TestFactsRetainsUnreadableCandidate(t *testing.T) {
	root := t.TempDir()
	path := writeFolderSpec(t, root, "unreadable", "Status: staged\n")
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip permissions: %v", err))
	}
	file, err := os.Open(path)
	if err == nil {
		_ = file.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}

	facts, err := Facts(root)
	if err != nil {
		t.Fatalf("Facts: %v", err)
	}
	want := []Fact{{Slug: "unreadable", Path: "specs/unreadable/spec.md"}}
	if !reflect.DeepEqual(facts, want) {
		t.Fatalf("Facts = %#v, want %#v", facts, want)
	}
}

func TestResolveRefusesLiveFlatForms(t *testing.T) {
	root := t.TempDir()
	flat := filepath.Join(root, "specs", "legacy.md")
	if err := os.MkdirAll(filepath.Dir(flat), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(flat, []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, _, ok, err := Resolve(root, "legacy")
	if ok || err == nil || !strings.Contains(err.Error(), flatLayoutInstruction) || !strings.Contains(err.Error(), flat) || !strings.Contains(err.Error(), filepath.Join(root, "specs", "legacy", "spec.md")) {
		t.Fatalf("flat-only resolve = ok %v, err %v", ok, err)
	}
	if err := os.Mkdir(filepath.Join(root, "specs", "legacy"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "legacy", "spec.md"), []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, _, ok, err = Resolve(root, "legacy")
	if ok || err == nil || !strings.Contains(err.Error(), flatLayoutInstruction) || !strings.Contains(err.Error(), flat) || !strings.Contains(err.Error(), filepath.Join(root, "specs", "legacy", "spec.md")) {
		t.Fatalf("collision resolve = ok %v, err %v", ok, err)
	}
}

func TestResolveFlatPrecedesIncompleteFolder(t *testing.T) {
	root := t.TempDir()
	flat := filepath.Join(root, "specs", "partial.md")
	folder := filepath.Join(root, "specs", "partial", "spec.md")
	if err := os.MkdirAll(filepath.Dir(folder), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(flat, []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, _, ok, err := Resolve(root, "partial")
	if ok || err == nil || !strings.Contains(err.Error(), flatLayoutInstruction) || !strings.Contains(err.Error(), flat) || !strings.Contains(err.Error(), folder) || strings.Contains(err.Error(), "spec folder is missing") {
		t.Fatalf("flat plus incomplete folder resolve = ok %v, err %v", ok, err)
	}
}

func TestResolveRefusesIncompleteFolderWithoutFlat(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "specs", "partial", "spec.md")
	if err := os.MkdirAll(filepath.Dir(missing), 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, _, ok, err := Resolve(root, "partial")
	if ok || err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("incomplete folder resolve = ok %v, err %v", ok, err)
	}
}

// retireRepo builds the minimal repository bench spec retire reads and deletes from: a
// folder spec committed at HEAD with body, plus any extraFiles (repo-relative path to
// content) committed alongside it — used to place a roadmap/FT<n>.md detail file. It
// returns a linked worktree of that repository, because retire refuses the primary
// checkout; retirePrimary returns the primary checkout instead.
func retireRepo(t *testing.T, slug, body string, extraFiles map[string]string) (root string) {
	t.Helper()
	return retireWorktree(t, retirePrimary(t, slug, body, extraFiles))
}

// retirePrimary builds the repository itself and returns its primary checkout.
func retirePrimary(t *testing.T, slug, body string, extraFiles map[string]string) (root string) {
	t.Helper()
	root = t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "a@b.c")
	runGit(t, root, "config", "user.name", "a")
	writeFolderSpec(t, root, slug, body)
	for path, content := range extraFiles {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-qm", "base")
	return root
}

// retireWorktree adds a linked worktree of primary and returns its resolved toplevel —
// the checkout shape a Bench phase runs a retire from.
func retireWorktree(t *testing.T, primary string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "linked")
	runGit(t, primary, "worktree", "add", "-q", "-b", "topic", root)
	out, err := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// runGit runs one git subcommand against root, failing the test on error.
func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// runRetire runs retireCommand from root, the way the CLI does: RepoBase resolves
// through the process cwd.
func runRetire(t *testing.T, root, arg string) (string, int) {
	t.Helper()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	return retireCommand([]string{arg})
}

// TestRetireNextLineNamesTheBoardRemainder covers FC1-FC4 and FC6, and the last-line
// and NBSP edges: the next: line's exact text as a function of the spec's Roadmap:
// value and whether roadmap/FT<n>.md exists on disk.
func TestRetireNextLineNamesTheBoardRemainder(t *testing.T) {
	generic := "next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/FT<n>.md detail file, commit as `spec-retire: s`\n"
	cases := []struct {
		name       string
		body       string
		extraFiles map[string]string
		want       string
	}{
		{
			name:       "FC1: FT7 with existing detail file names both",
			body:       "Status: implemented\nRoadmap: FT7\n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
		{
			name: "FC2: no Roadmap line prints the generic line",
			body: "Status: implemented\n",
			want: generic,
		},
		{
			name: "FC3: FT7 with no detail file names FT7 and no detail path",
			body: "Status: implemented\nRoadmap: FT7\n",
			want: "next: promote durable content, remove the ROADMAP row FT7, commit as `spec-retire: s`\n",
		},
		{
			name: "FC4: lowercase ft7 prints the generic line",
			body: "Status: implemented\nRoadmap: ft7\n",
			want: generic,
		},
		{
			name: "FC4: internal space FT 7 prints the generic line",
			body: "Status: implemented\nRoadmap: FT 7\n",
			want: generic,
		},
		{
			name: "FC4: empty Roadmap value prints the generic line",
			body: "Status: implemented\nRoadmap:\n",
			want: generic,
		},
		{
			name:       "FC6: fenced Roadmap line is not read",
			body:       "Status: implemented\n\n```\nRoadmap: FT7\n```\n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       generic,
		},
		{
			name:       "edge: Roadmap line is the last line with no trailing newline",
			body:       "Status: implemented\nRoadmap: FT7",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
		{
			name:       "edge: NBSP around the id trims like a space",
			body:       "Status: implemented\nRoadmap:  FT7 \n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := retireRepo(t, "s", tc.body, tc.extraFiles)
			out, code := runRetire(t, root, "s")
			if code != 0 {
				t.Fatalf("code = %d, want 0; out = %q", code, out)
			}
			lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
			got := lines[len(lines)-1] + "\n"
			if got != tc.want {
				t.Errorf("next line = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRetireStagedRefusesAndDeletesNothing pins FC5: a staged spec still refuses at
// exit 1, and the spec file survives untouched.
func TestRetireStagedRefusesAndDeletesNothing(t *testing.T) {
	root := retireRepo(t, "s", "Status: staged\nRoadmap: FT7\n", nil)
	specPath := filepath.Join(root, "specs", "s", "spec.md")
	before, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	after, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("spec file was deleted: %v", err)
	}
	if string(after) != string(before) {
		t.Errorf("spec file content changed on refusal")
	}
}

// TestRetireDeletesTheFolderAndExitsZero pins FA4: the surviving verb still removes a
// merged-implemented spec's whole folder at exit 0.
func TestRetireDeletesTheFolderAndExitsZero(t *testing.T) {
	root := retireRepo(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	out, code := runRetire(t, root, "s")
	if code != 0 {
		t.Fatalf("code = %d, want 0; out = %q", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "specs", "s")); !os.IsNotExist(err) {
		t.Errorf("specs/s still present: err = %v", err)
	}
}

// TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing pins the primary-checkout refusal:
// retire prints the one shared refusal and leaves the merged-implemented spec on disk.
func TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing(t *testing.T) {
	root := retirePrimary(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	if want := usage.PrimaryCheckoutRefusal() + "\n"; out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
	if _, err := os.Stat(filepath.Join(root, "specs", "s", "spec.md")); err != nil {
		t.Errorf("spec was deleted on refusal: %v", err)
	}
}
