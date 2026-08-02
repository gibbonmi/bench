# Ticket contracts: charges preserve blast radius, contracts, and evidence authorship

Status: staged

Decision source: ROADMAP.md row FT164 (the 2026-08-02 `/bench-what-next` drain,
commit `fbdcd1b`), a reviewed artifact assembling ten drained sources, read with
its FT175 delivery path and recommended sequence, which direct the split this
spec performs: the ticket-contract core lands here; the model-comparison,
inventory-currency, and shared-cache repair riders move to a follow-up. The
split audit for the kept clauses: the self-probe, probe-site, and
backup-isolation clauses are charge-side duties of the ticket contract itself —
they bind every write-delegate charge, not a specialized repair lane — so they
belong in the core rather than with the riders. The row was re-verified against
the current tree on 2026-08-02, with three corrections: the 25-ticket corpus
has shrunk to 12 (the pcgs tickets are reachable only via `git show 3972744:`);
the gate-cadence clause is already fully taught by `craft-tickets` and lands as
a no-op noted under Implementation decisions; and blast-radius classification
already exists as a section ahead of the breakdown — what is missing is the
procedure invoking it, fence-sized migrations, and the blocker-bound contract
ticket, so that clause lands as a sharpening, not an addition.

## Problem

A ticket written cold from the `craft-tickets` template cannot be assigned: the
specbuild parser requires `- [ ] [ID] text` acceptance rows and a one-line
backticked `Ownership fence:`, while the template teaches unlabeled rows, no
fence field, and a title-keyed `Blocked by:` — 17 tickets were hand-normalized
before the last lifecycle build could proceed. The conventions the last builds
actually used (fence, assumptions, requirement IDs, single-line rows) are
taught nowhere, so they decay the moment a ticket is written without reading
old examples. And the breakdown procedure misses failures the last three
builds paid for: its numbered method never invokes the classification section
sitting beside it, so a wide refactor was sliced with a deep unit grouped into
every consumer family and exhausted a delegate's context; two tickets with
correct disjoint fences were each green alone and red composed over an
undeclared domain mismatch; silent fixture-greens survived until a coordinator
probe differed in kind and site from the delegate's own; and a restore glob
clobbered four out-of-fence files.

## Solution

`craft-tickets` teaches the executable-contract ticket: the machine-parseable
field shapes, the breakdown procedure opening on blast-radius classification,
a contracts-discovery step that names what crosses each fence before ticket
files are written, and evidence-authorship rules. `craft-delegate` gains the
charge-side duties: the self-probe clause, probes differing in kind and site,
the enumerated-family tracing duty, and worktree-local backup isolation.
Every new rule is pinned by a section-scoped workflow anchor exercised by the
mutation harness, and the template's Good example is parsed by the real
specbuild parser inside the gate, so the taught shape and the accepted shape
cannot drift apart again.

## User stories

Stories 1–6 are **top** tier under the profile's skill-authoring leverage
override; story 7 is **mid**. Resolved ids are the `claude` column of the
profile's harness × tier table.

1. As a coordinator drafting tickets, I want the template to teach the
   machine-parseable executable-contract form — single-line acceptance rows
   `- [ ] [ID] <behavior>` with ticket-local IDs, a one-line backticked
   `Ownership fence:` naming every path the ticket will write, a one-line
   `Assumptions:` field, a basename-keyed `Blocked by:`, and a
   `## Red mutations` table binding each row to one concrete mutation, its
   independent owner, and the public operation sequence that proves it — so a
   ticket written cold from the template assigns without hand-normalization.
   Line: fable / high. Guidance prose compounds through every session that
   loads it, and the field shapes are load-bearing for the parser.

2. As a coordinator, I want the Good example to demonstrate every taught field
   and the Bad example replaced by an oversized-but-credible ticket, so the
   examples teach against the failure the corpus actually commits rather than
   one it does not. Line: fable / high. Examples are what the corpus copies,
   and a field taught but never demonstrated decays first.

3. As a coordinator slicing a build, I want the breakdown procedure's numbered
   method to open on the blast-radius branch — wide refactor takes
   expand–migrate–contract, migrations sized by one ownership fence each, the
   contract ticket's `Blocked by:` naming every migration — with step 2's
   independently-green check made mechanical by fence disjointness, the
   one-line test-harness ceiling noted beside the rule, and a cadence-changing
   ticket naming which command authors gate evidence and which phase consumes
   it. Line: fable / high. These are the procedure rules three builds paid to
   learn, and their phrasing decides whether they fire.

