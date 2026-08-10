# AXI query disclosure

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

After the spec-build and diff migrations, the remaining approved query surfaces — `anchors`, `learnings`, `maps`, `guards`, `coverage`, and `worktree list` — still end at their facts: a stale guard names no repair, an invalid map names no template, an unchecked coverage map names no check command. Help spellings are inconsistent (`worktree list` hand-parses `-h/--help` and rejects bare `help`), no gate check derives AXI conformance for the approved set as a set, and `craft-cli` still documents seven of the ten principles, so guidance cannot ask for what it never names.

## Solution

Complete contextual disclosure across the remaining approved surfaces as additive transitions: each surface keeps its primary bytes, streams, exits, and argv exactly and appends one `help[]` block composed from the `axi.Action` owner, with the honest empty block for terminal results. Route `worktree list` through `usage.Parse` so the shared help spellings apply, close one registry-derived conformance check over the whole approved set, and rewrite `craft-cli` to state all ten principles and the exact approved surface. The spec's action inventory starts from `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md` and a harness-log leverage review taken across the FT173 builds, with every observed opportunity given one disposition.

## User stories

1. As an agent using the approved query surfaces, I want each response to append state-derived `help[]` — `learnings` with open entries names the drain phase, `maps` names the shape phase on a frontier and `bench maps --template` plus the named path on an invalid map, `guards` names `bench link` on stale or unwired rows, `coverage` names `bench coverage --check <spec>` with the resolved spec on an unchecked map and the exact retry after a named repair, `worktree list` names `bench worktree path <id>` for active rows and `bench worktree clean <path>` for orphaned rows — while `anchors`, empty results, and complete terminal results append the honest empty block. Line: gpt-5.6-terra / high. The per-surface action cases are enumerable but omission-sensitive, and inventing busywork on terminal results is the equal and opposite failure.
2. As an agent asking any approved surface for help, I want the shared help spellings (`--help`, `-h`, bare `help`) accepted uniformly, with `worktree list` routed through `usage.Parse` and its accepted-argv delta named by an old-to-new fixture. Line: gpt-5.6-terra / medium. The migration is mechanical but changes the accepted-argv language, so the fixture must name exactly the new spellings.
3. As a kit maintainer, I want one registry-derived conformance check that grades the complete approved AXI set — the six root queries, `worktree list`, and the spec-build family — for structured stdout, definitive empties, honest exits, help spellings, and the help envelope, so surface membership and conformance cannot drift apart. Line: gpt-5.6-terra / high. The check must derive membership independently so flipping a surface in or out of the set is red in both directions.
4. As an agent author, I want `craft-cli` to state all ten principles, the help-block contract, and the exact approved surface, so future generation neither omits contextual disclosure nor widens AXI to operational commands. Line: gpt-5.6-sol / high. Generation-guiding kit prose takes the leverage override.
5. As the reviewer, I want the harness-log leverage review recorded as an evidence asset in which every observed opportunity — a missing or wrong `help[]`, a repeated tool-call sequence one query could replace, an output shaping agents perform themselves — receives exactly one disposition: folded into this spec, already owned, declined with a reason, or routed to a named roadmap item. Line: gpt-5.6-terra / high. An undispositioned observation is scope that leaks or value that silently drops.

## Implementation decisions

- Every addition is an additive transition per the compatibility contract: primary bytes, streams, exits, and accepted argv are unchanged except where a named old-to-new fixture states the delta (`worktree list` help spellings; the appended help block on each surface).
- Help derivation composes `axi.Action`; each surface's command function derives its actions from the same typed facts it already renders, and no shared renderer infers domain semantics. Harness phases (`/bench-what-next`, `/bench-shape-idea`) use the action owner's harness-phase kind introduced by `axi-spec-build-complete`, rendering in canonical form exactly as the CLI already prints phase names — never as shell commands and never by widening `internal/axi` here.
- Action derivation is per matching row, not per surface: every actionable row yields its own action with its row's values carried (each orphaned worktree row its own `bench worktree clean <path>`), exact duplicate templates dedupe, and row order is preserved.
- `anchors` remains terminal: success and empty results append the honest empty block. Empty-result help is the empty block everywhere — absence is the answer, not a prompt to invent work.
- `worktree list` moves onto `usage.Parse` with a declared grammar, satisfying the existing `subcommand-routing` posture without an exemption; no other accepted-argv change ships.
- The conformance check's production membership derives from the registry by the same AST route the existing `subcommand-routing` check uses; the approved-set expectation (the six root queries, `worktree list`, the spec-build family) is an independently authored conformance expectation — the AGENTS.md independence exception, necessary so flipping a member in either direction is red — and every other member is an explicit exemption. Per-member envelope grading for `cmd/bench`-owned surfaces lives in a package-main conformance test beside the registry. `projects/benchkit.md`'s AXI seam advertisement is updated in the same ticket.
- `craft-cli` states the ten principles with the per-surface application, the `help[N]{command}` contract including the honest empty form, and the approved surface; it does not advertise `--fields` or any behavior this sequence does not ship.
- The leverage review is the build's first ticket and a hard checkpoint: it runs before any QD1–QD4 implementation ticket, reads the named corpora (the Claude and Codex harness logs accumulated across the FT173 builds and representative recent Bench work), and its complete disposition ledger receives reviewer sign-off. A fold that would change locked coverage is a reviewer-approved spec amendment before implementation proceeds — never a silent build-time widening; the other dispositions are an "already owned" citation, a decline with reason, or a `bench idea` capture routed to the roadmap. An empty opportunity table is a valid outcome only when the corpus ledger names what was inspected and which sources were unavailable.

