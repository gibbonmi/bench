# One block reader for authored Markdown

Status: staged

Roadmap: FT358

Decision source: `roadmap/FT358.md`, a named reviewed artifact from drain `d-0bca6e72fedd`.

Verification log: 2 iteration(s) to accept — iteration 1 ran five parallel reviewers (Sonnet xhigh, Sonnet high, Opus xhigh, Fable high, and Fable xhigh). They returned 41 findings: 37 confirmed, 3 partly right, and 1 false. The fold added the rule-delta table and closed decisions 8 to 15. Iteration 2 ran Sonnet high, which closed 36 of 40 ids. The fold took the 4 partial ids and 7 new findings, with closed decision 16.

Amendment source: the independent source review of `37ac80b8`, findings MB-R1 to MB-R3.

Amendment verification log: MB-SPEC-R2 accepted SHA `c0e3ebdb676f09a621ffe0d7e1c4a3a38e31230d1e22a6975c8eda4f4251c249` with no material findings before ticket amendments. MB-TICKETS-R1 accepted the ticket graph at `c2e4719d05b523ac910b05cd69bbc92ddbc330eb` with no material findings.

Amendment status: spec and ticket amendments accepted, with MB-SPEC-R2 before ticket edits and MB-TICKETS-R1 after the planning checkpoint.

## Problem

The detector inventory below names the fence-rule copies, and they disagree.
`prose.stripFences` accepts tilde fences and closes by character and run length.
The other module detectors toggle state or match exact fence lines.
These production scanners have no fence rule:
`coverage.parse`, `learnings.Parse`, `retros.Parse`, `spec.FenceTokens`, and
`deriveImplemented`. A fenced example in a coverage map, a learning, a
retrospective, or a spec therefore reads as live grammar.

The `internal/spec` module tests a fence on a trimmed line in `LiveSpecSlugs` and
on an untrimmed line in `metadata`. An indented fence therefore hides a spec path
from one function and shows a `Status:` line to the other. Two comment rules and
two frontmatter rules also exist. A fix to one copy leaves the other copies wrong.

## Solution

One block reader, `internal/markdown`, owns the four block rules: frontmatter,
fenced blocks, HTML comments, and H2 headings. Each grammar module reads its lines
from the block reader and keeps only its own grammar. `internal/prose` stays the
owner of the sentence and paragraph rules, and it reads its prose lines from the
block reader.

The fence edges have one table test in the block reader. A conformance check,
`markdown-block-owner`, refuses a block-rule literal in production Go outside the
block reader. A new string-literal detector therefore reds the gate.

## User stories

Line: opus / high.
Implementation-line reason: MB-C1 is the hardest chunk, because the block reader must keep the prose findings and the anchor rejoin rule exact. MB-C4 is the second hard chunk, because the anchor walk maps runes to source positions. Each seam exists, and the package tests observe each row. On 2026-10-03 the reviewer approved Opus at xhigh effort for tickets 1 and 8, and high effort for each other ticket.
Harder chunks: MB-C1, MB-C4.

### The block reader owns the block rules

1. As a module author, I want one reader that classifies each line, so that no module derives a block rule again.
2. As a writer, I want every module to read a `~~~` block as fenced, so that no two modules disagree about code.
3. As a writer, I want a closer run at least as long as the opener, so that a shorter inner run stays inside.
4. As a writer, I want a closer to use the opener's character, so that a backtick line stays inside a tilde block.
5. As a writer, I want an opener after leading spaces or tabs to open a block, so that an indent changes nothing.
6. As a named-fence reader, I want the opener to carry its trimmed info string, so that I need no marker test.
7. As a writer, I want a run of fewer than three markers to open nothing, so that inline code stays body.
8. As a module, I want an unterminated fence to give one fault at its opener, so that I can refuse it.
9. As a skill author, I want the leading `---` block to read as frontmatter, so that no module reads it as grammar.
10. As a reader, I want an unterminated frontmatter block to give a fault on line 1, so that no text is swallowed.
11. As a writer, I want each comment span cut from the line text with line numbers kept, so that diagnostics name physical lines.
12. As a writer, I want a fence marker inside a comment to open no block, so that the comment hides the marker.
13. As the anchor evaluator, I want a rejoined comment removed too, so that decision DG43 holds under the block reader.
14. As a reader, I want an unterminated comment to give a fault at its opener, so that each module sees one answer.
15. As a section reader, I want a column-zero `## ` body line to be an H2 heading, so that each module reads one heading set.
16. As a section reader, I want a `## ` line inside a block to be no heading, so that a quoted example opens nothing.
17. As a section reader, I want each line to carry its open H2 section title, so that I need no walk.
18. As a writer of a CRLF file, I want carriage returns stripped and byte offsets kept, so that a span replacement stays exact.
19. As a module, I want empty content to give one empty body line, so that the empty case needs no branch.

### The prose mechanics check composes with the block reader

20. As the prose mechanics check, I want my lines from the block reader, so that I hold only the prose rules.
21. As the prose mechanics check, I want each delimiter fault to keep its kind and its line, so that the gate output does not change.
22. As the prose mechanics check, I want a rejoined comment to stay ungraded, so that prose and the anchors agree about comment text.
23. As the reviewer, I want the prose grade of the live tree unchanged, so that the migration moves no finding.

### The spec and coverage readers use the block reader

24. As the status reader, I want a tilde-fenced `Status:` line ignored, so that a quoted example keeps the status.
25. As the spec readers, I want `LiveSpecSlugs` and `metadata` to agree about an indented fence, so that the trim inconsistency closes.
26. As the landing, I want `Implemented` to ignore a fenced `Status: staged` copy, so that a quoting spec still flips.
27. As the fence-section reader, I want a fenced example to give no token, so that a quoted path grants no write.
28. As the fence-section reader, I want a fenced `## Ownership fences` line to open nothing, so that it declares no fences.
29. As the coverage check, I want fenced map rows and story numbers ignored, so that a quoted example adds nothing.
30. As the coverage check, I want a fenced `## ` line to leave the story list open, so that later stories count.

### The roadmap reader uses the block reader

31. As the roadmap index reader, I want a tilde-fenced row example ignored, so that it gives no row and no failure.
32. As the sequence reader, I want a tilde-fenced sequence heading to open nothing, so that the sequence stays put.
33. As the `Next:` grammar, I want a tilde-fenced `Next:` line to be no marker, so that it is no duplicate.

### The field scan and the ticket reader use the block reader

34. As the field scan, I want a field line inside a tilde fence marked fenced, so that a quoted field matches nothing.
35. As the ticket reader, I want an unterminated tilde fence refused, so that the refusal covers both marker characters.
36. As the ticket reader, I want a backtick line inside a four-backtick block to stay inside, so that the ticket stays balanced.

### The handoff document uses the block reader

37. As the handoff parser, I want a `## ` line inside a nested fence to stay body, so that it splits no section.
38. As a handoff writer, I want `OpenFence` to use the block reader's close rule, so that a mismatched closer is refused.
39. As the State scanner, I want `UnfencedLines` to skip a block that a longer run closes, so that a quote is no claim.

### The journal readers use the block reader

40. As the learnings reader, I want a fenced dated heading kept in its entry body, so that it splits no entry.
41. As the learnings reader, I want an unterminated fence reported as malformed, so that later entries do not vanish silently.
42. As the learnings reader, I want an unterminated comment reported as malformed, so that later entries do not vanish silently.
43. As the learnings reader, I want an unterminated frontmatter block reported as malformed, so that each fault kind has a record.
44. As the retrospective reader, I want a fenced required heading not to count, so that a quote does not satisfy the order.
45. As the recommendation reader, I want a fenced `## ` line to leave the improvements open, so that later items count.
46. As the recommendation reader, I want a fenced line to be no unit, so that a quoted item needs no `Feeds:` marker.

### The review record uses the block reader

47. As the review record reader, I want a named opener inside another fenced block ignored, so that a quoted record is not a duplicate.
48. As the review record reader, I want the unterminated and duplicate refusals kept, so that a broken record still refuses.
49. As the review record writer, I want `Render` to replace the span that the reader locates, so that the reader and the writer agree.

### The anchors use the block reader

50. As the anchor section walk, I want a tilde-fenced `## ` line to leave the section open, so that later needles stay in scope.
51. As the anchor section walk, I want a tilde-fenced `## <title>` line to open nothing, so that it is no duplicate.
52. As the anchor step walk, I want a tilde-fenced step opener to open nothing, so that a quoted step moves no body.
53. As the anchor evaluator, I want the block reader's comment spans with source positions kept, so that `Locate` names the same lines.
54. As the reviewer, I want every anchor in the live registry to keep its verdict, so that the migration moves no anchor.

### The skills index uses the block reader

55. As the skills index, I want `FrontmatterField` to read the block reader's frontmatter, so that a CRLF skill file indexes its fields.
56. As the skills index, I want an unclosed frontmatter block to give no field, so that swallowed body text authorizes nothing.

### The managed-block markers use the block reader

57. As `bench link`, I want a managed-block marker inside a tilde fence ignored, so that a quoted marker example is not rewritten.
58. As `bench unlink`, I want a tilde-fenced marker example kept, so that the strip removes only the live block.
59. As `bench link`, I want an unterminated tilde fence around Bench text refused, so that broken markers are never trusted.

### The gate refuses a new copy

60. As the reviewer, I want a check that reds a block-rule literal outside the block reader, so that no new detector lands.
61. As a test author, I want test files exempt from the check, so that a fixture can plant a fence.
62. As the reviewer, I want the check green on the live tree, so that the migration is complete when the check lands.

### Reviewed exclusions

63. As a maintainer, I want the block reader to omit CommonMark's indented code block, so that each module keeps its own indented-line rule.
64. As a maintainer, I want the block reader to omit container fences, so that no module gains a list or quote grammar.
65. As a maintainer, I want the block reader to omit the backtick info-string rule, so that the live-tree prose grade stays exact.
66. As a maintainer, I want each H3 boundary rule to stay in its module, so that the block reader owns only the H2 rule.

### Amendments from the consultation and the first review

