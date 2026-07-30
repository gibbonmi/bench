# Project: benchkit

The Bench kit itself — a harness-agnostic agent-development workflow shipped as the
npm package `redbench`. It is not an application: it is shell + markdown + JSON that
other repos consume. The deliverable is the `bench` CLI (`bin/bench.sh`), the working
agreement (`AGENTS.md`), the portable `.agents/` skills and commands, and harness
adapters that call shared `.bench/hooks/` scripts. Because the artifacts are plain
files, the kit must work identically under Claude Code, Codex, and any other
AGENTS.md harness — that portability is the product.

## Working branch

`main`. (The commit-on-green policy is canonical in `/bench-final-check`; the
default-branch guard is the pre-push hook, not commit-time — `bench commit` is
branch-agnostic. This line is only the binding.)

## Seams (test here; everything else is free to change)

- **The gate contract** (`.bench/gate.sh` / `bench gate`). The oracle surface. Everything
  in Bench routes through its exit code: 0 = shippable, non-zero = not done. The
  highest seam — if the gate is weak, the whole system is weak. Test by feeding it a
  conformant tree (green) and a broken one (red); never by trusting a reading of the
  diff. The gate package is the single deep owner of reusable-verdict authorization
  and durable execution; the accepted trust posture is recorded in ADR 0002.
- **The `bench` CLI subcommands.** The operational shell surface; the canonical
  inventory is `.bench/BENCH.md`'s CLI Inventory — this profile doesn't
  re-enumerate it. Stable command names and exit codes are the contract; the
  implementation behind each is free to change. Keep gate resolution
  (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) in one place.
  Missing-binary resolution is network-silent by default and names the explicit
  `bench repair` action; automation opts into the same repair path with exactly
  `BENCH_REPAIR=1`, while `BENCH_OFFLINE=1` and `BENCH_NO_REPAIR` suppress it.
- **The AXI query surface** (`bench learnings`, `bench maps`, `bench guards`,
  `bench diff`, `bench coverage`, `bench worktree list`, and the shared flat-table
  TOON emitter behind them). The agent-facing read-only
  surface, and the AXI-conformant half of the hybrid output contract: TOON stdout,
  definitive empty states, structured errors on stdout, exit 0/1/2. Gate-tested by
  the AXI contract fragments. Guard manifests come from each guard script's static
  leading-comment header (read as data, never executed), and
  `bench guards --brief` is the surface the SessionStart hook injects.
  `bench diff` is the single source of review-base truth: the shift loop records
  the pre-shift HEAD in `branch.<name>.benchBase`, `diff` resolves that key first
  and merge-base with the default branch as fallback. `bench coverage --check` is
  the one parser for the acceptance-coverage-map convention; the gate's docs
  fragment consumes it instead of carrying its own.
- **The ambient dashboard** (`bench status`). The single deterministic renderer the
  SessionStart hook and the user both call: it ranks the signals that fire on a fixed
  severity ladder and leads with the next action. Reads gate state from the **gate cache**
  (`<git-dir>/bench-last-gate`, written durably by gate execution) — never a cold gate run. The
  contract (gate-tested): show-only-on-signal, a five-row budget, a stale-green that is
  not a clean bill, and one combined capture-drain row (parked ideas + open learnings)
  pointing at `/bench-what-next`. A stale gate softens to `capture-only drift` /
  `re-run when convenient` only when every changed path is in the fixed, exact
  allowlist (`ROADMAP.md`, `IDEAS.md`, `.bench-notes.md`, `session-handoff.md` — no
  directory, suffix, or markdown-class matching; expanding it is a new decision); any
  mixed or untrusted diff fails closed to the strong stale row. Its severity-1 git
  signal reports dirty paths from the named/current checkout while aggregating
  unpushed commits and unique local branches across the repository; severity-2 intent
  joins the shared common-directory ledger, compact by default and expanded by `--all`.
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
- **The safe managed-asset lifecycle** (`bench link` and `bench unlink`). The adoption
  surface preflights and stages the complete write set, syncs durable content, and
  atomically promotes it or rolls the repository back. Relink reconciles old and new
  manifests: clean removed assets leave, while modified or project-owned collisions
  remain in place and produce a machine-readable partial result. Unlink removes only
  clean manifest-owned assets and reports any residuals with the same partial posture.
  The lifecycle preserves project-owned `AGENTS.md` text and installs the `.bench/bin/`
  local CLI set the shared hooks use when no global `bench` command is on PATH.
