# Durable caller migration

Status: staged

Decision source: named reviewed artifact `decisions/architecture-primitives.md`, resolved tickets `2.md` and `4.md` (2026-10-06).

Verification log: 2 original-spec iterations, 1 split-spec iteration, and 1 ticket iteration to accept — separate spec acceptance preceded reslicing.

## Problem

Replacement callers retain separate algorithms and inconsistent synchronization results.
Some full callers discard post-rename uncertainty or compensate a visible new ledger.
These accepted caller outcomes require a separate complete build after the durable leaf lands.

## Solution

Migrate the remaining classified callers to the completed durable-file-replacement prerequisite.
Preserve each caller's payload, mode, validation, destination authority, sentinel translation, and recovery policy.
Complete every caller failure witness and source-ownership closure without duplicating the leaf algorithm.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: failure ordering and real caller composition require stronger evidence than helper happy paths.
Author line: gpt-6.1-sol / high / one split pass plus at most two bounded repairs.
Harder chunks: D-B1, D-B5, D-B6, D-C4.

16. As a gate verdict caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
17. As a gate evidence and lane caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
18. As a prospective owner record caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
20. As a capture transaction files caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
21. As a handoff document caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
22. As an intent ledger caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
23. As a commitment repository caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
24. As a publication record caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
25. As a broker manifest caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
26. As a dashboard artifact caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
27. As a probe subject caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
28. As an assessment record caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
29. As a repair pilot document caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
30. As a caller, I want to retain its refusal policy, so that invalid input never reaches replacement.
31. As a maintainer, I want to remove independent replacement copies, so that the operation has one owner.
33. As a maintainer, I want to retain specialized protocols, so that file replacement does not erase transaction policy.
34. As a probe operator, I want to retain the preserved copy after failure, so that recovery remains available.
35. As a maintainer, I want to remove an unused replacement copy, so that an unreachable helper cannot preserve a second algorithm.
36. As a maintainer, I want to measure preserved caller behavior, so that the refactor keeps payload and refusal contracts.
37. As a gate operator, I want to see failed pending persistence, so that the oracle cannot start after uncertain persistence.
38. As a gate operator, I want to see failed terminal persistence, so that the failure cannot advertise durable success.
39. As a gate operator, I want to see failed evidence persistence, so that an oracle result cannot hide a storage failure.
40. As a gate operator, I want to see failed pending recovery, so that the failure cannot claim a rollback that failed.
41. As a probe operator, I want to keep recovery after failed mutation, so that the original bytes remain available.
42. As a probe operator, I want to restore after failed mutation, so that a failed durability step does not leave an avoidable mutation.
43. As a probe operator, I want to keep the pre-rename refusal, so that the migration preserves an unchanged subject.
44. As a transaction caller, I want to retain committed external effects, so that a visible new ledger has matching effects.
45. As a transaction caller, I want to restore uncommitted external effects, so that a failed preparation preserves the old transaction.
46. As an approval caller, I want to restore an uncertain first policy write, so that failed stage cannot leave an unapproved new policy.
47. As an approval caller, I want to restore a failed board stage, so that partial stage does not escape a successful rollback.
48. As an approval caller, I want to restore a refused ledger write, so that an unapproved receipt retains the original external state.
49. As an approval caller, I want to retain a published receipt with matching files, so that compensation cannot split the visible approved state.
50. As an approval caller, I want to see failed compensation, so that an incomplete rollback remains an error.
51. As an approval caller, I want to retain board absence after failed stage, so that a linked repository receives no invented board.
52. As a reauthorization caller, I want to retain matching request effects, so that a new ledger request cannot retain the old external effect.


## Implementation decisions

### Complete prerequisite

This spec depends on complete acceptance and landing of specs/durable-file-replacement/spec.md.
The prerequisite delivers D-A with D01-D15, D19, and D32.
No caller chunk starts from an accepted but unlanded leaf chunk.
Process-lifetime PL-C2 depends only on that complete prerequisite, not this successor.

The prerequisite owns internal/durablefile and its complete operation and error contract.
Consume Replace(path string, data []byte, mode fs.FileMode) error and the same-signature Replacer.Replace.
Consume *durablefile.Error through errors.As, including Stage, Published, and Unwrap for errors.Is.
Published means rename completed, and its failure message retains `new bytes may already be visible`.
Do not copy its algorithm or introduce another replacement owner.
Caller mode, validation, authority, parent creation, and transaction decisions stay with the existing caller.

The assessment and repair-pilot adapters retain their existing failure seams.
They compose those seams into the prerequisite operations.
Other caller tests attach through a narrow caller writer that receives the replacement operation.
The exported caller supplies the operating-system owner.
Do not add ambient package-variable swaps to subprocess tests.
If a new audited port appears, update its registry row with a real-producer test.

### Caller outcomes after replacement failure

The leaf classifies one file operation. The caller decides the transaction outcome.
Classify wrapped errors with errors.As, not message text.
An unclassified failure follows the existing pre-rename contract.
Keep the original failure even when a later recovery succeeds.

The gate execution caller retains its pending, timeout, evidence, or final diagnostic category.
For a replacement failure, append the cause and its visibility classification.
Pending persistence failure prevents the oracle start.
Evidence or final persistence failure returns ActionExit 1.
The returned Inspection must not advertise reusable green after a persistence failure.

