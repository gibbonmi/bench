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
  Reviewed spec builds route through all eight lifecycle `bench spec build`
  operations; harnesses never synthesize their commit, ref, replay, or worktree
  plumbing. `reclaim` is the one maintainer-run verb beside that lifecycle —
  plan/apply removal of a terminal run's provably dead provisional refs, never
  driven by a build harness.
  Their final attachment to the oracle is `promote` over the exact reviewed
  prospective composition, while every earlier ticket transition is provisional.
- **The AXI query surface** (`bench anchors`, `bench learnings`, `bench maps`, `bench guards`,
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
  pointing at `/bench-what-next`. A stale gate softens to `reduced-scope drift` /
  `re-run when convenient` only when the gate's reduced-scope declaration confines
  every changed path — the same declaration the gate's stripped-worktree construction reads, rendered
  once in the Gate section's reduced-scope table, so the board's advice and the
  oracle's behavior cannot name different files; any mixed or untrusted diff fails
  closed to the strong stale row. Its severity-1 git
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

```
.bench/gate.sh
```

The oracle for a kit runs in two tiers. The **dev tier** — `bench gate`, the
shift loop, `/bench-final-check`, and the pre-push hook — answers one question:
does the kit work from the tree? It runs the layers below (all green today).
The **ship tier** — `bench prep-release` — carries the release-evidence checks,
described after the layers. Every non-reused direct or prospective dev-tier run
owns one private temporary host Bench executable built through
`scripts/go-build.sh`. The owner passes its cleaned absolute path as
`BENCH_RUN_BINARY` to the shell entry, freshness check, phase table, contract and
conformance helpers, stripped subject, and nested canary gates. Nested consumers
validate and reuse that path; they have no build fallback. Once all descendant
processes have stopped, the owner removes the private directory on every terminal
outcome, so no run artifact is available to a later process. The `gate-phases`
plumbing subcommand runs every outer and inner phase table serially in stable
topological order. Primary and stripped-subject phases share that one schedule;
output shape, dependent and optional skips, run-all-and-aggregate red behavior,
and process-group cancellation are unchanged. A sibling knob,
`BENCH_REQUIRE_CAPABILITIES=1`,
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
   helper's slow-leg deadline. Ordinary dev gates resolve one content identity per
   registered check and run the always-on authorization checks plus only the ordinary
   checks whose identities lack valid retained evidence. The selected ordinary set stays
   in registry order and shares this one process; output, timings, and the durable verdict
   name every executed check and the exact identity and authorship time covering every
   inherited check. `bench gate --fresh`, prospective execution, and ship ignore these
   reusable slots and execute their complete applicable inventories.
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

**The reduced-scope declaration.** The kit root always runs the full phase
table — there is no whole-changeset reduced run (reviewer decision, 2026-08-03),
and a legacy on-disk reduced verdict fails the loader's exact-field-set
validation, refusing as an invalid cache record and forcing a fresh run. The
declaration itself survives: it bounds the stripped-worktree construction below
and feeds the status board's softening. For dev lifecycle entry, an exact-tip
per-component partial green is whole-tree green when the gate package
revalidates the inherited evidence for every skipped component (the
per-component ancestor lookup is content-addressed with no freshness bound —
the retained full green serves until the component's input identity moves,
because the components whose inputs moved run fresh either way; reviewer
decision, 2026-08-01). This does not make a narrow verdict
reusable evidence for a later run, and it does not relax the ship tier's full-run
precondition. The declaration is single-sourced in the gate package
(`gate.ReducedScope()`) and rendered here; the scope-binding conformance check
cross-checks this table against it, so drift between the two turns the gate red:

| reduced scope | declared |
|---|---|
| directories | `capture/`, `decisions/`, `specs/` |
| files | `.bench-notes.md`, `ROADMAP.md` |
| excludable phases | `gofmt`, `vet`, `test`, `race`, `contract`, `shellcheck`, `canary` |
| included phases | `conformance`, `conformance-suite` |

Membership is location: a file entry matches byte-for-byte, and a directory entry
covers every descendant, so a file that lands under `capture/`, `decisions/`, or
`specs/` tomorrow is declared by construction. All three directories are entirely
formatted documents whose graders are the included phases — `specs/` joined on exactly
that ground (reviewer decision, 2026-08-01), and `decisions/` joins because conformance
grades active maps and their owned research assets are documents (reviewer decision,
2026-08-02). `.agents/` is deliberately absent: its Markdown is a real input to the
contract and canary components (lifecycle contracts link the kit's asset tree, and
canary fixtures seed from it), so a guidance edit rides the per-component input
declarations below — the toolchain components skip, the consumers run — rather than
joining the declaration itself, whose stripped-worktree enforcement would refuse
it (reviewer decision, 2026-08-02). The run owner builds before phase selection,
so no build phase belongs to either list. Excludability is enforced by
construction, to the construction's exact width: every full gate on the kit's own
root runs the excludable phases against a stripped worktree the declared paths are
absent from (a root that is not the kit runs unsplit — the declaration is the
kit's own, and a linked repo never made it), so an excludable phase that hard-fails
on a missing declared file, or degrades into an environment-kind skip — the kit's
idiom for an absent subject file — reds the next full gate and moves to the
included set. The construction proves nothing beyond those two signatures. A phase
that reads a declared path but stays green with the file gone is invisible to it,
and a mis-filed file is graded only by the included phases: a `.go` file committed
under `capture/`, `decisions/`, or `specs/` is seen by no excludable phase at all, and the
included phases do not grade Go formatting, so a full gate can pass a tree that
the same file outside the declaration would have redded. The declared directories
hold formatted documents; code does not belong in them. Narrowing is selected by
the changeset, never by a flag — per-component scoping is the only narrowing the
kit root takes; `bench gate --fresh` is the escape to a whole-tree run.

**Per-component input declarations.** Beneath the reduced-scope declaration's
floor, each evidence-skipped component declares its own input set, so a
changeset touching only one component's inputs runs that component and skips
the rest on retained ancestor evidence rather than switching every excludable
phase together. The declarations are single-sourced in the gate package
(`gate.ComponentInputSources()`) and rendered here; the per-component
conformance check cross-checks this table against them, so drift between the
two turns the gate red:

| component | declares | provenance |
|---|---|---|
| gofmt | `module-test-closure`, `manifest` | `derived` |
| vet | `module-test-closure`, `manifest` | `derived` |
| test | `module-test-closure`, `manifest` | `derived` |
| race | `module-test-closure`, `manifest` | `derived` |
| conformance-suite | `module-test-closure`, `manifest` | `derived` |
| contract | `module-test-closure`, `manifest`, `consumer-document-inventory` | `derived` |
| shellcheck | `shellcheck-argv` | `derived` |
| canary | `hand-declared` | `hand-written` |

`derived` means the component's inputs are computed from a named derivation on
every resolution — the module-wide `go list -deps -test ./...` closure plus the
module manifest for the toolchain and contract components (never the binary's
narrower `./cmd/bench` closure, which excludes the packages they grade), the
documents resolved from the consumer inventory added for `contract` because it
executes the selected binary and grades managed-asset lifecycle behavior,
and shellcheck's own argv enumeration for `shellcheck` — so a hand-copied path
list can never survive as the declaration. `canary` is the registry's one
`hand-declared` entry: `internal/canary/`, `tests/canary/`, `.agents/` (its
fixtures seed from the kit tree, so guidance edits move sweep expectations), and
the wrapper scripts its phase execs, named directly because it has no derivable
source.

**Per-check conformance inputs.** The lower conformance registry is the single source for
each check's name, tier, subject, executable binding, declared input source, and canary
ownership. The gate resolves those declarations against its exact Git subject; uncertain
checks use the complete subject as an explicit catch-all. Exact-file absence differs from
a present empty file, and a declared symlink contributes its canonical in-repository target
content. Broken, escaping, or unavailable targets widen execution. Every identity also
binds the shared conformance implementation closure and the invocation schema, so drift in
selection machinery cannot inherit an older green. A changed owning canary family moves its
check identity; a changed bound check implementation runs that family, while shared or
unattributable conformance implementation drift runs every conformance canary family.

| conformance check | input source |
|---|---|
| `kit-compliance` | `catch-all` |
| `canary-inner-compliance` | `catch-all` |
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
| `skip-ownership` | `go-source` |
| `decision-map-integrity` | `decision-documents` |
| `example-agreement` | `catch-all` |
| `component-honesty-prose` | `benchkit-profile` |
| `contract-capture-reads` | `go-source` |
| `injected-port-registry` | `go-source` |

The ordered outer selector is authored only by gate phase construction after ambient
singular and plural selectors are stripped. The singular selector remains the canary-owned
inner control. Unknown, duplicate, tier-invalid, out-of-order, incomplete, or overlapping
partitions red and widen rather than producing empty green. Meta checks never retain slots;
a green aggregate authors only executed ordinary slots, a red retires only those it
executed, and interruption authors none.

Declaration-honesty width, stated with the same candor the stripped-worktree
construction prose carries above: the stripped-worktree construction proves
capture-surface blindness only. For these per-component declarations, honesty
rests on mandatory derivation plus this binding, and a component that reads an
undeclared non-capture path skips wrongly — that residual is recorded here,
not hidden. `canary`'s row carries the reviewer's 2026-08-01 narrowing as its
own accepted gap: the published binary's digest is excluded from its declared
inputs, so two changes graded separately may land together with the canary
never run against their combined tree. `bench gate --fresh` and the ship tier
are what re-prove the tripwire in that case, not a component-scoped run. One more
recorded residual, accepted at the build's review: the slot and attestation
field-set slices exist for the record-class family registry only, so a field
added to a record struct without updating its slice is unobservable — the
registry's disjointness check cannot see a collision a stale slice hides.

The **ship tier** — `bench prep-release`, maintainer-run once per release —
carries what the dev tier deliberately does not run: the release-evidence probe
(the four-platform artifact matrix build, the reproducibility rebuild, and a
real `release-preflight.sh --mode verify`), the cross-compile matrix
(`-tags stress`), the release-only package suites (`internal/preflight`,
`internal/releaseevidence`, `internal/publication` — excluded from the dev
tier's inner `go test`), and the ship-tier canary fixtures. It refuses to run
without a current dev-green verdict for the exact tree, so a dev-tier failure
reds the ship tier too — and a partial verdict is refused the same way, with the
refusal naming the skipped components and pointing at `bench gate --fresh`,
because a partial verdict graded only the components whose inputs moved, never
the whole tree a release answers for. A legacy reduced record fails the loader's
exact-field-set validation and refuses as an invalid cache record, forcing a
fresh run. Exit 0 is ship green, with evidence at
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
- **Spec authoring** → **mid model, fresh session by default**.
  `/bench-write-spec` accepts exactly one of three sources: a ready compiled
  map, a reviewer-confirmed current conversation, or a named reviewed artifact, then
  derives engineering seams, tests, coverage, hostile-input attachment, and
  gate attachment from that source and the current tree. Top + high remains a
  reviewer-approved escalation. Distinct from the doc-authoring leverage
  override above: that spends the top tier on the kit's guidance prose.
- **Spec-build guidance cadence** → **`gpt-5.6-sol / high`**. Run `bench
  structure` before and after the guidance cut; the inherited findings are the
  baseline and the cut adds none. Both dogfood traces use the public porcelain:
  three ownership-safe tickets fill three slots, and integrating one unlocks a
  fourth assignment while another delegate remains active. The final composed
  gate runs only through `promote`; existing shift and ordinary-commit runtime
  contracts remain positive controls.
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
- **Ticket-breakdown review pass** (`/bench-implement-spec`, after the ticket
  files are written and before `bench spec build start`) → **mid model, medium
  effort, 1 iteration**, read-only. Standing grant like the falsification pass:
  every breakdown gets it, spawned without asking, and charge time never
  re-derives this routing. Charged at `craft-tickets`' consolidated target list,
  never an open review; its findings are reslices the coordinator repairs before
  the lifecycle starts.
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
