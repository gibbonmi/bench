# gate-test-concurrency implementation evidence

This evidence grades spec-build run
`8851335a0084158b50aea39c5715de18d24cdf243727fc94ae0b44a57ec340ab`.
The repaired candidate base is commit
`07047319c93cc556aa37345cdc38a0452630f387`. The measured Go execution
subject is prospective tree `d1c2537c0d91656d5453c01a7bb9577eac3e05e7`:
the base candidate with `internal/gate` subtree
`d7a594bba8d648ad22759e681467e59eb80aebd6` and runner-test blob
`e64197acbd9f949d3a24bb5c7daf0595f1c56645`. Evidence-only Markdown does
not enter that Go execution subject. All current measurements below use
exactly those Go bytes.

The reference host is WSL2 Linux/amd64, Go 1.25.0, on a 13th Gen Intel
Core i7-13620H with 12 logical CPUs, six cores, and two threads per core.
Every command ran alone; no two tests or measurements overlapped.

## Exact census

An AST audit over every package-level `Test*` function in
`internal/gate/*_test.go` reports:

- Top-level tests: 245
- First-statement `t.Parallel` calls: 192
- Structurally serial tests: 53

The audit requires each eligible test to contain exactly one `t.Parallel`
call as its first statement, each serial test to contain none, and the set
below to equal the discovered serial set exactly.

## Complete serial inventory

