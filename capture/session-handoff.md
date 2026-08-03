# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2276b27`, clean tree, 5 unpushed commits
Spec: `specs/ft164-ticket-contracts/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `2a8f9f0` — stale, work tree `ada7d85`

## State

- **FT164 `--full` build: implementation phase, in flight.** Lifecycle
  `bench spec build` run active for `ft164-ticket-contracts`, subject `2276b27`.
  Seven executable-contract tickets committed under
  `specs/ft164-ticket-contracts/tickets/` (opus-authored, opus-reviewed, one
  repair round applied). Two assignments out on the frontier: request
  `ft164-export-parser` (ticket `export-the-ticket-parser.md`) and request
  `ft164-teach-template` (ticket
  `teach-the-executable-contract-template-and-examples.md`); the remaining five
  tickets are dependency-blocked (chain teach→open→add→charge→point; land
  blocked by export+teach). `bench spec build status ft164-ticket-contracts`
  is authoritative for assignment state.
- **Line (reviewer override, 2026-08-02):** opus/high for all write delegates
  and all reviews — opus followability is the build's purpose — overriding the
  spec's fable/high on stories 1–6; fable coordinates and verifies. Research
  fan-out (6 opus delegates) complete; key facts encoded in the tickets'
  Assumptions lines.
- **Coordinator decisions flagged for reviewer veto:** headings inside the
  craft-tickets fenced template/examples demote to `###` (the anchor harness's
  section resolver is fence-blind); the test-harness ceiling lands beside step
  2 of `## Draft the breakdown` instead of the unscopable preamble.
- **Reviewer-directed side work, not started here:** a Codex hotfix prompt was
  handed to the reviewer to relax `bench spec build start`'s exact-green
  precondition (accept reduced/partial verdicts whose inherited evidence
  composes to whole-tree green) and fix its remediation text. It edits
  `internal/specbuild` + `internal/gate` — must not run while the
  `export-the-ticket-parser` assignment is open in the same package. FT173's
  row now records the contextual-disclosure finding and the reviewer's
  per-command acceptance conditions.
- One open learnings entry (gate-subject mutation miss). Nothing pushed; push
  is the reviewer's call.

## Next command

`/bench-implement-spec --full specs/ft164-ticket-contracts/spec.md`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
