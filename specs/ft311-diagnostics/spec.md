# FT311 diagnostic improvements

Status: staged

Decision source: specs/ft311-diagnostics/decisions/ft311-coordinator-work.md

Verification log: 2 iteration(s) to accept — Astra/medium reviewed through codex exec. Iteration one returned four blocking findings, the author folded all seven findings, and iteration two accepted with one accounting finding.

## Problem

A probe verdict can cite a test that was already red before the mutation, so a `bit` proves nothing.
A `--check` probe prints the conformance package and not the check it graded or the test it ran.
The probe help does not say that a named check compiles from the run binary, so a stale executable surprises a delegate.
A long-paragraph finding names the paragraph and not the sentence where a repair can split it.
A pre-commit prose check needs a hand-built path list, and it grades the working file rather than the staged bytes.
An anchor query lists needles without a line, so a before-and-after check costs a second read.

## Solution

Extend the three diagnostic verbs in place, with no new command and no change to the gate.
`bench probe` runs the selection once before the mutation, names the selection and the tests that ran, and explains its forms in help.
`bench gate-prose` prints the sentence starts of a long paragraph and takes a `--staged` form that grades the index.
`bench anchors` prints the line of each needle.

This is the second of five FT311 specs.
It serves this kit and repositories that link it.
It changes neither tier defaults, nor the landing, nor the `bench test` tables.
A probe verdict stays one of the four glossary words.

## User stories

### Name the probe selection

Line: gpt-5.6-terra / medium.
The projection composes the focused-run owner, and the help must derive from that owner.

1. As a coordinator, I want `bench probe --check <name>` to echo the check name, so that the evidence names the check that graded the mutation.
2. As a coordinator, I want the probe to name the tests it selected, so that I can tell selection from execution.
3. As a coordinator, I want the probe to count the tests that ran, so that a run with no execution cannot read as evidence.
4. As a coordinator, I want a package probe to name its package expression and run pattern, so that one row serves both selection forms.
5. As a delegate, I want the probe help to explain the `--check` and `--package` forms, so that I pick the right form.
6. As a delegate, I want the help to state that `--check` compiles from the run binary after `bench worktree build`, so that no stale executable surprises me.
7. As a maintainer, I want the help notes derived from the selection owner, so that help cannot drift from the behavior.
8. As a delegate, I want a form refusal to stay one usage line, so that a mistyped invocation gets the grammar and not the notes.

### Prove the baseline before the mutation

Line: gpt-5.6-terra / medium.
The baseline changes what a verdict means, so the oracle-adjacent effort applies.

9. As a coordinator, I want the probe to run the selection unmutated before it mutates, so that a `bit` verdict rests on a green baseline.
10. As a coordinator, I want a red baseline reported as `invalid` with the red tests named, so that an old red is never a bite.
11. As a coordinator, I want a red baseline to leave the subject and the Bench home untouched, so that a refused probe costs no restore.
12. As a coordinator, I want the selection row to state the baseline outcome, so that a `bit` verdict shows the baseline it rests on.
13. As a coordinator, I want a baseline with no verdict reported by its kind, so that a build failure is not a red suite.

### Show where a long paragraph breaks

Line: gpt-5.6-terra / medium.
The parser is the prose-mechanics oracle, so its edit takes the mid tier.

14. As an author, I want `bench gate-prose` to print each sentence start of a long paragraph with its line, so that I know where to split.
15. As an author, I want a sentence start to be the first three words of the sentence, so that two sentences on one line differ.
16. As a repository owner, I want the whole-tree prose check output unchanged, so that document text never reaches the gate output.
17. As an author, I want sentence starts quoted the way the sentence text is quoted, so that a control byte cannot split the diagnostic line.

### Grade the staged bytes

Line: gpt-5.6-terra / medium.
The form composes the Git adapter and the prose grader, and the index shapes are hostile.

18. As a coordinator, I want `bench gate-prose <root> --staged` to grade the staged Markdown without a path list, so that a pre-commit check costs one call.
19. As a coordinator, I want the staged form to grade the index bytes when the working file differs, so that it grades the commit's bytes.
20. As a coordinator, I want the staged form to read the exclusion policy from the index, so that policy and subjects come from one tree.
21. As a coordinator, I want a staged deletion and a staged symlink left out of the subjects, so that the form grades regular files only.
22. As a coordinator, I want a hostile staged path graded whole, so that Git quoting cannot hide a subject.
23. As a coordinator, I want an empty staged set to pass with an empty pass table, so that a clean index is a definitive answer.
24. As a coordinator, I want `--staged` with a path list refused as a usage error, so that the two subject sources cannot mix.
25. As an existing caller, I want the named-path form and the lane's prose check unchanged, so that the new form breaks no commit.

### Locate an anchor needle

Line: gpt-5.6-luna / low.
The grammar is exact, the seam is the existing query, and same-package tests cover it.