4. As a coordinator, I want a contracts-discovery step between drafting the
   breakdown and writing ticket files — every value crossing a fence names its
   type, its membership or domain rule, its ordering, and its absence
   semantics; each such invariant lands as an acceptance row on the consumer
   ticket asserted against the real producer and the whole enumerated family;
   a junction ticket when neither side can assert it alone; claims re-derived
   from the tree after earlier tickets land — so a cross-fence mismatch is a
   sentence at slicing time instead of a composed red six tickets later.
   Line: fable / high. This is the row's core new procedure and its wording
   carries the whole defense.

5. As a coordinator charging a write-delegate, I want the charge duties in
   `craft-delegate` — the charge names the central-property-breaking mutation
   and the delegate applies it to its own finished work, reports the observed
   result, and adds the missing row on a silent green; the coordinator's probe
   differs in kind and in site from any the delegate ran, with the
   omission/swap probe-kind vocabulary in the charge template; a
   family-extending charge names every registry found by tracing an existing
   sibling; transient backups live inside the delegate's own worktree under a
   unique name and restores name exact files, never globs — so silent greens
   and cross-fence clobbers are caught at return time. Line: fable / high.
   Charge prose is the only surface that reaches a low-context delegate.

6. As a spec author, I want the spec-side halves — a slicing-time sentence in
   `craft-spec` naming the value contracts at each fence by pointer to the
   ticket rule, and a process-boundary lifecycle edge class in the canonical
   edge walk plus the profile's hostile-input checklist — so cross-seam
   mismatches and serialize-then-reload defects are visible before a build
   starts. Line: fable / high. Both are one-sentence additions whose owner
   placement is the decision the roadmap flags for review.

7. As the gate's maintainer, I want the enforcement wired: each needle in this
   spec's anchor inventory registered section-scoped in the anchor check with
   a mutation-table row proving it bites, the two existing anchor artifacts
   that move (the template placeholder needle and the probe-kind sentence's
   hard-wrapped mutation entry) updated in the same diff as their text, and
   the example-agreement check landed beside the other conformance checks.
   Line: opus / medium. The profile caches mid effort for conformance logic;
   the tier is mid because needle and table wiring is mechanical at a known
   seam — the needle *text* is authored by stories 1–6 at the top tier, and
   story 7 adds no prose of its own.

## Implementation decisions

**The template teaches the shape the parser accepts; parsing behavior does not
change.** FT174's closing question — template teaches the parseable shape, or
`resolveTicket` accepts the template's — is decided here in the template's
direction, as the roadmap's dependency table already assumes ("build the
parser against the identifier form the template teaches"). The one
`internal/specbuild` edit this build makes is an export: the existing ticket
resolution gains an exported entry point (the same parse, same refusals,
export only) so the conformance check can invoke it cross-package — the
precedent is the passlist check importing an exported symbol rather than
re-implementing a parse. Refusals for wrapped fields, the silent
`internal/<pkg>` fence fallback, blocked-by walking, fence disjointness at
assign, and requirement-ID completeness all ride FT174.

**Field shapes taught, matching the live parser exactly.** Acceptance rows are
single-line `- [ ] [ID] <behavior>`; the ID is ticket-local, a short uppercase
tag plus number, unique within the ticket, with one caveat taught beside it:
only `R`-prefixed IDs range-expand (`[R1-R3]`), so the template teaches
explicit per-row IDs and names ranges as `R`-only. `Ownership fence:` is one
line, comma-separated, each entry backticked, entries are path prefixes — a
package directory or an exact file — and the fence enumerates every path the
ticket will write, because checkpoint enforcement is a whitelist and scoping
prose ("the shellcheck entries of…") is unrepresentable. `Assumptions:` is one
line; clauses separate with semicolons because the parser splits the line on
commas, and a wrapped continuation line is silently dropped. `Blocked by:` is
one line: `none`, or sibling ticket file basenames.

