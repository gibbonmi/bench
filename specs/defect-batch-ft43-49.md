# defect-batch-ft43-49

Status: implemented

Source: reviewer-directed batch drain off the 2026-07-08 platform assessment
(no per-defect decision maps — every defaulted decision below is **flagged ⚑
for post-hoc veto**).

## Problem

Seven verified defects from the 2026-07-08 assessment: two enforcement
fail-opens (FT43 Codex Stop-hook timeout, FT49 pre-push guarding a fabricated
default branch), one concurrency race (FT45 lease reclaim), one perpetual
no-op recommendation (FT44 salvage-branch sweep), one uninstall residual
(FT47 dangling CLAUDE.md), and two record-drift defects (FT46 ADR 0002
posture 5 now false, FT48 CHANGELOG missing two mandated entries with no
owner for the duty).

## Solution

Fix each at its existing seam, smallest diff per defect, staged FT43 first
(XS, enforcement fail-open). Records are amended, never reverted: the verdict
reuse and the sweep's keep-behavior are sound decisions whose descriptions
drifted.

## User stories

1. As a Codex-harness user with an armed shift, I want the Stop-hook gate run
   to never be killed by a hook timeout, so that the completion oracle cannot
   fail open under exactly one harness. ⚑ Default: **drop** the `timeout`
   keys from `.codex/hooks.json` (both hook blocks) for parity with Claude's
   settings, which set none — rather than raising to a guessed margin.
   Line: claude-sonnet-5 / low. This is a two-key JSON deletion at a known
   file with the conformance guard written in story 2.

2. As the reviewer, I want a conformance check asserting no Codex hook block
   sets a `timeout`, so that the fail-open cannot silently return in a future
   template edit. ⚑ Default: absence check, not a margin rule — a margin
   needs a maintained gate-runtime constant, which is a new decision.
   Line: claude-opus-4-8 / medium. Gate and conformance logic is the
   profile's cached mid-effort routing because a wrong oracle check is worse
   than a slow one.

3. As a reviewer running `bench worktree clean`, I want a salvage branch
   whose content already landed in the default branch deleted automatically,
   so that I never hand-inspect a branch that provably has nothing to
   salvage. ⚑ Default: widen the "landed" proof from ancestry to patch
   containment via `git cherry <default> <branch>` — all commits reported
   already-applied (`-`) means landed; any `+` commit keeps the branch. The
   existing unresolvable-default false-empty guard is preserved unchanged.
   Line: claude-opus-4-8 / medium. The proof-of-landed semantics carry real
   correctness risk (deleting a branch is destructive) beyond cheap plumbing.

4. As a reviewer reading `bench status`, I want a kept (genuinely un-landed)
   salvage branch to get an honest action instead of the no-op
   `bench worktree clean`, so that following the recommended action always
   changes something. ⚑ Default action text: `inspect salvage branch(es) —
   bench worktree clean keeps them` (exact wording implementable-free).
   Line: claude-sonnet-5 / low. One status row's action string once story 3
   defines which branches are kept.

5. As a crash-recovery reclaimer of a dead-pid lease, I want `Claim`'s
   takeover to verify the lease it renamed is the lease it judged
   reclaimable, so that two concurrent reclaimers cannot both win one
   worktree. ⚑ Default: after the takeover rename succeeds, read the renamed
   stale file and compare to the judged content; on mismatch the rename stole
   another reclaimer's fresh lease — rename it back (best-effort) and return
   false. Line: claude-opus-4-8 / high. Concurrency interleavings are the one
   place in this batch where the cheapest plausible diff is likely wrong.

6. As a future debugger, I want a two-reclaimer stress case in the
   concurrent-acquire contract suite, so that the "cannot both win"
   guarantee is a tested claim instead of a comment. Line:
   claude-opus-4-8 / medium. The existing concurrent-acquire test is the
   prior art and the natural home, but forcing the dead-pid interleaving
   deterministically takes care.

7. As a teammate reading ADR 0002, I want posture 5 amended to record the
   shipped verdict-reuse decision (exact tree-hash key, fresh-only, single
   writer), so that the decision record describes current behavior. Amend,
   don't revert — the reuse is sound. ⚑ Default: the amendment also names the
   existing regression tests as the posture's bite-proof; a new test row is
   added only if `internal/commit` coverage has no reuse-refusal case (the
   assessment says both directions are already tested — verify, don't
   assume). Line: claude-fable-5 / high. Doc authoring is the profile's
   cached leverage-override routing: ADR prose compounds through every
   session that reads it.

8. As the `/bench-update-kit` baseline, I want the two missing 2026-07-07
   promotion entries backfilled in CHANGELOG.md (commits `453599a` FT35/36
   rule edits; `2a72310` verification-probe rule), so that the next
   re-synthesis diffs against a complete record. Line: claude-fable-5 /
   high. Same cached doc-authoring routing as story 7.

