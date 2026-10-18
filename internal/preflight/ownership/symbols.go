package ownership

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"path"
	"reflect"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/consumers"
)

// symbolPremises is every symbol and rule premise kind and its collector, which keeps the
// sizingPremises contract.
var symbolPremises = map[string]func(*census, PremiseItem) (string, error){
	"symbol":      (*census).symbol,
	"import":      (*census).importEdge,
	"anchor-kind": (*census).anchorKind,
	"sentinel":    (*census).sentinel,
	"error":       (*census).errorText,
	"bound":       (*census).bound,
}

// declared is one owner identity's exact declaration at the source tip: the object, the
// loaded package whose syntax declares it, and the top-level node that declares it. node is
// nil when no top-level node declares the object, as for a field.
type declared struct {
	obj  types.Object
	pkg  *consumers.Package
	node ast.Node
}

// symbol grades one symbol premise. Its subject is the owner identity, and its expected
// kind and signature repeat that identity's values.
func (c *census) symbol(item PremiseItem) (string, error) {
	owner, err := ownerSubject(item)
	if err != nil {
		return "", err
	}
	var expected struct {
		Kind      *string `json:"kind"`
		Signature *string `json:"signature"`
	}
	if _, err := typed(item.Expected, "expected: ", &expected, "kind", "signature"); err != nil {
		return "", err
	}
	if expected.Kind == nil || expected.Signature == nil || *expected.Kind != owner.Kind || *expected.Signature != owner.Signature {
		return "", errors.New("expected: the kind and signature differ from the subject identity")
	}
	if _, err := c.exact(owner); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s resolves in %s as the %s that its signature names at the source tip", owner.Symbol, owner.Path, owner.Kind), nil
}

// ownerSubject decodes the subject of one symbol or bound premise, which is one owner
// identity.
func ownerSubject(item PremiseItem) (Owner, error) {
	return identity(item.Subject, "subject: ")
}

// exact resolves one owner identity at the source tip and refuses each part that differs:
// the one declaration that its symbol names in its exact package, the file that declares
// it, its consumers declaration kind, and its declaration syntax as the Go formatter
// prints it. A missing or ambiguous declaration names no owner.
func (c *census) exact(owner Owner) (declared, error) {
	pkgs, err := c.packages()
	if err != nil {
		return declared{}, err
	}
	name := localName(owner)
	matches, _ := consumers.Resolve(pkgs, owner.Symbol)
	matches = slices.DeleteFunc(matches, func(m consumers.Match) bool { return m.PkgPath != owner.Package || m.Name != name })
	switch {
	case len(matches) == 0:
		return declared{}, fmt.Errorf("%s declares no %s at the source tip", owner.Package, name)
	case len(matches) > 1:
		return declared{}, fmt.Errorf("%s names %d declarations of %s at the source tip", owner.Symbol, len(matches), owner.Package)
	}
	match := matches[0]
	file, _, err := c.position(pkgs, match.Obj)
	switch {
	case err != nil:
		return declared{}, err
	case file != owner.Path:
		return declared{}, fmt.Errorf("%s is declared in %s, not %s", owner.Symbol, file, owner.Path)
	case match.Kind != owner.Kind:
		return declared{}, fmt.Errorf("%s is a %s, not a %s", owner.Symbol, match.Kind, owner.Kind)
	}
	decl := declarationNode(pkgs, match.PkgPath, match.Obj, name)
	signature, err := decl.signature()
	if err != nil {
		return declared{}, err
	}
	if signature != owner.Signature {
		return declared{}, fmt.Errorf("%s declares %q, not %q", owner.Symbol, signature, owner.Signature)
	}
	return decl, nil
}

// declarationNode finds the top-level node of a loaded package that declares obj under
// name, by the consumers declaration-name forms. Its package is the one whose syntax holds
// that node, so the type facts of that package describe the node.
func declarationNode(pkgs []*consumers.Package, pkgPath string, obj types.Object, name string) declared {
	at := obj.Pos()
	holds := func(node ast.Node) bool {
		names, ok := consumers.DeclNames(node)
		return ok && node.Pos() <= at && at < node.End() && slices.Contains(names, name)
	}
	for _, pkg := range pkgs {
		if pkg.PkgPath != pkgPath {
			continue
		}
		for _, file := range pkg.Files {
			if at < file.Pos() || at >= file.End() {
				continue
			}
			for _, decl := range file.Decls {
				nodes := []ast.Node{decl}
				if gen, ok := decl.(*ast.GenDecl); ok {
					nodes = nodes[:0]
					for _, spec := range gen.Specs {
						nodes = append(nodes, spec)
					}
				}
				for _, node := range nodes {
					if holds(node) {
						return declared{obj: obj, pkg: pkg, node: node}
					}
				}
			}
		}
	}
	return declared{obj: obj}
}

