# Repair collection pilot

Status: staged

Decision source: ready compiled map at specs/repair-collection-pilot/decisions/ft232-repair-loop.md

Verification log: pending independent Sol/high review

## Problem

Existing records cannot distinguish stalled repairs from productive repairs with a repeated first failure.
The reviewer needs attributed examples before deciding whether a detector is useful.
The bounded repair policy already applies and remains outside this build.

## Solution

Provide an explicitly activated collection pilot for ordinary work in the Bench kit.
Collect attributed observations and evidence through a small local command.
Audit proposed progress labels against verification before using those labels in the report.
Stop collection after ten completed repair sequences or fourteen calendar days, whichever occurs first.
Report the evidence, unsupported labels, incomplete sequences, and missing example classes.

FT232 stays open after this implementation.
The implementation delivers collection and report capability, plus the operating protocol.
The reviewer activates the real pilot separately and evaluates its eventual report.
Neither spec approval nor implementation starts that observation period.

## User stories

Line: gpt-5.6-sol / high
Implementation-line reason: RP-C2 combines explicit identity, partial failure evidence, and audited progress at a new domain behind the established CLI seam
Harder chunks: RP-C2

### Activation and custody

1. As a reviewer, I want explicit activation, so that routine work does not start an experiment.
2. As a reviewer, I want collection limited to the Bench kit, so that linked repositories retain their defaults.
3. As an operator, I want durable local records, so that worktree release does not erase evidence.
4. As an operator, I want repeat activation to preserve the start, so that retries cannot extend the pilot.
5. As an operator, I want honest refusals for invalid inputs, so that bad records cannot become evidence.
6. As an operator, I want safe concurrent updates, so that overlapping assignments do not erase observations.

### Repair evidence

7. As a reviewer, I want each sequence bound to one chunk, so that a new process cannot reset its identity.
8. As a reviewer, I want each observation bound to an assignment and source, so that timing cannot invent attribution.
9. As a reviewer, I want several findings in one sequence, so that finding counts do not inflate sequence counts.
10. As a reviewer, I want pre-review and post-review observations, so that the sequence covers the chunk's repairs.
11. As a reviewer, I want comparable failures, so that repeated first failures can be examined.
12. As a reviewer, I want partial failure evidence marked, so that a first red cannot impersonate the complete red set.
13. As a reviewer, I want unchanged reruns separate from repairs, so that repetition alone does not prove a stalled repair.
14. As a reviewer, I want verified closure, so that an author assertion cannot complete a successful sequence.
15. As a reviewer, I want reviewer handoff to end a sequence, so that unresolved work has an honest endpoint.
16. As an operator, I want idempotent record imports, so that retries do not duplicate evidence.

### Bounds and evidence audit

17. As a reviewer, I want collection to stop at ten completed sequences, so that the experiment remains bounded.
18. As a reviewer, I want collection to stop after fourteen days, so that sparse work cannot prolong the experiment.
19. As a reviewer, I want open sequences reported as incomplete, so that the deadline does not manufacture closure.
20. As a reviewer, I want each progress claim tied to verified improvement, so that green checks alone cannot establish progress.
21. As a reviewer, I want evidence audits separate from proposals, so that an unsupported label remains unknown.
22. As a reviewer, I want contradictory evidence retained, so that later reports cannot hide uncertainty.
23. As a reviewer, I want audits after collection stops, so that the frozen observations can receive a complete assessment.

### Report and adoption decision

24. As a reviewer, I want stalled repair examples, so that the report examines the detector's intended target.
25. As a reviewer, I want productive repair examples, so that the report examines plausible false positives.
26. As a reviewer, I want unchanged rerun examples, so that unchanged content has its own explanation.
27. As a reviewer, I want overlapping assignment examples, so that attribution errors remain visible.
28. As a reviewer, I want missing classes to yield an inconclusive report, so that absent evidence cannot support adoption.
29. As a reviewer, I want a reproducible evidence report, so that another session can inspect its basis.
30. As a reviewer, I want explicit operational instructions, so that the real pilot produces evidence after activation.
31. As a reviewer, I want collection independent of the oracle, so that observations cannot change checks or model routing.

