# Bench Operating Guide

Bench is installed in this repo as a local agent-development workflow. The short
project instruction block in `AGENTS.md` points here instead of inlining the full
operating guide.

## Core Rules

- The **gate** is the oracle. Work is done only when `bench gate` exits zero.
- The reviewer owns merge decisions. Agents build and report; they do not decide
  what ships.
- Use one small change at a time, and keep the repo green.
- Before a long run, declare the line: model, effort, rough token cap, and why.

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

1. `/start-ideation` for unresolved decisions.
2. `/spec` to lock stories, seams, and gate expectations.
3. `/build` to implement at the chosen seams.
4. `/prep-shift` for semantic review against standards and spec.
5. `/verify-gate` to report the gate result.

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
