# Native record operations

Status: staged

Related roadmap: FT318 and the result-projection face of FT168.

Decision source: the named reviewed ready map `decisions/architecture-planning.md`, answers 3 and 5.

Verification log: 1 iteration(s) to accept — the independent ticket review accepted all three chunks without repair.
Independent SPEC confirmation closed NR-R1 and NR-R2 at e7c03a6a0927ace66e4af5f9f89b53c04e47e564 before slicing.
The preserved first draft is response artifact 1791348252437980913-46812d8caa9e1acc.out in the primary Bench response store.

## Problem

The existing record forms write completion evidence, but dispatch history still needs hand-built JSON.
A changed plan can reach another chunk write before its amendment is recorded.
Probe evidence also makes authors copy fields with different exit meanings.

## Solution

Add a native assignment form over a complete existing version-2 completion plan.
Reject invalid candidates and stale inputs before writing.
Name the applicable record form when an entry is missing, and expose the required completion order.
Import a current probe result through its producer owner while retaining exact native provenance.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: the existing validators and fixtures define the contracts, while complete-plan writes and evidence ordering require careful composition.
Harder chunks: none. The requested gpt-6.1-sol / high spec author line is separate from this implementation recommendation.

### Dispatch

1. As an orchestrator, I want to record the first dispatch in an existing version-2 plan, so that an empty history becomes usable.
2. As an orchestrator, I want explicit execution fields checked against the source plan, so that a dispatch cannot change the run identity.
3. As an orchestrator, I want the complete candidate plan validated, so that an unrelated invalid ticket cannot survive a valid dispatch.
4. As an orchestrator, I want stale source writes refused, so that a dispatch cannot overwrite a newer plan.
5. As an orchestrator, I want replacement history appended, so that the preceding author remains visible to independence checks.
6. As an orchestrator, I want replacement evidence validated, so that a returned restore function cannot stand for stopped writer evidence.
7. As an orchestrator, I want duplicate or cross-ticket author sessions refused, so that each author keeps one ticket.
8. As an orchestrator, I want only the plan payload replaced, so that surrounding spec prose stays exact.
9. As an orchestrator, I want failed assignment writes to retain the original file, so that retry starts from the same source.
10. As an orchestrator, I want assignment forms in native help, so that dispatch needs no general JSON editor.
### Ordering

11. As an orchestrator, I want a changed plan amended before recording another chunk, so that recorded history uses one digest chain.
12. As an orchestrator, I want prose changes to retain their digest consequence, so that an unapproved semantic classifier grants no exemption.
13. As an orchestrator, I want a missing chunk refusal to name the chunk form, so that I can repair the missing entry.
14. As an author, I want a missing verification refusal to name the verification form, so that I can retain the required result.
15. As a reviewer, I want a missing review refusal to name the review form, so that I can retain the required axis result.
16. As an orchestrator, I want missing completion evidence to name its form, so that reconciliation has a direct repair route.
17. As an orchestrator, I want undispatched ownership to name the assignment form, so that verification uses an actual author.
18. As an orchestrator, I want completion help to show the required order, so that I do not discover it through repeated refusals.
### Probe

19. As an author, I want the normalized focused-execution exit retained by the probe, so that a successful mutation command cannot imply passing tests.
20. As an author, I want a probe result imported through its owner, so that I do not copy three outcome fields by hand.
21. As an author, I want import provenance bound to exact native bytes, so that the retained excerpt identifies the imported result.
22. As an author, I want unsupported or malformed result envelopes refused, so that an old result cannot supply a guessed exit.
23. As an author, I want non-biting and unrestored results retained honestly, so that import cannot turn them into completion evidence.
24. As an author, I want scalar result flags to remain available, so that existing explicit evidence calls keep their behavior.
25. As an author, I want copied flag help to name each producer field, so that command exit and focused-execution exit cannot be confused.
### Preservation

26. As a caller, I want existing record error identities preserved, so that callers can still inspect the same refusal condition.
27. As a coordinator, I want existing chunk and amendment history preserved, so that the new forms cannot rewrite prior evidence.
28. As a coordinator, I want capability-blocked completion kept outside this change, so that FT317 retains its decision authority.
29. As an orchestrator, I want version-1 conversion kept explicit and separate, so that ticket verification mapping remains an authoring obligation.
30. As a reviewer, I want independent entry expectations to fail on an omitted writer guard, so that owner tests cannot hide an adapter bypass.

## Implementation decisions

### Source pin and current owners

Research uses production `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68` and planning source `49a28e8124b5d919b4b39d4c346af962e1c1fa26`.
These are source observations, not execution evidence.
`ReadPlan` in `internal/reviewrecord/plan.go:35` validates the full plan and reads every committed ticket.
It validates execution, requirement ownership, dependency order, cycles, and acceptance-row derivation.
Its digest includes complete spec bytes and ticket bytes.

`validateExecution` and `validateAssignment` in `internal/reviewrecord/delegated.go` own assignment history.
Empty histories already represent undispatched version-2 tickets.
The first assignment has no replacement fields.
Later assignments name the preceding session, an accepted trigger, stopped writer evidence, and preserved work.
Only a no-progress replacement requires reassessment.
The schema accepts explicit nonempty model and effort values, including `unknown`.

`write` and `replace` in `internal/reviewrecord/write.go` own validation before atomic record replacement.
Reuse `replace` for the assignment file effect.
Use the same named-fence span owner that `fenced` and `Render` use.
Keep JSON strictness, input size, safe path, source digest, and evidence rules in their existing owners.

The source sweep covers qualified and unqualified calls to `ReadPlan`, `RecordChunk`, `RecordVerification`, `CheckTrees`, and the repair renderer.
Production consumers are `reviewrecord.atSource`, `reviewrecord.checkSource`, `recordcmd`, `preflight.withCompletionPlan`, and `gate.applyCheckpoint`.
Fixture consumers include both `recordtest` loaders and the existing delegated worktree journey.
Public adapters retain their signatures unless the accepted Markdown delivery changes their internal composition.
No `.mjs` or workflow consumer of these Go entry points was found.

### Canonical Markdown prerequisite

