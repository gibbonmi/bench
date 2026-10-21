# Simplify guard inventory scanning with a synchronous loop

Status: staged
Decision source: `decisions/architecture-guards.md`, answer 3, reviewed 2026-10-06
Subject: `a9c395e77fec36d1f60f6d057c654aecf5720e24` with production source `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`
Audience: every repository that links the kit
Related roadmap: FT366
Verification log: independent SPEC acceptance precedes slicing. The high reviewer accepted the decoder/source repairs at 02c799aa with zero findings and confidence 10. The xhigh placement confirmer accepted 6c13b4b3 with zero findings and confidence 9.
Independent graph review accepted aa86c64d on all three axes with zero findings and confidence 9.

## Problem

`internal/guards/guards.go` (`Scan`) launches enumeration and inspection workers.
Cancellation then waits for those workers to finish before Scan returns.
The channels and selects therefore add indirection without a hard deadline for blocking filesystem work.
`TestScanWaitsForCancelledWorkerCleanup` already requires cleanup before return.

## Solution

Run the same enumeration and inspection operations in a synchronous, context-aware loop.
Preserve completed rows, honest counts, static manifest reads, command output, and return codes.
State cancellation as a checkpoint between operations.
Do not promise to interrupt an in-flight filesystem call.

## User stories

Line: `gpt-5.6-sol` / high
Implementation-line reason: the scan shape is known, but cancellation boundaries and exact output preservation require careful work at existing injected seams
Harder chunks: GS1, use xhigh effort on the declared model
Author line: `gpt-6.1-sol` / high, one draft and at most two bounded review repairs

1. As an agent, I want complete scans preserved, so that the guard inventory keeps its current meaning.
2. As an agent, I want completed rows retained on cancellation, so that I can use the inspected prefix.
3. As an agent, I want honest omitted counts, so that an incomplete scan does not claim completeness.
4. As an agent, I want unknown enumeration counts retained, so that an interrupted inventory does not fabricate totals.
5. As an agent, I want cleanup before return, so that the next command cannot race unfinished scan work.
6. As an agent, I want cancellation boundaries stated honestly, so that a filesystem wait is not advertised as a hard deadline.
7. As an agent, I want absent and empty hook directories handled, so that linked repositories retain their supported inventory states.
8. As an agent, I want static header reads retained, so that inventory cannot execute a hook.
9. As an agent, I want wiring facts retained, so that Claude and Codex rows report the same configuration.
10. As an agent, I want complete metadata retained, so that the command keeps its machine-readable contract.
11. As an agent, I want incomplete metadata retained, so that a timeout remains visible in full and brief output.
12. As an agent, I want action help retained, so that stale and unwired guards keep their existing next commands.
13. As an agent, I want incomplete scans to suppress action help, so that partial data cannot recommend a misleading repair.
14. As an agent, I want command errors retained, so that invalid arguments and missing repositories remain distinguishable.
15. As a maintainer, I want the old private worker loop rejected, so that simplification cannot leave a second scan implementation.
16. As a maintainer, I want binding closure retained, so that an internal refactor preserves command registry checks.
17. As a maintainer, I want deterministic before-and-after evidence, so that unchanged behavior is demonstrated at the real command seam.

## Implementation decisions

### Loop and cancellation contract

Keep Scan(ctx, root), Rows(root), Command, ScanResult, Row, and their exported signatures unchanged.
Keep enumerateGuards and inspectGuard as the existing package-variable test seams.
Do not add an injected parameter or a new interface.
The concrete loop performs enumeration once, then inspects candidates in their current order.
GS01, GS02, and GS22 observe these contracts.

Check ctx.Err before enumeration and after it returns.
A pre-cancelled context starts no operation.
Cancellation observed after enumeration returns an incomplete result with unknown total and omitted counts.
It does not publish the returned candidates as a completed inventory.
GS03 and GS04 distinguish these cases.

Check ctx.Err before each inspection and after it returns.
Count an inspection only when the post-operation check confirms that cancellation has not occurred.
An inspection that returns no rows still counts as inspected.
An interrupted inspection contributes no rows and does not increment Inspected.
Already completed rows remain in order.
GS05 through GS09 observe those boundaries.

