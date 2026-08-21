# A capture line the parser cannot see is reported, not counted as zero

Status: implemented

Decision source: `roadmap/FT243.md` — the reviewed roadmap row "a capture entry the parser cannot see is reported as zero, not as a failure", drained 2026-08-21 and re-observed one drain later — plus the reviewer-confirmed decision of 2026-08-21 that the source's clause covers undated content as well as dated lines, taken in the `/bench-implement-spec` entry round.

Verification log: 2 iteration(s) to accept — the round returned three blocking findings and five non-blocking ones. Folded: the preamble anchor now excludes the scaffold's own worked-example heading (it had made DL27 and DL30 mutually unsatisfiable, because the shipped scaffold prints the marker below that line), and disposes of the no-real-heading and two-marker cases; `strings.TrimSpace` is pinned as the exact blank predicate with DL33 asserting it; DL32 pins a dated line flush at column one, which every other dated row left unguarded; the two-copies-of-the-literal count is corrected; the reason string is reworded to the writer's terms; and ticket 02's literal-uniqueness claim is scoped to Go sources. Round 2 verified each fold against the tree and returned no blocking finding; its three non-blocking ones are folded too — the Solution's window is qualified to "above the first heading", the body-content decision no longer claims a body walk for a malformed heading, and DL34 pins the second-marker case that had been disposed only in prose. The round's own partition answer — ticket 01 is a narrower capability that could ship on its own gate — is raised for reviewer disposition rather than folded.

## Problem

`capture/learnings.md` has one grammar: a `## <date> — <title> [open]` heading.
A writer who appends an entry as a bullet instead of a heading loses it. The
parser inspects only lines that start with `## `, so a bullet produces no entry,
no malformed record, and no parse failure. Every surface then agrees on a false
zero: `bench learnings` prints `learnings[0]` at exit 0, `bench roadmap
--context` reports `capture/learnings.md,parsed,1886` beside `learnings[0]` and
`parse_failures[0]`, and `bench status` shows zero open learnings.

The consequence is a silent drop. The drain reads that inventory as complete, so
it discards every entry the writer got wrong. The writer gets no signal, because
the file it appended to looks exactly like the file it read. This face occurred
twice: both entries of the 2026-08-21 drain were invisible, and two more
bullet-shaped entries were invisible one drain later.

Both observed occurrences were dated bullets, but the defect is not about dates.
A writer who appends an undated note below the scaffold's `<!-- entries below
-->` marker loses it the same way and for the same reason. The scaffold puts
that marker there to say where entries live; nothing reads it, so the promise it
makes to the writer is not kept.

## Solution

A capture source reports what it could not interpret. Two rules carry that.

A line that leads with a date but is not a well-formed dated heading becomes a
malformed record. A line below the entries marker that belongs to no entry
becomes a malformed record too. Both are the same kind of record the parser
already produces for a broken `## ` heading, and each names its source line and
its reason, so the writer can repair the exact line and the drain must verdict
it.

The two rules divide the file by what the writer declared. Above the marker the
journal is documentation — the scaffold's prose preamble and its worked example
live there — so only a dated line is judged. Below the marker and above the
first heading the writer said entries live here, so every line that is not one
is reported. From the first heading on, the existing walk already owns the
lines, and only the dated rule reaches inside them.

One parser carries both. `bench learnings`, `bench roadmap --context`, and
`bench status` all read that parser's malformed list today, so all three
surfaces stop reporting a false zero without a change of their own. The marker
becomes a parser-owned export, because the adoption scaffold and the
docs-reference check each hold that literal today, and neither is the rule that
will enforce it.

The quiet postures stay quiet: a freshly scaffolded journal, a drained journal
holding only its schema heading, a journal of well-formed entries, and a journal
with no marker at all read exactly as they do now.

## User stories

### Group 1 — a lost dated line becomes visible

Line: `opus` / medium. This is parser logic at a known seam, which the cached
routing prices at the cheap tier, but the row's whole value is the honesty of an
oracle surface, so it gets the mid tier. The reviewer accepted that deviation on
2026-08-21.

1. As a capture writer, I want a dated bullet I append to `capture/learnings.md`
   to be reported as malformed, so that my entry is never silently discarded.
2. As a drain reviewer, I want `bench learnings` to exit non-zero when the
   journal holds a dated line that is not a well-formed heading, so that a green
   read never certifies a lossy journal.
3. As a drain reviewer, I want each lost dated line to render as its own row
   naming its source line, so that I can repair the exact line.
