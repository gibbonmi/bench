# Conformance observation follow-on

Status: staged

Decision source: Named reviewed artifact `decisions/architecture-observations.md`, resolved tickets 1–3, 2026-10-06.

Verification log: 3 iterations to accept — two independent spec passes and one independent ticket pass.
Retained xhigh reviewer `deepen_execution` accepted the full spec after COF-R1 repair, before ticket slicing.
The same independent reviewer accepted the graph at `c04ccd1bbac2274fb8197d592ea7aef21beffd94` on all three axes, with 0 findings and confidence 9.
Reviewed spec SHA-256: `0c3a7dd1616cb5b1caef7a62e51187da0074b3e242fcfe9475ef6e327bec8684`.
This metadata close preserves all 38 rows, fences, and ticket bytes and claims no runtime or performance proof.

Review evidence: [spec-review.md](assets/spec-review.md).

## Problem

Several conformance checks walk and read overlapping Go source independently.
Sharing a helper does not prove that registered checks share their observations.
A stale observation can also hide a newly planted violation.

The approved landing-test-efficiency spec already owns the source observation module and its first four visitors.
This follow-on migrates six remaining visitors without changing their policies.
It does not reopen the LTE design or claim measured speed savings.

## Solution

Each selected visitor consumes the observation owner from the complete, accepted, and landed LTE outcome.
One dispatcher invocation owns its observations.
Visitors retain their own scope, exclusions, predicates, diagnostics, ordering, and error posture.
Later invocations and direct check calls observe fresh source.

Operation counts prove reuse through the actual dispatcher and registered bindings.
Independent policy checks remain independent from the observed source.
Every existing assertion keeps its purpose and expected behavior.

## User stories

Line: gpt-6.1-sol / high.

Implementation-line reason: Six legacy policies share one future owner, with different traversal and refusal rules.
Existing canaries and real dispatch tests provide strong witnesses.
The strict VCS reader makes COF-C3 the hardest chunk.
This recommendation does not change LTE's approved gpt-5.6-sol / high line.

Harder chunks: COF-C3.

1. As a maintainer, I want the complete LTE outcome first, so that this migration uses one accepted owner.
2. As a maintainer, I want canonical-path checks to keep their scope, so that excluded source stays excluded.
3. As a maintainer, I want canonical-path diagnostics preserved, so that a second derivation still turns red.
4. As a maintainer, I want cancellation checks to keep their scope, so that owner and test files stay excluded.
5. As a maintainer, I want cancellation registration rules preserved, so that alternate signal sets still turn red.
6. As a maintainer, I want cancellation prefilters preserved, so that irrelevant malformed source keeps its existing disposition.
7. As a maintainer, I want wait checks to keep their test-only scope, so that production source receives no new policy.
8. As a maintainer, I want deadline classifications preserved, so that polling and derived windows remain valid.
9. As a maintainer, I want default-branch checks to retain their byte predicate, so that migration adds no AST policy.
10. As a maintainer, I want control-escaping checks to retain their package count, so that ownership drift still turns red.
11. As a maintainer, I want VCS checks to keep their source scope, so that hidden and test source remain graded.
12. As a maintainer, I want VCS import resolution preserved, so that aliases and dot imports keep their current meaning.
13. As a maintainer, I want VCS refusals preserved, so that special files cannot become blocking reads.
14. As a maintainer, I want VCS failure order preserved, so that the same first failure wins.
15. As a maintainer, I want same-subject visitors to share observations, so that repeated work is removed through real dispatch.
16. As a maintainer, I want distinct subjects isolated, so that root violations cannot become kit diagnostics.
17. As a maintainer, I want later runs to observe mutation and restoration, so that cached green cannot hide drift.
18. As a maintainer, I want direct calls to observe mutation and restoration, so that fixture proofs remain effective.
19. As a maintainer, I want visitor order and membership preserved, so that dispatch never drops a check.
20. As a maintainer, I want timing state to reset, so that a second run cannot inherit the first run's rows.
21. As a maintainer, I want reader failures to remain visitor policy, so that sharing facts does not change refusals.
22. As a maintainer, I want old walkers detected, so that helper reuse cannot conceal duplicate source work.
23. As a maintainer, I want immutable source observations, so that one visitor cannot alter another visitor's input.
24. As a maintainer, I want hostile paths preserved, so that a refactor keeps the existing path and diagnostic contracts.
25. As a maintainer, I want retained canaries to bite through registered owners, so that helper-only proofs cannot hide an omitted consumer.
26. As a maintainer, I want metadata walks excluded, so that this outcome does not become a universal analyzer.
27. As a maintainer, I want registry and anchor closures checked, so that the planning fence covers actual enforcement consumers.
28. As a maintainer, I want bounded modules, so that a migration cannot spend a structure grant without review.

## Implementation decisions

### Closed dependency and owner contract

LTE means `specs/landing-test-efficiency/spec.md`, including LTE-C1 through LTE-C8 and final acceptance.
Implementation starts only after that entire outcome lands and its independent reviews accept it.
A partial LTE chunk is not an admissible dependency.
This spec adds no source owner and no dispatcher engine.

