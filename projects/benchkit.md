# Project: benchkit

The Bench kit itself — a harness-agnostic agent-development workflow shipped as the
npm package `benchkit`. It is not an application: it is shell + markdown + JSON that
other repos consume. The deliverable is the `bench` CLI (`bin/bench.sh`), the working
agreement (`AGENTS.md`), the portable `.agents/` skills and commands, and harness
adapters that call shared `.bench/hooks/` scripts. Because the artifacts are plain
files, the kit must work identically under Claude Code, Codex, and any other
AGENTS.md harness — that portability is the product.

## Working branch

`main`. (The commit-on-green policy and the default-branch guard are
canonical in `/bench-final-check`; this line is only the binding.)

## Seams (test here; everything else is free to change)

- **The gate contract** (`.bench/gate.sh` / `bench gate`). The oracle surface. Everything
  in Bench routes through its exit code: 0 = shippable, non-zero = not done. The
  highest seam — if the gate is weak, the whole system is weak. Test by feeding it a
  conformant tree (green) and a broken one (red); never by trusting a reading of the
  diff.
- **The `bench` CLI subcommands** (`gate`, `worktree`, `shift`, `init`, `link`,
  `models`, `structure`, `idea`, `roadmap`, `status`, `learnings`, `maps`,
  `guards`, `diff`, `coverage`). The operational shell surface.
  Stable command names and exit codes are the contract; the implementation behind each is
  free to change. Keep gate resolution (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) in
  one place.
- **The AXI query surface** (`bench learnings`, `bench maps`, `bench guards`,
  `bench diff`, `bench coverage`, and the shared flat-table TOON emitter behind
  them). The agent-facing read-only
  surface, and the AXI-conformant half of the hybrid output contract: TOON stdout,
  definitive empty states, structured errors on stdout, exit 0/1/2. Gate-tested by
  the AXI contract fragments. Guard manifests come from each guard script's own
  `--describe` (generated from the rules it enforces — never a registry), and
  `bench guards --brief` is the surface the SessionStart hook injects.
  `bench diff` is the single source of review-base truth: the shift loop records
  the pre-shift HEAD in `branch.<name>.benchBase`, `diff` resolves that key first
  and merge-base with the default branch as fallback. `bench coverage --check` is
  the one parser for the acceptance-coverage-map convention; the gate's docs
  fragment consumes it instead of carrying its own.
- **The ambient dashboard** (`bench status`). The single deterministic renderer the
  SessionStart hook and the user both call: it ranks the signals that fire on a fixed
  severity ladder and leads with the next action. Reads gate state from the **gate cache**
  (`<git-dir>/bench-last-gate`, written by the Stop hook) — never a cold gate run. The
  contract (gate-tested): show-only-on-signal, a five-row budget, a stale-green that is
  not a clean bill, and one combined capture-drain row (parked ideas + open learnings)
  pointing at `/bench-what-next`. A stale gate softens to `capture-only drift` /
  `re-run when convenient` only when every changed path is in the fixed, exact
  allowlist (`ROADMAP.md`, `IDEAS.md`, `.bench-notes.md` — no directory, suffix, or
  markdown-class matching; expanding it is a new decision); any mixed or untrusted
  diff fails closed to the strong stale row.
- **The capture inbox and working roadmap** (`bench idea` → `IDEAS.md`;
  `bench roadmap` → `ROADMAP.md`). Capture-and-forget: park an out-of-scope idea,
  commit to nothing; ideas graduate only through a `/bench-what-next` drain into the
  working roadmap. The contract (gate-tested in a throwaway repo): `idea` appends one
  dated line and creates the inbox; a no-arg `idea` errors without appending;
  `roadmap` prints the working document plus drain status, or its
  `## Recommended sequence` when nothing needs draining. `IDEAS.md` and `ROADMAP.md`
  are per-consumer content — never in the kit's `package.json` `files[]`.
- **The kit content surface** (`.agents/skills/*/SKILL.md`, `.agents/commands/*.md`).
  Portable harness-facing content. The contract is structural: every skill carries YAML
  frontmatter (name + description) and follows progressive disclosure; every command
  is a phase the index names. The `.bench/BENCH.md` skills index is generated
  from each skill's `index:` frontmatter (`.bench/skills-index.sh --write`);
  craft skills' visible names use `craft-*` so `$bench` menus show
  only human-run phase adapters. Codex command-adapter skills are derived from
  `.agents/commands/` and documented in `.bench/BENCH.md`. The gate's conformance
  layer enforces those contracts so disk, docs, and adapters do not drift.
  `.claude/` paths are adapters, not a second source of truth.