| test | load-bearing reason |
|---|---|
| `TestClosedSubjectReuseAndMutations` | Mutates the process environment with `t.Setenv`. |
| `TestComposedGreenAcceptsOnlyCompleteExactTipEvidence` | Mutates the process environment with `t.Setenv`. |
| `TestExecuteDeadlineRecordsDistinctTimeout` | Swaps the package-level `gateTimeout` variable. |
| `TestExecuteHealthyGateJustBelowDeadlineRecordsOrdinaryGreen` | Swaps the package-level `gateTimeout` variable. |
| `TestExecuteReusingFreshGreenAnswersFromEvaluationOwnedGeneration` | Mutates the process environment with `t.Setenv`. |
| `TestFT78Story3ProofLedger` | Runs a registered heterogeneous subtest ledger whose drivers include process-global mutations. |
| `TestFT78Story4ProofLedger` | Runs a registered heterogeneous subtest ledger whose drivers include process-global mutations. |
| `TestFreshFlagForcesARealRunPastAReusableGreen` | Mutates the process environment with `t.Setenv`. |
| `TestGateEnvStripsWrapperRoutingInternals` | Mutates the process environment with `t.Setenv`. |
| `TestGateEnvironmentIsPasslisted` | Mutates the process environment with `t.Setenv`. |
| `TestGateEvaluationBoundsOrdinarySourceWorkAcrossAllIdentityFamilies` | Mutates the process environment with `t.Setenv`. |
| `TestGateEvaluationBoundsProspectiveSourceWorkAndRunsFullComponentInventory` | Mutates the process environment with `t.Setenv`. |
| `TestGateEvaluationProspectiveValidationRejectsCheckoutDriftWithoutMaterializing` | Mutates the process environment with `t.Setenv`. |
| `TestGateGoConformanceSuitePreservesCache` | Mutates the process environment with `t.Setenv`. |
| `TestGateGoCoreTestUsesFreshVerdict` | Mutates the process environment through both testing and OS environment APIs. |
| `TestGateGoStepReportsASpawnFailure` | Mutates the process environment with `t.Setenv`. |
| `TestGateGoWithoutRootOutsideARepo` | Changes the process working directory with `t.Chdir`. |
| `TestGateRunDeadlineTermGraceThenKill` | Swaps the package-level `gateTimeout` variable. |
| `TestInjectedKitSurvivesHostileAmbientThroughEvaluationAndComposition` | Mutates `BENCH_KIT` with `t.Setenv`. |
| `TestInnerCanarySingularSelectionRemovesPluralSelection` | Mutates the process environment with `t.Setenv`. |
| `TestManifestDanglingSymlinkIsRed` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestManifestDirResolvesAgainstGradedRoot` | Swaps `benchkitPhasesForCommand` and mutates the process environment. |
| `TestManifestEmptyIsRed` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestManifestEntryLimitOverrideIsTightenOnly` | Mutates the process environment through testing and OS APIs. |
| `TestManifestMalformed` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestManifestSpecialFileIsRed` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestPhaseTableNoToolchainNoPhases` | Mutates the process environment with `t.Setenv`. |
| `TestPhasesCommandAbsentManifestFallsBack` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestPhasesCommandEmptyKitFallsBackToRoot` | Swaps `benchkitPhasesForCommand` and mutates `BENCH_KIT`. |
| `TestPhasesCommandInnerModeManifestDagOrder` | Swaps `benchkitPhasesForCommand` and mutates the process environment. |
| `TestPhasesCommandLoadsManifest` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestPhasesCommandManifestFieldsEndToEnd` | Swaps the package-level `benchkitPhasesForCommand` provider. |
| `TestPhasesCommandNonEmptyKitSelectsBuiltinTableAtKit` | Swaps `benchkitPhasesForCommand` and mutates `BENCH_KIT`. |
| `TestPhasesCommandRoutesCanaryToOwningPhase` | Swaps `benchkitPhasesForCommand` and mutates the process environment. |
| `TestPhasesCommandSignalHelper` | Is an `os.Exit` re-exec helper and reaches the package-level phase provider. |
| `TestPhasesCommandStragglerHelper` | Is an `os.Exit` re-exec helper and reaches the package-level phase provider. |
| `TestPhasesCommandVetPhaseReds` | Mutates the process environment with `t.Setenv`. |
| `TestPhasesForModeAcceptsPhaseTableNames` | Mutates the process environment with `t.Setenv`. |
| `TestPhasesForModeFallsBackForAbsentOwner` | Mutates the process environment with `t.Setenv`. |
| `TestPinCommandDeclineAndSecondPinOverwrite` | Changes the process working directory through testing and OS APIs. |
| `TestPinCommandNotInRepo` | Changes the process working directory through testing and OS APIs. |
| `TestPinCommandWritesCommittedBenchTree` | Changes the process working directory through testing and OS APIs. |
| `TestProspectiveFullInventoryHelper` | Is the environment-selected prospective re-exec helper. |
| `TestPublicDocumentClassesProjectTheirExactCheckPartition` | Mutates the process environment with `t.Setenv`. |
| `TestRetiredContractPackageExportChangesNothing` | Mutates the process environment with `t.Setenv`. |
| `TestRunnerFinalLineAndExitCodes` | Changes the process working directory in its runner subtests. |
| `TestRunnerOptionalBrokenSymlinkSkips` | Mutates the process environment with `t.Setenv`. |
| `TestRunnerOptionalUnexecutableStubGoesRed` | Mutates the process environment with `t.Setenv`. |
| `TestRunnerPhaseEnvStripsThenSets` | Mutates the process environment with `t.Setenv`. |
| `TestRunnerRunsPhasesConcurrently` | Measures `runPhases` overlap; scheduler overlap would make its timing assertion unattributable. |
| `TestStrictCapabilityFailure` | Mutates `BENCH_REQUIRE_CAPABILITIES` with `t.Setenv`. |
| `TestUnlockErrorClearsSameProcessOwnership` | Mutates the package-level `executionLockOwners` map. |
| `TestUnreadableSkipLogIsFatalOnlyUnderStrictMode` | Mutates the process environment with `t.Setenv`. |

## Same-command timing

