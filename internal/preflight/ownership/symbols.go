package ownership

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
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
// loaded package whose syntax declares it, the top-level node that declares it, and that
// node's declared spelling. node is nil when no top-level node declares the object, as for
// a field.
type declared struct {
	obj  types.Object
	pkg  *consumers.Package
	node ast.Node
	name string
}

// symbol grades one symbol premise. Its subject is the owner identity, and its expected
// kind and signature repeat that identity's values. The census keeps the premise's tip
// identity, and a fault that names no discrepancy class leaves it unread.
func (c *census) symbol(item PremiseItem) (reason string, err error) {
	identity := SymbolIdentity{Premise: item.ID}
	defer func() {
		if err != nil {
			identity.Reason = err.Error()
			identity.Discrepancy = cmp.Or(identity.Discrepancy, DiscrepancyUnread)
		}
		c.symbols = append(c.symbols, identity)
	}()
	if identity.Owner, err = ownerSubject(item); err != nil {
		return "", err
	}
	owner := identity.Owner
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
	if _, identity.Current, identity.Discrepancy, err = c.identify(owner); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s resolves in %s as the %s that its signature names at the source tip", owner.Symbol, owner.Path, owner.Kind), nil
}

// ownerSubject decodes the subject of one symbol or bound premise, which is one owner
// identity.
func ownerSubject(item PremiseItem) (Owner, error) {
	return identity(item.Subject, "subject: ")
}

// exact resolves one owner identity at the source tip and refuses the first part that
// differs, as identify classifies it.
func (c *census) exact(owner Owner) (declared, error) {
	decl, _, _, err := c.identify(owner)
	return decl, err
}

// identify resolves one owner identity at the source tip: the one declaration that its
// symbol names in its exact package, through the consumers origin of an alias, and the
// file, consumers kind, and normalized syntax that declare it. It answers the current
// declaration and, as err, the first discrepancy with its class. A bare spelling that names
// several declarations is ambiguous, because a package filter that selects one proves no
// cited owner.
func (c *census) identify(owner Owner) (declared, Declaration, Discrepancy, error) {
	pkgs, err := c.packages()
	if err != nil {
		return declared{}, Declaration{}, DiscrepancyUnread, err
	}
	name := localName(owner)
	all, _ := consumers.Resolve(pkgs, owner.Symbol)
	matches := slices.DeleteFunc(slices.Clone(all), func(m consumers.Match) bool { return m.PkgPath != owner.Package || m.Name != name })
	switch {
	case name == owner.Symbol && len(all) > 1:
		return declared{}, Declaration{}, DiscrepancyAmbiguous, fmt.Errorf("the bare name %s names %d declarations at the source tip", owner.Symbol, len(all))
	case len(matches) > 1:
		return declared{}, Declaration{}, DiscrepancyAmbiguous, fmt.Errorf("%s names %d declarations of %s at the source tip", owner.Symbol, len(matches), owner.Package)
	case len(matches) == 0:
		class, reason := owner.Mismatch(Declaration{}, "source tip")
		if context, why := c.context(pkgs, owner.Path, name); context != "" {
			class, reason = context, reason+"; "+why
		}
		return declared{}, Declaration{}, class, errors.New(reason)
	}
	match := matches[0]
	file, _, err := c.position(pkgs, match.Obj)
	if err != nil {
		return declared{}, Declaration{}, DiscrepancyUnread, err
	}
	pkgPath := match.PkgPath
	if match.Obj.Pkg() != nil {
		pkgPath = match.Obj.Pkg().Path()
	}
	decl := declarationNode(pkgs, pkgPath, match.Obj, name)
	signature, signatureErr := decl.signature()
	current := Declaration{Path: file, Name: decl.name, Kind: match.Kind, Signature: signature}
	switch class, reason := owner.Mismatch(current, "source tip"); {
	case signatureErr != nil && (class == "" || class == DiscrepancySignature):
		return declared{}, current, DiscrepancyUnread, signatureErr
	case class != "":
		return declared{}, current, class, errors.New(reason)
	}
	return decl, current, "", nil
}

