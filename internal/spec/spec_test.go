package spec

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// writeFolderSpec writes content to <dir>/specs/<slug>/spec.md and returns the path.
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
	path := writeFolderSpec(t, dir, "x", "# spec\n\nStatus: staged\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out, code, _ := Command([]string{"implemented", "x"})
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
// offers show, retire, history, and widen only.
func TestSpecUsageNamesTheSurvivingSubcommands(t *testing.T) {
	out, code, _ := Command(nil)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(out, "expected a subcommand: show, retire, history, widen)") {
		t.Errorf("out = %q, want the show/retire/history/widen usage line", out)
	}
}

func TestSpecHelpPrintsThePublicSubcommandInventory(t *testing.T) {
	const want = "usage: bench spec <subcommand>\n\nsubcommands:\n  show <slug-or-path> [--section stories|status|<literal-H2-title>] [--offset <byte-offset> --source <identity>]\n  retire <slug>\n  history <slug>\n  widen <slug> --ticket <basename> [--decision <reference>] <entry>...\n"
	for _, arg := range []string{"--help", "-h"} {
		out, code, _ := Command([]string{arg})
		if code != 0 {
			t.Fatalf("Command(%q) code = %d, want 0; out = %q", arg, code, out)
		}
		if out != want {
			t.Errorf("Command(%q) out = %q, want %q", arg, out, want)
		}
	}
}

func TestSpecSubcommandHelpUsesDeclaredUsage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"retire", "usage: bench spec retire <spec.md | slug>\n"},
		{"history", "usage: bench spec history <spec.md | slug>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code, _ := Command([]string{tt.name, "help"})
			if code != 0 {
				t.Errorf("help code = %d, want 0; out = %q", code, out)
			}
			if out != tt.want {
				t.Errorf("help out = %q, want %q", out, tt.want)
			}
		})
	}
}

func TestSpecSubcommandUsageRefusals(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{"retire unknown flag", "retire", []string{"--unknown"}},
		{"retire missing operand", "retire", nil},
		{"history unknown flag", "history", []string{"--unknown"}},
		{"history missing operand", "history", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code, _ := Command(append([]string{tt.command}, tt.args...))
			if code != 2 {
				t.Errorf("code = %d, want 2; out = %q", code, out)
			}
			if !strings.HasPrefix(out, "usage: bench spec ") {
				t.Errorf("out = %q, want usage", out)
			}
		})
	}
}

func TestResolveConvention(t *testing.T) {
	dir := t.TempDir()
	slugPath := writeFolderSpec(t, dir, "mine", "Status: staged\n")

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
	writeFolderSpec(t, root, "anchored", "Status: staged\n")
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
	writeFolderSpec(t, root, "flat", "Status: staged\n")
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
					capability.Unavailable(t, capability.Symlink, err.Error())
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
		capability.Unavailable(t, capability.Symlink, err.Error())
	}

	fifo := filepath.Join(root, "specs", "fifo", "spec.md")
	if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		capability.Unavailable(t, capability.Fifo, err.Error())
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
