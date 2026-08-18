# Evidence log 02 — the conformance gap (headline finding)

## E6 — Bench's largest subsystem does not run in Bench's oracle (OBSERVED)

### What conformance is

`internal/conformance` is **13,250 LOC (12,838 of it test code)** — the single
largest package in the repo — implementing **29 registered checks**
(`internal/conformance/registry/registry.go:108-136`). These are the mechanism
that distinguishes Bench from a prose-only skill kit: they mechanically enforce
guidance invariants that Pocock's skills leave to model compliance. Examples:

```
docs-currency-workflow        skills-index-command-adapters   handoff-shape-single-source
guidance-prose-budgets        decision-map-integrity          roadmap-detail-integrity
bench-sh-routes               axi-query-registry              default-branch-single-source
```

### The gap

`bench gate`'s built-in phase table (`internal/gate/phases.go:83-87` + the system
and shellcheck phases) is exactly:

```
gofmt · vet · go test -count=1 ./... · race · system · shellcheck
```

There is **no conformance phase**. The 29 checks reach the live tree only through
`TestRootConformance` (`internal/conformance/gate_entry_test.go:16`), which opens:

```go
root := os.Getenv("BENCH_CONFORMANCE_ROOT")
if root == "" {
    capability.Environment(t, "BENCH_CONFORMANCE_ROOT not set")
}
```

`go test ./...` **matches** that test, so `[test] ok internal/conformance 7.087s`
appears in the gate output — but the test **skipped**, and what actually ran was
only the fixture-based conformance tests. The gate's own footer discloses it as
one line of aggregate noise:

```
capability-skips: 7 (capability=6 environment=1)
```

The single `environment=1` skip is the entire live-root conformance suite.

A repository-wide search finds exactly one setter of `BENCH_CONFORMANCE_ROOT`
outside test files:

```
internal/preprelease/preprelease.go:100:
    Env: []string{"BENCH_CONFORMANCE_ROOT=" + root, ConformanceTierEnv + "=ship"}
```

— i.e. **`bench prep-release`**, the release rehearsal. Between releases, every
invariant these 29 checks defend is unenforced.

### It is not by design

`checkGateEntryContract` (`gate_entry_test.go:53-74`) lists the shell-fragment
invocation as a **retired** fragment `.bench/gate.sh` must not contain — that is
about *how* it ran, not *whether*. And the roadmap records the exact trap twice,
in `roadmap/FT133.md`:

```
Occurrence: 2026-07-26 conformance review — `TestRootConformance` matched but
            skipped without `BENCH_CONFORMANCE_ROOT`.
Occurrence: FT126 recurrence — a scoped conformance command again omitted
            `BENCH_CONFORMANCE_ROOT` and printed `ok`.
```

FT133 is still an open MEDIUM row on a 73-row board.

### When it broke

```
72c037a1  2026-07-05  "Implement gate phase-level concurrency"
          removed from .bench/gate.sh:
          -if ! (cd "$realkit" && BENCH_CONFORMANCE_ROOT="$root" \
          -      go test -count=1 ./internal/conformance -run '^TestRootConformance$'); then
```

**43 days** before HEAD (2026-08-17). No replacement phase was added.

### What is red right now at HEAD

```
$ BENCH_CONFORMANCE_ROOT="$PWD" go test -count=1 ./internal/conformance -run '^TestRootConformance$'
--- FAIL: TestRootConformance (3.20s)
  gate: stale Codex adapter reference $bench-finalize-spec in decisions/spec-build-review-gate-cadence.md:123
  gate: stale Codex adapter reference $bench-finalize-spec in decisions/spec-build-review-gate-cadence.md:129
  gate: .agents/commands/bench-implement-spec.md missing Entry orientation
  gate: .agents/commands/bench-implement-spec.md missing Exit handoff
  gate: .agents/commands/bench-final-check.md dropped whole-file implementation-retro replacement
  gate: .bench/BENCH.md dropped the implementation-retro drain owner
  gate: .bench/BENCH-reference.md:106 carries removed command token "spec build"
  gate: internal/canary/canary.go does not consume bounds.CanaryInnerWidth
  gate: decisions/spec-build-review-gate-cadence.md: Sources Path internal/specbuild/checkpoint.go must name an existing regular file
  gate: decisions/spec-build-review-gate-cadence.md: Sources Path internal/specbuild/lifecycle.go must name an existing regular file
FAIL
```

**10 violations**, while `bench gate` reports green.

### The demonstration that matters

`.agents/commands/bench-implement-spec.md` currently carries these headings:

```
## Entry   ## Declare the line, validate the tickets, route the venue   ## Build
## Land    ## When the build stops short   ## `--full <spec>`
```

The contract (`internal/conformance/docs_workflow_checks_test.go:537-541`)
requires literal `^## Entry orientation$` and `^## Exit handoff$`. The headings
were lost in:

```
fa4e1f02  2026-08-11  "doctrine: slim craft-tickets to the tracer contract,
                       implement-spec to the pointer (ticket 04)"
```

That is **a Bench build, run through Bench's own pipeline, that violated a Bench
conformance contract and landed green — six days ago.** Every landing since has
inherited it.

This is the audit's central evidence for *"do not assume a green test proves the
requirement"*: the check exists, is well-written, is registered, has tests, and
is bypassed by an environment variable nobody sets.

### Secondary consequence

`decisions/spec-build-review-gate-cadence.md` names `internal/specbuild/checkpoint.go`
and `internal/specbuild/lifecycle.go`. `internal/specbuild` was deleted by
`dae240df` (2026-08-11, "remove spec-build lifecycle core"). `bench maps` exits 1
on it; the gate is green; nothing routes anyone to the failure. Six days stale.

## E7 — three unreconciled verdict tiers

An agent in this repo faces three "red" surfaces with no stated precedence:

| surface | at HEAD | enforced by | reaches the agent |
|---|---|---|---|
| `bench gate` | **green** | invariant 1, Stop hook, `bench commit` | yes — authoritative |
| `bench status` gate row | **red** (stale store, E5) | nothing | yes — SessionStart headline |
| `bench structure` | **red**, 62 issues | nothing (not a gate phase) | yes — `bench status` row |
| `bench maps` | **red**, 1 invalid map | nothing | yes — `bench status` row |
| root conformance | **red**, 10 checks | `bench prep-release` only | **no** — silent skip |

Four of the five surfaces an agent can see are red; the one that decides "done"
is green; and the one carrying the most substantive violations is invisible.
`bench structure`'s 62 issues include 4 rows of rot in its own suppression file
(`structure-accept: stale accept row (not a scanned source file):
internal/contract/runtime/...`) — a suppression list pointing at deleted paths.
