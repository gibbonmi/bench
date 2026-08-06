package specbuild

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

const (
	optedInCoversHeader = "| row | story | behavior | seam | red signal | why it catches the failure |"
	legacyCoversHeader  = "| story | behavior | seam | red signal | why it catches the failure |"
)

// coversSpec renders a build-demo spec declaring one story under the given
// coverage-map header and rows, so each case below differs only in the map its
// author opted into. The separator line is header-width-agnostic: the parser
// skips an all-dashes row whatever its cell count.
func coversSpec(header string, rows ...string) string {
	return "# Build demo\n\nStatus: staged\n\n## User stories\n\n1. the story\n\n### Acceptance coverage map\n\n" +
		header + "\n|---|---|---|---|---|---|\n" + strings.Join(rows, "\n") + "\n"
}

// coversRepo seeds a started build run whose spec is spec and whose one.md
// charges rows, both committed: assign reads the covers policy off the run's
// real spec file, so the fixture has to be the real artifact on disk.
func coversRepo(t *testing.T, spec, rows string) *Service {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "spec.md"), spec)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\n\n"+rows)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "covers fixture")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return service
}

func TestAssignRefusesCoversViolationsUnderAnOptedInSpec(t *testing.T) {
	spec := coversSpec(optedInCoversHeader,
		"| AB1 | 1 | behavior | seam | red | why |",
		"| AB2 | 1 | behavior | seam | red | why |")
	for _, tc := range []struct{ name, rows, want string }{
		{
			"missing", "- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] unannotated\n",
			"spec build assign requires a covers annotation on acceptance row R11 of ticket one.md under an opted-in coverage map",
		},
		{
			// A lowercase operand fails the ID grammar, so ParseTicket leaves the row
			// unannotated and the missing-covers rule is what refuses it.
			"malformed", "- [ ] [R10] (covers ab1) lowercase\n",
			"spec build assign requires a covers annotation on acceptance row R10 of ticket one.md under an opted-in coverage map",
		},
		{
			// A range line carries no per-row provenance, so its expanded rows are
			// unannotated by construction and refuse under the same rule.
			"range", "- [ ] [R10-R11] (covers AB1) ranged\n",
			"spec build assign requires a covers annotation on acceptance row R10 of ticket one.md under an opted-in coverage map",
		},
		{
			"dangling", "- [ ] [R10] (covers ZZ9) names no map row\n",
			"spec build assign requires acceptance row R10 of ticket one.md to name a declared coverage map row, but ZZ9 names none",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := coversRepo(t, spec, tc.rows)
			assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", tc.name)
			if err == nil {
				t.Fatalf("Assign leased %#v for a %s covers annotation", assigned, tc.name)
			}
			if err.Error() != tc.want {
				t.Fatalf("error = %q, want %q", err, tc.want)
			}
		})
	}
}

func TestAssignAcceptsMappedAndLocalCoversAndLegacySpecs(t *testing.T) {
	optedIn := coversSpec(optedInCoversHeader,
		"| AB1 | 1 | behavior | seam | red | why |",
		"| AB2 | 1 | behavior | seam | red | why |")
	legacy := coversSpec(legacyCoversHeader, "| 1 | behavior | seam | red | why |")
	for _, tc := range []struct{ name, spec, rows string }{
		{"mapped and local", optedIn, "- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] (covers AB2) mapped\n- [ ] [R12] (covers local) ticket-time repair\n"},
		{"legacy map", legacy, "- [ ] [R10] unannotated under a five-cell map\n"},
		{"no map", "# Build demo\n\nStatus: staged\n", "- [ ] [R10] unannotated under no map\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := coversRepo(t, tc.spec, tc.rows)
			assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", tc.name)
			if err != nil {
				t.Fatalf("Assign: %v", err)
			}
			if assigned.ID == "" || assigned.Path == "" {
				t.Fatalf("assignment = %#v", assigned)
			}
		})
	}
}

// An author who has opted into row IDs gets a fail-closed assign: a map the
// checker rejects cannot resolve anything, and treating it as legacy would let
// a broken map silently disable the policy.
func TestAssignRefusesAnOptedInMapThatFailsIDValidation(t *testing.T) {
	spec := coversSpec(optedInCoversHeader,
		"| AB1 | 1 | behavior | seam | red | why |",
		"| AB1 | 1 | behavior | seam | red | why |")
	service := coversRepo(t, spec, "- [ ] [R10] (covers AB1) mapped\n")
	assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", "invalid map")
	if err == nil {
		t.Fatalf("Assign leased %#v under a map that fails ID validation", assigned)
	}
	const want = "spec build assign requires the spec's opted-in coverage map to validate, but it reports coverage map row 2 has a duplicate row id 'AB1' (first used at row 1)"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
}

