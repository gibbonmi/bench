# Transactional managed-asset lifecycle (FT84)

Status: staged

> Map provenance: this spec is compiled from `decisions/transactional-asset-lifecycle.md`, whose map was recorded under the `/bench-write-spec` **same-session exception** — every ticket restates a decision the reviewer already approved (the FT84 roadmap row through reviewed drains, and the FT76 grill of 2026-07-21). The map, and this spec's three Handoff-flag-7 proposals, are **flagged for reviewer veto** rather than re-asked.

## Problem

`bench link`, relink, and `bench unlink` are not transactional over the managed-asset set, and three concrete gaps follow.

- **Link writes sequentially with no rollback.** `Link` runs `writeAgentsFile` → `installClaudeMD` → `installPlan` (a per-file `os.Remove`-then-write loop) → `installGitHook`, in that order (verified `internal/adopt/link.go` `Link`, lines 53–68; `installPlan` 308–331; `installPlannedFile` 333–353). A failure partway leaves the already-written files on disk. The manifest is written **last** (`link.go:330`), so after a mid-sequence failure the on-disk files have no manifest rows, and a later relink classifies them as project-owned conflicts (`preflightLink` treats an un-manifested existing file as a project-owned conflict, `link.go:213–222`). A failed link does **not** leave the repo as it found it.

- **Relink never reconciles removed assets.** `installPlan` builds manifest rows only from the *current* kit plan and rewrites the manifest wholesale (`link.go:309–319`); nothing in `link.go` reads the prior manifest to find assets the new plan dropped. Relinking against a kit with fewer or renamed assets leaves the stale files on disk **and** drops them from the manifest — an active-but-unowned skill, exactly the state map #3 forbids.

- **Unlink is silent about residuals and speaks only prose.** `Unlink` returns `0` unconditionally (`internal/adopt/unlink.go:57`); the `manifestKept` flag is set when files are kept or refused but never changes the exit code (`unlink.go:126–131`). The residual report is human prose — `refused: <rel>` / `kept (modified): <rel>` lines (`writeUnlinkReport`, `unlink.go:272–311`) — not a machine-readable list. Link's conflict report is likewise prose on stderr (`link.go:186–223`), and link **aborts wholesale** on any conflict, writing nothing (`preflightLink` returns false → `Link` returns 1 before any write, `link.go:50–52`).

