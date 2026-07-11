# Roadmap context

Status: implemented

## Problem

`/bench-what-next` must currently rediscover roadmap, capture, spec, git, structure,
gate-cache, and promotion evidence through separate readers whose completeness and
failure postures differ. A fresh harness can therefore make a roadmap judgment from
different local evidence than another harness, and the phase has no single fail-closed
query it can trust.

## Solution

Add a fast, read-only `bench roadmap --context [--full]` AXI query that gathers every
local input `/bench-what-next` needs into one deterministic schema-1 TOON snapshot.
Each fact remains owned by its existing Go engine, the snapshot preserves malformed and
oversized evidence instead of silently dropping it, and the command performs no network
or cold verification work. Make the top-level wrapper prefer a distinct repo-local
wrapper so humans and every harness receive the repository's own CLI bytes, then make
the canonical phase consume this query exactly once and stop on query failure.

## User stories

1. As an agent maintaining the roadmap, I want one complete deterministic context
   snapshot, so that I can reconcile and drain from all relevant local evidence without
   follow-up discovery calls. Line: `gpt-5.6-luna` / low. The schema, seam, and gate
   assertions are exact and the public AXI contract observes the result.

2. As a kit maintainer, I want snapshot facts to come from typed owner APIs while
   existing renderers project from those same facts, so that adding context does not
   create a second parser or derivation that can drift. Line: `gpt-5.6-luna` / low.
   The ownership inventory fixes the known module shape and focused package tests cover
   each extracted fact.

3. As a harness or human caller, I want the top-level `bench` wrapper to select a
   distinct repo-local wrapper and preserve the full invocation, so that source,
   linked, global, deep-CWD, symlink, and worktree calls use the repository's CLI
   version without recursion. Line: `gpt-5.6-luna` / low. Wrapper behavior is exact,
   the seam is established, and the surface contract observes bytes, stderr, and exits.

4. As a `/bench-what-next` user, I want the canonical phase to invoke the context query
   once and stop if it fails, so that no harness reconstructs a partial snapshot or
   silently substitutes its own evidence. Line: `gpt-5.6-sol` / high. Editing shared
   command guidance has the leverage override because its semantics steer every future
   harness session and are only structurally gate-observable.

5. As a reviewer, I want the new owner, AXI, routing, and phase checks defended by
   targeted canaries, so that a future implementation cannot go green by dropping
   completeness, repo-local forwarding, or phase consumption. Line: `gpt-5.6-terra` /
   medium. The project caches gate and conformance logic at the mid line because an
   incorrect oracle would weaken every later shift.

## Implementation decisions

- `internal/roadmap` owns a typed roadmap `Document`, the deep `ContextSnapshot`
  aggregation, source-state and raw-fallback policy, truncation, deterministic ordering,
  and AXI rendering. It extends the existing roadmap seam rather than introducing a
  pass-through package.
- `internal/learnings`, `internal/structure`, `internal/spec`, and `internal/git`
  expose typed facts from their current engines. Their existing commands, dashboard
  rows, status signals, and counts become projections of those facts. No context code
  parses another command's stdout, and no fact gains two derivations.
- `status.GateVerdict` remains the sole gate-cache reader. Promotion-record parsing
  remains local to roadmap context until a second consumer earns its own seam.
- `bench roadmap --context [--full]` is a read-only, offline query. Success emits exit
  0 with empty stderr; an unreadable required source or failed derivation emits one
  structured AXI error on stdout with exit 1 and no partial snapshot; bad or conflicting
  arguments emit usage on stdout with exit 2; `-h` and `--help` emit usage on stdout
  with exit 0.