The command for every row in both samples is exactly `go test -count=1
./internal/gate`; its default Go timeout is 10 minutes. The original serial
subject is commit `14529f4c547f3b3815a03c5b48aa732e829b6b05`, tree
`42eb5d799be1053a2505bf1fa743c979bd9b9a14`. The after subject is the
prospective tree named above. Width is the host default in both samples:
`GOMAXPROCS` and `-parallel` were not overridden.

| sample | run | subject | exit | package time | wall time | max RSS | filesystem outputs |
|---|---:|---|---:|---:|---:|---:|---:|
| before | 1 | `14529f4c…` | 0 | 149.802 s | 151.06 s | not retained | not retained |
| before | 2 | `14529f4c…` | 0 | 139.741 s | 139.99 s | not retained | not retained |
| before | 3 | `14529f4c…` | 0 | 150.556 s | 150.85 s | not retained | not retained |
| after | 1 | `d1c2537c…` | 0 | 56.621 s | 56.95 s | 231,228 KiB | 1,359,816 blocks |
| after | 2 | `d1c2537c…` | 0 | 56.261 s | 56.72 s | 226,736 KiB | 1,264,936 blocks |
| after | 3 | `d1c2537c…` | 0 | 54.624 s | 55.15 s | 234,392 KiB | 1,264,896 blocks |

The before wall median is 150.85 seconds. The repaired exact-candidate wall
median is 56.72 seconds, below the 90-second target. The previously retained
TP candidate sample was 61.55, 63.61, and 55.40 seconds (61.55-second
median); the new sample supersedes it for the repaired candidate.

## Focused validation and local cost

Every row uses the repaired Go subject `d1c2537c…`, exit 0, a 600-second
timeout, and `/usr/bin/time -v`. Width is stated explicitly; `default` means
neither `GOMAXPROCS` nor `-parallel` was overridden.

| check | exact command | width | wall time | max RSS | filesystem outputs |
|---|---|---|---:|---:|---:|
| race | `go test -race -count=1 -timeout 600s ./internal/gate` | default | 69.99 s | 230,712 KiB | 1,321,920 blocks |
| two-core | `GOMAXPROCS=2 go test -count=1 -timeout 600s ./internal/gate` | two Go processors | 97.88 s | 238,488 KiB | 1,262,592 blocks |
| width one | `go test -count=1 -parallel=1 -timeout 600s ./internal/gate` | one test | 138.01 s | 226,172 KiB | 1,262,936 blocks |
| width two | `go test -count=1 -parallel=2 -timeout 600s ./internal/gate` | two tests | 100.98 s | 223,732 KiB | 1,262,720 blocks |
| hostile ambient | `BENCH_KIT=/nonexistent go test -count=1 -timeout 600s ./internal/gate` | default | 55.92 s | 238,248 KiB | 1,264,528 blocks |

GNU time reports filesystem-output blocks; at 512 bytes per block these runs
wrote about 617–664 MiB apiece. Width two reduces wall time relative to width
one without reducing writes. The remaining local-development bottleneck is
therefore fixture materialization/build output, not a reason to widen test
parallelism beyond two.

## Mutation ledger

Each entry records a green baseline, the applied mutation, its attributable
red, and exact restoration. Checkpoint tree IDs are the restored assignment
trees retained by the spec-build lifecycle.

RPC1's complete prospective-carrier source audit is the independently authored,
bounded Go source below. It walks every production Go file under `internal/gate`,
enumerates each `gateEvaluation` composite literal whose `prospective` field is
the literal `true`, reports its containing function, and exits nonzero when the
same literal has any `kit` field regardless of its assigned expression. The
source was written to `/tmp/rpc1-prospective-carrier-audit.go`; the exact command
was `gofmt -w /tmp/rpc1-prospective-carrier-audit.go && go run
/tmp/rpc1-prospective-carrier-audit.go internal/gate`.

