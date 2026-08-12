# Repair disclosure value proof seam

Blocked by: none
Writes: `internal/axi/action.go`, `internal/axi/action_test.go`, `internal/axitest`, `internal/coverage/coverage_test.go`, `internal/worktree/list_actions_test.go`

## What to build

Make honest-empty fallback cover every disclosure-derived known argument that the action renderer refuses, and collapse the repeated official-TOON decode plus POSIX-shell argv recovery harness behind one test-support seam. Keep each public surface's expected primary response, exit, action shape, and source value independently authored.

## Acceptance

- [ ] [DV1] (covers QD1) an empty known argument or one containing `<` or `>` never replaces a computed response or changes its exit; `RenderHelp` returns exactly `help[0]{cmd,why}:\n` instead of an error. Existing unsupported-control fallback remains exact.
- [ ] [DV2] (covers QD1) coverage and worktree public-command cases use legal paths containing `<` or `>`, preserve their checked primary bytes and exit 0, and append honest empty help. Mutating the fallback predicate makes each public case red.
- [ ] [DV3] (covers QD1) one test-support interface decodes the emitted help command with official `toon-go`, executes it with a POSIX shell, and recovers argv without embedding a surface-specific positional index. Action-owner, coverage, and worktree tests call that seam but retain independent expected argv bytes.
- [ ] [DV4] (covers local) malformed action definitions outside the exact empty, `<`/`>`, and unsupported-control disclosure-value set still return errors, so fallback cannot hide dropped arguments, guessed placeholders, prose commands, or noncanonical phases.
