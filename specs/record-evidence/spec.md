# `bench record` writes review-record evidence

Status: staged

Decision source: reviewer-delegated coordinator decisions, 2026-10-01, on the learning "Review record evidence needs a tool, not hand-built JSON". Related roadmap row: FT318.

This spec answers the FT318 question "a verb, or a smaller payload?" in favour of
a verb. After this spec, FT318 keeps only the completion form. The next drain
narrows FT318 to that form and records the name choice `bench record` over
`bench review record`.

Verification log: 1 iteration(s) to accept — iteration 1 (fable, high) approved with 11 findings and no blocker. The author folded all 11 findings and two advice items on the coordinator's decisions.

## Problem

Each spec-backed build retains its evidence in the `bench-review-record` fence of
`reviews/<slug>.md`. The checkpoint grades that fence. No verb writes it, so
authors and coordinators build the JSON entries by hand or with a script.

In the worktree-verb-runner build, `bench gate --checkpoint` refused hand-built
entries several times. An excerpt digest hashed the whole command output, not
the embedded excerpt. A source digest named the full tree, not the tree without
the record file. A verification entry sat in the review list. Plan digests were
computed by hand from the `ReadPlan` rule. Each error cost a record commit, and
some cost a confirming review round.

## Solution

One new verb, `bench record`, writes each evidence entry. It has four forms:

- `bench record chunk` writes the frozen pair of one chunk with its digests and acceptance rows.
- `bench record verification` appends one planned verification result to a chunk or to the completion list.
- `bench record review` appends one independent review result to a recorded chunk.
- `bench record amendment` records a plan-digest change.

The verb computes every digest with the functions the checkpoint uses. It puts
each entry in the correct list of the correct chunk. It parses the whole record
before it writes, so a refusal leaves the file unchanged. It writes only
`reviews/<slug>.md` in the current worktree, and the caller commits the file
with `bench commit`. The phase guidance names the verb in place of hand-built JSON.

## User stories

Line: opus / medium.
Implementation-line reason: RE-C2 is the hardest material chunk, because it adds the write transaction, the command family, and the registry closure that later tickets consume. The spec fixes each form, each derived field, and each output header. The seams exist: the `recordtest` fixtures build real plans, and `reviewrecord.Check` is an independent oracle that reds every wrong derived field.
Harder chunks: RE-C2.

Chunk entries:

1. As a coordinator, I want the chunk form to store full commit IDs, so that I never paste a short or symbolic revision.
2. As a coordinator, I want the chunk source digest computed from the tip tree without the record file, so that it matches the checkpoint.
3. As a coordinator, I want the chunk plan digest computed from the plan at the tip, so that I never hash plan bytes by hand.
4. As a coordinator, I want the acceptance rows copied from the plan at the tip, so that the rows match the `Covers:` lines.
5. As a coordinator, I want the verb to create a version 2 record when none exists, so that the first chunk needs no hand-written header.
6. As a coordinator, I want a version 1 plan refused at record creation, so that a new record never takes the legacy shape.
7. As a coordinator, I want a new chunk appended after the recorded chunks, so that the record keeps the chain order.
8. As a coordinator, I want a repeated chunk form to update that chunk and keep its results, so that a tip move loses no evidence.
9. As a coordinator, I want a chunk that the plan does not name refused, so that the record never holds an unmapped chunk.
10. As a coordinator, I want an unresolvable revision refused, so that no record names a missing commit.

Verification results:

11. As an author, I want `bench record verification` to put my result in the chunk's verification list, so that it never lands in the review list.
12. As an orchestrator, I want `--final` to put my result in the completion list, so that final verification uses the same writer.
13. As an author, I want the source to default to the chunk tip, so that a chunk result names the graded source.
14. As an author, I want the command copied from the plan, so that a retyped command never differs from the planned one.
15. As an author, I want a requirement that the plan does not name refused, so that the record holds only graded obligations.
16. As an author, I want the role taken from the plan, so that an author result and an integration result never swap roles.
17. As an author, I want the outcome derived from the exit code, so that a failed command never reads as a pass.
18. As an author, I want a planned probe written with the mutation from the plan, so that the probe matches its obligation.
19. As an author, I want the probe flags required for a planned probe and refused otherwise, so that a probe entry always has a plan.
20. As an author, I want a result for an unrecorded chunk refused with the `bench record chunk` route, so that I know the first step.

Review results:

21. As a coordinator, I want the review form to copy the chunk's frozen pair and digest, so that no result names a stale pair.
22. As a coordinator, I want the supersession link derived from the previous result of the same axis, so that a later round keeps the chain.
23. As a coordinator, I want the outcome derived from the finding IDs, so that a result with findings never reads as a pass.
24. As a coordinator, I want the role written as independent review, so that a review result never reads as author verification.

Native results:

25. As an author, I want the verb to embed the excerpt file bytes exactly, so that the excerpt is the text I received.
26. As an author, I want the excerpt digest computed over the embedded bytes, so that the digest never hashes another text.
27. As an author, I want an absent, linked, special, empty, oversized, or malformed excerpt file refused, so that no broken excerpt enters the record.

Plan amendments:

28. As a coordinator, I want `bench record amendment` to compute the old and the new plan digests, so that I never hash a plan by hand.
29. As a coordinator, I want each recorded chunk mapped to itself by default, so that an ordinary plan commit needs no mapping input.
30. As a coordinator, I want `--map` for a split or renamed chunk, so that a chunk change keeps its review chain.
31. As a coordinator, I want an unchanged plan refused, so that the record holds no empty amendment.
32. As a coordinator, I want a recorded chunk that the new plan lacks refused without a map, so that no chain breaks silently.

Safe writes:

33. As a writer, I want the verb to parse the whole record before it writes, so that an invalid record never reaches the tree.
34. As a writer, I want a refusal to leave the record directory unchanged, so that a failed call needs no cleanup.
35. As a writer, I want the prose outside the record fence kept byte for byte, so that the review pickup sections survive each write.
36. As a writer, I want the verb to change only `reviews/<slug>.md`, so that a record call never dirties another path.
37. As a writer, I want a linked, special, or empty record file refused, so that the verb never writes through a link.
38. As a writer, I want an unreadable existing record refused, so that the verb never rewrites a record that it cannot parse.
39. As a writer, I want the primary checkout refused, so that `main` receives a record only through a landing.
40. As a writer, I want a control character in a single-line flag value refused, so that no identity field splits a line.
41. As a writer, I want a duplicate evidence ID refused with that ID named, so that a repeated call never records one result twice.

One encoder:

42. As a maintainer, I want the test fixtures to render the record through the production renderer, so that the record fence has one encoder.
43. As a maintainer, I want the reader and the renderer to share one fence name and one fence locator, so that the two never disagree.

Command surface:

44. As an agent, I want `bench help` and each help spelling of `bench record` to show each form, so that I find the grammar without a trial.
45. As an agent, I want each grammar error refused at exit 2 with the usage line, so that a mistyped call writes nothing.
46. As an agent, I want one TOON row with the written entry and its digests, so that I can cite them in the prose.

Guidance:

47. As a coordinator, I want step 6 to route each completed entry to `bench record`, so that I build no record JSON.
48. As an author, I want the Land paragraph to name the chunk and verification forms, so that I record results with the verb.

Reviewed exclusions:

49. As the reviewer, I want the completion state, performer, and reconciliation kept out, so that FT317 decides capability-blocked rows first.
50. As the reviewer, I want `bench probe --names` and `bench test --list` kept out, so that this spec stays one capability.
51. As the reviewer, I want the checkpoint messages unchanged, so that the existing routes and tests keep their bytes.

Added at the review fold:

52. As an orchestrator, I want `--final` to digest the `--source` commit, so that a final result names the final source.
53. As an author, I want a `--source` whose digest differs from the chunk's refused, so that no chunk result names a stale source.

## Implementation decisions

### Decisions closed on 2026-10-01

The reviewer delegated these decisions to the coordinator, and the coordinator closed them:

1. The verb has four forms: `chunk`, `verification`, `review`, and `amendment`. The completion state, the completion performer, the reconciliation map, `bench probe --names`, and `bench test --list` stay out.
2. The spec carries no `Roadmap:` line. FT318 stays open for the completion form, and the next drain narrows it and records the name choice.
3. The grammar is `bench record <chunk|verification|review|amendment> <slug> …`.
4. `--excerpt <file>` supplies the excerpt as exact bytes. The file is regular, not a symlink, not empty, and valid UTF-8. The verb embeds the bytes and hashes them. `--ref <text>` names the native result.
5. The verification form accepts only a requirement that the plan at the source commit names for that chunk or for final verification. The verb copies the command from the plan. No `--command` flag exists.
6. Step 6 of the review phase and the Land paragraph of the implement phase name `bench record`. The step 6 text replaces "Preflight supplies the source and plan digests." The implement phase file does not grow.
7. The verb writes only `reviews/<slug>.md` in the current worktree and refuses the primary checkout. The caller commits the file with `bench commit`.

