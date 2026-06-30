---
description: Run the external gate (types, tests, lint, project conformance) and report findings. Use to check whether work is actually done, or to diagnose a red gate. Never use the model's own judgment as a substitute.
---

# /bench-qa — the gate is the oracle

Run the gate and report. The gate is the oracle; this command does not form an
opinion about whether the work is good, it reports what the gate says.

## Run it

```sh
bench gate
```

`bench gate` runs the project's gate command (from `projects/<name>.md`, falling
back to the kit default: typecheck → test → lint, plus any project conformance
check such as AXI conformance for gl-axi).

## Report

- **Green:** state it plainly. The work is eligible to be done. Hand back for me
  to merge.
- **Red:** report each failing check in the order it fails, with the smallest
  reproduction. Do not propose weakening the check. Diagnose the cause, propose a
  fix at the seam, and — if I approve — fix it and re-run the gate. A fix is only
  real when the gate is green again.

If a check itself looks wrong (a flaky test, an over-tight lint rule), say so
explicitly and stop. Changing a gate check is my call, not a step inside `/bench-qa`.

## Findings that the gate can't see

The gate catches regressions and conformance, not whether the design is right. If
while verifying you notice a real design problem the tests pass through — a leaky
seam, a story the spec missed, a shortcut that will cost later — name it
separately as a finding. Don't fold it silently into a fix.