NR-C1 uses the accepted delivery of `specs/markdown-block-reader/spec.md`, especially MB-C4's review-record reader and replacement spans.
That delivery keeps `ReadPlan`, `fenced`, and record rendering on one block owner.
The approved reader returns classified physical lines, byte offsets, and at most one fault.
Native assignment consumes that contract through reviewrecord and adds no fence detector.
This is an actual owner prerequisite, independent of the advisory C10 order.
The later chunks inherit it through NR-C1's reviewed interface.

### Assignment contract

The prospective domain entry is `RecordAssignment(root, spec, AssignmentCall)` in `internal/reviewrecord`.
`AssignmentCall` carries the source selector, ticket basename, one `Assignment`, and explicit execution assertions.
Keep `ReadPlan(root, tree, spec)` as the committed-tree adapter.
Extract its complete validation into one internal operation that accepts spec bytes and the canonical ticket reader.
Both committed reads and prospective candidate validation call that operation.

The form is:

```text
bench record assignment <slug> --source <commit> --ticket <basename> --performer <session> --assignment <id> --model <model> --effort <effort> --ref <native-dispatch-ref> --run <run-id> --orchestrator <session> --author-limit <positive> [--review-mode unified] [--predecessor <session> --trigger <trigger> --stopped <ref> --preserved <ref> [--reassessment <ref>]]
```

Every execution assertion must equal the source plan.
Omitting `--review-mode` asserts the plan's existing default review mode.
No flag initializes or changes execution identity.
Resolve `--source` to a full commit and require that commit to remain HEAD for this assignment write.
Require the working spec and each planned ticket to equal their committed source bytes.
Check those preimages again before calling the replacement owner.

The orchestrator remains the serialized plan writer.

Append a first assignment only to an empty target history.
A nonempty history requires a validated replacement whose predecessor is its last session.
Store the resolved source commit in the new assignment's source field.
Keep the explicit native dispatch reference as metadata, never as proof of a completed verification.
Validate the complete rendered candidate, including all other histories and tickets, before creating a temporary file.
Replace only the plan payload and preserve surrounding bytes.

The operation neither dispatches an agent nor commits the spec.
It does not write the completion record or automatically amend it.
After committing changed plan bytes, an existing record needs its normal amendment before another chunk entry.
A version-1 plan refuses this form.
Mapping every existing verification obligation to its owning ticket is a separate authoring obligation, not an inferred conversion.

A real first-dispatch example starts from this complete version-2 plan and its valid committed `1.md` ticket:

```bench-example
{"version":2,"execution":{"mode":"delegate","run_id":"native-demo","orchestrator_session":"native-coordinator","author_limit":2,"assignments":{"1.md":[]}},"chunks":[{"id":"C1","tickets":["1.md"],"verification":[{"id":"tests","ticket":"1.md","command":"bench test --package ./internal/example"}]}],"final_verification":[{"id":"final","command":"bench test --package ./internal/example"}]}
```

```text
bench record assignment example --source HEAD --ticket 1.md --performer native-author-1 --assignment native-dispatch-1 --model gpt-5.6-sol --effort high --ref native:dispatch-1 --run native-demo --orchestrator native-coordinator --author-limit 2
```

The result contains one first-dispatch occurrence for `1.md` and the original execution identity.
A replacement adds the predecessor, trigger, stopped, and preserved flags required by the existing history validator.
No-progress also adds reassessment.

### Ordering and repair ownership

`RecordChunk` at `internal/reviewrecord/write.go:125` currently omits the record-versus-current plan digest guard.
`checkSource` at `internal/reviewrecord/coverage.go:20` already detects the stale digest later.
Share that condition so the chunk writer refuses before appending or updating evidence.
Retain the missing record's initial chunk creation behavior.
A valid amendment remains the only writer that advances an existing record's plan digest.

Add typed repair information at the existing reviewrecord refusal sites.
Preserve original errors through wrapping and `errors.Is`.
Do not parse human error strings to classify missing obligations.

Reviewrecord owns typed repair form names and known context values.
The existing recordcmd registry uses those names for its form declarations and derives their complete help grammar.
Gate emits the corresponding `bench record <form> --help` route from the typed name.
It does not copy flags or import recordcmd.

The import graph remains acyclic: recordcmd uses probe, probe uses gate, and gate uses reviewrecord.
`gate.routedRefusal` retains `bench preflight review <slug>` for structural or unclassified refusals.
The gate does not become an evidence writer.

The repair families are assignment, amendment, chunk, verification, review, and completion.
Name unresolved future inputs as placeholders.
Use already-known slug, chunk, requirement, axis, and source values where the owner supplies them.
A refused occurrence that needs new evidence must not advertise a fabricated completed result.

Completion help derives its ordering note from this owner contract.
Amend a changed plan, record the chunk, and retain its verification and reviews.
Then retain final verification and write completion.

### Probe result import

`probe.probe` at `internal/probe/probe.go:163` currently discards `testreport.Execute`'s third return.
That return is the normalized focused-execution exit, not a raw child-process status.
`runGoTest` in `internal/testreport/command.go` normalizes a nonzero completed run to exit 1.
A biting probe command returns 0 after observing that failing execution and proving restoration.
The normalized focused-execution exit and probe command exit must remain distinct.

Preserve the current leading probe and selection table schemas and order.
Add one producer-owned `probe_execution` table immediately after that pair, before the execution report.
Its single row contains `schema`, `command_exit_code`, nullable `focused_exit_code`, and integer `trailing_bytes`.
Keep the existing verdict, subject, mutation, cause, failed-tests, restored, and selection meanings.
The original probe and selection field sets remain byte-identical.

Populate `focused_exit_code` directly from the mutated selection's Execute return.
If no mutated execution occurs, the cell is null, never an integer's default zero.
Baseline execution is a separate observation and cannot supply this field.
Baseline or mutation refusals therefore supply no mutated-execution exit.
All render paths use one producer result projection.
Keep restoration precedence, preserved-copy reporting, and current command exits unchanged.

Add `probe.ReadResult` over the producer-owned probe, selection, and probe_execution envelope.
Accept either a complete canonical tree-prefixed stdout response or the complete spill body selected by the caller.
The body form has no required tree lead, because responsebound excludes Owner.Lead from its spill file.
The current dispatcher derives that lead before running the verb.

Set trailing_bytes to the exact byte length of everything after the probe_execution row.
Compute it from the final suffix before rendering the envelope.
Include preserved-copy and restoration diagnostics and the execution report in that suffix.
A result without a suffix records zero; a baseline report still counts its bytes without supplying a mutated exit.