67. As `bench link`, I want an unterminated comment around Bench text refused, so that a hidden fence never frees a marker example.
68. As a handoff writer, I want an unterminated comment in the State file refused, so that it hides no heading.
69. As the coverage check, I want the historical marker read from raw text, so that a comment marker still exempts a spec.
70. As the learnings reader, I want the entries marker and the region below it read from raw text, so that DL31 holds.
71. As the anchor section walk, I want an indented `## ` line to open nothing, so that one H2 rule holds.
72. As `bench link`, I want a CRLF AGENTS.md rewrite to keep each carriage return, so that the file keeps its line endings.
73. As `bench link`, I want an unterminated comment without Bench text accepted, so that a stray comment does not block a link.
74. As the landing, I want a CRLF spec to flip and keep each carriage return, so that the file keeps its line endings.
75. As a fragment caller, I want no frontmatter pass on a section or State body, so that a leading `---` stays body.
76. As a module that ignores frontmatter, I want an unclosed `---` read as body, so that each refusal message stays exact.
77. As the reviewer, I want the three conformance readers on the block reader, so that no test helper keeps a block rule.

## Implementation decisions

### The block reader

`internal/markdown` is a new leaf package that imports only the standard library.
It has two entry points over one pipeline. `Read` reads a whole document and
applies the four rules below. `ReadFragment` reads a section or a State body and
skips the frontmatter rule, because frontmatter is a rule of document position.
Each one returns the classified lines,
the removed comment spans, and at most one fault.

The reader never refuses. Each
caller decides what a fault means for its own grammar.

Each line carries these facts:

- the 1-based physical line number and the byte offset of the line start
- the raw text without the newline and without one trailing carriage return
- the text with every removed comment byte cut out
- one block class: frontmatter, fence opener, fenced line, fence closer, or body
- the trimmed info string of a fence opener
- whether the line opens an H2 heading, and the heading title
- the title of the open H2 section, or no section before the first heading

The reader applies its rules in one fixed order. This order is the current order
of the two owners, `internal/prose` and `internal/anchors`, so neither one changes
its verdict on the live tree.

1. Frontmatter. A first line whose raw text is exactly `---` opens it. The next
   line whose raw text is exactly `---` closes it. The raw text has its carriage
   return stripped, so a CRLF file opens frontmatter. A leading or a trailing space
   opens nothing. With no closer, the reader reports a frontmatter fault on line 1
   and reads every line as body.
2. Comments. The reader removes each comment opener through the next closer over
   the text outside the frontmatter. A removal that rejoins a new comment removes
   that comment too, with the result of `stripCommentsMapped` in
   `internal/anchors/locate.go`. The earliest open opener wins, and the search
   for its closer starts after the opener's four bytes. One linear pass gives
   that result.
   The removed bytes include newlines, but each line keeps its number.

   An unterminated comment removes the rest of the document and gives a comment
   fault at its opener line.
3. Fences. On the comment-free text, the reader skips leading bytes that are a
   space (U+0020) or a tab (U+0009). A line whose next text is three or more
   backticks or tildes opens a fenced block. Any other leading space, such as
   U+00A0, opens nothing. The closer is a line whose marker run uses the same
   character and is at least as long.

   Text after the closer's run does not stop the close, as in `stripFences` in
   `internal/prose/prepare.go`. An unterminated block gives a fence fault at its
   opener line, and every later line is fenced.
4. H2 headings. A body line whose comment-free text starts with `## ` at column
   zero opens an H2 section. The title is the rest of the line, trimmed by
   `strings.TrimSpace`. That predicate is `unicode.IsSpace`, so U+00A0 trims and
   U+200B stays.

The reader reports at most one fault, in the order frontmatter, comment, fence.
The input three backticks, a newline, then a comment opener holds an open fence
and an open comment. The reader reports the comment fault at line 2 for it, as
`prose.prepare` does today. The package also exports the fence marker and the H2
prefix as constants. A writer that renders a fence or a heading uses them, so no
block literal stays outside the package.

### Raw text or comment-free text

The reader gives each line its raw text and its comment-free text. Each module
matches its own grammar on one of them, and this table pins the choice. A fence
marker and an H2 heading come from the comment-free text in every module.

| module | grammar text | why |
| --- | --- | --- |
| `internal/prose` | comment-free | the prose rules grade only visible text today |
| `internal/anchors` | comment-free | the evaluator hides commented needles today |
| `internal/spec` | raw | no spec grammar reads a comment today |
| `internal/coverage` | raw | the historical marker is itself a comment |
| `internal/roadmap` | raw | no roadmap grammar reads a comment today |
| `internal/maps` and `internal/tickets` | raw | no field grammar reads a comment today |
| `internal/handoffdoc` and `internal/handoff` | raw | no handoff grammar reads a comment today |
| `internal/learnings` | raw | the entries marker is itself a comment |
| `internal/retros` | raw | no retrospective grammar reads a comment today |
| `internal/reviewrecord` | raw | the payload bytes must stay exact |
| `internal/skillsindex` | raw | frontmatter holds no comment |
| `internal/adopt` | raw | each managed-block marker is itself a comment |

A rebuilder writes its output from the source bytes by the reader's line
offsets, never from the line text. The two rebuilders are `RewriteAgentsBlock`
and `StripAgentsBlock` in `internal/adopt/marker.go`, and `deriveImplemented` in
`internal/spec/spec.go`. Each carriage return therefore survives a rebuild.

### Fault posture

A module that refuses an unterminated fence today keeps that refusal and its exact
message bytes. These modules are `internal/prose`, `internal/tickets`,
`internal/handoffdoc`, `internal/reviewrecord`, and `internal/adopt`. Each one
now reads the reader's fence fault. Only `internal/prose` refuses all three kinds,
as it does today.

`internal/adopt` also refuses the reader's comment fault, with the message at
`internal/adopt/marker.go:68`. It refuses under the same condition as the fence
fault: the file holds Bench text. Its markers are comments, and the reader removes
comments before it classifies fences. An unterminated comment therefore hides each
later fence marker, and a fenced marker example would read as live.

`internal/handoff` refuses the reader's comment fault in `readStateFile` with a
new fault beside `faultStateFileFence`. The fault names the opener line. The
handoff file is untracked, so the prose grade never reads it.

`internal/learnings` adds a malformed record for each fault kind at the fault
line. The journal is untracked, so no tracked-path prose grade refuses it. The
prose mechanics check prints the journal's malformed records, so those records are
the journal's fail-closed channel. Every other module reads an unterminated block
to the end of the document. The prose grade refuses such a document on each
path that the prose grade covers.

Three modules ignore the frontmatter fault and read each line as body:
`internal/adopt`, `internal/tickets`, and `internal/reviewrecord`. This is
their behavior today, and it keeps each refusal message exact. MB95, MB96, and
MB97 cover them.

The fragment callers are `handoffdoc.OpenFence` and `handoffdoc.UnfencedLines`,
which read a State body, and the anchor scope walks `sectionRunesMapped` and
`stepRunesMapped`, which read a section body. Each other caller reads a whole
file through `Read`, `handoffdoc.Parse` and the anchor file read included.

### Module migrations

Each module replaces its own fence, frontmatter, comment, and H2 tests with the
reader's classification. It keeps its own grammar: the field table, the labels,
the row grammar, the section titles it selects, and any H3 boundary rule. The
exported signatures stay the same. `handoffdoc.OpenFence`, `handoffdoc.UnfencedLines`,
`spec.FenceTokens`, `spec.LiveSpecSlugs`, `anchors.MarkdownH2Sections`,
`anchors.StripHTMLComments`, `anchors.MarkdownNumberedSteps`, and
`skillsindex.FrontmatterField` keep their callers unchanged.

The `internal/spec` trim inconsistency closes by construction. `LiveSpecSlugs`,
`metadata`, `deriveImplemented`, and `FenceTokens` read one fence classification.
`FenceTokens` opens its section on the reader's H2 heading, so an indented
`## Ownership fences` line no longer opens it. It keeps its rule that any `#{2,} `
heading ends the section, applied to unfenced lines only.

`internal/anchors` asks the reader for the removed comment spans and the line
classes. It keeps its rune-origin mapping and its step grammar. Its section opener
compares the reader's H2 title with the anchor title, so an indented opener no
longer opens a section. `stripCommentsMapped` becomes a projection of the reader's
spans. Decision DG43 therefore holds under the reader, and the reader is the one
place that states it.

`internal/reviewrecord` finds its record as the one fence opener whose info string
equals the record name. `Render` replaces the span between that opener and its
closer by the reader's byte offsets.

`ParseTicket` calls the reader on its content for the fence fault. `FieldScan.Scan`
keeps its return shape, the lines and the duplicate diagnostics, so the decision
map output does not change.

`internal/roadmap` also migrates `sequenceBounds` and `sectionBounds` in
`internal/roadmap/sequence_projection.go`. The parser, `ProjectSequence`, and
`closeDependencies` then share the reader's H2 and fence classification.
The projection keeps its outcome grammar, and the dependency closure keeps its
table grammar. MB99 and MB100 observe the composed writers at their public seams.

### Structure budgets

The structure check reads its file limit in `internal/structure/structure.go`.
A file at or over that limit does not grow. Each migration ends with no net
growth in these production files:

- `internal/spec/spec.go` and `internal/coverage/coverage.go`
- `internal/skillsindex/skillsindex.go`
- `internal/anchors/locate.go` and `internal/handoffdoc/document.go`
- `internal/roadmap/tree.go`

The three conformance files of ticket 11 also end with no net growth:
`skill_description_budgets_test.go`, `skills_index_checks_test.go`,
and `docs_workflow_checks_test.go`. A crowded directory takes no new file:
`internal/adopt`, `internal/anchors`, `internal/reviewrecord`, and
`internal/roadmap`. Each ticket names its test file in `Writes:`. The tests use these existing homes:

- The spec rows go in `internal/spec/fences_test.go`, because `internal/spec/spec_test.go` is over the file limit.
- The marker rows go in `internal/adopt/link_plan_test.go`, because `internal/adopt/adopt_test.go` is over the file limit.
- The learnings rows go in `internal/learnings/learnings_test.go`, because `internal/learnings/entry_test.go` tests only the formatter.
- The roadmap rows go in `internal/roadmap/commitment_test.go`, beside its sequence and dependency closure tests.

The roadmap tests do not recreate the deleted helper file or grow
`internal/roadmap/tree_test.go`.