- **The distributable artifact contract** (wrapper and native package tarballs).
  The exact tarballs are the acceptance subject: one private-staging builder consumes
  the canonical asset manifest and platform matrix, artifact tests inspect and install
  its output, and the native workflow executes the same host smoke used locally. The
  installed shim keeps maintenance on the installed kit and routes operations through
  the linked repository's tracked launcher; staging details remain free to change.
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
- control bytes (ESC, BEL) in git-sourced text — commit subjects, branch names,
  paths — which `toon.Table` refuses rather than renders
- control bytes a sink *permits* but cannot survive: `toon.Representable` allows
  tab, newline, and return because the encoder escapes them, so a line-structured
  markdown or single-line sink borrowing that predicate accepts a value that
  splits its own field. Assert the permitted bytes, not only the refused ones —
  a test that exercises only what the predicate rejects proves nothing about the
  half that reaches the document
- a command whose own write changes a fact it reports: an artifact-rewriting
  command that also states tree cleanliness, staleness, or a derived next step
  falsifies its own output the moment it lands. Assert repeated application in
  the *tracked* configuration, and decide per field whether it excludes its own
  write or states the post-write truth — an untracked fixture holds the tree
  still and passes either way
- hand-edited files whose last line lacks a trailing newline
- absent file vs present-but-empty file (distinct behaviors, both asserted)
- special files in any discovered path — script discovery, control-record reads,
  the spec and decision sweeps (FIFOs, devices, sockets) — must be rejected
  before reading so neither static inspection nor an ambient command can block
- a dangling symlink where a file is expected: a plain read reports it as
  not-found, so a reader that does not stat first classifies a broken link as an
  authoritative empty state
- unquoted multi-word arguments (`$*` vs `$1`)
- a flag's value read as a positional: a parser resolving a subcommand or
  positional as "first token not `-`-prefixed" skips flags but not their values
  or anything after `--`, so `git stash -m list` and `git stash -- list` resolve
  to `list` and an allow verdict lands on the mutation the guard exists to
  refuse. The assertion must supply the sought token only as a flag value or
  pathspec — a case that also spells the positional explicitly passes for the
  wrong reason
- required tool missing from PATH (no global `bench`, no `readlink -f`)
- invocation through a symlink rather than the real path
- invocation through every shipped surface: real kit CLI, linked-repo by-path
  CLI, hooks, and adapters must all reach the same routed implementation
- destructive worktree state: foreign or identity-mismatched registrations, reused
  paths, the primary checkout, ignored residuals, dirty nested repositories, and
  plan/apply drift all fail closed without losing recovery state
- interrupt (SIGINT) mid-loop: leftover scratch state, leases, worktrees
- re-run idempotency: relink, reused worktree, second `init`, second `setup`
- cwd deeper than the repo root when the command assumes root
- non-TTY stdin on a prompting command must fail closed naming its
  non-interactive flags; `/dev/null` stdin reads as a character device, so
  TTY-detection contracts must drive a pipe, not the default null device
- host-backed filesystems under host-side I/O pressure: on WSL2, ext4 lives
  behind a VHDX whose file and directory `fsync` calls can stall for seconds
  even when guest-side CPU, memory, and `fsync` stress stay green

Known residual risk: `bench setup`'s real-TTY confirm wiring is one untested
constructor line binding stdin — testing it needs a pty dependency, which is a
reviewer decision the FT76 spec deliberately left open.

## Gate (`.bench/gate.sh`)

```
.bench/gate.sh
```

The oracle for a kit runs in two tiers. The **dev tier** — `bench gate`, the
shift loop, `/bench-final-check`, and the pre-push hook — answers one question:
does the kit work from the tree? It runs the layers below (all green today).
The **ship tier** — `bench prep-release` — carries the release-evidence checks,
described after the layers. `.bench/gate.sh` is a thin
exec into the `gate-phases` plumbing subcommand, which runs the layers below as
four concurrent phases in outer mode (`[phase]`-prefixed output, per-phase
verdicts, run-all-and-aggregate) and sequentially, unprefixed, sweep-skipped in
inner mode (`BENCH_CANARY_INNER=1`). A sibling knob, `BENCH_REQUIRE_CAPABILITIES=1`,
turns a nonzero capability-skip count red; absent, or set to anything else, the
`capability-skips` rows stay informational, because a developer's host legitimately
lacks capabilities and an unconditional red would make the gate unusable locally.
Both release workflows wire it on. An empty or legacy-flat canary run exercises
every non-canary phase; a nested fixture exercises only the phase its family
binds, which today means conformance. A `behavior-owned` fixture nests no gate
at all — it is graded by the compiled test binary of the one contract package
that owns its EXPECT, invoked against the fixture's tree. Canary EXPECT matching
is substring-based against that run's output, so the byte-shape of whichever
runs — inner gate or contract binary — is load-bearing:

