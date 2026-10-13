# Probe mutant sets

Status: staged

Related roadmap: FT168, partial. This spec settles the mutant-set, multi-package, and mutant-set record-import faces. It closes no roadmap row. The canary and failing-name faces stay on FT168.

Decision source: named reviewed artifact — the approved commitment outcome FT168-probe-batch in .bench/commitment.json with roadmap/FT168.md at 1b1d3c62b, and the reviewer decisions D1 and D2 of 2026-10-13.

Verification log: 4 iteration(s) to accept — the spec review and the ticket review each took 2 iterations on opus / xhigh. The spec round found an interrupt that read `refused` and a lost prose `--expect` rule. The ticket round found a shared seam chunk, a missing helper file, and the missing commitment file.

## Problem

An author who proves several test expectations with `bench probe` pays two focused runs for each expectation.
Each call runs one probe baseline and one mutated run, and each run builds its own private run binary (`internal/testreport/command.go:162`).
In a large package such as `internal/adopt`, each call costs about one minute, and the FT217 ticket 01 author ran 54 serial probes (`docs/known-issues.md:11`).
The FT168 occurrence of the same ticket counts 63 probes in about 17 minutes (`roadmap/FT168.md:71`). The two sources disagree on the count and the cost, and this spec relies on neither number.

Repeated `--swap` and `--omit` edits make one combined mutant, so the grammar has no form for several independent mutants (`internal/probe/command.go:17-20`).
The probe takes one `--package` value, so a mutation that two reader packages must catch needs two calls (`internal/probe/command.go:34`, `internal/probe/probe.go:146`).
`bench record verification --probe-result` reads one complete single-form answer only (`internal/probe/result.go:33-49`).

## Solution

`bench probe <file> --mutants <path> <selection>` reads a list of independent mutants from one strict JSON file.
The verb runs one probe baseline for the whole list, and then preserves, mutates, runs, and restores each mutant in turn.
Each mutated run reuses the baseline's run binary when the mutant leaves the run binary's source closure unchanged.
A mutant inside that closure builds its own run binary.
The answer prints one row for each mutant, a shared selection row with the build count, and an execution row under its own schema token.

`--package` repeats in both forms, and all the named packages run in one `go test` call for each run.
`bench record verification --probe-result <file> --mutant <n>` imports one row of a mutant-set answer as one verification entry.

## Reviewer decisions

- D1 (2026-10-13): A mutant set shares the baseline's private run-binary build with every mutant whose subject lies outside the run binary's source closure. Each mutant inside that closure builds its own run binary, because a shared stale binary grades unmutated code and can answer a false `silent`. The answer reports the build count. This rewords the criterion `FT168-probe-batch.batch` ("one private build") by reviewer decision.
- D2 (2026-10-13): This slice builds the mutant-set record import in `bench record verification`, although the criterion `FT168-probe-batch.record-form` names FT318 as the owner of that form. The spec-stage binding records this decision beside the criterion.

## User stories

Line: opus / high.
Implementation-line reason: PM-B adds a run-owner seam that holds one run binary across several focused runs. PM-C2 adds the mutant-set loop, its stop rules, and a second envelope. The other chunks change one flag, one form matrix, one record form, and guidance at known seams. The canned Go runner and the build-records fixture give strong tests.
Harder chunks: PM-B, PM-C2.

The scorecard decision of 2026-10-12 routes every coding delegate to Opus/high.
The reviewer directed the spec review round and the ticket review round to opus / xhigh.

### Mutant set input

1. As an author, I want `--mutants <path>` to read a list of mutants from one file, so that one call grades several mutants.
2. As an author, I want each entry to take the single form's mutation flags, so that an entry means what the same argv means alone.
3. As an author, I want each entry to take its own `--expect`, so that each mutant names its own expected witness.
4. As an author, I want an entry's repeated `--swap` and `--omit` edits to make one combined mutant, so that sets keep that form.
5. As an author, I want a selection token, a positional, or `--mutants` inside an entry refused, so that entries never change the selection.
6. As an author, I want a mutation flag or `--expect` beside `--mutants` refused as usage, so that each mutation has one source.
7. As an author, I want a hostile or empty mutants file refused before any write, so that a bad path never mutates the tree.
8. As an author, I want a file that is not a strict JSON array of string arrays refused, so that no malformed list runs.
9. As an author, I want an empty list refused, so that a set always grades a mutant.
10. As an author, I want a list over the entry cap refused with the cap, so that one call stays bounded.
11. As an author, I want an entry's mutation refusal to name the entry number before any run, so that I fix the right entry.
12. As an author, I want an unmatched entry to read `invalid` with `substring-miss` while others run, so that one stale entry costs nothing more.
13. As an author, I want an all-miss set to start no run and take no lock, so that a stale list costs nothing.

### Shared run

14. As an author, I want one probe baseline for the whole set, so that N mutants cost N+1 focused runs, not 2N.
15. As an author, I want a failed baseline to end the set before any write, so that a red tree gives no verdict.
16. As an author, I want one subject lock for the whole set, so that no other probe mutates the subject between two mutants.
17. As an author, I want each mutant to preserve, write, run, and restore in turn, so that each mutated run sees only its own mutant.
18. As an author, I want the subject restored byte-exact after each mutant, so that the next mutant starts from the start bytes.
19. As an author, I want a restore failure to keep the copy and stop the set, so that no mutant overwrites a damaged subject.
20. As an author, I want a gate run that starts mid-set to stop the remaining mutants, so that a gate never grades a mutant.
21. As an author, I want a signal in any mutant phase to stop the set after its restore, so that interrupts end the call.
22. As an author, I want a failed mutant write to stop the set with its recovery, so that the tree never stays mutated.

### Private build

23. As an author, I want the baseline's run binary reused for each mutant outside its source closure, so that fixture sets build once.
24. As an author, I want a mutant inside that closure to build its own run binary, so that binary tests grade the mutant.
25. As a maintainer, I want closure membership to come from the run binary's freshness verification, so that no second closure list can drift.
26. As an author, I want a stale inherited run binary to refuse that mutant's run, so that inheritance never becomes permission to build.
27. As an author, I want the answer to report the build count, so that I see what the set cost.
28. As an author, I want each per-mutant build removed after its run, so that a set leaves no private directory.
29. As an author, I want a prose set to build nothing, so that a Markdown set costs no Go build.
30. As a maintainer, I want the single form to keep its two independent builds, so that the closed mutant-build behavior stays.

### Envelope

31. As a consumer, I want a mutant-set answer to carry its own schema token, so that a single-form reader refuses it by name.
32. As an author, I want one row for each mutant in entry order, so that I read each mutant's evidence in one place.
33. As an author, I want each row's exit from the verdict-to-exit map, so that it matches the single form's exit for that verdict.
34. As an author, I want the command exit 0 only when every mutant bit, so that one exit states the whole set.
35. As an author, I want the command exit 2 when a restore failed, so that a damaged subject wins over every other exit.
36. As an author, I want one shared selection row with the build count, so that the shared facts print once.
37. As an author, I want each mutated run's report after the rows under its mutant number, so that I can read each red.
38. As a consumer, I want the execution row to count the trailing bytes, so that a cut answer is visible.
39. As an author, I want an `--expect` miss to change only its own row, so that one wrong witness does not hide the other bites.
40. As an author, I want a mutant that fails to build to name its first build line, so that I see the compile fault.

### Packages

41. As an author, I want `--package` to repeat in `bench probe`, so that one mutation grades several packages in one call.
42. As an author, I want repeated packages to share one `go test` call and one build per run, so that they cost one run.
43. As an author, I want the selection target to name every package expression in argv order, so that the row states what ran.
44. As a reader, I want the probe package note to state the repeat form, so that the help names the route the grammar takes.
45. As an author, I want a mutant set to take repeated `--package` values, so that each mutant grades every named package.

### Record import