Keep the existing pending-record recovery attempt after terminal persistence failure.
Report a failed recovery attempt beside the original failure.
Do not claim that pending bytes survived when the final rename completed and recovery failed.
The classified error describes possible visibility, not a universal rollback result.
The interrupted-green demotion path must return its persistence error to the diagnostic caller.
Sources: `internal/gate/run_transaction.go:175`, `internal/gate/run_transaction.go:213`, and `internal/gate/run_transaction.go:230`.

A failed probe mutation before rename retains its existing refusal and cleanup behavior.
After a completed mutation rename, keep the preservation copy until restore succeeds and read-back matches the original bytes.
Attempt restore before any focused mutation run starts.
If restore succeeds, release the copy and return the original mutation failure at exit 1.
If restore fails, retain the copy and return exit 2 with the preserved path.
The reason carries both the mutation failure and the restore failure.

A completed restore rename with failed synchronization still counts as failed restore.
The preserved bytes remain available even when the visible subject already matches them.
The existing TOON refusal policy must not hide the preserved path.
Sources: `internal/probe/probe.go:148`, `internal/probe/probe.go:197`, and `internal/probe/outcome_test.go:44`.

Intent compensates a failed ledger write only when its rename did not complete.
That compensation runs once under the existing lock.
When ledger rename completes, retain the decision's matching external effects.
Return the classified durability failure without compensation or a successful transaction result.
Release the ledger lock on both outcomes.
Sources: `internal/intent/transaction.go:46` and `internal/intent/transaction_test.go:184`.

Commitment stage snapshots policy and board before the first policy write.
Construct the complete restore closure before that write.
A failed policy or board stage invokes restore, including a failed initial policy write after rename.
With successful restore, the original policy and board bytes or absence return before approval fails.
The receipt remains unapproved because the ledger write never starts.
Sources: `internal/commitment/repository/repository.go:275` and `internal/commitment/repository/repository.go:297`.

After complete stage, a pre-rename ledger failure restores the external snapshots once.
A post-rename ledger failure retains the new policy, projected board, and approved receipt together.
Approve returns an error even when that coherent new state is visible.
No compensation may pair an approved new receipt with old policy or board bytes.
The board-absent branch retains absence in each outcome.
Sources: `internal/commitment/repository/repository.go:165` and `internal/commitment/store_test.go:89`.

Restore returns an error instead of discarding it.
Approve captures compensation failure without changing the public Compensation callback signature.
Join that failure with the initiating error and name each unrestored path.
Report recovery as incomplete until snapshot read-back confirms bytes, mode, or original absence.
A failed restore promises neither rollback nor durable recovery.
These rules address errors returned during this invocation, not interruption between separate files.

### Caller failure seams

The exported gate entry supplies the native replacer to its internal execution form.
A test replaces that operation within the same execution call.
The internal form still runs the real pending, oracle, evidence, and terminal branches.
Each fault identifies the destination and operation ordinal so recovery can fail independently.
Do not inject only at durableReplace or inspect a fabricated Result.

The probe entry supplies the native replacer to its internal probe form.
Its test fixture counts baseline and focused starts separately.
A mutation failure permits the completed baseline but no focused mutation start.
The test injects separate mutation and restore failures through the same caller invocation.
Do not call replaceAtomic alone for these outcomes.

Add intent.TransactWithReplacer with the existing Transact operands and one durablefile.Replacer operand.
The existing Transact delegates to that form with the native replacer.
Both forms share address resolution, lock, decision, write, compensation, and release logic.
Commitment Approve's internal form passes its invocation replacer into stage and that real ledger transaction.
The exported Approve supplies the native replacer.
This seam permits one destination-specific failure across policy, board, and ledger writes.

The reauthorization internal form uses that transaction seam without changing its exported behavior.
Its test exercises the real request transition and a reversible external-effect witness.
A post-rename failure keeps the next request digest and its matching external effect.
The existing pre-rename rollback test remains unchanged.
Source: `internal/intent/assignment.go:158` and `internal/intent/transaction_test.go:236`.

### Current caller inventory

The sources below pin the tree at `a382f4848d432815968f044820fad593e856e297`.
A source change invalidates its inventory row.
The row IDs below each require caller composition evidence.

| caller | source and symbol | caller mode | caller payload | row |
|---|---|---|---|---|
| gate verdict | `internal/gate/verdict.go:88` (`durableReplace`) | 0600 | Marshal plus final newline | D16 |
| gate evidence and lane | `internal/gate/engine.go:186` (`durableReplaceRecordAt`) | 0600 | Caller record bytes plus final newline | D17 |
| prospective owner record | `internal/gate/prospectiveartifact/prospectiveartifact.go:229` (`Publish`) | RecordMode | Marshal plus final newline | D18 |
| capture transaction files | `internal/capturetx/transaction.go:235` (`replace`) | 0644 | Caller document and manifest bytes | D20 |
| handoff document | `internal/handoffdoc/store.go:191` (`replace`) | 0644 | Caller rendered document | D21 |
| intent ledger | `internal/intent/intent.go:189` (`writePath`) | 0600 | Sorted ledger plus final newline | D22 |
| commitment repository | `internal/commitment/repository/repository.go:355` (`atomicWrite`) | Caller mode | Stage bytes and snapshot restore bytes | D23 |
| publication record | `internal/publication/record.go:80` (`SaveRecord`) | 0644 | Indented record plus final newline | D24 |
| broker manifest | `internal/brokermanifest/brokermanifest.go:34` (`Write`) | 0644 | Validated digest and tab-framed fields | D25 |
| dashboard artifact | `internal/dashboard/dashboard.go:138` (`atomicWrite`) | 0600 | Rendered HTML | D26 |
| probe subject | `internal/probe/subject.go:177` (`replaceAtomic`) | Subject mode | Mutation or preserved subject bytes | D27 |
| assessment record | `internal/assessment/store.go:116` (`replace`) | 0600 | Validated indented record plus final newline | D28 |
| repair pilot document | `internal/repairpilot/command.go:261` (`writeDocument`) | 0600 | Indented document plus final newline | D29 |