The accepted LTE contract is the interface authority.
Its future exported names are not present at the research pin.
Planning and slicing may proceed against LTE's approved value contracts.
After complete accepted LTE lands, resolve concrete API names and recheck landed-tree headroom before any implementation charge.
That name binding may clarify the plan but cannot weaken these value contracts.

| Boundary | Required value contract | Existing or accepted authority |
| --- | --- | --- |
| RunConformanceSelection to snapshot owner | One invocation owns lazy observations indexed by selected subject identity | LTE-C4 through LTE-C6 |
| checkBinding to visitor | Named implementation identity, tier, and subject remain independently recorded | check_bindings_test.go:15–46 and tier_test.go:103–148 |
| Snapshot to visitor | Ordered path facts, bytes or reader errors, and matching AST/token positions are immutable | LTE sourcefiles contract |
| Direct check to owner | Each legacy root-taking entry creates a fresh invocation adapter | LTE direct-adapter contract |
| VCS visitor to owner | No-follow classification and readable source facts preserve strict refusal order | build_contracts_test.go:138–211,293–299 |
| Visitor to dispatcher | Existing diagnostics pass through unchanged in existing order | checks_test.go:96–144 |

A visitor can request bytes without requesting an AST.
Parsing is lazy and occurs only when a selected visitor's current predicate requires it.
The shared AST must preserve the current parser mode and diagnostic token positions.
A visitor never edits shared bytes, nodes, or token data.
It uses a private policy result, not a shared verdict.

The VCS seam requires bounded no-follow producer observations.
The capability returns canonical bounds.Classified facts: State, Stream, Data, and Reason for the selected path.
It refuses links and other non-regular objects before requesting bytes.
It retains ControlRecordLimit and bounds' existing UTF-8, empty, oversized, and failed-read classification.

COF-C3 supplies this capability in the accepted sourcefiles owner if complete LTE does not already expose it.
The authorized bridge is a narrow reader injection port in internal/bounds/classify.go, with the existing ClassifyNoFollow adapter retained.
The existing bounds classifier keeps no-follow and classification policy.
The observation owner supplies invocation-cached regular-file reads through that port.

The port is passed explicitly per call and cannot use a package-global swap.
Its lifetime ends with the invoking snapshot or fresh direct adapter.
Do not add a private classifier, walker, or parallel cache.

The strict observation can reuse already complete regular-file bytes through bounds' existing bounded byte classifier.
A strict-first read stays bounded and cannot claim truncated bytes as a complete legacy observation.
A later unbounded legacy request may require a separate physical read for that path.
Cache identity therefore includes required read posture and completeness, while subject identity remains LTE-owned.
No additional read is allowed for overlapping complete regular source within the producer bound.

For the oversized count witness, select actual go-build-vcs and wait-deadline-literals bindings over internal/sample/large_test.go.
The file has a valid package, padding beyond ControlRecordLimit, and a literal NewTimer deadline after that padding.
Registry order requests the bounded VCS observation before the unbounded test-source observation.
Expect two physical reads: one bounded producer read, then one complete legacy read.
Expect one AST parse for the wait visitor and the existing VCS unavailable plus wait-literal diagnostics.
A truncated legacy view or a raised producer limit cannot satisfy both diagnostics and operation counts.

After complete accepted LTE lands, bind the port's concrete names against that tree before any implementation charge.
An existing compatible capability removes the bridge work rather than creating another API.
A source-staleness amendment records a different landed interface without weakening the accepted capability or preservation witnesses.

### Six policy owners

The current definitions below were read in full at the research pin.
Each root-taking function remains a fresh direct adapter with the same public-to-package behavior.
Its snapshot visitor receives the already selected subject and applies its own policy.

| Registered name and current owner | Eligibility and order | Policy and reader behavior retained |
| --- | --- | --- |
| canonical-path-owner: checkCanonicalPathOwner at canonical_path_owner_test.go:23 | cmd/internal non-test Go, excluding internal/canonicalpath/, uniqueSorted | Same FuncDecl containing literal filepath.Abs and filepath.EvalSymlinks produces the existing function diagnostic |
| cancel-signal-registrations: checkCancelSignalRegistrations at cancel_signal_registrations_test.go:29 | cmd/internal non-test Go, excluding internal/subprocess/, uniqueSorted | signal.Notify substring prefilter precedes parsing, exact CancelSignals spread rule retains all current diagnostics |
| wait-deadline-literals: checkWaitDeadlineLiterals at wait_deadline_literal_test.go:31 | cmd/internal test Go only, uniqueSorted | Current mentionsWaitDeadline prefilter and AST classification retain polling, backdated time, and TestDeadline exemptions |
| default-branch-single-source: checkDefaultBranchSingleSource at git_facts_checks_test.go:25 | Whole tree non-test Go, skip .git/node_modules/tests/dist directories, sorted diagnostics | Existing DefaultBranch call/declaration regexp consumes bytes, including comment/string matches |
| single-control-escaper: checkSingleControlEscaper at data_handling_test.go:114 | Whole tree non-test Go, skip .git/node_modules/tests/dist and exact internal/toon directory, sorted unique package set | Existing hexadecimal control-escape regexp grades zero/one/many packages without requiring a sanitize package name |
| go-build-vcs: checkGoBuildVCS at build_contracts_test.go:138 | Whole tree Go including tests, skip nested .git/.logs/dist/node_modules/vendor directories, lexical walk order | readBuildContractSource refuses non-parsed/non-empty classifications, parse errors continue, read/traversal errors abort |

