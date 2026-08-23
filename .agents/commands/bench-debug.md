---
description: The bug path. Build a tight, red-capable repro loop first, then fix against it. Use whenever something is broken, throwing, failing, or slow — instead of /bench-write-spec, which is the feature path. Reach for this the moment the work is "fix" rather than "build".
---

# /bench-debug — the repro loop is the oracle

## Entry orientation

This is the bug path. Use it when something is broken, throwing, failing, or slow.
It builds a tight red-capable repro loop first. It then uses that loop as the
external signal for diagnosis and for the fix.

The debug work may write something — a throwaway harness, an instrumented
copy, an edit under test. In that case, create or select its isolated worktree
*before* the first repro artifact exists. If a clean checkout picks up writes
during the debug work, it becomes dirty in an unattributable way. Do
not decide isolation after Phase 1 already produced artifacts; that decision
comes too late.

## Exit handoff

Close by reporting the repro command, the confirmed cause, and the fix state. Report
whether the repro loop and the regression check are green, and report any architecture
finding that would have prevented the bug. Once the loop and the project gate are
ready, the recommended next command is `/bench-final-check`. If the fix needs semantic
review against a written spec, run `/bench-review-implementation` first.

Bugs do not go through `/bench-write-spec`. A bug already has a spec — the thing should work and
does not. The whole discipline builds an external signal that goes red on *this*
bug; that signal becomes the gate the fix shift runs against. The same invariant governs
everything else in Bench: the agent never decides the bug is fixed; a check it
did not author does.

## Phase 1 — build the loop (this is the whole skill)

Before you read code to form a theory, build **one command** that you have
**already run once** (paste the invocation and its output).
`.agents/skills/bench-debug/references/loop-constructions.md` lists ten ways to
construct one, tried in roughly that order, plus the structured
human-in-the-loop last resort. The right loop does most of the fix.

### Tighten the loop

Treat the loop as a product. Make it faster: cache the setup, skip unrelated
init, narrow the scope. Make it sharper: assert the specific symptom, not "did not
crash." Make it more deterministic: pin the time, seed the RNG, isolate the filesystem.

Phase 1 is complete when you check every box. The command is:

- [ ] **red-capable** — drives the real bug path through the accused command: the exact surface
  the claim names, invoked as reported. It is never a lookalike (a raw `git add` is not
  `bench commit`). It asserts the user's *exact* symptom, so it goes red now and green when fixed.
  Not "runs without erroring."
- [ ] **deterministic** — same verdict every run. For a flaky bug, raise the reproduction rate
  (loop the trigger, run it in parallel, inject stress) until you can debug it. 50% is
  workable; 1% is not.
- [ ] **fast** — seconds, not minutes. A 2-second loop is a debugging superpower.
- [ ] **agent-runnable** — runs unattended, with no human click.

If you catch yourself theorizing before this command exists, stop. A jump to a hypothesis
without a red loop is the exact failure this rule prevents. If you genuinely cannot build a
loop, say so and list what you tried. Ask for an environment, a captured artifact (HAR, log
dump, trace), or permission to instrument.
No red-capable command, no Phase 2.

## Phase 2 — reproduce and minimise

Run it red. Confirm:

- [ ] the loop produces the failure mode the *user* described — not a nearby
  one. The wrong bug means the wrong fix.
- [ ] the failure reproduces across runs — or, for a flaky bug, at a rate high
  enough to debug against.
- [ ] you captured the exact symptom (error message, wrong output, slow timing),
  so later phases can verify the fix addresses it.

A green proxy only narrows a hypothesis — it never confirms one. A load- or
environment-sensitive failure is reproduced through the accused command. This happens under the
conditions that expose it, before any other stand-in is trusted.

Then shrink the repro one element at a time. Re-run it after each cut, until
every remaining piece is load-bearing. A minimal repro shrinks the suspect
space and becomes the regression test.
Do not proceed until you have reproduced and minimised.

## Phase 3 — hypothesise

