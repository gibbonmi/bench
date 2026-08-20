package outline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// H14: --full emits symbol rows repository-wide, uncapped, with completeness metadata.
func TestFullFormEmitsSymbolRowsRepositoryWide(t *testing.T) {
	root := outlineRepo(t)
	var source strings.Builder
	source.WriteString("package x\n")
	for i := 0; i < 201; i++ {
		fmt.Fprintf(&source, "func Symbol%d() {}\n", i)
	}
	writeOutlineFile(t, root, "many.go", source.String())
	gitAddOutline(t, root)

	full, code := Command([]string{"--full"})
	if code != 0 || !strings.HasPrefix(full, "outline[201]{file,line,kind,name}:\n") {
		t.Fatalf("full code/output = %d\n%s", code, headOf(full))
	}
	for _, want := range []string{"tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated", `  "1","1","0","201","201","0","false"`} {
		if !strings.Contains(full, want) {
			t.Fatalf("full metadata missing %q:\n%s", want, headOf(full))
		}
	}
}

// H12: the bare form emits meta and one row per scanned top-level directory carrying
// that directory's whole-subtree symbol count, and no symbol rows. The nested file makes
// the collapse observable: a per-scanned-directory roll-up would report a/deep separately.
func TestBareFormSummarizesTopLevelDirectories(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "a/deep/x.go", "package a\nfunc One() {}\nfunc Two() {}\n")
	writeOutlineFile(t, root, "a/y.go", "package a\nfunc Three() {}\n")
	writeOutlineFile(t, root, "b/z.go", "package b\nfunc Four() {}\n")
	writeOutlineFile(t, root, "root.go", "package main\nfunc Five() {}\n")
	gitAddOutline(t, root)

	out, code := Command(nil)
	if code != 0 || !strings.HasPrefix(out, "outline_dirs[3]{dir,symbols}:\n") {
		t.Fatalf("bare form is not a top-level summary: code=%d\n%s", code, headOf(out))
	}
	for _, want := range []string{`  a,"3"`, `  b,"1"`, `  .,"1"`, "outline_meta[1]{"} {
		if !strings.Contains(out, want) {
			t.Fatalf("bare output missing %q:\n%s", want, headOf(out))
		}
	}
	if strings.Contains(out, "outline[") || strings.Contains(out, "a/deep") {
		t.Fatalf("bare form did not collapse to top level:\n%s", headOf(out))
	}
}

// H13: a path argument keeps emitting symbol rows scoped to that path, including a path
// carrying spaces and glob characters.
func TestPathFormEmitsSymbolRowsForThatPath(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "a/x.go", "package a\nfunc One() {}\n")
	writeOutlineFile(t, root, "weird [dir] name/z.go", "package w\nfunc Odd() {}\n")
	gitAddOutline(t, root)

	scoped, code := Command([]string{"a"})
	if code != 0 || !strings.HasPrefix(scoped, "outline[1]{file,line,kind,name}:\n") {
		t.Fatalf("path form code/output = %d\n%s", code, headOf(scoped))
	}
	if !strings.Contains(scoped, "One") || strings.Contains(scoped, "Odd") {
		t.Fatalf("path form did not scope to the path:\n%s", headOf(scoped))
	}

	globbed, code := Command([]string{"weird [dir] name"})
	if code != 0 || !strings.HasPrefix(globbed, "outline[1]{file,line,kind,name}:\n") {
		t.Fatalf("glob-charactered path code/output = %d\n%s", code, headOf(globbed))
	}
	if !strings.Contains(globbed, "Odd") {
		t.Fatalf("glob-charactered path lost its rows:\n%s", headOf(globbed))
	}
}

// H15: the bare form's row count equals the scanned top-level directory count, with no
// cap applied. The fixture crosses the retired 200-row limit so a reintroduced cap is
// observable, and nests each file so the count is a subtree roll-up.
func TestBareFormRowCountEqualsTopLevelDirectoryCount(t *testing.T) {
	root := outlineRepo(t)
	for i := 0; i < 201; i++ {
		writeOutlineFile(t, root, fmt.Sprintf("d%03d/nested/x.go", i), "package x\nfunc Only() {}\n")
	}
	gitAddOutline(t, root)

	out, code := Command(nil)
	if code != 0 || !strings.HasPrefix(out, "outline_dirs[201]{dir,symbols}:\n") {
		t.Fatalf("bare row count is bounded rather than complete: code=%d\n%s", code, headOf(out))
	}
	if !strings.Contains(out, `  d200,"1"`) {
		t.Fatalf("the last scanned top-level directory is missing from the summary:\n%s", headOf(out))
	}
}

// H16: absent and present-but-empty are distinct definitive empty states, and the
// row-bearing forms keep their own typed zero-row table.
func TestEmptyStatesDistinguishAbsentFromPresentButEmpty(t *testing.T) {
	outlineRepo(t)
	absent, code := Command(nil)
	if code != 0 || !strings.HasPrefix(absent, "outline_dirs[0]{dir,symbols}:\n") {
		t.Fatalf("absent tree code/output = %d\n%s", code, headOf(absent))
	}
	if !strings.Contains(absent, `  "0","0","0","0","0","0","false"`) {
		t.Fatalf("absent tree lost its zeroed metadata:\n%s", headOf(absent))
	}

	root := outlineRepo(t)
	writeOutlineFile(t, root, "notes.rs", "fn main() {}\n")
	gitAddOutline(t, root)
	present, code := Command(nil)
	if code != 0 || !strings.HasPrefix(present, "outline_dirs[1]{dir,symbols}:\n") {
		t.Fatalf("present-but-empty tree code/output = %d\n%s", code, headOf(present))
	}
	if !strings.Contains(present, `  .,"0"`) || !strings.Contains(present, `  "1","1","0","0","0","0","false"`) {
		t.Fatalf("present-but-empty tree is not distinct from absent:\n%s", headOf(present))
	}

	full, code := Command([]string{"--full"})
	if code != 0 || !strings.HasPrefix(full, "outline[0]{file,line,kind,name}:\n") {
		t.Fatalf("row-bearing empty state = %d\n%s", code, headOf(full))
	}
}

