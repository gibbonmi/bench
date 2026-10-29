package conformance

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// guardProjectionOwnerDir holds the shared projection and envelope owner files. The guard packages in
// guardProjectionGoverned must use them and hold no private copy of their facts.
var (
	guardProjectionOwnerDir   = "internal/shellcommand"
	guardProjectionOwnerFiles = []string{guardProjectionOwnerDir + "/projection.go", guardProjectionOwnerDir + "/envelope.go"}
	guardProjectionGoverned   = []string{"internal/gitguard", "internal/benchguard", "internal/writeguard"}
)

// guardProjectionDiagPrefix starts each diagnostic of checkGuardProjectionOwner. The other constants
// name the diagnostic pieces for the check and its tests.
const (
	guardProjectionDiagPrefix  = "guard projection owner "
	guardProjectionUnreadable  = guardOwnerUnreadable
	guardProjectionIrregular   = guardOwnerIrregular
	guardProjectionUnparsed    = guardOwnerUnparsed
	guardProjectionEnvelope    = "hook envelope JSON decoding"
	guardProjectionShellName   = "shell-name predicate"
	guardProjectionWrapperFlag = "wrapper -c flag recognizer"
	guardProjectionExecChild   = "worktree exec child extraction"
	guardProjectionFlagWord    = "-c" // the wrapper flag word that the shell-name and flag tests use
)

// checkGuardProjectionOwner keeps the guard command grammar and the hook envelope decoding in
// internal/shellcommand. It engages only on a root with internal/conformance and a governed guard
// package, so a linked consumer passes. There, a missing or unreadable owner or source is a diagnostic.
// It reports each private copy of an owner fact in the production Go of a governed package.
func checkGuardProjectionOwner(root string) []string {
	if info, err := os.Stat(filepath.Join(root, "internal", "conformance")); err != nil || !info.IsDir() {
		return nil
	}
	var production, diags []string
	engaged := false
	for _, dir := range guardProjectionGoverned {
		files, dirDiags, present := guardProjectionSources(root, dir)
		engaged = engaged || present
		production = append(production, files...)
		diags = append(diags, dirDiags...)
	}
	if !engaged {
		return nil
	}
	owned, ownerDiags, _ := guardProjectionSources(root, guardProjectionOwnerDir)
	diags = append(diags, ownerDiags...)
	for _, rel := range guardProjectionOwnerFiles {
		if !slices.Contains(owned, rel) {
			diags = append(diags, guardProjectionDiagPrefix+rel+guardProjectionUnreadable)
		}
	}
	return append(diags, guardOwnerParse(root, guardProjectionDiagPrefix, append(owned, production...), func(fset *token.FileSet, rel string, file *ast.File) []string {
		if strings.HasPrefix(rel, guardProjectionOwnerDir+"/") {
			return nil
		}
		return guardProjectionCopies(fset, rel, file)
	})...)
}

// guardProjectionSources lists the production Go files beneath one governed package. Only an absent
// path, or a path that is not a directory, marks the package absent.
func guardProjectionSources(root, dir string) (production, diags []string, present bool) {
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(dir)))
	if errors.Is(err, fs.ErrNotExist) || (err == nil && !info.IsDir()) {
		return nil, nil, false
	}
	if err != nil {
		return nil, []string{guardProjectionDiagPrefix + dir + guardProjectionUnreadable}, true
	}
	production, diags = guardOwnerSources(root, dir, guardProjectionDiagPrefix)
	return production, diags, true
}

// guardProjectionCopies reports each private copy of an owner fact in one parsed file. It reads syntax
// nodes, so a comment never matches. The shell-name and exec-word tests ask the owner itself.
func guardProjectionCopies(fset *token.FileSet, rel string, file *ast.File) []string {
	var diags []string
	inner := map[ast.Node]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		construct := ""
		switch n := node.(type) {
		case *ast.CallExpr:
			construct = guardProjectionCallCopy(file, n)
		case *ast.BasicLit:
			if word, ok := guardProjectionString(n); ok && len(shellcommand.ShellChildren([]string{word, guardProjectionFlagWord, ""}, 0)) > 0 {
				construct = guardProjectionShellName
			}
		case *ast.BinaryExpr:
			if !inner[n] && (n.Op == token.LAND || n.Op == token.LOR) && guardProjectionExecWords(guardProjectionCompared(n, inner)) {
				construct = guardProjectionExecChild
			}
		}
		if construct != "" {
			diags = append(diags, fmt.Sprintf("%s%s:%d: %s", guardProjectionDiagPrefix, rel, fset.Position(node.Pos()).Line, construct))
		}
		return true
	})
	return diags
}

