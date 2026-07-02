# Working agreement

This is the canonical instruction file for project-owned content, read by every harness
(Claude Code via `@AGENTS.md` and `@.bench/BENCH.md` imports in `CLAUDE.md`; Codex,
OpenCode, and other AGENTS.md harnesses natively). The shared platform rules — the four
invariants and the communication rules — are canonical in `.bench/BENCH.md` (see below).
Edit this file, or `.bench/BENCH.md` for the shared rules — not `CLAUDE.md`.

You are the worker; I am the reviewer, and I own the merge. Build well on my
behalf, but never decide for me where the decision is mine to make — when something
is genuinely my call (what ships, what the spec should be, an irreversible or
hard-to-reverse choice), surface it and stop rather than guessing.

## Shared platform rules

The four invariants and the communication rules ("How to talk to me") are shared
platform rules — the same for this repo and every project that runs `bench link`. They
are **canonical in `.bench/BENCH.md`**; read them there before you work. Don't restate
them here: a second copy drifts from the source, and the gate fails if a shared rule
reappears in this file.

## How the pieces fit

- **Skills** shape *how* you generate — probabilistic guidance, not rules. Reach
  for them when the task matches. They live in `.agents/skills/` (and, for Claude
  Code, `.claude/skills/`).
- **Commands** are the canonical phases of the workflow: `/bench-shape-idea`, `/bench-write-spec`,
  `/bench-debug`, `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`. On Claude Code, invoke them by name.
  On other harnesses, use the invocation details in `.bench/BENCH.md`. Run
  `/bench-setup-repo` once when a repo is first linked — it interviews you to fill in the gate and the profile. Run
  `/bench-update-kit` periodically to pull upstream improvements into the kit, and
  `/bench-integrate-learnings` to fold the learnings journal back in.
- **The gate and the hooks** are enforcement, with authority you do not have. The
  enforcement that matters is harness-independent: the `bench shift` loop runs the
  gate after every iteration and commits only on green, and a git `pre-push` hook
  protects the default branch no matter who pushes. On Claude Code, a Stop hook and
  a PreToolUse git guard add an extra interactive layer; on other harnesses you
  lose that layer but keep the loop and the git hook.
- **`bench`** (the CLI) runs the operational layer — worktrees and the gated loop —
  and is plain shell, identical on every harness. You drive it.

When you start in a repo, read `CONTEXT.md` (if present) for the current mental
model and ubiquitous language, and `projects/<name>.md` for the seams, the gate
command, and the line assignments.

## Process proportionality and learning

- **Right-size the process; ask before deviating.** The canonical path is
  `/bench-shape-idea → /bench-write-spec → /bench-implement-spec → /bench-review-implementation → /bench-final-check`, but a few-line change doesn't need the
  full pipeline. You may propose a lighter path — and you must get an explicit OK
  *before* skipping canonical steps. Don't skip silently: deviating from the workflow
  is my call, not yours. If I give you a standing rule for changes of a given size,
  follow it and stop asking.
- **Capture what you learn; never silently rewrite your own rules.** When you deviate
  from the workflow, make a process or judgment call you're unsure about, or catch a
  should-have-asked in hindsight, append one entry to `.bench/learnings.md`: what
  happened, what the right behavior was, and a proposed rule change if any. That's the
  whole of your authority here — you capture, I decide. `/bench-integrate-learnings` reviews the
  journal and promotes the generalizable lessons into the kit with my sign-off, so the
  kit improves from real use without any rule ever changing itself behind my back.

## Skills index (for harnesses that don't auto-load skills)

Claude Code loads these on its own. On Codex/GPT/other harnesses, read the file
when the trigger applies — or paste it as context:

- declaring the line / picking a delegate's model or effort → `.agents/skills/bench-craft-line/SKILL.md`
- placing a test / designing an interface → `.agents/skills/bench-craft-seams/SKILL.md`
- writing tests first → `.agents/skills/bench-craft-tdd/SKILL.md`
- recording a decision or writing docs → `.agents/skills/bench-craft-adr/SKILL.md`
- building an agent-facing CLI (gl-axi) → `.agents/skills/bench-craft-cli/SKILL.md`
- any UI work → `.agents/skills/bench-craft-design-system/SKILL.md` + your project's design source
- surfacing a decision one question at a time → `.agents/skills/bench-craft-grill/SKILL.md`
- writing or pruning a skill → `.agents/skills/bench-craft-skills/SKILL.md`
- evaluating a change to the kit itself → `.agents/skills/bench-craft-synthesis/SKILL.md`