## Implementation decisions

### Command and owner

Add `bench repair-pilot activate`, `bench repair-pilot record --input <file>`, and `bench repair-pilot report [--full]`.
These operational commands use the production command registry and the existing wrapper fallback.
They use `usage.Grammar`, compact TOON, and exits 0 for success, 1 for refusal, and 2 for usage.
Help accepts `help`, `--help`, and `-h` without requiring a repository.
RP1, RP2, RP5, RP6, and RP45 cover this route through fixed command fixtures.

One deep module, `internal/repairpilot`, owns the pilot document, transitions, and report.
The command adapter supplies the current time, canonical repository key, Bench home, and kit-source disposition.
Reuse `poolkey.Key`, `gate.KitSourceCheckout`, `bounds.ClassifyNoFollow`, `jsonfile`, `usage`, and `toon`.
Compare the canonical repository roots when resolving kit identity across Bench worktrees.
Reuse the assessment reference type where useful, without changing assessment records or their consumers.

Keep command registration in the existing command files.
Move `boundaryRoot` beside the registry's root-resolution helpers to create headroom in `main.go`.
Do not add a file to the crowded command or assessment directories.
Keep the new package within the current structure budgets.
RP46 covers route integration and the unchanged assessment command family.

### Local document and activation

One versioned document lives below the Bench home at `repair-pilot/<repo-key>/pilot.json`.
The module derives that path once, outside the disposable worktree pool.
It uses a private directory, mode 0700, and private files, mode 0600.
It retains observations until the operator explicitly removes the named document outside this feature.
RP3 and RP4 cover the local custody contract.

The document records its version, repository key, activation time, observations, and audits.
Activation creates the document with the current UTC time.
An identical repeat returns the existing activation without changing it.
A stopped pilot cannot restart through `activate`.
RP7 and RP8 cover repeat activation and the terminal state.

Each mutation locks the repository's pilot document, reads it, validates the proposed append, and replaces it atomically.
An occupied lock refuses immediately and names a retry after the active writer completes.
Write and replacement failures retain the previous document.
A stale lock requires an operator decision after inspection, not automatic deletion.
RP9 through RP12 cover these failure positions.

No background worker, hook, timer, or ambient dashboard signal participates.
The next command evaluates elapsed time, but the cutoff remains the original deadline even if no command runs then.
The command never launches checks, models, evidence references, or network requests.
RP44 and the review-owned RP47 protect this separation.

### Sequence and observation contract

Use one sequence key composed from the implementation source identity, spec identity, and stable chunk ID.
A source identity names the retained integration source across ticket commits and process changes.
Assignment IDs and session IDs identify contributors, not new sequences.
Several contributors can belong to one sequence when the chunk retains that identity.
RP13 through RP17 exercise these distinctions.

The first blocking observation starts a sequence.
A record retains the observed time, observation ID, sequence key, assignment ID, session ID, and source revision or tree.
It also retains the work stage, observation kind, known failures, and native evidence references.

The supported work stages are pre-review and post-review; RP17 exercises each stage.
The supported observation kinds are repair and rerun.
A repair observation names its hypothesis, intended change, and verification reference.
A rerun records verification without claiming an implementation attempt.
RP17, RP18, and RP22 cover these producer-derived shapes.

A failure names its check, defect or requirement identity when known, diagnostic, and ownership.
Ownership is diff-owned, inherited, spec-predicted, or unknown.
A producer's first failure retains its order and native reference.
The record states whether the failure set is complete, first-only, or unknown.

Unknown identity never compares equal merely because two diagnostics are generic.
The collector accepts explicitly selected evidence and never assigns observations through timestamp containment.
RP19 through RP21 and RP48 cover the failure contract.

A sequence ends through verified closure of all its blockers or an evidenced reviewer handoff.
Both endpoints count as completed sequences, but handoff never claims successful repair.
A closure needs references that account for every recorded blocker.
A restarted review cannot reopen a completed sequence as a different sequence for the same chunk.
RP23 through RP25 cover those endpoints.

Imports append identified observations or audits.
The same ID with identical content returns success without a second append.
The same ID with different content refuses without changing stored evidence.
RP26 and RP27 cover retry behavior.

