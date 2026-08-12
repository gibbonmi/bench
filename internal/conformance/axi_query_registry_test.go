package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var approvedAXIQueries = map[string][]string{
	"anchors": nil, "learnings": nil, "maps": nil, "guards": nil, "diff": nil, "coverage": nil,
	"worktree": {"list"},
}

var axiChildName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var profileAXICommand = regexp.MustCompile("`bench ([^`]+)`")
var guidanceAXICommand = regexp.MustCompile("(?m)^\\| `bench ([^`]+)` \\| ([^|\\n]+) \\|$")
var guidanceAXIPrinciple = regexp.MustCompile(`(?m)^([0-9]+)\. \*\*([^*\n]+)\.\*\*`)

var expectedAXIPrinciples = []string{
	"Token-efficient output",
	"Minimal default schemas",
	"Content truncation",
	"Pre-computed aggregates",
	"Definitive empty states",
	"Structured errors and honest exit codes",
	"Ambient context",
	"Content first",
	"Contextual disclosure",
	"Consistent way to get help",
}

type parsedAXIRegistry struct {
	members map[string][]string
}

func checkAXIQueryRegistry(root string) []string {
	reasons, err := axiReasonConstants(filepath.Join(root, "cmd", "bench", "command_registry.go"))
	if err != nil {
		return []string{"AXI query registry cannot read dispositions: " + err.Error()}
	}
	mainPath := filepath.Join(root, "cmd", "bench", "main.go")
	registry, err := parseAXIRegistry(mainPath, readIfExists(mainPath), reasons)
	if err != nil {
		return []string{"AXI query registry invalid: " + err.Error()}
	}
	var diags []string
	if !reflect.DeepEqual(registry.members, approvedAXIQueries) {
		diags = append(diags, fmt.Sprintf("AXI query registry declares %#v, want %#v", registry.members, approvedAXIQueries))
	}
	profilePath := filepath.Join(root, "projects", "benchkit.md")
	profileQueries, err := profileAXIQueries(readIfExists(profilePath))
	if err != nil {
		diags = append(diags, "AXI profile seam invalid: "+err.Error())
	} else if registryQueries := flattenAXIQueries(registry.members); !reflect.DeepEqual(profileQueries, registryQueries) {
		diags = append(diags, fmt.Sprintf("AXI profile seam advertises %q, registry declares %q", profileQueries, registryQueries))
	}
	guidancePath := filepath.Join(root, ".agents", "skills", "bench-craft-cli", "SKILL.md")
	diags = append(diags, checkAXIGuidance(readIfExists(guidancePath), registry.members)...)
	return uniqueSorted(diags)
}

func checkAXIGuidance(guidance string, registryMembers map[string][]string) []string {
	var diags []string
	principles := guidanceAXIPrinciple.FindAllStringSubmatch(guidance, -1)
	gotPrinciples := make([]string, 0, len(principles))
	for i, principle := range principles {
		wantNumber := fmt.Sprint(i + 1)
		if principle[1] != wantNumber {
			diags = append(diags, fmt.Sprintf("AXI guidance principle %d is numbered %s, want %s", i+1, principle[1], wantNumber))
		}
		gotPrinciples = append(gotPrinciples, principle[2])
	}
	if !reflect.DeepEqual(gotPrinciples, expectedAXIPrinciples) {
		diags = append(diags, fmt.Sprintf("AXI guidance principles are %q, want %q", gotPrinciples, expectedAXIPrinciples))
	}

	guidanceQueries, err := guidanceAXIQueries(guidance)
	if err != nil {
		diags = append(diags, "AXI guidance approved-query table invalid: "+err.Error())
	} else if registryQueries := flattenAXIQueries(registryMembers); !reflect.DeepEqual(guidanceQueries, registryQueries) {
		diags = append(diags, fmt.Sprintf("AXI guidance advertises %q, registry declares %q", guidanceQueries, registryQueries))
	}
	if strings.Contains(guidance, "--fields") {
		diags = append(diags, "AXI guidance advertises the forbidden field-selection flag")
	}
	required := []string{
		"`help[N]{cmd,why}:`",
		"`help[0]{cmd,why}:`",
		"per matching row",
		"one check action per mapped coverage row",
		"Successful default extraction",
		"repairable mapped rows",
		"that exact command is the retry",
		"`retry after repairing coverage map`",
		"`coverage --check` and every refusal retain their error contracts and append no disclosure",
		"state-derived action",
		"stable source order",
		"known argument",
		"unknown future-input slots",
		"Terminal results offer no busywork",
		"`--help`, `-h`, and bare `help`",
	}
	contractText := strings.Join(strings.Fields(guidance), " ")
	for _, phrase := range required {
		if !strings.Contains(contractText, phrase) {
			diags = append(diags, fmt.Sprintf("AXI guidance omits required contract %q", phrase))
		}
	}
	return diags
}

