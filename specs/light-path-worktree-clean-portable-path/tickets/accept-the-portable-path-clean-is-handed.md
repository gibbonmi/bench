# Accept the portable path every other worktree command emits

Blocked by: none

Ownership fence: `internal/worktree`
Assumptions: `targetPath` already owns the `~` grammar for `path` and `exec`; `CleanCommand` is the only entry point that takes an operator-supplied operand

## What to build

`bench worktree path` prints a portable, `~`-prefixed path by design, and `targetPath`
expands that form back for `bench worktree path` and `bench worktree exec`. `bench
worktree clean` is the one command that never learned the convention: it hands the raw
operand to `canonicalPath`, which resolves `~/...` as a *relative* path against the repo
root and produces `/home/devuser/workspace/bench/~/.bench/worktrees/...` — a target that
cannot exist.

So the two commands an operator composes to recover a worktree do not compose. Reading
the path out of `bench worktree path` and feeding it to `bench worktree clean` is the
obvious move, it is what the help rows steer toward, and it silently could not work. A
sibling fix made that refusal loud rather than silent, which is right but still leaves
the operator to hand-expand a path the tool just printed.

Expand the home form once, where the operand enters the command, so the plan and the
apply that follows it agree on one canonical target — a fingerprint computed against the
unexpanded path and applied against the expanded one would be a worse bug than the one
being fixed. The expansion rule itself stays single-sourced in `targetPath`'s grammar
rather than pasted a second time; `clean` composes it. The automatic sweep addresses
registrations it already holds and never sees an operator operand, so it does not change.

`~user/...` stays unsupported, as it is for `path` and `exec` today, but it earns its own
honest reason instead of being canonicalized into a nonexistent repo-root path and
reported as merely unregistered.

## Acceptance

- [ ] [CP1] A `~`-prefixed path taken verbatim from `bench worktree path` plans a removal when handed to `bench worktree clean`.
- [ ] [CP2] Plan and apply resolve one identical target for the home form, so a fingerprint taken from the plan applies without restating the path.
- [ ] [CP3] `~user/...` is refused as an unsupported home target rather than reported as unregistered.
- [ ] [CP4] Absolute operands, relative operands, and the automatic sweep behave exactly as they do today.