26. As a coordinator, I want `bench anchors <path>` to print the line of each needle, so that one call checks a needle before and after an edit.
27. As a coordinator, I want an absent needle to print line 0, so that a needle that disappeared is visible in the same call.
28. As a coordinator, I want a section-scoped needle located inside its section only, so that a match outside the section does not read as satisfied.
29. As a coordinator, I want the line of the match's first character, so that a split or comment-shifted needle still points into the file.
30. As a coordinator, I want a refused anchor file to answer a structured refusal, so that a link cannot read as an absent needle set.
31. As a maintainer, I want the anchors row to keep its AXI shape plus one column, so that the empty state stays definitive.

### Keep the guidance current

Line: gpt-5.6-terra / high.
Guidance prose compounds through every session, so the leverage override applies.

32. As a teammate, I want the reference and the glossary to describe the new behavior, so that a cold session learns it.
33. As a reviewer, I want the other three FT311 capabilities kept out of this build, so that each approved spec has one outcome.

## Implementation decisions

### Probe selection and help

Keep the probe grammar unchanged.
The help answer is the usage line, then a `notes:` block that the focused-run owner renders.
The owner's notes state three facts.

`--package <expr>` takes a Go package expression, as `bench test --package` does.
`--check <name>` names a conformance check from the `bench test --help` inventory, and `prose` and `system` are not probe targets.
A named check compiles from the run binary's source, so an edited tree needs `bench worktree build <target>` first.

The rebuild token and the excluded check names come from the constants the runner already uses.
A form refusal prints the one usage line and no notes.

The verdict row keeps its six cells and its four words.
A `selection[1]{form,target,run,baseline,ran}` row follows the verdict row and precedes the report.
The `form` cell is `check` or `package`.

The `target` cell is the check name or the package expression.
The `run` cell is the exact `-run` pattern the focused run passed to Go, or `all` when it passed none.
The `baseline` cell is the observed baseline outcome kind, and it joins the row with the baseline run.
The selection ticket renders the row without that cell, so no committed slice prints a baseline it did not run.
The `ran` cell counts the distinct tests that emitted a run event in the mutated run.
It reads 0 when that run did not start.

The focused-run owner exposes the selection facts and the run count.
One producer spells the named-check run pattern for both the Go argv and the selection row.
The owner's report tables do not change, because FT290 owns them.
The owner's command file is over the line budget, so the new code lands in a new file and that file does not grow.

### Probe baseline

The probe runs the prepared selection once over the unmutated tree after the gate-lock check and before preservation.
A baseline whose kind is `passed` lets the probe continue as today.
Any other kind ends the probe before any write.

The verdict is then `invalid`, the cause is `baseline-` followed by the kind, and the `failed_tests` cell counts the baseline's failing rows.
The `restored` cell reads `untouched`, because no mutation was written.
The report after the rows is the baseline report, so its failures table names the already-red tests.
`--full` reaches both runs.

The baseline holds the shared cache lock like every focused run.
The five non-passed kinds are `failed`, `build-failed`, `no-test-run`, `refused`, and `interrupted`, and each one has its own case.
A baseline that reaches no Go verdict answers `baseline-refused`.
An interrupt during the baseline answers `baseline-interrupted` and leaves the tree as it was.
The probe writes no record and reads no gate verdict.

### Paragraph sentence starts

A paragraph finding carries the line and the start of each sentence in the paragraph.
A start is the first three whitespace-separated words of the sentence as written, with an inline code span kept verbatim.
A sentence with fewer than three words gives the words it has.

The named-path render appends `: sentences ` and a comma-separated list of `<line> "<start>"` items, one per sentence, in document order.
Each start is quoted the way the sentence text is quoted, so a control byte is escaped and the diagnostic stays one line.
The whole-tree grade and the `prose` named check strip the starts as they strip the sentence text.
The parser file is over the line budget, so the paragraph grading moves to a new file and the parser file does not grow.

### Staged prose form

The grammar becomes `bench gate-prose <root> (--staged | [--] [path...])`.
The help keeps the single-file example and adds `example: bench gate-prose . --staged`.
`--staged` with a path or with `--` is a usage error at exit 2.
The root operand must be the top of a Git working tree.
A root that is not answers a `prose:` refusal line on stdout at exit 1.

The subject list is every index entry that differs from HEAD with a regular-file mode and a `.md` name.
The differing kinds are added, copied, modified, renamed, and type-changed.
On an unborn branch every index entry is staged.
The subject bytes come from the index blob, never from the working file.

The exclusion policy comes from the index blob of the exclusion file, and an absent entry answers the existing absent-policy diagnostic.
The policy constructor validates each exclusion target against the index entry list, never against the working tree.
A target that is an index entry is a file row, and a target that is a prefix of index entries is a directory.
A target that is neither reds as an absent path, and a directory row without a trailing slash reds as today.

A staged deletion and a staged symlink are not subjects.
The Git adapter owns the two index reads.
It lists paths with NUL framing, so a space, a quote, a tab, or a newline survives.
The prose grader gains a bytes entry point and a policy-from-bytes constructor, and it keeps one parser.

