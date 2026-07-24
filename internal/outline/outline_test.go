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

func TestCommandBoundsRowsAndFullRetainsMetadata(t *testing.T) {
	root := outlineRepo(t)
	var source strings.Builder
	source.WriteString("package x\n")
	for i := 0; i < 201; i++ {
		fmt.Fprintf(&source, "func Symbol%d() {}\n", i)
	}
	writeOutlineFile(t, root, "many.go", source.String())
	gitAddOutline(t, root)

	bounded, code := Command(nil)
	if code != 0 || !strings.HasPrefix(bounded, "outline[200]{") {
		t.Fatalf("bounded code/output = %d\n%s", code, bounded)
	}
	for _, want := range []string{"tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated", `  "1","1","0","201","200","1","true"`} {
		if !strings.Contains(bounded, want) {
			t.Fatalf("bounded metadata missing %q:\n%s", want, bounded)
		}
	}

	full, code := Command([]string{"--full"})
	if code != 0 || !strings.HasPrefix(full, "outline[201]{") || !strings.Contains(full, `  "1","1","0","201","201","0","false"`) {
		t.Fatalf("full code/output = %d\n%s", code, full)
	}
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
	out, code := Command(nil)
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