### Collection cutoff

The deadline is activation plus fourteen calendar days in UTC.
The cutoff is the earlier of that deadline and the accepted endpoint of the tenth completed sequence.
The command accepts the observation that completes sequence ten, then closes collection.
At the deadline itself, the command accepts no new observations.
Open sequences at the cutoff remain incomplete, including when the count limit stops collection.
RP28 through RP31 cover both bounds and their intersection.

Records must describe work within the active collection window and cannot carry future observed times.
After the cutoff, a record can append audits of existing observations but cannot append observations, even with earlier claimed times.
This rule prevents a delayed import from changing the frozen sample.
The report retains known collection gaps rather than treating missing observations as zero attempts.
RP32 and RP33 cover this distinction.

### Progress and example audit

The author can propose a progress label, but a proposal is not an accepted label.
An audit names the observations examined, the auditor, the evidence references, and the reason for its conclusion.
A progress conclusion identifies a repaired requirement or defect and its successful verification.
Changed files, changed trees, green check outcomes, and unsupported author claims are insufficient on their own.
Missing or conflicting support leaves the effective label unknown.
RP34 through RP37 cover the structural evidence contract.

The software validates references and the audit structure, not the truth of arbitrary prose.
A separate reviewer or read-only review session examines the native evidence before supplying a supported audit.
RP49 is review-owned and checks that the operating protocol requires that examination during the real pilot.
The report distinguishes proposed labels, audited labels, and unresolved evidence.
An audit never overwrites an earlier proposal or audit.
Conflicting supported audits require an explicit resolution audit that cites both earlier audits.

The required example classes describe observations or intervals, not mutually exclusive sequence buckets.
A sequence can contain both productive and stalled repair intervals.

A stalled example needs an actual repair attempt and audited evidence that its target blockers did not improve.
A productive example needs an audited, verified repair even when the first failure repeats.
An unchanged rerun needs identical source content without a repair attempt.
An overlapping example needs explicit distinct assignment identities with intersecting observed intervals.
Neither overlap nor a repeated first failure determines progress.
RP38 through RP41 cover each class separately.

### Report and operating protocol

The default report gives the pilot state, cutoff, completed and incomplete counts, and required-class coverage.
It also reports the count of unknown labels and unresolved evidence gaps.
The full report lists the retained observations, proposals, audits, citations, and comparison intervals in stable order.
Both forms derive from the same validated document.
RP42 and RP43 cover complete reporting over canned observations.

Before cutoff, the report marks the evidence provisional.
At cutoff, any missing required example class makes the result inconclusive.
Complete class coverage means only that the sample covers the required cases.
It grants no detector threshold, accuracy claim, warning, model change, or gate override.
RP42 and RP47 cover these limits.

Add a repository-only operating guide at `docs/repair-collection-pilot.md`.
It explains activation, collection at verification points, evidence audit, cutoff, report export, and the reviewer decision.
It names the exact commands and gives one synthetic input marked as synthetic.
Synthetic fixtures never count as real pilot examples.
RP50 covers this guide and the later evidence report at `capture/reports/repair-collection-pilot.md`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RP-C1 / 1-activate-pilot.md | Explicit kit-only activation and durable storage | RP1-RP12, RP45-RP46 | TestRepairPilotActivation, TestRepairPilotStorage, TestRepairPilotGrammar, TestRepairPilotRoute | no |
| RP-C2 / 2-collect-repair-evidence.md | Attributed bounded repair evidence with audit inputs | RP13-RP37, RP48 | TestRepairPilotIdentity, TestRepairPilotFailures, TestRepairPilotEndpoints, TestRepairPilotCutoff, TestRepairPilotAudit | yes |
| RP-C3 / 3-report-pilot-evidence.md | Evidence report and real-pilot operating protocol | RP38-RP44, RP47, RP49-RP50 | TestRepairPilotReport, TestRepairPilotIsolation, protocol review | no |

Each ticket is one serial commit checkpoint on the retained integration source.
After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review.
Close its blocking findings before starting the next chunk.
After RP-C3, reconcile all acceptance rows and integration evidence before landing.
No ticket runs the real fourteen-day experiment as an implementation prerequisite.

