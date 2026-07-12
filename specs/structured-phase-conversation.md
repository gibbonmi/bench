# Structured Bench phase conversation

Status: implemented

## Problem

Bench tells agents to lead with the result, format for scanning, and stay concise,
but it does not give phase conversation a stable visual shape. Progress and exit
handoffs therefore vary by session: some arrive as dense prose, while others become
a stack of one-line bullets or headings that makes the answer feel choppy. The
reviewer has to reconstruct which text is the result, supporting detail, and next
action.

Duplicating a template into every phase command would replace inconsistency with a
new drift problem. Phase commands already own the content of their handoffs; the
shared communication rules are the one source that can govern how that content is
rendered across harnesses.

## Solution

Give Bench phase conversation two proportional patterns. A substantial in-progress
update uses compact bold **Status:** and **Next:** labels. A completed phase uses
`## Result`, `## Details`, and `## Next` headings. Empty groups disappear; related
sentences remain together; bullets and tables are reserved for facts that are truly
parallel rather than used as punctuation for every sentence.

Place the contract once in the shared communication rules. Individual phase commands
continue to define what their exit handoff contains, and the existing harness-native
invocation rule continues to govern the `Next` content. Root conformance and a
targeted canary prevent the shared rendering contract from disappearing, while a
fresh-session dogfood pass checks that the result reads coherently in real phase
conversation.

## User stories

1. As a reviewer following a substantial Bench phase run, I want progress grouped
   under compact **Status:** and **Next:** labels, so that I can see what is true now
   and what the agent will do without reading an unstructured update. Line:
   `gpt-5.6-sol` / high. Shared platform guidance steers every phase session that
   loads Bench's shared rules, so the leverage override routes this prose to the top
   tier.

2. As a reviewer receiving a Bench phase handoff, I want it led by `## Result` and
   organized with non-empty `## Details` and `## Next` sections as applicable, so
   that the outcome, evidence, and exact next action are predictable at a glance.
   Line: `gpt-5.6-sol` / high. The cross-phase output contract is semantic guidance
   with compounding impact and only structural automated coverage.

3. As a reviewer reading either form, I want related prose kept together and lists
   used only for parallel facts, so that stronger structure does not fragment the
   conversation into one-line sections or bullets. Line: `gpt-5.6-sol` / high. The
   cohesion rule is the judgment-heavy half of the shared communication contract and
   a superficially compliant but choppy rendering would pass a keyword-only check.

4. As a kit maintainer, I want root conformance and a permanent canary to guard the
   progress pattern, exit pattern, omit-empty rule, and cohesion rule, so that the
   shared contract cannot silently drift or vanish. Line: `gpt-5.6-terra` / medium.
   This uses the project’s known workflow-guidance conformance seam, whose oracle work
   is cached on the mid tier.

## Implementation decisions

The shared communication section owns rendering once for every Bench phase session
that loads Bench's shared rules. Phase commands remain responsible only for content
such as verdict counts, coverage state, or recommended commands. Harness adapters add
no output policy, and no phase command receives a copied template. The safe-link
contract still preserves project-owned bootstrap files; a session whose preserved
bootstrap does not import Bench's shared rules is outside this guidance contract.

A **substantial progress update** is one that reports a meaningful intermediate
result and continued work. It renders compact bold **Status:** and **Next:** labels,
each followed by cohesive prose. Routine one-sentence acknowledgements do not need a
template. A short update that carries meaningful state and continued work is still
substantial regardless of sentence count, and a phase that is finished uses the exit
form instead.

A phase exit renders:

- `## Result` first, stating the outcome directly;
- `## Details` only when supporting evidence or parallel facts materially help the
  reviewer understand or decide; and
- `## Next` when an action remains, naming the exact invocation this harness accepts
  under the existing command-translation rule.

Empty sections are omitted rather than printed as placeholders. A section may contain
a short paragraph, a compact list, or a small table, but related sentences stay in
one paragraph and lists or tables are used only for genuinely parallel facts. The
contract does not impose a machine-readable schema, require all three exit headings
for a one-line result, or turn every sentence into its own bullet.

The scope is conversational output produced while running an explicit Bench phase.
It does not restyle CLI or TOON output and does not govern repository artifacts such
as specs, decision maps, ADRs, reviews, or changelog entries. Ordinary conversation
outside a Bench phase continues to use the general communication rules without a
mandatory template.