The coordinator also accepted nine author calls:

- The verb creates only a version 2 record. It also appends to an existing version 1 record.
- The plan supplies the plan digest, the acceptance rows, and the role. `reviewrecord.SourceDigest` supplies each source digest. The chunk entry supplies the review base, tip, and source digest. The previous result of the same axis supplies `supersedes`.
- The verb writes only `completed` results.
- The caller supplies `--id`.
- The caller states the performer. The checkpoint grades the owed performer, and the verb does not.
- A planned probe takes `--probe-outcome`, `--probe-exit-code`, and `--probe-restore`. The plan supplies the mutation.
- The amendment maps each recorded chunk to itself by default, and a repeatable `--map` handles a split or a rename.
- The verb renders the record, parses it again, and then replaces the file.
- The checkpoint messages stay unchanged.

### Decisions closed at the review fold, 2026-10-01

An independent review approved the draft with 11 findings and no blocker. The
coordinator closed each finding:

1. Step 6 names `bench record` for each completed chunk, verification, review, and amendment entry. The completion entry and each result in the failed, skipped, or pending state stay hand-written.
2. The implement phase file holds 80 lines by the budget count, against a budget of 81. A review-owned row holds the count at 80.
3. The verb checks an `--id` against the record before it renders, and its refusal names the ID. `parse.go` keeps its message.
4. `parse.go` owns the one fence locator, and `fenced` and `Render` both use it. A duplicate fence has its own row.
5. A gate row proves that `Save` renders through `Render`. The `recordFence` helper stays review-owned.
6. The control-character row grades `--performer`, `--model`, `--effort`, and `--id`, because `Parse` grades none of their characters.
7. The form changes live in package `reviewrecord`, so `verifier`, `mappedIDs`, and `findChunk` stay private.
8. The RE5 why-clause names the escape in words.
9. Forbid rows retire the two hand-recording sentences of step 6.
10. The header states the FT318 answer.
11. Rows cover the oversized excerpt, the `-h` and `help` spellings, and a failed temporary write.

The coordinator also folded two pieces of advice. `--source` defaults to the
chunk entry's tip, and a chunk result whose source digest differs from the
chunk's refuses. The implement phase states when each entry is written.

### Decisions closed at the confirming pass, 2026-10-01

A confirming review approved the fold with two findings and no blocker. The
reviewer pre-approved the spec and the ticket graph, and the coordinator closed
both findings:

1. The orchestrator records the chunk entry and its tip when it freezes the chunk after the last ticket. Each ticket author then records its verification at that chunk source. A tip move after each ticket would stale the earlier entries of a two-ticket chunk.
2. The RE7 search covers the production files and `internal/reviewrecord/recordtest`. A `_test.go` file that builds a malformed document on purpose is exempt by rule.

### The forms and their flags

Each form takes one `<slug>` operand. The spec path is `specs/<slug>/spec.md`,
and `reviewrecord.Slug` grades it. These are the four grammars:

- `bench record chunk <slug> --chunk <id> --base <commit> --tip <commit>`
- `bench record verification <slug> (--chunk <id> [--source <commit>] | --final --source <commit>) --requirement <id> --id <id> --performer <session> --model <model> --effort <effort> --exit-code <n> --ref <ref> --excerpt <file> [--probe-outcome <verdict> --probe-exit-code <n> --probe-restore pass|fail]`
- `bench record review <slug> --chunk <id> --axis Standards|Spec|Coverage --id <id> --performer <session> --model <model> --effort <effort> --ref <ref> --excerpt <file> [--finding <id>]...`
- `bench record amendment <slug> --source <commit> [--map <old>=<new>[,<new>...]]...`

`bench record --help`, `-h`, and `help` print one `usage:` line for each
implemented form at exit 0. A bare `bench record` prints the same lines at exit 2.
Each form parses its flags through `usage.Parse`.

### The derived fields

The verb resolves each revision through `git rev-parse --verify <rev>^{commit}`
and stores the full object ID.

- **Chunk:** the base and the tip are full IDs. `reviewrecord.SourceDigest` of the tip tree gives the source digest. `reviewrecord.ReadPlan` at the tip tree gives the plan digest and the planned chunk. The planned chunk's `Rows` give the acceptance rows.
- **Verification:** the source commit is `--source`, or the chunk entry's tip when `--chunk` has no `--source`. `reviewrecord.SourceDigest` of the source tree gives the source digest. A chunk result whose source digest differs from the chunk entry's refuses and names `bench record chunk`. The plan at the source tree gives the requirement, its command, and its probe mutation. The plan's `verifier` rule gives the role. The state is `completed`, and the outcome is `pass` for exit code 0 and `fail` otherwise.
- **Review:** the chunk entry gives the base, the tip, and the source digest. The role is `independent-review`. The state is `completed`. The outcome is `pass` with no `--finding`, and `fail` otherwise. The finding IDs keep argv order. `supersedes` holds the ID of the last result of the same axis in that chunk, or no ID.
- **Native result:** the excerpt is the exact file bytes, and `reviewrecord.Digest` of those bytes is the digest. A probe entry carries the same native result as its verification entry.
- **Amendment:** `from` is the record's current `plan_digest`. `to` is the `ReadPlan` digest at the `--source` tree. The record's `plan_digest` then becomes `to`.

The amendment keys are the recorded chunk IDs under `from`. Each recorded chunk
maps through the earlier amendments, through the same rule that the checkpoint
uses, to its ID under `from`. Each key maps to itself unless a `--map` names it.
A key without a map must be a chunk of the new plan. Each `--map` key must be a
recorded key, and each target must be a chunk of the new plan.

A new record has version 2, the spec path, the tip plan digest, an empty
implementation session, one chunk, and a `pending` completion. The verb refuses
to create a record from a version 1 plan. A chunk form for a recorded chunk
replaces its base, tip, digests, and rows, and it keeps its results. A chunk
form for a new chunk appends it after the recorded chunks.

### The write transaction

One write function in `internal/reviewrecord` owns the transaction. It reads the
record file through the existing reader, which refuses a linked component, a
special file, and an empty file. It parses the existing record through the
production parser, applies the form's change, and renders the document. It then
parses the rendered bytes through the same reader path. Only then does it write
a temporary file in `reviews/` and rename it over the record.

Only the chunk form may start from an absent record file. The function then
creates the `reviews/` directory when it is absent. A refusal at any step leaves
the record bytes unchanged and leaves no temporary file. A failed temporary write
is such a refusal.

Before it renders, a verification or review change checks its `--id` against
every evidence ID in the record, and the refusal names that ID. `Parse` stays the
final guard, and its message stays unchanged.

### The renderer

`reviewrecord.Render` takes the current document bytes and a record, and it
returns the new document bytes. It encodes the record with two-space indentation
and without HTML escaping. It does not validate the record, because the test
fixtures render invalid records on purpose. The write transaction validates.

`parse.go` owns the one fence locator. `fenced` returns the payload through it,
and `Render` replaces the payload lines that it finds. No second fence scanner
exists.

- A document with one fence keeps every byte outside the payload lines.
- A document with no fence keeps its bytes. The renderer adds a newline when the document does not end with one, then a blank line and the fence.
- A nil document becomes `# Review outcomes`, a blank line, and the fence.
- A document with an unterminated or a duplicate fence is an error.

One constant names the fence for the reader and the renderer.
`recordtest.Fixture.Save` and the `recordFence` helper of the preflight tests
render through `Render`.

### The excerpt file

The verb reads `--excerpt` through `bounds.ClassifyNoFollow`. Every state other
than `parsed` refuses with the `toon.RecordError` line for the operand path. So an
absent file, a symlink, a special file, an empty file, invalid UTF-8, and a file
above `bounds.ControlRecordLimit` each refuse at exit 1. A relative path resolves
against the working directory of the process.

### Refusals and exit codes

A grammar error exits 2 with the usage line of the form. These are grammar errors:

- an unknown form, a missing operand, or a missing required flag
- both `--chunk` and `--final`, or neither
- `--final` without `--source`
- one or two of the three probe flags without the third
- an `--exit-code` or a `--probe-exit-code` that is not a base-10 integer
- an `--axis` outside `Standards`, `Spec`, and `Coverage`
- a `--map` value with no `=`, an empty side, or a key that an earlier `--map` names

A content refusal exits 1. An unreadable excerpt prints the `toon.RecordError`
line, and each other content refusal after step 2 prints
`error: bench record <form> refused — <cause>`. The checks run in this order,
and the first refusal wins:

