---
description: The bug path. Build a tight, red-capable repro loop first, then fix against it. Use whenever something is broken, throwing, failing, or slow — instead of /spec, which is the feature path. Reach for this the moment the work is "fix" rather than "build".
---

# /diagnose — the repro loop is the oracle

Bugs don't go through `/spec`. A bug already has a spec — the thing should work and
doesn't. The whole discipline is building an external signal that goes red on *this*
bug, and that signal becomes the gate the fix shift runs against. Same invariant as
everything else in Bench: the agent never decides the bug is fixed; a check it
didn't author does.

## Phase 1 — build the loop (this is the whole skill)

Before reading code to form a theory, build **one command** — a test, a curl, a CLI
invocation, a replay of a captured trace, a throwaway harness — that you have
**already run once** (paste the invocation and its output) and that is:

- **red-capable** — drives the real bug path and asserts the user's *exact* symptom,
  so it goes red now and green when fixed. Not "runs without erroring."
- **deterministic** — same verdict every run. Flaky bug? Raise the reproduction rate
  (loop the trigger, parallelize, inject stress) until it's debuggable; 50% is
  workable, 1% is not.
- **fast** — seconds, not minutes. A 2-second loop is a debugging superpower.
- **agent-runnable** — unattended, no human click.

If you catch yourself theorizing before this command exists, stop — jumping to a
hypothesis without a red loop is the exact failure this prevents. If you genuinely
can't build a loop, say so, list what you tried, and ask for an environment, a
captured artifact (HAR, log dump, trace), or permission to instrument. Don't
proceed on a theory.

## Phase 2 — reproduce and minimise

Run it red. Confirm it's the *user's* symptom, not a nearby one. Then shrink the
repro one element at a time, re-running after each cut, until every remaining piece
is load-bearing. A minimal repro shrinks the suspect space and becomes the
regression test.

## Phase 3 — hypothesise

Write **3–5 ranked, falsifiable** hypotheses before testing any: "if X is the cause,
changing Y makes it disappear." No prediction → it's a vibe; sharpen or drop it.
Show the ranked list before testing — I often re-rank instantly ("we just shipped
#3"). Don't block if I'm away; proceed on your ranking.

## Phase 4 — instrument

One variable at a time. Prefer a debugger/REPL over logs; tag every debug log with a
unique prefix (`[DEBUG-a4f2]`) so cleanup is one grep. For perf regressions, measure
first (baseline + bisect), don't log.

## Phase 5 — fix at a correct seam

Write the regression test **before** the fix — but only if a correct seam exists
(one where the test exercises the real bug pattern at the call site). **If no correct
seam exists, that is itself the finding** — note it; the architecture is preventing
the bug from being locked down. If a seam exists: repro → failing test → fix → green
→ re-run the full Phase 1 loop.

## Phase 6 — close out

Done means: the Phase 1 loop no longer reproduces, the regression test passes (or its
absence is documented), all `[DEBUG-...]` logs are gone, throwaways deleted, and the
correct hypothesis is named in the commit so the next debugger learns. Then ask what
would have *prevented* this — if the answer is architectural (no seam, tangled
callers), say so as a separate finding, after the fix is in.

## How it meets the rest of Bench

The Phase 1 loop is the project gate for the fix shift: add the repro as a test the
gate runs (or point `.bench/gate` at the repro command), then run
`bench shift "fix <bug>"` so the loop physically gates "done" on the bug going
green. The seam decision in Phase 5 is the `seams` skill. Declare the line first — a
hard bug is a high-effort shift.
