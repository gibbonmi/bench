package preflight

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
)

// legacyChargeCase is one enumerated legacy build charge response, checked against the
// fixed baseline under testdata/legacy-charge. The permanent test compares a live run
// with that stored capture and never rewrites it.
type legacyChargeCase struct {
	name string
	run  func(t *testing.T) (out string, code int, root, base, tip string)
}

func legacyChargeRun(t *testing.T, root string, args []string) (string, int, string, string, string) {
	t.Helper()
	out, code := Command(args)
	return out, code, root, args[6], args[8]
}

// legacyCommitted commits every change and repins the charge arguments to the new tip.
func legacyCommitted(t *testing.T, root, slug, message string, full bool) []string {
	t.Helper()
	runGit(t, "add", "-A")
	runGit(t, "commit", "-q", "-m", message)
	args := chargeArgs(t, root, slug, full)
	args[8] = runGit(t, "rev-parse", "HEAD")
	return args
}

func legacyChargeCases() []legacyChargeCase {
	replacePhase := func(t *testing.T, prepare func(t *testing.T)) {
		t.Helper()
		if err := os.Remove(buildPhase); err != nil {
			t.Fatal(err)
		}
		prepare(t)
	}
	return []legacyChargeCase{
		{"compact", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			return legacyChargeRun(t, root, chargeArgs(t, root, slug, false))
		}},
		{"full", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			return legacyChargeRun(t, root, chargeArgs(t, root, slug, true))
		}},
		{"large-ticket-compact", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticketDoc("One", "PF1", "PF2")+strings.Repeat("large ticket evidence\n", 3000))
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "large ticket", false))
		}},
		{"unicode-without-final-newline-full", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			ticket := strings.Replace(ticketDoc("One", "PF1", "PF2"), "Writes: specs", "Writes: specs, internal/example", 1)
			mustWriteFile(t, "specs/"+slug+"/tickets/one.md", strings.TrimSuffix(ticket, "\n")+"\n\nRésumé 雪\t\"q\" \\ x")
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "unicode ticket", true))
		}},
		{"refusal-absent-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			replacePhase(t, func(*testing.T) {})
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "absent phase", false))
		}},
		{"refusal-empty-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			mustWriteFile(t, buildPhase, "")
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "empty phase", false))
		}},
		{"refusal-symlink-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			replacePhase(t, func(t *testing.T) {
				if err := os.Symlink("../skills/bench-craft-delegate/SKILL.md", buildPhase); err != nil {
					t.Fatal(err)
				}
			})
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "linked phase", false))
		}},
		{"refusal-directory-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			replacePhase(t, func(t *testing.T) { mustWriteFile(t, filepath.Join(buildPhase, "inner.md"), "# Inner\n") })
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "directory phase", false))
		}},
		{"refusal-control-byte-phase", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			mustWriteFile(t, buildPhase, "unsafe \x1b source\n")
			return legacyChargeRun(t, root, legacyCommitted(t, root, slug, "control phase", false))
		}},
		{"refusal-dirty-checkout", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n// dirty\n")
			return legacyChargeRun(t, root, args)
		}},
		{"refusal-missing-ticket", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			args[4] = "missing.md"
			return legacyChargeRun(t, root, args)
		}},
		{"refusal-source-tip-mismatch", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			args[8] = args[6]
			out, code := Command(args)
			return out, code, root, args[6], runGit(t, "rev-parse", "HEAD")
		}},
		{"refusal-inactive-assignment", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD")}
			return legacyChargeRun(t, root, args)
		}},
		{"refusal-foreign-assignment", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			activeAssignment(t, root, t.TempDir())
			return legacyChargeRun(t, root, args)
		}},
		{"refusal-persistent-movement", func(t *testing.T) (string, int, string, string, string) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				mustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+strconv.Itoa(calls)+"\n")
			})
			defer restore()
			return legacyChargeRun(t, root, args)
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
	for _, test := range legacyChargeCases() {
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
