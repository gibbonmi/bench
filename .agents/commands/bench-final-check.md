---
description: Run the external gate and commit work on green; after a spec's final landing, report the evidence and capture the retro. Never use the model's own judgment as a substitute for retained or freshly observed evidence.
---

# /bench-final-check — the gate is the oracle

## Entry orientation

This is the final verification phase. For a reviewed spec-backed integration
source, `bench worktree land` already ran the external gate, published the
implemented spec, and released the source. Final-check reports that landing's
retained evidence and captures the retro. Other work lands only on green. It
does not substitute the model's own judgment for tests, types, lint, or
conformance checks.

## Exit handoff

Close by reporting the applicable oracle result plainly. A reviewed spec's
published `bench worktree land` commit is the only `Status: implemented` author.
It gets no second gate or landing mutation afterward. Do not re-run `bench gate`
over its unchanged tree.

Report the reviewed source pair, the landing commit, and the retained exact
green evidence. Then capture the retro below.

Everything else takes the gate-then-commit path. On green, land the named
paths with `bench commit -m "<msg>" <path>...`. This command gates and commits
them atomically. When there is nothing to commit, the honest no-op runs
`bench gate` and reports its verdict.

If the command refuses because of an
unexplained working-tree file, surface that file. Do not commit or revert it.
Then hand back to the reviewer to merge or to decide what ships. On red,
report the first failing check and the smallest reproduction. Then recommend
the fitting repair command: usually `/bench-implement-spec` for feature work,
or `/bench-debug` for a bug.

**The post-merge tail (exit duty).** After the green landing reaches the
default branch, read `bench status` and run the housekeeping rows it flags
before you close. A merged spec awaiting retirement gets `bench spec retire <slug>`
and its `spec-retire: <slug>` commit. Promote durable content first,
for example a decision to an ADR or a hostile edge to the profile. Retirement
of the whole `specs/<slug>/` folder removes its compiled decision provenance
with its tickets, so there is no separate top-level decision-map delete.
Promote or delete an orphaned review pickup by hand.

Scratch branches go through `bench worktree clean`.
Leftover worktrees are retired by `bench worktree clean --landed`: run the plan, apply it, and carry the plan and apply result in the landing report.

Leave the roadmap and capture rows to
`/bench-drain`; that phase owns the reconcile and the drain, and this duty
never restates it. On a topic branch these duties defer by design: the rows
fire only on the default branch. The next default-branch session's
SessionStart status re-surfaces them. State the deferral in the close instead
of a silent skip.

## Capture the implementation retro

After any applicable post-merge tail, an implemented spec has two final exit
duties. First, rewrite `capture/retros/<spec-slug>.md` in full. Do this only after the
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

Record concrete evidence:

- what landed
- elapsed time for each measured gate stage
- how ticket-sized delegate charges performed against charges handed a spec slice
- what the coordinator caught while accepting delegate claims
- specific improvements to Bench CLI, skills, and process, with the friction and expected effect named

Write each improvement item as one list item. Give the item one sentence that
states the change to make. End the item with one line that reads `Feeds: FT<n>`,
`Feeds: new`, or `Feeds: none`. Use `FT<n>` for the roadmap row the change
feeds, `new` for a row the drain opens, and `none` for a change that needs
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

Second, read `capture/agent-performance/README.md` and both provider scorecards.
Refresh the scorecard for every provider whose models served as implementer,
reviewer, or orchestrator on this landing. Update its
last-incorporated-landing line and fold the new evidence into affected aggregates;
leave an uninvolved provider unchanged. Completion means you accounted for
every participating model/effort/role without a per-run diary row.

These files are pending capture for `/bench-drain`, not
a second roadmap. Do not run another gate or commit just to capture the retro.
The successful landing boundary is already the verdict. The retro leaves
through the next reviewer-approved capture drain.

Report the applicable oracle result. This command does not form an opinion
about whether the work is good. It reports the gate's retained or fresh
result.

## Run it

For work that has
paths to land, the oracle run and landing are one command:

```sh
bench commit -m "<msg>" <path>...
```

`bench commit` runs the gate and commits only on green. A red run reports its
own first failing phase and refuses to commit. Do not run `bench gate` first.
The commit already is the gate run; the gate reuses a fresh green verdict for
the identical tree and never re-pays it. Standalone `bench gate` has two jobs
here: report the honest no-op, when nothing is left to commit, and diagnose a
red run.

Exit 3 means the commit is published but the checkout did not reconcile. Paste
the `next=` restore command from the `committed{...}` record to repair the
checkout.

The gate itself is an executable `.bench/gate.sh` when present. Otherwise it
is the `$BENCH_GATE` command string. Otherwise it is stack auto-detect:
typecheck, then test, then lint. `projects/<name>.md` documents what the gate
covers; it never selects the gate. To change what runs, change
`.bench/gate.sh`.

## Report

- **Spec landed:** report the final landing commit and its
  retained exact green evidence plainly. Capture the retro without another
  gate or commit.
- **Ordinary green:** the work is committed. State it plainly, and add one line
  noting that ship-tier verification has not run. A dev green claim shows the
  kit works from the tree. Release-evidence checks run once per release under
  `bench prep-release`. This is a statement, not an approval prompt. Hand back to me to merge.
- **Red:** report each failing check in the order it fails, with the smallest
  reproduction. Do not propose a weaker check. Diagnose the cause, and propose a
  fix at the seam. If I approve, fix it and re-run the gate. A fix is real
  only when the gate is green again.

If a check itself looks wrong, for example a flaky test or an over-tight lint
rule, say so explicitly and stop. A change to a gate check is my call, not a
step inside `/bench-final-check`. When I approve one, the `craft-gate` skill
governs how you make it.

## Findings that the gate can't see

If verification surfaces a design problem that the tests pass through, name it
as a finding for `/bench-review-implementation`. That phase owns semantic
review. Do not fold it silently into a fix.
