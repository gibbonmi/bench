package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// benchModule is the import path prefix of this module's packages.
const benchModule = "github.com/gibbonmi/bench/"

// testrepoDir keeps the gate fixtures after its Git half moved to gittest.
const testrepoDir = "internal/testrepo"

// testrepoFormerExports are the snapshot functions that moved from testrepo to gittest.
var testrepoFormerExports = []string{"CommitWorkingTree", "CommitAll"}

// gitIdentityKeys are the Git configuration keys and the environment names that set a commit
// identity. A fixture literal equal to one, or one followed by "=", states an identity.
var gitIdentityKeys = []string{"user.name", "user.email", "GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL"}

// gitIdentityHeader matches a raw commit or tag header line that states an identity.
var gitIdentityHeader = regexp.MustCompile(`(?m)^(author|committer|tagger) [^<\n]*<[^>\n]*>`)

// gitFixtureExceptions is the exact set of fixture sites that keep their own Git process or
// identity outside gittest, by file, enclosing declaration, and kind. Each value names the
// contract that gittest cannot carry. A sibling declaration in the same file has no entry.
var gitFixtureExceptions = map[ownedSite]string{
	{"cmd/bench/wrapper_kit_test.go", "newWrapperKitFixture", "git-runner"}:                              "a linked worktree that the wrapper resolution grades",
	{"internal/commit/merge_parent_test.go", "startPendingMerge", "git-runner"}:                          "an expected merge conflict with a nonzero exit",
	{"internal/conformance/git_fixture_owner_test.go", "gitIdentityKeys", "git-identity"}:                "the identity rule's own key vocabulary",
	{"internal/conformance/handoff_single_source_test.go", "trackedPaths", "git-runner"}:                 "an error-tolerant read of the graded root",
	{"internal/conformance/validity_checks_test.go", "checkExecutableGitModes", "git-runner"}:            "a read of the graded root's index modes",
	{"internal/conformance/validity_checks_test.go", "checkGateInputPathsNotIgnored", "git-runner"}:      "an ignore query whose exit status is the answer",
	{"internal/diff/command_test.go", "rawUntrackedPatch", "git-runner"}:                                 "an expected nonzero exit of git diff --no-index",
	{"internal/diff/command_test.go", "rawWhitespaceOffenses", "git-runner"}:                             "an expected nonzero exit of git diff --check",
	{"internal/diff/compatibility_test.go", "seedCompatibilityRepo", "git-identity"}:                     "an identity that the pinned commit ids grade",
	{"internal/diff/compatibility_test.go", "runGitEnv", "git-runner"}:                                   "an explicit child environment of commit dates that the pinned commit ids grade",
	{"internal/diff/matrix_test.go", "assertRawLiveSnapshot", "git-runner"}:                              "a git diff --check probe whose exit status is the answer",
	{"internal/diff/matrix_test.go", "rawNoIndexPatch", "git-runner"}:                                    "an expected exit 1 of git diff --no-index",
	{"internal/env/git_policy_test.go", "TestGitPolicyHasOneSource", "git-runner"}:                       "a production-subject query of the source checkout",
	{"internal/env/git_policy_test.go", "gitPolicyFindings", "git-runner"}:                               "a production-subject query of the source checkout",
	{"internal/env/kit_run_test.go", "checkKitRunIsolation", "git-runner"}:                               "an explicit child environment",
	{"internal/handoff/state_scan_test.go", "writeBlob", "git-runner"}:                                   "custom standard input",
	{"internal/preflight/preflighttest/reviewfiles.go", "SeedReviewPair", "git-runner"}:                  "the fast-import input stream with its pinned committer and times",
	{"internal/roadmap/retro_cost_test.go", "costRepo.gitAuthored", "git-runner"}:                        "deliberate author and committer times",
	{"internal/worktree/reset_refusal_test.go", "TestResetRefusesAnAmbiguousCheckpoint", "git-identity"}: "raw commit headers with fixed bytes for a hash-prefix collision search",
	{"internal/worktree/test_run_test.go", "benchSourceRoot", "git-runner"}:                              "an error-returning query of the source checkout, outside the journey harness file",
}

