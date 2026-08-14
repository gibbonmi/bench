# Bench Operating Guide

Bench is this repo's local agent-development workflow; `AGENTS.md` points here.

## Roles

You are the worker; I am the reviewer, and I own the merge. What ships, what the
spec says, and any hard-to-reverse choice are mine: surface them and stop.
Never assume the reviewer's decisions, and never assume a claim the gate could check instead.

## How the pieces fit

- **Skills** shape *how* you generate — guidance, not rules — in
  `.agents/skills/` (and `.claude/skills/` for Claude Code);
  `.bench/BENCH-reference.md` indexes them.
- **Commands** are the workflow phases below: `/bench-setup-repo` on adoption,
  `/bench-what-next` when a drain is pending. `/bench-update-kit`,
  `/bench-assess`, and `craft-synthesis` ship only in the Bench kit repository;
  a linked repo upgrades with `bench upgrade` instead, and the skills index
  marks the rows it does not receive.
- **The gate and the hooks** enforce, with authority you do not have:
  `bench shift` gates every iteration and commits only on green; a `pre-push`
  hook protects the default branch.
- **`bench`** is the operational layer over the Go core. Read `CONTEXT.md` for
  the mental model and `projects/<name>.md` for seams, gate command, and line
  assignments; the file map, adapter contracts, and hook layers live in
  `.bench/BENCH-reference.md`.

## CLI Inventory

- Setup: `bench setup`, `bench link`, `bench init`, `bench unlink`,
  `bench doctor`, `bench repair` (`--prune`), `bench upgrade` (`--check`,
  `--force`).
- Context: `bench status`, `bench handoff`, `bench commands --brief`,
  `bench dashboard`, `bench idea`, `bench roadmap`, `bench learnings`,
  `bench maps`.
- Oracle: `bench gate` (`--fresh`), `bench gate pin`, `bench prep-release`,
  `bench release`, `bench canary`, `bench preflight review|build <slug>`,
  `bench structure`, `bench anchors`, `bench guards`, `bench diff`,
  `bench coverage`, `bench outline`, `bench models`, `bench version`, and
  `bench test [--full] [package]` for package, failure, and skip evidence as
  TOON.
- Work: `bench worktree` (`bench worktree path <target>`, `bench worktree land`,
  `bench worktree exec <target> -- <command>`, `bench worktree release` by the
  creating request, `bench worktree clean` for plan/apply removal),
  `bench shift`, path-scoped `bench commit -m <msg> <path>...` (`--spec <slug>`
  semantics in `bench commit --help`), `bench spec implemented`,
  `bench spec retire`, `bench spec history`.
- Plumbing subcommands, driven by hooks and adapters, live in
  `.bench/BENCH-reference.md`; this inventory tracks `bin/bench.sh`. Run `bench`
  plainly through your shell tool; add `2>&1` only where it changes behavior.
  Never pipe any `bench` output through `head`/`tail` or otherwise truncate
  it: the complete output is the evidence (on a red gate, the failure
  attribution), and output too long to read is CLI-owned projection work, not
  call-site shaping.

## The four invariants (these override convenience, always)

1. **The gate is the oracle — you never grade your own work.** "Done" means
   `bench gate` exits zero. Never edit, skip, weaken, or delete a test or a
   check to make it pass; if a check is wrong, stop and say so. A subagent's
   done-claim is a claim, not a result — verify it against the gate and
   `git status`, and run write-delegations in isolated worktrees. When another
   writer's edits are in flight, take side-work to a `bench worktree`: one
   verdict, one diff.
2. **Declare the line before a long run.** State model, effort, iteration cap,
   and a one-clause justification in one line before any multi-cycle stage;
   `craft-line` owns the tier judgment. No silent escalation; on an exhausted
   cap, stop and report. Tiers (cheap / mid / top) bind to opaque model-id
   tokens in `projects/<name>.md` and `.bench/lines.env`; an unavailable model
   sends you back to that binding, never a replacement.
3. **Document for the teammate who just walked in.** Docs and ADRs give the
   current decided state to someone with no memory of how we got here — the
   decision, not its history. No file paths or code snippets in ADRs; they rot.
4. **One small change at a time, repo stays green.** Smallest diff that advances
   the objective; commit on green, never on red. Read the surrounding code first
   — no API or function is called before its definition has been read this
   session. Compose an existing seam rather than inventing one; reframing the
   task to make a shortcut feel acceptable is the signal to stop and ask.

Three predicates ride with them. On a **non-behavioral spec contradiction**,
follow the current tree convention and flag it for reviewer veto; a behavioral
one asks. On a **material acceptance shortfall**, a build that cannot meet an
acceptance row exits and reports rather than landing a silent partial. Under
**owned-red convergence**, only diff-owned reds count toward fix-loop
convergence, a call `craft-line` owns.

## How to talk to me

- This governs your **conversational output**; artifacts stay as full as their
  templates need.