- Successful output uses the following fixed flat-table blocks and field order. Every
  block is present at zero rows, and the eight source rows are fixed and ordered as
  `ROADMAP.md`, `IDEAS.md`, `.bench/learnings.md`, `.bench/structure.budgets`,
  `.bench/structure-accept`, `specs/`, `CHANGELOG.md`, and the logical
  `.git/bench-last-gate` label.

  | block | fields |
  |---|---|
  | `context[1]` | `schema,full` |
  | `sources[N]` | `source,state,bytes` |
  | `roadmap_rows[N]` | `id,title,spec,spec_status,external_trigger,body,body_bytes,truncated` |
  | `roadmap_sequence[N]` | `rank,text,command` |
  | `ideas[N]` | `date,text,text_bytes,truncated` |
  | `learnings[N]` | `date,title,state,body,body_bytes,truncated` |
  | `structure[N]` | `kind,path,actual,limit,state,detail` |
  | `specs[N]` | `slug,status,roadmap_id` |
  | `spec_history[N]` | `slug,hash,date,kind,subject` |
  | `git[1]` | `branch,default_branch,dirty,ahead,behind` |
  | `git_changes[N]` | `status,path` |
  | `gate_cache[1]` | `present,status,cached_tree,work_tree,timestamp,stale` |
  | `promotion_records[N]` | `kind,date,scope,roadmap_ids,body,body_bytes,truncated` |
  | `parse_failures[N]` | `source,reason,raw,raw_bytes,truncated` |

- `context.schema` starts at integer `1`; `context.full` is boolean. Source state is
  exactly `absent`, `empty`, `parsed`, or `malformed`. Roadmap rows, sequence, ideas,
  learnings, and promotion records keep document order; structure, specs, and git
  changes use path order; history is newest-first per slug.
- Body-like fields (`body`, idea `text`, and parse-failure `raw`) expose at most 4096
  source bytes by default, cut back to a valid UTF-8 boundary. Their byte-count field
  always reports the full source size and `truncated` is true exactly when bytes were
  withheld. `--full` removes only this ceiling and sets `context.full=true`.
- Invalid or TOON-unrepresentable bytes fail closed rather than being replaced. An
  unchanged second invocation is byte-identical to the first.
- The top-level shell wrapper forwards the entire argv to a distinct repo-local wrapper
  before normal dispatch. It compares resolved paths to prevent self-recursion, reaches
  the main checkout's binary from linked worktrees, and falls back to its own existing
  resolution when no distinct local wrapper exists.
- The canonical `/bench-what-next` command invokes `bench roadmap --context` once at
  entry, treats it as the complete local evidence snapshot for reconcile and drain, and
  stops on a query error without manual reconstruction. Harness adapters gain no schema,
  parser, resolver ceremony, or separate query behavior.
- The implementation follows the handoff's green slices: expand typed owner APIs first;
  add aggregation and the AXI contract second; add wrapper forwarding and its surface
  contract third; then wire the canonical phase anchor and canaries. These slices share
  one branch but each preserves existing public behavior before the final integration.

## Testing decisions

- Tests exercise observable behavior through existing seams. Focused package tests are
  used only for typed facts and byte-boundary policy that are cheaper to isolate below
  the public command; no test parses a private implementation detail.
- The populated AXI fixture builder is the one source of the complete repository shape.
  Edge tests mutate that fixture instead of pasting multiple miniature context fixtures.
- The real command contract is the highest seam for completeness, output shape,
  failures, offline/read-only behavior, and determinism. The existing surface contract
  owns wrapper routing, and root conformance owns the command-phase anchor.
- A check is complete only after its targeted break has been observed red with its own
  message and the restored implementation is green. The three canary fixtures sabotage
  public surfaces rather than repeat schemas or derivations.
- The feature gate is `.bench/gate.sh`.

### Seam diagram

Typed owner and aggregation seam:

```text
trigger: focused Go package tests
    │
    ▼
fixture files/git/cache ──▶ [ typed owner APIs + ContextSnapshot ] ──▶ typed facts/errors
                              ◀ tests attach here: call public package APIs and compare facts
```

Public AXI command seam:

```text
trigger: agent or /bench-what-next invokes the real wrapper
    │
    ▼
repo state + argv ──▶ [ bench roadmap --context ] ──▶ TOON stdout + empty stderr + exit
                        ◀ tests attach here: fixture repo drives the real command
```

Wrapper routing seam:

```text
trigger: human, harness, hook, or headless shift invokes bench
    │
    ▼
cwd + PATH + argv ──▶ [ top-level wrapper routing ] ──▶ selected wrapper bytes/exits
                        ◀ tests attach here: source/linked/global/worktree probes
```

Canonical phase seam:

```text
trigger: harness loads /bench-what-next
    │
    ▼
canonical command text ──▶ [ root conformance anchor ] ──▶ one query + fail-closed rule
                             ◀ tests attach here: conformance reads the shipped command
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The real command emits all 14 fixed schema-1 blocks, their exact fields/order, the eight ordered source rows, and definitive empty tables. | Public AXI command | Observed 2026-07-10: `bin/bench.sh roadmap --context` returned exit 0 with the human roadmap and none of the required blocks. The implementation check is `go test -count=1 ./internal/contract/axi -run '^TestAXIRoadmapContextContracts$'`. | Ignoring `--context`, omitting a block, or hard-coding a partial snapshot cannot satisfy the complete public output assertion. |
| 1 | Roadmap, sequence, ideas, learnings, structure, specs/history, git changes, gate cache, and promotion records are lossless and deterministically ordered. | Public AXI command | The observed command returned only `ROADMAP.md`; the AXI completeness contract starts red for every other evidence class. | A cheapest implementation that wraps the current roadmap renderer lacks nine evidence classes and fails the populated fixture. |
| edge of 1 | Absent, present-empty, malformed, and newline-less sources remain distinguishable; malformed fragments appear in `parse_failures` with raw evidence. | Typed owner and public AXI seams | The current readers either collapse absent/empty or skip malformed content; the new focused parser cases and populated AXI malformed case must be observed red before implementation. | These cases reject both silent omission and a parser that normalizes distinct source states into one empty result. |
| edge of 1 | Default bodies obey the 4096-byte ceiling at 4095, 4096, and 4097 bytes, preserve full counts, cut at a UTF-8 boundary, and `--full` returns the complete bodies without changing schema/order. | Typed owner and public AXI seams | Observed 2026-07-10: `bin/bench.sh roadmap --context --full` produced the same human roadmap as the default call and no `context.full`, byte counts, or truncation fields. | A byte-slice off-by-one, rune split, silent omission, or `--full` schema fork fails an enumerated boundary case. |
| edge of 1 | Help/usage, unreadable sources, failed git derivations, missing git, and invalid or TOON-unrepresentable bytes produce structured stdout with honest 0/1/2 exits, empty stderr, and no partial snapshot. | Public AXI command | The current command accepts unknown `--context` arguments as exit 0 human output; the AXI error subtests must then be observed red per failure class. | A fail-open parser, stderr-only error, leaked partial snapshot, or dishonest exit cannot pass the external command assertions. |
| edge of 1 | Root and deep-CWD invocations are read-only and offline: repository files/modes/bytes, git status, index hash, HEAD, untracked files, and gate cache remain unchanged; sentinels for `bench`, `curl`, `wget`, `gh`, `glab`, `claude`, `codex`, `opencode`, and a planted gate are untouched. | Public AXI command | The read-only/offline contract is not TDD-able until the snapshot exists; its first implementation step must plant the enumerated sentinels and observe the real command fail on any attempted call or state delta. | Effect-based before/after assertions catch mutations and forbidden subprocesses without relying on an incomplete source scan. |
| edge of 1 | Two unchanged invocations emit byte-identical stdout. | Public AXI command | The populated AXI contract compares consecutive real-command runs; it must be added and observed red with the incomplete renderer before the command is implemented. | Any unstable map iteration, time-derived field, or nondeterministic ordering changes the bytes. |
| 2 | Existing roadmap/dashboard/status/learnings/structure/spec-history outputs are projections of the new typed facts and retain their established behavior. | Typed owner and aggregation seam | Focused owner tests are already the regression baseline; new typed-fact tests must first fail to compile or fail assertions before each extraction lands. | The combined old-output and new-fact assertions reject a second parser and catch compatibility regressions during extraction. |
| 2 | Promotion linkage is derived once inside roadmap aggregation, while every other snapshot fact comes from its named owner API. | Typed owner and aggregation seam | Not directly TDD-able as black-box behavior; repository review plus package dependency shape verifies this one-source design, while AXI fact cases verify its output. | This is an internal knowledge-ownership constraint; claiming a public red would overstate what the seam observes. |
| 3 | Source, linked by-path, global, symlink, linked-worktree, deep-CWD, stale-PATH, missing-local, and self-resolution cases select the correct wrapper and preserve exact stdout, stderr, and 0/1/2 exits. | Wrapper routing seam | The routing research observed a stale PATH `bench` return `stale-global`/exit 7 while the repo-local resolver returned canonical TOON/exit 0. The implementation check is `go test -count=1 ./internal/contract/surface -run '^TestGoRoutingContracts$'`. | A wrapper that stays PATH-owned, skips forwarding, or recurses cannot match the repo-local wrapper across the enumerated surfaces. |
| edge of 3 | `--context --full` remains two argv entries through forwarding, and a missing local wrapper falls back to the top-level wrapper's normal resolution. | Wrapper routing seam | The surface contract must first inject an argv-reporting local wrapper and be observed red against the current non-forwarding top-level wrapper. | Joining multiword argv or requiring a local wrapper changes the probe's bytes or exit. |
| 4 | The canonical phase invokes `bench roadmap --context` exactly once and states that query failure stops the phase without manual evidence reconstruction. | Canonical phase seam | Observed 2026-07-10: `rg -n 'bench roadmap --context' .agents/commands/bench-what-next.md` returned no matches. Root conformance must report `bench-what-next dropped the roadmap context query`. | Deleting the invocation or adding fallback gathering violates the anchored shipped command text. |
| 5 | Incomplete AXI context, skipped repo-local forwarding, and dropped phase consumption each have a distinct attributed gate failure and a targeted tripwire. | Gate canary seams | Not TDD-able until each owning check exists; after it lands, sabotage must be observed red as `AXI roadmap context completeness contract failed`, `repo-local wrapper forwarding contract failed`, and `bench-what-next dropped the roadmap context query` respectively. | The three distinct messages prove each public check bites without turning the canary into a second oracle. |
| edge of 5 | Canary fixtures join `behavior-owned` for AXI/routing and `workflow-guidance-anchors` for phase prose, and the real inner gate reports a fixture that stops biting. | Gate canary seams | Each new fixture must be run through `bench canary` once with its planted break and then through the restored green gate. | Family routing preserves attribution and catches an unwired or always-green owning check. |

Cheapest wrong implementations considered: leaving `RoadmapCommand` argument-blind,
wrapping only the current human roadmap renderer, reparsing other command stdout,
forwarding without resolved-path recursion protection, teaching only one harness a
resolver ceremony, anchoring prose without a behavior contract, and adding checks
without canary sabotage. The rows above turn each into a red or explicitly classify the
one internal ownership constraint that cannot honestly be black-box TDD.

### Edge inventory

The canonical error, empty/absent, boundary, malformed, partial, re-run, and hostile-
environment classes are covered above. The shell-CLI profile cases are covered as
follows: spaces/globs/control bytes and newline-less files at the AXI seam; absent versus
empty sources at the owner and AXI seams; multiword arguments, missing global tools,
symlinks, every shipped wrapper surface, and deep CWD at the routing seam; read-only
state and repeat-run identity at the AXI seam.

- **Won't handle:** recovery after SIGINT — the query creates no scratch state, lease,
  worktree, cache, or mutation, so forced termination can only discard stdout and the
  read-only assertion covers the persistent invariant.
- **Won't handle:** invalid UTF-8 or TOON-refused control bytes by replacement — replacing
  hostile evidence would make a lossless snapshot false, so the in-scope behavior is a
  structured fail-closed error.
- **Won't handle:** harness-specific output compatibility — all in-scope callers consume
  the same AXI/TOON contract, and the wrapper probes prove that portable calling
  convention directly.
- **Won't handle:** manual evidence reconstruction after a query error — the primary
  caller still exercises the feature on success, while failure deliberately stops the
  phase rather than amputating the successful interface.

## Out of scope

- **Roadmap mutation, prioritization, and drain verdicts** — these are a separate
  judgment capability retained by `/bench-what-next`, not evidence gathering (6 edits,
  3 gate runs).
- **Network or provider probes for external graduation triggers** — this is a separate
  provider-integration capability with authentication, latency, and availability
  decisions; context only marks rows that require it (8 edits, 3 gate runs).
- **Harness-specific schemas or rendering branches** — this would be a separate adapter
  compatibility capability and is intentionally unnecessary while every harness accepts
  the portable AXI contract (6 edits, 3 gate runs).
- **Manual fallback evidence gathering** — this would be a separate degraded-mode
  capability with its own completeness contract; the chosen posture is fail closed
  (5 edits, 3 gate runs).
- **Fresh gate, canary, or other cold verification execution** — this is a separate
  verification-runner capability whose cost and mutation surface conflict with the fast
  snapshot contract; context reports only the existing gate cache (5 edits, 3 gate runs).
