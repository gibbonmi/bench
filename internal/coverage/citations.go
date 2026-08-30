package coverage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// citationRe matches one seam-cell citation: a backticked repo-relative `_test.go`
// path immediately followed by a parenthesized list. The list holds the backticked
// names that file must declare. A backticked test path with no such list is a
// mention, not a citation, so it grades nothing.
var citationRe = regexp.MustCompile("`([^`]*_test\\.go)`[ \t]*\\(([^)]*)\\)")

// citedNameRe pulls each backticked name out of a citation's parenthesized list.
var citedNameRe = regexp.MustCompile("`([^`]+)`")

// reviewPickup renders the review file a folder spec's fences must authorize. The
// landing carries the review beside the build, so an unauthorized pickup blocks it.
func reviewPickup(slug string) string { return "reviews/" + slug + ".md" }

// CheckFiles is Check plus the two checks that read the tree the spec sits in: every
// seam-cell test citation must resolve to a declared function, and a declared
// ownership-fence section must authorize the spec's own review pickup. The repo root
// and the slug both derive from specPath, so no caller can supply an anchor that
// disagrees with the file it passed.
//
// A spec that is not mapped — historical, or missing its map — gets Check's answer
// alone, which keeps the historical opt-out covering every check.
func CheckFiles(p parsed, specPath string) []string {
	v := Check(p)
	if State(p) != "mapped" {
		return v
	}
	base, slug := specLocation(specPath)
	v = append(v, checkCitations(p, base)...)
	return append(v, checkReviewPickup(p, slug)...)
}

// specLocation splits a spec path into the repo root its citations resolve against
// and the slug its review pickup is named for. A folder spec, `<base>/specs/<slug>/
// spec.md`, yields both. Any other shape yields the directory that holds the spec and
// an empty slug, which leaves the pickup check inapplicable rather than guessing a
// name.
func specLocation(specPath string) (base, slug string) {
	if specPath == "" {
		return "", ""
	}
	clean := filepath.Clean(specPath)
	dir := filepath.Dir(clean)
	if filepath.Base(clean) == "spec.md" && filepath.Base(filepath.Dir(dir)) == "specs" {
		return filepath.Dir(filepath.Dir(dir)), filepath.Base(dir)
	}
	if filepath.Base(dir) == "specs" {
		return filepath.Dir(dir), ""
	}
	return dir, ""
}

// checkCitations resolves every seam-cell citation against the tree at base. It holds
// no state between runs, so a renamed test is resolved again on the next call rather
// than trusted from a recorded verdict.
func checkCitations(p parsed, base string) []string {
	var v []string
	s := p.projection()
	for idx, r := range p.dataRows {
		rn := idx + 1
		for _, m := range citationRe.FindAllStringSubmatch(s.cell(r, fieldSeam), -1) {
			v = append(v, checkCitation(rn, base, m[1], m[2])...)
		}
	}
	return v
}

// checkCitation grades one citation: the file at rel must exist under base and must
// declare each name the list cites. A missing file reports once, because every name
// under it fails for the same reason.
func checkCitation(rn int, base, rel, list string) []string {
	names := citedNameRe.FindAllStringSubmatch(list, -1)
	if len(names) == 0 {
		return nil
	}
	content, err := os.ReadFile(filepath.Join(base, filepath.FromSlash(rel)))
	if err != nil {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', which does not exist", rn, rel)}
	}
	var v []string
	for _, n := range names {
		// A subtest citation names its parent function before the slash. The function is
		// what the file declares, so that leading segment is what resolves.
		name := strings.TrimSpace(n[1])
		if i := strings.IndexByte(name, '/'); i >= 0 {
			name = name[:i]
		}
		if name == "" || !declaresFunc(content, name) {
			v = append(v, fmt.Sprintf("coverage map row %d cites '%s', which '%s' does not declare", rn, n[1], rel))
		}
	}
	return v
}

// declaresFunc reports whether content declares a function of this name, with or
// without a receiver. A citation resolves against the declaration, never against a
// call site, so a deleted test whose name still appears in a comment stays red.
func declaresFunc(content []byte, name string) bool {
	re, err := regexp.Compile(`(?m)^func (?:\([^)]*\)[ \t]*)?` + regexp.QuoteMeta(name) + `\(`)
	if err != nil {
		return false
	}
	return re.Match(content)
}

// checkReviewPickup requires the review pickup among the fence tokens of a folder
// spec that declares an ownership-fence section. A spec with no such section, and a
// path that names no slug, are both inapplicable rather than red.
//
// The fence parser ends the section at any level-2-or-deeper heading, so an entry
// written under a subsection is outside the section and does not satisfy this check.
func checkReviewPickup(p parsed, slug string) []string {
	if slug == "" || !p.fencesDeclared {
		return nil
	}
	want := reviewPickup(slug)
	for _, token := range p.fences {
		if token == want {
			return nil
		}
	}
	return []string{fmt.Sprintf("coverage map spec declares ownership fences without the review pickup '%s'", want)}
}
