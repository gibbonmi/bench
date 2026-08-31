package gate

// The executed-test census. It answers which Go test phases the gate actually runs on
// this host, and for each one the build tags it compiles with and the packages it
// selects, so a consumer can prove that a named test file is reachable by a gate run
// rather than merely declared. The answer comes from the resolved phase table itself;
// no second copy of the phase facts exists.

import (
	"path/filepath"
	"sort"
	"strings"
)

// TagSet is one build-tag set a gate test phase compiles with, sorted and deduplicated.
// The empty set is the untagged default, which the ordinary test phase carries.
type TagSet []string

// TestExecution is one resolved Go test phase. Its package operands and tags stay
// together, so one phase cannot authorize a file with another phase's facts.
type TestExecution struct {
	Name     string
	Tags     TagSet
	Packages []string
	Dir      string
	GoC      string
}

// ExecutedTestCensus reports every Go test phase the gate executes for the tree under
// grade. root is that tree; kit is the checkout that owns the Go tests, the same pair
// the runner resolves its phase table with. The table is the root's phase manifest when
// the root declares one, else the built-in kit table, so the census and the oracle
// cannot disagree.
//
// A table with no test phase yields an empty census, which the caller reads as
// "inapplicable": a root the gate never runs tests in has no execution claim to grade.
func ExecutedTestCensus(root, kit string) ([]TestExecution, error) {
	phases, err := phaseTable(root, kit)
	if err != nil {
		return nil, err
	}
	var census []TestExecution
	for _, phase := range phases {
		testAt := goTestIndex(phase.Argv)
		if testAt < 0 {
			continue
		}
		dir := phase.Dir
		if dir == "" {
			dir = root
		}
		goC := goCOf(phase.Argv, testAt)
		effectiveDir := dir
		if goC != "" {
			effectiveDir = goC
			if !filepath.IsAbs(effectiveDir) {
				effectiveDir = filepath.Join(dir, effectiveDir)
			}
			goC = filepath.Clean(effectiveDir)
		}
		census = append(census, TestExecution{
			Name:     phase.Name,
			Tags:     tagSetOf(phase.Argv),
			Packages: testPackagesOf(phase.Argv, testAt),
			Dir:      filepath.Clean(effectiveDir),
			GoC:      goC,
		})
	}
	return census, nil
}

// goTestIndex reports where the `test` subcommand sits in an argv that invokes
// `go test`, and -1 for any other argv. The kit's own producer inserts `-C <dir>`
// between the executable and the subcommand, so the scan steps over that pair rather
// than reading a fixed position. The index is what the package-operand reader needs,
// because every operand follows the subcommand.
func goTestIndex(argv []string) int {
	if len(argv) == 0 || goExecutable(argv[0]) == "" {
		return -1
	}
	for i := 1; i < len(argv); i++ {
		if argv[i] == "-C" {
			i++
			continue
		}
		if argv[i] == "test" {
			return i
		}
		return -1
	}
	return -1
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

// goCOf returns the Go process directory setting. It stays attached to the phase because
// Go resolves a relative package operand after it changes to that directory.
func goCOf(argv []string, testAt int) string {
	for i := 1; i < testAt; i++ {
		if argv[i] == "-C" && i+1 < testAt {
			return argv[i+1]
		}
	}
	return ""
}

// testPackagesOf reads package operands from one Go test argv. Test-binary arguments
// begin at -args or at a -test flag, so neither can become package evidence.
func testPackagesOf(argv []string, testAt int) []string {
	var packages []string
	for i := testAt + 1; i < len(argv); i++ {
		arg := argv[i]
		if arg == "-args" || strings.HasPrefix(arg, "-test.") {
			break
		}
		if strings.HasPrefix(arg, "-") {
			if testFlagTakesValue(arg) {
				i++
			}
			continue
		}
		packages = append(packages, arg)
	}
	return packages
}

func testFlagTakesValue(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}
	switch arg {
	case "-tags", "--tags", "-mod", "-modfile", "-overlay", "-p", "-pkgdir", "-covermode",
		"-coverpkg", "-coverprofile", "-vet", "-run", "-bench", "-benchtime", "-count", "-cpu",
		"-parallel", "-shuffle", "-timeout":
		return true
	default:
		return false
	}
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

// BaselinePolicyEnv names the environment variable that points the phase schedule at a
// landing's baseline rather than at the tree under grade. A census consumer's fixture
// must neutralize it, or a prospective landing gate answers every synthetic root with
// the baseline's phase table instead of the root's own.
const BaselinePolicyEnv = baselinePolicyEnv

// KitRoot resolves the checkout that owns the Go tests for a tree under grade, through
// the same rule every gate entry point uses. A census consumer outside this package
// must resolve the pair the same way, or it grades cited files in one tree against a
// census taken in another.
func KitRoot(root string) string { return kitRoot(root) }
