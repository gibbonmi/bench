# Simplify guard inventory scanning with a synchronous loop

Status: staged
Decision source: `decisions/architecture-guards.md`, answer 3, reviewed 2026-10-06
Subject: `f07508b2add49469d686614c70f575531e8ea445` (refresh of 2026-10-20)
Audience: every repository that links the kit
Related roadmap: FT366
Verification log: pending — refresh review at f07508b2a

## Problem

`internal/guards/guards.go` (`Scan`) launches enumeration and inspection workers (`guards.go:172-200`).
Cancellation then waits for those workers to finish before Scan returns.
The channels and selects therefore add indirection without a hard deadline for blocking filesystem work.
`TestScanWaitsForCancelledWorkerCleanup` already requires cleanup before return.

## Solution

Run the same enumeration and inspection operations in a synchronous, context-aware loop.
Preserve completed rows, honest counts, static manifest reads, command output, and return codes.
State cancellation as a checkpoint between operations.
Do not promise to interrupt an in-flight filesystem call.

## User stories

Line: opus / high
Implementation-line reason: GS1 is the only chunk and the harder one. The loop shape is known and every boundary has an exact predicate, but the cancellation checkpoints and the byte-exact output preservation need careful work at the existing package-variable seams. Test strength is high, because checked-in stdout fixtures and controlled barriers already exist.
Harder chunks: GS1, at xhigh effort on the same model.

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
18. As a maintainer, I want the scan timeout to stay with its bounds-policy owner, so that the bounds registry keeps grading the scan deadline.

## Implementation decisions

### Loop and cancellation contract

Keep Scan(ctx, root), Rows(root), Command, ScanResult, and their exported signatures unchanged.
Rows returns `[][]string`; the package defines no row type, and the refactor adds none.
Keep enumerateGuards and inspectGuard as the existing package-variable test seams (`guards.go:61`, `:81`).
Do not add an injected parameter or a new interface.
The concrete loop performs enumeration once, then inspects candidates in their current order.
GS01, GS02, and GS22 observe these contracts.

Check ctx.Err before enumeration and after it returns.
A pre-cancelled context starts no operation.
Today enumeration starts before any context check (`guards.go:172-176`), so GS03 is red on the current tree.
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
This tightens that race, because it withholds a candidate when the checkpoint observes cancellation.
Preservation evidence fixes event order and does not claim an identical winner for a nondeterministic select tie.

Once enumeration completes while the context remains active, Total is the enumerated candidate count.
For an incomplete scan after that point, Omitted is Total minus Inspected.
A complete scan has Omitted=0 and the existing complete reason.
An interrupted scan keeps Reason=timeout, including explicit context cancellation, as the current API does.
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
The current Scan ignores that error and completes with the returned candidate slice, normally empty (`guards.go:180-183`).
This refactor does not redesign error reporting.
GS12 fixes that edge at the current Scan seam.

### Inventory and output preservation

Keep enumerateGuards, inspectGuard, HeaderFields, WiresScript, and action derivation behavior unchanged.
A missing or empty hooks directory still leaves the existing pre-push candidate.
An absent or incomplete static header creates the fallback row with an empty boundary and denial text `no manifest`.
A symlink to a regular hook can supply a header.
A FIFO is rejected before it is opened.
GS13 through GS16 preserve these facts.

Keep Command's existing root resolution and its timeout policy.
The policy line is `bounds.VerdictWindow(bounds.GuardScanTimeout)` (`guards.go:59`), and Command applies it through `bounds.Context` (`guards.go:280`).
The `bounds-policy` check pins that exact text in guards.go (`internal/conformance/bounds_policy_test.go:38`).
GS29 keeps the pinned line.

Keep the full TOON guards table and the guard_scan metadata fields.
Keep the brief rows, the guard_scan footer, and the full-manifests guidance line.
Retain the exact complete and incomplete response fixtures.
GS17 through GS20 compare those outputs.