// guardProjectionCallCopy names the owner fact that a call to encoding/json or regexp copies, or "".
func guardProjectionCallCopy(file *ast.File, call *ast.CallExpr) string {
	switch guardProjectionImported(file, call, "encoding/json") {
	case "Unmarshal", "NewDecoder":
		return guardProjectionEnvelope
	}
	switch guardProjectionImported(file, call, "regexp") {
	case "Compile", "MustCompile", "CompilePOSIX", "MustCompilePOSIX":
		if len(call.Args) > 0 {
			pattern, literal := guardProjectionString(call.Args[0])
			if compiled, err := regexp.Compile(pattern); literal && err == nil && compiled.MatchString(guardProjectionFlagWord) {
				return guardProjectionWrapperFlag
			}
		}
	}
	return ""
}

// guardProjectionImported gives the name that a call selects from the package at path, or "". The
// qualifier must be an unresolved identifier that an import of path binds. A dot import binds a bare name.
func guardProjectionImported(file *ast.File, call *ast.CallExpr, path string) string {
	qualifier, name := ".", ""
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := fun.X.(*ast.Ident); ok && ident.Obj == nil {
			qualifier, name = ident.Name, fun.Sel.Name
		}
	case *ast.Ident:
		if fun.Obj == nil {
			name = fun.Name
		}
	}
	for _, spec := range file.Imports {
		imported, err := strconv.Unquote(spec.Path.Value)
		local := imported[strings.LastIndex(imported, "/")+1:]
		if spec.Name != nil {
			local = spec.Name.Name
		}
		if err == nil && imported == path && local == qualifier {
			return name
		}
	}
	return ""
}

// guardProjectionCompared lists the string literal that each comparison of a logical chain tests, in
// order. Any other operand gives "". It marks each nested chain, so only the outer chain reports.
func guardProjectionCompared(expr ast.Expr, inner map[ast.Node]bool) []string {
	if n, ok := ast.Unparen(expr).(*ast.BinaryExpr); ok {
		switch n.Op {
		case token.LAND, token.LOR:
			inner[n] = true
			return append(guardProjectionCompared(n.X, inner), guardProjectionCompared(n.Y, inner)...)
		case token.EQL, token.NEQ:
			word, ok := guardProjectionString(n.Y)
			if !ok {
				word, _ = guardProjectionString(n.X)
			}
			return []string{word}
		}
	}
	return []string{""}
}

// guardProjectionExecWords reports whether two adjacent words form the worktree exec head of the owner.
func guardProjectionExecWords(words []string) bool {
	for i := 0; i+1 < len(words); i++ {
		if shellcommand.WorktreeExec([]string{"", words[i], words[i+1]}, 0).Head {
			return true
		}
	}
	return false
}

// guardProjectionString gives the value of a string literal expression.
func guardProjectionString(expr ast.Expr) (string, bool) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		value, err := strconv.Unquote(lit.Value)
		return value, err == nil
	}
	return "", false
}

// guardProjectionDiags writes sources beneath root and runs the registered package-core-guard binding
// on root. It keeps only the projection owner diagnostics, so each test grades the registered owner.
func guardProjectionDiags(t *testing.T, root string, sources map[string]string) []string {
	t.Helper()
	for rel, source := range sources {
		writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(rel)), source)
	}
	return guardOwnerDiags(t, root, guardProjectionDiagPrefix)
}

// guardProjectionKit gives the sources of a kit root that engages the subcheck.
func guardProjectionKit() map[string]string {
	sources := map[string]string{"internal/conformance/conformance.go": "package conformance\n", "internal/gitguard/verdict.go": "package gitguard\n", "internal/gitguard/inner/inner.go": "package inner\n"}
	for _, rel := range guardProjectionOwnerFiles {
		sources[rel] = "package shellcommand\n"
	}
	return sources
}

// guardProjectionScanCopy is the retired Git wrapper set, flag recognizer, and exec child extraction,
// renamed. Each fixture packs its declarations tightly, because only its syntax nodes matter.
const guardProjectionScanCopy = `package gitguard; import "regexp"; var shellNames = map[string]bool{"sh": true, "bash": true, "zsh": true}
var childFlagRe = regexp.MustCompile("^-[A-Za-z]*c[A-Za-z]*$")
func execChildWords(words []string) ([]string, bool) { if len(words) < 3 || words[1] != "worktree" || words[2] != "exec" { return nil, false }; return words[3:], true }`

