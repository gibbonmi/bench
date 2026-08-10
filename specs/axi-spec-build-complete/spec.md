# AXI spec build complete

Status: implemented

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

`bench spec build` is the highest-frequency observed agent surface and already satisfies AXI on TOON output, ordered grammar, definitive empty status, and honest usage/refusal exits. It is not AXI-complete: no operation emits contextual `help[]`; the single `next` cell mixes invokable commands with orchestration prose (`release assignment <id>`, `retry promote`) an agent may try to execute; refusal remedies exist only inside `fmt.Errorf` prose, so fixed arguments the service already knows cannot be recovered by a renderer; and the family has no content-first home — `bench spec build` with no operation exits 2.

## Solution

One atomic output migration of the complete nine-operation family to the full ten-principle envelope, built directly on the existing `ParseBuild` → `buildService` → renderer seams. Introduce `internal/axi` as the typed-action owner — actions distinguish literal command tokens, fixed arguments, and open placeholders, and refuse non-invokable prose — plus one `help[]` renderer appended after each operation's primary response. Actions derive from the same typed lifecycle facts `record.status` already computes. The family home becomes a live retained-run listing. Every intentional delta is named by a reviewed old-to-new fixture; everything outside the named deltas keeps its current bytes, protected by the existing exact tests.

## User stories

1. As an agent driving a spec build, I want every operation response to end with `help[]` rows derived from the same typed lifecycle facts as the result — carrying the known slug, assignment, ticket, receipt, and fingerprint values — so the next command is copy-paste exact rather than reconstructed from prose. Line: gpt-5.6-terra / high. Omission and placeholder-guessing across the operation/state matrix are the principal risks.
2. As an agent whose operation is refused, I want the refusal's remedy derived from the typed precondition result the service already knows, so the advertised command can actually satisfy the observed state (the FT164 wrong-remedy class). Line: gpt-5.6-terra / high. Remedies today are flattened into error prose; recovering them by string inspection would advertise wrong or stale commands.
3. As an agent orienting cold, I want `bench spec build` with no operation to render a live content-first family home listing every retained run with a definitive empty state, so orientation costs one call instead of a usage error plus discovery. Line: gpt-5.6-terra / medium. The projection is small but replaces a pinned usage/2 result and needs its exact old-to-new fixture.
4. As an agent running plan/apply operations, I want each abandon/reclaim plan to emit the one authorized apply command carrying the target and fingerprint, and a stale-fingerprint refusal to emit the re-plan command, so two-call safety never costs command reconstruction. Line: gpt-5.6-terra / high. Fingerprint carry-forward is authority-bearing; a guessed or stale value must be impossible.

## Implementation decisions

