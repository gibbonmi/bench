# Independent expectation audit

Status: staged

Decision source: Named reviewed artifact `decisions/architecture-observations.md`, resolved tickets 1–3, 2026-10-06.

Verification log: 2 iterations to accept — one independent spec pass and one independent ticket pass.
Retained xhigh reviewer `deepen_execution` accepted the full spec before ticket slicing and accepted the graph at `c04ccd1bbac2274fb8197d592ea7aef21beffd94`.
Both checkpoints returned Standards 0, Spec 0, and Coverage 0, with confidence 9 on each axis.
Reviewed spec SHA-256: `a971cd21813f1f94efb1bee748ca4d846dfac295d5a9c886f881f2dd09104d29`.
This metadata close preserves all 35 rows, fences, and ticket bytes and claims no executed mutation or runtime proof.

Review evidence: [spec-review.md](assets/spec-review.md).

## Problem

A test expectation can duplicate implementation knowledge for a useful reason.
An independently authored expectation can detect an omitted producer row that a derived expectation would also omit.
Unnecessary copies instead drift from their source.

The current AXI, dispatch, routing, help, and bounds expectations need an explicit purpose census.
An independence comment alone is not the required demonstrated red.
Exact numeric copies need a named mutation and recorded red before their independence is retained as justified.

## Solution

Audit the named expectation families and their real consumers.
Classify each copy as an omission oracle, a policy-value oracle, a derived fact, or incidental fixture arithmetic.
Retain justified independence only with a named, demonstrated mutation and restoration record.
Keep producer policy, parsers, executable registries, fixture harnesses, and derived counts single-sourced.

The audit keeps every existing assertion and behavior.
It does not generate an expected membership set from the producer it grades.
If a proposed disposition would weaken an omission proof, stop that change for review.
No mutation, test, or benchmark is executed during this planning stage.

## User stories

Line: gpt-6.1-sol / high.

Implementation-line reason: The difficult work is proving that independent expectations still reject omissions through actual consumers.
Current mutation tests provide strong seams.
This research did not establish recorded red evidence for every exact numeric pair.
IEA-C3 carries that proof obligation.
This recommendation does not change LTE's approved gpt-5.6-sol / high line.

Harder chunks: IEA-C3.

1. As a maintainer, I want a current expectation census, so that historical counts cannot decide today's changes.
2. As a maintainer, I want each copy to have a purpose, so that incidental repetition is not mistaken for policy duplication.
3. As a maintainer, I want AXI membership independence preserved, so that a removed approved query still turns red.
4. As a maintainer, I want extra AXI approvals rejected, so that the registry cannot broaden the contract silently.
5. As a maintainer, I want AXI guidance checked independently, so that missing or malformed promises stay visible.
6. As a maintainer, I want real command envelope coverage preserved, so that a missing fixture cannot erase a query's proof.
7. As a maintainer, I want function binding independence preserved, so that swapped implementations still turn red.
8. As a maintainer, I want tier binding independence preserved, so that moved tier metadata still turns red.
9. As a maintainer, I want subject binding independence preserved, so that changed subject metadata still turns red.
10. As a maintainer, I want routing coverage independence preserved, so that a new or removed route cannot pass without its guard.
11. As a maintainer, I want grammar reachability preserved, so that a routing claim cannot grade a comment instead of executable parsing.
12. As a maintainer, I want help inventory independence preserved, so that a deleted public row cannot disappear from both sides.
13. As a maintainer, I want each duration policy copy proved, so that numeric independence has a demonstrated purpose.
14. As a maintainer, I want each byte-limit policy copy proved, so that arbitrary numeric changes still turn red.
15. As a maintainer, I want each shift policy copy proved, so that exact limits and defaults remain intentional.
16. As a maintainer, I want concrete deadline entry coverage preserved, so that a derived value cannot hide a changed command timeout.
17. As a maintainer, I want derived inventories kept derived, so that this audit does not introduce new hand-maintained lists.
18. As a maintainer, I want mutation restoration verified, so that a proof cannot leave altered source behind.
19. As a maintainer, I want evidence tied to its source, so that a later reader can reproduce the exact red.
20. As a maintainer, I want silent or invalid mutations rejected, so that comments cannot substitute for proof.
21. As a maintainer, I want hostile fixtures preserved, so that malformed source and guidance still fail at their existing boundaries.
22. As a maintainer, I want exact closure fences, so that a purpose audit cannot change production policy or command grammar.
23. As a maintainer, I want bounded proof code, so that the audit does not grow another registry or harness.
24. As a maintainer, I want adjacent policy domains left with their owners, so that this audit does not become a repository-wide rewrite.