- **The safe link contract** (`bench link`). The adoption surface. It must preserve
  project-owned `AGENTS.md` text, install only a managed Bench block plus Bench-owned
  assets, fail on same-named project-owned skills/commands/hooks, and be idempotent
  through the link manifest. It also installs the `.bench/bin/` local CLI set the
  shared hooks use when no global `bench` command is on PATH.
- **AGENTS.md** — the canonical working agreement for project-owned content. `CLAUDE.md`
  imports it (and `.bench/BENCH.md`); never duplicate content there. The four invariants
  and the communication rules are canonical in `.bench/BENCH.md`; the craft skill and
  command indexes live here. The gate checks those indexes, checks command-adapter skills
  against `.bench/BENCH.md`, and checks that the shared rules are not copied back into
  AGENTS.md.

No UI. There is **no design source** for this repo.

## Hostile-input checklist (shell CLI)

The edge classes `/bench-write-spec`'s edge inventory walks for this domain —
the hostile inputs shell CLIs actually meet. Walk every class before locking a
coverage map; a class skipped here returns as a regression.

- paths and directory names containing spaces or glob characters
- hand-edited files whose last line lacks a trailing newline
- absent file vs present-but-empty file (distinct behaviors, both asserted)
- unquoted multi-word arguments (`$*` vs `$1`)
- required tool missing from PATH (no global `bench`, no `readlink -f`)
- invocation through a symlink rather than the real path
- invocation through every shipped surface: real kit CLI, linked-repo by-path
  CLI, hooks, and adapters must all reach the same routed implementation
- interrupt (SIGINT) mid-loop: leftover scratch state, leases, worktrees
- re-run idempotency: relink, reused worktree, second `init`
- cwd deeper than the repo root when the command assumes root

## Gate (`.bench/gate.sh`)

```
.bench/gate.sh
```

The oracle for a kit, in layers (all green today). `.bench/gate.sh` is a thin
exec into the `gate-phases` plumbing subcommand, which runs the layers below as
four concurrent phases in outer mode (`[phase]`-prefixed output, per-phase
verdicts, run-all-and-aggregate) and sequentially, unprefixed, sweep-skipped in
inner mode (`BENCH_CANARY_INNER=1`) — canary EXPECT matching is substring-based
against inner-gate output, so the inner byte-shape is load-bearing:

1. **Go root conformance** — the conformance phase runs
   `go test -count=1 ./internal/conformance -run '^TestRootConformance$'` with
   `BENCH_CONFORMANCE_ROOT` set to the tree under grade. That suite owns parse and
   validity checks, JSON validity, executable git modes, package-file resolution,
   npm dry-run package shape, generated skills-index equality, Codex adapter
   metadata, Claude skill mirroring, shared-rule single-sourcing, stale command
   references, token-diet placement, workflow anchors, line-routing enforcement,
   compiled-core build/vet/test/cross-compile checks, release-workflow structure,
   guard `--describe` manifests, the profile's hostile-input checklist anchor, and
   acceptance-coverage map validation.
2. **shellcheck** — best-effort, runs only when installed (`-S warning`). Not a hard
   dependency; upgrades the shell lint automatically once present.
3. **Safe-link behavior** — the gate runs `bench link` against throwaway repos to
   prove fresh installs, existing `AGENTS.md` preservation, relink idempotence,
   same-name conflicts, modified-managed file protection, Codex/Claude hook adapters,
   shared hooks, and default copy mode.
4. **Runtime and behavior contracts** — the remaining shell fragments exercise
   version routing, platform-package generation, runtime hooks, shift/worktree
   behavior, doctor/postinstall/status/roadmap behavior, and AXI query-surface
   behavior through throwaway fixture repos.
5. **AXI query-surface contracts** — the query subcommands' hybrid-contract half:
   TOON-shaped stdout, definitive empty states, structured stdout errors, honest
   exit codes, and each guard's aggregation behavior, exercised in throwaway
   fixture repos.