Root conformance extends the existing workflow-guidance anchor checker with a named
clause declaration for progress labels, exit headings, empty-section omission, and
cohesive-prose/list restraint. The checker derives both the clause inventory and its
cardinality from the active communication guidance, preserving one source per fact,
and emits an attributable diagnostic when a declared clause is absent or empty. A
targeted workflow-guidance canary supplies a shared-rules fixture with the declared
progress clause removed; the existing fixture registry and bite inventory own its
discovery.

After a green gate, run fresh-session dogfood in three enumerated classes: one
tool-using phase that emits at least one substantial progress update, one compact
phase that reaches an exit with little supporting detail, and one short update that
contains both meaningful intermediate state and continued work. Confirm that the
first and third use the labeled progress form, the second omits an unnecessary
`Details` section, completed phases lead with their result, and any next command is
harness-native. Readability is a reviewer judgment; the gate only proves the active
shared instruction remains structurally present.

## Testing decisions

- A good automated test drives root conformance against the shared rules or a planted
  broken tree and observes one attributable diagnostic per missing clause class.
- The existing workflow-guidance anchor checker is the sole automated seam. It must
  not parse every phase command or create an adapter registry for output formatting.
- Real conversational quality is verified through fresh-session dogfood across the
  three representative phase classes. It is explicitly not a machine-parsed
  transcript contract.
- Gate command: `.bench/gate.sh`.

### Seam diagram

Automated shared-rule seam:

    trigger: root conformance grades the kit or a planted canary tree
        │
        ▼
    shared communication rules  ──▶  [ workflow-guidance anchor checker ]  ──▶  diagnostics
    broken shared-rules fixture ──▶  [                                  ]  ──▶  gate exit 0/1
                                       ◀ tests attach here: require the four rendering
                                         clauses; canary removes one and expects its message

Fresh-session conversation seam:

    trigger: reviewer invokes an explicit Bench phase
        │
        ▼
    phase state + handoff content  ──▶  [ shared communication contract ]  ──▶  progress / exit
    current harness invocation    ──▶  [                               ]  ──▶  exact next action
                                        ◀ tests attach here: fresh-session dogfood reads
                                          progress-heavy, compact-exit, and short
                                          substantial-update responses

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A substantial in-progress Bench phase update groups its meaningful current state and continued action under compact bold **Status:** and **Next:** labels. | active shared guidance plus fresh-session conversation | Observed red: `rg -q '\*\*Status:\*\*'` against the pre-change shared rules exits 1; conversational adherence is not TDD-able, so fresh-session dogfood supplies the behavioral verdict. | The red proves the positive instruction was absent; the active named clause prevents a comment or historical quotation from standing in for guidance, and dogfood catches an invented rendering. |
| 2 | A phase exit leads with `## Result`, uses `## Details` only for material support, and uses `## Next` for the exact remaining harness-native action. | active shared guidance plus fresh-session conversation | Observed red: `rg -q '## Result'` against the pre-change shared rules exits 1; conversational adherence is not TDD-able, so fresh-session dogfood supplies the behavioral verdict. | The active exit clause states the result-led pattern once; the compact-exit dogfood catches content-first or rigid three-section handoffs that structural conformance cannot grade. |
| 2 | Empty progress groups and exit sections are omitted rather than rendered as placeholders. | active shared guidance plus fresh-session conversation | Observed red: `rg -qi 'omit.*empty'` against the pre-change shared rules exits 1; proportional omission is not TDD-able, so fresh-session dogfood supplies the behavioral verdict. | A rigid always-three-sections template is the degenerate structured output; the compact-exit class distinguishes proportional structure from placeholder noise. |
| 3 | Related sentences stay together, and bullets or tables are reserved for genuinely parallel facts rather than one item per sentence. | active shared guidance plus fresh-session conversation | Observed red: `rg -qi 'related sentences.*together'` against the pre-change shared rules exits 1; readability is not TDD-able, so fresh-session dogfood supplies the behavioral verdict. | A named cohesion clause keeps the rule active and attributable; the readability pass catches choppy prose that any structural checker would miss. |
| 4 | The active communication guidance declares the progress, exit, omission, and cohesion clauses once; root conformance derives the inventory and cardinality from that declaration, diagnoses an absent or empty declared clause by name, and fails closed when the shared rules are missing or empty. | root conformance on the real tree and planted fixture | Red-first implementation: add the named-clause check before editing the shared rules and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; it must fail with an attributable diagnostic. | A source-derived inventory removes duplicated guidance and count knowledge while still making a missing declared clause attributable. Scoping the parser to active guidance rejects copies in comments, quotations, negating bullets, or other sections. |
| 4 | A workflow-guidance canary permanently proves the shared-output check bites. | canary fixture registry and fixture-bite inventory | Red-first implementation: register a planted shared-rules fixture without the progress contract and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; missing EXPECT/bite wiring must fail the fixture contract. | If the anchor is deleted or weakened until the planted rules pass, the known-broken fixture stops biting and turns the canary red. |
| edge of 1 | Fresh-session dogfood covers exactly three representative classes: a tool-using phase with a substantial progress update, a compact phase exit with no material details, and a short update containing both meaningful intermediate state and continued work. | real phase conversation in a fresh session | Not TDD-able: conversational cohesion is not a portable machine schema; record all three rendered responses and the reviewer readability verdict in the synthesis handoff. | The three classes catch under-structuring, over-structuring, and the boundary error that treats sentence count rather than informational weight as the exemption. |