- `internal/axi` owns `axi.Action` and the `help[]` renderer. An action is a literal template of one of two kinds — a shell command or a harness phase in its canonical `/bench-*` form — never prose: fixed arguments carry known values, genuinely unknown future input renders as an open `<placeholder>`, and non-invokable orchestration facts cannot construct an action. The renderer emits one `help[N]{command}` block appended after the primary response, and `help[0]{command}:` for a terminal result with no useful action. The single-field schema is deliberate: the primary response already carries the state.
- The action owns copy-paste-safe serialization of its own argv: a fixed value that needs quoting renders through the existing `sanitize.ShellQuote`, an open placeholder renders as the bare `<name>` token — the repository's existing usage-string convention, marking the command a template the agent completes — and is never quoted (quoting would pass it as a literal argument), and a value no quoted single-line form can carry (an embedded newline or other TOON-refused control) refuses action construction rather than emitting a corrupt template. A fixed value never renders as a placeholder, and TOON cell escaping is table encoding, not command escaping, never a substitute.
- The family home adds one bounded retained-run enumerator owned by `internal/specbuild` beside the existing named-slug reader: retained state files are digest-named, so the home enumerates the retained-state directory under a named entry cap with an honest at-least marker when exceeded. Record-suffix entries are records read fail-closed; the state owner's own lock-suffix companions are a named skipped class; any other entry, and any non-regular, malformed, or unreadable record, is its own named diagnostic row — never a silent skip and never an abort that hides the healthy rows.
- The operation/state matrix has one derivation: the operation axis is `buildOperationOrder`, consumed through an exported single-source accessor in `internal/spec` rather than a second list, and the state axis per operation is closed over the typed outcome and refusal constructors the service can return for that operation — success, no-op, empty, plan, apply, and every distinct typed refusal, including missing run, missing assignment, request conflict, stale or spent refresh receipt, invalid evidence receipt, moved candidate, capacity/cleanup, fence or ownership drift, dirty subject, incomplete journal, unreadable state, stale/spent fingerprint, and terminal-run. The build records the concrete table in the ticket ledger with a fixture or an exact not-applicable disposition per cell, and the SB7 conformance check derives the same axes, so a new constructor or an unclassified cell is red rather than unsampled.
- The family root accepts the shared help spellings: `bench spec build help`, `--help`, and `-h` return the operation catalog with exit 0 per the shared parser convention, alongside the live home; the home and each help spelling are named accepted-argv deltas in the old-to-new fixture set.
- The action owner and renderer land here because spec build is the first and widest consumer; `axi-coherent-diff` and `axi-query-disclosure` compose the same owner. `internal/axi` derives no domain meaning: every fact an action carries is supplied by the lifecycle owner.
- Action derivation extends `record.status`'s existing single derivation in `internal/specbuild`; no second next-step derivation appears at the render boundary. Orchestration-only states (release/delegate pending, gate-diagnosis dispositions) yield the family's inspection action (`bench spec build status <slug> --full`) instead of advertising prose as executable.
- The `next` cell becomes honest in the same atomic migration: it holds exactly the first help action's command when one exists and stays empty otherwise; non-command orchestration state stays visible through `state` and the `--full` projection. Its old-to-new fixture names this delta per lifecycle state.
- Refusal envelopes route through the same typed derivation: the service returns typed precondition facts (for example the exact-green evidence class) and the renderer emits the sanitized error line plus the help block, so the FT164-class remedy (`bench gate --fresh`) derives from the evidence class rather than embedded prose.
- `bench spec build` with no operation renders `spec_build_runs[N]{slug,state,next}` from retained state plus its help block — replacing the current missing-argument usage/2 result under its named fixture. Malformed operations and flag errors keep their usage/2 contract.
- The 0/1/2 exit taxonomy is unchanged: success/no-op 0, lifecycle refusal 1, usage 2. Promote's gate-result rendering is unchanged; FT185's gate payload is composed when it exists and is never re-derived here.
- Every operation and lifecycle state receives a reviewed old-to-new fixture pair naming exactly the appended `help[]` block, the `next`-cell change, and the family-home replacement; everything outside a named delta is byte-equal.

## Testing decisions

- TDD attaches at the `axi.Action` construction/rendering seam with literal counterexamples (prose refusal, placeholder rules, fixed-value carry) and at the specbuild render boundary through an in-process operation/state matrix driving the real service against bounded fixtures.
- Old-to-new fixtures live under `internal/specbuild` and `cmd/bench` testdata; each pins the old bytes, the new bytes, and the named delta, and a candidate change outside a delta turns the pair red.
- Action tests fail independently when a useful action is removed, a known value renders as a placeholder, an unknown value is guessed, a fixed flag is dropped, a stale fingerprint is carried, or prose is advertised as a command.
- One registry-derived conformance check enumerates the nine operations from `buildOperationOrder` and requires the envelope per operation; promotion remains the sole full gate.