A ticket that cannot meet this rule stops and reports, because a budget grant is
a reviewer decision.

### The copy guard

`markdown-block-owner` is a new `go-source` conformance check. It follows
`checkGitPlumbingOwner` in `internal/conformance/git_plumbing_owner_test.go`.
It parses each non-test Go file under `cmd/` and `internal/`, outside
`internal/markdown`. It refuses a string literal whose value contains a run of
three backticks or three tildes. It also refuses a literal whose value is exactly `---`,
`<!--`, `-->`, or `## `.

A rune-slice conversion of such a literal counts, because the check reads the
string literal itself. A rune literal does not count, so the check stops a new
string-literal detector only. The check carries a bite proof in its own file, as
`TestWaitDeadlineLiteralsBites` does. It adds no canary fixture. The
`git-plumbing-owner` check is graded by the fixture
`tests/canary/package-core-guard/git-flag-retyped`, but this check's bite proof
plants each literal kind itself.

Ticket 11 also moves three conformance readers onto the block reader:
`skillDescriptionFolded` in `skill_description_budgets_test.go`,
`frontmatterHasKey` in `skills_index_checks_test.go`, and `h2Headings` in
`docs_workflow_checks_test.go`. Each one is a test file, so the check exempts it,
but each one copies a block rule.

Ticket 11 also maintains the live-tree classification in
`internal/conformance/tier_live_tree_test.go`. Its planned live-tree assertion
enters `classifiedLiveTreeTests` when `testChecksLiveTree` detects its body.
The hidden-reader metadata check and the stale-classification test remain intact.
MB101 observes the detected assertion through `conformance-meta`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| MB-C1 / `1-add-block-reader.md` | The block reader exists, and the prose check reads its prose lines from it. | MB1, MB2, MB3, MB4, MB5, MB6, MB7, MB8, MB9, MB10, MB11, MB12, MB13, MB14, MB15, MB16, MB17, MB18, MB19, MB20, MB21, MB22, MB23, MB71, MB72, MB73, MB74, MB75, MB76, MB77, MB78, MB79, MB88, MB89, MB90, MB92 | `bench test --package ./internal/markdown`, `bench test --package ./internal/prose`, `bench test --check prose-mechanics` | yes |
| MB-C2 / `2-read-spec-and-coverage-blocks.md`, `3-read-roadmap-blocks.md`, `4-read-field-scan-blocks.md` | The spec, coverage, roadmap, field-scan, and ticket readers use the block reader. | MB24, MB25, MB26, MB27, MB28, MB29, MB30, MB31, MB32, MB33, MB34, MB35, MB36, MB37, MB80, MB81, MB83, MB95, MB99, MB100 | `bench test --package ./internal/spec`, `bench test --package ./internal/coverage`, `bench test --package ./internal/roadmap`, `bench test --package ./internal/maps`, `bench test --package ./internal/tickets` | no |
| MB-C3 / `5-read-handoff-blocks.md`, `6-read-journal-blocks.md`, `9-read-skills-frontmatter.md` | The handoff, journal, and skills-index readers use the block reader. | MB38, MB39, MB40, MB41, MB42, MB43, MB44, MB45, MB46, MB47, MB59, MB60, MB82, MB84, MB93 | `bench test --package ./internal/handoffdoc`, `bench test --package ./internal/handoff`, `bench test --package ./internal/learnings`, `bench test --package ./internal/retros`, `bench test --package ./internal/skillsindex` | no |
| MB-C4 / `7-read-review-record-blocks.md`, `8-read-anchor-blocks.md` | The review record and the anchors use the block reader. | MB48, MB49, MB50, MB51, MB52, MB53, MB54, MB55, MB56, MB85, MB94, MB96 | `bench test --package ./internal/reviewrecord`, `bench test --package ./internal/anchors`, `bench test --check docs-currency-workflow` | yes |
| MB-C5 / `10-read-agents-marker-blocks.md`, `11-forbid-block-rule-copies.md` | The managed-block markers use the block reader, and the gate refuses a new copy. | MB57, MB58, MB61, MB70, MB86, MB87, MB62, MB63, MB64, MB65, MB66, MB67, MB68, MB69, MB91, MB97, MB98, MB101 | `bench test --package ./internal/adopt`, `bench test --check markdown-block-owner`, `bench test --check skill-description-budgets`, `bench test --check skills-index-command-adapters`, `bench test --check conformance-meta` | no |

MB-C1 creates the seam that every later ticket consumes, so its chunk review
closes before any other ticket starts. Tickets 2, 3, 4, 6, 8, and 11 each name the
five command-binding files, because build preflight binds those files to their
packages. No ticket expects to edit them, but the shared names make those
tickets run in series. Ticket 11 waits for each module ticket, because the check
reds the live tree until the last copy goes.

## Testing decisions

- A good block-reader test passes exact bytes to the reader and compares the classified lines, the spans, and the fault. It does not read a module's grammar.
- The fence edges have one table test, `TestReadFenceEdges` in `internal/markdown`. A module test proves only that the module reads the reader. One tilde case and one nested-run case at each module seam are sufficient.
- Each module row uses the module's existing exported seam, so the row observes the composed path to the real producer. The prior art for each seam was read in this session: `TestFindings`, `TestFenceTokensEndTheSectionAtAnyHeading`, `TestParseTicketHostileInput`, `TestParseKeepsAFencedHeadingInsideState`, `TestOpenFenceNamesTheOpeningLine`, `TestRenderRefusesADuplicateFence`, `TestLocateStripsRejoinedComments`, and `TestFrontmatterFieldRequiresCompleteLeadingFence`.
- A module test that needs a tilde case adds it beside the existing backtick case in the named test file. No ticket copies a module's old detector into a test.
- The package tests and conformance checks run inside the ordinary `test` phase, as `internal/gate` and `projects/benchkit.md` specify.

