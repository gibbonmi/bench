package conformance

import (
	"fmt"
	"go/ast"
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
		// A family or a plumbing line carries no Scope field at all, so even a zero value
		// on such a line is a declaration that the rule refuses.
		written := len(entry.fields["Scope"]) != 0
		switch {
		case len(entry.fields["Leaves"]) != 0:
			if written {
				diags = append(diags, fmt.Sprintf("command family %q declares a scope", entry.name))
			}
		case parityInternalCommand(entry):
			if written {
				diags = append(diags, fmt.Sprintf("plumbing command %q declares a scope", entry.name))
			}
		// The help-inventory check refuses any classification other than the public
		// and the internal one, so a classified entry that is not internal is public.
		case len(entry.fields["Inventory"]) == 1 && !declaresScope(entry):
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
		if !declaresScope(leaf) {
			diags = append(diags, fmt.Sprintf("worktree leaf %q declares no scope", leaf.name))
		}
	}
	return diags
}

// declaredScopes names the values of the scope type in `cmd/bench`. The zero value means
// undeclared, so only a line that names one of these values declares a scope.
var declaredScopes = map[string]bool{"scopeTree": true, "scopeRepository": true}

func declaresScope(entry commandRegistryEntry) bool {
	values := entry.fields["Scope"]
	if len(values) != 1 {
		return false
	}
	value, ok := values[0].(*ast.Ident)
	return ok && declaredScopes[value.Name]
}

// TestCommandScopeCheckBites is the recorded bite proof for the four scope rules. It plants
// one registry line and one leaf row that break each rule, beside one line and one row that
// obey them, and runs the real subcommand-routing check over the planted tree.
func TestCommandScopeCheckBites(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(dispatchFile)), "package main\n\nvar commandRegistry = []commandDefinition{\n"+
		"\t{Name: \"x\", Inventory: publicInventory(helpRow{Order: 1, Description: \"x\"})},\n"+
		"\t{Name: \"zero\", Inventory: publicInventory(helpRow{Order: 2, Description: \"zero\"}), Scope: 0},\n"+
		"\t{Name: \"declared\", Inventory: publicInventory(helpRow{Order: 3, Description: \"declared\"}), Scope: scopeTree},\n"+
		"\t{Name: \"z\", Inventory: internalInventory, Scope: scopeRepository},\n"+
		"\t{Name: \"quiet\", Inventory: internalInventory},\n"+
		"\t{Name: \"f\", Inventory: publicInventory(helpRow{Order: 4, Description: \"f\"}), Scope: scopeRepository, Leaves: worktreeLeaves},\n"+
		"}\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(worktreeLeafFile)), "package main\n\nvar worktreeLeaves = []commandLeaf{\n"+
		"\t{Name: \"y\"},\n\t{Name: \"zeroleaf\", Scope: 0},\n\t{Name: \"leaf\", Scope: scopeRepository},\n}\n")

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
		`command "zero" declares no scope`,
		`worktree leaf "zeroleaf" declares no scope`,
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scope diagnostics = %q, want %q", got, want)
	}
}