## Implementation decisions

### Independence rule and evidence contract

AGENTS.md and ADR0006 are the authority for the expectation exception.
Independence must be necessary for a named omission or mutation to turn the gate red.
That red must be recorded and demonstrated.
The exception covers only the expectation-versus-implementation pair.
It does not authorize duplicate production policy, parser logic, harnesses, executable registries, or derived counts.

This audit produces a spec-local expectation evidence record during implementation.
Each record names the exact producer, independent expectation, real consumer, purpose, mutation, source pin, focused command, red result, and restoration result.
It includes complete or digest-pinned diagnostic evidence.
The record distinguishes historical claimed evidence from freshly demonstrated evidence.
A planning row is not a recorded red.

Retain all current assertions while collecting evidence.
A silent or invalid probe cannot justify independence.
If the named mutation does not bite, correct the proof within the approved purpose or return the exact contradiction for review.
Do not delete the expectation merely to make the audit green.
Do not weaken a numerical comparison or derive its expected value from the mutated constant.
An unresolved required red is a completion blocker, not an accepted independence exemption.

### Current copy and consumer census

Research pin: a9c395e77fec36d1f60f6d057c654aecf5720e24.
The census covers the named families below and their transitive real consumers.
It is not a census of every independently authored expectation in the repository.

| Family and producer | Current copies or counterpart facts | Actual consumer and disposition |
| --- | --- | --- |
| Production AXI disposition: cmd/bench/main.go:68–159 commandRegistry | internal/conformance/axi_query_registry_test.go:18 approvedAXIQueries | checkAXIQueryRegistry and TestAXIMembershipExpectationBitesInBothDirections retain independent membership omission and addition proofs |
| AXI approved query advertisement | craft-cli query table and projects/benchkit.md AXI query seam | parseGuidanceQueries and parseProfileAXIQueries grade independently authored advertisements against approved membership |
| AXI principles | axi_query_registry_test.go:31 expectedAXIPrinciples and craft-cli guidance headings | AXI principle omission and numbering remain independent policy checks |
| AXI command test stimuli | cmd/bench/command_registry_test.go:220 axiEnvelopeCases and help_inventory_test.go:294 assessmentEnvelopeCases | TestAXIRegistryBindsEachRealCommandEnvelope checks exact fixture membership and each real Command.Run envelope |
| Executable check binding versus registry metadata | check_bindings_test.go:15 implementation, tier, subject and checks_test.go binding rows | tier_test.go:103 registryAgreementDiags plus CM5/CM6/CM7 tests retain only the facts that catch the named mutations |
| Routing policy versus real dispatch | subcommand_routing_table_test.go:13 subcommandRouting with 64 current entries | checkSubcommandRouting parses production dispatch and requires independent route/exemption coverage plus actual usage.Parse reachability |
| Public command help versus runtime inventory | cmd/bench/help_inventory_test.go:62 TestHelpInventoryIsComplete independently authored full help text | Actual Command.Run help output retains the deleted-row omission oracle |
| Bounds production scalar policy | internal/bounds/bounds.go:27–94 and thirteen exact copies at bounds_test.go:85–94 | TestProductionPolicyValues retains exact values only after thirteen named mutations are demonstrated red |
| Concrete command deadline | sessioninspect.go providerTimeout alias and sessioninspect_test.go:70 concrete 9s–10s remaining interval | TestCommandInstallsTenSecondDeadline reaches real Command and retains its independently authored ten-second policy witness after a named red |
| Derived duration inventory | bounds_test.go:154 registryDurations reads and type-checks production scalar declarations | TestDeadlineExceedsEveryRegistryDuration stays derived, without a second numeric inventory |
| Derived fixture and timing universes | canary.Fixtures and registry.Checks | Fixture completion sets and timing membership/cardinality remain derived under complete LTE, with existing omission mutations intact |
| Local selected-query arithmetic | cmd/bench/selected_queries_test.go: concrete six-line and eight-line fixture responses | Fixture-size arithmetic is incidental repetition, not a copy of a production numeric policy |

The current AXI expectation covers anchors, learnings, maps, guards, diff, coverage, consumers, harnesses, and roadmap root queries.
It also covers assessment list/show/compare and worktree list.
The real runtime registry, conformance expectation, envelope stimuli, skill advertisement, and project advertisement are five distinct surfaces.
They have different roles and must not become one producer-derived oracle.