ReadResult first composes responsebound.RequireComplete(data []byte) error and treetarget.ReadLead(data []byte) (Identity, []byte, bool, error).
The completeness check rejects any producer projection or spill-failed frame, including malformed reserved frames.
It never opens a path named by a marker.
ReadLead accepts zero or one canonical initial identity table and returns a body view without changing the selected bytes.
It rejects an incomplete or malformed identity row; another identity envelope is not allowed.

Treetarget owns the row field set and its strict typed validation beside Row.
ReadLead reuses the pinned TOON decoder and Row's canonical rendering.
It preserves the current bool or unknown dirty value and head fallback.
The row names observed scope and grants no source or execution authority.

Promote the existing responseboundtest.ParseLine framing logic into responsebound.ProjectionLine(line string) (Projection, bool, error).
Projection carries Kind and Path; Kind is spilled or spill-failed.
The producer and classifier share the existing marker spellings in that owner.
Keep the test accessor's API as a thin projection of this owner.
RequireComplete checks complete physical lines through ProjectionLine; it does not parse TOON cells or follow a spill.

After the optional lead, require the canonical envelope to start the body.
Decode its cells with the same pinned strict decoder and require trailing_bytes to be a nonnegative integer.
The exact remaining suffix length must equal it.
A missing envelope boundary, short suffix, extra suffix, or count mismatch refuses before writing.
The count measures byte completeness only; forged bytes or a recalculated count do not prove origin or execution.

Use the pinned official `toon-go` strict decoder for cell syntax and scalar types.
Its strict mode checks table counts and widths but overwrites repeated object keys.
The owner must separately validate the exact envelope shape, field set, cardinality, and unique result envelopes.
Re-encode through the shared producer schema to check the canonical body envelope.
Never use a regular expression as a second TOON parser.

A trailing execution report is retained as provenance and grants no separate outcome authority.

The verification form gains `--probe-result <file>` as an alternative to `--excerpt` and the three scalar probe flags.
Read the explicitly selected result file once through the existing classified no-follow reader and existing evidence bound.
That selected file supplies the exact excerpt bytes, and `--ref` still names the native source.
Parse a body view only; retain and digest the complete original bytes, including an optional tree lead.
Selecting the spill body retains that body's bytes, not reconstructed stdout or its projection.
No implicit follow of a marker, fallback file, or response concatenation can change this provenance.

Retain explicit performer, model, effort, requirement, list selector, source rules, and verification `--exit-code`.
Require that exit flag to agree with the imported producer's command exit.
The import obtains outcome, normalized focused-execution exit, and restoration from `probe.ReadResult`.

It keeps the requirement's planned mutation text under the existing plan owner.
It does not claim to prove that opaque mutation prose matches the probe's argv.

Refuse an unsupported schema, absent exit, duplicate result envelope, incompatible field type, or scalar/import combination before writing.
Do not guess an exit from verdict, failed-test count, or command exit.
A current silent or restore-failed result can be recorded honestly as failed evidence.
The existing checkpoint must continue to refuse it.
Import neither runs a probe nor grants evidence authenticity beyond the existing native-reference contract.

Copied flag help names these sources exactly:

| record flag | source verb and field |
| --- | --- |
| `--exit-code` for a probe verification | `bench probe` command exit, checked against `probe_execution.command_exit_code` on import |
| `--probe-outcome` | `bench probe`, `probe.verdict` |
| `--probe-exit-code` | `bench probe`, `probe_execution.focused_exit_code`, the normalized focused-execution exit |
| `--probe-restore` | `bench probe`, `probe.restored`, through the producer owner's restoration projection |

### Headroom and closure

The selected lane runs `structure --growth` through `internal/gate/lane.go`.
Growth grades increased per-file length against the exact grant or default 400-line cap.
It does not grade directory crowding.
Existing full-tree crowding remains inherited structural debt, not a new required lane refusal.
This spec adds no count-equal grant and edits no structure budget.

`recordcmd/command.go` has 399 lines and one line of default headroom.
Move its existing verification validation and handler together into `recordcmd/verification.go` before adding import flags.
Move `cause` into the coherent repair renderer in `recordcmd/refusal.go`.
Keep the form registry and its help derivation in command.go.
The resulting source target is at most 370 lines, with at least 30 lines of growth headroom.
Each new sibling targets at most 350 lines.

`write.go` has 318 lines and 82 lines of headroom.
Keep the new assignment operation in its own sibling and add only the shared guard there.
`plan.go` has 163 lines, `coverage.go` 259, and `checkpoint.go` 159.
The shared validator and typed refusal additions fit without a grant.
`probe.go` has 254 lines, baseline.go 45, and outcome_test.go 276.
Put result parsing and its new witnesses in result.go and result_test.go.

`verification_test.go` has 391 lines and nine lines of headroom.
Put new import tests in verification_import_test.go and reuse its existing form and excerpt helpers.
Retain the existing file's behavior tests.
`recordtest/fixture.go` has 366 lines and 34 lines of headroom.
Use its existing `Option`, `NewLinked`, `WritePlan`, and committed ticket fixture without a copied builder.

The five command registry names below are co-owned because probe is a bound command package.
`command_registry_test.go` has 794 lines and no new default growth allowance.
New entry cases go into record_operations_test.go using existing command fixtures.
The co-owned AXI registry file also starts above default and needs no membership change for these mutating forms.
`help_inventory_test.go` has 323 lines and 77 lines of headroom.
The record form registry already derives root help, so no new main.go dispatch declaration is needed.

Guidance updates preserve its current anchor placements and replace only descriptions that incorrectly require hand-built dispatch JSON.
The canonical all-string-literal scan includes 11 current Go holders that pin the three guidance subjects.
Its census includes registry_calibration_test.go, registry_chunk_chain_test.go, and registry_data_test.go.
Each independent expectation remains effective; none is removed or derived from changed production values.

Large registry_data.go and docs_workflow_checks_test.go must not grow.
Place changed native-record needles in the 42-line registry_chunk_chain.go when an added rule needs headroom.
Existing unchanged rules keep their bytes and file paths.

