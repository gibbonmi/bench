package ownership

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"slices"

	"github.com/gibbonmi/bench/internal/consumers"
)

// errorConstructors names, by package path, the one standard library function of each
// package that constructs an error from its message argument.
var errorConstructors = map[string]string{"errors": "New", "fmt": "Errorf"}

// sentinel grades one sentinel premise, whose owner is a package-level error var.
func (c *census) sentinel(item PremiseItem) (string, error) { return c.rule(item, sentinelText) }

// errorText grades one error premise, whose owner is an error type.
func (c *census) errorText(item PremiseItem) (string, error) { return c.rule(item, errorMethodText) }

// rule grades one sentinel or error premise: the exact owner, the static construction
// text that construct reads from the owner's declaration, and each named producer's static
// use of the owner. Exact construction proves no path by which the failure reaches its
// intended caller; only a reachable producer witness proves that.
func (c *census) rule(item PremiseItem, construct func(Owner, declared) (string, error)) (string, error) {
	object, err := fields(item.Subject, "subject: ", "owner", "producers")
	if err != nil {
		return "", err
	}
	owner, err := identity(object["owner"], "subject: owner: ")
	if err != nil {
		return "", err
	}
	producers, err := identities(object["producers"], "subject: producers: ")
	if err != nil {
		return "", err
	}
	var expected struct {
		Value *string `json:"value"`
	}
	if _, err := typed(item.Expected, "expected: ", &expected, "value"); err != nil {
		return "", err
	}
	switch {
	case expected.Value == nil:
		return "", errors.New("expected: value is not a string")
	case len(producers) == 0:
		return "", fmt.Errorf("subject: %s names no producer", owner.Symbol)
	}
	decl, err := c.exact(owner)
	if err != nil {
		return "", err
	}
	text, err := construct(owner, decl)
	if err != nil {
		return "", err
	}
	if text != *expected.Value {
		return "", fmt.Errorf("%s constructs %s, not %s", owner.Symbol, text, *expected.Value)
	}
	uses := consumers.Rows(c.pkgs, decl.obj, c.root)
	for _, producer := range producers {
		if _, err := c.exact(producer); err != nil {
			return "", err
		}
		if !slices.ContainsFunc(uses, func(row consumers.Row) bool { return row.File == producer.Path && row.Enclosing == localName(producer) }) {
			return "", fmt.Errorf("%s does not statically use %s", producer.Symbol, owner.Symbol)
		}
	}
	return fmt.Sprintf("%s constructs %s, and its %d named producer(s) use it statically at the source tip; no witness proves that the failure reaches its caller", owner.Symbol, text, len(producers)), nil
}

// sentinelText is the construction text of a sentinel: the one initializer of a
// package-level var, which must call an errorConstructors function with one constant
// string. Any other initializer, such as a call with a runtime argument, is no static
// construction.
func sentinelText(owner Owner, d declared) (string, error) {
	spec, ok := d.node.(*ast.ValueSpec)
	if _, isVar := d.obj.(*types.Var); !ok || !isVar || len(spec.Values) != len(spec.Names) {
		return "", fmt.Errorf("%s is no package-level var with its own initializer", owner.Symbol)
	}
	value := spec.Values[slices.IndexFunc(spec.Names, func(id *ast.Ident) bool { return id.Name == d.obj.Name() })]
	static := false
	if call, ok := value.(*ast.CallExpr); ok && len(call.Args) == 1 && constantString(d.pkg, call.Args[0]) {
		if selector, isSelector := call.Fun.(*ast.SelectorExpr); isSelector {
			callee, isFunc := d.pkg.Info.Uses[selector.Sel].(*types.Func)
			static = isFunc && callee.Pkg() != nil && errorConstructors[callee.Pkg().Path()] == callee.Name()
		}
	}
	if !static {
		return "", fmt.Errorf("the initializer of %s is no static construction", owner.Symbol)
	}
	return render(d.pkg, value)
}

// errorMethodText is the construction text of an error type: the one constant string that
// its own Error method returns. A method promoted from another type, or a body that returns
// anything but one constant, such as a field that the constructor sets, is dynamic.
func errorMethodText(owner Owner, d declared) (string, error) {
	named, ok := d.obj.(*types.TypeName)
	if !ok {
		return "", fmt.Errorf("%s is no type, so it declares no Error method", owner.Symbol)
	}
	errorType := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	if !types.Implements(named.Type(), errorType) && !types.Implements(types.NewPointer(named.Type()), errorType) {
		return "", fmt.Errorf("%s does not implement error", owner.Symbol)
	}
	method, _, _ := types.LookupFieldOrMethod(named.Type(), true, named.Pkg(), "Error")
	var body *ast.BlockStmt
	if method != nil {
		if fn, isFunc := declarationNode([]*consumers.Package{d.pkg}, d.pkg.PkgPath, method, named.Name()+".Error").node.(*ast.FuncDecl); isFunc {
			body = fn.Body
		}
	}
	if body != nil && len(body.List) == 1 {
		if ret, isReturn := body.List[0].(*ast.ReturnStmt); isReturn && len(ret.Results) == 1 && constantString(d.pkg, ret.Results[0]) {
			return render(d.pkg, ret.Results[0])
		}
	}
	return "", fmt.Errorf("the Error method of %s returns no constant, so its text is dynamic", owner.Symbol)
}

// constantString reports whether the type checker gives expr one constant string value.
func constantString(pkg *consumers.Package, expr ast.Expr) bool {
	value, ok := pkg.Info.Types[expr]
	return ok && value.Value != nil && value.Value.Kind() == constant.String
}

// bound grades one bound premise: its owner must be a constant, whose exact Go value is the
// static bound. A var, a function, or any other declaration whose value a runtime input can
// change stays unproved, because no inference turns it into a limit.
func (c *census) bound(item PremiseItem) (string, error) {
	owner, err := identity(item.Subject, "subject: ")
	if err != nil {
		return "", err
	}
	var expected struct {
		Value *string `json:"value"`
	}
	if _, err := typed(item.Expected, "expected: ", &expected, "value"); err != nil {
		return "", err
	}
	if expected.Value == nil {
		return "", errors.New("expected: value is not a string")
	}
	decl, err := c.exact(owner)
	if err != nil {
		return "", err
	}
	value, ok := decl.obj.(*types.Const)
	if !ok {
		return "", fmt.Errorf("%s is a %s, so its value is no static bound", owner.Symbol, owner.Kind)
	}
	at := fmt.Sprintf("%s holds the constant value %s", owner.Symbol, value.Val().ExactString())
	if value.Val().ExactString() != *expected.Value {
		return "", fmt.Errorf("%s, not %s", at, *expected.Value)
	}
	return at + " at the source tip", nil
}