// systemProcessSites is the exact set of process and gittest sites that internal/systemtest
// holds. Every system Git child starts through owner.runAt, which joins the starts ledger
// and the repository budget, so a Git process anywhere else in the package fails the census.
var systemProcessSites = map[ownedSite]string{
	{"internal/systemtest/owner_test.go", "systemOwner.runWithInput", "process"}:                     "the one owned process start that runAt wraps",
	{"internal/systemtest/owner_test.go", "systemOwner.startProcessGroup", "process"}:                "an owned process group start on the starts ledger",
	{"internal/systemtest/owner_test.go", "systemOwner.timeoutProcessGroup", "process"}:              "a shell sleep tree that the timeout outcome kills",
	{"internal/systemtest/owner_test.go", "systemOwner.interruptProcessGroup", "process"}:            "the test binary as the interrupt child",
	{"internal/systemtest/owner_test.go", "runInterruptChild", "process"}:                            "the sleep descendant of the interrupt child",
	{"internal/systemtest/owner_land_race_test.go", "systemStartSelected", "process"}:                "the selected executable, started on the starts ledger",
	{"internal/systemtest/owner_land_race_test.go", "systemCommitArgs", "gittest"}:                   "an identity argv that owner.runAt runs",
	{"internal/systemtest/owner_land_race_test.go", "systemIdentity", "gittest"}:                     "identity config commands that systemGit runs",
	{"internal/systemtest/owner_artifact_children_test.go", "startArtifactLand", "process"}:          "the selected executable in its own process group, started on the starts ledger",
	{"internal/systemtest/owner_artifact_recovery_test.go", "plantLiveProspectiveBundle", "process"}: "a live sleep process as a recorded bundle owner",
	// A scratch gate repository under a test temporary directory sits outside the owner's
	// three-repository budget. The subject is the gate's child process, not the Git setup.
	{"internal/systemtest/process_lifetime_cache_test.go", "TestProcessLifetimeCacheSecondProcessClean", "gittest"}:         "a scratch gate repository",
	{"internal/systemtest/process_lifetime_preservation_test.go", "TestProcessLifetimeGateDifferential", "gittest"}:         "a scratch gate repository",
	{"internal/systemtest/process_lifetime_preservation_test.go", "TestProcessLifetimeStartFailureDifferential", "gittest"}: "a scratch gate repository",
}

// systemProcessDiags grades the internal/systemtest sites of the census against
// systemProcessSites. A site outside the list fails, and so does an entry with no site.
func systemProcessDiags(sites []architectureSite) []string {
	refused, diags := exactSites(sites, systemProcessSites, "systemProcessSites", func(string) bool { return true })
	for _, site := range refused {
		diags = append(diags, fmt.Sprintf("%s in %s; a system Git child starts through owner.runAt, and any other system process needs a systemProcessSites entry", site.diagnostic(), site.symbol))
	}
	return diags
}

// gitFixtureOwnershipDiags grades the fixture half of the Git ownership rule. Fixture code is
// a test file or a file that imports testing. Outside gittest, fixture code starts no Git
// process and states no Git identity unless gitFixtureExceptions lists the site. No file
// keeps a former testrepo Git export, and neither fixture package imports the other.
// The rule reads syntax, so a literal or a comment that quotes such source is data.
func gitFixtureOwnershipDiags(root string) []string {
	var diags []string
	var found []architectureSite
	for _, path := range moduleGoFiles(root) {
		rel := slashRel(root, path)
		fileDiags, sites := gitFixtureSites(rel, readIfExists(path))
		diags = append(diags, fileDiags...)
		found = append(found, sites...)
	}
	// A throwaway root without an entry's package grades none of that package's entries.
	// A tree that holds the package must hold each entry's live site, even after a file move.
	refused, stale := exactSites(found, gitFixtureExceptions, "gitFixtureExceptions", func(path string) bool {
		return exists(filepath.Join(root, filepath.FromSlash(filepath.Dir(path))))
	})
	diags = append(diags, stale...)
	for _, site := range refused {
		action := "run it through gittest.Run, MustRun, Output, OutputBytes, or Commit"
		if site.kind == "git-identity" {
			action = "use gittest.IdentityConfig or gittest.IdentityArgs"
		}
		diags = append(diags, fmt.Sprintf("%s:%d (%s) holds a %s outside gittest; %s, or record the specialized site in gitFixtureExceptions", site.path, site.line, site.symbol, site.kind, action))
	}
	return uniqueSorted(diags)
}