The pass table, the exit meanings, the named-path form, and the lane's prose argv do not change.

### Anchor line column

The anchors row becomes `anchors[N]{kind,section,needle,line}` with a typed integer `line`.
The line is the 1-based physical line, in the file as stored, of the first character of the first match.
A match follows the kind's own rule: whitespace collapse for every kind, and case folding for the section kinds.

The locate step maps a collapsed match back to a physical line by rune index.
So a case fold that changes byte length cannot shift the line.
HTML comments are stripped as the evaluator strips them, and a needle inside a comment has no match.
A section kind searches the section body only.

The line is 0 when the needle has no match, when the section is absent or duplicated, or when the file is absent.
For a forbid kind a non-zero line locates the violation.

A link, a special file, or an unreadable file at the path answers a structured refusal on stdout at exit 1.
The classification precedes the read.
The help description names the four cells.
The locate step lives below the conformance import edge, and the evaluator does not change.

### Guidance and glossary

The reference's probe paragraph states the baseline run, the selection row, and the staged prose form.
The glossary gains one term for the probe baseline and keeps the four verdict words.
The changelog records the four verb changes under the unreleased heading.

## Testing decisions

The primary seams are the three public command functions over temporary fixtures.
The probe keeps its stub Go and published run-binary harness, whose marker file counts every Go start.
The staged form uses a real temporary repository, because the index is the subject.
The anchors query moves its test to an isolated fixture repository, so it stops reading the live tree.
Unit seams cover the locate step and the paragraph starts inside their own packages.

Every new test runs in the ordinary Go phase of the project gate.
Guidance edits are review-owned, because no anchor pins a descriptive sentence and this spec adds no anchor.
The full landing gate remains the code oracle.

