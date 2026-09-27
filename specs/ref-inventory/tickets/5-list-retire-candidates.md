# 5. List the superseded candidates and the unique count on spec retire

Blocked by: 2-route-status-to-the-plan.md, 4-discard-a-unique-ref-by-target.md
Writes: internal/spec/spec.go, internal/spec/spec_test.go
Covers: RI48, RI49, RI50, RI51, RI52

## What to build

Chunk: RI-C2b.

After every removal of `bench spec retire` succeeds, and before the `next:` line, print the candidate lines.
One line per active or cleanup-pending assignment whose label or request token contains the slug reads `superseded candidate: <assignment id> <label> — bench worktree clean --discard-branch --target <assignment id>`.
One count line reads `unique refs: <n> — bench worktree clean --discard-branch --unclaimed`, where `<n>` comes from the class counts of ticket 2.
The listing discards nothing and changes no exit code.
Verify the new import of the worktree package against the import cycle check before the commit.

## Acceptance

- [ ] Retire of slug `s` with one active assignment labelled `s-build` prints one candidate line with its id and the exact `--target` command.
- [ ] Retire with two unique unrecorded refs prints `unique refs: 2 — bench worktree clean --discard-branch --unclaimed`.
- [ ] Retire with no matching assignment prints `unique refs: 0` and no candidate line.
- [ ] Retire with a candidate exits 0 and every branch ref survives.