// gitFixtureSites parses one module file. It returns the diagnostics that no exception can
// admit and the Git process and identity sites that the exact exceptions grade.
func gitFixtureSites(rel, body string) ([]string, []architectureSite) {
	if body == "" {
		return nil, nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, rel, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for Git fixture ownership: " + err.Error()}, nil
	}
	var diags []string
	edges := map[string]string{gitAdapterHarness: testrepoDir}
	if !strings.HasSuffix(rel, "_test.go") {
		edges[testrepoDir] = gitAdapterHarness
	}
	for owner, forbidden := range edges {
		if strings.HasPrefix(rel, owner+"/") && importsPath(file, benchModule+forbidden) {
			diags = append(diags, fmt.Sprintf("%s imports %s; the gittest and testrepo fixtures stay independent", rel, forbidden))
		}
	}
	if strings.HasPrefix(rel, gitAdapterHarness+"/") {
		return diags, nil
	}
	fixture := strings.HasSuffix(rel, "_test.go") || importsPath(file, "testing")
	var sites []architectureSite
	inspectDecls(file, func(symbol string, node ast.Node) bool {
		at := func(kind string) architectureSite {
			return architectureSite{path: rel, line: fset.Position(node.Pos()).Line, kind: kind, symbol: symbol}
		}
		switch n := node.(type) {
		case *ast.FuncDecl:
			if strings.HasPrefix(rel, testrepoDir+"/") && n.Recv == nil && slices.Contains(testrepoFormerExports, n.Name.Name) {
				diags = append(diags, fmt.Sprintf("%s:%d defines the former testrepo export %s; gittest owns it", rel, at("").line, n.Name.Name))
			}
		case *ast.SelectorExpr:
			if name := importedName(file, n, benchModule+testrepoDir); slices.Contains(testrepoFormerExports, name) {
				diags = append(diags, fmt.Sprintf("%s:%d (%s) calls the former testrepo export %s; call the gittest snapshot", rel, at("").line, symbol, name))
			}
		case *ast.CallExpr:
			if fixture && gitProcessCall(file, n) {
				sites = append(sites, at("git-runner"))
			}
		case *ast.BasicLit:
			if value, ok := stringLiteralValue(n); fixture && ok && gitIdentityLiteral(value) {
				sites = append(sites, at("git-identity"))
			}
		}
		return true
	})
	return diags, sites
}

// gitProcessCall reports whether call starts git through os/exec, under any import name.
func gitProcessCall(file *ast.File, call *ast.CallExpr) bool {
	program := 0
	switch guardProjectionImported(file, call, "os/exec") {
	case "Command":
	case "CommandContext":
		program = 1
	default:
		return false
	}
	value, ok := "", len(call.Args) > program
	if ok {
		value, ok = stringLiteralValue(call.Args[program])
	}
	return ok && value == "git"
}

func gitIdentityLiteral(value string) bool {
	for _, key := range gitIdentityKeys {
		if value == key || strings.HasPrefix(value, key+"=") {
			return true
		}
	}
	return gitIdentityHeader.MatchString(value)
}

// importedName gives the name that expr selects from the package at path, or "". It reads
// the selector through the call resolver, so an import alias and a shadowing local resolve
// the same way for a reference as for a call.
func importedName(file *ast.File, expr ast.Expr, path string) string {
	return guardProjectionImported(file, &ast.CallExpr{Fun: expr}, path)
}

func importsPath(file *ast.File, path string) bool {
	for _, spec := range file.Imports {
		if value, ok := stringLiteralValue(spec.Path); ok && value == path {
			return true
		}
	}
	return false
}

// inspectDecls walks each top-level declaration of file and passes visit the name of the
// declaration that encloses each node.
func inspectDecls(file *ast.File, visit func(symbol string, node ast.Node) bool) {
	for _, decl := range file.Decls {
		symbol := declSymbol(decl)
		ast.Inspect(decl, func(node ast.Node) bool { return node == nil || visit(symbol, node) })
	}
}

// declSymbol names a top-level declaration: a function, a method as Receiver.Name, or the
// first name of a var or const declaration. An import declaration has no name.
func declSymbol(decl ast.Decl) string {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Recv == nil || len(d.Recv.List) == 0 {
			return d.Name.Name
		}
		receiver := d.Recv.List[0].Type
		for {
			switch r := receiver.(type) {
			case *ast.StarExpr:
				receiver = r.X
				continue
			case *ast.IndexExpr:
				receiver = r.X
				continue
			case *ast.IndexListExpr:
				receiver = r.X
				continue
			case *ast.Ident:
				return r.Name + "." + d.Name.Name
			}
			return d.Name.Name
		}
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.ValueSpec:
				return s.Names[0].Name
			case *ast.TypeSpec:
				return s.Name.Name
			}
		}
	}
	return ""
}

// ownedSite keys an exact allowlist entry: a file, the top-level declaration that encloses
// the site, and the site kind. A file-wide or a kind-free entry would admit a sibling site.
type ownedSite struct{ path, symbol, kind string }

