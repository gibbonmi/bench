# asd-ste100-progressive-disclosure

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-22 — the closed decisions are listed under `## Further notes`; the reviewer pre-approved the author's judgment calls and the spec commit in the prepared worktree

Verification log: 3 iteration(s) to accept — `gpt-5.6-sol` at xhigh through `codex exec` and `opus` at medium through the Claude agent surface; iteration 1 caught un-anchored moved passages, an orchestrator-graded exclusion rule, four unrowed embedded templates, rows at a seam that could not red them, and oversized tickets; iteration 2 caught the core-side absence proof, the parser predicates, the generated `CLAUDE.md`, and overbroad fences; iteration 3's remaining findings — `Forbid` rows for absence, the label-run paragraph rule, the list-marker rule, the shell scaffold ownership — were folded into this acceptance under the reviewer's three-iteration cap

## Problem

The kit's written guidance carries two costs. First, `.bench/BENCH.md` loads into every session, but it mixes the rules a session needs every time with mechanics a session needs rarely. Second, the prose across the kit — skills, commands, platform files, docs, roadmap, decisions, and code comments — does not yet follow ASD-STE100 in form. Long sentences and dense paragraphs cost every reader a decode. The 2026-08-22 two-arm comparison showed the STE form is shorter, more exact, and easier to slice. Nothing in the gate enforces the form, so drift returns one edit at a time.

## Solution

Split `.bench/BENCH.md` into an always-loaded core and reference material. The core keeps six rule families. They are the roles, the invariants and predicates, the communication rules, the workflow routing, the fix and capture triggers, and the CLI execution contract. Mechanics move to `.bench/BENCH-reference.md`, the command files, or the skills. Every enforcement anchor moves with its rule, and each moved rule keeps planted-reason proof. The always-loaded CLI rule shows one valid call and one invalid call.

Add one fail-closed conformance check, `prose-mechanics`, for the two ASD-STE100 rules a program can grade reliably. Rule one: no sentence is longer than 25 words. Rule two: no paragraph has more than six sentences. The check reads every authored Markdown file under the graded root. A reviewer-owned exclusion file, `.bench/prose-exclusions`, names the paths the check does not grade. The semantic STE rules stay with top-tier review.

Rewrite every authored Markdown file and every explanatory Go and shell comment in ASD-STE100. Cheap-tier delegates author each batch. The top-tier orchestrator validates semantic preservation and cross-file consistency. A repair returns to a cheap-tier delegate. The skill batches run in the reviewer's order. FT100 and FT179 keep the work this spec does not do.

Canonical vocabulary: the **prose mechanics check** is the `prose-mechanics` conformance check. Avoid "STE lint", "prose lint", and "grammar check". A **prose exclusion row** is one line of `.bench/prose-exclusions`: a path and a one-clause reason. Avoid "allowlist" and "skip list". The **always-loaded core** is `.bench/BENCH.md`. Avoid "the guide" without a qualifier.

## User stories

### Group A — the always-loaded core discloses progressively

Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex. The disposition table is exact, the registry and fixture seams are known, and the gate covers every anchor move.

1. As a session that loads `.bench/BENCH.md`, I want only the six always-loaded rule families in it, so that every loaded line steers generation.
2. As a session that needs a mechanic, I want a reference, command, or skill to hold it, so that I pay only on demand.
3. As a maintainer, I want each moved rule to keep its anchor at the new location, so that a later deletion still reds the gate.
4. As a maintainer, I want each moved anchor to keep or gain planted-reason proof, so that the check is observed red before trust.
5. As a cold session, I want the CLI contract to show `bench gate` as valid and `bench gate 2>&1 | tail -20` as invalid, so that I imitate the example.
6. As a maintainer, I want no H2 heading shared between the core and the reference, so that one fact keeps one source.
7. As a maintainer, I want every shared-rule marker to stay in the core, so that the single-source check stays green.
8. As a linked-repo reader, I want the core to keep its pointer to the reference, so that the moved material stays reachable.

### Group B — the prose mechanics check

Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex. The oracle's correctness matters more than speed, and the parser follows the coverage and prose-budget precedents.

9. As a gate user, I want a sentence over 25 words to red the gate with file and line, so that drift cannot land.
10. As a gate user, I want a paragraph over six sentences to red the gate with file and line, so that dense blocks cannot land.
11. As an author, I want no code span, fence, table, heading, frontmatter, comment, URL, terminator-free label line, or indent graded, so that only prose counts.
12. As an author, I want a list item graded as its own paragraph, so that a long bullet is held to the paragraph rule.
13. As an author, I want a token-internal period, five fixed abbreviations, or an ellipsis to add at most one boundary, so that none is invented.
14. As a reviewer, I want `.bench/prose-exclusions` to hold every excluded path with a one-clause reason, so that an exclusion is a named decision.
15. As a reviewer, I want a stale, duplicate, or glob row to red the gate, so that the exclusion file cannot rot.
16. As a gate user, I want a link, special, non-UTF-8, or oversized `*.md` file to red the gate, so that no unreadable byte is graded.
17. As a gate user, I want an unterminated fence, comment, or frontmatter block to red the gate, so that a truncating read cannot hide prose.
18. As a gate user, I want the walk to skip `.git/`, `node_modules/`, `dist/`, and any `testdata/` directory, so that vendored and fixture trees never count.
19. As a maintainer, I want one canary fixture per planted failure mode, so that the check is proven to bite through its registered owner.
20. As a maintainer, I want the check to refuse a hostile skill file and name it, so that it composes with the hostile-reader row.
21. As a build delegate, I want one focused live-tree test, so that a batch goes red and then green without a full gate.
22. As a maintainer, I want `ste-prose.md` to keep the two thresholds under registry anchors, so that rule text and enforcement cannot drift apart.