// guardProjectionPrivateCopies maps each planted path to its source. No caller reaches a copy.
var guardProjectionPrivateCopies = map[string]string{
	"internal/gitguard/scan.go":          guardProjectionScanCopy,
	"internal/gitguard/scan_test.go":     guardProjectionScanCopy,
	"internal/gitguard/testdata/copy.go": guardProjectionScanCopy,
	"internal/gitguard/gitguard.go":      `package gitguard; import js "encoding/json"; func commandOf(data []byte) (e struct{ ToolInput struct{ Command string } }) { _ = js.Unmarshal(data, &e); return e }`,
	"internal/gitguard/verdict.go": `package gitguard; import "encoding/json"; var globalOptsWithArg = map[string]bool{"-C": true, "-c": true} // sh -c 'bench worktree exec t -- git push'
func verdict(sub string, args []string) string { switch sub { case "worktree": if len(args) > 0 && args[0] == "remove" { return "bash -c" } }; data, _ := json.Marshal(args); return string(data) }
func decodeLocal(data []byte, json interface{ Unmarshal([]byte, any) error }) error { return json.Unmarshal(data, nil) }`,
	"internal/benchguard/benchguard.go": `package benchguard; import ("encoding/json"; "path/filepath"; "regexp"); var childFlag = regexp.MustCompile("^-[A-Za-z]*c[A-Za-z]*$")
func commandOf(data []byte) (command string, err error) { var envelope struct{ ToolInput map[string]json.RawMessage }; if err = json.Unmarshal(data, &envelope); err == nil { err = json.Unmarshal(envelope.ToolInput["command"], &command) }; return command, err }
func shellHead(word string) bool { switch filepath.Base(word) { case "sh", "bash", "zsh": return true }; return false }
func execHead(words []string, index int) bool { return index+2 < len(words) && words[index+1] == "worktree" && words[index+2] == "exec" }`,
	"internal/benchguard/reviewer.go": `package benchguard; import ("path/filepath"; "regexp"; "slices"); var reviewerFlag = regexp.MustCompile("^-[A-Za-z]*c[A-Za-z]*$")
func reviewerChild(words []string, head int) []string { switch { case head+2 < len(words) && words[head+1] == "worktree" && words[head+2] == "exec": if child := slices.Index(words, "--"); child >= 0 { return words[child+1:] }; case slices.Contains([]string{"sh", "bash", "zsh"}, filepath.Base(words[head])): for i := head + 1; i+1 < len(words) && reviewerFlag.MatchString(words[i]); i++ { return words[i+1:] } }; return nil }`,
	"internal/benchguard/inner/inner.go": `package inner; var shells = []string{"zsh"}`,
	"internal/writeguard/writeguard.go":  `package writeguard; import ("bytes"; "encoding/json"); func pathOf(data []byte) (e struct{ Cwd string }, err error) { return e, json.NewDecoder(bytes.NewReader(data)).Decode(&e) }`,
	"internal/writeguard/field.go":       `package writeguard; import . "encoding/json"; func fieldOf(raw []byte) (value string, err error) { return value, Unmarshal(raw, &value) }`,
}

// TestGuardProjectionOwnerBitesOnPrivateCopies plants renamed, unused copies of each retired recognizer
// and envelope decoder, reviewer.go included. The registered package-core-guard check must report each
// copied fact. Git's own worktree verb, comments, test files, testdata, and a shadowed name stay clean.
func TestGuardProjectionOwnerBitesOnPrivateCopies(t *testing.T) {
	sources := guardProjectionKit()
	for rel, source := range guardProjectionPrivateCopies {
		sources[rel] = source
	}
	var want, got []string
	for rel, constructs := range map[string][]string{
		"internal/gitguard/scan.go":          {guardProjectionShellName, guardProjectionWrapperFlag, guardProjectionExecChild},
		"internal/gitguard/gitguard.go":      {guardProjectionEnvelope},
		"internal/benchguard/benchguard.go":  {guardProjectionEnvelope, guardProjectionWrapperFlag, guardProjectionShellName, guardProjectionExecChild},
		"internal/benchguard/reviewer.go":    {guardProjectionWrapperFlag, guardProjectionExecChild, guardProjectionShellName},
		"internal/benchguard/inner/inner.go": {guardProjectionShellName},
		"internal/writeguard/writeguard.go":  {guardProjectionEnvelope},
		"internal/writeguard/field.go":       {guardProjectionEnvelope},
	} {
		for _, construct := range constructs {
			want = append(want, rel+" "+construct)
		}
	}
	for _, diag := range guardProjectionDiags(t, t.TempDir(), sources) {
		rel, rest, _ := strings.Cut(strings.TrimPrefix(diag, guardProjectionDiagPrefix), ":")
		_, construct, _ := strings.Cut(rest, ": ")
		got = append(got, rel+" "+construct)
	}
	slices.Sort(want)
	slices.Sort(got)
	if got = slices.Compact(got); !slices.Equal(got, want) {
		t.Fatalf("want copies reported\n%s\ngot\n%s", strings.Join(want, "\n"), strings.Join(got, "\n"))
	}
}