9. As the reviewer, I want the CHANGELOG append duty anchored as an explicit
   step in the `/bench-what-next` drain checklist, so that the duty has an
   owner. ⚑ Default: a checklist line in the command prose (kit edit under
   `craft-synthesis`), not a conformance check — escalate to a check only if
   drift recurs, per the roadmap row's own framing. Line: claude-fable-5 /
   high. Command prose is kit guidance under the leverage override.

10. As a user whose repo's default branch is not `main`, linked before its
    remote existed, I want the pre-push hook to guard my real default
    branch, so that the backstop cannot silently never fire. ⚑ Default:
    resolve at push time — the hook template queries
    `git symbolic-ref refs/remotes/origin/HEAD` live and falls back to the
    baked branch token when unresolvable; no doctor row needed once the hook
    self-corrects. Line: claude-opus-4-8 / medium. Shell-hook semantics in
    the enforcement path are worth mid; the fix shape is small but the
    fail-open class is the worst kind.

11. As a user running `bench unlink`, I want a link-created CLAUDE.md removed
    with the rest of the install, so that no leave-behind imports
    just-deleted files. A pre-existing user CLAUDE.md is never recorded and
    never removed. README's leave-behind list is corrected in the same
    change. Line: claude-sonnet-5 / medium. Mechanical once the
    record-only-when-link-wrote-it decision above is fixed, with the safe-link
    gate layer already owning the round-trip fixtures.

## Implementation decisions

- **FT43:** delete both `timeout` keys from the kit's `.codex/hooks.json`;
  linked repos receive the same file (it is copied by the link plan, not
  templated), so no second copy to fix. Conformance gains a
  no-Codex-hook-timeout check.
- **FT44:** `sweepDelegateBranches` gains a second landed-proof after the
  ancestry test fails: `git cherry <default> <branch>` with every line `-`
  ⇒ delete (report as landed-by-content); any `+` ⇒ keep. The status
  worktree row splits its orphan action: landed branches disappear via the
  sweep; kept branches get the honest inspect action (story 4).
- **FT45:** `Claim` verifies identity post-rename against the judged bytes;
  mismatch → best-effort rename-back, return false. The doc comment's
  guarantee then matches the code.
- **FT46/FT48/FT49 records:** ADR amendment and CHANGELOG backfill are
  prose-only; the pre-push template adds live resolution with the baked
  token as fallback (the token stays, so `bench guards`' marker contract is
  untouched).
- **FT47:** `installClaudeMD` reports whether it wrote the file; only a
  link-written CLAUDE.md enters the manifest (as a recorded row), so a
  pre-existing user CLAUDE.md is never manifest-recorded and never removed
  by unlink. README's leave-behind list is corrected in the same change.

## Testing decisions

- Good tests here exercise external behavior at the shipped surfaces: the
  gate's conformance layer for template/docs shape, and throwaway fixture
  repos in the contract suites for link/unlink, sweep, status, lease, and
  hook behavior. No internals-poking.
- Prior art: `internal/contract/runtime` (worktree/lease/status fixtures),
  the safe-link gate layer and `internal/adopt/adopt_test.go` (link/unlink
  round-trips), `internal/conformance` (docs/template checks),
  `internal/contract/runtime`'s concurrent-acquire test (story 6's home).
- Gate: `.bench/gate.sh` (the project gate), green before each commit.

### Seam diagrams

**S1 — conformance layer (stories 2; docs shape for 7, 8)**

    trigger: bench gate → conformance phase
        │
        ▼
    tree under grade ──▶ [ TestRootConformance checks ] ──▶ pass / named failure
                              ◀ tests attach here: canary fixture plants the
                                defect (a timeout key), asserts the gate goes
                                red with the targeted substring

**S2 — worktree sweep + status (stories 3, 4)**

    trigger: bench worktree clean / bench status (fixture repo)
        │
        ▼
    repo with salvage   ──▶ [ sweepDelegateBranches   ] ──▶ deleted/kept report
    branches (landed,       [ + status worktree row   ] ──▶ status action text
    un-landed)
                              ◀ tests attach here: runtime contract fixture
                                builds both branch kinds, runs the CLI, asserts
                                stdout rows and branch existence after

**S3 — lease Claim (stories 5, 6)**

    trigger: two concurrent Acquire/Claim callers, dead-pid lease planted
        │
        ▼
    lease file (stale) ──▶ [ Claim takeover path ] ──▶ exactly one true
                              ◀ tests attach here: concurrent-acquire suite
                                forces the interleaving, asserts single winner
                                and lease content integrity

