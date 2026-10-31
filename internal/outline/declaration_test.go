package outline

import (
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/responsebound"
)

// declarationSource holds each declaration family, an unattached comment, adjacent
// declarations, a grouped declaration, and a trailing line comment.
const declarationSource = `// Package fixture holds declarations.
package fixture

import "fmt"

// An unattached comment that a blank line separates.

// Run starts the work.
// It has two doc lines.
func Run() { fmt.Println("run") }
func Adjacent() {}

// Store holds values.
type Store struct{ values map[string]int }

// Get reads a value.
func (s *Store) Get(key string) int { return s.values[key] }

func (s Store) Len() int { return len(s.values) }

// Push appends.
func (l *List[T]) Push(item T) { l.items = append(l.items, item) }

// Limits groups two constants.
const (
	// Low is low.
	Low  = 1
	High = 2
)

var Single = "single" // trailing
`

// TestDeclarationSpans grades each span against independently authored source text: it
// starts at the attached doc comment and ends at the declaration end, and it excludes the
// unattached comment, the adjacent declaration, and the trailing line comment. Each name of
// a group has the complete group.
func TestDeclarationSpans(t *testing.T) {
	group := "// Limits groups two constants.\nconst (\n\t// Low is low.\n\tLow  = 1\n\tHigh = 2\n)"
	want := map[string]string{
		"Run":          "func // Run starts the work.\n// It has two doc lines.\nfunc Run() { fmt.Println(\"run\") }",
		"Adjacent":     "func func Adjacent() {}",
		"Store":        "type // Store holds values.\ntype Store struct{ values map[string]int }",
		"(*Store).Get": "method // Get reads a value.\nfunc (s *Store) Get(key string) int { return s.values[key] }",
		"Store.Len":    "method func (s Store) Len() int { return len(s.values) }",
		"(*List).Push": "method // Push appends.\nfunc (l *List[T]) Push(item T) { l.items = append(l.items, item) }",
		"Low":          "const " + group,
		"High":         "const " + group,
		"Single":       "var var Single = \"single\"",
	}
	declarations, err := Declarations("fixture.go", []byte(declarationSource))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, d := range declarations {
		got[d.Qualified()] = d.Kind + " " + declarationSource[d.Start:d.End]
	}
	if !maps.Equal(got, want) {
		t.Fatalf("declaration spans =\n%q\nwant\n%q", got, want)
	}
}

// symbolPage runs one --symbol call and answers its decoded page row and its reconstructed
// content.
func symbolPage(t *testing.T, args ...string) (map[string]any, string) {
	t.Helper()
	out, code := Command(args)
	document, err := axitest.DecodeDocument(out)
	if code != 0 || err != nil {
		t.Fatalf("Command(%q) code=%d err=%v out=%q, want one page", args, code, err, out)
	}
	rows, _ := document.Rows("page")
	page := rows[0].(map[string]any)
	contents, _ := document.Rows("content")
	var text strings.Builder
	for _, content := range contents {
		data, err := responsebound.DecodeContent(page["encoding"].(string), content.(map[string]any)["data"].(string))
		if err != nil {
			t.Fatal(err)
		}
		text.Write(data)
	}
	return page, text.String()
}

// TestDeclarationSelection grades the narrowing of --symbol over the tracked Go scope:
// ambiguity answers its candidates in file order with no body, and a receiver, a path, a
// file, or a filter narrows to one body. Absence is a present-false page, a malformed
// selected file refuses, and --symbol with --full refuses as usage.
func TestDeclarationSelection(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "b/b.go", "package b\n\nfunc Get() {}\n")
	writeOutlineFile(t, root, "a/a.go", "package a\n\nfunc Get() {}\n\ntype Store struct{}\n\nfunc (s *Store) Get() {}\n")
	writeOutlineFile(t, root, "a/generic.go", "package a\n\ntype List[T any] struct{}\n\nfunc (l List[T]) Get() {}\n")
	writeOutlineFile(t, root, "b/b_test.go", "package b\n\nfunc Only() {}\n")
	writeOutlineFile(t, root, "bad/bad.go", "package bad\n\nfunc Broken( {\n")
	gitAddOutline(t, root)

	out, code := Command([]string{"--symbol", "Get"})
	first, second := strings.Index(out, "a/a.go,3,Get,func"), strings.Index(out, "b/b.go,3,Get,func")
	if code != 1 || !strings.HasPrefix(out, "error: --symbol Get matches 2 declarations") || first < 0 || second < first || strings.Contains(out, "content[") ||
		!strings.Contains(out, "bench outline a/a.go --symbol Get") || !strings.Contains(out, "bench outline b/b.go --symbol Get") {
		t.Fatalf("ambiguous Get code=%d out=%q, want both candidates in file order with their commands and no body", code, out)
	}
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"a", "--symbol", "Get"}, "func Get() {}"},
		{[]string{"b/b.go", "--symbol", "Get"}, "func Get() {}"},
		{[]string{"--symbol", "Store.Get"}, "func (s *Store) Get() {}"},
		{[]string{"--symbol", "(*Store).Get"}, "func (s *Store) Get() {}"},
		{[]string{"--symbol", "List.Get"}, "func (l List[T]) Get() {}"},
		{[]string{"--symbol", "Only", "--test"}, "func Only() {}"},
	} {
		page, text := symbolPage(t, tc.args...)
		if page["present"] != true || page["kind"] != "declaration" || page["complete"] != true || text != tc.want {
			t.Fatalf("Command(%q) page=%v text=%q, want %q", tc.args, page, text, tc.want)
		}
	}
	for _, args := range [][]string{{"--symbol", "Missing"}, {"--symbol", "(*List).Get"}, {"--symbol", "Only", "--production"}} {
		page, text := symbolPage(t, args...)
		if page["present"] != false || page["complete"] != true || text != "" {
			t.Fatalf("Command(%q) page=%v text=%q, want a present-false page", args, page, text)
		}
	}
	for _, tc := range []struct {
		args   []string
		code   int
		prefix string
	}{
		{[]string{"bad", "--symbol", "Broken"}, 1, "error: bad/bad.go is not valid Go"},
		{[]string{"--full", "--symbol", "Get"}, 2, symbolUsageLine + " (--symbol and --full are mutually exclusive)"},
		{[]string{"--symbol", "a.b.c"}, 2, symbolUsageLine + " (--symbol takes"},
		{[]string{"--symbol", "(*T.M"}, 2, symbolUsageLine + " (--symbol takes"},
		{[]string{"--offset", "0", "--source", "x"}, 2, symbolUsageLine + " (the cursor flags take --symbol)"},
	} {
		out, code := Command(tc.args)
		if code != tc.code || !strings.HasPrefix(out, tc.prefix) || strings.Contains(out, "content[") {
			t.Fatalf("Command(%q) code=%d out=%q, want exit %d beginning %q", tc.args, code, out, tc.code, tc.prefix)
		}
	}
}
