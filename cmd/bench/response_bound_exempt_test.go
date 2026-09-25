package main

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/usage"
)

// fixtureLines is the length of each fixture response here. It is well above the line
// value, so a bounded fixture spills and an exempt one prints every line.
const fixtureLines = 30

// helpFixtures answers a bounded plain command and a bounded family, each of which prints
// fixtureLines lines for every call.
func helpFixtures() []commandDefinition {
	return []commandDefinition{
		{Name: "fixture", Inventory: publicInventory(), Bound: boundResponse, Run: linesHandler(stdoutOf, fixtureLines)},
		{Name: "family", Inventory: publicInventory(), Bound: boundResponse,
			LeafUsage: func() string { return numberedLines(fixtureLines) },
			Leaves: []commandLeaf{{Name: "leaf", Bound: boundResponse, Run: func(c Command, _ string, args []string) int {
				return linesHandler(stdoutOf, fixtureLines)(c, args)
			}}}},
	}
}

// registryWith answers the production registry with the handler of the entry named name
// replaced by run. The entry keeps its declared disposition.
func registryWith(name string, run commandHandler) []commandDefinition {
	registry := slices.Clone(commandRegistry)
	for i := range registry {
		if registry[i].Name == name {
			registry[i].Run = run
		}
	}
	return registry
}

// requireComplete fails unless run printed want in full on stdout and nothing else.
func requireComplete(t *testing.T, argv []string, run boundRun, want string) {
	t.Helper()
	if run.code != 0 || run.stdout != want || run.stderr != "" {
		t.Errorf("%q = (%d, %q, %q), want the complete response %q", argv, run.code, run.stdout, run.stderr, want)
	}
}

// requireBounded fails unless run printed the projection of a fixtureLines response: 10
// stdout lines whose fifth is the spill line.
func requireBounded(t *testing.T, argv []string, run boundRun) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(run.stdout, "\n"), "\n")
	if len(lines) != 10 || !strings.HasPrefix(lines[4], fmt.Sprintf("spilled{lines=%d,", fixtureLines)) {
		t.Errorf("%q = (%d, %q, %q), want 10 lines with the spill line fifth", argv, run.code, run.stdout, run.stderr)
	}
}

// BO8 and BO9: each help form prints its complete inventory or grammar.
func TestHelpFormsStayComplete(t *testing.T) {
	for _, argv := range [][]string{{"help"}, {"--help"}, {"-h"}} {
		requireComplete(t, argv, runRegistry(t, commandRegistry, argv...), renderCommandHelp())
	}
	for _, argv := range [][]string{{"worktree", "--help"}, {"worktree", "-h"}, {"worktree", "help"}} {
		requireComplete(t, argv, runRegistry(t, commandRegistry, argv...), usage.WorktreeUsage())
	}
	for _, argv := range [][]string{
		{"fixture", "--help"}, {"fixture", "-h"}, {"fixture", "help"},
		{"family", "--help"}, {"family", "leaf", "--help"}, {"family", "leaf", "-h"}, {"family", "leaf", "help"},
	} {
		requireComplete(t, argv, runRegistry(t, helpFixtures(), argv...), numberedLines(fixtureLines))
	}
}

// BO10: a help argument after a second argument is no help form, so the call stays
// bounded.
func TestHelpExemptionNeedsOneArgument(t *testing.T) {
	for _, argv := range [][]string{
		{"fixture"}, {"fixture", "operand", "--help"}, {"fixture", "--help", "--help"},
		{"family", "leaf"}, {"family", "leaf", "operand", "--help"}, {"family", "leaf", "--", "--help"},
	} {
		requireBounded(t, argv, runRegistry(t, helpFixtures(), argv...))
	}
}

// BO11: `bench dashboard --stdout` prints the complete page, and the bare verb, which
// writes the page to a file, stays bounded.
func TestDashboardStdoutStaysComplete(t *testing.T) {
	registry := registryWith("dashboard", linesHandler(stdoutOf, fixtureLines))
	requireComplete(t, []string{"dashboard", "--stdout"}, runRegistry(t, registry, "dashboard", "--stdout"), numberedLines(fixtureLines))
	requireBounded(t, []string{"dashboard"}, runRegistry(t, registry, "dashboard"))

	t.Setenv(benchhome.Env, t.TempDir())
	page := runAXICommandAt(t, newAXIEnvelopeRepo(t), []string{"dashboard", "--stdout"})
	if page.code != 0 || !strings.HasPrefix(page.stdout, "<!DOCTYPE html>") || !strings.HasSuffix(page.stdout, "</html>\n") || strings.Contains(page.stdout, "spilled{") {
		t.Fatalf("dashboard --stdout = (%d, %q, %q), want the complete page", page.code, page.stdout, page.stderr)
	}
}

// BO12: the exempt set equals its closed members, and each exemption names a reason.
func TestBoundExemptionsAreClosed(t *testing.T) {
	exempt := map[string]string{"help forms": boundHelpForm.exempt}
	for _, entry := range boundEntries(commandRegistry) {
		if !entry.disposition.bounded {
			exempt[entry.name] = entry.disposition.exempt
		}
	}
	var members []string
	for member, reason := range exempt {
		members = append(members, member)
		if reason == "" {
			t.Errorf("exemption %q names no reason", member)
		}
	}
	sort.Strings(members)
	want := []string{"dashboard --stdout", "help forms", "prep-release", "release", "release-preflight", "repair", "setup", "worktree shell"}
	if !reflect.DeepEqual(members, want) {
		t.Fatalf("exempt set = %q, want %q", members, want)
	}
}

// BO13: a hook or internal plumbing entry stays outside the bound, so its envelope
// prints in full. The run leaves any repository, so no hook span records.
func TestPlumbingStaysOutsideBound(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, definition := range commandRegistry {
		if definition.Inventory.Visibility == inventoryPublic {
			continue
		}
		run := runRegistry(t, registryWith(definition.Name, linesHandler(stderrOf, fixtureLines)), definition.Name)
		if run.code != 0 || run.stdout != "" || run.stderr != numberedLines(fixtureLines) {
			t.Errorf("%s = (%d, %q, %q), want the complete stderr", definition.Name, run.code, run.stdout, run.stderr)
		}
	}
}

// BO70: each ship-tier command prints its complete response, so a CI log keeps the
// release evidence.
func TestShipTierStaysComplete(t *testing.T) {
	for _, name := range []string{"release-preflight", "prep-release", "release"} {
		registry := registryWith(name, linesHandler(stdoutOf, fixtureLines))
		requireComplete(t, []string{name}, runRegistry(t, registry, name), numberedLines(fixtureLines))
	}
}