### Group C — authored Markdown reads as ASD-STE100

Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex, for each batch. The orchestrator validates at `fable` / high under Claude Code, `gpt-5.6-sol` / high under Codex. The reviewer set this spec-specific override on 2026-08-22; the general leverage rule in `craft-line` is unchanged.

23. As a reader of a skill, I want its body in ASD-STE100, so that every tier reads the rule the same way.
24. As a maintainer, I want the skill batches migrated in the reviewer's order, so that each batch's delegates load already-migrated guidance.
25. As a reader of a command file, I want its prose in ASD-STE100 with every needle intact, so that the phase keeps its enforced clauses.
26. As a reader of a thin `bench-*` adapter skill, I want it migrated after its command, so that the adapter follows its canonical document.
27. As a reader of a root document, I want AGENTS, CLAUDE, CONTEXT, README, and the docs in ASD-STE100, so that the cold read is cheap.
28. As a reader of the roadmap, I want every row body in ASD-STE100 with its grammar lines intact, so that the drain still parses it.
29. As a reader of a decision map, I want its prose in ASD-STE100 with its ticket grammar intact, so that the map validator stays green.
30. As a maintainer, I want every rewrite to keep every code token, path, command, identifier, and example, so that meaning survives the form change.
31. As a maintainer, I want frontmatter, fenced examples, tables, and check-read markers left byte-identical in a rewrite, so that machine-read surfaces do not move.
32. As a maintainer, I want each batch to remove only its own exclusion rows and never add one, so that the exclusion file only shrinks.

### Group D — explanatory comments read as ASD-STE100

Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex, for each batch. The orchestrator validates at `fable` / high under Claude Code, `gpt-5.6-sol` / high under Codex. Same reviewer override as Group C.

33. As a Go reader, I want every explanatory comment in ASD-STE100 under the `craft-comments` register, so that the next reader learns the constraint fast.
34. As a reader of shell code, I want every explanatory comment in ASD-STE100, so that hooks and scripts read the same way as Go.
35. As a maintainer, I want a comment that restates its line deleted, so that no-value comments do not survive the pass.
36. As a maintainer, I want machine-readable directives and guard manifest headers left byte-identical, so that `bench guards` and the toolchain keep reading them.
37. As a maintainer, I want a comment batch to change comment lines only, so that behavior cannot move under a prose pass.
38. As a maintainer, I want no comment lint added, so that the register stays a review judgment.

### Group E — records and reconciliation

Line: `sonnet` / low under Claude Code, `gpt-5.6-luna` / low under Codex. The texts are decided in this spec, and the delegate transcribes them.

39. As a reader of FT100, I want its body to name what this spec delivered and what remains, so that the cut is not duplicated.
40. As a reader of FT179, I want its body to name what this spec delivered and what remains, so that the remainder is not duplicated.
41. As a linked-repo maintainer, I want one CHANGELOG entry, so that the new check and the core split are visible.
42. As a cold session, I want `CONTEXT.md` to define the new terms, so that the vocabulary does not drift.
43. As a cold session, I want `capture/session-handoff.md` rewritten at every phase close, so that resumption never depends on conversation history.

### Group F — Go-embedded Markdown templates

Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex. The orchestrator validates at `fable` / high under Claude Code, `gpt-5.6-sol` / high under Codex. Same reviewer override as Group C.

44. As a linked-repo maintainer, I want the Markdown the kit writes into a tree to meet both bounds, so that scaffolds read the same way.

### Group G — one term for the split

Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex. The rows are registry data with fixtures, gate-covered.

45. As a reader of the core or reference, I want the split named "progressive disclosure", never "progressive loading", so that one term holds one meaning.

## Implementation decisions

### The core and the reference

`.bench/BENCH.md` keeps these sections: the title and pointer, `Roles`, the Bench CLI contract, `The four invariants`, `How to talk to me`, `Workflow`, and `Capture`. The `CLI Inventory` heading becomes `The Bench CLI contract`. No H2 in the core appears in the reference.

The disposition of the current core text:

| core section | disposition | target |
|---|---|---|
| title and pointer | keep | core |
| Roles | keep | core |
| How the pieces fit (four bullets) | move | reference, new H2 `How the pieces fit` |
| CLI Inventory category bullets | move | reference, `Command Notes` |
| CLI execution contract paragraph | keep and rewrite | core, with the valid and invalid pair |
| The four invariants and three predicates | keep; split long items into paragraphs of six sentences or fewer | core |
| How to talk to me | keep | core |
| Workflow numbered list | keep | core |
| Right-size paragraph and table | keep | core |
| retained-integration-source paragraph | move; merge into the reference paragraph that already states it | reference, `Command Notes` |
| Fix, don't park; batch approval | keep | core |
| Capture what you learn | keep | core |
| Capture: the trigger, the sink, the graduation route | keep | core |
| Capture: the no-PATH fallback | move | reference, `Files`, beside the `capture/IDEAS.md` role |
| Retros are capture | move | reference, `Files`, as one sentence naming both owners |

