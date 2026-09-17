package preflight

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// preparationRefusalCase is one enumerated bounded build preparation refusal, checked
// against the fixed baseline under testdata/legacy-charge. The permanent test compares a
// live run with that stored capture and never rewrites it.
type preparationRefusalCase struct {
	name string
	run  func(t *testing.T) (out string, code int, root, base, tip string)
}

func preparationRefusalRun(t *testing.T, root string, args []string) (string, int, string, string, string) {
	t.Helper()
	out, code := Command(args)
	return out, code, root, args[6], args[8]
}

func replacePhase(t *testing.T, prepare func(t *testing.T)) {
	t.Helper()
	if err := os.Remove(chargesource.BuildPhase); err != nil {
		t.Fatal(err)
	}
	prepare(t)
}

// requiredSourceMutation is one required-source or preparation-refusal setup, shared by
// the enumerated legacy differential and the classified evidence checks in
// charge_evidence_test.go. Its mutate func applies the checkout or argument change and
// returns the resulting arguments; a case that only touches the checkout ignores args.
type requiredSourceMutation struct {
	name   string
	mutate func(t *testing.T, root, slug string, args []string) []string
}

func requiredSourceMutations() []requiredSourceMutation {
	return []requiredSourceMutation{
		{"absent phase", func(t *testing.T, root, slug string, args []string) []string {
			replacePhase(t, func(*testing.T) {})
			return args
		}},
		{"empty phase", func(t *testing.T, root, slug string, args []string) []string {
			preflighttest.MustWriteFile(t, chargesource.BuildPhase, "")
			return args
		}},
		{"symlink phase", func(t *testing.T, root, slug string, args []string) []string {
			replacePhase(t, func(t *testing.T) {
				if err := os.Symlink("../skills/bench-craft-delegate/SKILL.md", chargesource.BuildPhase); err != nil {
					t.Fatal(err)
				}
			})
			return args
		}},
		{"directory phase", func(t *testing.T, root, slug string, args []string) []string {
			replacePhase(t, func(t *testing.T) {
				preflighttest.MustWriteFile(t, filepath.Join(chargesource.BuildPhase, "inner.md"), "# Inner\n")
			})
			return args
		}},
		{"control byte phase", func(t *testing.T, root, slug string, args []string) []string {
			preflighttest.MustWriteFile(t, chargesource.BuildPhase, "unsafe \x1b source\n")
			return args
		}},
		{"dirty checkout", func(t *testing.T, root, slug string, args []string) []string {
			preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n// dirty\n")
			return args
		}},
		{"missing ticket", func(t *testing.T, root, slug string, args []string) []string {
			args[4] = "missing.md"
			return args
		}},
		{"source-tip mismatch", func(t *testing.T, root, slug string, args []string) []string {
			args[8] = args[6]
			return args
		}},
		{"no assignment", func(t *testing.T, root, slug string, args []string) []string {
			return []string{"build", slug, "--charge", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
		}},
		{"foreign assignment", func(t *testing.T, root, slug string, args []string) []string {
			preflighttest.ActiveAssignment(t, root, t.TempDir())
			return args
		}},
	}
}

// mutationNamed looks up one requiredSourceMutations case by name.
func mutationNamed(t *testing.T, name string) func(t *testing.T, root, slug string, args []string) []string {
	t.Helper()
	for _, m := range requiredSourceMutations() {
		if m.name == name {
			return m.mutate
		}
	}
	t.Fatalf("no required-source mutation %q", name)
	return nil
}

func preparationRefusalCases() []preparationRefusalCase {
	return []preparationRefusalCase{
		{"refusal-absent-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			mutationNamed(t, "absent phase")(t, root, slug, nil)
			return preparationRefusalRun(t, root, preflighttest.LegacyCommitted(t, root, slug, "absent phase"))
		}},
		{"refusal-empty-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			mutationNamed(t, "empty phase")(t, root, slug, nil)
			return preparationRefusalRun(t, root, preflighttest.LegacyCommitted(t, root, slug, "empty phase"))
		}},
		{"refusal-symlink-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			mutationNamed(t, "symlink phase")(t, root, slug, nil)
			return preparationRefusalRun(t, root, preflighttest.LegacyCommitted(t, root, slug, "linked phase"))
		}},
		{"refusal-directory-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			mutationNamed(t, "directory phase")(t, root, slug, nil)
			return preparationRefusalRun(t, root, preflighttest.LegacyCommitted(t, root, slug, "directory phase"))
		}},
		{"refusal-control-byte-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			mutationNamed(t, "control byte phase")(t, root, slug, nil)
			return preparationRefusalRun(t, root, preflighttest.LegacyCommitted(t, root, slug, "control phase"))
		}},
		{"refusal-dirty-checkout", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := mutationNamed(t, "dirty checkout")(t, root, slug, preflighttest.ChargeArgs(t, root, slug))
			return preparationRefusalRun(t, root, args)
		}},
		{"refusal-missing-ticket", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := mutationNamed(t, "missing ticket")(t, root, slug, preflighttest.ChargeArgs(t, root, slug))
			return preparationRefusalRun(t, root, args)
		}},
		{"refusal-source-tip-mismatch", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := mutationNamed(t, "source-tip mismatch")(t, root, slug, preflighttest.ChargeArgs(t, root, slug))
			out, code := Command(args)
			return out, code, root, args[6], preflighttest.RunGit(t, "rev-parse", "HEAD")
		}},
		{"refusal-inactive-assignment", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := mutationNamed(t, "no assignment")(t, root, slug, nil)
			return preparationRefusalRun(t, root, args)
		}},
		{"refusal-foreign-assignment", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := mutationNamed(t, "foreign assignment")(t, root, slug, preflighttest.ChargeArgs(t, root, slug))
			return preparationRefusalRun(t, root, args)
		}},
		{"refusal-persistent-movement", func(t *testing.T) (string, int, string, string, string) {
			root, slug := preflighttest.SeedConformant(t)
			args := preflighttest.ChargeArgs(t, root, slug)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+strconv.Itoa(calls)+"\n")
			})
			defer restore()
			return preparationRefusalRun(t, root, args)
		}},
	}
}

func legacyChargeRecord(out string, code int) string {
	return fmt.Sprintf("exit: %d\n%s", code, out)
}

// TestLegacyPreparedPackDifferential is CE147. Every enumerated build charge response
// keeps its captured exit and bytes after the charge reads its sources and metadata back
// from the validated in-memory pack.
func TestLegacyPreparedPackDifferential(t *testing.T) {
	sourceDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range preparationRefusalCases() {
		t.Run(test.name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join(sourceDir, "testdata", "legacy-charge", test.name+".toon"))
			if err != nil {
				t.Fatal(err)
			}
			out, code, root, base, tip := test.run(t)
			got := legacyChargeRecord(normalizeLegacy(out, root, base, tip), code)
			if got != string(want) {
				t.Fatalf("legacy charge %s changed:\ngot:\n%s\nwant:\n%s", test.name, got, want)
			}
		})
	}
}
