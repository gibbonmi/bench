// Package consumers resolves real Go reference edges. It answers "who consumes this
// symbol" with resolved type information rather than a textual sweep, so an agent can
// cite an enumeration instead of a sampled grep.
//
// The package splits in two. A pure analysis core (resolve.go, rows.go) consumes
// already-typed packages and returns rows; it never reads the tree and never spawns a
// process. A package-internal loader seam (loader.go) wraps go/packages and is the only
// site that shells out to the go tool. Core tests therefore type-check fixture source in
// process, and exactly one test drives the real loader.
//
// The tool IDENTIFIES resolved reference edges. It does not see reflection,
// go:linkname, plugin, or exec edges; a textual sweep stays the candidate-class
// citation for those.
package consumers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// Package is one typed Go package the core consumes. It is the core's whole input
// contract: the loader seam produces it from go/packages, and a test produces the same
// shape from go/parser plus go/types. Neither producer is privileged.
type Package struct {
	PkgPath string
	Fset    *token.FileSet
	Files   []*ast.File
	Types   *types.Package
	Info    *types.Info
}

// Row is one resolved reference to the queried symbol. File is repository-relative, so
// the cell names what an agent would open. Column is carried for the sort order only and
// is never printed: two references on one line still order deterministically.
type Row struct {
	File      string
	Line      int
	Column    int
	Via       string
	Enclosing string
}

// rowFields is the consumers table schema. The queried symbol is constant across the
// whole table, so it gets no column.
var rowFields = []string{"file", "line", "via", "enclosing"}

// Render emits rows as the AXI `consumers[N]{file,line,via,enclosing}:` block. An empty
// slice renders the definitive empty table rather than nothing.
func Render(rows []Row) (string, error) {
	cells := make([][]any, len(rows))
	for i, r := range rows {
		cells[i] = []any{r.File, r.Line, r.Via, r.Enclosing}
	}
	return toon.TableTyped("consumers", rowFields, cells)
}

// relPath renders a file-set filename relative to root. A path outside root, or an
// empty root, falls back to the filename as the loader reported it: a wrong-looking
// absolute path is more useful than a silently dropped row.
func relPath(root, filename string) string {
	if root == "" {
		return filename
	}
	rel, err := filepath.Rel(root, filename)
	if err != nil {
		return filename
	}
	return rel
}

// insideRoot reports whether filename names a file under root. It is the one test of
// "this file belongs to the tree the answer is about", and the loader applies it before
// any file reaches the core. An out-of-root file is never a tree the citation can
// replay: the citation promises a replay at one checkout sha, so a row pointing outside
// that checkout names a file no reviewer can open and no re-run can reproduce. An empty
// root asserts nothing, so every file passes.
func insideRoot(root, filename string) bool {
	if root == "" {
		return true
	}
	rel, err := filepath.Rel(root, filename)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Find resolves query to exactly one declaration and enumerates its references. It is
// the core's whole answer path for an unambiguous query. An ambiguous query is an error
// here; the caller that owns the candidates table calls Resolve directly.
func Find(pkgs []*Package, query, root string) ([]Row, error) {
	matches, err := Resolve(pkgs, query)
	if err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("symbol %q matches %d declarations", query, len(matches))
	}
	return Rows(pkgs, matches[0].Obj, root), nil
}
