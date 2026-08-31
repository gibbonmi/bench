package coverage

import (
	"fmt"
	"go/build"
	"go/build/constraint"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate"
)

// citationRe matches one seam-cell test token: a backticked repo-relative `_test.go`
// path, and the parenthesized list that may follow it. The list holds the backticked
// names that file must declare. The list is optional in the pattern so that one
// grammar sees both shapes: a token with no list is a mention, which claims evidence
// it never names, and the mention rule grades it.
var citationRe = regexp.MustCompile("`([^`]*_test\\.go)`[ \t]*(?:\\(([^)]*)\\))?")

// citedNameRe pulls each backticked name out of a citation's parenthesized list.
var citedNameRe = regexp.MustCompile("`([^`]+)`")

// runNameRe matches one t.Run call up to its first argument, so the argument itself
// decides whether the subtest name is a literal the grammar can read. The receiver is
// anchored to t, the name the toolchain's own template and every test in this repo give
// the *testing.T parameter. Any other Run call — cmd.Run, a server's Run — names no
// subtest, so it neither adds a name nor exempts the file.
var runNameRe = regexp.MustCompile(`\bt\.Run\([ \t\n]*`)

// citation is one seam-cell test token the grammar reads: the repo-relative path, and
// the parenthesized name list when the token carries one. A token with no list is a
// mention, which claims evidence it never names.
type citation struct {
	rel     string
	hasList bool
	list    string
}

// namesSomething reports whether the citation's list names at least one thing. It is the
// one derivation of that rule, so the report and the checks cannot disagree on which
// cell holds evidence.
func (c citation) namesSomething() bool { return c.hasList && citedNameRe.MatchString(c.list) }

// citationsIn projects one seam cell into the citations it holds.
func citationsIn(cell string) []citation {
	var out []citation
	for _, m := range citationRe.FindAllStringSubmatchIndex(cell, -1) {
		c := citation{rel: cell[m[2]:m[3]], hasList: m[4] >= 0}
		if c.hasList {
			c.list = cell[m[4]:m[5]]
		}
		out = append(out, c)
	}
	return out
}

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
	census := executedCensus(base)
	for idx, r := range p.dataRows {
		rn := idx + 1
		// Only the seam cell holds the row's evidence, so only it is graded. A test path
		// in the why cell is prose about the failure, not a claim of coverage.
		cell := s.cell(r, fieldSeam)
		for _, c := range citationsIn(cell) {
			if !c.hasList {
				v = append(v, fmt.Sprintf("coverage map row %d mentions '%s' without a cited name list", rn, c.rel))
				continue
			}
			v = append(v, checkCitation(rn, base, c.rel, c.list, census)...)
		}
	}
	return v
}

// executedCensus answers which build-tag sets the gate compiles the tree at base with.
// The gate package owns that derivation, so the census and the oracle cannot disagree,
// and the kit resolves through the gate's own rule so the cited files and the census
// come from one tree.
//
// A census the gate cannot derive leaves the execution check inapplicable rather than
// red. A root with no Go module, and a root whose phase manifest is a defect, are both
// answered by the gate elsewhere; a coverage check is not the surface that reports them.
func executedCensus(base string) []gate.TagSet {
	if base == "" {
		return nil
	}
	census, err := gate.ExecutedTagCensus(base, gate.KitRoot(base))
	if err != nil {
		return nil
	}
	return census
}

