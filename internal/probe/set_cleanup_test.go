package probe

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// TestSetSelectionCloseFailure is PL130. The set's owned build leaves a directory that its
// Close cannot remove, so the held selection's Close fails after the last run. Every row bit,
// but the set exits 1, and the cleanup line follows the reports. A render refusal keeps it too.
func TestSetSelectionCloseFailure(t *testing.T) {
	cleanupPrefix, _, _ := strings.Cut(testreport.SelectionCleanupFailure(errors.New("detail")), "detail")
	for _, tc := range []struct {
		name, run string
		refused   bool
	}{
		{"answer", "^TestClampNegative$", false},
		{"render refusal", "^TestClampNegative$|\a", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := clampFixture(t)
			stuck := filepath.Join(scratchDir(t), "stuck")
			installSealedBuild(t, f, "kept=\"$(dirname \"$2\")/kept\"\nmkdir \"$kept\"\n: > \"$kept/file\"\nchmod 0500 \"$kept\"\nprintf '%s\\n' \"$kept\" >> "+sanitize.ShellQuote(stuck)+"\n")
			t.Cleanup(func() {
				body, _ := os.ReadFile(stuck)
				for _, kept := range strings.Fields(string(body)) {
					_ = os.Chmod(kept, 0o700)
				}
			})
			installSet(t, f, bite, bite)
			out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap), "--package", "./", "--run", tc.run)...)
			if code != 1 {
				t.Fatalf("exit = %d, want 1\n%s", code, out)
			}
			at := strings.LastIndex(out, cleanupPrefix)
			if at < 0 || !strings.HasSuffix(out, "\n") || strings.Count(out[at:], "\n") != 1 || !strings.Contains(out[at:], "/kept") {
				t.Fatalf("stdout = %q, want the held selection's cleanup line %q last", out, cleanupPrefix)
			}
			if tc.refused {
				if refusal, _, _ := strings.Cut(toon.RenderError(errors.New("detail")), "detail"); !strings.HasPrefix(out, refusal) {
					t.Fatalf("stdout = %q, want the render refusal", out)
				}
				return
			}
			if report := strings.LastIndex(out, reportHeading(t, 2)); report < 0 || report > at {
				t.Fatalf("stdout = %q, want the cleanup line after the reports", out)
			}
			requireRow(t, out, code, 1, setRows(t, biteRow(1, "swap"), biteRow(2, "swap"))+setSelection(t, "passed", 1))
			requireTrailingBytes(t, out, 1)
		})
	}
}
