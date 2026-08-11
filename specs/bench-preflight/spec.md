# bench preflight

Status: staged

Decision source: `specs/remove-spec-build-lifecycle/decisions/pocock-alignment.md` ticket #7 (compiled ready map, reviewer-resolved 2026-08-11; check-4 predicate — row-ID scan with tag-scoped set equality — clarified with the reviewer in-session 2026-08-11)

## Problem

The retired lifecycle's failure mode was declared artifacts drifting from
reality: receipts claiming paths that never changed, tickets claiming rows a
spec never declared, reviews starting from a base that had moved. The gate
answers "is the tree done?" but nothing mechanical answers "do the artifacts
still describe reality?" at the moment a phase starts — so a build or review
begins on a stale base, an unauthorized write set, or an unowned coverage row,
and the drift is discovered late, by a human or not at all.

## Solution

One read-only command, `bench preflight <build|review> <slug>`, run at build
entry and review entry by the phase commands. It is the start-oracle over
artifacts-vs-reality — the gate stays the done-oracle over the tree. It runs
the map's five mechanical checks against the named spec and the working tree,
prints one definitive TOON verdict row per check, and exits non-zero on any
red; a red preflight stops the phase. It is fail-closed: an artifact it cannot
load or parse is a red, never a skip.

## User stories

1. As a coordinator at review entry, I want `bench preflight review <slug>`
   to run all five checks — base-current, paths-authorized, rows-owned,
   rows-membership, diff-nonempty — printing one verdict row per check and
   exiting 0 only when all are green, so a review never starts from artifacts
   that disagree with the tree.
   Line: opus / medium. The profile's cached gate/conformance routing — mid
   effort, because a wrong start-oracle verdict silently blesses the drift it
   exists to catch.
2. As a coordinator at build entry, I want `bench preflight build <slug>` to
   run the checks that are computable before tickets exist — base-current,
   paths-authorized, and the spec bootstrap — reporting the ticket checks as
   definitive not-applicable rows when `specs/<slug>/tickets/` is absent and
   running them for real when it is present, so a resumed half-built spec is
   validated and a fresh one is not falsely red.
   Line: opus / medium. Mode applicability is verdict-core logic, not shell
   plumbing, so it takes the same cached gate/conformance routing as story 1.
3. As an agent driving preflight against a broken input, I want every
   unloadable or malformed artifact to answer with a structured red naming
   the artifact — missing or symlink-dangling spec, coverage map invalid or
   not opted into row IDs, fences section absent or empty, tickets directory
   absent in review mode, a special file inside `tickets/`, not-in-repo,
   usage errors — never a green-by-omission, so fail-closed is a property of
   the command rather than a hope.
   Line: opus / medium. Fail-closed classification is verdict-core logic;
   same cached gate/conformance routing.
4. As a fresh session reading kit prose, I want the phase commands to invoke
   preflight at their entries, the CLI inventory and shell help to advertise
   it, the AXI seam list in the profile to claim it, and the changelog to
   carry the addition — with the two phase-entry clauses pinned by anchor
   rows so the prose cannot drift from the command.
   Line: fable / high. Generation-guiding kit prose takes `craft-line`'s
   leverage override at the profile's cached top/high routing.
5. As a kit maintainer, I want the subcommand-routing conformance check to
   derive dispatch names from the live `commandRegistry` composite literal —
   the `commands` map and `run()` function it parses today no longer exist —
   so the routing registry grades reality again and the new `preflight` row
   is actually enforced.
   Line: opus / medium. Gate-authoring under `craft-gate` discipline at the
   profile's cached gate/conformance mid-effort routing; the repair may not
   weaken any surviving predicate.

## Implementation decisions