```bench-completion-plan
{"version":1,"chunks":[{"id":"RP-C1","tickets":["1-activate-pilot.md"],"verification":[{"id":"activation-storage","command":"bench test --package ./internal/repairpilot","probe":"omit preservation on failed replacement"},{"id":"dispatch","command":"bench test --package ./cmd/bench --run RepairPilot"}]},{"id":"RP-C2","tickets":["2-collect-repair-evidence.md"],"verification":[{"id":"collection","command":"bench test --package ./internal/repairpilot","probe":"swap the inclusive deadline comparison"}]},{"id":"RP-C3","tickets":["3-report-pilot-evidence.md"],"verification":[{"id":"report","command":"bench test --package ./internal/repairpilot","probe":"omit one required example class from completeness"},{"id":"protocol","command":"review RP47, RP49, and RP50 against the operating guide and final integration diff"}]}],"final_verification":[{"id":"acceptance","command":"bench test --package ./internal/repairpilot"},{"id":"integration","command":"bench test --package ./cmd/bench --run RepairPilot"}]}
```

## Testing decisions

Use TDD at the production `repairpilot.Command` seam, with fixed time and real temporary files.
The tests feed JSON through the real decoder, transition owner, writer, and renderer.
Only time and filesystem failure injection substitute external dependencies.
No test mocks the internal identity, cutoff, or audit rules.

Prior art: `assessment.Command`, `Store.Record`, `TestAssessmentCommandUnsafe`, and `TestAssessmentRecordUpdates`
The new module keeps the pilot domain separate from workflow cost accounting.
Command dispatch tests call `cmd/bench.Command.Run` through the existing fixture harness.
The gate's ordinary `go test -trimpath -count=1 ./...` phase executes the new tests.
No new gate phase or conformance check is necessary.

All test names in the coverage map are planned additions.
The fixture inventory lives at `internal/repairpilot/testdata/cases.json`.
Each coverage row names a fixture case by its suffix after the test name.
A fixture can contribute to several cases without copied harness code.

Use canned times and references when comparing reports.
Each independent expectation needs a recorded omission or swap that makes its named test red.
A probe that does not compile or execute proves nothing.