The core keeps one sentence that points a reviewed spec-backed build at `bench worktree land` and at the reference for the landing shape. The `Command Notes` sentence that says the core gives category-level guidance changes to name the reference itself. The `.bench/BENCH.md` budget row stays at 180; absence from the core is proven by registry rows, not by the budget.

The CLI contract in the core has three parts. First, the inventory sentences: `bench help` is the complete executable inventory, plumbing lives in the reference, and Bench is run exactly as its help spells it. Second, the ownership sentences: Bench owns non-interactive input, complete output, and the required next actions, and the complete output is the evidence. Third, the pair and the rule: "`bench gate` is valid. `bench gate 2>&1 | tail -20` is not valid." Do not append a subcommand, `</dev/null`, `2>&1`, a pipeline, or a shell follow-on.

### Anchors that move

Five registry rows change `File:` from the core to the reference. They are the kit-only ship row, the `bench upgrade` row, the no-PATH fallback row, the retro capture owner row, and the retro drain owner row. The core's retained-integration-source row is deleted, because the reference already carries the sibling row for the same fact. The integration-source anchor count in `docs_workflow_helpers_test.go` moves from 12 to 11. The rows on `Roles`, `How to talk to me`, and `Workflow` stay, because those sections stay.

One new row requires the CLI pair in the core. Two new rows require the two thresholds in `ste-prose.md`: the 25-word sentence bound and the six-sentence paragraph bound. Seven new rows require the moved units that carried no anchor, one per independently deletable unit, each at the reference:

- "guidance, not rules" for the skills bullet
- "with authority you do not have" for the gate-and-hooks bullet
- "Read `CONTEXT.md` for the mental model" for the `bench` bullet
- "Setup and adoption connect a repository to the kit" for the first category bullet
- "Context commands expose current state" for the second
- "Oracle commands inspect or enforce readiness" for the third
- "Work commands own isolated execution" for the fourth

With those rows, every moved unit has one needle, so a cut instead of a move reds the gate. Thirteen `Forbid` rows on the core prove absence, one per moved needle. Those needles are the seven above, the five moved rows' needles, and the retained-integration-source needle the core loses. A copy instead of a move leaves a needle in the core and reds its `Forbid` row. Each `Forbid` row has a `files/`-form fixture that plants the needle in the core.

Two `Forbid` rows, one on the core and one on the reference, forbid the retired term "progressive loading". The build writes one `files/`-form fixture for each new or moved row under `tests/canary/workflow-guidance-anchors/`. `capture-sink-anchor` stays, because its needle stays in the core.

### The prose mechanics check

`internal/prose` is a new package and a deep module. It owns the document parser, the tree walk, the exclusion-file grammar, and the subject classification, and it returns diagnostics for one root. The conformance check `prose-mechanics` is the registered wrapper that calls the package on the graded root, as `structure-accept-currency` wraps `ValidateAcceptGrants`. Its subject is `root` and its input source is `catch-all`.

When the graded root holds no `*.md` subject, the check returns no diagnostics. A canary fixture that plants an absent exclusion file can then clear on restore. The live root holds hundreds of subjects, so that guard never fires there. The check joins seven registries:

- the check list in `internal/conformance/registry`
- the canary family table
- the Go-fixture registry
- the check map
- the hostile-reader table
- the live-tree test list
- the profile's conformance table

The parser works in this order:

1. Strip a leading YAML frontmatter block.
2. Strip every HTML comment.
3. Skip every fenced code block.
4. Skip headings, table rows, thematic breaks, link reference definitions, and HTML blocks.
5. Skip an indented block of four or more spaces that starts after a blank line.
6. Treat every label line — a physical line that opens with a short label and a colon — as its own paragraph. Skip it when it holds no sentence terminator; that is a field line.
7. Split the rest into paragraphs at blank lines, and start a new paragraph at each list item, with the list marker removed.
8. Replace each inline code span with one token, keep link text, drop link targets, and strip emphasis markers.
9. Split each paragraph into sentences at `.`, `!`, `?`, or an ellipsis followed by whitespace, a closing quote or bracket, or the end.
10. Count a token as one word when it holds a letter or a digit; split tokens at Unicode whitespace. A list marker is not a token.

A period inside a token is not a boundary. The tokens `e.g.`, `i.e.`, `etc.`, `vs.`, and `cf.` are not boundaries. A list-item continuation line that follows a non-blank line is prose, not code. A `\r` is whitespace. The predicate for a word boundary is Go's `unicode.IsSpace`; a no-break space splits and a zero-width space does not.

Field lines are template values, not prose. A short label is at most four words. `Writes:`, `Blocked by:`, `Next:`, `Occurrences:`, `Feeds:`, and `Status:` lines are the common cases. A label line that holds a sentence terminator is prose and is graded; an `Occurrence:` line is the common case. Every label line is its own paragraph, so a run of twenty `Occurrence:` lines is twenty paragraphs. A drain that appends one more never deepens a paragraph.

