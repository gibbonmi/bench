# Ownership-safe worktree cleanup review pickup

## Standards

3 findings. Worst issue: the safety-critical ownership ref grammars have
multiple production sources.

1. **High — centralize assignment and recovery ref grammars.** Assignment refs
   are independently constructed or validated in
   `internal/worktree/ownership.go:166`, `internal/intent/assignment.go:102`,
   `internal/intent/assignment.go:256`, and `internal/worktree/resume.go:317`.
   Recovery prefixes are independently built in
   `internal/intent/assignment.go:279`, `internal/worktree/classifier.go:197`,
   `internal/worktree/clean.go:301`, and `internal/worktree/subshell.go:333`.
   Collapse each grammar to one constructor. Rule: `AGENTS.md:25-29`, one
   source per fact.
2. **Medium — collapse the worktree command grammar.** The unused
   `internal/worktree/worktree.go:316` usage string omits `release`;
   `bin/bench.sh:269-273` advertises `release` but omits the supported
   `--discard-ignored` and `--full` clean flags; and
   `internal/worktree/worktree.go:341` repeats create syntax. Make one source
   feed help and usage errors. Rule: `AGENTS.md:25-29`, enforcement and its
   advertisement must not derive the same fact separately.
3. **Low — derive the ignored-entry overflow label.**
   `internal/worktree/ownership.go:24` defines the 1,000-entry limit while
   `internal/worktree/classifier.go:126` hardcodes `at-least=1001`. Derive the
   label from the limit. Rule: `AGENTS.md:25-29`, one source per fact.

## Spec

3 findings. Worst issue: the mandatory live Claude lifecycle acceptance pass
has no recorded verdict.

1. **High — complete and record the live Claude dogfood.** The spec requires a
   fresh Claude Code `--worktree` run and a recorded successful and
   injected-failure lifecycle at
   `specs/ownership-safe-worktree-cleanup.md:364-367` and
   `specs/ownership-safe-worktree-cleanup.md:439`. The landing commit contains
   only simulated hook tests in
   `internal/contract/runtime/runtime_shift_adapters_test.go:100-220`; no live
   verdict is recorded.
2. **Medium — classify recovery syntax before repository lookup.** Invocation
   errors must return 2 under `specs/ownership-safe-worktree-cleanup.md:248-254`
   and `specs/ownership-safe-worktree-cleanup.md:302-304`. The dispatch at
   `cmd/bench/main.go:324-330` resolves the repository first, so no-argument
   `worktree recovery` outside a repository returns 1 rather than usage exit 2.
   The existing usage test at
   `internal/contract/runtime/runtime_worktree_test.go:703-709` covers only an
   in-repository invocation.
3. **Medium — make the concurrency row race-gated.** The mapped row at
   `specs/ownership-safe-worktree-cleanup.md:434` requires synchronized
   apply/replay under `go test -race`. The concurrency test exists at
   `internal/worktree/resume_test.go:83-139`, but the gate invocations in
   `internal/gate/phases.go:69-71` and
   `internal/conformance/package_core_checks_test.go:148-154` omit `-race`.
   Wire the mapped race run into the oracle.

## Coverage

2 findings. Worst issue: an oversized malformed create event can mutate
repository state.

1. **Medium — reject oversized hook input before mutation.**
   `internal/harness/worktree.go:87-95` reads through a 1 MiB `LimitReader`
   without checking for overflow. A stream whose first 1 MiB is valid JSON plus
   whitespace and whose remaining bytes are invalid is accepted and creates a
   locked worktree. Add an oversized-envelope case at the simulated-hook seam
   from `specs/ownership-safe-worktree-cleanup.md:436`, asserting nonzero exit
   and no registration, assignment, branch, or marker. Existing hook coverage
   at `internal/contract/runtime/runtime_shift_adapters_test.go:100-139` uses
   only ordinary-size events.
2. **Medium — cover identical replay after successful WorktreeRemove.**
   `internal/harness/worktree.go:69-84` queries Git through the already-removed
   path before the shared release command can read its completed receipt, so an
   identical post-success remove event returns 1. This violates the identical
   release idempotency promised at
   `specs/ownership-safe-worktree-cleanup.md:474-476`. Existing retry coverage
   at `internal/contract/runtime/runtime_shift_adapters_test.go:156-220` retries
   only while the checkout survives. Add a post-success replay requiring exit 0
   and unchanged receipts and refs.