The missing-axis repair context is created in record.go's checkReviews, not inferred from its error text.
NR-C2 owns that file along with the other typed refusal sites.
NR-C3 owns refusal_test.go's additive invalid-mutation expectation update.
Its subject-byte, empty-home, and no-child assertions remain unchanged.

The shared runner moves existing executable construction out of the 383-line probe_test.go, which currently has 17 lines of default headroom.
Baseline_test.go has 251 lines, outcome_test.go 276, refusal_test.go 313, and omit_file_test.go 87.
Thin adapters preserve those callers without new fixture derivation or crowded sibling files.
The new probetest owner targets at most 150 lines in one file; internal/probe currently has 12 immediate Go files and gains no immediate file.
No structure grant or accepted debt grows for this extraction.

The framing projection fits responsebound/owner.go's current 234 lines without a new sibling.
Its multi-case marker tests fit the current 40-line lines_test.go; preserve existing owner_test.go witnesses.
The shared responseboundtest spill accessor calls ProjectionLine without owning another marker parser.

ReadLead fits treetarget/identify.go's current 56 lines, with tests in its 149-line identify_test.go.
Keep both below the default 400-line growth limit, without a grant.
The production call chain uses these public owners and imports no test helper.
No new conformance check, hook, injected global port, or system-tagged test is required.
The exact canary directories below include recursive BASE includes, dot-restored overlays, and mutation paths.

## Implementation chunks

Each accepted green outcome now has one complete vertical ticket and one independent chunk checkpoint.
NR-C1's accepted owner prerequisite must land before its author starts.
Shared writer, form-registry, and help writes order NR-C2 after NR-C1 and NR-C3 after NR-C2.
Each chunk retains all sufficient existing tests and includes its real command entry.

The 11 anchor holders join every affected guidance chunk; their existing independent assertions remain intact.
NR-C2 includes reviewrecord/record.go for missing-axis typed repair information.
NR-C3 includes probe/refusal_test.go, both framing owners and their listed tests, the shared spill accessor, and the real CLI journey.
NR-C3 also owns the shared canned-runner extraction and every listed existing probe caller, with both caller families migrated before its green checkpoint.
No chunk depends on the unlanded exact-reader or preflight policy specs.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| NR-C1 / 1-record-native-assignments.md | First dispatch and replacement write one fully validated version-2 plan through the CLI. | NR1-NR14, NR36 | planned assignment tests, existing delegated validators, planned TestNativeRecordCLIJourney | no |
| NR-C2 / 2-enforce-native-record-order.md | Stale chunk writes refuse early and checkpoint refusals disclose the matching record repair form. | NR15-NR22, NR33-NR35, NR37 | planned ordering and checkpoint repair tests, existing amendment and completion tests | no |
| NR-C3 / 3-import-complete-probe-results.md | A probe reports its normalized focused-execution exit and verification imports its complete owner result. | NR23-NR32, NR38-NR45 | planned result/import tests and CLI journey, existing probe restoration and verification tests | no |

### Ordering and execution boundary

The stable mapping is NR-C1 to NR-C1, NR-C2 to NR-C2, and NR-C3 to NR-C3.
Each chunk retains its original rows, delivered outcome, tests, and review boundary.
Ticket 2 blocks on ticket 1, and ticket 3 blocks on both predecessors.
The shared writer, form registry, guidance, holders, fixtures, and real CLI test require this serial order.
Each predecessor commits green and closes independent Standards, Spec, and Coverage review before its successor starts.

The Writes union has 145 exact paths, equal to the unchanged 146-entry fence minus its review pickup.
Each ticket includes its first-use command, anchor, and fixture closure, not only the spec-wide union.
The 78 canary units and 11 anchor holders remain available to every affected guidance checkpoint.
Each test expectation and headroom route closes in its introducing ticket; final reconciliation cannot borrow a later repair.
Co-owned files may remain unchanged when their existing assertions and closure suffice.

This version-1 completion plan is a supported planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator must bind actual run and session identities in a version-2 plan with author-limit 1.
Each ticket then receives a fresh author on the approved gpt-5.6-sol / high line and mapped verification ownership.
No staged status, plan-only pass, or empty execution history grants implementation authority.
The complete accepted Markdown-block-reader outcome, including MB-C4, must land before those interface and headroom bindings close.

### Completion plan

