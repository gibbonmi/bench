# Context — ubiquitous language for benchkit

The terms below have one canonical name each. Use these terms; do not invent
synonyms. A cold session reads this file first so the vocabulary does not drift.

## Core terms

- **gate** — the executable oracle at `.bench/gate.sh` (or `$BENCH_GATE`, or an
  auto-detected default). Exit 0 means the work is shippable. The gate is the
  *only* thing that can call work done. Not "the checks", not "CI", not "the
  suite" — the gate.
- **oracle** — the authority that decides done-ness. The gate is this repo's oracle.
  The model is never the oracle; the model does not grade its own work.
- **run binary** — the exact source-bound executable the gate authorizes and runs;
  `internal/runbinary` implements it as its `Selection`. Not "private binary", not
  "selected executable" — run binary.
- **canary inventory** — the non-empty set of immutable fixture bindings that
  `bench canary` validates. It proves the inventory shape and one accepted binding
  per fixture, not owner execution or a planted red. Not “canary dispatch”, not
  “the sweep”.
- **planted-reason proof** — a direct mutation test materializes one fixture
  and invokes its exact owner. It requires the fixture's `EXPECT` diagnostic,
  restores the subject, and requires that diagnostic to disappear. Not
  inventory validation, not an empty-tree collision screen.
- **linked repo** — a project that receives Bench's linked payload and local
  launcher, but not the kit's source-only ordinary test packages. Not “consumer
  repo”, not “the kit”.
- **shift** — one run of the gated loop (`bench shift "<objective>"`): the agent
  iterates, runs the gate, and commits only on green. The shift is the unit of
  autonomous work.
- **seam** — a stable interface where a test attaches and where you compose rather
  than invent. Each repo lists its seams in `projects/<name>.md`. Not "boundary",
  not "interface layer" — seam.
- **fresh test run** — one Go test run with successful test-result reuse disabled
  while the ordinary build and module caches remain available. For this kit, the
  whole-tree form is `go test -count=1 ./...`. Not "cold test" or "clean-cache
  run" — those also discard compilation or dependency state and measure a
  different workload.
- **real Git journey** — a behavioral test that materializes an actual Git
  repository and crosses a public command seam. It proves Git composition and
  process behavior, not every equivalence partition behind the command. Not
  "integration test" alone — that label does not state whether real Git runs.
- **line** — the declared model, effort, and rough token cap for a stage, with one
  clause of justification. "Declare the line" means state this before a long run.
  Not "budget" alone — the line is the whole routing decision.
- **worktree** — an isolated Git checkout. `bench shift` leases warm, reusable
  pooled worktrees. Interactive and harness lifecycle commands create exact owned,
  locked assignment worktrees and release them safely.
- **worktree admin entry** — one file or directory git keeps per registered
  worktree under `<git-common-dir>/worktrees/<id>/` (`gitdir`, `HEAD`,
  `commondir`, `locked`, …), read by `git worktree` subcommands. Git-owned state
  that Bench never writes. Not "metadata file", not "worktree config" — admin
  entry.
- **identity component** — one named check a lifecycle verb runs against an
  assignment's own records: the request token, the assignment state, the assignment
  path, the owner marker, the worktree registration, or the Bench lock. One registry
  declares the set and its precedence. A refusal names exactly one component, and the
  component that carries a sanctioned repair names the exact command. Not "identity
  check", not "mismatch reason", not "refusal cause" — identity component.
- **landed assignment** — an assignment worktree whose ledger state is `active`
  and whose branch has landed on the default branch. Its lease is not live:
  nobody released it. Bench derives this classification; it is never a ledger state.
  `bench worktree clean --landed` retires it. Not "orphan" (that names age
  alone), not "stale", "idle", "abandoned", or "unreleased" (true of every
  active row) — landed.
- **invariant** — one of the four non-negotiable rules (canonical in `.bench/BENCH.md`)
  that override convenience. Not "guideline", not "best practice" — invariant.
- **harness** — the agent runtime that reads `AGENTS.md` (Claude Code, Codex,
  OpenCode, …). The kit is harness-agnostic by design.
