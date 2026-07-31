# FT126 recurrence tallying

Status: implemented
Roadmap: FT126

## Problem

Repeated evidence for one roadmap item has no durable identity or count. The
maintenance phase must infer recurrence from prose, duplicate capture artifacts can
inflate that inference, and draining the artifact erases the only visible record.
This encourages duplicate roadmap rows and leaves prioritization dependent on session
recall.

## Solution

Let an idea, learning entry, or retrospective recommendation cite one current roadmap
owner and one stable incident key. Store each accepted owner/incident pair on the
roadmap row, derive its count from those keys, and expose pending, already-recorded,
and structurally invalid citations through the existing roadmap context snapshot.
The reviewed maintenance phase records each new key before removing its source and
uses recurrence as the first tie-breaker after stronger priority constraints.

Compiled decision map:
`specs/ft126-recurrence-tallying/decisions/recurrence-tallying.md`.

Same-session veto surface: the map and the spec-level defaults that compile it remain
explicitly vetoable at sign-off. Those defaults are the schema-3 block and field
names, the `sequence_trusted` encoding, the five discrepancy kinds, the retrospective
improvement-subsection unit scope, the exit-2 syntax versus exit-1 unknown-owner split,
fail-closed owner validation against an untrusted roadmap, capture-row sort order and
literal state names, schema-3-only maintenance compatibility, the `baseline-01`
migration names, and placement of recurrence after explicit pricing but before the
existing defect and cost tie-breakers.

The reviewer approved narrowing FT126 to this mapped recurrence capability. The
unmapped roadmap-parser, context-completeness, lifecycle-discrepancy, and claim-probe
work remains live as FT172.

## User stories

1. As a person capturing an idea, I want to attach one current roadmap owner and
   incident key without changing ordinary unowned capture, so that recurrence can be
   recorded at the source. Line: gpt-5.6-luna / low. The argument grammar and append
   seam are exact, and focused unit plus runtime contracts observe every result.

2. As a roadmap maintainer, I want each row to hold sorted unique incident keys and
   expose their derived count, so that the durable tally has one source and cannot
   drift from a heading number. Line: gpt-5.6-terra / medium. The seam is known and
   well covered, but parser correctness and the one-time legacy migration affect the
   evidence the maintenance phase trusts.

3. As a maintenance-phase operator, I want one context projection across ideas,
   learning entries, and retrospective recommendations, so that duplicate, recorded,
   malformed, and unknown citations are visible before I drain anything. Line:
   gpt-5.6-terra / high. Cross-source capture-unit association and fail-visible trust
   state are subtle oracle behavior even though the existing AXI seam is fixed.

4. As a reviewer prioritizing the roadmap, I want recurrence applied only after
   severity, actionability, dependencies, and explicit pricing, so that repeated
   evidence breaks real ties without overruling stronger decisions. Line:
   gpt-5.6-sol / high. This changes reusable phase guidance, so the kit's leverage
   override applies even though conformance can anchor its structure.

## Implementation decisions

- `internal/roadmap` remains the deep owner of the occurrence-token grammar, roadmap
  ledger, owner validation, cross-source normalization, deduplication, derived counts,
  discrepancy classification, and context trust result. The CLI and source-specific
  packages do not reconstruct recurrence policy.
- An occurrence token is the final non-whitespace text of one capture unit and has
  exactly this form: `[occurrence:<owner>/<incident>]`. The owner is an existing,
  current `FT` row ID. The incident is 1–64 ASCII bytes, contains only lowercase
  letters, digits, and hyphens, and begins and ends with a letter or digit. Text that
  resembles a token anywhere else in a unit is ordinary prose.
- `bench idea --owner <FT> --incident <key> <text...>` declares both value flags
  through the shared argument grammar and appends one canonical token to the dated
  idea line. Both flags must appear together. Syntax and incident errors exit 2;
  an owner value outside `FT[1-9][0-9]*` exits 2; and a syntactically valid owner
  absent from the current roadmap, retired from it, or unverifiable because the
  roadmap is structurally untrusted exits 1. Every refusal leaves `IDEAS.md`
  byte-identical. The no-flag form, its stdout, its multi-word join, its `--` escape
  for leading-dash text, and its one-write newline normalization remain unchanged.
- Capture-unit ownership stays with the source packages. One idea row is one unit;
  one parsed learning entry is one unit; and each recommendation paragraph or list
  item under a retrospective's `## Agent-experience improvements` subsections is one
  unit. Those packages expose unit bodies and stable source locations.
  `internal/roadmap` interprets only the final token and never scans a whole document
  globally.