- **The five checks and their exact predicates** (map #7's list, made
  precise):
  1. `base-current` — the resolved default branch's tip is an ancestor of
     HEAD: `merge-base <default> HEAD` equals `rev-parse <default>`. An
     unresolvable default branch is red. On the default branch itself this is
     trivially green.
  2. `paths-authorized` — every changed path since the resolved review base
     (committed, index, and tracked worktree — `bench diff`'s exact file
     set) is authorized by the spec's `## Ownership fences`: equal to a
     fence entry or under a fence prefix with a `/` separator. An empty
     changed set is green (0 paths).
  3. `rows-owned` — every row ID the spec's coverage map declares appears as
     a word-boundary token in at least one file under
     `specs/<slug>/tickets/`.
  4. `rows-membership` — every word-boundary token in the ticket files that
     matches a declared row tag (the alphabetic prefix of the spec's own row
     IDs, e.g. `PF` for `PF1..n`) names a declared row; tokens under other
     tags (`FT93`, cross-spec references) are ignored. With check 3 this
     forces set equality between declared and cited rows. Reviewer-confirmed
     predicate, 2026-08-11.
  5. `diff-nonempty` — the review base resolves (the rev-parse half) and the
     changed-file set is non-empty.
- **Mode applicability** (enumerated): `review` runs all five. `build` runs
  `base-current` and `paths-authorized` always; `rows-owned` and
  `rows-membership` run when `specs/<slug>/tickets/` exists (present but
  empty ⇒ red — rows unowned) and print `not-applicable` when it does not;
  `diff-nonempty` is `not-applicable` in build mode (a fresh build starts
  with an empty diff by design).
- **Bootstrap precondition, fail-closed** (map #10's posture): before any
  check runs, the spec must resolve via the existing spec resolver, carry
  `Status: staged`, and its coverage map must pass the coverage validator
  with row IDs opted in. A spec that fails bootstrap is a structured error
  naming what failed to load — a legacy 5-cell map is told to opt into row
  IDs, not silently passed.
- **Fence grammar, position-anchored**: a fence entry is a backticked token
  in the `## Ownership fences` section that is not inside parentheses —
  parenthetical prose is annotation, never authorization (the
  quoted-grammar-token hostile class). No backticked entries ⇒ red.
- **Reconciliation with map #3**: #3 demoted the per-*ticket* fence to an
  advisory `Writes:` note, "never for refusal". The spec-level
  `## Ownership fences` section is a different artifact — the
  reviewer-signed authorized write set, which #6's answer names as what
  review derives "from the spec". `paths-authorized` checks the spec fence
  only; it never reads or enforces anything ticket-level. #7 postdates #3
  and mandates the check; this is the exception line that records the two
  decisions as compatible.
- **Exact error contract**: every structured error is one `toon.Errorf`
  line — `error: <kind> — <hint>` on stdout, exit 1 — and every usage error
  is the usage parser's or `toon.Usage`/`toon.MissingArg`'s line with exit
  2, the same single-source constructors every AXI command uses. "Structured
  error" anywhere in this spec means exactly that shape.
- **Accepted trust assumption** (the `Bootstrap authority before execution`
  rule, applied): "a red preflight stops the phase" is prose obeyed by the
  agent running the phase command — there is no enforcement chain from a
  trust root, none is claimed, and map #7 deliberately placed invocation in
  the phase commands rather than the gate or hooks. The reviewer accepts
  agent compliance as the stop mechanism; the anchors pin the prose, not
  the obedience.
- **Routing-checker repair** (story 5): `dispatchNames` currently parses a
  `commands` map and a `run()` function that were replaced by the
  `commandRegistry` composite literal, so the real-tree check reports every
  registered name as no-longer-dispatched — a latent ship-tier red on
  today's main. The repair derives dispatch names from `commandRegistry`'s
  `Name:` fields, updates the fixture bite tests to the same shape, and
  weakens no predicate; the accepted "reaches the grammar helper" predicate
  stays the AST reference the checker already implements.
- **Where conformance bites**: root-subject conformance checks (routing,
  docs currency, anchors) execute under `bench prep-release` (ship tier),
  not the dev gate's phase table. Each conformance red in the coverage map
  is demonstrated during the build with the direct invocation
  `BENCH_CONFORMANCE_ROOT=<root> go test -count=1 -run '^TestRootConformance$'
  ./internal/conformance`; the dev-observable seams for the feature itself
  are the ordinary Go tests.
- **Package naming**: the existing `internal/preflight` (release
  authorization, CLI verb `release-preflight`) is renamed
  `internal/releasepreflight` for verb parity, and the new decision domain
  takes `internal/preflight` to match its verb. Mechanical: seven importing
  files plus the package's own six. Flagged for reviewer veto.
- **Compose, never re-derive** (one source per fact): row IDs and map
  violations via the coverage package's exported spec parser (its stale doc
  comment naming deleted lifecycle consumers is repaired in the same
  change); spec resolution and status via the spec package; the default
  branch via the git package's existing resolver; the review base via a
  newly exported resolution function in the diff package that `bench diff`
  itself consumes, so the two commands cannot disagree about the base —
  bare `bench diff` output stays byte-identical.