### Seam diagram

    trigger: an operation returns typed lifecycle facts (success, empty, refusal, plan)
        │
        ▼
    lifecycle owner facts ──▶ [ axi.Action derivation + help renderer ] ──▶ existing TOON response + help[]
                                    ◀ tests attach here: operation/state matrix and action counterexamples

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| SB1 | 1 | `axi.Action` distinguishes shell-command and harness-phase templates, literal tokens, fixed arguments, and open placeholders, refuses non-invokable prose, serializes its argv copy-paste-safe (`sanitize.ShellQuote` for hostile values, construction refusal for values no single-line form carries), and the renderer emits `help[N]{command}` with `help[0]{command}:` as the honest terminal form. | `internal/axi` action owner and renderer | observed red: `test -d internal/axi` exited 1 | No typed action or help renderer exists anywhere in the tree, and TOON cell escaping cannot stand in for argv serialization. |
| SB2 | 1 | Per cell of the derived operation/state matrix — operations from the exported `buildOperationOrder` accessor, states closed over the service's typed outcome and refusal constructors, every cell covered or carrying an exact not-applicable disposition in the ticket ledger — the response appends `help[]` derived from the typed facts with every known slug, assignment, ticket, receipt, and fingerprint carried forward. | specbuild render boundary over the real service | observed red: `rg -n 'help\[' cmd internal` exited 1 | Zero production emitters means every applicable cell starts red; constructor-closed axes reject a build that covers one state per operation and calls the rest done. |
| SB3 | 2 | Per typed refusal constructor the service exposes — including exact-green evidence, stale or spent refresh receipt, invalid evidence receipt, moved candidate, missing run, missing assignment, request conflict, capacity/cleanup, fence or ownership drift, dirty subject, incomplete journal, unreadable state, stale/spent fingerprint, and terminal-run — the remedy derives from the typed precondition result and advertises the command that can satisfy that exact class. | service refusal envelope | not TDD-able until SB1 supplies the action owner; today the remedy exists only as `fmt.Errorf` prose (`internal/specbuild/assign.go:608`) | A remedy recovered by string inspection cannot carry typed fixed arguments, and closure over the constructors rejects a build that proves two examples and samples the rest. |
| SB4 | 3 | `bench spec build` with no operation renders `spec_build_runs[N]{slug,state,next}` over every retained run via the bounded enumerator — non-regular, foreign, malformed, and unreadable entries as named diagnostic rows, lock companions as a named skipped class, the entry cap with its at-least marker — a definitive `spec_build_runs[0]` empty, exit 0, plus its help block; the home and the shared help spellings (`help`, `--help`, `-h` returning the catalog at exit 0) each replace the former usage/2 result under a named fixture. | family home enumerator in `internal/specbuild` and the render boundary | observed red: `bench spec build` exited 2 and `rg -n 'spec_build_runs' internal cmd` exited 1 | The pinned no-argument diagnostic tests prove the current usage/2 result, and the hostile-entry rows make a silent skip or whole-home abort observable. |
| SB5 | 4 | Abandon and reclaim plans emit the one authorized apply command carrying slug and fingerprint; a stale or spent fingerprint refusal emits the exact re-plan command. | plan/apply renderer | not TDD-able until SB1; the plan tables today carry counts and refs but no executable template | Fingerprint carry-forward is authority-bearing: a guessed, stale, or dropped value must turn the derivation test red before it can reach an agent. |
| SB6 | 1,2,3,4 | Per operation and lifecycle state, one reviewed old-to-new fixture names the exact delta — appended `help[]`, the `next`-cell change, the family-home replacement — and everything outside a named delta is byte-equal between old and new fixtures. | paired old-to-new fixture suite | observed red: `test -e internal/specbuild/testdata` exited 1 | Without the paired fixtures an unnamed output change rides the migration unseen; the pair makes every delta a reviewed artifact. |
| SB7 | 1,2 | A registry-derived conformance check consumes the exported operation-order accessor, derives the same operation/state axes as SB2, and requires structured stdout, the 0/1/2 taxonomy, the shared help spellings at each operation and the family root, and the help envelope per operation; omitting one operation or leaving a cell unclassified is red. | `internal/conformance` over the exported `internal/spec` accessor | observed red: `rg -n 'axi-disclosure' internal/conformance` exited 1 | Derivation from the production single source rejects a hand-maintained second list that drifts when an operation or refusal constructor is added. |
| SB8 | 1,2,4 | Removing a useful action, rendering a known value as a placeholder, guessing an unknown value, dropping a fixed flag, carrying a stale fingerprint, advertising prose as a command, and emitting a hostile value unquoted each independently turn a focused owner test red. | action derivation tests | not TDD-able until SB1 supplies the owner | Each mutation class has a distinct silent-failure mode; one aggregate assertion would let six of the seven regress unseen. |

### Ticket derivation

Every mapped row becomes ticket acceptance with `(covers <row>)`, atomic `Closure:` tokens, and a subject mutation under the approved fence. A row may split per operation batch by repeating its covers ID. Only an unforeseen local behavior may use `(covers local)`.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| SB1 | construct and render one action / shell and phase kinds / literal token / fixed value / open placeholder / prose refusal / quoted hostile value / newline refusal / empty help | `internal/axi` | mark an orchestration label executable, flatten a fixed argument into prose, and emit a space-bearing value unquoted | pure action tests; construct each class and require typed round-trip, exact serialized argv, or refusal |
| SB2 | append derived help per operation and state / operation / state class / carried value | `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go` | drop help from one enumerated operation or render one known value as a placeholder | operation/state matrix tests; invoke each operation through the real service under bounded fixtures |
| SB3 | derive refusal remedy from typed precondition / evidence class / stale class / exact command | `internal/specbuild`, `cmd/bench/specbuild.go` | advertise a remedy the observed evidence class cannot satisfy | refusal envelope tests; drive each precondition class and require the class-exact action |
| SB4 | render the live family home / enumerated runs rows / diagnostic rows for non-regular, malformed, unreadable entries / definitive empty / exit 0 / help block | `internal/spec`, `internal/specbuild`, `cmd/bench/specbuild.go` | restore the usage/2 result, omit a retained run, or silently skip a malformed entry | family-home tests plus the named old-to-new fixture; run `bench spec build` with retained, empty, and hostile retained-state fixtures |
| SB5 | emit the authorized apply template / slug / fingerprint / stale re-plan | `internal/specbuild`, `cmd/bench/specbuild.go` | carry a stale fingerprint or drop the target from the template | plan/apply tests; plan, apply, and stale-apply fixtures for abandon and reclaim |
| SB6 | pin one old-to-new pair per operation and state / old bytes / new bytes / named delta | `internal/specbuild`, `cmd/bench` | change output outside a named delta | paired fixture suite; compare each operation's old and new fixtures and the delta ledger |
| SB7 | derive the operation/state axes and require the envelope / operation / state cell / stdout / exit / help spellings / help shape | `internal/conformance`, `internal/spec/build.go`, `projects/benchkit.md` | remove one operation from the derived enumeration or leave one cell unclassified | conformance check; consume the exported operation-order accessor and grade each operation |
| SB8 | prove seven mutation classes independently red / removal / placeholder / guess / dropped flag / stale fingerprint / prose / unquoted hostile value | `internal/axi`, `internal/specbuild` | apply each named mutation one at a time | focused derivation tests; one test per mutation class |

### Edge inventory

- Error path — SB3 covers every refusal class; usage and malformed argv keep their pinned usage/2 contract.
- Empty or absent input — SB4's definitive empty home; SB2 covers the no-run `state=empty` status and `--full` zero-row detail empties.
- Boundary values — SB2 covers single-run versus many-run projections; help blocks are bounded by derivation, never truncated.
- Malformed input — `ParseBuild`'s unknown-operation, missing-flag, and empty-value behavior is unchanged and remains pinned by the existing parser tests reaching the new home path.
- Interrupted or partial state — SB5 covers spent partial reclaim receipts; recovery states carry their retained-ref actions.
- Re-run idempotency — repeated status and home calls return identical bytes; no-op transitions keep exit 0 with their existing rows.
- Process-boundary lifecycle — matrix fixtures reload retained state in fresh service instances; the family home reads only durable retained state.
- Hostile environment — control bytes in slugs, paths, and error text keep the existing `sanitize.Controls`/`toon.Representable` behavior; help command values additionally serialize through the action's own argv quoting (SB1), since TOON cell escaping is not command escaping.
- Command self-observation — help derivation reads the same typed facts as the response and cannot mutate lifecycle state.
- Special files and dangling symlinks — the family home's retained-run enumerator is the one new discovered-path reader: it stats before reading, refuses non-regular and symlinked entries into named diagnostic rows, and reports malformed or unreadable records without hiding healthy rows (SB4); every other retained-state reader keeps its current refusals.

**Won't handle:** an ambient spec-build projection in `bench status` and any generated skill remain separate capabilities (roadmap-owned, not part of this family's envelope).

## Out of scope

- Coherent Git inspection — `axi-coherent-diff` (~16 edits, 1 promotion gate).
- Contextual disclosure on the remaining query surfaces, conformance closure over the whole approved set, and the guidance rewrite — `axi-query-disclosure` (~14 edits, 1 promotion gate).
- FT185's gate-result payload; this spec composes it when available and never re-derives it.
- Widening any other operational family to AXI, `--fields`, a legacy mode, or a dual renderer.

## Ownership fences

- Action owner: `internal/axi/`.
- Family migration: `internal/specbuild/`, `internal/spec/build.go`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`, `cmd/bench/testdata/`.
- Conformance and advertisement: `internal/conformance/`, `projects/benchkit.md`.
- Gate authority, lifecycle authority, `internal/gate/`, and every other command family remain unchanged inputs.
