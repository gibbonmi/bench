# Working agreement

This is the canonical instruction file for project-owned content, read by every harness
(Claude Code via `@AGENTS.md` and `@.bench/BENCH.md` imports in `CLAUDE.md`; Codex,
OpenCode, and other AGENTS.md harnesses natively). The shared platform rules — roles,
the four invariants, how the pieces fit, the workflow and its proportionality rules,
the communication rules, and the skills index — are canonical in `.bench/BENCH.md`
(see below). Edit this file for project-owned content, or `.bench/BENCH.md` for the
shared rules — not `CLAUDE.md`.

## Shared platform rules

The platform rules are the same for this repo and every project that runs
`bench link`. They are **canonical in `.bench/BENCH.md`**; read them there before
you work. Don't restate them here: a second copy drifts from the source, and the
gate fails if a shared rule reappears in this file.

## This repo

This repo is the Bench kit itself — the platform files here are the source every
linked repo receives, so kit edits follow the `craft-synthesis` discipline and the
leverage override in `craft-line`. The project profile is `projects/benchkit.md`;
read it for the gate's shape, the tier binding, and cold-session notes.

**Code standard — one source per fact.** Knowledge duplication is a defect: two
derivations of the same fact (an enforcement and its advertisement, a parser and
its count, a fixture harness pasted N times) must collapse to one source, and
review grades diffs against this. Honest repetition of incidental text is fine
where an abstraction would be worse; it's duplicated *knowledge* that drifts.

**Shell conventions for agents in this repo.**

- Use `rg` (ripgrep, installed here) instead of `grep` in interactive Bash calls —
  plain `grep` only inside kit scripts, where the kit's own portability rules require
  POSIX grep. Prefer dedicated read tools over `cat`/`head`/`tail`/`sed`.
- Don't prefix Bash commands with `cd` into the working directory — the Bash tool's
  CWD already persists there, so the `cd` is a no-op and can trigger a needless
  permission prompt. Only `cd` when genuinely moving to a *different* directory.