## Testing decisions

- TDD attaches per surface at the existing command seams with real fixtures (a stale guard manifest, an invalid map, an unchecked coverage map, an orphaned worktree row); each surface's action cases and its honest-empty cases get independent assertions.
- The `worktree list` grammar migration is proved by an accepted/rejected-argv paired delta: old and new spellings, streams, and exits, with only the named new spellings differing.
- Conformance tests mutate membership in both directions — drop an approved surface, add an operational command — and require red; envelope grading reuses the real command outputs.
- Mutation discipline for actions follows the owner tests landed by `axi-spec-build-complete`; each surface adds only its derivation cases.

### Seam diagram

    trigger: an approved query surface renders its result
        │
        ▼
    surface-owned facts ──▶ [ per-surface action derivation → axi renderer ] ──▶ existing response + help[]
                                  ◀ tests attach here: per-surface fixtures and set-membership conformance

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| QD1 | 1 | Per approved surface — `learnings`, `maps`, `guards`, `coverage`, `worktree list` — and per actionable row within a surface, each named state appends its state-derived actions with that row's values carried (resolved spec path, worktree id, clean path, diagnostic path), duplicates deduped and row order preserved; `anchors` plus every empty or complete terminal result appends the honest empty block. | per-surface derivation at the existing command seams | observed red: `rg -n 'help\[' cmd internal` exited 1 at authoring time | Zero emitters on these surfaces; per-row enumeration with independent fixtures rejects a first-match-only derivation, a sampled subset, and invented busywork on terminal results. |
| QD2 | 2 | `worktree list` routes through `usage.Parse` with the shared help spellings; the accepted-argv delta is exactly the named new spellings and everything else is byte-equal. | `worktree list` grammar | observed red: `rg -n 'usage\.Parse' internal/worktree/list.go` exited 1 | The hand-rolled parser rejects bare `help` today; the paired delta proves parity everywhere except the named spellings. |
| QD3 | 3 | One conformance check derives the approved set independently and grades each member's structured stdout, definitive empty, exit taxonomy, help spellings, and help envelope; dropping an approved member or admitting an operational one is red. | registry-derived conformance | observed red: `rg -n 'axi-disclosure' internal/conformance` exited 1 | Independent set derivation makes membership drift observable in both directions instead of resting on the profile's prose. |
| QD4 | 4 | `craft-cli` states all ten principles, the `help[N]{command}` contract with its honest empty form, and the exact approved surface, without advertising unshipped behavior. | `craft-cli` guidance | observed red: `rg -n 'Content first' .agents/skills/bench-craft-cli/SKILL.md` exited 1 | The present document names seven principles; the conformance cross-check rejects guidance that claims more or less than production declares. |
| QD5 | 5 | The leverage review runs as the build's first ticket, before any QD1–QD4 implementation ticket; its asset names its exact corpora — log directories, transcript counts, session-id ranges, every unavailable source named rather than substituted — carries one disposition per observed opportunity (fold, already owned, decline with reason, or named roadmap route), and its complete ledger receives reviewer sign-off, with any fold landing as a reviewer-approved spec amendment before implementation proceeds. | first-ticket checkpoint and its evidence asset | observed red: `rg -l 'ft173-log-leverage' decisions` exited 1 | The filename probe is only the absence red; acceptance is the reviewer-signed ledger with its corpus census, so an empty correctly named file, a disposition-free row, or a silent build-time coverage widening each fails the checkpoint. |

### Ticket derivation

Every mapped row becomes ticket acceptance with `(covers <row>)`, atomic `Closure:` tokens, and a subject mutation under the approved fence. QD1 may split per surface by repeating its covers ID. Only an unforeseen local behavior may use `(covers local)`.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| QD1 | append state-derived help per surface and per actionable row / surface / state case / row multiplicity / dedupe / carried value / honest empty | `cmd/bench/main.go`, `internal/learnings`, `internal/maps`, `internal/guards`, `internal/coverage`, `internal/worktree` | drop one surface's action case, emit only the first matching row's action, guess a value, or advertise an action on a terminal result | per-surface tests; drive each named state fixture (including a many-actionable-rows fixture) through the public command |
| QD2 | accept shared help spellings on worktree list / `--help` / `-h` / bare `help` / rejected argv parity | `internal/worktree`, `internal/usage` | accept an argv the old parser rejected outside the named spellings | paired argv delta; run old and new grammars over the accepted/rejected matrix |
| QD3 | derive and grade the approved set / AST-derived membership / independent expectation / envelope / exemptions | `internal/conformance`, `cmd/bench/command_registry_test.go`, `projects/benchkit.md` | flip one approved member off and one operational member on | conformance check plus package-main registry test; derive the set and grade real command outputs |
| QD4 | publish ten principles / per-surface scope / help contract / no unshipped advertisement | `.agents/skills/bench-craft-cli/SKILL.md`, `CHANGELOG.md` | remove one principle or advertise `--fields` | docs conformance; compare guidance claims to production declarations |
| QD5 | record named corpora and one disposition per observed opportunity / corpus ledger / unavailable sources named / fold / owned / decline / route | `decisions/byte-preserving-axi-foundation/assets`, `capture/IDEAS.md` | drop the corpus ledger or leave one observation undispositioned | reviewer check of the asset against its corpus ledger and completeness rule |

### Edge inventory

- Error path — every surface's existing structured errors are unchanged; refusal help rows are added only where the surface already knows the remedy (QD1's named cases).
- Empty or absent input — the honest empty help block on every empty result; absent journal, no maps, and clean guards keep their existing definitive empties.
- Boundary values — surfaces with one row versus many derive identical action shapes; help blocks are derivation-bounded.
- Malformed input — invalid maps and malformed coverage rows keep their exit-1 diagnostics and gain only the named repair-then-retry action.
- Interrupted or partial state — `guards`' incomplete scan keeps its honest aggregate; its help names `bench link` only for the stale/unwired class, never for timeout incompleteness.
- Re-run idempotency — appended help derives from the same facts as the response, so repeated calls return identical bytes.
- Process-boundary lifecycle — fixtures run the public commands; conformance grades rendered output, not internal state.
- Hostile environment — control-bearing values in help commands render through the existing escaping rules; deep-cwd behavior per surface is unchanged and graded by QD3.
- Command self-observation — help derivation cannot mutate the facts it reads.
- Special files and dangling symlinks — no new discovered-path reader; existing surface refusals are unchanged.

