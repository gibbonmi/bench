# Learnings — usage journal

## 2026-08-19 — FT226 spec review took two iterations: fence missed the handoff file, residue unit unpinned [open]

**What happened.** The `/bench-write-spec` round for `ft226-test-home-isolation`
returned REVISE on two blocking findings. The ownership fence listed the test files
and the spec directory but not `capture/session-handoff.md`, which every phase close
rewrites, so `bench preflight build` was red before the build started. And the
`TestMain` residue predicate said "any entry under the home" while ticket 02's
acceptance counted top-level entries — two readings of the same report, and the
edge inventory's "empty `worktrees/` is residue" claim had no red under the
recursive one.

**Right behavior.** Run `bench preflight build <slug>` before the review round and
treat every `paths-authorized` red on a file the workflow itself writes (handoff,
retro) as a fence omission, not noise. And when a spec names a unit-tested
predicate, state its unit of report once in the implementation decision and derive
the ticket's planted cases from that sentence.

**Proposed rule change.** `craft-spec`'s fence rule: the fence lists every path the
workflow writes at phase close — name `capture/session-handoff.md` as the standing
example. Optionally `bench preflight build` could mark the handoff path as
workflow-owned so its dirt never reads as a fence breach.

## 2026-08-19 — `bench worktree path` prints a `~` path, so a composed `cd` failed and an edit landed in the main checkout [open]

**What happened.** Verifying a delegate's ticket, I ran
`cd "$(bench worktree path <label>)" && git show HEAD:<file> > <file> && <edit>`. The
verb printed `~/.bench/worktrees/...`; a quoted `~` is not expanded by bash, so the
`cd` failed — but the `&&` chain's later steps ran anyway, in the main checkout,
rewriting a skill file there. Caught by `git status`; restored from HEAD, no harm.

**Right behavior.** Never compose a `cd` from a printed path in one chain; use
`bench worktree exec <label> -- <cmd>` for anything inside the worktree, and run
`git status` on main after any worktree session. If a path must be used, `$HOME`-expand
it explicitly.

**Proposed rule change.** `bench worktree path` should print the absolute path it
resolved (`$HOME` expanded), the same form `worktree create` prints — an agent-facing
verb's output is meant to be pasted into a shell. If the `~` form is deliberate for
display, add `--absolute`; either way the two verbs should agree.