```go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type prospectiveLiteral struct {
	file     string
	line     int
	function string
	kitFile  string
	kitLine  int
}

func field(literal *ast.CompositeLit, name string) *ast.KeyValueExpr {
	for _, elt := range literal.Elts {
		keyed, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := keyed.Key.(*ast.Ident)
		if ok && ident.Name == name {
			return keyed
		}
	}
	return nil
}

func isGateEvaluation(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "gateEvaluation"
}

func isTrue(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "true"
}

func main() {
	root := "internal/gate"
	if len(os.Args) == 2 {
		root = os.Args[1]
	} else if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run audit.go [internal/gate]")
		os.Exit(2)
	}

	var files []string
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	sort.Strings(files)

	fset := token.NewFileSet()
	var literals []prospectiveLiteral
	for _, path := range files {
		parsed, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				literal, ok := node.(*ast.CompositeLit)
				if !ok || !isGateEvaluation(literal.Type) {
					return true
				}
				prospective := field(literal, "prospective")
				if prospective == nil || !isTrue(prospective.Value) {
					return true
				}
				position := fset.Position(literal.Pos())
				result := prospectiveLiteral{file: position.Filename, line: position.Line, function: fn.Name.Name}
				if kit := field(literal, "kit"); kit != nil {
					kitPosition := fset.Position(kit.Pos())
					result.kitFile = kitPosition.Filename
					result.kitLine = kitPosition.Line
				}
				literals = append(literals, result)
				return true
			})
		}
	}

	if len(literals) == 0 {
		fmt.Fprintln(os.Stderr, "RPC1 audit: no prospective gateEvaluation literals found")
		os.Exit(2)
	}
	sort.Slice(literals, func(i, j int) bool {
		if literals[i].file == literals[j].file {
			return literals[i].line < literals[j].line
		}
		return literals[i].file < literals[j].file
	})
	carriers := 0
	for _, literal := range literals {
		kit := "absent"
		if literal.kitFile != "" {
			kit = fmt.Sprintf("%s:%d", literal.kitFile, literal.kitLine)
			carriers++
		}
		fmt.Printf("RPC1 prospective literal: %s:%d function=%s kit=%s\n", literal.file, literal.line, literal.function, kit)
	}
	if carriers != 0 {
		fmt.Printf("RPC1 audit: stored kit in %d of %d prospective gateEvaluation literals\n", carriers, len(literals))
		os.Exit(1)
	}
	fmt.Printf("RPC1 audit: green prospective_literals=%d stored_kit=0 production_files=%d\n", len(literals), len(files))
}
```

