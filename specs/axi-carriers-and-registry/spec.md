# AXI carriers and registry

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

Bench advertises ten AXI principles through fragmented prose and command-local mechanics. There is no typed owner for outcome, action, projection, aggregate, or empty facts, and the executable registry cannot declare which public members implement which mechanics.

## Solution

After `axi-compatibility-oracle` is implemented, introduce `internal/axi` as the shared carrier and validation owner, enrich the root and nested production registries with inert AXI declarations, and document the complete ten-principle per-surface contract. This expansion adds the new form beside existing derivations; it migrates no production domain and changes no public bytes.

## User stories

1. As a command-domain contributor, I want small typed carriers that preserve owner-defined meaning, so later migrations share mechanics without normalizing semantics. Line: gpt-5.6-terra / high. The interfaces are known, but a shallow universal wrapper can remain syntactically valid and semantically wrong.
2. As a kit contributor, I want one executable declaration of all command members and their AXI disposition, so membership, renderer, grammar, detail, and empty facts cannot drift across catalogs. Line: gpt-5.6-terra / high. Registry validation is gate-observable but spans root, wrapper attachment, and nested-family composition.
3. As an agent author, I want the canonical CLI guidance to state all ten principles and the exact approved surface, so future generation does not widen AXI to operational commands. Line: gpt-5.6-sol / high. Generation-guiding kit prose takes the leverage override.

## Implementation decisions

- `internal/axi` owns generic outcomes with domain-supplied kind/exit policy, typed executable actions with fixed and open arguments, owner-supplied bounded projections, ordered typed aggregates, and explicit empty classifications. It derives no domain meaning and emits no public envelope.
- The production root registry declares all 48 members; nested family registries compose the same declaration type. Required metadata is attachment, AXI disposition, grammar/help owner, renderer family, empty class, default/detail modes, and applicable shared routes.
- The exact approved AXI set remains six root queries—`anchors`, `learnings`, `maps`, `guards`, `diff`, `coverage`—plus nested `worktree list`. All other members are explicitly classified and retain their separately reviewed contracts.
- Registry metadata is validated and observable to conformance but inert to command dispatch and rendering. Independent expectations make flipping an approved member or operational exemption red.
- `craft-cli` states all ten principles and their current per-surface application. It does not advertise `--fields`, rendered `help[]`, content-first operational homes, or any public behavior deferred to later specs.

## Testing decisions

- TDD attaches at the pure carrier interfaces with minimal domain policies and literal ordered facts.
- Registry tests enumerate real production roots and nested grammars, reject missing/duplicate/unclassified entries, and compare output before/after inert metadata mutations.
- Conformance derives the exact AXI set independently and keeps its project-profile advertisement in the same ticket that registers the check.
- The implemented `axi-compatibility-oracle` is the no-public-delta seam; promotion remains the sole full gate.

### Seam diagram

    trigger: a domain declares a result or a registry declares a command member
        │
        ▼
    owner facts ──▶ [ internal/axi carriers + registry validation ] ──▶ inert typed declaration
                            ◀ tests attach here: policy tables and complete production enumeration

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CR1 | 1 | Outcomes preserve domain-owned kinds, payloads, and exact exit policies, including AXI 0/1/2 and specialized operational mappings. | `internal/axi` outcome interface | observed red: `test -d internal/axi` exited 1 | Pure policy tables reject a directory stub and a universal exit map. |
| CR2 | 1 | Actions distinguish literal command tokens, fixed arguments, open placeholders, and non-invokable prose without rendering `help[]`. | `internal/axi` action interface | observed red: CR1's shared-owner probe exited 1 and production `help[]` emitters remain zero | Typed validation rejects flattened prose and premature public disclosure. |
| CR3 | 1 | Projections carry selected content plus owner-supplied total, emitted, omitted, truncated, completeness, and counting unit without inference. | `internal/axi` projection interface | observed red: CR1's probe exited 1 while four local policies remain | Counterexample tables reject a universal cap, unit, or visible-row total. |
| CR4 | 1 | Aggregates preserve owner names, order, scalar type, unknown values, and required zeros; empty classifications remain domain-specific. | `internal/axi` aggregate and empty interfaces | observed red: CR1's probe exited 1 | Ordered typed facts reject sorting, stringification, zero omission, and unknown coercion. |
| CR5 | 2 | Per production-registry member, one declaration enumerates all 48 roots and every declared nested member with complete required metadata derived from the root and nested executable owners. | root and nested production registries | observed red: `rg -n -e 'AXI' -e 'Schema' -e 'Empty' -e 'Help' cmd/bench/command_registry.go` exited 1 | Enumeration from executable owners rejects a parallel sample catalog. |
| CR6 | 2 | Validation rejects missing, duplicate, unclassified, or independently advertised members while metadata mutations change no command output. | registry validation and conformance | not TDD-able before CR5 defines the declaration; CR5 is the absence red | Mutations make incomplete membership red and the oracle rejects metadata leaking into bytes. |
| CR7 | 2 | The exact AXI set is six approved root queries plus `worktree list`; every other member is an explicit exemption. | independent disposition expectation | observed red: current scope exists only in guidance/profile prose | Flipping either direction fails even when output remains byte-identical. |
| CR8 | 3 | Canonical CLI guidance states all ten principles, the exact per-surface scope, and current deferrals without claiming operational commands conform automatically. | `craft-cli` and project-profile guidance | observed red: `rg -n '8\. \*\*Content first' .agents/skills/bench-craft-cli/SKILL.md` exited 1 | The present seven-principle document and any whole-binary widening fail the independent scope check. |