// exactSites returns the found sites that allowed does not list, and a sorted diagnostic for
// each entry of the list named list that no found site matches, when graded holds its path.
func exactSites(found []architectureSite, allowed map[ownedSite]string, list string, graded func(path string) bool) (refused []architectureSite, stale []string) {
	live := map[ownedSite]bool{}
	for _, site := range found {
		key := ownedSite{site.path, site.symbol, site.kind}
		if _, ok := allowed[key]; ok {
			live[key] = true
			continue
		}
		refused = append(refused, site)
	}
	for key := range allowed {
		if !live[key] && graded(key.path) {
			stale = append(stale, fmt.Sprintf("%s names %s %s (%s), which holds no such site", list, key.path, key.symbol, key.kind))
		}
	}
	slices.Sort(stale)
	return refused, stale
}

// TestGitFixtureOwnership is GF10. A retained testrepo export, a caller of one under an
// import alias, and an import edge between the two fixture packages each fail the
// registered check, which grades the live tree through the root conformance run.
func TestGitFixtureOwnership(t *testing.T) {
	requireOneGitDiag(t, map[string]map[string]string{
		"retained definition":      {testrepoDir + "/working_tree.go": "package testrepo\n\nfunc CommitAll(root, message string) error { return nil }\n"},
		"aliased caller":           {"internal/example/example_test.go": "package example\n\nimport (\n\t\"testing\"\n\n\ttr \"" + benchModule + testrepoDir + "\"\n)\n\nfunc TestX(t *testing.T) { _ = tr.CommitWorkingTree(\"a\", \"b\") }\n"},
		"gittest imports testrepo": {gitAdapterHarness + "/snapshot.go": "package gittest\n\nimport _ \"" + benchModule + testrepoDir + "\"\n"},
		"testrepo imports gittest": {testrepoDir + "/gate_fixture.go": "package testrepo\n\nimport _ \"" + benchModule + gitAdapterHarness + "\"\n"},
	}, map[string]string{
		"retained definition":      "defines the former testrepo export CommitAll",
		"aliased caller":           "calls the former testrepo export CommitWorkingTree",
		"gittest imports testrepo": "imports " + testrepoDir,
		"testrepo imports gittest": "imports " + gitAdapterHarness,
	})
}

// requireOneGitDiag builds one throwaway root for each case and requires exactly one
// registered Git ownership diagnostic that holds the case's wanted text.
func requireOneGitDiag(t *testing.T, cases map[string]map[string]string, want map[string]string) {
	t.Helper()
	for name, files := range cases {
		t.Run(name, func(t *testing.T) {
			diags := checkGitPlumbingOwner(throwawayRoot{files: files}.build(t))
			if len(diags) != 1 || !strings.Contains(diags[0], want[name]) {
				t.Fatalf("diagnostics = %v, want one that holds %q", diags, want[name])
			}
		})
	}
}

// gitFixtureSource renders one example test file that imports os/exec, testing, and gittest
// and holds body as the body of TestX.
func gitFixtureSource(body string) string {
	return "package example\n\nimport (\n\trun \"os/exec\"\n\t\"testing\"\n\n\t\"" + benchModule + gitAdapterHarness + "\"\n)\n\nvar _ = run.Command\nvar _ = gittest.MustRun\n\nfunc TestX(t *testing.T) {\n\troot := t.TempDir()\n\t" + body + "\n}\n"
}

// TestGitFixtureOwnershipBites is GF24. A generic Git process outside gittest fails the
// registered check, under an import alias, in a context form, and in a non-test fixture
// helper. A production file that starts git is the administration rule's subject alone.
func TestGitFixtureOwnershipBites(t *testing.T) {
	helper := "package exampletest\n\nimport (\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc Status(t testing.TB, root string) { _ = exec.Command(\"git\", \"-C\", root, \"status\").Run() }\n"
	requireOneGitDiag(t, map[string]map[string]string{
		"aliased runner":  {"internal/example/example_test.go": gitFixtureSource(`_, _ = run.Command("git", "-C", root, "status").Output()`)},
		"context runner":  {"internal/example/example_test.go": gitFixtureSource(`_ = run.CommandContext(t.Context(), "git", "-C", root, "status").Run()`)},
		"fixture helper":  {"internal/exampletest/fixture.go": helper},
		"runner in a var": {"internal/example/example_test.go": gitFixtureSource("") + "\nvar status = func(root string) error { return run.Command(\"git\", \"-C\", root, \"status\").Run() }\n"},
	}, map[string]string{
		"aliased runner":  "example_test.go:15 (TestX) holds a git-runner outside gittest",
		"context runner":  "(TestX) holds a git-runner outside gittest",
		"fixture helper":  "fixture.go:8 (Status) holds a git-runner outside gittest",
		"runner in a var": "(status) holds a git-runner outside gittest",
	})
	production := "package example\n\nimport \"os/exec\"\n\nfunc status(root string) error { return exec.Command(\"git\", \"-C\", root, \"status\").Run() }\n"
	if diags := checkGitPlumbingOwner(throwawayRoot{files: map[string]string{"internal/example/example.go": production}}.build(t)); len(diags) != 0 {
		t.Fatalf("production Git process: diagnostics = %v, want none", diags)
	}
}