### Seam diagram

    operator: bench repair-pilot <operation>
                         |
                         v
    explicit JSON --> [ Command: document transitions ] --> local pilot document
                         |                                  |
                         +---------- report <---------------+
                         |
                         v
                  TOON evidence report
    tests: fixed time, real files, selected write failures, same production entry

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RP1 | 1 | A missing pilot report returns inactive without creating a document. | TestRepairPilotActivation/inactive | Automatic activation creates forbidden state. |
| RP2 | 1 | Explicit activation creates a pilot at the supplied current time. | TestRepairPilotActivation/activate | A no-op has no stored activation. |
| RP3 | 3 | Two worktree contexts share the canonical repository's pilot document. | TestRepairPilotStorage/worktrees | Worktree-local storage produces two pilots. |
| RP4 | 3 | The stored pilot uses private permissions outside the disposable pool. | TestRepairPilotStorage/custody | A pool-local record fails the path assertion. |
| RP5 | 2 | A linked repository refuses activation without writing state. | TestRepairPilotActivation/linked | Unrestricted activation creates a foreign pilot. |
| RP6 | 5 | A malformed operation exits 2 before reading pilot state. | TestRepairPilotGrammar/usage | Eager storage access masks the grammar refusal. |
| RP7 | 4 | Repeat activation preserves the original timestamp and observations. | TestRepairPilotActivation/repeat | A reset changes the frozen activation. |
| RP8 | 4 | Activation after cutoff preserves the stopped pilot. | TestRepairPilotActivation/stopped | A new observation window changes the cutoff. |
| RP9 | 6 | An occupied pilot lock refuses without changing the document. | TestRepairPilotStorage/locked | A second writer overwrites the first writer's evidence. |
| RP10 | 6 | Failure to create a temporary document retains the prior bytes. | TestRepairPilotStorage/create-failure | Direct overwrite destroys the prior document. |
| RP11 | 6 | A partial temporary write retains the prior bytes. | TestRepairPilotStorage/write-failure | Premature publication exposes partial JSON. |
| RP12 | 6 | Replacement failure retains the prior bytes and removes the temporary file. | TestRepairPilotStorage/replace-failure | Failed publication leaves changed or orphan state. |
| RP13 | 7 | A fresh session for the same source and chunk retains one sequence. | TestRepairPilotIdentity/session | Session-based grouping doubles the sequence. |
| RP14 | 8 | Overlapping assignments keep their explicitly supplied observation identities. | TestRepairPilotIdentity/assignment | Time-based association assigns an observation to the wrong work. |
| RP15 | 9 | Two findings in one chunk share one sequence. | TestRepairPilotIdentity/findings | Finding-based grouping inflates the count. |
| RP16 | 7 | The same chunk ID in different implementation sources creates distinct sequences. | TestRepairPilotIdentity/sources | Chunk-only grouping conflates unrelated work. |
| RP17 | 10 | Pre-review and post-review observations remain in one sequence. | TestRepairPilotIdentity/stages | Review-stage grouping resets the sequence. |
| RP18 | 8 | Missing assignment or source evidence refuses an attributed observation. | TestRepairPilotIdentity/missing | Unattributed observations become apparently comparable evidence. |
| RP19 | 11 | A repeated first failure retains its check, identity, diagnostic, and reference. | TestRepairPilotFailures/repeated | Diagnostic-only storage loses the comparable defect identity. |
| RP20 | 12 | A first-only failure set remains first-only in the report. | TestRepairPilotFailures/partial | A first red becomes a complete-set assertion. |
| RP21 | 11 | Two generic diagnostics with unknown defect identity remain incomparable. | TestRepairPilotFailures/generic | Text equality fabricates a repeated defect. |
| RP22 | 13 | A rerun on unchanged content does not count as a repair attempt. | TestRepairPilotFailures/rerun | Every verification becomes a repair. |
| RP23 | 14 | Closure without verification for one recorded blocker refuses. | TestRepairPilotEndpoints/unverified | Partial closure silently completes the sequence. |
| RP24 | 14 | Verification covering all recorded blockers completes the sequence successfully. | TestRepairPilotEndpoints/closed | Always-open behavior never counts verified completion. |
| RP25 | 15 | An evidenced reviewer handoff completes the sequence with unresolved blockers visible. | TestRepairPilotEndpoints/handoff | A handoff erases the unresolved outcome. |
| RP26 | 16 | Reimport of an identical observation leaves one observation. | TestRepairPilotIdentity/retry | An append-only retry duplicates the record. |
| RP27 | 16 | Conflicting content for an existing ID refuses without changing evidence. | TestRepairPilotIdentity/conflict | Last-write-wins hides the conflict. |
| RP28 | 17 | The tenth accepted sequence endpoint closes collection. | TestRepairPilotCutoff/count | An eleventh endpoint can enlarge the sample. |
| RP29 | 18 | The deadline itself refuses a new observation. | TestRepairPilotCutoff/deadline | A greater-than comparison accepts a boundary observation. |
| RP30 | 17, 18 | Collection uses the earlier endpoint when count and time compete. | TestRepairPilotCutoff/earlier | A later limit extends the window. |
| RP31 | 19 | A sequence open at cutoff remains incomplete. | TestRepairPilotCutoff/incomplete | Deadline handling manufactures a successful endpoint. |
| RP32 | 18 | A stopped pilot refuses a backdated new observation. | TestRepairPilotCutoff/late-import | Backdating changes the frozen sample. |
| RP33 | 23 | A stopped pilot accepts an audit of an existing observation. | TestRepairPilotCutoff/audit | A blanket freeze prevents evidence assessment. |
| RP34 | 20 | Progress without a verified repaired requirement or defect remains unknown. | TestRepairPilotAudit/no-repair | A green check alone establishes progress. |
| RP35 | 21 | An author proposal without an evidence audit remains unknown. | TestRepairPilotAudit/proposal | Self-labeling becomes an accepted conclusion. |
| RP36 | 21 | A supported audit exposes its repaired target and verification references. | TestRepairPilotAudit/supported | A label-only report hides its proof. |
| RP37 | 22 | Conflicting audits retain unknown status until an explicit resolution cites both. | TestRepairPilotAudit/conflicting | The latest label silently wins. |
| RP38 | 24 | A stalled class requires an audited repair attempt with no verified improvement. | TestRepairPilotReport/stalled | Repetition alone supplies a stalled example. |
| RP39 | 25 | A verified repair supplies a productive example despite its repeated first failure. | TestRepairPilotReport/productive | First-red equality overrides verified progress. |
| RP40 | 26 | An unchanged rerun supplies only its evidenced rerun classification. | TestRepairPilotReport/unchanged | A rerun is mislabeled as a stalled repair. |
| RP41 | 27 | Explicitly attributed intersecting assignment intervals supply an overlap example. | TestRepairPilotReport/overlap | Timestamp-only grouping passes a falsely attributed example. |
| RP42 | 28 | Any missing required class makes the terminal report inconclusive. | TestRepairPilotReport/missing-class | Three classes pass a four-class requirement. |
| RP43 | 29 | The full report exposes every retained observation, audit, unknown label, and evidence gap. | TestRepairPilotReport/full | A missing record passes a summary-only assertion. |
| RP44 | 31 | Reporting reads the pilot without changing its stored bytes. | TestRepairPilotIsolation/read-only | A report mutates its own evidence. |
| RP45 | 5 | Invalid record files refuse before state mutation. | TestRepairPilotGrammar/hostile | A permissive decoder admits the hostile fixture inventory. |
| RP46 | 2, 30 | The public dispatcher reaches the pilot owner while assessment routes retain their behavior. | TestRepairPilotRoute/dispatch | A test-only route or hijacked assessment route fails. |
| RP47 | 31 | Pilot integration introduces no gate, warning, or model-routing effect. | Review-owned integration inspection | An automatic hook or routing write violates the declared fence. |
| RP48 | 12 | Failure ownership and completeness retain each producer-supplied vocabulary value. | TestRepairPilotFailures/vocabulary | A default coercion turns unknown ownership into diff-owned. |
| RP49 | 20, 21 | The operating protocol requires a separate native-evidence audit for each accepted progress label. | Review-owned protocol walkthrough | Structural validation alone cannot establish semantic truth. |
| RP50 | 30 | The operating guide gives activation, collection, audit, cutoff, and report steps with exact commands. | Review-owned protocol walkthrough | A command inventory alone leaves the pilot without an operator procedure. |