Gate verdict callers are execution and interrupted-pending persistence in `internal/gate/run_transaction.go:175`.
Their full execution failure outcomes belong to D37-D40.
Gate evidence also serves lane records through `internal/gate/lane.go:331`.

Capture replacement serves Append, Abort, recovery, and its manifest writer.
Commitment replacement serves stage and file-snapshot restore.
Intent transaction and commitment approval compose these operations through D44-D52.
Probe replacement serves mutation and restore.
The complete probe caller retains recovery authority through D41-D43.
The publication state machine retains its release lock and state transitions.
Sources: `internal/capturetx/public.go:49`, `internal/capturetx/transaction.go:77`, `internal/commitment/repository/repository.go:297`, and `internal/publication/statemachine.go:115`.

The unused `internal/adopt/manifest.go:76` writeManifest has no static consumer in the current tree.
The real adoption manifest path uses stageBytes in `internal/adopt/link_stage.go:101`.
Delete only the unused replacement helper.
Keep manifestBytes, manifestMode, and the lifecycle transaction unchanged.


Review-record replacement is already delivered by the prerequisite through D19.
D30 and D36 retain its existing validation and differential assertions through prerequisite evidence.
Successor execution must preserve that landed caller without expanding its implementation fence.

### Explicit behavior changes

| existing rule | new rule | rows |
|---|---|---|
| Review, capture, commitment, and pilot writers sync only the file. | The owner also syncs the directory. | D19, D20, D23, D29 |
| Handoff and intent writers ignore directory failures. | A directory failure returns an error. | D21, D22 |
| Publication, broker, dashboard, probe, and assessment writers omit file sync. | The owner syncs the file and directory. | D24-D28 |
| Some modes apply after file sync. | The mode precedes file sync. | D06 |
| Post-rename errors lack a visibility classification. | The error records completed rename. | D11-D13 |


## Implementation chunks

The former D-A maps to the complete prerequisite spec and remains its first stable capability.
D-B1 through D-B6, D-C1 through D-C6, and D-D retain their stable IDs here.
The historical fourteen-ticket graph remains preserved at the snapshot commit.
Independent split-spec review accepted the prerequisite and successor before this graph was authored.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| D-B1 / 02-preserve-gate-failures.md | Preserve durability failures through gate execution | D16, D17, D37, D38, D39, D40 | gate-tests | yes |
| D-B2 / 03-migrate-owner-records.md | Migrate prospective owner-record replacement | D18 | caller-tests | no |
| D-B3 / 04-migrate-capture-documents.md | Migrate capture transaction file replacement | D20 | caller-tests | no |
| D-B4 / 05-migrate-handoff-documents.md | Migrate handoff document replacement | D21 | caller-tests | no |
| D-B5 / 06-preserve-ledger-effects.md | Preserve matching effects after ledger publication | D22, D44, D45, D52 | caller-tests | yes |
| D-B6 / 07-recover-approval-stage.md | Recover failed approval stages and retain published receipts | D23, D46, D47, D48, D49, D50, D51 | repository-tests, approval-tests, command-contract, axi-registry, routing | yes |
| D-C1 / 08-migrate-publication-records.md | Migrate publication record replacement | D24 | caller-tests, port-registry | no |
| D-C2 / 09-migrate-broker-manifests.md | Migrate broker manifest replacement | D25 | caller-tests | no |
| D-C3 / 10-migrate-dashboard-artifacts.md | Migrate dashboard artifact replacement | D26 | caller-tests | no |
| D-C4 / 11-preserve-probe-recovery.md | Preserve probe recovery after failed mutation | D27, D34, D41, D42, D43 | caller-tests, command-contract, axi-registry, routing | yes |
| D-C5 / 12-migrate-assessment-records.md | Migrate assessment record replacement | D28 | caller-tests, port-registry, command-contract, axi-registry, routing | no |
| D-C6 / 13-migrate-repair-documents.md | Migrate repair pilot document replacement | D29 | caller-tests | no |
| D-D / 14-close-replacement-ownership.md | Remove the unused replacement owner and close caller coverage | D30, D31, D33, D35, D36 | adoption-tests, coverage | no |

D-B5 supplies the real classified ledger transaction seam to D-B6.
Shared command registries order D-B6, D-C4, and D-C5.
The audited port registry orders D-C1 before D-C5.
D-D waits for every successor caller migration and consumes the prerequisite's reviewed caller evidence.
Each chunk requires its own real caller tests and independent review before dependent successors.

### Completion plan

