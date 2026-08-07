# one-shot-feedback

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-07 — scope confirmed as the pre-assignment breakdown review, retro repair attribution with its drain tally, and the exact-predicate extension to the falsification charge and grill; the mechanical `assign` producer-path check was explicitly cut.

## Problem

The kit's pre-build pipeline aims every ticket at landing green in one pass, but
its reviews cluster after composition, where every accepted finding is already a
repair loop. The review obligations that would catch slicing defects before code
exist — scattered per-field in `craft-tickets` with no named moment or owner — so
they demonstrably fail to fire (a checked-in ticket fenced a producer file the
tree contradicts; an acceptance row hid a disjunction behind "stale or absent").
And the metric this all serves — repair rounds per ticket trending toward zero —
is not recorded anywhere, so prose changes meant to move it are steered on
anecdote.

## Solution

Close the feedback loop at both ends. Before any assignment, a named fresh-context
breakdown review grades the ticket files against the obligations `craft-tickets`
already states, including a producer-path witness against the tree. After every
promoted build, the retro attributes each repair round to a causal phase and the
`/bench-what-next` drain reports the running one-shot tally, so the trend becomes
a tracked number. Upstream, the spec falsification charge and the grill both
refuse outcome labels, and the falsification charge asks the story-partition
question, so ambiguity is caught in the phase that created it.

## User stories

The stories partition into four disjoint prose fences sharing only the anchors
registry's append-only row list. Per `craft-spec`'s partition rule that is a
split signal, and the reviewer chose the bundle (2026-08-07): one capability —
close the pre-build feedback loop — with the genuinely separate mechanical
`assign` check split out to Out of scope. The bundle is recorded, not defaulted.

1. As the build coordinator, after writing ticket files and before the first
   assignment, I charge one fresh read-only delegate with the breakdown review —
   the consolidated target list `craft-tickets` names — and when the harness
   cannot spawn one, the pass runs inline and is flagged in the build plan, so
   slicing defects become pre-code reslices instead of repair rounds.
   Ownership fence: `.agents/commands/bench-implement-spec.md`,
   `.agents/skills/bench-craft-tickets/SKILL.md`, `projects/benchkit.md`,
   `internal/anchors/registry_data.go`.
   Line: opus (claude mid binding) / medium. Reviewer-set routing (2026-08-07),
   overriding the doc-authoring leverage default; the story's registry data rows
   ride the same line.
2. As the reviewer draining capture, I read each promoted build's retro
   attribution table — one row per ticket, repair rounds and causal phase from
   the fixed vocabulary the retro template solely owns — and the drain's exit
   reports the one-shot tally, so I can see whether prose changes move the trend.
   Ownership fence: `.agents/commands/bench-final-check.md`,
   `.agents/commands/bench-what-next.md`, `internal/retros/`,
   `internal/anchors/registry_data.go`.
   Line: opus (claude mid binding) / medium. Same reviewer-set routing.
3. As the spec author, the falsification pass I charge asks two more questions —
   does any behavior, red signal, or decision answer name an outcome family
   instead of an exact predicate, and do the stories partition into disjoint
   package or fence sets such that a narrower capability could ship on its own
   gate — so label-shaped rows and defaulted bundles are caught before sign-off.
   Ownership fence: `.agents/commands/bench-write-spec.md`,
   `internal/anchors/registry_data.go`.
   Line: opus (claude mid binding) / medium. Same reviewer-set routing.
4. As the reviewer being grilled, each decision closes with the answer restated
   as the exact predicate it fixes, never an outcome label, so ambiguity cannot
   enter the pipeline at its source.
   Ownership fence: `.agents/skills/bench-craft-grill/SKILL.md`,
   `internal/anchors/registry_data.go`.
   Line: opus (claude mid binding) / medium. Same reviewer-set routing.

Build venue (reviewer decision, 2026-08-07): the orchestrating session itself
performs this build's pre-assignment breakdown review and the delegate-claim
reviews; write delegates still author the code, and the composed candidate still
takes the standard pre-promotion review route.

`internal/anchors/registry_data.go` is the one file every fence shares; it is an
append-only row list, so the build sequences those appends (or routes them
through one owning ticket) rather than fanning four writers into one file —
`craft-tickets` owns that build-time call, and the shared fence is recorded here
so it cannot surprise the slicing.

## Implementation decisions

- **Gate attachment is the anchors registry, and the needle is the predicate.**
  Each behavior's exact needle string is fixed in the coverage map below;
  deleting or paraphrasing the guidance later turns the conformance phase red.
  `Require` matching is case-sensitive whole-file substring after whitespace
  collapse, so each needle is chosen as a full operative clause — long enough
  that a file containing it states the rule — never a category label.