func guidanceAXIQueries(guidance string) ([]string, error) {
	const marker = "## Bench application"
	if strings.Count(guidance, marker) != 1 {
		return nil, fmt.Errorf("found %d Bench application sections, want exactly 1", strings.Count(guidance, marker))
	}
	section := guidance[strings.Index(guidance, marker)+len(marker):]
	if end := strings.Index(section, "\n## "); end >= 0 {
		section = section[:end]
	}
	const tableStart = "| approved query | contextual disclosure |\n| --- | --- |\n"
	if strings.Count(section, tableStart) != 1 {
		return nil, fmt.Errorf("found %d approved-query tables, want exactly 1", strings.Count(section, tableStart))
	}
	rows := strings.Split(section[strings.Index(section, tableStart)+len(tableStart):], "\n")
	if len(rows) == 0 || !strings.HasPrefix(rows[0], "|") {
		return nil, fmt.Errorf("approved-query table has no rows")
	}
	seen := make(map[string]bool, len(rows))
	queries := make([]string, 0, len(rows))
	for i, row := range rows {
		if !strings.HasPrefix(row, "|") {
			break
		}
		match := guidanceAXICommand.FindStringSubmatch(row)
		if match == nil {
			return nil, fmt.Errorf("approved-query table row %d is malformed", i+1)
		}
		query := match[1]
		if seen[query] {
			return nil, fmt.Errorf("approved-query table repeats %q", query)
		}
		seen[query] = true
		if strings.TrimSpace(match[2]) == "" {
			return nil, fmt.Errorf("approved-query table gives %q no application", query)
		}
		queries = append(queries, query)
	}
	sort.Strings(queries)
	return queries, nil
}

func TestAXIGuidanceContractBites(t *testing.T) {
	h := NewHarness(t)
	guidance := h.ReadRootFile(".agents", "skills", "bench-craft-cli", "SKILL.md")
	if diags := checkAXIGuidance(guidance, approvedAXIQueries); len(diags) != 0 {
		t.Fatalf("live guidance is not a valid mutation base: %v", diags)
	}
	const anchorsCell = "| `bench anchors` |"
	mutations := []struct {
		name, old, replacement, want string
	}{
		{"removed principle", "8. **Content first.**", "", "principles are"},
		{"renumbered principle", "8. **Content first.**", "11. **Content first.**", "numbered 11"},
		{"omitted member", "| `bench worktree list` |", "<!-- omitted -->", "guidance advertises"},
		{"additional member", anchorsCell, "| `bench status` | Report the operational dashboard. |\n" + anchorsCell, "guidance advertises"},
		{"malformed member", anchorsCell, "| bench anchors |", "row 1 is malformed"},
		{"forbidden option", "Minimal default schemas", "Minimal default schemas --fields", "forbidden field-selection"},
		{"missing help contract", "`help[0]{cmd,why}:`", "`help[0]:`", "omits required contract"},
		{"unchecked coverage classifier", "one check action per mapped coverage row", "one check action per unchecked coverage row", "omits required contract"},
		{"coverage retry moved to refusal", "Successful default extraction", "`coverage --check` refusal", "omits required contract"},
		{"coverage retry why changed", "`retry after repairing coverage map`", "`retry after checking coverage map`", "omits required contract"},
		{"coverage refusal gains disclosure", "`coverage --check` and every refusal retain their error contracts and append no disclosure", "`coverage --check` and every refusal retain their error contracts and append disclosure", "omits required contract"},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if strings.Count(guidance, mutation.old) != 1 {
				t.Fatalf("mutation anchor %q count = %d, want 1", mutation.old, strings.Count(guidance, mutation.old))
			}
			mutated := strings.Replace(guidance, mutation.old, mutation.replacement, 1)
			if diags := checkAXIGuidance(mutated, approvedAXIQueries); !diagnosticsContain(diags, mutation.want) {
				t.Fatalf("mutation diagnostics = %v, want substring %q", diags, mutation.want)
			}
		})
	}
}

func diagnosticsContain(diags []string, want string) bool {
	for _, diag := range diags {
		if strings.Contains(diag, want) {
			return true
		}
	}
	return false
}

