package staleness

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// Decide classifies one AuditFacts value. It performs no I/O, and the same facts always
// yield the same audit. spec-drift is red when any audit-scope member changed. A check that a
// diagnostic names is diagnostic, and any diagnostic makes the audit incomplete.
func Decide(facts AuditFacts) Audit {
	audit := Audit{
		Baseline:    facts.Baseline,
		Tip:         facts.Tip,
		Drift:       drift(facts.Paths),
		Diagnostics: facts.Diagnostics,
		Complete:    len(facts.Diagnostics) == 0,
	}
	audit.ManualRequired = !audit.Complete
	for _, name := range checkNames {
		state := StateGreen
		if name == CheckSpecDrift && len(audit.Drift) > 0 {
			state = StateRed
		}
		if slices.ContainsFunc(facts.Diagnostics, func(d Diagnostic) bool { return d.Subject == name }) {
			state = StateDiagnostic
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