The first five visitors currently ignore traversal errors.
Their existing readIfExists behavior treats failed reads as empty source.
The VCS visitor instead publishes its existing scan-unavailable diagnostic.
Shared observation errors do not impose one global fail posture.
Unselected and excluded source cannot introduce new diagnostics or eager failures.

Canonical and cancellation matching uses the current literal selector spellings.
Neither gains import-alias resolution.
VCS keeps importedPackage alias and dot-import handling for os/exec and the genuine gate constant.
Wait classification keeps its current TestDeadline spelling rule, including the existing selector-name treatment.

### Current reader and writer census

Research subject: direct `internal/conformance/*.go` files at `a9c395e77fec36d1f60f6d057c654aecf5720e24`.
The lexical query selected files containing parser.ParseFile, filepath.WalkDir, or filepath.Walk followed by an opening parenthesis.
It found 26 files.
This count describes that exact lexical query, not all scanners or all repository parser calls.

Exact lexical query, scoped to the source pin above:

```python
from pathlib import Path
needles = ("parser.ParseFile", "filepath.WalkDir", "filepath.Walk(")
files = [p for p in sorted(Path("internal/conformance").glob("*.go"))
         if any(needle in p.read_text() for needle in needles)]
```


The import and call census also checked go/parser imports, ParseDir, visitor helpers, and executable bindings.
These imports use the unaliased parser name in the inspected family.
ParseDir adds tier_live_tree_test.go to the named exclusions below.
Historical FT365 says 31 files.
Its historical count is context, not a current census or a savings estimate.

| Disposition | Current files and actual owners | Reason |
| --- | --- | --- |
| Complete LTE dependency | ordinary_build_census_test.go, git_plumbing_owner_test.go, bounds_policy_test.go, skip_ownership_test.go | Four accepted visitors and their walker removal already belong to LTE |
| LTE helper closure | bounds_waits_policy_test.go:96 waitPolicy.packageBindings | Bounds visitor package-binding facts belong to LTE's bounds migration, not a fifth new policy |
| This follow-on | canonical_path_owner_test.go, cancel_signal_registrations_test.go, wait_deadline_literal_test.go, git_facts_checks_test.go, data_handling_test.go, build_contracts_test.go | Six reusable source visitors with preserved local policy |
| Targeted registry parsers | axi_query_registry_test.go:279 parseAXIRegistry, command_registry_parse_test.go:35 parseRegistryTable, subcommand_routing_test.go:166 packageReachesGrammar | Named registry syntax and bounded package grammar checks retain separate semantics |
| Targeted policy evidence | guard_classifier_table_test.go:175 parseGuardClassifierRows, handoff_single_source_test.go:218 checkHarnessPrefix, injected_ports_test.go:132 derivedInjectedPorts | Explicit tables and bounded packages, not overlapping whole-source visitors |
| Targeted package evidence | cross_compile_default_test.go:75 packageTestFuncs, otel_seam_test.go:59 seamFunctionBody, skills_index_checks_test.go:373 TestSkillsIndexConformanceCarriesNoSecondReader | Named package/symbol and second-reader proofs retain their present callers |
| Fixture and dispatch metadata | fixture_bite_test.go:56 TestFixtureBiteProofArchitecture, tier_test.go:232 hiddenLiveTreeDiags | Enforcement reads its own proof architecture and classification |
| Non-Go or mixed asset walks | docs_workflow_checks_test.go:758 walkConformanceDocs, workflow_checks_test.go:117 checkRetiredReproducibilityRecord, package_shipped_surface_test.go:133 shippedFiles, help_inventory_single_source_test.go:50 checkHelpInventorySingleSource | Markdown, packaging, or mixed file policy requires a distinct capability |
| ParseDir-only metadata exclusion | tier_live_tree_test.go:74 TestClassifiedLiveTreeInventoryNamesDetectedTests | Package declaration inventory is outside the 26-file lexical subject |

The build-contract file also owns publishedExecutable.
Its targeted script/build-helper proof stays unchanged outside the VCS visitor.
Likewise, DATA_HANDLING and project pass-list checks in data_handling_test.go stay unchanged.
The census does not transfer their authority.

A resolved consumer query covered the six migrated entries, four migrated helpers, and three injected-port evidence symbols.
It returned 34 static edges across eight files from 327 loaded Go packages.
Its complete answer hash is `b2627a4850c274131aa4baa76cc13629820be782edabf526abbb1f9afc94e3d7`.

