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
  Reviewed spec-backed builds keep serial green ticket commits on one retained
  integration source. Semantic review binds its frozen base and tip;
  `bench worktree land` composes and gates that pair on the destination, and its
  published commit owns the spec's `Status: implemented` flip.
- **The AXI query surface** (`bench anchors`, `bench learnings`, `bench maps`, `bench guards`,
  `bench diff`, `bench coverage`, `bench roadmap`, and `bench worktree list`, and the
  shared flat-table TOON emitter behind them). The agent-facing read-only
  surface, and the AXI-conformant half of the hybrid output contract: TOON stdout,
  definitive empty states, structured errors on stdout, exit 0/1/2. Gate-tested by
  the AXI contract fragments. Guard manifests come from each guard script's static
  leading-comment header (read as data, never executed), and
  `bench guards --brief` is the surface the SessionStart hook injects.
  `bench diff` is the single coherent review snapshot: it reports the revision,
  aggregate, inventory, checkout, whitespace, and optional complete patch from
  one movement-checked read. Spec-backed review supplies the retained source's
  explicit frozen base; bare mode retains its recorded-base fallback for other
  work. `bench coverage --check` is
  the one parser for the acceptance-coverage-map convention; the gate's docs
  fragment consumes it instead of carrying its own.
- **The ambient dashboard** (`bench status`). The single deterministic renderer the
  SessionStart hook and the user both call: it ranks the signals that fire on a fixed
  severity ladder and leads with the next action. Reads gate state from the **gate cache**
  (`<git-dir>/bench-last-gate`, written durably by gate execution) — never a cold gate run. The
  contract (gate-tested): show-only-on-signal, a five-row budget, a stale-green that is
  not a clean bill, and one combined capture-drain row (parked ideas + open learnings)
  pointing at `/bench-what-next`. A stale exact verdict always remains the strong stale
  row; there is no path-based reduced-scope softening. Its severity-1 git
  signal reports dirty paths from the named/current checkout while aggregating
  unpushed commits and unique local branches across the repository; severity-2 intent
  joins the shared common-directory ledger, compact by default and expanded by `--all`.
- **The capture inbox and working roadmap** (`bench idea` → `capture/IDEAS.md`;
  `bench roadmap` → `ROADMAP.md`). Capture-and-forget: park an out-of-scope idea,
  commit to nothing; ideas graduate only through a `/bench-what-next` drain into the
  working roadmap. The contract (gate-tested in a throwaway repo): `idea` appends one
  dated line and creates the inbox; a no-arg `idea` errors without appending;
  `roadmap` prints the working document plus drain status, or its
  `## Recommended sequence` when nothing needs draining. `capture/IDEAS.md` and `ROADMAP.md`
  are per-consumer content — never in the kit's `package.json` `files[]`.
- **The kit content surface** (`.agents/skills/*/SKILL.md`, `.agents/commands/*.md`).
  Portable harness-facing content. The contract is structural: every skill carries YAML
  frontmatter (name + description) and follows progressive disclosure; every command
  is a phase the index names. The `.bench/BENCH.md` skills index is generated
  from each skill's `index:` frontmatter (`bench skills-index --write`);
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
  against `.bench/BENCH.md`, and checks that no shared rule's literal marker phrase
  reappears in AGENTS.md — a substring check, so it reds a verbatim copy; a
  paraphrased restatement is review's to catch, not the gate's.

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
- a *live* symlink where a generator's input file is expected: following it makes
  bytes outside the graded tree authoritative, and a generator that then rewrites
  that path writes through the link to a target the tree never named. Refusing a
  broken link is not enough — the working link is the destructive half
- unquoted multi-word arguments (`$*` vs `$1`)
- lifecycle guidance that names every sanctioned operation but routes one step
  through raw Git anyway; swap the route while preserving all command tokens, so
  a token-presence check cannot pass synthesized commit, ref, replay, or worktree
  plumbing