A consumer (and FT76's bootstrap, which composes these seams) cannot rely on link/unlink being atomic, reconciling, or machine-parseable.

## Solution

Make the whole lifecycle transactional over the managed-asset set, entirely inside `internal/adopt` (the wrapper dispatch is untouched; routes already exist).

- **Link and relink stage then promote.** Build the complete write set, stage it, preflight it, sync where durability matters, then atomically promote or roll back. A structural preflight failure or an I/O fault mid-promotion leaves the repo byte-identical to how link found it. Git hook and `.claude/settings.json` are part of the staged set, so they compose or fail *before* any partial write — never a semantic merge (map #5: never merge, never overwrite).
- **Relink reconciles old and new manifests.** Removed clean assets leave; a removed asset the project modified stays owned (kept on disk, kept as a manifest row) and is reported; no stale skill becomes active-but-unowned.
- **Conflicts are partial, not fatal.** Non-conflicting assets are written; a project-owned or locally-modified collision is preserved as project-owned, reported in a machine-readable conflict list, and returns a distinct partial status.
- **Unlink is honest about residuals.** It returns a nonzero partial status when residuals remain and emits a machine-readable residual list, and it never recommends raw deletion of an executable Bench does not own.

## User stories

1. As an adopter, I want a link that hits an I/O fault mid-write to leave my repo exactly as it found it, so that a failed install never leaves a half-written `.bench/` footprint I have to clean up by hand. Line: gpt-5.6-terra / medium. The promote-or-rollback engine has a wrong-but-green failure mode — a rollback that passes the happy path but corrupts on a fault at an untested promotion step — so it earns a tier over pure plumbing even though the gate observes the outcome.

2. As an adopter on durable media, I want staged writes synced before promotion, so that a crash right after link reports success cannot leave the managed set torn. Line: gpt-5.6-terra / medium. Durability sequencing is part of the same engine and shares its correctness risk, so it routes with story 1.

3. As an adopter upgrading or downgrading the kit, I want relink to reconcile the old and new manifests so that renamed or removed assets from the previous kit version leave and are never left active-but-unowned. Line: gpt-5.6-terra / medium. The removed-versus-modified-versus-clean classification is the reconciler's core judgment and is easy to get subtly wrong while still passing a happy-path relink.

4. As an adopter who edited a kit asset the new kit version drops, I want that modified asset kept on disk and kept owned in the manifest and reported, so that a downgrade never silently deletes my work or orphans it. Line: gpt-5.6-terra / medium. This is the reconciler's most error-prone branch and shares its tier with story 3.

5. As an adopter with one project-owned file colliding with a managed asset, I want link to write everything else, preserve my file, and tell me exactly what it skipped with a distinct partial status, so that one conflict does not block the entire install. Line: gpt-5.6-terra / medium. This changes link's current whole-abort contract and interleaves with the transaction engine, so it carries the semantic tier.

6. As a tool consuming link and unlink output, I want the conflict list and the residual list emitted as a machine-readable block distinct from the human prose, so that FT76's bootstrap can parse what was preserved without scraping sentences. Line: gpt-5.6-luna / medium. The list is rendered from the existing verdict struct through the shared `internal/toon` emitter, so it is mechanical once the struct exists.

7. As an adopter running unlink on a repo with residuals, I want a distinct nonzero exit status when anything is kept or refused, so that a script driving unlink can tell a clean removal from a partial one. Line: gpt-5.6-luna / medium. The ownership logic already exists; this adds the exit-code branch and updates the contracts that assert the old exit 0.

8. As an adopter, I want unlink to never print a raw-deletion recommendation for an executable Bench does not own, so that following the report can never destroy a foreign binary. Line: gpt-5.6-luna / low. Current unlink already emits no deletion advice, so this pins the invariant against regression rather than building new behavior.

9. As an adopter, I want repeated link/unlink cycles to converge idempotently, so that re-running either command leaves no accumulating residue. Line: gpt-5.6-luna / low. Single-cycle idempotence is already covered; this extends the matrix to a second full cycle.

10. As the contract suite, I want a test-only fault-injection seam on link's promotion loop, so that the I/O-failure contract can drive a deterministic mid-transaction failure through the built wrapper. Line: gpt-5.6-luna / low. The seam is a single env-gated branch in the promotion loop, inert unless the variable is set.

## Implementation decisions

**Module boundaries (Handoff #1).** All work lands in `internal/adopt` (`link.go`, `unlink.go`, `manifest.go`, plus a new transaction/reconcile unit). The wrapper dispatch and per-asset writers are thin; the transaction engine (stage → preflight → promote/rollback) and the manifest reconciler are the deep modules (Handoff #3). Both are **single-sourced**: link and unlink share one reconcile/verdict path and one machine-readable emitter — no duplicated policy between the two paths (code standard: one source per fact).

**Atomicity mechanism — per-file write-ahead journal (proposal — Handoff flag 7).** Managed targets are scattered and interleaved with project files (`AGENTS.md`, `CLAUDE.md`, `.bench/…`, `.agents/…`, `.claude/…`), so a single staged-directory swap is not available — you cannot swap the whole repo. The engine therefore stages each planned file beside its destination, preflights the complete set, then promotes each by atomic rename, recording a journal entry per target: for a fresh destination, rollback removes it; for a destination that pre-existed as managed-clean (including the bespoke `AGENTS.md`/`CLAUDE.md` in-place edits), the original bytes are saved and rollback restores them. On any error during promotion, the journal is replayed in reverse to restore the pre-link state, then link exits nonzero. The git hook install and `settings.json` write join this journaled promotion rather than running after it.

**Manifest reconciler.** On relink, read the prior manifest and diff its rels against the new plan. For a rel present in the old manifest but absent from the new plan: if the on-disk file still matches the old hash, remove it and drop its row (it leaves); if it was modified, keep it on disk, keep its row (still owned), and report it as a preserved residual. New and unchanged assets promote as today; the kit version is restamped. Reconciliation reuses `resolveInside`/traversal refusal so a hand-edited old manifest can never name a removal target outside the repo.

**Conflict posture.** Preflight classifies each planned asset. A **structural** failure (missing kit asset, symlink parent, non-directory parent, unparseable `AGENTS.md`, foreign non-managed pre-push hook) aborts the whole transaction → nothing written, repo as found, exit 1. A **conflict** (destination is project-owned or modified-managed) does not abort: that asset is preserved and added to the conflict list, the non-conflicting set promotes atomically, and link exits with the partial status. `.claude/settings.json` needs no special-casing — it is a managed asset under this posture (preserve-or-write, never merge).

**Partial-status exit code — 3 (proposal — Handoff flag 7).** `0` = full convergence; `1` = hard error (not-a-repo, structural preflight failure, I/O fault after rollback, absent/unreadable manifest); `2` = usage; `3` = partial (link preserved conflicts / unlink kept or refused residuals). `3` is free for both commands and distinct from the AXI `0/1/2` convention, which link/unlink do not follow.

**Fault-injection seam — env hook `BENCH_LINK_FAULT` (proposal — Handoff flag 7).** The contract drives the built wrapper as a black box (the fixture harness runs a separate process), so a Go-level filesystem shim is not observable across the process boundary — an env var is. `BENCH_LINK_FAULT=<k>` forces the k-th promotion to return an I/O error and is inert when unset. This follows the established test-env-var precedent in this repo (`BENCH_KIT` drives `testLinkMetacharKitPath`; `BENCH_CANARY_INNER`, `BENCH_REPAIR`, `BENCH_OFFLINE` gate production paths). Flagged for reviewer veto against the filesystem-shim alternative.

**Machine-readable lists.** The conflict list (link) and residual list (unlink) render from the single verdict struct through `internal/toon` as a TOON block on stdout with `path` and `reason` columns, distinct from the human prose. One struct, one emitter — link and unlink do not each grow their own format.

**Existing contracts this spec revises (flagged for veto).** The exit-code and partial-apply changes alter contracts that currently assert the old behavior; the implementation updates their expectations to the new intended behavior (not a weakening to pass): `testLinkConflictWithoutManifest` and `testLinkModifiedManaged` (whole-abort/exit-1/no-manifest → exit 3, non-conflicting assets written, conflict list); `testUnlinkKeepsModified`, `testUnlinkRefusesTraversal`, `testUnlinkReportsBothRuns` (exit 0 with residuals → exit 3). The absent-manifest hard error (`testUnlinkAbsentManifest`) stays exit 1.

## Testing decisions

- **A good test here** drives the **built wrapper** black-box through the existing contract fixture harness and asserts only external state: exit codes, target-repo file and git state, manifest contents, and the TOON lists on stdout — never `internal/adopt` internals (Handoff #4).
- **Seams (Handoff #5).** No new seam: coverage extends the existing surface contracts `internal/contract/surface/link_test.go` (`TestLinkContracts`) and `unlink_test.go` (`TestUnlinkContracts`), plus the packed-artifact contract's cold-run leg in `artifact/artifact_test.go`. The gate runs all contract packages as one phase, so no seam is gate-invisible. Prior art: the `contract.NewFixture` harness, the swapped-`BENCH_KIT` kit-copy pattern in `testLinkMetacharKitPath`, and the tree-snapshot/`git status --porcelain` assertions in `testUnlinkDryRunNoWrites`.
- **Gate command:** the project gate, `.bench/gate.sh` (the contract phase runs the surface and artifact packages).

### Seam diagram

Link transaction + reconcile (one seam, driven black-box):

    trigger: contract fixture runs `bench link` (built wrapper), BENCH_KIT set, BENCH_LINK_FAULT optional
        │
        ▼
    old manifest + kit plan   ──▶  [ stage → preflight → promote / rollback  ]  ──▶  repo files, new manifest, exit code
    BENCH_LINK_FAULT=k        ──▶  [ + reconcile removed assets              ]  ──▶  TOON conflict list on stdout
                                        (fault forces promote step k to fail)
                                        ◀ tests attach here: snapshot tree before link, run link, then assert
                                          bit-identical-on-fault / converged files / manifest rows / exit / conflict list

Unlink residuals (one seam, driven black-box):

    trigger: contract fixture runs `bench unlink [--dry-run]` (built wrapper)
        │
        ▼
    link manifest + disk   ──▶  [ classify: removed / kept-modified / refused → sweep ]  ──▶  removed files, exit code
                           ──▶  [ render residual verdict via internal/toon           ]  ──▶  TOON residual list on stdout
                                    ◀ tests attach here: assert exit 3 when residuals remain,
                                      residual TOON block parseable, no raw-deletion recommendation

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | I/O fault mid-promotion leaves the repo byte-identical to pre-link and exits 1 | link contract (`BENCH_LINK_FAULT`) | new `testLinkRollsBackOnFault`: snapshot tree, link with fault at step k, assert exit 1 and tree diff empty — red on the current sequential writer | current code leaves the first k files on disk with no rollback, so the post-link tree diff is nonempty and the assertion fails |
| 1 | fault at the first vs last promotion step both roll back fully (boundary) | link contract | same test parametrized over k∈{1, len-1}: both leave an empty diff — red on current | current code's residue grows with k, so a late-k fault leaves the most residue and fails hardest |
| 2 | staged files are synced before promotion | link contract | **not TDD-able**: fsync/crash durability is unobservable from a black-box contract without a real crash; asserted by review, not the gate | no process-level assertion can distinguish an fsync'd promotion from an un-synced one |
| 3 | relink across kit versions writes added assets, removes clean dropped assets, restamps version | link contract (swapped `BENCH_KIT`) | new `testLinkReconcilesAcrossKitVersions`: link kit A (asset X), relink kit B (X absent, Y added); assert Y present, X absent from disk and manifest, version=B — red on current | current `installPlan` never removes X, so X survives on disk and drops from the manifest; the "X absent from disk" assertion fails |
| 4 | a modified dropped asset is kept on disk, kept as a manifest row, and reported | link contract (swapped `BENCH_KIT`) | new `testLinkReconcileKeepsModifiedRemoved`: edit X locally, relink kit B; assert X on disk with the edit, X still a manifest row, X in the conflict/residual list — red on current | current code drops X from the new manifest while leaving it on disk unowned; the "still a manifest row" assertion fails |
| 5 | project-owned collision: preserve the file, write the rest, exit 3, emit conflict list | link contract | updated `testLinkConflictWithoutManifest`: assert exit 3, project file intact, manifest present for non-conflicting assets, conflict list names the path — red on current | current code exits 1, writes nothing, and writes no manifest; the exit-3, manifest-present, and conflict-list assertions all fail |
| 5 | modified-managed collision on relink: preserve, still converge the rest, exit 3 | link contract | updated `testLinkModifiedManaged`: assert exit 3 and conflict list names the modified file, file intact — red on current | current relink exits 1 whole-abort; the exit-3 and conflict-list assertions fail |
| 6 | conflict/residual lists are a machine-readable TOON block distinct from prose | link + unlink contracts | new `testUnlinkResidualListMachineReadable` (and link twin): assert a parseable TOON block with `path`/`reason` columns — red on current | current output is only prose `refused:`/`kept (modified):` lines; the structured-block parse fails |
| 7 | unlink exits 3 when it keeps a modified file or refuses a row | unlink contract | updated `testUnlinkKeepsModified` / `testUnlinkRefusesTraversal` / `testUnlinkReportsBothRuns`: assert exit 3 — red on current | current `Unlink` returns 0 unconditionally; the exit-3 assertion fails |
| 8 | unlink never recommends raw deletion of a non-owned executable | unlink contract | **already covered** — current unlink emits no deletion advice; new `testUnlinkNoRawDeletionAdvice` pins it as a regression guard, not a red-first row | a future report change that added `rm <foreign-exe>` advice would fail the "report contains no deletion recommendation" assertion |
| 9 | two full link/unlink cycles converge with zero residue | link + unlink contracts | new `testLinkUnlinkMatrix`: link;unlink;link;unlink and relink;relink, assert each exit and zero residue after cycle 2 — single-cycle idempotence already covered by `testLinkSafeFreshRelink`/`testUnlinkAbsentManifest` | a reconciler that left an orphan or a stale version stamp on the second cycle fails the "zero residue after cycle 2" assertion |
| 10 | the fault seam drives a deterministic mid-transaction failure and is inert unset | link contract | covered by story 1's rows (which drive the seam) plus every other link contract running with the variable unset | a fault branch that fired when unset would fail every green link contract; one that never fired would leave story 1's rows green-on-current |
| 5,7 | every link/unlink exit code is pinned: 0 convergence, 1 hard error, 2 usage, 3 partial | link + unlink contracts | 0/2 already covered (`linkOK`, usage tests); 1 by the fault row; 3 by the conflict/residual rows above | each code has at least one row that goes red if that outcome is misrouted |
| 1 | the packed wrapper's clean link/unlink round-trip stays green (exit 0, manifest gone) | `artifact/artifact_test.go` | **already covered** by `assertInstalledArtifactLifecycle` (link → … → unlink asserts manifest absent) — pins that FT84 does not regress the packed path | a transaction change that broke the packed clean round-trip fails the existing "unlink left link manifest" assertion |

Handoff black-box assertables (item 4) all land as rows: exit codes → the enumeration row; file+git state bit-identical after fault → story 1; manifest across upgrade/downgrade → stories 3–4; conflict and residual lists on stdout → stories 5–6.

### Edge inventory

Walked per behavior against the canonical classes and the profile's shell-CLI hostile-input checklist; each lands as a row above or a **Won't handle** line here.

- **Error path** — I/O fault mid-promotion → story 1 row.
- **Boundary** — fault at first/last promotion step → story 1 boundary row.
- **Empty/absent input** — absent/unreadable manifest on unlink stays exit 1 (already covered, `testUnlinkAbsentManifest`); absent `AGENTS.md`/`CLAUDE.md` on link (create-if-absent) already covered.
- **Malformed input** — a hand-edited old-manifest row that traverses outside the repo is refused by the reconciler via the shared `resolveInside` guard → covered by extending `testUnlinkRefusesTraversal`'s guard to the reconcile path.
- **Re-run idempotency** — repeated link/unlink → story 9 row.
- **Hostile environment** — spaces/glob metacharacters in the kit or repo root: the staging temp paths must tolerate them → covered by keeping `testLinkMetacharKitPath` green with an added assertion that staging under a `kit[1]` root still converges. Control bytes in the TOON lists: managed rels contain none, and `toon.Table` refuses control bytes by construction → safe.
- **Special files** — a FIFO/device at a managed destination: preflight already rejects non-regular kit assets (`link.go:195`); the reconciler rejects a non-regular dropped-asset target rather than removing it → covered by extending the reconcile classification.
- **Won't handle: SIGKILL / power loss mid-promotion** — the in-process journal rolls back returned I/O *errors*, not an untrappable process kill; recovery is a re-run of link, which the reconciler makes idempotent (partial files are managed-clean or conflicts). Crash-consistent cross-process journal replay is a separate capability (below).
- **Won't handle: semantic JSON merge of a modified `settings.json`** — a rejected alternative (map #5/#9: never merge); a modified `settings.json` is a preserved conflict, not a merge target.

## Out of scope

- **Crash-consistent recovery across process kills** — a persistent write-ahead journal that a *later* link replays to finish or undo an interrupted promotion after SIGKILL/power loss. A separate capability (durable on-disk WAL + recovery-on-next-link), not the rest of this feature, which handles returned I/O errors in-process. Estimate: ~6 edits, ~4 gate runs.
- **Relink plan preview** — surfacing the reconcile plan before applying it. Left to FT76, whose bootstrap preview may subsume it (map "Not yet specified"). Estimate: owned by FT76.
- **Shim and binary-cache lifecycle removal** — unlinking the doctor-generated `~/bin/bench` stable shim and pruning the binary cache. Owned by the shipped FT87 slice 2; this spec's unlink only refuses to recommend raw deletion of a non-owned executable. Estimate: owned by FT87.