Exact resolved-reference query:

```text
bench consumers conformance.checkCanonicalPathOwner conformance.canonicalPathDerivations conformance.checkCancelSignalRegistrations conformance.cancelSignalRegistrationDiags conformance.checkWaitDeadlineLiterals conformance.waitDeadlineLiteralDiags conformance.checkDefaultBranchSingleSource conformance.checkSingleControlEscaper conformance.controlEscaperPackages conformance.checkGoBuildVCS conformance.checkInjectedPortRegistry conformance.derivedInjectedPorts conformance.fileDeclaresTest --full
```

The query is a static reference census, with Bench consumers' documented reflection, build-tag, and dynamic-execution limits.

### Registered consumer and closure contract

All six current bindings select the root subject and retain their existing names.
Registry order comes from internal/conformance/registry/checks.go:22–31.
Selection, scope, ship partition, and timing remain the LTE dispatcher contract.

Retained fixture proofs reach checkBindings through runFixtureBite at fixture_bite_test.go:750–779.
Canonical and cancellation canaries retain their exact BASE, MUTATE, and EXPECT artifacts.
CM5, CM6, and CM7 retain independent function, tier, and subject expectations.
New observation tests must enter the live-tree classification where applicable.
No new executable check name or fixture family is introduced.

Anchor consumers remain read-only because no guidance, command, or anchored policy sentence changes.
The command registry closure is main.go, command_registry.go, command_registry_test.go, help_inventory_test.go, axi_query_registry_test.go, and subcommand_routing_table_test.go.
Its existing ticket-registry co-name requirement remains satisfied by keeping those paths unchanged.

The injected-port registry audits five named production packages, excluding both conformance and bounds.
This outcome does not add either package to its audited set.
The bounds read port closes through its local classification test and the actual registered VCS consumer.

Its classifier owner also inherits two existing fixture units through their BASE lists.
FixturePins at internal/canary/inventory.go:149–246 derives that closure from BASE, overlays, and mutation paths.
Preflight enforces the owning ticket's co-name requirement at internal/preflight/closure.go:54–73,118–123.
No port-registry or command-registry row changes.

## Implementation chunks

Independent full-spec acceptance preceded these three serial vertical tickets.
Each chunk delivers registered behavior and its real dispatch proof.
Ticket-graph acceptance remains pending.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| COF-C1 / 01-share-canonical-and-cancellation-observations.md | Canonical and cancellation visitors share complete LTE observations through real bindings | COF01–COF06, COF15–COF18, COF22–COF23, COF29 | Existing canary proof and planned TestResidualObservationDispatch | no |
| COF-C2 / 02-share-deadline-and-byte-observations.md | Deadline and byte visitors preserve their policies through shared dispatch | COF07–COF10, COF19–COF21, COF25–COF26, COF37 | Existing visitor tests and planned TestResidualObservationPolicyMatrix | no |
| COF-C3 / 03-share-vcs-observations-and-close-all-visitors.md | VCS strict reader consumes shared facts and proves the aggregate ten-visitor contract | COF11–COF14, COF24, COF27–COF28, COF30–COF36, COF38 | Existing VCS special-file tests and planned TestResidualObservationDispatch | yes |

COF-C2 depends on COF-C1.
COF-C3 depends on COF-C2 and binds the owner capability against the complete landed LTE source.
Every chunk depends on the entire landed LTE outcome.
There is no owner-only or helper-only checkpoint.

The classifier-owning COF-C3 slice must co-name tests/canary/package-core-guard/bounds-classify-limit-restated and tests/canary/package-core-guard/bounds-read-limit-restated when it writes internal/bounds/classify.go.
These are required existing fixture units, not new fixture work.
Their BASE, MUTATE.json, EXPECT, CHECK, meanings, and existing assertions remain unchanged.
The classify fixture retains the learnings ControlRecordLimit omission.
The read fixture retains the models ModelReadLimit omission.

### Stable chunk mapping and checkpoint contracts

COF24, COF28, COF30, and COF31 move from COF-C1 to COF-C3.
COF27 moves from COF-C2 to COF-C3.
All other row assignments and stable chunk IDs remain unchanged.
The terminal slice owns complete-family reconciliation, not the first valid checkpoint for an earlier consumer.
Each ticket applies freshness, hostile-input preservation, owner consumption, headroom, and valid evidence requirements to its own introduced behavior.
No earlier ticket borrows a later repair.

Planning and slicing use LTE's approved value contracts.
After complete accepted LTE lands, concrete API binding and the landed-tree headroom recheck must finish before any implementation charge.

## Testing decisions

Existing visitor assertions remain intact.
New tests use small fixture trees through RunConformanceSelection and the actual registered checkBindings.
Count proofs observe owner operations, not elapsed time or a standalone helper invocation.

The shared fixture contains root source, cmd/sample production source, and internal/sample test source.
It includes source selected by several AST visitors and both byte visitors.
For each complete regular fixture path within ControlRecordLimit, the count oracle expects one physical byte read across requesting visitors.
For each actual parsed path and parser mode, it expects one parse across requesting visitors.
Each requested directory observation occurs once per subject snapshot.