46. As an orchestrator, I want `bench record verification --probe-result <file> --mutant <n>` to import row n as one entry, so that a planned probe can cite one mutant of a set.
47. As an orchestrator, I want `--exit-code` checked against that row's own exit, so that a typed exit never stands in for the imported one.
48. As an orchestrator, I want the row's focused exit and restored cell in the entry, so that the checkpoint grades imported values.
49. As an orchestrator, I want a mutant-set answer without `--mutant` refused by name, so that a set never imports as one probe.
50. As an orchestrator, I want `--mutant` with a single-form answer refused, so that the flag never selects a row that is not there.
51. As an orchestrator, I want an out-of-range `--mutant` refused with the row count, so that a wrong number never imports another row.
52. As an orchestrator, I want a bad `--mutant` value or a `--mutant` without `--probe-result` refused as usage, so that the flag never stands alone.
53. As an orchestrator, I want a row whose mutated run never started refused, so that no null exit records.
54. As a consumer, I want the probe package to own both schemas and the record help's field names, so that no rename drifts.

### Guidance and term

55. As a reader, I want `CONTEXT.md` to define **mutant set** apart from **combined mutant**, so that the two terms never blur.
56. As a reader, I want the reference guide to state the set grammar, rows, exits, and build rule, so that it stays current.
57. As a reader, I want the delegate skill to name the mutant set, so that a charge can ask for one call.
58. As a reader, I want the probes guidance to confirm the `restored` cell of each mutant row, so that a set reads safely.
59. As a reader, I want `/bench-implement-spec` to route an author's witness probes of one subject through one mutant set, so that the probe cost drops.
60. As a maintainer, I want the probe-cost known-issues row removed with its anchor, so that the issues table lists only open issues.
61. As a reader, I want the help inventory's probe line and description to state several mutants and packages, so that the inventory states current behavior.
62. As a reader, I want the changelog to name the new forms, so that a release note states the change.

### Reviewed exclusions

63. As a reviewer, I want a mutant set across several files left to a later spec, so that one call keeps one subject.
64. As a reviewer, I want `bench canary <fixture>` and the `--names` comparison left on FT168, so that this spec changes only the probe and its import.
65. As a reviewer, I want the response bound unchanged, so that FT168's open bound decision stays the reviewer's.
66. As a reviewer, I want `first_failure` and `--expect` to keep bare test names across several packages, so that each cell keeps one meaning.

### Roadmap upkeep

67. As a reader, I want `roadmap/FT168.md` to drop its delivered faces and its stale spec reference, so that it names only open faces.

## Implementation decisions

### Term

**mutant set** — the list of independent mutants that one `bench probe` call grades against one probe baseline. Each mutant runs in its own mutated run and has its own restore.
`CONTEXT.md` gains this entry with the Avoid list "probe batch", "mutant batch", "multi-probe".
The **combined mutant** entry stays unchanged, so "mutation batch" and "batch probe" stay avoided terms for it, and no term lists one Avoid word twice.
Guidance never writes "batch probe".

### Grammar

The probe grammar adds `--mutants <path>` as one more mutation group, and `--package` becomes repeatable.
The usage line becomes `bench probe <file> ((<the current single forms>) [--expect <test>] | --mutants <path>) (--package <expr> [--package <expr>]... [--run <go-regex>] | --check <name> | --check system --run <go-regex>) [--full]`.
The line puts `--expect` inside the single-form group, because `gradeForms` refuses it beside `--mutants`.
`Usage` in `internal/probe/command.go` stays the one owner of that line, and the help inventory derives its suffix from it.
`gradeForms` refuses `--mutants` together with any other mutation group flag, `--with`, `--with-file`, `--occurrence`, or `--expect`, with the usage line at exit 2.

The mutation flag declarations move into one list that the top-level grammar and the entry grammar both read.
The entry grammar takes those mutation flags and `--expect`, with no positional.
`gradeForms` splits into a mutation half and a selection half.
The top-level argv runs both halves, and each entry runs the mutation half.
So the form matrix keeps one source.

One rule spans both halves: a prose run names no test, so `--expect` with `--check prose` refuses.
The verb applies that rule to each entry against the top-level selection.
An entry that holds `--expect` under a `--check prose` selection prints the entry-refused line and the usage line at exit 2 before any run.

### Mutants file

The verb reads the `--mutants` path through `readMutationFile`, the bounded reader that `--with-file` and `--omit-file` use.
That reader keeps its refusals: absent, a symlink, a directory, a special file, oversized, unreadable, and empty.
The bytes then decode through `jsonfile.DecodeDocument` into a list of string lists.
A decode error refuses with `error: probe mutant set malformed — <path>: <reason>` at exit 1.

`jsonfile.DecodeDocument` accepts a top-level `null` (`internal/jsonfile/decode.go:58-66`).
So the verb also refuses a top-level `null` and a `null` entry with the malformed line.
An empty entry `[]` reaches the entry grammar, which refuses it because it names no mutation group.
A `null` string inside an entry decodes to an empty string, which the entry grammar refuses through `NoEmptyValue`.

An empty list refuses with `error: probe mutant set empty — <path> holds no entry` at exit 1.
A list of more than 64 entries refuses with `error: probe mutant set too large — <path> holds <n> entries, the cap is 64` at exit 1.
One named constant in `internal/probe` holds the cap of 64.

Each entry parses through `usage.Parse` with the entry grammar.
A parse fault or a mutation-half fault prints `error: probe mutant set entry refused — entry <n> does not match the mutation grammar`, then the top-level usage line, at exit 2.
`<n>` counts from 1 in file order.
Each entry then resolves its mutation through `mutationInput` and `mutate` over the subject's start bytes.
A refusal from either keeps its single-form title, and its detail starts with `entry <n>: `.
Every entry refusal answers before `testreport.Prepare` and before any run child.

### Sequence

The set keeps the single form's refusal order for every per-call step.
That order is usage, repository root, subject, mutants file, entries, selection, prose subject, gate lock, subject lock, and reread.
An unmatched entry marks its row `invalid` with the cause `substring-miss` and the `restored` cell `untouched`.
That row also reads a null focused exit and `ran` 0.
When every entry misses, the rows read an empty `baseline` cell and `builds` 0 at exit 1.
That answer takes no lock and starts no run.

One subject lock and one reread cover the whole set.
The probe baseline runs once through the held runs below.
A failed baseline gives every row `invalid` with the cause `baseline-<kind>` and the `restored` cell `untouched`.
The baseline report then follows the rows at exit 1.

One cancel context, derived through `subprocess.NotifyCancel`, spans the set from the baseline to the last restore.
The held runs take that context for each selection, fresh build, and Go child.
So a signal of `subprocess.CancelSignals` never takes the default action between two runs.
The set restores the subject before it answers.
A SIGKILL and a failed restore stay outside this promise, as they do for the single form.

Then the verb takes each matching entry in file order.
Before each mutant write, the verb examines the cancel context and grades the gate lock again.
Each stop below marks that row and each later row `invalid` with the `restored` cell `untouched`, and the set stops.
A cancelled context gives the cause `set-stopped`.
A held gate lock gives the cause `gate-held`.
A gate state that the verb cannot read gives the cause `gate-unavailable`.

Without a stop, the verb preserves, writes, runs, and restores that mutant as the single form does.
The restore is deferred around the run.
A run that the signal ended reads the outcome kind `interrupted` in every phase, so a cancelled fresh build never reads `refused`.
An `interrupted` row keeps its own cells, and the set stops after its restore.

A failed mutant write gives its row `invalid` with the cause `mutation-failed`.
An unpublished failure releases the copy and reads `restored` as `untouched`.
A published failure restores the subject and reads `restored` as `yes`, and a failed restore keeps the copy as the single form does.
A restore failure gives its row `restore-failed` with the `restored` cell `no`, and the `preserved` table follows the execution row.
After a write failure, an interrupt, or a restore failure, each later row reads `invalid` with the cause `set-stopped` and the `restored` cell `untouched`.

### Held runs

`internal/testreport` gains one run-owner seam that holds one run-binary selection across several focused runs over one prepared request.
The caller passes the seam its cancel context, and each selection, build, and Go child of each held run runs under that context.
`freshness.VerifyAt` takes no context (`internal/freshness/freshness_verify.go:23`), so its rule is a predicate on its return.
A run whose context is cancelled when its validation or its build returns reads the outcome kind `interrupted`.