### Seam diagram

    trigger: coordinator or delegate command
        |
        v
    argv --> [ probe.Command ] --> verdict row, selection row, report
                 ^ tests attach here: stub Go, marker file, canned JSON
        |
        v
    argv --> [ gate.GateProseCommand ] --> pass table or findings
                 ^ tests attach here: temp root, real index, encoder-derived table
        |
        v
    argv --> [ anchorsCommand ] --> anchors table with line, help envelope
                 ^ tests attach here: isolated fixture repository, anchors.Locate

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DG1 | 1, 2 | A check probe prints a selection row with form `check`, the check name, and the root conformance run pattern. | New TestProbeNamesTheCheckSelection through probe Command | An omitted row or a package-shaped row fails the exact comparison. |
| DG2 | 4 | A package probe prints form `package`, the expression, and the `--run` pattern or `all`. | New TestProbeNamesThePackageSelection through probe Command | A row that copies the check form or drops the pattern fails the exact comparison. |
| DG3 | 3 | The `ran` cell counts the distinct tests with a run event in the mutated run. | New TestProbeCountsTheTestsThatRan through probe Command | Two run events print 2, and a no-test stream prints 0 beside verdict `invalid`. |
| DG4 | 5, 6 | `bench probe --help` prints the usage line and the three notes. | Existing TestProbeHelpSpellings extended through probe Command | A dropped note line reds the exact expectation. |
| DG5 | 7 | The probe help embeds the focused-run owner's notes rather than a copy. | Review-owned Standards axis over the help producer | A second copy of the note text is duplicated knowledge. |
| DG6 | 8 | A form refusal prints only the usage line. | Existing TestProbeUsageNamesOneSelectionAndOneMutation through probe Command | A refusal that prints the notes fails the one-line expectation. |
| DG7 | 9 | The probe starts one Go test before the mutation and one after it. | New TestProbeRunsTheBaselineFirst through probe Command | The marker records two starts, and the stub's first-start copy of the subject holds the unmutated bytes. |
| DG8 | 10 | A failing baseline answers verdict `invalid`, cause `baseline-failed`, and the baseline failures table. | New TestProbeRefusesARedBaseline through probe Command | A stub that fails on the first start must not answer `bit`. |
| DG9 | 11 | A red baseline writes no mutation and leaves the Bench home empty. | New TestProbeRefusesARedBaseline through probe Command | The subject bytes and the home listing red on any write. |
| DG10 | 12 | The selection row's baseline cell reads `passed` on a bite and `failed` on a red baseline. | New TestProbeNamesTheCheckSelection and TestProbeRefusesARedBaseline | A constant cell fails one of the two comparisons. |
| DG11 | 13 | A baseline with no verdict names its kind in the cause. | New TestProbeReportsABaselineWithoutAVerdict through probe Command | A build-failed stream without `--run` and a no-test stream give `baseline-build-failed` and `baseline-no-test-run`, never `baseline-failed`. |
| DG12 | 9 | `--full` reaches the baseline report. | Existing TestProbeForwardsFullToTheFocusedRun extended | A long baseline diagnostic is previewed without the flag and complete with it. |
| DG13 | 14 | A paragraph finding from the named form lists each sentence's line and start in order. | New TestGateProseCommandNamesTheSentenceStarts through GateProseCommand | Seven sentences across three lines give seven items with the right lines. |
| DG14 | 15 | A start is the first three words as written, and a short sentence gives its words. | New TestParagraphSentenceStarts through prose Findings | Two sentences on one line give two distinct starts, and a code span stays verbatim. |
| DG15 | 16 | The whole-tree grade prints the paragraph line with no starts. | Existing TestGrade extended through prose Grade | The diagnostic must equal the renderer's line exactly. |
| DG16 | 17 | A control byte in a start is escaped and the finding stays one line. | New TestParagraphStartsEscapeControlBytes through GateProseCommand | A raw ESC in the output reds. |
| DG17 | 18 | `--staged` grades the staged Markdown and prints its pass table for a clean index. | New TestGateProseStagedGradesTheIndex through GateProseCommand over a real repository | An empty pass table on a staged file fails the comparison. |
| DG18 | 19 | The staged form grades the index bytes when the working file differs. | New TestGateProseStagedIgnoresTheWorkingFile through GateProseCommand | A long index blob with a short working file reds, and the inverse passes. |
| DG19 | 20 | The staged form reads the exclusion policy from the index. | New TestGateProseStagedReadsTheStagedPolicy through GateProseCommand | An index policy that excludes the subject passes while the working policy would red, and the reverse reds. |
| DG20 | 21 | A staged deletion and a staged symlink are not subjects. | New TestGateProseStagedSkipsDeletionsAndLinks through GateProseCommand | Only the long regular file beside them is reported. |
| DG21 | 22 | A staged path with a space, a double quote, a tab, or a newline is graded whole. | New TestGateProseStagedHostilePaths through GateProseCommand | The pass table derives through the encoder, and a split or dropped path fails it. |
| DG22 | 23 | An empty staged set passes with an empty pass table. | New TestGateProseStagedEmptyIndex through GateProseCommand | The output must equal the empty table at exit 0. |
| DG23 | 24 | `--staged` with a path or a `--` is a usage error at exit 2. | New TestGateProseStagedRefusesAPathList through GateProseCommand | Any stdout or a graded subject reds. |
| DG24 | 25 | The named form, the help example, and the lane argv keep their bytes. | Existing TestGateProseCommandCleanList, TestGateProseCommandHelp, and lane_test expectations | A grammar change to the named form reds the pinned strings. |
| DG25 | 18 | A root that is not a working-tree top answers a `prose:` refusal at exit 1. | New TestGateProseStagedOutsideARepository through GateProseCommand | A plain temporary directory must not answer an empty pass. |
| DG26 | 26 | The anchors table carries a `line` cell with the physical line of each present needle. | New TestAnchorsReportsNeedleLines through anchorsCommand in an isolated fixture repository | A three-cell row or a zero line on a present needle fails the exact comparison. |
| DG27 | 27 | An absent needle prints 0, and an absent file prints 0 on every row. | New TestAnchorsReportsAbsentNeedles through anchorsCommand | A removed needle must read 0 while its siblings keep their lines. |
| DG28 | 28 | A section needle outside its section prints 0, and one inside prints its line. | New TestLocateHonorsSectionScope in the anchors package | A whole-file search reds the outside placement. |
| DG29 | 29 | A needle split across lines, after a multi-line comment, or after non-ASCII uppercase text prints its first character's line. | New TestLocateMapsCollapsedMatchesToLines in the anchors package | Byte-offset locating reds the comment and case-fold shapes. |
| DG30 | 30 | A symlink, a FIFO, or an unreadable anchor file answers a structured refusal at exit 1. | New TestAnchorsRefusesAnUnreadableFile through anchorsCommand | A plain read prints 0 on every row at exit 0 for the link, and blocks in open on the FIFO. |
| DG31 | 31 | The empty state is the four-cell header plus the empty help envelope, and the help names four cells. | Existing TestAXIRegistryBindsEachRealCommandEnvelope and TestHelpInventoryIsComplete updated | The pinned markers red on the old shape. |
| DG32 | 32 | The reference and the glossary describe the baseline, the selection row, the staged form, the starts, and the line. | Review-owned Spec axis over the two files | A missing paragraph leaves a cold session on the old behavior. |
| DG33 | 33 | No landing, reset, tier-default, or `bench test` table change enters the build. | Review-owned scope audit | An added report column or a landing edit is outside the fence. |
| DG34 | 20 | An exclusion target present in the index and absent from the working tree is honored, and a target absent from the index reds the policy. | New TestGateProseStagedValidatesTargetsAgainstTheIndex through GateProseCommand | A working-tree stat reds the first case and passes the second. |
| DG35 | 13 | A baseline that reaches no Go verdict answers `baseline-refused` with `untouched`. | New TestProbeReportsABaselineRefusal through probe Command | A stub Go that prints no event and exits 1 reaches the no-packages refusal, and a collapsed cause reds. |
| DG36 | 13 | An interrupt during the baseline answers `baseline-interrupted` and `untouched` with the subject and the home untouched. | Existing TestProbeRestoresOnInterrupt split into a first-start case and a second-start case | The first-start signal must not print `interrupted` with `yes`, and the second-start signal keeps the restored row. |
| DG37 | 18 | On an unborn branch the staged form grades every index entry. | New TestGateProseStagedUnbornBranch through GateProseCommand | A form that diffs against HEAD refuses or grades nothing. |
| DG38 | 20 | An index without the exclusion file answers the absent-policy diagnostic at exit 1. | New TestGateProseStagedAbsentPolicy through GateProseCommand | A form that treats an absent policy as empty passes. |
| DG39 | 28 | A duplicated owning section prints 0 for its section needle. | New TestLocateHonorsSectionScope in the anchors package | A locate that searches the first duplicate prints its line. |
| DG40 | 18 | A relative root that names a working-tree top grades the staged form. | New TestGateProseStagedAcceptsARelativeRoot through GateProseCommand | A top compared as spelled refuses `.`, which is the help's own example. |
| DG41 | 19 | A staged blob that is invalid UTF-8 or over the record limit answers the named form's unreadable-subject diagnostic at exit 1. | New TestGateProseStagedRefusesAnUnboundedBlob through GateProseCommand | An unbounded index read grades bytes the named form refuses. |
| DG42 | 22 | A staged path with an ESC byte answers the shared unrepresentable-cell refusal at exit 1. | New TestGateProseStagedRefusesAControlBytePath through GateProseCommand | A refusal that names no path is the documented choice, and a pass table would be a lie. |
| DG43 | 29 | A comment that the strip creates from rejoined text hides its needle from the locate step as it hides it from the evaluator. | New TestLocateStripsRejoinedComments in the anchors package | A single-pass strip prints a line for a needle the evaluator reads as absent. |

