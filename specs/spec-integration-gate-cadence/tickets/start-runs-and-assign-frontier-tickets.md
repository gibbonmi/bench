# Start runs and assign frontier tickets

Blocked by: none

## What to build

Create the deep `internal/specbuild` module with durable, versioned state in the
Git common directory. Its first vertical slice owns `Start`, `Assign`, and
`Status`: exact-base bootstrap, spec and ticket resolution, request-idempotent
owned worktrees rooted at the candidate, and one compact next-action projection.
The module depends on a small gate-owner interface so this ticket can land
independently of the concrete gate implementation.

## Acceptance

- [ ] [R01-R02] Start creates one run and its candidate/green identities only from validated exact green evidence, otherwise names `bench gate` then retry without mutation.
- [ ] [R04] Compatible start re-entry resumes; conflicting branch, tip, dirt, detached state, or run identity refuses without duplication.
- [ ] [R06-R09] Assign binds a real ticket's rows, fence, assumptions, candidate base, and request identity to one owned worktree; siblings may share a base and hostile ticket forms have no side effects.
- [ ] [R24] Status derives a definitive compact durable row, including the terminal/empty posture, without exposing private state.
- [ ] [R52] Slugs and tickets containing spaces or glob characters resolve literally while refs use opaque identities.

