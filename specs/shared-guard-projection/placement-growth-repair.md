# Guard conformance placement repair

Accepted source: 02c799aa8e7ba9f1dfcf38a14906f51ab4a6d8db
Shared spec SHA256: 1e33af8212b42d36fc311fe89e3075d626b78ec81085b4937826a56c5993fc4b
Scan spec SHA256: 0e99ad3f2f87f457ecb3b3f64fabd474381b4621293089befe112aa54fafb2cd
Prior confirmation: adoption_spec_review accepted both specs on all three axes, confidence 10
Placement author: gpt-6.1-sol/high, one bounded amendment pass
Placement confirmation: /root/deepen_execution, gpt-6.1-sol/xhigh, accepted 6c13b4b31eb4bd00b4a888e060f1c48a82f572c8 on all three axes, zero findings, confidence 9
Ticket graph: independently accepted at aa86c64d0e6227b4ba6d094518bbd295f329e0af by /root/deepen_tests_parsers, declared gpt-6.1-sol/xhigh, all three axes zero findings, confidence 9

## Learning and required lane

Check landing headroom through the actual lane and its granted file limits before choosing a source placement.
Directory crowding alone did not establish this landing blocker.
The verified blocker is growing an already oversized, ungranted orchestration file.

internal/structure/structure.go:59–156 evaluates full-tree file sizes and directory counts.
The default limits are 400 lines per file and 12 source files per directory.
.bench/structure-accept:8 grants the existing internal/conformance/ directory as one cohesive check family.
Existing crowded cmd/bench and systemtest directories remain inherited structural debt.
This amendment neither grants that debt nor changes the directory policy.

internal/structure/structure.go:283–343 owns Growth.
It reports FILE GREW only when a source file exceeds its effective limit and its base line count.
internal/gate/lane.go:82 selects structure --growth for the required commit lane.
The ordinary gate entry does not add a full-tree directory check.
No runtime gate was executed to establish these source facts.

## Exact placement and headroom

internal/conformance/package_core_checks_test.go has 470 lines at the accepted source.
Its default cap is 400; neither structure.budgets nor structure-accept grants that file.
A new subcheck call would increase its count beyond both the limit and the accepted base.
The required growth predicate would therefore report FILE GREW.

internal/conformance/checks_test.go has 641 lines.
.bench/structure.budgets:49 grants it 709 lines, leaving 68 lines of current headroom.
Move only checkPackageCoreAndGuards, whose declaration occupies lines 22–37, into that existing executable binding and runner owner.
The 16-line declaration leaves 52 lines before added separators and ownership subcheck calls.
The destination must remain within 709 lines after both outcomes, and each new source file must remain within 400 lines.

The original file shrinks by the declaration and any removed separator; its remaining inherited size is permitted by Growth.
Do not move package-core producers, fixture helpers, npm probes, or tests.
Do not create a new orchestration file or change a budget grant.
Both accepted prospective fences already contain the source and destination.
Acceptance rows, test witnesses, caller policy, and all other fences remain unchanged.

Shared GP3 and scan GS1 apply the move only if the declaration still resides in package_core_checks_test.go.
If the other outcome already landed the move, extend the sole declaration in checks_test.go.
Keep existing subcheck ordering and any already landed ownership subcheck.
Neither outcome depends on the other's unlanded implementation.
This is prospective placement guidance; no source was relocated during specification.

## Declaration, callers, and source readers

internal/conformance/checks_test.go:49 binds package-core-guard to the function by symbol.
internal/conformance/registry/checks.go:19 advertises the same implementation name, Dev tier, SubjectRoot, and InputCatchAll.
checkBinding.identity in check_bindings_test.go:21–28 derives the runtime function name, not its declaration file.
The move keeps that name, signature, registry row, and executable binding unchanged.
node_runtime_policy_test.go:143 and :147 call the same package symbol and need no source edit.

package_core_diagnostics_test.go:202–228 reads package_core_checks_test.go by path to count npm formatter call sites.
The npm producer and its one formatProbeFailure call remain in that file.
The moved orchestration contains no formatter call, so the reader's expected count and subject remain unchanged.
Do not rewrite this expectation merely because the orchestrator moves.

TestSkillsIndexConformanceCarriesNoSecondReader reads checks_test.go and skills_index_checks_test.go by path.
It forbids copied skills-index marker, allowlist path, and line-format literals.
The moved declaration has none of those literals; extend it with function calls rather than copied policy data.
The live-tree classification inventory parses the package and keys by test name, so this non-test declaration changes no classification.
The planned ownership tests still require their accepted live-tree rows.

The hidden symbol and path census also finds planning citations in degraded-guard-refusal and shared-test-fixtures.
Those citations describe their own pinned source; they are not executable declaration-location contracts.
No fixture or source reader requires this orchestration declaration to stay in package_core_checks_test.go.

## Transitive closure preserved

The canonical-path-owner/second-derivation canary targets internal/gitguard/gitguard.go's JSON import and Checker declaration.
The injected-ports/unregistered-port canary targets that same Checker declaration and includes checker_junction_test.go in its base.
Neither targets the conformance orchestration location.
Shared projection keeps both canary families in its existing fence because its Git migration changes their actual mutation subject.
Preserve their named red outcomes while rebasing only mutation needles changed by that implementation.

The three GitHook anchors remain attached to .bench/hooks/block-dangerous-git.sh in internal/anchors/registry_data.go:403–405.
They pin the honest-mistake boundary, one-level scan depth, and named external backstops.
This placement moves no anchored prose or hook path.
No anchor relocation or anchor expectation rewrite is required.

The five command registry and inventory closure paths remain in both accepted prospective fences.
Their source checks retain package-core-guard's existing name, tier, subject, and input source.
No new conformance check count, profile row, injected port, or parallel scanner is introduced.
No production, test, registry, fixture, roadmap, or metadata source changed in this amendment.

## Confirmed placement and graph use

The xhigh confirmer accepted the complete placement amendment before ticket slicing.
It rechecked actual lane/file limits, callers, source readers, canary families, hook anchors, command closure, and all 89 predicate bytes.
The complete 6c13b4b source remains retained; backend execution identity and usage were not observed.

GP3 and GS1 own the relocation at their first registered subcheck use, with either-plan-first composition.
The graph keeps the existing npm producer and its independent formatter-call expectation in the original file.
Canonical path fences and exact first-use ownership retain both source and destination, without a grant or unlanded cross-plan dependency.
The graph review records state the exact row partition and mechanical fence-commentary normalization.
No production declaration moved during this planning pass, and all runtime witnesses remain future obligations.
