package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// processLifetimePrefilter is the substring prefilter: a file that spells none of the
// private lifetime mechanics, none of the owner calls, and none of the vector readers carries none of them.
var processLifetimePrefilter = slices.Concat([]string{"Setpgid", "Kill", "Signal", "Wait"}, processLifetimeOwnerCalls, slices.Collect(maps.Keys(processLifetimeVectorReaders)))

// processLifetimeProtectionImport is the import path of the protection leaf of the owner.
const processLifetimeProtectionImport = processLifetimeOwnerImport + "/protection"

// processLifetimeVectorReaders maps each protection reader to the argument test that makes the
// call read the process vector that RunInherited reads. A read of another value is not counted.
var processLifetimeVectorReaders = map[string]func(file *ast.File, args []ast.Expr) bool{
	"Inherited": func(_ *ast.File, args []ast.Expr) bool { return len(args) == 0 },
	"InheritedFrom": func(file *ast.File, args []ast.Expr) bool {
		call, ok := processLifetimeOnlyCall(args)
		return ok && len(call.Args) == 0 && guardProjectionImported(file, call, "os") == "Environ"
	},
	"ParseVector": func(file *ast.File, args []ast.Expr) bool {
		call, ok := processLifetimeOnlyCall(args)
		return ok && len(call.Args) == 1 && guardProjectionImported(file, call, "os") == "Getenv" &&
			guardProjectionImported(file, &ast.CallExpr{Fun: call.Args[0]}, processLifetimeProtectionImport) == "Env"
	},
}

// processLifetimeOnlyCall gives the call that is the one argument of args.
func processLifetimeOnlyCall(args []ast.Expr) (*ast.CallExpr, bool) {
	if len(args) != 1 {
		return nil, false
	}
	call, ok := ast.Unparen(args[0]).(*ast.CallExpr)
	return call, ok
}

// processLifetimeBuildsRegistration reports whether decl of the file at file, or a same-package
// function that it calls, builds a subprocess.Registration literal.
func processLifetimeBuildsRegistration(fset *token.FileSet, file string, decl *ast.FuncDecl) (bool, error) {
	functions, err := processLifetimePackageFunctions(filepath.Dir(file), ".")
	if err != nil {
		return false, err
	}
	candidates := functions.named(decl.Name.Name)
	start := slices.IndexFunc(candidates, func(candidate *packageFunc) bool {
		return candidate.rel == filepath.Base(file) && candidate.line == fset.Position(decl.Pos()).Line
	})
	if start < 0 {
		return false, nil
	}
	builds := func(current *packageFunc) bool {
		found := false
		ast.Inspect(current.body, func(node ast.Node) bool {
			if literal, ok := node.(*ast.CompositeLit); ok && literal.Type != nil && guardProjectionImported(current.file, &ast.CallExpr{Fun: literal.Type}, processLifetimeOwnerImport) == "Registration" {
				found = true
			}
			return !found
		})
		return found
	}
	from := candidates[start]
	return builds(from) || functions.reaches(from, builds), nil
}

