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

// seamDiags grades the variable, setter, runtime, and environment rows against the census
// of the whole production source. A seam in a package that no row names is still derived,
// so a new replacement cannot hide behind a remembered package list. A row whose package
// the graded tree does not hold is not graded, as for named ports.
func (c *seamCensus) seamDiags(tree *sourceTree, policy Policy) []string {
	rows := policy.Rows
	controls := c.deriveControls(policy.Families)
	diags := controls.diags
	registered := map[string]Row{}
	for _, row := range rows {
		registered[rowKey(row)] = row
	}
	for _, seam := range sortedKeys(c.variables) {
		row, found := registered[string(KindVariable)+" "+seam]
		if !found {
			diags = append(diags, fmt.Sprintf(
				"%s: %s is replaced at %s — add a variable row to injectedPortRegistry naming its production owner, consumer tests, failure purpose, and isolation evidence",
				UnregisteredSeamMessage, seam, c.variables[seam]))
			continue
		}
		diags = append(diags, c.variableRowDiags(row, seam)...)
	}
	for _, setter := range sortedKeys(c.setters) {
		canonical := c.setters[setter]
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
	diags = append(diags, c.controlDiags(controls, registered)...)
	for _, row := range rows {
		if row.Kind == KindPort || tree.packages[row.Package] == nil {
			continue
		}
		identity := row.Package + "." + row.Port
		_, variable := c.variables[identity]
		_, setter := c.setters[identity]
		_, branch := controls.runtime[identity]
		derived := map[Kind]bool{KindVariable: variable, KindSetter: setter, KindRuntime: branch, KindEnvironment: controls.hooks[row.Port]}
		if !derived[row.Kind] {
			diags = append(diags, orphanDiag(row))
		}
	}
	return diags
}

// controlDiags grades the runtime-branch and environment-hook rows. A runtime consumer
// reaches its branch through production calls. An environment consumer sets the hook
// itself or through a helper of its own test package. A consumer that sets it for a child
// process is accepted here, and review owns whether that child runs the owner. A consumer
// that sets it in its own process must also reach the owner through production calls.
func (c *seamCensus) controlDiags(controls *controlCensus, registered map[string]Row) []string {
	var diags []string
	for _, branch := range sortedKeys(controls.runtime) {
		row, found := registered[string(KindRuntime)+" "+branch]
		if !found {
			diags = append(diags, UnregisteredRuntimeMessage+": "+branch+" branches on "+controls.runtime[branch]+
				" — add a runtime row to injectedPortRegistry naming the protection's consumer tests, failure purpose, and isolation evidence")
			continue
		}
		diags = append(diags, rowEvidence(row, branch, true, func(file, test string) (bool, bool) {
			return c.reachesControl(file, test, true, func(facts *funcFacts) bool { return facts == c.bodies[branch] })
		})...)
	}
	for _, hook := range sortedKeys(controls.hooks) {
		readers := controls.readers[hook]
		row, found := registered[string(KindEnvironment)+" "+hook]
		if !found {
			first := sortedKeys(readers)[0]
			diags = append(diags, UnregisteredHookMessage+": "+hook+" is read by "+first+" at "+readers[first]+
				" and set by test source — add an environment row to injectedPortRegistry naming a reader, consumer tests, failure purpose, and isolation evidence")
			continue
		}
		_, owns := readers[row.Package+"."+row.Owner]
		owner := c.bodies[row.Package+"."+row.Owner]
		diags = append(diags, rowEvidence(row, hook, owns, func(file, test string) (bool, bool) {
			declared, child := c.reachesControl(file, test, false, func(facts *funcFacts) bool { return controls.children[facts][hook] })
			if !declared || child {
				return declared, child
			}
			_, setenv := c.reachesControl(file, test, false, func(facts *funcFacts) bool { return controls.setenvs[facts][hook] })
			_, reached := c.reachesControl(file, test, true, func(facts *funcFacts) bool { return facts == owner })
			return true, setenv && reached
		})...)
	}
	return diags
}

// variableRowDiags grades the policy facts of one admitted variable seam.
func (c *seamCensus) variableRowDiags(row Row, seam string) []string {
	return rowEvidence(row, seam, c.readers[seam][row.Owner], func(file, test string) (bool, bool) {
		return c.reaches(file, test, seam)
	})
}

// rowEvidence grades the evidence facts that every variable, runtime, and environment
// row carries: its owner, purpose, isolation, and consumers that reach the seam.
func rowEvidence(row Row, seam string, owns bool, reach func(file, test string) (declared, reached bool)) []string {
	var diags []string
	if !owns {
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
		declared, reached := reach(file, test)
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
	if row.Kind == KindEnvironment {
		return string(row.Kind) + " " + row.Port
	}
	return string(row.Kind) + " " + row.Package + "." + row.Port
}

func kindName(kind Kind) string {
	if kind == KindPort {
		return "port"
	}
	return string(kind)
}