- **Architecture**: the new package follows the decision-domain discipline —
  a pure verdict core consuming immutable gathered facts, ordinary tests
  creating no repositories and starting no processes; one thin gatherer does
  the git/filesystem reads. Ticket-file enumeration stats each entry and
  refuses special files (FIFOs, devices, sockets) before reading, with the
  lstat-mode check as the named validator, so no ambient invocation can
  block.
- **CLI surface**: routed `route_porcelain` in `bin/bench.sh` with a help
  line; a plain `commandRegistry` row; grammar through the usage parser (the
  routing conformance check requires the package to call it); AXI output
  contract — TOON stdout (`phase:`, `spec:`, then a
  `checks[N]{check,verdict,detail}` table), structured errors, exit 0 green
  / 1 red or structured error / 2 usage. The profile's AXI seam bullet gains
  the command.
- **Docs in the same green change** (the docs-currency sweep is
  bidirectional): `.bench/BENCH.md` CLI inventory line, `bin/bench.sh` help,
  the two phase-command entry steps, `projects/benchkit.md` AXI bullet, one
  `CHANGELOG.md` line. Two new workflow-guidance anchor rows pin the
  phase-entry clauses.

## Testing decisions

- A good test drives the public CLI (`bench preflight …`) in a throwaway
  repository seeded with one defect per check, asserting the exact verdict
  row and exit code; prior art: the diff and coverage command tests and the
  testrepo harness. The verdict core gets table-driven tests over immutable
  fact values, prior art: the existing decision-domain tests.
- Seams receiving tests: the CLI contract seam (TDD-marked) and the verdict
  core (TDD-marked). The exported diff-resolution function keeps the diff
  package's existing tests as its regression net.
- The gate seam observing the feature: the ordinary Go suite in the dev
  gate; the conformance layer (subcommand routing, docs-currency sweep,
  anchors) at the ship tier, with each conformance red demonstrated during
  the build via the direct root-conformance invocation named in the
  implementation decisions.