// checkCitation grades one citation: the file at rel must exist under base and must
// declare each name the list cites. A missing file reports once, because every name
// under it fails for the same reason.
func checkCitation(rn int, base, rel, list string, census []gate.TagSet) []string {
	names := citedNameRe.FindAllStringSubmatch(list, -1)
	if len(names) == 0 {
		return nil
	}
	if !citedPathSafe(rel) {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', which is not a path inside the repo", rn, rel)}
	}
	path := filepath.Join(base, filepath.FromSlash(rel))
	// The path is classified before anything opens it. A FIFO planted at a cited path
	// would otherwise block the read forever, inside the oracle.
	if bounds.ClassifyNoFollow(path).State == bounds.StateWrongType {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', which is not a regular file", rn, rel)}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', which does not exist", rn, rel)}
	}
	var v []string
	literal, subtests := runNames(content)
	for _, n := range names {
		// A subtest citation names its parent function before the slash. The function is
		// what the file declares, so that leading segment is what resolves.
		name := strings.TrimSpace(n[1])
		segment := ""
		if i := strings.IndexByte(name, '/'); i >= 0 {
			// Only the first segment is graded. A deeper composed name needs the call
			// graph, which this grammar does not read.
			name, segment = name[:i], name[i+1:]
			if j := strings.IndexByte(segment, '/'); j >= 0 {
				segment = segment[:j]
			}
		}
		if name == "" || !declaresFunc(content, name) {
			v = append(v, fmt.Sprintf("coverage map row %d cites '%s', which '%s' does not declare", rn, n[1], rel))
			continue
		}
		// The parent resolves against a declaration; the segment resolves against the
		// t.Run names the file spells out. A file that computes a name spells none, so
		// its segments are unreadable here and stay unjudged.
		if segment == "" || !literal || hasSubtest(subtests, segment) {
			continue
		}
		v = append(v, fmt.Sprintf("coverage map row %d cites '%s', whose subtest '%s' is no t.Run name in '%s'", rn, n[1], segment, rel))
	}
	return append(v, checkExecution(rn, rel, path, content, census)...)
}

// citedPathSafe reports whether rel is a repo-relative forward-slash path that stays
// inside the tree the citation resolves against. An absolute path, a backslash
// separator, or a ".." segment would read outside that tree, so the check refuses it
// before the join rather than resolving it. This is the rule validatePayloadRows
// applies to an allowlist source.
func citedPathSafe(rel string) bool {
	if strings.HasPrefix(rel, "/") || strings.ContainsRune(rel, '\\') {
		return false
	}
	for _, segment := range strings.Split(rel, "/") {
		if segment == ".." {
			return false
		}
	}
	return true
}

// checkExecution grades the cited file against the executed tag census: a declared test
// is evidence only when a gate run on this host actually compiles the file that holds
// it. An empty census means the gate runs no test phase here, so there is no execution
// claim to grade and the check is inapplicable.
//
// One census set satisfying the file is enough, because the gate runs every set.
func checkExecution(rn int, rel, path string, content []byte, census []gate.TagSet) []string {
	if len(census) == 0 {
		return nil
	}
	line, err := buildLine(content)
	if err != nil {
		return []string{fmt.Sprintf("coverage map row %d cites '%s', whose //go:build expression does not parse: %s", rn, rel, err)}
	}
	dir, name := filepath.Split(path)
	for _, set := range census {
		ctxt := contextFor(set)
		if matched, err := ctxt.MatchFile(dir, name); err == nil && matched {
			return nil
		}
	}
	return []string{fmt.Sprintf("coverage map row %d cites '%s', which no executed tag set builds (%s)", rn, rel, constraintOf(line))}
}

// contextFor is the build context one executed tag set compiles in. It starts from the
// toolchain's own default, so the host GOOS and GOARCH, the release tags, and the
// toolchain-implied tags (unix, gc, cgo) all carry their real values rather than a
// second list kept here. MatchFile then applies the two rules that decide reachability:
// the //go:build expression and the GOOS and GOARCH filename suffix rule, which strips
// _test before it reads the suffix.
func contextFor(set gate.TagSet) build.Context {
	ctxt := build.Default
	ctxt.BuildTags = append(append([]string(nil), build.Default.BuildTags...), set...)
	return ctxt
}

