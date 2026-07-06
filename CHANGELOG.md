# Bench changelog

The synthesis record. `/bench-update-kit` reads this as its baseline — what Bench
already adopted from upstream, and what it deliberately rejected — so closed
decisions stay closed and each re-synthesis is diffed against a known state. Append
one entry per `/bench-update-kit` run or learnings-sourced promotion (queued by
`/bench-what-next`, built under the synthesis discipline); don't rewrite history.

## Unreleased

- **Learnings promotion (2026-07-06, L1: shared-tree worktree rule).** From the
  2026-07-05 shared-tree-contention entry (drained to the roadmap in a prior
  reconcile). Folded into invariant 1 in `.bench/BENCH.md` — one sentence
  generalizing the existing delegate worktree-isolation clause to your own
  side-work: when `git status` shows another writer's in-flight edits, take
  side-work to a `bench worktree` or wait, so every gate verdict answers for
  exactly one diff. Fold, not a new piece; prose-only, gate green.

- **Learnings run (2026-07-04, scope: journal close-out).** Drained three
  entries. Promoted (already shipped): review-findings persistence — the
  `/bench-review-implementation` pickup artifact at `reviews/<spec-slug>.md`
  landed with conformance anchors and a canary bite proof, resolving the
  chat-only-findings entry. Recommended-and-parked: stale-gate benign/real
  status split — a `bench status` semantics change routed to the roadmap, to
  be shaped in one session with the parked `gate_tree_hash` capture-scratch
  carve-out. Dismissed: spec-without-decision-map deviation — one-off context;
  the entry contract and ask-first rule held, no rule change.

- **Learnings run (2026-07-04, scope: review/spec guidance).** Drained six open
  entries. Promoted: `craft-synthesis` requires a fresh-session dogfood run when
  skill or command triggers changed; `craft-delegate` records that read-only
  review findings are verified and fixed by the invoking session, not a separate
  worktree; `/bench-review-implementation` has an inline-axis fallback when a
  harness forbids unsolicited sub-agents; `/bench-write-spec` now checks external
  format/library divergence and runnable byte/wire compatibility probes; Research
  assets that claim byte or wire compatibility carry their own probe; and the
  Bench profile's hostile-input checklist now includes real CLI, linked by-path
  CLI, hook, and adapter invocation surfaces. Left open: stale-gate
  classification and durable review-findings storage, both still need product
  decisions. Dogfood loop pending the current gate diagnosis; verified with
  consistency greps and targeted tests in this run.

- **Learnings run (2026-07-04, scope: line governance).** Drained six open
  entries and one roadmap item. Promoted: the line declaration now uses an
  iteration cap as the stop condition; `craft-line` keeps venue routing
  right-sized instead of requiring every story to delegate, with inline work
  allowed for tiny slices and atomic diffs when deviations are reported; ordinary
  spec authoring stays mid-tier, with a top-tier exception only for Handoff
  uncertainty plus reviewer approval; delegate charges should use the Handoff
  digest and line-ranged excerpts over whole-file read lists; and `/bench-debug`
  allows direct fix-and-gate for small single-seam fixes. Dismissed: mandatory
  facilitator delegation as too rigid for delegation overhead and atomic diffs.
  Dogfood loop pending the current gate diagnosis; verified with consistency
  greps and targeted tests in this run.

- **Learnings run (2026-07-04, scope: learnings).** Drained two open entries
  (seam-level test batching; skipped conformance test counted as a pass).
  Promoted into `craft-tdd`: the red step now says to stub minimal declarations
  in compiled languages so the one test compiles before confirming a behavioral
  red — batching the file is never the fix — and the oracle section warns that
  a skip-only run still prints `ok`, so read test output, not the summary line.
  Pruned the batching entry's own proposed "same seam" clause as a duplicate of
  prose already in the skill. Dogfood loop per proportionality: prose-only —
  consistency grep of conformance pins plus a green `bench gate`.

- **Learnings run (2026-07-02, scope: learnings).** Drained one open entry.
  Promoted: cached routing for review-axis delegates in `projects/benchkit.md`
  Lines — ~60k tokens each on mid (three axes ≈ 180k), sourced from actuals
  running 2x a ~30k declaration. Project-specific calibration, not a kit rule.
  Dogfood loop waived: prose-only routing note; verified via consistency grep
  and a green gate.

- **Learnings run (2026-07-02, scope: learnings).** Drained one open entry.
  Promoted: quantified oracle promises must name their granularity — when a
  spec's behavior or red-signal promise ranges over a set ("each check"), the
  spec enumerates the set or states per-item vs per-class explicitly
  (`/bench-write-spec` step 4). Sourced from the state-surface build reading
  "each check bites" as per-class and review catching ~9 missing canaries.
  Dogfood loop waived: prose-only spec-authoring guidance; the next real spec
  is the dogfood.

- **Learnings run (2026-07-02, scope: learnings).** Drained two open entries.
  Promoted: a batch approval covers per-spec sign-offs when the reviewer is
  unreachable — build-and-flag, with specs left as post-hoc veto surface and a
  hard stop absent a batch approval (`.bench/BENCH.md` Workflow); venue routing
  in `craft-line` — delegate a spec'd build to `bench shift` on the cheap
  binding when every story's line is cheap and the gate fully observes the
  coverage map, with `/bench-write-spec`'s handoff clauses now pointing at that
  test instead of "mechanical enough". Dogfood loop waived: prose-only guidance
  a shift can't observe; verified via consistency audit and a green gate.

- **Learnings run (2026-07-02, scope: learnings).** Dismissed: unbound-model
  delegation at the reviewer's word — should-have-applied, `craft-line`'s
  resolve-first rule and the line hook already cover it; no kit change. Per
  reviewer direction, the journal convention changed from mark-resolved to
  prune-on-resolve: `.bench/learnings.md` now holds open entries only, verdicts
  live in the CHANGELOG and integration commits (updated the journal header, the
  `bench init` scaffold, `/bench-integrate-learnings`, and `craft-synthesis`).

- **Learnings run (2026-07-02, scope: learnings).** Drained eight open entries from
  `.bench/learnings.md` (seven tagged, one untagged straggler). Promoted: derived
  out-of-scope estimates (`<n> edits, <n> gate runs`) in `/bench-write-spec` step 3
  and its template; a Won't-handle line must leave at least one in-scope caller able
  to exercise the interface (`/bench-write-spec` step 5); delegate done-claims are
  verified against the gate and `git status`, with write-delegations in isolated
  worktrees (`.bench/BENCH.md` invariant 1); explicit staging instead of blind
  `git add -A`, with unexplained working-tree files blocking the commit
  (`/bench-implement-spec`). Dismissed: cached-routing miss and Codex skill-menu leak
  (fixes already shipped), seam-widening (one instance, not a pattern), invisible
  guards (already parked on the roadmap). Parked: the shift loop's own `add -A`
  staging as a roadmap idea.
- **Learnings run (2026-07-01, scope: learnings).** Drained seven open entries from
  `.bench/learnings.md`. Promoted: value-ranked out-of-scope capture in
  `/bench-write-spec`; a standing direct fix-and-gate shortcut for concrete review
  findings in `.bench/BENCH.md`; rename/refactor hygiene in `/bench-implement-spec`;
  and command-currency hardening in the gate for `.agents/**`, non-historical
  `decisions/`, Codex `$bench-*` adapters, and exact historical markers, with
  targeted canaries. Recorded as already promoted: Bench harness invocation guidance
  belongs in `.bench/BENCH.md`, not `AGENTS.md`.
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
