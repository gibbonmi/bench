package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// checkCommandScopes grades where the command registry and the worktree leaf table
// declare a tree scope. A public definition that is not a family declares one, a family
// and a plumbing definition declare none, and each leaf row of the family declares one.
// The routing check already reports a registry that does not parse, so this check adds
// nothing for it. A partial tree without the leaf file grades the registry only, by the
// same rule the routing check applies to an absent package.
func checkCommandScopes(root, path, body string) []string {
	entries, err := parseCommandRegistry(path, body)
	if err != nil {
		return nil
	}
	var diags []string
	for _, entry := range entries {
		declared := len(entry.fields["Scope"]) != 0
		switch {
		case len(entry.fields["Leaves"]) != 0:
			if declared {
				diags = append(diags, fmt.Sprintf("command family %q declares a scope", entry.name))
			}
		case parityInternalCommand(entry):
			if declared {
				diags = append(diags, fmt.Sprintf("plumbing command %q declares a scope", entry.name))
			}
		// The help-inventory check refuses any classification other than the public
		// and the internal one, so a classified entry that is not internal is public.
		case len(entry.fields["Inventory"]) == 1 && !declared:
			diags = append(diags, fmt.Sprintf("command %q declares no scope", entry.name))
		}
	}
	leafPath := filepath.Join(root, filepath.FromSlash(worktreeLeafFile))
	leafBody := readIfExists(leafPath)
	if leafBody == "" {
		return diags
	}
	leaves, err := parseWorktreeLeaves(leafPath, leafBody)
	if err != nil {
		return append(diags, worktreeLeafFile+" cannot be parsed for worktree leaves: "+err.Error())
	}
	for _, leaf := range leaves {
		if len(leaf.fields["Scope"]) == 0 {
			diags = append(diags, fmt.Sprintf("worktree leaf %q declares no scope", leaf.name))
		}
	}
	return diags
}

// TestCommandScopeCheckBites is the recorded bite proof for the four scope rules. It plants
// one registry line and one leaf row that break each rule, beside one line and one row that
// obey them, and runs the real subcommand-routing check over the planted tree.
func TestCommandScopeCheckBites(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(dispatchFile, "package main\n\nvar commandRegistry = []commandDefinition{\n"+
		"\t{Name: \"x\", Inventory: publicInventory(helpRow{Order: 1, Description: \"x\"})},\n"+
		"\t{Name: \"declared\", Inventory: publicInventory(helpRow{Order: 2, Description: \"declared\"}), Scope: scopeTree},\n"+
		"\t{Name: \"z\", Inventory: internalInventory, Scope: scopeRepository},\n"+
		"\t{Name: \"quiet\", Inventory: internalInventory},\n"+
		"\t{Name: \"f\", Inventory: publicInventory(helpRow{Order: 3, Description: \"f\"}), Scope: scopeRepository, Leaves: worktreeLeaves},\n"+
		"}\n")
	write(worktreeLeafFile, "package main\n\nvar worktreeLeaves = []commandLeaf{\n\t{Name: \"y\"},\n\t{Name: \"leaf\", Scope: scopeRepository},\n}\n")

	var got []string
	for _, diag := range checkSubcommandRouting(root) {
		if strings.Contains(diag, "scope") {
			got = append(got, diag)
		}
	}
	want := []string{
		`command "x" declares no scope`,
		`command family "f" declares a scope`,
		`plumbing command "z" declares a scope`,
		`worktree leaf "y" declares no scope`,
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scope diagnostics = %q, want %q", got, want)
	}
}
