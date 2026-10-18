package staleness

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// Decide classifies one AuditFacts value. It performs no I/O, and the same facts always
// yield the same audit. spec-drift is red when any audit-scope member changed, cited-symbols
// is red when any cited symbol has a supported discrepancy, and cited-tests is red when any
// test seam has a supported discrepancy; each such discrepancy joins the drift. A check that
// a diagnostic names, or a check over a symbol or a test whose discrepancy no collector
// grades, is diagnostic unless it is red, and any diagnostic makes the audit incomplete.
func Decide(facts AuditFacts) Audit {
	paths := drift(facts.Paths)
	var symbolDrift, testDrift []DriftEntry
	undecided := map[string]bool{}
	for _, state := range facts.Symbols {
		switch {
		case slices.Contains(diagnosticDiscrepancies, state.Discrepancy):
			undecided[CheckCitedSymbols] = true
		case state.Discrepancy != "":
			symbolDrift = append(symbolDrift, symbolEntry(state))
		}
	}
	for _, state := range facts.Tests {
		switch {
		case slices.Contains(redTests, state.Discrepancy):
			testDrift = append(testDrift, testEntry(state))
		case state.Discrepancy != "":
			undecided[CheckCitedTests] = true
		}
	}
	audit := Audit{
		Baseline:    facts.Baseline,
		Tip:         facts.Tip,
		Drift:       append(append(paths, symbolDrift...), testDrift...),
		Diagnostics: facts.Diagnostics,
		Complete:    len(facts.Diagnostics) == 0,
	}
	audit.ManualRequired = !audit.Complete
	red := map[string]bool{CheckSpecDrift: len(paths) > 0, CheckCitedSymbols: len(symbolDrift) > 0, CheckCitedTests: len(testDrift) > 0}
	for _, name := range checkNames {
		state := StateGreen
		if undecided[name] || slices.ContainsFunc(facts.Diagnostics, func(d Diagnostic) bool { return d.Subject == name }) {
			state = StateDiagnostic
		}
		if red[name] {
			state = StateRed
		}
		audit.Checks = append(audit.Checks, Check{Name: name, State: state})
		audit.ManualRequired = audit.ManualRequired || state != StateGreen
	}
	return audit
}

// Red reports whether any check is red. A diagnostic check is not red.
func (a Audit) Red() bool {
	return slices.ContainsFunc(a.Checks, func(c Check) bool { return c.State == StateRed })
}

// Render prints the check states, the drift entries, the diagnostics, and the staleness
// decision as TOON tables. Each table prints even when it is empty, so an empty drift table
// is a definitive answer and never an omission.
func (a Audit) Render() (string, error) {
	var checks, drift, diagnostics [][]string
	for _, check := range a.Checks {
		checks = append(checks, []string{check.Name, string(check.State)})
	}
	for _, entry := range a.Drift {
		drift = append(drift, []string{entry.Before, entry.After, entry.Kind, entry.BeforeObject, entry.AfterObject})
	}
	for _, diagnostic := range a.Diagnostics {
		diagnostics = append(diagnostics, []string{diagnostic.Subject, diagnostic.Reason})
	}
	var b strings.Builder
	for _, table := range []struct {
		name   string
		fields []string
		rows   [][]string
	}{
		{"spec_checks", []string{"check", "state"}, checks},
		{"drift", []string{"before", "after", "kind", "before_object", "after_object"}, drift},
		{"audit_diagnostics", []string{"subject", "reason"}, diagnostics},
	} {
		text, err := toon.Table(table.name, table.fields, table.rows)
		if err != nil {
			return "", err
		}
		b.WriteString(text)
	}
	staleness, err := toon.TableTyped("staleness", []string{"complete", "manual_required", "baseline", "tip"},
		[][]any{{a.Complete, a.ManualRequired, a.Baseline, a.Tip}})
	if err != nil {
		return "", err
	}
	return b.String() + staleness, nil
}
