# Repair test-support fence

Blocked by: repair-post-review-evidence-closure.md
Writes: `internal/axitest/help.go`, `internal/axi/axitest/help.go`, `internal/axi/action_test.go`, `internal/coverage/coverage_test.go`, `internal/worktree/list_actions_test.go`, `capture/session-handoff.md`

## What to build

Move the shared TOON decode and POSIX argv test helper beneath the spec-authorized `internal/axi` subtree without changing its interface or any product behavior. Remove the out-of-fence package completely, update its three test consumers, and regenerate the review-ready handoff from main.

## Acceptance

- [ ] [TF1] (covers local) `internal/axitest` no longer exists in the candidate; the identical helper is owned at `internal/axi/axitest`, and the action-owner, coverage, and worktree tests import only that authorized package.
- [ ] [TF2] (covers QD1) the shared helper still decodes one rendered help command through official `toon-go`, executes it with a POSIX shell, and returns the complete argv vector; all three consumers retain independent expected bytes.
- [ ] [TF3] (covers local) `bench preflight review axi-query-disclosure` reports all five checks green, including `paths-authorized`, before a new candidate snapshot is captured.
- [ ] [TF4] (covers local) `capture/session-handoff.md` is regenerated from main with `$bench-review-implementation specs/axi-query-disclosure/spec.md --full` as the exact next command.