The check is fail-closed. It reds with one distinct message for each of these states:

- a sentence over the bound
- a paragraph over the bound
- an unterminated fence, comment, or frontmatter block
- a symbolic link or special file named `*.md`
- an unreadable, non-UTF-8, or oversized subject
- a missing exclusion file
- a malformed row, a duplicate row, or a glob row
- a row that names an absent path

A diagnostic names the file and line and the counts. It never quotes the text, so control bytes cannot reach the gate output.

The walk grades every regular file named `*.md` under the graded root. It skips `.git/`, `node_modules/`, `dist/`, and every directory named `testdata/`. It does not descend a symbolic link to a directory. A file under a prefix named by an exclusion row is not graded. The engine classifies each subject through `bounds.ClassifyNoFollow`, the one producer classifier the prose-budget check already composes; it owns shape, bounded bytes, and the refusal reasons.

### The exclusion file

`.bench/prose-exclusions` holds one row per line: a repo-relative path, one space, and a one-clause reason. A trailing `/` names a directory prefix. A `#` starts a comment. The file's permanent rows are:

| path | reason |
|---|---|
| `tests/canary/` | fixtures carry planted content |
| `docs/audits/` | the audit record keeps model inputs and outputs verbatim |
| `CHANGELOG.md` | append-only release history |
| `capture/IDEAS.md` | an inbox of verbatim parked text |

During the build, the file also holds one temporary row for every authored area a later ticket migrates. Each migration ticket removes its own rows in the commit that makes its files pass. No ticket adds a row. When a batch cannot make a file pass, the build exits and reports; it does not add a row. Ticket 29 leaves only the four permanent rows.

The gate holds that rule. A live-tree test in the conformance package pins the approved row set and asserts that the live rows are a subset of it. Ticket 01c writes the file and the approved set as the initial rows. A batch that removes rows stays inside the set; a batch that adds one reds. Ticket 29 narrows the approved set to the permanent rows, so a later addition needs a visible test edit. That expectation is independent on purpose: the added-row mutation reds only because the test does not read the file it grades.

### Authoring and validation

A cheap-tier delegate authors one batch. It rewrites prose sentences only. It leaves frontmatter, fenced blocks, tables, check-read HTML markers, headings, anchor needles, and every inline-code token byte-identical. It writes in ASD-STE100 per `ste-prose.md` and, for comments, inside the register `craft-comments` owns. The orchestrator validates each returned batch before it lands, with four checks:

- a token-set comparison over the diff shows the same inline-code tokens before and after, except those the ticket names
- a line-class comparison shows prose lines changed and nothing else
- a read of each rule sentence confirms its constraint survives
- a read across the batch confirms one term per meaning from `CONTEXT.md`

A repair returns to a cheap-tier delegate with the finding and a sentinel.

A comment batch changes comment lines only. The orchestrator verifies the line class: every changed `.go` line starts with `//` after whitespace, and every changed shell line starts with `#`. A comment that restates its line is deleted. A provenance tag on an edited comment line is removed, per FT179's own rule. Directive comments and guard manifest headers are out of the rewrite. The `session-start.sh` `denies:` value stays `nothing (informational)`.

Go-embedded Markdown and shell that the kit writes into a tree is authored prose. That prose is the handoff shape text, the AGENTS managed block, the generated `CLAUDE.md`, the learnings preamble, and the profile scaffold. The scaffolded gate comments and the pre-push hook are its shell half. Ticket 28b rewrites the Markdown strings and regenerates `capture/session-handoff.md` in the same commit, because the handoff-shape check requires byte equality. Ticket 28 rewrites the shell comments, the pre-push hook and the gate scaffold strings in `internal/adopt` included.

Frontmatter `description`, `index`, and `name` values are trigger text and stay byte-identical. A skill's prose budget is reviewer-owned; a rewrite that would exceed a budget exits and reports. The `.bench/BENCH.md` budget row stays at 180.

### FT100 and FT179

FT100 keeps the demonstrated-delta audit of every always-loaded clause and craft skill, the cut line, and the budget measure decision; it stays blocked on FT231. This spec does not consume FT100's scope and does not touch that block. It delivers the core-versus-reference partition with anchors moved and the STE form across the surface. FT179 keeps the high-stakes surface documentation, the `craft-comments` rule additions, and the `craft-review` update. This spec delivers every explanatory comment rewritten under the current register, no-value comments deleted, and provenance tags removed on edited lines. Ticket 29 writes those two body updates and one `Occurrence:` line each, after every deliverer has landed.

### No new surfaces

No new `bench` verb, no comment lint, no budget change, no new dependency. The check and the package are the only new code. The exclusion file is the only new tracked file outside `specs/`, `internal/prose/`, and `tests/canary/`.

## Testing decisions

- The engine seam is `internal/prose`: table-driven unit tests feed one document and assert the findings, with each edge in the inventory as one case.
- The check seam is the conformance registry: `tests/canary/prose-mechanics/` holds one `files/`-form fixture per planted red, and the live-tree test grades the kit root.
- The anchor seam is `internal/anchors` plus `tests/canary/workflow-guidance-anchors/`: each new or moved row has one fixture.
- The regression seams are the existing checks:
  - shared-rule single source, the structured phase contract, the cold-pickup CLI list, and the token diet
  - the skills index, the stale-command sweep, the guidance token sweep, and the prose budgets
  - the row grammar, the roadmap detail integrity, the decision-map integrity, the handoff shape, and the guard headers
  - `gofmt`, `vet`, `test`, and `shellcheck`
