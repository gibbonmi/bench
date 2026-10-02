package compatibility

import (
	"github.com/gibbonmi/bench/internal/toon"
	"strings"
)

// RecoveryAction states who can authorize a possible recovery action.
// It does not assert that the selected interface needs that action.
type RecoveryAction struct {
	Category  string
	Authority string
	Action    string
}

// RecoveryActions supplies the shared recovery authority inventory.
func RecoveryActions() []RecoveryAction {
	return []RecoveryAction{
		{"managed-integration", "automatic", "restore only missing or unchanged Bench-managed assets with retained preimages; preserve modified and foreign assets"},
		{"upstream-repair", "unresolved", "no supported automatic harness repair is established; retain the startup error and use supported external diagnostics before a normal-tool retest"},
		{"private-runtime", "decision", "if recovery requires private runtime edits, obtain an explicit reviewer decision; no automatic upstream repair is established"},
		{"process-interruption", "decision", "if recovery requires stopping or restarting a process, preserve active work and obtain an explicit reviewer decision before interruption"},
		{"security-policy", "decision", "if recovery requires permission or trust changes, obtain an explicit reviewer decision before applying them"},
	}
}

// RenderRecovery emits authority boundaries without asserting a diagnosis.
func RenderRecovery() (string, error) {
	rows := [][]string{}
	for _, item := range RecoveryActions() {
		rows = append(rows, []string{item.Category, item.Authority, item.Action})
	}
	return toon.Table("recovery", []string{"category", "authority", "action"}, rows)
}

// QualificationAction names the live checks that local installation cannot supply.
func QualificationAction() string {
	names := []string{}
	for _, selected := range Interfaces() {
		names = append(names, string(selected))
	}
	return "interface qualification remains pending; run normal shell and repository-wrapper checks through the actual tools in each interface: " + strings.Join(names, ", ")
}