4. As a drain reviewer, I want `bench roadmap --context` to carry each lost
   dated line as a `parse_failures` row sourced at `capture/learnings.md`, so
   that the inventory I drain from is honest.
5. As an operator, I want `bench status` to render the learnings component as
   unknown when the journal holds a lost dated line, so that the ambient
   dashboard never reports zero open learnings for a lossy journal.

### Group 2 — the mistake's shape does not decide whether it is seen

Line: `opus` / medium. The same seam and the same walk; these are the hostile
spellings of the same defect.

6. As a capture writer, I want a dated line to be recognized whatever list,
   quote, or heading marker precedes it, so that the shape of my mistake does
   not decide whether it is seen.
7. As a capture writer, I want a dated line separated from its marker by a
   non-ASCII space to still be recognized, so that a hand-edited NBSP cannot
   re-open the silent drop.
8. As a capture writer, I want a dated line inside an existing entry's body to
   be reported, so that appending a bullet under the last entry is caught like
   any other lost entry.
9. As a drain reviewer, I want a lost dated line's recorded text to carry no
   trailing carriage return, so that a CRLF journal cannot split its own field
   in `bench roadmap --context --full`.
10. As a capture writer, I want a dated line at end of file with no trailing
    newline to be reported, so that the last thing I appended is not the one
    thing that vanishes.

### Group 3 — content below the entries marker is accounted for

Line: `opus` / medium. The same parser walk gains a second rule, and the marker
it anchors on becomes one exported fact with three readers.

11. As a capture writer, I want an undated note I append below the entries
    marker to be reported as unaccounted content, so that a note I meant as an
    entry is never silently discarded.
12. As a drain reviewer, I want a contiguous run of unaccounted lines to be
    reported as one record naming its first line, so that one pasted mistake
    costs me one repair instead of one row per line.
13. As a capture writer, I want a blank or whitespace-only line below the marker
    to stay quiet, so that ordinary spacing is not reported as a defect.
14. As a capture writer, I want an undated line inside a well-formed entry's
    body to stay quiet, so that the entry body bullets the scaffold asks me for
    are not reported back at me.
15. As an operator, I want a journal that has no entries marker to keep exactly
    today's undated behavior, so that a pruned journal does not open red.
16. As a maintainer, I want the entries marker recognized only above the first
    real entry heading, and not closed by the scaffold's own worked example, so
    that the rule opens on a freshly adopted repo and a marker pasted into an
    entry's body cannot re-anchor it over the whole file.
17. As a drain reviewer, I want an unaccounted run's recorded text to carry no
    trailing carriage return, so that a CRLF journal cannot split its own field.

### Group 4 — one date grammar, one marker, and the good paths do not move

Line: `opus` / medium. Regression breadth for the states that must stay exactly
as they are, plus the assertions that keep the shared facts from drifting apart.

18. As a maintainer, I want a digit-shaped but non-calendar date to be judged
    the same way by the heading rule and by the line rule, so that the two
    cannot drift onto separate definitions of a date.
19. As a maintainer, I want the entries marker to have one definition that the
    parser, the adoption scaffold, and the docs-reference check all read, so
    that the boundary the parser enforces is the boundary the scaffold ships.
20. As a maintainer, I want a `## ` heading to keep exactly its current
    disposition, so that the two existing malformed reasons do not change.
21. As an adopter, I want a freshly scaffolded journal to read as the empty
    table at exit 0, so that a new repo does not open red.
22. As an operator, I want a drained journal holding only its schema heading to
    read as the empty table at exit 0, so that an empty inbox stays quiet.
23. As an operator, I want a well-formed open entry to keep its row at exit 0,
    so that the change adds a diagnostic without moving the good path.
24. As an adopter, I want the scaffold's `## <date>` worked example to stay a
    non-entry, so that the shipped template is still documentation.

### Group 5 — reviewed exclusions

Line: `opus` / medium. Each exclusion is pinned or excused rather than left to
drift.

25. As a reviewer, I want an undated, non-heading line above the entries marker
    to stay quiet, so that the scaffold's prose preamble does not red every
    adopted repo.
26. As a reviewer, I want a control byte inside a reported line to meet the
    pre-existing TOON refusal, so that this change adds no new sink exposure.
27. As a reviewer, I want a dated bullet inside a fenced code block to be
    reported like any other, so that the journal grammar gains no fence concept
    that could hide a real entry.

## Implementation decisions