- The embedded-template seam is a unit test in each owning package that renders the template and grades it through `internal/prose`.
- The exclusion-file seam is the live-tree subset test in the conformance package; its approved set is an independent expectation on purpose.
- The orchestrator's token-set and line-class comparisons are verification steps in each ticket charge; they are not gate checks.
- The project gate observes every row in its `test` phase; `shellcheck` observes the shell batch when present.

### Seam diagram

    trigger: `bench gate` runs the conformance registry over the graded root
        │
        ▼
    root + .bench/prose-exclusions ──▶ [ prose-mechanics: walk, exclude, parse, grade ] ──▶ diagnostics or none
                                            │
                                            └── internal/prose: one document ──▶ findings
                          ◀ unit tests attach at internal/prose; canary fixtures and the live-tree test attach at the check

    trigger: a migration ticket lands on the integration source
        │
        ▼
    batch diff ──▶ [ orchestrator: token-set, line-class, rule-read, term-read ] ──▶ accept or repair charge
                          ◀ the gate attaches through prose-mechanics, the anchors, and the regression checks

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PD1 | 1, 2 | every moved unit has one `Require` needle found at the reference and one `Forbid` needle absent from the core | anchors registry over both files | a build that cuts a unit instead of moving it loses its `Require` needle and reds; a build that copies a unit and leaves it in the core trips its `Forbid` row and reds |
| PD2 | 3 | each moved registry row names its new file and its needle is found there | `bench anchors` on both files and the docs-currency check | a moved passage with an unmoved row reds the live tree; a moved row with an unmoved passage reds too |
| PD3 | 4 | each new or moved row has one fixture whose planted absence reds with that row's diagnostic and goes green on restore | workflow-guidance-anchors canary family | a row with no fixture is never observed red |
| PD4 | 5 | the core contains the pair sentence and a registry row requires it | anchors registry and its fixture | deleting the pair reds the gate |
| PD5 | 6 | no `## ` heading appears in both the core and the reference | token-diet check | a moved section that keeps its heading in both files reds |
| PD6 | 7 | the eleven shared-rule markers stay in the core and appear in neither AGENTS.md nor README.md | shared-rule single-source check | a marker moved with its paragraph reds |
| PD7 | 8 | the core contains the string `BENCH-reference.md` and CLAUDE.md does not import the reference | token-diet check | a core that loses the pointer reds |
| PD8 | 1 | the core's invariant items and the predicate paragraph each hold six sentences or fewer | prose-mechanics on the live tree | the current items hold seven and eight sentences, so an unsplit core reds |
| PD9 | 9 | a document with one 26-word sentence returns one finding naming the line and count 26 | `internal/prose` unit test | a counter that stops at 25 or counts punctuation reds this case |
| PD10 | 9 | a document with one 25-word sentence returns no finding | `internal/prose` unit test | an off-by-one bound reds |
| PD11 | 10 | a paragraph of seven one-word sentences returns one finding naming the first line and count 7 | `internal/prose` unit test | a splitter that misses terminal punctuation passes it |
| PD12 | 10 | a paragraph of six sentences returns no finding | `internal/prose` unit test | an off-by-one bound reds |
| PD13 | 11 | a 40-token code span, a 40-word fenced block, a 40-cell table row, a 30-word heading, a 30-word frontmatter value, a 30-word HTML comment, and a 40-character link target each return no finding | `internal/prose` unit test | the cheapest parser grades every line as prose |
| PD14 | 11 | a block indented four spaces after a blank line returns no finding, and the same text as a list continuation after a non-blank line is graded | `internal/prose` unit test | a parser that treats every indent as code drops list prose |
| PD15 | 12 | a list item of seven sentences returns one finding, two adjacent list items of four sentences each return none, and a 25-word sentence under an ordered-list marker returns none | `internal/prose` unit test | a parser that joins list items into one paragraph reds the second case, and one that counts the marker reds the third |
| PD16 | 13 | `spec.md`, `1.2`, `e.g.`, `i.e.`, `etc.`, `vs.`, `cf.`, and `…` inside one sentence yield one sentence, and `word. Next` yields two | `internal/prose` unit test | a naive split on every period, or a table missing one of the five abbreviations, inflates the sentence count |
| PD17 | 13 | a sentence broken by a no-break space counts two words across it, and one broken by a zero-width space counts one | `internal/prose` unit test | an ASCII-only boundary or a zero-width boundary reds one side |
| PD18 | 14 | a root whose exclusion file names `docs/x/` with a reason grades no file under `docs/x/` | `internal/prose` unit test and the canary family | a check that ignores the file grades the excluded tree |
| PD19 | 15 | a row naming an absent path, a duplicate path, or a glob reds with its own message | prose-mechanics canary family | a lenient parser lets the file rot |
| PD20 | 16 | a symbolic link named `*.md`, a FIFO named `*.md`, a non-UTF-8 `*.md`, and a `*.md` one byte over `bounds.ControlRecordLimit` each red with the refusal message, and a `*.md` exactly at the limit is graded | `internal/prose` unit test | a plain read grades a link target, blocks on a FIFO, or loads an unbounded file |
| PD21 | 17 | a file with an unterminated fence, HTML comment, or frontmatter block reds naming the opening line | `internal/prose` unit test and the canary family | a truncating reader reports a clean file |
| PD22 | 18 | a long sentence under `tests/canary/`, `node_modules/`, `dist/`, `.git/`, or a `testdata/` directory is not graded | `internal/prose` unit test | a walk without the skip list reds the live kit |
| PD23 | 19 | every fixture in `tests/canary/prose-mechanics/` bites through the registered owner and clears on restore | fixture-bite test | a fixture bound to the wrong owner never bites |
| PD24 | 20 | the check over a hostile skill root names the refused skill path | hostile-reader composition test | a reader that swallows the refusal reports a clean skill |
| PD25 | 21 | the live-tree test reds when a temporary exclusion row is removed before its files pass, and goes green after the batch | live-tree test in the conformance package | a delegate with no focused seam runs the full gate or ships unchecked |
| PD26 | 22 | `ste-prose.md` carries the 25-word and six-sentence sentences and two rows require them | anchors registry and fixtures | editing the threshold in the rule file without the code reds |
| PD27 | 23, 24, 25, 26 | after each batch ticket lands, the live-tree test passes with the batch's rows removed | prose-mechanics on the live tree | an unrewritten file or a forgotten row keeps the red |
| PD28 | 25 | every anchor needle on a command or skill survives its rewrite | docs-currency check and the anchors registry | a rewrite that paraphrases a needle reds |
| PD29 | 27, 28, 29 | after the root, roadmap, and decisions batches land, the shared-rule, row-grammar, roadmap-detail, decision-map, and data-handling checks stay green | those checks | a rewrite that breaks a grammar line reds |
| PD30 | 30 | the inline-code token multiset of each rewritten file equals the original's, except tokens the ticket names | orchestrator token-set comparison | a rewrite that drops a path or command fails the comparison |
| PD31 | 31 | no changed line in a batch lies in frontmatter, a fence, a table, or a check-read marker | orchestrator line-class comparison | a rewrite that touches a machine-read surface fails the comparison |
| PD32 | 32 | the live rows of `.bench/prose-exclusions` are a subset of the approved set a live-tree test pins, and ticket 29 narrows that set to the permanent rows | live-tree subset test in the conformance package | a build that parks a hard file behind a new row reds the gate at that commit |
| PD33 | 33, 34, 35 | after each comment batch lands, `gofmt`, `vet`, `test`, and `shellcheck` stay green, and each `.go` and shell file in the batch reads in STE | gate phases and the orchestrator read | a comment edit that breaks a directive or a guard header reds the gate |
| PD34 | 36 | the guard manifest headers and `session-start.sh`'s `denies:` value are byte-identical after the shell batch | guard header check | a reflowed header key reds |
| PD35 | 37 | every changed line in a comment batch starts with `//` or `#` after whitespace | orchestrator line-class comparison | a behavior edit inside a comment pass fails the comparison |
| PD36 | 38 | the conformance registry gains exactly one check, `prose-mechanics` | registry inventory and the profile table | a second check added without its profile row reds the advertisement check; PD43 catches the careful variant that grades comment text |
| PD37 | 39, 40 | after ticket 29, `roadmap/FT100.md` and `roadmap/FT179.md` each hold the decided remaining-scope body and one new `Occurrence:` line naming this spec | roadmap-detail-integrity check and the orchestrator read | a record written before its deliverers land claims undelivered work, and a drain-only path leaves the old scope |
| PD38 | 41, 42, 43 | `CHANGELOG.md` holds one entry, `CONTEXT.md` holds the three terms, and `capture/session-handoff.md` is rewritten at each phase close | orchestrator read and the handoff-shape check | a build that skips the records leaves the next session cold |
| PD39 | 28 | every `roadmap/FT*.md` keeps its heading line, `Next:` line, and any `Occurrences:` ledger byte-identical | row-next-grammar and roadmap-detail-integrity checks | a reflowed `Next:` line reds |
| PD40 | 44 | the handoff shape constant and `capture/session-handoff.md` are byte-equal after ticket 28b | handoff-shape check | a rewritten template with a stale handoff reds |
| PD41 | 11 | `Writes:`, `Blocked by:`, `Next:`, `Feeds:`, and a four-word label line of thirty tokens with no terminator return no finding, an `Occurrence:` line and a five-word label line of thirty words ending in a period return one each, and twenty contiguous one-sentence `Occurrence:` lines return none | `internal/prose` unit test | a literal-label parser, a parser with no label-length bound, one that skips every label line, or one that joins a label run into one paragraph each fail one of the cases |
| PD42 | 44 | the rendered handoff shape and scaffold texts, the AGENTS managed block, the generated `CLAUDE.md`, the learnings preamble, and the profile scaffold each pass `internal/prose` with no finding | unit tests in `internal/handoff` and `internal/adopt` over `internal/prose` | a delegate that rewrites only the handoff template passes PD40 and leaves five templates long |
| PD43 | 38 | a `.go` file and a `.sh` file under the root, each with a 40-word comment sentence, are not graded | `internal/prose` unit test | a walk keyed on content instead of the `.md` name becomes the comment lint the reviewer excluded |
| PD44 | 45 | the core and the reference each carry a `Forbid` row for "progressive loading" with a fixture that bites | anchors registry and its fixtures | a later edit that reintroduces the retired term passes every other check |

