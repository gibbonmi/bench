package outline

import (
	"maps"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
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

func (m Map[K, V]) Get(key K) V { return m.values[key] }

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
		"Map.Get":      "method func (m Map[K, V]) Get(key K) V { return m.values[key] }",
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

// candidateRoutes runs one ambiguous --symbol call, and answers the argv of each help
// command after `bench`. It requires exit 1 and no body.
func candidateRoutes(t *testing.T, args ...string) (string, [][]string) {
	t.Helper()
	out, code := Command(args)
	document, err := axitest.DecodeDocument(out)
	if code != 1 || err != nil || strings.Contains(out, "content[") {
		t.Fatalf("Command(%q) code=%d err=%v out=%q, want an ambiguity refusal with no body", args, code, err, out)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatal(err)
	}
	var routes [][]string
	for _, action := range actions {
		help, err := toon.Table("help", []string{"cmd", "why"}, [][]string{{action.Cmd, action.Why}})
		if err != nil {
			t.Fatal(err)
		}
		argv, err := axitest.RecoverHelpCommandArgv(help)
		if err != nil || len(argv) < 2 || argv[0] != "bench" {
			t.Fatalf("route %q: %v", action.Cmd, err)
		}
		routes = append(routes, argv[1:])
	}
	return out, routes
}

// TestDeclarationSelection grades the narrowing of --symbol over the tracked Go scope:
// ambiguity answers its candidates in file order with no body, and a receiver, a path, a
// file, or a filter narrows to one body. A bare name that only methods declare is
// ambiguous, and two candidates that share a file and a qualified name each route to the
// lines of their own span in that file.
// Absence is a present-false page, a malformed or oversized file in scope refuses, an
// invalid fixture does not, and --symbol with --full refuses as usage.
func TestDeclarationSelection(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "b/b.go", "package b\n\nfunc Get() {}\n")
	writeOutlineFile(t, root, "a/a.go", "package a\n\nfunc Get() {}\n\ntype Store struct{}\n\nfunc (s *Store) Get() {}\n\nfunc (s *Store) Put() {}\n")
	writeOutlineFile(t, root, "a/generic.go", "package a\n\ntype List[T any] struct{}\n\nfunc (l List[T]) Get() {}\n\nfunc (l List[T]) Put() {}\n")
	writeOutlineFile(t, root, "a/testdata/broken.go", "package x\n\nfunc Get( {\n")
	writeOutlineFile(t, root, "b/b_test.go", "package b\n\nfunc Only() {}\n")
	writeOutlineFile(t, root, "bad/bad.go", "package bad\n\nfunc Broken( {\n")
	writeOutlineFile(t, root, "d/init.go", "package d\n\nfunc init() {}\n\n// init runs last.\nfunc init() {\n}\n")
	gitAddOutline(t, root)

	out, code := Command([]string{"--symbol", "Get"})
	first, second := strings.Index(out, "a/a.go,3,Get,func"), strings.Index(out, "b/b.go,3,Get,func")
	if code != 1 || !strings.HasPrefix(out, "error: --symbol Get matches 2 declarations") || first < 0 || second < first || strings.Contains(out, "content[") ||
		!strings.Contains(out, "bench outline a/a.go --symbol Get") || !strings.Contains(out, "bench outline b/b.go --symbol Get") {
		t.Fatalf("ambiguous Get code=%d out=%q, want both candidates in file order with their commands and no body", code, out)
	}
	// A bare name that no package-level declaration has lists each same-name method, and
	// each route reads its method.
	out, routes := candidateRoutes(t, "--symbol", "Put")
	if !strings.Contains(out, "a/a.go,9,(*Store).Put,method") || !strings.Contains(out, "a/generic.go,7,List.Put,method") || len(routes) != 2 {
		t.Fatalf("method-only Put out=%q routes=%q, want both methods as candidates", out, routes)
	}
	for i, want := range []string{"func (s *Store) Put() {}", "func (l List[T]) Put() {}"} {
		if _, text := symbolPage(t, routes[i][1:]...); text != want {
			t.Fatalf("route %q read %q, want %q", routes[i], text, want)
		}
	}
	// Two init functions share a file and a qualified name, so each names its line and its
	// route reads the lines of its own span, from its doc comment through its end.
	out, routes = candidateRoutes(t, "--symbol", "init")
	if !strings.Contains(out, "d/init.go,3,init,func") || !strings.Contains(out, "d/init.go,6,init,func") ||
		!reflect.DeepEqual(routes, [][]string{{"read", "d/init.go", "--lines", "3:3"}, {"read", "d/init.go", "--lines", "5:7"}}) {
		t.Fatalf("same-file init out=%q routes=%q, want each line and a bench read route of each span", out, routes)
	}
	// From a subdirectory, each route names its file from that directory and runs there.
	t.Chdir(filepath.Join(root, "b"))
	_, routes = candidateRoutes(t, "--symbol", "Get")
	if !reflect.DeepEqual(routes, [][]string{{"outline", "../a/a.go", "--symbol", "Get"}, {"outline", "b.go", "--symbol", "Get"}}) {
		t.Fatalf("routes from b = %q, want cwd-relative operands", routes)
	}
	for _, route := range routes {
		if _, text := symbolPage(t, route[1:]...); text != "func Get() {}" {
			t.Fatalf("route %q from b read %q", route, text)
		}
	}
	if _, routes = candidateRoutes(t, "--symbol", "init"); !reflect.DeepEqual(routes, [][]string{{"read", "../d/init.go", "--lines", "3:3"}, {"read", "../d/init.go", "--lines", "5:7"}}) {
		t.Fatalf("same-file route from b = %q, want the cwd-relative file", routes)
	}
	t.Chdir(root)
	for _, tc := range []struct {
		args []string
		want string
	}{
		// a holds an invalid testdata Go file, and a fixture is outside the Go scope.
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
		{[]string{"--offset", "0"}, 2, symbolUsageLine + " (the cursor flags take --symbol)"},
		{[]string{"--source", "x"}, 2, symbolUsageLine + " (the cursor flags take --symbol)"},
	} {
		out, code := Command(tc.args)
		if code != tc.code || !strings.HasPrefix(out, tc.prefix) || strings.Contains(out, "content[") {
			t.Fatalf("Command(%q) code=%d out=%q, want exit %d beginning %q", tc.args, code, out, tc.code, tc.prefix)
		}
	}
	// An oversized Go file in scope can hold the declaration, so the call refuses.
	writeOutlineFile(t, root, "big/big.go", "package big\n"+strings.Repeat(" ", int(bounds.OutlineFileLimit)))
	gitAddOutline(t, root)
	if out, code := Command([]string{"big", "--symbol", "Get"}); code != 1 || !strings.HasPrefix(out, "error: big/big.go is oversized") || strings.Contains(out, "content[") {
		t.Fatalf("oversized scope code=%d out=%q, want a refusal", code, out)
	}
}