The current select can choose either ready branch when completion and cancellation race.
The loop uses its explicit post-operation context check in that simultaneous case.
This tightens that race by withholding a candidate when cancellation is observed at the checkpoint.
Preservation evidence fixes event order and does not claim an identical winner for a nondeterministic select tie.

Once enumeration completes while the context remains active, Total is the enumerated candidate count.
For an incomplete scan after that point, Omitted is Total minus Inspected.
A complete scan has Omitted=0 and the existing complete reason.
An interrupted scan keeps Reason=timeout, including explicit context cancellation as the current API does.
GS06, GS07, and GS10 observe the exact result fields.

After the final inspection passes its completion checkpoint, Scan returns the complete result.
A later cancellation does not reopen completed work or change that result.
GS28 observes the final completion boundary.

Run each operation directly on the caller goroutine.
Scan returns only after its current operation has returned and performed cleanup.
Cancellation cannot interrupt os.ReadDir, os.Stat, os.ReadFile, or an injected function that ignores its context.
This limitation is explicit, not a new blocking behavior.
GS08 and GS11 prove the contract with controlled release barriers.

Preserve the existing non-cancellation enumeration-error result.
The current Scan ignores that error and completes with the returned candidate slice, normally empty.
This refactor does not redesign error reporting.
GS12 fixes that edge at the current Scan seam.

### Inventory and output preservation

Keep enumerateGuards, inspectGuard, HeaderFields, WiresScript, and action derivation behavior unchanged.
A missing or empty hooks directory still leaves the existing pre-push candidate.
An absent or incomplete static header creates the fallback row with an empty boundary and denial text `no manifest`.
A symlink to a regular hook can supply a header.
A FIFO is rejected before opening.
GS13 through GS16 preserve these facts.

Keep Command's existing root resolution and bounds.ContextGuardScanTimeout policy.
Keep the full TOON guards table and guard_scan metadata fields.
Keep brief rows, the guard_scan footer, and the full-manifests guidance line.
Retain exact complete and incomplete response fixtures.
GS17 through GS20 compare those outputs.

Keep action help order and deduplication for complete stale or unwired inventories.
Suppress action help for incomplete inventories.
Keep current exit 0 for an incomplete inventory, rather than introducing an error exit.
Keep existing invalid-argument and outside-repository command errors.
GS19 through GS21 observe those entry decisions.

### Single loop enforcement

Add checkGuardScanLoopOwner under the existing package-core-guard registered check.
Grade production Go files beneath internal/guards, including unused private function bodies.
Reject GoStmt, SelectStmt, channel creation, and send/receive pipeline constructs there.
The scan package has no other current production goroutine pipeline to preserve.
GS23 restores a renamed copy of the old Scan pipeline and requires the registered check to turn red.

Use AST nodes and real import bindings rather than text in comments or fixture strings.
The check does not ban goroutines in test code, including TestScanWaitsForCancelledWorkerCleanup's calling goroutine.
A kit source surface with a missing or unreadable scan owner produces a named diagnostic.
An absent kit source surface in a linked consumer produces no diagnostic.
GS24 and GS25 cover these dispositions.

Register the planned live-tree owner test in tier_live_tree_test.go.
Keep package-core-guard's existing name, tier, input source, and subject.
Conformance still runs inside the ordinary test phase.
GS26 prevents an unbound or hidden live-tree check.

### Conformance orchestration placement

Before extending the registered check, move only checkPackageCoreAndGuards from package_core_checks_test.go to checks_test.go.
Both existing ownership fences include these files.
Keep its name, signature, existing subcheck order, and callers unchanged.
Keep the other producers, npm formatter call, and their tests in package_core_checks_test.go.
Add the new ownership subcheck to the relocated function.

The accepted source has 470 lines in package_core_checks_test.go against the default 400-line limit.
It has no file grant, so adding an orchestration call there would fail the required growth lane.
checks_test.go has 641 lines against its existing 709-line grant, leaving 68 lines for the move and extensions.
The moved declaration spans 16 lines; preserve the destination cap and keep each new file within 400 lines.
Do not add a budget or acceptance grant.

