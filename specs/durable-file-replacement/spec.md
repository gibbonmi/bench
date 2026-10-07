# Durable file replacement

Status: staged

Decision source: named reviewed artifact `decisions/architecture-primitives.md`, resolved tickets `2.md` and `4.md` (2026-10-06).

Verification log: 2 original-spec iterations, 1 split-spec iteration, and 1 ticket iteration to accept — separate spec acceptance preceded reslicing.

## Problem

File replacement has separate owners with different synchronization steps.
The review-record writer omits parent-directory synchronization.
A leaf capability must complete and land before its independent process-policy consumer can use it.
Sources: internal/reviewrecord/write.go:92 and the named PL-R1-DELIVERY review.

## Solution

Own replacement of one regular file through one dependency leaf.
Deliver the leaf with its real review-record caller and native qualification.
Keep payload, mode, validation, destination authority, and sentinel translation with the caller.
Treat synchronization failure as failure, including uncertainty after rename.
The remaining callers migrate through the separately deliverable durable-caller-migration successor.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: failure ordering and real caller composition require stronger evidence than helper happy paths.
Author line: gpt-6.1-sol / high / one split pass plus at most two bounded repairs.
Harder chunks: D-A.

1. As a writer, I want to retain the exact replacement bytes, so that readers receive the complete payload.
2. As a writer, I want to create an empty file, so that an empty payload stays distinct from absence.
3. As a caller, I want to select the file mode, so that the owner keeps caller permission policy.
4. As a reader, I want to observe whole files, so that concurrent reads cannot see partial payloads.
5. As a writer, I want to use an adjacent temporary file, so that rename stays within the destination directory.
6. As a writer, I want to synchronize before rename, so that publication follows the complete file write.
7. As a writer, I want to synchronize after rename, so that success includes the directory step.
8. As a writer, I want to detect a short write, so that success cannot hide a partial payload.
9. As an operator, I want to retain old bytes before rename, so that a refused preparation does not publish.
10. As an operator, I want to retain absence before rename, so that a refused first write creates no destination.
11. As an operator, I want to receive post-rename uncertainty, so that a failed sync cannot imply rollback.
12. As an operator, I want to observe new bytes after rename, so that failure evidence states the visible result.
13. As a caller, I want to retain the underlying error, so that sentinel checks remain available.
14. As an operator, I want to remove temporary residue, so that a refused preparation leaves no stray file.
15. As an operator, I want to receive cleanup failures, so that a leak cannot disappear from the error.
19. As a review record caller, I want to use the shared replacement owner, so that the caller receives the declared durability result.
32. As an operator, I want to qualify native durability, so that the guarantee reflects supported host behavior.

## Implementation decisions

The new package is `internal/durablefile`.
Its production entry is `Replace(path string, data []byte, mode fs.FileMode) error`.
The parent directory must already exist.
The owner does not create a parent directory or validate the destination authority.
The caller creates directories through its current policy before it calls the owner.

The owner creates a unique temporary regular file beside the destination.
It applies the requested mode before file synchronization.
It writes the complete payload, synchronizes the file, closes it, renames it, and synchronizes the parent directory.
It closes the directory before it returns success.
A short write returns `io.ErrShortWrite`.
The owner removes an uninstalled temporary file on failure.

The typed error is *durablefile.Error, with Stage string and Published bool.
Its underlying cause is available through Unwrap() error.
Published reports whether rename completed.
The stage values name create, mode, write, file-sync, file-close, rename, directory-open, directory-sync, and directory-close.
After rename, its message includes `new bytes may already be visible`.

Cleanup failures join the primary cause.
The owner promises no automatic rollback.

Tests use the same replacement operation through a `Replacer` value with injected file operations.
`Replacer.Replace` has the production entry signature.
The default entry uses the operating-system adapter.
The operation ports cover temporary creation, file operations, rename, directory open, directory operations, and removal.
No caller copies the ordered sequence.
The injected file interface carries Write, Chmod, Sync, Close, and Name.

The review-record entry supplies the native replacer to its internal write transaction.
Tests supply the replacement operation through that same transaction.
Keep record lookup, change, rendering, bounded grading, parse validation, and parent creation in the caller.
The replace wrapper supplies 0644 and returns the leaf error without losing errors.Is or errors.As.
After rename, failure must retain the phrase `new bytes may already be visible`.
The caller promises no automatic rollback after that failure.

### Current caller inventory

The pinned production source is a382f4848d432815968f044820fad593e856e297.