// context classifies an owner path that no loaded package holds but the committed tip
// holds: a Go source that declares name there sits outside the default build context, and a
// path that holds no Go source is unread. Any other path answers no class.
func (c *census) context(pkgs []*consumers.Package, path, name string) (Discrepancy, string) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if loaded, err := c.relative(pkg.Fset.Position(file.Pos()).Filename); err == nil && loaded == path {
				return "", ""
			}
		}
	}
	src, err := git.ReadTreeFile(c.root, c.source.Tip, path)
	if err != nil {
		return "", ""
	}
	_, found, err := SourceDeclaration(path, src, name)
	switch {
	case err != nil:
		return DiscrepancyUnread, path + " holds no Go source: " + err.Error()
	case found:
		return DiscrepancyContext, path + " declares it only outside the default build context"
	}
	return "", ""
}

// declarationNode finds the top-level node of a loaded package that declares obj, by the
// consumers declaration-name forms: the node that spans obj's position and declares name or
// obj's own name, so an alias spelling reaches its origin's declaration. Its package is the
// one whose syntax holds that node, so the type facts of that package describe the node.
func declarationNode(pkgs []*consumers.Package, pkgPath string, obj types.Object, name string) declared {
	at := obj.Pos()
	spelling := func(node ast.Node) string {
		names, _ := consumers.DeclNames(node)
		if node.Pos() > at || at >= node.End() {
			return ""
		}
		for _, declared := range names {
			if declared == name || declared == obj.Name() || strings.HasSuffix(declared, "."+obj.Name()) {
				return declared
			}
		}
		return ""
	}
	for _, pkg := range pkgs {
		if pkg.PkgPath != pkgPath {
			continue
		}
		for _, file := range pkg.Files {
			if at < file.Pos() || at >= file.End() {
				continue
			}
			for _, top := range topLevel(file) {
				if name := spelling(top.node); name != "" {
					return declared{obj: obj, pkg: pkg, node: top.node, name: name}
				}
			}
		}
	}
	return declared{obj: obj}
}

// signature is the normalized declaration syntax of the declaring node.
func (d declared) signature() (string, error) {
	_, isConst := d.obj.(*types.Const)
	var fset *token.FileSet
	if d.pkg != nil {
		fset = d.pkg.Fset
	}
	return normalize(fset, d.node, isConst, d.obj.Name())
}

// normalize is the normalized syntax of one top-level declaring node: the Go formatter's
// print of the node without its comments or a function body. A constant's syntax is empty,
// because its value belongs to a bound premise. Any other node, such as a field's, has no
// supported normalization; name names the declaration in that fault.
func normalize(fset *token.FileSet, node ast.Node, isConst bool, name string) (string, error) {
	prefix := ""
	switch n := node.(type) {
	case *ast.FuncDecl:
		trimmed := *n
		trimmed.Body = nil
		node = &trimmed
	case *ast.TypeSpec:
		prefix = "type "
	case *ast.ValueSpec:
		if isConst {
			return "", nil
		}
		prefix = "var "
	default:
		return "", fmt.Errorf("%s has no top-level declaration whose syntax normalizes", name)
	}
	text, err := printed(fset, node)
	return prefix + text, err
}

// render is the Go formatter's print of one node of pkg without its comments.
func render(pkg *consumers.Package, node ast.Node) (string, error) {
	return printed(pkg.Fset, node)
}

// printed is the Go formatter's print of one node of fset without its comments. The printer
// prints each comment group that a node holds, such as a field's doc, so printed detaches
// every group before the print. The loaded syntax is shared, so each group returns to its
// node before printed returns.
func printed(fset *token.FileSet, node ast.Node) (string, error) {
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
	if err := format.Node(&buf, fset, node); err != nil {
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