Keep action help order and deduplication for complete stale or unwired inventories.
The deduplication lives in `internal/axi/action.go`, which this refactor does not touch.
Suppress action help for incomplete inventories.
Keep the current exit 0 for an incomplete inventory, rather than a new error exit.
Keep the existing invalid-argument and outside-repository command errors.
GS19 through GS21 observe those entry decisions.

### Single loop enforcement

Add checkGuardScanLoopOwner under the existing registered package-core-guard check.
It grades production Go files beneath internal/guards, including unused private function bodies.
It rejects GoStmt, SelectStmt, channel creation, and send or receive pipeline constructs there.
guards.go is the only production file in internal/guards, and Scan holds every such construct today.
GS23 restores a renamed copy of the old Scan pipeline and requires the registered check to turn red.

Use AST nodes and real import bindings, not text in comments or fixture strings.
The check does not ban goroutines in test code, including the calling goroutine of TestScanWaitsForCancelledWorkerCleanup.
A kit source surface with a missing or unreadable scan owner produces a named diagnostic.
An absent kit source surface in a linked consumer produces no diagnostic.
GS24 and GS25 cover these dispositions.

Register the planned live-tree owner test in tier_live_tree_test.go.
Keep package-core-guard's existing name, tier, input source, and subject (`internal/conformance/registry/checks.go:19`).
Conformance still runs inside the ordinary test phase.
GS26 prevents an unbound or hidden live-tree check.

The `stdlib-reimplementation` check (`internal/conformance/stdlib_reimplementation_test.go:39`) also grades guards.go.
The new loop must not hand-write a loop that a standard-library function replaces.

### Conformance orchestration placement

Before you extend the registered check, move only checkPackageCoreAndGuards from package_core_checks_test.go to checks_test.go.
Keep its name, signature, existing subcheck order, and callers unchanged.
Keep the other producers, the npm formatter call, and their tests in package_core_checks_test.go.
Add the new ownership subcheck to the relocated function.

package_core_checks_test.go has 470 lines against the default 400-line limit, and no grant covers it.
The required commit lane runs `structure --growth` (`internal/gate/lane.go:83`).
That ratchet reports FILE GREW for a source file over its limit that also grew (`internal/structure/structure.go:309-331`).
An added orchestration call there would therefore turn the lane red.
checks_test.go has 642 lines against its 739-line grant (`.bench/structure.budgets:51`), which leaves 97 lines.
The moved declaration spans 16 lines (`package_core_checks_test.go:22-37`).
Keep checks_test.go within 739 lines after both guard outcomes, and keep each new file within 400 lines.
Do not add a budget or acceptance grant.

The move changes no reader contract.
`checks_test.go:49` binds package-core-guard by function symbol, so the declaration file does not matter.
`node_runtime_policy_test.go:143` and `:147` call the same package symbol and need no edit.
`package_core_diagnostics_test.go:202-228` counts formatter calls in package_core_checks_test.go; the moved function holds none.
`TestSkillsIndexConformanceCarriesNoSecondReader` reads checks_test.go for copied skills-index literals; the moved function holds none.

Apply this placement only if the declaration remains in the original file when implementation begins.
If the shared-guard-projection outcome already moved it, extend that sole declaration in checks_test.go.
Retain both ownership subchecks if both outcomes have landed.
Either outcome can land first; neither requires the other's unlanded implementation.

## Implementation chunks

The preferred sequence follows the shared grammar checkpoint described in answer 3.
That sequence coordinates review and does not require an unlanded grammar implementation.
GS1 does not depend on a new projection API.
It remains separately reviewable and can stop after its own green vertical slice.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| GS1 / 1-simplify-guard-inventory-loop.md | synchronous scan preserves rows, cancellation, command output, and single ownership | GS01–GS29 | scan-loop, cleanup, command fixtures, bounds policy, and registered ownership mutation | yes, xhigh |

The chunk includes source claim verification and closure maintenance before final integration.
The one ticket owns all 29 rows and closes their first-use contracts before its independent chunk review.

### Graph ordering and execution boundary

The stable old-to-new mapping is GS1 to GS1.
GS1 keeps its independent Standards, Spec, and Coverage checkpoint.
GS1 is independently useful and consumes no unlanded projection interface.