### Seam diagram

    trigger: each grammar module, the prose mechanics check, the anchor evaluator
        │
        ▼
    document bytes  ──▶  [ internal/markdown block reader ]  ──▶  classified lines, comment spans, one fault
                      ◀ tests attach here: TestReadFenceEdges passes exact bytes
                        and compares each line class, span, and fault
        │
        ▼
    classified lines  ──▶  [ module grammar: spec, coverage, roadmap, maps, tickets,
                             handoffdoc, learnings, retros, reviewrecord, anchors,
                             skillsindex, adopt, prose ]  ──▶  the module's own result
                      ◀ tests attach here: each module's exported seam, with one
                        tilde case and one nested-run case

    production Go source  ──▶  [ markdown-block-owner check ]  ──▶  one diagnostic per literal
                      ◀ tests attach here: the bite proof plants each literal kind

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| MB1 | 2 | The lines between a `~~~` opener and a `~~~` closer read as fenced lines | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A backtick-only detector reads the inner line as a body line. |
| MB2 | 3 | Inside a four-backtick block, a three-backtick line reads as a fenced line, and a four-backtick line closes the block | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A toggle detector closes the block at the shorter run. |
| MB3 | 4 | Inside a `~~~` block, a three-backtick line reads as a fenced line | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A detector that ignores the marker character closes the block on the backtick run. |
| MB4 | 5 | An opener after two spaces and an opener after one tab each open a fenced block | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A column-zero prefix test reads the indented opener as a body line. |
| MB5 | 6 | The opener line ```` ```go  ```` carries the info string `go` | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A reader without the info string forces the review record to keep its own marker test. |
| MB6 | 7 | A line that starts with two backticks reads as a body line and opens no block | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A detector that counts any backtick run opens a block on inline code. |
| MB7 | 8 | An opener with no closer gives one fence fault at the opener line, and each later line reads as fenced | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A reader that drops the fault lets a module read a truncated document as complete. |
| MB8 | 9 | Lines 1 to 3 of `---`, `index: x`, `---` read as frontmatter, and a `---` on line 5 opens nothing | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that opens frontmatter on any `---` line hides a horizontal rule's paragraph. |
| MB74 | 9 | A first line ` ---` with a leading space opens no frontmatter | planned TestReadBlockRules in internal/markdown/markdown_test.go | The old `TrimSpace` test in `prose.prepare` opens frontmatter on the indented line. |
| MB75 | 9 | A first line `---` with a trailing space opens no frontmatter | planned TestReadBlockRules in internal/markdown/markdown_test.go | The old `TrimSpace` test in `prose.prepare` opens frontmatter on the padded line. |
| MB9 | 10 | A first line `---` with no closer gives a frontmatter fault on line 1, and each line reads as body | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that swallows the document hides every heading below the opener. |
| MB10 | 11 | A comment from line 2 to line 4 is cut from the text of those lines, and line 5 keeps the number 5 | planned TestReadBlockRules in internal/markdown/markdown_test.go | A removal that drops the newlines shifts each later line number. |
| MB11 | 12 | A three-backtick line inside a comment opens no fenced block | planned TestReadBlockRules in internal/markdown/markdown_test.go | A fence pass that runs before the comment pass opens a block inside the comment. |
| MB12 | 13 | The text `<!<!-- z -->-- needle -->` leaves no `needle` in the line text | planned TestReadBlockRules in internal/markdown/markdown_test.go | A single-pass comment strip leaves the rejoined comment's text visible. |
| MB71 | 14 | The input of three backticks, a newline, and a comment opener gives the comment fault at line 2 and no fence fault | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that reports the fence fault first changes the prose finding for this input. |
| MB79 | 13 | Reading 1 MiB that holds 10,000 comments allocates at most four times the input size | planned TestReadBlockRules in internal/markdown/markdown_test.go | The per-comment tail copy of `stripCommentsMapped` allocates near 10,000 times the input size. |
| MB88 | 13 | A line of a comment opener, ` a `, a second opener, ` b `, a closer, and ` c` keeps only ` c`, because the first opener runs through the first closer | planned TestReadBlockRules in internal/markdown/markdown_test.go | A stack that closes the innermost opener leaves the first opener open and gives a comment fault. |
| MB89 | 14 | A comment opener followed only by `>` gives a comment fault at line 1, because the closer search starts after the opener's four bytes | planned TestReadBlockRules in internal/markdown/markdown_test.go | A search that starts inside the opener reads its own dashes as a closer and removes nothing. |
| MB90 | 10 | The lines `---`, `a: b`, and `---` with a trailing space give a frontmatter fault on line 1 | planned TestReadBlockRules in internal/markdown/markdown_test.go | The old `TrimSpace` closer test closes the block on the padded line. |
| MB92 | 75 | `ReadFragment` on `---`, three backticks, `x`, `---` classifies line 2 as a fence opener and gives a fence fault at line 2 | planned TestReadBlockRules in internal/markdown/markdown_test.go | A frontmatter pass on a fragment swallows lines 1 to 4 and hides the open fence. |
| MB13 | 14 | A comment with no closer gives a comment fault at its opener line, and each later line has empty text | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that drops the fault lets a module read the commented rest as live text. |
| MB14 | 15 | The line `## Title  ` is an H2 heading with the title `Title` | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that keeps the trailing spaces gives a title that no exact compare matches. |
| MB15 | 15 | The lines `  ## Title`, `##Title`, and `### Title` are not H2 headings | planned TestReadBlockRules in internal/markdown/markdown_test.go | A trimmed prefix test reads an indented line or an H3 line as an H2 heading. |
| MB16 | 16 | A `## Title` line inside a fenced block is not an H2 heading | planned TestReadBlockRules in internal/markdown/markdown_test.go | A heading test that skips the fence class opens a section inside a quoted example. |
| MB17 | 17 | Each line after `## A` carries the section `A` until the line `## B`, and each line before `## A` carries no section | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that resets the section on each line gives the body lines no section. |
| MB18 | 18 | For the content `a\r\n```\r\n`, line 2 opens a fenced block, its raw text holds no carriage return, and its byte offset is 3 | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that keeps the carriage return misreads the marker, and a wrong offset breaks a span replacement. |
| MB19 | 19 | Empty content gives one empty body line and no fault | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that returns no line breaks a caller that reads line 1. |
| MB72 | 19 | The content `a` with no newline gives one body line `a` and no fault | planned TestReadBlockRules in internal/markdown/markdown_test.go | A reader that drops a last line without a newline loses the document's only line. |
| MB73 | 3 | Inside a three-backtick block, the line ```` ```go ```` closes the block | planned TestReadFenceEdges in internal/markdown/markdown_test.go | A closer that refuses trailing text keeps the block open and fences the rest of the document. |
| MB76 | 15 | The lines `## Title` with a trailing U+00A0 and `##` plus a space, a U+00A0, and `Title` each give the title `Title` | planned TestReadBlockRules in internal/markdown/markdown_test.go | An ASCII-only trim keeps the U+00A0 and gives a title that no exact compare matches. |
| MB77 | 15 | The line `## Title` with a trailing U+200B gives a title that ends with U+200B | planned TestReadBlockRules in internal/markdown/markdown_test.go | A trim widened to zero-width characters hides a byte that the writer must repair. |
| MB78 | 5 | A line of U+00A0 and three backticks opens no fenced block | planned TestReadFenceEdges in internal/markdown/markdown_test.go | An indent test by `unicode.IsSpace` opens a block that `prose.prepare` does not open today. |
| MB20 | 20 | The prose findings for the PD13 and PD21 documents stay the same | `internal/prose/parse_test.go` (`TestFindings`) | A composition that skips one block rule grades a fenced, commented, or frontmatter line. |
| MB21 | 21 | An unterminated frontmatter block, comment, and fence each give the same finding kind and line as before | `internal/prose/parse_test.go` (`TestFindings`) | A wrong fault mapping changes the gate output for an unterminated delimiter. |
| MB22 | 22 | A 30-word sentence inside a rejoined comment gives no sentence finding | planned TestFindingsSkipsARejoinedComment in internal/prose/prose_test.go | The old per-line comment strip grades the rejoined comment's words. |
| MB23 | 23 | The prose grade of the live tree gives no finding at the ticket tip | `internal/conformance/prose_mechanics_test.go` (`TestProseMechanicsHoldsOnTheLiveTree`) | A rule change that reclassifies a live line gives a finding on a tracked document. |
| MB24 | 24 | A spec with a real `Status: staged` line and a later tilde-fenced `Status: implemented` line reads as staged | planned TestMetadataSkipsAFencedStatus in internal/spec/fences_test.go | The old backtick-only toggle reads the tilde example last and returns implemented. |
| MB25 | 25 | A two-space-indented fence that holds a spec path and a column-zero `Status: implemented` line hides both from `LiveSpecSlugs` and from `Facts` | planned TestMetadataSkipsAFencedStatus in internal/spec/fences_test.go | The old untrimmed test in `metadata` reads the indented fence's status line. |
| MB80 | 74 | `Implemented` on a CRLF spec flips the status line and keeps each carriage return | planned TestImplementedSkipsAFencedStatus in internal/spec/fences_test.go | The old `stagedRe` match fails on the carriage return, and a rebuild from line text drops it. |
| MB26 | 26 | `Implemented` on a spec with one real `Status: staged` line and one fenced copy flips only the real line | planned TestImplementedSkipsAFencedStatus in internal/spec/fences_test.go | The old scan counts two matches and refuses the spec. |
| MB27 | 27 | A fenced example that holds a backticked path inside `## Ownership fences` gives no token | planned TestFenceTokensSkipAFencedExample in internal/spec/fences_test.go | The old fence-free scan returns the quoted path as a fence token. |
| MB28 | 28 | A fenced `## Ownership fences` line reports no declared section | planned TestFenceTokensSkipAFencedExample in internal/spec/fences_test.go | The old scan opens the section on the quoted heading. |
| MB29 | 29 | A fenced map row and a fenced story-number line add no data row and no story | planned TestParseSkipsFencedCoverageLines in internal/coverage/coverage_test.go | The old fence-free scan counts the quoted row and the quoted story. |
| MB83 | 69 | A spec whose line is the coverage historical marker comment keeps its historical exemption | planned TestParseSkipsFencedCoverageLines in internal/coverage/coverage_test.go | A parse of comment-free text drops the marker and grades the historical map. |
| MB30 | 30 | A story after a fenced `## Example` line inside `## User stories` counts as a declared story | planned TestParseSkipsFencedCoverageLines in internal/coverage/coverage_test.go | The old scan ends the story list at the quoted heading. |
| MB31 | 31 | A tilde-fenced roadmap row example gives no row and no malformed-row failure | planned TestParseDocumentSkipsTildeFences in internal/roadmap/commitment_test.go | `ParseDocument` matches a row on every line today, fenced or not, so the old scan parses the quoted row. |
| MB81 | 31 | A backtick-fenced roadmap row example gives no row and no malformed-row failure | planned TestParseDocumentSkipsTildeFences in internal/roadmap/commitment_test.go | `ParseDocument` matches a row on every line today, so the old scan parses the quoted row. |
| MB32 | 32 | A tilde-fenced `## Recommended sequence` line opens no sequence section | planned TestParseDocumentSkipsTildeFences in internal/roadmap/commitment_test.go | The old backtick-only toggle starts the sequence at the quoted heading. |
| MB33 | 33 | A tilde-fenced `Next:` line after a real `Next:` line gives no duplicate diagnostic | planned TestParseDocumentSkipsTildeFences in internal/roadmap/commitment_test.go | The old backtick-only toggle counts the quoted marker as a second marker. |
| MB99 | 32 | `ProjectSequence` replaces only the live sequence after a tilde-fenced sequence heading | planned TestProjectSequenceSkipsFencedHeadings in internal/roadmap/commitment_test.go | A surviving backtick-only `sectionBounds` replaces the quoted section instead of the real sequence. |
| MB100 | 16 | `Close` removes FT1 from only the live dependency section after a tilde-fenced dependency heading | planned TestCommitmentDependencyClosureSkipsFencedHeadings in internal/roadmap/commitment_test.go | A surviving backtick-only `sectionBounds` selects the quoted section and leaves the live dependency table unchanged. |
| MB34 | 34 | A `Writes:` line inside a tilde fence is marked fenced and matches no field | planned TestFieldScanMarksTildeFencedLines in internal/maps/tickets_test.go | The old backtick-only toggle matches the quoted field. |
| MB35 | 35 | A ticket with an unterminated `~~~` block gives the `unterminated fence` diagnostic | planned TestParseTicketReadsTheBlockReaderFault in internal/tickets/fence_test.go | The old count tests only backtick lines and gives no diagnostic. |
| MB36 | 36 | A ticket with a four-backtick block that holds a three-backtick line gives no `unterminated fence` diagnostic | planned TestParseTicketReadsTheBlockReaderFault in internal/tickets/fence_test.go | The old odd-count test reports the balanced block as unterminated. |
| MB37 | 36 | A `Covers:` line inside the four-backtick block of MB36 sets no `Covers` value | planned TestParseTicketReadsTheBlockReaderFault in internal/tickets/fence_test.go | The old toggle closes the block at the inner run and reads the quoted field. |
| MB95 | 76 | A ticket whose first line is an unclosed `---` gives the same diagnostics as it does today | planned TestParseTicketReadsTheBlockReaderFault in internal/tickets/fence_test.go | A ticket reader that refuses the frontmatter fault adds a diagnostic that no ticket grammar names. |
| MB38 | 37 | A State body with a four-backtick block that holds a three-backtick line and a `## X` line parses as one section | planned TestParseKeepsANestedFenceInsideState in internal/handoffdoc/document_test.go | The old toggle closes the block at the inner run and splits a section at `## X`. |
| MB39 | 38 | `OpenFence` on the text `~~~` then a three-backtick line reports line 1 | planned TestOpenFenceUsesTheCloseRule in internal/handoffdoc/document_test.go | The old toggle pairs the two markers and reports no open fence. |
| MB40 | 39 | `UnfencedLines` yields no line from inside a four-backtick block that holds a three-backtick line | planned TestOpenFenceUsesTheCloseRule in internal/handoffdoc/document_test.go | The old toggle yields the lines after the inner run. |
| MB93 | 75 | `OpenFence` on the State body of MB92 reports line 2 | planned TestOpenFenceUsesTheCloseRule in internal/handoffdoc/document_test.go | A caller that reads a State body through `Read` hides the open fence inside frontmatter. |
| MB82 | 68 | `readStateFile` on a State that holds an unterminated comment refuses with a fault that names the opener line | planned TestStateFileRefusesAnUnterminatedComment in internal/handoff/state_file_test.go | Without the fault, the comment hides each later heading from the State walk. |
| MB41 | 40 | An entry body with a fenced dated heading gives one entry whose body holds the fenced line | planned TestParseKeepsAFencedHeadingInTheBody in internal/learnings/learnings_test.go | The old scan splits a second entry at the quoted heading. |
| MB42 | 41 | An unterminated fence in an entry body gives the malformed reason `unterminated fenced block` at the opener line | planned TestParseReportsAnUnterminatedBlock in internal/learnings/learnings_test.go | Without the record, the fence swallows each later entry and no reader sees it. |
| MB43 | 42 | An unterminated comment gives the malformed reason `unterminated HTML comment` at the opener line | planned TestParseReportsAnUnterminatedBlock in internal/learnings/learnings_test.go | Without the record, the comment swallows each later entry and no reader sees it. |
| MB44 | 43 | An unterminated frontmatter block gives the malformed reason `unterminated frontmatter block` on line 1 | planned TestParseReportsAnUnterminatedBlock in internal/learnings/learnings_test.go | A fault kind without a record leaves one silent path. |
| MB84 | 70 | A non-entry line below the entries marker keeps its unaccounted-content record | `internal/learnings/learnings_test.go` (`TestParseReportsUnaccountedContentBelowTheEntriesMarker`) | A region read from comment-free text loses the marker line and drops the record. |
| MB45 | 44 | A retrospective whose only `## Outcome` line is fenced fails `Parse` with the missing-heading error | planned TestParseIgnoresAFencedHeading in internal/retros/retros_test.go | The old scan counts the quoted heading and accepts the retrospective. |
| MB46 | 45 | An item after a fenced `## Notes` line inside the improvements is a recommendation | planned TestRecommendationsSkipFencedLines in internal/retros/recommendations_test.go | The old scan ends the section at the quoted heading. |
| MB47 | 46 | A fenced `- item` line inside the improvements gives no recommendation | planned TestRecommendationsSkipFencedLines in internal/retros/recommendations_test.go | The old scan reads the quoted item as a unit without a `Feeds:` marker. |
| MB48 | 47 | A four-backtick block that quotes the record fence before the real record reads the real payload with no duplicate error | planned TestRecordSkipsAQuotedFence in internal/reviewrecord/record_test.go | The old exact-line scan reports the quoted opener as a duplicate. |
| MB49 | 48 | A record fence with no closer refuses as unterminated | `internal/reviewrecord/record_test.go` (`TestRenderRefusesAnUnterminatedFence`) | A migration that drops the fault reads a truncated record. |
| MB50 | 48 | Two record fences refuse as a duplicate | `internal/reviewrecord/record_test.go` (`TestRenderRefusesADuplicateFence`) | A migration that takes the first opener hides the second record. |
| MB51 | 49 | `Render` keeps the prose around the replaced fence byte for byte | `internal/reviewrecord/record_test.go` (`TestRenderKeepsProseAroundTheFence`) | A wrong offset from the reader cuts or duplicates the surrounding prose. |
| MB96 | 76 | A review record whose first line is an unclosed `---` reads its payload as it does today | planned TestRecordSkipsAQuotedFence in internal/reviewrecord/record_test.go | A record reader that refuses the frontmatter fault drops a readable record. |
| MB52 | 50 | A needle after a tilde-fenced `## Other` line inside section `S` locates inside `S` | planned TestLocateReadsTheBlockReader in internal/anchors/locate_test.go | The old backtick-only walk ends `S` at the quoted heading. |
| MB53 | 51 | A tilde-fenced `## S` line before the real `## S` gives a section count of 1 | planned TestLocateReadsTheBlockReader in internal/anchors/locate_test.go | The old walk counts the quoted heading and reports a duplicate section. |
| MB85 | 71 | An indented `  ## S` line gives no section `S` | planned TestLocateReadsTheBlockReader in internal/anchors/locate_test.go | The old opener compares `trimSpaceRunes` of the line and opens the section on the indented line. |
| MB54 | 52 | A tilde-fenced step opener gives no step body | planned TestLocateReadsTheBlockReader in internal/anchors/locate_test.go | The old walk opens the step at the quoted line. |
| MB94 | 75 | `MarkdownNumberedSteps` on a section body whose first line is `---` finds step 1 | planned TestLocateReadsTheBlockReader in internal/anchors/locate_test.go | A scope walk that reads a section body through `Read` swallows the step as frontmatter. |
| MB55 | 53 | The rejoined comment needle locates nothing, and the text after it keeps line 3 | `internal/anchors/locate_test.go` (`TestLocateStripsRejoinedComments`) | A strip that loses the rejoin or the source position breaks DG43 or the line number. |
| MB56 | 54 | Each anchor in the live registry keeps its verdict at the ticket tip | review-owned: `bench test --check docs-currency-workflow` green at the tip | A rule change that moves a live section boundary reds or frees a live anchor. |
| MB57 | 57 | A marker example inside a `~~~` block stays, and `RewriteAgentsBlock` replaces only the live block | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | The old backtick-only toggle reads the quoted marker as a live start marker. |
| MB58 | 58 | `StripAgentsBlock` keeps a tilde-fenced marker example and removes the live block | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | The old toggle strips the quoted example. |
| MB59 | 55 | A CRLF skill file `---`, `index: x`, `---` gives the field value `x` | planned TestFrontmatterFieldReadsCRLF in internal/skillsindex/frontmatter_test.go | The old exact compare reads `---` plus a carriage return as no opener. |
| MB60 | 56 | An opener with no closer gives an empty field value | `internal/skillsindex/skillsindex_test.go` (`TestFrontmatterFieldRequiresCompleteLeadingFence`) | A migration that reads the fault as frontmatter authorizes swallowed text. |
| MB61 | 59 | An unterminated `~~~` block around Bench text refuses with the message `conflict: AGENTS.md has an unclosed code fence around Bench markers; marker detection cannot be trusted` | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | The old backtick-only toggle sees no open fence and trusts the markers. |
| MB70 | 67 | A file that holds only an unterminated HTML comment opener and a fenced marker example refuses with the unclosed-fence message of MB61 | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | A fence-only fault check sees no open fence and rewrites the hidden example as the live block. |
| MB86 | 73 | A file that holds an unterminated comment opener and no Bench text rewrites with no refusal | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | A comment fault that ignores the Bench-text condition refuses each file with a stray comment opener. |
| MB87 | 72 | `RewriteAgentsBlock` on a CRLF file keeps each carriage return outside the replaced block | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | A rebuild from the reader's line text drops each carriage return. |
| MB91 | 72 | `StripAgentsBlock` on a CRLF file keeps each carriage return outside the removed block | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | A rebuild from the reader's line text drops each carriage return. |
| MB97 | 76 | An AGENTS.md whose first line is an unclosed `---` rewrites as it does today, with no refusal | planned TestMarkersSkipTildeFences in internal/adopt/link_plan_test.go | A marker reader that refuses the frontmatter fault blocks a link on a file that links today. |
| MB62 | 60 | A production Go literal that contains three backticks after other text gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits the backtick kind, or tests only a prefix, lets an embedded fence survive. |
| MB63 | 60 | A production Go literal that contains three tildes after other text gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits the tilde kind lets a tilde detector survive. |
| MB64 | 60 | A production Go literal equal to `---` gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits frontmatter lets a frontmatter reader survive. |
| MB65 | 60 | A production Go literal equal to the four-byte HTML comment opener gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits the comment opener lets a comment strip survive. |
| MB66 | 60 | A production Go literal equal to the three-byte HTML comment closer gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits the comment closer lets a comment strip survive. |
| MB67 | 60 | A production Go literal equal to `## ` gives one diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check that omits the H2 prefix lets a section walk survive. |
| MB68 | 61 | The same literals in a `_test.go` file and in `internal/markdown` give no diagnostic | planned TestMarkdownBlockOwnerBites in internal/conformance/markdown_block_owner_test.go | A check without the exemptions reds each fixture and the owner itself. |
| MB69 | 62 | `bench test --check markdown-block-owner` is green on the live tree | planned TestMarkdownBlockOwnerHoldsOnTheLiveTree in internal/conformance/markdown_block_owner_test.go | A surviving copy in any module reds the live-tree run. |
| MB101 | 62 | The detected `TestMarkdownBlockOwnerHoldsOnTheLiveTree` assertion has its entry in `classifiedLiveTreeTests` | review-owned: ticket 11 traces the assertion through testChecksLiveTree and runs `bench test --check conformance-meta` | Removing the required entry gives the unregistered-live-tree assertion diagnostic. |
| MB98 | 77 | `skillDescriptionFolded`, `frontmatterHasKey`, and `h2Headings` each call the block reader and hold no block-rule literal | review-owned: the ticket 11 review reads the three function bodies | The check exempts test files, so only a review sees a surviving copy in these helpers. |

