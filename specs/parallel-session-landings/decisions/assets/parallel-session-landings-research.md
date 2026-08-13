# Parallel-session landing research

Local-code research performed 2026-08-13 against `main` at `b1390c9f`. This
asset answers decision ticket #1 in
`specs/parallel-session-landings/decisions/parallel-session-landings.md`.

## Public symptom

`bench preflight build roadmap-progressive-index` exits 1 with every check green
except `paths-authorized`. With the recorded `branch.main.benchBase` at
`5dda59db`, the refusal includes historical canary, assessment, lifecycle, and
capture paths outside FT198's ownership fences.

Testing candidate bases through invocation-local Git configuration, without
changing repository config, minimized the failure:

| base | build preflight | review preflight |
|---|---|---|
| `b6ef66e7` (staged FT198 spec) | ownership red only on `capture/session-handoff.md` | ownership red; diff non-empty |
| `0924e02e` (reviewed FT198 tickets) | ownership red only on `capture/session-handoff.md` | ownership red; diff non-empty |
| `c46b135a` (third landed FT198 ticket) | ownership red only on `capture/session-handoff.md` | ownership red; diff non-empty |
| `b1390c9f` (later roadmap/handoff drain) | green | ownership green; diff empty |

No scalar base simultaneously retains a non-empty complete FT198 review and
excludes the later phase-owned handoff change.

## Producer and consumers

- `internal/shift/loop.go` writes `branch.<name>.benchBase` from the branch tip
  recorded at shift entry.
- `internal/diff/diff.go` owns review-base resolution and prefers that recorded
  key over the default-branch merge base.
- `internal/preflight/gather.go` deliberately consumes the same resolved base
  and passes its complete changed-path set to ownership checking.
- `internal/preflight/decision.go` compares every changed path with the staged
  spec's ownership fences and makes an empty review diff red.
- `internal/handoff/handoff.go` rewrites `capture/session-handoff.md` as tracked
  phase metadata. `.agents/commands/bench-implement-spec.md` requires that write
  at full-run phase boundaries but does not create a separate review identity.

The components are internally consistent. The missing fact is the build's
explicit source identity: one branch-relative contiguous range cannot name a
non-contiguous authorship set after phase-owned commits interleave.

## Retired assumptions

The earlier version of this asset researched the provisional `bench spec-build`
lifecycle. That command family and `internal/specbuild` were removed wholesale
in `dae240df`; `CHANGELOG.md` records zero compatibility and serial
commit-on-green ticket landing as the replacement. Those deleted run-state,
receipt, promotion, and status surfaces are not current implementation evidence
and no longer support this map.