1. An empty root prints the `toon.NotInRepo` line.
2. The primary checkout prints the `usage.PrimaryCheckoutRefusal` line.
3. A control character in the value of a single-line flag names that flag. The single-line flags are `--chunk`, `--id`, `--performer`, `--model`, `--effort`, `--ref`, `--requirement`, `--finding`, `--probe-outcome`, `--map`, and `--excerpt`.
4. An invalid slug, an unreadable excerpt, an unreadable record, an unresolvable revision, or an unreadable plan refuses with its cause.
5. A missing chunk entry or a missing record names `bench record chunk`.
6. A form rule refuses with its cause. Examples are an unplanned requirement, a source mismatch, a recorded `--id`, a probe mismatch, and an unmapped chunk.
7. A rendered record that fails the reader refuses with the parser message.
8. A failed temporary write or rename refuses with its cause.

### The output

The dispatcher prints its tree row first, because the verb is tree-scoped. Each
form then prints one TOON table with one row and no help envelope:

- `chunk[1]{id,action,base,tip,source_digest,plan_digest,rows}`. The `action` cell is `created`, `added`, or `updated`. The `rows` cell is the count of acceptance rows.
- `verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}`. The `list` cell is `chunk` or `completion`, and the `chunk` cell is empty for `completion`.
- `review[1]{chunk,id,axis,outcome,supersedes,source_digest,excerpt_digest}`. The `supersedes` cell is empty for a first result.
- `amendment[1]{from,to,chunks}`. The `chunks` cell is the count of amendment keys.

### Registry and structure budgets

The registry row is `record`, with the mutation AXI exemption, the tree scope,
and the bounded response. Its help rows read `recordcmd.HelpRows`, as the
preflight rows read `evidencecmd.HelpRows`, at order 40. The subcommand-routing
table records `record` as routed through `internal/reviewrecord/recordcmd`.

`cmd/bench/main.go` is 446 lines against a budget of 400, and the growth lane
reds a line that it gains. Ticket 2 moves the inline adapter of the `assessment`
row into `cmd/bench/command_registry.go`, so that row takes one line. The
`record` row then keeps `main.go` at or below 446 lines.

`internal/reviewrecord/` holds 11 source files against a directory budget of 12.
The build adds one file, `write.go`, and puts its tests in `record_test.go` and
`source_test.go`. Each form change lives in package `reviewrecord`, so `verifier`,
`mappedIDs`, and `findChunk` stay private:

- `write.go` holds `Render`, the transaction, and the chunk, verification, and review changes. It stays at or below 400 lines.
- The amendment change goes in `coverage.go`, beside the `mappedIDs` rule that it reuses.
- The new package `internal/reviewrecord/recordcmd` parses the flags, reads the excerpt, and prints the output.

### The guidance

Ticket 6 replaces lines 208 to 214 of step 6 in
`.agents/commands/bench-review-implementation.md`. No anchor names a sentence on
those lines. The new step 6 text reads:

    Retain every terminal return in one fenced `bench-review-record` JSON payload.
    The `internal/reviewrecord` types own the schema.

    Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.
    `bench record chunk` writes the frozen pair and its digests, and `bench record review` appends one axis result.
    `bench record verification` appends one planned result, and `bench record amendment` records a plan-digest change.
    Supply the performer, model, effort, and native result. The verb derives the role, the source, and each digest.

    The completion entry and each result in the failed, skipped, or pending state stay hand-written, because `bench record` has no form for them.
    A hand-written entry names its performer, role, model, effort, source digest, state, and native result.
    It embeds the minimal native excerpt with the SHA-256 digest of that excerpt.

    Use explicit `unknown` for unavailable model or effort metadata.
    Pass the minimal native excerpt to `bench record` in a file, and the verb computes its digest.
    Local logs are supplemental evidence.
    Keep author verification separate from independent review.

Ticket 6 inserts two sentences after the first sentence of the Land paragraph in
`.agents/commands/bench-implement-spec.md`, on the same physical line 56. No
anchor names a sentence on that line. The new sentences read:

    When the orchestrator freezes a chunk after its last ticket, it records the chunk entry with `bench record chunk`.
    Each ticket author then writes its verification entries at that chunk source with `bench record verification`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RE-C1 / `1-render-record-fence.md` | One renderer and one fence locator write the record fence, and the shared fixtures use them. | RE1, RE2, RE3, RE4, RE5, RE6, RE7, RE8, RE9, RE10, RE102, RE103 | `bench test --package ./internal/reviewrecord/...`, `bench test --package ./internal/preflight` | no |
| RE-C2 / `2-record-chunk-entry.md` | `bench record chunk` writes a chunk entry through the safe write transaction. | RE11, RE12, RE13, RE14, RE15, RE16, RE17, RE18, RE19, RE20, RE21, RE22, RE23, RE24, RE25, RE26, RE27, RE28, RE29, RE30, RE31, RE32, RE33, RE34, RE35, RE37, RE38, RE39, RE40, RE41, RE105, RE106 | `bench test --package ./internal/reviewrecord/...`, `bench test --package ./cmd/bench`, `bench test --package ./internal/conformance` | yes |
| RE-C3 / `3-record-verification-result.md`, `4-record-review-result.md` | Authors and coordinators append verification and review results with the verb. | RE36, RE42, RE43, RE44, RE45, RE46, RE47, RE48, RE49, RE50, RE51, RE52, RE53, RE54, RE55, RE56, RE57, RE58, RE59, RE60, RE61, RE62, RE63, RE64, RE65, RE66, RE67, RE68, RE69, RE70, RE71, RE104, RE107, RE108, RE72, RE73, RE74, RE75, RE76, RE77, RE78, RE79, RE80, RE81, RE82, RE83 | `bench test --package ./internal/reviewrecord/...`, `bench test --package ./cmd/bench` | no |
| RE-C4 / `5-record-plan-amendment.md`, `6-name-bench-record-in-guidance.md` | The verb records plan amendments, and the phase guidance routes each completed entry to the verb. | RE84, RE85, RE86, RE87, RE88, RE89, RE90, RE91, RE92, RE93, RE94, RE95, RE96, RE97, RE98, RE99, RE100, RE101, RE109, RE110, RE111 | `bench test --package ./internal/reviewrecord/...`, `bench test --package ./cmd/bench`, `bench test --package ./internal/anchors`, `bench test --check guidance-prose-budgets` | no |

Ticket 1 creates the renderer that ticket 2 consumes, so RE-C1 stays small and its
review closes first. Ticket 2 creates the transaction and the command family that
tickets 3, 4, and 5 consume, so RE-C2 is its own chunk too. Ticket 4 reuses the
private excerpt reader of ticket 3 inside RE-C3.

## Testing decisions

- A good test calls `recordcmd.Command` with the root of a linked worktree and with argv. It compares the exit code and the output. It then reads the record back through `reviewrecord.Read`, and it grades the result with `reviewrecord.Check` where a row names the checkpoint.
- The seam is the new `recordcmd.Command`, and its oracle is the existing checkpoint reader. The prior art is `TestReviewRecordSource` in `internal/reviewrecord/source_test.go`, which builds a fixture with `recordtest` and grades it with `CheckSource`. It was read in this session.
- Ticket 2 adds `recordtest.NewLinked`. It prepares a `recordtest` fixture in a linked worktree of a new repository, so the verb sees a worktree and the primary checkout stays testable.
- `Render` has its own tests in `internal/reviewrecord/record_test.go`. The `Save` row lives in `internal/reviewrecord/source_test.go`, because an internal test of the package cannot import `recordtest`.
- A test derives each expected TOON row through `toon.Table`, because a digest that starts with `0` and a digit renders quoted.
- The FIFO rows take `capability.Fifo`. The read-only row takes `capability.Privilege` under root, as `internal/probe/omit_file_test.go` does.
- The package tests run in the gate's `test` phase, so that phase observes each row. The anchor rows run in the `internal/anchors` package tests in the same phase.