// processLifetimeFileDiags reports each private lifetime mechanic in one production file, each
// owner call that no census entry reaches, and each function that reads the process vector and
// builds a subprocess.Registration by hand. Comments and strings are not call sites.
func processLifetimeFileDiags(file, rel string, reached map[string]bool) []string {
	body := readIfExists(file)
	if !slices.ContainsFunc(processLifetimePrefilter, func(token string) bool { return strings.Contains(body, token) }) {
		return nil
	}
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, file, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for the process-lifetime closure: " + err.Error()}
	}
	adapter, _ := filepath.Match(processLifetimeSetpgidAdapters, rel)
	at := func(node ast.Node) string { return rel + ":" + strconv.Itoa(fset.Position(node.Pos()).Line) }
	var diags []string
	var stack []ast.Node
	ast.Inspect(parsed, func(node ast.Node) bool {
		if node == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, node)
		switch n := node.(type) {
		case *ast.KeyValueExpr:
			if key, ok := n.Key.(*ast.Ident); ok && key.Name == "Setpgid" && !adapter {
				diags = append(diags, at(n)+" sets up a process group outside "+processLifetimeOwner)
			}
		case *ast.SelectorExpr:
			if n.Sel.Name == "Setpgid" && !adapter {
				diags = append(diags, at(n)+" sets up a process group outside "+processLifetimeOwner)
			}
		}
		parent := stack[max(len(stack)-2, 0)]
		if selector, ok := parent.(*ast.SelectorExpr); ok && selector.Sel == node {
			return true
		}
		if spelling := processLifetimeKillSpelling(parsed, node); spelling != "" {
			call, called := parent.(*ast.CallExpr)
			switch {
			case !called || call.Fun != node:
				diags = append(diags, at(node)+" uses "+spelling+" as a value outside "+processLifetimeOwner)
			case len(call.Args) != 2:
			case !processLifetimeSignalZero(parsed, call.Args[1]):
				diags = append(diags, at(call)+" sends "+expressionText(fset, call.Args[1])+" to "+expressionText(fset, call.Args[0])+" outside "+processLifetimeOwner)
			case processLifetimeNegative(call.Args[0], processLifetimeEnclosing(stack), nil) && processLifetimeInLoop(stack):
				diags = append(diags, at(call)+" polls process group "+expressionText(fset, call.Args[0])+" in a private drain loop outside "+processLifetimeOwner)
			}
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
			if process, ok := selector.X.(*ast.SelectorExpr); ok && process.Sel.Name == "Process" && (selector.Sel.Name == "Kill" || selector.Sel.Name == "Signal") {
				diags = append(diags, at(call)+" signals a child through "+expressionText(fset, selector)+" outside "+processLifetimeOwner)
			}
			if selector.Sel.Name == "Wait" && len(call.Args) == 0 && processLifetimeCommand(parsed, selector.X) && processLifetimeSpawned(stack) {
				diags = append(diags, at(call)+" waits on "+expressionText(fset, selector.X)+" in a private goroutine outside "+processLifetimeOwner)
			}
		}
		if name := guardProjectionImported(parsed, call, processLifetimeOwnerImport); slices.Contains(processLifetimeOwnerCalls, name) {
			caller := processLifetimeDeclName(stack)
			if !reached[path.Dir(rel)+"#"+caller] {
				diags = append(diags, at(call)+" calls subprocess."+name+" from "+caller+", which no process-lifetime census entry reaches")
			}
		}
		name := guardProjectionImported(parsed, call, processLifetimeProtectionImport)
		if reads, counted := processLifetimeVectorReaders[name]; counted && reads(parsed, call.Args) {
			if decl := processLifetimeDecl(stack); decl != nil {
				builds, err := processLifetimeBuildsRegistration(fset, file, decl)
				switch {
				case err != nil:
					diags = append(diags, rel+" package cannot be parsed for the hand-built registration rule: "+err.Error())
				case builds:
					diags = append(diags, at(call)+" "+decl.Name.Name+" reads the process vector through protection."+name+" and builds a subprocess.Registration by hand outside "+processLifetimeOwner)
				}
			}
		}
		return true
	})
	return diags
}

// processLifetimeKillSpelling names the Kill of a signal package that node selects, or "".
func processLifetimeKillSpelling(file *ast.File, node ast.Node) string {
	switch node.(type) {
	case *ast.SelectorExpr, *ast.Ident:
	default:
		return ""
	}
	for _, pkg := range processLifetimeSignalPackages {
		if guardProjectionImported(file, &ast.CallExpr{Fun: node.(ast.Expr)}, pkg) == "Kill" {
			return pkg + ".Kill"
		}
	}
	return ""
}

// processLifetimeSignalZero reports whether expr is signal zero: the literal 0, or the
// conversion of 0 to the Signal type of a signal package. Signal zero only observes.
func processLifetimeSignalZero(file *ast.File, expr ast.Expr) bool {
	expr = ast.Unparen(expr)
	if call, ok := expr.(*ast.CallExpr); ok && len(call.Args) == 1 {
		for _, pkg := range processLifetimeSignalPackages {
			if guardProjectionImported(file, call, pkg) == "Signal" {
				return processLifetimeSignalZero(file, call.Args[0])
			}
		}
		return false
	}
	literal, ok := expr.(*ast.BasicLit)
	return ok && literal.Kind == token.INT && literal.Value == "0"
}

// processLifetimeNegative reports whether expr is a negative operand: a unary minus, through
// parentheses, or an identifier that a same-function assignment binds to one.
func processLifetimeNegative(expr ast.Expr, function ast.Node, visited map[string]bool) bool {
	switch n := ast.Unparen(expr).(type) {
	case *ast.UnaryExpr:
		return n.Op == token.SUB
	case *ast.Ident:
		if function == nil || visited[n.Name] {
			return false
		}
		if visited == nil {
			visited = map[string]bool{}
		}
		visited[n.Name] = true
		negative := false
		ast.Inspect(function, func(node ast.Node) bool {
			assign, ok := node.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) != len(assign.Rhs) {
				return !negative
			}
			for i, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name == n.Name && processLifetimeNegative(assign.Rhs[i], function, visited) {
					negative = true
				}
			}
			return !negative
		})
		return negative
	}
	return false
}