// headOf keeps a failure message readable when the subject is a long table.
func headOf(out string) string {
	lines := strings.Split(out, "\n")
	if len(lines) > 12 {
		lines = append(lines[:12], "...")
	}
	return strings.Join(lines, "\n")
}

func TestCommandNamesSizeBinaryAndNonregularSkips(t *testing.T) {
	root := outlineRepo(t)
	exact := "package x\nfunc Exact() {}\n"
	exact += strings.Repeat(" ", (2<<20)-len(exact))
	writeOutlineFile(t, root, "exact.go", exact)
	writeOutlineFile(t, root, "large.go", exact+"x")
	writeOutlineFile(t, root, "binary.go", "package x\x00func Hidden() {}\n")
	writeOutlineFile(t, root, "target.go", "package x\nfunc Target() {}\n")
	if err := os.Symlink("target.go", filepath.Join(root, "link.go")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlink unavailable: %v", err))
	}
	gitAddOutline(t, root)
	// --full is the row-bearing repository-wide form; skips ride the same envelope.
	out, code := Command([]string{"--full"})
	if code != 0 {
		t.Fatalf("code = %d; out=%s", code, out)
	}
	for _, want := range []string{"exact.go,\"2\",func,Exact", "binary.go,binary", "large.go,oversized", "link.go,nonregular"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestCommandNamesUnreadableSkip(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "unreadable.go", "package x\nfunc Hidden() {}\n")
	gitAddOutline(t, root)
	old := openOutlineFile
	openOutlineFile = func(path string) (*os.File, error) {
		if strings.HasSuffix(path, "unreadable.go") {
			return nil, os.ErrPermission
		}
		return os.Open(path)
	}
	t.Cleanup(func() { openOutlineFile = old })
	out, code := Command(nil)
	if code != 0 || !strings.Contains(out, "unreadable.go,unreadable") {
		t.Fatalf("unreadable file not classified: code=%d\n%s", code, out)
	}
}

func outlineRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q", root)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	return root
}

func writeOutlineFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitAddOutline(t *testing.T, root string) {
	t.Helper()
	if out, err := exec.Command("git", "-C", root, "add", "-A").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
}

// symbolsCase drives the pure parser seam: feed one buffer for a language and assert
// the exact (line,kind,name) rows, which is where the per-language pattern table's
// correctness lives.
type symbolsCase struct {
	name    string
	path    string
	content string
	want    []Symbol
}

func TestSymbols(t *testing.T) {
	cases := []symbolsCase{
		{
			name: "go func method and type",
			path: "x.go",
			content: "package x\n" +
				"\n" +
				"type Widget struct{}\n" +
				"\n" +
				"func New() *Widget { return nil }\n" +
				"\n" +
				"func (w *Widget) Name() string { return \"\" }\n",
			want: []Symbol{
				{Line: 3, Kind: "type", Name: "Widget"},
				{Line: 5, Kind: "func", Name: "New"},
				{Line: 7, Kind: "func", Name: "Name"},
			},
		},
		{
			name: "shell paren form and function keyword",
			path: "x.sh",
			content: "#!/usr/bin/env bash\n" +
				"route_porcelain() {\n" +
				"  :\n" +
				"}\n" +
				"function run_gate {\n" +
				"  :\n" +
				"}\n",
			want: []Symbol{
				{Line: 2, Kind: "function", Name: "route_porcelain"},
				{Line: 5, Kind: "function", Name: "run_gate"},
			},
		},
		{
			name: "markdown atx headings",
			path: "x.md",
			content: "# Title\n" +
				"\n" +
				"intro prose\n" +
				"\n" +
				"### A sub heading\n" +
				"not a #heading without space\n",
			want: []Symbol{
				{Line: 1, Kind: "heading", Name: "Title"},
				{Line: 5, Kind: "heading", Name: "A sub heading"},
			},
		},
		{
			name: "python def and class",
			path: "x.py",
			content: "import os\n" +
				"\n" +
				"class Widget:\n" +
				"    def render(self):\n" +
				"        return 1\n",
			want: []Symbol{
				{Line: 3, Kind: "class", Name: "Widget"},
				{Line: 4, Kind: "def", Name: "render"},
			},
		},
		{
			name: "js/ts function class and exported const arrow",
			path: "x.ts",
			content: "function plain() {}\n" +
				"class Widget {}\n" +
				"export const make = (a) => a + 1;\n",
			want: []Symbol{
				{Line: 1, Kind: "function", Name: "plain"},
				{Line: 2, Kind: "class", Name: "Widget"},
				{Line: 3, Kind: "const", Name: "make"},
			},
		},
		{
			name: "final line without trailing newline is still scanned",
			path: "x.go",
			content: "package x\n" +
				"func Tail() {}", // no trailing \n
			want: []Symbol{
				{Line: 2, Kind: "func", Name: "Tail"},
			},
		},
		{
			name:    "unknown extension yields no rows",
			path:    "x.rs",
			content: "fn main() {}\n",
			want:    nil,
		},
		{
			name:    "present-but-empty buffer yields no rows",
			path:    "x.go",
			content: "",
			want:    nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Symbols(c.path, []byte(c.content))
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Symbols(%q) = %#v\nwant %#v", c.path, got, c.want)
			}
		})
	}
}