- **The breakdown review is a named moment, not new obligations.** `craft-tickets`
  gains a section (`## Review the breakdown before assignment`) naming the moment
  (ticket files written, nothing assigned), the owner (one fresh read-only
  delegate), and the charge — pointing by name at the six per-field obligations
  the skill already states (the `Integration surfaces:` `none` re-search, the
  keep-together split attempt, label-shaped rows and closure tokens, mutation
  subject-not-assertion honesty, `covers` honesty, fence-versus-advertisement) —
  plus the one genuinely new item: the producer-path witness, verifying each
  `Contracts:` and `Integration surfaces:` producer path exists in the tree and
  holds the named value. `/bench-implement-spec` adds the step between ticket
  derivation and lifecycle start; findings are reslices repaired before
  `bench spec build start`. When the harness cannot spawn a delegate,
  `craft-delegate`'s capability-aware policy applies and the pass runs inline,
  flagged in the build plan, never silently skipped.
- **Attribution rides the retro template, which solely owns the vocabulary.**
  `/bench-final-check`'s retro gains a `## Repair attribution` heading with one
  table row per ticket: repair rounds landed, and per round one cause from the
  vocabulary `shaping-ambiguity`, `spec-row`, `ticket-slicing`, `tree-drift`,
  `delegate-error`, `other`; a zero-round ticket records `none`. The template is
  the vocabulary's single source: `/bench-what-next` names "the cause vocabulary
  the retro template owns" and never restates the values, so the two files
  cannot drift. Placement: the heading sits before `## Agent-experience
  improvements`. No anchors kind expresses ordering and section resolution skips
  fenced headings, so placement is a reviewer-graded exception backed by the
  `internal/retros` regression fixture below, which proves the chosen placement
  is invisible to the recommendations parser.
- **The drain reports the tally.** `/bench-what-next` step 4 gains one duty: when
  drained retros carry attribution tables, the exit reports tickets total,
  one-shots, and per-cause counts, reading causes only from the drained tables.
  Reporting only — no new roadmap grammar, no CLI change.
- **Falsification charge and grill extend in place.** `/bench-write-spec` step 9
  adds the two questions to the existing charge sentence; `craft-grill` adds one
  discipline bullet. No new files, no new skills.
- **The cached routing row.** `projects/benchkit.md`'s Lines section gains the
  breakdown-review pass routing (mid model, medium effort, one iteration,
  read-only, standing grant like the falsification pass) so charge time never
  re-derives it. The binding table itself is untouched.

## Testing decisions

- The external behavior under test is presence and permanence of guidance: each
  behavior's quoted needle is required by the anchors registry, and the
  conformance phase (`go test ./internal/conformance -run '^TestRootConformance$'`)
  is where absence goes red. Semantic quality of the prose stays with review —
  the gate proves the operative clause exists, review grades its surroundings.
