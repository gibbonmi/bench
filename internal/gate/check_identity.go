package gate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	pathpkg "path"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const checkPolicyVersion = componentPolicyVersion + "/conformance-check-v1"

// checkIdentityMaterial is every authority input that makes one ordinary check's prior
// green reusable. Keeping the frame in one value makes field-omission mutations exhaustive.
type checkIdentityMaterial struct {
	Check          registry.Check
	Ordinal        int
	Inputs         []treeEntry
	Implementation []treeEntry
	CanaryFamilies []string
	CanaryInputs   []treeEntry
	Invocation     []string
}

func conformanceCheckIdentity(material checkIdentityMaterial) string {
	h := sha256.New()
	frameEach(h, "policy", checkPolicyVersion)
	frameEach(h, "name", material.Check.Name)
	frameEach(h, "tier", string(material.Check.Tier))
	frameEach(h, "ordinal", fmt.Sprint(material.Ordinal))
	frameEach(h, "implementation", material.Check.Implementation)
	frameEach(h, "subject", string(material.Check.Subject))
	frameEach(h, "input-source", string(material.Check.Inputs))
	frameEntries(h, "input", material.Inputs)
	frameEntries(h, "implementation-closure", material.Implementation)
	frameEach(h, "canary-family", material.CanaryFamilies...)
	frameEntries(h, "canary-input", material.CanaryInputs)
	frameEach(h, "invocation", material.Invocation...)
	return hex.EncodeToString(h.Sum(nil))
}

func frameEntries(h io.Writer, tag string, entries []treeEntry) {
	for _, entry := range entries {
		frameEach(h, tag+"-path", entry.Path)
		frameEach(h, tag+"-content", entry.Metadata)
	}
}

func resolveConformanceCheckIdentities(root string, tier registry.Tier, generation *treeGeneration) (map[string]string, error) {
	snapshot := generation.snapshot
	implementation := snapshotEntriesUnder(snapshot, "internal/conformance/")
	if len(implementation) == 0 {
		return nil, errors.New("conformance implementation closure is empty")
	}
	externalImplementations, err := resolveExternalImplementationEntries(root, generation)
	if err != nil {
		return nil, err
	}
	invocation := []string{
		registry.ConformanceTierEnv,
		registry.ConformanceCheckEnv,
		registry.ConformanceChecksEnv,
		registry.ConformanceInheritedEnv,
	}
	identities := map[string]string{}
	for ordinal, check := range registry.Checks {
		if check.Meta || !check.RunsAt(tier) {
			continue
		}
		inputs, err := resolveCheckInputsGeneration(root, check, generation)
		if err != nil {
			return nil, fmt.Errorf("identify conformance check %s: %w", check.Name, err)
		}
		families := registry.CanaryFamilies(check.Name)
		var canaryInputs []treeEntry
		for _, family := range families {
			canaryInputs = append(canaryInputs, snapshotEntriesUnder(snapshot, "tests/canary/"+family+"/")...)
		}
		checkImplementation := append([]treeEntry(nil), implementation...)
		checkImplementation = append(checkImplementation, externalImplementations[check.Implementation]...)
		identities[check.Name] = conformanceCheckIdentity(checkIdentityMaterial{
			Check:          check,
			Ordinal:        ordinal,
			Inputs:         inputs,
			Implementation: checkImplementation,
			CanaryFamilies: families,
			CanaryInputs:   canaryInputs,
			Invocation:     invocation,
		})
	}
	return identities, nil
}

func resolveExternalImplementationEntries(root string, generation *treeGeneration) (map[string][]treeEntry, error) {
	snapshot := generation.snapshot
	wanted := map[string]bool{}
	for _, check := range registry.Checks {
		wanted[check.Implementation] = true
	}
	packages := map[string]string{}
	for _, entry := range snapshot.entries {
		if strings.HasPrefix(entry.Path, "internal/conformance/") || !strings.HasSuffix(entry.Path, ".go") {
			continue
		}
		raw, err := generation.blob(entry)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(token.NewFileSet(), entry.Path, raw, 0)
		if err != nil {
			continue
		}
		for _, declaration := range file.Decls {
			fn, ok := declaration.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !wanted[fn.Name.Name] {
				continue
			}
			if _, duplicate := packages[fn.Name.Name]; duplicate {
				return nil, fmt.Errorf("conformance implementation %s is declared more than once", fn.Name.Name)
			}
			packages[fn.Name.Name] = pathpkg.Dir(entry.Path)
		}
	}
	resolved := map[string][]treeEntry{}
	if len(packages) == 0 {
		return resolved, nil
	}
	modulePath, err := snapshotModulePath(generation)
	if err != nil {
		return nil, err
	}
	for function, packageDir := range packages {
		closure, err := snapshotGoPackageClosure(generation, modulePath, packageDir)
		if err != nil {
			return nil, fmt.Errorf("resolve conformance implementation %s: %w", function, err)
		}
		resolved[function] = closure
	}
	return resolved, nil
}

