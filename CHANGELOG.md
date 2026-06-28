# Bench changelog

The synthesis record. `/resynthesize` reads this as its baseline — what Bench
already adopted from upstream, and what it deliberately rejected — so closed
decisions stay closed and each re-synthesis is diffed against a known state. Append
one entry per `/resynthesize` run; don't rewrite history.

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
