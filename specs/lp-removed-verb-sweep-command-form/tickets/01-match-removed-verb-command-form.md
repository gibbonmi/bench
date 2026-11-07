# Match only the command form in the removed-verb sweep

Blocked by: none
Writes: internal/conformance/docs_workflow_checks_test.go
Covers: none

## What to build

The removed-verb sweep finds the tokens `spec build` and `worktree recovery` in kit guidance. It used a substring match, so ordinary prose such as "this spec builds" made the sweep red. One matcher now finds each token only with a word boundary on each side. The sweep still finds each token with or without a `bench ` prefix.

## Acceptance

- [ ] The sweep gives no diagnostic for the prose "this spec builds" or "the worktree recovers".
- [ ] The sweep gives a diagnostic for "bench spec build" and for "`worktree recovery`".
- [ ] All tokens go through one matcher.
