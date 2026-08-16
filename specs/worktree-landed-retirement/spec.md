# worktree landed-assignment visibility and bulk retirement

Status: staged

Decision source: `specs/worktree-landed-retirement/decisions/worktree-landed-retirement.md` (ready compiled map, FT210, all fifteen tickets reviewer-closed 2026-08-16; #9 re-closed the same day during spec authoring)

Verification log: spec 2 + tickets 2 iteration(s) to accept — reviewer-capped at two rounds each with `gpt-5.6-sol`/high via `codex exec` (round 2 = fix verification, not a fresh hunt). Spec: round 1 BLOCK (10 blocking, 1 advisory; one finding reopened map #9 with the reviewer), round 2 BLOCK on 5 partials, all folded after the cap (per-row tuple check on apply, `--discard-branch` a no-op under `--landed`, `ForbidInSection`, socket/device rows, unreadable-lease exclusion). Tickets: round 1 BLOCK (11 blocking, 1 advisory → re-sliced 7 to 17), round 2 BLOCK on 2 partials, both folded after the cap (LR11/LR12 split, apply pre-mutation refusal criterion). No terminal advisory ACCEPT verdict; reviewer sign-off is the hard stop.

## Problem

An assignment worktree whose branch has landed and whose writer is gone still reads
`state=active`, so the session-start line `retained active=N` looks healthy when it
means nobody released N trees. Nothing makes the release happen: a multi-ticket build
releases its integration source through `bench worktree land` and leaves every
per-ticket worktree behind, and clearing them means running `bench worktree clean
<path>` once per tree from a hand-written loop.

## Solution

Classify a **landed assignment** — ledger state `active`, branch landed on the default
branch, lease not live — as a derived reason, never a new ledger state. `bench worktree
list` and the resume summary project it (`retained landed=N` beside `active=N`) and
advertise one remedy, `bench worktree clean --landed`: a selector on the existing
plan/apply verb that prints every landed assignment in the repository under one set
fingerprint and, on `--apply`, removes the clean ones, settles their ledger records,
and leaves dirty or uncertain rows listed with their per-path remedy. Workflow prose
makes release mechanical: the coordinator releases each accepted independent slice,
and `/bench-final-check` runs the sweep as a named step before its landing report.

## User stories

**Seeing landed assignments** — Line: `opus` / medium. Lifecycle classification at a
known seam, mid because it decides what a destructive verb later selects.

1. As a session starting up, I want the resume summary to count landed assignments separately (`retained landed=N`), so that `active=N` no longer reads as healthy when nobody released N trees.
2. As an agent reading `bench worktree list`, I want landed rows identifiable from the existing `state`, `landed`, and `lease` columns without a new column, so that the schema stays stable.
3. As an agent reading `bench worktree list`, I want the bare `bench worktree clean --landed` invocation advertised as an action when any row is landed, so that I know the remedy without consulting docs.
4. As a session starting up, I want the summary to advertise the same bare invocation when landed > 0 and never a discard flag, so that the remedy line cannot manufacture the next residue.
5. As an operator, I want a landed tree older than the orphan window listed as landed rather than as an orphan candidate, so that one tree is not reported twice with two remedies.
6. As an operator, I want a landed tree with a live lease counted as `live-lease`, not landed, so that a running writer is never targeted.
7. As an operator, I want a landed tree with a dead or absent lease counted landed, so that a crashed or departed writer does not immortalize the tree.
8. As an operator, I want a landed tree with an unparseable lease counted landed but retained `uncertain` by the sweep, so that visibility and destruction are decided by different bars.
9. As an operator, I want dirty tracked state and ignored residue to leave the classification unchanged, so that residue cannot hide a landed tree from the count.
10. As an operator, I want unknown landedness (no default branch, errored proof) left as plain `active`, so that a guess never feeds a destructive verb.
11. As an operator, I want `cleanup-pending`, `recovered`, and `complete` records never classified landed, so that recovery state stays a per-row decision.
12. As a maintainer, I want the classification derived, not a new ledger state, so that no migration or state-machine change is needed.
13. As a maintainer, I want one classifier function shared by list, summary, and the selector, so that the three cannot disagree.

**Retiring them in bulk** — Line: `opus` / high. Composes the cleanup transaction over a
set; the drift and interruption edges are where a cheaper line ships a partial apply.

14. As an operator, I want `bench worktree clean --landed` to print one plan over every landed assignment repository-wide, so that a build's trees are not retired by a hand-written loop.
15. As the closing phase, I want the selection independent of which request or session created the tree, so that no build-scoped bookkeeping is required.
16. As an operator, I want every plan row to carry one set fingerprint, so that I apply the whole set with one value.
17. As an operator, I want `--apply <fp>` refused entirely if any row drifted (new member, dirty tree, live lease, moved HEAD), so that a partial apply never happens silently.
18. As an operator, I want a later row that drifts after set validation refused before its transaction, so that a plan the set never bound — especially a preserving one — never executes.
19. As an operator, I want removed rows' ledger records settled in the same apply, so that `list` and the summary are correct immediately.
20. As an operator, I want dirty or unknown-lease rows retained and listed with `bench worktree clean <path>`, so that bulk never authors recovery refs.
21. As an operator, I want an empty landed set to exit 0 with an empty plan and no fingerprint, so that nothing-to-do is not an error.
22. As an operator, I want `--apply` on an empty set to be an invocation error, so that a fingerprint that was never issued cannot be applied.
23. As an operator, I want `--landed` with a path operand, or a malformed `--apply` value, refused as an invocation error, so that selector and operand stay exclusive and the fingerprint grammar holds.
24. As an operator, I want `--discard-ignored`, `--full`, and `--discard-branch` applied to every row, so that modifiers behave as on per-path clean.
25. As an operator, I want proven-landed branches deleted with their trees exactly as per-path clean does, so that bulk and per-path never dispose of the same tree differently.
26. As an operator, I want an interrupted apply to leave completed rows settled and the rest untouched, with the old fingerprint refused, so that recovery is a re-plan and nothing more.
27. As an operator, I want paths with spaces, globs, or control bytes rendered safely in the plan and help, so that a pasted remedy works and a control byte cannot forge a line.
28. As an operator, I want a recorded path that is now a FIFO, socket, device, or dangling symlink retained without any git call, so that the sweep cannot block.
29. As an agent, I want the bare plan to be non-destructive, so that I can inspect before applying.

**Making release mechanical** — Line: `fable` / high. Guidance prose compounds through
every session (the profile's authoring routing); the anchor rows are the gate-observable half.

30. As a coordinator, I want `craft-delegate` to tell me to run `bench worktree release` for each accepted independent slice, so that release is a step, not a memory.
31. As `/bench-final-check`, I want the landed sweep as a named step before the landing report, so that a build's trees are retired while they are still landed-unreleased.
32. As a reviewer, I want the sweep result carried in the landing report, so that the phase's evidence includes what it retired.
33. As a maintainer, I want both sentences anchored, so that their loss turns the gate red.
34. As an operator, I want a changelog entry for the classification and the verb, so that the shipped surface is discoverable.

**Explicitly not wanted** (see Out of scope)

35. As a maintainer, I do not want a request-derivation override or bulk `release`, so that the request digest stays the ownership proof.
36. As a maintainer, I do not want automatic removal of landed trees at session start, so that the automatic path stays retain-only.
37. As a maintainer, I do not want `recovered` rows or ref inventory in this sweep, so that FT98/FT199 decisions stay open.

Story partition: 1–29 share the landed classifier and the `worktree_cleanup` plan seam;
30–34 name the verb. One bundle, chosen by the reviewer (map #5).

## Implementation decisions

- **One classifier, three consumers.** A single function in `internal/worktree`
  answers "is this assignment landed?" from the assignment record, its landedness
  proof, and its lease probe; `list`, the resume summary, and the `--landed` selector
  all call it. The predicate is exactly `state=active ∧ landed=true ∧ lease≠live`,
  enumerated over what the producers emit:
  - state: only `active` qualifies; `cleanup-pending`, `recovered`, and `complete`
    never do, whatever their landedness or lease;
  - landedness: only a proof of `true` qualifies; `false`, an unresolvable default
    branch, an empty branch, and a proof that errored (`unknown:<err>`) all read as
    not landed and leave the row plain `active`;
  - lease: the classifier reads the lease probe; `none`, `dead`, and `unknown` (a
    regular lease file with unparseable content) all satisfy `lease≠live`; only
    `live` fails it. Unknown-lease rows are therefore *landed* for counting and
    advertising, and the bulk planner separately retains them as `uncertain`
    (map #6). A lease file the per-path planner cannot read is a planner error, and
    the bulk row carries that error as `retain`/`uncertain` exactly as the automatic
    planner already wraps a per-path error; no fixture is promised for it (see Won't
    handle).
  Dirty tracked state, ignored residue, and nested-repository state never enter the
  predicate — they shape the row's plan action.
- **`landed` is a `CleanupReason`.** The automatic planner returns retain reason
  `landed` for a landed assignment ahead of every other retain reason it would
  otherwise report for the same tree (`active`, `orphaned`, `ignored`, `dirty`,
  `uncertain`), so the resume summary's `retained` counts read `landed=N` for exactly
  the set the bulk verb would plan. `live-lease` stays a distinct count because a live
  lease fails the predicate. The automatic path still removes nothing for these rows
  (FT148 design; map #4).
- **Orphan interplay.** The age-based orphan candidate list excludes landed
  assignments; the summary prints one line advertising the bare `bench worktree clean
  --landed` invocation when the landed count is above zero, in the same
  operator-safe form the orphan lines use (plan half only, never a discard flag).
- **`list` advertising.** `bench worktree list` renders one AXI action `bench
  worktree clean --landed` when at least one row is landed, in addition to the
  existing per-row `path`/`exec` actions for active rows. No new column: the row's
  `state`, `landed`, and `lease` columns already carry the three facts.
- **Grammar.** `usage.WorktreeClean` becomes `bench worktree clean [--discard-ignored]
  [--discard-branch] [--full] (<path> | --landed) [--apply <fingerprint>]`. `--landed`
  with a path operand (either order), `--apply` on an empty landed set, and an
  `--apply` value that is not 64 lowercase hex characters are the existing invocation
  error (exit 2, `worktree_cleanup` error row naming the usage). The `.bench/BENCH.md`
  CLI-inventory sentence follows the constant.
- **Selection is repository-wide.** The selector reads every ledger assignment,
  whatever request, label, or session created it, and keeps those passing the
  classifier, sorted by assignment id. Before planning, each selected path is
  shape-classified (`ClassifyPathShape`); only a checkout directory reaches the
  per-path planner. Any other shape — absent, dangling symlink, non-directory, decayed
  directory, special metadata (FIFO, device, socket, symlinked `.git`), unknown — is
  planned `retain`, reason `uncertain`, with the shape name in `detail`, and no git
  command is invoked against that path.
- **Bulk plan composes per-path planning.** Each checkout row runs the existing
  explicit per-path planner with the given options. Every per-path fingerprint binds
  the whole lifecycle ledger, so it cannot be frozen across sibling rows: settling
  row one rewrites the ledger row two was planned against. The **set fingerprint**
  therefore binds, per row in id order, (assignment id, target, planned action, HEAD
  OID, tracked state, ignored count, lease probe state), plus the three option flags
  and a distinct version tag. Any drift — a new landed row, a row that stopped
  qualifying, a HEAD that moved even while staying landed, a tree that got dirty, a
  lease that went live — changes it.
- **Plan output.** The bare form prints the existing `worktree_cleanup` table with one
  row per selected assignment; each row's `fingerprint` column carries the set
  fingerprint, and per-row `action`, `tracked`, `ignored`, `recovery`, and `detail`
  are the per-path planner's verdict (or the shape retention above). Retained rows
  keep their reason in `detail`. An AXI help block follows: `bench worktree clean
  --landed [flags] --apply <fingerprint>` when at least one row's action removes, and
  one `bench worktree clean <path>` line per retained row, rendered through the
  existing action quoting so a path with spaces or glob characters is pasteable and a
  path with control bytes is replaced by a pointer rather than emitted. The bulk plan
  never plans a preserving removal: a row whose per-path action would preserve
  (`recover-remove`, or `discard-remove` with non-clean tracked state, or a detached
  registration) is retained with reason `dirty` and the per-path remedy (map #6).
  Zero selected rows print the empty table, no help, exit 0.
- **Apply.** `--apply <fp>` re-selects and re-plans the whole set, refuses with the
  existing stale-fingerprint error and removes nothing if the recomputed set
  fingerprint differs, then applies each removable row in id order: it re-plans that
  row immediately before its transaction, recomputes that row's digest tuple (the same
  seven fields the set fingerprint bound for it), refuses the row and stops the apply
  if the tuple differs from the validated one, and otherwise passes the fresh per-path
  fingerprint into the existing per-path cleanup transaction, whose own re-plan refuses
  the row as stale if it moved in between. A row can therefore never execute a plan
  the set fingerprint did not bind — in particular never a preserving one. Each transaction's terminal step marks the assignment
  `complete` so the transaction's own completion deletes the record — the ledger is
  settled inside the apply (map #8), the same seam release uses. Retained rows are
  reported unchanged. Exit 0 when every removable row completed. The two mid-apply
  failure shapes this spec fixes are LR14's (a lifecycle fault after the first row's
  transaction completes) and LR20's (a later row's tuple drifted after set
  validation): both exit non-zero with every completed row already settled, the
  drifted or later rows untouched, and the original fingerprint refused on re-run; no
  other mid-apply failure behavior is promised here.
- **Branch disposition matches per-path clean** (map #9 as re-closed): a branch the
  tool's own landedness proof licenses is deleted with its tree. Because the selector
  admits only rows whose proof is `true`, every bulk row's branch is already licensed
  for deletion; `--discard-branch` therefore composes but changes no bulk row's action
  or branch outcome — its only visible effect is the plan row's `detail` naming the
  operator assertion, exactly as per-path clean renders it today. A branch the
  derivation cannot prove landed is never selected by `--landed` and stays on the
  per-path verb. `--discard-ignored` and `--full` pass through to every row's planner
  unchanged.
- **Interruption.** An apply killed between rows leaves completed rows settled and the
  rest untouched; the same fingerprint is then refused as drift and the operator
  re-plans. No set-level receipt is introduced — the per-row receipts are the
  recovery state.
- **Prose.** `craft-delegate`'s acceptance sentence in `## Verifying the done-claim`
  becomes: "Acceptance closes an independent worktree after its slice lands: the
  coordinator runs `bench worktree release --request <opaque-id> <path>` for it."
  `/bench-final-check`'s post-merge tail in `## Exit handoff` replaces "leftover
  worktrees and scratch branches go through `bench worktree clean`" with "scratch
  branches go through `bench worktree clean`; leftover worktrees are retired by `bench
  worktree clean --landed`: run the plan, apply it, and carry the plan and apply result
  in the landing report".
- **Anchors.** Three rows join `internal/anchors` in the group `checkWorkflowAnchors`
  evaluates (`AfterImplementSpec`), exactly:
  - `RequireInSection`, file `.agents/commands/bench-final-check.md`, section
    `Exit handoff`, needle "leftover worktrees are retired by `bench worktree clean
    --landed`: run the plan, apply it, and carry the plan and apply result in the
    landing report", diagnostic ".agents/commands/bench-final-check.md post-merge
    tail dropped the landed worktree sweep step";
  - `ForbidInSection`, file `.agents/commands/bench-final-check.md`, section
    `Exit handoff`, needle "leftover worktrees and scratch branches go through",
    diagnostic
    ".agents/commands/bench-final-check.md still routes leftover worktrees to a bare
    per-path clean";
  - `RequireInSection`, file `.agents/skills/bench-craft-delegate/SKILL.md`, section
    `Verifying the done-claim`, needle "Acceptance closes an independent worktree
    after its slice lands: the coordinator runs `bench worktree release --request
    <opaque-id> <path>` for it", diagnostic
    ".agents/skills/bench-craft-delegate/SKILL.md dropped release-at-acceptance".
- **Changelog.** One `Added` entry under `[Unreleased]` for the classification and the
  bulk verb (operator-visible; `craft-synthesis` policy).
- **Bootstrap authority is not applicable.** No new trusted-execution or refusal-
  before-execution claim; the verb runs the already-built binary.

## Testing decisions

- The classifier and the automatic planner are exercised through the real
  session-start command (`ResumeCleanCommand`) and `ListCommand` over a fixture
  repository holding assignment worktrees in each partition the producers emit:
  landed with no lease, dead lease, live lease, unknown lease (regular lease file
  with unparseable content); landed aged past `bounds.AssignmentStale`; non-landed
  active; non-landed aged; landed with undeclared ignored residue; landed dirty;
  landed under an unresolvable default and under an erroring proof; and
  `cleanup-pending`, `recovered`, `complete` records over landed branches. Prior art:
  `orphan_sweep_test.go`, `list_actions_test.go`, `resume_test.go`.
- The bulk verb is exercised through `CleanCommand` (bare then `--apply`) on the same
  fixture family, asserting the table rows, help block, exit codes, tree presence,
  branch presence, and the ledger via `ListCommand` after apply. Drift rows mutate the
  pool between plan and apply. Prior art: `clean_operand_test.go`,
  `clean_branch_test.go`, `recovery_retry_test.go` (fault seam).
- Interruption uses the existing lifecycle `Fault` seam (`hit(localFault, step)`) to
  fail after the first row's completion, then re-plans; no sleep or polling.
- Anchor rows are proven red the way every workflow anchor in the kit is: one canary
  fixture per row under `tests/canary/workflow-guidance-anchors/<name>/` with `BASE`
  naming the live file, `MUTATE.json` replacing the sentence (or, for the
  `ForbidInSection` row, re-inserting the retired one), and `EXPECT` holding the exact
  diagnostic; the existing universal fixture-bite proof
  (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) asserts each bites after the
  mutation and clears after restoration.
- Gate seam: `bench gate` runs the Go suite (`internal/worktree`, `internal/anchors`,
  `cmd/bench` help iteration), the anchor evaluation over the edited prose, and the
  canary fixture-bite proof.

### Seam diagram

    session start / `bench worktree list` / `bench worktree clean --landed`
        │
        ▼
    ledger assignments + git landedness + lease probe
        │
        ▼
    [ landed classifier ] ──▶ resume summary `retained landed=N` + remedy line
                          ──▶ list AXI action `clean --landed`
                          ──▶ selector ─▶ shape check ─▶ [ per-path explicit planner ×N ] ─▶ set fingerprint
                                              │
                          `--apply <fp>` ─▶ re-plan set, compare set fp ─▶ per row: fresh plan ─▶ [ cleanup transaction ] ─▶ record settled
                      ◀ tests attach here: run the commands over a fixture pool; read table rows, help,
                        exit code, tree and branch presence, and `list` afterwards

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| LR1 | 1 | Over a pool with one landed assignment (active, merged branch, no lease) and one non-landed active, `bench resume-clean` prints `retained active=1 landed=1`, removes nothing, and both trees remain present. | `ResumeCleanCommand` on a fixture pool | Not observed: today prints `retained active=2`. | A classifier that never fires, or one that removes, fails the count or the presence check. |
| LR2 | 6 | A landed tree with a live lease counts `live-lease=1` and not `landed`. | `ResumeCleanCommand` | Not observed: `lifecycle_test.go` probes reclaim through `Release`, never the summary count (`TestResumeSummaryLeasePartition` will assert the `live-lease=1` token). | Pins the `live` side of `lease≠live`. |
| LR21 | 7 | A landed tree with a dead lease counts `landed=1`. | `ResumeCleanCommand` | Not observed: `TestResumeSummaryLeasePartition` will assert `landed=1` for a dead-pid lease file. | A predicate reading only `none` fails. |
| LR22 | 8 | A landed tree whose regular lease file holds unparseable content counts `landed=1`. | `ResumeCleanCommand` | Not observed: `TestResumeSummaryLeasePartition` will assert `landed=1` for an unparseable lease. | Pins that `unknown` satisfies `lease≠live`. |
| LR3 | 9 | A landed assignment with undeclared ignored residue, and one with dirty tracked state, both count as `landed`, not `ignored`, `active`, or `orphaned`. | `ResumeCleanCommand` | Not observed: today the automatic planner reports `ignored` for the residue tree and `active`/`orphaned` for the dirty one. | Residue must not hide a landed tree from the count the bulk verb acts on. |
| LR4 | 5 | A landed assignment aged past `bounds.AssignmentStale` counts as `landed` and prints no orphan candidate line. | `ResumeCleanCommand` | Not observed: today `sweepOrphanAssignments` lists it (`TestOrphanSweepSkipsLanded` will assert zero orphan lines). | Catches double-listing. |
| LR23 | 5 | A non-landed assignment aged past `bounds.AssignmentStale` still prints its orphan candidate line. | `ResumeCleanCommand` | Already covered: `TestResumeSummaryListsOrphans` in `orphan_sweep_test.go` asserts the orphan line for an aged active row; retained as the control that the landed exclusion does not swallow it. | The exclusion must be landed-only. |
| LR24 | 4 | A summary with landed > 0 prints exactly one line advertising `bench worktree clean --landed` and never `--discard-ignored`. | `ResumeCleanCommand` | Not observed: `TestResumeSummaryAdvertisesLandedSweep` will assert the exact line and the flag's absence. | A remedy line naming a destructive flag manufactures the next residue. |
| LR5 | 10 | With no resolvable default branch, a landed-looking assignment counts as `active` and `list` renders no `clean --landed` action. | `ResumeCleanCommand`, `ListCommand` | Not observed: `TestLandedClassifierUnknownDefault` will assert `active=1` and no action. | An unknown proof that classifies as landed lets the bulk verb select it. |
| LR25 | 10 | With a branch whose landedness proof errors (ref deleted while registration and record remain), the assignment counts as `active` and `list` renders no `clean --landed` action. | `ResumeCleanCommand`, `ListCommand` | Not observed: `TestLandedClassifierErroredProof` will assert `active=1` and no action. | An errored proof is not landed. |
| LR26 | 11 | `cleanup-pending`, `recovered`, and `complete` records over landed branches are never counted `landed` and render no `clean --landed` action. | `ResumeCleanCommand`, `ListCommand` | Not observed: `TestLandedClassifierNonActiveStates` will assert no `landed=` token and no action for each state. | Only `active` qualifies. |
| LR6 | 3 | `bench worktree list` renders exactly one `bench worktree clean --landed` help action when at least one row is landed, none when zero, and keeps the per-row `path`/`exec` actions for every active row. | `ListCommand` | Not observed: `list_actions_test.go` asserts only path/exec actions. | A missing or duplicated action fails the exact help block. |
| LR7 | 14, 15 | Bare `clean --landed` over three landed assignments created under three different requests and labels (two clean, one dirty) prints three `worktree_cleanup` rows in assignment-id order with actions `remove`, `remove`, `retain` (reason `dirty`) and no `recovery` ref planned. | `CleanCommand` | Not observed: `--landed` is an invocation error today (`TestCleanLandedPlansRepositoryWideSet` will assert the three rows). | Fails a build-scoped selection, a preserving plan, or an unordered set. |
| LR27 | 16 | Every row of the LR7 plan carries the same set fingerprint in its `fingerprint` column. | `CleanCommand` | Not observed: `TestCleanLandedPlanSharesOneFingerprint` will assert one distinct value across rows. | Fails per-row fingerprints. |
| LR28 | 29 | The LR7 plan's help block names the `--apply <fp>` invocation and one `clean <path>` line for the dirty row, and the bare plan removes nothing. | `CleanCommand` | Not observed: `TestCleanLandedPlanAdvertisesApplyAndRemedies` will assert the help lines and tree presence. | A plan that mutates, or omits the remedy, fails. |
| LR8 | 19, 25 | `clean --landed --apply <fp>` on LR7's pool removes both clean trees, reports their rows as `removed`, deletes their proven-landed branches, leaves the dirty tree present and its row `retain`, exits 0, and `bench worktree list` immediately shows no row for the removed assignments and still shows the dirty one. | `CleanCommand`, `ListCommand` | Not observed. | Catches an apply that removes but leaves the ledger record, freezes a ledger-bound fingerprint across rows, or touches the retained row. |
| LR9 | 17 | Between plan and apply, (a) a new landed assignment appears, (b) a planned tree gains an uncommitted file, (c) a planned lease goes live, (d) a planned branch fast-forwards to the default tip (still landed, still `remove`, new HEAD OID): each makes `--apply <fp>` exit 1 with the stale-fingerprint error, remove nothing, and settle nothing. | `CleanCommand` | Not observed. | Membership, action, and OID must all be bound; a digest over membership alone passes (d). |
| LR10 | 8, 20 | A landed assignment whose regular lease file holds unparseable content is planned `retain` with reason `uncertain`, listed with its per-path remedy, counted `landed` by LR2's summary, and left in place by apply. | `CleanCommand`, `ResumeCleanCommand` | Not observed. | Pins map #6 for the uncertain half, distinct from dirty, and the classification/plan split for the same tree. |
| LR11 | 21 | With zero landed assignments, bare `clean --landed` prints the empty `worktree_cleanup` table, no help block, and exits 0 on repeated runs. | `CleanCommand` | Not observed: `TestCleanLandedEmptySetExitsClean` will assert the empty table and exit 0 twice. | Fails an implementation that treats empty as an error. |
| LR29 | 22 | `clean --landed --apply <any 64-hex>` over zero landed assignments is the invocation error, exit 2. | `CleanCommand` | Not observed: `TestCleanLandedApplyOnEmptySetRefused` will assert exit 2 and the usage row. | No fingerprint was issued for nothing. |
| LR12 | 23 | `clean --landed <path>` and `clean <path> --landed` are the invocation error, exit 2, naming the new usage. | `CleanCommand` | Not observed: `TestCleanLandedRefusesPathOperand` will assert exit 2 and the usage text (`clean_operand_test.go` has no selector case). | The selector and the operand must be exclusive. |
| LR30 | 23 | With `--landed`, `--apply` values that are short, long, non-hex, or uppercase hex are the invocation error, exit 2. | `CleanCommand` | Not observed: `TestCleanLandedRefusesMalformedFingerprint` will assert exit 2 for each shape (`clean_operand_test.go` has no malformed-fingerprint case). | The fingerprint grammar must hold for the selector form. |
| LR31 | 23 | `bench worktree --help` output contains the full new `usage.WorktreeClean` grammar string. | `cmd/bench` help route | Not observed: `command_registry_test.go` asserts only the `bench worktree clean` prefix (`TestWorktreeHelpNamesLandedGrammar` will assert the full string). | A stale usage constant hides the selector. |
| LR13 | 24 | Over two removable rows where only the second has undeclared ignored residue, bare `--landed` retains the second (`ignored`) and `--landed --discard-ignored` plans it `discard-remove` with `--full` widening its ignored preview. | `CleanCommand` | Not observed: `TestCleanLandedModifiersReachEveryRow` will assert the second row's action under each flag set. | Modifiers must reach every row's planner, not only the first. |
| LR32 | 25 | Applying the LR13 pool removes both rows and deletes both proven-landed branches without `--discard-branch`, and with `--discard-branch` the outcome is identical while each removed row's plan `detail` names the assertion. | `CleanCommand` | Not observed: `TestCleanLandedBranchDispositionMatchesPerPath` will assert branch absence in both runs and the `detail` text. | Branch disposition must match per-path clean. |
| LR33 | 25 | A squash-landed branch the derivation cannot prove is not selected by `--landed`. | `CleanCommand` | Not observed: `TestCleanLandedSkipsUnprovableBranch` will assert the row is absent from the plan. | Selection admits only proof-licensed rows. |
| LR14 | 26 | With a lifecycle fault injected after the first row's transaction completes, apply exits non-zero, the first tree is removed and its record settled, the second tree remains, re-running `--apply` with the original fingerprint is refused as drift, and a fresh bare plan lists only the remaining row. | `CleanCommand` with the lifecycle `Fault` seam | Not observed. | Interruption must leave per-row truth, not a half-set the same fingerprint could re-apply. |
| LR15 | 30, 33 | The three anchor tuples above are registered exactly once each, and the live `.agents/commands/bench-final-check.md` and `.agents/skills/bench-craft-delegate/SKILL.md` satisfy all three on the clean tree. | `internal/anchors` registry, `checkWorkflowAnchors` on the live tree | Not observed: needles absent, tuples unregistered (`TestLandedRetirementAnchorTuples` will assert exactly-once registration). | A tuple registered twice or a live file missing its sentence fails. |
| LR34 | 31, 32, 33 | Each of the three anchor tuples has one canary fixture whose mutation (removing the required sentence, or re-inserting the forbidden one) makes the anchor evaluation emit exactly the registered diagnostic and whose restoration clears it. | `tests/canary/workflow-guidance-anchors/` fixtures through `TestEveryRetainedFixtureBitesThroughRegisteredOwner` | Not observed: fixtures absent. | The gate, not review, holds the prose; a needle planted outside its section fails the section-scoped tuple. |
| LR16 | 34 | `CHANGELOG.md` `[Unreleased]` gains one `Added` entry naming the `landed` classification and `bench worktree clean --landed`. | none | Not TDD-able: prose the gate does not grade. | Ticket acceptance only. |
| LR17 | 27 | A landed assignment whose path contains a space and `*` is planned `remove`, its per-path remedy (when retained in a dirty variant) and the apply invocation render quoted and pasteable, and apply removes it. | `CleanCommand` | Not observed. | The bulk table and help must survive the paths the pool actually produces. |
| LR18 | 27 | A landed assignment whose path contains ESC is planned `retain` (`uncertain`, unsafe control bytes), its `target` cell is the pointer form, no help line emits the byte, exit 0, and apply leaves it in place. | `CleanCommand` | Not observed: `TestCleanLandedControlBytePathRetained` will assert the pointer target, byte-free help, exit 0, and tree presence. | Asserts the refused control-byte half at the new sink. |
| LR35 | 27 | A landed assignment whose path contains a tab is rendered escaped by the table encoder as one row. | `CleanCommand` | Not observed: `TestCleanLandedTabPathRendersOneRow` will assert one parseable row. | Asserts the permitted control-byte half. |
| LR19 | 28 | A landed assignment whose recorded path is now a FIFO, one whose path is a dangling symlink, one whose path is a unix socket (the `internal/bounds` `requireSocket` pattern, capability-guarded), and one whose path resolves to the `/dev/null` device node (guarded as `internal/bounds` guards it), are each planned `retain` (`uncertain`, shape named), the command completes without opening the special file, and apply removes none. | `CleanCommand` with the existing capability guards | Not observed. | A selector that hands every ledger path to git blocks on the FIFO or socket. |
| LR20 | 18 | With a fault callback that dirties the second row's tree at the first row's terminal step and returns nil, apply settles the first row, refuses the second as drifted, removes and preserves nothing for it, exits non-zero, and the original fingerprint is refused on re-run. | `CleanCommand` with the lifecycle `Fault` seam | Not observed. | Per-row tuple comparison is what keeps a plan the set never bound — a preserving one — from executing. |

Not covered: story 2 — the existing `state`, `landed`, and `lease` columns are unchanged; LR6 and LR34 exercise what `list` adds.
Not covered: story 12 — a design constraint (no new ledger state); enforced by the absence of an `AssignmentState` edit, not by a behavior.
Not covered: story 13 — a design constraint (one classifier); LR1–LR6 and LR7 exercise the three consumers, and review verifies the single source.
Not covered: story 35 — bulk release is out of scope (map #3); no in-scope behavior to assert.
Not covered: story 36 — automatic removal is out of scope (map #4); LR1's "removes nothing" is the retained control.
Not covered: story 37 — recovered rows are out of scope (map #4); LR26 asserts they are never counted `landed`.

### Edge inventory

- Error path: LR9, LR10, LR12, LR14, LR18, LR19, LR20.
- Empty/absent: LR11 (zero landed rows), LR5 (absent default, absent ref), LR19 (absent tree behind a dangling link).
- Boundary values: LR4 (age past `AssignmentStale` versus landed precedence); LR2 (lease live/dead/unknown); LR12 (fingerprint length).
- Malformed input: LR12 (selector plus operand, malformed fingerprint), LR2/LR10 (unparseable lease).
- Interrupted/partial state: LR14, LR20.
- Re-run idempotency: LR9 (a spent fingerprint refuses); LR11 (repeated bare plan over an empty pool).
- Process-boundary lifecycle: LR8 reads the ledger through a second command (`list`) after apply.
- Hostile environment: LR17 (spaces, glob), LR18 (refused and permitted control bytes), LR19 (FIFO, dangling symlink, socket and device where creatable).
- Project checklist: spaces/globs LR17; control bytes in sinks LR18; an artifact-rewriting command reporting a fact it changed — LR8 asserts the post-write truth through `list`; absent vs empty LR5/LR11; special files and dangling links LR19; interrupt LR14; rerun LR9/LR11; fresh process LR8; destructive worktree state and plan/apply drift LR9/LR10/LR14/LR19.

Won't handle: a live symlink planted at an assignment path — the per-path planner
canonicalizes the operand, which then names a target the registration does not, and the
existing "target is not registered" retention holds (LR19's dangling case is the
surviving in-scope caller through the same shape check). Non-TTY stdin — the verb never
prompts. Deep cwd — `git.Root` discovery is unchanged. An unreadable lease file —
the per-path planner's read error is wrapped as `retain`/`uncertain` by the existing
error path, and no new code decides it — a permission-denied fixture would only run
behind `capability.Privilege` and would prove the existing error wrap, not this spec;
LR10's unparseable lease is the surviving in-scope caller of the same `uncertain`
disposition.

## Ownership fences

- `internal/worktree/` (production and test files; the coordinator commits per ticket)
- `internal/usage/worktree.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/` (three new fixture directories only)
- `cmd/bench/command_registry_test.go`
- `.bench/BENCH.md` (the CLI-inventory sentence for `bench worktree clean` only)
- `.agents/commands/bench-final-check.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `CHANGELOG.md`

## Out of scope

- Settling the ledger record when a request-less **per-path** `clean` removes a young
  active assignment (FT93b residue) — 2 edits, 1 gate run; a separate reviewer decision
  because it changes an established per-path contract.
- Bulk `release` or a request-derivation override (FT148 #1 stands).
- Automatic removal of landed trees at session start.
- `recovered` rows and their recovery refs (FT98, FT199).
- Repository-wide ref inventory and branch classification (FT199).
- A gate check on ledger state.
- Any mid-apply failure behavior beyond LR14's and LR20's fixed shapes.

## Further notes

First spec authored in the extensive-story shape (`craft-spec` remade on to-spec, 2026-08-16);
its 5 tickets were sliced under the tracer-bullet rule (`craft-tickets` remade on to-tickets)
after a 17-ticket cut under the retired smallest-independently-green rule was rejected by
the reviewer.
