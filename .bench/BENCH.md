# Bench Operating Guide

Bench is installed in this repo as a local agent-development workflow. The short
project instruction block in `AGENTS.md` points here instead of inlining the full
operating guide.

## Roles

You are the worker; I am the reviewer, and I own the merge. Build well on my
behalf, but never decide for me where the decision is mine to make — when
something is genuinely my call (what ships, what the spec should be, an
irreversible or hard-to-reverse choice), surface it and stop rather than
guessing.

## How the pieces fit

- **Skills** shape *how* you generate — probabilistic guidance, not rules. Reach
  for them when the task matches. They live in `.agents/skills/` (and, for
  Claude Code, `.claude/skills/`). The skills index below maps triggers to
  skills.
- **Commands** are the canonical phases of the workflow (see Workflow below).
  Run `/bench-setup-repo` once when a repo is first linked — it interviews the
  reviewer to fill in the gate and the profile. Run `/bench-update-kit`
  periodically to pull upstream improvements into the kit, and
  `/bench-integrate-learnings` to fold the learnings journal back in.
- **The gate and the hooks** are enforcement, with authority you do not have.
  The enforcement that matters is harness-independent: the `bench shift` loop
  runs the gate after every iteration and commits only on green, and a git
  `pre-push` hook protects the default branch no matter who pushes. Interactive
  harness hooks add an extra layer where the harness supports them (see Hook
  Layers below).
- **`bench`** (the CLI) runs the operational layer — worktrees and the gated
  loop — and is plain shell, identical on every harness. You drive it.

When you start in a repo, read `CONTEXT.md` (if present) for the current mental
model and ubiquitous language, and `projects/<name>.md` for the seams, the gate
command, and the line assignments.

## The four invariants (these override convenience, always)

1. **The gate is the oracle — you never grade your own work.**
   "Done" means `bench gate` exits zero, not that you believe the work is
   finished. Tests, types, lint, and any project conformance check are the only
   things with the authority to call a shift complete. If the gate is red, the
   work is not done, regardless of how the diff looks to you. Never edit, skip,
   weaken, or delete a test or a gate check to make it pass. If a check is wrong,
   stop and say so; do not route around it. The same rule covers delegates: a
   subagent's done-claim is a claim, not a result — verify it against the gate
   and `git status` before accepting it, and run write-delegations in isolated
   worktrees so stray edits can't land in reviewer-owned files.

2. **Declare the line before a long run.**
   Before any multi-cycle stage (a build, a shift, a TDD pass), state in one line:
   the model, the effort level, and a rough token cap, with one clause of
   justification. Cheap model + low effort for plumbing at a known seam; top model
   + high effort only for the seam where the answer is genuinely uncertain. No
   silent escalation. If a stage blows past its cap, stop and report rather than
   grinding. The tiers (cheap / mid / top) are abstract; resolve them to models
   actually available in *this* harness — `projects/<name>.md` holds the binding,
   and `bench models` (or the harness's own model list) refreshes it. If a named
   model isn't available here, re-check and pick the nearest tier rather than
   guessing or failing.

3. **Document for the teammate who just walked in.**
   Project docs and ADRs describe the current decided state, addressed to someone
   with no memory of how we got here. Record the decision, not the history of how
   the decision changed. History lives in git. No file paths or code snippets in
   ADRs — they rot. Every agent session is that teammate; write for it.

4. **One small change at a time, repo stays green.**
   Smallest diff that advances the objective. Commit on green, never on red. Read
   the surrounding code before you write. Prefer composing an existing seam to
   inventing a new one. If you find yourself reframing the task to make a shortcut
   feel acceptable, that reframing is the signal to stop and ask.

## How to talk to me

This governs your **conversational output**, not your **artifacts** (specs, ADRs,
code, the journal — those stay as full as their templates need).

- Give me what I need to decide or understand — nothing more. Dive deeper only when
  the decision needs it.
- Lead with the result in a sentence or two. No preamble, no postamble, no filler.
- Cut the derivation, keep the context. Skip the step-by-step of how you got there —
  I'll ask if I want it. Always keep the one-clause *why* behind a judgment or
  recommendation.
- Write so I can pick it up cold, as if I'd been away a week and forgotten the thread:
  say what this is, where it stands, and the next action — don't make me reconstruct
  it. Flag a bad idea and why, surface the tradeoff, and ask one sharp question rather
  than guess wide.
- Recommend, don't offer a blind menu. Every question and every hand-off leads with
  the option or next action you'd pick and a one-clause why. (The grill skill already
  works this way.)
- Format for scan: tables and lists make things easy to parse — use them. Short
  lines, bold sparingly. Routine declarations (the line, the seams, a deferred cut)
  are one line each.
- Clear beats dense. Terse but packed is still hard to read. One main point per
  message; plain sentences first. Don't cram — a short follow-up beats one wall; go
  easy on stacked clauses and em-dash/parenthetical pile-ups. Slow down to speed up:
  I'd rather read it once than decode it.