Expected path sets are authored from the fixture contract, not the owner's resulting counters.
The test names the selected binding set independently and checks exact execution membership.

The full proof selects LTE's four visitors and this outcome's six visitors.
A per-binding snapshot mutant increments repeated directory/read/parse counts and fails.
A restored private walker fails a structural owner check even if counters cannot see its private I/O.
The check examines each migrated entry and its helper closure for direct walking, reading, parsing, or hidden replacement helpers.
Named metadata exclusions remain outside that rule.

For refusal order, select only go-build-vcs over a fixture with FIFO a.go and a later z.go missing its VCS flag.
The sole winning diagnostic is `Go build VCS scan unavailable: <root>/a.go: wrong-type: not a regular file: p---------`.
The root placeholder is the explicit fixture root, not a value obtained from the new diagnostic producer.
A whole-inventory collector that reports z.go first fails this exact comparison.
A separate malformed a.go case keeps its parse diagnostic and the later z.go flag diagnostic in lexical order.

A second invocation follows a planted source mutation, then a restored third invocation.
Diagnostics must show clean, violation, clean with fresh operation counts each time.
Direct adapters perform the same sequence independently.
Distinct root and kit fixtures contain different violations and cannot exchange observations or diagnostics.
Identical selected subjects may share only within the current invocation under LTE's identity contract.

### Seam diagram

