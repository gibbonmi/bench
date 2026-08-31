package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestHookSeamsMatchTheHookDispatchRows covers OT21 at the registry. Every hook row
// records through one symbol, so the conformance check cannot see a dropped Hook flag or
// a registry row that names no verb. The two lists are authored independently — the
// dispatch table declares the hooks, and the registry declares the recorded seams — so
// set equality is the only reconciliation between them.
func TestHookSeamsMatchTheHookDispatchRows(t *testing.T) {
	var dispatched []string
	for _, definition := range commandRegistry {
		if definition.Hook {
			dispatched = append(dispatched, otelHookSeamPrefix+definition.Name)
		}
	}
	var registered []string
	for _, entry := range otelrecord.Registry {
		if strings.HasPrefix(entry.Seam, otelHookSeamPrefix) {
			registered = append(registered, entry.Seam)
		}
	}
	sort.Strings(dispatched)
	sort.Strings(registered)
	if len(dispatched) == 0 {
		t.Fatal("the dispatch table declares no hook row")
	}
	if strings.Join(dispatched, ",") != strings.Join(registered, ",") {
		t.Fatalf("hook seams: dispatch %v, registry %v", dispatched, registered)
	}
}
