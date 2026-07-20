package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/env"
)

// DATA_HANDLING.md carries a machine-parseable passlist listing between these
// markers; the derivation check reads the same env constants the adapter launch
// enforces from, so the advertisement cannot drift from the enforcement.
const (
	passlistBegin = "<!-- passlist:begin -->"
	passlistEnd   = "<!-- passlist:end -->"
)

// passlistTokenRe matches a backtick-quoted environment name or PREFIX* glob,
// the token form the variable listing documents each passlist pattern as.
var passlistTokenRe = regexp.MustCompile("`([A-Z][A-Z0-9_]*\\*?)`")

// checkDataHandlingDerivation asserts DATA_HANDLING.md's variable listing
// documents every pattern the internal/env passlist admits. The constants are
// the compiled-in enforcement values (env.AgentPasslist, which is SharedBasics
// plus the agent additions), so a pattern added to the code without a row in the
// doc turns the gate red with a diagnostic naming that pattern. The check fails
// loudly — not vacuously — when the marked region is absent or empty, so a doc
// that dropped the listing cannot pass by carrying nothing to check.
func checkDataHandlingDerivation(root string) []string {
	doc := readIfExists(filepath.Join(root, "DATA_HANDLING.md"))
	region, ok := passlistRegion(doc)
	if !ok {
		return []string{fmt.Sprintf("DATA_HANDLING.md passlist derivation region missing: expected the variable listing between %s and %s so the doc can be checked against the internal/env constants", passlistBegin, passlistEnd)}
	}
	documented := passlistTokens(region)
	if len(documented) == 0 {
		return []string{"DATA_HANDLING.md passlist derivation region empty: no `PATTERN` tokens found between the passlist markers"}
	}
	var diags []string
	for _, pattern := range env.AgentPasslist {
		if !documented[pattern] {
			diags = append(diags, fmt.Sprintf("DATA_HANDLING.md passlist derivation: internal/env pattern %q is not documented in the variable listing", pattern))
		}
	}
	return diags
}

func passlistRegion(doc string) (string, bool) {
	i := strings.Index(doc, passlistBegin)
	if i < 0 {
		return "", false
	}
	j := strings.Index(doc, passlistEnd)
	if j < 0 || j < i {
		return "", false
	}
	return doc[i+len(passlistBegin) : j], true
}

func passlistTokens(region string) map[string]bool {
	set := map[string]bool{}
	for _, m := range passlistTokenRe.FindAllStringSubmatch(region, -1) {
		set[m[1]] = true
	}
	return set
}

// controlEscapeRe is the fingerprint of control-rune escaping: rendering a
// control rune into a readable \uXXXX form. Detection code (shift.hasControlByte,
// toon's cell predicate) reports a bool and never emits this escape, so keying on
// the emission separates the single escaping policy from every detector.
var controlEscapeRe = regexp.MustCompile(`\\u%04[xX]`)

// checkSingleControlEscaper asserts exactly one non-test package outside
// internal/toon implements control-rune escaping, so the policy stays
// single-sourced in internal/sanitize instead of being re-copied for a new
// render path. internal/toon is excluded because its cell policy *refuses* a
// control-bearing cell (a closed AXI decision) rather than escaping it.
func checkSingleControlEscaper(root string) []string {
	pkgs := controlEscaperPackages(root)
	switch len(pkgs) {
	case 1:
		return nil
	case 0:
		return []string{"control-rune escaping policy has no owner: no non-test package outside internal/toon implements it (expected internal/sanitize)"}
	default:
		return []string{fmt.Sprintf("control-rune escaping policy has more than one owner: %s — internal/sanitize must be the single source outside internal/toon", strings.Join(pkgs, ", "))}
	}
}

func controlEscaperPackages(root string) []string {
	seen := map[string]bool{}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "tests", "dist":
				return filepath.SkipDir
			}
			if rel, rerr := filepath.Rel(root, path); rerr == nil && filepath.ToSlash(rel) == "internal/toon" {
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		if controlEscapeRe.MatchString(readIfExists(path)) {
			seen[slashRel(root, filepath.Dir(path))] = true
		}
		return nil
	})
	pkgs := make([]string, 0, len(seen))
	for dir := range seen {
		pkgs = append(pkgs, dir)
	}
	sort.Strings(pkgs)
	return pkgs
}