func closureTicket(closure, mutations string) string {
	return "# Close every ticket fact\n\n" +
		"Ownership fence: `internal/specbuild`\n" +
		"Integration surfaces: parsed closure inventory→`internal/specbuild`\n" +
		"Contracts: none crosses\n" + closure +
		"\n## Acceptance\n\n" +
		"- [ ] [R10] empty and malformed input are refused.\n" +
		"- [ ] [R11] a rerun is idempotent.\n" +
		"\n## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" + mutations
}

func closureRepo(t *testing.T, ticket string) *Service {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), ticket)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "closure fixture")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return service
}

func TestAssignRequiresAtomicClosureForModernTickets(t *testing.T) {
	validMutations := "| R10/empty | pass empty input | parser test | run parser test |\n" +
		"| R10/malformed | pass malformed input | parser test | run parser test |\n" +
		"| R11/rerun | create a second record | lifecycle test | run lifecycle test |\n"
	for _, tc := range []struct {
		name, closure, mutations, want string
	}{
		{
			name:      "missing closure field",
			mutations: validMutations,
			want:      "spec build assign requires ticket one.md to declare an atomic Closure inventory",
		},
		{
			name:      "malformed fact identity",
			closure:   "Closure: R10/Empty, R10/malformed, R11/rerun\n",
			mutations: validMutations,
			want:      "spec build assign requires Closure fact \"R10/Empty\" of ticket one.md to use <acceptance-ID>/<fact-name>",
		},
		{
			name:      "duplicate fact identity",
			closure:   "Closure: R10/empty, R10/empty, R10/malformed, R11/rerun\n",
			mutations: validMutations,
			want:      "spec build assign requires unique Closure facts in ticket one.md, but R10/empty is repeated",
		},
		{
			name:      "acceptance row owns no fact",
			closure:   "Closure: R10/empty, R10/malformed\n",
			mutations: "| R10/empty | pass empty input | parser test | run parser test |\n| R10/malformed | pass malformed input | parser test | run parser test |\n",
			want:      "spec build assign requires every acceptance row of ticket one.md to own a Closure fact, but R11 owns none",
		},
		{
			name:      "fact names no acceptance row",
			closure:   "Closure: R10/empty, R10/malformed, R12/rerun\n",
			mutations: "| R10/empty | pass empty input | parser test | run parser test |\n| R10/malformed | pass malformed input | parser test | run parser test |\n| R12/rerun | create a second record | lifecycle test | run lifecycle test |\n",
			want:      "spec build assign requires Closure fact R12/rerun of ticket one.md to name an acceptance row, but R12 names none",
		},
		{
			name:      "fact has no mutation",
			closure:   "Closure: R10/empty, R10/malformed, R11/rerun\n",
			mutations: "| R10/empty | pass empty input | parser test | run parser test |\n| R11/rerun | create a second record | lifecycle test | run lifecycle test |\n",
			want:      "spec build assign requires every Closure fact of ticket one.md to have a Red mutations row, but R10/malformed has none",
		},
		{
			name:      "mutation names no fact",
			closure:   "Closure: R10/empty, R10/malformed, R11/rerun\n",
			mutations: validMutations + "| R11/interrupt | stop halfway | lifecycle test | run lifecycle test |\n",
			want:      "spec build assign requires every Red mutations criterion of ticket one.md to name a Closure fact, but R11/interrupt names none",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := closureRepo(t, closureTicket(tc.closure, tc.mutations))
			assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", tc.name)
			if err == nil {
				t.Fatalf("Assign leased %#v for an open closure graph", assigned)
			}
			if err.Error() != tc.want {
				t.Fatalf("error = %q, want %q", err, tc.want)
			}
		})
	}
}

func TestAssignAcceptsClosedModernTicketAndLegacyTicket(t *testing.T) {
	modern := closureTicket(
		"Closure: R10/empty, R10/malformed, R11/rerun\n",
		"| R10/empty | pass empty input | parser test | run parser test |\n"+
			"| R10/malformed | pass malformed input | parser test | run parser test |\n"+
			"| R11/rerun | create a second record | lifecycle test | run lifecycle test |\n",
	)
	legacy := "# Legacy ticket\n\nOwnership fence: `internal/specbuild`\n\n- [ ] [R10] old ticket remains assignable.\n"
	for _, tc := range []struct{ name, ticket string }{{"modern", modern}, {"legacy", legacy}} {
		t.Run(tc.name, func(t *testing.T) {
			service := closureRepo(t, tc.ticket)
			assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", tc.name)
			if err != nil {
				t.Fatalf("Assign: %v", err)
			}
			if assigned.ID == "" {
				t.Fatalf("assignment = %#v", assigned)
			}
		})
	}
}
