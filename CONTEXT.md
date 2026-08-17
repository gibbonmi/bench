# Context — ubiquitous language for benchkit

The terms below have one canonical name each. Use these; don't invent synonyms. A
cold session reads this first to avoid drifting the vocabulary.

## Core terms

- **gate** — the executable oracle at `.bench/gate.sh` (or `$BENCH_GATE`, or an
  auto-detected default). Exit 0 means shippable. It is the *only* thing that can
  call work done. Not "the checks", not "CI", not "the suite" — the gate.
- **oracle** — the authority that decides done-ness. The gate is this repo's oracle.
  The model is never the oracle; it does not grade its own work.
- **canary inventory** — the non-empty set of immutable fixture bindings that
  `bench canary` validates. It proves inventory shape and one accepted binding
  per fixture, not owner execution or a planted red. Not “canary dispatch”, not
  “the sweep”.
- **planted-reason proof** — a direct mutation test materializes one fixture,
  invokes its exact owner, requires the fixture's `EXPECT` diagnostic, restores
  the subject, and requires that diagnostic to disappear. Not inventory
  validation, not an empty-tree collision screen.
- **linked repo** — a project that receives Bench's linked payload and local
  launcher, but not the kit's source-only ordinary test packages. Not “consumer
  repo”, not “the kit”.
- **shift** — one run of the gated loop (`bench shift "<objective>"`): iterate, run
  the gate, commit only on green. The unit of autonomous work.
- **seam** — a stable interface where a test attaches and where you compose rather
  than invent. Listed per repo in `projects/<name>.md`. Not "boundary", not
  "interface layer" — seam.
- **line** — the declared model + effort + rough token cap for a stage, with one
  clause of justification. "Declare the line" = state this before a long run. Not
  "budget" alone — the line is the whole routing decision.
- **worktree** — an isolated Git checkout. `bench shift` leases warm, reusable
  pooled worktrees. Interactive and harness lifecycle commands create exact owned,
  locked assignment worktrees and release them safely.
- **worktree admin entry** — one file or directory git keeps per registered
  worktree under `<git-common-dir>/worktrees/<id>/` (`gitdir`, `HEAD`,
  `commondir`, `locked`, …), read by `git worktree` subcommands. Git-owned state
  Bench never writes. Not "metadata file", not "worktree config" — admin entry.
- **landed assignment** — an assignment worktree whose ledger state is `active`,
  whose branch has landed on the default branch, and whose lease is not live: nobody
  released it. A derived classification, never a ledger state; retired by
  `bench worktree clean --landed`. Not "orphan" (that is age alone), not "stale",
  "idle", "abandoned", or "unreleased" (true of every active row) — landed.
- **invariant** — one of the four non-negotiable rules (canonical in `.bench/BENCH.md`)
  that override convenience. Not "guideline", not "best practice" — invariant.
- **harness** — the agent runtime that reads `AGENTS.md` (Claude Code, Codex,
  OpenCode, …). The kit is harness-agnostic by design.
- **kit** — benchkit itself: the shipped workflow (CLI + working agreement + skills +
  commands + hooks). Not "framework", not "tool" — kit.
- **profile** — the per-repo file at `projects/<name>.md` a cold session reads to
  learn the seams, lines, gate command, and (for UI repos) the design source.
- **skill** — probabilistic guidance that shapes *how* the model generates
  (`.agents/skills/*/SKILL.md`; `.claude/skills/` mirrors them for Claude Code).
  Reached for when the task matches; not a rule.
- **command** — a canonical phase of the workflow (`/bench-shape-idea`, `/bench-write-spec`, `/bench-debug`,
  `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`, plus `/bench-setup-repo`, `/bench-update-kit`, `/bench-what-next`). Not "slash command
  template" in prose — command.
- **roadmap** — the working prioritization document at `ROADMAP.md` (repo root):
  assessed open work in priority order, ending in a `## Recommended sequence`
  that names the next actions. Printed with `bench roadmap`. A row leaves when
  the work ships or a reconcile removes it. Not "icebox", not "backlog" — roadmap.
- **roadmap index / roadmap detail** — schema-4 projections of the **roadmap** and
  its capture evidence. The index inventories every row and capture unit with
  true body sizes but no bodies; detail is a complete body fetched for named rows
  or read from the capture paths the index names. Not a truncated preview.
