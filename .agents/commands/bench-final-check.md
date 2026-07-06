---
description: Run the external gate (types, tests, lint, project conformance), report findings, and commit the verified work to the active branch on green. Use to check whether work is actually done, or to diagnose a red gate. Never use the model's own judgment as a substitute.
---

# /bench-final-check — the gate is the oracle

## Entry orientation

This is the final verification phase. It runs the external gate and reports the
oracle result; it does not substitute model judgment for tests, types, lint, or
project conformance.

## Exit handoff

Close by reporting the gate result plainly. On green, land the verified work with
`bench commit`, naming the files the work touched — it gates and commits them
atomically, and it enforces the commit discipline so this phase doesn't restate
it. Use `bench commit --spec <slug>` when the commit finishes a spec, so the
status flip rides in the same commit. If it refuses over an unexplained
working-tree file, handle the refusal as `/bench-implement-spec`'s close-on-green
describes. Then hand back for the reviewer to merge or decide what ships. On red, report
the first failing check and the smallest reproduction, then recommend the command
that fits the failure: usually `/bench-implement-spec` for a feature regression or
`/bench-debug` for a bug.

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

- **Green:** state it plainly, land the verified work with `bench commit` per the
  exit handoff, and hand back for me to merge.
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