- a flag's value read as a positional: a parser resolving a subcommand or
  positional as "first token not `-`-prefixed" skips flags but not their values
  or anything after `--`, so `git stash -m list` and `git stash -- list` resolve
  to `list` and an allow verdict lands on the mutation the guard exists to
  refuse. The assertion must supply the sought token only as a flag value or
  pathspec — a case that also spells the positional explicitly passes for the
  wrong reason
- a grammar token quoted in surrounding prose: a documented example of an inline
  annotation (a backticked `(covers AB1)` in a row's own text) parses as live
  syntax unless the grammar anchors the annotation to its position
- non-ASCII whitespace in hand-edited markdown: an NBSP where a
  position-anchored grammar permits only space and tab silently unanchors the
  token, so the diagnostic must stay fail-closed rather than false-positive
- required tool missing from PATH (no global `bench`, no `readlink -f`)
- invocation through a symlink rather than the real path
- invocation through every shipped surface: real kit CLI, linked-repo by-path
  CLI, hooks, and adapters must all reach the same routed implementation
- destructive worktree state: foreign or identity-mismatched registrations, reused
  paths, the primary checkout, ignored residuals, dirty nested repositories, and
  plan/apply drift all fail closed without losing recovery state
- interrupt (SIGINT) mid-loop: leftover scratch state, leases, worktrees
- re-run idempotency: relink, reused worktree, second `init`, second `setup`
- state serialized by one process and reloaded by a fresh one: the writer's
  in-memory value and the reader's re-parse agree at unit level and diverge
  across the boundary, so the assertion drives a second process rather than
  reusing the first's structures. Recomposition and recovery suites that stop
  at the first success prove one path and leave every other recomposition
  route unwalked
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

The dev gate answers one question for the exact current subject: does the complete kit
work from this tree? A non-reused run owns one selected Bench executable built through
`scripts/go-build.sh`, passes its cleaned absolute path as `BENCH_RUN_BINARY`, and
removes the private build after every terminal outcome. Ordinary commands route through
the production Go `Command` registry; wrapper, executable identity, freshness, and
process behavior remain at the bounded system seam.

The kit phase table is exactly:

| phase | authoritative argv |
|---|---|
| `gofmt` | `bench gate-go gofmt <root>` |
| `vet` | `go -C <root> vet ./...` |
| `test` | `go test -count=1 ./...` |
| `race` | one `go test -race -count=1 -v` invocation derived from `internal/racetests.Tests` |
| `system` | `go test -count=1 -tags=system ./internal/systemtest` |
| `shellcheck` | the stable shell-file inventory, optional when shellcheck is absent |

Go owns package scheduling inside the one ordinary test driver. There is no separate
contract or conformance dev driver, per-package loop, nested Go test, fixture-executing
canary phase, component partition, or stripped-subject phase schedule. The race runner
verifies every registry sentinel executed. The tagged system package has one
`TestMain` owner, at most three disposable repositories, one selected executable
identity ledger, teardown on green/red/interrupt/timeout, and exactly one
stripped-distribution journey.

The five command decision domains—gate, adopt, preflight, canary,
and freshness—consume immutable values in process. Their ordinary tests create no
repositories and start no operating-system processes. `internal/git` owns the one
ordinary repository adapter; `internal/gate` owns the one ordinary controlled process
group adapter.

Canary fixtures are immutable inputs to registered conformance checks. Each retained kit
fixture has exactly one check owner; its ordinary mutation test calls that owner directly,
requires the fixture-specific red, restores the fixture subject, and requires that red to
disappear. The top-level `bench canary` command validates and aggregates the complete
canary inventory without invoking owners or starting a successor process. Linked repos
receive that inventory validation and own planted-reason proof in their native tests. One
tagged system journey proves the selected executable reaches the production inventory path.

The workflow-guidance family pins the spec-to-ticket handoff from identified rows and
approved ownership fences through ticket evidence, ledger review, and fence-drift
repair. Its auto-discovered mutations keep each clause in
the section where a fresh agent acts on it.

The conformance registry remains the single source for check order, subject, input
derivation, implementation, and canary family ownership. This table is the profile's
current-state advertisement of its non-meta input bindings:

| conformance check | input source |
|---|---|
| `kit-compliance` | `catch-all` |
| `canary-fixture-compliance` | `catch-all` |
| `load-validity-metadata` | `catch-all` |
| `skills-index-command-adapters` | `catch-all` |
| `docs-currency-workflow` | `catch-all` |
| `gate-entry-contract` | `gate-entry` |
| `ordinary-build-census` | `catch-all` |
| `offline-smoke-proof` | `offline-smoke` |
| `handoff-shape-single-source` | `catch-all` |
| `harness-prefix-single-source` | `go-source` |
| `package-shipped-surface` | `catch-all` |
| `line-routing` | `catch-all` |
| `package-core-guard` | `catch-all` |
| `release-evidence-probe` | `catch-all` |
| `bench-sh-routes` | `bench-routes` |
| `default-branch-single-source` | `go-source` |
| `data-handling-derivation` | `go-source+data-handling` |
| `single-control-escaper` | `go-source` |
| `bounds-policy` | `catch-all` |
| `marker-wait-deadlines` | `go-source` |
| `subcommand-routing` | `go-source` |
| `axi-query-registry` | `catch-all` |
| `skip-ownership` | `go-source` |
| `decision-map-integrity` | `decision-documents` |
| `injected-port-registry` | `go-source` |
| `guidance-prose-budgets` | `benchkit-profile` |

A green verdict records the exact whole subject and oracle. Reuse is allowed only for a
current exact green; partial/component and reduced-scope records are legacy input classes
that fail closed and are never authored. Prospective execution uses the same complete
phase architecture. A stale exact verdict stays stale rather than being softened by path
classification.

`BENCH_REQUIRE_CAPABILITIES=1` makes capability skips fatal. Without it, capability
rows remain informational because a developer host may legitimately lack optional
facilities. The release workflows enable strict capability posture.

The ship tier remains `bench prep-release`, run once per release. It requires a current
exact dev-green verdict and owns release-evidence verification, cross-platform artifacts,
reproducibility, stress/cross-compile coverage, and publication/preflight rehearsal. Dev
green proves the complete branch-native source architecture; it does not grant publish
authority or claim release artifacts were reproduced.

### Guidance prose budgets

Guidance prose that outgrows a session's attention stops being read. This table is the
one source for how long each subject may be, and the `guidance-prose-budgets` check
parses it rather than repeating its numbers. An exact row beats the glob row a subject
also matches, so raising or lowering a budget is an edit here and nowhere else.

| subject | limit |
|---|---|
| `.bench/BENCH.md` | 180 |
| `.agents/commands/bench-implement-spec.md` | 60 |
| `.agents/commands/bench-write-spec.md` | 73 |
| `.agents/skills/bench-craft-tickets/SKILL.md` | 100 |
| `.agents/skills/bench-craft-spec/SKILL.md` | 150 |
| `.agents/skills/*/SKILL.md` | 120 |

The glob row is what classifies a newly added skill, so a skill arrives budgeted without
anybody editing the checker. Every other `.agents/commands/*.md` file is outside the
reviewed universe, and the `.claude/skills/*` adapter symlinks are distribution surfaces
rather than subjects — a symbolic link or special file found where a subject belongs is
refused unread.

## Lines (model + effort routing)

The routing rubric — the three-signal decision table and the escalation ladder —
is the `craft-line` skill. This section holds what is project-specific: the
binding, the cached routings, and the escalation policy.

**Harness × tier binding** (advisory candidates from `bench models`; set
2026-07-10): tier is the only shared identity, and each harness holds its own
column — no family is canonical.

| tier | codex | claude | opencode |
|---|---|---|---|
| top | `gpt-5.6-sol` | `fable` | unbound |
| mid | `gpt-5.6-terra` | `opus` | unbound |
| cheap | `gpt-5.6-luna` | `sonnet` | unbound |

These opaque safe tokens are this repo's current choices, not a namespace rule;
the token grammar and discovery posture live in `craft-line`. OpenCode is
unadopted here, so its column stays unbound and its adapter refuses to launch
rather than borrowing another harness's ids. Machine-readable source:
`.bench/lines.env`, read by the Agent-tool hook and the shift adapters, each
naming its own harness; the `checkLineBinding` conformance check cross-checks
this table against it, so drift between the two turns the gate red. Code-author
venue follows `craft-delegate`; line choice does not override its threshold. A
headless shift declares the tier — `BENCH_MODEL=cheap` — and the adapter's own
column supplies the id.

**Escalation policy:** no standing top-tier opt-out — any bump to the table's top
row pauses and asks the reviewer, whichever harness column you are running in
(the ladder is in `craft-line`). Tier moves still get declared — no silent
escalation.

- **Skill / command / doc authoring** → **top model, high effort**. This is the
  leverage override in `craft-line`: guidance prose compounds through every
  session that loads it while the edit costs few tokens. The `craft-skills` and
  `craft-adr` skills apply. Spend here.
- **Spec and ticket authoring** → **the session holding the decision source, at
  whatever tier it runs**.
  `/bench-write-spec` accepts exactly one of three sources: a ready compiled
  map, a reviewer-confirmed current conversation, or a named reviewed artifact,
  and authors the spec and tickets from that source and the current tree. Top +
  high remains a reviewer-approved escalation. After ticket approval, a fresh
  mid-tier session starts the build. Distinct from the doc-authoring leverage
  override above: that spends the top tier on the kit's guidance prose.
- **`bench` CLI shell plumbing** → cheap model, low–medium effort at the known seam.
  Mechanical once the gate-resolution and worktree-pool shapes exist.
- **Gate / conformance logic** → mid effort. Correctness of the oracle matters more
  than speed — a wrong gate is the worst class of bug in a kit whose whole premise is
  "the gate is the oracle."
- **Spec-and-tickets review round** (`/bench-write-spec`) → **mid model, high
  effort**, read-only and same-family through the harness's native agent surface.
  It reviews the spec and its ticket breakdown together against `craft-tickets`
  after `/bench-write-spec` slices it; `/bench-write-spec` owns the round's
  operating protocol.
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
- To exercise or measure a durable worktree artifact directly, invoke that
  worktree's own `./dist/bench`. `bench` on PATH resolves to the main checkout's
  wrapper and may belong to a different source tree. Gate and `bench test` runs
  do not reuse that artifact: their owner builds one private exact-source binary,
  propagates it through every ordinary child, and removes it after teardown.
- A direct `./dist/bench` invocation needs `BENCH_HOME` exported (the wrapper
  exports it; the gate-inputs closure declares it). Without it the prospective
  subject stays open, so a landing's gate can run fully green and still refuse
  as `prospective authorization refused: infrastructure`.
- Never mutate the repository while a gate is running. The gate binds its
  verdict to the starting subject and rejects a run whose subject changes.
- Canary mutation tests are ordinary in-process checks. Do not add a gate, wrapper,
  `go test`, or `go run` constructor to a fixture owner; the architecture census treats
  that as a regression.
- Never stop a gate by killing only its shell wrapper. Signal `gate-run`, which owns
  teardown of the gate script's process group and descendants.
- Nothing under `.claude/` is a copy. `.claude/commands` is a git-tracked
  symlink to `../.agents/commands`; `.claude/skills` is a real directory whose
  every entry is a symlink to `../../.agents/skills/<name>`. So editing
  `.agents/` is the whole edit — there is no mirror to sync and no content drift
  for a check to guard — but *adding* a skill still needs its new
  `.claude/skills/` symlink entry created. Read "mirror" in any artifact as
  "symlink"; FT152's spec and map both assumed copied trees and specified
  mirror work that did not exist.