| fact | green baseline | mutation and attributable red | exact restoration |
|---|---|---|---|
| IK1/component-identity | Subject `0c432fcb…`; `go test -count=1 ./internal/gate ./cmd/bench`, exit 0. | Re-read ambient `BENCH_KIT` in component identity; `BENCH_KIT=/nonexistent go test -count=1 -run '^TestInjectedKitSurvivesHostileAmbientThroughEvaluationAndComposition$' ./internal/gate`, exit 1 with differing identities. | Restored exact bytes; focused command exit 0; checkpoint tree `0c432fcb9ef114e76323c502d11c7d5e3871d4ca`. |
| IK2/empty-fallback | Subject `0c432fcb…`; focused empty-fallback test exit 0. | Accepted the empty environment value as kit; `go test -count=1 -run '^TestPhasesCommandEmptyKitFallsBackToRoot$' ./internal/gate`, exit 1 on the fallback assertion. | Restored exact predicate; same focused command exit 0; checkpoint tree `0c432fcb…`. |
| IK4/subject-strip-wrapper | Subject `0c432fcb…`; focused subject-environment test exit 0. | Stopped stripping `BENCH_WRAPPER`; `go test -count=1 -run '^TestGateEnvStripsWrapperRoutingInternals$' ./internal/gate`, exit 1 on wrapper leakage. | Restored exact bytes; same focused command exit 0; checkpoint tree `0c432fcb…`. |
| RP1/entry-migration | Subject `519d812d…`; `BENCH_KIT=/nonexistent go test -count=1 ./internal/gate`, exit 0, package 142.242 s, wall 142.55 s. | Routed `TestPhasesCommandAbsentManifestFallsBack` through exported `PhasesCommand` without a pin; hostile focused run had exit 1 with `PhasesCommand = 1` after built-in phase spawn failure, package 0.003 s, wall 1.09 s. | Restored `manifest_test.go` SHA-256 `ea76208ebf40be80ca8d692991b507eed5ee6d1b810f6f4afb0a0e116b485054`; hostile package command exit 0; checkpoint tree `519d812d0de967ada3d27eb39fe4bc6595d381ff`. |
| RP1/injected-evaluation | Subject `519d812d…`; hostile ambient guard exit 0. | Recoupled injected evaluation to hostile ambient kit; the focused witness had exit 1 across the expected fixture families. | Exact restoration; focused witness exit 0 in 1.0 s; checkpoint tree `519d812d…`. |
| TP1/serial-list | Subject `7be6308f…`; `go test -count=1 -run '^TestStrictCapabilityFailure$' ./internal/gate`, exit 0, package 0.002 s. | Added `t.Parallel` to the `t.Setenv`-pinning test; the same command had exit 1 with `testing: test using t.Setenv or t.Chdir can not use t.Parallel`. | Restored exact bytes; same focused command exit 0 at 0.002 s; checkpoint tree `7be6308f048666ee0622928f1358cefc59961b8e`. |
| RKW1/distinct-kit-entry | Subject `f0708759…`; `go test -count=1 -run '^TestPhasesCommandNonEmptyKitSelectsBuiltinTableAtKit$' ./internal/gate`, exit 0. | Made `kitRoot` ignore non-empty `BENCH_KIT`; the focused command had exit 1 because root `001` reached the provider instead of kit `002`. | Production SHA-256 restored to `816fe9730903d2515f64ce4158ab6f25ce3bd04ec2d4c43665197c142193c4f2`; same focused command exit 0; checkpoint tree `f0708759d977ffe4ad7d1b287fd70306aec8cb0f`. |
| RAK1/absent-tool-root | Subject `afa1c560…`; exact source audit and `GOMAXPROCS=2 go test -p 1 -count=1 -run '^TestAbsentOptionalToolLeavesNoSlot$' ./internal/gate`, exit 0 in 0.724 s. | Restored `phaseTable(fixture.root, kitRoot(fixture.root))`; the exact source audit had exit 1 and named the redundant ambient read. | Restored `phaseTable(fixture.root, fixture.root)` exactly; audit and focused test exit 0; checkpoint tree `afa1c56017a929db0f35db838cd731f80de2e909`. |
| RDW1/dead-wrapper-removal | Subject `298e15cd…`; zero-caller audit exit 0 and focused evaluation set exit 0 in 3.160 s. | Reintroduced `newEngineEvaluation`; the zero-caller audit had exit 1 and named its declaration. | Exact restoration; both removed wrappers absent, remaining explicit functions have callers, and the focused set had exit 0; checkpoint tree `298e15cd3f9d1252507dfbff1cbb09602ec120f1`. |
| RPS1/prospective-seam | Subject `1a3465b7…`; zero-`executeTreeAtKit`-carrier audit and focused `TestExecuteTree` set exit 0, final run 6.561 s. | Reintroduced `executeTreeAtKit` with an unconsumed kit parameter; the engine-seam source audit had exit 1 while compile-only remained green in 0.221 s. | Exact restoration; zero `executeTreeAtKit` carrier hits and focused tests exit 0; checkpoint tree `1a3465b7324d0be6a1021121e1647cd2cd69843c`. |
| RPC1/prospective-constructor | Base `77e98c42d94c92409e91df38e2b551f0df2d0ac4`; repaired production tree `3fe6de7f3709bed39c53faeeca3e0c50edaab820`, `internal/gate` tree `6276e46a8e2eea525abeda54f9a284d1509828ed`, and `evaluation.go` blob `70fec55267a33ec13882560d77d74c51e7508a0e`. `GOMAXPROCS=2 go test -p 1 -parallel 1 -count=1 -run '^(TestGateEvaluationBoundsProspectiveSourceWorkAndRunsFullComponentInventory|TestGateEvaluationProspectiveValidationRejectsCheckoutDriftWithoutMaterializing)$' ./internal/gate` exited 0 in 0.290 s before repair. The semantic audit exited 0 after repair with `function=newProspectiveTreeEvaluation kit=absent` and `green prospective_literals=1 stored_kit=0 production_files=25`. | Mutation 1 restored the four-argument wrapper, forwarder, and `kit: kit` assignment (base mutation blob `ff93a5882e1a7301e4708b67526748f5c4266f89`); the semantic audit exited 1 with `evaluation.go:114 function=newProspectiveTreeEvaluationAtKit kit=evaluation.go:117` and `stored kit in 1 of 1`, while `GOMAXPROCS=2 go test -p 1 -parallel 1 -run '^$' ./internal/gate` exited 0 with no tests to run. After exact restoration to blob `70fec552…`, mutation 2 added `kit: identityRoot` directly to the surviving constructor; the audit exited 1 with `evaluation.go:110 function=newProspectiveTreeEvaluation kit=evaluation.go:113` and `stored kit in 1 of 1`, while the same compile-only command exited 0 in 0.002 s. | Restored blob `70fec55267a33ec13882560d77d74c51e7508a0e` with `apply_patch`; the audit exited 0 with `function=newProspectiveTreeEvaluation kit=absent` and `green prospective_literals=1 stored_kit=0 production_files=25`; the focused command exited 0 in 0.323 s. |