- A roadmap row has zero or one physical line matching
  `Occurrences: <key>[, <key>...]`. Keys are sorted bytewise, unique, and use the same
  incident grammar as capture tokens. No line means count zero. An empty line,
  malformed key, duplicate key, unsorted key set, or second line rejects the ledger
  as a structural discrepancy; the count is always derived from accepted keys.
- The context snapshot advances from schema 2 to schema 3. Its `context` block adds
  `sequence_trusted`; structural occurrence discrepancies set it to false while the
  snapshot still renders all evidence and the recommended sequence. The
  `roadmap_rows` block adds `occurrence_count` and the canonical `occurrence_keys`
  string. No numeric count is stored or rendered from a second field.
- The schema adds
  `capture_occurrences[]{owner,incident,source,capture_unit,state}`. Rows are sorted by
  those five fields. Multiple source rows may share one owner/incident identity, but
  the normalized pending set contains that pair once. State is `pending` until the
  row contains the key and `already-recorded` afterward.
- The schema also adds
  `occurrence_discrepancies[]{source,capture_unit,kind,owner,incident,structural}`.
  Kinds are `malformed-token`, `unknown-owner`, `multiple-tokens`,
  `malformed-ledger`, and `already-recorded`. The first four are structural;
  `already-recorded` is false because it is valid evidence that the duplicate source
  may be removed without incrementing the row.
- The same incident key on separate capture units for different owners is one
  occurrence for each owner. Multiple tokens on one unit are never treated as
  multiple owners: the whole unit is structurally invalid until split.
- The one-time migration removes `evidence supplied` count claims from the nine
  current row headings and writes canonical `Occurrences:` keys that preserve these
  counts:

  | row | migrated count |
  |---|---:|
  | FT71 | 1 |
  | FT158 | 3 |
  | FT128 | 1 |
  | FT98 | 3 |
  | FT169 | 1 |
  | FT126 | 1 |
  | FT141 | 1 |
  | FT94 | 1 |
  | FT125 | 1 |

  A distinct incident already named by its row receives a grammar-valid descriptive
  key. Any remaining unnamed evidence uses `baseline-01`, `baseline-02`, and
  `baseline-03` in order; the migration commit supplies the historical pin without
  inventing a source narrative.
- `.agents/commands/bench-what-next.md` consumes schema 3 only. During a reviewed
  drain it adds each new owner/incident pair to the owning row before removing any
  source unit, removes an already-recorded duplicate without adding a key, and stops
  before mutation when `sequence_trusted` is false.
- Sequence precedence is: severity, actionable versus blocked state, literal
  dependencies, explicit reviewer pricing, then descending occurrence count. The
  existing defect/feature and cost rules operate only after those inputs remain tied.
  No CLI command sorts or rewrites `ROADMAP.md`; judgment and mutation remain in the
  reviewed maintenance phase.

## Testing decisions

- A good test drives document bytes or the shipped CLI and observes typed facts,
  stdout, exit code, and file bytes. It does not assert private regular expressions or
  helper call order.
- The primary seam is `bench roadmap --context`, which makes ledger parsing,
  cross-source normalization, discrepancy state, deterministic order, and schema
  version externally visible. Focused `internal/roadmap` tests attach at the same
  exported document/context interface for hostile parser inputs.
- The `bench idea` runtime seam owns valid append and fail-without-mutation behavior.
  Both the kit binary and a linked repository's by-path CLI reach the same core.
- The conformance seam owns the schema-3 and maintenance-procedure prose anchors plus
  the absence of legacy heading counts after migration. The one-time migrated counts
  are checked against the enumerated table above and then derived through the context
  seam.
- Prior art is `internal/roadmap/context_test.go`,
  `internal/contract/axi/axi_roadmap_context_test.go`,
  `internal/roadmap/roadmap_test.go`, and the `bench idea/roadmap` runtime contract in
  `internal/contract/runtime`.
- Focused package, AXI, runtime, and conformance tests precede the feature gate. The
  gate that defines done is `bench gate`.

### Seam diagram: roadmap occurrence projection

```text
trigger: bench roadmap --context or an internal/roadmap contract test
    |
    v
ROADMAP + ideas + learnings + retros --> [ internal/roadmap ] --> schema-3 TOON + trust
                                        [ parse / normalize ]
                                           ^ tests attach here: fixture bytes in,
                                             typed facts or CLI bytes out
```

### Seam diagram: owned idea capture

```text
trigger: bench idea
    |
    v
argv + trusted ROADMAP rows --> [ idea CLI over internal/roadmap grammar ] --> IDEAS append or refusal
                                  ^ tests attach here: invoke binary and compare
                                    exit, stdout, and before/after file bytes
```