| caller | source and symbol | caller mode | caller payload | row |
|---|---|---|---|---|
| review record | `internal/reviewrecord/write.go:92` (`replace`) | 0644 | Graded rendered document | D19 |

The source census for other replacement callers belongs to specs/durable-caller-migration/spec.md.
This prerequisite migrates only the review-record caller.
The leaf imports no process, resource-policy, transaction, or artifact owner.

### Explicit behavior change

Review records currently synchronize only the file.
D19 adds parent-directory synchronization and exposes its failures through the actual record-writing transaction.
Mode application precedes file synchronization.
Existing successful rendering, mode, validation, sentinel, and pre-rename refusal assertions remain required.

## Implementation chunks

D-A is the complete owner-plus-real-caller capability.
Its complete spec must pass all acceptance rows, native qualification, independent review, reconciliation, and landing.
It does not wait for any caller-migration successor chunk.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| D-A / 01-own-replacement-and-review.md | Own durable replacement through review records | D01, D02, D03, D04, D05, D06, D07, D08, D09, D10, D11, D12, D13, D14, D15, D19, D32 | owner-tests, caller-tests, linux-native, macos-native | yes |

The process-lifetime PL-C2 prerequisite is this complete accepted and landed spec.
The caller-migration successor also starts only after this complete prerequisite lands.
Ticket authoring followed independent split-spec acceptance.
The version 1 completion plan records future evidence, not observed implementation results.

### Completion plan