Not covered: story 1 — MB1 to MB19 cover the classification facts one at a time, and story 1 is their union.
Not covered: story 63 — Won't handle: the indented code block.
Not covered: story 64 — Won't handle: container fences.
Not covered: story 65 — Won't handle: the backtick info-string rule.
Not covered: story 66 — Won't handle: H3 boundaries.

### Edge inventory

The canonical edge classes at the reader seam each have a row:

- empty input: MB19
- one line without a newline: MB72
- a CRLF document: MB18, MB59, MB80, and MB87
- a nested run and a mixed marker: MB2 and MB3
- an indent, a non-ASCII indent, and a short run: MB4, MB78, and MB6
- a closer with trailing text: MB73
- an unterminated block of each kind: MB7, MB9, and MB13
- two open blocks at the end of the document: MB71
- a block inside another block: MB11 and MB16
- each side of the frontmatter predicate: MB8, MB74, and MB75
- each side of the title predicate: MB14, MB76, and MB77

The hostile-input classes are a rejoined comment (MB12), a large comment-dense
input (MB79), a quoted heading (MB16), and a quoted record opener (MB48).

Each module's absent-versus-empty pair stays with the module, because the reader
reads bytes only. This repository and each linked repository receive the same
reader, because each module reads the files of the repository it runs in.

**Won't handle** — an indented code class in the block reader — `prose.walkParagraphs` keeps its existing indented-code rule in `internal/prose/parse.go`.