Apply this placement only if the declaration remains in the original file when implementation begins.
If the other accepted guard outcome already moved it, extend that sole declaration in checks_test.go.
Retain both ownership subchecks if both outcomes have landed.
Either outcome can land first; neither requires the other's unlanded implementation.
[Placement evidence](../shared-guard-projection/placement-growth-repair.md) records the lane mechanism, source readers, and closure disposition.

## Implementation chunks

The preferred sequence follows the shared grammar checkpoint described in answer 3.
That sequence coordinates review and does not require an unlanded grammar implementation.
Its implementation does not depend on a new projection API.
It remains separately reviewable and can stop after its own green vertical slice.
The independently accepted GS1 outcome now has one complete vertical ticket and review checkpoint.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| GS1 / 1-simplify-guard-inventory-loop.md | synchronous scan preserves rows, cancellation, command output, and single ownership | GS01–GS28 | scan-loop, cleanup, command fixtures, and registered ownership mutation | yes, xhigh |

The chunk includes source claim verification and closure maintenance before final integration.
Overlapping conformance and command registry fences require serial implementation handoff, regardless of which outcome lands first.
The one ticket owns all 28 rows and closes their first-use contracts before its independent chunk review.

### Graph ordering and execution boundary

The stable old-to-new mapping is GS1 to GS1.
Each ticket retains its accepted outcome and independent Standards, Spec, and Coverage checkpoint.
Shared projection tickets run GP1, then GP2, then GP3; each successor waits for its predecessor's green commit and independent review.
GS1 is independently useful and consumes no unlanded projection interface.

Cross-spec conformance and command-registry writes require a serial integrated-source handoff.
The coordinator chooses either complete outcome first and refreshes source and headroom before overlapping successor charges.
The first ticket extending package-core-guard performs the accepted orchestration move; the other extends that sole declaration.
Neither plan replaces the other outcome's already landed subcheck or its independent npm producer assertions.

This graph's Writes union has 12 paths, equal to its unchanged path fence minus the review pickup.
Each introducing ticket owns its canonical fixture, binding, and live-tree closure and closes its own headroom obligations.
GP04 closes with both consumers in GP2, while GP1 still proves the first-use Git selection contract.
Every GP47–GP61 decoder case closes in GP2 and remains a mandatory GP3 actual-entry regression.

The version-1 completion plan is a supported planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator binds actual run/session identities and ticket verification ownership in version 2 with author-limit 1.
Each ticket receives a fresh author on gpt-5.6-sol/high, with GP3 and GS1 at the accepted xhigh effort.
Refresh current owner interfaces, source closure, and actual headroom before that charge.
The complete degraded repair and independently captured pre-migration executable precede GP1, not only final integration.
Staged status and planning validation grant no build authority.

### Completion plan

These commands and mutations are future obligations; none ran during authoring.
Before retaining any independent expectation, pin its mutation source and exact diagnostic and demonstrate a compiling behavioral red.
Require byte-identical restoration and the same focused green; invalid or restore-failed proof closes no obligation.
Final verification runs only after this plan's complete independently reviewed outcome.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "GS1",
      "tickets": [
        "1-simplify-guard-inventory-loop.md"
      ],
      "verification": [
        {
          "id": "scan-and-command",
          "command": "bench test --package ./internal/guards"
        },
        {
          "id": "registered-owner",
          "command": "bench test --check package-core-guard"
        },
        {
          "id": "conformance-closure",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "command-closure",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "post-inspection-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopPostInspectionCancellation$'",
          "probe": "Omit only the post-inspection context checkpoint; a returned interrupted candidate must publish rows and red the completed-prefix contract."
        },
        {
          "id": "copy-check-omission",
          "command": "bench test --package ./internal/conformance --run '^TestGuardScanLoopOwnerBitesOnPrivateCopy$'",
          "probe": "Omit only the registered checkGuardScanLoopOwner call while retaining its function and tests; an unused old worker pipeline must stop producing the expected registered diagnostic and red this witness."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check guard-scan-loop"
    },
    {
      "id": "scan-and-command",
      "command": "bench test --package ./internal/guards"
    },
    {
      "id": "command-closure",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "registered-owner",
      "command": "bench test --check package-core-guard"
    },
    {
      "id": "conformance-closure",
      "command": "bench test --package ./internal/conformance"
    }
  ]
}
```

## Testing decisions

Use the existing enumerateGuards and inspectGuard substitutions in the internal/guards test process.
Use channels only in test fixtures to control when an injected operation finishes.
Restore each package variable with test cleanup and do not run those substitutions in parallel.
The production loop has no channel or worker orchestration.

Keep TestScanWaitsForCancelledWorkerCleanup unchanged as the existing cleanup seam.
Keep the existing command stdout fixtures and their current exact assertions.
Use the ordinary registered ownership check for a surviving private loop.
No new subprocess harness, benchmark, or timing threshold is required.

### Seam diagram

```text
bench guards -> Command -> timeout context -> Scan
                                             |
                       enumerateGuards -> candidate order
                                             |
                       direct inspectGuard loop -> ScanResult
                                             |
                          existing formatter -> stdout and exit
