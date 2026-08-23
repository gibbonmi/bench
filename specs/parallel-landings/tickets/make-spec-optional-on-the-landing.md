# Make the spec optional on the landing

Blocked by: none
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/worktree/land_surface_test.go, internal/landing/landing.go, internal/landing/landing_test.go, internal/usage/worktree.go
Line: opus / medium — a known seam with exact rows, and the surface grants landing authority.

## What to build

`bench worktree land` and its `--resume` form accept a call with no `--spec`.
Without it, the source proof skips the spec resolve, the transition, and the
fence authorization; every other proof runs unchanged. The landing request
carries no spec path, the composition applies no transition, and the resume
path skips only the published-spec authentication. The `landed{...}` record,
the refusal record, and the exit codes do not change.

The usage lines show the
optional flag; the reference paragraph belongs to the guidance ticket. A
`--spec` that names a staged `spec.md` behaves exactly as today on the first
run. This ticket and `union-merge-the-phase-owned-journals.md` both write
`internal/landing/landing.go`, so they run serially, this one first.

## Acceptance

- [ ] A land without `--spec`, with a correct request, base, and tip, exits 0 with `worktree=released` (covers WL1).
- [ ] The spec-less published commit has the destination and the source tip as its two parents, and the green marker advances to it (covers WL2).
- [ ] A spec-less landing whose composed tree fails the gate refuses and publishes nothing (covers WL3).
- [ ] A spec-less source whose commits touch a path no fence names lands (covers WL4).
- [ ] A spec-less landing publishes every `specs/` file byte-identical to the composition (covers WL5).
- [ ] A spec-less land with a `--source-tip` that differs from the worktree HEAD refuses with both identities (covers WL6).
- [ ] A spec-less landing interrupted after publication resumes without `--spec` and exits 0 (covers WL7).
- [ ] A land with `--spec` that names a staged `spec.md` still publishes `Status: implemented` (covers WL9).
- [ ] A spec-less landing's `landed{...}` record carries the same fields and exit 0 (covers WL10).
- [ ] A land with `--spec` that names a staged `spec.md` still refuses an out-of-fence path (covers WL21).
- [ ] A spec-less refusal still exits 1 (covers WL22).
- [ ] A spec-less land whose `--base` is not an ancestor of the source tip refuses before the gate (covers WL23).
- [ ] A spec-less landing's `landed{source_base=...}` value is the resolved review base (covers WL24).
- [ ] A `--resume` without `--spec` on a published spec-backed landing completes the marker and the release without a second publication (covers WL25).
- [ ] `--spec ""` stays a usage error with exit 2 (edge under WL1).
