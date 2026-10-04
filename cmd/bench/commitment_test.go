package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
)

// TestCommitmentNextProjection runs the three readers through the root dispatcher on one
// adopted repository. Outcome A is blocked. When B is independent, the roadmap board, its
// context snapshot, the status board, and the dashboard each name B and A's blocker. When B
// depends on A, none of them names a next outcome. The repository also holds an unrelated
// staged spec and an old recommended sequence, so a surviving legacy reader names other work.
func TestCommitmentNextProjection(t *testing.T) {
	const start = "bench commitment start --outcome B --request <request> --deliverable specs/B/spec.md"
	for _, tc := range []struct {
		name      string
		dependent bool
		next      string
		state     string
		command   string
	}{
		{name: "independent successor", next: "B", state: "eligible", command: start},
		{name: "dependent successor", dependent: true, state: "all-blocked", command: "bench commitment plan --input <file>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := commitmenttest.Staged(t, "A", "B")
			if tc.dependent {
				commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
					policy.Milestones[0].Outcomes[1].Dependencies = []string{"A"}
				})
			}
			commitmenttest.Write(t, root, "specs/unrelated/spec.md", commitmenttest.StagedBody)
			commitmenttest.Write(t, root, "ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. unrelated - /bench-implement-spec specs/unrelated/spec.md\n")
			commitmenttest.Commit(t, root, "add unrelated work")
			t.Chdir(root)
			t.Setenv("BENCH_HOME", t.TempDir())
			run := func(args ...string) string {
				t.Helper()
				var stdout, stderr bytes.Buffer
				if code := (Command{Stdout: &stdout, Stderr: &stderr}).Run(args); code != 0 {
					t.Fatalf("bench %s = %d: %s%s", strings.Join(args, " "), code, stdout.String(), stderr.String())
				}
				return spilledResponse(t, stdout.String())
			}
			run("commitment", "block", "--outcome", "A", "--reason", "vendor fix")

			deliverable := ""
			if tc.next != "" {
				deliverable = "specs/" + tc.next + "/spec.md"
			}
			want := map[string]any{"state": tc.state, "active_milestone": "M1", "next_outcome": tc.next, "deliverable": deliverable, "blocked": "A", "command": tc.command}
			for _, args := range [][]string{{"roadmap"}, {"roadmap", "--context"}} {
				document, err := axitest.DecodeDocument(run(args...))
				if err != nil {
					t.Fatal(err)
				}
				outlook, err := document.Rows(commitment.OutlookTable)
				if err != nil || len(outlook) != 1 {
					t.Fatalf("%v outlook = %v, %v", args, outlook, err)
				}
				for field, value := range want {
					if got := outlook[0].(map[string]any)[field]; got != value {
						t.Errorf("%v %s = %v, want %v", args, field, got, value)
					}
				}
				blockers, err := document.Rows(commitment.BlockerTable)
				if err != nil || !reflect.DeepEqual(blockers, []any{map[string]any{"outcome": "A", "reason": "vendor fix"}}) {
					t.Errorf("%v blockers = %v, %v", args, blockers, err)
				}
			}

			detail := strings.TrimSuffix(tc.state+" "+tc.next, " ") + " in M1; blocked A"
			board := run("status", "--all")
			found := false
			for _, line := range strings.Split(board, "\n") {
				found = found || strings.HasPrefix(strings.TrimSpace(line), "commitment "+detail) && strings.HasSuffix(line, "→ "+tc.command)
			}
			if !found || strings.Contains(board, "staged spec") {
				t.Errorf("status board lacks the commitment row %q → %q or names the staged spec:\n%s", detail, tc.command, board)
			}

			page := run("dashboard", "--stdout")
			for _, item := range []string{"<dt>state</dt><dd>" + tc.state + "</dd>", "<dt>A</dt><dd>vendor fix</dd>", "<dt>command</dt><dd>" + strings.NewReplacer("<", "&lt;", ">", "&gt;").Replace(tc.command) + "</dd>"} {
				if !strings.Contains(page, item) {
					t.Errorf("dashboard lacks %q", item)
				}
			}
			if names := strings.Contains(page, "<dt>next_outcome</dt>"); names != (tc.next != "") || strings.Contains(page, "<dt>next_outcome</dt><dd>A</dd>") {
				t.Errorf("dashboard next outcome is wrong for %q", tc.next)
			}
		})
	}
}

