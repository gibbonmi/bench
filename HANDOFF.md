# Bench handoff

Pick-up doc for a fresh session or harness. Package consumers should read
`README.md`, `AGENTS.md`, `.bench/BENCH.md`, and this file first. Kit development in
this repository also uses `CONTEXT.md` and `projects/benchkit.md`; those are
repo-only development context, not npm package contents.

## Current shape

Bench is the npm package `benchkit`: a harness-agnostic kit made of shell,
markdown, JSON, portable commands, portable skills, shared hooks, and adapter files.
The product is portability: Claude Code, Codex, and AGENTS.md harnesses should land on
the same working agreement and the same operational substrate.

The shared invariants and communication rules are canonical in `.bench/BENCH.md`.
`AGENTS.md` is the project-owned working agreement and points there. `CLAUDE.md` is an
adapter import shim, not a source of truth.

## Shipped surfaces

- Commands: `/bench-setup-repo`, `/bench-shape-idea`, `/bench-write-spec`,
  `/bench-debug`, `/bench-implement-spec`, `/bench-review-implementation`,
  `/bench-final-check`, `/bench-update-kit`, `/bench-integrate-learnings`
- Skills: `craft-seams`, `craft-tdd`, `craft-adr`, `craft-cli`,
  `craft-design-system`, `craft-skills`, `craft-grill`, `craft-synthesis`,
  `craft-line`, `craft-review`, `craft-delegate`, `craft-gate`
- CLI: `bench link`, `bench init`, `bench models`, `bench structure`, `bench idea`,
  `bench roadmap`, `bench status`, `bench learnings`, `bench maps`, `bench guards`,
  `bench diff`, `bench coverage`, `bench doctor`, `bench gate`, `bench worktree`,
  `bench shift`, `bench version`; plumbing the hooks and shell adapters call:
  `bench tree-hash`, `bench worktree-pool`, `bench worktree-lease-file`,
  `bench guard-git`
- Hooks: shared `.bench/hooks/stop.sh`, `.bench/hooks/session-start.sh`, and
  `.bench/hooks/block-dangerous-git.sh`, with Claude and Codex adapters
- Profiles: linked repos own their profile. The kit-development profiles under
  `projects/` are repo-only development context, not npm package contents.

## How to continue

Use `bench gate` as the oracle. For repo setup, run `bench link`, `bench init`, then
`/bench-setup-repo`. For ordinary kit work, use the canonical workflow in
`.bench/BENCH.md` unless the reviewer explicitly approves a lighter path. For parked
ideas, use `bench idea "<text>"` and promote only through `/bench-shape-idea`.

The gate is the load-bearing safety net. It checks shell/JSON validity, package
contents, link behavior, roadmap/status contracts, shift/worktree contracts, command
and skill conformance, command-name currency in living docs, and canary fixtures that
prove the checks still bite.

## Maintenance

`/bench-update-kit` compares the kit against upstream sources and records adopted
changes in `CHANGELOG.md`. `/bench-integrate-learnings` drains `.bench/learnings.md`
with reviewer sign-off. Both use `craft-synthesis`: respect closed decisions, assess
each candidate, run legibility/consistency/dogfood loops, and propose rather than
merge.
