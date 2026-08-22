---
description: Run the external gate and commit work on green; after a spec's final landing, report the evidence and capture the retro. Never use the model's own judgment as a substitute for retained or freshly observed evidence.
---

# /bench-final-check — the gate is the oracle

## Entry orientation

This is the final verification phase. For a reviewed spec-backed integration
source, `bench worktree land` already ran the external gate, published the
implemented spec, and released the source; final-check reports that landing's
retained evidence and captures the retro. Other work lands only on green. It
does not substitute model judgment for tests, types, lint, or conformance.

## Exit handoff

Close by reporting the applicable oracle result plainly. A reviewed spec's
published `bench worktree land` commit is the `Status: implemented` author and
gets no second gate or landing mutation afterward: do not re-run `bench gate`
over its unchanged tree or run `bench spec implemented` on top of it. Report the
reviewed source pair, landing commit, and retained exact green evidence, then
capture the retro below.

Everything else takes the gate-then-commit path. On
green, land the named paths with `bench commit -m "<msg>" <path>...`; it gates
and commits them atomically. The honest no-op runs `bench gate` and reports its
verdict when there is nothing to commit. If the command refuses over an
unexplained working-tree file, surface that file without committing or reverting
it. Then hand back for the reviewer to merge or decide what ships. On red, report
the first failing check and smallest reproduction, then recommend the fitting
repair command: usually `/bench-implement-spec` for feature work or `/bench-debug`
for a bug.

**The post-merge tail (exit duty).** After the green landing
reaches the default branch, read `bench status` and run the housekeeping rows it flags before
closing: a merged spec awaiting retirement gets `bench spec retire <slug>` and
its `spec-retire: <slug>` commit — promoting durable content first (a decision
to an ADR, a hostile edge to the profile); retiring the whole
`specs/<slug>/` folder removes its compiled decision provenance with its tickets,
so there is no separate top-level decision-map delete; an orphaned review pickup is promoted or
deleted by hand; scratch branches go through `bench worktree clean`; leftover worktrees are retired by `bench worktree clean --landed`: run the plan, apply it, and carry the plan and apply result in the landing report. Leave the roadmap and capture rows to
`/bench-drain` — that phase owns the reconcile and the drain, and this duty
never restates it. On a topic branch these duties defer by design: the rows
fire only on the default branch, and the next default-branch session's
SessionStart status re-surfaces them — state the deferral in the close instead
of silently skipping it.

## Capture the implementation retro

After any applicable post-merge tail, an implemented spec has two final exit
duties: rewrite `capture/retros/<spec-slug>.md` in full. Do this only after the
spec's final green landing commit has flipped it to `Status: implemented`. A
re-run replaces that slug's whole file; it never appends,
and it leaves other pending retros untouched.

Use these headings exactly:

```markdown
## Outcome

## Gate-stage timings

## Ticket-versus-spec-slice and delegate performance

## Coordinator catches

## Repair attribution

## Agent-experience improvements

### Bench CLI

### Skills

### Process
```

Record concrete evidence: what landed; elapsed time for each measured gate
stage; how ticket-sized delegate charges performed against charges handed a
spec slice; what the coordinator caught while accepting delegate claims; and
specific improvements to Bench CLI, skills, and process, with the friction and
expected effect named.

Write each improvement item as one list item. Give the item one sentence that
states the change to make. End the item with one line that reads `Feeds: FT<n>`,
`Feeds: new`, or `Feeds: none`. Use `FT<n>` for the roadmap row the change
feeds, `new` for a row the drain must open, and `none` for a change that needs
no row. The gate reds an item that ends without that line. Each item under the
three improvement headings takes this shape:

```markdown
- <one sentence that states the change to make>
  Feeds: none
```

Under that repair-attribution heading, write one table row per ticket in the
build: the ticket, how many repair rounds it took to land, and one cause per
round. A ticket that landed in one pass records `none`. Causes come from this
vocabulary and no other, one term per round: `shaping-ambiguity`, `spec-row`,
`ticket-slicing`, `tree-drift`, `delegate-error`, `other`. This template is the
single guidance source for that vocabulary; the anchors registry needle pinning
those terms is its enforcement copy, not a second source. A later reader of the
drained tables reads the terms from the tables themselves.

Then read `capture/agent-performance/README.md` and both provider scorecards.
Refresh the scorecard for every provider whose models served as implementer,
reviewer, or orchestrator on this landing. Update its
last-incorporated-landing line and fold the new evidence into affected aggregates;
leave an uninvolved provider unchanged. Completion means every participating
model/effort/role is accounted for without adding a per-run diary row.

These files are pending capture for `/bench-drain`, not
a second roadmap. Do not run another gate or commit just to capture the retro;
the successful landing boundary is already the verdict, and the retro leaves
through the next reviewer-approved capture drain.

Report the applicable oracle. This command does not form an opinion about whether
the work is good: it reports the gate's retained or fresh result.

## Run it

For work that has
paths to land, the oracle run and landing are one command:

```sh
bench commit -m "<msg>" <path>...
```

`bench commit` runs the gate itself and commits only on green — a red run
reports its own first failing phase and refuses. Don't run `bench gate` first;
the commit is the gate run (a fresh green verdict for the identical tree is
reused, never re-paid). Standalone `bench gate` has two jobs here: the honest
no-op (nothing left to commit, the verdict still gets reported) and diagnosing
a red run.

The gate itself is an executable `.bench/gate.sh` when present, else the
`$BENCH_GATE` command string, else stack auto-detect (typecheck → test →
lint). `projects/<name>.md` documents what the gate covers — it never selects the
gate; to change what runs, change `.bench/gate.sh`.

## Report

- **Spec landed:** report the final landing commit and its retained exact green
  evidence plainly; capture the retro without another gate
  or commit.
- **Ordinary green:** the work is committed; state it plainly, and add one line
  noting that ship-tier verification has not run — dev green claims the kit
  works from the tree; release-evidence checks run once per release under
  `bench prep-release`. A statement, not an approval prompt. Hand back for me to merge.
- **Red:** report each failing check in the order it fails, with the smallest
  reproduction. Do not propose weakening the check. Diagnose the cause, propose a
  fix at the seam, and — if I approve — fix it and re-run the gate. A fix is only
  real when the gate is green again.

If a check itself looks wrong (a flaky test, an over-tight lint rule), say so
explicitly and stop. Changing a gate check is my call, not a step inside `/bench-final-check`;
when I approve one, the `craft-gate` skill governs how it's made.

## Findings that the gate can't see

If verifying surfaces a design problem the tests pass through, name it as a
finding for `/bench-review-implementation` — that phase owns semantic review;
don't fold it silently into a fix.