**One parser owns both rules.** `learnings.Parse` is the single source
`bench learnings`, `bench roadmap --context`, and `bench status` already read —
directly for the first two, through `roadmap.learningCount` for the third. Both
rules land there and nowhere else. No consumer gains a second derivation of
"what the journal could not interpret".

**The record type does not change.** A lost line becomes a
`learnings.Malformed` with `Reason`, `Raw`, and `Line` — the same shape the two
existing malformed reasons use. Every consumer's rendering therefore needs no
edit, which is what keeps the production change to the parser plus the scaffold
that now reads its marker.

**The reason names the repair.** A lost dated line carries the reason
`dated learning entry is not a heading`. Unaccounted content below the marker
carries `learning content below the entries marker is not an entry`. Both are
distinct from `malformed learning heading` and
`dated learning heading must end with [open]`, so a reader can tell "you used the
wrong marker" from "your heading is broken" from "this belongs to no entry".

**The date predicate is extracted, not copied.** `isDatedHeading` already
decides whether ten bytes are `YYYY-MM-DD`. That test becomes one named helper
that both the heading rule and the new line rule call, so the grammar of a date
has one definition. It stays digit-shape only; neither rule gains a calendar
parse.

**A line leads with a date when a bounded prefix walk reaches one.** The walk
strips, from the start of the line, a run of runes where each rune either
satisfies `unicode.IsSpace` or is one of the markdown marker bytes `-`, `*`,
`+`, `>`, and `#`. What remains must open with a date. `unicode.IsSpace` is the
exact predicate, not a family: it carries the whole `White_Space` property, so
every `Zs` separator — U+00A0, U+2000 through U+200A, U+3000 — is stripped, and
the zero-width characters U+200B and U+FEFF are not. A line already starting
`## ` is excluded from the rule outright, because the two existing heading
reasons own it and a second record would double-report it.

This is deliberately wider than `learnings.isSpace`, which is ASCII-only because
it serves a TOON whitespace class. The prefix walk serves a hand-edited markdown
file instead, so it does not borrow that predicate.

**The entries marker becomes an exported parser constant.** The literal
`<!-- entries below -->` lives in `internal/adopt`'s scaffold and again in
`internal/conformance`'s docs-reference check today. `internal/learnings`
exports it and both read it from there, so the boundary the parser enforces is
the boundary the scaffold ships and the docs check splits on.

**The marker is a preamble terminator, not a search hit.** It opens the
unaccounted rule only when it appears above the first *real* entry heading — a
line starting `## ` that `isTemplatePlaceholder` does not claim. The exclusion
is load-bearing rather than tidy: the shipped scaffold prints its worked example
`## <date> - <short title>  [open]` above the marker, so an anchor that counted
that line would never open the rule on the one journal shape it exists to serve.
A journal that holds a marker and no real entry heading at all opens the rule
from the marker to end of file. A journal that holds no marker keeps today's
behavior for undated content entirely: only the dated rule applies to it. This
is what keeps a pruned journal and every pre-marker journal quiet. A second
marker below the first is an ordinary line and joins the unaccounted run, which
is the honest reading: the boundary was already declared once. The marker is
matched after `strings.TrimSpace`, so a marker carrying trailing whitespace still
opens the rule: requiring an exact match would let one invisible space disable
the diagnostic silently, which is the failure class this spec exists to close.

**Unaccounted content is a contiguous run, reported once.** Below the marker and
above the first real entry heading, each line is classified in order: a blank
line is quiet and ends any open run, a dated line takes the dated reason and
ends any open run, and any other line opens or continues a run. Blank is
`strings.TrimSpace(line) == ""` exactly, so a whitespace-only line is blank:
reporting a run whose recorded text is invisible spaces gives the writer nothing
to repair. Each run yields
one record carrying the run's first line and that line's 1-based number. A
pasted paragraph is one mistake, so it is one repair.

**Body content is never unaccounted.** Once a heading line appears, no line
before the next `## ` is unaccounted, so the entry body bullets the scaffold
asks for are untouched by the new rule. Only the dated rule reaches inside a
body.

**No fence concept.** A dated bullet inside a fenced code block is reported like
any other. Adding a fence concept to suppress the diagnostic would trade a lost
entry for a hidden one, and the journal grammar has no fence today.

**The recorded text is the line with its trailing carriage return removed**, the
same normalization the existing malformed-heading record applies before storing
`Raw`.

**Records stay in ascending source-line order** across all four reasons, so the
rendered rows read against the file top to bottom.

