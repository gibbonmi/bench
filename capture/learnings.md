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

## 2026-08-19 — `bench worktree land` refuses unless source and destination carry identical staged spec bytes [open]

**What happened.** The FT226 build wrote its verification log into the spec inside
the integration worktree, which the spec's own ownership fence authorizes. At
landing, `bench worktree land` refused with `source and destination do not carry
identical staged spec bytes`. The fix was to copy the spec directory to the
destination and commit it there first, spending a whole extra gate run, after
which the land succeeded. Two earlier refusals in the same sequence cost another
two attempts: `--base`/`--source-tip` must be full 40-character SHAs (an
abbreviated pair returns the unrelated-sounding `worktree source tip mismatch`),
and `bench worktree reauthorize` requires `--base` to be an ancestor of the source
tip, so a destination that moved during the build is reauthorized against the
frozen review base, not the new destination HEAD.

**Right behavior.** Treat the staged spec as destination-owned during a
worktree-backed build: write the verification log to the destination and commit it
there before calling `land`, or accept the extra gate run. Always pass full SHAs to
`land`. When the destination moves mid-build, reauthorize against the frozen review
base and let `land` compose.

**Proposed rule change.** Two candidates, both reviewer's call. First, `land`'s
spec-bytes refusal should name the remedy it wants — it knows both paths and could
print the exact destination commit to make first. Second, `/bench-implement-spec`
could say where the verification log is written, since the spec's ownership fence
currently implies the source is fine and the landing verb disagrees. The
abbreviated-SHA case is a smaller separate fix: `land` should either accept what
`git rev-parse` accepts or say that it wants full SHAs, rather than reporting a
mismatch.

## 2026-08-19 — an ownership fence that lists a file but not its paired golden reds preflight for a correct edit [open]

**What happened.** Twice in one session. FT234 ticket 01 fenced `cmd/bench/main.go`
but not `cmd/bench/main_test.go`, which holds the hand-written help-inventory
golden. Advertising the new `reclaim` verb in the `worktree --help` row is
correct and required, and it necessarily changes the golden — so
`bench preflight build` went red with `not authorized by any ownership fence:
cmd/bench/main_test.go` on a diff that was right. The FT226 entry above is the
same class with a different paired path (`capture/session-handoff.md`, which
every phase close writes).

**Right behavior.** When a spec fences a file that has a paired expectation the
tree maintains by hand — a golden, an inventory fixture, an anchor registry — fence
the pair, not the file. The write delegate should also report a red preflight
rather than describing its diff as clean; this one did not, and the coordinator
caught it only by running preflight independently.

**Proposed rule change.** `craft-spec`'s fence rule gains one sentence: a fenced
file's hand-maintained paired expectation is fenced with it, with the help-inventory
golden and the session handoff as the two standing examples. Guidance, not
enforcement — `bench preflight` already catches the violation cheaply and
precisely, so the cost of the current state is one red preflight and a two-line
spec edit, which does not justify teaching preflight to infer pairings.