The first run selects the run binary as `Execute` does today, through `selectRunBinary`.
Before each later run, the owner validates the held binary against the current tree through the run-binary owner's existing validation.
That validation runs `freshness.VerifyAt` over the build inputs that `buildInputs` lists (`internal/freshness/freshness_buildinputs.go:50-151`).
Those inputs are the build files of `./cmd/bench` and its dependencies, `go.mod`, `go.sum`, and each manifest path.
So the freshness digest stays the one source of closure membership, and the owner writes no second closure list.

`internal/runbinary` exports one accessor on `Selection` that reports whether the factory built the binary or inherited it.
The held runs read that accessor, so no code reads `BENCH_RUN_BINARY` a second time.
A test drives the held runs through a factory-valued seam with a stub builder and a stub verifier.
Today `installTestSelectionFactory` installs only the `ReuseOrOwn` function value (`internal/testreport/runbinary_test.go:126-137`).

A held binary that still validates serves the run.
When a held binary that the owner built no longer validates, the owner keeps it and builds one fresh run binary for that run only.
The owner removes that fresh binary's private directory when the run returns.
A held binary from an inherited `BENCH_RUN_BINARY` that no longer validates refuses that run, as `ReuseOrOwn` refuses today.
A signal that ends a validation or a fresh build gives that run the outcome kind `interrupted`, not `refused`.
The owner counts each binary that it builds, and its close removes the held private directory.

A prose request selects no run binary, so its build count is 0.

The single form keeps its two `testreport.Execute` calls, so each of its runs builds its own run binary.
The mutant set alone uses the held runs.

### Envelope

A mutant-set answer has its own schema token, `bench-probe-set-v1`.
The single-form envelope and its token `bench-probe-result-v2` stay unchanged.
The answer prints three tables:

- `mutants[<n>]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}`, with one row for each entry in file order.
- `selection[1]{form,target,run,baseline,builds}`, with the shared selection cells, the baseline kind, and the build count.
- `probe_set_execution[1]{schema,command_exit_code,trailing_bytes}`.

`mutant` counts from 1.
The `mutation`, `cause`, `failed_tests`, `first_failure`, and `ran` cells mean what the single form's cells of the same names mean for that mutant's mutated run.
`first_build_line` names the first build line of a mutated run that failed to build, and it is null otherwise.
`exit_code` follows the single form's verdict-to-exit map: `bit` 0, `silent` 1, `invalid` 1, and `restore-failed` 2.
`command_exit_code` is 2 when a row reads `restore-failed`, 0 when every row reads `bit`, and 1 otherwise.
Each typed cell renders through `toon.TableTyped`.

After the execution row, the `preserved` table follows when a restore failed.
Then each mutant whose mutated run started prints `mutant_report[1]{mutant}` and that run's `bench test` report, in entry order.
`trailing_bytes` counts every byte after the execution row.
The answer spills over the 10-line response bound, and the complete answer is in the spill file.

### Packages

The probe's `--package` flag becomes repeatable, and `selectionArgs` passes each value of `parsed.Repeated["--package"]` in argv order.
`bench test` already runs all its package expressions in one `go test` call (`internal/testreport/command.go:188-201`), so each run keeps one private build.
The `target` cell already joins the expressions with a space (`internal/testreport/selection_facts.go:30-37`).
`probePackageNote` retires, and `ProbeNotes` states `PackageExpressionNote`, because the two grammars now take the same repeat form.

This spec reads the criterion `FT168-probe-batch.packages` ("one private build") as one build for each focused run, shared by every selected package.
A single-form call therefore keeps two builds, one for each run.
The closed row PM33 requires the mutated run to build from the mutant.
The literal reading of one build for each call conflicts with PM33, so the reviewer can veto this reading at sign-off.

### Record import

The verification form adds `--mutant <n>`, valid only with `--probe-result`.
`internal/probe` exports `ReadMutant(data []byte, n int)`, which reads one complete mutant-set answer and returns row n as a `MutantResult`.
`MutantResult` carries `Verdict`, `Restored`, `RowExit`, and `FocusedExit`, so no field named for the command exit carries a row exit.
`ReadResult` refuses a mutant-set answer with the exported sentinel `ErrMutantSet`, and `ReadMutant` refuses a single-form answer with the exported sentinel `ErrSingleForm`.
The two readers share the completeness, spill-frame, strict-decode, and canonical-render rules of today's reader.

With `--mutant`, the import binds `--exit-code` to the row's `exit_code` and the probe outcome to its `verdict`.
It binds the probe exit to the row's `focused_exit_code` and the restore to its `restored` cell.
A null `focused_exit_code` refuses as today.

Without `--mutant`, a mutant-set answer refuses with `--probe-result is a mutant-set answer; name one row with --mutant <n>`.
With `--mutant`, a single-form answer refuses with `--mutant needs a mutant-set answer; --probe-result holds a bench-probe-result-v2 answer`.
An out-of-range row refuses with `--mutant <n> names no row; the answer holds <count> rows`.
A zero, negative, or non-decimal `--mutant`, or a `--mutant` without `--probe-result`, prints the form's usage at exit 2.
The help text names each set field through a probe-owned field-source function, as `probeSources` names the single-form fields today.

### Guidance