**Requirement IDs are ticket-local.** The row's claim that `bench coverage
<spec>` is the ID source does not survive the tree: the coverage map has no ID
column, and FT180's spec-less folders will have no map at all. Declaring IDs
ticket-local keeps one owner for the convention, costs nothing now, and leaves
FT174 free to add a completeness check over ticket files alone. **Confirmed
2026-08-02**, with one amendment: a change that goes straight to tickets
without a spec still carries its FT roadmap ID as provenance — only a genuine
hotfix goes without one.

**`Blocked by:` keys on file basenames, not titles and not `#N`.** Titles are
what FT174 exists to retire (a retitle silently breaks the edge), and `#N`
requires an ordering ticket folders do not have. Basenames are stable,
already the assignment argument (`--ticket <file>`), and give FT174's walker a
resolvable identifier. **Confirmed 2026-08-02** — FT174's row gestures at the
decision-map `#N` grammar, but what it directs is extending that validator's
cycle/dangling logic, which walks basename keys just as well.

**The template gains a `## Red mutations` table** — one row per acceptance ID:
the concrete mutation, its independent owner, and the exact public operation
sequence that proves it red. Nothing parses its content in this build; it is
the executable-contract half the pcgs tickets carried by hand. Re-derivation
is taught beside it: ticket claims are re-derived from the tree after earlier
tickets land, never from the spec's account of the base.

**The anchor inventory — the enumerated set behind every "every".** Anchors
are substring checks over a whole file, so an unscoped needle is satisfiable
from an appendix paragraph; every needle below therefore registers through the
check's existing section-scoping facility, pinned to the section that owns it.
What anchors cannot pin is ordinal position within a section — that residual
is review-graded and named here rather than claimed. The set, one needle per
listed fact, section in parentheses:

- craft-tickets, template block (5): the acceptance-row placeholder in the new
  labeled form; the fence line; the assumptions line; the basename blocked-by
  line; the red-mutations table header.
- craft-tickets, cadence paragraph (2): the existing gate-checkbox prohibition
  sentence (newly pinned — it already exists, so its red capability comes from
  the mutation table, not from anchor-first); the evidence-authorship sentence.
- craft-tickets, breakdown procedure (4): the method's opening classification
  branch; fence-sized migrations; the contract ticket's blocked-by-every-
  migration rule; step 2's fence-disjointness method.
- craft-tickets, definition paragraph (1): the one-line test-harness ceiling.
- craft-tickets, contracts-discovery step (4): the four-fact sentence; the
  consumer-row-against-real-producer-and-whole-family sentence; the junction
  rule; the re-derivation rule.
- craft-delegate (5): the self-probe clause; the site-differs sentence; the
  omission/swap vocabulary in the charge template; the registry-tracing
  sentence; the backup-isolation sentence.
- craft-spec (2): the slicing-time contract-pointer sentence; the
  process-boundary class token in the edge-walk sentence.
- profile checklist (1): the process-boundary entry's lead phrase.

24 needles. Two existing anchor artifacts move in the same diff as their text:
the template placeholder needle (both occurrences of the old placeholder row
are replaced, the old needle retired, the new one section-scoped to the
template block) and the probe-kind sentence, whose mutation-table entry
hard-codes the current line wrapping and fails on the hoist unless edited with
it. The craft-tickets↔craft-spec cross-pointer anchors are not touched: the
craft-spec sentence is an addition beside them, not an edit.

**Each needle gets one mutation-table row; set completeness is review-graded.**
The bite harness's table-driven mutation test is the per-needle proof: a table
row lands before the skill text and fails on its anchor count until the
sentence exists, then proves the diagnostic fires when the sentence is
mutated. That makes each row red-capable at unit cost rather than one canary
fixture per needle; canary fixtures stay at their existing per-family
granularity. No completeness assertion exists that counts needles against
table rows — building one is priced out of scope — so "all 24 covered" is
graded by review against this enumeration, and the coverage map says so.

**The example-agreement check.** The Good example sits between explicit
begin/end markers (the passlist check's shape: fail-loud when the markers are
absent, and grading only the marked block, so the template block and the Bad
example can never be graded by mistake). The check materializes the block as a
real `tickets/` file beside a temp spec path — the parser takes paths, not
bytes — and runs the exported specbuild parse over it, asserting: at least two
acceptance rows with distinct IDs in the taught grammar; the exact multi-entry
fence the example states; the exact assumptions; a `Blocked by:` line present;
and a `## Red mutations` section naming every acceptance ID. The expected
values are independently authored literals in the check, per the
independent-expectation ADR: their independence is what lets the named
omission — the example drifting off the parseable shape while the prose stays
plausible — turn the gate red, and that red is demonstrated by the check's own
mutation fixtures.

