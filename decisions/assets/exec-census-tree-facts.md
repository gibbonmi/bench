# Tree facts for the exec census shaping

Read on 2026-08-25 at `main` HEAD `969461e8`. A read-only delegate gathered
the citations; the main session verified findings 1 and 2 against the tree.

1. **The exec verb already passes stdin to the child.** `cmd/bench/main.go:69`
   builds the command with `Stdin: os.Stdin`, and `internal/worktree/exec.go`
   sets `cmd.Stdin` from that reader. A heredoc on stdin reaches the child
   today.
2. **No exec refusal prints the worktree path.** `printTargetRefusal`
   (`internal/worktree/path.go:115`) prints the verb and a detail keyed by
   the assignment id. A child start failure prints only the error.
3. **Assignment path shape.** The pool directory is
   `$BENCH_HOME/worktrees/<repo-base>-<crc32-of-root>/`
   (`internal/worktree/worktree.go:47`). Each assignment lives at
   `<pool>/<ownerID>-<assignmentID>` (`internal/worktree/ownership.go:185`).
4. **No Bench verb runs `git cherry-pick`.** The landing composes with
   `git merge-tree --write-tree` (`internal/landing/composition.go:250`) and
   refuses content it cannot merge. No verb leaves a conflicted index; an
   agent's own raw `git cherry-pick` does.
5. **Hook events.** Claude Code wires `WorktreeCreate`, `WorktreeRemove`,
   `SessionStart`, `Stop`, and `PreToolUse` (`Bash`, `Agent`). Codex wires
   `SessionStart`, `Stop`, and `PreToolUse` (`Bash`). Nothing in the kit uses
   `PostToolUse`.
6. **Hook state lives in the Go core under `$BENCH_HOME`.** Every hook script
   is a shim that pipes its envelope to a plumbing verb. The repository key
   is the pool directory name from finding 3.
7. **Landing record keys.** `landed{source_base,source_tip,destination_base,
   published_commit,tree,worktree[,next]}` (`internal/worktree/land_refusal.go:15`).
8. **One shared tokenizer.** `internal/benchguard` and `internal/gitguard`
   both call `shellcommand.Parse` on the envelope's `tool_input.command`.
   A census reads the verb head from that same parse.