- **ideas inbox** — the capture-and-forget sink at `capture/IDEAS.md` (repo root):
  out-of-scope ideas parked with `bench idea`, committing to nothing.
  Append-only, no status or lifecycle; drained to zero into the **roadmap** by
  the maintenance phase.
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
- **coverage row** — one row of a spec's acceptance coverage map, tying one
  story to one observable behavior at a seam: story, behavior, seam, why it
  catches the failure, behind an optional leading row ID. Not "acceptance
  criterion", not "test case", not "requirement" — coverage row.
- **acceptance row** — the **coverage row** an implementation ticket claims as
  its own acceptance. Not "the ticket's checklist", not "definition of done" —
  acceptance row.
- **session handoff** — the phase-close continuation artifact at
  `capture/session-handoff.md`, rewritten in full by `bench handoff`. It hands the whole
  repository to the next session and is rewritten, never appended to.
- **pin block** — the header `bench handoff` prints to stdout and writes at the top of the
  session handoff: repository, path, branch and HEAD, tree state, staged spec, gate
  verdict, and the next command. One derivation feeds both sinks, so the printed block and
  the file cannot state different facts. Not "the header", not "the summary" — pin block.
- **ambient dashboard** — what `bench status` prints: the cold-session + on-demand view
  of what needs attention. The feature that renders it is the *ambient-feedback surface*;
  the printed thing is the ambient dashboard. Not "status report", not "summary".
- **the dashboard page** — the standalone HTML artifact `bench dashboard` renders, distinct
  from the **ambient dashboard** (`bench status`): the page is the shareable rendered view,
  the ambient dashboard is the terminal print. Not "the dashboard" unqualified — that name
  ambiguously means both.
- **signal** — one ranked line on the **ambient dashboard** (gate, git, worktree, intent, guards,
  drain, structure, decisions, specs, reviews, roadmap). Shown only when it fires. Not
  "check", not "alert".
- **severity ladder** — the fixed rank order that decides which **signal** leads the
  dashboard and which drop under the five-row budget. Not "priority queue".
- **gate cache** — the durable ready or pending verdict that gate execution binds to
  the closed oracle subject in the Git dir, so read-only consumers can project gate
  state without a cold run. Not "gate log".
- **landing source** — a build-owned Git integration branch identified by its
  frozen base and current source tip. Serial green tickets accumulate there;
  semantic review binds to that pair and `bench worktree land` consumes it. Not
  a mutable review base, not a reconstructed path list.
- **landing destination** — the expected tip of the branch that receives a
  landing source. Movement requires a new composition and gate verdict, not a
  new source review by itself. Not the source's frozen base.
- **prospective landing tree** — the exact Git tree produced by composing a
  landing source onto an expected landing destination, including any authorized
  final transition. It is the whole-project gate subject. Not the source diff,
  not the ambient working tree.
- **frontier round** — one numbered round of `craft-grill`: every question whose
  prerequisites are already settled, asked together with a recommendation, before
  the skill waits and recomputes the next round. Not "one question at a time".
- **disposition** — the exactly-one next-action label a review finding carries:
  `no-op` (the concern is refuted, no repair target remains), `auto-fix` (a
  deterministic rule or exact spec predicate, repairable in already-approved
  scope), or `ask-user` (needs judgment, scope, authority, or an oracle change).
  A repair-routing label, not permission for the read-only review phase to edit.
- **prose budget** — the per-surface line-count ceiling in `projects/benchkit.md`'s
  one mechanically parseable budget table, enforced by a fail-closed conformance
  check. Not "line limit" in prose — prose budget.
- **green marker** — the ref `refs/bench/green/<branch>` the gate advances to the tip
  it graded green, read through one marker module that peels, classifies a dangling
  ref, and answers whether the marker authorizes a tip. Not "green ref", not "the
  marker" unqualified — green marker.
- **record class** — one registered shape of a gate verdict record: its name, exact
  field set, and validator, enumerated only in the record-class registry. Not "verdict
  kind", not "field set" — record class.
- **eligibility verdict** — the one decided answer to "is this worktree Bench-owned and
  safe to remove", carrying its evidence, with refusal precedence declared as ordered
  data; automatic cleanup is a stricter reading of it. Not "cleanup plan" (that is the
  explicit plan/apply fingerprint), not "classification" — eligibility verdict.
- **objective projection** — a per-surface rendering of a shift objective (banner,
  prompt, scratch, predicate argument, durable commit subject) issued by the objective
  module; no surface receives the unprojected text. Not "the objective string" —
  objective projection.

## Avoid

- "CI" / "the build" when you mean **the gate**.
- "task" / "session" when you mean a **shift**.
- "boundary" / "abstraction point" when you mean a **seam**.
- "framework" / "tooling" when you mean the **kit**.
- "icebox" / "backlog" when you mean the **roadmap**.