**This build touches the gate, so `craft-synthesis`'s prose-only substitute is
unavailable.** Completion evidence is a dogfood run in a fresh session: write
one ticket cold from the new template, assign it through `bench spec build
assign` in a throwaway fixture repo, and record the result in the final-check
evidence — plus the mutation-table reds above and a read of every steered
surface.

**Gate cadence is already taught; only its evidence-authorship clause is
new.** The template carries no gate checkbox and the three-beat cadence is in
place — no edit. The new sentence: a ticket that changes gate cadence names
which command authors evidence (`bench gate`, the canonical producing entry)
and which phase consumes it; a bare `gate-run --fresh` prints a valid phase
result without publishing the project-green evidence promotion consumes. When
story 2 replaces the Bad example, the gate-checkbox prohibition survives in
the cadence paragraph's prose, newly pinned per the inventory.

**Owner decisions the roadmap reserved, recommended here.** (1) The
process-boundary lifecycle class — defects visible only after state is
serialized and a fresh process reloads it, and recomposition suites that stop
at first success — lands as a canonical edge class in `craft-spec`'s walk and
as a concrete entry in the profile's hostile-input checklist; it is an edge
class, not a charge rule. (2) The spec-side half of contracts-discovery lands
as one sentence in `craft-spec`'s "Slicing a build for delegates" pointing at
the `craft-tickets` rule by name, mirroring the existing anchored
cross-pointer pair, so the procedure keeps one owner. Both owner placements
are **confirmed 2026-08-02**.

**Sequencing against the craft-research skill: none.** Closed 2026-08-02 in
that session's grill and recorded in `decisions/craft-research.md`:
craft-research owns read-side research discipline only, carves nothing from
`craft-delegate`'s charge, probe, or isolation sections, touches neither
`craft-tickets` nor the "fan-out search" frontmatter line, and states no
tiers. The two builds are disjoint and land in either order.

**The nine older non-parseable tickets are not reformatted** — retrofitting
closed work rewrites history for no gain, per the source row. The follow-up
riders (model-comparison, inventory-currency, shared-cache charge rules) stay
on the roadmap under FT164's row for the next drain to re-home; this spec
does not edit `ROADMAP.md`.

## Testing decisions

- The external behavior under test is enforcement, not prose quality: a rule
  sentence deleted or mutated turns the gate red, and a Good example that the
  live parser cannot read turns the gate red. Prose *correctness* is graded by
  review and the falsification pass, and the coverage map says so honestly.
- Seams: the workflow-anchor check with its section-scoping facility and
  table-driven mutation harness as prior art, and the new example-agreement
  check beside the other conformance checks, invoking the exported specbuild
  parse.
- The feature gate is `bench gate` — a full run every time, since `.agents/`
  is outside the reduced scope.

### Seam diagram: section-scoped anchors

```text
trigger: bench gate (conformance phase, docs-currency-workflow)
    |
    v
skill body text --> [ anchor check: needle within its owning section ] --> gate diagnostic
                        ^ tests attach here: the mutation table rewrites the
                          pinned sentence and asserts the named red
```

### Seam diagram: example-agreement check

```text
trigger: bench gate (conformance phase)
    |
    v
craft-tickets SKILL.md --> [ extract marker-delimited Good example ] --> block bytes
                                     | fail-loud when markers absent
                                     v
                    temp spec dir + tickets/<file> --> [ exported specbuild parse ]
                                     |
                                     v
        rows, fence, assumptions, blocked-by, mutations table == authored literals
        ^ tests attach here: mutate the example (wrapped fence, unlabeled row,
          dropped marker) and assert the named red
```

### Acceptance coverage map

