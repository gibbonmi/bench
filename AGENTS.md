# Working agreement

This file is the canonical instruction file for project-owned content. Every
harness reads this file. Claude Code reads it through the `@AGENTS.md` and
`@.bench/BENCH.md` imports in `CLAUDE.md`. Codex, OpenCode, and other AGENTS.md
harnesses read the file directly.

The shared platform rules are canonical in `.bench/BENCH.md` (see below).
These rules cover the roles, the four invariants, how the pieces fit, the
workflow and its proportionality rules, the communication rules, and the
skills index. Edit this file for project-owned content. Edit
`.bench/BENCH.md` for the shared rules. Do not edit `CLAUDE.md`.

## Shared platform rules

The platform rules are the same for this repo and for every project that runs
`bench link`. These rules are **canonical in `.bench/BENCH.md`**; read them
there before you start work. Do not restate the rules here, because a second
copy drifts from the source. If a shared rule's exact marker phrase reappears
in this file, the gate fails.

A paraphrase of the same rule can pass the gate's substring check, but a
paraphrase still violates the one-source rule. Review catches this
violation, not the gate.

## This repo

This repo is the Bench kit itself. The platform files here are the source
that every linked repo receives. Kit edits therefore follow the
`craft-synthesis` discipline and the leverage override in `craft-line`. The
project profile is `projects/benchkit.md`. Read it for the gate's shape, the
tier binding, and cold-session notes.

**Code standard — one source per fact.** Knowledge duplication is a defect.
Two derivations of the same fact must collapse into one source. For
example: an enforcement and its advertisement, a parser and its count, or a
fixture harness pasted N times. Review grades diffs against this rule. Honest
repetition of incidental text is fine where an abstraction would be worse,
because duplicated *knowledge*, not incidental text, drifts.

An independently authored test expectation counts as duplicated
implementation knowledge, unless two conditions hold. First, the
expectation's independence must be necessary so that a named omission or
mutation turns the gate red. Second, someone must record and demonstrate
that red. This exception applies only to the expectation-versus-implementation
pair. Production policy, parsers, fixture harnesses, executable registries,
and derived counts remain single-sourced.

**Dependency standard.** A third-party Go dependency must follow the
precedent set by the first dependency: an official-org source, an
MIT-compatible license, and a build-time-only footprint. This footprint means
the dependency links into the binary, with no runtime service and no
install-time fetch. A dependency outside this shape needs a reviewer
decision; it is not a default.

**Phase-close handoff.** A Bench phase closes when a drain commits, a spec
stages, a build lands, or a review delivers. The closing message must then
do one of two things. It must emit a copy-paste fresh-session
continuation prompt, or it must update `capture/session-handoff.md`.

The update must pin the repository, the branch, and the commit. It must also
pin the spec path and status, the decisions that stay closed, and the exact
harness-native next command. This way, resumption never depends on
conversation history.

Rewrite `capture/session-handoff.md` in full; do not append to it. The file
is pruned, not accreted, because a fresh session pays for every line it reads
cold. Follow the shape already in the file. When the handoff and the tree
disagree, the tree wins. `bench status` reports how many commits have landed
since the last handoff update. So a stale handoff is only ambient
information; the next session does not need to confirm it.

**Shell conventions for agents in this repo.**

- Use `rg` (ripgrep, installed here) instead of `grep` in interactive Bash
  calls. Use plain `grep` only inside kit scripts, where the kit's own
  portability rules require POSIX grep. Prefer dedicated read tools over
  `cat`, `head`, `tail`, and `sed`.
- Do not prefix Bash commands with `cd` into the working directory. The Bash
  tool's CWD already persists there, so the `cd` is a no-op. It can also
  trigger a needless permission prompt. Use `cd` only when you move to a
  genuinely *different* directory.
- Wait on a PID or a sentinel file, never by polling a self-matching pattern.
  `pgrep -f` of your own command text matches the command that is doing the
  polling.
- Run a preserve step (a copy aside) as its own command, and make sure the copy
  exists before any discard step. If a hook refuses a chain, no step of the
  chain runs.
- A destructive or bulk-rewrite script runs plan-before-apply: print the exact
  target list, sample it, then apply.
- A repository-wide sweep uses `rg --hidden` (excluding `.git/`) so
  dot-directories are enumerated.
- Discover Bench verbs non-interactively. Use `bench help` for the inventory
  and `bench <verb> --help` for a grammar. Note that `bench commands --brief`
  is a three-verb liveness probe, not an inventory. Never discover a verb by
  trying a bare unknown verb.