## Testing decisions

- **External behavior a good test exercises:** feed the parser a journal in the
  exact shape a writer produced — a pre-drain journal of dated bullets, and a
  freshly scaffolded journal with a note appended below its marker — and require
  the lost content to appear as records with their lines; then drive the three
  consumer surfaces and require each to stop reporting zero.
- **Seams and their prior art:** `internal/learnings` unit tests on `Parse`
  (prior art: the `TestRows` table); the `internal/learnings` command test with
  checked-in byte-exact stdout fixtures (prior art:
  `TestRefusalEvidencePreservesPrimaryAndDisclosesOnlyRepairableMalformed` and
  its `testdata/candidate-*.stdout` files); `internal/adopt`'s scaffold test for
  the marker's one definition; `internal/roadmap`'s context test for the
  `parse_failures` block; `internal/status`'s status test for the drain row.
- **Gate seam:** the kit gate's `test` phase (`go test -count=1 ./...`), which
  runs all five packages. No new phase, and no conformance registry entry.

### Seam diagram

    trigger: a writer appends to capture/learnings.md
        │
        ▼
    journal bytes ──▶ [ learnings.Parse ] ──▶ entries[] + malformed[]
                          ▲ anchors on learnings.JournalEntriesMarker
                          ◀ tests attach here: feed bytes, assert both lists
        │
        ├──▶ [ learnings.Command ]   ──▶ bench learnings rows + exit code
        ├──▶ [ roadmap.BuildContext ] ──▶ parse_failures block
        └──▶ [ roadmap.learningCount ] ──▶ bench status drain row state
                          ◀ tests attach here: assert each surface stops reporting zero

    internal/adopt scaffold ──▶ reads learnings.JournalEntriesMarker
                          ◀ test attaches here: scaffolded bytes parse with zero records

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DL1 | 1 | `Parse` returns one malformed record for a `- <date> — <title>` bullet under the schema heading | `internal/learnings` unit test on `Parse` | a parser that inspects only `## ` lines returns zero records, which is the shipped defect |
| DL2 | 3 | the record carries the bullet's 1-based source line | `internal/learnings` unit test on `Parse` | a record without its line cannot be repaired, and a hardcoded zero passes any presence-only assertion |
| DL3 | 1 | the record's reason is `dated learning entry is not a heading` | `internal/learnings` unit test on `Parse` | a reason borrowed from the heading path misdescribes the defect and sends the writer to the wrong repair |
| DL4 | 2 | `bench learnings` exits 1 on the two-bullet journal that produced the 2026-08-21 drop | `internal/learnings` command test | that journal exits 0 today, and only the exit code stops a drain from trusting the read |
| DL5 | 3 | `bench learnings` renders one `line <n>` row per lost dated line | `internal/learnings` command test with a byte-exact stdout fixture | a byte-exact fixture refuses a rendering that drops one of the two lost entries |
| DL6 | 6 | a dated line preceded by `-`, `*`, `+`, `>`, or `#` is reported | `internal/learnings` unit test on `Parse` | a detector anchored to `- ` alone misses the other four markers |
| DL7 | 7 | a dated line whose marker is separated by U+00A0 is reported | `internal/learnings` unit test on `Parse` | ASCII-only whitespace stripping re-opens the silent drop on hand-edited markdown |
| DL8 | 7 | a dated line whose marker is separated by U+3000 is reported, and one separated by U+200B is not | `internal/learnings` unit test on `Parse` | pins `unicode.IsSpace` as the predicate, which a literal U+00A0 special case and a zero-width-inclusive walk both fail |
| DL9 | 8 | a dated bullet appended inside an existing open entry's body is reported | `internal/learnings` unit test on `Parse` | body lines are consumed wholesale until the next `## `, so an appended entry is the likeliest live miss |
| DL10 | 20 | `## broken` still yields exactly one record reading `malformed learning heading` | `internal/learnings` unit test on `Parse` | a rewrite routing every non-heading line through the new rules would relabel or double-report it |
| DL11 | 20 | `## 2026-01-01 — x` without `[open]` still yields exactly one record reading `dated learning heading must end with [open]` | `internal/learnings` unit test on `Parse` | the same rewrite would take this line through the new rules too |
| DL12 | 18 | `- 2026-88-88 — x` is reported as a lost dated line | `internal/learnings` unit test on `Parse` | the heading rule is digit-shape only, so a line rule written with a calendar parse diverges from it while every valid-date row still passes |
| DL13 | 4 | `bench roadmap --context` renders one `parse_failures` row sourced at `capture/learnings.md` for the lost line | `internal/roadmap` context test | the drain reads this block as its complete inventory |
| DL14 | 5 | `bench status` renders the drain row's learnings component as `unknown (capture/learnings.md is malformed)` | `internal/status` status test | a fabricated `0 open learning(s)` is the exact face the decision source names |
| DL15 | 21, 24 | a freshly scaffolded journal renders the empty table at exit 0 | `internal/learnings` command test | a rule that flags preamble prose or the worked example reds every adopted repo on its first read |
| DL16 | 22 | a journal holding only its schema heading renders the empty table at exit 0 | `internal/learnings` command test | the drained inbox is the quiet posture, and a diagnostic here would be a false alarm every drain |
| DL17 | 23 | a well-formed open entry keeps its `date,title` row at exit 0 | `internal/learnings` command test | a detector that fires on a heading's own date destroys the good path |
| DL18 | 25 | an undated, non-heading line above the entries marker produces no record | `internal/learnings` unit test on `Parse` | pins the reviewed exclusion, so a rule anchored at the file start instead of the marker goes red here |
| DL19 | 3 | records render in ascending source-line order across all four reasons | `internal/learnings` unit test on `Parse` | a record appended after the walk makes a mixed journal's rows unreadable against the file |
| DL20 | 10 | a dated bullet on a final line with no trailing newline is reported | `internal/learnings` unit test on `Parse` | the split's last element is the standard off-by-one in a line walker |
| DL21 | 9 | a CRLF-terminated dated bullet's record carries the line without its trailing carriage return | `internal/learnings` unit test on `Parse` | `toon.Representable` permits a carriage return, so an unstripped record splits its own field in `bench roadmap --context --full` |
| DL22 | 11 | an undated line below the entries marker yields one record reading `learning content below the entries marker is not an entry` | `internal/learnings` unit test on `Parse` | the line produces nothing today, which is the widened half of the defect |
| DL23 | 12 | three contiguous undated lines below the marker yield exactly one record, carrying the first line's text and its 1-based number | `internal/learnings` unit test on `Parse` | a per-line record turns one pasted mistake into three rows the drain must verdict separately |
| DL24 | 13 | a blank line below the marker produces no record | `internal/learnings` unit test on `Parse` | a rule that counts blank lines reds every scaffolded journal, whose marker is followed by exactly one |
| DL33 | 13 | a whitespace-only line below the marker produces no record | `internal/learnings` unit test on `Parse` | pins `strings.TrimSpace` as the blank predicate, which a `line == ""` implementation fails by reporting a run whose recorded text is invisible |
| DL25 | 14 | an undated bullet inside a well-formed open entry's body produces no record | `internal/learnings` unit test on `Parse` | the scaffold asks for exactly those bullets, so reporting them reds the documented entry shape |
| DL26 | 15 | a journal with no entries marker produces no record for its undated lines | `internal/learnings` unit test on `Parse` | a rule anchored at the schema heading instead of the marker reds every pruned and pre-marker journal |
| DL27 | 16 | a marker appearing below a real entry heading does not open the rule | `internal/learnings` unit test on `Parse` | a whole-file search for the literal re-anchors on a marker pasted into an entry body |
| DL34 | 16 | a second marker below the first joins the open run's record rather than starting a new region | `internal/learnings` unit test on `Parse` | a `strings.LastIndex` anchor restarts the preamble at the second marker and silently drops every line between the two |
| DL35 | 12 | a run still open at end of file, in a journal with no trailing newline, still yields its record | `internal/learnings` unit test on `Parse` | every other run is closed by the empty final element a trailing newline produces, so the walk's own end-of-input flush is otherwise unasserted and can be deleted whole |
| DL36 | 16 | a marker carrying trailing whitespace still opens the rule | `internal/learnings` unit test on `Parse` | an exact-match anchor lets one invisible space disable the diagnostic silently, which is the failure class this spec exists to close |
| DL28 | 17 | a CRLF-terminated unaccounted run's record carries the line without its trailing carriage return | `internal/learnings` unit test on `Parse` | the new rule is a second `Raw` writer, so it needs its own normalization or it splits its own field |
| DL29 | 19 | the bytes `internal/adopt` scaffolds parse with zero records, and the scaffold's marker is `learnings.JournalEntriesMarker` | `internal/adopt` scaffold test | a scaffold holding its own copy of the literal drifts from the parser's boundary, and the first fresh repo then reads red |
| DL30 | 11 | `bench learnings` exits 1 and renders a `line <n>` row for a scaffolded journal with one undated note appended below its marker | `internal/learnings` command test with a byte-exact stdout fixture | that journal exits 0 today; it also goes red on an anchor that lets the scaffold's own worked example close the preamble, which is the cheapest wrong reading of the marker rule |
| DL32 | 6 | a dated line flush at column one, with no leading marker or whitespace, is reported | `internal/learnings` unit test on `Parse` | a prefix walk that requires a non-empty marker run passes every other dated row while the plainest lost entry of all stays silent |
| DL31 | 27 | a dated bullet inside a fenced code block below the marker is still reported | `internal/learnings` unit test on `Parse` | pins the no-fence-concept decision, so a later fence suppression has to be a deliberate spec change |

