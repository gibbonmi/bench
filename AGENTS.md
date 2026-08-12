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
you work. Don't restate them here: a second copy drifts from the source. A shared
rule's literal marker phrase reappearing in this file fails the gate; a paraphrased
restatement slips past that substring check but still violates the one-source rule
and is review's to catch, not the gate's.

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
An independently authored test expectation is not duplicated implementation
knowledge only when its independence is necessary for a named omission or mutation
to make the gate red and that red is recorded and demonstrated. This exception
applies only to the expectation-versus-implementation pair; production policy,
parsers, fixture harnesses, executable registries, and derived counts remain
single-sourced.

**Dependency standard.** A third-party Go dependency follows the precedent set
by the first one: official-org source, MIT-compatible license, build-time-only
footprint (linked into the binary, no runtime service or install-time fetch).
Anything outside that shape is a reviewer decision, not a default.

**Phase-close handoff.** When a Bench phase closes (a drain committed, a spec
staged, a build landed, a review delivered), the closing message must either
emit a copy-paste fresh-session continuation prompt or update
`capture/session-handoff.md` — pinning repository, branch, commit, spec path and
status, decisions that stay closed, and the exact harness-native next command —
so resumption never depends on conversation history.

Rewrite `capture/session-handoff.md` in full rather than appending to it — it is pruned,
not accreted, because a fresh session pays for every line it reads cold. The file
carries its own shape; follow what is already in it. Where the handoff and the
tree disagree, the tree wins: `bench status` reports how many commits have landed
since the handoff was last written, so a stale one is ambient rather than
something the next session has to think to check.

**Shell conventions for agents in this repo.**

- Use `rg` (ripgrep, installed here) instead of `grep` in interactive Bash calls —
  plain `grep` only inside kit scripts, where the kit's own portability rules require
  POSIX grep. Prefer dedicated read tools over `cat`/`head`/`tail`/`sed`.
- Don't prefix Bash commands with `cd` into the working directory — the Bash tool's
  CWD already persists there, so the `cd` is a no-op and can trigger a needless
  permission prompt. Only `cd` when genuinely moving to a *different* directory.
- Wait on a PID or a sentinel file, never by polling a self-matching pattern —
  `pgrep -f` of your own command text matches the command that's doing the polling.
- A destructive or bulk-rewrite script runs plan-before-apply: print the exact
  target list, sample it, then apply.
- A repository-wide sweep uses `rg --hidden` (excluding `.git/`) so dot-directories
  are enumerated.
- Discover Bench verbs non-interactively — `bench commands --brief` or the source —
  never by trying a bare unknown verb; pipe stdin from `/dev/null` where a command
  might prompt.