// processLifetimeEnclosing is the innermost function on the stack, or nil at package scope.
func processLifetimeEnclosing(stack []ast.Node) ast.Node {
	for i := len(stack) - 1; i >= 0; i-- {
		switch stack[i].(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			return stack[i]
		}
	}
	return nil
}

// processLifetimeInLoop reports whether the innermost function runs the node inside a loop.
func processLifetimeInLoop(stack []ast.Node) bool {
	for i := len(stack) - 1; i >= 0; i-- {
		switch stack[i].(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			return true
		case *ast.FuncDecl, *ast.FuncLit:
			return false
		}
	}
	return false
}

// processLifetimeSpawned reports whether the innermost function is the body of a go statement,
// or the node is the call of one.
func processLifetimeSpawned(stack []ast.Node) bool {
	for i := len(stack) - 1; i >= 1; i-- {
		switch stack[i].(type) {
		case *ast.GoStmt:
			return true
		case *ast.FuncLit:
			call, ok := stack[i-1].(*ast.CallExpr)
			return ok && i >= 2 && call.Fun == stack[i] && isGoStmt(stack[i-2])
		case *ast.FuncDecl:
			return false
		}
	}
	return false
}

func isGoStmt(node ast.Node) bool {
	_, ok := node.(*ast.GoStmt)
	return ok
}

// processLifetimeCommand reports whether expr is an identifier that the file declares as an
// os/exec command: a *exec.Cmd parameter or variable, or the result of exec.Command or
// exec.CommandContext.
func processLifetimeCommand(file *ast.File, expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok || ident.Obj == nil {
		return false
	}
	isCmdType := func(typ ast.Expr) bool {
		star, ok := typ.(*ast.StarExpr)
		return ok && guardProjectionImported(file, &ast.CallExpr{Fun: star.X}, "os/exec") == "Cmd"
	}
	switch decl := ident.Obj.Decl.(type) {
	case *ast.Field:
		return isCmdType(decl.Type)
	case *ast.ValueSpec:
		if decl.Type != nil {
			return isCmdType(decl.Type)
		}
		for i, name := range decl.Names {
			if name.Name == ident.Name && i < len(decl.Values) {
				return processLifetimeCommandCall(file, decl.Values[i])
			}
		}
	case *ast.AssignStmt:
		for i, lhs := range decl.Lhs {
			if name, ok := lhs.(*ast.Ident); ok && name.Name == ident.Name && i < len(decl.Rhs) && len(decl.Lhs) == len(decl.Rhs) {
				return processLifetimeCommandCall(file, decl.Rhs[i])
			}
		}
	}
	return false
}

func processLifetimeCommandCall(file *ast.File, expr ast.Expr) bool {
	call, ok := ast.Unparen(expr).(*ast.CallExpr)
	if !ok {
		return false
	}
	name := guardProjectionImported(file, call, "os/exec")
	return name == "Command" || name == "CommandContext"
}

// processLifetimeDeclName names the top-level function declaration on the stack, with its
// receiver-free name, or "package scope".
func processLifetimeDeclName(stack []ast.Node) string {
	if fn := processLifetimeDecl(stack); fn != nil {
		return fn.Name.Name
	}
	return "package scope"
}

// processLifetimeDecl is the top-level function declaration on the stack, or nil at package scope.
func processLifetimeDecl(stack []ast.Node) *ast.FuncDecl {
	for _, node := range stack {
		if fn, ok := node.(*ast.FuncDecl); ok {
			return fn
		}
	}
	return nil
}