### Seam diagram

    trigger: /bench-implement-spec entry (build) · /bench-review-implementation entry (review)
        │
        ▼
    mode + slug ──▶ [ bench preflight: gatherer (git / diff base / spec /
                      coverage / tickets) → pure verdict core ] ──▶ TOON checks table, exit 0/1/2
                        ◀ tests attach here: CLI contract test seeds one defect per
                          check in a throwaway repo and asserts its row + exit;
                          verdict-core table tests drive the pure function directly

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| PF1 | 1 | `bench preflight review <slug>` on a conformant tree (branch on current default tip, staged spec, valid opted-in map, all rows cited, non-empty authorized diff) prints five green rows and exits 0; a second run is byte-identical | CLI contract | observed red 2026-08-11: `bench preflight review remove-spec-build-lifecycle` → `bench: unknown subcommand: "preflight"`, exit 2 | the whole surface is absent today; this is the green tracer every other row cuts against |
| PF2 | 1 | with the default branch advanced past the branch point, `base-current` is the red row and exit is 1 | CLI contract | same observed red as PF1 — the verb is unknown, so the assertion fails today | an always-green stub passes PF1; a seeded stale base catches it |
| PF3 | 1 | with a tracked change outside every fence entry, `paths-authorized` is red naming the offending path; a changed path equal to a fence entry or under a fence prefix stays green | CLI contract | same observed red as PF1 | catches both the stub and an over-permissive prefix match (`internal/git2` must not match fence `internal/git`) |
| PF4 | 1 | with one declared row ID cited by no ticket file, `rows-owned` is red naming the uncited ID | CLI contract | same observed red as PF1 | an uncited row is exactly the declared-artifact drift the lifecycle's receipts hid |
| PF5 | 1 | a ticket token under the spec's own tag naming no declared row (`PF99`) makes `rows-membership` red; a foreign-tag token (`FT93`) is ignored — both cases asserted | CLI contract | same observed red as PF1 | phantom citations are the other half of set equality; tag scoping keeps cross-spec references from false-positives |
| PF6 | 1 | with an empty changed set in review mode, `diff-nonempty` is red | CLI contract | same observed red as PF1 | reviewing nothing is the rev-parse+non-empty-diff failure the map's check 5 names |
| PF7 | 1 | a changed path carrying a control byte (ESC) makes the command exit 1 with the unrepresentable-TOON-cell error, never a mangled table; a path with a space or glob character renders escaped and authorizes correctly | CLI contract | same observed red as PF1 | the TOON sink refuses control bytes; asserting only clean paths would prove nothing about the refusal half (profile hostile classes) |
| PF8 | 1 | with `branch.<name>.benchBase` recorded past an out-of-fence commit, `paths-authorized` is green; with the key removed (merge-base fallback) the same tree is red — the CLI observably consumes the recorded-key resolution `bench diff` uses, and bare `bench diff` output stays byte-identical after the export | CLI contract + diff package | same observed red as PF1 | an implementation that calls the exported resolver and discards its result, re-deriving by merge-base alone, answers wrong under a recorded key — consumption is asserted behaviorally, not by compilation |
| PF9 | 2 | `bench preflight build <slug>` with no `tickets/` directory prints `not-applicable` for the two row checks and `diff-nonempty`, runs the rest, exits 0 when they are green | CLI contract | same observed red as PF1 | a fresh build must not be falsely red, and silence is not a verdict — the not-applicable rows are printed, definitive states |
| PF10 | 2 | `bench preflight build <slug>` with a present `tickets/` directory runs the row checks for real; present-but-empty `tickets/` is red (declared rows unowned) | CLI contract | same observed red as PF1 | absent and empty are distinct behaviors (profile hostile class); an empty directory passing would let a resumed build skip row accounting |
| PF11 | 2 | `base-current` red in build mode exits 1 even while ticket checks are not-applicable | CLI contract | same observed red as PF1 | not-applicable rows must not soften real reds into a green exit |
| PF12 | 3 | a missing spec, or a dangling symlink where the spec should be, answers a structured error naming the spec path, exit 1 | CLI contract | same observed red as PF1 | a reader that does not stat first classifies a broken link as an authoritative empty state (profile hostile class) |
| PF13 | 3 | a spec whose coverage map fails validation, or a legacy 5-cell map without row IDs, answers a structured error carrying the coverage validator's message and naming the row-ID opt-in, exit 1 | CLI contract | same observed red as PF1 | fail-closed bootstrap: row accounting over an invalid map would be a derived count with no source |
| PF14 | 3 | a spec whose `## Ownership fences` section is absent, empty, or contains no backticked entry outside parentheses answers a structured error, exit 1; a backticked token inside parentheses is never an authorization | CLI contract | same observed red as PF1 | an empty fence section is incomplete, not unrestricted authority; the parenthesized-token rule is the anchored-grammar hostile class |
| PF15 | 3 | review mode with `tickets/` absent is a structured red; a FIFO or other special file inside `tickets/` is refused before reading, named in the error, and the command does not block | CLI contract | same observed red as PF1 | review without tickets has nothing to account; an unread special file would otherwise hang every phase entry |
| PF16 | 3 | outside a git repository the command answers the standard not-in-repo error; a missing or unknown mode, missing slug, or unknown flag exits 2 with usage | CLI contract | same observed red as PF1 | the AXI error contract (exit 0/1/2, structured errors) is what makes the command safe to wire into phase prose |
| PF17 | 3 | a spec or ticket file whose last line lacks a trailing newline parses identically to one that has it | verdict core | same observed red as PF1 | hand-edited markdown routinely drops the final newline (profile hostile class); a scanner keyed to `\n` silently drops the last row citation |
| PF18 | 4 | the subcommand-routing registry carries `"preflight": routed("internal/preflight")`, the package carries an AST reference to `usage.Parse` (the checker's accepted predicate), and `checkColdPickupCLILists` passes with the new `.bench/BENCH.md` inventory line | conformance (root subject, demonstrated via the direct root-conformance invocation) | observed red by construction 2026-08-11: the route without the doc line fails the cold-pickup list check, and a dispatched name without a registry row fails the routing check — both asserted after story 5's repair makes the routing check see dispatch again | the conformance layer is bidirectional — code without docs and docs without code both red |
| PF19 | 4 | each phase command's entry step carries the sentence naming `bench preflight build <slug>` (implement) / `bench preflight review <slug>` (review) and the clause "a red preflight stops the phase"; one anchor row per file pins exactly that sentence | anchors / conformance | observed red by construction 2026-08-11: the anchor rows are authored against prose that does not exist yet, red until the phase-command edits land | un-pinned phase prose is exactly what drifted from the CLI in the lifecycle era; anchors make the drift a gate red |
| PF20 | 4 | `CHANGELOG.md` carries one line naming the new command | conformance / changelog | observed red 2026-08-11: no such entry exists | the changelog is the one advertised migration surface for CLI additions (Spec A precedent, RM12) |
| PF21 | 2 | `bench preflight build <slug>` with a tracked change outside every fence entry makes `paths-authorized` red, exit 1 | CLI contract | same observed red as PF1 | hard-coding build-mode `paths-authorized` green passes PF9–PF11; this row is the build-mode red the story's promise requires |
| PF22 | 3 | a spec whose `Status:` is anything but `staged` (e.g. `implemented`) answers a structured error naming the found status, exit 1 | CLI contract | same observed red as PF1 | the bootstrap's staged-status requirement was otherwise prose-only — an implementation accepting any status would pass every other row |
| PF23 | 5 | the routing check derives dispatch names from `commandRegistry` and, on today's tree shape, reports real violations instead of "no longer dispatches" for every name; its fixture bite tests drive the registry shape | conformance (root subject) | observed red 2026-08-11: `BENCH_CONFORMANCE_ROOT=<root> go test -count=1 -run '^TestRootConformance$' ./internal/conformance` emits "the subcommand argument-routing registry names \"X\", which cmd/bench/main.go no longer dispatches" for every registered name | a routing check blind to the live dispatch surface is a vacuous oracle — it can neither pin `preflight` nor catch the next unregistered verb |
| PF24 | 1 | the release-authorization domain lives at `internal/releasepreflight` with the `release-preflight` verb and its exported surface unchanged, and the literal `internal/preflight` appears in no live surface — imports, conformance census and release-only registry strings, the structure accept list, the package-core-guard canary baseline, prose comments, `.bench/BENCH-reference.md` — outside the exempt historical set (`specs/`, `CHANGELOG.md`, `capture/`, `ROADMAP.md`, `decisions/`; exemption reviewer-approved 2026-08-11) | rename prefactor (Go build, conformance census, structure accept list, canary baseline) | observed red 2026-08-11: `internal/releasepreflight` does not exist and the literal is live in 13 surfaces | a half-renamed tree leaves the census, accept list, and canary baseline grading a package that no longer exists, and squats the path the new domain needs |
| PF25 | 1 | the diff package exports one review-base resolution — recorded `branch.<name>.benchBase` key first, default-branch merge-base fallback, structured error when neither resolves — that `resolveBranchRange` itself consumes; bare `bench diff` output stays byte-identical | diff package | observed red 2026-08-11: `resolveBase` is unexported and handles only the recorded key; the fallback and error semantics live inline in `resolveBranchRange`, so no complete resolution exists for a second consumer | two consumers deriving the base independently is exactly the disagreement the compose-never-re-derive decision forbids |

### Edge inventory

- Error path — PF2–PF6, PF21 (per-check reds), PF12–PF16, PF22 (bootstrap
  and usage reds).
- Empty/absent input — PF6 (empty diff), PF9/PF10 (absent vs
  present-but-empty `tickets/`), PF14 (empty fences), PF15 (absent
  `tickets/` in review).
- Boundary values — PF3's prefix boundary (`internal/git` vs
  `internal/git2`). No numeric boundaries otherwise; **Won't handle:**
  none apply.
