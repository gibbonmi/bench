package gate

// The executed-tag census. It answers which build-tag sets the gate actually compiles
// with on this host, so a consumer can prove that a named test file is reachable by a
// gate run rather than merely declared. The answer comes from the resolved phase table
// itself; no second copy of the tag list exists.

import (
	"sort"
	"strings"
)

// TagSet is one build-tag set a gate test phase compiles with, sorted and deduplicated.
// The empty set is the untagged default, which the ordinary test phase carries.
type TagSet []string

// ExecutedTagCensus reports every build-tag set the gate executes for the tree under
// grade. root is that tree; kit is the checkout that owns the Go tests, the same pair
// the runner resolves its phase table with. The table is the root's phase manifest when
// the root declares one, else the built-in kit table, so the census and the oracle
// cannot disagree.
//
// Each `go test` phase contributes the set its argv declares, and a test phase with no
// `-tags` operand contributes the untagged set. A table with no test phase yields an
// empty census, which the caller reads as "inapplicable": a root the gate never runs
// tests in has no execution claim to grade.
func ExecutedTagCensus(root, kit string) ([]TagSet, error) {
	phases, err := phaseTable(root, kit)
	if err != nil {
		return nil, err
	}
	var census []TagSet
	seen := make(map[string]bool)
	for _, phase := range phases {
		if !isGoTestArgv(phase.Argv) {
			continue
		}
		set := tagSetOf(phase.Argv)
		key := strings.Join(set, ",")
		if seen[key] {
			continue
		}
		seen[key] = true
		census = append(census, set)
	}
	return census, nil
}

// isGoTestArgv reports whether argv invokes `go test`. The kit's own producer inserts
// `-C <dir>` between the executable and the subcommand, so the scan steps over that
// pair rather than reading a fixed position.
func isGoTestArgv(argv []string) bool {
	if len(argv) == 0 || goExecutable(argv[0]) == "" {
		return false
	}
	for i := 1; i < len(argv); i++ {
		if argv[i] == "-C" {
			i++
			continue
		}
		return argv[i] == "test"
	}
	return false
}

// goExecutable answers the argv[0] forms that name the Go toolchain: the bare name a
// phase table writes, and an absolute path a manifest may declare instead.
func goExecutable(arg string) string {
	if arg == "go" || strings.HasSuffix(arg, "/go") {
		return arg
	}
	return ""
}

// tagSetOf collects the `-tags` operand of one test argv. Both spellings Go accepts are
// read, because a project manifest writes the argv by hand. The returned set is empty
// when the argv carries no operand, which is the untagged default.
func tagSetOf(argv []string) TagSet {
	var tags []string
	for i, arg := range argv {
		operand, ok := tagsOperand(argv, i, arg)
		if !ok {
			continue
		}
		tags = append(tags, splitTags(operand)...)
	}
	return normalizeTags(tags)
}

func tagsOperand(argv []string, i int, arg string) (string, bool) {
	for _, flag := range []string{"-tags", "--tags"} {
		if value, found := strings.CutPrefix(arg, flag+"="); found {
			return value, true
		}
		if arg == flag && i+1 < len(argv) {
			return argv[i+1], true
		}
	}
	return "", false
}

// splitTags accepts the comma form and the deprecated space form of the operand.
func splitTags(operand string) []string {
	return strings.FieldsFunc(operand, func(r rune) bool {
		return r == ',' || r == ' '
	})
}

func normalizeTags(tags []string) TagSet {
	set := make(TagSet, 0, len(tags))
	seen := make(map[string]bool, len(tags))
	for _, tag := range tags {
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		set = append(set, tag)
	}
	sort.Strings(set)
	return set
}