- **environment closure** — the harness process carries the toolchain effects
  implied by its inherited initialization state. Not "loaded environment" or
  "working PATH" — both names hide partial propagation.
- **kit** — benchkit itself: the shipped workflow (CLI + working agreement + skills +
  commands + hooks). Not "framework", not "tool" — kit.
- **profile** — the per-repo file at `projects/<name>.md` a cold session reads to
  learn the seams, lines, gate command, and (for UI repos) the design source.
- **skill** — probabilistic guidance that shapes *how* the model generates
  (`.agents/skills/*/SKILL.md`; `.claude/skills/` mirrors them for Claude Code).
  A session reaches for a skill when the task matches; a skill is not a rule.
- **command** — a canonical phase of the workflow (`/bench-shape-idea`, `/bench-write-spec`, `/bench-debug`,
  `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`, plus `/bench-setup-repo`, `/bench-update-kit`, `/bench-drain`). Not "slash command
  template" in prose — command.
- **roadmap** — the working prioritization document, split between `ROADMAP.md`
  (repo root) and `roadmap/`: assessed open work in priority order, ending in a
  `## Recommended sequence` that names the next actions. `ROADMAP.md` is the
  index — one physical heading line per row, no bodies. `roadmap/FT<n>.md`
  owns that row's body, its `Occurrence:` ledger, and its `Sources:` line.
  `bench roadmap` prints it. A row leaves — index line and detail file together — when the
  work ships or a reconcile removes it. Not "icebox", not "backlog" — roadmap.
- **roadmap index / roadmap detail** — schema-4 projections of the **roadmap** and
  its capture evidence. The index inventories every row and capture unit with
  true body sizes but no bodies. Detail is a complete body, read for a named row
  from its `roadmap/FT<n>.md` owner or from the capture paths the index names.
  Detail is not a truncated preview.
- **ideas inbox** — the capture-and-forget sink at `capture/IDEAS.md` (repo root):
  out-of-scope ideas parked with `bench idea`, committing to nothing.
  It is append-only, carries no status or lifecycle, and the maintenance phase
  drains it to zero into the **roadmap**.
- **park** — to capture an idea in the **ideas inbox** without committing to it
  (`bench idea "<text>"`). A *parked idea* graduates into committed work only when
  `/bench-shape-idea` pulls it into a decision map. Not "stash", not "file" — park.
- **decision map** — a situational working map for a multi-session unresolved
  decision tree at `decisions/<topic>.md`. Its **decision tickets** record
  reviewer choices, constraints, exclusions, research objects, and bounded
  discretion; a ready map is compiled beside its spec. It does not choose the
  spec's engineering seams or tests. Not "PRD", not "design doc" — decision map.
- **decision ticket** — one reviewer choice or evidence-producing question in
  a decision map, linked to other decision tickets by `Blocked by`. Distinct
  from an **implementation ticket**, the independently-green build unit under
  `specs/<slug>/tickets/`.
- **coverage row** — one row of a spec's acceptance coverage map. It ties one
  story to one observable behavior at a seam, behind an optional leading row
  ID. It carries four parts: story, behavior, seam, why it catches the
  failure. Not "acceptance criterion", not "test case", not "requirement" —
  coverage row.
- **acceptance row** — the **coverage row** an implementation ticket claims as
  its own acceptance. Not "the ticket's checklist", not "definition of done" —
  acceptance row.
- **session handoff** — the phase-close continuation artifact at
  `capture/session-handoff.md`, which `bench handoff` rewrites in full. It hands the whole
  repository to the next session; the next session never appends to it.
- **pin block** — the header `bench handoff` prints to stdout and writes at the top of the
  session handoff: repository, path, branch and HEAD, tree state, staged spec, gate
  verdict, and the next command. One derivation feeds both sinks, so the printed block and
  the file cannot state different facts. Not "the header", not "the summary" — pin block.
- **ambient dashboard** — what `bench status` prints: the cold-session + on-demand view
  of what needs attention. The feature that renders it is the *ambient-feedback surface*;
  the printed thing is the ambient dashboard. Not "status report", not "summary".
