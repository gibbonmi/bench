package preflight

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestChargeTextSinksUseTOONEscapingAndRefusal grades the source-byte refusal the
// preparation shares with every other required source. The evidence read path owns the
// exact-text rows, so this file keeps only the refusal direction.
func TestChargeTextSinksUseTOONEscapingAndRefusal(t *testing.T) {
	for _, control := range []struct {
		name string
		byte byte
	}{{"ESC", 0x1b}, {"BEL", 0x07}} {
		t.Run(control.name, func(t *testing.T) {
			root, slug := preflighttest.SeedConformant(t)
			preflighttest.MustWriteFile(t, chargesource.BuildPhase, "unsafe "+string(control.byte)+" source\n")
			preflighttest.RunGit(t, "add", chargesource.BuildPhase)
			preflighttest.RunGit(t, "commit", "-q", "-m", "source with "+control.name)
			out, code := Command(preflighttest.ChargeArgs(t, root, slug, false))
			if code != 1 || !strings.Contains(out, "source required") ||
				!strings.Contains(out, chargesource.BuildPhase) ||
				!strings.Contains(out, "cannot represent") || strings.Contains(out, "prepared[") {
				t.Fatalf("%s source = (%d):\n%s", control.name, code, out)
			}
		})
	}
}

func TestChargeQuotesNumericLookingSourceTip(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	t.Setenv("GIT_COMMITTER_DATE", "2026-01-02T03:04:05Z")
	tip := ""
	for i := 0; i < 512; i++ {
		preflighttest.RunGit(t, "commit", "--amend", "-q", "-m", "numeric tip "+strconv.Itoa(i))
		tip = preflighttest.RunGit(t, "rev-parse", "HEAD")
		if tip[0] == '0' && tip[1] >= '0' && tip[1] <= '9' {
			break
		}
	}
	if tip[0] != '0' || tip[1] < '0' || tip[1] > '9' {
		t.Fatalf("fixture did not produce a numeric-looking tip: %s", tip)
	}
	out, code := Command(preflighttest.ChargeArgs(t, root, slug, false))
	if code != 0 || !strings.Contains(out, ",\""+tip+"\",") {
		t.Fatalf("numeric-looking source tip = (%d, %s):\n%s", code, tip, out)
	}
}