### Seam diagram

    trigger: an author, a coordinator, or an orchestrator runs `bench record <form> <slug> …`
        │
        ▼
    argv, worktree root, commits, plan, excerpt file  ──▶  [ recordcmd.Command → reviewrecord write transaction ]  ──▶  reviews/<slug>.md, one TOON row, exit code
                      ◀ tests attach here: `Command` over a linked recordtest worktree, then
                        `reviewrecord.Read` and `reviewrecord.Check` read the written record

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RE1 | 35 | `Render` over a document with prose before and after one fence changes only the payload lines | planned TestRenderKeepsProseAroundTheFence in internal/reviewrecord, through `Render` | A renderer that writes a new document drops the prose after the fence. |
| RE2 | 35 | `Render` over a document with no fence and no final newline adds a newline, a blank line, and the fence after the unchanged prose | planned TestRenderAppendsAFenceAfterUnterminatedProse in internal/reviewrecord, through `Render` | A fence glued to the last prose line is not a fence line, so `Read` reports a missing fence. |
| RE3 | 5 | `Render` of a nil document returns `# Review outcomes`, a blank line, and the fence | planned TestRenderStartsANewDocument in internal/reviewrecord, through `Render` | A renderer that returns only the fence changes the bytes that every fixture consumer reads today. |
| RE4 | 25 | An excerpt with a tab, a return, U+001B, a newline, and `<&>` reads back byte-exact through `Read` after `Render` | planned TestRenderRoundTripsAHostileExcerpt in internal/reviewrecord, through `Render` and `Read` | A renderer that writes raw control bytes splits a payload line or fails the bounded read. |
| RE5 | 35 | The rendered payload has two-space indentation and holds `&` unescaped | planned TestRenderWritesAReadablePayload in internal/reviewrecord, through `Render` | The default JSON encoder escapes `&` as a six-character Unicode sequence. |
| RE6 | 38 | `Render` over a document with an unterminated fence returns an error | planned TestRenderRefusesAnUnterminatedFence in internal/reviewrecord, through `Render` | A renderer that appends a fence leaves the old opening line, and the reader then refuses the document. |
| RE102 | 38 | `Render` over a document with two record fences returns an error | planned TestRenderRefusesADuplicateFence in internal/reviewrecord, through `Render` | A renderer that replaces the first fence leaves the second, and the reader then refuses the document. |
| RE7 | 42 | A search for the fence-opening text `bench-review-record\n` in the non-test Go files, `internal/reviewrecord/recordtest` included, finds no hit | review-owned: the search command in Further notes | A fixture or a production writer that keeps its own encoder adds a hit. |
| RE8 | 43 | A search for the text `"bench-review-record"` in non-test Go code finds one constant declaration | review-owned: the search command in Further notes | A renderer with its own literal adds a second hit. |
| RE103 | 42 | `recordtest.Fixture.Save` of a record whose excerpt holds `&` writes a raw `&` in `reviews/example.md` | planned TestFixtureSaveRendersThroughRender in internal/reviewrecord/source_test.go, through `Save` | A `Save` that keeps `json.MarshalIndent` writes the Unicode escape for `&`. |
| RE9 | 42 | `TestReviewRecordSource` passes with `recordtest.Fixture.Save` rendering through `Render` | `internal/reviewrecord/source_test.go` (`TestReviewRecordSource`), with no assertion changed | A renderer whose payload the reader cannot parse reds each fixture consumer. |
| RE10 | 42 | `TestDelegatedEvidenceProjection` passes with the preflight `recordFence` helper rendering through `Render` | `internal/preflight/delegated_evidence_test.go` (`TestDelegatedEvidenceProjection`), with no assertion changed | The helper renders a version 9 record, so a renderer that validates refuses that case. |
| RE11 | 1 | `bench record chunk example --chunk 1 --base <base> --tip HEAD` writes the 40-hex IDs that `git rev-parse` gives for both revisions | planned TestRecordChunkWritesFullCommitIDs in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that stores the operand stores `HEAD`, and `Parse` refuses it. |
| RE12 | 2 | Over a tip whose tree holds `reviews/example.md`, the chunk `source_digest` equals `reviewrecord.SourceDigest` of the tip tree and differs from the tip tree ID | planned TestRecordChunkSourceDigestExcludesTheRecord in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A digest of the whole tip tree equals the tree ID, which is the worktree-verb-runner error. |
| RE13 | 3 | The chunk `plan_digest` equals the `Digest` field that `reviewrecord.ReadPlan` returns at the tip tree | planned TestRecordChunkPlanDigestReadsThePlan in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A digest of the spec bytes alone omits the ticket bytes. |
| RE14 | 4 | The chunk `acceptance_rows` hold `E1`, the `Covers:` row of ticket `1.md` | planned TestRecordChunkCopiesAcceptanceRows in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that leaves the rows empty fails the match. |
| RE15 | 5 | With no `reviews/` directory and a version 2 plan, `reviewrecord.Read` of the new file returns version 2, the spec path, the tip plan digest, an empty implementation session, and one chunk | planned TestRecordChunkCreatesAVersionTwoRecord in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | No writer exists today, so the call fails as an unknown verb. |
| RE16 | 5 | With an empty `reviews/` directory, the chunk form creates `reviews/example.md` | planned TestRecordChunkCreatesInAnEmptyDirectory in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A create path that calls `os.Mkdir` fails when the directory exists. |
| RE17 | 6 | With a version 1 plan and no record, the chunk form exits 1, the output names `version 2`, and no record file exists | planned TestRecordChunkRefusesAVersionOnePlan in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that copies the plan version writes a version 1 record that `Parse` refuses with another message. |
| RE18 | 7 | A chunk form for chunk `2` after chunk `1` leaves chunk `1` first and chunk `2` second | planned TestRecordChunkAppendsAfterTheRecordedChunk in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that prepends the chunk breaks the chain order that the checkpoint reads. |
| RE19 | 8 | A chunk form for recorded chunk `1` with a new tip keeps its two verification results and three review results | planned TestRecordChunkUpdateKeepsResults in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that replaces the whole chunk entry drops the results. |
| RE20 | 8 | Two runs of one chunk form with the same operands exit 0 and leave byte-identical record files | planned TestRecordChunkRepeatIsByteIdentical in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that appends on each run adds a second chunk `1`, which `Parse` refuses. |
| RE21 | 9 | A chunk form for chunk `9`, which the plan does not name, exits 1, names chunk `9`, and leaves the record bytes unchanged | planned TestRecordChunkRefusesAnUnplannedChunk in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that skips the plan lookup writes a chunk that the checkpoint cannot map. |
| RE22 | 10 | A `--tip` of `no-such-rev` exits 1 and leaves the record bytes unchanged | planned TestRecordChunkRefusesAnUnknownRevision in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A resolver error that falls through to an empty ID writes an invalid chunk. |
| RE23 | 46 | A chunk form that creates the record prints `chunk[1]{id,action,base,tip,source_digest,plan_digest,rows}:` with the action `created` | planned TestRecordChunkReportsCreated in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The verb prints nothing today. |
| RE24 | 46 | A chunk form that adds chunk `2` prints the action `added` | planned TestRecordChunkReportsAdded in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A fixed action cell prints `created`. |
| RE25 | 46 | A chunk form for recorded chunk `1` prints the action `updated` | planned TestRecordChunkReportsUpdated in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A fixed action cell prints `created`. |
| RE26 | 35 | A chunk form over a record with prose before and after the fence keeps every prose byte | planned TestRecordKeepsTheProseAroundTheFence in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that writes a new document drops the review pickup sections. |
| RE27 | 36 | After a chunk form, `git status --porcelain` lists only `reviews/example.md` | planned TestRecordChangesOnlyTheRecordPath in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A temporary file left in `reviews/` adds a second path. |
| RE28 | 37 | A dangling symlink at `reviews/example.md` makes the chunk form exit 1, and the link target stays absent | planned TestRecordRefusesADanglingRecordLink in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that reads with `os.ReadFile` sees absence and creates the target through the link. |
| RE29 | 37 | A live symlink at `reviews/example.md` makes the chunk form exit 1 and leaves the link target bytes unchanged | planned TestRecordRefusesALiveRecordLink in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that follows the link rewrites a file outside the tree. |
| RE30 | 37 | A FIFO at `reviews/example.md` makes the chunk form exit 1 before the test deadline | planned TestRecordRefusesASpecialRecordFile in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A plain open blocks on the FIFO. |
| RE31 | 37 | An empty `reviews/example.md` makes the chunk form exit 1 and leaves the file empty | planned TestRecordRefusesAnEmptyRecordFile in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that reads empty as absent writes a new record over the file. |
| RE32 | 38 | A record file with an unterminated fence makes the chunk form exit 1 and leaves its bytes unchanged | planned TestRecordRefusesAnUnterminatedFence in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that appends a fence leaves two fence openings in the file. |
| RE33 | 38 | A record whose payload holds a duplicate JSON key makes the chunk form exit 1 and leaves its bytes unchanged | planned TestRecordRefusesAnInvalidRecord in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that decodes without the production parser keeps one value and drops the other. |
| RE34 | 39 | In the primary checkout the chunk form exits 1 with the `usage.PrimaryCheckoutRefusal` text and writes no file | planned TestRecordRefusesThePrimaryCheckout in internal/reviewrecord/recordcmd, through `Command` over the fixture's primary checkout | A verb without the check writes a record on `main`. |
| RE35 | 39 | In the primary checkout with a `--tip` of `no-such-rev`, the output holds the primary-checkout refusal and no revision refusal | planned TestRecordPrimaryRefusalComesFirst in internal/reviewrecord/recordcmd, through `Command` over the fixture's primary checkout | A verb that resolves revisions first prints the revision refusal. |
| RE106 | 34 | With `reviews/` read-only, the chunk form exits 1, and the record bytes and the directory listing stay unchanged | planned TestRecordFailedTemporaryWriteChangesNothing in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that rewrites the record in place succeeds in a read-only directory and changes the bytes. |
| RE37 | 45 | With an empty root the chunk form exits 1 with the `toon.NotInRepo` text | planned TestRecordRefusesOutsideARepository in internal/reviewrecord/recordcmd, through `Command` | A verb that resolves the record path first reports a file error. |
| RE38 | 45 | Each of `bench record`, `bench record nosuch`, `bench record chunk` with no operand, and `bench record chunk example` with no `--tip` exits 2 with a line that starts `usage: bench record` | planned TestRecordGrammarRefusals in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A parser that ignores a missing flag runs a partial form. |
| RE39 | 44 | `bench record --help` exits 0 and prints the usage line of each implemented form from `recordcmd.HelpRows` | planned TestRecordHelpPrintsEachForm in internal/reviewrecord/recordcmd, through `Command` | A help text kept apart from the form grammars omits a form. |
| RE105 | 44 | `bench record -h` and `bench record help` each print the `--help` lines at exit 0 | planned TestRecordHelpSpellings in internal/reviewrecord/recordcmd, through `Command` | A help check that matches only `--help` sends the other spellings to the unknown-form refusal. |
| RE40 | 44 | The real dispatcher answers `bench record --help` with the `recordcmd` usage text at exit 0 | planned TestRecordRouteAnswersItsUsage in cmd/bench/help_inventory_test.go, through `Command.Run` | A registry row wired to another handler prints other text. |
| RE41 | 44 | `bench help` prints the chunk row that Further notes quotes | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), extended in place | The independent golden expectation lacks the row until the ticket adds it. |
| RE36 | 40 | Each of `--performer`, `--model`, `--effort`, and `--id` with a value that holds U+001B makes the verification form exit 1, name that flag, and leave the record bytes unchanged | planned TestRecordRefusesAControlCharacterInAFlag in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | `Parse` grades none of these characters, so a writer that relies on it writes the byte. |
| RE42 | 11 | A verification result for chunk `1` adds one entry to the chunk's `verification` list and none to its `reviews` list | planned TestRecordVerificationLandsInTheChunkList in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that appends to the review list repeats the worktree-verb-runner error. |
| RE43 | 12 | `--final --requirement acceptance` adds one entry to `completion.verification` | planned TestRecordFinalVerificationLandsInTheCompletionList in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that needs a chunk refuses the final form. |
| RE44 | 13 | With `HEAD` at a later commit that changes a source file, a chunk result with no `--source` takes the source digest of the chunk tip | planned TestRecordVerificationDefaultsToTheChunkTip in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that defaults to `HEAD` names a later source. |
| RE107 | 53 | A chunk result whose `--source` names a commit with another source digest exits 1, names `bench record chunk`, and leaves the record bytes unchanged | planned TestRecordVerificationRefusesAStaleSource in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that records any source writes a result that the checkpoint calls stale. |
| RE108 | 52 | `--final --source` on a commit past the last chunk tip writes the `source_digest` of that commit's tree | planned TestRecordFinalVerificationDigestsTheNamedSource in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that copies the last chunk's digest fails the match. |
| RE45 | 14 | The entry's `command` equals `go test ./...`, the planned command of requirement `tests` | planned TestRecordVerificationCopiesThePlannedCommand in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer with no command source writes an empty command, which `Parse` refuses. |
| RE46 | 15 | `--requirement nosuch` exits 1, names `tests` and `additional`, and leaves the record bytes unchanged | planned TestRecordVerificationRefusesAnUnplannedRequirement in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that skips the plan lookup writes an obligation that the checkpoint never grades. |
| RE47 | 16 | A chunk entry in a version 2 record has the role `author-verification` | planned TestRecordVerificationWritesTheAuthorRole in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | One fixed integration role fails this list. |
| RE48 | 16 | A `--final` entry in a version 2 record has the role `integration-verification` | planned TestRecordFinalVerificationWritesTheIntegrationRole in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | One fixed author role fails the completion list. |
| RE49 | 16 | When the plan gives requirement `tests` to a ticket with an empty assignment history, the form exits 1 and names `undispatched ticket` | planned TestRecordVerificationRefusesAnUndispatchedTicket in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that skips the `verifier` rule records an obligation that no author owns. |
| RE50 | 17 | `--exit-code 0` writes the outcome `pass` and `exit_code` 0 | planned TestRecordVerificationPassesOnExitZero in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that leaves the outcome empty writes an entry that `Parse` refuses. |
| RE51 | 17 | `--exit-code 3` writes the outcome `fail` and `exit_code` 3, and the verb exits 0 | planned TestRecordVerificationFailsOnNonzeroExit in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A fixed `pass` outcome hides a failed command. |
| RE52 | 18 | For requirement `tests`, the entry holds a `probe` with the mutation `omit source check`, the given outcome, exit code, and restore, and the entry's native result | planned TestRecordVerificationWritesThePlannedProbe in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that takes the mutation from another source can differ from the plan. |
| RE53 | 19 | For requirement `tests`, a call without the probe flags exits 1, names `--probe-outcome`, and leaves the record bytes unchanged | planned TestRecordVerificationRequiresThePlannedProbe in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that skips the probe writes an entry that the checkpoint refuses later. |
| RE54 | 19 | For requirement `additional`, which plans no probe, a call with the probe flags exits 1 and leaves the record bytes unchanged | planned TestRecordVerificationRefusesAnUnplannedProbe in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that keeps the flags writes a probe that the plan does not name. |
| RE55 | 20 | A result for chunk `2` when the record holds only chunk `1` exits 1 and names `bench record chunk` | planned TestRecordVerificationNamesTheChunkForm in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that creates the chunk writes it with no frozen pair. |
| RE56 | 20 | `--final` with no record file exits 1, names `bench record chunk`, and creates no file | planned TestRecordFinalVerificationNeedsARecord in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that creates a record here writes one with no chunk, which the checkpoint refuses. |
| RE57 | 25 | The entry's `native_ref.excerpt` equals the excerpt file bytes for a file with a tab, a return, U+001B, and a newline | planned TestRecordVerificationEmbedsTheExcerptBytes in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that trims or escapes the text changes the bytes. |
| RE58 | 26 | The entry's `native_ref.digest` equals `reviewrecord.Digest` of the excerpt file bytes | planned TestRecordVerificationDigestsTheExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A digest over trimmed text or over the path fails the match. |
| RE59 | 27 | An absent `--excerpt` file exits 1 with `is absent` and leaves the record bytes unchanged | planned TestRecordRefusesAnAbsentExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A reader that treats absence as empty text writes an entry that `Parse` refuses with another message. |
| RE60 | 27 | A symlink to a regular file at `--excerpt` exits 1 with `is wrong-type` | planned TestRecordRefusesALinkedExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A reader that follows the link embeds bytes that the caller did not name. |
| RE61 | 27 | A FIFO at `--excerpt` exits 1 with `is wrong-type` before the test deadline | planned TestRecordRefusesASpecialExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A plain open blocks on the FIFO. |
| RE62 | 27 | An empty `--excerpt` file exits 1 with `is empty` | planned TestRecordRefusesAnEmptyExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A reader that skips the state check reaches the `Parse` refusal, which names no file. |
| RE63 | 27 | An `--excerpt` file that holds byte 0xFF exits 1 with `is malformed` | planned TestRecordRefusesAMalformedExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The JSON encoder replaces the byte with U+FFFD, so the embedded text differs from the file. |
| RE104 | 27 | An `--excerpt` file one byte above `bounds.ControlRecordLimit` exits 1 with `is unreadable` | planned TestRecordRefusesAnOversizedExcerpt in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A reader with no bound embeds the bytes, and the record then fails its own bounded read. |
| RE64 | 25 | An `--excerpt` path that holds a space and `[glob]*` reads as written | planned TestRecordReadsAnExcerptPathWithGlobCharacters in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A reader that expands the path finds no file. |
| RE65 | 41 | A second result with an `--id` that the record holds exits 1 before the render, names that ID, and leaves the record bytes unchanged | planned TestRecordVerificationRefusesADuplicateID in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The `Parse` refusal names no ID, so a writer that relies on it fails the name match. |
| RE66 | 34 | After the parser refuses a `--ref` that holds `../`, `reviews/` holds only the unchanged `example.md` | planned TestRecordRefusalLeavesNoTemporaryFile in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that writes before it parses leaves its temporary file or a changed record. |
| RE67 | 46 | A chunk result prints `verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}:` with the list `chunk` | planned TestRecordVerificationReportsTheChunkList in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The verb prints nothing for this form until the ticket adds it. |
| RE68 | 46 | A `--final` result prints the list `completion` and an empty `chunk` cell | planned TestRecordFinalVerificationReportsTheCompletionList in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A fixed list cell prints `chunk`. |
| RE69 | 11 | After the chunk form, one verification result for each planned requirement of chunk `1`, and three fixture review results, the committed record passes `reviewrecord.Check` for chunk `1` | planned TestRecordedVerificationPassesTheCheckpoint in internal/reviewrecord/recordcmd, through `Command` and then `reviewrecord.Check` | The checkpoint is an independent reader, so a wrong derived field reds it. |
| RE70 | 45 | Each of `--chunk 1 --final`, no list flag, `--final` without `--source`, one probe flag alone, and `--exit-code x` exits 2 with a line that starts `usage: bench record verification` | planned TestRecordVerificationGrammarRefusals in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A parser that accepts a partial form writes an incomplete entry. |
| RE71 | 44 | `bench help` prints the verification row that Further notes quotes | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), extended in place | The golden expectation lacks the row until the ticket adds it. |
| RE72 | 21 | A review entry's `base`, `tip`, and `source_digest` equal those of the chunk entry | planned TestRecordReviewCopiesTheChunkPair in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that resolves `HEAD` names a later tip, and the checkpoint calls the result stale. |
| RE73 | 22 | The first Standards result has an empty `supersedes` list | planned TestRecordReviewFirstResultSupersedesNothing in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that names the last result of any axis names an earlier axis, and `Parse` refuses it. |
| RE74 | 22 | A second Standards result written after a Spec result supersedes the first Standards ID | planned TestRecordReviewSupersedesTheSameAxis in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that names the last result of any axis names the Spec ID. |
| RE75 | 23 | A result with no `--finding` has the outcome `pass` and an empty `finding_ids` list | planned TestRecordReviewWithoutFindingsPasses in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A fixed `fail` outcome with no finding is refused by `Parse`. |
| RE76 | 23 | `--finding R1 --finding R2` writes the outcome `fail` and the `finding_ids` `R1` and `R2` in that order | planned TestRecordReviewWithFindingsFails in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that sorts or drops an ID fails the match. |
| RE77 | 24 | The entry has the role `independent-review` and lands in the chunk's `reviews` list | planned TestRecordReviewLandsInTheReviewList in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that appends to the verification list fails the match. |
| RE78 | 20 | A review result for an unrecorded chunk exits 1 and names `bench record chunk` | planned TestRecordReviewNamesTheChunkForm in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that creates the chunk writes it with no frozen pair. |
| RE79 | 33 | In a version 1 record, a result whose `--performer` is the implementation session exits 1 with `invalid independent review axis or performer` and leaves the record bytes unchanged | planned TestRecordReviewRefusesTheImplementationSession in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that skips the parse before the write records a self-review. |
| RE80 | 46 | The review form prints `review[1]{chunk,id,axis,outcome,supersedes,source_digest,excerpt_digest}:` | planned TestRecordReviewReportsItsRow in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The verb prints nothing for this form until the ticket adds it. |
| RE81 | 21 | After the chunk form, two verification results, and three review results, all from the verb and committed, the record passes `reviewrecord.Check` for chunk `1` | planned TestRecordedChunkPassesTheCheckpoint in internal/reviewrecord/recordcmd, through `Command` and then `reviewrecord.Check` | The checkpoint is an independent reader, so a wrong review field reds it. |
| RE82 | 45 | Each of `--axis Security`, a missing `--axis`, and a missing `--excerpt` exits 2 with a line that starts `usage: bench record review` | planned TestRecordReviewGrammarRefusals in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A parser that accepts any axis writes an entry that `Parse` refuses. |
| RE83 | 44 | `bench help` prints the review row that Further notes quotes | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), extended in place | The golden expectation lacks the row until the ticket adds it. |
| RE84 | 28 | After a plan commit, `bench record amendment example --source HEAD` appends an amendment whose `from` is the prior record `plan_digest` and whose `to` is the `ReadPlan` digest at `HEAD` | planned TestRecordAmendmentComputesBothDigests in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer with a digest of the spec bytes alone fails the `to` match. |
| RE85 | 28 | After the amendment, the record `plan_digest` equals the amendment's `to` | planned TestRecordAmendmentMovesThePlanDigest in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that leaves the digest makes the checkpoint refuse `stale plan digest`. |
| RE86 | 29 | The amendment maps each recorded chunk ID to a list that holds only that ID | planned TestRecordAmendmentMapsEachChunkToItself in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that writes no mapping fails the checkpoint's mapping rule. |
| RE87 | 29 | A second amendment's `from` equals the first amendment's `to` | planned TestRecordAmendmentChainsFromTheLastDigest in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that reads `from` from a chunk entry names the first plan digest twice. |
| RE88 | 30 | `--map 1=1a,1b` writes the `chunk_ids` entry `1` as `1a` and `1b` for a plan that split chunk `1` | planned TestRecordAmendmentWritesAMappedSplit in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that ignores `--map` refuses the split as an unmapped chunk. |
| RE89 | 31 | With an unchanged plan the amendment form exits 1, names `unchanged`, and leaves the record bytes unchanged | planned TestRecordAmendmentRefusesAnUnchangedPlan in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer with no check appends an amendment whose `from` equals its `to`. |
| RE90 | 32 | When the new plan lacks recorded chunk `1` and no `--map` names it, the amendment form exits 1 and names chunk `1` | planned TestRecordAmendmentRefusesAnUnmappedChunk in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that maps by identity names a chunk that the new plan lacks. |
| RE91 | 30 | `--map 9=1`, where chunk `9` is not recorded, exits 1 and names chunk `9` | planned TestRecordAmendmentRefusesAMapForAnUnrecordedChunk in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that copies each map writes a key that the chain never reads. |
| RE92 | 30 | `--map 1=9`, where the new plan has no chunk `9`, exits 1 and names chunk `9` | planned TestRecordAmendmentRefusesAMapToAnUnplannedChunk in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that copies each map writes a target that the checkpoint refuses. |
| RE93 | 28 | With no record file the amendment form exits 1 and names `bench record chunk` | planned TestRecordAmendmentNeedsARecord in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A writer that creates a record here writes one with no chunk. |
| RE94 | 46 | The amendment form prints `amendment[1]{from,to,chunks}:` with the count of amendment keys | planned TestRecordAmendmentReportsItsRow in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | The verb prints nothing for this form until the ticket adds it. |
| RE95 | 29 | After chunk `1`, a plan commit, and chunk `2` work, the amendment and the verb-written chunk `2` evidence pass `reviewrecord.Check` for chunk `2` | planned TestRecordedAmendmentPassesTheCheckpoint in internal/reviewrecord/recordcmd, through `Command` and then `reviewrecord.Check` | Without the amendment the check returns `stale plan digest`, so a wrong amendment reds the row. |
| RE96 | 45 | Each of `--map 1`, `--map =1`, `--map 1=`, and `--map 1=2 --map 1=3` exits 2 with a line that starts `usage: bench record amendment` | planned TestRecordAmendmentGrammarRefusals in internal/reviewrecord/recordcmd, through `Command` over a linked fixture worktree | A parser that splits loosely writes an empty chunk ID. |
| RE97 | 44 | `bench help` prints the amendment row that Further notes quotes | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), extended in place | The golden expectation lacks the row until the ticket adds it. |
| RE98 | 47 | Step 6 of `.agents/commands/bench-review-implementation.md` holds the sentence "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new require-in-step rule | A removal of the sentence reds the anchor. |
| RE99 | 47 | `.agents/commands/bench-review-implementation.md` does not hold the sentence "Preflight supplies the source and plan digests." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new forbid rule | A restored sentence reds the anchor. |
| RE109 | 47 | `.agents/commands/bench-review-implementation.md` does not hold the sentence `Record the performer, role, model, effort, frozen base and tip, source, state, and native result.` | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new forbid rule | A restored hand-recording sentence reds the anchor. |
| RE110 | 47 | `.agents/commands/bench-review-implementation.md` does not hold the sentence `Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.` | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new forbid rule | A restored hand-recording sentence reds the anchor. |
| RE100 | 48 | The Land section of `.agents/commands/bench-implement-spec.md` holds the sentence "Each ticket author then writes its verification entries at that chunk source with `bench record verification`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new require-in-section rule | A removal of the sentence reds the anchor. |
| RE111 | 48 | The Land section of `.agents/commands/bench-implement-spec.md` holds the sentence "When the orchestrator freezes a chunk after its last ticket, it records the chunk entry with `bench record chunk`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`), through a new require-in-section rule | A removal of the sentence reds the anchor. |
| RE101 | 48 | `.agents/commands/bench-implement-spec.md` holds 80 lines by the `proseBudgetLineCount` rule, its count at `0c95c944` | review-owned: a line count of the file at the ticket 6 tip | The budget check reds only at 82 lines, so only a count catches one added line. |