```text
RunConformanceSelection(root, kit, tier, selection)
    -> actual registered bindings and selected subject
    -> invocation-owned LTE observations
    -> six named policy visitors
    -> unchanged diagnostics and timing rows
Tests observe binding execution, diagnostic bytes, and owner operation counts.
Direct root-taking checks create a fresh adapter for each call.
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| COF01 | 1 | Build admission requires complete accepted and landed LTE | review-owned dependency check against LTE final source and review records | A partial chunk cannot supply an admitted dependency |
| COF02 | 2 | Canonical eligibility retains cmd/internal non-test and owner exclusions | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Violations in each included and excluded partition expose broadened filters |
| COF03 | 3 | Same-function Abs plus EvalSymlinks keeps its exact function diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The canonical canary must bite through its real binding |
| COF04 | 4 | Cancellation eligibility retains owner and test exclusions | `internal/conformance/cancel_signal_registrations_test.go` (`TestCancelSignalRegistrationsBites`) | Included violation and excluded test source prevent silent scope changes |
| COF05 | 5 | Notify and NotifyContext retain spread, argument, and diagnostic rules | `internal/conformance/cancel_signal_registrations_test.go` (`TestCancelSignalRegistrationsBites`) | Valid spread and invalid sets must remain distinguishable |
| COF06 | 6 | Cancellation prefilter preserves malformed source without signal.Notify | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Eager AST failure would add a diagnostic to the silent partition |
| COF07 | 7 | Wait visitor grades only cmd/internal test Go | planned internal/conformance/wait_deadline_literal_test.go (`TestWaitDeadlineObservationScope`) | Matching production literals stay silent while matching test literals bite |
| COF08 | 8 | Wait AST rules retain literal deadlines and valid derived or polling windows | `internal/conformance/wait_deadline_literal_test.go` (`TestWaitDeadlineLiteralsBites`) | Existing red and green cases prevent a new projection from rewriting purpose |
| COF09 | 9 | DefaultBranch byte matching retains all directory and file partitions | `internal/conformance/git_facts_checks_test.go` (`TestDefaultBranchSweepScopesToTheFunction`) and planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Comment/string byte matches and excluded directories expose AST or filter substitution |
| COF10 | 10 | Control escaping retains sorted zero, one, and multiple package dispositions | `internal/conformance/data_handling_test.go` (`TestSingleControlEscaperBites`) | Empty and duplicate owners still produce the original diagnostics |
| COF11 | 11 | VCS scan retains root, hidden, tagged, and test source eligibility | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | A missing flag in every included partition exposes shared-filter leakage |
| COF12 | 12 | VCS command family retains aliases, dot imports, context calls, and genuine gate constants | `internal/conformance/build_contracts_test.go` (`TestGoBuildVCSCallFamily`) | Wrong constants and conflicting flags remain red beside valid calls |
| COF13 | 13 | VCS refuses special files through the no-follow reader | `internal/conformance/build_contracts_test.go` (`TestBuildContractChecksRefuseSpecialFiles`) and planned internal/conformance/build_contracts_test.go (`TestVCSObservationRefusals`) | FIFO and symlink paths cannot become ordinary blocking or followed reads |
| COF14 | 14 | VCS's first unavailable path wins before a later source violation | planned internal/conformance/build_contracts_test.go (`TestVCSObservationRefusalOrder`) | A fixture containing both failures rejects eager whole-inventory error substitution |
| COF15 | 15 | Actual selected canonical and cancellation bindings share one subject snapshot | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | Executing only an owner helper cannot satisfy actual binding membership |
| COF16 | 15 | Overlapping complete regular source within the producer bound produces one byte read per path | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | Per-binding snapshots repeat reads in the independent fixture path count |
| COF17 | 15 | Overlapping AST requests produce one parse per path and mode | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | A parser hidden in a helper increments counts or fails the old-reader guard |
| COF18 | 15 | Shared traversal produces one observation per requested directory | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | Restoring visitor-local traversal violates the count or structural guard |
| COF19 | 19 | Selection and scope retain exact registered execution order and membership | `internal/conformance/tier_test.go` (`TestScopedRunExecutesOnlyTheNamedCheck`, `TestOrderedSetRunsMetaAndSelectedOrdinaryChecksInRegistryOrder`) | An omitted or extra binding differs from the independent expected sequence |
| COF20 | 19 | Dev and ship partitions retain their current binding dispositions | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`) and LTE dispatch regression tests | Tier remapping and missing implementation mutations remain red |
| COF21 | 20 | Timing cardinality and per-run reset retain the complete LTE contract | `internal/conformance/tier_test.go` (`TestTimingOrderStable`) | The second full run cannot inherit timing rows or lose a selected check |
| COF22 | 16 | Distinct root and kit subjects retain separate diagnostics and observations | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | Two fixtures with different violations expose a snapshot keyed only by first root |
| COF23 | 17 | Later dispatch observes mutation and restoration as clean, red, clean | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationFreshness`) | A retained snapshot hides the middle violation or retains the restored violation |
| COF24 | 18,25 | Direct adapters retain mutation and restoration freshness | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) and planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationFreshness`) | Repeated calls through each adapter cannot reuse a stale AST or byte buffer |
| COF25 | 21 | First five visitors retain silent traversal and empty-read failure behavior | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Shared failure policy would publish a new refusal in the legacy silent partition |
| COF26 | 21 | Parse diagnostics preserve prefilter, path, and line behavior | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Malformed matching source distinguishes cached parse facts from changed formatting |
| COF27 | 24 | Paths with spaces, control bytes, quotes, glob characters, and non-ASCII retain relative diagnostics | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Wrong root resolution or path normalization changes the exact diagnostic subject |
| COF28 | 22 | No migrated entry or helper retains an independent walker, reader, or parser | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationBindingsUseSharedOwner`) | Restoring the old canonical helper or hiding its parser behind an alias must turn red |
| COF29 | 23 | Visitor processing leaves shared observations unchanged | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | Reversing visitor order over the same input exposes mutable shared nodes or bytes |
| COF30 | 26,27 | Only the six named visitor families migrate with their enforcement closures intact | review-owned census and `internal/conformance/tier_live_tree_test.go` (`TestClassifiedLiveTreeInventoryNamesDetectedTests`) | A hidden live-tree test or unrelated parser migration violates the named boundary |
| COF31 | 28 | Changes fit landed-tree structure limits without new grants | review-owned structure/headroom check before ticket acceptance | An extra conformance file or oversized proof cannot silently consume a grant |
| COF32 | 13 | VCS accepts empty regular source and refuses oversized or invalid text with current classifier states | planned internal/conformance/build_contracts_test.go (`TestVCSObservationRefusals`) | A bytes-only owner would erase the existing classification distinction |
| COF33 | 14 | VCS parse failure continues to later files in lexical order | planned internal/conformance/build_contracts_test.go (`TestVCSObservationRefusalOrder`) | Aborting on an AST error loses the independently expected later flag diagnostic |
| COF34 | 15,19 | All ten actual LTE and follow-on visitors satisfy aggregate operation and membership counts | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationDispatch`) | A restored old walker or one private snapshot per binding fails the complete selection |
| COF35 | 27 | Canonical and cancellation canary artifacts retain exact mutation and expectation bytes | `internal/conformance/fixture_bite_test.go` (`TestFixtureBiteProofArchitecture`, `TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Replacing a canary expectation with owner output cannot preserve the omission proof |
| COF36 | 1,13 | VCS consumes canonical no-follow observations through the landed owner or its authorized reader bridge | planned internal/bounds/classify_bytes_test.go (`TestNoFollowObservationReaderPreservesClassification`) and planned internal/conformance/build_contracts_test.go (`TestVCSObservationRefusals`) | A missing owner capability, duplicate classifier, or unbounded strict read fails its preserved state and entry witnesses |
| COF38 | 13,21 | Strict-first oversized source retains two-read policy counts and the later legacy violation | planned internal/conformance/canonical_path_owner_test.go (`TestResidualObservationPolicyMatrix`) | Reusing truncated bytes would hide a legacy byte violation beyond the strict bound |
| COF37 | 19 | Invalid ordered selection retains its diagnostic and full-tier widening | `internal/conformance/tier_test.go` (`TestInvalidOrderedSetRedsAndWidensToTheFullTier`, `TestShipTierIgnoresOuterOrderedSelection`) | Replacing legacy widening with refusal changes the observed execution set |

### Edge inventory

| Class | Disposition at this seam |
| --- | --- |
| Missing or empty source | Preserve each visitor's current read and prefilter disposition, including VCS empty acceptance |
| Malformed Go and absent final newline | Preserve exact parser behavior and visitor-specific prefilter or continuation policy |
| Symlink, dangling symlink, FIFO, other special file | Preserve VCS no-follow refusal before reading, with separate legacy dispositions for the other visitors |
| Oversized or invalid UTF-8 source | Preserve VCS bounds classification, without applying its strict policy to legacy byte readers |
| Directory traversal errors | First five visitors remain silent, VCS retains scan-unavailable and first failure order |
| Comment, string, alias, dot import | Byte visitors retain byte matching, literal-selector visitors retain literal matching, VCS retains import resolution |
| Deep root, spaces, glob, quotes, slash, control bytes, non-ASCII | Observe the explicit selected subject, preserve relative paths and exact diagnostic bytes |
| Mutation, restoration, repeated invocation | Create fresh observations for each direct call and dispatcher invocation |
| Reversed visitor processing | Preserve immutable facts and each visitor's own sorted or lexical result order |
| Unknown scope or scope combined with ordered selection | Preserve dispatcher refusal and empty timing state before visitor execution |
| Invalid ordinary ordered selection | Preserve the existing red diagnostic and widening to the full tier |

**Won't handle:** Universal type analysis — current VCS importedPackage and wait AST rules remain their policy owners.
**Won't handle:** Markdown, archive, shell, and mixed asset observation — walkConformanceDocs and shippedFiles retain their current callers.
**Won't handle:** Targeted registry and metadata parsers — the census names their existing owners and purpose.
**Won't handle:** New cancellation or transaction protocols — the current source readers are synchronous and retain their existing call contract.
**Won't handle:** Concurrent filesystem mutation during one invocation — LTE's invocation snapshot remains the accepted consistency boundary.

## Ownership fences

The implementation fence permits the following exact paths and one accepted owner prefix.
It grants no authority to change the excluded policy, production commands, fixture meanings, or structure budgets.

- `internal/conformance/canonical_path_owner_test.go`
- `internal/conformance/cancel_signal_registrations_test.go`
- `internal/conformance/wait_deadline_literal_test.go`
- `internal/conformance/git_facts_checks_test.go`
- `internal/conformance/data_handling_test.go`
- `internal/conformance/build_contracts_test.go`
- `internal/conformance/check_bindings_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/sourcefiles/`
- `internal/bounds/classify.go`
- `internal/bounds/classify_bytes_test.go`
- `tests/canary/package-core-guard/bounds-classify-limit-restated`
- `tests/canary/package-core-guard/bounds-read-limit-restated`
- `internal/conformance/tier_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `specs/conformance-observation-follow-on/assets/`
- `reviews/conformance-observation-follow-on.md`

The sourcefiles prefix permits only the additive bounded no-follow observation capability and its landed interface binding.
The two bounds paths permit the narrow injected reader bridge and its preservation tests, without changing classifier policy or existing adapters.
It does not authorize a second owner, parser, classifier, or cache lifetime.
The tier files permit closure updates and existing contract preservation, not assertion removal.
The new reader port requires local classification preservation and real VCS caller witnesses, with no hidden global state.
Registry metadata, injected-port rows, anchors, and command registry paths remain read-only closure pins.

The two package-core-guard fixture units are read-only co-name pins for the classifier-owning slice.
Their complete artifacts and assertion purposes remain unchanged, including BASE, MUTATE.json, EXPECT, and CHECK.
Adding the reader port does not authorize either fixture mutation to target the port or change its expected diagnostic.

Canary artifacts remain read-only at tests/canary/canonical-path-owner/ and tests/canary/cancel-signal-registrations/.
The exact directories are second-derivation and own-signal-set, respectively.
Their BASE, MUTATE, and EXPECT files remain unchanged.
No ticket may claim a planning-only map or spec movement as implementation authority.

## Out of scope

The four LTE visitors and dispatch consolidation remain one approved separate outcome, with their existing nine tickets and gate plan.
They are not re-estimated or duplicated here.

Universal semantic/type analysis is a separate capability: provisional 8 edits, 3 gate runs.
The estimate comes from the five targeted parser families and three registry/integration surfaces in the census.
It is not a commitment or FT373 implementation plan.

Mixed asset observation is a separate capability: provisional 6 edits, 2 gate runs.
The estimate covers four named asset walkers, one owner, and one integration proof.

## Further notes

### Research questions and source status

| Question | Durable answer | Source and limit |
| --- | --- | --- |
| What owner is already approved? | LTE owns lazy source observations and its first four visitors | Full LTE spec, accepted design source, future API names not yet landed |
| Which residual visitors fit? | Six named policy visitors can consume shared source facts | Full six source files, binding and helper census at the research pin |
| What remains separate? | Targeted parsers and asset/metadata walks keep their actual owners | Exact 26-file lexical census plus ParseDir import/call census |
| What proves savings? | Real dispatch membership and exact observation operation counts | LTE46 precedent and planned aggregate fixture, no benchmark executed |
| How does VCS preserve its strict reader? | COF-C3 adds a bounded no-follow owner capability and a narrow existing-classifier reader bridge if needed | readBuildContractSource plus bounds/classify.go:125–210, concrete names bind after complete LTE |

Source pin: author worktree HEAD a9c395e77fec36d1f60f6d057c654aecf5720e24.
Production pin: main 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68, equal production at the author pin.
Sources include the complete architecture-observations map and its three resolved tickets, full FT365, full LTE spec, and ADR0006.

Structure census at this pin: conformance has 77 direct Go files, bounds has 6, and cmd/bench has 39.
Conformance has no directory headroom for another file.
Canonical has 74 lines, cancellation 178, wait 294, default-branch 69, control 322, and build-contract 328.
The ordinary file cap is 400.
The classifier has 210 lines and its byte-classification test file has 72 at this pin.

Checks has 641 lines under its 709 grant, and tier has 893 under its 973 grant.
Fixture-bite already has 872 lines against an 826 grant and is not granted for writing here.

After complete accepted LTE lands, recheck that tree's headroom before any implementation charge, because LTE can change these premises.
Shrinking or placing proofs in existing smaller visitor files is the default.
No budget or acceptance-list edit is authorized.

### Approval checkpoint

| Item | Author disposition | Independent reviewer disposition |
| --- | --- | --- |
| Implementation line | Recommend gpt-6.1-sol / high without changing LTE | Full spec accepted; ticket graph pending |
| Seam | Complete LTE owner through actual registered dispatch | Full spec accepted; ticket graph pending |
| Acceptance and hostile edges | COF01–COF38 preserve legacy policies and require omission witnesses | Full spec accepted; ticket graph pending |
| Ownership fence | Six visitors, dispatch adapters, bounded no-follow owner capability, narrow bounds reader port, enforcement closure | Full spec accepted; ticket graph pending |
| Scope cuts | Targeted metadata and universal analysis remain separate | Full spec accepted; ticket graph pending |

Reviewer attention: COF-C3 explicitly adds the bounded no-follow owner capability and narrow existing-classifier reader port if complete LTE lacks them.
This is part of the approved observation outcome, not a private fallback or a new classifier.

This is a planning artifact.
No mutation proof, test suite, benchmark, or implementation gate has been run for its future behavior.
Independent spec acceptance is recorded at the frozen hash in the review record.
Independent ticket review must accept this graph before implementation approval.

### Completion plan

The version 1 plan declares future verification obligations, not executed evidence.
Before approved implementation dispatch, the coordinator records fresh ticket authors through the required version 2 amendment.
Each chunk requires author verification and independent Standards, Spec, and Coverage acceptance before its successor.
After complete accepted LTE lands, concrete API binding and the landed-tree headroom recheck must finish before any implementation charge.
Independent ticket approval remains a prerequisite.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "COF-C1",
      "tickets": [
        "01-share-canonical-and-cancellation-observations.md"
      ],
      "verification": [
        {
          "id": "registered-canonical-cancellation",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "fresh-shared-owner",
          "command": "bench test --package ./internal/conformance",
          "probe": "create per-binding snapshots, retain an alias-hidden old canonical reader, and omit mutation/restoration freshness"
        }
      ]
    },
    {
      "id": "COF-C2",
      "tickets": [
        "02-share-deadline-and-byte-observations.md"
      ],
      "verification": [
        {
          "id": "registered-wait-and-bytes",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "legacy-policy-preservation",
          "command": "bench test --package ./internal/conformance",
          "probe": "replace byte matching with AST matching, remove prefilters, or apply strict VCS error policy to a legacy visitor"
        }
      ]
    },
    {
      "id": "COF-C3",
      "tickets": [
        "03-share-vcs-observations-and-close-all-visitors.md"
      ],
      "verification": [
        {
          "id": "classifier-preservation",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "observation-owner",
          "command": "bench test --package ./internal/conformance/sourcefiles"
        },
        {
          "id": "actual-ten-bindings",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "vcs-first-refusal",
          "command": "bench test --package ./internal/conformance",
          "probe": "collect the later z.go source violation before the unavailable FIFO a.go"
        },
        {
          "id": "bounded-and-legacy-reads",
          "command": "bench test --package ./internal/conformance",
          "probe": "reuse truncated strict bytes for oversized legacy source or raise the producer limit"
        },
        {
          "id": "whole-family-omissions",
          "command": "bench test --package ./internal/conformance",
          "probe": "restore a private source walker, omit one real binding, or retain one snapshot per binding"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check specs/conformance-observation-follow-on/spec.md"
    },
    {
      "id": "registered-consumers",
      "command": "bench test --package ./internal/conformance"
    },
    {
      "id": "bounds",
      "command": "bench test --package ./internal/bounds"
    },
    {
      "id": "source-owner",
      "command": "bench test --package ./internal/conformance/sourcefiles"
    }
  ]
}
```