The version 1 plan records future implementation evidence.
It claims no current implementation pass, red, native qualification, or benchmark.
The orchestrator adds required version 2 author sessions before dispatch.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "D-B1",
      "tickets": [
        "02-preserve-gate-failures.md"
      ],
      "verification": [
        {
          "id": "gate-tests",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/gate",
          "probe": "Discard the leaf cause, allow oracle start after pending failure, hide evidence failure, or discard the recovery error. Each mutation must fail its real execution witness."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/gate",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-B2",
      "tickets": [
        "03-migrate-owner-records.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/gate/prospectiveartifact"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "probe": "Restore the former local replacement sequence. The real owner-record junction test must detect the missing shared failure."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-B3",
      "tickets": [
        "04-migrate-capture-documents.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/capturetx"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/capturetx",
          "probe": "Bypass the shared owner in capture replace. The capture operation failure witness must turn red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/capturetx",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-B4",
      "tickets": [
        "05-migrate-handoff-documents.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/handoffdoc"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/handoffdoc",
          "probe": "Restore ignored directory errors or bypass the leaf. The section-writing composition test must detect the missing failure."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/handoffdoc",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-B5",
      "tickets": [
        "06-preserve-ledger-effects.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/intent"
        },
        {
          "id": "existing-canary-owner",
          "command": "bench test --check bounds-policy"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/intent",
          "probe": "Compensate every write error or omit the reauthorization effect check. The new-ledger and external-witness test must detect their mismatch."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/intent",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-B6",
      "tickets": [
        "07-recover-approval-stage.md"
      ],
      "verification": [
        {
          "id": "repository-tests",
          "command": "bench test --package ./internal/commitment/repository"
        },
        {
          "id": "approval-tests",
          "command": "bench test --package ./internal/commitment"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/commitment/repository",
          "probe": "Create restore only after policy write, compensate a published ledger, restore only one file, or discard restore failure. Each mutation must fail the corresponding real approval witness."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/commitment/repository",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C1",
      "tickets": [
        "08-migrate-publication-records.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/publication"
        },
        {
          "id": "port-registry",
          "command": "bench test --check injected-port-registry"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/publication",
          "probe": "Restore the former unsynchronized writer. The publication record junction test must detect the missing failure."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/publication",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C2",
      "tickets": [
        "09-migrate-broker-manifests.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/brokermanifest"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/brokermanifest",
          "probe": "Bypass the leaf with the former manifest writer. The manifest composition test must report the injected missing-failure red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/brokermanifest",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C3",
      "tickets": [
        "10-migrate-dashboard-artifacts.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/dashboard"
        },
        {
          "id": "existing-canary-owner",
          "command": "bench test --check git-plumbing-owner"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/dashboard",
          "probe": "Select 0644 by assumption or restore the former writer. The mode or injected-failure witness must turn red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/dashboard",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C4",
      "tickets": [
        "11-preserve-probe-recovery.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/probe"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/probe",
          "probe": "Release preservation on every mutation error or treat failed restore sync as success. The preservation-path and focused-start witnesses must detect each mutation."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/probe",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C5",
      "tickets": [
        "12-migrate-assessment-records.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "port-registry",
          "command": "bench test --check injected-port-registry"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/assessment",
          "probe": "Retain the FileOps local replacement copy or bypass shared sync. The Record failure-composition witness must turn red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/assessment",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-C6",
      "tickets": [
        "13-migrate-repair-documents.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/repairpilot"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/repairpilot",
          "probe": "Restore the local file-only sequence or weaken the short-write assertion. The corresponding real Store.Replace witness must turn red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/repairpilot",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "D-D",
      "tickets": [
        "14-close-replacement-ownership.md"
      ],
      "verification": [
        {
          "id": "adoption-tests",
          "command": "bench test --package ./internal/adopt"
        },
        {
          "id": "coverage",
          "command": "bench coverage --check durable-caller-migration"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/adopt",
          "probe": "Alter the surviving adoption stage bytes or its lifecycle transaction result. The existing real adoption assertion must fail, then pass after restoration. The source algorithm census remains independently review-owned."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ordinary-integration",
      "command": "bench test --package ./..."
    },
    {
      "id": "coverage",
      "command": "bench coverage --check durable-caller-migration"
    },
    {
      "id": "ticket-grammar",
      "command": "bench test --check ticket-grammar"
    }
  ]
}
```

## Testing decisions

Test purpose takes precedence over test count.
The prerequisite supplies the owner ordering, real-file visibility, and native qualification evidence.
This successor consumes that completed contract and proves its real caller composition.
Injected caller failures prove exact states before and after rename.

Drive each caller through its real encoder and owner call.
Inject a directory-sync failure through that call and inspect the caller error.
For each caller family, bypass the shared owner with its former sequence.
That mutation must turn its named caller row red.
An owner-only failure test does not prove caller adoption.

Preserve existing assertions and sentinel categories.
Use old-to-new differential fixtures for each enumerated caller payload family.
Keep independently authored payload expectations only when a named omission produces recorded red evidence.
Do not derive an expected payload from the new encoder under test.

Prior art: `internal/repairpilot/command_test.go:107` exercises create, short-write, and rename failures through Store.Replace.
Prior art: `internal/assessment/record_test.go:184` exercises Record write and rename failures.
Prior art: `internal/handoffdoc/store_test.go:19` exercises competing writers through WriteSection.
Read enforcement: `internal/conformance/injected_ports_registry_test.go:6` and `internal/conformance/injected_ports_test.go:40`.
The dev gate unit phase runs the affected Go packages.
The whole-project gate remains the acceptance oracle.

### Seam diagram

    caller validation and encoder
        |
        v
    payload + mode + authorized path --> durablefile.Replace --> visible file or classified error
                                             ^
                                             |
                                    injected filesystem operations

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| D16 | 16 | The gate verdict caller exposes an injected directory-sync failure. | planned internal/gate/durable_replacement_test.go (TestVerdictReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D17 | 17 | The gate evidence and lane caller exposes an injected directory-sync failure. | planned internal/gate/durable_replacement_test.go (TestEvidenceReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D18 | 18 | The prospective owner record caller exposes an injected directory-sync failure. | planned internal/gate/prospectiveartifact/replacement_test.go (TestOwnerReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D20 | 20 | The capture transaction files caller exposes an injected directory-sync failure. | planned internal/capturetx/replacement_test.go (TestCaptureReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D21 | 21 | The handoff document caller exposes an injected directory-sync failure. | planned internal/handoffdoc/replacement_test.go (TestHandoffReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D22 | 22 | The intent ledger caller exposes an injected directory-sync failure. | planned internal/intent/replacement_test.go (TestIntentReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D23 | 23 | The commitment repository caller exposes an injected directory-sync failure. | planned internal/commitment/repository/replacement_test.go (TestCommitmentReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D24 | 24 | The publication record caller exposes an injected directory-sync failure. | planned internal/publication/replacement_test.go (TestPublicationReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D25 | 25 | The broker manifest caller exposes an injected directory-sync failure. | planned internal/brokermanifest/replacement_test.go (TestBrokerReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D26 | 26 | The dashboard artifact caller exposes an injected directory-sync failure. | planned internal/dashboard/replacement_test.go (TestDashboardReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D27 | 27 | The probe subject caller exposes an injected directory-sync failure. | planned internal/probe/replacement_test.go (TestProbeReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D28 | 28 | The assessment record caller exposes an injected directory-sync failure. | planned internal/assessment/replacement_test.go (TestAssessmentReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D29 | 29 | The repair pilot document caller exposes an injected directory-sync failure. | planned internal/repairpilot/replacement_test.go (TestPilotReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D30 | 30 | Each migrated caller retains its pre-write validation and sentinel category. | planned caller tests from the migration table | A refused payload or destination cannot trigger the replacement observer. |
| D31 | 31 | Migrated callers contain no independent temporary-write and rename algorithm. | review-owned caller inventory and implementation diff | The reviewer rejects a surviving algorithm even when its happy path passes. |
| D33 | 33 | Excluded protocols retain their current assertions and operation owners. | review-owned exclusion inventory and focused existing tests | An indiscriminate migration changes the excluded protocol or its assertions. |
| D34 | 34 | A failed probe restore retains its preserved copy. | planned internal/probe/replacement_test.go (TestProbeFailedRestore) | A post-rename error cannot authorize deletion of the recovery copy. |
| D35 | 35 | The unused adopt writeManifest helper no longer defines a replacement algorithm. | review-owned internal/adopt/manifest.go caller census | Keeping the unused implementation fails the source inspection. |
| D36 | 36 | Differential fixtures preserve each migrated caller payload, mode, and pre-rename refusal result. | planned caller tests from the migration table | A helper-only test cannot detect a changed caller encoder or mode. |
| D37 | 37 | Pending persistence failure prints its cause and prevents oracle execution. | planned internal/gate/durable_replacement_test.go (TestGatePendingDurabilityFailure) | A caller that discards the leaf error loses the visibility phrase or starts the oracle. |
| D38 | 38 | Final persistence failure retains its diagnostic cause at ActionExit 1. | planned internal/gate/durable_replacement_test.go (TestGateTerminalDurabilityFailure) | The real green and red branches must expose the injected post-rename cause. |
| D39 | 39 | Green-evidence persistence failure prints the classified cause at ActionExit 1. | planned internal/gate/durable_replacement_test.go (TestGateEvidenceDurabilityFailure) | A narrow writer test misses the caller that replaces the error with generic text. |
| D40 | 40 | Failed pending recovery reports both persistence errors. | planned internal/gate/durable_replacement_test.go (TestGatePendingRecoveryFailure) | Independent final-sync and recovery-rename faults expose a discarded recovery cause. |
| D41 | 41 | A post-rename mutation failure with failed restore retains the preserved copy at exit 2. | planned internal/probe/replacement_test.go (TestProbeMutationRecoveryFailure) | The existing unconditional release deletes the byte-exact recovery copy. |
| D42 | 42 | A post-rename mutation failure releases preservation only after verified restore. | planned internal/probe/replacement_test.go (TestProbeMutationRecoverySuccess) | A successful real restore must precede copy release and an exit 1 failure. |
| D43 | 43 | A pre-rename mutation failure retains the original subject and existing cleanup result. | planned internal/probe/replacement_test.go (TestProbeMutationPreparationFailure) | The old-byte and home-state witnesses preserve the existing refusal contract. |
| D44 | 44 | A post-rename ledger failure skips compensation. | planned internal/intent/replacement_test.go (TestTransactionPublishedFailure) | A compensation counter and new ledger bytes detect the former unconditional rollback. |
| D45 | 45 | A pre-rename ledger failure compensates exactly once. | `internal/intent/transaction_test.go` (`TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes`) | The existing test compares old bytes and the compensation count under the real lock. |
| D46 | 46 | An initial post-rename policy failure restores the original policy state before approval returns. | planned internal/commitment/repository/replacement_test.go (TestApprovalInitialPolicyFailure) | The initial policy rename succeeds before sync fails, which reaches the previously unfurnished restore path. |
| D47 | 47 | A post-rename board failure restores both original file states. | planned internal/commitment/repository/replacement_test.go (TestApprovalBoardFailure) | Distinct original bytes detect a restore that handles only one file. |
| D48 | 48 | A pre-rename approval-ledger failure restores the original policy and board. | planned internal/commitment/repository/replacement_test.go (TestApprovalUnpublishedLedgerFailure) | The real ledger transaction must run the stage restore exactly once. |
| D49 | 49 | A post-rename approval-ledger failure retains the new policy and board with its approved receipt. | planned internal/commitment/repository/replacement_test.go (TestApprovalPublishedLedgerFailure) | Read-back of all three artifacts detects an approved receipt paired with old external files. |
| D50 | 50 | A failed stage or compensation restore reports its cause and unrestored path. | planned internal/commitment/repository/replacement_test.go (TestApprovalRestoreFailure) | A second independent fault detects discarded snapshot-restore errors. |
| D51 | 51 | Failed initial policy stage in a board-absent repository restores policy without creating a board. | planned internal/commitment/repository/replacement_test.go (TestApprovalFailureWithoutBoard) | The fixture verifies original policy or absence and an unapproved receipt with no board. |
| D52 | 52 | A post-rename reauthorization failure retains the next request and its matching external effect. | planned internal/intent/replacement_test.go (TestReauthorizePublishedFailure) | The real transition witness detects rollback of an effect whose ledger swap is already visible. |

### Edge inventory

The audience includes this kit and repositories that receive its executable.
An absent destination and an empty payload have separate rows D02 and D10.
Spaces, embedded quotes, glob characters, and Unicode paths belong to the real-file fixtures.
Caller tests retain special-file, symlink, payload-size, and permission refusal cases through D30.

Failure fixtures cover every stage named by the owner, including directory close.
A failing cleanup returns both error causes through D15.
A post-rename failure retains the new bytes through D12.
The caller lock remains held until its operation and permitted compensation return.

D37-D40 drive the complete gate execution path.
D41-D43 distinguish failed mutation from failed restore.
D44-D52 drive the compound ledger and external-file outcomes.

**Won't handle:** untrusted parent races — the surviving gate caller owns destination authority before replacement.
**Won't handle:** directory creation durability — the surviving review writer owns parent creation before replacement.
**Won't handle:** unsupported or remote filesystems — native qualification covers the supported Linux and macOS environments.
**Won't handle:** permanent storage failure recovery — the surviving probe caller retains its preservation copy.

| excluded operation | current owner | reason and surviving in-scope caller |
|---|---|---|
| Exclusive immutable creation | capturetx.writeNew and chargeevidence.Staged | Exclusive creation and hard-link publication differ from replacement. capturetx.replace migrates. |
| Directory exchange and promotion | releaseevidence/evidence_promotion.go and capturetx/public.go | Whole-directory transactions retain their protocols. Gate evidence migrates. |
| Adoption transaction images | adopt/transaction/publication.go | Identity checks, symlinks, and undo require staged publication. Broker manifests migrate. |
| Executable and seal publication | freshness/freshness_publish.go | Pair rollback retains its transaction. Broker-manifest file writes migrate. |
| Record rotation | otelrecord/writer.go | Rotation retains its writer lifecycle. Publication records migrate. |
| Lease and lock records | worktree/lifecycle.go, worktree/ownership.go, intent.acquire | Exclusion and ownership protocols retain their writers. Intent ledger files migrate. |
| Cancellable reference publication | skillsindex/skillsindex.go:316 | The pre-rename cancellation barrier retains its protocol. Review records migrate. |
| Git tree construction | landing/spectree.go | Object construction does not replace a working file. Probe subjects migrate. |



The prerequisite already supplies D02, D10, D12, D15, and D32 evidence.
References to those rows consume its completed contract without transferring their ownership here.
The global D31 census includes the landed prerequisite and every successor caller.
D36 consumes the prerequisite's review-record differential evidence beside each successor family's evidence.

## Ownership fences

These exact fences bound future implementation writes.
Review pickup is phase-owned.
The implementation ticket Writes union must match this union after reviewed slicing.

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/adopt/manifest.go`
- `internal/assessment/record_test.go`
- `internal/assessment/replacement_test.go`
- `internal/assessment/store.go`
- `internal/brokermanifest/brokermanifest.go`
- `internal/brokermanifest/replacement_test.go`
- `internal/capturetx/replacement_test.go`
- `internal/capturetx/transaction.go`
- `internal/commitment/repository/replacement.go`
- `internal/commitment/repository/replacement_test.go`
- `internal/commitment/repository/repository.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/injected_ports_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/dashboard/dashboard.go`
- `internal/dashboard/replacement_test.go`
- `internal/gate/durable_replacement_test.go`
- `internal/gate/engine.go`
- `internal/gate/prospectiveartifact/prospectiveartifact.go`
- `internal/gate/prospectiveartifact/replacement_test.go`
- `internal/gate/run_durability.go`
- `internal/gate/run_transaction.go`
- `internal/gate/verdict.go`
- `internal/handoffdoc/replacement_test.go`
- `internal/handoffdoc/store.go`
- `internal/intent/assignment.go`
- `internal/intent/intent.go`
- `internal/intent/replacement_test.go`
- `internal/intent/transaction.go`
- `internal/probe/probe.go`
- `internal/probe/replacement_test.go`
- `internal/probe/subject.go`
- `internal/publication/record.go`
- `internal/publication/replacement_test.go`
- `internal/repairpilot/command.go`
- `internal/repairpilot/command_test.go`
- `internal/repairpilot/replacement_test.go`
- `reviews/durable-caller-migration.md`
- `tests/canary/package-core-guard/bounds-intent-window-fixed`
- `tests/canary/package-core-guard/git-flag-retyped`

## Out of scope

No production implementation, test execution, benchmark, ticket slicing, or roadmap edit belongs to this author phase.
FT354 remains open for the strict-JSON outcome and the shell-argument outcome.
This spec has no Roadmap retirement field.

A cancellable replacement API is a separate capability: 3 edits, 1 gate run.
A new multi-file crash-recovery protocol is a separate capability: 6 edits, 2 gate runs.
The caller error outcomes above are in scope for the existing transaction protocol.
A remote-filesystem durability policy is a separate capability: 4 edits, 2 gate runs.
These estimates name future planning units, not implementation promises.


## Further notes

### Source-sentence coverage

This table reconciles the original complete source contract.
D01-D15, D19, and D32 refer to completed prerequisite evidence.
The remaining row IDs name this successor's unchanged owned requirements.

| reviewed source clause | coverage |
|---|---|
| One leaf owner writes an adjacent temporary regular file. | D01-D05, D31 |
| Apply mode, sync file, close, rename, and sync directory. | D06-D08 |
| Keep caller mode, validation, authority, and sentinel translation. | D16-D30, D36 |
| Failed synchronization is failure. | D09-D13 and each caller row |
| After rename, new bytes may already be visible. | D11-D12, D34 |
| No universal old-byte promise. | D09-D12, D37-D52 |
| Strengthen review-record directory durability explicitly. | D19 and the behavior-change table |
| Keep directory exchange, locks, and evidence promotion separate. | D33 and the exclusion inventory |
| Inject failures before and after rename. | D09-D15, D37-D52 |
| Qualify supported platform behavior before acceptance. | D32 |
| Split primitives and attach caller bypass mutations. | D16-D29 and the chunk table |

### Pre-review proof checklist

- Cited symbols: the definition reads and static consumer census resolve the named owner and caller symbols.
- Import edges: the current packages resolve through go list. The new leaf import remains a build-time validation obligation.
- Source-row clauses and occurrences: the source-sentence table consumes tickets 2 and 4 without moving their parent map.
- Promised field labels: Stage, cause, and completed rename belong to the proposed error contract, not a persisted schema.
- Changed-function callers: the inventory includes complete gate, probe, approval, ledger, and reauthorization transaction paths.
- Compensation callers: commitment approval and assignment reauthorization retain effects only after a completed ledger rename.
- Copy survival: D31 and D35 grade remaining algorithms. Each migrated family also requires its bypass mutation.
- Rendered-shape readers: none for success output. Caller diagnostics preserve context and add the post-rename uncertainty phrase.
- Pin operators: exact payload bytes, permission equality, errors.Is, and trace ordering belong to the planned tests.
- Entry reads: none introduced. Existing directory authority stays with caller entries.
- Derived expectations: caller fixture inputs supply expected bytes. The owner cannot derive its own payload oracle.
- Consolidated rules: the behavior-change and caller tables identify each old and new rule.
- Quantified obligations: every stage and migrated family receives the table-driven failure matrix.
- Workflow-step writes: none introduced.

### Evidence status and validation plan

Facts come from the pinned source tree and the named reviewed artifacts.
The caller census used bench consumers and a hidden production-Go rename and sync search.
The search also examined immutable, staged, cancellation, rotation, and lease owners.
No native durability result or byte-compatibility probe ran during authoring.

Native Linux and macOS qualification must complete in the prerequisite before this successor starts.
The spec-only checks are prose mechanics and acceptance coverage.
The proposed caller failures are static consequences, not executed reproductions.
Independent review must approve the provisional fences and chunks before ticket creation.

### Review repair dispositions

| finding | disposition | rows and fences |
|---|---|---|
| DFR-R1-GATE | Preserve causes through complete execution and failed pending recovery. | D37-D40. Gate run_transaction.go joins its exact fence. |
| DFR-R1-PROBE | Classify failed mutation before copy release and restore before a focused run. | D41-D43. Probe probe.go joins its exact fence. |
| DFR-R1-TRANSACTION | Compensate only unpublished ledger writes and restore failed initial stage. | D44-D52. Intent transaction.go and assignment.go join exact fences. |

The repair retains every existing successful and pre-rename failure assertion.
The new typed error determines the new caller outcomes.
No runtime probe or production edit supports these prospective cases yet.
A confirming independent review must accept the repaired spec before ticket creation.


### Delivery partition

This successor owns exactly D16-D18, D20-D31, and D33-D52: 35 original rows.
The prerequisite owns exactly D01-D15, D19, and D32: 17 original rows.
No row changes its accepted predicate, seam, omission witness, or ID.
D31, D33, and D36 reconcile both completed capabilities without duplicating prerequisite implementation ownership.

### Snapshot and delivery basis

The original accepted spec SHA256 is 354285d94746b5584e132159fc87acc545d995d6d81abb384541a6db63552bab.
Snapshot commit 8e3f3c0b18b9af191cbe1d59f3e5c5f16da69f45 preserves the full spec and all fourteen pending ticket drafts.
Those ticket files are historical drafts, not an approved graph for either split spec.
Both split specs received independent acceptance before this graph was authored.
Independent ticket review accepted this graph before spec-stage close.

The split addresses PL-R1-DELIVERY without changing accepted behavior.
Whole-spec completion, acceptance reconciliation, and landing require all planned chunks and rows.
Sources: internal/reviewrecord/coverage.go:132, internal/reviewrecord/check.go:64, and the independent delivery review.
The shared parent decision map stays in place.

### Ticket fence closure additions

| added exact path | stable chunk | source and preserved contract |
|---|---|---|
| tests/canary/package-core-guard/bounds-intent-window-fixed | D-B5 | Fixture BASE pins intent.go. Keep its VerdictWindow-versus-FixedWindow omission policy. |
| tests/canary/package-core-guard/git-flag-retyped | D-C3 | Fixture BASE pins dashboard.go. Keep its Git administration flag ownership policy. |

These paths expand the former fence to its existing fixture closure.
The additions preserve EXPECT and MUTATE.json and change no accepted product behavior.
The bound package registry also adds its five canonical command-registry files to the stable chunk union.


### Source headroom

The gate execution source has 398 lines, and repository.go has 400 at the pinned source.
D-B1 co-owns run_durability.go for cohesive persistence-failure and recovery helpers.
D-B6 co-owns replacement.go for snapshot replacement and restore.
Each extraction lands with its complete caller outcome and preserves the accepted interface.
The existing registry and fixture files receive no unrelated growth or policy rewrite.

## Split approval

| item | proposed contract | spec reviewer disposition |
|---|---|---|
| Implementation line | gpt-5.6-sol / high, D-B1, D-B5, D-B6, and D-C4 harder | Accepted at frozen split hash |
| Seams | Real caller encoding and transaction, recovery, or replacement entry | Accepted at frozen split hash |
| Acceptance and edges | Thirty-five unchanged original rows and full exclusion census | Accepted at frozen split hash |
| Ownership fences | Remaining exact caller, registry, fixture, and review-pickup union | Accepted at frozen split hash |
| Scope cuts | Existing specialized protocols and future crash recovery remain excluded | Accepted at frozen split hash |
| Delivery | Complete prerequisite acceptance and landing before every caller chunk | Accepted at frozen split hash |

## Ticket approval

Author line: gpt-6.1-sol / high / one reslicing pass plus at most two bounded repairs.
Accepted split SHA256: 32b04d2f6d1ce773170d1641eb3408a998d84656104c9fbf9702474f92c2ea66
Independent ticket review accepted the frozen graph before this spec-stage close.
No version 1 plan authorizes dispatch.

| numbered ticket | Blocked by | delivered outcome |
|---|---|---|
| 1. 02-preserve-gate-failures.md | none | Preserve durability failures through gate execution |
| 2. 03-migrate-owner-records.md | none | Migrate prospective owner-record replacement |
| 3. 04-migrate-capture-documents.md | none | Migrate capture transaction file replacement |
| 4. 05-migrate-handoff-documents.md | none | Migrate handoff document replacement |
| 5. 06-preserve-ledger-effects.md | none | Preserve matching effects after ledger publication |
| 6. 07-recover-approval-stage.md | 06-preserve-ledger-effects.md | Recover failed approval stages and retain published receipts |
| 7. 08-migrate-publication-records.md | none | Migrate publication record replacement |
| 8. 09-migrate-broker-manifests.md | none | Migrate broker manifest replacement |
| 9. 10-migrate-dashboard-artifacts.md | none | Migrate dashboard artifact replacement |
| 10. 11-preserve-probe-recovery.md | 07-recover-approval-stage.md | Preserve probe recovery after failed mutation |
| 11. 12-migrate-assessment-records.md | 11-preserve-probe-recovery.md, 08-migrate-publication-records.md | Migrate assessment record replacement |
| 12. 13-migrate-repair-documents.md | none | Migrate repair pilot document replacement |
| 13. 14-close-replacement-ownership.md | 02-preserve-gate-failures.md, 03-migrate-owner-records.md, 04-migrate-capture-documents.md, 05-migrate-handoff-documents.md, 06-preserve-ledger-effects.md, 07-recover-approval-stage.md, 08-migrate-publication-records.md, 09-migrate-broker-manifests.md, 10-migrate-dashboard-artifacts.md, 11-preserve-probe-recovery.md, 12-migrate-assessment-records.md, 13-migrate-repair-documents.md | Remove the unused replacement owner and close caller coverage |

All thirteen tickets also depend on complete prerequisite acceptance and landing.
The sibling blocker field records only this successor's local ordering.
D31 and D35 require independent source-census review, not an unrelated package-test probe.
D30 and D36 include the prerequisite's recorded review-record evidence without reopening its implementation fence.

### Accepted ticket-review source

Accepted graph commit: fb38ea93ea98cd273bc3865f6e82af458711a7f2
Accepted spec SHA256 before closure metadata: 2d606f650bb99d3c6daa3d040f504f432210d7e0cdf5f428b53ddbb1878be1ee
Reviewer: /root/primitive_specs_review
Reviewer line: gpt-6.1-sol / high
Judgment: Standards 0 blockers, Spec 0 blockers, Coverage 0 blockers
Review-axis confidence: 10 each, static planning only

The version 1 plan remains an authored implementation plan, not completed implementation evidence.
No runtime test, native qualification, manual gate, or benchmark supports this stage.
Ticket bytes, acceptance rows, ownership fences, and implementation decisions remain at the accepted source.