Not covered: story 49 — the verb has no completion-state form, and FT318 keeps that work open until FT317 decides.
Not covered: story 50 — no ticket writes those projections, and the Out of scope section prices them.
Not covered: story 51 — no ticket edits a checkpoint message, and the existing checkpoint tests pin each one.

### Edge inventory

Audience: each behavior serves every repository that links the kit, because a
linked repository runs spec-backed builds with the same checkpoint. A linked
repository can have no `reviews/` directory: RE15 creates it. An empty
`reviews/` directory: RE16. An absent record file: RE15. An empty record file: RE31.

The shell CLI hostile-input checklist, class by class:

- Paths with spaces or glob characters: RE64.
- Control bytes in text that reaches a TOON cell: RE36 refuses the single-line flag values. Each other cell is a computed hex value or a fixed word.
- Control bytes that a sink permits: the excerpt is the one multi-line value. RE4 and RE57 prove a tab, a return, U+001B, and a newline round-trip through the JSON payload.
- A numeric-looking cell: each test derives its expected row through `toon.Table`.
- A command whose own write changes a fact it reports: each digest reads a commit tree, so the write cannot change it. RE20 proves a repeated run. The dispatcher's tree row names the tree before the write, by its own contract.
- A file that lacks a final newline: RE2.
- Absent versus empty files: RE15, RE31, RE59, and RE62. An oversized excerpt: RE104.
- Special files: RE30 and RE61.
- A dangling symlink where a file is expected: RE28.
- A live symlink where a file is expected: RE29 and RE60.
- Symbolic refs: RE11 resolves `HEAD` to a full object ID.
- An unterminated or a duplicate delimiter: RE6, RE102, and RE32.
- A flag value read as a positional: `usage.Parse` owns flag values, and RE38 grades the missing-flag refusal.
- The primary checkout: RE34 and RE35.