- Malformed input — PF13 (invalid map), PF14 (fence grammar), PF17
  (missing trailing newline), PF5 (phantom row token), PF22 (wrong spec
  status).
- Interrupted/partial state — **Won't handle:** the command is read-only
  and writes no state; an interrupt loses only stdout already printed.
- Re-run idempotency — PF1's byte-identical second run.
- Process-boundary lifecycle — **Won't handle:** no state is serialized or
  reloaded; every run re-derives from the tree.
- Hostile environment — PF7 (control bytes, spaces, globs in git-sourced
  paths), PF12 (dangling symlink), PF15 (special files), PF16
  (not-in-repo). Non-TTY stdin: **Won't handle** — the command never
  prompts. WSL2 fsync pressure: **Won't handle** — the command performs no
  writes.
- Self-falsifying write — **Won't handle:** read-only command; no field it
  reports can be changed by its own execution.
- **Won't handle:** a slugless light-path mode (checks 1 and 5 only, no
  spec) — separate capability, priced in Out of scope.
- **Won't handle:** mechanical enforcement (a hook blocking a phase on red)
  — map #7 places invocation in the phase commands, not the gate or hooks;
  a red stops the phase by the phase's own contract.

## Ownership fences

- `specs/bench-preflight/` (this spec and its future tickets)
- `internal/preflight/` (new decision domain), `internal/releasepreflight/`
  (renamed release-authorization package)
