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
   `gpt-5.6-sol` / high. Shared platform guidance steers every future phase session,
   so the leverage override routes this prose to the top tier.

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

The shared communication section owns rendering once for every Bench phase. Phase
commands remain responsible only for content such as verdict counts, coverage state,
or recommended commands. Harness adapters add no output policy, and no phase command
receives a copied template.

A **substantial progress update** is one that reports a meaningful intermediate
result and continued work. It renders compact bold **Status:** and **Next:** labels,
each followed by cohesive prose. Routine one-sentence acknowledgements do not need a
template, and a phase that is finished uses the exit form instead.

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

Root conformance extends the existing workflow-guidance anchor checker with four
distinct diagnostics: missing progress labels; missing exit headings; missing
omit-empty behavior; and missing cohesive-prose/list restraint. The checker reads the
shared rules directly, preserving one source per fact. A targeted workflow-guidance
canary supplies a shared-rules fixture with the rendering contract removed and
expects the progress-contract diagnostic; the existing fixture registry and bite
inventory own its discovery.

After a green gate, run fresh-session dogfood in two enumerated classes: one
tool-using phase that emits at least one substantial progress update, and one compact
phase that reaches an exit with little supporting detail. Confirm that the first uses
the labeled progress form, the second omits an unnecessary `Details` section, both
lead with their result, and any next command is harness-native. Readability is a
reviewer judgment; the gate only proves the shared instruction remains present.

## Testing decisions

- A good automated test drives root conformance against the shared rules or a planted
  broken tree and observes one attributable diagnostic per missing clause class.
- The existing workflow-guidance anchor checker is the sole automated seam. It must
  not parse every phase command or create an adapter registry for output formatting.
- Real conversational quality is verified through fresh-session dogfood across the
  two representative phase classes. It is explicitly not a machine-parsed transcript
  contract.
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
                                          one progress-heavy and one compact phase response

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A substantial in-progress Bench phase update groups its meaningful current state and continued action under compact bold **Status:** and **Next:** labels. | workflow-guidance anchor checker | Observed red: `rg -q '\*\*Status:\*\*'` against the current shared rules exits 1; during TDD the new progress anchor must fail with its distinct diagnostic before the rule edit. | Unstructured prose and invented label names both lack the canonical pair, so the checker rejects the cheapest inconsistent rendering. |
| 2 | A phase exit leads with `## Result`, uses `## Details` only for material support, and uses `## Next` for the exact remaining harness-native action. | workflow-guidance anchor checker | Observed red: `rg -q '## Result'` against the current shared rules exits 1; the new exit-pattern anchor must fail before the headings land. | A generic summary, wrong heading vocabulary, or content-first handoff cannot satisfy the result-led canonical pattern. |
| 2 | Empty progress groups and exit sections are omitted rather than rendered as placeholders. | workflow-guidance anchor checker | Observed red: `rg -qi 'omit.*empty'` against the current shared rules exits 1; the omit-empty anchor must fail before that behavior is added. | A rigid always-three-sections template is the degenerate structured output; requiring omission distinguishes proportional structure from placeholder noise. |
| 3 | Related sentences stay together, and bullets or tables are reserved for genuinely parallel facts rather than one item per sentence. | workflow-guidance anchor checker | Observed red: `rg -qi 'related sentences.*together'` against the current shared rules exits 1; the cohesion anchor must fail before the shared rule carries both constraints. | Keyword-only headings can still produce choppy prose; the explicit cohesion and list-restraint clause rejects that cheapest wrong interpretation. |
| 4 | Four clause classes have distinct root-conformance diagnostics: progress labels, exit headings, empty-section omission, and cohesive prose/list restraint. Missing or empty shared rules fail closed. | root conformance on the real tree and planted fixture | Red-first implementation: add the four anchor checks before editing the shared rules and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; it must fail with the targeted diagnostics. | Distinct diagnostics keep a red attributable and prevent one surviving keyword from masking loss of another contract class. |
| 4 | A workflow-guidance canary permanently proves the shared-output check bites. | canary fixture registry and fixture-bite inventory | Red-first implementation: register a planted shared-rules fixture without the progress contract and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; missing EXPECT/bite wiring must fail the fixture contract. | If the anchor is deleted or weakened until the planted rules pass, the known-broken fixture stops biting and turns the canary red. |
| edge of 1 | Fresh-session dogfood covers exactly two representative classes: a tool-using phase with a substantial progress update and a compact phase exit with no material details. | real phase conversation in a fresh session | Not TDD-able: conversational cohesion is not a portable machine schema; record both rendered responses and the reviewer readability verdict in the synthesis handoff. | The two classes catch both over-structuring and under-structuring: missing labels in ongoing work, and an unnecessary empty or trivial Details section at exit. |

**Degenerate-implementation check.** The cheapest wrong edit adds the three exit
headings everywhere, prints empty placeholders, and turns each sentence into a bullet.
Rows 2–4 reject that rendering. The cheapest wrong gate checks one keyword without a
canary; rows 5–6 reject it. The dogfood row catches prose that satisfies every anchor
but still reads mechanically.

### Edge inventory

- **Error path — phase stops short:** covered by the same exit pattern; `Result`
  states the stop, `Details` carries material cause/state, and `Next` gives the exact
  resume route. The phase-specific command still owns which facts are required.
- **Empty/absent details:** covered by story 2's omit-empty row; a compact exit may use
  only `Result` and `Next`, or only `Result` when no action remains.
- **Boundary values — one-sentence update vs substantial update:** covered by the
  implementation decision. A routine one-sentence acknowledgement stays untemplated;
  meaningful intermediate state plus continued work uses **Status:** and **Next:**.
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
- **Every shipped harness surface:** covered through the single shared communication
  source loaded by the harness working agreement. Thin command and skill adapters do
  not restyle output.
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
