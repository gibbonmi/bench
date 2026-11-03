package evidencecmd_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestEvidenceArgumentCommands is S17: each evidence next command carries its identity,
// source, and cursor as always-quoted shell arguments, and its fixed tokens stay bare. The
// wanted bytes are authored from the prepared identity, so a bare value, a dropped flag
// pair, or a quoted fixed token reds the exact command. A POSIX shell then reads each
// command back, so a value that loses its quotes or splits reds the recovered argv. The
// read refuses a hostile identity, source, or cursor before any command renders, and a
// sentinel file reds a value that the shell executes.
func TestEvidenceArgumentCommands(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	ticket := preflighttest.TicketDoc("One", "PF1", "PF2") + strings.Repeat("paged ticket\n", 1000)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
	identity, prepared, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "paged source"))
	hex := strings.TrimPrefix(identity, "sha256:")
	sourcePage, code := preflight.Command([]string{"evidence", identity, "--source", "s2"})
	if code != 0 {
		t.Fatalf("source read = (%d):\n%s", code, sourcePage)
	}
	sourceRows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, sourcePage), "page")
	if len(sourceRows) != 1 {
		t.Fatalf("source page rows = %d:\n%s", len(sourceRows), sourcePage)
	}
	for _, row := range []struct {
		name, next, want string
		argv             []string
	}{
		{"preparation", prepared["next"].(string), "bench preflight evidence '" + identity + "'",
			[]string{"bench", "preflight", "evidence", identity}},
		{"summary", readSummary(t, identity)["next"].(string),
			"bench preflight evidence '" + identity + "' --cursor 'v1." + hex + ".m.0.0'",
			[]string{"bench", "preflight", "evidence", identity, "--cursor", "v1." + hex + ".m.0.0"}},
		{"source page", sourceRows[0]["next"].(string),
			"bench preflight evidence '" + identity + "' --source 's2' --cursor 'v1." + hex + ".s.2.1'",
			[]string{"bench", "preflight", "evidence", identity, "--source", "s2", "--cursor", "v1." + hex + ".s.2.1"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			if row.next != row.want {
				t.Fatalf("next = %q, want %q", row.next, row.want)
			}
			argv, err := axitest.HelpAction{Cmd: row.next}.Argv()
			if err != nil || !slices.Equal(argv, row.argv) {
				t.Fatalf("argv of %q = (%q, %v), want %q", row.next, argv, err, row.argv)
			}
		})
	}
	sentinel := filepath.Join(t.TempDir(), "executed")
	for _, row := range []struct {
		name string
		args []string
	}{
		{"identity with a command substitution", []string{"evidence", "sha256:$(touch " + sentinel + ")"}},
		{"source with an embedded quote", []string{"evidence", identity, "--source", "s2'; touch " + sentinel}},
		{"cursor with a command list", []string{"evidence", identity, "--cursor", "v1." + hex + ".m.0.0; touch " + sentinel}},
	} {
		t.Run(row.name, func(t *testing.T) {
			if out, code := preflight.Command(row.args); code != 2 || strings.Contains(out, "bench preflight evidence") {
				t.Fatalf("hostile read = (%d), want the usage refusal and no next command:\n%s", code, out)
			}
		})
	}
	if _, err := os.Stat(sentinel); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("an evidence operand ran as a command: %v", err)
	}
}