### Ticket derivation

Every mapped row becomes ticket acceptance, atomic closure, and subject mutation under the approved fence. `already covered` never substitutes the old positive control for a shared-route mutation, and `not TDD-able` rows begin red immediately after their named prerequisite exists. Only an unforeseen local behavior may use `(covers local)`.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| CR1 | validate one domain policy / kind / payload / exact exit | `internal/axi` | replace a specialized exit map with universal 0/1/2 and admit an undeclared kind | pure outcome policy test; construct each policy and require exact validation |
| CR2 | validate one executable action / literal token / fixed argument / open placeholder / prose disposition | `internal/axi` | flatten fixed/open arguments into prose and mark orchestration prose executable | pure action test; construct each class and require typed round-trip or refusal |
| CR3 | preserve one owner projection / selected content / total / emitted / omitted / truncated / completeness / unit | `internal/axi` | infer total from emitted content or normalize byte and code-point units | projection counterexample test; construct owner facts and require exact typed values |
| CR4 | preserve one ordered aggregate and empty declaration / name / order / scalar type / unknown / zero / empty class | `internal/axi` | sort facts, stringify a number, coerce unknown to zero, omit zero, or default empty | aggregate/empty tests; validate and enumerate exact supplied facts |
| CR5 | declare every production root and nested member / attachment / disposition / grammar owner / renderer / empty / detail / routes | `cmd/bench`, `internal/usage`, `internal/spec`, `internal/publication`, `internal/preflight`, `internal/gate`, `internal/harness` | delete one derived member or clear one required metadata field | registry enumeration test; compare real executable owners to declarations per member |
| CR6 | reject missing / duplicate / unclassified / independently advertised / byte-active metadata | `cmd/bench`, `internal/conformance`, `projects/benchkit.md` | duplicate one registry name, add one parallel advertisement, then route metadata into output | registry/conformance and compatibility tests; validate declarations and compare the affected command observation |
| CR7 | retain six approved roots / nested worktree list / every explicit exemption | `cmd/bench`, `internal/conformance`, `projects/benchkit.md` | flip one approved query off and one operational member on | independent exact-set test; derive production dispositions and require the approved set |
| CR8 | publish ten principles / per-surface scope / approved set / deferrals / changelog | `.agents/skills/bench-craft-cli/SKILL.md`, `projects/benchkit.md`, `CHANGELOG.md`, `internal/conformance` | remove principle 8, widen AXI to all commands, or advertise `--fields`/`help[]` | docs-currency conformance; compare guidance claims to production declarations and current emitters |

### Edge inventory

- Error path — CR1 and CR6 reject impossible policies and malformed declarations.
- Empty or absent input — CR4 distinguishes explicit empty classes, unknown, and no-semantic-empty.
- Boundary values — CR3 preserves zero, equal, and emitted-less-than-total facts under distinct units.
- Malformed input — CR2 rejects invalid actions; CR6 rejects incomplete and duplicate declarations.
- Interrupted or partial state — CR1 and CR4 retain owner-supplied partial/unknown facts; no process owner is added.
- Re-run idempotency — validation is pure and repeated metadata changes remain byte-inert.
- Process-boundary lifecycle — the prerequisite oracle runs candidate executable cases; these carriers hold no process-local authority.
- Hostile environment — registry enumeration is independent of cwd and PATH; the prerequisite oracle retains wrapper hostile cases.
- Command self-observation — registry metadata cannot alter the commands used to validate its inertness.
- Special files and dangling symlinks — no new discovered-path reader is introduced.

## Out of scope

- Production outcome/action routing — 16 edits, 1 promotion gate run in `axi-outcome-action-migration`.
- Production bounded projections — 8 edits, 1 promotion gate run in `axi-bounded-projection-migration`.
- Production aggregates and empty routing — 14 edits, 1 promotion gate run in `axi-aggregate-empty-migration`.
- Public `--fields`, `help[]`, content-first operational homes, schema versions, or sink changes.

## Ownership fences

- Shared owner: `internal/axi/**` excluding `internal/axi/compatibility/**`.
- Registry declaration: `cmd/bench/main.go`, `cmd/bench/command_registry.go`, `internal/usage/**`, `internal/spec/**`, `internal/publication/**`, `internal/preflight/**`, `internal/gate/**`, `internal/harness/**`.
- Conformance and advertisements: `internal/conformance/**`, `.agents/skills/bench-craft-cli/SKILL.md`, `projects/benchkit.md`, `CHANGELOG.md`.
- The compatibility oracle and all production domain renderers remain unchanged inputs.
