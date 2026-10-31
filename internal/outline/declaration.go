package outline

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
)

// declarationKind is the page selection kind of a --symbol read.
const declarationKind = "declaration"

// Declaration is one top-level Go declaration that a --symbol read can select. Name is the
// declared identifier. Receiver is the base type name of a method's receiver, and Pointer
// marks a pointer receiver. Line is the 1-based line of the name. Start and End are the byte
// offsets of the span: from the attached Doc group, or from the keyword when no Doc group is
// attached, through the declaration End. Each name of a grouped declaration has the span of
// the complete group.
type Declaration struct {
	Name, Receiver, Kind string
	Pointer              bool
	Line, Start, End     int
}

// Qualified answers the selector spelling of the declaration: the name of a package-level
// declaration, T.M for a value-receiver method, and (*T).M for a pointer-receiver method.
func (d Declaration) Qualified() string {
	switch {
	case d.Receiver == "":
		return d.Name
	case d.Pointer:
		return "(*" + d.Receiver + ")." + d.Name
	}
	return d.Receiver + "." + d.Name
}

// Declarations parses content as one Go file with its comments, and answers each top-level
// declaration except the imports, in source order. A syntax error refuses the whole file,
// because the parser then cannot establish a complete span.
func Declarations(name string, content []byte) ([]Declaration, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, content, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	positions := fset.File(file.Pos())
	var declarations []Declaration
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			item := Declaration{Name: d.Name.Name, Kind: "func", Line: positions.Line(d.Name.Pos())}
			if d.Recv != nil && len(d.Recv.List) > 0 {
				item.Kind = "method"
				item.Receiver, item.Pointer = receiver(d.Recv.List[0].Type)
			}
			item.Start, item.End = span(positions, d.Doc, d)
			declarations = append(declarations, item)
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				continue
			}
			start, end := span(positions, d.Doc, d)
			for _, name := range specNames(d) {
				declarations = append(declarations, Declaration{Name: name.Name, Kind: d.Tok.String(), Line: positions.Line(name.Pos()), Start: start, End: end})
			}
		}
	}
	return declarations, nil
}

// span answers the byte offsets of decl: from its attached Doc group when it has one, and
// otherwise from its keyword, through its End. A comment that the parser does not attach,
// such as one that a blank line separates, stays outside.
func span(positions *token.File, doc *ast.CommentGroup, decl ast.Decl) (int, int) {
	start := decl.Pos()
	if doc != nil {
		start = doc.Pos()
	}
	return positions.Offset(start), positions.Offset(decl.End())
}

// specNames answers the declared names of a type, const, or var declaration in source order.
func specNames(d *ast.GenDecl) []*ast.Ident {
	var names []*ast.Ident
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			names = append(names, s.Name)
		case *ast.ValueSpec:
			names = append(names, s.Names...)
		}
	}
	return names
}

// receiver answers the base type name of a receiver type, and whether the receiver is a
// pointer. The type parameters of a generic receiver are not part of the name.
func receiver(expr ast.Expr) (string, bool) {
	pointer := false
	for {
		switch e := expr.(type) {
		case *ast.StarExpr:
			pointer, expr = true, e.X
		case *ast.ParenExpr:
			expr = e.X
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		case *ast.Ident:
			return e.Name, pointer
		default:
			return "", pointer
		}
	}
}

// symbolRead is one parsed --symbol call. A bare name selects a package-level declaration. A
// receiver-qualified name selects a method: T.M selects M on T or on *T, and (*T).M selects
// only a pointer-receiver method. The cursor fields are the page cursor.
type symbolRead struct {
	text, receiver, name string
	pointer, continued   bool
	offset               int64
	token                string
}

// parseSymbolCall reads the --symbol form from the parsed flags. It reports whether the call
// selects a declaration, or a usage refusal reason: the cursor flags without --symbol,
// --symbol with --full, a malformed cursor, or a selector that is not a Go name.
func parseSymbolCall(flags map[string]string, full bool) (symbolRead, bool, string) {
	text, symbolic := flags[SymbolFlag]
	_, hasOffset := flags[responsebound.OffsetFlag]
	_, hasSource := flags[responsebound.SourceFlag]
	switch {
	case !symbolic && (hasOffset || hasSource):
		return symbolRead{}, false, "the cursor flags take " + SymbolFlag
	case !symbolic:
		return symbolRead{}, false, ""
	case full:
		return symbolRead{}, false, SymbolFlag + " and --full are mutually exclusive"
	}
	read := symbolRead{text: text, name: text}
	var err error
	if read.offset, read.token, read.continued, err = responsebound.ParseCursor(flags); err != nil {
		return symbolRead{}, false, err.Error()
	}
	if receiver, name, qualified := strings.Cut(text, "."); qualified {
		read.name, read.receiver = name, receiver
		if inner, pointer := strings.CutPrefix(receiver, "(*"); pointer {
			read.receiver, read.pointer = strings.TrimSuffix(inner, ")"), true
			if read.receiver == inner {
				read.receiver = ""
			}
		}
		if !token.IsIdentifier(read.receiver) {
			read.name = ""
		}
	}
	if !token.IsIdentifier(read.name) {
		return symbolRead{}, false, SymbolFlag + " takes a Go name, T.Method, or (*T).Method"
	}
	return read, true, ""
}