1. **Go root conformance** — the conformance phase runs
   `go test -count=1 ./internal/conformance -run '^TestRootConformance$'` with
   `BENCH_CONFORMANCE_ROOT` set to the tree under grade. That suite owns parse and
   validity checks, JSON validity, executable git modes, package-file resolution,
   npm dry-run package shape, generated skills-index equality, Codex adapter
   metadata, Claude skill mirroring, shared-rule single-sourcing, stale command
   references, token-diet placement, workflow anchors, line-routing enforcement,
   compiled-core build/vet/test/cross-compile checks, release-workflow structure,
   static guard-header manifests, the profile's hostile-input checklist anchor,
   acceptance-coverage map validation, no bare `t.Skip`/`t.Skipf`/`t.SkipNow`
   outside the capability helper package, every subcommand entry point recorded in
   the routing registry, and no numeric duration literal passed to the marker-wait
   helper's slow-leg deadline.
2. **shellcheck** — best-effort, runs only when installed (`-S warning`). Not a hard
   dependency; upgrades the shell lint automatically once present.
3. **Managed-asset lifecycle behavior** — the gate runs link, relink, and unlink
   against throwaway repos to prove fresh installs, transactional rollback, manifest
   reconciliation across upgrade and downgrade, modified and stale asset handling,
   repeated lifecycle cycles, existing `AGENTS.md` preservation, same-name conflicts,
   Codex/Claude hook adapters, shared hooks, and default copy mode.
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
   is the gate guarding the gate. A nested fixture keeps the real gate entry path but
   runs only the phase that owns its failure, avoiding unrelated whole-gate work,
   while the baseline that rejects vacuous EXPECTs carries no phase pin and so still
   exercises every inner phase. A
   `behavior-owned` fixture nests nothing: its owning contract package is compiled
   once per package group and that binary is invoked per fixture tree.
   An EXPECT for a `behavior-owned` fixture is never checked for mutation-specificity:
   its contract group's empty-tree baseline is only a collision screen against
   infrastructure noise, so a generic banner that any failure prints passes it forever.
   Write the EXPECT as the owning contract test's own failure message for the specific
   mutated fact; what the comparison does and does not establish is stated on the
   baseline comparison in `internal/canary`.
   Fixtures hide dot-dirs behind a `dot-` prefix (e.g. `dot-claude`) so the harness
   doesn't load fixture skills as real ones; the canary restores them at run time.
   The sweep's aggregate concurrency is budgeted, not left to either factor: every
   run it makes — inner gate or contract binary — is invoked with an explicit
   `GOMAXPROCS` equal to
   `bounds.CanaryInnerWidth`, stripped-then-set so an inherited value cannot leak
   past the cap, and the worker pool derives as `runtime.GOMAXPROCS(0)` divided by
   that width, floored at one and capped at the fixture count. There is no
   Bench-specific knob — `GOMAXPROCS=8 bench gate` is the operator lever, and it
   shrinks the whole canary budget. The tripwire decision is recorded in
   `docs/adr/0001-working-tree-gate-tripwire.md`.