The version 1 plan records future implementation evidence.
It claims no current implementation pass, red, native qualification, or benchmark.
The orchestrator adds required version 2 author sessions before dispatch.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "D-A",
      "tickets": [
        "01-own-replacement-and-review.md"
      ],
      "verification": [
        {
          "id": "owner-tests",
          "command": "bench test --package ./internal/durablefile"
        },
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/reviewrecord"
        },
        {
          "id": "linux-native",
          "command": "bench test --package ./internal/durablefile --run 'TestNativeReplacement'"
        },
        {
          "id": "macos-native",
          "command": "bench test --package ./internal/durablefile --run 'TestNativeReplacement'"
        },
        {
          "id": "owner-omission-proof",
          "command": "bench test --package ./internal/durablefile",
          "probe": "Omit file sync, omit directory sync, rename before file close, or accept a short write. Each named owner assertion must fail, then pass after restoring the source."
        },
        {
          "id": "owner-preservation",
          "command": "bench test --package ./internal/durablefile",
          "probe": "Corrupt supplied payload bytes, requested mode, temporary cleanup, or wrapped error cause. The corresponding independent owner assertion must fail, then pass after source restoration."
        },
        {
          "id": "caller-bypass-proof",
          "command": "bench test --package ./internal/reviewrecord",
          "probe": "Restore the former review-record writer or bypass the shared replacement operation in the real record transaction. TestReviewReplacement must fail its directory-sync cause assertion, then pass after source restoration."
        },
        {
          "id": "caller-preservation",
          "command": "bench test --package ./internal/reviewrecord",
          "probe": "Corrupt the review-record rendered payload, 0644 mode, or pre-write refusal policy. The real caller assertion must fail, then pass after source restoration."
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
      "command": "bench coverage --check durable-file-replacement"
    }
  ]
}
```

## Testing decisions

Test purpose takes precedence over test count.
The owner trace proves mode, write, file synchronization, close, rename, directory synchronization, and directory close ordering.
Real temporary-directory fixtures prove exact bytes, modes, regular-file placement, cleanup, and before-or-after-rename visibility.
Injected failures cover every declared stage and joined cleanup causes.
Native Linux and macOS tests qualify supported host and filesystem behavior before prerequisite acceptance.
They claim no power-loss simulation or universal filesystem guarantee.

Drive the actual review-record write transaction through its encoder, bounded grading, parser, and replacement operation.
Inject a directory-sync failure and inspect the classified error, exact rendered bytes, and 0644 mode.
Pre-rename failures retain existing bytes or absence and the existing refusal category.
A review-writer bypass mutation must turn TestReviewReplacement red in the reviewrecord package.
An owner-only test cannot supply this caller evidence.
Owner omissions and caller bypasses use their respective package commands and independent red witnesses.

Preserve every existing review-record assertion and independent payload expectation.
Record a named omission red, restore the source, and record the pass before accepting independent expectations.
Prior art: internal/reviewrecord/write.go:44 and internal/repairpilot/command_test.go:107.
The dev gate grades affected packages, and the whole-project gate remains the acceptance oracle.

### Seam diagram

    record change -> rendering -> bounded grading -> parser -> authorized path + bytes + 0644
                                                                  |
                                                                  v
                                                        durablefile.Replace
                                                                  ^
                                                        injected file operations

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| D01 | 1 | Replacement of old bytes installs the exact supplied bytes. | planned internal/durablefile/replace_test.go (TestReplaceBytes) | A truncated write differs from the supplied payload. |
| D02 | 2 | An empty payload creates a present regular file. | planned internal/durablefile/replace_test.go (TestReplaceEmpty) | Skipping an empty write leaves the destination absent. |
| D03 | 3 | The installed file has the requested permission bits. | planned internal/durablefile/replace_test.go (TestReplaceMode) | The default temporary mode cannot satisfy the executable fixture. |
| D04 | 4 | A concurrent reader observes only the old or new complete payload. | planned internal/durablefile/replace_test.go (TestReplaceVisibility) | An in-place write exposes a partial payload. |
| D05 | 5 | The temporary regular file belongs to the destination parent. | planned internal/durablefile/replace_test.go (TestAdjacentTemporary) | A temporary file in the system directory fails the recorded path predicate. |
| D06 | 6 | Rename follows the successful mode, write, file-sync, and file-close operations. | planned internal/durablefile/replace_test.go (TestReplaceOrder) | An operation trace detects an omitted or reordered step. |
| D07 | 7 | Success follows parent-directory synchronization and close. | planned internal/durablefile/replace_test.go (TestReplaceDirectorySync) | An omitted directory step cannot produce the required trace. |
| D08 | 8 | A short write returns an error before rename. | planned internal/durablefile/replace_test.go (TestReplaceShortWrite) | A nil-error partial write cannot pass as complete. |
| D09 | 9 | Each pre-rename failure leaves the existing destination bytes unchanged. | planned internal/durablefile/replace_test.go (TestReplaceBeforeRenameFailures) | The old-byte oracle detects an early destination write. |
| D10 | 10 | Each pre-rename failure leaves an absent destination absent. | planned internal/durablefile/replace_test.go (TestReplaceAbsentFailures) | An eager destination create violates the absence predicate. |
| D11 | 11 | Each post-rename failure reports that new bytes may already be visible. | planned internal/durablefile/replace_test.go (TestReplaceAfterRenameFailures) | A plain successful return or old-byte claim fails the error predicate. |
| D12 | 12 | A post-rename sync failure leaves the new destination bytes visible. | planned internal/durablefile/replace_test.go (TestReplaceAfterRenameState) | A fake that fails before rename cannot satisfy the state assertion. |
| D13 | 13 | Replacement errors preserve errors.Is access to the injected cause. | planned internal/durablefile/replace_test.go (TestReplaceErrorCause) | A text-only wrapper loses the injected sentinel. |
| D14 | 14 | A pre-rename failure removes its temporary file when removal succeeds. | planned internal/durablefile/replace_test.go (TestReplaceTemporaryCleanup) | A planted write failure detects a leaked sibling file. |
| D15 | 15 | A failed temporary cleanup preserves the primary error and the cleanup cause. | planned internal/durablefile/replace_test.go (TestReplaceCleanupFailure) | Discarding the removal error loses the cleanup sentinel. |
| D19 | 19 | The review record caller exposes an injected directory-sync failure. | planned internal/reviewrecord/replacement_test.go (TestReviewReplacement) | Replacing the owner call with the former rename sequence hides the injected failure. |
| D32 | 32 | Native Linux and macOS qualification establishes the declared synchronization result. | planned internal/durablefile/platform_test.go (TestNativeReplacement) | A simulated filesystem cannot establish native directory-sync support. |

### Edge inventory

An absent destination and an empty payload have separate rows D02 and D10.
Real-file fixtures include spaces, embedded quotes, glob characters, and Unicode paths.
Every owner stage, including directory close, has an injected failure witness.
Cleanup failure keeps both causes through D15.
Post-rename failure retains visible new bytes through D12.
Review-record fixtures preserve existing payload, mode, validation, authority, and sentinel policy through D19.

**Won't handle:** untrusted parent races — the review-record caller retains destination authority.
**Won't handle:** directory creation durability — the review-record caller retains parent creation.
**Won't handle:** unsupported or remote filesystems — D32 qualifies the supported Linux and macOS environments.
**Won't handle:** other replacement callers — durable-caller-migration owns their original accepted obligations.
**Won't handle:** specialized publication or lock protocols — the successor retains their named exclusion inventory without migrating them.

## Ownership fences

These exact fences bound future implementation writes.
Review pickup is phase-owned.
The implementation ticket Writes union must match this union after reviewed slicing.

- `internal/durablefile`
- `internal/reviewrecord/replacement_test.go`
- `internal/reviewrecord/write.go`
- `reviews/durable-file-replacement.md`

## Out of scope

No production implementation, test execution, benchmark, ticket slicing, or roadmap edit belongs to this split-author phase.
FT354 remains open for caller migration, strict JSON, and shell-argument outcomes.
This spec has no Roadmap retirement field.

A cancellable replacement API is a separate capability: 3 edits, 1 gate run.
A new multi-file crash-recovery protocol is a separate capability: 6 edits, 2 gate runs.
A remote-filesystem durability policy is a separate capability: 4 edits, 2 gate runs.

## Further notes

### Source-sentence coverage

| reviewed source clause | owned coverage |
|---|---|
| One leaf owner writes an adjacent temporary regular file. | D01-D05 |
| Apply mode, sync file, close, rename, and sync directory. | D06-D08 |
| Preserve failure state, cause, cleanup, and post-rename uncertainty. | D09-D15 |
| Strengthen review-record directory durability explicitly. | D19 |
| Qualify supported native platform behavior before acceptance. | D32 |

The successor owns D16-D18, D20-D31, and D33-D52 unchanged.
This prerequisite owns D01-D15, D19, and D32 unchanged.
Their disjoint union is the original 52-row acceptance map.

### Evidence status

Facts come from the pinned tree, named decision artifact, and accepted original spec.
No native test, runtime reproduction, production edit, or benchmark ran during authoring.
Prose and coverage checks validate this planning artifact.
Independent review accepted this split before this ticket graph was authored.

### Delivery partition

The original accepted spec SHA256 is 354285d94746b5584e132159fc87acc545d995d6d81abb384541a6db63552bab.
Snapshot commit 8e3f3c0b18b9af191cbe1d59f3e5c5f16da69f45 preserves the full spec and all fourteen pending ticket drafts.
Those ticket files are historical drafts, not an approved graph for either split spec.
Both split specs received independent acceptance before this graph was authored.
Independent ticket review accepted this graph before spec-stage close.

The split addresses PL-R1-DELIVERY without changing accepted behavior.
Whole-spec completion, acceptance reconciliation, and landing require all planned chunks and rows.
Sources: internal/reviewrecord/coverage.go:132, internal/reviewrecord/check.go:64, and the independent delivery review.
The shared parent decision map stays in place.

## Split approval

| item | proposed contract | spec reviewer disposition |
|---|---|---|
| Implementation line | gpt-5.6-sol / high, D-A harder | Accepted at frozen split hash |
| Seam | Real record write through the leaf, with injected operations | Accepted at frozen split hash |
| Acceptance and edges | Seventeen unchanged original rows and native qualification | Accepted at frozen split hash |
| Ownership fences | Leaf, exact review writer and fixture, phase review pickup | Accepted at frozen split hash |
| Scope cuts | Other callers belong to the successor, specialized protocols remain excluded | Accepted at frozen split hash |
| Delivery | Complete prerequisite lands before PL-C2 or any successor caller chunk | Accepted at frozen split hash |

## Ticket approval

Author line: gpt-6.1-sol / high / one reslicing pass plus at most two bounded repairs.
Accepted split SHA256: 05c79801d3746e5250fb23d7261f66827f971e6d26c0f9661ec7aaff48a0c2f8
Independent ticket review accepted the frozen graph before this spec-stage close.
No version 1 plan authorizes dispatch.

| numbered ticket | Blocked by | delivered outcome |
|---|---|---|
| 1. 01-own-replacement-and-review.md | none | Own durable replacement through review records |

### Accepted ticket-review source

Accepted graph commit: fb38ea93ea98cd273bc3865f6e82af458711a7f2
Accepted spec SHA256 before closure metadata: 4eacd9370ea59bf7003f897aa14e3e2702382b393b3db115727c094e61243333
Reviewer: /root/primitive_specs_review
Reviewer line: gpt-6.1-sol / high
Judgment: Standards 0 blockers, Spec 0 blockers, Coverage 0 blockers
Review-axis confidence: 10 each, static planning only

The version 1 plan remains an authored implementation plan, not completed implementation evidence.
No runtime test, native qualification, manual gate, or benchmark supports this stage.
Ticket bytes, acceptance rows, ownership fences, and implementation decisions remain at the accepted source.