// TestGuardProjectionOwnerOptionalSurface grades roots without kit guard source: a consumer without
// internal/conformance, a conformance tree without a guard package, and a guard path that is a file.
// Each gets no projection owner diagnostic, even beside a private copy and with no owner.
func TestGuardProjectionOwnerOptionalSurface(t *testing.T) {
	conformance := "internal/conformance/conformance.go"
	for name, sources := range map[string]map[string]string{
		"consumer":             {"internal/gitguard/scan.go": guardProjectionScanCopy},
		"no guard package":     {conformance: "package conformance\n", "internal/worker/scan.go": guardProjectionScanCopy},
		"guard path is a file": {conformance: "package conformance\n", "internal/gitguard": guardProjectionScanCopy},
	} {
		if diags := guardProjectionDiags(t, t.TempDir(), sources); len(diags) != 0 {
			t.Errorf("%s: want no projection owner diagnostic, got %v", name, diags)
		}
	}
}

// guardProjectionStrip removes every permission of path. It registers the restore before the strip,
// so the removal of TempDir can descend again.
func guardProjectionStrip(t *testing.T, path string) error {
	t.Cleanup(func() { _ = os.Chmod(path, 0o755) })
	if err := os.Chmod(path, 0); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip permissions: %v", err))
	}
	if file, err := os.Open(path); err == nil {
		_ = file.Close()
		capability.Capability(t, capability.Privilege, "a mode 0 path is still readable by this user")
	}
	return nil
}

// TestGuardProjectionOwnerMissingOwner damages one source of a kit root per case. Each absent,
// unreadable, unparseable, or special source gets a named diagnostic, never a clean pass. The test
// calls the subcheck itself, because a sibling package-core-guard subcheck opens a FIFO.
func TestGuardProjectionOwnerMissingOwner(t *testing.T) {
	owner, git := guardProjectionOwnerFiles[len(guardProjectionOwnerFiles)-1], "internal/gitguard/"
	for _, tc := range []struct {
		name, rel, want string
		plant           func(t *testing.T, path string) error
	}{
		{"absent owner", owner, owner + guardProjectionUnreadable, func(t *testing.T, path string) error { return os.Remove(path) }},
		{"unreadable owner", owner, guardProjectionUnparsed + owner, guardProjectionStrip},
		{"unparseable source", git + "verdict.go", guardProjectionUnparsed + git + "verdict.go", func(t *testing.T, path string) error { return os.WriteFile(path, []byte("func ("), 0o644) }},
		{"unstatable package", "internal/writeguard", "internal/writeguard" + guardProjectionUnreadable, func(t *testing.T, path string) error {
			if err := os.Symlink("writeguard", path); err != nil {
				capability.Unavailable(t, capability.Symlink, err.Error())
			}
			return nil
		}},
		{"unreadable subdirectory", git + "inner", git + "inner" + guardProjectionUnreadable, guardProjectionStrip},
		{"special source", git + "x.go", git + "x.go" + guardProjectionIrregular, func(t *testing.T, path string) error {
			if err := syscall.Mkfifo(path, 0o644); err != nil {
				capability.Unavailable(t, capability.Fifo, err.Error())
			}
			return nil
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for rel, source := range guardProjectionKit() {
				writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(rel)), source)
			}
			if err := tc.plant(t, filepath.Join(root, filepath.FromSlash(tc.rel))); err != nil {
				t.Fatal(err)
			}
			result := make(chan []string, 1)
			go func() { result <- checkGuardProjectionOwner(root) }()
			select {
			case diags := <-result:
				if want := guardProjectionDiagPrefix + tc.want; !slices.Contains(diags, want) {
					t.Fatalf("want %q, got %v", want, diags)
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatalf("the check blocked on the %s case", tc.name)
			}
		})
	}
}

// TestGuardProjectionOwnerHoldsOnTheLiveTree runs the registered package-core-guard binding on the kit
// root, which holds the owner and the governed guard packages. It requires no projection owner
// diagnostic, a governed package that engages the subcheck, and its own live-tree classification.
func TestGuardProjectionOwnerHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	root := h.KitRoot
	if !exists(filepath.Join(root, filepath.FromSlash(guardProjectionGoverned[0]))) || !classifiedLiveTreeTest(t.Name()) {
		t.Fatalf("%s needs %s in the kit root and a live-tree classification", t.Name(), guardProjectionGoverned[0])
	}
	if diags := guardProjectionDiags(t, root, nil); len(diags) != 0 {
		t.Fatalf("the kit guards hold a private copy of an owner fact:\n%s", strings.Join(diags, "\n"))
	}
}