### Edge inventory

Audience: the Bench kit repository only

The hostile record inventory for RP45 includes absent input, empty input, and a final line without a newline.
It includes malformed JSON, duplicate keys, unknown fields, an unsupported version, and a foreign repository key.
It includes a FIFO, a live symlink, a dangling symlink, a linked parent, and an oversized document.
It includes unsafe IDs, ESC, BEL, tabs, newlines, and numeric-looking identifiers.
Absent input refuses; an absent pilot reports inactive; an empty stored document refuses.
Human input may omit a final newline, while the machine document requires one.

Input paths with spaces and glob characters reach the same decoder without shell expansion.
IDs use a bounded ASCII token grammar, so Unicode whitespace cannot create another identity.
Evidence references remain opaque strings with a declared producer and native location.
The command does not read, execute, trim, or resolve evidence references.
A multiline evidence location follows the existing reference contract and TOON escaping.
Hostile fixtures exercise the newly introduced command, not only the shared parser.

RP9-RP12 cover concurrent update refusal and interrupted persistence through injected filesystem failures.
RP26-RP33 cover retries, partial state, exact limits, delayed collection, and audit-only continuation.
RP34-RP43 cover unsupported, contradictory, and missing evidence.
Tests do not swap a package variable.
The adapter supplies time as an explicit value; all report comparisons use the same canned time.


Won't handle: remote evidence fetch — the operator's evidence audit reads selected native sources separately.

Won't handle: a background deadline notification — every pilot command applies the same fixed cutoff when called.

Won't handle: automated semantic verification — the evidence auditor owns the truth of the progress conclusion.

Won't handle: shell or patch-header parsing — the record command consumes structured JSON through the existing argument parser.

Won't handle: destructive worktree operations — collection uses canonical identity and never edits or releases a worktree.