**Won't handle** — a fence inside a list item or a block quote — the reader classifies lines, not containers, and each module's grammar stays line-based.

**Won't handle** — the rule that a backtick info string holds no backtick — `internal/prose` omits it today, and adding it moves the live-tree grade.

**Won't handle** — an H3 or deeper section boundary — three modules keep their own rule over unfenced lines: `spec.FenceTokens`, `retros.Recommendations`, and `mapEndRe` at `internal/coverage/coverage.go:52`.

**Won't handle** — a managed-block marker inside an enclosing HTML comment — the markers are comments, and `RewriteAgentsBlock` reads the raw line, as it does today.

**Won't handle** — a rune-literal fence test — `markdown-block-owner` reads string literals only. The legitimate rune users stay: the code-span fold at `internal/prose/parse.go:63-71` and `backtickRun` at `internal/prose/parse.go:97`, `internal/coverage/citations.go:294`, and `internal/spec/fences.go:67`. Ticket 1 deletes the rune test at `internal/prose/prepare.go:104`, so no fence test stays in rune form.

**Won't handle** — the block markers in `findBlock` at `internal/skillsindex/skillsindex.go:466` — the scan fails closed, because a second marker of either kind leaves the file unchanged.

**Won't handle** — `internal/outline/outline.go:100` — lines 15 to 17 document a fence-blind line-regex indexer, and a fenced heading indexes only as a benign candidate.

**Won't handle** — the CHANGELOG heading regex at `internal/releasepreflight/identity.go:74` — it reads release headings of one generated file, not authored Markdown blocks.

**Won't handle** — the skill example reader at `internal/tickets/example_test.go:51` — it reads a skill example in a test, not authored grammar.

## Ownership fences