// signature is the normalized declaration syntax: the Go formatter's print of the
// declaring node without its comments or a function body. A constant's signature is empty,
// because its value belongs to a bound premise. A declaration with no top-level node, such
// as a field, has no supported normalization.
func (d declared) signature() (string, error) {
	node, prefix := d.node, ""
	switch n := d.node.(type) {
	case *ast.FuncDecl:
		trimmed := *n
		trimmed.Body = nil
		node = &trimmed
	case *ast.TypeSpec:
		prefix = "type "
	case *ast.ValueSpec:
		if _, isConst := d.obj.(*types.Const); isConst {
			return "", nil
		}
		prefix = "var "
	default:
		return "", fmt.Errorf("%s has no top-level declaration whose syntax normalizes", d.obj.Name())
	}
	text, err := render(d.pkg, node)
	return prefix + text, err
}

// render is the Go formatter's print of one node of pkg without its comments. The printer
// prints each comment group that a node holds, such as a field's doc, so render detaches
// every group before the print. The loaded syntax is shared, so each group returns to its
// node before render returns.
func render(pkg *consumers.Package, node ast.Node) (string, error) {
	var restore []func()
	defer func() {
		for _, reattach := range restore {
			reattach()
		}
	}()
	detach := func(groups ...**ast.CommentGroup) {
		for _, group := range groups {
			held := *group
			*group = nil
			restore = append(restore, func() { *group = held })
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Field:
			detach(&n.Doc, &n.Comment)
		case *ast.ValueSpec:
			detach(&n.Doc, &n.Comment)
		case *ast.TypeSpec:
			detach(&n.Doc, &n.Comment)
		case *ast.ImportSpec:
			detach(&n.Doc, &n.Comment)
		case *ast.GenDecl:
			detach(&n.Doc)
		case *ast.FuncDecl:
			detach(&n.Doc)
		}
		return true
	})
	var buf bytes.Buffer
	if err := format.Node(&buf, pkg.Fset, node); err != nil {
		return "", fmt.Errorf("the declaration syntax does not format: %w", err)
	}
	return buf.String(), nil
}

// importEdge grades one import premise from its calling package: the target must declare
// the owner, the owner must be exported to another package, the Go internal-package rule
// must admit the caller, and the target must not already reach the caller through the
// loaded import graph, tests included, because the edge would then close a cycle.
func (c *census) importEdge(item PremiseItem) (string, error) {
	subject, err := decodeImportSubject(item)
	if err != nil {
		return "", err
	}
	var expected struct {
		Reachable bool `json:"reachable"`
	}
	if _, err := typed(item.Expected, "expected: ", &expected, "reachable"); err != nil {
		return "", err
	}
	caller, target, owner := subject.Caller, subject.Target, subject.Owner
	switch {
	case !expected.Reachable:
		return "", errors.New("expected: reachable is not true")
	case owner.Package != target:
		return "", fmt.Errorf("subject: the owner package %s is not the target %s", owner.Package, target)
	}
	if _, err := c.exact(owner); err != nil {
		return "", err
	}
	parent, internal := internalParent(target)
	switch {
	case !slices.ContainsFunc(c.pkgs, func(pkg *consumers.Package) bool { return pkg.PkgPath == caller }):
		return "", fmt.Errorf("the calling package %s does not load at the source tip", caller)
	case caller == target:
		return "", fmt.Errorf("%s cannot import its own package, which is a cycle", caller)
	case !exported(localName(owner)):
		return "", fmt.Errorf("%s is unexported, so %s cannot reach it from another package", owner.Symbol, caller)
	case internal && caller != parent && !strings.HasPrefix(caller, parent+"/"):
		return "", fmt.Errorf("%s sits outside the tree of %s, so it cannot import the internal package %s", caller, parent, target)
	case reaches(c.pkgs, target, caller):
		return "", fmt.Errorf("%s already imports %s through its own imports, so an edge from %s closes a cycle", target, caller, caller)
	}
	return fmt.Sprintf("%s reaches %s through an import of %s with no cycle at the source tip", caller, owner.Symbol, target), nil
}

// importSubject is the decoded subject of one import premise. Caller and Target are package
// import paths, and Owner is the declaration identity that the target declares.
type importSubject struct {
	Caller string `json:"caller"`
	Target string `json:"target"`
	Owner  Owner  `json:"-"`
}

