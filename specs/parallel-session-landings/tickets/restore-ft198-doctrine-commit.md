# Restore the retained FT198 doctrine commit byte for byte

Blocked by: reauthorize-retained-assignments.md, land-reviewed-sources-atomically.md, resume-published-landings.md, route-workflow-through-integration-sources.md
Writes: .agents/commands/bench-what-next.md, CONTEXT.md, internal/anchors/registry_data.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/recurrence_maintenance_contract_test.go, internal/roadmap

## What to build

Recover retained FT198 assignment `f70c5f8dc98f203fac19bdd6e07df1d3`
without reconstructing its work. Record byte hashes for its nine existing dirty
paths, then run the reauthorization command from a clean destination checkout
using the current parallel-session-landings integration source's own
`./dist/bench` after the generic implementation and tests are green. Build that
binary first with
`bash scripts/go-build.sh <source-worktree-root> <source-worktree-root>/dist/bench`
and never substitute a PATH `bench`. Invoke its absolute path from the clean
destination checkout for reauthorization. Invoke the same absolute path with the
retained FT198 source worktree as cwd for
`bench preflight build roadmap-progressive-index --base 0924e02e`. Replace the
lost request token at source tip `c46b135a`. Return the new opaque token to the
coordinator as ephemeral hand-back data for the dependent ticket; never commit
it or write it into a handoff file.

The retained handoff requires one bounded repair before commit: collapse the
duplicate default body-omission enforcement between
`internal/roadmap/context_render.go` and `internal/roadmap/context_types.go` to
one owner. Only those two paths may change, and only for that collapse; the other
seven retained path hashes stay byte-identical. Do not regenerate, format, copy,
or discard the retained work, and do not touch `main`, either foreign assignment,
or ambient phase-handoff state. This ticket restores one independently green
FT198 doctrine commit; the remaining AXI ticket and public landing are owned by
its dependent ticket.

## Acceptance

- [ ] The retained assignment's path, branch, start, tip, marker, and old request
      digest are verified before replacement, and the new request authenticates
      the same worktree; the exact new token is returned ephemerally for the
      dependent ticket and is absent from repository and handoff files (covers
      the real-assignment half of PL29).
- [ ] The versioned feature binary is built by `scripts/go-build.sh` from the
      parallel-session-landings source and invoked by absolute path from the clean
      destination checkout for reauthorization.
- [ ] Default body omission has one enforcement owner after the named two-file
      repair, with the mutation proof still red when that owner is removed; the
      other seven retained path hashes are byte-identical before reauthorization
      and at the doctrine commit.
- [ ] The committed source remains rooted at frozen review base `0924e02e` with
      destination-only history excluded, proved by invoking the feature binary's
      absolute path for `bench preflight build roadmap-progressive-index --base
      0924e02e` with the retained FT198 source worktree as cwd, never PATH `bench`
      (covers PL2).