Each refusal leaves the record unchanged: RE17, RE21, RE22, RE31, RE32, RE33,
RE36, RE46, RE53, RE54, RE65, RE79, RE89, RE106, and RE107. RE66 proves that a
refusal at the parse leaves no temporary file, and RE106 proves the same for a
failed temporary write.

Tests that swap a package variable: none. The verb takes its root as a parameter.

**Won't handle:**

- Pending, failed, and skipped results — the checkpoint refuses each until a completed result exists, and step 6 keeps them hand-written.
- A performer that the plan does not owe — `bench gate --checkpoint` grades the owed performer and role.
- Unplanned verification evidence, such as an author's own probe — it stays in the record prose, because the checkpoint grades only planned requirements.
- A chunk result at a source older than the chunk tip — the checkpoint grades only the current source.
- The completion state, performer, and reconciliation — `bench gate --checkpoint --complete` keeps grading the hand-written completion until FT318 adds a form.
- Two `bench record` calls at once in one worktree — one writer owns a worktree, and the last rename wins.
- A crash between the temporary write and the rename — the hidden temporary file stays untracked, and `bench commit` commits only the paths that it names.
- A chunk recorded before its planned predecessor — the checkpoint names the missing predecessor, and serial chunks record in plan order.
- A base that is not an ancestor of the tip — the checkpoint refuses the chain and names the gap.
- A record file with CRLF line ends — the reader already finds no fence in such a file, and Bench writes LF.
- A record file mode other than 0644 — Git keeps only the executable bit, and a record is never executable.

