---
description: The bug path. Build a tight, red-capable repro loop first, then fix against it. Use whenever something is broken, throwing, failing, or slow — instead of /bench-write-spec, which is the feature path. Reach for this the moment the work is "fix" rather than "build".
---

# /bench-debug — the repro loop is the oracle

## Entry orientation

This is the bug path. Use it when something is broken, throwing, failing, or slow.
It produces a tight red-capable repro loop first, then uses that loop as the
external signal for diagnosis and the fix.

## Exit handoff

Close by reporting the repro command, the confirmed cause, the fix state, whether
the repro loop and regression check are green, and any architecture finding that
would have prevented the bug. The recommended next command is `/bench-final-check`
once the loop and project gate are ready; if the fix needs semantic review against
a written spec, run `/bench-review-implementation` first.

Bugs don't go through `/bench-write-spec`. A bug already has a spec — the thing should work and
doesn't. The whole discipline is building an external signal that goes red on *this*
bug, and that signal becomes the gate the fix shift runs against. Same invariant as
everything else in Bench: the agent never decides the bug is fixed; a check it
didn't author does.

## Phase 1 — build the loop (this is the whole skill)

Before reading code to form a theory, build **one command** — a test, a curl, a CLI
invocation, a replay of a captured trace, a throwaway harness — that you have
**already run once** (paste the invocation and its output) and that is:

- **red-capable** — drives the real bug path through the accused command — the
  surface the claim names, invoked as reported, never a lookalike (a raw
  `git add` is not `bench commit`) — and asserts the user's *exact* symptom,
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

The regression-test seam and the edit owner may differ. Before editing, enumerate
the relevant callers and trace the affected paths to the narrowest shared function
or module where the invariant is uniform. Fix it once there; keep the regression
test at the highest seam where the reported failure remains observable. A shared
helper whose callers require different behavior does not own that invariant. If
the paths have no honest shared owner, record the tangled ownership as the Phase 6
architecture finding instead of patching every caller.

## Phase 6 — close out

Done means: the Phase 1 loop no longer reproduces, the regression test passes (or its
absence is documented), all `[DEBUG-...]` logs are gone, throwaways deleted, and the
correct hypothesis is named in the commit so the next debugger learns. Then ask what
would have *prevented* this — if the answer is architectural (no seam, tangled
callers), say so as a separate finding, after the fix is in.

## Finding a retired spec

The feature you're debugging may have no live spec — specs are promote-then-deleted
on merge, so a shipped feature's spec is in git history, not the working tree. Don't
assume the behavior was never specified: `git log --diff-filter=D -- specs/` lists
every deleted spec, and `git log --grep=spec-retire` finds the retirement commits
(and any decision it promoted). Recover the origin spec there before hypothesising.
For a single known slug, `bench spec history <slug>` runs both queries, merges and
dedupes them, and renders one newest-first table — the sanctioned shortcut instead
of hand-running the two commands above.

## How it meets the rest of Bench

This phase is reviewer-invoked. A spec-build write delegate never charges it:
when its repro proves the defect lives outside its ticket fence, the delegate
stops implementation edits, keeps its in-fence work dirty in its owned worktree,
and returns a bounded blocked report — repro command with red exit and output
digest, the failing surface it observed, assignment ID, recorded base, in-fence
dirty paths. The reviewer then runs this skill against that surface, and the
run produces the debug receipt: the report's evidence plus confirmed cause, the
paths the repair must own, and whether the ticket can proceed after the repair.
The coordinator's route from that receipt — repair ticket, recomposition,
`assign --refresh` — is `.agents/commands/bench-implement-spec.md`'s "When a
delegate is blocked outside its fence"; this skill only produces the receipt's
evidence.

The Phase 1 loop joins the project gate for the fix. If the fix launches a shift,
add the repro as a test the gate runs alongside its existing checks, committed
in the project's expected-failure form — a quarantine marker naming the bug (a
skip or expected-fail annotation in the project's test framework) — so the
committed tree stays green while the repro survives shift rollback: an
iteration that ends red rolls the worktree back to the last commit, which
destroys any uncommitted repro test. The fix's green commit removes the marker,
turning the repro into the live regression test. A project with no
expected-failure form keeps the repro out of the shift and runs it by hand
against the fix — state that fallback in the close rather than committing a red
tree; invariant 4 (commit only on green) has no red-commit exception. Route code
authorship through `craft-delegate`, including a diagnosed single-seam fix; it
owns the inline threshold, worktree isolation, and verification discipline.
Replacing the gate with the repro weakens the oracle — that is my call, never a
debug step. The seam decision in Phase 5 is the `craft-seams` skill. Declare the
line first — a hard bug is a high-effort shift. Any delegation along the way (a
fan-out search, a scoped fix) carries its own line: an explicit bound model alias
on the Agent call — never the inherited default — with effort and iteration cap
stated in the charge, per `craft-line` and `craft-delegate`.
