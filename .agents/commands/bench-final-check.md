---
description: Run the external gate (types, tests, lint, project conformance), report findings, and commit the verified work to the active branch on green. Use to check whether work is actually done, or to diagnose a red gate. Never use the model's own judgment as a substitute.
---

# /bench-final-check — the gate is the oracle

## Entry orientation

This is the final verification phase. It runs the external gate and reports the
oracle result; it does not substitute model judgment for tests, types, lint, or
project conformance.

## Exit handoff

Close by reporting the gate result plainly. This phase owns the landing commit
and the spec's `Status: implemented` transition — `/bench-implement-spec` ends
at its last green build commit and hands off here. On green, land the verified
work with `bench commit -m "<msg>"`, naming the files the work touched — it
gates and commits them atomically, and it enforces the commit discipline so
this phase doesn't restate it. Use `bench commit -m "<msg>" --spec <slug>` when
the commit finishes a spec, so the status flip rides in the same commit;
`implemented` honestly means *built, gate-green, awaiting review/merge* — the
state `bench status`'s retirement signal keys on once the spec reaches the
default branch. The honest no-op: a branch that arrives with nothing left to
commit is reported green all the same, and the status flip is still performed
via `bench spec implemented <slug>` when the spec has not already flipped —
that command is the single source of the edit; never hand-write a line-start
`Status: implemented` into any `specs/*.md`, or the retirement detector fires
on a spec that is not done. When the reviewer accepts residual review
risk and skips the fix pass, delete `reviews/<spec-slug>.md` and name it in the
landing commit — the pickup must not outlive the decision it captured. If it
refuses over an unexplained
working-tree file, handle the refusal as `/bench-implement-spec`'s close-on-green
describes. Then hand back for the reviewer to merge or decide what ships. On red, report
the first failing check and the smallest reproduction, then recommend the command
that fits the failure: usually `/bench-implement-spec` for a feature regression or
`/bench-debug` for a bug.

**The post-merge tail (exit duty).** After landing on green on the default
branch, read `bench status` and run the housekeeping rows it flags before
closing: a merged spec awaiting retirement gets `bench spec retire <slug>` and
its `spec-retire: <slug>` commit — promoting durable content first (a decision
to an ADR, a hostile edge to the profile) and deleting the shipped decision map
behind the spec in the same pass; an orphaned review pickup is promoted or
deleted by hand; leftover worktrees and scratch branches go through
`bench worktree clean`. Leave the roadmap and capture rows to
`/bench-what-next` — that phase owns the reconcile and the drain, and this duty
never restates it. On a topic branch these duties defer by design: the rows
fire only on the default branch, and the next default-branch session's
SessionStart status re-surfaces them — state the deferral in the close instead
of silently skipping it.

Run the gate and report. The gate is the oracle; this command does not form an
opinion about whether the work is good, it reports what the gate says.

## Run it

```sh
bench gate
```

`bench gate` runs the project's gate: an executable `.bench/gate.sh` when present,
else the `$BENCH_GATE` command string, else stack auto-detect (typecheck → test →
lint). `projects/<name>.md` documents what the gate covers — it never selects the
gate; to change what runs, change `.bench/gate.sh`.

## Report

- **Green:** state it plainly, land the verified work with `bench commit -m "<msg>"`
  per the exit handoff, and hand back for me to merge.
- **Red:** report each failing check in the order it fails, with the smallest
  reproduction. Do not propose weakening the check. Diagnose the cause, propose a
  fix at the seam, and — if I approve — fix it and re-run the gate. A fix is only
  real when the gate is green again.

If a check itself looks wrong (a flaky test, an over-tight lint rule), say so
explicitly and stop. Changing a gate check is my call, not a step inside `/bench-final-check`;
when I approve one, the `craft-gate` skill governs how it's made.

## Findings that the gate can't see

The gate catches regressions and conformance, not whether the design is right. If
while verifying you notice a real design problem the tests pass through — a leaky
seam, a story the spec missed, a shortcut that will cost later — name it
separately as a finding. Don't fold it silently into a fix.