The cheapest wrong implementations and the rows that red on them:

- Cut a core passage instead of moving it: PD1 and PD2.
- Move a passage and leave its row on the core: PD2.
- Skip the fixture for a moved row: PD3.
- Write the rule without the pair: PD4.
- Grade every Markdown line as prose: PD13.
- Treat every indent as code: PD14.
- Split on every period: PD16.
- Grade a field line as prose: PD41.
- Use an ASCII-only word boundary: PD17.
- Read the link target or block on a FIFO: PD20.
- Truncate at an unterminated delimiter: PD21.
- Walk the fixture tree: PD22.
- Park a hard file behind a new exclusion row: PD32.
- Rewrite only the handoff template and leave the scaffolds long: PD42.
- Grade comment text and become a comment lint: PD43.
- Reintroduce "progressive loading": PD44.
- Edit code under a comment pass: PD35.
- Reflow a guard header: PD34.

### Edge inventory

- Paths with spaces or glob characters in an exclusion row: a glob row reds. A path with a space is not expressible; the build names the first such path in its report.
- Hand-edited file with no trailing newline: the last paragraph still grades, and the last exclusion row still parses.
- Absent versus empty: an empty `*.md` returns no finding. An absent exclusion file reds when the root holds any `*.md` subject. An empty exclusion file is valid and excludes nothing. A root with no `*.md` subject yields no diagnostics.
- Special files and links: a FIFO or device named `*.md` reds before any read, and a link to a file reds. A link to a directory is not descended and not reported, because `.claude/` holds such links by design.
- Non-ASCII whitespace: the boundary predicate is `unicode.IsSpace`. A no-break space splits a word and a zero-width space does not; both sides are asserted.
- A grammar token quoted in prose: a period inside a code span never splits, because the span becomes one token first.
- Field lines: a label with no terminator is a field and is not graded; a label line with a terminator is prose. Every label line is its own paragraph, so a roadmap row's run of `Occurrence:` lines never forms one paragraph.
- List markers: `-`, `*`, `+`, and an ordered marker such as `3.` are stripped before the count, so a 25-word list item stays at 25.
- Unterminated delimiters: fence, HTML comment, and frontmatter each red naming their opening line.
- CRLF endings: `\r` is whitespace and never a word.
- Re-run idempotency: the check reads and writes nothing; two runs over one tree give one verdict.
- Cwd deeper than the root: the check receives the root as an argument, like every registered check.
- Control bytes in a document: the diagnostic carries counts and positions only, so the gate output stays representable.
- Control bytes in a path: the walk grades the file, and a diagnostic renders its path with Go's `%q`, so one diagnostic stays one line.
- In-flight tree states: a staged spec of another build sits under a temporary row. The check has no other special case. A pending retro, a learnings entry, a review pickup file, a light-path ticket, and a handoff are graded on their next commit. Their authors write them in ASD-STE100.
- **Won't handle:** the 20-word procedural bound — a program cannot tell an instruction from a description. Review holds that bound; the 25-word bound is the outer edge of both.
- **Won't handle:** raw HTML blocks as prose — a block whose first line starts with `<` is skipped. No authored file carries one today; the deepen report keeps its HTML inside fences.
- **Won't handle:** Markdown in `.mjs` and `.cjs` comments — the decision names Go and shell; the release scripts stay readable to their maintainers.
- **Won't handle:** a rewrite of `docs/audits/` and past `CHANGELOG.md` entries — they are records; the exclusion rows name them.
- **Won't handle:** grading a linked repo — the conformance registry runs only where the entry test is declared, which is the kit.