FT365's 62 routing rows are historical.
The current independent routing map has 64 entries at this pin.
This audit preserves that complete current membership and every routed/exempt reason.
It does not turn grammar policy into runtime registry metadata.

A resolved Go consumer query covered all thirteen numeric symbols in TestProductionPolicyValues.
It found 106 static reference edges across 54 files and 327 loaded packages.
Its complete answer hash is `dc950daca3791d536f8975a24096d2cebcb6d5477b7e791ac8116161148be7f4`.

Exact resolved-reference query:

```text
bench consumers bounds.ProviderTimeout bounds.EnvironmentDiscoveryTimeout bounds.GitRefreshTimeout bounds.GuardScanTimeout bounds.GateTimeout bounds.ModelReadLimit bounds.OutlineFileLimit bounds.ControlRecordLimit bounds.IterationMin bounds.IterationMax bounds.MainIterationsDefault bounds.RefactorIterationsDefault bounds.MaxWall --full
```

This exposed the additional concrete session-inspection deadline assertion.
Static references do not prove the absence of every unrelated numeric literal.
The audit names this limit rather than claiming repository-wide completeness.

The numeric production consumers include models, sessioninspect, refresh, guards, gate, outline, control-record readers, and shift.loop.
No consumer behavior, policy value, or production alias changes in this outcome.
ProviderTimeout and EnvironmentDiscoveryTimeout aliases retain their existing command integration.
ControlRecordLimit's many storage readers retain their independent hostile-input fixtures.

### Exact numeric mutation matrix

The following thirteen pairs are the complete exact-value table in TestProductionPolicyValues.
Each future probe changes one production declaration, runs that real test, records red, and verifies restored source and green.
Changing more than one declaration in a probe is not evidence for each individual pair.

| Symbol | Independent expected value | Named future mutation | Required red category |
| --- | --- | --- | --- |
| ProviderTimeout | 10*time.Second | 10 to 11 seconds | duration policy changed |
| EnvironmentDiscoveryTimeout | 2*time.Second | 2 to 3 seconds | duration policy changed |
| GitRefreshTimeout | 30*time.Second | 30 to 31 seconds | duration policy changed |
| GuardScanTimeout | 5*time.Second | 5 to 6 seconds | duration policy changed |
| GateTimeout | 45*time.Minute | 45 to 46 minutes | duration policy changed |
| ModelReadLimit | 5<<20 | 5 to 6 MiB | read/output policy changed |
| OutlineFileLimit | 2<<20 | 2 to 3 MiB | read/output policy changed |
| ControlRecordLimit | 2<<20 | 2 to 3 MiB | read/output policy changed |
| IterationMin | 1 | 1 to 2 | shift policy changed |
| IterationMax | 100 | 100 to 101 | shift policy changed |
| MainIterationsDefault | 12 | 12 to 13 | shift policy changed |
| RefactorIterationsDefault | 4 | 4 to 5 | shift policy changed |
| MaxWall | 24*time.Hour | 24 to 25 hours | shift policy changed |

ProviderTimeout's 11-second mutation also runs TestCommandInstallsTenSecondDeadline through sessioninspect.Command.
Its remaining deadline exceeds the independent ten-second upper bound and must turn red.
That entry test has a distinct consumer purpose from the scalar table.
An alias or TestDeadline-derived expected window cannot replace it.

### Named omission witnesses

| Independent purpose | Future witness and real observed boundary | Cheapest wrong change rejected |
| --- | --- | --- |
| AXI missing approved member | Existing TestAXIMembershipExpectationBitesInBothDirections removes anchors approval in parsed production source | Deriving approvedAXIQueries from commandRegistry would omit anchors from both sides |
| AXI extra approval | Same test promotes status from exempt to approved root | An expectation broadened from runtime membership would accept the unapproved query |
| AXI guidance promises | Existing TestAXIGuidanceContractBites omits a principle/query and supplies malformed table rows | Guidance rendered only from the runtime list would hide an omitted authored obligation |
| AXI fixture membership | Omit the anchors entry from axiEnvelopeCases, then run TestAXIRegistryBindsEachRealCommandEnvelope | Generating fixture keys from the runtime set would conceal missing envelope stimulus |
| CM5 function independence | Existing TestConformanceMetaBites swaps line-routing and bounds implementations | Registry-derived binding identity would agree with the wrong registry metadata |
| CM6 tier independence | Existing TestConformanceMetaBites changes advertised dev tier while binding stays unchanged | One shared tier fact would hide metadata drift |
| CM7 subject independence | Existing TestConformanceMetaBites changes the advertised meta subject while binding stays unchanged | One shared subject fact would hide root/kit drift |
| Routing exhaustive membership | Existing routing tests delete a guard entry and add a dispatch name without a table row | A dispatch-derived table would hide the missing policy |
| Grammar reachability | Existing TestSubcommandRoutingRoutedClaimBites removes executable usage.Parse and leaves comment/string decoys | Lexical advertisement is not executable grammar coverage |
| Public help inventory | Omit the public models row from production commandRegistry, then run TestHelpInventoryIsComplete through real help | Deriving the expected help inventory would erase the deleted row from both sides |
| Numeric policy | Each matrix mutation runs TestProductionPolicyValues | Expected values obtained from the mutated constants would agree with the wrong policy |
| Concrete deadline | ProviderTimeout 10s to 11s runs sessioninspect.Command through TestCommandInstallsTenSecondDeadline | A producer-derived deadline expectation would accept the new duration |