6. **Canary (meta)** — the gate runs itself against deliberately-broken fixtures in
   `tests/canary/` and asserts each goes red with its targeted error substring. Proves
   the checks above still *bite*: a check rotted into an always-pass fails here. This
   is the gate guarding the gate. Fixtures hide dot-dirs behind a `dot-` prefix (e.g.
   `dot-claude`) so the harness doesn't load fixture skills as real ones; the canary
   restores them at run time. The tripwire decision is recorded in
   `docs/adr/0001-working-tree-gate-tripwire.md`.

The gate file lives outside `package.json` `files[]`, so it never ships to consumers.

## Lines (model + effort routing)

The routing rubric — the three-signal decision table and the escalation ladder —
is the `craft-line` skill. This section holds what is project-specific: the
binding, the cached routings, and the escalation policy.

**Tier → model** (this harness; refresh with `bench models`; set 2026-07-01):
top = Fable 5 (`claude-fable-5`) · mid = Opus 4.8 (`claude-opus-4-8`) · cheap =
Sonnet 5 (`claude-sonnet-5`). Machine-readable source: `.bench/lines.env`,
read by the Agent-tool hook and the shift adapters — keep it in sync with this
paragraph. Haiku 4.5 leaves the rotation (no `effort` support, and a fourth
tier adds a distinction the routing signals can't reliably make). Caveat for
Claude Code delegation: the Agent tool addresses models by alias only, so
`lines.env` also declares which aliases bind (`BENCH_ALIAS_TOP=fable`,
`BENCH_ALIAS_MID=opus`, `BENCH_ALIAS_CHEAP=sonnet` — bare `sonnet` resolves to
Sonnet 5, the cheap tier). Cheap-rated in-session work runs inline or bumps
to mid (declared); headless shift runs target `claude-sonnet-5` via
`BENCH_MODEL` through the adapter (adapters take exact ids, not aliases).

**Escalation policy:** no standing top-tier opt-out — any bump to Fable 5
pauses and asks the reviewer (the ladder is in `craft-line`). Tier moves still
get declared — no silent escalation.

- **Skill / command / doc authoring** → **top model, high effort**. This is the
  leverage override in `craft-line`: guidance prose compounds through every
  session that loads it while the edit costs few tokens. The `craft-skills` and
  `craft-adr` skills apply. Spend here.
- **Spec authoring** → **mid model, fresh session by default**. Every spec is
  compiled from a closed map's Handoff; `/bench-write-spec` owns the venue
  mechanic. Top + high is allowed only when the Handoff carries uncertainty flags
  and the reviewer approves the escalation. Distinct from the doc-authoring
  leverage override above: that spends the top tier on the kit's guidance prose; a
  normal spec is decided content transcribed off a Handoff.
- **`bench` CLI shell plumbing** → cheap model, low–medium effort at the known seam.
  Mechanical once the gate-resolution and worktree-pool shapes exist.
- **Gate / conformance logic** → mid effort. Correctness of the oracle matters more
  than speed — a wrong gate is the worst class of bug in a kit whose whole premise is
  "the gate is the oracle."
- **Review-axis delegate** (`/bench-review-implementation`, one per axis) → mid
  model, medium effort, **~1 iteration each** (three axes can run in parallel).
  Read-heavy: each takes the full diff plus standards docs and runs verification
  commands.

## Notes for cold sessions

- Read `AGENTS.md` first — the working agreement. The four invariants and the
  communication rules are canonical in `.bench/BENCH.md` (AGENTS.md points there); read
  that too. `CLAUDE.md` imports both; edit `AGENTS.md` or `.bench/BENCH.md`, not it.
- `CONTEXT.md` pins the ubiquitous language (gate, oracle, shift, seam, line, …). Use
  those terms exactly; don't invent synonyms.
- The kit's portability across harnesses is a closed decision. Claude and Codex hook
  adapters are interactive layers on top of shared `.bench/hooks/` scripts and the
  harness-independent substrate (the `bench shift` loop + the git `pre-push` hook) —
  never the only thing enforcing an invariant.
- The two `projects/*.md` example profiles (gl-axi, regroup) are shipped templates,
  not live projects. This file (`benchkit.md`) is the profile for this repo.