Not covered: story 26 — the exposure is pre-existing and identical on the
existing malformed-heading path; this spec changes neither the sink nor its
refusal, so a row here would assert someone else's behavior.

### Edge inventory

The walk covers the shell-CLI hostile-input classes that reach a markdown line
parser: hand-edited files with no trailing newline (DL20), CRLF line endings
(DL21, DL28), hand-edited files with no trailing newline while an unaccounted run
is still open (DL35), non-ASCII whitespace in hand-edited markdown (DL7, DL8), control
bytes a sink permits but cannot survive (excluded below), and the absent-versus-empty
pair (unchanged — classification runs before the parser and its cases are already
asserted). The classes about paths, processes, worktrees, and TTYs do not reach
this seam: `Parse` takes bytes and touches no filesystem.

**Won't handle** — undated, non-heading content above the entries marker — the
shipped scaffold's prose preamble and worked example have exactly that shape,
and the surviving in-scope caller is `bench learnings` on a freshly adopted
repo, which must read green.

**Won't handle** — undated content in a journal that carries no entries marker —
without the writer's declared boundary there is no honest way to tell preamble
from stray content, and the surviving in-scope caller is this repo's own pruned
journal, which must keep reading green.

**Won't handle** — a control byte inside a reported line's recorded text — the
existing malformed-heading path carries the identical exposure and this spec
changes neither the sink nor `toon.Table`'s refusal; the surviving in-scope
caller is `bench roadmap --context --full`.