**S4 — safe-link round-trip (stories 10, 11)**

    trigger: bench link / bench unlink in a throwaway repo
        │
        ▼
    repo (with/without ──▶ [ link plan + manifest;  ] ──▶ files on disk,
    CLAUDE.md, remote,     [ unlink consuming it    ] ──▶ manifest rows,
    odd default branch)                                    hook content
                              ◀ tests attach here: adopt/link gate fragments
                                assert presence/absence after round-trip, and
                                a push against the fixture's real default
                                branch is blocked

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | no `timeout` key in any `.codex/hooks.json` hook block | S1 | story 2's conformance check red on today's tree before the edit | today's file has two `timeout: 30` keys, so the check fails until the deletion ships |
| 2 | conformance rejects a Codex hook timeout | S1 | canary fixture with a planted `timeout` key stays green before the check exists | a rotted or missing check lets the planted defect pass silently |
| 3 | content-landed salvage branch is deleted by the sweep | S2 | fixture: branch cherry-picked into default (non-ancestor, patch-contained) survives `bench worktree clean` today | ancestry-only proof keeps exactly this branch; the new proof deletes it |
| 3 | un-landed salvage branch is kept | S2 | fixture: branch with a unique patch gets deleted (must stay kept) | guards the destructive direction of the new proof |
| 3 | unresolvable default branch still refuses loudly | S2 | already covered (existing false-empty guard test) | regression guard on the preserved refusal path |
| 4 | kept salvage branch's status action is not `bench worktree clean` | S2 | fixture: status row for a kept branch still prints `bench worktree clean` | the no-op recommendation is the defect itself |
| 5 | second reclaimer of the same dead-pid lease returns false | S3 | two-reclaimer interleaving yields two `true` returns today | directly falsifies the "cannot both win" comment pre-fix |
| 6 | single-winner guarantee is gate-run | S3 | new stress case absent from the suite (test-presence is the signal; it must fail red against unfixed `Claim`) | locks the guarantee so the comment can't rot silently again |
| 7 | ADR 0002 posture 5 describes verdict reuse | S1 | not TDD-able (prose accuracy) — reviewer reads the amendment; existing reuse regression tests named in the posture | record drift is graded by the reviewer, not the gate |
| 8 | CHANGELOG contains both 2026-07-07 promotion entries | S1 | not TDD-able (prose) — reviewer checks the two entries against `453599a`/`2a72310` | the baseline is complete only if both adopted rules appear |
| 9 | drain checklist owns the CHANGELOG duty | S1 | not TDD-able (prose) — reviewer reads the added checklist step | duty anchoring is a process rule, deliberately not a gate check yet |
| 10 | pre-push blocks a push to the repo's *real* default branch (`master` fixture, linked pre-remote) | S4 | fixture: link before remote, set `master` default, push to master succeeds today | the baked `main` guard provably never fires on this fixture |
| 10 | no-remote repo still gets a working baked fallback | S4 | fixture: hook with unresolvable origin/HEAD blocks push to the baked branch | live resolution must not break the existing offline path |
| 11 | link-created CLAUDE.md is removed by unlink | S4 | fixture: link (no prior CLAUDE.md) → unlink leaves CLAUDE.md on disk today | the dangling import is the defect |
| 11 | pre-existing user CLAUDE.md survives unlink untouched | S4 | fixture: user CLAUDE.md present pre-link is absent or altered after unlink | guards the destructive direction of manifest-recording |

### Edge inventory

Walked against the profile's hostile-input checklist per behavior; each lands
as a row above or a **Won't handle** here.

- Branch names with control bytes reaching `toon.Table` via the status row —
  **already covered** by the existing git-sourced-text refusal contract; the
  new action text adds no new git-sourced field.
- SIGINT mid-`Claim` leaving a `.stale.<pid>` file — **Won't handle**: the
  existing reclaim path already treats an orphaned stale file as an ordinary
  dead lease on the next scan; no new state is introduced.
- Re-run idempotency: relink must not duplicate the CLAUDE.md manifest row —
  row folded into the FT47 fixtures (relink asserted single-row).
- Absent vs present-but-empty CLAUDE.md at link time — empty file is not
  `legacyClaudeMD()`, so it is user content: covered by the survives-unlink
  row.
- Repo whose default branch contains a `/` (e.g. `release/2026`) in the
  pre-push live resolution — covered implicitly: `symbolic-ref --short`
  output is used verbatim; fixture uses `master`, and a slash adds no new
  parse. **Won't handle** a dedicated slash fixture — same code path,
  one-clause risk.
- `git cherry` on a branch identical to default (zero commits out) — treated
  as landed (empty cherry output = nothing unique); asserted inside the
  landed-branch fixture.
- Missing `git` / no `readlink -f` — **Won't handle**: unchanged from the
  surrounding hook's existing dependency posture.

## Out of scope

- **Gate-wall margin rule for hook timeouts** — a maintained runtime budget
  is a separate capability (a measured constant with its own drift problem);
  ~6 edits, ~4 gate runs. Parked unless a timeout ever legitimately returns.
- **Doctor row comparing baked vs live default branch** — superseded by
  push-time resolution; only worth building if the live-resolution default
  is vetoed. ~4 edits, ~3 gate runs.
- **Conformance check on CHANGELOG completeness** — deliberately deferred by
  story 9's own rule: escalate only if the drift recurs. ~5 edits, ~3 gate
  runs.
