---
description: Report a terminal spec-build promotion and capture its retro, or run the external gate and commit ordinary non-lifecycle work on green. Never use the model's own judgment as a substitute for retained or freshly observed evidence.
---

# /bench-final-check — the gate is the oracle

## Entry orientation

This is the final verification phase. For a reviewed spec build it reports the
terminal result already authored by promotion; for light-path and ordinary work
it runs the external gate and lands only on green. It does not substitute model
judgment for tests, types, lint, or project conformance.

When a spec-build slug is in scope, begin with
`bench spec build status <slug> --full`. A terminal promoted run takes the retained-evidence route below. An
active nonterminal run stops here and reports the durable next action from that
projection; return to assignment, checkpoint, integration, review, or promotion
as named instead of entering the ordinary landing path. An empty lifecycle is
not permission to land reviewed spec-backed work through an older command: start
or resume it through `/bench-implement-spec`. Only light-path and ordinary
non-lifecycle work take the fresh gate-then-commit route.

## Exit handoff

Close by reporting the applicable oracle result plainly.
`bench spec build promote` is the sole spec-backed gate, commit, and `Status: implemented` author.
A terminal promoted run gets no second gate or landing mutation: do not run
`bench gate`, `bench commit --spec`, or `bench spec implemented`. Report the
promotion subject, published working-branch commit, and retained exact green
evidence from `status --full`, then capture the retro below. If the run is active
but nonterminal, report its durable next action and stop without changing state.

Light-path and ordinary non-lifecycle work retain the gate-then-commit path. On
green, land the named paths with `bench commit -m "<msg>" <path>...`; it gates
and commits them atomically. The honest no-op runs `bench gate` and reports its
verdict when there is nothing to commit. If the command refuses over an
unexplained working-tree file, surface that file without committing or reverting
it. Then hand back for the reviewer to merge or decide what ships. On red, report
the first failing check and smallest reproduction, then recommend the fitting
repair command: usually `/bench-implement-spec` for feature work or `/bench-debug`
for a bug.

**The post-merge tail (exit duty).** After the promoted or ordinary green landing
reaches the default branch, read `bench status` and run the housekeeping rows it flags before
closing: a merged spec awaiting retirement gets `bench spec retire <slug>` and
its `spec-retire: <slug>` commit — promoting durable content first (a decision
to an ADR, a hostile edge to the profile); retiring the whole
`specs/<slug>/` folder removes its compiled decision provenance with its tickets,
so there is no separate top-level decision-map delete; an orphaned review pickup is promoted or
deleted by hand; leftover worktrees and scratch branches go through
`bench worktree clean`. Leave the roadmap and capture rows to
`/bench-what-next` — that phase owns the reconcile and the drain, and this duty
never restates it. On a topic branch these duties defer by design: the rows
fire only on the default branch, and the next default-branch session's
SessionStart status re-surfaces them — state the deferral in the close instead
of silently skipping it.

## Capture the implementation retro

After any applicable post-merge tail, a promoted spec build has one last exit
duty: rewrite `.bench/retros/<spec-slug>.md` in full. Do this only after
`bench spec build status <slug> --full` reports a terminal promoted run and its retained
exact green evidence. A re-run replaces that slug's whole file; it never appends,
and it leaves other pending retros untouched.

Use these headings exactly:

```markdown
## Outcome

## Gate-stage timings

## Ticket-versus-spec-slice and delegate performance

## Coordinator catches

## Agent-experience improvements

### Bench CLI

### Skills

### Process
```

Record concrete evidence: what landed; elapsed time for each measured gate
stage; how ticket-sized delegate charges performed against charges handed a
spec slice; what the coordinator caught while accepting delegate claims; and
specific improvements to Bench CLI, skills, and process, with the friction and
expected effect named. This file is pending capture for `/bench-what-next`, not
a second roadmap. Do not run another gate or commit just to capture the retro;
the successful landing boundary is already the verdict, and the retro leaves
through the next reviewer-approved capture drain.

Report the applicable oracle. This command does not form an opinion about whether
the work is good: it reports retained promotion evidence or a fresh ordinary
gate result.

## Run it

For a spec-build terminal report, the read-only command is:

```sh
bench spec build status <slug> --full
```

Do not reconstruct or rerun its promotion evidence. For ordinary work that has
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

- **Promoted:** report the terminal promotion subject, published commit, and
  retained exact green evidence plainly; capture the retro without another gate
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
