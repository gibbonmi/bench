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
