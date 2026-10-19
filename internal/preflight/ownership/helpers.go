package ownership

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
)

// census reads the static references of declared owners. Its tip side is the typed
// packages of the checkout, which must be clean and at the source tip, loaded once on first
// use; its base side reads the source base without types. A snapshot keeps only the
// resolved identities that the census reads, never the packages themselves. tracked holds
// the tracked paths of each commit that a sizing premise read. symbols holds the identity
// of each symbol premise that the census graded.
type census struct {
	root    string
	source  Source
	loaded  bool
	pkgs    []*consumers.Package
	err     error
	tracked map[string][]string
	symbols []SymbolIdentity
}

// packages is the typed tip side, which reads the checkout at the source tip.
func (c *census) packages() ([]*consumers.Package, error) {
	if !c.loaded {
		c.loaded = true
		if c.err = c.atTip(); c.err == nil {
			c.pkgs, c.err = consumers.Load(c.root)
		}
	}
	return c.pkgs, c.err
}

// atTip refuses a checkout that is not the source tip, because a read of it would enumerate
// bytes that the source pair does not hold.
func (c *census) atTip() error {
	head, err := git.Output("-C", c.root, "rev-parse", "HEAD")
	switch {
	case err != nil:
		return fmt.Errorf("the checkout HEAD does not resolve: %w", err)
	case head != c.source.Tip:
		return fmt.Errorf("the checkout HEAD %s is not the source tip %s", head, c.source.Tip)
	}
	return nil
}

// declaration resolves one owner identity at the source tip: the declaration that its
// symbol names in its exact package, and the repository path that declares it. An absent
// declaration answers a nil object and no error.
func (c *census) declaration(owner Owner) (types.Object, string, error) {
	pkgs, err := c.packages()
	if err != nil {
		return nil, "", err
	}
	matches, err := consumers.Resolve(pkgs, owner.Symbol)
	if err != nil {
		return nil, "", nil
	}
	for _, match := range matches {
		if match.PkgPath == owner.Package {
			file, _, err := c.position(pkgs, match.Obj)
			return match.Obj, file, err
		}
	}
	return nil, "", nil
}

// position is the repository path and line that declare one loaded object.
func (c *census) position(pkgs []*consumers.Package, obj types.Object) (string, int, error) {
	position := pkgs[0].Fset.Position(obj.Pos())
	rel, err := c.relative(position.Filename)
	if err != nil {
		return "", 0, err
	}
	return rel, position.Line, nil
}

// relative is the repository path of one loaded file name.
func (c *census) relative(filename string) (string, error) {
	rel, err := filepath.Rel(c.root, filename)
	return filepath.ToSlash(rel), err
}

// localName is an owner's declaration spelling inside its package, by the consumers query
// grammar. Owner decoding refuses a symbol with no such spelling, so that read cannot fail
// here.
func localName(owner Owner) string {
	name, _ := consumers.LocalName(owner.Package, owner.Symbol)
	return name
}

// premises grades each current premise of a kind this census owns: a reused owner, a reused
// test, the retirement, posture, and printed-field member censuses, the sizing premises, and
// the symbol and rule premises. It answers the graded row of each such premise by ID. A
// fault leaves its premise unknown with that fault as the reason.
func (c *census) premises(items []PremiseItem) map[string]Premise {
	graded := map[string]Premise{}
	for _, item := range items {
		var reason string
		var err error
		shape, member := memberCensuses[item.Kind]
		collect, collected := sizingPremises[item.Kind]
		if !collected {
			collect, collected = symbolPremises[item.Kind]
		}
		switch {
		case item.State != "current":
			continue
		case item.Kind == "reused-owner":
			reason, err = c.reusedOwner(item)
		case item.Kind == "reused-test":
			reason, err = c.reusedTest(item)
		case member:
			reason, err = c.memberCensus(item, shape)
		case collected:
			reason, err = collect(c, item)
		default:
			continue
		}
		row := Premise{ID: item.ID, State: StateGreen, Source: c.source.Tip, Reason: reason}
		if err != nil {
			row.State, row.Reason = StateUnknown, item.Kind+" premise: "+err.Error()
		}
		graded[item.ID] = row
	}
	return graded
}

// reusedOwner verifies one reused owner location at its pinned source and answers the
// reason that locates its rule source.
func (c *census) reusedOwner(item PremiseItem) (string, error) {
	var expected struct {
		Required bool `json:"required"`
	}
	loc, err := reusedOwnerSubject(item)
	if err == nil {
		_, err = typed(item.Expected, "expected: ", &expected, "required")
	}
	if err == nil && !expected.Required {
		err = fmt.Errorf("expected: required is not true")
	}
	if err == nil {
		_, err = c.locate(loc)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("reused owner %s line %d at the source %s holds its needle", loc.Path, loc.Line, loc.Source), nil
}

// reusedOwnerSubject decodes the subject of one reused-owner premise: its owner location.
func reusedOwnerSubject(item PremiseItem) (Location, error) {
	var subject struct {
		Owner Location `json:"owner"`
	}
	object, err := typed(item.Subject, "subject: ", &subject, "owner")
	if err == nil {
		err = location(object, "owner", "subject: ", true)
	}
	return subject.Owner, err
}

// locate verifies that one owner location's line holds its needle exactly once at its
// pinned source, and answers that source's commit. A moved or repeated needle is a fault,
// never a nearby match.
func (c *census) locate(loc Location) (string, error) {
	rev, err := c.rev(loc.Source, loc.Path)
	if err != nil {
		return "", err
	}
	data, err := git.ReadTreeFile(c.root, rev, loc.Path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	if loc.Line > len(lines) || strings.Count(strings.TrimSuffix(lines[loc.Line-1], "\r"), loc.Needle) != 1 {
		return "", fmt.Errorf("%s line %d does not hold its needle exactly once at the source %s", loc.Path, loc.Line, loc.Source)
	}
	return rev, nil
}

// rev is the commit that a declared source of subject names: the explicit base or tip.
func (c *census) rev(source, subject string) (string, error) {
	if err := enum("", "source", source, "base", "tip"); err != nil {
		return "", err
	}
	rev := c.source.Tip
	if source == "base" {
		rev = c.source.Base
	}
	if rev == "" {
		return "", fmt.Errorf("no source %s pins %s", source, subject)
	}
	return rev, nil
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

// topNode is one top-level node that a declaration name can name: a function declaration,
// or one spec of a general declaration with that declaration's token.
type topNode struct {
	node ast.Node
	tok  token.Token
}

// topLevel is each top-level node of file, in source order.
func topLevel(file *ast.File) []topNode {
	var nodes []topNode
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			nodes = append(nodes, topNode{decl, token.FUNC})
			continue
		}
		for _, spec := range gen.Specs {
			nodes = append(nodes, topNode{spec, gen.Tok})
		}
	}
	return nodes
}
