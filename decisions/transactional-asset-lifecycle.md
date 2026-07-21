# Transactional managed-asset lifecycle (FT84)

## Destination

Link, relink, and unlink become transactional over the managed-asset set:
stage and preflight the complete write set, then atomically promote or roll
back; upgrade and downgrade reconcile old and new manifests; unlink is honest
about residuals. FT76's bootstrap composes these seams
([[repo-aware-bootstrap]] #1), so this ships first. Sources: `RR:A-05`,
`RR:A-07`, `RR:A-10`, `RR:A-13`; `RC:M-04`.

Recorded under the same-session exception in `/bench-write-spec`: every ticket
below restates a decision the reviewer already approved (the FT84 roadmap row
through reviewed drains, or the FT76 grill of 2026-07-21) — flagged in the
spec for veto rather than re-asked.

## #1: Relationship to FT76

Blocked by: none
Type: Grill

### Question

Who owns write-transaction semantics when the bootstrap also "transactionally
seeds"?

### Answer

FT84 owns them and builds first; FT76 is a composition layer over these seams
and adds no second transaction implementation. (Resolved 2026-07-21, FT76 map
#1.)

## #2: Failure posture for link/relink writes

Blocked by: none
Type: Grill

### Question

Today `Link` preflights, then writes sequentially (instructions file, CLAUDE.md,
plan, git hook); a mid-sequence failure leaves a partial footprint. What is the
required posture?

### Answer

Stage and preflight the complete write set, sync where durability matters,
then atomically promote or roll back. Settings and Git hooks compose or fail
*before* any partial write. A failed link leaves the repo as it found it.
(Roadmap row, reviewer-approved.)

## #3: Manifest reconciliation on upgrade and downgrade

Blocked by: none
Type: Grill

### Question

Relinking against a different kit version: what happens to assets the new
manifest no longer carries, and to assets the project has modified?

### Answer

Old and new manifests reconcile: removed clean assets leave; modified assets
remain owned until explicitly resolved; no stale skill becomes
active-but-unowned. (Roadmap row, reviewer-approved.)

## #4: Unlink and shim-removal posture

Blocked by: none
Type: Grill

### Question

What does unlink guarantee about what it removes and what it reports?

### Answer

Unlink and shim removal verify ownership markers, return a nonzero partial
status when residuals remain, and emit a machine-readable residual list. They
never recommend raw deletion of an executable Bench does not own. (Roadmap
row, reviewer-approved.)

## #5: Conflict posture shared with the bootstrap

Blocked by: none
Type: Grill

### Question

What happens to project-owned or modified content the write set would touch?

### Answer

Never merge, never overwrite. Non-conflicting assets are written; conflicting
assets are preserved as project-owned and reported with a distinct partial
status and a machine-readable conflict list. (Resolved 2026-07-21, FT76 map
#5.)

## #6: Closure matrix

Blocked by: none
Type: Grill

### Question

What scenario coverage closes this row?

### Answer

Upgrade, downgrade, I/O failure mid-transaction, modified and stale assets,
and repeated link/unlink matrices. (Roadmap row, reviewer-approved.)

## Handoff

1. **Module boundaries.** `internal/adopt` owns the whole lifecycle: plan
   build, preflight, staging, atomic promotion, rollback, manifest
   reconciliation, and unlink. The wrapper dispatch is untouched (routes
   exist). Contract coverage extends the existing surface contracts for link
   and unlink.
2. **Contracts.** `bench link [copy|symlink]`: exit 0 on full convergence;
   nonzero on preflight or I/O failure with the repo left exactly as found;
   distinct partial status plus a machine-readable conflict list when
   conflicting assets were preserved. `bench unlink [--dry-run]`: 0 on clean
   removal; nonzero partial status plus a machine-readable residual list
   otherwise; never a raw-deletion recommendation for a non-owned executable.
3. **Deep vs thin.** The transaction engine (stage → preflight → promote /
   roll back) and the manifest reconciler are the deep modules. The wrapper
   route and the per-asset writers are thin.
4. **Black-box assertables.** Exit codes; target-repo file and git state
   (bit-identical to pre-link state after an induced mid-transaction
   failure); manifest contents across upgrade/downgrade; the conflict and
   residual lists on stdout.
5. **Gate attachment.** The existing link/unlink surface contracts plus the
   packed-artifact contract's cold-run leg; the gate runs all contract
   packages as one phase. No seam is gate-invisible.
6. **Hostile-input owners.** Induced I/O failure mid-write — the transaction
   contract; modified managed assets and stale removed-asset entries — the
   reconcile contract; spaced/metachar paths and the ephemeral-npx-cache kit —
   the existing link/artifact contracts; repeated link/unlink cycles — the
   matrix contract.
7. **Uncertainty flags.** The atomicity mechanism (per-file journal vs
   staged-directory swap) and the exact partial-status exit code are spec
   proposals, not settled decisions. The fault-injection seam for the
   I/O-failure contract (env hook vs filesystem shim) is unsettled — the spec
   proposes, the reviewer vetoes.
8. **Rejected alternatives.** Absorbing this into FT76 (#1); best-effort
   continue-on-error writes (#2); auto-deleting modified assets on downgrade
   (#3); merging conflicting content (#5).
9. **Domain watch-outs.** An ephemeral npx-cache kit forces copy mode today —
   symlinks would dangle. `.claude/settings.json` is an adapter target a
   project may have modified; composition must preserve, never clobber. The
   supported platform matrix is darwin/linux only.

Dependency order: n/a — single spec; within it, transaction engine →
manifest reconcile → unlink residuals.

## Not yet specified

- Whether relink surfaces a plan preview (FT76's preview may subsume it).

## Out of scope

- The bootstrap's explore/present/confirm flow — [[repo-aware-bootstrap]].
- Repair/binary-cache lifecycle — owned by the shipped FT87 slice 2.
