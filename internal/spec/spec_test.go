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

func TestFlipRewritesOnlyTheStatusLine(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "space separator, trailing newline",
			in:   "# spec\n\nStatus: staged\n\n## body\n",
			want: "# spec\n\nStatus: implemented\n\n## body\n",
		},
		{
			name: "tab separator preserved",
			in:   "Status:\tstaged\nbody\n",
			want: "Status:\timplemented\nbody\n",
		},
		{
			name: "no trailing newline preserved",
			in:   "# spec\nStatus: staged",
			want: "# spec\nStatus: implemented",
		},
		{
			name: "trailing whitespace on status line preserved",
			in:   "Status:  staged  \n",
			want: "Status:  implemented  \n",
		},
		{
			name: "other staged word untouched",
			in:   "Status: staged\nThe work is staged elsewhere.\n",
			want: "Status: implemented\nThe work is staged elsewhere.\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := writeSpec(t, dir, "s", tc.in)
			resolved, err := Flip(dir, "s")
			if err != nil {
				t.Fatalf("Flip: %v", err)
			}
			if resolved != path {
				t.Errorf("resolved = %q, want %q", resolved, path)
			}
			got, _ := os.ReadFile(path)
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

func TestFlipErrors(t *testing.T) {
	cases := []struct {
		name    string
		content string // "" means don't create the file (not-found)
		wantSub string
	}{
		{"not found", "", "not found"},
		{"no staged line", "# spec\nStatus: draft\n", "no `Status: staged`"},
		{"already implemented", "# spec\nStatus: implemented\n", "no `Status: staged`"},
		{"more than one staged line", "Status: staged\nStatus: staged\n", "expected exactly one"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.content != "" {
				writeSpec(t, dir, "s", tc.content)
			}
			_, err := Flip(dir, "s")
			if err == nil {
				t.Fatal("Flip: expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantSub)
			}
		})
	}
}

func TestFlipReRunIsNonDestructive(t *testing.T) {
	dir := t.TempDir()
	path := writeSpec(t, dir, "s", "Status: staged\n")
	if _, err := Flip(dir, "s"); err != nil {
		t.Fatalf("first flip: %v", err)
	}
	after, _ := os.ReadFile(path)
	if _, err := Flip(dir, "s"); err == nil {
		t.Fatal("second flip: expected error on already-implemented spec")
	}
	again, _ := os.ReadFile(path)
	if string(again) != string(after) {
		t.Errorf("second flip mutated the file: %q -> %q", after, again)
	}
}

func TestResolveConvention(t *testing.T) {
	dir := t.TempDir()
	slugPath := writeSpec(t, dir, "mine", "Status: staged\n")

	// bare slug resolves to <base>/specs/<slug>/spec.md
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