- `internal/markdown/`
- `internal/prose/prepare.go`
- `internal/prose/prose_test.go`
- `internal/spec/spec.go`
- `internal/spec/fences.go`
- `internal/spec/fences_test.go`
- `internal/coverage/coverage.go`
- `internal/coverage/coverage_test.go`
- `internal/roadmap/tree.go`
- `internal/roadmap/sequence_projection.go`
- `internal/roadmap/commitment_test.go`
- `internal/maps/fields.go`
- `internal/maps/tickets_test.go`
- `internal/tickets/tickets.go`
- `internal/tickets/fence_test.go`
- `internal/handoffdoc/fence.go`
- `internal/handoffdoc/document.go`
- `internal/handoffdoc/document_test.go`
- `internal/learnings/learnings.go`
- `internal/handoff/state_file.go`
- `internal/handoff/state_file_test.go`
- `internal/learnings/entry.go`
- `internal/learnings/learnings_test.go`
- `internal/retros/retros.go`
- `internal/retros/retros_test.go`
- `internal/retros/recommendations.go`
- `internal/retros/recommendations_test.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/write.go`
- `internal/reviewrecord/record_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/anchors/locate.go`
- `internal/anchors/match.go`
- `internal/anchors/locate_test.go`
- `internal/skillsindex/skillsindex.go`
- `internal/skillsindex/frontmatter_test.go`
- `internal/adopt/marker.go`
- `internal/adopt/link_plan_test.go`
- `internal/conformance/markdown_block_owner_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/tier_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/conformance/skill_description_budgets_test.go`
- `internal/conformance/skills_index_checks_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/preflight/preflighttest/fixture.go`
- `projects/benchkit.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `tests/canary/package-core-guard/bounds-classify-limit-restated`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `reviews/markdown-block-reader.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-add-block-reader.md` | none | MB-C1 |
| `2-read-spec-and-coverage-blocks.md` | `1-add-block-reader.md` | MB-C2 |
| `3-read-roadmap-blocks.md` | `1-add-block-reader.md` | MB-C2 |
| `4-read-field-scan-blocks.md` | `1-add-block-reader.md` | MB-C2 |
| `5-read-handoff-blocks.md` | `1-add-block-reader.md` | MB-C3 |
| `6-read-journal-blocks.md` | `1-add-block-reader.md` | MB-C3 |
| `7-read-review-record-blocks.md` | `1-add-block-reader.md` | MB-C4 |
| `8-read-anchor-blocks.md` | `1-add-block-reader.md` | MB-C4 |
| `9-read-skills-frontmatter.md` | `1-add-block-reader.md` | MB-C3 |
| `10-read-agents-marker-blocks.md` | `1-add-block-reader.md` | MB-C5 |
| `11-forbid-block-rule-copies.md` | tickets 2 to 10 | MB-C5 |

## Out of scope

- A full CommonMark block grammar with lists, block quotes, indented code, and setext headings. Estimate: 20 edits, 4 gate runs.
- One shared H3 boundary rule for the fence section and the retrospective units. Estimate: 8 edits, 2 gate runs.
- A Markdown reader for a TypeScript, JavaScript, or Python project. No such reader exists today, and the parked language-support idea owns that work. Estimate: 0 edits, 0 gate runs in this spec, because no such reader exists to change.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| "Fourteen fence detectors live in nine modules." | MB24 to MB61, one module migration each, and MB62 to MB69 |
| "Two accept a tilde fence and twelve do not." | MB1, MB3, MB24, MB31 to MB35, MB39, MB52 to MB54, MB57, MB58, MB61 |
| "Three line scanners (coverage, learnings, and retros) have no fence rule." | MB29, MB30, MB41 to MB47; the tree holds two more, covered by MB26 to MB28 |
| "`internal/spec` trims a line at one site and does not trim it at a second site." | MB25 |
| "One deep reader owns the block rules: fences, frontmatter, comments, and H2 sections." | MB1 to MB19, MB71 to MB79, MB88 to MB90, MB92 |
| "Each grammar module reads unfenced lines from it." | MB20 to MB61, MB70, MB80 to MB87, MB91 to MB97 |
| "`internal/prose` stays the prose-rule owner, and the reader composes with it." | MB20 to MB23 |
| "The fence edges then have one test." | MB1 to MB7, MB73, and MB78 in `TestReadFenceEdges` |

### Detector inventory

The inventory names the fence detectors at `37ac80b8`. Each row is one fence
test that toggles or classifies a line. `sequenceBounds` composes the
`sectionBounds` detector, which also serves dependency closure.

| module | site | tilde |
| --- | --- | --- |
| `internal/prose` | `prepare.go:76` `stripFences` | yes |
| `internal/handoffdoc` | `fence.go:15` `isFence` | yes |
| `internal/spec` | `spec.go:45` `LiveSpecSlugs` | no |
| `internal/spec` | `spec.go:69` `metadata` | no |
| `internal/maps` | `fields.go:58` `FieldScan.Scan` | no |
| `internal/adopt` | `marker.go:42` `scanMarkers` | no |
| `internal/adopt` | `marker.go:85` `RewriteAgentsBlock` | no |
| `internal/adopt` | `marker.go:137` `StripAgentsBlock` | no |
| `internal/roadmap` | `tree.go:89` `ParseDocument` | no |
| `internal/roadmap` | `tree.go:268` `parseSequence` | no |
| `internal/roadmap` | `sequence_projection.go:17` `sectionBounds` | no |
| `internal/roadmap` | `tree.go:339` `rowNextDiagnostics` | no |
| `internal/tickets` | `tickets.go:96` `ParseTicket` | no |
| `internal/anchors` | `locate.go:100` `scopeRunesMapped` | no |
| `internal/reviewrecord` | `parse.go:249` `locate` | no |

The five production scanners with no fence rule are `coverage.parse` at
`coverage.go:187`, `learnings.Parse` at `learnings.go:49`, `retros.Parse` at
`retros.go:63`, `spec.FenceTokens` at `fences.go:28`, and `deriveImplemented` at
`spec.go:156`. The helpers `hasAnyHeading` at `learnings.go:264` and
`Recommendations` at `recommendations.go:16` belong to two of them. The second
frontmatter reader is `skillsindex.FrontmatterField` at `skillsindex.go:170`. The
second comment strip is `stripCommentsMapped` at `locate.go:76`. Ticket 11 moves
three conformance readers, and a Won't handle line names `findBlock`.

### Rule delta by site

The table records the module rules and the current roadmap section detector at `37ac80b8`. "None" in the last column means the
site keeps its verdict. Each other entry names a row, a flagged addition (FA1 to
FA9 below), or a Won't handle line.

| site | indent today | tilde today | CR today | closer today | trim predicate today | comments today | changes under the reader |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `prose/prepare.go:76` `stripFences` | spaces and tabs | yes | ignored after the run | same character, run at least as long, trailing text kept | `TrimLeft` of space and tab | stripped first, one pass per line | none for fences: MB1 to MB7, MB73 |
| `prose/prepare.go:26` `stripFrontmatter` | any `TrimSpace` space | not applicable | trimmed | next `---` line | `TrimSpace` | not applicable | a padded `---` opens or closes nothing: MB74, MB75, MB90, FA1 |
| `prose/prepare.go:43` `stripComments` | not applicable | not applicable | kept | next closer | none | one pass per line | a rejoined comment goes: MB22 |
| `handoffdoc/fence.go:15` `isFence` | spaces and tabs | yes | callers replace CRLF, `UnfencedLines` keeps a CR | any marker line toggles | `TrimLeft` and `TrimRight` of space and tab | none | close by character and length: MB38 to MB40; FA2 |
| `handoff/state_file.go:69` `readStateFile` | through `OpenFence` | through `OpenFence` | through `OpenFence` | through `OpenFence` | `TrimSpace` for the repair text | none | an open comment refuses: MB82 |
| `spec/spec.go:45` `LiveSpecSlugs` | any `TrimSpace` space | no | trimmed | any marker line toggles | `TrimSpace` | none | tilde: MB24; U+00A0 indent: MB78, FA3; FA2 |
| `spec/spec.go:69` `metadata` | column zero | no | `TrimSuffix` of CR | any marker line toggles | none | none | tilde: MB24; indent: MB25, FA4; FA2 |
| `spec/spec.go:156` `deriveImplemented` | no fence rule | no fence rule | `stagedRe` refuses a CR | no fence rule | none | none | fenced status: MB26; CRLF flips: MB80 |
| `spec/fences.go:28` `FenceTokens` | no fence rule; heading by `TrimSpace` | no fence rule | trimmed | no fence rule | `TrimSpace` | none | fenced example: MB27, MB28; column-zero heading: FA5 |
| `coverage/coverage.go:187` `parse` | no fence rule | no fence rule | `TrimSuffix` of CR | no fence rule | none; exact heading | raw marker | fenced lines: MB29, MB30; marker stays raw: MB83; FA2 |
| `coverage/coverage.go:52` `mapEndRe` | column zero | not applicable | not applicable | not applicable | none | none | none: Won't handle H3 |
| `maps/fields.go:58` `FieldScan.Scan` | column zero | no | CRLF replaced | any marker line toggles | none | none | tilde: MB34; indent: FA4; FA2 |
| `tickets/tickets.go:96` `ParseTicket` | column zero | no | CRLF replaced | odd count of marker lines | none | none | reader fault: MB35 to MB37; indent: FA4; FA2 |
| `adopt/marker.go:42`, `:85`, `:137` | column zero | no | a CR stays in each line | any marker line toggles | none | markers read raw | tilde: MB57, MB58, MB61; comment fault: MB70, MB86; CRLF rebuild: MB87; indent: FA4; FA2 |
| `roadmap/tree.go:89` `ParseDocument` | column zero | no | `TrimRight` of CR | any marker line toggles; rows match inside a fence | `TrimRight` of space, tab, and CR | none | fenced rows skip: MB31, MB81, FA6; indent: FA4; FA2 |
| `roadmap/tree.go:268` `parseSequence` | column zero | no | `TrimRight` of CR | any marker line toggles | `TrimRight` of space, tab, and CR | none | tilde: MB32; indent: FA4; FA2 |
| `roadmap/sequence_projection.go:17` `sectionBounds` | column zero | no | `TrimRight` of CR | any marker line toggles | `TrimRight` of space, tab, and CR | none | tilde and quoted boundaries: MB32, MB99, MB100; indent: FA4; FA2 |
| `roadmap/tree.go:339` `rowNextDiagnostics` | column zero | no | `TrimRight` of CR | any marker line toggles | `TrimRight` of CR, then of space and tab | none | tilde: MB33; indent: FA4; FA2 |
| `anchors/locate.go:100` `scopeRunesMapped` | any `unicode.IsSpace` space | no | trimmed | any marker line toggles | `unicode.IsSpace` | stripped first, with rejoin | tilde: MB52 to MB54; U+00A0 indent: MB78, FA3; indented heading: MB85 |
| `anchors/locate.go:76` `stripCommentsMapped` | not applicable | not applicable | kept | next closer, then a rescan | none | rejoin, one copy per comment | linear work: MB79; same result: MB55 |
| `reviewrecord/parse.go:249` `locate` | exact line | no | a CR stops the match | exact closer line | none | none | quoted opener: MB48; indent, tilde, and CR open: FA7; FA2 |
| `learnings/learnings.go:49` `Parse` | no fence rule | no fence rule | `TrimSuffix` of CR | no fence rule | `TrimSpace` for blank lines | raw marker | fenced heading: MB41; faults: MB42 to MB44; region raw: MB84; FA2 |
| `retros/retros.go:63` `Parse` | no fence rule | no fence rule | `TrimSuffix` of CR | no fence rule | none; exact heading | none | fenced heading: MB45; FA2 |
| `retros/recommendations.go:16` `Recommendations` | no fence rule | no fence rule | `TrimSuffix` of CR | no fence rule | `TrimSpace` for the title | none | fenced heading and item: MB46, MB47; FA2 |
| `skillsindex/skillsindex.go:170` `FrontmatterField` | exact `---` | not applicable | a CR stops the match | exact `---` | none | not applicable | CRLF opens: MB59 |
| `skillsindex/skillsindex.go:466` `findBlock` | exact marker line | no fence rule | a CR stops the match | not applicable | none | markers read raw | none: Won't handle |
| `conformance/skill_description_budgets_test.go:136` | exact `---` | not applicable | a CR stops the match | exact `---` | none | not applicable | CRLF opens: FA8 |
| `conformance/skills_index_checks_test.go:249` | exact `---` | not applicable | a CR stops the match | exact `---` | none | not applicable | CRLF opens: FA8 |
| `conformance/docs_workflow_checks_test.go:778` `h2Headings` | no fence rule | no fence rule | kept | no fence rule | none | none | a fenced or commented heading drops: FA8 |

Each changes cell names its coverage rows, flagged addition, or exclusion.
FA9 is the check itself, rather than a detector site.

### Reader sweep

`bench consumers` at `6d8191dc` lists each production caller of a changed function.

- `prose.stripFences`: `prose.prepare`.
- `handoffdoc.isFence`: `splitSections`, `parseSection`, `OpenFence`, and `UnfencedLines`. `handoffdoc.OpenFence`: `handoff.readStateFile` and `handoffdoc.Parse`. `handoffdoc.UnfencedLines`: `handoff.readStateFile` and `handoff.scanState`.
- `spec.LiveSpecSlugs`: `roadmap.ParseDocument` and `status.roadmapReconcileCounts`. `spec.metadata`: `Facts`, `AwaitsRetirement`, and `retireCommand`. `spec.FenceTokens`: `coverage.parse` and `preflight.gather`. `spec.Implemented`: `gate.completionTree`, `landing.LandReviewed`, `worktree.landingSource`, and `worktree.specTransitionFacts`.
- `maps.FieldScan`: `decisionMapSchema.fieldScan`, `decisionMapSchema.ticketFileScan`, and `tickets.fieldScan`.
- `adopt.scanMarkers`: `evalAgentsRow`, `validateAgentsContent`, `RewriteAgentsBlock`, `StripAgentsBlock`, and `stagedRepairAgents`. `RewriteAgentsBlock`: `stagedAgents` and `stagedRepairAgents`. `StripAgentsBlock`: `stripAgentsForUnlink`. The link transaction calls `stagedRepairAgents` at `internal/adopt/link_transaction.go:236`, so ticket 10 reads that composed repair path.
- `roadmap.ParseDocument`: `BuildContext`, `validateOccurrenceOwner`, `RoadmapCommand`, `RoadmapText`, `ValidateRoadmapTree`, and `roadmapflow.openMass`. `parseSequence`: `RecommendedSequence`, `ParseDocument`, and `SequenceText`. `rowNextDiagnostics`: `ParseDocument`.
- `roadmap.sequenceBounds`: `parseSequence` and `ProjectSequence`. `sectionBounds`: `sequenceBounds` and `closeDependencies`. `closeDependencies`: `Close`. `ProjectSequence`: `Close`, `Store.protectedCandidate`, and `Store.stage` in `internal/commitment/repository`. The current `bench consumers` sweep at `37ac80b8` resolves these edges.
- `tickets.ParseTicket`: `preflight.gradeTickets` and `reviewrecord.ReadPlan`.
- `anchors.scopeRunesMapped`: `sectionRunesMapped` and `stepRunesMapped`. `MarkdownH2Sections`: `MarkdownH2Section` and `resolveSection`. `StripHTMLComments`: `read`. `MarkdownNumberedSteps`: `resolveStep`.
- `reviewrecord.locate`: `fenced` and `Render`.
- `skillsindex.FrontmatterField`: `Entries`.
- `learnings.Parse`: `Entries`, `Command`, `roadmap.BuildContext`, and `roadmap.learningCount`.
- `retros.Parse`: `roadmap.RetroCommand`. `retros.Recommendations`: `Parse`, `ValidateImprovementMarkers`, and `roadmap.BuildContext`.

The review named five behavior-visible callers one level further out. `bench
consumers` at `ba78dabc` lists each one:

- `reviewrecord.fenced` serves `ReadPlan` at `plan.go:47` and `parseRecord` at `files.go:99`. `ReadPlan` serves `preflight.withCompletionPlan`, `reviewrecord.checkSource`, `reviewrecord.atSource`, and the two `recordtest` fixture loaders. `parseRecord` serves `Read`, `ReadTree`, and `write`.
- `coverage.parse` serves `ParseSpec` and `Command`. `ParseSpec` serves `preflight.gather`.
- `spec.AwaitsRetirement` serves `retireCommand`, `implementedAtHEAD`, `status.retirementCount`, and `status.roadmapReconcileCounts`.
- `spec.Facts` serves `handoff.liveSpecs`, `preflight.specStatus`, `roadmap.BuildContext`, and `status.stagedSpecCount`.

Each caller keeps its call, because each exported signature stays the same. A
caller sees only the reclassified lines that the rows name.

The test callers of a changed export keep their calls. `prose.Paragraphs` has
three test callers in `internal/conformance/charge_evidence_guidance_test.go`, at
lines 121, 147, and 392. The test callers of the anchor exports are these:

- `anchors.StripHTMLComments`: `internal/conformance/docs_workflow_helpers_test.go:322`, which binds it to the alias `stripHTMLComments`, and line 501 of that file; `internal/anchors/match_test.go:209`; and `internal/anchors/locate_test.go:112` and `:113`.
- `anchors.MarkdownH2Section`: `internal/conformance/docs_workflow_helpers_test.go:181` and `internal/anchors/match_test.go:100`.
- `anchors.MarkdownNumberedSteps`: `internal/anchors/match_test.go:114`.

The conformance callers read live guidance text. Ticket 1 records equal
`bench gate-prose` output over each tracked Markdown file at the base and at the
tip, and ticket 8 keeps `docs-currency-workflow` green. The anchor package tests
keep their expectations without an edit, as MB55 states.

### Fence closures

Build preflight adds the fence paths after `projects/benchkit.md` as closures.
The five command-binding files follow the bound packages of tickets 2, 3, 4, 6,
8, and 11. The anchor registry files and the canary fixtures follow
`projects/benchkit.md` in ticket 11. The fixture `bounds-classify-limit-restated`
follows `internal/learnings/learnings.go` in ticket 6. No ticket expects to edit
a closure path.

`CONTEXT.md` already defines the block reader, the fenced block, and the unfenced
line. The glossary decision is satisfied, and this spec needs no glossary write.

### Fence overlap with in-flight work

- The `cli-desktop-consistency` prerequisite is satisfied at `37ac80b8`. `stagedRepairAgents` and its link-transaction caller are present. Ticket 10 reads that current path before its marker migration.
- `specs/ft290-test-projection` ticket 4 writes `internal/prose/walk.go` and `internal/prose/walk_test.go`. Ticket 1 writes `internal/prose/prepare.go` and `internal/prose/prose_test.go` only, so the two fences do not overlap.
- The open drain `d-007c40a25f47` owns `ROADMAP.md` and `roadmap/`. This spec writes neither, and the FT358 row retires at this spec's landing.

### Flagged additions

- MB42, MB43, and MB44 add malformed records to `internal/learnings`. The source names no refusal, but without them a stray fence makes later entries vanish silently.
- MB59 changes `FrontmatterField` for a CRLF skill file. The change follows from the one frontmatter rule.
- MB22 changes the prose grade of a rejoined comment. The change follows from one comment rule, and DG43 already chose the rejoin rule for the anchors.
- MB70 and MB86 add a comment-fault refusal to `internal/adopt`, under the Bench-text condition. Closed decision 1 requires it.
- MB82 adds a comment-fault refusal to `readStateFile`. Closed decision 12 requires it.
- MB80 lets `Implemented` flip a CRLF spec. Closed decision 11 requires it.
- FA1: `prose.prepare` stops opening or closing frontmatter on a padded `---` line. Closed decision 8 requires it, and MB23 shows that the live tree does not move.
- FA2: a fence marker or a `## ` line inside an HTML comment no longer opens a block or a section. This is new in each module that read comments as text:  `handoffdoc`, `spec`, `coverage`, `maps`, `tickets`, `adopt`, `roadmap`, `reviewrecord`, `learnings`, and `retros`.
- FA3: a U+00A0 indent stops opening a fence in `spec.LiveSpecSlugs` and `anchors.scopeRunesMapped`, because the reader's indent is space and tab only. MB78 pins the rule.
- FA4: an indented fence opener now opens a block in each column-zero site. The roadmap sites are `tree.go:89`, `:268`, and `:339`, plus `sequence_projection.go:17`. The other sites are `maps/fields.go:58`, `tickets/tickets.go:96`, `adopt/marker.go:42`, `:85`, and `:137`, and `spec/spec.go:69`. MB25 covers the spec site. The `reviewrecord` site is FA7.
- FA5: `FenceTokens` opens its section only on a column-zero `## Ownership fences` line. Closed decision 4 requires it.
- FA6: `ParseDocument` stops parsing a roadmap row inside any fence, a backtick fence included. MB81 covers it.
- FA7: `reviewrecord.locate` now opens a record on an indented opener, a tilde opener, or a CRLF opener. `ReadPlan` and `parseRecord` read records that Bench writes itself, so the output of each Bench writer does not move.
- FA8: the three conformance readers of ticket 11 read CRLF frontmatter, and `h2Headings` drops a fenced or commented heading. Closed decision 9 requires the move.
- FA9: `markdown-block-owner` is a new conformance check. The source asks for one reader, and the check is the red-capable copy-survival row.

