# Review pickup — defect-batch-ft43-49 (FT45 landing commit 94de8e2)

Semantic three-axis review of `94de8e2` (stories 5–6: lease `Claim` takeover
identity check + two-reclaimer test). Findings needing a fix pass or a
reviewer decision.

## Standards

1 finding; worst: doc-comment knowledge duplication.

- **[Sev 2] `Claim` header restates the recovery mechanism the inline comment
  owns.** `internal/worktree/lifecycle.go:85-89` (header) vs `:108-112`
  (inline): both derive "bytes differ → renamed back, best-effort → concede".
  AGENTS.md one-source-per-fact — a future mechanism change must edit two
  prose sites or leave one stale. Fix: keep the header at the guarantee
  ("two concurrent reclaimers cannot both win"), let the inline own the how.

## Spec

1 finding; worst: story 6's test home deviates from the spec's stated seam.

- **[Sev 2 — reviewer veto call, not a defect] Story 6 test landed at an
  internal seam, not the contract suite.** Spec (story 6, testing decisions
  "no internals-poking", S3 diagram) names `internal/contract/runtime`'s
  concurrent-acquire suite as the home; the test lives in
  `internal/worktree/lifecycle_test.go` driving the `claimTakeoverGap`
  package var. The guarantee is genuinely locked and deterministic — the
  tradeoff is that deterministic forcing of the interleaving needs an
  in-process seam (the `restoreClean` precedent). Flagged at FT45 handoff;
  decision: accept the placement, or direct a (non-deterministic) contract-
  surface variant.

## Coverage

1 finding (verified by reproduction) + 1 minor note; worst: three-party
rename-back clobber.

- **[Sev 2 — CONFIRMED] Best-effort rename-back clobbers a fresh
  first-writer's lease → double-win.** `internal/worktree/lifecycle.go:115`:
  `os.Rename(stale, leasePath)` replaces the destination unconditionally.
  Interleaving: reclaimer A judges dead lease, pauses; reclaimer B completes
  takeover (fresh lease `B`, returns true); A renames `B` away; first-writer
  C `tryCreate`s the now-empty slot (returns true); A's mismatch rename-back
  overwrites C's lease with B's bytes. B and C both hold one worktree and the
  on-disk lease names neither loser. Reproduced deterministically in a
  scratch module (verbatim takeover logic + a scaffold seam for the
  rename→rename-back window; in production the window is one ReadFile+compare
  wide — narrow but real). The spec's edge inventory never lists a
  third-party first-writer in the steal window, so this is an undecided edge,
  not an accepted residual. Candidate fixes for the decision: (a) no-clobber
  restore (`os.Link(stale, leasePath)` + remove-stale; on EEXIST discard the
  stolen bytes) — kills the lease destruction, leaves a narrower B-vs-C
  double-use residual; (b) a reclaim-lock protocol serializing takeovers —
  full single-winner guarantee, but a new protocol decision (lock staleness,
  crash-safety) beyond the batch spec's flagged default.

- **Minor (spec prose):** the edge inventory claims a leftover
  `.stale.<pid>` is "treated as an ordinary dead lease on the next scan";
  actually it sits inside `.git` where the pool scan never reads it —
  harmless litter, never reclaimed. Amend the sentence when next editing the
  spec.
