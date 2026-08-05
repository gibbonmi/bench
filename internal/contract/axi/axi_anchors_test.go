package axi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/toon"
)

func testAXIAnchorsRegistryRows(t *testing.T) {
	contract.NoteContractFailure(t, "AXI anchors registry-row contract failed")
	f := contract.NewFixture(t)
	for _, path := range []string{".bench/BENCH.md", ".agents/skills/bench-craft-tickets/SKILL.md"} {
		want, sectionScoped := expectedAnchorTable(t, path)
		if !sectionScoped {
			t.Fatalf("registry path %q has no section-scoped anchor", path)
		}

		out := f.Bench("anchors", path)
		out.RequireExit(0)
		if out.Stdout != want {
			t.Fatalf("bench anchors %s output differs from registry order\nwant:\n%s\ngot:\n%s", path, want, out.Stdout)
		}
	}
}

func testAXIAnchorsEmpty(t *testing.T) {
	f := contract.NewFixture(t)
	for _, path := range []string{"not/pinned.md", "space [glob]*.md"} {
		out := f.Bench("anchors", path)
		out.RequireExit(0)
		if out.Stdout != "anchors[0]{kind,section,needle}:\n" {
			t.Fatalf("definitive empty anchors output for %q = %q", path, out.Stdout)
		}
	}
}

func testAXIAnchorsUsage(t *testing.T) {
	f := contract.NewFixture(t)
	for _, args := range [][]string{
		{"anchors"},
		{"anchors", ".bench/BENCH.md", "extra.md"},
		{"anchors", "--unknown"},
	} {
		out := f.Bench(args...)
		out.RequireExit(2)
		if out.Stderr != "" {
			t.Fatalf("bench %v usage wrote stderr:\n%s", args, out.Stderr)
		}
		requireContainsFold(t, out.Stdout, "usage: bench anchors")
	}
}

func testAXIAnchorsSubdirectoryRoutes(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/BENCH.md", "fixture target exists\n")
	deep := filepath.Join(f.Root, "sub", "deeper")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	want, _ := expectedAnchorTable(t, ".bench/BENCH.md")
	binary := filepath.Join(contract.SubjectRoot(t), "dist", "bench")
	for _, arg := range []string{
		".bench/BENCH.md",
		filepath.Join("..", "..", ".bench", "BENCH.md"),
	} {
		launcher := runBenchInDir(t, f, deep, "anchors", arg)
		launcher.RequireExit(0)
		if launcher.Stdout != want {
			t.Fatalf("subdirectory launcher output for %q differs from repo-relative query\nwant:\n%s\ngot:\n%s", arg, want, launcher.Stdout)
		}

		direct := contract.RunAt(t, f, deep, nil, binary, "anchors", arg)
		direct.RequireExit(0)
		if direct.Stdout != want {
			t.Fatalf("subdirectory binary output for %q differs from repo-relative query\nwant:\n%s\ngot:\n%s", arg, want, direct.Stdout)
		}
	}
}

func expectedAnchorTable(t *testing.T, path string) (string, bool) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(contract.SubjectRoot(t), filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("read anchored file %s: %v", path, err)
	}

	var rows [][]string
	sectionScoped := false
	for _, anchor := range anchors.Entries() {
		if anchor.File != path {
			continue
		}
		kind := anchorKindName(t, anchor.Kind)
		rows = append(rows, []string{kind, anchor.Section, anchor.Needle})
		text := string(content)
		switch anchor.Kind {
		case anchors.Forbid:
			text = anchors.StripHTMLComments(text)
		case anchors.RequireInSection, anchors.ForbidInSection:
			sectionScoped = true
			text = anchors.MarkdownH2Section(anchors.StripHTMLComments(text), anchor.Section)
		}
		if !anchors.Satisfied(anchor.Kind, text, anchor.Needle) {
			t.Fatalf("registry needle %q does not match real content in %s section %q", anchor.Needle, path, anchor.Section)
		}
	}
	if len(rows) == 0 {
		t.Fatalf("registry has no anchors for pinned path %s", path)
	}
	out, err := toon.Table("anchors", []string{"kind", "section", "needle"}, rows)
	if err != nil {
		t.Fatalf("render expected anchors table: %v", err)
	}
	return out, sectionScoped
}

func anchorKindName(t *testing.T, kind anchors.Kind) string {
	t.Helper()
	switch kind {
	case anchors.Require:
		return "require"
	case anchors.Forbid:
		return "forbid"
	case anchors.RequireInSection:
		return "require-in-section"
	case anchors.ForbidInSection:
		return "forbid-in-section"
	default:
		t.Fatalf("unknown anchor kind %d", kind)
		return ""
	}
}