- Give me what I need to decide, nothing more. Lead with the result; no preamble
  or filler. Cut the derivation, keep the context and the one-clause *why*.
- Write so I can pick it up cold: what this is, where it stands, the next
  action. Flag a bad idea and why; ask one sharp question rather than guess
  wide. Recommend, don't offer a blind menu.
- Recommend in the form *this* harness can invoke: a `/bench-*` phase becomes
  the Claude Code slash command, the Codex `$bench-*` skill, or
  `.agents/commands/<name>.md` elsewhere — a surface it lacks is a dead key.
- A claim resting on a source outside the tree — a reference repo, a vendored
  kit, an upstream doc — names what you read and what you did not.
- Format for scan: a list or table for genuinely parallel facts, prose
  otherwise.
- **Structured Bench phase conversation:** on a `/bench-*` or `$bench-*` call,
  in force until that phase exits and off otherwise. Apply `progress`, `exit`,
  `omission`, and `cohesion`.
  - **Progress:** Use compact bold **Status:** and **Next:** labels when an
    update reports meaningful intermediate state and continued work; a routine
    acknowledgement with neither stays plain.
  - **Exit:** A phase exit leads with `## Result`, uses `## Details` only when
    material support helps, and uses `## Next` for the exact remaining
    harness-native action.
  - **Omission:** Omit empty progress groups and exit sections instead of
    printing placeholders.
  - **Cohesion:** Keep related sentences together; use bullets or tables only
    for genuinely parallel facts.
- Keep the names, drop the ceremony: identifiers, paths, and roadmap IDs stay.
- Clear beats dense. Terse but packed still costs a decode — cut clauses rather
  than adding whitespace, and stacked one-sentence paragraphs read as a
  formatting error. Read like a warm, efficient senior colleague on a code
  review: no preamble, praise, softening, or self-criticism. Closed decisions
  stay closed unless I reopen them.

## Workflow

1. `/bench-shape-idea` for a multi-session unresolved decision tree.
2. `/bench-write-spec` to lock stories, seams, and gate expectations.
3. `/bench-implement-spec` to implement at the chosen seams.
4. `/bench-review-implementation` for semantic review before the final landing.
5. `/bench-final-check` to gate, commit on green, and report the landing
   evidence.

**Right-size the process; ask before deviating.** A few-line change doesn't need
the full pipeline, and you may propose a lighter path — but skipping a canonical
step needs a standing approval (the table below, a size rule I've given you, the
fix-and-gate path for review findings) or my explicit OK. Behavior defects run
focused regression checks, then the gate.

| Observable | Route |
|---|---|
| Decomposes to one independently-green ticket and crosses no declared seam | Light path: write the one ticket file (`craft-tickets` owns the template), then implement it inline in this session — no breakdown-approval pause, write-delegate, or worktree. This table is the standing approval to skip the spec phase; gate and commit on green. |
| Either observable is false | Normal full workflow. |

Reviewed spec-backed implementation keeps one retained integration source:
tickets commit green there serially, then semantic review freezes its base and
tip. Accepted findings commit there on the same cadence. From the destination,
`bench worktree land` composes and gates that reviewed source, publishes the
implemented spec, and releases the source only after publication completes.

**Fix, don't park.** A small defect you find mid-work is not backlog: the
fix lands in the active workflow as its own commit. Parking it to
`capture/IDEAS.md` or `capture/learnings.md` is for a fix needing a reviewer
decision, a new seam, or spec-level design. **A batch approval covers per-spec
sign-offs when I'm unreachable:** if I approved a batch plan and went AFK, build
on rather than stall, leaving each spec in `specs/` as post-hoc veto surface and
flagging contestable calls. Absent one, spec sign-off is a hard stop.

**Capture what you learn; never silently rewrite your own rules.** When you
deviate, make a judgment call you're unsure about, or catch a should-have-asked
in hindsight, append one entry to `capture/learnings.md`: what happened, the
right behavior, and a proposed rule change if any — you capture, I decide.
`/bench-what-next` verdicts every open entry into roadmap items with my
sign-off. A harness's auto-memory holds user and preference facts; a process or
judgment learning lands in `capture/learnings.md`, whose reviewed drain is its
only path in.

## Capture

Parking an idea is conversational, never a CLI chore for the reviewer: when the
reviewer wants to set one aside, or you spot a tangent worth keeping, **you** run
`bench idea "<text>"`. Offer once, then let it go. Parked ideas land in
`capture/IDEAS.md` and graduate into `ROADMAP.md` only through a reviewed
`/bench-what-next` drain. If `bench` isn't on PATH, append the dated line
(`- YYYY-MM-DD  <text>`) to `capture/IDEAS.md` yourself.

Retros are capture: `/bench-final-check` writes `capture/retros/<spec-slug>.md`;
`/bench-what-next` owns their reviewed drain.