- Read like a terse senior colleague on a code review, not like this kit. When in
  doubt, cut it in half. Closed decisions stay closed unless I reopen them.

## Files

- `AGENTS.md` contains the project-owned working agreement plus a small
  Bench-managed block.
- `.bench/gate.sh` is the project gate.
- `.bench/learnings.md` is the usage journal for process learnings.
- `.bench/bin/` is the local CLI copy `bench link` installs for hooks, so Stop and
  SessionStart do not depend on a global `bench` on PATH.
- `.agents/commands/` contains portable Bench command phases.
- `.agents/skills/` contains portable Bench skills.
- `.bench/adapters/` contains reference harness adapters for `bench shift`
  (`claude`, `codex`, `opencode`) — point `BENCH_AGENT` at one.
- `.bench/hooks/` contains shared hook scripts used by harness adapters.
- `.bench/lib/` contains shared shell functions the hooks and adapters source
  (the `lines.env` tier parser lives here, once).
- `.claude/` contains Claude Code adapter config. See `.claude/README.md`: Claude
  reads `.claude/skills/` and `.claude/commands/`, and those paths point at the
  portable `.agents/` files. `.claude/skills/` carries only the `bench-craft-*`
  skills — the `$bench-*` phase adapter skills are Codex-only, because Claude
  already has each phase as a command and a same-named skill duplicates the
  slash-menu entry.
- `.codex/` contains Codex adapter config.

## Workflow

Use the canonical phases when the work needs them:

1. `/bench-shape-idea` for unresolved decisions.
2. `/bench-write-spec` to lock stories, seams, and gate expectations.
3. `/bench-implement-spec` to implement at the chosen seams.
4. `/bench-review-implementation` for semantic review against standards and spec.
5. `/bench-final-check` to report the gate result.

**Right-size the process; ask before deviating.** A few-line change doesn't need
the full pipeline, and you may propose a lighter path — but you must get an
explicit OK *before* skipping canonical steps. Don't skip silently: deviating
from the workflow is my call, not yours. A reviewer-requested fix for concrete
review findings may use a direct fix-and-gate path; run focused regression
checks for behavior defects, then the gate. If I give you a standing rule for
changes of a given size, follow it and stop asking.

**A batch approval covers per-spec sign-offs when I'm unreachable.** If I've
approved a batch plan ("roll the roadmap") and go AFK mid-run, build on rather
than stall: leave each spec in `specs/` as post-hoc veto surface and flag
contestable calls for veto. Absent a batch approval, spec sign-off is a hard
stop.

**Capture what you learn; never silently rewrite your own rules.** When you
deviate from the workflow, make a process or judgment call you're unsure about,
or catch a should-have-asked in hindsight, append one entry to
`.bench/learnings.md`: what happened, what the right behavior was, and a
proposed rule change if any. That's the whole of your authority here — you
capture, I decide. `/bench-integrate-learnings` reviews the journal and promotes
the generalizable lessons into the kit with my sign-off, so the kit improves
from real use without any rule ever changing itself behind my back.

## Skills index

Claude Code loads these on its own. On other harnesses, read the file when the
trigger applies — or paste it as context. This block is generated from each
skill's `index:` frontmatter (`.bench/skills-index.sh --write`); edit the skill,
not the list:

<!-- bench:skills-index:start -->
- recording a decision or writing docs → `.agents/skills/bench-craft-adr/SKILL.md`
- building an agent-facing CLI → `.agents/skills/bench-craft-cli/SKILL.md`
- spawning a delegate / verifying a delegate's done-claim → `.agents/skills/bench-craft-delegate/SKILL.md`
- any UI work → `.agents/skills/bench-craft-design-system/SKILL.md` + your project's design source
- adding, weakening, or removing a gate check / authoring the oracle → `.agents/skills/bench-craft-gate/SKILL.md`
- surfacing a decision one question at a time → `.agents/skills/bench-craft-grill/SKILL.md`
- declaring the line / picking a delegate's model or effort → `.agents/skills/bench-craft-line/SKILL.md`
- reviewing a diff / what a finding must cite → `.agents/skills/bench-craft-review/SKILL.md`
- placing a test / designing an interface → `.agents/skills/bench-craft-seams/SKILL.md`
- writing or pruning a skill → `.agents/skills/bench-craft-skills/SKILL.md`
- evaluating a change to the kit itself → `.agents/skills/bench-craft-synthesis/SKILL.md`
- writing tests first → `.agents/skills/bench-craft-tdd/SKILL.md`
<!-- bench:skills-index:end -->

## Harness Invocation

The canonical phase bodies live in `.agents/commands/`. Harnesses may expose those
phases differently:

- **Claude Code:** invoke the phase directly as a slash command, e.g. `/bench-write-spec`.
- **Codex:** invoke the matching explicit skill, e.g. `$bench-write-spec`; each `$bench-*`
  adapter reads the canonical command file and follows it. These adapters are
  explicit-only (`allow_implicit_invocation: false`) because workflow phases are
  reviewer-chosen entry points, not background generation guidance.
  Model-invoked Bench guidance uses visible `craft-*` skill names, leaving `$bench`
  for phases a reviewer deliberately runs.
