# Resume an incomplete published landing without republishing

Blocked by: land-reviewed-sources-atomically.md
Writes: internal/landing, internal/worktree, internal/intent, internal/gate/authorization, internal/usage, cmd/bench/main.go, bin/bench.sh, internal/systemtest

## What to build

Add `bench worktree land --resume <published-commit>` as the only continuation
for a first run that won destination publication but did not finish marker,
checkout reconciliation, or release. Authenticate the published commit, its
source parent and staged spec, the destination relationship, request, path, and
terminal receipt. Resume performs no composition, gate, commit, or destination
CAS.

Later destination commits are allowed only when the published commit is an
ancestor; reconciliation targets the current destination tree, and a project-
green marker at the published commit or a descendant is already satisfied.
Repeat after completed cleanup returns authenticated `already-complete`; an
evicted bounded receipt refuses as `missing-terminal-receipt` without recreating
state or entering first-run behavior.

## Acceptance

- [ ] Every injected post-CAS incomplete step is resumable without a second
      publication, rollback, or gate run, preserving the first run's authoritative
      `landed` result (covers PL21, PL25).
- [ ] Concurrent marker movement resumes from its captured state without an
      unconditional write (covers PL28).
- [ ] Later destination movement reconciles forward and never regresses a marker;
      behind, absent, or divergent marker state refuses unchanged (covers PL30).
- [ ] Repeated completed resume exits 0 as `already-complete`; FIFO receipt
      eviction exits 1 as `missing-terminal-receipt`, with neither path requiring
      the removed assignment (covers PL25).
- [ ] A `--resume` value cannot masquerade as the required positional path, with
      `--` separator controls (covers the resume half of PL27).