// TestDataHandlingDerivationBites is the recorded bite proof for
// checkDataHandlingDerivation (per craft-gate): a doc whose region lists every
// env.AgentPasslist pattern passes clean; dropping one pattern fires a diagnostic
// naming exactly that pattern; and removing the marked region fails loudly rather
// than passing vacuously.
func TestDataHandlingDerivationBites(t *testing.T) {
	writeDoc := func(t *testing.T, body string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "DATA_HANDLING.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}
	region := func(patterns []string) string {
		var b strings.Builder
		b.WriteString(passlistBegin + "\n\n| Pattern | Family | Doc |\n|---|---|---|\n")
		for _, p := range patterns {
			fmt.Fprintf(&b, "| `%s` | fam | doc |\n", p)
		}
		b.WriteString("\n" + passlistEnd + "\n")
		return b.String()
	}

	if diags := checkDataHandlingDerivation(writeDoc(t, region(env.AgentPasslist))); len(diags) != 0 {
		t.Fatalf("complete listing: want no diagnostics, got %v", diags)
	}

	dropped := env.AgentPasslist[len(env.AgentPasslist)-1]
	partial := append([]string(nil), env.AgentPasslist[:len(env.AgentPasslist)-1]...)
	diags := checkDataHandlingDerivation(writeDoc(t, region(partial)))
	want := fmt.Sprintf("internal/env pattern %q is not documented", dropped)
	if !containsDiagnostic(diags, want) {
		t.Fatalf("dropped pattern %q: want diagnostic %q, got %v", dropped, want, diags)
	}
	for _, present := range partial {
		if containsDiagnostic(diags, fmt.Sprintf("pattern %q is not documented", present)) {
			t.Fatalf("documented pattern %q wrongly reported missing: %v", present, diags)
		}
	}

	if diags := checkDataHandlingDerivation(writeDoc(t, "# no region here\n")); !containsDiagnostic(diags, "passlist derivation region missing") {
		t.Fatalf("absent region: want a loud region-missing diagnostic, got %v", diags)
	}
	emptyRegion := passlistBegin + "\n\n" + passlistEnd + "\n"
	if diags := checkDataHandlingDerivation(writeDoc(t, emptyRegion)); !containsDiagnostic(diags, "passlist derivation region empty") {
		t.Fatalf("empty region: want a loud region-empty diagnostic, got %v", diags)
	}
}

// TestDataHandlingDerivationFixtureBite is the canary-fixture bite proof for
// story 18: the fixture omits one passlist pattern from an otherwise-complete
// listing, and the derivation check fires that pattern's own diagnostic under a
// full RunConformance pass.
func TestDataHandlingDerivationFixtureBite(t *testing.T) {
	fixture := "undocumented-passlist-var"
	root := materializeConformanceFixture(t, fixture)
	h := NewHarness(t)
	expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

	diags := RunConformance(root, h.KitRoot)

	if !containsDiagnostic(diags, expect) {
		t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
	}
}

// TestSingleControlEscaperBites is the recorded bite proof for
// checkSingleControlEscaper (per craft-gate): a tree with one escaper package
// passes; a second package that grows its own escaper fires a diagnostic naming
// it; a package that only *detects* control bytes does not count as an escaper;
// and a tree with no escaper at all fails loudly.
func TestSingleControlEscaperBites(t *testing.T) {
	// escaperBody carries the escape-emission idiom (the \uXXXX format verb); it
	// must not be passed through fmt.Sprintf, whose verbs would consume the %04x.
	escaperBody := "\n\nimport \"fmt\"\n\n// escape renders a control rune in readable form.\nfunc escape(r rune) string {\n\treturn fmt.Sprintf(`" + `\u%04x` + "`, r)\n}\n"
	detectionBody := "\n\n// hasControl reports whether s carries a control byte; it never escapes.\nfunc hasControl(s string) bool {\n\tfor i := 0; i < len(s); i++ {\n\t\tif s[i] < 0x20 {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n"

	writePkg := func(t *testing.T, root, rel, body string) {
		t.Helper()
		dir := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		pkg := filepath.Base(rel)
		if err := os.WriteFile(filepath.Join(dir, pkg+".go"), []byte("package "+pkg+body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Sole escaper plus a detector and an excluded toon escaper: exactly one owner.
	single := t.TempDir()
	writePkg(t, single, "internal/sanitize", escaperBody)
	writePkg(t, single, "internal/shift", detectionBody)
	writePkg(t, single, "internal/toon", escaperBody)
	if diags := checkSingleControlEscaper(single); len(diags) != 0 {
		t.Fatalf("single escaper: want no diagnostics, got %v", diags)
	}

	// A second package grows its own escaper: red, naming the extra owner.
	dupe := t.TempDir()
	writePkg(t, dupe, "internal/sanitize", escaperBody)
	writePkg(t, dupe, "internal/dashboard", escaperBody)
	diags := checkSingleControlEscaper(dupe)
	if !containsDiagnostic(diags, "more than one owner") || !containsDiagnostic(diags, "internal/dashboard") {
		t.Fatalf("second escaper: want a more-than-one-owner diagnostic naming internal/dashboard, got %v", diags)
	}

	// No escaper anywhere: loud, not vacuous.
	none := t.TempDir()
	writePkg(t, none, "internal/shift", detectionBody)
	if diags := checkSingleControlEscaper(none); !containsDiagnostic(diags, "has no owner") {
		t.Fatalf("no escaper: want a no-owner diagnostic, got %v", diags)
	}
}
