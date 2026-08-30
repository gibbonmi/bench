package preflight

import (
	"strings"
	"testing"
)

// TestSourceTipOmittedKeepsTodaysVerdict is H30's control half: without
// the flag both modes render exactly the five rows they render today, in
// order, and no tip row appears. The flag is an addition, not a new
// requirement.
func TestSourceTipOmittedKeepsTodaysVerdict(t *testing.T) {
	for _, mode := range []string{"review", "build"} {
		t.Run(mode, func(t *testing.T) {
			_, slug := seedConformant(t)
			out, code := Command([]string{mode, slug})
			if code != 0 {
				t.Fatalf("bare %s = (%d):\n%s", mode, code, out)
			}
			if !strings.Contains(out, "checks[5]{check,verdict,detail,next}") {
				t.Fatalf("bare %s did not render today's five-row table:\n%s", mode, out)
			}
			if strings.Contains(out, "tip-current") {
				t.Fatalf("bare %s rendered a tip row without --source-tip:\n%s", mode, out)
			}
		})
	}
}

// TestSourceTipAcceptedByBothModes is H30's accepting half: `review` and
// `build` both take the flag, in bare and explicit-base form. A pin that
// agrees with the derived tip is green.
func TestSourceTipAcceptedByBothModes(t *testing.T) {
	for _, mode := range []string{"review", "build"} {
		t.Run(mode, func(t *testing.T) {
			_, slug := seedConformant(t)
			base := runGit(t, "rev-parse", "main")
			tip := runGit(t, "rev-parse", "HEAD")

			out, code := Command([]string{mode, slug, "--source-tip", tip})
			if code != 0 || !strings.Contains(out, "tip-current,green") {
				t.Fatalf("bare %s --source-tip = (%d):\n%s", mode, code, out)
			}
			// WF41: the pinned-tip form carries the next column too, so the added
			// row and the added column are pinned by the same fixture.
			if !strings.Contains(out, "checks[6]{check,verdict,detail,next}") {
				t.Fatalf("pinned %s did not add exactly one row:\n%s", mode, out)
			}
			// The pin is verified against the derived tip, not compared literally: a
			// revision spelling that resolves to the same commit is green.
			out, code = Command([]string{mode, slug, "--source-tip", "HEAD"})
			if code != 0 || !strings.Contains(out, "tip-current,green") {
				t.Fatalf("bare %s --source-tip HEAD = (%d):\n%s", mode, code, out)
			}
			out, code = Command([]string{mode, slug, "--base", base, "--source-tip", tip})
			if code != 0 || !strings.Contains(out, "tip-current,green") {
				t.Fatalf("explicit-base %s --source-tip = (%d):\n%s", mode, code, out)
			}
		})
	}
}

// TestSourceTipMismatchRendersRedRow is H31: a resolvable pin that names a
// different commit than the derived tip is a verdict row. A flag that
// were parsed and ignored could not pass.
func TestSourceTipMismatchRendersRedRow(t *testing.T) {
	for _, mode := range []string{"review", "build"} {
		t.Run(mode, func(t *testing.T) {
			_, slug := seedConformant(t)
			base := runGit(t, "rev-parse", "main")
			tip := runGit(t, "rev-parse", "HEAD")

			out, code := Command([]string{mode, slug, "--source-tip", base})
			if code != 1 || !strings.Contains(out, "tip-current,red") {
				t.Fatalf("bare %s stale pin = (%d):\n%s", mode, code, out)
			}
			if !strings.Contains(out, base) || !strings.Contains(out, tip) {
				t.Fatalf("stale pin red did not name both commits:\n%s", out)
			}
			if strings.Contains(out, "error: cannot resolve --source-tip") {
				t.Fatalf("a stale pin reported as an unresolvable one:\n%s", out)
			}
			out, code = Command([]string{mode, slug, "--base", base, "--source-tip", base})
			if code != 1 || !strings.Contains(out, "tip-current,red") {
				t.Fatalf("explicit-base %s stale pin = (%d):\n%s", mode, code, out)
			}
		})
	}
}

// TestSourceTipUnresolvableIsAGrammarErrorNotAMismatch is H32: a pin that names no
// commit never becomes a verdict row, so a typo and a drift stay different
// diagnoses. The control-byte value carries the spec-TOON refusal edge: a
// `--source-tip` value reaches a rendered cell, so an unrepresentable one is refused
// rather than rendered.
func TestSourceTipUnresolvableIsAGrammarErrorNotAMismatch(t *testing.T) {
	t.Run("unreachable revision", func(t *testing.T) {
		_, slug := seedConformant(t)
		out, code := Command([]string{"review", slug, "--source-tip", "missing"})
		if code != 1 || !strings.HasPrefix(out, "error: cannot resolve --source-tip") {
			t.Fatalf("unreachable pin = (%d):\n%s", code, out)
		}
		if strings.Contains(out, "tip-current") || strings.Contains(out, "checks[") {
			t.Fatalf("unreachable pin leaked into the verdict table:\n%s", out)
		}
	})
	t.Run("unreachable revision under an explicit base", func(t *testing.T) {
		_, slug := seedConformant(t)
		base := runGit(t, "rev-parse", "main")
		out, code := Command([]string{"build", slug, "--base", base, "--source-tip", "missing"})
		if code != 1 || !strings.HasPrefix(out, "error: cannot resolve --source-tip") {
			t.Fatalf("unreachable pin under explicit base = (%d):\n%s", code, out)
		}
		if strings.Contains(out, "snapshot drift") || strings.Contains(out, "checks[") {
			t.Fatalf("unreachable pin under explicit base misclassified:\n%s", out)
		}
	})
	t.Run("control byte", func(t *testing.T) {
		_, slug := seedConformant(t)
		out, code := Command([]string{"review", slug, "--source-tip", "feature\x1b"})
		if code != 1 || !strings.Contains(out, "unrepresentable TOON cell") {
			t.Fatalf("control-byte pin = (%d):\n%s", code, out)
		}
		if strings.Contains(out, "\x1b") {
			t.Fatalf("control-byte pin was rendered rather than refused:\n%q", out)
		}
		if strings.Contains(out, "tip-current") || strings.Contains(out, "checks[") {
			t.Fatalf("control-byte pin leaked into the verdict table:\n%s", out)
		}
	})
}