Anchor-first is the red mechanism for newly taught rules: the section-scoped
needle and its mutation-table row land before the skill text and stay red
until the sentence exists in the owning section. The one pinned sentence that
already exists (the gate-checkbox prohibition) takes its red from the mutation
table alone, and rows whose subject is semantic quality say so rather than
pretending gate coverage.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the template teaches single-line `- [ ] [ID] <behavior>` acceptance rows | section-scoped anchors | TDD-able: the old placeholder's needle pins the current template, so a template rewrite alone reds it, and the new needle is scoped to the template block, so an anchor-first landing reds until both placeholder rows change. | The current placeholder shape parses to zero rows and refuses assignment. |
| 1 | the template carries a one-line backticked `Ownership fence:` that enumerates every path the ticket will write | section-scoped anchors | TDD-able: anchor-first in the template block. | A fenceless ticket either refuses assignment or, when its prose happens to mention `internal/<pkg>`, assigns on that silently inferred fence. |
| 1 | the template carries a one-line `Assumptions:` field taught with semicolon-separated clauses | section-scoped anchors | TDD-able: anchor-first in the template block. | The parser reads only the prefixed line and splits it on commas, so a wrapped continuation drops silently and comma-joined prose shatters into fragments frozen by drift checks. |
| 1 | `Blocked by:` is taught keyed by sibling ticket file basenames | section-scoped anchors | TDD-able: anchor-first in the template block. | Title keys break silently on retitle, and FT174's walker needs a resolvable identifier. |
| 1 | the template carries a `## Red mutations` table binding each row ID to a mutation, its owner, and the proof sequence | section-scoped anchors | TDD-able: anchor-first in the template block. | An umbrella claim is satisfied by delete-only probes that miss swaps, stale identities, and second-call failures. |
| 1 | a ticket written verbatim from the Good example parses to at least two distinct-ID rows, its exact fence, and its exact assumptions | example-agreement check | TDD-able: the check lands against today's example and reds; it goes green only when the example carries the parseable shape. | Anchors pin presence, not parseability — only the live parser proves a cold-written ticket assigns. |
| 2 | the Good example demonstrates every taught field | example-agreement check plus review | Partially TDD-able: rows, fence, assumptions, the blocked-by line, and per-ID mutations-table coverage are graded by the agreement check; field *credibility* is not TDD-able — review grades it. | The example is what the corpus copies; a taught-but-undemonstrated field decays first, and a one-row example demonstrates almost nothing. |
| 2 | the Bad example is an oversized-but-credible ticket | review | Not TDD-able: example quality is semantic; the falsification pass and reviewer grade it. | The current Bad example teaches against a failure the corpus does not commit, spending its slot. |
| 2 | the gate-checkbox prohibition survives the Bad example's replacement | section-scoped anchors | TDD-able via the mutation table only: the sentence already exists, so its new needle starts green; the table row mutating it proves the red. | Replacing the example without pinning the prose deletes the anti-pattern's only remaining illustration. |
| 3 | the breakdown procedure's numbered method opens on the blast-radius branch | section-scoped anchors plus review | TDD-able for the wording: anchor-first in the procedure section. Ordinal position within the section is not anchor-expressible — review grades that the branch is the opening step. | The classification section already exists beside the procedure; what failed in FT154's build is that the method never invoked it. |
| 3 | migrations are sized by one ownership fence each | section-scoped anchors | TDD-able: anchor-first. | The unfenced first slice grouped a deep unit with every consumer family and exhausted the delegate's context. |
| 3 | the contract ticket's `Blocked by:` names every migration | section-scoped anchors | TDD-able: anchor-first. | Prose sequencing without a blocker edge lets the contract ticket run while a migration is still open. |
| 3 | step 2 names fence disjointness as the mechanical concurrent-eligibility check | section-scoped anchors | TDD-able: anchor-first. | "Confirm independently green" without a method is a drafting guess where parallel assignment makes it a correctness precondition. |
| 3 | the one-line test-harness ceiling is noted beside the independently-green rule | section-scoped anchors | TDD-able: anchor-first. | Without the ceiling a one-line change pays a worktree, a fresh delegate, and a full gate by default. |
| 3 | a cadence-changing ticket names the evidence-authoring command and the consuming phase | section-scoped anchors | TDD-able: anchor-first. | Treating producing and non-publishing gate entries as interchangeable cost several redundant full runs. |
| 4 | the contracts-discovery step requires every fence-crossing value to name type, membership, ordering, and absence semantics | section-scoped anchors | TDD-able: anchor-first in the new step's section. | The pcgs build shipped two tickets green alone and red composed over exactly such an undeclared mismatch. |
| 4 | the discovered invariant lands as a consumer-ticket acceptance row asserted against the real producer and the whole enumerated family | section-scoped anchors | TDD-able: anchor-first. | Fixture-satisfied rows let both halves look green for six tickets, and worktree-local membership hides absent family members. |
| 4 | the junction rule: neither side able to assert alone means a junction ticket, and a junction row more than one ticket downstream moves a narrower copy to the junction | section-scoped anchors | TDD-able: anchor-first. | The composed red surfaced six tickets after the junction it should have been pinned at. |
| 4 | ticket claims are re-derived from the tree after earlier tickets land | section-scoped anchors | TDD-able: anchor-first. | One ticket asserted three defects the preceding ticket had already fixed. |
| 5 | the charge names the central-property mutation and the delegate self-probes its finished work, reporting the result and adding the missing row on a silent green | section-scoped anchors | TDD-able: anchor-first. | Three fixture-shaped silent greens were caught only because the charge said run the mutation, not reason about it. |
| 5 | the coordinator probe differs in kind and in site, with omission/swap vocabulary in the charge template | section-scoped anchors; the kind sentence's hard-wrapped mutation entry moves with the hoist | TDD-able: anchor-first for the site sentence and vocabulary; the hoist is red-capable because the mutation table's anchor count fails when the sentence's wrapping moves, forcing the entry's edit in the same diff. | Three same-site probes were vacuous on the first try, and a vacuous probe is indistinguishable from a pass. |
| 5 | a family-extending charge names every registry found by tracing an existing sibling | section-scoped anchors | TDD-able: anchor-first. | FT152's canary work missed the classification registry and paid a repair round. |
| 5 | transient backups live inside the delegate's worktree under a unique name; restores name exact files, never globs | section-scoped anchors | TDD-able: anchor-first. | A stale scratchpad swept into a later delegate's restore glob clobbered four out-of-fence files. |
| 6 | `craft-spec`'s slicing section names value contracts at each fence by pointer to the ticket rule | section-scoped anchors | TDD-able: anchor-first, beside the untouched cross-pointer pair. | A restated procedure is the drift this repo's code standard names as a defect. |
| 6 | the canonical edge walk carries the process-boundary lifecycle class and the profile checklist carries its concrete entry | section-scoped anchors | TDD-able for both needles: anchor-first. The entries' wording beyond the pinned phrases is not TDD-able — review grades it. | Unit-level success hid defects that appeared only after a fresh process reloaded serialized state. |
| 7 | each of the 24 inventoried needles has a mutation-table row proving its diagnostic fires | mutation-table harness | TDD-able per needle: the table row lands before the skill text and fails on its anchor count until the sentence exists. Completeness of the set has no gate signal — no needle-counting assertion exists — and is review-graded against this spec's enumeration. | An anchor without a mutation proof can rot into an always-pass, which is the decay this spec exists to stop. |
| 7 | the two moved anchor artifacts keep biting after their moves | mutation-table harness | TDD-able: the placeholder swap reds the retired needle's row until it is replaced, and the kind sentence's hard-wrapped entry fails on the hoist until edited; the untouched fixtures run unchanged as the regression control. | A move that silently drops an anchor entry passes every presence check while deleting the enforcement. |
| edge of 7 | the example-agreement check fails closed when the begin/end markers are absent, duplicated, or enclose an empty block | example-agreement check | TDD-able: fixtures for each malformation assert the named red. | A fail-open extractor grades nothing forever — or grades the template block instead of the example — while reporting green. |
| edge of 7 | a wrapped fence or continuation assumptions line in the example reds the agreement check | example-agreement check | TDD-able: mutate the example into the wrapped form and assert the red. | Wrapping is the exact corpus failure, and the live parser truncates it silently rather than refusing. |
| edge of 7 | an example block ending at end-of-file without a trailing newline still parses identically | example-agreement check | TDD-able: a fixture with the marker block at EOF, final newline stripped. | Hand-edited markdown loses trailing newlines, and a last-line field silently vanishing from the parse would pass the presence anchors. |