### Closed decisions

The reviewer defers these decisions to the Fable consultation of 2026-10-03.

1. Fault posture, amended (2026-10-03, Fable consultation). Each existing refusal keeps the fence fault, and the learnings records stay. `internal/adopt` also maps the comment fault to its message at `marker.go:68`, under the same Bench-text condition, because an unterminated comment hides later fence markers. MB70 and MB86 cover it.
2. Ticket 10 sequencing, accepted (2026-10-03, Fable consultation). Ticket 10 requires the landed `cli-desktop-consistency` marker path. That prerequisite is satisfied at `37ac80b8`.
3. Guard scope, accepted (2026-10-03, Fable consultation). The check covers the `## ` literal, and each renderer uses the exported constants.
4. H2 opener, accepted (2026-10-03, Fable consultation). A heading opens at column zero only, so the anchors and `FenceTokens` stop accepting an indented opener.
5. Comment rule, accepted (2026-10-03, Fable consultation). The block reader applies the DG43 rejoin rule for each module.
6. Guard literal rule, amended (2026-10-03, Fable consultation). The check refuses a literal that contains a run of three backticks or three tildes, not only one that starts with the run. The one hit outside the migrated modules is `internal/preflight/preflighttest/fixture.go`, which ticket 11 routes through the fence constant.
7. Ticket 11 closures, accepted (2026-10-03, Fable consultation). Ticket 11 keeps each closure path in `Writes:`, because the ticket grammar has no closure field and build preflight reads closures from `Writes:`.
8. Frontmatter predicate, decided (2026-10-03, Fable consultation). The opener and the closer are exactly `---` after the carriage-return strip, in every module. This fails closed for a hostile skill file, and no loader accepts an indented opener. MB59, MB74, and MB75 cover it.
9. Conformance readers, decided (2026-10-03, Fable consultation). Ticket 11 moves `skillDescriptionFolded`, `frontmatterHasKey`, and `h2Headings` onto the block reader. The skill example reader at `internal/tickets/example_test.go:51` takes a Won't handle line.
10. Linear rejoin, decided (2026-10-03, Fable consultation). The reader removes comments with a stack over its output, in one pass. MB79 bounds the work, and MB12, MB55, MB88, and MB89 pin the DG43 result.
11. Rebuild from source bytes, decided (2026-10-03, Fable consultation). The adopt rebuilders and `deriveImplemented` write from the source bytes by the reader's line offsets. MB80 and MB87 cover a CRLF round trip.
12. Handoff comment fault, decided (2026-10-03, Fable consultation). `readStateFile` refuses an unterminated comment with a new fault that names the opener line. MB82 covers it.
13. Glossary, decided (2026-10-03, Fable consultation). The glossary stays outside this spec's writes. `CONTEXT.md` already holds the three terms, so the separate glossary prerequisite is satisfied.
14. Rune literals, decided (2026-10-03, Fable consultation). The check stops a new string-literal detector only, and a Won't handle line names the legitimate rune users. Ticket 1 deletes the rune test at `internal/prose/prepare.go:104`.
15. Learnings region, decided (2026-10-03, Fable consultation). The region below the entries marker reads raw lines, because its rule is that nothing belongs there. DL31 stays green without an edit, and MB84 covers it.
16. Fragments and the frontmatter fault, decided (2026-10-03, Fable consultation). The reader exports `Read` for a document and `ReadFragment` for a section or State body, over one pipeline. The fragment callers are the two handoff State readers and the two anchor scope walks. `internal/adopt`, `internal/tickets`, and `internal/reviewrecord` ignore the frontmatter fault and read each line as body, as they do today. MB92 to MB97 cover it.

### Pre-review proof checklist

- `Cited symbols`: the refreshed inventory resolves at `37ac80b8`, including `sequenceBounds` and `sectionBounds`. `TestWaitDeadlineLiteralsBites`, `checkGitPlumbingOwner`, and each test function in the map's seam cells resolve in the named files.
- `Import edges`: each module gains the edge to `internal/markdown`. That package imports only the standard library, so no edge closes a cycle. `go list` runs at ticket 1, because the package does not exist before it.
- `Source-row clauses and occurrences`: the source trace above quotes each clause of `roadmap/FT358.md`. Card 02 of the 2026-09-29 quality survey is not retained in the tree or under `~/.bench`. The detector inventory names the current sites instead of repeating its counts.
- `Promised field labels`: the malformed reasons `unterminated fenced block`, `unterminated HTML comment`, and `unterminated frontmatter block` in MB42 to MB44, the handoff comment fault of MB82, and the check name `markdown-block-owner`.
- `Changed-function callers`: the reader sweep lists each production caller and the named test callers.
- `Copy survival`: MB62 to MB69 fail when any string-literal block rule survives outside `internal/markdown`.
- `Rendered-shape readers`: no existing message text changes. The `unterminated fence` diagnostic in `internal/tickets`, the two AGENTS.md conflict messages, the handoff refusals, and the review record refusals keep their bytes. The handoff comment fault of MB82 is a new message.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"MB-C1","tickets":["1-add-block-reader.md"],"verification":[{"id":"markdown","command":"bench test --package ./internal/markdown"},{"id":"prose","command":"bench test --package ./internal/prose"},{"id":"prose-mechanics","command":"bench test --check prose-mechanics"}]},{"id":"MB-C2","tickets":["2-read-spec-and-coverage-blocks.md","3-read-roadmap-blocks.md","4-read-field-scan-blocks.md"],"verification":[{"id":"spec","command":"bench test --package ./internal/spec"},{"id":"coverage","command":"bench test --package ./internal/coverage"},{"id":"roadmap","command":"bench test --package ./internal/roadmap"},{"id":"maps","command":"bench test --package ./internal/maps"},{"id":"tickets","command":"bench test --package ./internal/tickets"}]},{"id":"MB-C3","tickets":["5-read-handoff-blocks.md","6-read-journal-blocks.md","9-read-skills-frontmatter.md"],"verification":[{"id":"handoffdoc","command":"bench test --package ./internal/handoffdoc"},{"id":"handoff","command":"bench test --package ./internal/handoff"},{"id":"learnings","command":"bench test --package ./internal/learnings"},{"id":"retros","command":"bench test --package ./internal/retros"},{"id":"skillsindex","command":"bench test --package ./internal/skillsindex"}]},{"id":"MB-C4","tickets":["7-read-review-record-blocks.md","8-read-anchor-blocks.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"docs-currency-workflow","command":"bench test --check docs-currency-workflow"}]},{"id":"MB-C5","tickets":["10-read-agents-marker-blocks.md","11-forbid-block-rule-copies.md"],"verification":[{"id":"adopt","command":"bench test --package ./internal/adopt"},{"id":"markdown-block-owner","command":"bench test --check markdown-block-owner"},{"id":"skill-description-budgets","command":"bench test --check skill-description-budgets"},{"id":"skills-index-command-adapters","command":"bench test --check skills-index-command-adapters"},{"id":"conformance-meta","command":"bench test --check conformance-meta"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/markdown-block-reader/spec.md"},{"id":"markdown","command":"bench test --package ./internal/markdown"},{"id":"prose-mechanics","command":"bench test --check prose-mechanics"},{"id":"docs-currency-workflow","command":"bench test --check docs-currency-workflow"},{"id":"markdown-block-owner","command":"bench test --check markdown-block-owner"},{"id":"conformance-meta","command":"bench test --check conformance-meta"}]}
```