Both guard outcomes write checks_test.go, package_core_checks_test.go, and tier_live_tree_test.go.
These shared writes require a serial integrated-source handoff between the two outcomes.
The coordinator chooses either complete outcome first and refreshes source and headroom before an overlapping successor charge.
The first ticket that extends package-core-guard performs the orchestration move; the other extends that sole declaration.

Derived Writes scope for the slicer, at f07508b2a: `internal/guards/guards.go`, `internal/guards/scan_loop_test.go` (new), `internal/conformance/guard_scan_owner_test.go` (new), `internal/conformance/package_core_checks_test.go`, `internal/conformance/checks_test.go`, `internal/conformance/tier_live_tree_test.go`.

The refresh drops six paths from the 2026-10-06 fence.
Commit 1e8846bb0 deleted the ownership closure that required the five command-registry co-names: `cmd/bench/command_registry.go`, `cmd/bench/command_registry_test.go`, `cmd/bench/help_inventory_test.go`, `internal/conformance/axi_query_registry_test.go`, and `internal/conformance/subcommand_routing_table_test.go`.
No row changes `internal/conformance/registry/checks.go`, because the registry row keeps its metadata.
The check floor still requires the `./cmd/bench` run, because internal/guards is a command package (`internal/tickets/registry_data.go:16`).
The slicer's research confirms or restores each dropped path.

The existing guards_test.go, guards_cleanup_test.go, and testdata stdout fixtures stay unchanged witnesses at f07508b2a.
The staged, unbound `shell-argument-ownership` spec also writes `internal/guards/guards_test.go`.
If that spec lands first, build entry refreshes these witnesses and their fixtures against the integrated source.

The version-1 completion plan is a planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator binds actual run and session identities and ticket verification ownership in version 2 with author-limit 1.
The ticket receives a fresh author on opus at xhigh effort.
Refresh current owner interfaces, source closure, and actual headroom before that charge.
Staged status and planning validation grant no build authority.

### Completion plan