func snapshotModulePath(generation *treeGeneration) (string, error) {
	snapshot := generation.snapshot
	entry, found := snapshot.entry("go.mod")
	if !found {
		return "", errors.New("module manifest is absent")
	}
	raw, err := generation.blob(entry)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "module" {
			modulePath := fields[1]
			if unquoted, err := strconv.Unquote(modulePath); err == nil {
				modulePath = unquoted
			}
			if modulePath != "" {
				return modulePath, nil
			}
		}
	}
	return "", errors.New("module path is unavailable")
}

// snapshotGoPackageClosure follows module-local imports from the package that declares an
// external check. Package content alone is not its implementation: a helper package can
// change the check's verdict without moving the declaring function.
func snapshotGoPackageClosure(generation *treeGeneration, modulePath, startDir string) ([]treeEntry, error) {
	snapshot := generation.snapshot
	queue := []string{startDir}
	seen := map[string]bool{}
	included := map[string]bool{}
	for len(queue) > 0 {
		dir := queue[0]
		queue = queue[1:]
		if seen[dir] {
			continue
		}
		seen[dir] = true
		packageEntries := snapshot.entries
		if dir != "." {
			packageEntries = snapshotEntriesUnder(snapshot, strings.TrimSuffix(dir, "/")+"/")
		}
		if len(packageEntries) == 0 {
			return nil, fmt.Errorf("module-local package %s is absent", dir)
		}
		for _, entry := range packageEntries {
			included[entry.Path] = true
			if pathpkg.Dir(entry.Path) != dir || !strings.HasSuffix(entry.Path, ".go") {
				continue
			}
			raw, err := generation.blob(entry)
			if err != nil {
				return nil, err
			}
			file, err := parser.ParseFile(token.NewFileSet(), entry.Path, raw, parser.ImportsOnly)
			if err != nil {
				return nil, fmt.Errorf("parse imports from %s: %w", entry.Path, err)
			}
			for _, spec := range file.Imports {
				importPath, err := strconv.Unquote(spec.Path.Value)
				if err != nil {
					return nil, fmt.Errorf("parse import in %s: %w", entry.Path, err)
				}
				if importPath == modulePath {
					queue = append(queue, ".")
				} else if strings.HasPrefix(importPath, modulePath+"/") {
					queue = append(queue, strings.TrimPrefix(importPath, modulePath+"/"))
				}
			}
		}
	}
	closure := make([]treeEntry, 0, len(included))
	for _, entry := range snapshot.entries {
		if included[entry.Path] {
			closure = append(closure, entry)
		}
	}
	return closure, nil
}

type conformanceCanaryIdentities struct {
	Shared string
	Bound  map[string]string
}

func resolveConformanceCanaryIdentitiesFromGeneration(root string, tier registry.Tier, generation *treeGeneration) (conformanceCanaryIdentities, error) {
	snapshot := generation.snapshot
	owners := map[string]string{}
	for _, check := range registry.Checks {
		if !check.Meta && check.RunsAt(tier) && len(registry.CanaryFamilies(check.Name)) > 0 {
			owners[check.Implementation] = check.Name
		}
	}
	bound := map[string]string{}
	shared := sha256.New()
	for _, entry := range snapshotEntriesUnder(snapshot, "internal/conformance/") {
		raw, err := generation.blob(entry)
		if err != nil {
			return conformanceCanaryIdentities{}, err
		}
		remaining := raw
		if strings.HasSuffix(entry.Path, ".go") {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, entry.Path, raw, 0)
			if err != nil {
				return conformanceCanaryIdentities{}, fmt.Errorf("parse %s: %w", entry.Path, err)
			}
			type sourceRange struct{ start, end int }
			var ownedRanges []sourceRange
			for _, declaration := range file.Decls {
				fn, ok := declaration.(*ast.FuncDecl)
				if !ok || fn.Recv != nil {
					continue
				}
				check, owned := owners[fn.Name.Name]
				if !owned {
					continue
				}
				start := fset.Position(fn.Pos()).Offset
				end := fset.Position(fn.End()).Offset
				if _, duplicate := bound[check]; duplicate || start < 0 || end > len(raw) || start >= end {
					return conformanceCanaryIdentities{}, fmt.Errorf("ambiguous conformance implementation %s", fn.Name.Name)
				}
				h := sha256.New()
				frameEach(h, "path", entry.Path)
				frameEach(h, "function", string(raw[start:end]))
				bound[check] = hex.EncodeToString(h.Sum(nil))
				ownedRanges = append(ownedRanges, sourceRange{start: start, end: end})
			}
			if len(ownedRanges) > 0 {
				var sharedSource bytes.Buffer
				cursor := 0
				for _, source := range ownedRanges {
					sharedSource.Write(raw[cursor:source.start])
					sharedSource.WriteString("<canary-owned-function>")
					cursor = source.end
				}
				sharedSource.Write(raw[cursor:])
				remaining = sharedSource.Bytes()
			}
		}
		frameEach(shared, "path", entry.Path)
		frameEach(shared, "content", string(remaining))
	}
	external, err := resolveExternalImplementationEntries(root, generation)
	if err != nil {
		return conformanceCanaryIdentities{}, err
	}
	for function, check := range owners {
		if bound[check] != "" || len(external[function]) == 0 {
			continue
		}
		h := sha256.New()
		frameEntries(h, "package", external[function])
		bound[check] = hex.EncodeToString(h.Sum(nil))
	}
	for function, check := range owners {
		if !isContentAddress(bound[check]) {
			return conformanceCanaryIdentities{}, fmt.Errorf("conformance implementation %s is unavailable", function)
		}
	}
	return conformanceCanaryIdentities{Shared: hex.EncodeToString(shared.Sum(nil)), Bound: bound}, nil
}