The **ship tier** — `bench prep-release`, maintainer-run once per release —
carries what the dev tier deliberately does not run: the release-evidence probe
(the four-platform artifact matrix build, the reproducibility rebuild, and a
real `release-preflight.sh --mode verify`), the cross-compile matrix
(`-tags stress`), the release-only package suites (`internal/preflight`,
`internal/releaseevidence`, `internal/publication` — excluded from the dev
tier's inner `go test`), and the ship-tier canary fixtures. It refuses to run
without a current dev-green verdict for the exact tree, so a dev-tier failure
reds the ship tier too. Exit 0 is ship green, with evidence at
`dist/preflight/release-index.json` and `dist/artifacts`.

**What dev green claims — and does not.** Dev green means the kit works from
the tree: every static conformance check (including the static half of the
release-preflight check), the lifecycle, contract, and AXI phases, and the
dev-tier canary passed. It does not claim the release artifacts build,
reproduce, or pass preflight verify — those are ship-tier facts, restaged to
run once per release instead of once per commit, with no check losing
authority. The dev contract suites drive `scripts/build-artifacts.sh` under
the shared-build-cache opt-in, so their green proves the generator's logic —
the planned artifact set, idempotently, with reproducible pins — not
byte-reproducibility across independent builds, which stays with the ship
tier above. Neither tier grants publish authority: the publish path's
`VerifyPublishAuthority` refusal demands a publish-mode index that only the
release workflow's own preflight produces, while `prep-release` emits
verify-mode evidence — a rehearsal, never a substitute for that boundary.

The gate file lives outside `package.json` `files[]`, so it never ships to consumers.

## Lines (model + effort routing)

The routing rubric — the three-signal decision table and the escalation ladder —
is the `craft-line` skill. This section holds what is project-specific: the
binding, the cached routings, and the escalation policy.

**Tier → model** (Codex; advisory candidates from `bench models`; set
2026-07-10): the reviewer-owned binding is top = Sol (`gpt-5.6-sol`) · mid =
Terra (`gpt-5.6-terra`) · cheap = Luna (`gpt-5.6-luna`). These opaque safe
tokens are this repo's current Codex choices, not a namespace rule; the token
grammar and discovery posture live in `craft-line`. Machine-readable source:
`.bench/lines.env`, read by the Agent-tool hook and the shift adapters; the
`checkLineBinding` conformance check cross-checks this paragraph against it, so
drift between the two turns the gate red. Claude Code delegation addresses
models by alias only, so `lines.env` separately declares the corresponding tier
aliases (`BENCH_ALIAS_TOP=fable`, `BENCH_ALIAS_MID=opus`,
`BENCH_ALIAS_CHEAP=sonnet`); those aliases do not resolve to the Codex model ids.
Code-author venue follows `craft-delegate`; line choice does not override its
threshold. Headless shift runs target `gpt-5.6-luna` via `BENCH_MODEL` through
the adapter (adapters take exact ids, not aliases).

**Escalation policy:** no standing top-tier opt-out — any bump to the top binding
(`gpt-5.6-sol`; `fable` in Claude Code) pauses and asks the reviewer (the ladder
is in `craft-line`). Tier moves still get declared — no silent escalation.

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
- **Spec falsification pass** (`/bench-write-spec` step 9) → **mid model, high
  effort, 1 iteration**, read-only. Standing grant: every draft gets the pass,
  spawned without asking, because at the mid binding it is not a top-tier bump.
  The step's signals nominate a draft for a top-binding pass instead — that is
  an ordinary bump and pauses and asks; never escalate silently. Charged at
  falsification questions, never an open review; its verdict is advisory and
  sign-off stays the reviewer's.
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
- A symbolic link inside an allowlisted kit payload tree is refused, not followed.
  Closed decision (2026-07-23): following the link would ship bytes the allowlist
  never named, so the allowlist would stop being the complete statement of what a
  consumer receives. Don't reopen it as a link/upgrade ergonomics fix.
- Never build `dist/bench` with plain `go build`; use
  `bash scripts/go-build.sh <root> <out>` so the binary carries the package
  version required by the version and upgrade contracts.
- To exercise or measure a worktree's build, invoke that worktree's own
  `./dist/bench`. `bench` on PATH resolves to the main checkout's wrapper,
  which runs the main checkout's `dist/bench` — a different binary, usually
  older than the one under test, so the timings and behavior belong to the
  wrong subject. The gate is exempt: it hands off to the worktree binary for
  its phases.
- Never mutate the repository while a gate is running. The gate binds its
  verdict to the starting subject and rejects a run whose subject changes.
- `internal/canary`'s own tests run nested. The conformance phase runs the kit's
  `go test` over core packages as a subprocess inheriting the phase environment,
  so inside a fixture's inner gate they run at `GOMAXPROCS=2`, where the derived
  worker bound is one. A concurrency expectation keyed to machine width does not
  merely fail there — it deadlocks until the phase timeout turns conformance red.
  Key every such expectation to the derived bound, and gate the ones that need
  real overlap through `capability.CPU`. Prove both directions before believing
  it: `GOMAXPROCS=2 go test -timeout 120s ./internal/canary` green with the
  expected `bench-skip kind=capability class=cpu` lines under `-v`, and the same
  run at full width green with none. A deleted assertion and an honest skip both
  look green; only the emitted line tells them apart.
- Never stop a gate by killing only its shell wrapper. Signal `gate-run`, which
  owns teardown of the gate script's process group, so canary and nested
  `gate-phases` children cannot outlive the run.
- Nothing under `.claude/` is a copy. `.claude/commands` is a git-tracked
  symlink to `../.agents/commands`; `.claude/skills` is a real directory whose
  every entry is a symlink to `../../.agents/skills/<name>`. So editing
  `.agents/` is the whole edit — there is no mirror to sync and no content drift
  for a check to guard — but *adding* a skill still needs its new
  `.claude/skills/` symlink entry created. Read "mirror" in any artifact as
  "symlink"; FT152's spec and map both assumed copied trees and specified
  mirror work that did not exist.