// matches reports whether d is a declaration that the selector names.
func (s symbolRead) matches(d Declaration) bool {
	return d.Name == s.name && d.Receiver == s.receiver && (!s.pointer || d.Pointer)
}

// candidate is one matching declaration with the file that holds it.
type candidate struct {
	Declaration
	file    string
	content []byte
}

// page selects the declaration among the tracked Go files in scope and answers one bounded
// page of its original source bytes. A file that the read policy refuses or that is not
// valid Go refuses the call, because it can hold the declaration. No match is one
// present-false page. Several matches refuse with each candidate and its narrowing command,
// and with no declaration body. Candidates order by file, source offset, and qualified name.
func (s symbolRead) page(root string, files []string, command []string) (string, int) {
	var found []candidate
	scope := sha256.New()
	var scopeBytes int64
	for _, rel := range files {
		if path.Ext(rel) != ".go" || IsFixture(rel) {
			continue
		}
		content, reason, ok := ReadScannable(filepath.Join(root, filepath.FromSlash(rel)))
		if !ok {
			return toon.Errorf(rel+" is "+reason, "it can hold the declaration; pass a path or a filter that excludes it") + "\n", 1
		}
		fmt.Fprintf(scope, "%d:%s,%x,", len(rel), rel, sha256.Sum256(content))
		scopeBytes += int64(len(content))
		if !bytes.Contains(content, []byte(s.name)) {
			continue
		}
		declarations, err := Declarations(rel, content)
		if err != nil {
			return toon.Errorf(rel+" is not valid Go", err.Error()) + "\n", 1
		}
		for _, d := range declarations {
			if s.matches(d) {
				found = append(found, candidate{d, rel, content})
			}
		}
	}
	slices.SortFunc(found, func(a, b candidate) int {
		return cmp.Or(strings.Compare(a.file, b.file), cmp.Compare(a.Start, b.Start), strings.Compare(a.Qualified(), b.Qualified()))
	})
	request := responsebound.ReadRequest{Command: command, Kind: declarationKind, Selection: s.text, Scope: root, Offset: s.offset, Token: s.token, Continued: s.continued}
	switch len(found) {
	case 0:
		request.Absent, request.SourceSize, request.Payload = true, scopeBytes, bytes.NewReader(nil)
		copy(request.Digest[:], scope.Sum(nil))
	case 1:
		match := found[0]
		selected := match.content[match.Start:match.End]
		request.Selection, request.Path = match.Qualified(), match.file
		request.Digest, request.SourceSize = sha256.Sum256(match.content), int64(len(match.content))
		request.Payload, request.Size = bytes.NewReader(selected), int64(len(selected))
	default:
		return ambiguous(root, s.text, found)
	}
	page, err := responsebound.PageRead(request)
	if err != nil {
		return err.Error() + "\n", 1
	}
	return page.Output, 0
}

// ambiguous refuses a selector that names several declarations. It lists each candidate and
// offers one command per candidate that names its file from the caller's directory and its
// qualified name.
func ambiguous(root, text string, found []candidate) (string, int) {
	cwd, _ := os.Getwd()
	rows := make([][]any, 0, len(found))
	var actions []axi.Action
	for _, c := range found {
		rows = append(rows, []any{c.file, c.Line, c.Qualified(), c.Kind})
		operand, err := filepath.Rel(cwd, filepath.Join(root, filepath.FromSlash(c.file)))
		if err != nil {
			operand = c.file
		}
		var arguments []axi.InvocationArgument
		for _, argument := range commandArgv(operand, SymbolFlag, c.Qualified()) {
			arguments = append(arguments, axi.KnownArgument(argument))
		}
		actions = append(actions, axi.ExecutableInvocation("read this declaration", arguments...))
	}
	table, err := toon.TableTyped("candidates", []string{"file", "line", "name", "kind"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return toon.Errorf(fmt.Sprintf("%s %s matches %d declarations", SymbolFlag, text, len(found)), "read one candidate through its command") + "\n" + table + help, 1
}

// symbolCommand is the caller's own --symbol call after `bench`, without the cursor flags.
// Each continuation repeats it, so it names the same scope from the caller's directory.
func symbolCommand(positionals []string, flags map[string]string) []string {
	command := append(commandArgv(positionals...), SymbolFlag, flags[SymbolFlag])
	for _, filter := range []string{"--production", "--test"} {
		if _, ok := flags[filter]; ok {
			command = append(command, filter)
		}
	}
	return command
}

// commandArgv answers the outline argv after `bench`, with extra after the command name.
func commandArgv(extra ...string) []string {
	return append(strings.Fields(grammar.Cmd)[1:], extra...)
}
