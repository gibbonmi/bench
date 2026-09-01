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

// TestExecution is one resolved Go test phase. Its package operands, tags, and
// environment stay together, so one phase cannot authorize a file with another phase's
// facts. Env carries the phase's own KEY=VALUE overrides, which decide what the
// toolchain selects and compiles as surely as the argv does.
type TestExecution struct {
	Name     string
	Tags     TagSet
	Packages []string
	Dir      string
	GoC      string
	Env      []string
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
		testAt, goC := goTestIndex(phase.Argv)
		if testAt < 0 {
			continue
		}
		dir := phase.Dir
		if dir == "" {
			dir = root
		}
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
			Env:      phase.Env,
		})
	}
	return census, nil
}

// goTestIndex reports where the `test` subcommand sits in an argv that invokes
// `go test`, and -1 for any other argv. The kit's own producer inserts `-C <dir>`
// between the executable and the subcommand, so the scan steps over that pair rather
// than reading a fixed position. The index is what the package-operand reader needs,
// because every operand follows the subcommand.
//
// The same walk answers the Go process directory setting, because both facts live in the
// one region between the executable and the subcommand. That setting stays attached to
// the phase: Go resolves a relative package operand after it changes to that directory.
func goTestIndex(argv []string) (testAt int, goC string) {
	if len(argv) == 0 || goExecutable(argv[0]) == "" {
		return -1, ""
	}
	for i := 1; i < len(argv); i++ {
		switch {
		case argv[i] == "-C":
			// Go accepts -C only as the first argument, so a second occurrence makes the
			// phase fail on its own terms. The first is the one this scan reports.
			if goC == "" && i+1 < len(argv) {
				goC = argv[i+1]
			}
			i++
		case argv[i] == "test":
			return i, goC
		default:
			return -1, ""
		}
	}
	return -1, ""
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

// testPackagesOf reads package operands from one Go test argv. Only -args ends the
// operand region, because every argument after it reaches the test binary
// uninterpreted. A -test.-prefixed argument is an ordinary flag the toolchain forwards,
// so the scan steps over it and keeps reading the operands behind it.
func testPackagesOf(argv []string, testAt int) []string {
	var packages []string
	for i := testAt + 1; i < len(argv); i++ {
		arg := argv[i]
		if arg == "-args" {
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

// testFlagTakesValue reports whether a `go test` flag consumes the argument after it.
// The set is every value-taking flag `go help build` and `go help testflag` document for
// `go test`, less the ones that carry no separated value: a boolean flag, and any flag
// written in the `-flag=value` form. An omission here reads a flag's own value as a
// package operand, which makes `go list` reject the phase and false-reds every citation
// that phase could have supplied.
//
// `-C` is absent because Go accepts it only as the first argument of the whole command
// line, which is ahead of the subcommand this scan starts from. `-gocoverdir` is absent
// because cmd/go withholds it from the end-user flag set.
func testFlagTakesValue(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}
	// Go's flag package accepts one or two leading dashes for the same flag, so the bare
	// name is what this table keys on.
	switch strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-") {
	// Build flags, which `go test` shares with `go build`.
	case "asmflags", "buildmode", "compiler", "covermode", "coverpkg", "gccgoflags",
		"gcflags", "installsuffix", "ldflags", "mod", "modfile", "overlay", "p", "pgo",
		"pkgdir", "tags", "toolexec":
		return true
	// Flags `go test` itself handles.
	case "coverprofile", "exec", "o", "vet":
		return true
	// Flags `go test` forwards to the test binary.
	case "bench", "benchtime", "blockprofile", "blockprofilerate", "count", "cpu",
		"cpuprofile", "fuzz", "fuzzminimizetime", "fuzztime", "list", "memprofile",
		"memprofilerate", "mutexprofile", "mutexprofilefraction", "outputdir", "parallel",
		"run", "shuffle", "skip", "timeout", "trace":
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