```

Controlled test barriers attach to enumeration and inspection completion.
Command fixtures attach to the existing formatter output.
The registered conformance check attaches to the production scan ownership boundary.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GS01 | 1 | A complete canned scan retains candidate order and row order | planned TestScanLoopCompleteOrder in internal/guards/scan_loop_test.go | A reordered loop changes inventory meaning |
| GS02 | 1 | Enumeration runs once per Scan call | planned TestScanLoopEnumeratesOnce in internal/guards/scan_loop_test.go | A retry or per-candidate enumeration duplicates work |
| GS03 | 4 | A pre-cancelled Scan starts no operation | planned TestScanLoopPreCancelled in internal/guards/scan_loop_test.go | Cancellation is checked only after starting filesystem work |
| GS04 | 4 | Cancellation during enumeration retains unknown Total and Omitted fields | `internal/guards/guards_test.go` (`TestScanEnumerationTimeoutUsesUnknownCounts`) | Returned partial candidates fabricate known totals |
| GS05 | 2 | Cancellation before the next inspection retains the completed row prefix | planned TestScanLoopBetweenCandidates in internal/guards/scan_loop_test.go | The loop drops completed rows or starts the next candidate |
| GS06 | 3 | A cancelled inspection yields Omitted=Total-Inspected | `internal/guards/guards_test.go` (`TestScanTimeoutPreservesPartialRowsAndHonestCounts`) | The interrupted candidate is counted as inspected |
| GS07 | 3 | An active-context inspection with no rows still increments Inspected | planned TestScanLoopEmptyInspectionCounts in internal/guards/scan_loop_test.go | Counting rows instead of candidates reports false omissions |
| GS08 | 5 | Scan returns only after cancelled inspection cleanup completes | `internal/guards/guards_cleanup_test.go` (`TestScanWaitsForCancelledWorkerCleanup`) | Returning on cancellation races unfinished work |
| GS09 | 2 | Cancellation observed after inspection prevents that candidate's rows from publication | planned TestScanLoopPostInspectionCancellation in internal/guards/scan_loop_test.go | The loop publishes an interrupted candidate |
| GS10 | 3, 4 | Complete and cancelled scans retain all ScanResult metadata values | planned TestScanLoopResultMetadata in internal/guards/scan_loop_test.go | A renamed or defaulted field changes the result contract |
| GS11 | 6 | An operation that ignores cancellation keeps Scan pending until release | planned TestScanLoopBlockingOperationHasNoHardDeadline in internal/guards/scan_loop_test.go | A false deadline returns while filesystem work remains active |
| GS12 | 1 | A non-cancellation enumeration error retains the current complete result | planned TestScanLoopEnumerationError in internal/guards/scan_loop_test.go | A refactor silently introduces a new refusal state |
| GS13 | 7 | Absent and empty hook directories retain the pre-push candidate | planned TestScanLoopAbsentAndEmptyHooks in internal/guards/scan_loop_test.go | Optional hook surfaces change the inventory total |
| GS14 | 8 | Static header, symlink, missing, and incomplete cases retain their current rows | `internal/guards/guards_test.go` (`TestGuardRowReadsStaticHeader`) | The loop begins executing hooks or broadens header recognition |
| GS15 | 8 | FIFO header inspection returns without opening the FIFO | `internal/guards/guards_test.go` (`TestGuardRowRejectsFIFOWithoutOpening`) | A moved read blocks on a special file |
| GS16 | 9 | Claude and Codex wiring retain the same row facts | `internal/guards/guards_test.go` (`TestRowsReportFollowOnManifestAndBothHarnessWires`) | Enumeration migration loses a harness configuration |
| GS17 | 10 | Complete clean-primary stdout matches its existing fixture | `internal/guards/guards_test.go` (`TestCommandPreservesCheckedInCleanPrimaryResponse`) | Field presence alone permits changed output |
| GS18 | 11 | Enumeration timeout stdout matches its existing unknown-count fixture | `internal/guards/guards_test.go` (`TestCommandPreservesCheckedInEnumerationTimeoutPrimaryResponse`) | The command renders fabricated counts |
| GS19 | 12 | Complete stale and unwired stdout retains the existing action fixture | `internal/guards/guards_test.go` (`TestCommandPreservesCheckedInStaleAndUnwiredPrimaryResponse`) | The new loop changes help order or duplicates actions |
| GS20 | 11, 13 | Inspection timeout stdout retains its existing fixture without action help | `internal/guards/guards_test.go` (`TestCommandPreservesCheckedInIncompleteTimeoutPrimaryResponse`) | Partial data emits misleading complete-scan guidance |
| GS21 | 14 | Invalid arguments and a non-repository retain their command error tuples | planned TestScanLoopCommandErrors in internal/guards/scan_loop_test.go | Simplification changes Command's operand or root failure policy |
| GS22 | 1 | Rows retains Scan(background).Rows behavior | planned TestScanLoopRowsAdapter in internal/guards/scan_loop_test.go | A second public path keeps the old scanner |
| GS23 | 15 | A renamed unused worker-pipeline copy makes package-core-guard red | planned TestGuardScanLoopOwnerBitesOnPrivateCopy in internal/conformance/guard_scan_owner_test.go | New loop tests pass while the old implementation survives |
| GS24 | 15 | A linked consumer without kit scan source passes the ownership subcheck | planned TestGuardScanLoopOwnerOptionalSurface in internal/conformance/guard_scan_owner_test.go | The new invariant rejects ordinary consumer repositories |
| GS25 | 15 | A kit source tree with a missing or unreadable scan owner receives a named diagnostic | planned TestGuardScanLoopOwnerMissingOwner in internal/conformance/guard_scan_owner_test.go | A missing owner makes the check vacuously green |
| GS26 | 15 | The live-tree check runs through its registered owner with explicit classification | planned TestGuardScanLoopOwnerHoldsOnTheLiveTree in internal/conformance/guard_scan_owner_test.go | A private helper or hidden test never grades the ordinary root |
| GS27 | 16, 17 | Deterministic scan family S produces identical command tuples before and after migration | planned TestScanLoopDifferentialCommandFamilies in internal/guards/scan_loop_test.go | Existing partial fixtures miss an unasserted command change |
| GS28 | 1, 6 | Cancellation after the final completion checkpoint leaves the complete result unchanged | planned TestScanLoopAfterFullCompletion in internal/guards/scan_loop_test.go | A late context read discards fully completed work |

GS27 runs preserved pre-change and candidate owner-package test executables over the same controlled enumeration and inspection events.
Prepare the contract driver before replacing the loop, then retain its pre-change executable as the reference.
The reference capture is independently reviewed and remains separate from candidate code.
Both executions feed results through the real Command formatter in their own test process.
Compare full and brief stdout, stderr, exit, and ScanResult values.
No elapsed-time equality is required, and no production filesystem race defines the expected result.

### Edge inventory

| Family S member | Controlled event or fixture | Result boundary |
|---|---|---|
| complete | ordered candidates with zero, one, and multiple rows | current complete metadata and row order |
| pre-cancelled | cancelled context before Scan | no enumeration, unknown counts |
| enumeration cancellation | barrier releases only after cancellation | unknown total and omitted counts |
| between candidates | cancel after publishing a completed candidate | completed prefix and known omitted arithmetic |
| inspection cancellation | candidate returns after cancellation | no interrupted candidate publication |
| no rows | inspection completes with an empty row slice | candidate still counted |
| enumeration error | injected error while context is active | existing complete result |
| blocking operation | injected operation ignores context until explicit release | no hard return deadline |
| optional surfaces | absent and empty hooks directory | pre-push candidate retained |
| static headers | regular file, symlink, FIFO, missing file, incomplete header | current HeaderFields contract |
| wiring | absent, empty, invalid JSON, valid Claude and Codex configuration | current WiresScript facts |
| command | clean, stale, unwired, enumeration timeout, inspection timeout, full and brief | current stdout and help contracts |
| command errors | invalid flag and outside-repository cwd | current command errors |

Use bounded fixture waits only to detect an unexpected early return or fixture deadlock.
Do not treat a timing measurement as the product's deadline contract.
Existing file-system failure behavior stays with HeaderFields and enumeration.
No new cleanup-error state is introduced by the direct loop.

**Won't handle** an identical winner for simultaneous cancellation and completion — Scan uses explicit context checkpoints rather than Go select arbitration.
**Won't handle** a hard filesystem deadline — Scan waits for its current operation to return.
**Won't handle** parallel guard inspection — Scan retains ordered candidate inspection.
**Won't handle** new enumeration-error diagnostics — Scan retains its existing active-context result.
**Won't handle** new hook execution — HeaderFields remains a static reader.
**Won't handle** guard grammar or envelope policy — the separate shared-guard-projection outcome owns those migrations.

**Won't handle** degraded hook repair — the accepted degraded-guard-refusal outcome remains its own checkpoint.

## Ownership fences

These 13 exact path entries preserve the independently accepted implementation fence.
Current author writes planning artifacts only; the accepted graph grants no implementation authority.

- `internal/guards/guards.go`
- `internal/guards/scan_loop_test.go` (new)
- `internal/conformance/guard_scan_owner_test.go` (new)
- `internal/conformance/package_core_checks_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/tier_live_tree_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/guard-scan-loop.md` (new)

BoundFiles in internal/tickets/registry_data.go derives command-registry co-names from the internal/guards package prefix.
Keep those registries and their assertions intact.
Their presence in the fence permits necessary binding consistency repair, not new commands or weakened assertions.
The existing guards_test.go, guards_cleanup_test.go, and testdata stdout fixtures remain unchanged witnesses.
The source BASE and MUTATE path sweep found no fixture pin for internal/guards/guards.go.
No anchored guidance file moves.

## Out of scope

Shared grammar migration remains its own spec and checkpoint.
Its price follows specs/shared-guard-projection/spec.md, not a copied estimate.
Hard-deadline filesystem isolation would require a different capability and failure contract, estimated at 3 production edits and 2 gate runs.
Parallel inventory is a separate performance capability, estimated at 2 production edits and 2 gate runs after its own design.
Neither capability is authorized here.

## Further notes

### Source-clause accounting

| Source clause in answer 3 | Disposition |
|---|---|
| “Use three separately reviewable outcomes under FT366.” | separate staged specs and future checkpoints |
| “First, repair malformed degraded input and its diagnostic.” | accepted predecessor spec |
| “Second, migrate common grammar and envelope facts.” | shared-guard-projection spec |
| “Third, simplify the guard inventory scan if its behavior stays fixed.” | GS01–GS28 and GS27 differential exit row |
| “The scan waits for its goroutines after cancellation” | current Scan and unchanged cleanup witness, GS08 |
| “A context-aware loop removes that indirection.” | direct loop and GS23 private-copy refusal |
| “It does not establish a hard deadline for a blocking filesystem operation.” | GS11 and explicit Won't handle |
| “retain tests through actual hook and guard entries.” | grammar spec actual-entry rows, scan Command fixtures GS17–GS21 |
| “Cover attached separators, unclosed quotes, wrapper depth, Git options, and worktree exec children.” | grammar spec GP02–GP10, GP34–GP35, GP46 |
| “Cover missing library, missing core, core error, and valid recovery commands.” | grammar spec GP36–GP41 |
| “Prove that each guard retains its own uncertainty and diagnostic policy.” | grammar spec caller preservation rows, scan metadata and output rows |

Each source clause occurs once in the ready answer.
The shared map stays in place for the other outcomes and its explicit rg residual.
Flagged additions: none in product behavior.
The private-loop ownership check enforces the approved removal of worker indirection.

### Pre-review proof checklist

- Cited symbols: Scan, Rows, Command, enumerateGuards, inspectGuard, HeaderFields, WiresScript, and each cited test were read in the subject tree.
- Import edges: none introduced. Existing internal/guards imports remain within their current dependency direction.
- Source-row clauses and occurrences: the table covers answer 3. Other outcomes retain their named source clauses in their own specs.
- Promised field labels: Rows, Status, Inspected, Total, Omitted, and Reason remain unchanged. The emitted guard_scan and guards table shapes remain unchanged.
- Changed-function callers: bench consumers enumerates Scan callers in Rows, Command, guards_test.go, and guards_cleanup_test.go.
- Copy survival: GS23 plants a renamed unused worker pipeline and observes the registered package-core-guard diagnostic.
- Rendered-shape readers: Command, checked-in stdout fixtures, SessionStart's brief query, AXI query registry, and command help retain the current shapes.
- Pin operators: BoundFiles selects internal/guards by tickets.Covers prefix containment. Preflight requires each bound path to be co-named.
- Entry reads: Command reads cwd and Git root before Scan. Enumeration reads hook and harness configuration paths, and inspection reads static header files.
- Derived expectations: existing fixture files and independently reviewed canned reference results own output expectations. Candidate scan output is not the oracle.
- Consolidated rules: only Scan's enumeration and per-candidate worker orchestration become direct calls. Enumeration and inspection remain their current single owners.
- Quantified obligations: family S enumerates complete, cancelled, empty, absent, unreadable, and blocking operation states.
- Workflow-step writes: none. The command remains a read-only query with no new durable artifact.

The caller sweep is tied to subject a9c395e77fec36d1f60f6d057c654aecf5720e24 and citation c1ce54b3baafa91dfcefa87a26b1c453c53e17d67cdc95d394c0b3d7e165803b.
Scan callers are Rows, Command, TestScanWaitsForCancelledWorkerCleanup, TestScanTimeoutPreservesPartialRowsAndHonestCounts, and TestScanEnumerationTimeoutUsesUnknownCounts.
Their functions retain signatures and behavior, so existing caller files require no source edit.
The new test file avoids crowding existing guarded test files.

The standing check grades optional kit source only.
It tolerates the separate staged specs without implementation and enables its invariant when GS1 completes.
Its own omission tests distinguish absent consumer source from amputated kit source.
No new conformance registry count, profile row, or injected-port row is needed.

### Review repair record

[GS-S1 repair evidence](review-r1-repairs.md) records the two corrected source facts and their existing witnesses.
Independent review accepted the repaired spec at 02c799aa8e7ba9f1dfcf38a14906f51ab4a6d8db with no material findings.
The synchronous boundaries and coverage map remain unchanged.
The xhigh placement confirmer accepted the bounded amendment at 6c13b4b before ticket slicing.
Independent graph review accepted aa86c64d; no runtime preservation proof is claimed.

### Review disposition

| Item | Proposed disposition |
|---|---|
| Implementation line | gpt-5.6-sol/high, GS1 xhigh, accepted at 02c799aa |
| Seam and coverage | synchronous scan, controlled cancellation, exact Command output, registered ownership, accepted at 02c799aa |
| Ownership fences | exact path fence above, placement accepted at 6c13b4b; graph accepted at aa86c64d |
| Exclusions | no hard deadline, no parallel inspection, no new diagnostics, accepted at 02c799aa |
| Ticket graph | complete vertical graph authored after independent SPEC acceptance; graph accepted at aa86c64d |

### Numbered breakdown

| ticket | title | Blocked by | delivered outcome |
|---|---|---|---|
| 1-simplify-guard-inventory-loop.md | Simplify guard inventory scanning | none | Synchronous scan, exact command output, cancellation, and registered loop ownership |