func resolveCheckInputsGeneration(root string, check registry.Check, generation *treeGeneration) ([]treeEntry, error) {
	snapshot := generation.snapshot
	var entries []treeEntry
	switch check.Inputs {
	case registry.InputCatchAll:
		entries = append(entries, snapshot.entries...)
	case registry.InputGoSource:
		entries = goSourceEntries(snapshot)
	case registry.InputGoAndDataHandling:
		entries = goSourceEntries(snapshot)
	default:
		if !check.Inputs.Valid() {
			return nil, fmt.Errorf("unknown input source %q", check.Inputs)
		}
	}
	files, directories := declaredCheckInputPaths(check.Inputs)
	for _, file := range files {
		entries = append(entries, declaredCheckEntry(snapshot, file)...)
	}
	for _, directory := range directories {
		entries = append(entries, snapshotEntriesUnder(snapshot, directory)...)
	}
	return resolveCheckSymlinks(root, generation, entries)
}

func declaredCheckInputPaths(source registry.InputSource) (files, directories []string) {
	switch source {
	case registry.InputGoAndDataHandling:
		files = []string{"DATA_HANDLING.md"}
	case registry.InputGateEntry:
		files = []string{".bench/gate.sh"}
	case registry.InputOfflineSmoke:
		files = []string{"scripts/smoke-offline.sh"}
	case registry.InputBenchRoutes:
		files = []string{"bin/bench.sh"}
	case registry.InputDecisionDocuments:
		directories = []string{"decisions/", "specs/"}
	case registry.InputBenchkitProfile:
		files = []string{"projects/benchkit.md"}
	}
	return files, directories
}

func goSourceEntries(snapshot treeSnapshot) []treeEntry {
	var entries []treeEntry
	for _, entry := range snapshot.entries {
		if strings.HasSuffix(entry.Path, ".go") || entry.Path == "go.mod" || entry.Path == "go.sum" {
			entries = append(entries, entry)
		}
	}
	return entries
}

func declaredCheckEntry(snapshot treeSnapshot, name string) []treeEntry {
	if entry, found := snapshot.entry(name); found {
		return []treeEntry{entry}
	}
	return []treeEntry{{Path: name, Metadata: "absent"}}
}

func resolveCheckSymlinks(root string, generation *treeGeneration, entries []treeEntry) ([]treeEntry, error) {
	resolved := append([]treeEntry(nil), entries...)
	for _, entry := range entries {
		targets, err := checkSymlinkTargets(root, generation, entry, map[string]bool{}, 0)
		if err != nil {
			return nil, fmt.Errorf("declared input %q: %w", entry.Path, err)
		}
		resolved = append(resolved, targets...)
	}
	return resolved, nil
}

func checkSymlinkTargets(root string, generation *treeGeneration, entry treeEntry, seen map[string]bool, depth int) ([]treeEntry, error) {
	snapshot := generation.snapshot
	fields := strings.Fields(entry.Metadata)
	if len(fields) != 3 || fields[0] != "120000" {
		return nil, nil
	}
	if depth >= 16 || seen[entry.Path] {
		return nil, errors.New("cyclic or over-deep symlink")
	}
	seen[entry.Path] = true
	defer delete(seen, entry.Path)
	raw, err := generation.blob(entry)
	if err != nil {
		return nil, errors.New("symlink target unavailable")
	}
	target := string(raw)
	if target == "" || strings.ContainsRune(target, '\x00') || pathpkg.IsAbs(target) {
		return nil, errors.New("symlink target is invalid")
	}
	canonical := pathpkg.Clean(pathpkg.Join(pathpkg.Dir(entry.Path), target))
	if canonical == ".." || strings.HasPrefix(canonical, "../") {
		return nil, errors.New("symlink target escapes the subject")
	}
	targetEntries := snapshotEntriesUnder(snapshot, canonical)
	if len(targetEntries) == 0 {
		targetEntries = snapshotEntriesUnder(snapshot, canonical+"/")
	}
	if len(targetEntries) == 0 {
		return nil, errors.New("symlink target is absent")
	}
	var resolved []treeEntry
	for _, targetEntry := range targetEntries {
		resolved = append(resolved, treeEntry{Path: entry.Path + "->" + targetEntry.Path, Metadata: targetEntry.Metadata})
		nested, err := checkSymlinkTargets(root, generation, targetEntry, seen, depth+1)
		if err != nil {
			return nil, err
		}
		for _, nestedEntry := range nested {
			resolved = append(resolved, treeEntry{Path: entry.Path + "->" + nestedEntry.Path, Metadata: nestedEntry.Metadata})
		}
	}
	return resolved, nil
}