// TestProcessLifetimeMigrationClosure is PL91. It grades the private-mechanics half on
// synthetic trees: each rejected mechanic turns red through its alias, dot import, and
// negative-operand spellings; read-only observations, platform adapters, the owner tree,
// comments, strings, and test files stay green.
func TestProcessLifetimeMigrationClosure(t *testing.T) {
	const rel = "internal/example/loop.go"
	grade := func(body string, extra map[string]string) []string {
		t.Helper()
		files := map[string]string{"internal/example/run.go": processLifetimeOwned, rel: body}
		for path, content := range extra {
			files[path] = content
		}
		return processLifetimeClosureDiags(throwawayRoot{files: files}.build(t), processLifetimeFixtureCensus)
	}
	red := []struct{ name, body, want string }{
		{"private scanner loop", "package example\n\nimport (\n\t\"os/exec\"\n\t\"syscall\"\n)\n\nfunc scan(cmd *exec.Cmd) {\n\tcmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}\n\tdone := make(chan error, 1)\n\tgo func() { done <- cmd.Wait() }()\n\t_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)\n\t_ = cmd.Process.Kill()\n}\n",
			rel + ":9 sets up a process group outside internal/subprocess"},
		{"private waiter", "package example\n\nimport \"os/exec\"\n\nfunc wait() chan error {\n\tcmd := exec.Command(\"true\")\n\tdone := make(chan error, 1)\n\tgo func() { done <- cmd.Wait() }()\n\treturn done\n}\n",
			rel + ":8 waits on cmd in a private goroutine outside internal/subprocess"},
		{"aliased negative group signal", "package example\n\nimport sys \"syscall\"\n\nfunc stop(pid int) {\n\tgroup := -pid\n\t_ = sys.Kill((group), sys.SIGKILL)\n}\n",
			rel + ":7 sends sys.SIGKILL to (group) outside internal/subprocess"},
		{"dot-imported unix signal", "package example\n\nimport . \"golang.org/x/sys/unix\"\n\nfunc stop(pid int) { _ = Kill(pid, SIGTERM) }\n",
			rel + ":5 sends SIGTERM to pid outside internal/subprocess"},
		{"private drain loop", "package example\n\nimport \"syscall\"\n\nfunc drain(pid int) {\n\tgroup := -pid\n\tfor syscall.Kill((group), 0) == nil {\n\t}\n}\n",
			rel + ":7 polls process group (group) in a private drain loop outside internal/subprocess"},
		{"kill as a value", "package example\n\nimport \"syscall\"\n\nvar signal = syscall.Kill\n",
			rel + ":5 uses syscall.Kill as a value outside internal/subprocess"},
		{"child signal", "package example\n\nimport (\n\t\"os\"\n\t\"os/exec\"\n)\n\nfunc stop(cmd *exec.Cmd) { _ = cmd.Process.Signal(os.Interrupt) }\n",
			rel + ":8 signals a child through cmd.Process.Signal outside internal/subprocess"},
	}
	for _, c := range red {
		if diags := grade(c.body, nil); !slices.Contains(diags, c.want) {
			t.Errorf("%s: want %q, got %q", c.name, c.want, diags)
		}
	}
	green := []struct {
		name, body string
		extra      map[string]string
	}{
		{"read-only observations", "package example\n\nimport \"syscall\"\n\nfunc present(pid, pgid int) bool {\n\tfor i := 0; i < 1; i++ {\n\t\t_ = syscall.Kill(pid, 0)\n\t}\n\treturn syscall.Kill(-pgid, syscall.Signal(0)) == nil\n}\n", nil},
		{"comments and strings", "package example\n\n// stop once called syscall.Kill(-pid, syscall.SIGKILL) with Setpgid.\nconst stop = \"syscall.Kill(-pid, syscall.SIGKILL)\"\n", nil},
		{"owner tree, platform adapter, and test file", "package example\n", map[string]string{
			"internal/subprocess/control.go":          "package subprocess\n\nimport \"syscall\"\n\nfunc stop(pid int) { _ = syscall.Kill(-pid, syscall.SIGKILL) }\n",
			"internal/runbinary/sysprocattr_linux.go": "package runbinary\n\nimport \"syscall\"\n\nfunc attrs() *syscall.SysProcAttr { return &syscall.SysProcAttr{Setpgid: true} }\n",
			"internal/example/loop_test.go":           "package example\n\nimport \"syscall\"\n\nfunc stop(pid int) { _ = syscall.Kill(-pid, syscall.SIGKILL) }\n",
			"internal/example/synchronous.go":         "package example\n\nimport \"os/exec\"\n\nfunc wait() error {\n\tcmd := exec.Command(\"true\")\n\t_ = cmd.Start()\n\treturn cmd.Wait()\n}\n",
		}},
	}
	for _, c := range green {
		if diags := grade(c.body, c.extra); len(diags) != 0 {
			t.Errorf("%s: want no diagnostics, got %q", c.name, diags)
		}
	}
}