### Edge inventory

The canonical walk covers empty input, boundaries, errors, repetition, process boundaries, and hostile environments.
The attached profile is the shell CLI checklist in projects/benchkit.md.
Every behavior serves this repository and every repository that links the kit, except the named-check note, which names a kit-only rebuild verb.

| Class | Concrete disposition | Rows |
|---|---|---|
| Absent versus empty | An absent anchor file gives 0 on every row, and an empty index gives an empty pass table. An absent exclusion entry in the index keeps the absent-policy diagnostic. | DG22, DG27, DG38 |
| Sections | A section needle outside its section, or under a duplicated heading, prints 0. | DG28, DG39 |
| File kinds | A link, a FIFO, or a device at the anchor path refuses before the read. A staged symlink is not a prose subject. | DG20, DG30 |
| Paths | A space, a double quote, a tab, and a newline in a staged path survive NUL framing and the encoder. An ESC byte in a staged path refuses the run. | DG21, DG42 |
| Encoding and size | A staged blob that is invalid UTF-8 or over the record limit is refused as the named form refuses it. | DG41 |
| Text sinks | A control byte in a sentence start is escaped, and the whole-tree grade carries no document text. | DG15, DG16 |
| Grammar | `--staged` beside a path or a `--`, and a probe form refusal, answer one usage line at exit 2. | DG6, DG23 |
| Index versus working tree | The staged bytes, the staged policy, and the policy's target validation win over the working tree. An unborn branch stages every entry. | DG18, DG19, DG34, DG37 |
| Whitespace and case | A needle split across lines and a case fold that changes byte length still map to the first character's line. | DG29 |
| Comments | A needle inside an HTML comment has no match, text after a multi-line comment keeps its physical line, and a comment the strip creates hides its needle. | DG27, DG29, DG43 |
| No execution | A baseline or a mutated run with no test run cannot read as evidence, and a baseline with no Go verdict is refused. | DG3, DG11, DG35 |
| Already red | A red baseline is `invalid` and names the red tests, and it writes nothing. | DG8, DG9 |
| Interruption | An interrupt during the baseline answers `baseline-interrupted` and leaves the tree untouched. | DG36 |
| Repetition | Two probes on one repository keep separate preserved copies, as today, and the baseline adds no record. | DG9 |
| Process boundary | The stub Go proves the baseline start order from outside the verb's process. | DG7 |
| Root | A root that is not a working-tree top refuses the staged form, and a relative root that is a top grades it. | DG25, DG40 |

**Won't handle:** A named green baseline from the gate cache instead of a baseline run is a separate capability. Every probe runs its baseline, so no caller loses a verdict.

**Won't handle:** A per-package `tests_run` column and a `--fixtures` family projection belong to FT290. The probe's `ran` cell reads the same report.

**Won't handle:** A default root for `bench gate-prose` under `bench worktree exec` belongs to FT254. The staged form takes an explicit root.

**Won't handle:** A per-file anchorHarness diagnostic slice belongs to FT120. The anchors query serves the before-and-after read.

