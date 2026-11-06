package injectedports

import (
	"fmt"
	"slices"
	"strings"
)

// These prefixes name the variable and setter admission failures. Like the named-port
// prefixes, each one is quoted by the tests that prove it fires.
const (
	UnregisteredSeamMessage   = "test seam has no registry row"
	UnregisteredSetterMessage = "test seam setter has no registry row"
	SeamOwnerMessage          = "test seam registry row names no production owner that reads the seam"
	SeamConsumerMessage       = "test seam registry row names a consumer test the tree does not declare"
	SeamUnreachedMessage      = "test seam registry row names a consumer test that never installs the seam"
	SeamPurposeMessage        = "test seam registry row names no failure purpose"
	SeamIsolationMessage      = "test seam registry row names no isolation evidence"
	UnresolvedAliasMessage    = "test seam setter row does not resolve to its canonical seams"
	DuplicateRowMessage       = "injected port registry repeats a source identity"
)

// seamDiags grades the variable and setter rows against a census of the whole production
// source of root. A seam in a package that no row names is still derived, so a new
// replacement cannot hide behind a remembered package list. A row whose package the
// graded tree does not hold is not graded, as for named ports.
func seamDiags(root string, rows []Row) []string {
	tree, diags := loadSourceTree(root)
	census := deriveSeams(tree)
	registered := map[string]Row{}
	for _, row := range rows {
		registered[rowKey(row)] = row
	}
	for _, seam := range sortedKeys(census.variables) {
		row, found := registered[string(KindVariable)+" "+seam]
		if !found {
			diags = append(diags, fmt.Sprintf(
				"%s: %s is replaced at %s — add a variable row to injectedPortRegistry naming its production owner, consumer tests, failure purpose, and isolation evidence",
				UnregisteredSeamMessage, seam, census.variables[seam]))
			continue
		}
		diags = append(diags, census.variableRowDiags(row, seam)...)
	}
	for _, setter := range sortedKeys(census.setters) {
		canonical := census.setters[setter]
		row, found := registered[string(KindSetter)+" "+setter]
		if !found {
			diags = append(diags, fmt.Sprintf(
				"%s: %s replaces %s — add a setter row to injectedPortRegistry whose Canonical names those seams",
				UnregisteredSetterMessage, setter, strings.Join(canonical, ", ")))
			continue
		}
		if !slices.Equal(row.Canonical, canonical) {
			diags = append(diags, fmt.Sprintf(
				"%s: %s names %q but replaces %q — a setter cannot borrow the evidence of a seam it does not replace",
				UnresolvedAliasMessage, setter, row.Canonical, canonical))
		}
	}
	for _, row := range rows {
		if row.Kind == KindPort || tree.packages[row.Package] == nil {
			continue
		}
		identity := row.Package + "." + row.Port
		_, variable := census.variables[identity]
		_, setter := census.setters[identity]
		if !(row.Kind == KindVariable && variable) && !(row.Kind == KindSetter && setter) {
			diags = append(diags, orphanDiag(row))
		}
	}
	return diags
}

// variableRowDiags grades the policy facts of one admitted variable seam.
func (c *seamCensus) variableRowDiags(row Row, seam string) []string {
	var diags []string
	if !c.readers[seam][row.Owner] {
		diags = append(diags, fmt.Sprintf("%s: %s names owner %q, which is not a production function of its package that reads it", SeamOwnerMessage, seam, row.Owner))
	}
	if strings.TrimSpace(row.Purpose) == "" {
		diags = append(diags, fmt.Sprintf("%s: %s — state the production failure or control that its consumers observe", SeamPurposeMessage, seam))
	}
	if strings.TrimSpace(row.Isolation) == "" {
		diags = append(diags, fmt.Sprintf("%s: %s — name the serial ownership, the cleanup owner, and how live consumers finish first", SeamIsolationMessage, seam))
	}
	if len(row.Consumers) == 0 {
		diags = append(diags, fmt.Sprintf("%s: %s names no consumer test", SeamConsumerMessage, seam))
	}
	for _, consumer := range row.Consumers {
		file, test, _ := strings.Cut(consumer, "::")
		declared, reached := c.reaches(file, test, seam)
		switch {
		case !declared:
			diags = append(diags, fmt.Sprintf("%s: %s advertises %s, which is not a top-level test of that file", SeamConsumerMessage, seam, consumer))
		case !reached:
			diags = append(diags, fmt.Sprintf("%s: %s advertises %s, which neither assigns it nor calls its setter through its own helpers", SeamUnreachedMessage, seam, consumer))
		}
	}
	return diags
}

// duplicateRows refuses a source identity that two rows claim, so a later row cannot
// silently replace the evidence of an earlier one.
func duplicateRows(rows []Row) []string {
	seen := map[string]bool{}
	var diags []string
	for _, row := range rows {
		key := rowKey(row)
		if seen[key] {
			diags = append(diags, fmt.Sprintf("%s: %s %s.%s has more than one row", DuplicateRowMessage, kindName(row.Kind), row.Package, row.Port))
		}
		seen[key] = true
	}
	return diags
}

func rowKey(row Row) string {
	return string(row.Kind) + " " + row.Package + "." + row.Port
}

func kindName(kind Kind) string {
	if kind == KindPort {
		return "port"
	}
	return string(kind)
}