// TestGitFixtureIdentityBites is GF25. A default identity outside gittest fails the
// registered check even when gittest runs the command: a configuration key, a per-command
// option, an environment name, and a raw commit header each count.
func TestGitFixtureIdentityBites(t *testing.T) {
	requireOneGitDiag(t, map[string]map[string]string{
		"config key":     {"internal/example/example_test.go": gitFixtureSource(`gittest.MustRun(t, root, "config", "user.name", "bench")`)},
		"command option": {"internal/example/example_test.go": gitFixtureSource(`gittest.MustRun(t, root, "-c", "user.email=bench@local", "commit", "-qm", "x")`)},
		"environment":    {"internal/example/example_test.go": gitFixtureSource(`t.Setenv("GIT_COMMITTER_NAME", "bench")`)},
		"raw header":     {"internal/example/example_test.go": gitFixtureSource(`_ = "tree x\nauthor bench <bench@local> 1 +0000\n"`)},
	}, map[string]string{
		"config key":     "(TestX) holds a git-identity outside gittest",
		"command option": "(TestX) holds a git-identity outside gittest",
		"environment":    "(TestX) holds a git-identity outside gittest",
		"raw header":     "(TestX) holds a git-identity outside gittest",
	})
}

// TestFixtureOwnershipLiteralData is GF27. Source that a literal or a comment quotes is
// data, so neither the Git nor the skip ownership check reports it.
func TestFixtureOwnershipLiteralData(t *testing.T) {
	quoted := "package example\n\n" +
		"// A comment that quotes exec.Command(\"git\", \"status\") and testrepo.CommitAll is data.\n" +
		"var examples = []string{\n" +
		"\t\"exec.Command(\\\"git\\\", \\\"-C\\\", root, \\\"config\\\", \\\"user.name\\\", \\\"x\\\")\",\n" +
		"\t\"testrepo.CommitAll(root, \\\"snapshot\\\")\",\n" +
		"\t\"capability.Capability(t, capability.Fifo, \\\"no fifo\\\")\",\n" +
		"}\n"
	root := throwawayRoot{files: map[string]string{"internal/example/example_test.go": quoted}}.build(t)
	if diags := append(checkGitPlumbingOwner(root), checkSkipOwnership(root)...); len(diags) != 0 {
		t.Fatalf("quoted source: diagnostics = %v, want none", diags)
	}
}

// TestFixtureOwnershipExceptionScope is GF28. The newWrapperKitFixture exception admits
// that function's Git process alone. A sibling runner in the same file fails, and so does
// an identity inside the excepted function, as the exception names one kind. An exception
// whose site is gone from a package the tree holds fails as stale.
func TestFixtureOwnershipExceptionScope(t *testing.T) {
	const rel = "cmd/bench/wrapper_kit_test.go"
	allowed := "package main\n\nimport (\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc newWrapperKitFixture(t *testing.T, root string) {\n\t_ = exec.Command(\"git\", \"worktree\", \"add\", root).Run()\n%s}\n"
	if diags := checkGitPlumbingOwner(throwawayRoot{files: map[string]string{rel: fmt.Sprintf(allowed, "")}}.build(t)); len(diags) != 0 {
		t.Fatalf("the excepted site alone: diagnostics = %v, want none", diags)
	}
	sibling := fmt.Sprintf(allowed, "") + "\nfunc plantedGitRunner(dir string) ([]byte, error) { return exec.Command(\"git\", \"-C\", dir, \"status\").Output() }\n"
	requireOneGitDiag(t, map[string]map[string]string{
		"sibling runner":   {rel: sibling},
		"identity in site": {rel: fmt.Sprintf(allowed, "\t_ = exec.Command(\"ls\", \"user.name\")\n")},
		"site gone":        {"cmd/bench/main_test.go": "package main\n"},
	}, map[string]string{
		"sibling runner":   "(plantedGitRunner) holds a git-runner outside gittest",
		"identity in site": "(newWrapperKitFixture) holds a git-identity outside gittest",
		"site gone":        "gitFixtureExceptions names " + rel + " newWrapperKitFixture (git-runner), which holds no such site",
	})
}
