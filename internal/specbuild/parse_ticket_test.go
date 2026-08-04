package specbuild_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/specbuild"
)

// These tests live in the external test package so the export stays
// load-bearing: an unexported parse would not compile here.

func ticketFixture(t *testing.T, body string) (specPath, ticketPath string) {
	t.Helper()
	root := t.TempDir()
	specPath = filepath.Join(root, "spec.md")
	if err := os.WriteFile(specPath, []byte("# spec\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "tickets"), 0o755); err != nil {
		t.Fatalf("create tickets directory: %v", err)
	}
	ticketPath = filepath.Join(root, "tickets", "one.md")
	if err := os.WriteFile(ticketPath, []byte(body), 0o644); err != nil {
		t.Fatalf("write ticket: %v", err)
	}
	return specPath, ticketPath
}

func TestParseTicketMatchesAssignPath(t *testing.T) {
	body := "# Export the ticket parser as one cross-package entry point\n" +
		"\n" +
		"Blocked by: none\n" +
		"Ownership fence: `internal/specbuild`, `internal/conformance`\n" +
		"Assumptions: the parse is shared; the fence entries are path prefixes\n" +
		"\n" +
		"## Acceptance\n" +
		"\n" +
		"- [ ] [EP1] the exported entry parses every declared field.\n" +
		"- [x] [EP2] the exported entry refuses a fenceless ticket.\n"
	specPath, ticketPath := ticketFixture(t, body)

	parsed, err := specbuild.ParseTicket(specPath, "one.md")
	if err != nil {
		t.Fatalf("ParseTicket: %v", err)
	}
	sum := sha256.Sum256([]byte(body))
	if parsed.Path != ticketPath {
		t.Errorf("Path = %q, want %q", parsed.Path, ticketPath)
	}
	if parsed.Title != "Export the ticket parser as one cross-package entry point" {
		t.Errorf("Title = %q", parsed.Title)
	}
	if parsed.Digest != hex.EncodeToString(sum[:]) {
		t.Errorf("Digest = %q, want the sha256 of the ticket bytes", parsed.Digest)
	}
	if want := []string{"EP1", "EP2"}; !slices.Equal(parsed.Rows, want) {
		t.Errorf("Rows = %q, want %q", parsed.Rows, want)
	}
	if want := []string{"internal/specbuild", "internal/conformance"}; !slices.Equal(parsed.Fence, want) {
		t.Errorf("Fence = %q, want %q", parsed.Fence, want)
	}
}

// Each refusal is asserted by its own message: the lifecycle fixtures are
// malformed in more than one way at once, so a ticket refused for the wrong
// reason reads there as a pass.
func TestParseTicketRefusesZeroRowsAndFencelessTickets(t *testing.T) {
	for _, row := range []struct{ name, body, want string }{
		{
			name: "zero acceptance rows",
			body: "# Fenced but uncharged\n\nOwnership fence: `internal/specbuild`\n",
			want: "spec build ticket declares no charged rows",
		},
		{
			name: "no fence and no package mention",
			body: "# Charged but unfenced\n\n- [ ] [EP1] the ticket names no path it may write.\n",
			want: "spec build ticket declares no ownership fence",
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			specPath, _ := ticketFixture(t, row.body)
			parsed, err := specbuild.ParseTicket(specPath, "one.md")
			if err == nil {
				t.Fatalf("ParseTicket accepted %q as %#v", row.name, parsed)
			}
			if err.Error() != row.want {
				t.Errorf("error = %q, want %q", err, row.want)
			}
		})
	}
}

