# Bench changelog

The synthesis record. `/bench-update-kit` reads this as its baseline — what Bench
already adopted from upstream, and what it deliberately rejected — so closed
decisions stay closed and each re-synthesis is diffed against a known state. Append
one entry per `/bench-update-kit` or `/bench-integrate-learnings` run; don't rewrite history.

## Unreleased

- **Safe link dogfood slice (2026-06-28).** Made `bench link` copy by default,
  preserve project-owned `AGENTS.md` content through a managed Bench block, fail on
  same-name project-owned skills/commands/hooks, record installed assets in a link
  manifest, install portable `.agents/` content with Claude/Codex adapters, and gate
  the contract through throwaway-repo link checks plus npm package inspection. Added
  `.claude/README.md` so users can see that Claude paths are adapters to `.agents/`
  and shared `.bench/hooks/`.
- **Learnings run (2026-06-28, scope: learnings).** Drained two open entries from
  `.bench/learnings.md`. Promoted: a maintainer rule into HANDOFF — any `.bench/*`
  file the kit's prose references must be scaffolded by `bench init` and locked by a
  behavioral gate check (the `learnings.md` bug; executable fix already in `724bf8c`).
  Dismissed: "gate green claimed without a run" — already governed by invariant 1 and
  `/bench-final-check`; a pre-commit gate run would be a fourth check surface. Skipped:
  generalizing gate check 1d to every `.bench/*` file — speculative, only two exist.
- **Communication made first-class + portable; four learnings drained (2026-06-30,
  scope: learnings).** Sharpened AGENTS.md "How to talk to me" — clarity over density,
  tables/lists encouraged, the finding-vs-context tension scoped (cut the derivation,
  keep the one-clause why and enough context to resume cold), and "recommend, don't
  offer a blind menu." Shipped it to consumers via a new Communication section in
  `.bench/BENCH.md`. Promoted from `.bench/learnings.md`: recommend-at-every-question-
  and-hand-off; approval-table-before-build (`/bench-write-spec` exit); scan-for-unwritten-answers-
  before-closing-a-map (`craft-grill` + `/bench-shape-idea` exits). Dismissed: the persistent
  task-list progress tracker — `TaskCreate/TaskUpdate` are Claude-Code-only, so
  mandating them in harness-shared rules leaks one harness into the core.

## 0.2.0 — app-agnostic, npx-distributable, self-maintaining

- Made the kit app-agnostic: removed the hockey/coach/puck framing and all
  Regroup-specific content from core files (AGENTS.md, every skill and command, the
  CLI). Project-specific rules now live only in `projects/<name>.md`.
- Packaged for `npx` (`package.json`, dual `bench`/`benchkit` bins). `bench link`
  detects an ephemeral npx cache and copies instead of symlinking.
- Generalized the `design-system` skill and made the Claude-Design ↔ other-harness
  transition an explicit property: consumption reads committed artifacts, so the
  authoring tool is never a workflow dependency.
- Added `/setup` (interactive per-repo configuration, mirrors
  `setup-matt-pocock-skills`) and a Pocock-structure migration path.
- Added `/resynthesize` (this maintenance command) with three quality loops.

## 0.1.0 — initial synthesis

Baseline of what Bench incorporates from upstream. Adopted, with provenance in
`README.md`:

- **From Pocock:** `/start-ideation` (decision-mapping), `/spec` (to-prd), `/build` (implement),
  `/prep-shift` (review), `/fix-bug` (diagnosing-bugs), and the `seams` (codebase-design
  + design-it-twice), `tdd-at-seams` (tdd), `adr` (domain-modeling ADR format),
  `grill` (grilling), `writing-great-skills` skills, plus the `block-dangerous-git`
  hook (git-guardrails).
- **From kunchenguid:** the `axi` skill (AXI spec), the `.bench/gate.sh` + Stop-hook
  external-oracle pattern (no-mistakes), the `bench shift` gated loop with
  notes-between-iterations (gnhf), and `bench worktree` (treehouse).
- **Deliberately rejected:** firstmate's fleet orchestrator (overkill for solo work),
  the strict iron-law unattended TDD run (the self-grading failure it produces),
  lavish-axi and the tool-specific AXI binaries (build to the spec instead). These
  are closed unless a future `/resynthesize` finds a material upstream change.
- **Bench-native (neither repo):** the declared line (model + effort), stateless-reader
  docs, the design system as visual oracle, and the gate-as-oracle / never-self-grade
  invariant that binds it together.