- `CONTEXT.md` gains the **mutant set** entry.
- `.bench/BENCH-reference.md` states the `--mutants` form, the repeated `--package`, the rows, the exits, the held-runs build rule, and the stop rules.
- `.agents/skills/bench-craft-delegate/SKILL.md` names the mutant set beside the mutation forms.
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` adds one sentence after line 125: a coordinator confirms the `restored` cell of each mutant row. The anchored sentences at lines 125 and 170 stay byte-identical.
- `.agents/commands/bench-implement-spec.md` adds one sentence after line 27: the author probes several witnesses of one subject in one `bench probe --mutants` call. The anchored witness-cap sentences stay byte-identical.
- `docs/known-issues.md` loses the probe-cost row, because D1's reworded retire condition holds after this spec. The `Require` anchor on its workaround sentence and that anchor's test row go in the same change. The `Forbid` anchor on "Cap the witnesses in a charge" stays.
- `cmd/bench/main.go` changes the probe description to state one or several mutants of one file and a restore after each run.
- `CHANGELOG.md` gains an Unreleased entry.
- `roadmap/FT168.md` drops the delivered mutant-set and multi-package faces and the stale staged-spec reference, and it keeps the canary and failing-name faces. The row stays open.

### Trust chain

The trust chain is unchanged.
The held runs select the run binary through the same owner and validate it through the same freshness check before each reuse.
A held binary never authenticates itself: each reuse recomputes the source digest from the tree.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| PM-A / 1-repeat-the-probe-package.md | Repeated `--package` in the single form, the merged package note, and the help line | PM45-PM50 | planned package tests, existing help and notes tests | no |
| PM-B / 2-hold-one-run-binary-across-runs.md | The held-runs seam with one caller context, reuse on a fresh verification, a per-run build when stale, an interrupted kind for a cancelled build, and a build count | PM27-PM29, PM31, PM70, PM76, PM79, PM80 | planned held-runs tests in `internal/testreport` | yes |
| PM-C1 / 3-share-the-mutation-flag-grammar.md | The shared mutation flag list, the split form matrix, and the named cross-half rule, with no new behavior | none | existing form-matrix tests, the help inventory test | no |
| PM-C2 / 4-grade-a-mutant-set.md | The `--mutants` grammar, file, entry refusals, the set cancel context, sequence, stop rules, envelope, repeated packages in a set, and the probe description | PM1-PM26, PM30, PM32, PM33, PM35-PM41, PM43, PM44, PM51, PM67, PM69, PM71, PM72, PM77, PM78 | planned set and set-input tests, the sequenced canned runner, the help inventory test | yes |
| PM-D / 5-import-one-mutant-row.md | `ReadMutant`, the set refusal of `ReadResult`, and `--mutant` in the verification form | PM34, PM42, PM52-PM60, PM74 | planned import and reader tests, the help inventory test | no |
| PM-E / 6-document-the-mutant-set.md | Term, reference guide, delegate guidance, implement-spec route, known-issues row, roadmap upkeep, and changelog | PM61-PM66, PM68, PM73, PM75 | review-owned reads, the live anchors check, prose gate | no |

PM-A comes first, because PM-C1 and PM-C2 read its repeated `--package` grammar.
PM-B comes before PM-C2, because the set runs through the held runs.
PM-D follows PM-C2, because it reads the set envelope.
PM-E comes last, because its guidance names every form.
PM-A and PM-B both write `internal/testreport`, but they write disjoint files, so they can author at the same time only under `--delegate`.

PM-C1 is a prefactor with no new behavior: the mutation flags move into one list, and `gradeForms` splits into its two halves.
PM-C1 creates the seam that PM-C2 consumes, so its own review closes before PM-C2 starts.
The slice renamed the spec-review chunk PM-C to PM-C1 and PM-C2.

### Completion plan

The authored plan uses version 1, as `references/completion-plan.md` states. The build adds the version 2 execution data.
Each code chunk names one central-property probe.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"probe-mutant-sets-build-20261013","orchestrator_session":"claude:pm-orchestrator-20261013","author_limit":1,"assignments":{"1-repeat-the-probe-package.md":[{"session":"claude:pm_t01","assignment":"a29dbcc8fb1fdd433ff5a8009d1bd4f2","model":"opus","effort":"high","source":"b84b9af315005214ee7a2b768787b7e633dab83f","native_ref":"claude-agent:pm_t01"},{"session":"claude:pm_t01_rep1","assignment":"a29dbcc8fb1fdd433ff5a8009d1bd4f2","model":"opus","effort":"high","source":"9fd00961132dceccafefe1b5942cee6ec0cdaac6","native_ref":"claude-agent:pm_t01_rep1","predecessor":"claude:pm_t01","trigger":"user-directed","stopped":"claude:pm_t01 ended after its ticket commit 868975d7 and its PM-A record commit 14bec38c","preserved":"no uncommitted work; ticket commit 868975d7 and records through 9fd00961 on the probe-mutant-sets-build branch"}],"2-hold-one-run-binary-across-runs.md":[{"session":"claude:pm_t02","assignment":"a29dbcc8fb1fdd433ff5a8009d1bd4f2","model":"opus","effort":"high","source":"e5ca6b0c50fd357d5233499495af5774ae566689","native_ref":"claude-agent:pm_t02"},{"session":"claude:pm_t02_rep1","assignment":"a29dbcc8fb1fdd433ff5a8009d1bd4f2","model":"opus","effort":"high","source":"2c88de8f9ee99166def2bcdd584c1d369d6224c8","native_ref":"claude-agent:pm_t02_rep1","predecessor":"claude:pm_t02","trigger":"user-directed","stopped":"claude:pm_t02 ended after its ticket commit 27346d11 and its PM-B record commit 1b8f4f66","preserved":"no uncommitted work; ticket commit 27346d11 and records through 2c88de8f on the probe-mutant-sets-build branch"}],"3-share-the-mutation-flag-grammar.md":[{"session":"claude:pm_t03","assignment":"a29dbcc8fb1fdd433ff5a8009d1bd4f2","model":"opus","effort":"high","source":"32dc3858c6e91b286869293d5a86ecec7f3ece79","native_ref":"claude-agent:pm_t03"}],"4-grade-a-mutant-set.md":[],"5-import-one-mutant-row.md":[],"6-document-the-mutant-set.md":[]}},"chunks":[{"id":"PM-A","tickets":["1-repeat-the-probe-package.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"1-repeat-the-probe-package.md"},{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"1-repeat-the-probe-package.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"1-repeat-the-probe-package.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"1-repeat-the-probe-package.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"1-repeat-the-probe-package.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Pass only the last --package value from selectionArgs. TestProbeRepeatsThePackage must turn red.","ticket":"1-repeat-the-probe-package.md"}]},{"id":"PM-B","tickets":["2-hold-one-run-binary-across-runs.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"2-hold-one-run-binary-across-runs.md"},{"id":"runbinary","command":"bench test --package ./internal/runbinary","ticket":"2-hold-one-run-binary-across-runs.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"2-hold-one-run-binary-across-runs.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"2-hold-one-run-binary-across-runs.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"2-hold-one-run-binary-across-runs.md"},{"id":"central-property","command":"bench test --package ./internal/testreport","probe":"Reuse the held run binary without its validation before a later run. TestHeldRunsOwnAFreshBuildWhenStale must turn red.","ticket":"2-hold-one-run-binary-across-runs.md"}]},{"id":"PM-C1","tickets":["3-share-the-mutation-flag-grammar.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"3-share-the-mutation-flag-grammar.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"3-share-the-mutation-flag-grammar.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"3-share-the-mutation-flag-grammar.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"3-share-the-mutation-flag-grammar.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Drop the prose --expect case from the cross-half rule. TestExpectRequiresTheNamedWitness must turn red.","ticket":"3-share-the-mutation-flag-grammar.md"}]},{"id":"PM-C2","tickets":["4-grade-a-mutant-set.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"4-grade-a-mutant-set.md"},{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"4-grade-a-mutant-set.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"4-grade-a-mutant-set.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"4-grade-a-mutant-set.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"4-grade-a-mutant-set.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Skip the restore between two mutants, so the second mutant stacks on the first. TestMutantSetRunsEachMutantAlone must turn red.","ticket":"4-grade-a-mutant-set.md"}]},{"id":"PM-D","tickets":["5-import-one-mutant-row.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"5-import-one-mutant-row.md"},{"id":"record-forms","command":"bench test --package ./internal/reviewrecord/...","ticket":"5-import-one-mutant-row.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"5-import-one-mutant-row.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"5-import-one-mutant-row.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"5-import-one-mutant-row.md"},{"id":"central-property","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"Check --exit-code against the command exit instead of the selected row's exit. TestVerificationChecksTheRowExit must turn red.","ticket":"5-import-one-mutant-row.md"}]},{"id":"PM-E","tickets":["6-document-the-mutant-set.md"],"verification":[{"id":"guidance","command":"bench test --check docs-currency-workflow","ticket":"6-document-the-mutant-set.md"},{"id":"conformance-meta","command":"bench test --check conformance-meta","ticket":"6-document-the-mutant-set.md"},{"id":"prose","command":"bench test --check prose","ticket":"6-document-the-mutant-set.md"},{"id":"anchors-package","command":"bench test --package ./internal/anchors","ticket":"6-document-the-mutant-set.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"6-document-the-mutant-set.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"6-document-the-mutant-set.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"6-document-the-mutant-set.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/probe-mutant-sets/spec.md"},{"id":"probe-family","command":"bench test --package ./internal/probe/..."},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"runbinary","command":"bench test --package ./internal/runbinary"},{"id":"record-forms","command":"bench test --package ./internal/reviewrecord/..."},{"id":"anchors-package","command":"bench test --package ./internal/anchors"},{"id":"command-package","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good test drives `probe.Command` in process over the fixture repository of `internal/probe/probe_test.go` (`newFixture`) with the canned Go runner.
It observes the stdout rows, the exit, the subject bytes during and after each run, the Go start count, and the Bench home listing.
Expected rows derive through `toon.TableTyped`, as `probeRow` does.

`probetest` gains one installer that answers the first `go test` start as the baseline and the k-th later start as the k-th mutated answer.
`InstallStarts` and its odd-and-even rule stay for the single form.
The two installers share one script prelude in `probetest`, so the start-count rule has one owner.

A held-runs row drives the new seam in `internal/testreport` through its factory-valued test seam.
That seam takes a counting builder and a verifier that refuses once the tree changes.
The cancelled-build row reuses the parking builder of `internal/testreport/cancel_test.go:276-285`.
A probe-level build row uses `buildEachRun` from `internal/probe/subject_check_fixture_test.go` over `clamp.go`, outside the fixture closure.
A passing in-closure rebuild needs a build script that seals its output against the mutated tree.
So the reuse-and-rebuild rows sit at the `testreport` seam.

PM71 is the one in-closure row at the probe seam.
Its subject is the fixture's `cmd/bench/main.go`, which lies inside the fixture closure.
Its build script copies the sealed pair on its first call and parks on its second call.
The signal ends that parked build, so the row needs no seal for the mutant.
The held runs read the `Selection` accessor to tell an owned binary from an inherited one.
PM29 observes that read: a stale inherited binary refuses with no build.

A record row builds a mutant-set answer through an exported probe renderer or a probe run, never a hand-typed envelope.
The probe run of `cmd/bench/record_operations_test.go:167-179` is the precedent for the run route.
A gate-lock row holds the gate lock from a second process, as `TestProbeRefusesUnderALiveGateRun` does.
An interrupt row signals a helper process during a numbered start, as `interruptedProbe` in `internal/probe/baseline_test.go` does.
The gate observes every row through `bench test` of `./internal/probe/...`, `./internal/testreport`, `./internal/reviewrecord/...`, `./internal/anchors`, and `./cmd/bench`.

### Seam diagram

    trigger: an author or a coordinator runs bench probe <file> --mutants <path> <selection>
        │
        ▼
    argv, mutants file, start bytes  ──▶  [ internal/probe: entries, selection, subject lock, held runs, per-mutant preserve, write, run, restore ]  ──▶  mutants, selection, probe_set_execution, reports
                      ◀ tests attach here: probe.Command over a fixture repository with the sequenced canned runner and the build-records fixture

    trigger: a mutant set runs its second and later focused runs
        │
        ▼
    held selection, current tree  ──▶  [ testreport held runs → runbinary validation → reuse or one fresh build ]  ──▶  outcome, report, build count
                      ◀ tests attach here: the held-runs seam with installTestSelectionFactory, a counting builder, and a tree-aware verifier

    trigger: an orchestrator runs bench record verification ... --probe-result <file> --mutant <n>
        │
        ▼
    answer bytes, row number  ──▶  [ probe.ReadMutant → importedProbe → reviewrecord.RecordVerification ]  ──▶  one verification entry
                      ◀ tests attach here: recordcmd verification tests over a rendered mutant-set answer

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PM1 | 1 | A two-entry set over `clamp.go` prints two `mutants` rows that both read `bit` at command exit 0 | planned TestMutantSetGradesEachEntry in internal/probe/set_test.go | A set that runs only the first entry prints one row. |
| PM2 | 2 | During each mutated run, the subject holds the bytes that the single form makes from the same entry argv | planned TestMutantSetEntryMatchesTheSingleForm in internal/probe/set_test.go | A second entry parser that drops `--occurrence` or `--with-file` writes other bytes, and the hook refuses them. |
| PM3 | 3 | Entry 1 with `--expect TestA` whose run fails only `TestB` reads `invalid` with the cause `expect-missed` | planned TestMutantSetExpectIsPerEntry in internal/probe/set_test.go | An `--expect` read once for the call grades row 1 by no witness and reads `bit`. |
| PM4 | 4 | An entry with two `--swap` pairs makes one combined mutant whose row reads the mutation cell `swap+swap` | planned TestMutantSetEntryTakesACombinedMutant in internal/probe/set_test.go | An entry parser that keeps one swap writes one edit and prints `swap`. |
| PM5 | 5 | An entry that holds `--package`, `--check`, `--run`, `--full`, `--mutants`, or a positional, and the empty entry `[]`, each print the entry-refused line for that entry and the usage line at exit 2 with no Go start | planned TestMutantSetEntryRefusesSelectionTokens in internal/probe/set_input_test.go | An entry grammar copied from the top-level grammar accepts a selection token and changes the shared run. |
| PM6 | 6 | `--mutants` beside `--swap`, `--omit`, `--unwrap`, `--insert-before`, `--omit-file`, `--delete`, `--with`, `--with-file`, `--occurrence`, or `--expect` prints the usage line at exit 2 | planned TestMutantsExcludesTopLevelMutation in internal/probe/set_input_test.go | A grader that skips the new group lets a top-level flag silently join every entry. |
| PM7 | 7 | A mutants path that is absent, a symlink, a directory, a FIFO, oversized, or empty refuses with `probe mutation file unavailable` and leaves the subject bytes unchanged | planned TestMutantsFileRefusesHostilePaths in internal/probe/set_input_test.go | A plain `os.ReadFile` follows the symlink or blocks on the FIFO. |
| PM8 | 8 | A file that holds an object, a string entry, a number token, a top-level `null`, a `null` entry, or a trailing second value refuses with `probe mutant set malformed` before any Go start | planned TestMutantsFileRefusesMalformedJSON in internal/probe/set_input_test.go | A lenient decoder accepts the trailing value and runs the first list. |
| PM9 | 9 | The file `[]` refuses with `probe mutant set empty` | planned TestMutantsFileRefusesMalformedJSON in internal/probe/set_input_test.go | A set with no entry prints a baseline and no row. |
| PM10 | 10 | A list of 65 entries refuses with `probe mutant set too large` and names the cap 64 | planned TestMutantSetEntryCap in internal/probe/set_input_test.go | A list with no cap runs 65 mutated runs. |
| PM11 | 10 | A list of 64 entries that all miss prints 64 rows | planned TestMutantSetEntryCap in internal/probe/set_input_test.go | A cap compared with `>=` refuses the 64th entry. |
| PM12 | 11 | An entry whose `--swap` matches twice refuses with the ambiguous title and a detail that starts with `entry 2: ` before any Go start | planned TestMutantSetEntryRefusalNamesTheEntry in internal/probe/set_input_test.go | A refusal that drops the entry number names no entry, and a late refusal starts the baseline. |
| PM13 | 11 | An entry whose `--with-file` is absent refuses with `probe mutation file unavailable` and a detail that starts with `entry 2: ` | planned TestMutantSetEntryRefusalNamesTheEntry in internal/probe/set_input_test.go | A reader that resolves files after the baseline mutates the tree before the refusal. |
| PM14 | 12 | Entry 1 misses and entry 2 bites, so row 1 reads `invalid`, `substring-miss`, `untouched`, a null focused exit, and `ran` 0, and the Go start count is 2 | planned TestMutantSetMissIsPerEntry in internal/probe/set_test.go | A single-form miss rule refuses the whole call and starts no run. |
| PM15 | 13 | A set whose every entry misses prints every row as `substring-miss` and `builds` 0 at exit 1 with no Go start and no lock entry | planned TestMutantSetOfMissesStartsNoRun in internal/probe/set_test.go | A sequence that runs the baseline first starts one Go child. |
| PM16 | 14 | A three-entry set starts exactly four `go test` children | planned TestMutantSetSharesOneBaseline in internal/probe/set_test.go | A loop that calls the single form starts six children. |
| PM17 | 15 | A failing baseline gives every row `invalid` with the cause `baseline-failed` and `untouched`, with one Go start at exit 1 | planned TestMutantSetBaselineEndsTheSet in internal/probe/set_test.go | A loop that ignores the baseline kind writes the first mutant. |
| PM18 | 16 | A second process that holds the subject lock makes the set refuse with `probe subject in use` before any Go start | planned TestMutantSetTakesOneSubjectLock in internal/probe/set_test.go | A set that skips the lock runs beside another probe of the subject. |
| PM19 | 16 | A lock attempt on the subject from the mutant-2 hook refuses as in use | planned TestMutantSetTakesOneSubjectLock in internal/probe/set_test.go | A lock released between mutants lets the hook's attempt succeed. |
| PM20 | 17 | During mutant 2's run, the subject holds mutant 2's bytes and none of mutant 1's edit | planned TestMutantSetRunsEachMutantAlone in internal/probe/set_test.go | A loop that skips the restore between mutants stacks the edits, and the hook refuses them. |
| PM21 | 18 | After the set, the subject holds its start bytes and start mode, and the Bench home holds no probe entry | planned TestMutantSetRestoresAfterEachMutant in internal/probe/set_test.go | A loop that releases only the last copy leaves a preserved directory. |
| PM22 | 19 | A restore that fails after mutant 1 gives row 1 `restore-failed`, row 2 `invalid` with `set-stopped`, a `preserved` table, and exit 2 with one mutated start | planned TestMutantSetStopsOnARestoreFailure in internal/probe/set_test.go | A loop that continues writes mutant 2 over a damaged subject and starts a third Go child. |
| PM23 | 20 | A gate lock taken from the mutant-1 hook gives row 2 `invalid` with `gate-held` and `untouched`, and no third Go start | planned TestMutantSetStopsUnderAGateRun in internal/probe/set_test.go | A gate check made once per call writes mutant 2 under the gate. |
| PM24 | 21 | An interrupt during mutant 1's run gives row 1 `invalid` with `interrupted` and `yes`, row 2 `set-stopped`, and the start bytes | planned TestMutantSetStopsOnInterrupt in internal/probe/set_test.go | A loop that ignores the interrupted kind starts mutant 2's run. |
| PM25 | 22 | An unpublished write failure on mutant 1 gives row 1 `invalid` with `mutation-failed` and `untouched`, row 2 `set-stopped`, and exit 1 | planned TestMutantSetStopsOnAMutationFailure in internal/probe/set_test.go | A loop that drops the failure runs the unmutated subject and reads `silent`. |
| PM26 | 23 | With an owned run binary, a two-entry set over `clamp.go` records exactly one build | planned TestMutantSetReusesTheBuildOutsideTheClosure in internal/probe/set_test.go | A set over two `Execute` calls for each mutant records three builds. |
| PM27 | 24 | After the verifier refuses the held binary, the next held run builds one fresh run binary whose builder reads the changed tree | `internal/testreport/held_runs_test.go` (`TestHeldRunsOwnAFreshBuildWhenStale`) | An owner that reuses the held binary without validation records one build. |
| PM28 | 25 | Each held run after the first calls the run-binary verifier on the held path, and an accepting verifier gives reuse | `internal/testreport/held_runs_test.go` (`TestHeldRunsReuseOnAFreshVerify`) | An owner with its own file list never calls the verifier stub. |
| PM29 | 26 | With an inherited `BENCH_RUN_BINARY` that the verifier refuses on the second run, that run's outcome kind reads `refused` and the builder records no call | `internal/testreport/held_runs_test.go` (`TestHeldRunsRefuseAStaleInheritedBinary`) | An owner that treats a stale inherited binary as owned builds one. |
| PM30 | 27 | The `builds` cell of a set over an inherited run binary reads 0, and over an owned binary outside the closure it reads 1 | planned TestMutantSetReportsTheBuildCount in internal/probe/set_test.go | A cell that counts runs reads 3. |
| PM31 | 28 | After a stale run returns, its fresh private directory is gone, and after close no held directory remains | `internal/testreport/held_runs_test.go` (`TestHeldRunsRemoveEachOwnedBuild`) | An owner that closes only the held selection leaves the fresh directory. |
| PM32 | 29 | A two-entry prose set prints `builds` 0 | planned TestMutantSetOfProseBuildsNothing in internal/probe/set_test.go | An owner that selects a run binary before the prose route reads 1. |
| PM33 | 30 | The single form builds twice and the second build reads the mutant bytes | `internal/probe/subject_check_test.go` (`TestProbeBuildsTheMutantRunBinary`) | A single form moved onto the held runs records one build. |
| PM34 | 31 | `probe.ReadResult` of a mutant-set answer returns `ErrMutantSet` | planned TestReadResultRefusesAMutantSet in internal/probe/result_test.go | A reader that answers the first-table miss names no mutant set. |
| PM35 | 32 | The first line of a two-entry answer is the exact `mutants[2]{...}` header from the envelope design | planned TestMutantSetRowsKeepTheirSchema in internal/probe/set_test.go | A row that drops `focused_exit_code` changes the header. |
| PM36 | 33 | In a set with one biting and one silent mutant, the rows read `exit_code` 0 and 1 | planned TestMutantSetCommandExit in internal/probe/set_test.go | A row exit copied from the command exit reads 1 on both rows. |
| PM37 | 34 | A set whose two rows read `bit` exits 0 | planned TestMutantSetCommandExit in internal/probe/set_test.go | A command exit that reads 1 for any set with more than one row reads 1. |
| PM38 | 34 | A set whose rows read `bit` and `silent` exits 1 | planned TestMutantSetCommandExit in internal/probe/set_test.go | A command exit taken from the first row reads 0. |
| PM39 | 35 | A set whose rows read `bit` and `restore-failed` exits 2 | planned TestMutantSetCommandExit in internal/probe/set_test.go | A command exit that ranks `invalid` above `restore-failed` reads 1. |
| PM40 | 36 | The second table is one `selection[1]{form,target,run,baseline,builds}` row | planned TestMutantSetRowsKeepTheirSchema in internal/probe/set_test.go | A selection row printed for each mutant repeats the table. |
| PM41 | 37 | After the execution row, each started mutant prints `mutant_report[1]{mutant}` and its report in entry order, and a missed entry prints none | planned TestMutantSetPrintsEachReport in internal/probe/set_test.go | A render that keeps only the last report loses the first red. |
| PM42 | 38 | `trailing_bytes` equals the byte count after the execution row, and `ReadMutant` refuses an answer cut by one byte | planned TestReadMutantRequiresTheCompleteAnswer in internal/probe/result_test.go | A count that skips the reports passes a cut answer. |
| PM43 | 39 | In a set where row 1 reads `expect-missed`, row 2 still reads `bit` | planned TestMutantSetExpectIsPerEntry in internal/probe/set_test.go | A set that stops at the first `invalid` row marks row 2 `set-stopped`. |
| PM44 | 40 | A mutated run whose canned build output starts with `clamp.go:3:2: "os" imported and not used` gives its row the cause `build-failed` and that exact `first_build_line` cell, and the other row's cell is null | planned TestMutantSetNamesABuildFailure in internal/probe/set_test.go | A build line kept on the shared rows names the wrong mutant. |
| PM45 | 41 | A single-form probe with `--package ./clamp --package ./other` passes both expressions in argv order to each `go test` start | `internal/probe/package_test.go` (`TestProbeRepeatsThePackage`) | `selectionArgs` that reads `Flags` passes only `./other`. |
| PM46 | 42 | That probe starts exactly two `go test` children and records one build for each run | `internal/probe/package_test.go` (`TestProbePackagesShareOneBuildPerRun`) | A loop over packages starts four children and four builds. |
| PM47 | 43 | That probe's `target` cell reads `./clamp ./other` | `internal/probe/package_test.go` (`TestProbeRepeatsThePackage`) | A target from the last flag value reads `./other`. |
| PM48 | 44 | Every help spelling prints the usage line and the notes `PackageExpressionNote` and the check note | `internal/probe/refusal_test.go` (`TestProbeHelpSpellings`) | The retired note's "run one call per package" survives in the help. |
| PM49 | 44 | `ProbeNotes` starts with `PackageExpressionNote` and holds two facts | `internal/testreport/selection_facts_test.go` (`TestProbeNotesDeriveTheirTokens`) | A second probe note copy keeps the old remedy. |
| PM50 | 44 | The help inventory's probe line names `[--package <expr>]...` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A hand-copied suffix keeps the one-package grammar. |
| PM51 | 45 | A two-entry set with two `--package` values passes both expressions to each of its three Go starts | planned TestMutantSetTakesRepeatedPackages in internal/probe/set_test.go | A set that builds its request from one value drops `./clamp`. |
| PM52 | 46 | `--mutant 2` of a two-row answer writes one entry whose probe outcome is row 2's verdict | planned TestVerificationImportsOneMutantRow in internal/reviewrecord/recordcmd/verification_import_test.go | An import that reads row 1 records row 1's verdict. |
| PM53 | 47 | `--exit-code 1` against a row whose `exit_code` is 0 refuses with a line that names `mutants.exit_code` | planned TestVerificationChecksTheRowExit in internal/reviewrecord/recordcmd/verification_import_test.go | A check against the command exit passes when the command exit is 1. |
| PM54 | 48 | The entry's probe exit equals row 2's `focused_exit_code`, and its restore reads `pass` | planned TestVerificationImportsOneMutantRow in internal/reviewrecord/recordcmd/verification_import_test.go | An import that takes the row exit as the probe exit records 0. |
| PM55 | 49 | A mutant-set answer without `--mutant` refuses with the exact line `--probe-result is a mutant-set answer; name one row with --mutant <n>` | planned TestVerificationRefusesAMutantSetWithoutARow in internal/reviewrecord/recordcmd/verification_import_test.go | A refusal from the first-table miss names no route. |
| PM56 | 50 | `--mutant 1` with a single-form answer refuses with the exact line `--mutant needs a mutant-set answer; --probe-result holds a bench-probe-result-v2 answer` | planned TestVerificationMutantRefusals in internal/reviewrecord/recordcmd/verification_import_test.go | An import that ignores `--mutant` records the single row. |
| PM57 | 51 | `--mutant 3` of a two-row answer refuses with the line that names 2 rows | planned TestVerificationMutantRefusals in internal/reviewrecord/recordcmd/verification_import_test.go | A clamp to the last row imports row 2. |
| PM58 | 52 | `--mutant 0`, `--mutant x`, and `--mutant 1` without `--probe-result` each print the form usage at exit 2 | planned TestVerificationMutantRefusals in internal/reviewrecord/recordcmd/verification_import_test.go | A parser that accepts zero selects no row and reports out of range instead. |
| PM59 | 53 | `--mutant 1` of a row that reads `substring-miss` refuses with the null focused-exit line | planned TestVerificationMutantRefusals in internal/reviewrecord/recordcmd/verification_import_test.go | An import that maps null to 0 records a probe that never ran. |
| PM60 | 54 | The help inventory's record verification line names `--mutant <n>` and `mutants.exit_code` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A help text that spells the set field by hand drifts from the probe owner. |
| PM61 | 55 | `CONTEXT.md` defines **mutant set** and keeps the **combined mutant** entry byte-identical | review-owned | A merged entry blurs the two terms. |
| PM62 | 56 | The reference guide's probe notes state the `--mutants` form, the rows, the exits, the build rule, and the stop rules | review-owned | A guide that keeps "one package or one check" misstates the grammar. |
| PM63 | 57 | The delegate skill's probe sentence names the mutant set | review-owned | A skill that lists only single forms never routes a charge to the set. |
| PM64 | 58 | The probes guidance adds the per-row `restored` sentence, and the anchored needles at lines 125 and 170 stay byte-identical | review-owned | A rewritten anchored sentence reds the anchors check and the `delegate-probe-mutated-bytes` canary. |
| PM65 | 59 | `/bench-implement-spec` adds the mutant-set witness sentence, and the anchored witness-cap sentences stay byte-identical | review-owned | A reworded cap sentence reds `TestWriteSpecChunkedSlicing`. |
| PM66 | 60 | The `Require` anchor on the probe-cost workaround sentence and that anchor's test row are gone, and the `Forbid` anchor on "Cap the witnesses in a charge" stays | `bench test --check docs-currency-workflow` (the live anchors over `docs/known-issues.md`) | A removed row with a kept `Require` anchor reds the live anchors check. |
| PM73 | 60 | `docs/known-issues.md` no longer holds the probe-cost row | review-owned | A kept row whose anchor is gone reds no check, so review reads the file. |
| PM67 | 61 | The help inventory's probe description states one or several mutants of one file and a restore after each run | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A description that keeps "mutate one file once" misstates the set. |
| PM69 | 61 | The help inventory's probe line names `--mutants <path>` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A suffix that predates the set form omits `--mutants`. |
| PM68 | 62 | `CHANGELOG.md` holds an Unreleased entry for the mutant set, the repeated `--package`, and `--mutant` | review-owned | A missing entry hides the change from a release note. |
| PM70 | 21 | A cancel of the caller context while a stale held run's fresh build parks gives that run the outcome kind `interrupted`, and the parked builder's directory is gone | `internal/testreport/held_runs_test.go` (`TestHeldRunsInterruptACancelledBuild`) | An owner that maps the cancelled build through `selectionFailure` reads `refused`. |
| PM71 | 21 | A signal during mutant 1's fresh build gives row 1 `invalid` with `interrupted`, row 2 `set-stopped` with no write, exactly one Go start in all, and the subject at its start bytes | planned TestMutantSetStopsOnAnInterruptedBuild in internal/probe/set_test.go | A set that opens a new context for each run reads row 1 `refused` and writes mutant 2. |
| PM76 | 21 | A caller context cancelled while the held validation of an inherited binary runs gives that run the outcome kind `interrupted`, not `refused` | `internal/testreport/held_runs_test.go` (`TestHeldRunsInterruptACancelledValidation`) | An owner that maps every failed validation through `selectionFailure` reads `refused`. |
| PM79 | 24 | When a stale owned held binary's fresh build fails with no cancel, that run reads the outcome kind `build-failed` and starts no Go child | `internal/testreport/held_runs_test.go` (`TestHeldRunsFailAStaleRunWhoseBuildFails`) | An owner that falls back to the stale held binary runs the unmutated code and can read `passed`. |
| PM80 | 21 | A cancel of the caller context while the first held run's build parks gives that run the outcome kind `interrupted`, removes the parked directory, and leaves the build count at 0 | `internal/testreport/held_runs_test.go` (`TestHeldRunsInterruptACancelledFirstBuild`) | An owner that applies the interrupt rule only to later runs reads `refused` for a cancelled first build. |
| PM77 | 20 | A gate state that the verb cannot read before mutant 2's write gives row 2 `invalid` with `gate-unavailable` and `untouched`, and no second mutated start | planned TestMutantSetStopsOnAnUnreadableGate in internal/probe/set_test.go | A loop that reads only the held answer of the gate check writes mutant 2. |
| PM78 | 22 | A published write failure on mutant 1 whose restore succeeds gives row 1 `mutation-failed` with `restored` `yes`, and the subject holds its start bytes | planned TestMutantSetStopsOnAMutationFailure in internal/probe/set_test.go | A row that copies the unpublished cell reads `untouched` after a write that reached the subject. |
| PM72 | 2 | Under `--check prose`, an entry that holds `--expect` prints the entry-refused line for that entry and the usage line at exit 2 with no run | planned TestMutantSetEntryRefusesProseExpect in internal/probe/set_input_test.go | An entry that runs only the mutation half accepts the pair and answers every bite `expect-missed`. |
| PM74 | 54 | `ReadMutant` refuses a spill frame, a non-canonical row, and a repeated `mutants` table before it returns a row | planned TestReadMutantSharesTheReaderRules in internal/probe/result_test.go | A second reader that skips the shared rules imports a projection or a hand-edited row. |
| PM75 | 67 | `roadmap/FT168.md` no longer names the mutant-set or multi-package face or the retired staged spec, and it keeps the canary and failing-name faces | review-owned | A stale row tells the next drain to plan a delivered face. |

Not covered: story 63 — a reviewed exclusion whose cost the Out of scope section prices.
Not covered: story 64 — a reviewed exclusion, because FT168 keeps the canary and failing-name faces.
Not covered: story 65 — a reviewed exclusion, because FT168 keeps the open response-bound decision.
Not covered: story 66 — a reviewed exclusion whose bare-name rule the edge inventory states.

### Edge inventory

The canonical edge classes at the probe seam:

- Empty input: an empty mutants file (PM7), an empty list (PM9), an empty entry `[]` (PM5), and a `null` list or entry (PM8).
- Boundary: 64 and 65 entries (PM10, PM11), and `--mutant 0` and an out-of-range row (PM57, PM58).
- Hostile path: a symlink, a directory, a FIFO, and an oversized file (PM7).
- Malformed structure: an object, a string entry, a number token, and a trailing value (PM8).
- Ordering: entry refusals before any run (PM12, PM13), the miss before the lock (PM15), and the gate lock before each write (PM23).
- Partial failure: a write failure, an interrupt, and a restore failure stop the set (PM22, PM24, PM25, PM70, PM71, PM76). A fresh build that fails reads `build-failed`, and a cancelled first build reads `interrupted` (PM79, PM80).
- Cross-half rule: an entry `--expect` under a prose selection (PM72).
- Concurrency: a second probe and a gate run during a set (PM18, PM19, PM23).
- Audience: this repository's kit checkout and every linked repository read the same verb. A linked repository with no `cmd/bench` closure gets a run-binary refusal at the first run, as the single form does today.

**Won't handle:**

- A mutant set across several subject files — one call keeps one subject, one lock, and one preserved copy. A caller runs one set for each file.
- A same-named test in two selected packages — `first_failure` and `--expect` keep bare names as the single form does. The failures table of each report names the package.
- A wider response bound for a set answer — FT168 holds that reviewer decision, and the spill file holds the complete answer.
- A held run binary for the single form — the closed mutant-build row PM33 keeps two builds, and a one-entry mutant set gives the reuse.
- A duplicate entry in a set — each entry runs as written, and the rows show both verdicts.

## Ownership fences

The fence is the union of the six ticket `Writes:` lines and the review pickup.

- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.bench/BENCH-reference.md`
- `.bench/commitment.json`
- `CHANGELOG.md`
- `CONTEXT.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/record_operations_test.go`
- `cmd/bench/test_command.go`
- `docs/known-issues.md`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/probe/baseline.go`
- `internal/probe/baseline_test.go`
- `internal/probe/combined_test.go`
- `internal/probe/command.go`
- `internal/probe/edit.go`
- `internal/probe/edit_test.go`
- `internal/probe/lock_holder_test.go`
- `internal/probe/lock_test.go`
- `internal/probe/mutation_input.go`
- `internal/probe/outcome_test.go`
- `internal/probe/package_test.go`
- `internal/probe/probe.go`
- `internal/probe/probe_test.go`
- `internal/probe/probetest/fixture.go`
- `internal/probe/refusal_test.go`
- `internal/probe/replacement_test.go`
- `internal/probe/result.go`
- `internal/probe/result_test.go`
- `internal/probe/set.go`
- `internal/probe/set_input.go`
- `internal/probe/set_input_test.go`
- `internal/probe/set_test.go`
- `internal/probe/subject.go`
- `internal/probe/subject_check_fixture_test.go`
- `internal/probe/unwrap_test.go`
- `internal/probe/verdict.go`
- `internal/probe/verdict_test.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/verification.go`
- `internal/reviewrecord/recordcmd/verification_import_test.go`
- `internal/reviewrecord/recordcmd/verification_test.go`
- `internal/runbinary/runbinary.go`
- `internal/runbinary/runbinary_test.go`
- `internal/testreport/cancel_test.go`
- `internal/testreport/command.go`
- `internal/testreport/held_runs.go`
- `internal/testreport/held_runs_test.go`
- `internal/testreport/outcome.go`
- `internal/testreport/runbinary_test.go`
- `internal/testreport/selection_facts.go`
- `internal/testreport/selection_facts_test.go`
- `projects/benchkit.md`
- `reviews/probe-mutant-sets.md`
- `roadmap/FT168.md`
- `tests/canary/claude-agent-definitions/`
- `tests/canary/docs-currency-token-diet/`
- `tests/canary/guidance-prose-budgets/`
- `tests/canary/line-routing/`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skill-description-budgets/`
- `tests/canary/skills-index-command-adapters/`
- `tests/canary/workflow-guidance-anchors/`