Existing tests that plant synthetic source remain valuable proofs at their declared seam.
A future bench probe verifies each retained family against its actual gate consumer where source mutation is necessary.
The mutation is temporary, focused, and restored by the supported probe path.
No mutation changes the committed production tree.

### Closure and headroom contract

The command registry closure includes main.go, command_registry.go, command_registry_test.go, help_inventory_test.go, axi_query_registry_test.go, and subcommand_routing_table_test.go.
Production registry and grammar paths are read-only mutation targets, not committed-write authority.
The internal/tickets/registry_data.go commandRegistries rule remains unchanged.
No CLI name, argv grammar, help contract, AXI envelope, exit code, or scope changes.

Anchor registries remain read-only because guidance and advertisement bytes do not change.
Canary inventory, the universal fixture caller, retained EXPECT files, and injected-port registry remain unchanged.
If a future proof needs a new live-tree reader, its classification must close in tier_live_tree_test.go.
No second parser or fixture harness is permitted.

Conformance has 77 direct Go files and no file-count headroom.
The AXI source has 445 lines, command_registry_test has 794, and help_inventory_test has 323.
Do not grow oversized source files to add evidence narration.
Use spec-local evidence and existing smaller tests when a new proof is required.

Bounds_test has 312 lines, selected_queries has 186, and sessioninspect_test has 145.
Recheck limits at implementation source binding, including complete LTE changes to shared dispatch evidence.
No budget or acceptance-list change is authorized.

## Implementation chunks

Independent full-spec acceptance preceded these three serial vertical tickets.
Each chunk delivers retained consumer behavior plus its demonstrated purpose record.
A documentation-only inventory is not an intermediate green checkpoint for unproved independence.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| IEA-C1 / 01-prove-axi-and-binding-independence.md | AXI and binding independence remain live with named omission evidence | IEA01–IEA11, IEA27, IEA31 | Existing AXI, envelope, and CM5/6/7 consumers | no |
| IEA-C2 / 02-prove-routing-and-help-independence.md | Routing and public help keep complete omission proofs | IEA12–IEA15, IEA28–IEA30, IEA32 | Existing routing and real help consumers | no |
| IEA-C3 / 03-prove-numeric-policies-and-close-evidence.md | Every named numeric copy has demonstrated red and complete restoration evidence | IEA16–IEA26, IEA33–IEA35 | TestProductionPolicyValues and real sessioninspect.Command entry | yes |

IEA-C2 depends on IEA-C1.
IEA-C3 depends on IEA-C2 for the common evidence contract.
The complete landed LTE outcome is required before changing or grading shared dispatch/fixture architecture.
This audit does not duplicate LTE's implementation or alter its approved line.

### Stable chunk mapping and checkpoint contracts

IEA23–IEA26 move from IEA-C1 to IEA-C3.
All other row assignments and stable chunk IDs remain unchanged.
The terminal slice owns complete-family reconciliation, not the first valid checkpoint for an earlier consumer.
Each ticket applies freshness, hostile-input preservation, owner consumption, headroom, and valid evidence requirements to its own introduced behavior.
No earlier ticket borrows a later repair.

## Testing decisions

Tests attach to existing real consumers.
Keep independently authored expected sets apart from producer lists and AST-derived inventories.
Exact membership proofs check both missing and extra rows.
Synthetic values remain explicitly authored fixtures.

The future evidence record contains one named mutation per justified pair or omission purpose.
Focused reds include the expected diagnostic category and actual test identity.
Restoration includes the exact restored digest and a passing repeat of the same focused test.
A compile failure, malformed mutation, unrelated failure, or restoration failure is not a successful omission proof.
No elapsed-time savings claim belongs to this audit.