func TestAXIMembershipExpectationBitesInBothDirections(t *testing.T) {
	h := NewHarness(t)
	body := h.ReadRootFile("cmd", "bench", "main.go")
	reasons := mustAXIReasonConstants(t, h.RootPath("cmd", "bench", "command_registry.go"))
	mutations := []struct {
		name, old, replacement string
	}{
		{"approved root removed", `{Name: "anchors", AXI: axiApprovedRoot`, `{Name: "anchors", AXI: axiExempt(axiReasonOperational)`},
		{"operational root admitted", `{Name: "status", AXI: axiExempt(axiReasonOperational)`, `{Name: "status", AXI: axiApprovedRoot`},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if strings.Count(body, mutation.old) != 1 {
				t.Fatalf("mutation anchor %q count = %d, want 1", mutation.old, strings.Count(body, mutation.old))
			}
			mutated := strings.Replace(body, mutation.old, mutation.replacement, 1)
			registry, err := parseAXIRegistry("mutated-main.go", mutated, reasons)
			if err != nil {
				t.Fatal(err)
			}
			if reflect.DeepEqual(registry.members, approvedAXIQueries) {
				t.Fatal("membership mutation remained green against the independent expectation")
			}
		})
	}
}

func TestAXIRegistryParserFailsClosed(t *testing.T) {
	entry := func(disposition string) string {
		return "package main\nvar commandRegistry = []commandDefinition{{Name: \"x\", AXI: " + disposition + "}}\n"
	}
	cases := []struct {
		name, body, want string
	}{
		{"absent registry", "package main\n", "found 0 commandRegistry declarations"},
		{"multiple registries", entry("axiApprovedRoot") + "var commandRegistry = []commandDefinition{}\n", "found 2 commandRegistry declarations"},
		{"missing disposition", "package main\nvar commandRegistry = []commandDefinition{{Name: \"x\"}}\n", `command "x" has 0 AXI dispositions`},
		{"conflicting dispositions", "package main\nvar commandRegistry = []commandDefinition{{Name: \"x\", AXI: axiApprovedRoot, AXI: axiApprovedChildren(\"list\")}}\n", `command "x" has 2 AXI dispositions`},
		{"malformed disposition", entry("true"), `command "x" has malformed AXI disposition`},
		{"empty exemption", entry(`axiExempt("")`), `command "x" has an empty AXI exemption`},
		{"unknown exemption", entry(`axiExempt(missingReason)`), `command "x" has an unresolved AXI exemption reason`},
		{"empty child set", entry(`axiApprovedChildren()`), `command "x" has an empty approved child set`},
		{"invalid child", entry(`axiApprovedChildren("bad child")`), `command "x" has invalid AXI child`},
		{"duplicate child", entry(`axiApprovedChildren("list", "list")`), `command "x" repeats AXI child`},
		{"missing name", "package main\nvar commandRegistry = []commandDefinition{{AXI: axiApprovedRoot}}\n", "registry entry 1 has 0 Name fields"},
		{"duplicate command", "package main\nvar commandRegistry = []commandDefinition{{Name: \"x\", AXI: axiApprovedRoot}, {Name: \"x\", AXI: axiApprovedRoot}}\n", `commandRegistry repeats command "x"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseAXIRegistry("fixture.go", tc.body, map[string]string{"emptyReason": ""})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("parse error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func parseAXIRegistry(path, body string, reasons map[string]string) (parsedAXIRegistry, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, body, 0)
	if err != nil {
		return parsedAXIRegistry{}, err
	}
	var registries []ast.Expr
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range value.Names {
				if name.Name == "commandRegistry" {
					if len(value.Names) != 1 || len(value.Values) != 1 || i >= len(value.Values) {
						return parsedAXIRegistry{}, fmt.Errorf("commandRegistry must be one named value with one literal")
					}
					registries = append(registries, value.Values[i])
				}
			}
		}
	}
	if len(registries) != 1 {
		return parsedAXIRegistry{}, fmt.Errorf("found %d commandRegistry declarations, want exactly 1", len(registries))
	}
	literal, ok := registries[0].(*ast.CompositeLit)
	if !ok {
		return parsedAXIRegistry{}, fmt.Errorf("commandRegistry is not a composite literal")
	}
	result := parsedAXIRegistry{members: make(map[string][]string)}
	seen := make(map[string]bool, len(literal.Elts))
	for i, element := range literal.Elts {
		entry, ok := element.(*ast.CompositeLit)
		if !ok {
			return parsedAXIRegistry{}, fmt.Errorf("registry entry %d is not a composite literal", i+1)
		}
		names := registryFields(entry, "Name")
		if len(names) != 1 {
			return parsedAXIRegistry{}, fmt.Errorf("registry entry %d has %d Name fields, want exactly 1", i+1, len(names))
		}
		name, ok := stringLiteral(names[0])
		if !ok || name == "" {
			return parsedAXIRegistry{}, fmt.Errorf("registry entry %d has malformed or empty Name", i+1)
		}
		if seen[name] {
			return parsedAXIRegistry{}, fmt.Errorf("commandRegistry repeats command %q", name)
		}
		seen[name] = true
		dispositions := registryFields(entry, "AXI")
		if len(dispositions) != 1 {
			return parsedAXIRegistry{}, fmt.Errorf("command %q has %d AXI dispositions, want exactly 1", name, len(dispositions))
		}
		approved, children, err := parseAXIDisposition(name, dispositions[0], reasons)
		if err != nil {
			return parsedAXIRegistry{}, err
		}
		if approved {
			result.members[name] = children
		}
	}
	return result, nil
}

func registryFields(entry *ast.CompositeLit, name string) []ast.Expr {
	var values []ast.Expr
	for _, element := range entry.Elts {
		pair, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := pair.Key.(*ast.Ident)
		if ok && key.Name == name {
			values = append(values, pair.Value)
		}
	}
	return values
}

func parseAXIDisposition(command string, expr ast.Expr, reasons map[string]string) (bool, []string, error) {
	if ident, ok := expr.(*ast.Ident); ok && ident.Name == "axiApprovedRoot" {
		return true, nil, nil
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false, nil, fmt.Errorf("command %q has malformed AXI disposition", command)
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil, fmt.Errorf("command %q has malformed AXI disposition", command)
	}
	switch fn.Name {
	case "axiApprovedChildren":
		if len(call.Args) == 0 {
			return false, nil, fmt.Errorf("command %q has an empty approved child set", command)
		}
		seen := make(map[string]bool, len(call.Args))
		children := make([]string, 0, len(call.Args))
		for _, arg := range call.Args {
			child, ok := stringLiteral(arg)
			if !ok || !axiChildName.MatchString(child) {
				return false, nil, fmt.Errorf("command %q has invalid AXI child", command)
			}
			if seen[child] {
				return false, nil, fmt.Errorf("command %q repeats AXI child %q", command, child)
			}
			seen[child] = true
			children = append(children, child)
		}
		return true, children, nil
	case "axiExempt":
		if len(call.Args) != 1 {
			return false, nil, fmt.Errorf("command %q has malformed AXI exemption", command)
		}
		reason, ok := stringLiteral(call.Args[0])
		if !ok {
			ident, isIdent := call.Args[0].(*ast.Ident)
			if !isIdent {
				return false, nil, fmt.Errorf("command %q has malformed AXI exemption reason", command)
			}
			reason, ok = reasons[ident.Name]
			if !ok {
				return false, nil, fmt.Errorf("command %q has an unresolved AXI exemption reason", command)
			}
		}
		if strings.TrimSpace(reason) == "" {
			return false, nil, fmt.Errorf("command %q has an empty AXI exemption", command)
		}
		return false, nil, nil
	default:
		return false, nil, fmt.Errorf("command %q has malformed AXI disposition", command)
	}
}

func axiReasonConstants(path string) (map[string]string, error) {
	body := readIfExists(path)
	if body == "" {
		return nil, fmt.Errorf("%s is absent or empty", path)
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, body, 0)
	if err != nil {
		return nil, err
	}
	reasons := map[string]string{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range value.Names {
				if !strings.HasPrefix(name.Name, "axiReason") {
					continue
				}
				if i >= len(value.Values) {
					continue
				}
				if reason, ok := stringLiteral(value.Values[i]); ok {
					reasons[name.Name] = reason
				}
			}
		}
	}
	return reasons, nil
}

func mustAXIReasonConstants(t *testing.T, path string) map[string]string {
	t.Helper()
	reasons, err := axiReasonConstants(path)
	if err != nil {
		t.Fatal(err)
	}
	return reasons
}

func flattenAXIQueries(members map[string][]string) []string {
	var queries []string
	for root, children := range members {
		if children == nil {
			queries = append(queries, root)
			continue
		}
		for _, child := range children {
			queries = append(queries, root+" "+child)
		}
	}
	sort.Strings(queries)
	return queries
}

func profileAXIQueries(profile string) ([]string, error) {
	const marker = "- **The AXI query surface**"
	start := strings.Index(profile, marker)
	if start < 0 {
		return nil, fmt.Errorf("project profile has no AXI query surface seam")
	}
	paragraph := profile[start:]
	const endMarker = ", and the\n  shared flat-table"
	if end := strings.Index(paragraph, endMarker); end >= 0 {
		paragraph = paragraph[:end]
	} else {
		return nil, fmt.Errorf("project profile AXI seam has no approved-list terminator")
	}
	matches := profileAXICommand.FindAllStringSubmatch(paragraph, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("project profile AXI seam names no bench queries")
	}
	seen := make(map[string]bool, len(matches))
	queries := make([]string, 0, len(matches))
	for _, match := range matches {
		query := match[1]
		if seen[query] {
			return nil, fmt.Errorf("project profile AXI seam repeats %q", query)
		}
		seen[query] = true
		queries = append(queries, query)
	}
	sort.Strings(queries)
	return queries, nil
}