## Out of scope

- A mutant set across several subject files: 8 edits, 3 gate runs. It needs a subject list, per-subject locks in a fixed order with a rollback, and a preserved directory for each subject.
- `bench canary <fixture>`: FT168 keeps it.
- The `--names` failing-name comparison: FT168 keeps it.
- The response-bound shape of a biting answer: FT168 keeps that reviewer decision.

## Further notes

### Flagged additions

- The gate-lock check before each mutant write (story 20) is a safety addition beyond the source text. A gate run can start during a long set, and the single form checks the gate once.
- The 64-entry cap (story 10) is an addition from the hostile-input walk.
- The removal of the probe-cost known-issues row (story 60) follows the known-issues header rule once D1 holds.
- The completion plan uses version 1 as `references/completion-plan.md` states, and not the version 2 shape that the coordinator suggested.
- The `/bench-implement-spec` sentence of story 59 is a workflow addition beyond the source text. It routes an author's witness probes of one subject through one mutant set.
- The reading of "one private build" for repeated packages as one build for each run (the Packages decision) is open to reviewer veto.
- The gate-unavailable stop cell and the set cancel context come from the round 1 spec review.

### Source sentence to row

| source sentence | rows |
| --- | --- |
| `FT168-probe-batch.batch`: one `bench probe` call grades several mutants against one shared baseline run and one private build (reworded by D1) | PM1, PM16, PM26, PM27, PM30 |
| `FT168-probe-batch.packages`: one `bench probe` call grades several packages in one private build | PM45, PM46, PM47, PM51 |
| `FT168-probe-batch.record-form`: the spec stage decides how `bench record` imports a batch probe result (built here by D2) | PM52-PM59 |
| `FT168-probe-batch.open`: the delivery closes no roadmap row, and the canary and failing-name faces stay on FT168 | the header's `Related roadmap` line and the Out of scope section |
| `roadmap/FT168.md:10` "One probe takes several mutations." | PM1, PM2 |
| `roadmap/FT168.md:11` "One probe grades several packages in one private build." | PM45, PM46 |
| `roadmap/FT168.md:38-40` the batch face | PM16, PM26 |
| `docs/known-issues.md:11` retire condition | PM66 |

