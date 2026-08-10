# AXI compatibility oracle

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

Superseded by the 2026-08-10 forward-build ruling (see `ROADMAP.md` FT173): do
not assign, integrate, or promote. The active run resolves through `bench spec
build abandon`, then `reclaim`; this spec then retires. The four migration
specs this file references were deleted in the same restructure.

## Problem

Bench has exact-output tests for individual commands but no independent system that proves an internal AXI refactor preserves the complete public language. A migration can keep sampled stdout green while moving an error to stderr, changing an exit, accepting new argv, normalizing an empty state, or drifting a default/`--full` pair.

## Solution

Create one executable compatibility oracle before any AXI carrier or migration work. It freezes observations from subject `8ae1512f95e64588487430aefa5b02c288d7de3a` — the current production tree, since every commit since has touched only `specs/**`, so the pinned production bytes are the ones a candidate is built from — derives complete case membership from the current root, wrapper, and nested grammar owners, and compares a candidate through the production `Command` seam or one exact selected executable. The oracle is independently useful: every later byte-preserving migration can depend on it without carrying its own golden inventory.

## User stories

1. As a contributor, I want a provenance-checked baseline manifest whose expected observations cannot be refreshed from the candidate, so the oracle remains independent. Line: gpt-5.6-terra / high. The file and builder authority are exact but fail-closed special-file and preimage behavior need high-effort falsification.
2. As an agent using Bench, I want every current command and argv class compared on stdout, stderr, exit, and acceptance, so an internal migration cannot change the language I consume. Line: gpt-5.6-terra / high. The public grammar is distributed across root, wrapper, and nested owners and omission is the principal risk.

## Implementation decisions

- The compatibility module is a test-facing deep owner: manifest validation, case closure, observation capture, and delta reporting sit behind one interface. It does not render production command output.
- The comparator's exported identifier is fixed: `compatibility.Compare` in package `internal/axi/compatibility`, taking the authenticated baseline record and a candidate's paired fresh-state observations and returning the per-observation delta report. Dependent specs and tickets reference this literal symbol; renaming it is a spec change, not an implementation choice.
- The baseline records the pinned subject, canonical-builder seal, stable case ID, argv, fixture identity, raw stdout and stderr, exit, and accepted/rejected classification. The reader rejects an unpinned subject, a seal whose input preimage does not match, a non-regular or symlinked fixture, duplicate IDs, missing observations, and candidate-authored refreshes.
- Membership derives from all 48 Go roots, wrapper no-argument/help, `--version`, `-v`, wrapper-only `repair`, and the existing nested grammar owners. Every applicable required argv class is present or has a specific not-applicable reason.
- Ordinary cases use the production `Command` seam. Wrapper routing, process identity, environment, signal, or stream cases use baseline and candidate executables built once each through `scripts/go-build.sh`, with absolute paths and bounded process lifetimes.
- `specs/axi-compatibility-oracle/testdata` is the repository's first durable testdata directory: it holds the frozen baseline manifest as TOON records, which by design cannot be regenerated from the candidate, while hostile fixtures — FIFOs, symlinks, special files — are created at test time rather than stored.
- Census and membership tests live in `cmd/bench/axi_compatibility_test.go` (package `main`) because the 48-root `commandRegistry` enumerator is package-private in `cmd/bench/main.go`; `internal/axi/compatibility` never imports or re-enumerates it.
- Each case runs twice from fresh state. Raw equality includes ordering, quoting, controls, final newlines, exact empty forms, every default/`--full` pair, and the four existing bounds policies below/at/above their limits.

## Testing decisions

- TDD attaches at manifest validation and complete case derivation, using independently authored literal fixtures rather than candidate renderers.
- Paired execution attaches at `Command` and the exact executable; tests mutate candidate producers, never frozen expected bytes, to prove stdout, stderr, exit, acceptance, controls, and newlines all bite.
- The Go test phase is the gate owner. Later specs consume this oracle as a prerequisite; this spec changes no public bytes.