// decodeImportSubject decodes the subject of one import premise.
func decodeImportSubject(item PremiseItem) (importSubject, error) {
	var subject importSubject
	object, err := typed(item.Subject, "subject: ", &subject, "caller", "target", "owner")
	if err == nil {
		subject.Owner, err = identity(object["owner"], "subject: owner: ")
	}
	return subject, err
}

// exported reports whether each part of a declaration spelling is exported.
func exported(name string) bool {
	for _, part := range strings.Split(name, ".") {
		if !token.IsExported(part) {
			return false
		}
	}
	return true
}

// internalParent is the parent of the last internal element of an import path. By the Go
// internal-package rule, only a package rooted at that parent can import the path.
func internalParent(importPath string) (string, bool) {
	parts := strings.Split(importPath, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "internal" {
			return strings.Join(parts[:i], "/"), true
		}
	}
	return "", false
}

// reaches reports whether from imports to, directly or through the imports of each loaded
// form of each package on the way, test forms included.
func reaches(pkgs []*consumers.Package, from, to string) bool {
	seen := map[string]bool{from: true}
	next := []string{from}
	for len(next) > 0 {
		at := next[0]
		next = next[1:]
		for _, pkg := range pkgs {
			if pkg.PkgPath != at {
				continue
			}
			for _, imported := range pkg.Types.Imports() {
				if imported.Path() == to {
					return true
				}
				if !seen[imported.Path()] {
					seen[imported.Path()] = true
					next = append(next, imported.Path())
				}
			}
		}
	}
	return false
}

// anchorKind grades one anchor kind against the anchors kind inventory at the source tip:
// the constants of the compiled anchors Kind type, read through the typed loader in the
// package that declares that type. The owner location must be the declaration of the kind
// constant, so a familiar spelling elsewhere or another constant of the package names no
// supported kind.
func (c *census) anchorKind(item PremiseItem) (string, error) {
	var expected struct {
		Supported bool `json:"supported"`
	}
	subject, err := decodeAnchorSubject(item)
	if err == nil {
		_, err = typed(item.Expected, "expected: ", &expected, "supported")
	}
	switch {
	case err != nil:
	case !expected.Supported:
		err = errors.New("expected: supported is not true")
	case subject.Owner.Source != "tip":
		err = errors.New("subject: the anchor kind inventory is read only at the source tip")
	default:
		_, err = c.locate(subject.Owner)
	}
	var pkgs []*consumers.Package
	if err == nil {
		pkgs, err = c.packages()
	}
	if err != nil {
		return "", err
	}
	inventory := reflect.TypeFor[anchors.Kind]()
	var kind types.Object
	matches, _ := consumers.Resolve(pkgs, path.Base(inventory.PkgPath())+"."+subject.Kind)
	for _, match := range matches {
		named, ok := match.Obj.Type().(*types.Named)
		if _, isConst := match.Obj.(*types.Const); isConst && ok && match.PkgPath == inventory.PkgPath() && match.Name == subject.Kind &&
			named.Obj().Pkg() != nil && named.Obj().Pkg().Path() == inventory.PkgPath() && named.Obj().Name() == inventory.Name() {
			kind = match.Obj
		}
	}
	if kind == nil {
		return "", fmt.Errorf("%q is no constant of the anchor kind inventory %s.%s at the source tip", subject.Kind, inventory.PkgPath(), inventory.Name())
	}
	file, line, err := c.position(pkgs, kind)
	switch {
	case err != nil:
		return "", err
	case file != subject.Owner.Path || line != subject.Owner.Line:
		return "", fmt.Errorf("the anchor kind inventory declares %s at %s line %d, not line %d of %s", subject.Kind, file, line, subject.Owner.Line, subject.Owner.Path)
	}
	return fmt.Sprintf("the anchor kind inventory declares the anchor kind %s at %s line %d at the source tip", subject.Kind, file, line), nil
}

// anchorSubject is the decoded subject of one anchor-kind premise: the owner location that
// declares the kind constant, and the kind name.
type anchorSubject struct {
	Owner Location `json:"owner"`
	Kind  string   `json:"kind"`
}

// decodeAnchorSubject decodes the subject of one anchor-kind premise.
func decodeAnchorSubject(item PremiseItem) (anchorSubject, error) {
	var subject anchorSubject
	object, err := typed(item.Subject, "subject: ", &subject, "owner", "kind")
	if err == nil {
		err = location(object, "owner", "subject: ", true)
	}
	return subject, err
}
