# landing-refusal-diagnostics review pickup

Frozen base: `850fd2677c4fd56ff2a0e8f2f1c5ef1698ba1148`

Reviewed tip: `a302a01fc2362a97a0ee58fbce1bcc9f9fe4ae59`

## Standards

Finding count: 1. Worst issue: P2 malformed tracer-bullet ticket.

- **auto-fix — repair ticket omits the required ticket shape.**
  `specs/landing-refusal-diagnostics/tickets/repair-review-findings.md` lacks
  `Blocked by:`, `## What to build`, and `## Acceptance`. `craft-tickets`
  requires those fields for every ticket; add the dependency and observable
  acceptance rows while retaining the advisory write set.

## Spec

Finding count: 2. Worst issue: P1 mapped behavior is exercised at the wrong input.

- **auto-fix — LR18 does not diagnose the specified non-ancestor `--base`.**
  LR18 requires that input to name the assignment's recorded start
  (`specs/landing-refusal-diagnostics/spec.md:284`). `ReauthorizeCommand`
  returns the generic range error before it loads and proves the assignment
  (`internal/worktree/reauthorize.go:62`), while the current test mutates stored
  `Start` and supplies a valid base (`internal/worktree/reauthorize_test.go:80`).
  Drive the specified bad-base input red, then move the typed wanted-value
  refusal behind the otherwise-valid assignment proofs.

- **auto-fix — lost-token recovery is absent from `land --resume`.**
  Story 21 promises recovery for a fresh-session landing refusal when exactly
  one active assignment occupies the target path
  (`specs/landing-refusal-diagnostics/spec.md:138`). Resume uses
  `resumeAssignment`, whose no-match path returns no recovery hint
  (`internal/worktree/land.go:282`). Add the same typed reauthorize recovery
  without breaking the terminal-receipt path for already-complete resumes.

## Coverage

Finding count: 1. Worst issue: P2 hostile-input matrix gap.

- **auto-fix — safe target plus control-bearing request is untested.**
  `releaseNext` has distinct safe/unsafe target branches
  (`internal/worktree/ownership.go:434`). Tests cover a safe target with a
  line-safe token and an unsafe target with both token classes, but not a safe
  target with `release\n\x1brequest`. Add that cell and assert exit 1, stderr
  channel, `<request>`, absence of caller-token bytes, and no raw controls.