### Seam diagram: maintenance guidance conformance

```text
trigger: conformance phase
    |
    v
phase command + ROADMAP copy --> [ workflow anchor checks ] --> targeted diagnostics
                                  ^ tests attach here: remove one required clause
                                    or retain one legacy heading claim
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | valid `--owner FT98 --incident 2026-07-30-scoped-roadmap-commit` capture appends one final canonical token and preserves the existing parked stdout | idea CLI | Observed red: the current CLI probe exits 2 with `unknown argument: --owner`. | The old unowned append path cannot produce the required owner/incident evidence. |
| 1 | the no-flag form appends exactly the current dated free-text line and output | idea CLI | Already covered by `TestIdeaCreatesDatedLine` and the runtime idea/roadmap contract; keep them unchanged as positive controls. | A metadata-only rewrite that breaks ordinary capture fails the existing contract. |
| 1 | either metadata flag alone, an empty value, a malformed owner spelling, and a malformed incident each exit 2 without changing `IDEAS.md` | idea CLI | Not TDD-able before the new flags exist; future table and runtime tests must first fail on the current unknown-flag path and compare before/after bytes. | It pins syntax and token-grammar refusal separately from semantic owner lookup. |
| 1 | a well-formed owner absent from the current roadmap exits 1 without changing `IDEAS.md` | idea CLI | Not TDD-able before owner lookup exists; future runtime fixture first fails because the current CLI exits 2 at the unknown flag. | An implementation that accepts any well-formed FT token cannot satisfy the current-owner lookup. |
| 1 | an owner present in Git history but retired from the current roadmap exits 1 without changing `IDEAS.md` | idea CLI | Not TDD-able before owner lookup exists; future fixture commits the row, removes it, then first fails on the current unknown-flag path. | A stale or all-time FT registry accepts the retired ID while a wholly nonexistent-ID test would still pass. |
| 1 | an otherwise present owner in a structurally untrusted roadmap exits 1 without changing `IDEAS.md` | idea CLI | Not TDD-able before ledger trust exists; future fixture first fails because current argument parsing never reaches roadmap validation. | Partial row recovery cannot authorize mutation when the document also says its ledger is malformed. |
| 1 | metadata flags before multi-word text and `--` before leading-dash text preserve the complete text and place the token last | idea CLI | Observed red: the current CLI rejects the first metadata flag before it can append either text form. | A parser that treats a flag value as text or truncates variadic text cannot match the expected line. |
| 2 | a missing `Occurrences:` line yields count zero; one canonical sorted line yields exactly its key count and key string | roadmap context | Observed red: `bench roadmap --context` has no `occurrence_count` or `occurrence_keys` columns. | A stored numeric count, ignored ledger, or constant zero cannot produce both derived fields. |
| 2 | empty, malformed, duplicate, unsorted, or repeated ledger lines produce `malformed-ledger` and make `sequence_trusted` false | roadmap context | Observed red: the current context has neither the discrepancy block nor a sequence-trust field. | A permissive parser or dedup-on-read implementation would incorrectly preserve trust. |
| 2 | migration removes legacy heading count claims and preserves counts for FT71=1, FT158=3, FT128=1, FT98=3, FT169=1, FT126=1, FT141=1, FT94=1, and FT125=1 | conformance plus roadmap context | Observed red: `! rg -q '^\*\*FT.*evidence supplied' ROADMAP.md` exits 1 on the nine enumerated rows. | The exact inventory catches a skipped row or collapsed multi-occurrence baseline instead of accepting a partial migration. |
| 2 | no heading or second ledger stores a numeric recurrence count | conformance | Observed red: current headings still carry singular and `three times` count claims. | It enforces the key set as the sole durable source of the derived count. |
| 3 | schema 3 renders `sequence_trusted`, the two occurrence columns, `capture_occurrences`, and `occurrence_discrepancies` in fixed block order | roadmap context | Observed red: current probes for all three new blocks and the trust field exit 1. | A partial schema advance cannot satisfy the complete ordered header inventory. |
| 3 | capture rows sort by owner, incident, source, capture unit, and state and use only `pending` or `already-recorded` | roadmap context | Observed red: schema 2 has no capture-occurrence rows or state values. | Stable byte output can still pass with an unintended order or vocabulary unless both are asserted directly. |
| 3 | two source units with the same owner/incident expose both source locations but normalize to one pending pair | roadmap context | Not TDD-able before occurrence projection exists; future AXI fixture first fails because current output has no pair identity. | Counting artifacts rather than unique pairs inflates the pending tally and fails the grouped identity assertion. |
| 3 | after the owning row records the key, every matching source row becomes `already-recorded`, adds a non-structural discrepancy, and contributes no pending pair | roadmap context | Not TDD-able before occurrence state exists; extend the same fixture and first observe the absent state transition. | An increment-on-every-drain implementation or a silent duplicate cannot produce the required state and advisory row. |
| 3 | malformed final token, unknown owner, multiple tokens on one unit, and malformed ledger each emit their named structural discrepancy and set trust false | roadmap context | Observed red: `bench roadmap --context` exposes no occurrence discrepancy or trust result today. | Dropping any invalid unit silently leaves the sequence trusted and misses its kind-specific assertion. |
| 3 | one incident cited on separate units for two owners remains one pending pair per owner | roadmap context | Not TDD-able before normalization exists; future fixture first fails with no occurrence table. | Deduplicating by incident alone incorrectly erases one independently evidenced FT. |
| 3 | a learning entry with token-like prose before its end is not cited, while a final token belongs only to that entry | roadmap context | Not TDD-able before unit projection exists; future learning fixture first fails because current context does not associate tokens with entries. | A global text scan fabricates a citation from prose or assigns it to the wrong entry. |
| 3 | one retro with separate cited recommendation paragraphs and list items emits one unit per recommendation under the improvement subsections | roadmap context | Not TDD-able before retrospective units exist; future retro fixture first fails because current context exposes only whole-file bodies. | Whole-document matching cannot preserve multiple recommendation owners or their source locations. |
| 3 | unchanged repeated context invocations are byte-identical and read-only, including from a nested cwd | roadmap context | Already covered by `testRoadmapContextComplete` and `testRoadmapContextReadOnlyOffline`; extend their expected schema. | A nondeterministic sort, cwd-relative read, or hidden mutation changes existing byte and repository-state assertions. |
| 4 | the maintenance phase adds every pending key before removing its source and removes already-recorded sources without another key | conformance | Observed red: `.agents/commands/bench-what-next.md` contains no occurrence or ledger procedure. | Removing first loses durable identity; always adding makes a repeated drain increment twice. |
| 4 | the maintenance phase requires schema 3 and stops on any other context schema instead of guessing missing recurrence facts | conformance | Observed red: the current command explicitly consumes schema 2 and has no recurrence compatibility clause. | A silent schema fallback can run the drain without the trust or occurrence blocks this procedure requires. |
| 4 | equal-severity, equally actionable rows order by dependencies and explicit pricing before descending count, then apply the existing defect and cost rules | conformance plus equal-class fixture | Observed red: the current command contains no occurrence-count tie-break anchor. | An unconditional highest-count sort fails the stronger-input controls, while ignoring count fails the equal-class case. |
| 4 | structural occurrence discrepancies stop the maintenance phase before its batch mutation while leaving the complete context evidence visible | conformance plus roadmap context | Observed red: schema 2 has no `sequence_trusted` fact for the phase to gate on. | A prose-only warning that still mutates cannot satisfy the explicit stop and unchanged-source assertions. |
| edge of 1 | incident length 1 and 64 accept; length 0 and 65, uppercase, non-ASCII, separators, leading/trailing hyphen, and control bytes reject without mutation | idea CLI and roadmap context | Not TDD-able before the shared grammar exists; future boundary table must first fail on the current unknown-flag or ignored-ledger behavior. | It pins every byte and length boundary rather than testing one representative invalid key. |
| edge of 1 | repeated flags, missing values, flag-looking values, `--`, and text containing spaces or glob characters follow the shared argument grammar | idea CLI | Already covered for the shared parser; extend the idea runtime table for both value flags and the post-terminator text form. | It catches a local parser that reads flag values as positionals or lets the shell reshape text. |
| edge of 1 | invocation from a repository path containing spaces or glob characters and from a symlinked linked CLI writes only the repository-root inbox | idea CLI | Already covered in part by nested-cwd and linked CLI contracts; extend the metadata case to the hostile repository path. | A cwd-relative or unquoted path implementation writes the wrong file or fails before validation. |
| edge of 1 | a newline-less existing inbox receives the complete text and token in the same append write after one normalizing newline | idea CLI | Already covered by `TestIdeaNewlineNormalization`; extend its appended entry with metadata. | A separate metadata write can be interrupted between text and token or merge onto the prior line. |
| edge of 2 | CRLF and a missing final newline preserve a valid ledger; tabs, returns, or newlines inside a key reject rather than split the line sink | roadmap context | Not TDD-able before ledger parsing exists; future byte-table test first fails because current parser ignores the reserved line. | It tests both the line endings the document permits and control bytes the key grammar cannot survive. |
| edge of 2 | absent versus empty ROADMAP and capture files retain their existing distinct source states while empty occurrence sets render definitive zero-row blocks | roadmap context | Already covered by context source-state contracts; extend their schema-3 block expectations. | A new occurrence reader cannot collapse a failed or empty control record into trustworthy absence. |
| edge of 3 | FIFO, device, socket, oversized file, and dangling symlink capture sources are classified before reads and never block occurrence projection | roadmap context | Already covered by roadmap, retrospective, and AXI source-state tests; extend them to require no fabricated occurrence rows. | Re-reading a source globally for tokens bypasses the classified readers and can block or follow unintended bytes. |
| edge of 3 | TOON-refused control bytes in a source location fail through the existing structured render error; escapable tab, newline, and return in source text do not split a table record | roadmap context | Already covered at the TOON seam; add one occurrence-row projection for each accepted and refused class. | Exercising only rejected controls would miss permitted bytes that reach the flat-table sink. |
| edge of 3 | repeated source artifacts and repeated drains keep one normalized owner/incident pair | roadmap context | Not TDD-able before normalization exists; the cross-source and already-recorded fixtures supply the first red. | It makes re-run idempotency a set-identity contract instead of a best-effort prose promise. |

Degenerate implementations considered: storing a mutable number fails the derived-key
columns; counting capture artifacts fails the same-pair multi-source row; scanning whole
documents fails learning and retro unit association; silently dropping invalid tokens
leaves `sequence_trusted` true; sorting solely by count fails the stronger-input controls;
and appending text and metadata separately fails the one-write newline-less fixture.

### Edge inventory

- Error path — resolved by paired-flag, unknown-owner, malformed-ledger, discrepancy,
  and untrusted-sequence rows.
- Empty/absent input — resolved by the source-state and definitive-empty-block row.
- Boundary values — resolved by the 1/64-byte incident and nine-row migration rows.
- Malformed input — resolved by token, ledger, control-byte, and argument-grammar rows.
- Interrupted or partial state — the idea append remains one write, covered by the
  newline-less metadata row; abrupt termination during a kernel partial write is
  **Won't handle** because transactional capture durability is a separate capability.
- Re-run idempotency — resolved by same-pair multi-source, already-recorded, and repeated
  drain rows.
- Hostile environment — resolved by hostile repository paths, nested cwd, classified
  special files, dangling symlinks, bounded files, and TOON byte classes.
- **Won't handle:** reporting post-write cleanliness or a derived count from `bench idea`
  — the command reports only the captured text, so its own append cannot falsify that
  field.
- **Won't handle:** a missing external tool — recurrence parsing and rendering add no
  runtime or install-time dependency.
- **Won't handle:** hooks and adapters — they do not invoke either affected CLI command;
  the real kit and linked-repository CLI surfaces remain in scope.
- **Won't handle:** destructive worktree registrations — recurrence reads and appends
  control records without creating, releasing, or deleting a worktree.
- **Won't handle:** non-TTY stdin — neither command prompts or reads stdin.
- **Won't handle:** durable acknowledgment under host-side I/O pressure — the feature
  preserves `bench idea`'s existing append durability and adds no `fsync` promise.

## Out of scope

- **FT172 explicit roadmap-row header grammar** — choosing and enforcing the exact
  bold-heading form is a separate parser-hardening capability with an unresolved
  reviewer fork (`3 edits, 2 gate runs`).
- **FT172 last-drain workload boundary** — deriving the last drain, commits since, and
  code-touch set is a separate Git-evidence capability (`4 edits, 3 gate runs`).
- **FT172 spec/row lifecycle discrepancies** — missing, retired, and unowned spec
  cross-checks are a separate reconciliation capability (`3 edits, 2 gate runs`).
- **FT172 capture diagnosis verification** — deciding between the general claim probe
  and the cheaper symptom/claim rule is a separate guidance-policy capability
  (`2 edits, 2 gate runs`).
- **Automatic roadmap mutation or global sorting** — a separate unreviewed-mutation
  capability rejected by the map; `/bench-what-next` remains the owner
  (`4 edits, 3 gate runs`).
- **A second occurrence history file or event database** — a separate persistence
  capability rejected because row keys plus Git already own current and historical
  state (`6 edits, 3 gate runs`).
- **Multiple primary owners on one capture unit** — a separate multi-owner citation
  capability rejected by the map; separate units may cite the same incident for
  different owners (`3 edits, 2 gate runs`).
