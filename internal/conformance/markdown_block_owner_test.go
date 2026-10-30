package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/markdown"
)

// markdownBlockOwner is the block reader, the one package that may spell a block-rule literal.
const markdownBlockOwner = "internal/markdown"

// blockRuleKinds are the block-rule literals that production Go may not spell. A fence run
// refuses a literal that contains it anywhere, so a detector cannot hide the run after other
// text. Each other kind refuses only a literal that equals it, so a whole heading title or a
// managed-block marker stays legal.
var blockRuleKinds = []struct {
	name, value string
	anywhere    bool
}{
	{"backtick fence run", markdown.FenceMarker, true},
	{"tilde fence run", "~~~", true},
	{"frontmatter marker", "---", false},
	{"comment opener", "<!--", false},
	{"comment closer", "-->", false},
	{"H2 prefix", markdown.H2Prefix, false},
}

// checkMarkdownBlockOwner refuses each string literal of a non-test Go file of the module,
// outside the block reader, that spells a block rule. Such a literal is a second detector of a
// fence, a frontmatter block, a comment, or an H2 heading, and it drifts from the reader. Test
// files are exempt, so a fixture can plant a block. A rune literal is not a string literal, so
// the check does not read it.
func checkMarkdownBlockOwner(root string) []string {
	var diags []string
	for _, path := range moduleGoFiles(root) {
		if rel := slashRel(root, path); !strings.HasSuffix(rel, "_test.go") && !strings.HasPrefix(rel, markdownBlockOwner+"/") {
			diags = append(diags, blockRuleLiterals(path, rel)...)
		}
	}
	return diags
}

// blockRuleLiterals gives one diagnostic for each string literal of the file at path that
// spells a block rule, and it names the first kind that the literal matches. A file that does
// not parse gives one diagnostic, so it cannot pass unread.
func blockRuleLiterals(path, rel string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, readIfExists(path), 0)
	if err != nil {
		return []string{rel + " cannot be parsed for markdown-block ownership: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok {
			return true
		}
		value, ok := stringLiteralValue(literal)
		if !ok {
			return true
		}
		for _, kind := range blockRuleKinds {
			if value == kind.value || kind.anywhere && strings.Contains(value, kind.value) {
				diags = append(diags, fmt.Sprintf("%s:%d spells the %s outside %s; read or render the block through the block reader", rel, fset.Position(literal.Pos()).Line, kind.name, markdownBlockOwner))
				break
			}
		}
		return true
	})
	return diags
}

// blockRuleKindsIndependentExpectation is this test's own list of the block-rule kinds, in the
// order that blockRuleProduction plants them. It stays independent of blockRuleKinds, so a kind
// that the check omits still reds the bite proof.
var blockRuleKindsIndependentExpectation = []string{"backtick fence run", "tilde fence run", "frontmatter marker", "comment opener", "comment closer", "H2 prefix"}

// blockRuleProduction plants one literal of each block-rule kind, one on each line from line 4.
// The two fence runs sit after other text, so a check that reads only a prefix misses them.
const blockRuleProduction = "package example\n\nvar (\n" +
	"\tbacktick = \"text ```\"\n" +
	"\ttilde = []rune(\"text ~~~\")\n" +
	"\tfrontmatter = \"---\"\n" +
	"\topener = `<!--`\n" +
	"\tcloser = \"-->\"\n" +
	"\theading = \"## \"\n" +
	")\n"

// blockRuleNearMisses holds the legal neighbors of each kind: a short run, a whole title, a
// padded marker, a managed-block marker, and a rune literal.
const blockRuleNearMisses = "package example\n\nvar (\n" +
	"\tshort = \"``\"\n" +
	"\ttitle = \"## Recommended sequence\"\n" +
	"\tpadded = \"--- \"\n" +
	"\tmarker = \"<!-- bench:start -->\"\n" +
	"\ttick = '~'\n" +
	")\n"