- `internal/preprelease/`, `internal/publication/`,
  `internal/conformance/`, `internal/conformance/registry/` (rename
  imports; routing-registry and docs-check rows)
- `internal/diff/` (exported base resolution), `internal/coverage/` (stale
  doc-comment repair only), `internal/anchors/` (two new anchor rows)
- `tests/canary/workflow-guidance-anchors/` (fixtures for the new anchors,
  if the anchor family requires them)
- `cmd/bench/` (registry row, grammar)
- `.bench/structure-accept` (renamed package entry),
  `tests/canary/package-core-guard/` (renamed baseline paths) — added by
  reviewer-approved fence repair 2026-08-11
- `bin/bench.sh` (route + help)
- `.bench/BENCH.md`, `.bench/BENCH-reference.md` (CLI inventory; renamed
  plumbing mention)
- `.agents/commands/bench-implement-spec.md`,
  `.agents/commands/bench-review-implementation.md` (entry steps)
- `projects/benchkit.md` (AXI seam bullet)
- `CHANGELOG.md`, `capture/session-handoff.md`

## Out of scope

- **Slugless light-path preflight** (`bench preflight review` with no spec:
  base-current + diff-nonempty only): ~5 edits, 1 gate run. Separate
  capability — it serves work the entry contract of this feature explicitly
  keys to a spec.
- **`bench status` preflight signal** (a dashboard row when the staged
  spec's preflight is red): ~4 edits, 1 gate run. Separate capability with
  its own severity-ladder decision.
- **Hook-enforced phase blocking** (~8 edits, 2 gate runs): contradicts map
  #7's placement of invocation in the phase commands; recorded here so the
  exclusion is visible, not silent.
- **Spec C — doctrine adoption** (map #4, #5, #6, #8, #9, #10): the review
  re-derivation mandate (#6) that consumes preflight's guarantees lands
  there, ~40 edits, 4 gate runs.
- **Other ship-tier conformance reds discovered alongside story 5** (the
  same root-conformance run): `decision-map-integrity` violations from
  `decisions/gate-budget.md` and `decisions/spec-build-review-gate-cadence.md`
  citing deleted files (map retirement is already queued for the
  `/bench-what-next` drain), and the injected-port derivation reporting
  "found no ports" for `internal/canary` (~unknown cause, needs its own
  diagnosis; parked to the drain). Repairing them here would cross this
  spec's capability; they are named so the exclusion is visible.