func TestParseTicketRefusesDuplicateAcceptanceIDs(t *testing.T) {
	for _, row := range []struct{ name, body, want string }{
		{
			name: "literal repeat",
			body: "# Twice charged\n\nOwnership fence: `internal/specbuild`\n\n- [ ] [EP1] first obligation.\n- [ ] [EP1] second obligation under the same ID.\n",
			want: "spec build ticket one.md declares duplicate acceptance ID EP1",
		},
		{
			name: "range overlap",
			body: "# Overlapping ranges\n\nOwnership fence: `internal/specbuild`\n\n- [ ] [R10-R12] the range charges three rows.\n- [ ] [R11] one of them is charged again.\n",
			want: "spec build ticket one.md declares duplicate acceptance ID R11",
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			specPath, _ := ticketFixture(t, row.body)
			parsed, err := specbuild.ParseTicket(specPath, "one.md")
			if err == nil {
				t.Fatalf("ParseTicket accepted %q as %#v", row.name, parsed)
			}
			if err.Error() != row.want {
				t.Errorf("error = %q, want %q", err, row.want)
			}
		})
	}
}

func TestParseTicketRangeExpansionUnchanged(t *testing.T) {
	body := "# Range expansion\n" +
		"\n" +
		"Ownership fence: `internal/specbuild`\n" +
		"\n" +
		"- [ ] [R10-R12] three rows expand from one R-prefixed range.\n" +
		"- [ ] [EP1] a plain identifier stays literal.\n" +
		"- [ ] [A1-A3] a non-R range stays literal.\n"
	specPath, _ := ticketFixture(t, body)

	parsed, err := specbuild.ParseTicket(specPath, "one.md")
	if err != nil {
		t.Fatalf("ParseTicket: %v", err)
	}
	if want := []string{"R10", "R11", "R12", "EP1", "A1-A3"}; !slices.Equal(parsed.Rows, want) {
		t.Errorf("Rows = %q, want %q", parsed.Rows, want)
	}
}

// A ticket may only maintain what its fence lets it write, so a declared
// crossing that names no fenced path is a repair scoped to the lines that
// prompted it rather than to the artifacts it changes. The parse stays
// permissive — its other consumers grade the grammar — and the assign path
// refuses on this answer.
func TestTicketContractsAnchoredToTheOwnershipFence(t *testing.T) {
	const charged = "\n- [ ] [CA1] the crossing names a path the ticket writes.\n"
	for _, row := range []struct {
		name, contracts string
		anchored        bool
	}{
		{
			name:      "concept operands on both sides",
			contracts: "Contracts: every registry row's port name crosses registry\u2192derived inventory\n",
		},
		{
			name:      "one operand inside the fence",
			contracts: "Contracts: the gate invocation crosses `internal/specbuild`\u2192bash, asserted by CA1 against the real subprocess\n",
			anchored:  true,
		},
		{
			name:      "far side names no path at all",
			contracts: "Contracts: the derived inventory crosses `internal/specbuild`\u2192every audited package, asserted by CA1 against the real derivation\n",
			anchored:  true,
		},
		{
			name:      "explicitly declared empty",
			contracts: "Contracts: none crosses\n",
			anchored:  true,
		},
		{
			name:     "no Contracts line",
			anchored: true,
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			body := "# Anchor the crossing\n\nOwnership fence: `internal/specbuild`\n" + row.contracts + charged
			specPath, _ := ticketFixture(t, body)
			parsed, err := specbuild.ParseTicket(specPath, "one.md")
			if err != nil {
				t.Fatalf("ParseTicket refused a well-formed ticket: %v", err)
			}
			if got := parsed.ContractsAnchored(); got != row.anchored {
				t.Errorf("ContractsAnchored() = %v, want %v", got, row.anchored)
			}
		})
	}
}

func TestParseTicketReadsTheContractsLine(t *testing.T) {
	body := "# Carry the crossing into the parsed ticket\n" +
		"\nOwnership fence: `internal/specbuild`\n" +
		"Contracts: the parsed operands cross `internal/specbuild`\u2192callers, asserted by CA1 against the real parser\n" +
		"\n- [ ] [CA1] the parsed ticket carries its declared crossing.\n"
	specPath, _ := ticketFixture(t, body)
	parsed, err := specbuild.ParseTicket(specPath, "one.md")
	if err != nil {
		t.Fatalf("ParseTicket: %v", err)
	}
	const want = "the parsed operands cross `internal/specbuild`\u2192callers, asserted by CA1 against the real parser"
	if parsed.Contracts != want {
		t.Errorf("Contracts = %q, want %q", parsed.Contracts, want)
	}
}