**Won't handle** — a zero-width character (U+200B, U+FEFF) standing between a
marker and a date — it is not `White_Space`, and treating an invisible character
as a separator would make the predicate unreadable to the writer repairing the
line; the surviving in-scope caller is DL8, which pins the exclusion.

**Won't handle** — `roadmap.learningCount` zeroing the well-formed open count
whenever any record exists — that fail-closed posture is pre-existing, and the
surviving in-scope caller is `bench status`, which renders unknown rather than a
number.

## Ownership fences

- `internal/learnings/learnings.go`
- `internal/learnings/learnings_test.go`
- `internal/learnings/testdata/`
- `internal/adopt/init.go`
- `internal/adopt/adopt_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/roadmap/context_test.go`
- `internal/status/status_test.go`

No production file outside `internal/learnings/learnings.go` and
`internal/adopt/init.go` is in the fence. If the build finds a consumer that
needs a production edit to surface a record, that is a seam the spec did not
anticipate: stop and report rather than widen the fence.

## Out of scope

- **Adopt the entries marker in this repo's own `capture/learnings.md`.** The
  kit's journal has been pruned past its marker, so the unaccounted rule is
  inert here until it is restored. Restoring it also narrows
  `internal/conformance`'s stale-slash-reference scan from the whole journal to
  the preamble alone, so it must land together with that check's re-anchoring.
  2 edits, 2 gate runs.
- **Apply the account-for-your-content property to `capture/retros/*.md`.** A
  separate parser with a separate grammar and its own consumer set.
  4 edits, 3 gate runs.
- **Report the well-formed open count alongside a malformed state in
  `bench status`.** Changes `roadmap.Drain`'s shape and the status row's
  rendering contract. 2 edits, 2 gate runs.

## Further notes

The observed red the build must flip, reproduced 2026-08-21 in a throwaway
repository holding the pre-drain journal: `bench learnings` printed
`learnings[0]{date,title}:` and exited 0 while the file held two dated entries.

The unaccounted rule has no matching observed red, because this repo's journal
carries no marker. Its red is the freshly scaffolded journal of DL30, which the
build produces from `internal/adopt`'s own scaffold rather than from a
hand-written fixture.