These commands and mutations are future obligations; none ran during authoring.
Before you retain an independent expectation, pin its mutation source and exact diagnostic and demonstrate a compiling behavioral red.
Require byte-identical restoration and the same focused green; an invalid or restore-failed proof closes no obligation.
Final verification runs only after the complete, independently reviewed outcome.

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
          "id": "bounds-policy",
          "command": "bench test --check bounds-policy"
        },
        {
          "id": "conformance-closure",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "registry-closure",
          "command": "bench test --package ./internal/conformance/registry"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check root-conformance"
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
    },
    {
      "id": "root-conformance",
      "command": "bench test --check root-conformance"
    }
  ]
}
```

## Testing decisions

Use the existing enumerateGuards and inspectGuard substitutions in the internal/guards test process.
Use channels only in test fixtures to control when an injected operation finishes.
Restore each package variable with test cleanup, and do not run those substitutions in parallel.
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
The registered conformance checks attach to the production scan ownership boundary and the bounds-policy pin.

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
| GS14 | 8 | Static header, symlink, and incomplete-header cases retain their current rows | `internal/guards/guards_test.go` (`TestGuardRowReadsStaticHeader`) | The loop begins executing hooks or broadens header recognition |
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
| GS26 | 15, 16 | The live-tree check runs through its registered owner with explicit classification | planned TestGuardScanLoopOwnerHoldsOnTheLiveTree in internal/conformance/guard_scan_owner_test.go | A private helper or hidden test never grades the ordinary root |
| GS27 | 16, 17 | Deterministic scan family S produces identical command tuples before and after migration | planned TestScanLoopDifferentialCommandFamilies in internal/guards/scan_loop_test.go | Existing partial fixtures miss an unasserted command change |
| GS28 | 1, 6 | Cancellation after the final completion checkpoint leaves the complete result unchanged | planned TestScanLoopAfterFullCompletion in internal/guards/scan_loop_test.go | A late context read discards fully completed work |
| GS29 | 18 | guards.go keeps the `bounds.VerdictWindow(bounds.GuardScanTimeout)` text that bounds-policy pins | `internal/conformance/bounds_policy_test.go` (`checkBoundsPolicy`, registered as bounds-policy) | A rewrite that drops or respells the policy line reds bounds-policy with "does not consume" |

GS27 runs preserved pre-change and candidate owner-package test executables over the same controlled enumeration and inspection events.
Prepare the contract driver before you replace the loop, then retain its pre-change executable as the reference.
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

No existing test drives a nonexistent header path; GS27's family S owns the missing-file member.
Use bounded fixture waits only to detect an unexpected early return or a fixture deadlock.
Do not treat a timing measurement as the product's deadline contract.
Existing filesystem failure behavior stays with HeaderFields and enumeration.
The direct loop introduces no new cleanup-error state.

**Won't handle** an identical winner for simultaneous cancellation and completion — Scan uses explicit context checkpoints rather than Go select arbitration.
**Won't handle** a hard filesystem deadline — Scan waits for its current operation to return.
**Won't handle** parallel guard inspection — Scan retains ordered candidate inspection.
**Won't handle** new enumeration-error diagnostics — Scan retains its existing active-context result.
**Won't handle** new hook execution — HeaderFields remains a static reader.
**Won't handle** guard grammar or envelope policy — the separate shared-guard-projection outcome owns those migrations.
**Won't handle** historical line pins to Scan in `ASSESSMENT.md:42` and `docs/research/bench-efficiency-assessment-2026-10-06.md:122` — they describe a dated tree and are not contracts.
**Won't handle** the stale column list in the package comment (`guards.go:1-2`) — it predates this spec, and the refactor does not change the output.

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
| “Use three separately reviewable outcomes under FT366.” | separate specs and checkpoints; the degraded outcome landed |
| “First, repair malformed degraded input and its diagnostic.” | landed predecessor outcome |
| “Second, migrate common grammar and envelope facts.” | shared-guard-projection spec |
| “Third, simplify the guard inventory scan if its behavior stays fixed.” | GS01–GS29 and the GS27 differential row |
| “The scan waits for its goroutines after cancellation” | current Scan and the unchanged cleanup witness, GS08 |
| “A context-aware loop removes that indirection.” | direct loop and GS23 private-copy refusal |
| “It does not establish a hard deadline for a blocking filesystem operation.” | GS11 and the explicit Won't handle |
| “retain tests through actual hook and guard entries.” | grammar spec actual-entry rows, scan Command fixtures GS17–GS21 |
| “Cover attached separators, unclosed quotes, wrapper depth, Git options, and worktree exec children.” | grammar spec GP02–GP10, GP34–GP35, GP46 |
| “Cover missing library, missing core, core error, and valid recovery commands.” | grammar spec GP36–GP41, GP71 |
| “Prove that each guard retains its own uncertainty and diagnostic policy.” | grammar spec caller preservation rows, scan metadata and output rows |

Each source clause occurs once in the ready answer.
The shared map stays in place for the projection outcome and its explicit rg residual.

Flagged additions for reviewer veto:

- GS29 and story 18 are new. They name an existing pin that the 2026-10-06 spec missed; they add no product behavior.
- The refresh drops six paths from the 2026-10-06 fence, as Graph ordering states.
- The spec keeps `Related roadmap: FT366` instead of a `Roadmap:` line. FT366 closes only when both bound specs land, so a retiring spec must not remove the shared row.

The private-loop ownership check enforces the approved removal of worker indirection.

### Seam research citations

- `checkGuardScanLoopOwner` extends the package-core-guard orchestration at `internal/conformance/package_core_checks_test.go:22-37`, registered at `internal/conformance/registry/checks.go:19`.
- Its subject is the worker pipeline at `internal/guards/guards.go:172-200`, which answer 3 cites at `guards.go:185` and `:198`.

### Build-process constraints

These rows of `docs/known-issues.md` bind the build of this spec:

- Row 10 (`docs/known-issues.md:10`): rebuild the worktree binary after the Go ticket commit, before preflight.
- Row 13 (`docs/known-issues.md:13`): after each assignment record, run the plan amendment before the chunk record.
- Row 14 (`docs/known-issues.md:14`): a test-only repair needs a confirming round on all three axes.
- Row 15 (`docs/known-issues.md:15`): run the evidence currency check through `bench worktree exec`.
- Row 16 (`docs/known-issues.md:16`): the seam research citations above satisfy it.

Structure obligations: package_core_checks_test.go must not grow, and it shrinks by the moved declaration.
guards.go has 327 lines and must stay within 400.
checks_test.go must stay within its 739-line grant.
Each new file must stay within 400 lines.
No file in the derived scope needs a new grant.

### Pre-review proof checklist

- Cited symbols: the refresh research read Scan, Rows, Command, enumerateGuards, inspectGuard, HeaderFields, WiresScript, and each cited test at f07508b2a.
- Import edges: none introduced. Existing internal/guards imports keep their current dependency direction.
- Source-row clauses and occurrences: the table covers answer 3. The projection outcome keeps its own source clauses.
- Promised field labels: Rows, Status, Inspected, Total, Omitted, and Reason remain unchanged. The emitted guard_scan and guards table shapes remain unchanged.
- Changed-function callers: the sweep below lists the Scan and Rows callers.
- Copy survival: GS23 plants a renamed unused worker pipeline and observes the registered package-core-guard diagnostic.
- Rendered-shape readers: Command, the checked-in stdout fixtures, the SessionStart brief query (`internal/sessioninspect/sessioninspect.go:152`, `internal/systemtest/session_start_test.go:198`), the isolated command fixtures (`cmd/bench/isolated_command_fixtures_test.go:58-70`), and command help keep the current shapes.
- Pin operators: `bounds-policy` pins the timeout text in guards.go (GS29). The check floor adds `./cmd/bench` because internal/guards is a command package.
- Entry reads: Command reads cwd and the Git root before Scan. Enumeration reads hook and harness configuration paths, and inspection reads static header files.
- Derived expectations: the existing fixture files and the independently reviewed canned reference results own output expectations. Candidate scan output is not the oracle.
- Consolidated rules: only Scan's enumeration and per-candidate worker orchestration become direct calls. Enumeration and inspection remain their current single owners.
- Quantified obligations: family S lists the complete, cancelled, empty, absent, unreadable, and blocking operation states.
- Workflow-step writes: none. The command remains a read-only query with no new durable artifact.

The caller sweep is tied to f07508b2add49469d686614c70f575531e8ea445.
Scan callers are Rows (`guards.go:164`), Command (`guards.go:282`), TestScanWaitsForCancelledWorkerCleanup, TestScanTimeoutPreservesPartialRowsAndHonestCounts, and TestScanEnumerationTimeoutUsesUnknownCounts.
Command is also called by the SessionStart brief query (`internal/sessioninspect/sessioninspect.go:152`), which passes the brief output through unparsed.
Their functions keep signatures and behavior, so the existing caller files need no source edit.
The new test file avoids crowding the existing guarded test files.

The standing check grades optional kit source only.
It enables its invariant when GS1 completes.
Its own omission tests distinguish absent consumer source from amputated kit source.
No new conformance registry count, profile row, or injected-port row is needed.

### Review disposition

| Item | Proposed disposition |
|---|---|
| Implementation line | opus/high, GS1 at xhigh |
| Seam and coverage | synchronous scan, controlled cancellation, exact Command output, registered ownership, bounds-policy pin |
| Derived Writes scope | six paths, six paths dropped from the prior fence |
| Exclusions | no hard deadline, no parallel inspection, no new diagnostics, historical pins |
| Ticket graph | the slicer refreshes the one ticket against this spec |

### Numbered breakdown

| ticket | title | Blocked by | delivered outcome |
|---|---|---|---|
| 1-simplify-guard-inventory-loop.md | Simplify guard inventory scanning | none | Synchronous scan, exact command output, cancellation, bounds-policy pin, and registered loop ownership |