### Seam diagram

    trigger: compatibility test enumerates a production command and argv class
        │
        ▼
    frozen manifest + case ──▶ [ compatibility oracle ] ──▶ exact four-observation delta
                                      ◀ tests attach here: malformed manifests and mutated candidate producers

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CO1 | 1 | The manifest accepts only pinned subject `8ae1512f95e64588487430aefa5b02c288d7de3a`, the matching canonical-builder seal preimage, regular non-symlink fixtures, unique IDs, and complete four-observation records. | manifest reader | observed red: `rg -n -e 'baseline.*candidate' -e 'candidate.*baseline' internal cmd --glob '*_test.go'` exited 1 | No current reader can authenticate or reject a baseline fixture. |
| CO2 | 1 | Expected bytes are produced only by the pinned baseline executable and cannot be refreshed from the candidate under test. | baseline capture and candidate comparison | observed red: the paired-harness search in CO1 exited 1 | A candidate-authored golden fails provenance even when candidate and expected bytes agree. |
| CO3 | 2 | Per production member, case membership closes over all 48 roots, wrapper-only surfaces, every declared nested grammar member, and every applicable required argv class or exact not-applicable reason derived from `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`. | production registry, wrapper, and nested grammar census | observed red: the current evidence is prose-only in the FT173 assets | Derivation from all three production owners rejects a sampled or manually incomplete matrix. |
| CO4 | 2 | Every case compares raw stdout, raw stderr, exit, and accepted/rejected status twice from fresh state. | production `Command` and exact executable | observed red: no paired four-observation harness exists | Dropping any observation or reusing process state makes a targeted candidate mutation silently green. |
| CO5 | 2 | Every default/`--full` pair, empty class, bounds edge, quoting/control byte, and final newline remains exact. | exact-output matrix | not TDD-able until CO1 supplies the independent baseline; current command tests are positive controls | Producer mutations at each class make a byte delta without modifying the assertion. |
| CO6 | 2 | Root, wrapper, and nested help, malformed argv, `--`, no-op, refusal, outside-repository, hostile cwd/PATH, symlink invocation, and bounded interruption cases retain exact acceptance, streams, and exits. | exact executable and dispatcher matrix | not TDD-able until CO3 derives complete membership | The hostile cases exercise public behavior that in-process happy paths cannot observe. |

### Ticket derivation

Every row below becomes one ticket acceptance row with `(covers <row>)`. Its slash-separated facts become distinct `Closure:` tokens and distinct subject-mutation rows; the named owner runs the public operation. A ticket may split a row per enumerated member by repeating the same covers ID. None of these rows may become `(covers local)`.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| CO1 | load one authenticated regular manifest / subject / seal preimage / file kind / unique ID / four fields | `internal/axi/compatibility`, `specs/axi-compatibility-oracle/testdata` | drop one builder-seal preimage field; separately replace the regular fixture with a symlink and duplicate one ID | manifest-reader tests; load each hostile manifest under a named test timeout and require the exact refusal |
| CO2 | capture one baseline and compare one separately built candidate / baseline-only authorship / immutable expected bytes | `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata` | point expected capture at the candidate executable | paired-capture test; build both pinned subjects once, run the case, require provenance refusal before equality |
| CO3 | derive one case for every root, wrapper, nested member, and applicable class from the production owners / member / class / not-applicable reason | `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go` | omit one enumerated production member or one required class at a time | membership-closure test; enumerate real registries plus wrapper and compare to the manifest index |
| CO4 | compare stdout / stderr / exit / acceptance on two fresh executions | `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go` | move one candidate error between streams, change its exit, then widen one rejected argv spelling | paired-execution test; invoke `Command` or bounded exact processes twice and require four raw equalities |
| CO5 | preserve every helper-census default/full pair, empty class, cap edge, control byte, and final newline / pair / empty / bound / byte | `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata` | mutate the producing renderer to normalize one empty, cap, escape, or newline class | exact-matrix tests; run the derived case IDs and require raw stdout/stderr equality |
| CO6 | preserve every hostile argv/environment class `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md` declares / grammar / cwd / PATH / symlink / interruption | `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata` | accept one malformed argv, assume root cwd, trust ambient PATH, or leave a child past the named deadline | hostile exact-executable tests; run each derived case in a fresh bounded process and require exact observations plus teardown |

### Edge inventory

- Error path — CO4 and CO6 compare refusal bytes, both streams, and exits.
- Empty or absent input — CO5 covers zero-row, one-row, prose-clean, absent, and unreadable forms.
- Boundary values — CO5 covers every current cap below, at, above, and across multibyte boundaries.
- Malformed input — CO1 rejects malformed fixtures; CO6 covers malformed argv and empty required values.
- Interrupted or partial state — CO6 bounds every child, signal, and wait and compares partial observations.
- Re-run idempotency — CO4 runs every case twice from fresh state.
- Process-boundary lifecycle — CO4 and CO6 use fresh exact executables where serialization or wrapper identity matters.
- Hostile environment — CO6 covers spaces/globs in paths, deep cwd, stripped PATH, symlink invocation, missing optional tools, and control-bearing Git text.
- Command self-observation — baseline and candidate outputs are captured in separate immutable observations.
- Special files and dangling symlinks — CO1 rejects them before reading.

## Out of scope

- Shared AXI carriers and registry metadata — 8 edits, 1 promotion gate run in `axi-carriers-and-registry`.
- Outcome/action production migration — 16 edits, 1 promotion gate run in `axi-outcome-action-migration`.
- Bounded-projection production migration — 8 edits, 1 promotion gate run in `axi-bounded-projection-migration`.
- Aggregate/empty production migration — 14 edits, 1 promotion gate run in `axi-aggregate-empty-migration`.
- Any public schema, flag, `help[]`, family home, sink, or exit change.

## Ownership fences

- Oracle and fixtures: `internal/axi/compatibility/**`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/**`.
- Existing production registries, wrapper, builder, and command renderers are read-only inputs.
- `ROADMAP.md`, `capture/**`, every other spec, and all unrelated paths remain foreign.
