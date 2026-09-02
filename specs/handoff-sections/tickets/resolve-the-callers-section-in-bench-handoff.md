# Resolve the caller's section in bench handoff

Blocked by: add-the-handoff-document-leaf-package.md
Writes: internal/handoff/handoff.go, internal/handoff/facts.go, internal/handoff/render.go, internal/handoff/sections.go, internal/handoff/render_test.go, internal/handoff/facts_test.go, internal/handoff/sections_test.go, internal/intent/assignment.go, internal/intent/assignment_lookup_test.go (new), internal/preflight/gather.go
Covers: HS2, HS3, HS6, HS7, HS8, HS9, HS10

## What to build

Verify the premise first: `bench handoff` reads no request, `assignmentTarget`
in internal/preflight/gather.go is the only worktree-to-assignment lookup, and
`liveSpecs` roots at the primary checkout. Then promote that lookup to an
exported `intent.AssignmentForWorktree`, and make preflight call it. In
`Command`, own `main` when `git.IsPrimaryCheckout` is true. Otherwise own the
matching active assignment's section. When neither holds, refuse with exit 1
and no write. Rewrite only the owned section through the leaf package, and
re-emit every other section byte for byte.

Render the six pins from the assignment. The request digest is the key, and
the plain token, the label, and the recorded base come from the record. The
tip comes from `git rev-parse HEAD` in the assignment's worktree. One spec path
and status pair per live spec comes from `spec.Facts` over that worktree. The
header carries the repository, the path, the `main` HEAD, and the gate verdict.

## Acceptance

- [ ] From assignment A's worktree, the verb rewrites A's section and leaves B's section byte-identical.
- [ ] From the primary checkout, the verb writes `main`.
- [ ] From a path with no active assignment that is not primary, the verb exits 1 and the file is unchanged.
- [ ] A section renders `Request token:` with the plain token byte for byte beside the label, the recorded base, and the spec pairs.
- [ ] A section's tip equals the worktree HEAD after a commit inside the worktree, and a worktree with two staged specs renders two pairs.
- [ ] The header pins the `main` HEAD, not the worktree HEAD.
- [ ] Self-probe: match the section by label instead of the digest, and report the stale-section test red.