## Ownership fences

- `specs/asd-ste100-progressive-disclosure/`
- `tests/canary/prose-mechanics/`
- `tests/canary/workflow-guidance-anchors/`
- `.bench/prose-exclusions`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `.agents/`
- `.claude/README.md`
- `.claude/output-styles/`
- `AGENTS.md`
- `CLAUDE.md`
- `CONTEXT.md`
- `README.md`
- `CHANGELOG.md`
- `DATA_HANDLING.md`
- `SECURITY.md`
- `ASSESSMENT.md`
- `skills-assessment.md`
- `projects/benchkit.md`
- `projects/gl-axi.md`
- `docs/adr/`
- `docs/greenfield-build-sequence.md`
- `docs/release-runbook.md`
- `docs/reporesident-distillation.md`
- `ROADMAP.md`
- `roadmap/`
- `decisions/`
- `specs/inherited-toolchain-environment/`
- `tickets/`
- `capture/agent-performance/`
- `capture/audits/`
- `capture/FIXES.md`
- `capture/parallel-session-friction.md`
- `capture/learnings.md`
- `capture/session-handoff.md`
- `cmd/`
- `internal/`
- `consumer_payload.go`
- `consumer_payload_test.go`
- `bin/bench.sh`
- `bin/bench-postinstall.sh`
- `.bench/gate.sh`
- `.bench/gate-prospective.sh`
- `.bench/hooks/`
- `.bench/adapters/`
- `.bench/lib/`
- `scripts/aggregate-native-proofs.sh`
- `scripts/build-artifacts.sh`
- `scripts/build-offline-archives.sh`
- `scripts/compare-artifacts.sh`
- `scripts/gen-platform-packages.sh`
- `scripts/go-build.sh`
- `scripts/gremlins-diff.sh`
- `scripts/install-govulncheck.sh`
- `scripts/native-proof.sh`
- `scripts/release-preflight.sh`
- `scripts/smoke-artifacts.sh`
- `scripts/smoke-offline.sh`
- `scripts/lib/search.sh`