// TestMarkdownBlockOwnerBites is the bite proof for checkMarkdownBlockOwner. Each failure
// message names the case and the diagnostics, never the planted literal.
func TestMarkdownBlockOwnerBites(t *testing.T) {
	t.Run("each kind in internal", func(t *testing.T) {
		const rel = "internal/example/example.go"
		diags := checkMarkdownBlockOwner(throwawayRoot{files: map[string]string{rel: blockRuleProduction}}.build(t))
		if len(diags) != len(blockRuleKindsIndependentExpectation) {
			t.Fatalf("diagnostics = %q, want one for each of the %d kinds", diags, len(blockRuleKindsIndependentExpectation))
		}
		for i, kind := range blockRuleKindsIndependentExpectation {
			if want := fmt.Sprintf("%s:%d ", rel, i+4); !strings.HasPrefix(diags[i], want) || !strings.Contains(diags[i], " "+kind+" ") {
				t.Errorf("diagnostic %d = %q, want it to name %s and the %s", i, diags[i], want, kind)
			}
		}
	})
	t.Run("two kinds in one cmd literal", func(t *testing.T) {
		const rel = "cmd/example/main.go"
		diags := checkMarkdownBlockOwner(throwawayRoot{files: map[string]string{rel: "package main\n\nconst both = \"~~~ and ```\"\n"}}.build(t))
		if len(diags) != 1 || !strings.HasPrefix(diags[0], rel+":3 ") || !strings.Contains(diags[0], " backtick fence run ") {
			t.Fatalf("diagnostics = %q, want one that names %s:3 and the backtick fence run", diags, rel)
		}
	})
	t.Run("unparseable file", func(t *testing.T) {
		diags := checkMarkdownBlockOwner(throwawayRoot{files: map[string]string{"internal/example/example.go": "package example\n\nvar (\n"}}.build(t))
		if len(diags) != 1 || !strings.Contains(diags[0], "internal/example/example.go cannot be parsed") {
			t.Fatalf("diagnostics = %q, want one parse diagnostic", diags)
		}
	})
	t.Run("near misses", func(t *testing.T) {
		if diags := checkMarkdownBlockOwner(throwawayRoot{files: map[string]string{"internal/example/example.go": blockRuleNearMisses}}.build(t)); len(diags) != 0 {
			t.Fatalf("diagnostics = %q, want none", diags)
		}
	})
	t.Run("test file and owner exempt", func(t *testing.T) {
		root := throwawayRoot{files: map[string]string{
			"internal/example/example_test.go": blockRuleProduction,
			markdownBlockOwner + "/extra.go":   blockRuleProduction,
		}}.build(t)
		if diags := checkMarkdownBlockOwner(root); len(diags) != 0 {
			t.Fatalf("diagnostics = %q, want none", diags)
		}
	})
}

// TestMarkdownBlockOwnerHoldsOnTheLiveTree is the check's live-tree assertion. No production Go
// file outside the block reader spells a block-rule literal.
func TestMarkdownBlockOwnerHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if diags := checkMarkdownBlockOwner(h.KitRoot); len(diags) != 0 {
		t.Fatalf("a block-rule literal survives outside the block reader:\n%s", strings.Join(diags, "\n"))
	}
}

// TestConformanceReadersUseTheBlockReader pins the block rules that three conformance readers
// take from the block reader: a CRLF frontmatter block opens, and a fenced or commented heading
// is no heading.
func TestConformanceReadersUseTheBlockReader(t *testing.T) {
	const crlf = "---\r\nname: x\r\ndescription: a\r\n  folded\r\ndisable-model-invocation: true\r\n---\r\n"
	if !skillDescriptionFolded(crlf) {
		t.Error("skillDescriptionFolded misses the folded description of a CRLF frontmatter block")
	}
	if !frontmatterHasKey(crlf, "disable-model-invocation") {
		t.Error("frontmatterHasKey misses the key of a CRLF frontmatter block")
	}
	heads := h2Headings("## Kept\n" + markdown.FenceMarker + "\n## Fenced\n" + markdown.FenceMarker + "\n<!--\n## Commented\n-->\n")
	if len(heads) != 1 || !heads["## Kept"] {
		t.Errorf("h2Headings reads %d headings, want only the heading outside the fence and the comment", len(heads))
	}
}