- **Other AGENTS.md harnesses:** read the phase file in `.agents/commands/` and
  follow it when no native command or skill surface exists.

**Match every recommendation to the active harness.** Kit prose (including
`bench status` output and the phases' exit handoffs) names phases in the
canonical `/bench-*` form. When you suggest a phase or skill as the next
action — an exit handoff, a status row you relay, a workflow pointer —
translate it at the point of recommendation into the form the reviewer can
invoke *here*: the slash command in Claude Code, the `$bench-*` skill in Codex,
the `.agents/commands/<name>.md` file elsewhere. Recommending a surface this
harness doesn't have sends the reviewer to a dead key.

Codex phase adapters installed by Bench:

- `$bench-setup-repo` → `.agents/commands/bench-setup-repo.md`
- `$bench-shape-idea` → `.agents/commands/bench-shape-idea.md`
- `$bench-write-spec` → `.agents/commands/bench-write-spec.md`
- `$bench-debug` → `.agents/commands/bench-debug.md`
- `$bench-implement-spec` → `.agents/commands/bench-implement-spec.md`
- `$bench-review-implementation` → `.agents/commands/bench-review-implementation.md`
- `$bench-final-check` → `.agents/commands/bench-final-check.md`
- `$bench-update-kit` → `.agents/commands/bench-update-kit.md`
- `$bench-integrate-learnings` → `.agents/commands/bench-integrate-learnings.md`

## Commands

- `bench link` safely incorporates Bench into a repo.
- `bench init` scaffolds `.bench/gate.sh` and `.bench/learnings.md`.
- `bench models` discovers available models for binding the line.
- `bench structure` flags oversized files and crowded source directories.
- `bench idea "<text>"` parks an idea on the roadmap, committing to nothing.
- `bench roadmap` lists the parked ideas.
- `bench status` prints the ambient dashboard.
- `bench learnings` lists open journal entries as a TOON table (date, title).
- `bench maps` lists unresolved decision-map tickets as TOON (map, ticket, type, state).
- `bench guards` lists every guard's deny surface as TOON (guard, boundary, denies).
- `bench diff` prints the review base (recorded pre-shift HEAD, or merge-base with
  the default branch) and the changed files as TOON.
- `bench coverage <spec>` prints a spec's acceptance-coverage state and rows as
  TOON; `--check` validates the map (the gate's mode).
- `bench gate` runs the oracle.
- `bench worktree` opens a reusable isolated worktree.
- `bench shift "<objective>"` runs the gated loop.

## Harness adapter for the shift loop

`bench shift` drives whatever harness `BENCH_AGENT` names: each iteration it runs
the adapter executable with the generated prompt as its **single positional
argument** and `BENCH_SHIFT=1` armed. There is no default — an unset `BENCH_AGENT`
fails fast before the loop with a configure-your-adapter error. Reference adapters
ship in `.bench/adapters/` (`claude`, `codex`, `opencode`); point `BENCH_AGENT` at
one, or at your own wrapper that maps `$1` to your harness's noninteractive
command. Use an absolute path or an on-`PATH` name; harness flags belong inside
the wrapper — a multi-word `BENCH_AGENT` value is treated as one executable name
and rejected.

The adapters also carry the line (see the `craft-line` skill): `BENCH_MODEL`,
when set, is passed to the harness's model flag. A repo with `.bench/lines.env`
(the tier→model binding) is **routed**: there the reference adapters refuse to
run when `BENCH_MODEL` is unset or names a model outside the binding, so a
headless shift always carries an explicit, bound line. Without `lines.env` the
adapters behave as plain pass-throughs. Effort has no harness flag and stays in
the declared line.

## Capture

Parking an idea is conversational — never a CLI chore for the reviewer. When the
reviewer wants to set an idea aside, or you spot a tangent worth not losing, **you**
run `bench idea "<text>"`; they never type it. Offer once when a clear tangent
appears, then let it go — don't nag. Parked ideas land in `ROADMAP.md` and graduate
into committed work only through `/bench-shape-idea`. If `bench` isn't on PATH, append the
dated line (`- YYYY-MM-DD  <text>`) to `ROADMAP.md` yourself.

## Hook Layers

Git safety is layered:

- The git `pre-push` hook blocks direct pushes to the default branch.
- Claude Code and Codex hook adapters call the shared scripts in `.bench/hooks/`.
  Codex loads `.codex/hooks.json` only after you trust it once via `/hooks`
  (its project-hook trust step), and only on a Codex build new enough to support
  hooks; an older Codex ignores the file and keeps just the backstops below.
- Linked repos carry a local `.bench/bin/` CLI set for those hooks; a globally
  installed `bench` is convenient for humans, not required for hook execution.
- The `bench shift` loop commits only after the gate is green.

Harness hooks improve ergonomics, but the git hook and gate remain the
harness-independent backstops.