### Seam diagram

```text
Production registry or numeric policy
    -> actual parser, dispatcher, or command consumer
    -> independently authored expectation
    -> named omission or mutation turns the existing focused test red
    -> supported restoration and same focused test green
Evidence binds producer, expectation, purpose, source, command, red, and restoration.
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| IEA01 | 1,2 | Census records each named copy's producer, consumer, and distinct purpose | review-owned expectation evidence census | A historical count or untraced copy cannot receive a justified disposition |
| IEA02 | 3 | Removing anchors AXI approval retains independent membership red | `internal/conformance/axi_query_registry_test.go` (`TestAXIMembershipExpectationBitesInBothDirections`) | Producer-derived membership would remove the member from both sides |
| IEA03 | 4 | Promoting status to AXI approved root retains extra-member red | `internal/conformance/axi_query_registry_test.go` (`TestAXIMembershipExpectationBitesInBothDirections`) | A broadened expectation would accept runtime metadata as policy authority |
| IEA04 | 5 | Missing principle and approved guidance query retain their existing diagnostics | `internal/conformance/axi_query_registry_test.go` (`TestAXIGuidanceContractBites`) | Independent authored guidance promises cannot disappear with a generated table |
| IEA05 | 5,21 | Malformed, duplicate, missing, or conflicting AXI declarations remain fail-closed | `internal/conformance/axi_query_registry_test.go` (`TestAXIRegistryParserFailsClosed`) | A rewritten parser or relaxed evidence rule cannot silently accept malformed policy |
| IEA06 | 6 | Every actual approved query retains exact envelope fixture membership | `cmd/bench/command_registry_test.go` (`TestAXIRegistryBindsEachRealCommandEnvelope`) | Omitting the anchors fixture differs from independently checked runtime membership |
| IEA07 | 6,21 | Every envelope stimulus still reaches actual Command.Run | `cmd/bench/command_registry_test.go` (`TestAXIRegistryBindsEachRealCommandEnvelope`) | A helper-only record cannot prove real stdout, refusal, empty, and deep-cwd behavior |
| IEA08 | 7 | CM5 retains independently bound function identity | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`) | Swapped line-routing and bounds implementations cannot agree with derived identity |
| IEA09 | 8 | CM6 retains independently bound tier | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`) | Advertised dev-to-ship drift still differs from the executable binding |
| IEA10 | 9 | CM7 retains independently bound subject | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`) | Root/kit metadata drift still differs from the actual executable subject |
| IEA11 | 17 | Order, inputs, meta membership, and derived counts remain registry-owned | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`, `TestTimingOrderStable`) and review-owned copy census | A new copied policy field would exceed the narrow binding independence exception |
| IEA12 | 10 | Routing retains all 64 current entries and existing exemption reasons | `internal/conformance/subcommand_routing_test.go` (`TestSubcommandRoutingRegistryBites`, `TestSubcommandRoutingRequiresFollowOnGuardEntry`) | New and orphan route mutations cannot disappear from a producer-derived guard table |
| IEA13 | 11 | Routed claims require executable usage.Parse reachability | `internal/conformance/subcommand_routing_test.go` (`TestSubcommandRoutingRoutedClaimBites`) | Comment/string decoys remain red after the real grammar call is omitted |
| IEA14 | 10,11 | Doctor routing and worktree dispatcher exemption retain their existing purpose | `internal/conformance/subcommand_routing_test.go` (`TestSubcommandRoutingGradesDoctorLeafThroughUsageParse`, `TestSubcommandRoutingKeepsWorktreeDispatcherExempt`) | Sharing a name set cannot erase independent grammar/exemption policy |
| IEA15 | 12 | Deleting the models public row turns real help inventory red | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A runtime-derived expected string would omit the same public row |
| IEA16 | 13 | ProviderTimeout 10s to 11s produces recorded duration-policy red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | A constant-derived expected duration would agree with the mutant |
| IEA17 | 13 | EnvironmentDiscoveryTimeout 2s to 3s produces recorded duration-policy red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | The exact two-second policy must remain independently observed |
| IEA18 | 13 | GitRefreshTimeout, GuardScanTimeout, and GateTimeout each produce their separate matrix red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | One multi-value probe cannot stand in for three independent policy witnesses |
| IEA19 | 14 | ModelReadLimit, OutlineFileLimit, and ControlRecordLimit each produce their separate matrix red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | Each exact byte limit requires its own demonstrated purpose |
| IEA20 | 15 | IterationMin and IterationMax each produce their separate matrix red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | A producer-derived range would accept either changed endpoint |
| IEA21 | 15 | MainIterationsDefault, RefactorIterationsDefault, and MaxWall each produce their separate matrix red | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | Each exact default and wall limit must reject its own mutant |
| IEA22 | 16 | ProviderTimeout 11s also turns the real command deadline test red | `internal/sessioninspect/sessioninspect_test.go` (`TestCommandInstallsTenSecondDeadline`) | Deriving the expected interval from the policy would hide changed command behavior |
| IEA23 | 18 | Every future probe verifies byte-identical source restoration | review-owned evidence record with supported bench probe restoration and digest | A red followed by altered source cannot count as demonstrated purpose |
| IEA24 | 18 | The same focused consumer passes after restoration | review-owned evidence record with focused command and restored green | An unrelated baseline failure cannot masquerade as a successful mutation proof |
| IEA25 | 19 | Each record pins source, producer, expectation, command, diagnostic, and restoration | review-owned spec-local expectation evidence record | A historical claim without reproducible source identity remains unproved |
| IEA26 | 20 | Silent, invalid, compile-only, unrelated, and restore-failed probes are refused as proof | review-owned evidence classification against actual test diagnostic | Merely running a mutator cannot justify duplicate knowledge |
| IEA27 | 3,7,8,9 | No retained omission oracle is derived from the producer it grades | review-owned diff plus existing AXI and CM mutation consumers | An oracle rewrite that agrees with omission violates the approved purpose |
| IEA28 | 12,17 | Runtime help remains a single registry projection | `cmd/bench/help_inventory_test.go` (`TestHelpRendersPublicCommandRegistryRows`, `TestHelpRendersTreeTargetFromScope`) | Retaining test independence cannot authorize a second production help list |
| IEA29 | 17 | Shared command registry syntax remains parseRegistryTable-owned | `internal/conformance/axi_query_registry_test.go` (`TestAXIRegistryParserFailsClosed`) and review-owned parser diff | A new audit parser would duplicate production-source interpretation |
| IEA30 | 2,17 | Local six-line and eight-line fixture arithmetic stays distinct from numeric policy | `cmd/bench/selected_queries_test.go` (`TestSelectedWorktreeWithinResponseBound`, `TestSelectedHistoryWithinResponseBound`) and review-owned purpose census | An incidental count cannot become a blanket exception for copied policy |
| IEA31 | 22 | AXI and binding closure retains command, anchor, canary, and port registries unchanged | review-owned exact path diff and existing registry agreement tests | A purpose audit cannot silently change executable or advertised policy |
| IEA32 | 24 | Adjacent retired-grammar and phase-invocation oracles remain with their existing owners | `cmd/bench/command_registry_test.go` (`TestRemovedGrammarsRefuseThroughTheirFamily`) and `internal/conformance/skills_index_checks_test.go` (`TestCommandInvocationPolicyGradesTableCompleteness`) | Closed omission policies cannot be rewritten as an audit side effect |
| IEA33 | 17 | registryDurations still discovers all scalar duration declarations | `internal/bounds/bounds_test.go` (`TestDeadlineExceedsEveryRegistryDuration`) | Replacing it with the thirteen exact-policy copies loses future duration coverage |
| IEA34 | 23 | Evidence adds no second harness, scalar inventory, parser, or oversized module | review-owned fence and landed-source headroom check | A registry-shaped evidence implementation would duplicate the knowledge being audited |
| IEA35 | 19,20 | Every retained independent pair has demonstrated named red before completion | review-owned complete evidence table checked against the named census | A missing numeric red blocks completion rather than receiving an unresolved exemption |

### Edge inventory

| Class | Disposition at this seam |
| --- | --- |
| Missing or extra membership | Preserve both-direction AXI and routing omission witnesses |
| Duplicate/conflicting/unknown declaration | Keep the existing fail-closed named registry parser |
| Missing, malformed, oversized, symlink, or invalid guidance | Keep existing bounds classification and its early diagnostic winner |
| Comments, strings, aliases, dot imports | Retain each current parser's exact semantics and grammar-decoy refusals |
| Deep cwd and hostile argv | Existing real AXI envelope tests preserve stdout, usage, refusal, and root scope |
| Empty query results and bounded response spilling | Preserve current real-command fixtures and complete spill reads |
| Mutation fails to compile or changes unrelated source | Classify invalid rather than demonstrated policy red |
| Silent mutation or missing diagnostic | Block independence justification and return the exact purpose gap |
| Restoration failure or digest mismatch | Refuse evidence and restore source before continuing |
| Several numeric declarations mutated together | Reject as proof for individual pairs |
| Existing expected assertion disagrees with a new projection | Preserve the original assertion and investigate the projection |

**Won't handle:** Repository-wide literal deduplication — the named consumer census defines this audit's complete scope.
**Won't handle:** Independent production policy copies — the ADR exception never authorizes them, and current policy owners remain canonical.
**Won't handle:** Universal Go type analysis — registryDurations retains its bounded scalar-registry type check under its existing owner.
**Won't handle:** New CLI grammar or guidance publication — current command definitions, anchors, and authored guidance remain unchanged.
**Won't handle:** Reopening retired grammar and phase-invocation decisions — their existing named omission tests remain the surviving callers.

## Ownership fences

The implementation fence permits only purpose-preserving test clarification, focused proof support, and spec-local evidence.
Temporary producer mutations use supported probes and must never join a commit.

- `internal/conformance/check_bindings_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/tier_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/bounds/bounds_test.go`
- `internal/sessioninspect/sessioninspect_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/selected_queries_test.go`
- `specs/independent-expectation-audit/assets/`
- `reviews/independent-expectation-audit.md`

The registry table files permit annotations and proof-preserving clarification only.
Their member sets, expected values, routed policies, exemption reasons, and independent identities remain fixed.
The shared binding files remain subject to complete LTE's accepted contract.
The command_registry.go entry is a required command-package co-name pin and must remain byte-identical in every delivered diff.
No production policy, grammar, advertisement, anchor registry, canary fixture, injected-port row, or budget file may change.
The named command registry closure must be checked even when its production paths remain read-only.

## Out of scope

Source observation reuse is a separate outcome with its own conformance-observation-follow-on spec.
LTE owns the first four source visitors and dispatch consolidation.
Neither is duplicated or re-estimated here.

A repository-wide independent-expectation inventory is a separate capability: provisional 8 edits, 2 gate runs.
The estimate covers six additional policy-family inventories, one evidence assembly surface, and one integration proof.
It is not an authorized blanket rewrite.

New production registry schema or grammar consolidation is a separate capability: provisional 6 edits, 2 gate runs.
The estimate covers the production registry, shared syntax owner, routing consumer, help consumer, AXI consumer, and integration proof.

## Further notes

### Research questions and source status

| Question | Durable answer | Source and limit |
| --- | --- | --- |
| Which copies currently exist? | Named AXI, binding, routing, help, numeric, and entry-deadline pairs have different purposes | Full owner/test definitions and current resolved consumer queries |
| Which expectations must stay independent? | CM5/CM6/CM7, AXI membership/guidance, routing, and public-help omission oracles | ADR0006 and existing named mutation assertions, no red executed in this planning stage |
| What requires new evidence? | Thirteen numeric policy pairs and the concrete deadline entry expectation need demonstrated red records | Exact current table and transitive consumer census |
| What is already derived? | Duration, fixture, timing, and runtime help inventories retain their existing owners | Current derivation definitions and complete LTE contract |
| What is not proved? | This spec does not show current reds, universal literal completeness, or measured speed savings | Future implementation evidence and bounded static-query limitations |

Source pin: author worktree HEAD a9c395e77fec36d1f60f6d057c654aecf5720e24.
Production pin: main 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68, equal production at the author pin.
Sources include the full architecture-observations map and resolved tickets, FT365, LTE spec, and ADR0006.
Primary code sources are cited in the census and mutation matrix.
The source-map artifact stays in place because these two outcomes share it.

### Approval checkpoint

| Item | Author disposition | Independent reviewer disposition |
| --- | --- | --- |
| Implementation line | Recommend gpt-6.1-sol / high without changing LTE | Full spec accepted; ticket graph pending |
| Seam | Existing real consumers with independent expected sets and exact temporary mutations | Full spec accepted; ticket graph pending |
| Acceptance and hostile edges | IEA01–IEA35 preserve existing assertions and require recorded red/restoration | Full spec accepted; ticket graph pending |
| Ownership fence | Named tests and spec-local evidence, production and guidance remain read-only | Full spec accepted; ticket graph pending |
| Scope cuts | Repository-wide audit and new registry schema remain separate | Full spec accepted; ticket graph pending |

This is an independently accepted spec with a planning ticket graph.
It is not a demonstrated red record.
Independent spec acceptance is recorded at the frozen hash in the review record.
Independent ticket review must accept this graph before implementation approval.

### Completion plan

The version 1 plan declares future verification obligations, not executed evidence.
Before approved implementation dispatch, the coordinator records fresh ticket authors through the required version 2 amendment.
Each chunk requires author verification and independent Standards, Spec, and Coverage acceptance before its successor.
Complete LTE source binding and independent ticket approval remain prerequisites.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "IEA-C1",
      "tickets": [
        "01-prove-axi-and-binding-independence.md"
      ],
      "verification": [
        {
          "id": "axi-and-meta-consumers",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "actual-envelopes",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "axi-omissions",
          "command": "bench test --package ./internal/conformance",
          "probe": "remove anchors approval, promote status approval, and omit principle/query guidance separately"
        },
        {
          "id": "envelope-fixture-omission",
          "command": "bench test --package ./cmd/bench",
          "probe": "omit the anchors axiEnvelopeCases stimulus"
        },
        {
          "id": "cm5-cm6-cm7",
          "command": "bench test --package ./internal/conformance",
          "probe": "separately swap implementation, change advertised tier, and change advertised subject"
        }
      ]
    },
    {
      "id": "IEA-C2",
      "tickets": [
        "02-prove-routing-and-help-independence.md"
      ],
      "verification": [
        {
          "id": "routing-consumers",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "actual-help",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "routing-omissions",
          "command": "bench test --package ./internal/conformance",
          "probe": "delete a guard entry, add an uncovered dispatch, and omit executable usage.Parse while leaving comment/string decoys"
        },
        {
          "id": "help-row-omission",
          "command": "bench test --package ./cmd/bench",
          "probe": "temporarily remove the public models production registry row"
        }
      ]
    },
    {
      "id": "IEA-C3",
      "tickets": [
        "03-prove-numeric-policies-and-close-evidence.md"
      ],
      "verification": [
        {
          "id": "scalar-policy-consumer",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "actual-deadline",
          "command": "bench test --package ./internal/sessioninspect"
        },
        {
          "id": "numeric-ProviderTimeout",
          "command": "bench test --package ./internal/bounds",
          "probe": "ProviderTimeout: 10s to 11s; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-EnvironmentDiscoveryTimeout",
          "command": "bench test --package ./internal/bounds",
          "probe": "EnvironmentDiscoveryTimeout: 2s to 3s; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-GitRefreshTimeout",
          "command": "bench test --package ./internal/bounds",
          "probe": "GitRefreshTimeout: 30s to 31s; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-GuardScanTimeout",
          "command": "bench test --package ./internal/bounds",
          "probe": "GuardScanTimeout: 5s to 6s; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-GateTimeout",
          "command": "bench test --package ./internal/bounds",
          "probe": "GateTimeout: 45m to 46m; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-ModelReadLimit",
          "command": "bench test --package ./internal/bounds",
          "probe": "ModelReadLimit: 5MiB to 6MiB; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-OutlineFileLimit",
          "command": "bench test --package ./internal/bounds",
          "probe": "OutlineFileLimit: 2MiB to 3MiB; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-ControlRecordLimit",
          "command": "bench test --package ./internal/bounds",
          "probe": "ControlRecordLimit: 2MiB to 3MiB; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-IterationMin",
          "command": "bench test --package ./internal/bounds",
          "probe": "IterationMin: 1 to 2; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-IterationMax",
          "command": "bench test --package ./internal/bounds",
          "probe": "IterationMax: 100 to 101; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-MainIterationsDefault",
          "command": "bench test --package ./internal/bounds",
          "probe": "MainIterationsDefault: 12 to 13; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-RefactorIterationsDefault",
          "command": "bench test --package ./internal/bounds",
          "probe": "RefactorIterationsDefault: 4 to 5; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "numeric-MaxWall",
          "command": "bench test --package ./internal/bounds",
          "probe": "MaxWall: 24h to 25h; require the accepted exact policy category, restored digest, and same focused green"
        },
        {
          "id": "concrete-deadline-policy",
          "command": "bench test --package ./internal/sessioninspect",
          "probe": "ProviderTimeout 10s to 11s must exceed the independent ten-second command deadline upper bound"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check specs/independent-expectation-audit/spec.md"
    },
    {
      "id": "axi-meta-routing",
      "command": "bench test --package ./internal/conformance"
    },
    {
      "id": "command-consumers",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "scalar-policy",
      "command": "bench test --package ./internal/bounds"
    },
    {
      "id": "command-deadline",
      "command": "bench test --package ./internal/sessioninspect"
    }
  ]
}
```