## Ownership fences

- `internal/reviewrecord/files.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/write.go`
- `internal/reviewrecord/record_test.go`
- `internal/reviewrecord/source_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/reviewrecord/recordcmd/`
- `internal/preflight/delegated_evidence_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-implement-spec.md`
- `internal/anchors/`
- `tests/canary/workflow-guidance-anchors/`
- `reviews/record-evidence.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-render-record-fence.md` | none | RE-C1 |
| `2-record-chunk-entry.md` | `1-render-record-fence.md` | RE-C2 |
| `3-record-verification-result.md` | `2-record-chunk-entry.md` | RE-C3 |
| `4-record-review-result.md` | `3-record-verification-result.md` | RE-C3 |
| `5-record-plan-amendment.md` | `4-record-review-result.md` | RE-C4 |
| `6-name-bench-record-in-guidance.md` | `5-record-plan-amendment.md` | RE-C4 |

Tickets 2 to 5 each add one row to the pinned help inventory and one form to
`recordcmd`, so they land in that order. Ticket 4 reuses the excerpt reader of
ticket 3. Ticket 6 names every form, so it lands last.

Tickets 2 to 5 mark `write.go` and `recordcmd/` as new paths. The build
preflight requires that marker while a path is absent from the tree it reads. It
accepts the marker after the path exists.

## Out of scope

- A completion form for the state, the performer, and the reconciliation map, after FT317 decides capability-blocked rows. Estimate: 8 edits, 2 gate runs.
- A `--names` failure-list projection for `bench probe`. Estimate: 5 edits, 1 gate run.
- `bench test --list`. Estimate: 6 edits, 1 gate run.
- Checkpoint refusals that name `bench record` as the repair route. Estimate: 6 edits, 2 gate runs.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| Learning: excerpt digests hashed the whole output | RE57, RE58, RE69 |
| Learning: one source digest included the record file | RE12, RE44, RE107, RE108 |
| Learning: one verification entry sat in the review list | RE42, RE77 |
| Learning: plan digests were computed by hand | RE13, RE84, RE87 |
| Learning: the verb parses the payload before it writes | RE33, RE66, RE79 |
| Q1: four forms, with the completion entry and the projections out | RE11 to RE97, stories 49 and 50 |
| Q2: no `Roadmap:` line | the spec header |
| Q3: the grammar `bench record <form> <slug> …` | RE38, RE39, RE40, RE41, RE71, RE83, RE97, RE105 |
| Q4: `--excerpt` exact bytes, regular, no symlink, not empty, valid UTF-8 | RE57 to RE64, RE104 |
| Q5: planned requirements only, and the command from the plan | RE45, RE46 |
| Q6: the guidance names `bench record` with no implement-phase growth | RE98 to RE101, RE109 to RE111 |
| Q7: only `reviews/<slug>.md`, and the primary checkout refused | RE27, RE34, RE35 |
| Call: version 2 creation, and appends to version 1 | RE15, RE17, RE79 |
| Call: the derived fields | RE12 to RE14, RE44, RE45, RE47, RE48, RE72 to RE74 |
| Call: completed results only | RE50, RE75 |
| Call: the caller supplies `--id` | RE65 |
| Call: the probe flags and the planned mutation | RE52, RE53, RE54 |
| Call: identity mapping and `--map` | RE86, RE88, RE90 to RE92 |
| Call: render, parse again, then replace | RE66, RE106 |
| Call: checkpoint messages unchanged | story 51 |
| Fold 1: no overclaim in step 6 | RE98 |
| Fold 2: the implement phase holds 80 lines | RE101 |
| Fold 3: the duplicate-ID refusal names the ID | RE65 |
| Fold 4: one fence locator, and a duplicate fence refused | RE102 |
| Fold 5: `Save` renders through `Render` | RE103 |
| Fold 6: the control-character row on ungraded flags | RE36 |
| Fold 9: the hand-recording sentences retired | RE109, RE110 |
| Fold 11: oversized excerpt, help spellings, failed write | RE104, RE105, RE106 |
| Advice: `--source` defaults to the chunk tip, and a stale source refuses | RE44, RE107, RE108 |
| Advice and confirming pass 1: the orchestrator records the chunk entry at the freeze, and each author then verifies at that source | RE100, RE111 |
| Confirming pass 2: the RE7 search scope | RE7 |

### Reader and writer sweeps

Readers of `reviews/<slug>.md` and of the fence:

- `internal/reviewrecord/files.go` reads the fence through `Read` and `ReadTree`. The gate checkpoint and `internal/preflight/review.go` call them. They parse JSON, so the new indentation changes no answer.
- `internal/reviewrecord/plan.go` reads the plan fence through `fenced`, which keeps its answer when the locator moves beside it.
- `internal/preflight/fence_writes.go` and `internal/coverage/citations.go` read only the record path.
- `internal/status/status.go` and `internal/spec/spec.go` read only the existence of the file.
- `.bench/BENCH-reference.md` line 327 names the fence and stays true.
- The `rg --hidden` sweep found no `.mjs` script and no `.github` workflow that reads the record or calls `bench record`.

Writers of `reviews/<slug>.md`: hand edits by coordinators and authors,
`recordtest.Fixture.Save` in tests, `bench spec retire`, which deletes the file,
and the new verb. One worktree has one writer at a time, which is invariant 1 of
the operating guide.

Copies of the fence encoding that ticket 1 replaces:

- `internal/reviewrecord/recordtest/fixture.go` line 331, in `Save`.
- `internal/preflight/delegated_evidence_test.go` line 19, in `recordFence`.

The RE7 search covers the production files and `internal/reviewrecord/recordtest`.
A `_test.go` file that builds a malformed document on purpose is exempt by rule,
because no encoder can produce that document. The exempt files include these:

- `internal/reviewrecord/source_test.go` line 166, a duplicate JSON key.
- `internal/reviewrecord/source_test.go` line 228, a record at a path that holds a newline.
- `internal/preflight/delegated_evidence_test.go` line 43, an unterminated fence.
- the planned tests of RE6, RE102, RE32, and RE33, which build malformed documents.

The `recordFence` helper sits in a `_test.go` file, so the RE7 search does not see
it. Review confirms that it renders through `Render`, and RE10 runs it.

The RE7 search is `rg -n -F 'bench-review-record\n' --glob '*.go' --glob '!*_test.go'`.
The RE8 search is `rg -n -F '"bench-review-record"' --glob '*.go' --glob '!*_test.go'`.

### Proof checklist

- Cited symbols: each symbol resolves in the tree at `0c95c944`. The merge of `main` at `9c5d0981` changed only one craft-delegate reference file.
  - In `internal/reviewrecord`: `Read`, `ReadTree`, `Parse`, `ReadPlan`, `SourceDigest`, `Digest`, `RecordPath`, `Slug`, `Axes`, `Check`, `CheckSource`, `ErrMissing`, and the private `fenced`, `readFile`, `parseRecord`, `verifier`, `findChunk`, and `mappedIDs`.
  - In `internal/reviewrecord/recordtest`: `Fixture`, `Save`, `New`, `NewDelegated`, `Attach`, `Prepare`, `Delegate`, and `Native`.
  - In other packages: `bounds.ClassifyNoFollow`, `bounds.ControlRecordLimit`, `toon.RecordError`, `toon.Errorf`, `toon.NotInRepo`, `toon.Table`, `usage.Parse`, `usage.PrimaryCheckoutRefusal`, `git.IsPrimaryCheckout`, `git.ResolveCommit`, `evidencecmd.HelpRows`, `capability.Fifo`, and `capability.Privilege`.
  - In `cmd/bench`: `commandRegistry`, `preflightHelpRows`, `repairPilotCommand`, `boundaryRoot`, `axiReasonMutation`, and `scopeTree`.
  - In `internal/conformance`: `proseBudgetLineCount`.
- Import edges: `internal/reviewrecord/recordcmd` imports `internal/reviewrecord`, `internal/usage`, `internal/toon`, `internal/bounds`, and `internal/git`. `cmd/bench` imports `internal/reviewrecord/recordcmd`. `go list` at `0c95c944` shows that none of those packages imports `cmd/bench` or the new package, so no cycle exists.
- Source-row clauses and occurrences: the source trace table above.
- Promised field labels: `chunk{id,action,base,tip,source_digest,plan_digest,rows}`, `verification{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}`, `review{chunk,id,axis,outcome,supersedes,source_digest,excerpt_digest}`, and `amendment{from,to,chunks}`.
- Changed-function callers: `fenced` keeps its signature, and its callers `parseRecord` and `ReadPlan` keep their answers. `Save` keeps its signature and its bytes. Its callers are in `cmd/bench`, `internal/gate`, `internal/landing`, `internal/worktree`, `internal/systemtest`, and `internal/reviewrecord`. `recordFence` has its callers in its own file.
- Copy survival: RE7, RE8, and RE103.
- Rendered-shape readers: the help golden in `cmd/bench/help_inventory_test.go` is the one reader of the new help rows. Tickets 2 to 5 each extend it. The new TOON tables have no reader before this spec.

Sources re-read in the authoring session:

- `capture/learnings.md` and `roadmap/FT318.md`
- the integration record `reviews/worktree-verb-runner.md`, read in place in its worktree
- each production file of `internal/reviewrecord`, and `recordtest/fixture.go`
- `internal/preflight/review.go`, `internal/preflight/plan.go`, `internal/preflight/decision.go`, and `internal/gate/checkpoint.go`
- `cmd/bench/main.go`, `cmd/bench/command_registry.go`, and `cmd/bench/tree_scope.go`
- `internal/bounds/classify.go`, `internal/structure/structure.go`, `internal/conformance/prose_budget_test.go`, and the two phase files

Not re-read: `internal/conformance/axi_query_registry_test.go` past line 90.

### Help rows

Ticket 2 to ticket 5 add these rows to `bench help`, in this order, at order 40:

- `bench record [--in <label|primary>] chunk <slug> --chunk <id> --base <commit> --tip <commit>` — "write one chunk's frozen pair, digests, and acceptance rows into reviews/<slug>.md"
- `bench record [--in <label|primary>] verification <slug> (--chunk <id> [--source <commit>] | --final --source <commit>) --requirement <id> --id <id> --performer <session> --model <model> --effort <effort> --exit-code <n> --ref <ref> --excerpt <file> [--probe-outcome <verdict> --probe-exit-code <n> --probe-restore pass|fail]` — "append one planned verification result with its computed digests"
- `bench record [--in <label|primary>] review <slug> --chunk <id> --axis Standards|Spec|Coverage --id <id> --performer <session> --model <model> --effort <effort> --ref <ref> --excerpt <file> [--finding <id>]...` — "append one independent review result to a recorded chunk"
- `bench record [--in <label|primary>] amendment <slug> --source <commit> [--map <old>=<new>[,<new>...]]...` — "record the plan-digest change at a source commit"

### Fence disposition

Reviewer disposition of the ownership fences: open. The preflight test file is
in the fence only for the copy that ticket 1 replaces.

The build preflight binds more paths than the tickets edit:

- The registry closure binds `cmd/bench/command_registry_test.go`, `internal/conformance/axi_query_registry_test.go`, and the subcommand-routing table to each `cmd/bench` and anchor registry write.
- The fixture closure binds `tests/canary/package-core-guard/unrouted-subcommand` to `cmd/bench/main.go`. It also binds the `tests/canary/workflow-guidance-anchors/` fixtures to the two phase files.
- The anchor closure binds every anchor registry file in `internal/anchors/` to the two phase files.

Ticket 6 edits only `registry_chunk_chain.go` and `registry_chunk_chain_test.go` in
`internal/anchors/`. No ticket edits a canary fixture, because no fixture mutation
names an edited sentence.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"RE-C1","tickets":["1-render-record-fence.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"preflight","command":"bench test --package ./internal/preflight"}]},{"id":"RE-C2","tickets":["2-record-chunk-entry.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"conformance","command":"bench test --package ./internal/conformance"}]},{"id":"RE-C3","tickets":["3-record-verification-result.md","4-record-review-result.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"cmd","command":"bench test --package ./cmd/bench"}]},{"id":"RE-C4","tickets":["5-record-plan-amendment.md","6-name-bench-record-in-guidance.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/record-evidence/spec.md"},{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"conformance","command":"bench test --package ./internal/conformance"}]}
```

### Flagged additions

Each addition below is not in the decision source. The reviewer can veto each one.

- The heading `# Review outcomes` of a new record file, which keeps the fixture bytes.
- Two-space indentation and no HTML escaping in the payload.
- The four TOON headers, the `action` and `list` cells, and the refusal title `bench record <form> refused`.
- The control-character rule for the single-line flags and `--excerpt`.
- The refusal of an empty record file, and the creation of an absent `reviews/` directory.
- The flag spellings `--final`, `--finding`, and `--map <old>=<new>[,<new>...]`.
- The grammar refusals at exit 2 for a malformed or repeated `--map`, and for `--final` without `--source`.
- The move of the `assessment` adapter, for the headroom of `cmd/bench/main.go`.
- The migration of `Save` and `recordFence` to `Render`, under the one-source standard.
- The anchor rules of ticket 6.
- The amendment change in `coverage.go`, beside `mappedIDs`.
- The refusal of a failed temporary write, with no fallback to an in-place write.
