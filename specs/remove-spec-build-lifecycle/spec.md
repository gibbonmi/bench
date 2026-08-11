# Remove the provisional spec-build lifecycle

Status: staged

Decision source: `specs/remove-spec-build-lifecycle/decisions/pocock-alignment.md` (compiled ready map, tickets #1, #3, #11, #12; reviewer-resolved 2026-08-11)

## Problem

The provisional spec-build lifecycle defers contact with the composed tree to
promotion and substitutes declared artifacts — receipts, fences, closure
graphs — in between. In practice it manufactured the failures it was built to
prevent: an abandoned run with mutually stranded ownership fences, 24 repair
tickets against a 126-line spec, eight review rounds that sampled instead of
re-derived, and a learnings journal consumed by receipt schemas and abandon
preconditions. The machinery verifies consistency; completeness never had a
home in it.

## Solution

Delete the lifecycle wholesale — no compatibility surface, per the map's #12.
Spec-backed work lands the same way light-path work always has: serial
commit-on-green per ticket on a real branch, `bench commit` as the sole
landing path, `--spec` on the green landing commit as the sole author of
`Status: implemented`. The worktree pool, gate, and capture surfaces stay;
everything that existed to manage work-not-yet-on-the-branch goes.

## User stories

1. As an agent driving a build, I want the `bench spec build` grammar gone —
   any `spec build` invocation answers with the CLI's standard
   unknown-subcommand structured error — with `internal/specbuild`, its cmd
   wiring, and its conformance-registry rows deleted and the suite green.
   Line: opus / high. The deletion crosses the CLI dispatcher and four
   conformance registries, and a half-removed registry row is a silent gate
   weakening.
2. As a session starting in this repo or any linked repo, I want the
   provisional and preservation ref machinery gone: `bench worktree recovery`
   out of the grammar, `bench resume` no longer authoring preservation refs,
   its reconcile deleting the lifecycle namespaces `refs/bench/specbuild/`
   and `refs/bench/recovery/` — never the gate's `refs/bench/green/` or
   diagnostic refs — and purging legacy lifecycle assignments from the intent
   ledger, idempotently. Line: opus / high. Ref and ledger reconcile
   semantics run unattended at every session start, so a defect here corrupts
   state silently.
3. As a coordinator landing a spec's final ticket, I want `bench commit
   --spec <slug>` to stop consulting run state and be the sole
   staged→implemented author, with `bench commit --help` stating exactly that.
   Line: sonnet / medium. The change is small and fully gate-observable at the
   commit seam.
4. As a kit maintainer, I want the binary to carry no ticket parser —
   `Blocked by:` stays a documented convention agents read — and the
   example-agreement conformance check retired with the parser it exercised.
   Line: opus / medium. Retiring a gate check rides `craft-gate` discipline:
   legitimate only because its subject is deleted, and the removal must not
   soften any surviving check.
5. As a fresh session reading kit prose, I want every reference to removed
   commands repaired — `.bench/BENCH.md` CLI inventory and workflow,
   `.bench/BENCH-reference.md`'s lifecycle lookup, the phase commands,
   `craft-tickets`/`craft-delegate` lifecycle mentions, the 33 spec-build
   anchor rows in `internal/anchors/registry_data.go` with their canary
   fixtures, `projects/benchkit.md`, `ROADMAP.md`, and `CHANGELOG.md` (one
   removal line) — proven by a new red-capable removed-verb sweep, because
   the existing sweeps match only slash-command and top-level-verb tokens and
   would pass a leftover `bench spec build promote`. Line: fable / high.
   Generation-guiding kit prose takes `craft-line`'s leverage override at the
   profile's cached top/high routing; this story repairs references only —
   the doctrine rewrite is Spec C's.

## Implementation decisions

- `internal/specbuild` is deleted whole, with its cmd dispatch file and the
  compile-time port pins. The shell wrapper drops the `spec build` grammar and
  help block; `bench spec implemented|retire|history` survive unchanged.
- `internal/commit` drops its specbuild status consult; the active-build
  refusal is deleted rather than stubbed. `--spec` semantics land in the help
  text (closes the open learnings entry about `--spec`).
- `internal/coverage` keeps the acceptance-coverage-map parser and `--check`;
  its covers-annotation reference is a comment only (no specbuild import) and
  is retired with the annotation convention.
- The worktree pool (`create/path/exec/release/clean`) and the intent ledger
  stay. `recovery` leaves the grammar; `bench resume`'s reconcile becomes the
  standing cleaner: it deletes refs under `refs/bench/specbuild/` and
  `refs/bench/recovery/` only — `refs/bench/green/<branch>` and diagnostic
  refs are the gate's verdict store (`internal/gate`) and survive — and drops
  lifecycle-typed or unknown-typed assignments from the ledger, making the
  one-time debris cleanup (the stale share of today's 190 refs, ~20 stale
  assignments) the same code path as steady-state idempotent hygiene.
  Zero-compat per map #12: no shim, no migration messaging beyond the
  changelog line.
- The existing docs sweeps cannot see removed multi-token verbs (they match
  `/bench-*` tokens and top-level verb names, over a file set that excludes
  `projects/`, `README.md`, `CHANGELOG.md`, and `ROADMAP.md`), so the build
  adds one removed-verb sweep — literal `spec build` and `worktree recovery`
  tokens over the full kit-prose set — that is red on today's tree and stays
  as the standing guard.
- The 33 anchor rows pinning spec-build prose and their
  `tests/canary/workflow-guidance-anchors/spec-build-*` fixtures retire with
  their subjects under `craft-gate` discipline; no surviving anchor or
  check's predicate is weakened.
- Conformance registries (injected ports, ordinary-build census, fixture
  registries) lose their specbuild rows; the example-agreement check retires
  with `ParseTicket`. No surviving check's predicate is weakened.
- Ticket files keep `Blocked by:` as their one machine-readable line, read by
  coordinators, not by the binary.

## Testing decisions

- Good tests here drive the public CLI surface (`bench …` via the shell
  wrapper or `cmd/bench`) and the resume reconcile seam, asserting observable
  outputs: structured errors, ref-namespace emptiness, ledger contents, spec
  status lines. Prior art: the existing CLI contract tests and
  `internal/worktree`/`internal/intent` test suites.
- The gate seam observing the feature is the ordinary dev gate: the Go suite
  plus the conformance layer (docs-contracts sweeps, registry checks).

### Seam diagram

    trigger: any bench invocation / SessionStart resume
        │
        ▼
    argv ──▶ [ bin/bench.sh → cmd/bench dispatcher ] ──▶ structured output / exit
                  ◀ tests attach here: invoke removed grammar, assert the
                    standard unknown-subcommand error
        │
        ▼
    repo state ──▶ [ resume reconcile (refs/bench/*, intent ledger) ] ──▶ clean state, report
                  ◀ tests attach here: seed legacy refs + ledger entries in a
                    throwaway repo, run resume twice, assert empty + idempotent

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| RM1 | 1 | `bench spec build <verb>` returns the standard unknown-subcommand structured error | CLI dispatcher | observed red 2026-08-11: `bench spec build status <slug>` answers with a live lifecycle precondition error, so an unknown-subcommand assertion fails today | a surviving dispatch route means the lifecycle is reachable, not removed |
| RM2 | 1 | `rg -l specbuild` over the tree is empty: the package dir, `cmd/bench/specbuild.go`, `wiring_pins.go` pins, the specbuild port rows in `injected_ports_test.go`, the census entry in `ordinary_build_census_test.go`, and the fixture rows in `fixture_bite_test.go` are all gone, and `go test ./...` is green | Go suite / conformance | observed red 2026-08-11: `rg specbuild internal/conformance` lists census, injected-ports, and fixture rows | an orphaned registry row is a gate check grading a deleted subject — red or vacuous, both defects |
| RM3 | 2 | after resume's reconcile, `git for-each-ref refs/bench/specbuild refs/bench/recovery` is empty (any ref name inside those namespaces), stays empty across a create/release cycle, and refs outside them are untouched | resume reconcile | observed red 2026-08-11: 190 refs under `refs/bench/`, the stale share in `specbuild/` and `recovery/` | leftover provisional refs are the "looks durable, isn't landed" attractor the map's #1 closes |
| RM10 | 2 | `refs/bench/green/<branch>` and diagnostic refs survive a resume reconcile byte-identical | resume reconcile | not TDD-able as a pre-build red — the reconcile that could eat them doesn't exist yet; the implementing ticket authors this guard test first at the new seam | the gate's verdict store shares the `refs/bench/` prefix; an over-broad delete silently destroys green evidence at every session start |
| RM4 | 2 | `bench worktree recovery` answers with the unknown-verb structured error and help omits it | CLI dispatcher | observed red 2026-08-11: `bench worktree --help` line 8 prints the recovery grammar | an advertised dead verb sends agents to a removed surface |
| RM5 | 2 | resume over a ledger holding lifecycle-typed or unknown-typed assignments exits 0, purges them, authors no preservation refs; second run is a no-op | resume reconcile | observed red 2026-08-11: this session's SessionStart preserved 20 recovery refs from open assignments | the reconcile runs unattended in every repo; a crash or re-preservation loop is silent corruption |
| RM6 | 3 | `bench commit --spec <slug>` flips staged→implemented on green with no run-state consult; `--help` states the semantics | commit seam | observed red 2026-08-11: `internal/commit/commit.go:37` consults `specbuild.New(...).Status`; help lacks the line (open learnings entry) | the lifecycle's status authority must move, not dangle — a dangling consult panics or refuses on deleted state |
| RM7 | 4 | no ticket-parser symbol remains; example-agreement check retired | Go suite / conformance | observed red 2026-08-11: `example_agreement_test.go` parses the craft-tickets example via `specbuild.ParseTicket` | a surviving parser re-imposes the deleted schema on tickets through the gate |
| RM8 | 5 | a new removed-verb sweep — literal `spec build` and `worktree recovery` tokens over `.bench/`, `.agents/`, `projects/`, `README.md`, `CHANGELOG.md`, `ROADMAP.md`, and staged specs — is green post-removal and stays as a standing check | conformance sweeps | observed red by construction 2026-08-11: the tokens appear across `.bench/BENCH.md`, `.bench/BENCH-reference.md:99-106`, the phase commands, and skills today, so the authored sweep is red until the prose repair lands | the existing sweeps match only `/bench-*` and top-level-verb tokens (`docs_workflow_checks_test.go:214-308`) and pass a leftover `bench spec build promote` |
| RM9 | 2 | interrupted cleanup converges: killing resume mid-reconcile and re-running leaves the two lifecycle ref namespaces empty and every surviving ledger entry worktree-pool-typed with an existing worktree path — no entry names a specbuild assignment | resume reconcile | not TDD-able before the reconcile exists; the implementing ticket authors the interruption test first at the new seam | partial deletion plus a crash is the one path back to phantom provisional state |
| RM11 | 2 | every kept verb — `worktree create/path/exec/release/clean`, `gate`, `commit`, `status`, `guards`, `idea`, `roadmap`, `spec implemented/retire/history` — still answers `--help` with exit 0, pinned as an enumerated keep-list check | CLI dispatcher | cannot start red — a protective pin over surface that must survive; its red is demonstrated during authoring by temporarily dropping one kept route and observing the failure | the routing registry (`subcommand_routing_test.go:57-105`) goes silently green when a route is deleted, so over-deletion needs its own enumerated guard |
| RM12 | 5 | `CHANGELOG.md` carries the one removal entry naming the deleted `spec build` family and `worktree recovery` | removed-verb sweep / changelog | observed red 2026-08-11: no such entry exists in `CHANGELOG.md` | map #12 allows exactly one migration surface — the changelog line — and nothing else advertises the break |

### Edge inventory

- Error path — RM1, RM4 (structured errors on removed grammar).
- Empty/absent input — resume over an already-clean tree: RM5's no-op half;
  empty ledger already covered by existing resume tests, which run the same
  reconcile entry the change edits.
- Boundary values — **Won't handle:** none apply; no numeric or size
  boundaries in scope.
- Malformed input — unknown-typed ledger entries: RM5 (purge granularity is
  per entry, enumerated as lifecycle-typed and unknown-typed).
- Interrupted/partial state — RM9.
- Re-run idempotency — RM3, RM5 (second-run no-op).
- Process-boundary lifecycle — RM5 (state written by the old binary, read by
  the new one across the upgrade boundary).
- Hostile environment — arbitrary ref names inside the two lifecycle
  namespaces: RM3's any-name quantifier; refs outside them — including the
  gate's `refs/bench/green/<branch>` — are never touched (RM3, RM10).
- **Won't handle:** linked-repo backwards compatibility and migration
  tooling — map #12 records the zero-compat decision; removed verbs simply
  become unknown commands.
- **Won't handle:** rescoping the three parked specs (`axi-coherent-diff`,
  `axi-query-disclosure`, `single-build-serial-gate`) — they reference the
  retired spec's slug in prose, not removed commands; their rescope is their
  own future work per map #2.

## Ownership fences

- `internal/specbuild` (deletion), `cmd/bench`, `internal/commit`,
  `internal/coverage`, `internal/conformance`, `internal/intent`,
  `internal/worktree`, `internal/harness`, `internal/status`, `internal/spec`,
  `internal/anchors`, `internal/shift`, `internal/gate`
- `tests/canary/workflow-guidance-anchors` (spec-build fixture deletion)
- `bin/bench.sh`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`
- `.agents/commands/bench-implement-spec.md`,
  `.agents/commands/bench-final-check.md`,
  `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`,
  `.agents/skills/bench-craft-delegate/SKILL.md`
- `projects/benchkit.md`, `ROADMAP.md`, `README.md`, `CHANGELOG.md`,
  `capture/session-handoff.md`, `capture/learnings.md` (close the `--spec`
  entry)

## Out of scope

- **Spec B — `bench preflight`** (map #7): ~15 edits, 3 gate runs. Separate
  capability with its own future spec.
- **Spec C — doctrine adoption and prose diet** (map #4, #5, #6, #8, #9,
  #10): ~40 edits, 4 gate runs. Separate capability; this spec repairs
  references only.
- **Parked-spec rescopes** (map #2): ~6 edits, 1 gate run each, after the
  reshape.