### Pre-review proof checklist

- `Cited symbols`: each symbol below was read in this session.
  - `internal/probe`: `Usage`, `grammar`, `gradeForms`, `namedGroup`, `mutationInput`, `readMutationFile`, and `selectionArgs`.
  - `internal/probe`: `runWith`, `probe`, `preserve`, `restore`, `render`, `envelope`, and `resultSchema`.
  - `internal/probe`: `ReadResult`, `FieldSource`, `verdictFor`, `mutatedCells`, `firstFailure`, and `witnessed`.
  - `internal/testreport`: `Prepare`, `Execute`, `runFocusedRequest`, `selectRunBinary`, `ProbeNotes`, and `PackageExpressionNote`.
  - `internal/testreport`: `probePackageNote`, `installTestSelectionFactory`, and `OutcomeInterrupted`.
  - `internal/runbinary`: `ReuseOrOwn`, `Factory`, `Inherit`, `Own`, and `Selection.Close`.
  - Other owners: `freshness.VerifyAt`, `jsonfile.DecodeDocument`, `importedProbe`, `probeSources`, and `verificationForm`.
  - Test owners: `InstallStarts` in `internal/probe/probetest`, and `buildEachRun` in `internal/probe/subject_check_fixture_test.go`.
- `Import edges`: `internal/probe` → `internal/jsonfile` is new. `go list` shows `internal/jsonfile` imports only the standard library, so the edge makes no cycle. Every other edge exists.
- `Source-row clauses and occurrences`: the source-sentence table above. The FT168 occurrences that drive this spec are the 2026-10-09 FT217 ticket 01 and build lines, the 2026-09-30 tree-targets line, and the 2026-09-27 drain line.
- `Promised field labels`: `mutant`, `verdict`, `mutation`, `cause`, `failed_tests`, `first_failure`, `ran`, `focused_exit_code`, `exit_code`, `first_build_line`, `restored`, `form`, `target`, `run`, `baseline`, `builds`, `schema`, `command_exit_code`, `trailing_bytes`.
- `Changed-function callers`: the slicer reruns `bench consumers` for each function below.
  - `selectionArgs`: `runWith` and `internal/probe/outcome_test.go:204-212`.
  - `gradeForms`: `Command` and `internal/probe/outcome_test.go:204-212`.
  - `ProbeNotes`: `probeGrammar` and `TestProbeNotesDeriveTheirTokens`.
  - `importedProbe`: `verification`.
  - `ReadResult`: `importedProbe` and `internal/probe/result_test.go`.