// buildLine returns the file's //go:build line, and an error when that line holds an
// expression the toolchain cannot parse. A file with no such line returns the empty
// string, which is the always-satisfied case. The scan stops at the package clause,
// because a later line is a comment rather than a constraint.
func buildLine(content []byte) (string, error) {
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "package ") {
			return "", nil
		}
		if !constraint.IsGoBuild(line) {
			continue
		}
		if _, err := constraint.Parse(line); err != nil {
			return "", err
		}
		return line, nil
	}
	return "", nil
}

// constraintOf names the constraint that kept a file out of every executed set, so the
// repair starts at the constraint rather than at a hand census. A file with no build
// line was refused by its own name.
func constraintOf(line string) string {
	if line == "" {
		return "the filename's GOOS or GOARCH suffix"
	}
	return line
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

// runNames reads the subtest names a file spells out. Only a t.Run call is read, so a
// Run call on any other receiver neither names a subtest nor exempts the file. literal
// is false when a t.Run call names its subtest with an expression rather than a string
// literal, because a table-driven suite names its cases at run time and no scan can see
// them.
func runNames(content []byte) (literal bool, names []string) {
	body := string(content)
	literal = true
	for _, m := range runNameRe.FindAllStringIndex(body, -1) {
		rest := body[m[1]:]
		switch {
		case strings.HasPrefix(rest, `"`):
			name, err := strconv.Unquote(quotedPrefix(rest))
			if err != nil {
				literal = false
				continue
			}
			names = append(names, name)
		case strings.HasPrefix(rest, "`"):
			// A raw literal is as readable as an interpreted one, so it names a subtest
			// rather than exempting the file.
			if i := strings.IndexByte(rest[1:], '`'); i >= 0 {
				names = append(names, rest[1:1+i])
				continue
			}
			literal = false
		default:
			literal = false
		}
	}
	return literal, names
}

// quotedPrefix returns the quoted string literal that starts s, or s itself when the
// literal never closes. An unclosed literal fails to unquote, which marks the file
// unreadable rather than inventing a name for it.
func quotedPrefix(s string) string {
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '"':
			return s[:i+1]
		case '\n':
			return s
		}
	}
	return s
}

// hasSubtest reports whether a cited segment names one of the file's t.Run literals. Go
// replaces a space with an underscore when it runs a subtest, so a citation written in
// either form resolves.
func hasSubtest(names []string, segment string) bool {
	want := strings.ReplaceAll(segment, " ", "_")
	for _, n := range names {
		if strings.ReplaceAll(n, " ", "_") == want {
			return true
		}
	}
	return false
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

// reviewOwnedPrefix marks a seam cell whose evidence a review axis owns rather than a
// test. It is the one marker grammar for that state, matched literally on the lowercase
// form, so a deliberate exclusion reads the same to a human and to this check.
const reviewOwnedPrefix = "review-owned:"

// uncitedRows names each mapped row whose seam cell holds no citation and carries no
// review-owned marker. It reads the same seam-cell projection and the same citation
// grammar the violation checks read, so a cell cannot be a citation to one and prose to
// the other.
//
// An opted-in map names its rows by their row IDs, which survive a reorder. A non-opt-in
// map has only the positional row number to name.
func uncitedRows(p parsed) []string {
	if State(p) != "mapped" {
		return nil
	}
	s := p.projection()
	var names []string
	for idx, r := range p.dataRows {
		cell := strings.TrimSpace(s.cell(r, fieldSeam))
		if strings.HasPrefix(cell, reviewOwnedPrefix) || citesATest(cell) {
			continue
		}
		if s.optIn() {
			names = append(names, s.cell(r, fieldRow))
			continue
		}
		names = append(names, strconv.Itoa(idx+1))
	}
	return names
}

// citesATest reports whether a seam cell holds at least one citation: a test path with a
// name list that names something. A bare path is a mention and an empty list names
// nothing, so neither is evidence, and both are graded elsewhere.
func citesATest(cell string) bool {
	for _, c := range citationsIn(cell) {
		if c.namesSomething() {
			return true
		}
	}
	return false
}