**Degenerate-implementation check.** The cheapest wrong edit adds the three exit
headings everywhere, prints empty placeholders, and turns each sentence into a bullet.
Rows 2–4 plus the dogfood row reject that rendering. The cheapest wrong gate checks
one keyword anywhere in the file or hardcodes a second copy of the guidance; rows 5–6
reject it with an active source-derived declaration and canary.

### Edge inventory

- **Error path — phase stops short:** covered by the same exit pattern; `Result`
  states the stop, `Details` carries material cause/state, and `Next` gives the exact
  resume route. The phase-specific command still owns which facts are required.
- **Empty/absent details:** covered by story 2's omit-empty row; a compact exit may use
  only `Result` and `Next`, or only `Result` when no action remains.
- **Boundary values — one-sentence update vs substantial update:** covered by the
  third dogfood class. A routine one-sentence acknowledgement stays untemplated; a
  short update with meaningful intermediate state plus continued work uses
  **Status:** and **Next:** regardless of sentence count.
- **Malformed structure — headings present but prose fragmented:** covered by story
  3's cohesion row and the fresh-session readability pass.
- **Interrupted or partial phase:** covered by the phase's existing stopped-work
  content contract rendered through the shared exit form; no new persistent state is
  introduced.
- **Re-run idempotency:** conversational rendering mutates no state, so repeating a
  phase does not accumulate output artifacts. Repository idempotency remains owned by
  the phase itself.
- **Hostile environment — harness command vocabulary differs:** covered by the
  existing harness-native recommendation rule, consumed by `Next`; this spec does not
  duplicate its mapping.
- **Every shipped harness surface that loads Bench's shared rules:** covered through
  the single shared communication source. Thin command and skill adapters do not
  restyle output. A project-owned bootstrap that safe-link preserves without a shared
  rules import is explicitly outside the claim; changing that preservation contract
  would be a separate adoption decision.
- **Missing or empty shared rules:** covered by root conformance's required-file
  posture and the four new anchors.
- **Paths with spaces/globs, control bytes, missing trailing newline, unquoted
  arguments, missing CLI tools, symlink invocation, SIGINT, and cwd below root:**
  **Won't handle** — no shell, path, byte, process, or CLI input surface is added.

## Out of scope

- **CLI and TOON output schemas** — a separate agent-facing data-contract capability;
  conversational headings must not leak into stdout. Estimated later cost:
  `~6 edits, 3 gate runs` for any specific CLI redesign.
- **Repository artifact templates** — specs, maps, ADRs, reviews, and changelog entries
  keep their own full templates. Extending this rendering contract to artifacts would
  be a separate documentation-system capability. Estimated later cost: `~5 edits,
  2 gate runs` per artifact family.
- **Machine-parsed conversational transcripts** — a separate telemetry/schema
  capability requiring a stable event source rather than markdown styling. Estimated
  later cost: `~8 edits, 4 gate runs`.
- **Mandatory structured output for ordinary non-phase conversation** —
  reviewer-directed scope cut; the general communication rules remain proportional.
  Estimated reversal cost: `1 shared-rule edit, 1 gate run`.
- **Mandatory implementation delegation** — the first approved slice has its own spec
  and gate surface; this output slice neither duplicates nor changes it. Estimated
  work is owned by that sibling spec.
