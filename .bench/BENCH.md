# Bench Operating Guide

Bench is installed in this repo as a local agent-development workflow. The short
project instruction block in `AGENTS.md` points here instead of inlining the full
operating guide.

## The four invariants (these override convenience, always)

1. **The gate is the oracle — you never grade your own work.**
   "Done" means `bench gate` exits zero, not that you believe the work is
   finished. Tests, types, lint, and any project conformance check are the only
   things with the authority to call a shift complete. If the gate is red, the
   work is not done, regardless of how the diff looks to you. Never edit, skip,
   weaken, or delete a test or a gate check to make it pass. If a check is wrong,
   stop and say so; do not route around it.

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
- `.agents/commands/` contains portable Bench command phases.
- `.agents/skills/` contains portable Bench skills.
- `.bench/hooks/` contains shared hook scripts used by harness adapters.
- `.claude/` contains Claude Code adapter config. See `.claude/README.md`: Claude
  reads `.claude/skills/` and `.claude/commands/`, and those paths point at the
  portable `.agents/` files.
- `.codex/` contains Codex adapter config.

## Workflow

Use the canonical phases when the work needs them:

1. `/bench-ideate` for unresolved decisions.
2. `/bench-spec` to lock stories, seams, and gate expectations.
3. `/bench-build` to implement at the chosen seams.
4. `/bench-review` for semantic review against standards and spec.
5. `/bench-qa` to report the gate result.

Small mechanical changes can use the lighter path, but skipping the canonical
workflow is a reviewer decision.

## Commands

- `bench link` safely incorporates Bench into a repo.
- `bench init` scaffolds `.bench/gate.sh` and `.bench/learnings.md`.
- `bench gate` runs the oracle.
- `bench worktree` opens a reusable isolated worktree.
- `bench shift "<objective>"` runs the gated loop.

## Hook Layers

Git safety is layered:

- The git `pre-push` hook blocks direct pushes to the default branch.
- Claude Code and Codex hook adapters call the shared scripts in `.bench/hooks/`.
- The `bench shift` loop commits only after the gate is green.

Harness hooks improve ergonomics, but the git hook and gate remain the
harness-independent backstops.