- `Copy survival`: PM48 and PM49 fail if the retired probe package note survives.
- `Rendered-shape readers`: the probe help line (`cmd/bench/help_inventory_test.go:122`), the record verification help line (`cmd/bench/help_inventory_test.go:148`), the probe notes copy (`internal/probe/refusal_test.go:221-227`), and the notes test (`internal/testreport/selection_facts_test.go:148-167`).
- `Pin operators`: `TestHelpInventoryIsComplete` compares each line for equality. `TestProbeHelpSpellings` compares the whole help for equality. `TestProbeNotesDeriveTheirTokens` applies a prefix check and a count of 2.
- `Entry reads`: none.
- `Derived expectations`: the expected rows derive through `toon.TableTyped` from the row design, not from the renderer under test.
- `Consolidated rules`: the probe package note merges into `PackageExpressionNote`. Old rule: "run one call per package or use a parent ./... pattern". New rule: "repeat --package to run several expressions in one call, or use a parent ./... pattern". PM48 and PM49 cover the change.
- `Quantified obligations`: "each entry" refusals (PM5, PM12, PM13), "each mutant" restore (PM20, PM21), "each later run" validation (PM28), and "every row" exits (PM36-PM39).
- `Workflow-step writes`: none. The record import writes one verification entry through the existing transaction.