- **Degenerate implementation, named.** The cheapest wrong implementation of
  stories 1, 3, and 4 pastes each needle verbatim while gutting or contradicting
  the guidance around it; no anchors row can red on it, because `Require` is
  presence-only. The map accepts this as its stated exception: the mitigation is
  needle choice (each needle is the operative clause itself, so the paste *is*
  the rule) plus review, and the exception is recorded here rather than implied.
  The composition degenerate for story 2 — the template emitting one cause
  vocabulary while the drain tallies another — is killed structurally: the
  template solely owns the vocabulary and the drain reads causes only from
  drained tables, so there is no second list to drift (RA1's needle pins the
  vocabulary line; RA3's pins the read-from-tables duty).
- Seam with prior art: the workflow-anchors registry already guards exactly this
  class of drift for these same files — `bench-final-check.md`'s existing retro
  headings each carry a `Require` row, and RA1 extends that enumerated set.
- The retro-parser non-interference claim is asserted at the existing
  `internal/retros` seam: a regression fixture carrying `## Repair attribution`
  before the improvements heading proves `Recommendations` returns unchanged
  units. It cannot start red (the parser already keys on the exact improvements
  heading) and is landed as a regression guard, honestly classified.
- The gate seam that observes the whole feature is the conformance phase of
  `bench gate`; `bench coverage --check specs/one-shot-feedback/spec.md` guards
  this map itself.

### Seam diagram

    trigger: conformance phase of `bench gate`
        │
        ▼
    kit guidance files  ──▶  [ workflow-anchors registry check ]  ──▶  green, or a named missing-anchor diagnostic
                                 ◀ tests attach here: add the Require needle, run conformance
                                   before the prose lands, observe the missing-anchor red

    trigger: `/bench-final-check` writes a retro
        │
        ▼
    retro body with `## Repair attribution`  ──▶  [ retros.Recommendations parser ]  ──▶  unchanged recommendation units
                                 ◀ tests attach here: fixture with the new heading in
                                   `internal/retros` asserts the returned units are unchanged

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| BR1 | 1 | `/bench-implement-spec` carries the needle `one fresh read-only delegate grades the ticket breakdown before any assignment; a harness that cannot spawn one runs the pass inline and flags it in the build plan` between ticket derivation and lifecycle start | anchors registry needle on the command file | Require row lands first; `go test ./internal/conformance -run '^TestRootConformance$'` reds with the missing-anchor diagnostic until the prose lands | a build following the command cannot skip a step the file must state verbatim to stay green |
| BR2 | 1 | `craft-tickets` carries the section heading `## Review the breakdown before assignment` and the needle `every named producer path exists in the tree and holds the named value`, with the charge pointing at all six enumerated per-field obligations | anchors registry needles on the skill file | same needle-first conformance red per needle | the consolidated charge is what makes the six scattered obligations fire; the witness needle is the one new obligation stated verbatim |
| BR3 | 1 | `projects/benchkit.md` Lines carries the needle `Ticket-breakdown review pass` with its routing row | anchors registry needle on the profile | same needle-first conformance red | an unrouted delegate gets an improvised line; the cached row is the enforcement |
| RA1 | 2 | the retro template carries `## Repair attribution` and the vocabulary needle `shaping-ambiguity`, `spec-row`, `ticket-slicing`, `tree-drift`, `delegate-error`, `other` as one line the template solely owns | anchors registry needles on the final-check command | same needle-first conformance red per needle | attribution absent from the template means the metric is never recorded; the vocabulary needle pins the single source |
| RA2 | 2 | the retros parser returns unchanged recommendation units for a body carrying `## Repair attribution` before the improvements heading | `internal/retros` package test | already covered — the parser keys on the exact improvements heading, so the assertion cannot start red; the regression fixture lands with RA1 anyway | a future parser rewrite that keys on position instead of the heading breaks the fixture |
| RA3 | 2 | `/bench-what-next` carries the needle `report tickets total, one-shots, and per-cause counts, reading causes only from the drained tables` | anchors registry needle on the what-next command | same needle-first conformance red | a drain that never reports the tally leaves the metric unread; reading only from tables is the anti-drift half of the clause |
| FC1 | 3 | the step-9 charge carries the needle `does any behavior, red signal, or decision answer name an outcome family instead of an exact predicate` | anchors registry needle on the write-spec command | same needle-first conformance red | the charge is the only fresh context that grades the draft; an unasked question is an unguarded miss class |
| FC2 | 3 | the step-9 charge carries the needle `could a narrower capability ship on its own gate` | anchors registry needle on the write-spec command | same needle-first conformance red | partition today is self-graded by the author at story-locking; the charge is the adversarial moment |
| GP1 | 4 | `craft-grill` carries the needle `close each decision by restating the answer as the exact predicate it fixes — never an outcome label` | anchors registry needle on the grill skill | same needle-first conformance red | shaping is where label-shaped answers enter; the close rule is the earliest catch point |

### Edge inventory

- **Error path — harness cannot spawn the breakdown delegate:** the inline
  fallback is inside BR1's needle text, so its absence reds conformance.
- **Empty/absent — a build with zero repair rounds:** the vocabulary's `none`
  value covers it; a table of all-`none` rows is the one-shot ideal and tallies
  cleanly (RA1).
- **Empty/absent — a retro predating the template change:** RA3's needle scopes
  the tally to drained retros that carry attribution tables, so old retros drain
  exactly as today. **Won't handle** beyond that — no back-fill of historical
  retros; none are pending.
- **Malformed input — a cause value outside the fixed vocabulary:** **Won't
  handle** — the retro is reviewed prose and the drain is a reviewer-approved
  batch; a novel cause is a review catch, not a parser's.
- **Malformed input — non-ASCII whitespace inside a needle phrase in a target
  file:** **Won't handle.** The anchors matcher collapses whitespace with
  `strings.Fields`, which treats NBSP and its relatives as whitespace, so a
  space-corrupted phrase still satisfies its needle — the check neither catches
  nor false-fails on it. Hand-corrupted guidance prose is review's to catch;
  no mechanical guard is claimed.
- **Re-run idempotency:** the retro is rewrite-in-full by existing contract;
  attribution rides that unchanged. Anchors rows are static data with no state.
- **Boundary values, interrupted state, process-boundary lifecycle, hostile
  environment:** **Won't handle** — every deliverable is static guidance prose
  plus one static registry row set; no runtime state crosses a process boundary.

## Out of scope

- **Mechanical producer-path witness in `bench spec build assign`** — the Go
  lifecycle check refusing a ticket whose `Contracts:` producer path is absent
  from the tree. Separate capability with its own gate story; ~6 edits, 3 gate
  runs. The breakdown review's manual witness covers the intent meanwhile.
- **CLI-computed one-shot tally** (`bench roadmap` deriving the trend from retro
  tables mechanically). Separate capability; ~8 edits, 3 gate runs. The drain's
  reported tally is the manual bridge.
- **A falsification pass at `/bench-shape-idea` exit.** Separate capability;
  ~2 edits, 1 gate run. The grill's predicate close (story 4) covers the
  highest-leverage slice of it.
- **An anchors kind that expresses ordering or reaches inside fenced code
  blocks** (would make RA1's placement gate-observable). Separate capability in
  the anchors matcher; ~4 edits, 2 gate runs. The `internal/retros` fixture and
  review cover placement meanwhile.