Won't handle: hostile writers that replace files outside the command's lock protocol — the local operator owns the Bench home.

Won't handle: unrequested hook or adapter entry points — the operator invokes the existing CLI wrapper explicitly.

## Ownership fences

Reviewer disposition: proposed for the spec-and-ticket sign-off

- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `internal/repairpilot/` (new)
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `DATA_HANDLING.md`
- `CONTEXT.md`
- `docs/repair-collection-pilot.md` (new)
- `capture/reports/repair-collection-pilot.md` (new)
- `specs/repair-collection-pilot/`
- `reviews/repair-collection-pilot.md` (new)

The command and AXI registry fences include required co-named consumers.
The approved AXI query set stays unchanged because the pilot has its own operational output contract.
Co-naming those registries does not require edits to unchanged facts.
The report artifact fence supports the later real-pilot report, not fabricated implementation evidence.

## Out of scope

The existing bounded repair policy remains the authority for repair allowance and blocker classification.
A detector, automatic warning, model routing, and gate override remain excluded under decision ticket #14.

Detector estimate: 3 edits, 1 gate run for a later independent capability, subject to the pilot's evidence and reviewer decision.
Those edits represent a signal owner, its tests, and an operator-facing integration.
Warning estimate: 2 edits, 1 gate run for presentation and its tests after a detector decision.
Linked-repository collection estimate: 3 edits, 1 gate run for scope policy, compatibility coverage, and linked operator guidance.
These estimates do not authorize those capabilities.

## Further notes

### Source verification and reader sweep

Source revision: a23e5b20be0fdf3cf4d2e9a4d37f261e05f7b9e3

The author reread every structured source from the compiled map.
The roadmap remains an experiment, and its no-signature premise is stale for lane spans.
`beginLaneSpan` retains the first check and diagnostic; `writeLaneRecord` still omits them.
`beginCommitSpan` discards its tracing context, and phase spans omit assertion identities and failure ownership.
The record processor drops write failures, so missing spans remain unknown evidence.

The existing report's historical snapshot is evidence at its pinned revision, not a current completeness claim.
The author did not reread the 101609090-byte historical trace snapshot.

The reader sweep covered hidden files, excluding `.git`, including scripts, workflows, fixtures, and conformance checks.
It searched the moved source path, lane record writer, declared span attributes, and assessment record fields.
No lane record, OTEL attribute, assessment schema, or assessment selector changes in this design.
Their readers remain excluded from writes: gate record tests, OTEL encode tests, system OTEL tests, and assessment consumers.
The pilot stores explicit references to those records without parsing an inferred assignment from them.

The moved map's index, tickets, asset, and bounded-policy retrospective receive their new paths in this spec phase.
The map remains one source, and FT232 remains on the board.
No `Roadmap:` retirement binding attaches this spec to FT232 because the later detector decision remains open.

The command reader inventory is `commandRegistry`, `Command.Run`, `renderCommandHelp`, and `TestHelpInventoryIsComplete`.
Its enforcement includes `subcommandRouting`, `approvedAXIQueries`, and the ticket registry's `commandRegistries`.
`boundaryRoot` callers are the assessment route, `leafRoot`, and `resumeCleanCommand`.
Its relocation changes no API or result.
Spec-phase map moves precede the build baseline and do not enlarge the implementation fence.
The wrapper fallback already forwards unknown names to the production binary.

The author read `internal/coverage/coverage.go`, `internal/tickets/tickets.go`, and `internal/tickets/registry_data.go` for artifact enforcement.
The author read the command registry, wrapper route, help inventory, AXI registry, and routing table before choosing fences.
The ordinary-build census and conformance registry establish the test execution route.
The new package needs no separate live-tree conformance registration.
The co-named canary fixtures pin existing command, data-handling, and glossary facts.
Their expected rules stay unchanged.

Shipped-surface claim words: none for repository-only documentation
The compiled command is available in the binary, but its activation scope excludes linked repositories.

### Source-sentence-to-row table