Degenerate implementations are pinned per story. Teaching the fields in an
appendix while leaving the template block unchanged (1) fails the
section-scoped needles and the retired placeholder's mutation row — the
appendix route is closed by scoping, not by substring presence. A Good example
updated to *look* right but not parse, or shrunk to one row (1, 2), fails the
agreement check's distinct-ID and per-ID mutations-table assertions. Restating
the classification rule outside the procedure section (3) fails its scoped
needle; making it a step but not the opening one is the named review-graded
residual. A contracts step reduced to "check the fences agree" (4) fails the
four-fact needle. Hoisting the kind sentence by deleting it from the
spec-build lane (5) fails its mutation entry's anchor count. Restating the
full contracts procedure in `craft-spec` (6) is caught by review against the
cross-pointer precedent, not by the gate — named here as a residual. Adding
twenty needles and one mutation row (7) passes the gate and is caught only by
review against the enumeration — the honest limit of this build's
enforcement, priced under Out of scope.

### Edge inventory

- Error path — resolved by the fail-closed marker rows and the wrapped-form
  agreement rows; the anchor check's own error path is the existing harness's.
- Empty or absent input — resolved by the absent/duplicated/empty-marker row,
  which distinguishes a missing block from a present-but-empty one.
- Boundary values — resolved by the EOF-without-newline row and the
  single-line field shapes the parser enforces.