Write **3–5 ranked, falsifiable** hypotheses before you test any: "if X is the cause,
changing Y makes it disappear." A hypothesis with no prediction is a vibe; sharpen it or drop it.
Show the ranked list before you test — I often re-rank it instantly ("we just shipped
#3"). If I am away, do not block; proceed on your ranking.

## Phase 4 — instrument

Change one variable at a time. Prefer a debugger or a REPL over logs. Tag every debug log with a
unique prefix (`[DEBUG-a4f2]`) so cleanup is one grep. For performance regressions, measure
first (set a baseline, then bisect); do not log.

## Phase 5 — fix at a correct seam

If a correct seam exists, write the regression test **before** the fix. A correct seam
is one where the test exercises the real bug pattern at the call site. **If no correct seam
exists, that absence is itself the finding** — note it; the architecture leaves the bug
without a lock. If a seam exists, follow this order: repro, failing test, fix, green, then
re-run the full Phase 1 loop.

The regression-test seam and the edit owner may differ. Before you edit, list
the relevant callers and trace the affected paths to the narrowest shared function
or module where the invariant stays uniform. Fix it once there. Keep the regression
test at the highest seam where the reported failure stays observable. A shared
helper whose callers require different behavior does not own that invariant. If
the paths have no honest shared owner, record the tangled ownership as the Phase 6
architecture finding instead of patching every caller.

## Phase 6 — close out

Done means:

- [ ] the Phase 1 loop no longer reproduces.
- [ ] the regression test passes, or you documented its absence.
- [ ] all `[DEBUG-...]` logs are gone (one grep on the prefix) and you deleted the throwaways.
- [ ] the commit names the correct hypothesis, so the next debugger learns.

Then ask what would have *prevented* this bug. If the answer is architectural (no
seam, tangled callers), report it as a separate finding, after the fix lands.

## Finding a retired spec

The feature you are debugging may have no live spec. Bench promotes a spec, then deletes
it on merge, so a shipped feature's spec lives in git history, not in the working tree. Do not
assume nobody specified the behavior: `git log --diff-filter=D -- specs/` lists
every deleted spec, and `git log --grep=spec-retire` finds the retirement commits
(and any decision it promoted). Recover the origin spec there before you hypothesise.
For a single known slug, `bench spec history <slug>` runs both queries, merges and dedupes
them, and renders one newest-first table. This shortcut replaces a hand run of the two
commands above.

## How it meets the rest of Bench

The reviewer invokes this phase; a write delegate never charges it. When a write delegate's
repro proves the defect lives outside its ticket fence, the delegate stops implementation edits.
It keeps its
in-fence work dirty in its owned worktree. It returns a bounded blocked report: the repro
command, the red output digest, the failing surface it observed, and its in-fence dirty paths.
The reviewer runs this skill against that report to confirm the cause. The coordinator then
validates the report and reslices repair tickets per
`.agents/commands/bench-implement-spec.md`'s "When the build stops short"; this skill only
produces the report's evidence.

The Phase 1 loop joins the project gate for the fix. If the fix launches a shift, add the repro
as a test the gate runs alongside its existing checks. It is committed in the project's
expected-failure form: a quarantine marker naming the bug (a skip or expected-fail annotation
in the project's test framework). This way the committed tree stays green while the repro
survives shift rollback. An iteration that ends red rolls the worktree back to the last commit,
which destroys any uncommitted repro test. The fix's green commit removes the marker, which
turns the repro into the live regression test.

A project with no expected-failure form keeps the repro out of the shift and runs it by hand
against the fix. State that fallback in the close rather than commit a red tree; invariant
4 (commit only on green) has no red-commit exception. Route code authorship through
`craft-delegate`, including a diagnosed single-seam fix; it owns the inline threshold, worktree
isolation, and verification discipline. A repro that replaces the gate weakens the oracle;
that replacement is my call, never a debug step. The `craft-seams` skill owns the seam decision in Phase 5.
Declare the line first — a hard bug runs a high-effort shift.

Any delegation along the way (a fan-out search, a scoped fix) carries its own line. State an
explicit bound model alias on the Agent call, never the inherited default. State the effort and
the iteration cap in the charge, per `craft-line` and `craft-delegate`.