// admissionRoute is one command route and the declaration on that route that consumes
// the commitment owner. Use names an owner method that the declaration calls, or the
// adapter type that the declaration constructs.
type admissionRoute struct {
	route, file, declaration, use string
	construct                     bool
}

// TestCommitmentRouteInventory holds each admission route to its consumer. The inventory
// is written by hand, independently of the consumers: if one declaration drops its call
// to the commitment owner, its row fails. Each named method must exist on the owner, and
// each route's verb must be a registered command.
func TestCommitmentRouteInventory(t *testing.T) {
	routes := []admissionRoute{
		{route: "bench commitment start", file: "internal/commitment/commitcmd/admission.go", declaration: "admission", use: "Start"},
		{route: "bench worktree create --from", file: "internal/worktree/tree_target.go", declaration: "createSiblingStart", use: "Inheritance"},
		{route: "bench worktree create --from", file: "internal/worktree/pool_root.go", declaration: "createAt", use: "RegisterSibling"},
		{route: "bench shift --outcome", file: "internal/shift/loop.go", declaration: "loop", use: "ReadyOutcome"},
		{route: "bench preflight build", file: "internal/preflight/command.go", declaration: "verdictCommand", use: "Ready"},
		{route: "bench preflight build --charge", file: "internal/preflight/command.go", declaration: "prepareEvidenceCommand", use: "Ready"},
		{route: "bench preflight evidence --check-current", file: "internal/preflight/charge_pack.go", declaration: "currentEvidenceCommand", use: "Ready"},
		{route: "bench commit", file: "internal/commit/commit.go", declaration: "commitAttributed", use: "AuthorizeCandidate"},
		{route: "bench worktree land", file: "internal/worktree/land.go", declaration: "landAttributed", use: "commitmentAdmission", construct: true},
		{route: "bench worktree land", file: "internal/worktree/land.go", declaration: "commitmentAdmission.Check", use: "AdmitPublication"},
		{route: "bench worktree land", file: "internal/worktree/land.go", declaration: "commitmentAdmission.Publish", use: "PublishAdmitted"},
	}
	owner := reflect.TypeOf(commitrepo.Store{})
	for _, route := range routes {
		verb := strings.Fields(route.route)[1]
		registered := false
		for _, definition := range commandRegistry {
			registered = registered || definition.Name == verb
		}
		if !registered {
			t.Errorf("%s: verb %q is not registered", route.route, verb)
		}
		if _, found := owner.MethodByName(route.use); !route.construct && !found {
			t.Errorf("%s: the commitment owner has no method %s", route.route, route.use)
		}
		if !declarationUses(t, route.file, route.declaration, route.use) {
			t.Errorf("%s: %s %s no longer uses %s", route.route, route.file, route.declaration, route.use)
		}
	}
}

// declarationUses reports whether the named declaration in a repository file calls use as a
// method or constructs use as a composite literal. A call on an imported package name is a
// package function, not a method, so it does not count.
func declarationUses(t *testing.T, file, declaration, use string) bool {
	t.Helper()
	parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join("..", "..", filepath.FromSlash(file)), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	imported := map[string]bool{}
	for _, spec := range parsed.Imports {
		name := path.Base(strings.Trim(spec.Path.Value, `"`))
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imported[name] = true
	}
	for _, decl := range parsed.Decls {
		function, ok := decl.(*ast.FuncDecl)
		if !ok || function.Body == nil || declarationName(function) != declaration {
			continue
		}
		found := false
		ast.Inspect(function.Body, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.CallExpr:
				if selector, ok := node.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == use {
					receiver, named := selector.X.(*ast.Ident)
					found = found || !named || !imported[receiver.Name]
				}
			case *ast.CompositeLit:
				if name, ok := node.Type.(*ast.Ident); ok && name.Name == use {
					found = true
				}
			}
			return !found
		})
		return found
	}
	t.Fatalf("%s declares no %s", file, declaration)
	return false
}

// declarationName is a function's name, or a method's receiver type and name joined by a dot.
func declarationName(function *ast.FuncDecl) string {
	if function.Recv == nil || len(function.Recv.List) == 0 {
		return function.Name.Name
	}
	receiver := function.Recv.List[0].Type
	if star, ok := receiver.(*ast.StarExpr); ok {
		receiver = star.X
	}
	if name, ok := receiver.(*ast.Ident); ok {
		return name.Name + "." + function.Name.Name
	}
	return function.Name.Name
}