The commands and mutations below are future obligations, not executed proof.
At implementation, use bench probe for each named mutation and pin its exact landed source operand and diagnostic.
Require a compiling behavioral red, byte-identical restoration, and the same focused green before retaining independent evidence.
The final requirements run on the complete reconciled source after all three chunk checkpoints close.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "NR-C1",
      "tickets": [
        "1-record-native-assignments.md"
      ],
      "verification": [
        {
          "id": "record-domain",
          "command": "bench test --package ./internal/reviewrecord"
        },
        {
          "id": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd"
        },
        {
          "id": "command-package",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "native-assignment-journey",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'"
        },
        {
          "id": "guidance",
          "command": "bench test --check docs-currency-workflow"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check conformance-meta"
        },
        {
          "id": "whole-plan-omission",
          "command": "bench test --package ./internal/reviewrecord/recordcmd --run '^TestAssignmentWholePlanValidation$'",
          "probe": "Omit only the complete candidate-validation call in RecordAssignment, retaining the validator used by ReadPlan; an unrelated invalid ticket must red TestAssignmentWholePlanValidation."
        }
      ]
    },
    {
      "id": "NR-C2",
      "tickets": [
        "2-enforce-native-record-order.md"
      ],
      "verification": [
        {
          "id": "record-domain",
          "command": "bench test --package ./internal/reviewrecord"
        },
        {
          "id": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd"
        },
        {
          "id": "checkpoint-repair",
          "command": "bench test --package ./internal/gate --run 'Test.*(Checkpoint|Review)'"
        },
        {
          "id": "command-package",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "native-ordering-journey",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'"
        },
        {
          "id": "guidance",
          "command": "bench test --check docs-currency-workflow"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check conformance-meta"
        },
        {
          "id": "ordinary-census",
          "command": "bench test --check ordinary-build-census"
        },
        {
          "id": "chunk-guard-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit only the shared current-plan digest guard at RecordChunk, retaining the guard owner and other callers; the real record entry must write stale evidence and red TestNativeRecordCLIJourney."
        }
      ]
    },
    {
      "id": "NR-C3",
      "tickets": [
        "3-import-complete-probe-results.md"
      ],
      "verification": [
        {
          "id": "probe-family",
          "command": "bench test --package ./internal/probe/..."
        },
        {
          "id": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd"
        },
        {
          "id": "response-framing",
          "command": "bench test --package ./internal/responsebound/..."
        },
        {
          "id": "tree-framing",
          "command": "bench test --package ./internal/treetarget/..."
        },
        {
          "id": "command-package",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "native-import-journey",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'"
        },
        {
          "id": "guidance",
          "command": "bench test --check docs-currency-workflow"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check conformance-meta"
        },
        {
          "id": "focused-exit-substitution",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Substitute the successful probe command exit for Execute's normalized focused-execution exit in the result producer; retain both executions and red the independent command-0/focused-1 journey assertion."
        },
        {
          "id": "tree-lead-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit optional tree-lead handling in probe.ReadResult while retaining its public owner; complete ordinary Command.Run stdout must red the native import journey."
        },
        {
          "id": "projection-refusal-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit only responsebound.RequireComplete from ReadResult while preserving envelope and suffix parsing; the projected-stdout refusal and outside-path sentinel must red the real import journey."
        },
        {
          "id": "suffix-comparison-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit only exact trailing_bytes comparison from ReadResult; a complete envelope with a truncated declared suffix must red the real import journey."
        },
        {
          "id": "suffix-renderer-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit a rendered report suffix byte without changing its declared trailing_bytes; the actual producer/consumer complete-report journey must red."
        },
        {
          "id": "runner-alternation-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Make the shared runner always emit its baseline answer while retaining both supplied Start values; existing baseline/mutated cases and the independent real journey must red."
        },
        {
          "id": "runner-passthrough-omission",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "probe": "Omit the shared runner's real-Go passthrough for non-test arguments while retaining test answers; the actual Go selection/freshness path in the native journey must red."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check native-record-operations"
    },
    {
      "id": "record-family",
      "command": "bench test --package ./internal/reviewrecord/..."
    },
    {
      "id": "probe-family",
      "command": "bench test --package ./internal/probe/..."
    },
    {
      "id": "response-framing",
      "command": "bench test --package ./internal/responsebound/..."
    },
    {
      "id": "tree-framing",
      "command": "bench test --package ./internal/treetarget/..."
    },
    {
      "id": "command-package",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "native-complete-journey",
      "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'"
    },
    {
      "id": "checkpoint-repair",
      "command": "bench test --package ./internal/gate --run 'Test.*(Checkpoint|Review)'"
    },
    {
      "id": "guidance",
      "command": "bench test --check docs-currency-workflow"
    },
    {
      "id": "root-conformance",
      "command": "bench test --check conformance-meta"
    },
    {
      "id": "ordinary-census",
      "command": "bench test --check ordinary-build-census"
    }
  ]
}
```

## Testing decisions

Test purpose precedes test count.
Use table cases in a few existing fixture families, rather than one new test per row.
Assignment tests cross `recordcmd.Command` and inspect the actual spec bytes and committed-plan validation.
Ordering tests cross the existing writer, then the real checkpoint.
Probe tests use the shared canned Go runner below to distinguish probe command 0 from focused-execution 1 deterministically.

### Shared canned runner

The current runner is private in internal/probe/probe_test.go:134-181.
Its installStubStarts owns executable construction, real-Go passthrough, the invocation marker, and baseline/mutated parity.
Recordtest.Prepare at internal/reviewrecord/recordtest/fixture.go:165 supplies the committed record fixture, not that runner.
Gittest.StubGitDir owns git transport, and kittest.WriteTree owns a minimal build-input tree.
Neither supplies the two-start Go transport needed by the real probe journey.

NR-C3 extracts only that transport into internal/probe/probetest/fixture.go.
The concrete interface is Start{Hook string, Events string, Exit int} and InstallStarts(t testing.TB, marker string, baseline, mutated Start).
Start values carry caller-authored bytes and exits; the helper owns no expected verdict, diagnostic, TOON row, or record schema.
The helper imports sanitize.ShellQuote and standard-library test/file/process utilities, without importing probe or recordcmd.
Only test callers import this support package; production dispatch and policy remain unchanged.

InstallStarts retains the existing real-Go lookup before PATH replacement, marker reset, executable mode, and resolved temporary directory.
It logs every Go start, passes non-test arguments to that real executable, and alternates the two supplied answers by test-start count.
Each hook receives the existing start ordinal before its quoted event heredoc and supplied exit.
The caller owns the marker path, while t.TempDir and t.Setenv bound files and PATH restoration to that test.
No package variable, runbinary factory replacement, injected production port, or second runner enters the plan.

In the same green NR-C3 chunk, remove the private runner and stubAnswer derivation from probe_test.go.
Its existing convenience adapters may forward caller-owned stubStart values into probetest.Start without constructing another executable.
All current callers in probe_test.go, baseline_test.go, outcome_test.go, refusal_test.go, and omit_file_test.go then reach InstallStarts.
Keep cannedPass, cannedFailure, quietRun, twoRuns, hooks, expected rows, and effect assertions with their current probe tests.
The cmd/bench journey authors its own Go event streams and independent expected exits instead of importing those expectations.

Compose recordtest.Prepare on the journey's actual temporary Go repository for its version-2 plan and committed tickets.
Reuse kittest.WriteTree for build inputs and freshness.Publish for the sealed run binary, following the existing probe fixture's selection path.
Caller-owned subject and test sources identify the focused selection; neither shared helper invents them.
The ordinary case and --full long-report case invoke the unchanged production command registry through Command.Run.
A long caller-authored failing stream must actually spill; responseboundtest.Find selects the complete body only after asserting that real marker.

The same-chunk migration must retain every existing probe assertion and prove both caller families use the one runner.
A future omission of non-test passthrough must red the real Go selection/freshness path rather than silently substituting a planted handler.
Omitting answer alternation must red the existing baseline/mutated cases and the real command journey's independent command-0/focused-1 assertion.
These are implementation-time mutation witnesses, not reds executed during this planning repair.

### Native command journey

The CLI journey crosses `Command.Run` for the producer and consumer and grades the retained record through its real owner.
Its bounded producer captures ordinary stdout with the canonical tree lead and imports that explicitly selected file.
A longer canned execution report uses the existing probe --full flag to force an actual spill.
The test selects its complete body and imports those exact bytes.
Both journeys use the production registry, not a planted handler or a domain-only response.
Assert each retained excerpt and digest against the selected file, rather than a re-encoded result or reconstructed stdout.

The same journey supplies projected stdout, spill-failed stdout, and separately truncated envelope and suffix files to the consumer.
Each refusal leaves record bytes unchanged and never opens a marker-named outside sentinel.
A complete suffix with its count changed in either direction also refuses.
A responsebound-owned failure fixture supplies the partial spill whose complete envelope precedes the shortened suffix.
These cases distinguish complete-report import from an envelope-only parser.

Existing amendment, primary-checkout, no-follow, oversized-record, temporary-write, history, and restoration tests remain sufficient for their unchanged behavior.
Do not add tests for the existence of new declarations.
Place a new lower seam only when the command cannot expose its failure.
Instance-local file effects can drive a failure case where existing filesystem fixtures are insufficient.
Any such seam requires the existing injected-port audit's exact registry closure before implementation uses it.
The current plan does not require a new injected port.

The existing dispatcher framing stays unchanged in main.go, tree_scope.go, and census_output.go.
Those sources were read to derive the actual journey; no production change there is needed or authorized.

Independent expectations for NR23 and NR37 are necessary because matching producer and consumer derivations can agree on a wrong exit or order.
During implementation, record future red proof by substituting the probe command exit for Execute's returned exit.
Also omit the chunk digest guard at the real entry.
The corresponding named tests must fail, then pass after exact restoration.

Also omit tree-lead handling, ignore a reserved projection frame, and bypass the exact trailing_bytes comparison.
The complete ordinary and spill journeys and the incomplete-file sentinels must then fail.
Change one renderer suffix without updating its count to exercise producer/consumer disagreement.
These reds are planned, not observed during specification.

### Seam diagram

    explicit native argv + frozen plan + committed tickets
        │
        ▼
    recordcmd form registry ──▶ reviewrecord complete validator ──▶ canonical atomic replacement
                                      │
                                      ▼
    checkpoint refusal ──▶ typed repair name ──▶ matching native form help

    actual focused execution ──▶ probe result owner ──▶ native excerpt ──▶ record verification writer

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| NR1 | 1 | First dispatch appends one assignment to the target empty history | planned TestAssignmentFirstDispatch in internal/reviewrecord/recordcmd/assignment_test.go | A writer that only replaces an existing occurrence leaves the target undispatched. |
| NR2 | 2 | A mismatched run identifier refuses without changing the spec | planned TestAssignmentExecutionIdentity in internal/reviewrecord/recordcmd/assignment_test.go | Changing the execution identity instead of checking it produces a different run. |
| NR3 | 2 | A mismatched orchestrator, author limit, or review mode refuses without changing the spec | planned TestAssignmentExecutionIdentity in internal/reviewrecord/recordcmd/assignment_test.go | Checking only the run identifier admits a different execution posture. |
| NR4 | 3 | An invalid non-target ticket makes assignment refuse before creating a temporary file | planned TestAssignmentWholePlanValidation in internal/reviewrecord/recordcmd/assignment_test.go | Validating only the target history writes an invalid complete plan. |
| NR5 | 3 | A missing or cyclic ticket dependency makes assignment refuse before replacement | planned TestAssignmentWholePlanValidation in internal/reviewrecord/recordcmd/assignment_test.go | A copied assignment validator misses the committed ticket graph. |
| NR6 | 4 | An assignment source older than HEAD refuses without replacing the spec | planned TestAssignmentSourceRefusals in internal/reviewrecord/recordcmd/assignment_test.go | Resolving a valid old commit alone overwrites current dispatch history. |
| NR7 | 4 | A changed working spec or planned ticket refuses before any write | planned TestAssignmentSourceRefusals in internal/reviewrecord/recordcmd/assignment_test.go | Reading committed inputs while replacing divergent working bytes loses authored changes. |
| NR8 | 5 | Replacement appends after the named preceding session | planned TestAssignmentReplacementHistory in internal/reviewrecord/recordcmd/assignment_test.go | Replacing the last array element erases a participant that review independence must retain. |
| NR9 | 6 | A replacement lacking the required stopped or preserved evidence refuses | planned TestAssignmentReplacementHistory in internal/reviewrecord/recordcmd/assignment_test.go | Treating a predecessor or cleanup callback as sufficient bypasses the existing replacement contract. |
| NR10 | 6 | A no-progress replacement lacking reassessment refuses | planned TestAssignmentReplacementHistory in internal/reviewrecord/recordcmd/assignment_test.go | A generic replacement path misses the trigger-specific requirement. |
| NR11 | 7 | A session already assigned to another ticket refuses | planned TestAssignmentWholePlanValidation in internal/reviewrecord/recordcmd/assignment_test.go | A target-only check permits one session to own two tickets. |
| NR12 | 8 | A successful assignment changes only the completion-plan payload | planned TestAssignmentPreservesDocument in internal/reviewrecord/recordcmd/assignment_test.go | Rendering a new whole spec drops the prefix, suffix, or original line endings. |
| NR13 | 9 | A failed temporary write leaves the original spec bytes intact | planned TestAssignmentWriteFailure in internal/reviewrecord/recordcmd/assignment_test.go | Writing the destination before validation destroys the retry source. |
| NR14 | 10 | Record family help exposes the complete assignment form | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | A domain-only writer cannot be invoked through the real command entry. |
| NR15 | 11 | A changed plan digest refuses a chunk write until an amendment is recorded | planned TestChunkNeedsCurrentAmendment in internal/reviewrecord/recordcmd/ordering_test.go | Current RecordChunk can append structurally valid evidence under a stale record plan digest. |
| NR16 | 12 | A spec prose edit changes the plan digest used by the chunk refusal | planned TestChunkNeedsCurrentAmendment in internal/reviewrecord/recordcmd/ordering_test.go | A semantic exemption would silently remove the whole-spec input from ReadPlan. |
| NR17 | 13 | A missing chunk checkpoint refusal names the chunk repair form | planned TestCheckpointEntryRepairForms in internal/gate/record_repair_test.go | A generic preflight route alone makes the operator rediscover the missing writer. |
| NR18 | 14 | A missing planned verification checkpoint refusal names the verification repair form | planned TestCheckpointEntryRepairForms in internal/gate/record_repair_test.go | Routing every missing entry to chunk writes the wrong artifact. |
| NR19 | 15 | A missing axis review checkpoint refusal names the review repair form | planned TestCheckpointEntryRepairForms in internal/gate/record_repair_test.go | A verification-only route cannot satisfy the independent review obligation. |
| NR20 | 16 | A missing completion checkpoint refusal names the completion repair form | planned TestCheckpointEntryRepairForms in internal/gate/record_repair_test.go | A generic entry route cannot write reconciliation. |
| NR21 | 17 | An undispatched ticket refusal names the assignment repair form | planned TestCheckpointEntryRepairForms in internal/gate/record_repair_test.go | The operator otherwise must edit the execution history by hand. |
| NR22 | 18 | Completion help derives amendment-before-chunk order from the repair owner | planned TestCompletionHelpOrdering in internal/reviewrecord/recordcmd/ordering_test.go | An independent copied help sequence can disagree with the actual refusal order. |
| NR23 | 19 | A biting probe reports focused-execution exit 1 while its command returns 0 | planned TestProbeReportsNormalizedExecutionExit in internal/probe/result_test.go | Substituting the successful probe command exit reports passing mutated execution. |
| NR24 | 19 | A probe that never executes the mutated selection reports no focused-execution exit | planned TestProbeReportsNormalizedExecutionExit in internal/probe/result_test.go | Guessing from verdict words invents execution evidence for a refused run. |
| NR25 | 20 | Result import fills the planned probe from the probe owner projection | planned TestVerificationImportsProbeResult in internal/reviewrecord/recordcmd/verification_import_test.go | A CLI that still requires copied outcome flags does not deliver native import. |
| NR26 | 21 | Imported evidence retains the exact supplied native excerpt and its computed digest | planned TestVerificationImportsProbeResult in internal/reviewrecord/recordcmd/verification_import_test.go | Re-encoding the receipt loses its original provenance bytes. |
| NR27 | 22 | An older envelope without the focused-execution field refuses import | planned TestProbeResultEnvelopeRefusals in internal/probe/result_test.go | Deriving an exit from failed_tests accepts a producer that did not report it. |
| NR28 | 22 | A duplicate or substituted result envelope refuses import without writing evidence | planned TestVerificationImportRefusals in internal/reviewrecord/recordcmd/verification_import_test.go | The pinned decoder otherwise overwrites repeated object keys. |
| NR29 | 23 | An imported silent result remains a failed probe at the checkpoint | planned TestVerificationImportsProbeResult in internal/reviewrecord/recordcmd/verification_import_test.go | Import must not interpret successful restoration as a biting mutation. |
| NR30 | 23 | An imported failed restoration cannot pass the checkpoint | planned TestVerificationImportsProbeResult in internal/reviewrecord/recordcmd/verification_import_test.go | A biting test result alone does not prove restoration. |
| NR31 | 24 | The existing scalar verification form retains its explicit probe values | `internal/reviewrecord/recordcmd/verification_test.go` (`TestRecordVerificationWritesThePlannedProbe`) | The import must not reinterpret a caller-supplied scalar exit as the command exit. |
| NR32 | 25 | Verification help names the producer verb and field for each copied probe flag | planned TestVerificationHelpFieldSources in internal/reviewrecord/recordcmd/verification_import_test.go | Labels such as exit alone do not distinguish the three exit meanings. |
| NR33 | 26 | Missing-record and missing-chunk errors remain reachable through errors.Is | planned TestRecordRepairPreservesErrorIdentity in internal/reviewrecord/repair_test.go | Wrapping a repair hint as plain text destroys the existing sentinel identity. |
| NR34 | 27 | A valid amendment retains preceding evidence and the mapped digest chain | `internal/reviewrecord/recordcmd/amendment_test.go` (`TestRecordAmendmentChainsFromTheLastDigest`) | Replacing history instead of appending breaks the second amendment predecessor. |
| NR35 | 28 | Completion continues to refuse capability-blocked reconciliation values | review-owned scope comparison at internal/reviewrecord/completion.go and roadmap/FT317.md | Adding a new completion value would decide the separate unresolved capability. |
| NR36 | 29 | Assignment refuses version-1 plans without converting their verification mappings | planned TestAssignmentWholePlanValidation in internal/reviewrecord/recordcmd/assignment_test.go | An implicit conversion cannot infer which ticket owes every existing requirement. |
| NR37 | 30 | The real record entry refuses a stale plan before its writer changes any bytes | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | Owner-only checks stay green when command dispatch calls the writer before the guard. |
| NR38 | 20 | The real record entry imports complete canonical tree-prefixed probe stdout | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | A parser that requires probe at byte zero refuses the actual dispatcher response. |
| NR39 | 21 | Import of the complete spill body retains exactly that selected file's bytes | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | Reconstructing the omitted tree lead or using projected stdout changes the excerpt digest. |
| NR40 | 22 | Projected stdout refuses import without changing record bytes | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | Ignoring spilled or cut output accepts evidence the owner marked incomplete. |
| NR41 | 22 | Spill-failed stdout refuses import without opening its marker path | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | A complete-looking prefix cannot permit marker following or fallback to a different file. |
| NR42 | 22 | A file truncated within a mandatory envelope refuses import | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | Prefix-only field extraction misses the unfinished selected result. |
| NR43 | 22 | A file truncated within the declared suffix refuses import | planned TestNativeRecordCLIJourney in cmd/bench/record_operations_test.go | A partial spill can retain all outcome cells while losing later native evidence. |
| NR44 | 22 | A complete suffix that disagrees with trailing_bytes refuses import | planned TestProbeResultEnvelopeRefusals in internal/probe/result_test.go | Count changes and renderer suffix omissions detect unchecked completeness metadata. |
| NR45 | 22 | A malformed or repeated tree identity envelope refuses import | planned TestProbeResultEnvelopeRefusals in internal/probe/result_test.go | Blindly stripping two lines or overwriting a second tree object accepts noncanonical framing. |

### Edge inventory

- Empty assignment history accepts first dispatch, while an absent history key remains invalid under the full validator.
- Empty, duplicate, controlled, unknown, and repeated CLI operands retain the existing usage and safe-output contracts.
- Invalid non-target tickets, dependency cycles, unplanned tickets, reused sessions, and mismatched execution assertions refuse before replacement.
- Stale HEAD, divergent working inputs, symlink components, FIFO, empty spec, malformed JSON, duplicate keys, and oversized records remain refusal states.
- Replacement triggers enumerate the existing `Triggers()` values, with no-progress's additional reassessment obligation.
- Unicode prose and CRLF outside the payload retain exact bytes.
- Existing historical verification may still name its recorded source, rather than current HEAD.
- A source-bound assignment write does not tighten that independent historical evidence policy.
- Probe envelopes cover bit, silent, invalid, restore-failed, and no mutated execution.
- Scalar flags and imported fields cannot both supply the same result.
- Complete stdout and complete spill bodies are separate producer forms with exact selected-file provenance.
- Projection markers, failed spills, truncated mandatory cells, truncated suffixes, and incorrect byte counts refuse.
- A canonical tree row remains descriptive data and does not authenticate the native source.

**Won't handle** — concurrent unauthorized file writers after the final preimage check — the orchestrator remains the serialized plan writer.
**Won't handle** — a proof that native dispatch metadata stopped an external process — the existing replacement schema requires explicit stopped writer evidence.
**Won't handle** — matching free-form planned mutation prose to probe argv — the current requirement owner supplies that text without a structured mutation language.
**Won't handle** — capability-blocked reconciliation — FT317 still owns that unresolved completion value.

## Ownership fences

These prospective files include the owner, test, registry, anchor, and fixture closure.
Co-owned existing tests remain untouched when their current witnesses suffice.
No implementation write is authorized by this planning artifact alone.

- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.bench/BENCH-reference.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/record_operations_test.go`
- `cmd/bench/response_bound_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/gate/checkpoint.go`
- `internal/gate/record_repair_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/probe/baseline.go`
- `internal/probe/baseline_test.go`
- `internal/probe/omit_file_test.go`
- `internal/probe/outcome_test.go`
- `internal/probe/probe.go`
- `internal/probe/probe_test.go`
- `internal/probe/probetest/fixture.go`
- `internal/probe/refusal_test.go`
- `internal/probe/result.go`
- `internal/probe/result_test.go`
- `internal/responsebound/lines_test.go`
- `internal/responsebound/owner.go`
- `internal/responsebound/owner_test.go`
- `internal/responsebound/responseboundtest/spill.go`
- `internal/reviewrecord/assignment.go`
- `internal/reviewrecord/assignment_test.go`
- `internal/reviewrecord/check.go`
- `internal/reviewrecord/completion.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/delegated.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/plan.go`
- `internal/reviewrecord/record.go`
- `internal/reviewrecord/recordcmd/amendment_test.go`
- `internal/reviewrecord/recordcmd/assignment.go`
- `internal/reviewrecord/recordcmd/assignment_test.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/completion.go`
- `internal/reviewrecord/recordcmd/ordering_test.go`
- `internal/reviewrecord/recordcmd/refusal.go`
- `internal/reviewrecord/recordcmd/refusal_test.go`
- `internal/reviewrecord/recordcmd/verification.go`
- `internal/reviewrecord/recordcmd/verification_import_test.go`
- `internal/reviewrecord/recordcmd/verification_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/reviewrecord/repair.go`
- `internal/reviewrecord/repair_test.go`
- `internal/reviewrecord/write.go`
- `internal/treetarget/identify.go`
- `internal/treetarget/identify_test.go`
- `reviews/native-record-operations.md`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-step-move`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-axis-exclusions`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`

## Out of scope

- Version-1 conversion: estimated five edits and three gate runs. It needs explicit verification-to-ticket mapping and a separate authoring contract.
- Capability-blocked completion: estimated three edits and two gate runs after FT317's decision. Existing completion validation remains the surviving caller.
- Arbitrary prose semantic-digest exemption: estimated five edits and three gate runs. ReadPlan continues to include complete spec bytes.
- FT168's replacement-file, system, multi-mutation, and multi-package forms: estimated six edits and four gate runs. Existing probe selections remain the surviving caller.
- Raw child-process status export: estimated two edits and two gate runs. The accepted producer field reports Execute's normalized focused-execution exit.
- General JSON editing, agent dispatch, live qualification, and automatic source commits are separate capabilities and receive no write path here.

## Further notes

The ready shared decision map remains in place.
No map answer, FT290 field, roadmap priority, commitment, or other approved implementation line changes.
FT318 remains partial because this spec excludes its FT317-dependent value.
FT168 remains partial because only the result projection supports native import here.
C10's 293-to-375-to-318-to-125 sequence remains advice.
The separate preflight specs own closure and staleness grading, while this spec consumes their existing interfaces.

The verification record is static: full source bodies, owner definitions, caller sweeps, fixture closure, and growth policy were inspected.
Prose, coverage, link checks, and the required planning lane freeze this spec before independent review.
No runtime test, omission probe, benchmark, or implementation result is claimed.
The graph retains exact per-ticket closure and shared-write blockers from the accepted source.
Independent ticket review accepted the graph; build entry must refresh the landed provider, source pin, and headroom before any charge.

| approval item | proposed disposition |
| --- | --- |
| implementation line | gpt-5.6-sol / high, retained accepted line |
| seams and test purpose | Existing validators, producer result owner, real record entry, and gate checkpoint |
| coverage and exclusions | NR1-NR45, with explicit FT317, conversion, and digest-policy exclusions |
| ownership | Unchanged 146-entry fence; exact per-ticket closure independently accepted |
| graph | Three serial vertical tickets accepted at 0d7fa110cc3869a4b089cba73cdb2aa435d38983 after SPEC acceptance |

The numbered breakdown is:

| ticket | title | Blocked by | delivered outcome |
| --- | --- | --- | --- |
| 1-record-native-assignments.md | Record validated native assignments | none | First and replacement dispatch through the fully validated real assignment entry |
| 2-enforce-native-record-order.md | Enforce amendment order and native repair routes | 1-record-native-assignments.md | Early stale-plan refusal and matching repair forms through the actual checkpoint |
| 3-import-complete-probe-results.md | Import complete native probe results | 1-record-native-assignments.md, 2-enforce-native-record-order.md | Complete actual producer/consumer import with one shared runner and preserved scalar behavior |