The wide `internal/` and `cmd/` prefixes authorize ticket 01's new package and check, and comment-only edits everywhere else. PD35 and the ticket acceptance rows hold each comment batch to comment lines. `capture/session-handoff.md` covers the phase-close handoff write and ticket 28b's regeneration; no other ticket names the handoff. `capture/IDEAS.md` and the `.mjs` and `.cjs` scripts sit outside the fence on purpose. `specs/inherited-toolchain-environment/` is in the fence only while that spec's folder exists at charge time.

## Out of scope

- FT100's demonstrated-delta audit, cut line, and budget measure — about 40 edits and 6 gate runs after FT231 lands.
- FT179's high-stakes surface documentation and `craft-comments` rule additions — about 25 edits and 4 gate runs.
- A `bench prose` AXI query over the engine — about 12 edits and 2 gate runs; the live-tree test is the delegate surface.
- Comment lint in the gate — about 15 edits and 2 gate runs; the reviewer excluded it.
- A prose mechanics check for linked repos — about 10 edits and 2 gate runs; needs a decision on where a consumer's exclusions live.
- ASD-STE100 rewrite of `.mjs` and `.cjs` comments — about 9 edits and 1 gate run.
- Lowering the `.bench/BENCH.md` budget to lock the split's gain — 1 edit and 1 gate run; FT100 owns the measure.

## Further notes

The closed decisions from the 2026-08-22 conversation:

- Audit all authored Markdown and explanatory Go and shell comments.
- Exclude fixtures, generated or vendored content, and machine-readable directives.
- Do not lint code comments in this spec.
- Add a hermetic, fail-closed full-tree Markdown check for mechanically reliable ASD-STE100 rules.
- Leave semantic rules to top-tier review.
- Use "progressive disclosure", not "progressive loading".
- Keep the named rule families in the always-loaded core and move detailed mechanics out with their anchors and proof.
- Cheap-tier delegates author, the top-tier orchestrator validates, and repairs return to cheap-tier delegates, as a spec-specific override.
- Reconcile FT100 and FT179 without duplicating their remaining work.
- Migrate the skills in the stated order, with thin adapters after their commands.
- Show `bench gate` as valid and `bench gate 2>&1 | tail -20` as invalid.
- Review with the same-family round before staging.
- Build in a worktree separate from `main`.

The reviewer pre-approved the spec commit and the author's judgment calls on 2026-08-22. These calls are flagged for post-hoc veto:

- The exclusion rows live in `.bench/prose-exclusions`, on the `.bench/structure-accept` line grammar, instead of a profile table; ninety-odd temporary rows do not suit a Markdown table.
- `CHANGELOG.md` and `docs/audits/` are records and stay excluded.
- `capture/IDEAS.md` stays excluded: `bench idea` appends the reviewer's words verbatim, so grading the inbox would make that command a gate hazard.
- Frontmatter values stay byte-identical, because they are trigger text.
- The universal sentence bound is 25 words; review holds the 20-word procedural bound.
- Ticket 29 edits `roadmap/FT100.md` and `roadmap/FT179.md` outside a drain, with the text this spec decides, after every deliverer has landed.
- `specs/inherited-toolchain-environment/` is rewritten by ticket 19 when its folder exists; when the folder is gone, its row is stale and leaves in the same commit.
- The roadmap bodies are one ticket of about 1,800 lines in 73 small files; the delegate works file by file.
- The `.bench/BENCH.md` budget row stays at 180; FT100 owns the measure and the cut line. Absence of a moved unit from the core is proven by one `Forbid` row per unit, not by the budget.
- Ticket 01 is a prefactor by design; `craft-tickets` puts prefactoring first as its own ticket. Ticket 01c is the tracer that registers the engine and makes it observable at the gate.
- Ticket 02, the partition and anchor moves, runs at mid tier against the leverage override. The disposition table is exact and every move is gate-graded. The prose rewrite of both files is ticket 02b, under the reviewer's cheap-author, top-validate override.

The reviewer-approved ticket frontier is the 41 tickets under `tickets/`. Ticket 01 opens the chain; ticket 29 closes it. Tickets 01 and 01c alone ship a working oracle. Every later batch is then a green follow-on that shrinks the exclusion file one commit at a time.

The review round's narrower-capability answer: the sixteen comment batches and the shell batch could ship alone on the existing gate phases. The reviewer scoped them into this spec. The review round used `gpt-5.6-sol` at xhigh through `codex exec` and `opus` at medium through the Claude agent surface, per the reviewer's 2026-08-22 instruction.