`repair-ambient-red-contract` supplied corrected metadata after debug; it made
no candidate behavior claim. The promotion-recomposition repair was a separate
main-lifecycle prerequisite, not a gate-test candidate mutation.

## RPE structural probes

| criterion | green baseline | mutation red | restoration |
|---|---|---|---|
| RPE1/overlap-serialization | AST audit exit 0 at 245/192/53; focused overlap test exit 0, package 0.255 s, wall 1.52 s. | Re-added `t.Parallel` to `TestRunnerRunsPhasesConcurrently`; audit exited 1: `serial test has 1 parallel calls`. | Runner blob restored to `e64197ac…`; audit exit 0 at 245/192/53. |
| RPE1/census | AST audit exit 0 at 245/192/53. | Removed the first-statement marker from `TestMergeEnvStripsThenSets`; audit exited 1: `eligible test calls=0 first=false`. | Runner blob restored to `e64197ac…`; audit exit 0 at 245/192/53. |
| RPE1/serial-inventory | Exact list comparison exit 0. | Omitted `TestRunnerRunsPhasesConcurrently` from this inventory; comparison exited 1 naming the missing serial test. | Restored the row exactly; comparison exit 0. |
| RPE1/timing | Exact-candidate median is 56.72 s. | The retained pre-adoption subject is the structural marker-removal mutation: three same-command runs produced a 150.85 s median, above 90 s. No duplicate high-I/O rerun was needed. | Repaired exact candidate retained all eligible markers and the 56.72 s median. |
| RPE1/mutation-ledger | Ledger-shape audit exit 0 with baseline, mutation, red exit/symptom, restoration, and restored exit for every required fact. | Removed the restored exit from the TP row; audit exited 1 naming `TP1/serial-list` as incomplete. | Restored the field exactly; ledger-shape audit exit 0. |
| RPE1/narrow-timeout | Width-one baseline exits 0 at 138.01 s under 600 s. | Restored `-timeout 120s` in the narrow-width row; validation exited 1 because the timeout is below the observed unmutated baseline. The test was not misreported as a deadlock red. | Restored `-timeout 600s`; validation exit 0. |