**Won't handle:** contextual disclosure for operational families outside the approved set (status, handoff, worktree mutations, release, adoption, plumbing) — each needs its own evidence and reviewer decision; a routed roadmap item from QD5 is the only path in.

## Out of scope

- The spec-build family envelope and the action owner — `axi-spec-build-complete` (prerequisite).
- The coherent diff snapshot and its actions — `axi-coherent-diff` (prerequisite for schema stability, ~16 edits, 1 promotion gate).
- `--fields`, default-schema width changes, family-home changes on any surface in this spec, or rewriting the already-migrated spec-build and diff contracts.
- Widening AXI to any operational family.

## Ownership fences

- Surface derivations: `cmd/bench/main.go`, `internal/learnings/`, `internal/maps/`, `internal/guards/`, `internal/coverage/`, `internal/worktree/`, `internal/usage/`.
- Conformance and advertisement: `internal/conformance/`, `cmd/bench/command_registry_test.go`, `projects/benchkit.md`.
- Guidance: `.agents/skills/bench-craft-cli/SKILL.md`, `CHANGELOG.md`.
- Evidence: `decisions/byte-preserving-axi-foundation/assets/`, `capture/IDEAS.md` (routed captures only).
- `internal/axi`, `internal/specbuild`, and `internal/diff` are consumed inputs and remain unchanged.
