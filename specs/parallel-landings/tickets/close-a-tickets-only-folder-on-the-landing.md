# Close a tickets-only folder on the landing

Blocked by: make-spec-optional-on-the-landing.md
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/landing/landing.go, internal/landing/close.go, internal/landing/landing_test.go
Line: opus / medium — reuses the removal the commit path already owns, at the reviewed landing seam.

## What to build

A `--spec <slug>` that names a tickets-only folder (a direct child of `specs/`
with no `spec.md`) closes that folder on the reviewed landing. The source proof
treats that slug as a close, not as a staged spec: no transition and no fence
authorization. The composition removes the folder from the composed tree with
the removal the commit path already owns, and reconcile removes it from the
destination checkout. A folder the destination already removed composes as a
no-op. The refusal for a slug that is neither a staged `spec.md` nor a
tickets-only folder stays the existing unreadable-spec refusal.

## Acceptance

- [ ] A land with `--spec` that names a tickets-only folder publishes a tree without that folder and releases (covers WL8).
- [ ] A land with `--spec` that names a tickets-only folder the destination already removed lands (edge under WL8).
- [ ] A land with `--spec` that names a folder absent from the source refuses with the existing unreadable-spec detail (edge under WL8).