- **the dashboard page** — the standalone HTML artifact `bench dashboard` renders, distinct
  from the **ambient dashboard** (`bench status`). The page is the shareable rendered view;
  the ambient dashboard is the terminal print. Not "the dashboard" unqualified — that name
  ambiguously means both.
- **signal** — one ranked line on the **ambient dashboard** (setup, gate, git, worktree, intent, guards,
  drain, structure, decisions, specs, reviews, roadmap). The dashboard shows a signal only
  when it fires. Not "check", not "alert".
- **severity ladder** — the fixed rank order that decides which **signal** leads the
  dashboard and which drop under the five-row budget. Not "priority queue".
- **gate cache** — the durable ready or pending verdict that gate execution binds to
  the closed oracle subject in the Git dir. Read-only consumers can then project gate
  state without a cold run. Not "gate log".
- **landing source** — a build-owned Git integration branch identified by its
  frozen base and current source tip. Serial green tickets accumulate there;
  semantic review binds to that pair and `bench worktree land` consumes it. Not
  a mutable review base, not a reconstructed path list.
- **landing destination** — the expected tip of the branch that receives a
  landing source. To move it, compose the tree again and get a new gate verdict; a
  new source review alone is not enough. Not the source's frozen base.
- **prospective landing tree** — the exact Git tree produced by composing a
  landing source onto an expected landing destination, including any authorized
  final transition. It is the whole-project gate subject. Not the source diff,
  not the ambient working tree.
- **frontier round** — one numbered round of `craft-grill`: every question whose
  prerequisites are already settled, asked together with a recommendation. The
  skill then waits and recomputes the next round. Not "one question at a time".
- **disposition** — the exactly-one next-action label a review finding carries.
  `no-op` means the concern is refuted and no repair target remains. `auto-fix`
  means a deterministic rule or exact spec predicate, repairable in
  already-approved scope. `ask-user` means it needs judgment, scope, authority,
  or an oracle change. It is a repair-routing label, not permission for the
  read-only review phase to edit.
- **prose budget** — the per-surface line-count ceiling in `projects/benchkit.md`'s
  one mechanically parseable budget table, enforced by a fail-closed conformance
  check. Not "line limit" in prose — prose budget.
- **green marker** — the ref `refs/bench/green/<branch>` the gate advances to the tip
  it graded green. One marker module reads it: it peels, classifies a dangling
  ref, and answers whether the marker authorizes a tip. Not "green ref", not "the
  marker" unqualified — green marker.
- **record class** — one registered shape of a gate verdict record: its name, exact
  field set, and validator, enumerated only in the record-class registry. Not "verdict
  kind", not "field set" — record class.
- **eligibility verdict** — the one decided answer to "is this worktree Bench-owned and
  safe to remove", carrying its evidence, with refusal precedence declared as ordered
  data. Automatic cleanup is a stricter reading of it. Not "cleanup plan" (that is the
  explicit plan/apply fingerprint), not "classification" — eligibility verdict.
- **objective projection** — a per-surface rendering of a shift objective (banner,
  prompt, scratch, predicate argument, durable commit subject) issued by the objective
  module. No surface receives the unprojected text. Not "the objective string" —
  objective projection.
- **prose mechanics check** — the `prose-mechanics` conformance check, which grades the
  two ASD-STE100 rules a program can measure: sentence length and paragraph length. Not
  "STE lint", not "prose lint", not "grammar check" — prose mechanics check.
- **prose exclusion row** — one line of `.bench/prose-exclusions`: a path the prose
  mechanics check does not grade, and a one-clause reason. The reviewer owns that file.
  Not "allowlist", not "skip list" — prose exclusion row.
- **always-loaded core** — `.bench/BENCH.md`, which holds the six rule families every
  session loads; the mechanics live in `.bench/BENCH-reference.md`. Not "the guide"
  unqualified, and not "progressive loading" (the split is progressive disclosure) —
  always-loaded core.

## Avoid

- "CI" / "the build" when you mean **the gate**.
- "task" / "session" when you mean a **shift**.
- "boundary" / "abstraction point" when you mean a **seam**.
- "framework" / "tooling" when you mean the **kit**.
- "icebox" / "backlog" when you mean the **roadmap**.
- "loaded environment" / "working PATH" when you mean **environment closure**.