- Malformed input — resolved by the wrapped-fence and unlabeled-row agreement
  rows.
- Interrupted or partial state — **Won't handle**: the build lands as ordinary
  commits graded atomically by the gate; the only new runtime write is a temp
  spec-and-ticket pair inside the Go test's own temp directory, whose cleanup
  the test framework owns.
- Re-run idempotency — **Won't handle**: no skills-index or generated surface
  is touched (`index:` frontmatter is unchanged), and the new check is
  read-only over the tree, writing only its own temp directory per run.
- Hostile environment — **Won't handle, with the gap named**: the conformance
  harness reads tracked files with a bare read — no symlink stat, no
  special-file rejection — and this build inherits that posture for one more
  tracked path rather than hardening a shared helper mid-build; the hardening
  is priced under Out of scope. Special files under `tests/` fixtures and the
  temp directory are the test framework's domain.
- Process-boundary lifecycle — the class this spec *adds* to the walk; for
  this build itself, **Won't handle**: nothing this build ships serializes
  state for a later process — the agreement check's temp files live and die
  inside one test process.
- Profile shell-CLI checklist, remaining classes (paths with spaces, control
  bytes in git-sourced text, TTY detection, symlink invocation, worktree
  state, unquoted arguments, I/O pressure) — **Won't handle**: this build
  ships no command surface, flag, or prompt; its one new code path is a
  conformance check over tracked bytes and a test-owned temp directory.

## Out of scope

- The repair riders — model-comparison charges, inventory-currency repairs,
  shared-cache/hermetic repair distinctions — as one follow-up `craft-delegate`
  visit: 4 edits, 2 gate runs. Separate capability: specialized repair and
  experiment lanes, not the ticket contract.
- A needle-completeness assertion counting anchor entries against
  mutation-table rows — 3 edits, 2 gate runs. Separate capability: it would
  grade the harness itself, retrofits ~150 existing needles of which most
  have no mutation row today, and its absence is named honestly in story 7's
  rows rather than papered over.
- Hardening the conformance harness's file reads (symlink stat, special-file
  rejection) — 3 edits, 2 gate runs. Separate capability: a shared-helper
  posture change affecting every existing check, owned by `craft-gate`.
- Parser and validation enforcement — blocked-by grammar and graph walking,
  fence disjointness at assign, wrapped-field refusals, fence-fallback
  refusal, requirement-ID completeness, non-`R` range expansion — is FT174,
  already a roadmap row that depends on this spec's identifier decisions.
- The spec-optional route (FT180) and the refactor lane (FT108) consume these
  conventions; both are existing roadmap rows.
- Reformatting the nine older non-parseable ticket files under `specs/` is
  rejected by the source, not deferred.
- Building recomposition test harnesses for the process-boundary class in any
  project is per-project work; this spec only makes the class nameable.
- Adding `.agents/` to the reduced gate scope to cheapen iteration: a scope
  change is `craft-gate` work with its own honesty argument — 3 edits, 2 gate
  runs — and prior scope decisions (2026-08-01) deliberately kept the declared
  set minimal.