| source clause | disposition |
| --- | --- |
| Ticket #9: "The reviewer authorizes a bounded collection experiment." | RP1-RP12 |
| Ticket #9: "It links repair attempts to assignments and records comparable failures with verified progress evidence." | RP13-RP27, RP34-RP37, RP48-RP49 |
| Ticket #9: "After the experiment, the reviewer reevaluates whether a detector is justified." | RP47, RP50 |
| Ticket #10: "The pilot requires explicit activation and covers only the Bench kit repository." | RP1, RP2, RP5 |
| Ticket #10: "It observes ordinary implementation and repair work." | RP13-RP25, RP50 |
| Ticket #10: "It changes no collection default for linked repositories." | RP5, RP46 |
| Ticket #11: "A progress label requires a repaired requirement or defect with verification attached." | RP34, RP36, RP49 |
| Ticket #11: "The agent may propose the label, but the experiment audits that label against its evidence." | RP35-RP37, RP49 |
| Ticket #11: "Check outcomes alone, changed files, and unsupported author claims do not establish progress." | RP34-RP35 |
| Ticket #11: "Unsupported labels remain unknown." | RP34-RP37, RP43 |
| Ticket #12: "Stop the pilot at ten completed repair sequences or fourteen calendar days after activation, whichever occurs first." | RP28-RP33 |
| Ticket #12: "The report must cover stalled repairs, productive repairs, unchanged reruns, and overlapping assignments." | RP38-RP41 |
| Ticket #12: "If a required example class is absent, report that gap and an inconclusive result." | RP42-RP43 |
| Ticket #13: "One repair sequence covers one implementation chunk." | RP13-RP17 |
| Ticket #13: "It starts at that chunk's first blocking failure." | RP13, RP18 |
| Ticket #13: "It ends when verification closes all its blockers or work stops for reviewer handoff." | RP23-RP25 |
| Ticket #13: "The sequence can include several findings and both pre-review and post-review repairs." | RP15, RP17 |
| Ticket #13: "A fresh process or review does not create a new sequence." | RP13, RP17 |
| Ticket #13: "Report sequences still open at the deadline as incomplete." | RP31 |
| Ticket #14: "The second implements the explicitly enabled Bench collection pilot and its evidence report." | RP1-RP50 |
| Ticket #14: "A detector and any warning remain outside both specifications." | RP47 and Out of scope |

Policy-only decisions #1-#4, #7, and #8 remain consumed by the landed policy.
Tickets #5 and #6 explain the need for real examples and the limits of the old records.
They authorize no additional detector behavior.

### Flagged additions

- A local operational command and versioned document implement the collection capability: RP1-RP12, RP45-RP46.
- UTC defines the calendar-day calculation: RP29-RP30.
- Frozen observations allow later audits but refuse delayed imports: RP32-RP33.
- Structured audit references support a separate semantic audit: RP34-RP37, RP49.
- Idempotent imports and atomic replacement protect evidence custody: RP9-RP12, RP26-RP27.
- Overlapping example classes apply to intervals within a sequence: RP38-RP41.

These are engineering choices proposed for this sign-off, not additional detector policy.

### Pre-review proof checklist

- Cited symbols: `assessment.Command`, `Store.Record`, `poolkey.Key`, `poolkey.Canonical`, `gate.KitSourceCheckout`, `bounds.ClassifyNoFollow`, `jsonfile.DecodeDocument`, `usage.Grammar`, and `toon.Table` resolve in the tree.
- Import edges: the new module consumes shared bounds, JSON, identity, and rendering owners without changing their APIs.
- Source-row clauses and occurrences: the table above quotes the applicable clauses from tickets #9-#14, once per canonical answer.
- Promised field labels: version, repository key, activation time, observation ID, sequence key, assignment ID, session ID, source identity, stage, kind, failures, references, proposals, and audits.
- Changed-function callers: only the command registry gains a route; the three existing `boundaryRoot` callers keep the same function.
- Copy survival: no production owner replaces copied behavior; the compiled decision map moves as one unit.

RP45's hostile inventory must run through the new command before its implementation ticket closes.
RP12 has an explicit temporary-file omission case, so cleanup cannot pass by checking only the final document.
No coverage row cites a test-only helper across a package seam.
Review-owned rows state their semantic or integration limit instead of claiming an automatic truth oracle.