**Won't handle:** A `--staged` subject for the lane's prose check is not needed, because the lane already grades the composed tree. The lane argv keeps its bytes.

No new Git flag beyond the observed index reads, no package-variable substitution, and no deletion of a tree file is required by this design.
If implementation needs one, its ticket must attach the corresponding hostile case before it claims its row.

## Ownership fences

These entries are the union of the ticket Writes.
They authorize implementation only after the reviewer approves this spec and its ticket graph.

- `internal/bounds/classify.go`
- `internal/bounds/classify_bytes_test.go`
- `internal/anchors/locate.go`
- `internal/anchors/locate_test.go`
- `cmd/bench/anchors_command.go`
- `cmd/bench/anchor_help_test.go`
- `cmd/bench/testdata/anchors`
- `cmd/bench/main.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/command_registry.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/docs-currency-token-diet`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/workflow-guidance-anchors`
- `internal/prose/prose.go`
- `internal/prose/subject.go`
- `internal/prose/starts.go`
- `internal/prose/parse.go`
- `internal/prose/parse_test.go`
- `internal/prose/subject_test.go`
- `internal/prose/walk_test.go`
- `internal/prose/walk.go`
- `internal/prose/exclusions.go`
- `internal/gate/gate_prose.go`
- `internal/gate/gate_prose_test.go`
- `internal/gate/gate_prose_staged_test.go`
- `internal/gate/gate_prose_staged_root_test.go`
- `internal/git/staged.go`
- `internal/git/staged_test.go`
- `internal/probe/command.go`
- `internal/probe/probe.go`
- `internal/probe/baseline.go`
- `internal/probe/probe_test.go`
- `internal/probe/baseline_test.go`
- `internal/probe/refusal_test.go`
- `internal/probe/outcome_test.go`
- `internal/testreport/selection_facts.go`
- `internal/testreport/selection_facts_test.go`
- `internal/testreport/outcome.go`
- `internal/testreport/outcome_test.go`
- `internal/testreport/testreport.go`
- `internal/testreport/command.go`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `CHANGELOG.md`
- `reviews/ft311-diagnostics.md`
- `ROADMAP.md`
- `roadmap/FT311.md`
- `decisions/ft311-coordinator-work.md`
- `decisions/ft311-coordinator-work`

The last four entries are the spec phase's own change, which the landed range carries.
They are the two roadmap references to the moved map and the map's source paths, and no ticket writes them.
The command registry, its two pinned conformance tests, and the four canary entries are closure declarations from the preflight proposal.
The registry is named because the probe and the anchors query are bound commands, and the canary families pin the guidance files.
The expected edit inside them is the anchors description and the three-cell pins, and no guidance fixture changes.

The coordinator owns the conditional review pickup, so no ticket names it.

The build cannot edit this spec, its acceptance rows, or its ticket graph without the existing spec-change authority.
Four fenced files are over the line budget: `cmd/bench/main.go`, `cmd/bench/command_registry_test.go`, `internal/prose/parse.go`, and `internal/testreport/command.go`.
Each of them changes in place or sheds lines in the same commit, and none of them gains a line.

## Out of scope

These are planning estimates of file edits, not measured implementation costs.
Each future capability gets its own specification and one whole-project landing gate.

| Capability | Estimated price | Derivation |
|---|---|---|
| Named green baseline from the gate cache | 4 edits, 1 gate runs | A gate-cache reader in the probe, a `named` baseline cell, and two tests |
| Landing completion | 16 edits, 1 gate runs | Landing effects, broker refresh, cleanup, handoff state input, and retrospective facts with tests |
| Recoverable reset | 10 edits, 1 gate runs | Worktree grammar, checkpoint preservation, plan validation, restore, and closure tests |
| Lower-tier trials | 6 edits, 1 gate runs | Task definitions, twelve matched runs' evidence projection, comparison, and reviewer decision record |

The remaining order is landing completion, recoverable reset, then trials.
FT120, FT254, FT258, and FT290 retain their neighboring subjects.
The existing merge lane and the landing gate remain unchanged.

## Further notes

### Review round and dogfood runs

The initial review at tip `9eed6663e1c2fc9ae656b1ffe9903f1393fa10bc` over base `f41917040c24637ca0c5d69441398361559ce319` returned 15 findings across the three axes and 9 repair targets.
Ticket 7 carries the accepted repairs and rows DG40 through DG43.
The build extended three ticket fences in range for test files that did not fit the line budget. It flags each extension for reviewer veto.

The coordinator ran three dogfood runs on the reviewed tip.
A `bench probe` swap of one word in the owner's notes reddened the three help spellings, so the independent help expectation is the mutation catch.
A `bench probe` swap that dropped the `$` from the named-check pattern reddened the facts test and the argv test. So the independent pattern expectation is the mutation catch.
A hand run of `bench gate-prose . --staged` from the worktree top refused the relative root, which is DG40.

### Decision provenance

The whole ready map and its topic folder moved together into this spec's decisions directory.
The two roadmap references to the old path moved with it in the same change.
Future FT311 specs consume this same source in place.
Before this spec is retired, move the shared map back to the top-level decisions folder and repair its references in the same change.
Do not retire the FT311 roadmap row when only this second capability lands.

The three structured map sources were reread on 2026-09-09.
The research report's implementation claims were checked against the current probe, focused-run, prose, and anchors owners.
The report's session counts remain observational evidence, not a target for this spec.
The FT311 detail file's diagnostics group is the fold list this spec settles.
The kit profile's tier binding and lane table were reread, and neither changes here.

### Source-sentence-to-row accounting

References below identify resolved tickets in the single decision source and the FT311 detail file's diagnostics group.
Each table entry paraphrases its assigned clause without creating another decision source.

| Source clause | Disposition |
|---|---|
| FT311: `--check` names the tests it selected and echoes the check name | DG1, DG2, DG3 |
| FT311: help states that `--check` compiles from the run binary after `bench worktree build` | DG4, DG5 |
| FT311: help states that `--check` names a conformance check and `--package` takes a Go expression | DG4, DG5 |
| FT311: the probe runs or names a green baseline, or its verdict names an already-red suite | DG7 through DG12, and the named baseline is priced out of scope |
| FT311: `bench gate-prose` prints the sentence starts of a long paragraph | DG13 through DG16 |
| FT311: `bench gate-prose --staged` grades the staged files without a path list | DG17 through DG25 |
| FT311: `bench anchors` prints a line column for a before-and-after check | DG26 through DG31 |
| Ticket 1: probe output identifies selection and execution, help explains build and selection, prose exposes starts and staged contents, anchors supply lines | DG1 through DG31 |
| Ticket 3: the CLI derives facts and invents none | DG3, DG5, DG7, DG10 |
| Ticket 17: staged bytes win over working files | DG18, DG19 |
| Ticket 17: long-paragraph diagnostics identify the paragraph and its sentence starts | DG13, DG14 |
| Ticket 17: anchors report source line locations | DG26, DG29 |
| Ticket 17: probe output distinguishes selected checks from tests that ran | DG1, DG3 |
| Ticket 17: no execution cannot appear as successful evidence | DG3, DG11 |
| Ticket 17: help describes the actual behavior without a second implementation | DG4, DG5 |
| Tickets 10, 18, and 19: five specs, one source, fixed order | DG33 and Out of scope |
| Out of scope: no fast-lane or landing-gate change | DG24, DG33 |

### Reader sweep and enforcement reads

The repository sweep used hidden-file coverage over Go, Markdown, fixtures, scripts, and workflow files.
Historical audits under docs and the closed FT303 assessment name the verbs without a shape claim, so they need no edit.

| Reader or owner | Read evidence and disposition |
|---|---|
| probe Command, run, render, and verdictFor | The verdict row, the refusal order, and the exit mapping. DG1 through DG12 extend them without a new verdict word. |
| testreport Prepare, Execute, Outcome, and runNamedCheck | The one selection grammar and the run-binary source root. DG1 through DG5 expose facts it already holds. |
| testreport report render | The `packages`, `failures`, and `skips` tables. They do not change, and FT290 owns them. |
| worktree merge caller of testreport Command | An unchanged caller of the unchanged report. |
| prose Findings, Render, GradeSubject, GradeNamed, GradeNamedResults | The one parser and the two renders. DG13 through DG16 extend the named render only. |
| prose loadExclusions | The policy loader. DG19 adds a bytes constructor beside it. |
| testreport runProseCheck | The `prose` named check reads Grade, whose output does not change. |
| gate GateProseCommand and parseGateProseArgs | The grammar and the pass table. DG17 through DG25 add one flag. |
| gate BenchkitLane and resolveLane | The lane's prose argv and its named-markdown token. They do not change, and lane_test pins them. |
| git Output, Raw, OK, and RootAt | The adapter forms the staged reads compose. |
| anchorsCommand, anchorQueryPath, and anchors Entries | The query and its registry read. DG26 through DG31 add the locate step. |
| anchors EvaluateGroup, Satisfied, and read | The evaluator's match rule and comment strip. The locate step reuses both and the evaluator does not change. |
| cmd/bench anchors testdata and anchor_help_test | Pinned three-cell responses over the live tree. Ticket 1 replaces them with an isolated fixture. |
| command_registry_test AXI matrix and help_inventory_test | The empty marker and the help description. Ticket 1 updates both pins. |
| axi_query_registry_test and craft-cli table | Membership and disclosure only. No inventory change. |
| subcommand_routing_table_test and entry_point_parity_test | Route ownership for `gate-prose` and `probe`. No new command. |
| projects/benchkit.md lane table | The prose lane argv. Unchanged. |
| CONTEXT.md probe verdict term | The four words stay. Ticket 6 adds the baseline term. |
| BENCH-reference probe paragraph | Ticket 6 states the baseline and the selection row. No anchor pins the paragraph. |
| craft-delegate probe sentences | They name the verb and its restore. They stay true and need no edit. |

### Pre-review proof checklist

- Cited symbols: every symbol in the reader sweep table above was read in the current tree.
- Import edges: probe already imports testreport and gate. gate already imports prose and the git adapter. cmd/bench already imports anchors. The prose and anchors packages import neither git nor conformance, and this design keeps that.
- Source-row clauses and occurrences: the accounting table enumerates the seven fold sentences and the ticket clauses. No source clause is rewritten by this spec.
- Promised field labels: `selection`, `form`, `target`, `run`, `baseline`, `ran`, `untouched`, `baseline-failed`, `baseline-build-failed`, `baseline-no-test-run`, `baseline-refused`, `baseline-interrupted`, `notes:`, `--staged`, `sentences`, `line`.
- Changed-function callers: probe render and run are called inside the probe only. testreport Prepare and Execute are called by the probe, and Command by the registry and the merge verb. prose Render is called by RenderNamedResult only. GateProseCommand is called by the registry only. anchorsCommand is called by the registry and its test.
- Copy survival: none. No copy is replaced by a new owner.
- Git flags: `diff --cached --name-only -z --diff-filter=ACMRT --no-renames`, `ls-files --stage -z`, and `show :<path>` ran on 2026-09-09 under Git 2.43.0 over the hostile shapes. On an unborn branch the listing named the one added file. With the index at `Two.` and the working file at `Three working.`, the blob read answered `Two.`. The deletion was absent under the filter, and the rename listed only its new name. A link-to-file type change listed under `T`, and the symlink entry carried mode 120000 and answered the link target. The listing framed a space, a double quote, a tab, and a newline whole, and the blob reads for those paths answered their bytes.

### Flagged additions and engineering choices

No new product promise beyond the reviewed map is intended.
The `selection` row and its cells, the `baseline-` causes, the `untouched` cell, the `notes:` block, the `sentences` suffix, and the `line` cell are engineering choices.
The always-run baseline is a choice inside the source's "runs or names" clause, and the named form is priced out of scope.
The index-sourced exclusion policy follows the lane, which grades the composed tree's policy.
The three-word start is a choice, and a reviewer can widen it without a row change.

### Implementation ticket approval table

| Ticket | Blocked by | Delivered outcome |
|---|---|---|
| 1. Locate anchor needles | none | `bench anchors` prints a typed line per needle and refuses an unreadable file |
| 2. Print the sentence starts of a long paragraph | none | The named prose form lists each sentence's line and start |
| 3. Grade the staged Markdown | 2.md | `bench gate-prose <root> --staged` grades the index bytes and policy |
| 4. Name the probe selection and its help | 1.md | The four-cell selection row and the owner-derived help notes |
| 5. Run the baseline before the mutation | 4.md | The baseline run, its cell in the selection row, and the five non-passed kinds |
| 6. Fold the guidance and the changelog | 3.md, 5.md | The reference, the glossary, and the changelog state the landed behavior |

Tickets 1 and 2 form the first frontier and run in parallel, one worktree each.
Ticket 4 follows ticket 1 only because both name the command registry closure, and tickets 3 and 4 run in parallel.
Ticket 3 carries the prose and gate invariants, and ticket 5 carries the probe invariant.
Ticket 4 carries the focused-run owner invariant.
Ticket 6 is the last writer of the shared guidance files.
The review pickup is created only for actionable findings.

### Approval surface

| Subject | Disposition requested |
|---|---|
| Stories and lines | Approve the six outcome groups and their bound model efforts. |
| Seams | Approve the three public command seams and the two unit seams. |
| Acceptance and edges | Approve DG1 through DG39 and the five explicit exclusions. |
| Ownership fences | Approve the exact union above for implementation. |
| Scope and tickets | Approve the six-ticket graph within the second FT311 capability. |

### Verification outcome

The reviewer-named model reviewed the spec and the six tickets twice through `codex exec` in a read-only sandbox.
Iteration one, at `1e2700c07deb61210a491f0a54ab5fa5cbd90631`, returned `revise` with four blocking findings and three non-blocking findings.
The blocking findings were a constant baseline cell in ticket 4 and working-tree validation of staged exclusion targets.
The other two were two baseline kinds without a case and three stated states without a fixture.

The author folded all seven findings into rows DG34 through DG39, the two probe tickets, ticket 3, the fence paragraph, and the Git-flag note.
Iteration two, at `2de046fdce479cc3c465e2a5390b7940061e5331`, verified every fold closed and accepted with one accounting finding, which the author folded.
The journal carries the two-iteration learning under the title `ft311-diagnostics spec round needed two iterations`.

The build preflight was green on every check after the closure fold.
The coverage map validated at 39 rows, and the prose lane passed on the spec and the six tickets.
The reviewed graph, fences, rows, and lines are the approval surface above.
