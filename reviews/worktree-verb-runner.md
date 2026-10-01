# Worktree verb runner review record

## VR-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `vr-t1-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `7872c5e4` and committed `137fb303` on a lane pass in the first attempt. The author then committed this record in a second commit.

The chunk pair is `0c95c944..137fb303`. The base is the `main` tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The coordinator moves the tip to the record commit. The source digest does not change, because the record file is outside the graded source.

The author wrote each test before the code that it grades, in three stages. The red and green log for each row follows:

- Stage 1: every runner function and reader was a stub, and the runner returned exit -1. All 22 planned tests failed, and VR1 failed in each of its 14 key subtests, which gives 36 failure rows. The log records 35 failures. The stub source is not in the history, so the count is not reproducible.
- Stage 2: the runner and `checkVerbCall` were complete, and the readers and must forms were stubs. VR1 to VR5, VR17 to VR19, VR21, and VR22 passed. VR6 to VR14, VR16, and VR59 failed for the reader reason, with 11 failures. The log does not record the stage 2 result of VR15, and that result is not reproducible.
- Stage 3: after the readers and the must forms, all 22 tests passed.
- VR15: in stage 1, the test failed in its fixture only, because the stub runner gave no plan. A probe of the finished must form gave the row its own red. The probe made every reader error fail the recorder, and the test failed on the `none` plan.
- VR20 is review-owned. The call value declares the `kit` and `clock` fields beside `root` and `home`.

The author reported these deviations from the ticket. The Spec axis grades each one.

- The `show` and `build` grammars refuse an extra argument with the command name only. So the VR1 expectation for those two keys reads the command name from the first three words of the usage constant.
- The VR59 test builds the faulted unclaimed set with the four fixture lines of `TestCleanUnclaimedErrorRowRefusesTheSet`, as the ticket tells. That file is outside the fence, so the author did not extract a shared builder. This is a second copy of one fixture.
- The runner tests do not call `newOwnedAssignment`. That builder returns a positional tuple, ticket 3 changes it, and the ticket 3 `Writes:` line does not hold the runner files. A one-value builder, `runnerRepo`, returns a call value, and each test creates its assignment with `mustCreate`. A first draft had a local tuple builder, and the duplicated-facts sweep removed it, because VR40 refuses a tuple in any test file.
- VR17 stubs the `buildSubject` field of the joins value, so the default build is never reached.
- VR5 compares a direct `ExecCommand` call with the runner, because a child that writes both streams and exits 3 grades all three values at once.
- The runner also fails the test for an unknown key. The fingerprint reader returns an error for a cell that does not decode as text. No test grades that branch.
- No check required an edit to the five registry paths on the `Writes:` line, so the diff leaves them unchanged.

The author computed the plan digest from the `ReadPlan` rule. The input is the JSON array of the spec bytes and the eleven ticket bytes, in plan order. The same rule gives the digest `7355203a` for the tree targets record at `a08359f9`, which confirms the method. The source digest is the tree of `137fb303`, which does not hold this record file.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The second row is the plan probe `1-probe`, and the third row is the author's own probe of the joins dispatch. The JSON payload holds the exact command and output of those two rows. The first row ran on the source before the ticket commit, so only this table holds it.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/worktree/verb_runner_test.go` | swap: `if errors.Is(err, errNoVerbFingerprint) {` to `if err == nil {` | TestVerbResultMustNoFingerprintAcceptsAnErrorPlan | VR15 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `unapplicableFingerprint` to `"probe-placeholder"` | TestVerbResultFingerprintTreatsAPlaceholderAsAbsent | VR59 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `if call.joins != nil {` to `if false {` | TestVerbRunnerPassesTheJoinsValue | VR17 | bit |

### Verification

The author ran each VR-C1 verification on the source of `137fb303`, and each passed. The JSON payload holds each result. The package excerpt omits its two skip rows. Each skip is an environment capability skip for unix sockets.

`TestPackageTestCountPin` passed at 686 tests, and `TestSerialSetStaysBelowTheCeiling` passed at the ceiling of 46. The author also ran these checks, and each passed: `bench test --package ./internal/conformance` and `bench structure --growth 0c95c944`. The commit chain reported 15 green checks, 1 check that does not apply, and 0 red checks.

### Coordinator verification

The coordinator found a clean tree, and the ticket commit touches only paths on the `Writes:` line. The coordinator ran an independent omission probe at a site that no author probe used. The probe removed the two-fingerprints refusal in `readVerbFingerprint`, and `TestVerbResultFingerprintRefusesConflictingCells` failed. The restore reads `yes`.

## VR-C1 chunk review, round 1

The frozen pair is base `0c95c9447c20189f3f2155719ef965bffc339856` and tip `4479b2fa62bcb334f4b65aedfcd155cde31d66a6`. The coordinator moved the chunk tip from the ticket commit to the record commit, and the source digest stays the same. The shared evidence is `sha256:e3159ba832b81092e567ee0470a52d17f5da030e7694a5fb0fa17e4318cc1ef4`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 9: Standards 4, Spec 3, and Coverage 2. The Standards literal finding and the Spec literal finding name the same fix. Two findings are `no-op`. So 6 repair targets remain, and they take repair cycle 1 of 2.

The chunk record for VR44 to VR46 follows. The PASS set at the tip is the base set plus the 22 added tests, and all 686 top-level tests pass. The SKIP set holds the two socket capability subtests, and the base has the same two. No existing test changed except the `worktreeTestCount` value, so no assertion fell.

## Standards

Findings: 4. The worst issue is a second copy of the faulted unclaimed fixture.

- R1: `internal/worktree/verb_runner_check_test.go:268-278` copies the five fixture lines and the reason comment of `clean_unclaimed_test.go:312-321`. `AGENTS.md` names a fixture harness pasted N times as duplicated knowledge. Extract one builder beside `addUnclaimedBranch` in `clean_set_apply_test.go`, and let both tests call it. Both files are in the spec fence, so the repair takes a plan commit for the ticket 1 `Writes:` line. `auto-fix`. Confidence 8.
- R2: `internal/worktree/verb_runner_check_test.go:83` restates the reason text `"invalid invocation; run "` from the producer at `worktree.go:247`. The reader sweep says that no literal moves into a test. The `clean` subtest can render `cleanInvocationError` into a buffer and compare the whole output. `auto-fix`. Confidence 7.
- R3: `internal/worktree/verb_runner_check_test.go:39`, the `usageCommand` helper, derives the refusal command name a second time from the first three words of the usage constant. The verbs take that name from `worktreeShowGrammar.Cmd` and `buildGrammar.Cmd`. Read those two fields instead. The production duplication between the `Cmd` literals and the usage constants is outside the fence and goes to the ideas inbox. `auto-fix`. Confidence 5.
- R4: `internal/worktree/verb_runner_check_test.go:206` asserts that the fingerprint is the last record cell, which is a layout fact of `reset.go:67` that the reader contract does not need. Assert that the record holds the cell instead. `auto-fix`. Confidence 4.

## Spec

Findings: 3. The worst issue is a why-clause mutant that the VR9 test does not kill alone.

- `internal/worktree/verb_runner_check_test.go:182`: a reader that takes the first 64-hex run passes VR9, because a real plan has no 64-hex text before the cell. VR14 kills that mutant, so the suite still catches it. `no-op`. Confidence 7.
- R2: the same literal as the Standards finding R2. `auto-fix`. Confidence 6.
- R5: `reviews/worktree-verb-runner.md:11-12`: the stage 2 log does not account for VR15. The stage 1 count reads 35, but 22 tests and 14 subtests give 36. Correct the log from the run evidence. This correction is evidence-only. `auto-fix`. Confidence 5.

## Coverage

Findings: 2. The worst issue is the untested record branch of the fingerprint reader.

- R6: `internal/worktree/verb_runner_test.go:192-203`: no test grades the absent rule or the conflict rule on a record line. A no-op `reset` plan writes `fingerprint=none` in its record. A swap that returns the first record cell was silent, and the restore reads `yes`. Add a real no-op `reset` record test that expects the no-fingerprint error, add a conflicting-record case, and raise `worktreeTestCount` for each added test. `auto-fix`. Confidence 8.
- `internal/worktree/verb_runner_test.go:95-108`: only VR17 grades a `joined:` entry. A swap of the `land` joins form was silent. Tickets 2 to 10 run stubbed joins values through each of these keys, so a wrong form fails those migrated tests. `no-op`. Confidence 6.

### Advice

- The `"usage: "` prefix in the VR1 expectations repeats each grammar's help composition. A test can read the grammar's help field where a grammar variable exists.
- The length check of 64 in the VR9 test restates the digest width.
- A shared record encoder in `toon` would give the producer and the reader one source.

### Command contribution

The implementation command did not contribute. The ticket told the author to build the fault as an outside test does, but its `Writes:` line did not hold the owning file. Ticket slicing owns that fix: the first ticket that needs a shared builder names it and writes its owning file.

## VR-C1 repair 1

A fresh `bench-writer` repair session, `vr-t1-repair-1`, ran on opus at medium effort, with a cap of 2 attempts. The session started at `c71d932b` and committed `3a318142` on a lane pass in the first attempt. The lane passed with 15 green checks, 1 check that does not apply, and 0 red checks. The source digest of the repair is `0dee07c82fe40e1420c5629289bd195d4e8ff5cd`, the tree of `3a318142` without this record file. This record is a second commit.

The repair changes only test files. The red and green route for each target follows:

- R1: a new builder, `addBrokenUnclaimedBranch`, owns the broken unclaimed ref and its reason comment. `TestCleanUnclaimedErrorRowRefusesTheSet` and `TestVerbResultFingerprintTreatsAPlaceholderAsAbsent` both call it, and the assertions of the clean test do not change. The builder is in `verb_runner_check_test.go`, not beside `addUnclaimedBranch`. In `clean_set_apply_test.go`, the builder made the file 417 lines, and the structure lane refused it at its budget of 400. Red: a swap in the builder that makes the ref name a commit failed both tests.
- R2: the `clean` row of VR1 renders `cleanInvocationError` into a buffer, and the expectation is that whole output. Red: a swap that sends the `clean` key to the reclaim entry failed the `clean` subtest.
- R3: the `show` and `build` rows read `worktreeShowGrammar.Cmd` and `buildGrammar.Cmd`. The `usageCommand` helper is deleted. Red: a swap that sends the `show` key to the build entry failed the `show` subtest.
- R4: the VR10 test asserts that the record holds the `fingerprint=` cell with the value, and does not assert its position. Red: a reader that cuts the cell at the first `0` failed the test.
- R5: the stage 1 and stage 2 lines of the author evidence now state what this record cannot reproduce. The stub source is not in the history.
- R6: two new tests read the record branch. `TestVerbResultFingerprintReadsANoOpRecordAsAbsent` reads a real no-op `reset` plan and expects the no-fingerprint error. `TestVerbResultFingerprintRefusesConflictingRecords` adds a second copy of a real `reset` record with a different fingerprint and expects a conflict error. A shared builder, `resetRunnerPlan`, gives both tests and VR10 a real plan. `worktreeTestCount` is now 688.

### Repair probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first two rows ran on the committed source, and the JSON payload holds the plan probe.

| File | Mutation | Test | Target | Verdict |
|---|---|---|---|---|
| `internal/worktree/verb_runner_test.go` | swap: `values = recordFingerprints(stdout)` to a return of the first record cell | both R6 tests | R6 | bit, 2 failed |
| `internal/worktree/verb_runner_test.go` | swap: `unapplicableFingerprint` to `"probe-placeholder"` | TestVerbResultFingerprintTreatsAPlaceholderAsAbsent | VR59 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `if value != values[0] {` to `if false {` | TestVerbResultFingerprintRefusesConflictingRecords | R6 | bit |
| `internal/worktree/verb_runner_check_test.go` | swap: `"hash-object", "-w", "tracked.txt"` to `"rev-parse", "HEAD"` | both R1 tests | R1 | bit, 2 failed |
| `internal/worktree/verb_runner_test.go` | swap: the `clean` entry to `ReclaimCommand` | TestVerbRunnerKeyReachesItsOwnVerb/clean | R2 | bit |
| `internal/worktree/verb_runner_test.go` | swap: the `show` entry to `BuildCommand` | TestVerbRunnerKeyReachesItsOwnVerb/show | R3 | bit |
| `internal/worktree/verb_runner_test.go` | swap: the record cell cut at `,` to a cut at `0` | TestVerbResultFingerprintReadsTheRecordCell | R4 | bit |

### Repair verification

The session ran each check on the source of `3a318142`, and each passed: `bench test --package ./internal/worktree`, the plan probe `1-probe`, `bench test --package ./internal/conformance`, and `bench structure --growth 0c95c944`. `TestPackageTestCountPin` passed at 688 tests, and `TestSerialSetStaysBelowTheCeiling` passed at the ceiling of 46. The two new tests call `t.Parallel()` and bind no environment. The digest of the `1-worktree` result is of the text that the session received, because that output did not spill to a file.

## VR-C1 chunk review, round 2

This round confirms repair 1. The frozen pair is base `0c95c9447c20189f3f2155719ef965bffc339856` and tip `d44d28ce6ff1964ffc4cb77d969cecdb881356eb`. The shared evidence is `sha256:565cf66f6d2860a4079406d195130860fd2993adebbac2089ebf7de385f56926`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, and each read the repair delta `179e2eb8..d44d28ce`. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 2: Standards 0, Spec 2, and Coverage 0. Every fold of R1 to R6 holds. Both Spec findings name spec text, not code, so they use no repair cycle. The consumed allowance stays at 1 of 2 repair cycles.

Plan commit `f3b6ad76` closed both Spec findings, but it changed the source past the reviewed tip. Commit `b02ef7a3` restored the reviewed spec bytes. Plan commit `51f0951e` then applied the same amendment inside VR-C1, and round 3 reviewed it.

The coordinator made three evidence-only corrections to this record. The repair verification entries named the full tree `bef4351236549fad5b2cef6e5958d23b685aabb8` as their source digest. The source digest excludes this record file, so the entries now name `0dee07c82fe40e1420c5629289bd195d4e8ff5cd`. The plan digest reads the plan at the reviewed tip `d44d28ce`. Each author and repair result digest hashed the whole command output, so each digest now hashes its embedded excerpt, as the record parser requires.

The chunk record for VR44 to VR46 follows. All 688 top-level tests pass, which is the base 664 plus 24 added tests. The SKIP set holds the two socket capability subtests. No existing test lost a failure call. The only edits to existing files are setup lines in `clean_unclaimed_test.go` and the count constant.

### Standards, round 2

Findings: 0. All five folds hold. The shared fault builder sits in the runner check file, because `clean_set_apply_test.go` is at its 400-line budget. The axis rates that placement a judgment call and not a binding defect.

### Spec, round 2

Findings: 2. The worst issue is a coverage map that did not name the two record-branch tests.

- R7: `specs/worktree-verb-runner/spec.md` had no row for the record-branch absent rule and conflict rule. Plan commit `51f0951e` added VR60 and VR61, and ticket 1 covers both. `auto-fix`. Confidence 5.
- R8: the VR1 row said that every expectation reads its usage constant, but the `show` and `build` rows read the grammar `Cmd` field. Plan commit `51f0951e` stated that source in VR1, in the ticket, and in the reader sweep. This change is non-behavioral, and the reviewer can veto it. `auto-fix`. Confidence 6.

### Coverage, round 2

Findings: 0. Each fold bit under a new probe site. The probes removed the placeholder rule on the record branch and broke the shared fault builder. They also moved the `show` and `clean` keys to other verbs and cut the record cell at the wrong byte. Each restore reads `yes`.

### Advice, round 2

- A new file can hold the whole unclaimed branch fixture family and free `clean_set_apply_test.go` from its budget.
- The VR1 `clean` row does not refuse an empty expectation. A check that the rendered refusal is not empty closes that gap.
- In the conflict test, a second `value` shadows the first one.

## VR-C1 verification rerun

A fresh `bench-writer` session, `vr-t1-verify-2`, ran on opus at medium effort, with a cap of 1 attempt. The session changed no code and no spec. Plan commit `51f0951e` added the coverage rows VR60 and VR61 after the last verification, so the ticket 1 checks ran again on that source. The source digest is `9b1948029d85ed2ec2df74290d084ae136196359`, the tree of `51f0951e` without this record file.

- `1-worktree`: `bench test --package ./internal/worktree` passed. The SKIP set holds the two socket capability subtests.
- `1-probe`: the swap of `unapplicableFingerprint` to `"probe-placeholder"` bit, with 1 failed test. The restore reads `yes`.

## VR-C1 chunk review, round 3

This round confirms plan commit `51f0951e` and the verification rerun. The frozen pair is base `0c95c9447c20189f3f2155719ef965bffc339856` and tip `ef2cd35f1052f7003c1f7647539d656be7585821`. The shared evidence is `sha256:a2e41f11506f672a638765f9167e9c3453cbde7d4ffcd22cacfbd62347e03db8`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, and each read the delta `d44d28ce..ef2cd35f`. That delta changes no test code. Only the Coverage axis ran a probe, and it left the tree clean.

The raw finding count is 5: Standards 4 and Spec 1. One Standards finding and the Spec finding name one binding fix, R9. Of the three Standards judgment items, the coordinator fixed one tense error in place, and the advice list below keeps the other two.

The round 2 prose of this record routed the amendment to VR-C2, but plan commit `51f0951e` applied it inside VR-C1. The coordinator corrected that prose, and each issuing axis reaffirmed the correction with a pass. The correction is evidence-only, so the consumed allowance stays at 1 of 2 repair cycles.

### Standards, round 3

Binding findings: 1. R9: `reviews/worktree-verb-runner.md:128,142,143` placed the amendment in VR-C2. One source per fact requires one account of where the amendment lives. `auto-fix`. Confidence 8. The axis reaffirmed the correction.

### Spec, round 3

Findings: 1. R7 and R8 are closed: VR60 and VR61 match their tests, and the VR1 wording matches the `show` and `build` rows. R9 is the same record defect as the Standards finding. `auto-fix`. Confidence 7. The axis reaffirmed the correction.

### Coverage, round 3

Findings: 0. VR60 and VR61 name real tests. The round 2 omission of the placeholder rule covers VR60. A swap of the two-fingerprints check bit VR61, and the restore reads `yes`. The chunk record for VR44 to VR46 is unchanged from round 2, because no test code changed.

### Advice, round 3

- The VR1 why-cell in the spec runs to 38 words and says "command name" where the ticket says "`Cmd` field". The VR-C2 enabling plan commit can shorten it.
- Ticket 1 says "a real no-op `reset` record", and VR60 says "a no-op `reset` plan whose record carries `none`". One term serves both.

## VR-C2 ticket 2 author evidence

Ticket 2 had a fresh `bench-writer` author, `vr-t2-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `371873a6` on the chunk base `ef2cd35f`. The author committed `c56f6107` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `c56f6107` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. The red for each row is the named regression probe, which bit before and after the change. The static rows read the tree at `c56f6107`.

- VR23: the VR23 command printed no line. The four helpers `runReset`, `runResetWith`, `resetFingerprint`, and `restoreFingerprint` are deleted.
- VR24: `restoreFixture` returns `restoredAssignment`, which `verb_fixture_test.go` declares. The tuple scan lists `restoreFixture` at the base and omits it at `c56f6107`.
- VR25: the verb form command over the seven reset files and `verb_fixture_test.go` printed no line.
- VR47: no reset file declares a `bytes.Buffer`, and no reset file calls a core reader or `checkVerbCall`.

The author reported these deviations and choices. The Spec axis grades each one.

- Each deleted helper checked the plan exit code. Each call site keeps that check as a `requireTest` line before `mustFingerprint`.
- `restoreFingerprint` also checked the `mode=restore` and `envelope=` cells. The predicate `isRestorePlanOf` holds that check once, and it runs no verb and reads no fingerprint.
- The `restoredAssignment.call` method builds a call value from the root and the home, and it runs no verb.
- `newOwnedAssignment` keeps its tuple until ticket 3, so its call sites build each call value inline.
- The `runMerge` call in `reset_repair_test.go` stays for ticket 7.
- No check required an edit to the five registry paths, so the diff leaves them unchanged.

### Probe verdicts

Each probe ran through `bench probe` with the same mutation, and each restore reads `yes`. The mutation in `internal/worktree/reset_apply.go` swaps `plan.action == "none" || fingerprint != plan.fingerprint` to `plan.action == "none"`. The run selects `TestReset`. The failing test set is the same before and after the change.

| Source | Verdict | Failed tests |
|---|---|---|
| `371873a6`, before the first edit | bit | `TestResetApplyRefusesAStalePlan`, `TestResetFingerprintTracksTheIndex`, `TestResetRestoreRefusesAStaleIndex`, `TestResetRestoreRefusesAStalePlan` |
| `c56f6107` | bit | `TestResetApplyRefusesAStalePlan`, `TestResetFingerprintTracksTheIndex`, `TestResetRestoreRefusesAStaleIndex`, `TestResetRestoreRefusesAStalePlan` |

### VR46 pre-check

The author counted the calls in each function of the seven reset files at `ef2cd35f` and at `c56f6107`. The count includes `t.Fatal`, `t.Fatalf`, `t.Error`, and `t.Errorf`. Most tests fail through `requireTest` and `mustNoError`, so the author also counted those two helpers.

| Count | `ef2cd35f` | `c56f6107` |
|---|---|---|
| Functions that both trees hold | 76 | 76 |
| Direct `t.Fatal` family calls | 2 | 2 |
| `requireTest` and `mustNoError` calls | 140 | 187 |
| `mustFingerprint` calls | 0 | 47 |

No test count fell. The four deleted helpers held 5 `requireTest` calls. The exit and cell checks move to each call site, and `mustFingerprint` replaces the found, non-empty, and non-`none` checks.

### Tuple scan

The scan applies the positional-tuple predicate of the spec to each `_test.go` file of `internal/worktree`. The author ran the program from the scratch area with its own `go.mod`, which holds `module tuplescan` and `go 1.22`.

```go
// Command tuplescan lists each positional fixture tuple in the _test.go files of one
// directory. A positional fixture tuple is a function that returns at least two values
// and no error, where one value is a Creation or a []Creation, or where it returns three
// or more values. A function that returns exactly a joins value and a probe is excluded.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func typeText(fset *token.FileSet, expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + typeText(fset, e.Elt)
		}
	case *ast.StarExpr:
		return "*" + typeText(fset, e.X)
	case *ast.SelectorExpr:
		return typeText(fset, e.X) + "." + e.Sel.Name
	case *ast.MapType:
		return "map[" + typeText(fset, e.Key) + "]" + typeText(fset, e.Value)
	}
	return fmt.Sprintf("%T", expr)
}

func main() {
	dir := os.Args[1]
	paths, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
	if err != nil {
		panic(err)
	}
	sort.Strings(paths)
	fset := token.NewFileSet()
	for _, path := range paths {
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			panic(err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Results == nil {
				continue
			}
			var types []string
			for _, field := range fn.Type.Results.List {
				count := len(field.Names)
				if count == 0 {
					count = 1
				}
				for i := 0; i < count; i++ {
					types = append(types, typeText(fset, field.Type))
				}
			}
			if len(types) < 2 {
				continue
			}
			hasError, hasCreation := false, false
			for _, typ := range types {
				hasError = hasError || typ == "error"
				hasCreation = hasCreation || typ == "Creation" || typ == "[]Creation"
			}
			if hasError {
				continue
			}
			if len(types) == 2 && (types[0] == "joins" || types[1] == "joins") && !hasCreation {
				continue
			}
			if hasCreation || len(types) >= 3 {
				fmt.Printf("%s:%d: %s (%s)\n", filepath.Base(path), fset.Position(fn.Pos()).Line, fn.Name.Name, strings.Join(types, ", "))
			}
		}
	}
}
```

At `ef2cd35f` the scan printed 35 lines. At `c56f6107` it printed the 32 lines below. The three lines that went are `runReset`, `runResetWith`, and `restoreFixture`. Each remaining fixture builder is in the tuple list of the spec, and each remaining run wrapper goes under VR41.

```text
clean_branch_test.go:20: unprovableLandedAssignment (string, Creation)
clean_discard_test.go:67: runDiscard (string, int, map[string][]string)
clean_discard_test.go:93: planAndApply (string, int, map[string][]string)
clean_landed_test.go:18: landedSetFixture (string, string, Creation, Creation, Creation)
clean_landed_test.go:32: runCleanup (string, string, int)
clean_landed_test.go:38: runCleanupWith (string, string, int)
clean_set_apply_test.go:26: removableSetFixture (string, string, []Creation, map[string]string)
clean_set_apply_test.go:236: retainedMemberFixture (string, string, Creation, Creation)
exec_test.go:253: execAtOwnedTarget (string, string, string, int)
land_effects_cleanup_test.go:26: foldLandingSibling (Creation, string)
land_effects_test.go:21: brokerChangingLanding (string, Creation, string, string, string)
land_effects_test.go:160: brokerDestinationFixture (string, Creation, string, string, string)
land_fixtures_test.go:29: publicLandingFixture (string, Creation, string, string, string, string)
land_fixtures_test.go:36: publicLandingFixtureAtHome (string, Creation, string, string, string)
land_fixtures_test.go:44: specLessLandingFixture (string, Creation, string, string, string, string)
land_fixtures_test.go:55: foldedLandingFixture (string, Creation, string, string, string, string, string)
land_fixtures_test.go:69: landingFixtureAtHome (string, Creation, string, string, string)
land_fixtures_test.go:195: ticketsOnlyLandingFixture (string, Creation, string, string, string, string)
land_freshness_test.go:134: redProspectiveGateLanding (string, Creation, string, string, string, string)
land_surface_test.go:14: landSurface (string, Creation, string, string)
land_surface_test.go:20: landIn (int, string, string)
live_binary_test.go:21: newResidueGuardFixture (string, Creation, string)
merge_test.go:23: mergeFixture (joins, string, string, string, []Creation)
merge_test.go:50: runMerge (int, string, string)
pool_reclaim_test.go:22: newReclaimPool (string, string, string)
pool_root_test.go:40: poolRootFixture (string, string, string)
reauthorize_test.go:264: reauthorizeFixture (string, Creation, string, string, string)
resume_test.go:526: newOwnedSubmoduleAssignment (string, Creation, string)
resume_test.go:545: newOwnedAssignment (string, Creation, string)
resume_test.go:552: newPendingAssignment (string, Creation, string)
unlanded_route_test.go:16: refusedUnlandedRelease (string, Creation, string, string)
worktree_test.go:750: runCreate (int, string, string)
```

### Verification

The author ran each check on the source of `c56f6107`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The package excerpt omits its two skip rows, and each skip is a unix socket capability skip.

`TestPackageTestCountPin` passed with `worktreeTestCount` at 688, and `TestSerialSetStaysBelowTheCeiling` passed at the ceiling of 46. The diff adds no `t.Setenv` call. The author also ran `bench test --package ./internal/conformance` and `bench structure --growth ef2cd35f`, and each passed.

## VR-C2 ticket 3 author evidence

Ticket 3 had a fresh `bench-writer` author, `vr-t3-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `8331ef2b` on the chunk base `ef2cd35f`. The author committed `1c202cd6` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `1c202cd6` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. The red for VR26 is the tuple scan at the base, and the regression probe bit before and after the change.

- VR26: `newOwnedAssignment`, `newPendingAssignment`, and `newOwnedSubmoduleAssignment` each return `ownedAssignment`, which `verb_fixture_test.go` declares. At `8331ef2b` the tuple scan printed 32 lines, and three of them were these builders. At `1c202cd6` the scan printed 29 lines, and none of them was one of these builders.

The author reported these choices. The Spec axis grades each one.

- Each call site that reads two or more parts reads one value `f`, and each read names its field. A call site that reads one part reads that field directly from the builder result.
- `restoredAssignment` embeds `ownedAssignment`, so one `call` method builds the call value from the root and the home.
- A reset call value that sets only the root, the home, and the arguments now comes from `f.call`. A call value that also sets `joins` keeps its literal, because a split would add a line.
- One test in `resume_test.go` builds two assignments in one scope. The first value is `o`, and the second value is `f`.
- No check required an edit to the five registry paths, so the diff leaves them unchanged.

### Ticket 3 probe verdicts

The probe used the same mutation and run as ticket 2, and each restore reads `yes`.

| Source | Verdict | Failed tests |
|---|---|---|
| `8331ef2b`, before the first edit | bit | `TestResetApplyRefusesAStalePlan`, `TestResetFingerprintTracksTheIndex`, `TestResetRestoreRefusesAStaleIndex`, `TestResetRestoreRefusesAStalePlan` |
| `1c202cd6` | bit | `TestResetApplyRefusesAStalePlan`, `TestResetFingerprintTracksTheIndex`, `TestResetRestoreRefusesAStaleIndex`, `TestResetRestoreRefusesAStalePlan` |

### Ticket 3 line counts and VR46 pre-check

Each over-budget file that ticket 3 writes has the same line count at `ef2cd35f` and at `1c202cd6`. The counts are `worktree_test.go` 1031, `resume_test.go` 598, `identity_component_test.go` 509, `lifecycle_test.go` 436, `exec_test.go` 432, and `ownership_test.go` 421. Only `verb_fixture_test.go` grew, by 6 lines. `bench structure --growth ef2cd35f` passed.

The author counted the `t.Fatal`, `t.Fatalf`, `t.Error`, and `t.Errorf` calls and the `requireTest` and `mustNoError` calls in each function of each package test file. The count was 2534 in 1027 functions at `8331ef2b` and at `1c202cd6`. No function count changed.

### Ticket 3 verification

The author ran each check on the source of `1c202cd6`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The package excerpt omits its two skip rows, and each skip is a unix socket capability skip. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

## VR-C2 ticket 4 author evidence

Ticket 4 had a fresh `bench-writer` author, `vr-t4-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `1b56289c` on the chunk base `ef2cd35f`. The author committed `d9f08fdf` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `d9f08fdf` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. The red for VR27 is the tuple scan at the base, and the regression probe bit before and after the change.

- VR27: `unprovableLandedAssignment` and `newResidueGuardFixture` each return `ownedAssignment`. `newReclaimPool` and `poolRootFixture` each return `poolFixture`. `verb_fixture_test.go` declares both types. At `d9f08fdf` the tuple scan printed 25 lines, and none of them was one of these four builders.

The author reported these choices. The Spec axis grades each one.

- `unprovableLandedAssignment` and `newResidueGuardFixture` return the root, the registration, and the home, which is the shape of `ownedAssignment`. Thus the author reused that type and declared no new type for them.
- `newReclaimPool` and `poolRootFixture` both return a root, a home, and a pool path, so they share the new type `poolFixture`. The pool of `newReclaimPool` is the parent of the pool keys. The pool of `poolRootFixture` is the pool root of one repository. The comment of each builder states which pool it prepares.
- A scratch tool rewrote each call site. The tool replaced each destructured call with one value `f`, and it replaced each resolved use of a part with its field. The tool changed no string literal.
- No check required an edit to the five registry paths, so the diff leaves them unchanged.

### Ticket 4 probe verdicts

The probe swapped the drift check in `pool_reclaim.go` for a check that is always false, and it ran the `Reclaim` tests. Each restore reads `yes`.

| Source | Verdict | Failed tests |
|---|---|---|
| `1b56289c`, before the first edit | bit | `TestReclaimApplyRefusesAFingerprintThePoolNoLongerMatches` |
| `d9f08fdf` | bit | `TestReclaimApplyRefusesAFingerprintThePoolNoLongerMatches` |

### Ticket 4 line counts and VR46 pre-check

The over-budget file `pool_reclaim_test.go` has 546 lines at `ef2cd35f` and at `d9f08fdf`. Only `verb_fixture_test.go` grew, by 8 lines. `bench structure --growth ef2cd35f` passed.

The author counted the `t.Fatal`, `t.Fatalf`, `t.Error`, and `t.Errorf` calls and the `requireTest` and `mustNoError` calls in each function of each package test file. The count was 2534 in 1027 functions at `1b56289c` and at `d9f08fdf`. No function count changed.

### Ticket 4 verification

The author ran each check on the source of `d9f08fdf`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The package excerpt omits its two skip rows, and each skip is a unix socket capability skip. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

## VR-C2 chunk review, round 1

The frozen pair is base `ef2cd35f1052f7003c1f7647539d656be7585821` and tip `80eac02e1f04599f86cbe40b6374592b75fee5c0`. The shared evidence is `sha256:b1dc8ea63c04cb761a949c260828bab28c4a72bd8586a40d0c2822c24b2935be`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 3: Standards 3, Spec 0, and Coverage 0. R12 is an evidence-only correction to this record, which the coordinator makes in this commit. So 2 repair targets remain, and they take repair cycle 1 of 2 for VR-C2.

The chunk record for VR44 to VR47 follows. All 688 top-level tests pass, and the SKIP set holds the two socket capability subtests. No test function lost a failure call, so no drop needs an account. The side-by-side read of the seven reset files found each base check on exit code, streams, refs, HEAD, status, and file bytes. No reset file holds a buffer pair for a verb call.

### Standards, VR-C2 round 1

Findings: 3. The worst issue is one field name for two directories.

- R10: `internal/worktree/verb_fixture_test.go:20-25` declares `poolFixture.pool`. That field holds the pool parent in `newReclaimPool` and one repository's pool in `poolRootFixture`. The STE rule gives one word to one thing. Split the type, or rename the field for one builder. `auto-fix`. Confidence 7.
- R11: `internal/worktree/reset_apply_test.go:71,85,99,104,152,271` holds six `verbCall` literals that restate how `ownedAssignment.call` maps the root and the home. `reset_restore_test.go` uses `f.call` and then sets the joins value, so the chunk has two ways to build the same call. `AGENTS.md` asks for one source per fact. `auto-fix`. Confidence 6.
- R12: this record said that the VR-C1 round 3 Standards axis found one finding, but its excerpt says four findings with one binding. This commit states both counts. `auto-fix`. Confidence 4.

The axis rates the `isRestorePlanOf` predicate and the embedded `ownedAssignment` as correct. The predicate returns a bool, reads no fingerprint, and splits no rows.

### Spec, VR-C2 round 1

Findings: 0. Rows VR23 to VR27 hold at the frozen tip. The VR23 command and the verb form command print no line, and the tuple scan omits all eight builders. The over-budget files stay at their base line counts, and every ticket commit stays inside its `Writes:` line.

### Coverage, VR-C2 round 1

Findings: 0. Four new probes bit. They disabled the recapture check, the off-branch guard, the staged index layer, and the envelope tip check. Each restore reads `yes`.

### Advice, VR-C2 round 1

- About 30 call sites run a reset plan, check its exit code, and read its fingerprint. A later change can let `mustFingerprint` refuse a nonzero exit.
- `requireResetRefusal` runs a verb and returns its stdout. The Enumerations list does not name it, but it has the shape of a run wrapper.
- The comment on `newResidueGuardFixture` still says that the builder returns the private home.
- `reset_repair_test.go` still calls `runMerge`, which ticket 7 removes.

### Repair route, VR-C2

R10 sits on the ticket 4 `Writes:` line, and R11 sits on `reset_apply_test.go`, which tickets 2 and 3 write. A plan commit adds `reset_apply_test.go` to the ticket 4 `Writes:` line. One fresh repair session for ticket 4 then repairs both. After that repair, the ticket 2 and ticket 3 authors rerun their verification at the final chunk source.

## VR-C2 repair 1

A fresh `bench-writer` repair session, `vr-t4-repair-1`, did repair cycle 1 of 2 for VR-C2 on opus at medium effort, with a cap of 2 attempts. The session started at `abfca3e2` on the chunk base `ef2cd35f`. It committed `53c499c9` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `53c499c9` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

- R10: `verb_fixture_test.go` now declares two types. `reclaimPoolFixture` holds the pool parent in its `pool` field, and `newReclaimPool` creates that directory. `repoPoolFixture` holds the pool root of one repository in its `poolRoot` field, and `poolRootFixture` computes that path but does not create it. The comment on each type states this. `pool_root_test.go` reads `f.poolRoot` in each place where it read `f.pool`.
- R11: the six calls in `reset_apply_test.go` that set a joins value now build the call with `f.call` and then set `call.joins`, as `reset_restore_test.go` does. No assertion changed.
- The comment on `newResidueGuardFixture` now states that the builder returns the repository, the owned registration, and the private home.

### Repair 1 probe verdicts

Each probe ran on the source of `53c499c9`, and each restore reads `yes`.

| Mutation | Verdict | Failed tests |
|---|---|---|
| The reset plan check in `reset_apply.go` ignores the fingerprint | bit | `TestResetApplyRefusesAStalePlan`, `TestResetFingerprintTracksTheIndex`, `TestResetRestoreRefusesAStaleIndex`, `TestResetRestoreRefusesAStalePlan` |
| The drift check in `pool_reclaim.go` is always false | bit | `TestReclaimApplyRefusesAFingerprintThePoolNoLongerMatches` |
| `TestResetApplyExitsThreeWithoutAnEnvelope` does not set `call.joins` | bit | `TestResetApplyExitsThreeWithoutAnEnvelope` |

The third probe is the repair session's own probe. It shows that the joins value still reaches the verb after the R11 change.

### Repair 1 verification

The session ran each check on the source of `53c499c9`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The two skips are unix socket capability skips. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

`pool_reclaim_test.go` has 546 lines, and `bench structure --growth ef2cd35f` passed. The tuple scan printed 25 lines, and none of them is one of the eight VR-C2 builders. The session counted the `t.Fatal` family calls and the `requireTest` and `mustNoError` calls in each test function of the package at `abfca3e2` and at `53c499c9`. The count was 2232 in 703 functions at each commit, and no function count dropped.

## VR-C2 ticket 2 verification rerun

The ticket 2 author ran `bench test --package ./internal/worktree` again on the final chunk source at `e8a24a1e`, after repair commit `53c499c9`. The run passed with the two unix socket capability skips. The JSON payload holds the result as `vr-c2-2-worktree-r2`.

## VR-C2 ticket 3 verification rerun

The ticket 3 author ran `bench test --package ./internal/worktree` again on the final chunk source at `35f6f1fd`, after repair commit `53c499c9`. The run passed with the two unix socket capability skips. The JSON payload holds the result as `vr-c2-3-worktree-r2`.

## VR-C2 chunk review, round 2

This round confirms repair 1 of VR-C2. The frozen pair is base `ef2cd35f1052f7003c1f7647539d656be7585821` and tip `f03e7fb9b46e5cdbe0ae55d94f171b209fe27447`. The shared evidence is `sha256:99f35c57a149ab5971deb83e0328f3b0709df8889832aa12e198d85534e4bfc9`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, and each read the repair delta `80eac02e..f03e7fb9`. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 2: Standards 1 and Coverage 1. R10 and R11 hold. The Standards axis found that the R12 correction still disagreed with the advice list. The coordinator corrected that prose, and the issuing axis reaffirmed the correction with a pass. The consumed allowance stays at 1 of 2 repair cycles.

The Coverage axis found a weak lock check that the chunk base already had. The spec asks this migration to keep assertions, not to strengthen them. So this record keeps that item as advice, and the learnings inbox holds it.

The chunk record for VR44 to VR46 follows. The package run at the tip passes, and the SKIP set holds the two socket capability subtests. No test function in the repair delta lost a failure call. Each ticket reran its verification at the final source `b7f16bc7bc9bb2827e92686bc09084c902334231`. The VR-C2 plan commits added author assignments and wording only. So the payload maps each chunk ID of the VR-C1 plan to the same ID in the current plan.

### Advice, VR-C2 round 2

- The lock check in `reset_apply_test.go` near line 104 passes without its joins value. A positive check that the plan run reached the recorder would close that gap.
- Six sites build a call and then set the joins value. A method on the fixture can collapse that pattern.
- Several reclaim tests pass the parent of the pool where they mean the home.

## VR-C3 ticket 5 author evidence

Ticket 5 had a fresh `bench-writer` author, `vr-t5-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `393ed757` on the chunk base `f03e7fb9`. The author committed `b5e755a9` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `b5e755a9` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. For each row, the red is the scan output at the base, and the green is an empty scan at `b5e755a9`. Both regression probes bit before and after the change, with equal failing sets.

- VR28: at `393ed757` the VR28 command printed 36 lines over the ticket files. At `b5e755a9` it printed no line.
- VR29: at `393ed757` the tuple scan printed 25 lines, and the four cleanup builders were among them. At `b5e755a9` it printed 17 lines, and none of them is one of the four builders.
- VR30: at `393ed757` the verb form command printed 59 lines over the ticket files. At `b5e755a9` it printed no line.

The author reported these choices. The Spec axis grades each one.

- `landedSetFixture` returns `landedSet`, and `removableSetFixture` returns `removableSet`. `retainedMemberFixture` returns `retainedMemberSet`, and `refusedUnlandedRelease` returns `refusedRelease`. `verb_fixture_test.go` declares the four types.
- `verb_fixture_test.go` declares `repoHome`, which is a root and a home. Its method `call` builds a verb call, and its method `callWith` also sets the joins value. `ownedAssignment`, `reclaimPoolFixture`, and the four new types embed `repoHome`, so one method builds each call.
- `verb_fixture_test.go` also declares `cleanupTable`, `textRow`, `textRows`, and `textRowsBy`. They read the rows that `mustRows` decoded as text cells. They run no verb and parse no rendered text.
- `requireRowFor` fails a test when no decoded row names a target. `requireCleanupColumns` holds the column check of the deleted `unclaimedVerdicts`. `reclaimApplyHead` is the one source of the reclaim apply command.
- Each site that read `reclaimFingerprint` now reads `mustFingerprint`. It also keeps the check that the plan prints the apply command with that fingerprint.
- Each former `mustReclaim` and `runDiscard` call keeps its empty stderr check at the call site. Each former `runResume` call keeps its exit check and its `chdir`.
- `TestRetiringVerbSpillsToPrimary` writes the release streams to the response owner after the release returns. The spill starts before the release, so a spill in the retired scope still fails the test.
- No check required an edit to the five registry paths, so the diff leaves them unchanged.

### Ticket 5 probe verdicts

Each probe ran on the named source, and each restore reads `yes`.

| Mutation | Source | Verdict | Failed tests |
|---|---|---|---|
| The set apply in `clean_set.go` skips its fingerprint match | `393ed757`, before the first edit | bit | `TestCleanSetPreexistingDrift/explicit`, `TestCleanSetSpentPlan`, `TestCleanSetStaleReplanAction/explicit` |
| The set apply in `clean_set.go` skips its fingerprint match | `b5e755a9` | bit | `TestCleanSetPreexistingDrift/explicit`, `TestCleanSetSpentPlan`, `TestCleanSetStaleReplanAction/explicit` |
| The receipt checkpoint check in `resume.go` is always false | `393ed757`, before the first edit | bit | `TestResumeSupersedesPlannedReceiptOverAbsentTarget/present_target` |
| The receipt checkpoint check in `resume.go` is always false | `b5e755a9` | bit | `TestResumeSupersedesPlannedReceiptOverAbsentTarget/present_target` |

The first probe ran the `Clean` tests, and the second probe ran the `Resume|Clean` tests.

### Ticket 5 line counts and VR46 pre-check

From `f03e7fb9` to `b5e755a9`, `resume_test.go` went from 598 to 585 lines. `pool_reclaim_test.go` went from 546 to 533 lines, and `lifecycle_test.go` went from 436 to 433 lines. `clean_discard_test.go` has 396 lines and `clean_set_apply_test.go` has 399 lines, so both stay in their budget. `verb_fixture_test.go` grew from 41 to 126 lines. `bench structure --growth f03e7fb9` passed.

The author counted the `t.Fatal`, `t.Fatalf`, `t.Error`, and `t.Errorf` calls and the `requireTest` and `mustNoError` calls in each function of each package test file. The count was 2534 in 1031 functions at `393ed757`, and 2541 in 1032 functions at `b5e755a9`. Three test functions lost one call each, and each drop is a must-form replacement.

| Test | Base | Tip | Replacement |
|---|---|---|---|
| `TestCleanExplicitSetPlan` | 9 | 8 | `mustFingerprint` replaces the check that two 64-hex runs agree |
| `TestCleanCommandAcceptsTheAbsolutePathThatPathPrints` | 7 | 6 | `mustFingerprint` replaces the check that the plan carries a fingerprint |
| `TestCleanSetHostileOperand` | 7 | 6 | `mustNoFingerprint` replaces the check that each fingerprint cell is `none` |

The reader in `mustNoFingerprint` also accepts an empty cell, so the last check is weaker than the base check. The nine deleted helpers lost their calls, and each check moved to the callers or into a must form. Other functions gained calls from the exit and stderr checks that the deleted wrappers made.

### Ticket 5 VR47 and verification

No ticket file declares an output buffer pair for a verb call. The three pairs left in `subshell_test.go` feed `subshellAt`, which the spec does not handle.

The author ran each check on the source of `b5e755a9`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The two skips are unix socket capability skips. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

## VR-C3 ticket 6 author evidence

Ticket 6 had a fresh `bench-writer` author, `vr-t6-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `a50d414a` on the chunk base `f03e7fb9`. The author committed `a531f85b` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `a531f85b` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. For each row, the red is the scan output at the start, and the green is an empty scan at `a531f85b`.

- VR31: at `a50d414a` the VR31 command printed 4 lines over the ticket files. The lines were the `runCreate` and `execAtOwnedTarget` declarations and the two inline fingerprint matches. At `a531f85b` it printed no line.
- VR32: at `a50d414a` the verb form command printed 102 lines over the ticket files. At `a531f85b` it printed no line.

The author reported these choices. The Spec axis grades each one.

- Each migrated call uses `runVerb` with a `repoHome` call, a fixture `call`, or a `verbCall` value. A call with a joins value uses `callWith`. No new builder exists.
- `requireCreateFromRefusal` takes the verb result in place of the exit code and the two streams.
- `selectedRows` decoded the list output and read one table, which is what `mustRows` does. The author deleted it, and the seven readers call `mustRows` with the new `selectedTable` constant.
- Seven failure messages named a verb entry in call form, for example `ListCommand(%q)`. The verb form command matched those message texts, so the author changed each one to name the verb, for example `list %q`. No assertion changed.
- The `runCreate` row in `capture/restructure-backlog.md` now names `requireCreateFromRefusal` as the create section anchor. The author also updated the line count and the other anchors in that row to the lines at `a531f85b`.
- No check required an edit to the five registry paths, so the diff leaves them unchanged. `childFailure` stays, and the `runMerge` and `mergeFixture` uses stay for ticket 7.

### Ticket 6 probe verdicts

Each probe omitted the `PWD` entry that `execEnv` appends in `exec.go`. Each probe ran the `Exec|PWD` tests, and each restore reads `yes`.

| Source | Verdict | Failed tests |
|---|---|---|
| `a50d414a`, before the first edit | bit | `TestExecEnvironmentContainsOneCanonicalPWD`, `TestExecPWDMatchesChildDirectory/absent`, `TestExecPWDMatchesChildDirectory/inherited`, `TestExecPWDMatchesChildDirectory/repeated_overrides` |
| `a531f85b` | bit | `TestExecEnvironmentContainsOneCanonicalPWD`, `TestExecPWDMatchesChildDirectory/absent`, `TestExecPWDMatchesChildDirectory/inherited`, `TestExecPWDMatchesChildDirectory/repeated_overrides` |

The two failing sets are equal.

### Ticket 6 line counts and VR46 pre-check

From `f03e7fb9` to `a531f85b`, `worktree_test.go` went from 1031 to 985 lines, and `exec_test.go` went from 432 to 421 lines. `list_actions_test.go` went from 395 to 393 lines, so it stays in its budget. `bench structure --growth f03e7fb9` passed.

The author counted the `t.Fatal`, `t.Fatalf`, `t.Error`, `t.Errorf`, `t.FailNow`, `t.Fail`, `requireTest`, and `mustNoError` calls in each top-level test of the package. The count was 2243 in 703 tests at `a50d414a`, and 2241 in 703 tests at `a531f85b`. Two tests lost one call each, and each drop is a must-form replacement.

| Test | Base | Tip | Replacement |
|---|---|---|---|
| `TestCleanApplyAcceptsAFingerprintPrefix` | 5 | 4 | `mustFingerprint` replaces the check that the plan carries a fingerprint |
| `TestCleanDropsTheCensusRecords` | 4 | 3 | `mustFingerprint` replaces the check that the plan carries a fingerprint |

The two deleted helpers lost their calls. Each exit check and each stream check that `execAtOwnedTarget` and `runCreate` returned to a caller stays at that caller.

### Ticket 6 VR47 and verification

No ticket file declares an output buffer pair for a verb call. The buffers left in `exec_test.go` feed `runWorktreeChild`, which is the exec child runner and not a verb form. The buffer in `show_test.go` takes the stderr of a direct `git cat-file`, and the buffer in `worktree_test.go` is the advisory writer of a joins value.

The author ran each check on the source of `a531f85b`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The two skips are unix socket capability skips. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

## VR-C3 ticket 7 author evidence

Ticket 7 had a fresh `bench-writer` author, `vr-t7-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `5b61f38f` on the chunk base `f03e7fb9`. The author committed `c6252480` on a lane pass in the first attempt, and then committed this record alone. The build preflight on `c6252480` reported 13 green checks, 2 checks that do not apply, and 0 red checks.

The ticket is a pure migration, so no row has a new test. For each row, the red is the scan output at the start, and the green is the scan output at `c6252480`.

- VR33: at `5b61f38f` the VR33 command printed 1 line, the `runMerge` declaration in `merge_test.go`. At `c6252480` it printed no line.
- VR34: at `5b61f38f` the tuple scan listed `mergeFixture` and `reauthorizeFixture`. At `c6252480` the scan lists neither builder. The scan still lists the landing builders that tickets 8 and 9 own.
- VR35: at `5b61f38f` the verb form command printed 13 lines over the ticket files. At `c6252480` it printed 1 line, the `LandCommand` call in `delegated_integration_test.go` that ticket 9 moves.

The author reported these choices. The Spec axis grades each one.

- `mergeFixture` returns a `mergeSet` value. The value carries a `repoHome`, the joins value, the tally file, and the registrations. Its `merge` method builds the joins form call through `callWith`. Thus each merge call is one runner line, and each `mergeLane` and `mergeReconcile` stub still runs.
- `reauthorizeFixture` returns a `reauthorizeSet` value. The value embeds `ownedAssignment` and adds the reviewed base and tip. The rollback test passes its stubbed joins value through `callWith`.
- `requireMergeRefusal` and `requireMergeLaneRefusal` take the verb result in place of the exit code and the streams. `requireMergeUnchanged` takes the merge set, the target, and the previous tip. Each assertion in the three helpers stays.
- The `fold` method in `delegated_integration_test.go` builds its call from a `repoHome` value with the journey's joins value. The `LandCommand` call in that file stays for ticket 9.
- In `reset_repair_test.go`, the two reset calls use the merge set's `call` method in place of a `verbCall` literal.
- One failure message named `ReauthorizeCommand` in call form. The author changed it to `reauthorize %q`, as ticket 6 did. No assertion changed.
- No check required an edit to the five registry paths, so the diff leaves them unchanged. `mergedRecord` stays.

### Ticket 7 probe verdicts

Each probe ran the `Merge` tests, and each restore reads `yes`. The first probe drops the conflict paths from the merge conflict refusal. The second probe makes the unreadable fingerprint refusal return 0.

| Probe | Source | Verdict | Failed tests |
|---|---|---|---|
| Conflict paths | `5b61f38f`, before the first edit | bit | `TestMergeRefusalEscapesAControlBytePath`, `TestMergeRefusesAConflictingNonCapturePath` |
| Conflict paths | `c6252480` | bit | `TestMergeRefusalEscapesAControlBytePath`, `TestMergeRefusesAConflictingNonCapturePath` |
| Unreadable fingerprint | `5b61f38f`, before the first edit | silent | none |
| Unreadable fingerprint | `c6252480` | silent | none |

For each probe, the two failing sets are equal. The second probe is silent at the base, so its silent result at the tip is not a loss.

### Ticket 7 line counts and VR46 pre-check

| File | `f03e7fb9` | `5b61f38f` | `c6252480` |
|---|---|---|---|
| `merge_test.go` | 932 | 932 | 922 |
| `worktree_test.go` | 1031 | 985 | 985 |

Each file stays at or below its line count at `f03e7fb9`. `bench structure --growth f03e7fb9` passed.

The author counted the `t.Fatal`, `t.Fatalf`, `t.Error`, `t.Errorf`, `t.FailNow`, `t.Fail`, `requireTest`, and `mustNoError` calls in each top-level test of the package. The count was 2241 in 703 tests at `5b61f38f`, and 2241 in 703 tests at `c6252480`. No test lost a call, so the VR46 table has no row. The deleted `runMerge` helper held no failure call.

### Ticket 7 VR47 and verification

No ticket file declares an output buffer pair for a verb call outside the excepted landing call. In `delegated_integration_test.go`, one buffer pair feeds the `LandCommand` call that ticket 9 moves. The other pair feeds `gate.RunCommand`, which is not a worktree verb. The buffer in `worktree_test.go` is the advisory writer of a joins value.

The author ran each check on the source of `c6252480`, and each passed. `bench test --package ./internal/worktree` passed, and the JSON payload holds the result. The two skips are unix socket capability skips. `worktreeTestCount` stays at 688, and the serial ceiling stays at 46.

## VR-C3 chunk review, round 1

The frozen pair is base `f03e7fb9b46e5cdbe0ae55d94f171b209fe27447` and tip `ab9d305b314690959e5f68d3c751ff90169b8f29`. The shared evidence is `sha256:217102f69e8b2a67ce032c8e3fb64c5271c1238b82d142301f75de4a3f560a10`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 6: Standards 3, Spec 2, and Coverage 1. R15 and R17 are evidence-only corrections to this record, which this commit makes. R16 is a spec text correction, and R14 and R18 need a code repair. So 3 repair targets remain, and they take repair cycle 1 of 2 for VR-C3.

The chunk record for VR44 to VR47 follows. The package run passes with 688 top-level tests, and the SKIP set holds the two socket capability subtests. Five tests lost one failure call each, and the axes read each test side by side with its base form.

The accepted drops under VR46 follow. Four tests replaced a fingerprint check with `mustFingerprint`: `TestCleanExplicitSetPlan`, `TestCleanCommandAcceptsTheAbsolutePathThatPathPrints`, `TestCleanApplyAcceptsAFingerprintPrefix`, and `TestCleanDropsTheCensusRecords`.

Eight checks for no 64-hex text in stdout became `mustNoFingerprint`, which reads only the two fingerprint cells. The spec allows that replacement. Those sites are `clean_set_test.go:122` and `:340`, `clean_discard_test.go:137`, `:157`, and `:169`, `clean_discard_transaction_test.go:148`, `clean_classes_test.go:205`, and `clean_unclaimed_test.go:318`. The shared-value check in `clean_landed_test.go:64-67` no longer asserts 64 hex characters, because `mustFingerprint` reads the cell. The drop in `TestCleanSetHostileOperand` is not accepted; R14 restores it.

### Standards, VR-C3 round 1

Findings: 3. The worst issue is a record that called a weaker check narrower.

- R17: `reviews/worktree-verb-runner.md:564` called the `mustNoFingerprint` check narrower, and line 585 counted five readers where seven exist. This commit corrects both. `auto-fix`. Confidence 7.
- R18: `verb_fixture_test.go:109` and `list_selected_test.go:94` declare the table names `worktree_cleanup` and `selected`, but `verb_runner_check_test.go` still spells both literals. The repair makes the runner tests read the two constants. The production literals stay, because this spec changes no production file, and the ideas inbox holds that pair. `ask-user` for the production pair; `auto-fix` for the test side. Confidence 5.

The axis found the ticket 5 helpers clean: each reads rows that `mustRows` already decoded, and none runs a verb or splits rendered text.

### Spec, VR-C3 round 1

Findings: 2. The worst issue is a record that logged one looser check where nine exist.

- R15: the ticket 5 record named only one looser no-fingerprint check. This commit logs each accepted drop above. `auto-fix`. Confidence high.
- R16: `specs/worktree-verb-runner/spec.md:401` and ticket 5 say that `mustFingerprint` replaces the `cleanupRowsField` read in `clean_set_command_test.go`. That read checked that every cell is `none`, so `mustFingerprint` must fail there. The tree uses `mustNoFingerprint`, which the fingerprint rules require. A plan commit corrects the spec text, and the reviewer can veto it. `auto-fix`. Confidence high.

Rows VR28 to VR35 hold. Each review-owned command prints no line, except the one `LandCommand` line that VR35 allows.

### Coverage, VR-C3 round 1

Findings: 1. The worst issue is a weakened refusal test.

- R14: `internal/worktree/clean_set_command_test.go:157`: at the base, every fingerprint cell had to equal `none`. Now `mustNoFingerprint` also accepts an empty cell. A probe set the placeholder constant to the empty string, and no clean test failed. Restore a check that each cell reads `none`. `auto-fix`. Confidence high.

Three new probes bit on the selected list, the show stream route, and the inside-tree release. A probe of the reclaim advertisement stayed silent, and the base has the same gap.

### Advice, VR-C3 round 1

- `runVerb` runs inside goroutines in two concurrency tests. If its `t.Fatalf` fires there, the test can hang.
- `delegatedJourneyJoins` repeats the tally lane that `mergeFixture` builds.
- About 40 failure messages name a verb entry, and eight name the verb.
- The record counted 703 tests in two places, where the pin is 688. That count includes subtests or helpers.

### Repair route, VR-C3

R14 and R18 go to one fresh repair session for ticket 5. A plan commit adds `verb_runner_check_test.go` to the ticket 5 `Writes:` line and corrects the spec text for R16. After the repair, the ticket 6 and ticket 7 authors rerun their verification at the final chunk source.

```bench-review-record
{
  "version": 2,
  "spec": "specs/worktree-verb-runner/spec.md",
  "plan_digest": "sha256:fc14637963d16b0fca7a4058f4eac565d31a02cb563e2fcdc1331d0547ba5c47",
  "implementation_session": "",
  "chunks": [
    {
      "id": "VR-C1",
      "base": "0c95c9447c20189f3f2155719ef965bffc339856",
      "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
      "plan_digest": "sha256:095153eee4b7547b074ceddeb588e1a6a77a4e4663f33dd70ed0b00e55144d39",
      "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
      "acceptance_rows": [
        "VR1",
        "VR2",
        "VR3",
        "VR4",
        "VR5",
        "VR6",
        "VR7",
        "VR8",
        "VR9",
        "VR10",
        "VR11",
        "VR12",
        "VR13",
        "VR14",
        "VR15",
        "VR16",
        "VR17",
        "VR18",
        "VR19",
        "VR20",
        "VR21",
        "VR22",
        "VR59",
        "VR60",
        "VR61"
      ],
      "verification": [
        {
          "id": "vr-c1-1-worktree-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/1-worktree@137fb303",
            "digest": "sha256:e0a99d19673fa071a4aab6ce31a04e581a6233c21ef1905353526dc0ec9e9c6e",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,51282\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c1-1-probe-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/1-probe@137fb303",
            "digest": "sha256:cb76d16fbf406599002d07b0d92643b9337e8419f0362e61c16114efdfd45de0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,43"
          },
          "requirement": "1-probe",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-author-20261001/1-probe@137fb303",
              "digest": "sha256:cb76d16fbf406599002d07b0d92643b9337e8419f0362e61c16114efdfd45de0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,43"
            }
          }
        },
        {
          "id": "vr-c1-author-probe-vr17-joins-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/vr17-joins-probe@137fb303",
            "digest": "sha256:52ab180cc662e1d83145275c16aae2873f295ae08b8395e8fcf063087facffc7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbRunnerPassesTheJoinsValue,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,60"
          },
          "requirement": "author-probe-VR17-joins",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'if call.joins != nil {' --with 'if false {' --package ./internal/worktree --run TestVerbRunnerPassesTheJoinsValue",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-author-20261001/vr17-joins-probe@137fb303",
              "digest": "sha256:52ab180cc662e1d83145275c16aae2873f295ae08b8395e8fcf063087facffc7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbRunnerPassesTheJoinsValue,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,60"
            }
          }
        },
        {
          "id": "vr-c1-1-worktree-r2",
          "performer": "claude:bench-writer/vr-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-repair-1-20261001/1-worktree@3a318142",
            "digest": "sha256:72cb722cbb113d863e4c1bf0439bc8e8c0c50a0b6739059a5f5af238f7370600",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,50325\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c1-1-probe-r2",
          "performer": "claude:bench-writer/vr-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-repair-1-20261001/1-probe@3a318142",
            "digest": "sha256:d4c62cc4f91e4a5a9a0f48e10c988ef86c408c95989af5576fede4b1d6e87bda",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,36"
          },
          "requirement": "1-probe",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-repair-1-20261001/1-probe@3a318142",
              "digest": "sha256:d4c62cc4f91e4a5a9a0f48e10c988ef86c408c95989af5576fede4b1d6e87bda",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,36"
            }
          }
        },
        {
          "id": "vr-c1-1-worktree-r3",
          "performer": "claude:bench-writer/vr-t1-verify-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-verify-2-20261001/1-worktree@51f0951e",
            "digest": "sha256:5520df147faeccab0adfafbe80226e4f8bde0fff627ed689ee9fc7b7cf751b25",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,49455\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c1-1-probe-r3",
          "performer": "claude:bench-writer/vr-t1-verify-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-verify-2-20261001/1-probe@51f0951e",
            "digest": "sha256:6689340540929a0f4d1d45b2e180f465fe41003a75e8ee8cdccd8f01214d0e78",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,33"
          },
          "requirement": "1-probe",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-verify-2-20261001/1-probe@51f0951e",
              "digest": "sha256:6689340540929a0f4d1d45b2e180f465fe41003a75e8ee8cdccd8f01214d0e78",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,33"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "vr-c1-standards-r1",
          "performer": "claude:bench-reviewer/vr-c1-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards@4479b2fa",
            "digest": "sha256:19f02fddf42f5f19396b117b6a2e9a6b6ef8edf11c1b72e04404a52e97b6b28f",
            "excerpt": "Standards: 4 findings. Worst: the VR59 test pastes a second copy of the faulted-unclaimed fixture and its comment from clean_unclaimed_test.go."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R1",
            "R2",
            "R3",
            "R4"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-spec-r1",
          "performer": "claude:bench-reviewer/vr-c1-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec@4479b2fa",
            "digest": "sha256:6f57a4c2208052b9c43203b4c35da81b90400c0409e5d68580a15592d70382c3",
            "excerpt": "Spec: 3 findings. Worst: the VR9 first-64-hex mutant survives its own test; VR14 kills it at suite level. All 22 planned tests exist with exact names."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R2",
            "R5"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-coverage-r1",
          "performer": "claude:bench-reviewer/vr-c1-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-coverage@4479b2fa",
            "digest": "sha256:b054f944c1d41ef4c3d785f4257ce6521b29fb826df474d43a963f49e4013569",
            "excerpt": "Coverage: 2 findings. Worst: the record branch of readVerbFingerprint has no test for the absent rule or the conflict rule; a swap that returns the first record cell is silent."
          },
          "axis": "Coverage",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R6"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-standards-r2",
          "performer": "claude:bench-reviewer/vr-c1-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards-r2@d44d28ce",
            "digest": "sha256:02c08c774b1445068209559be3c4397268999c530c936632f712774dd155e382",
            "excerpt": "Standards: 0 findings. All five folds hold; the shared fault builder in the runner check file is a judgment call, not a binding defect."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-standards-r1"
          ]
        },
        {
          "id": "vr-c1-spec-r2",
          "performer": "claude:bench-reviewer/vr-c1-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec-r2@d44d28ce",
            "digest": "sha256:333dfc09f314f4275d1726d255de0f5ae78f45ede30466295b171d7182d8a1d5",
            "excerpt": "Spec: 2 findings on spec text. Folds R2, R5, R6 and the plan commit hold; the R6 tests lacked coverage rows, and VR1 wording names only the usage constant."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [
            "R7",
            "R8"
          ],
          "supersedes": [
            "vr-c1-spec-r1"
          ]
        },
        {
          "id": "vr-c1-coverage-r2",
          "performer": "claude:bench-reviewer/vr-c1-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-coverage-r2@d44d28ce",
            "digest": "sha256:a9eacdfca84623ee5f8281963a2cf3e4333538d549aea8f4dabf925f6464fe75",
            "excerpt": "Coverage: 0 findings. All folds hold, every probe bit and restored, 688 top-level tests pass, and the skip set is the two socket subtests."
          },
          "axis": "Coverage",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-coverage-r1"
          ]
        },
        {
          "id": "vr-c1-standards-r3",
          "performer": "claude:bench-reviewer/vr-c1-standards-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards-r3@ef2cd35f",
            "digest": "sha256:95c91bb50538e1abe32bfb54592c2b749a03c41a964247a833fc627985927a91",
            "excerpt": "Standards: 4 findings, 1 binding. Worst: the round 2 record placed the VR60 and VR61 amendment in VR-C2, but plan commit 51f0951e put it in VR-C1."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "finding_ids": [
            "R9"
          ],
          "supersedes": [
            "vr-c1-standards-r2"
          ]
        },
        {
          "id": "vr-c1-spec-r3",
          "performer": "claude:bench-reviewer/vr-c1-spec-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec-r3@ef2cd35f",
            "digest": "sha256:eaea315b7fdc944a892d3dc17c6e04cd5e46ae1a020e90b61f1b632b81725b74",
            "excerpt": "Spec: 1 finding. R7 and R8 are closed; the record still routed the amendment to VR-C2."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "finding_ids": [
            "R9"
          ],
          "supersedes": [
            "vr-c1-spec-r2"
          ]
        },
        {
          "id": "vr-c1-coverage-r3",
          "performer": "claude:bench-reviewer/vr-c1-coverage-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-coverage-r3@ef2cd35f",
            "digest": "sha256:090d8a9500dcd3c629503251722907e5de457e3be849c12f72ea1c945b43678b",
            "excerpt": "Coverage: 0 findings. Test code is unchanged since round 2; VR60 and VR61 name real tests, and the VR61 swap probe bit and restored."
          },
          "axis": "Coverage",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-coverage-r2"
          ]
        },
        {
          "id": "vr-c1-standards-r3-reaffirm",
          "performer": "claude:bench-reviewer/vr-c1-standards-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards-r3/reaffirm@ef2cd35f",
            "digest": "sha256:4d3ca07d628f48a0eff4da8517babfaf98a0ca1bc6df79ddf375e2560cb487b4",
            "excerpt": "Standards reaffirm: pass. Record lines 128, 142 and 143 now name 51f0951e and VR-C1, which closes the binding finding."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-standards-r3"
          ]
        },
        {
          "id": "vr-c1-spec-r3-reaffirm",
          "performer": "claude:bench-reviewer/vr-c1-spec-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b1948029d85ed2ec2df74290d084ae136196359",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec-r3/reaffirm@ef2cd35f",
            "digest": "sha256:6bb7953d1db89b0279f6df23466c12d8879cab66c031e901da5f9f244447f79c",
            "excerpt": "Spec reaffirm: pass. Record lines 128, 142 and 143 now match the git log, so the finding is closed and the other verdicts stand."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-spec-r3"
          ]
        }
      ]
    },
    {
      "id": "VR-C2",
      "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
      "tip": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
      "plan_digest": "sha256:fc14637963d16b0fca7a4058f4eac565d31a02cb563e2fcdc1331d0547ba5c47",
      "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
      "acceptance_rows": [
        "VR23",
        "VR24",
        "VR25",
        "VR26",
        "VR27"
      ],
      "verification": [
        {
          "id": "vr-c2-2-worktree-r1",
          "performer": "claude:bench-writer/vr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1c354a766a4d09d2673b5151f806eed82f378a14",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t2-author-20261001/2-worktree@c56f6107",
            "digest": "sha256:929611328f629f2c705026835eb872158f9c91dc7d0bd0d0cf70ab51d5d15943",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,50660\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c2-3-worktree-r1",
          "performer": "claude:bench-writer/vr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6e41a0cc15f3366a8dac08ef745a6ab39d1f59e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t3-author-20261001/3-worktree@1c202cd6",
            "digest": "sha256:6f671bc8e81d4b1cc7f709d80f29d761ab5e98f01aee96ab4509415334ce598d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,51143\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "3-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c2-4-worktree-r1",
          "performer": "claude:bench-writer/vr-t4-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "36e7b63fe9a09b930fd31563be5ea77eadf841f4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t4-author-20261001/4-worktree@d9f08fdf",
            "digest": "sha256:152b3e7832dee3346d958f133579e16c04106e00fc47dd07f52220938552966b",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,51513\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "4-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c2-4-worktree-r2",
          "performer": "claude:bench-writer/vr-t4-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t4-repair-1-20261001/4-worktree@53c499c9",
            "digest": "sha256:df38e02b8a91ac993227b17ec008de9cd2df8336326dcd8d1c26acf0fccc6345",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,49252\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "4-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c2-2-worktree-r2",
          "performer": "claude:bench-writer/vr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t2-author-20261001/2-worktree@e8a24a1e",
            "digest": "sha256:b7f0ed0485343c12c4f4ba12a06ce8d998246fb0f26aecce5cec56cbe5b8cfa7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,49388\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c2-3-worktree-r2",
          "performer": "claude:bench-writer/vr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t3-author-20261001/3-worktree@35f6f1fd",
            "digest": "sha256:83bf934371342760370e11f417971015bdc0fb27e1133db919d5fc7b69fab42c",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,49219\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "3-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "vr-c2-standards-r1",
          "performer": "claude:bench-reviewer/vr-c2-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "36e7b63fe9a09b930fd31563be5ea77eadf841f4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c2-standards@80eac02e",
            "digest": "sha256:28ca210d82b39720e00adbdfe02aef3c9a602eeb4ebc22c5d0c5d155f4b43737",
            "excerpt": "Standards: 3 findings. Worst: poolFixture.pool names two different directories; six reset_apply literals restate ownedAssignment.call."
          },
          "axis": "Standards",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "80eac02e1f04599f86cbe40b6374592b75fee5c0",
          "finding_ids": [
            "R10",
            "R11",
            "R12"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c2-spec-r1",
          "performer": "claude:bench-reviewer/vr-c2-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "36e7b63fe9a09b930fd31563be5ea77eadf841f4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c2-spec@80eac02e",
            "digest": "sha256:83c24dd6a0e32b19682e7cc04ce8d6a08df053fd99404e497eb28db70fe9c0a0",
            "excerpt": "Spec: 0 findings. VR23 to VR27 hold at 80eac02e; the tuple scan omits all eight builders and every ticket stays inside its Writes line."
          },
          "axis": "Spec",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "80eac02e1f04599f86cbe40b6374592b75fee5c0",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "vr-c2-coverage-r1",
          "performer": "claude:bench-reviewer/vr-c2-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "36e7b63fe9a09b930fd31563be5ea77eadf841f4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c2-coverage@80eac02e",
            "digest": "sha256:ebc6b8a790351c418e660c2c9bc5fa50e6ce7f4f5fb48497c3d48a3b7ab9f06a",
            "excerpt": "Coverage: 0 findings. Four new probes bit; 688 tests pass, the skip set is the two socket subtests, and no assertion count fell."
          },
          "axis": "Coverage",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "80eac02e1f04599f86cbe40b6374592b75fee5c0",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "vr-c2-standards-r2",
          "performer": "claude:bench-reviewer/vr-c2-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c2-standards-r2@f03e7fb9",
            "digest": "sha256:18c557db4cda967014bddddd0c5a3dc920e1cd3eac9256e41458fb5346557844",
            "excerpt": "Standards: 1 finding. R10 and R11 hold; the R12 record prose still disagreed with its advice list."
          },
          "axis": "Standards",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "finding_ids": [
            "R13"
          ],
          "supersedes": [
            "vr-c2-standards-r1"
          ]
        },
        {
          "id": "vr-c2-spec-r2",
          "performer": "claude:bench-reviewer/vr-c2-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c2-spec-r2@f03e7fb9",
            "digest": "sha256:27cce7e4f1cce26ecfc264fb2054f5c6f6a9fbff8b3472ea1633799b88312eb2",
            "excerpt": "Spec: 0 findings. The plan commit, the repair commit, and the three rerun entries match the tree and the plan."
          },
          "axis": "Spec",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "finding_ids": [],
          "supersedes": [
            "vr-c2-spec-r1"
          ]
        },
        {
          "id": "vr-c2-coverage-r2",
          "performer": "claude:bench-reviewer/vr-c2-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c2-coverage-r2@f03e7fb9",
            "digest": "sha256:403db01666cb3b0747caac3e8555aed881c7d8465cd552f61a530c4f98a1b98d",
            "excerpt": "Coverage: R10 and R11 hold, no assertion dropped, and the package passes. One weak lock check predates the chunk and stays advice."
          },
          "axis": "Coverage",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "finding_ids": [],
          "supersedes": [
            "vr-c2-coverage-r1"
          ]
        },
        {
          "id": "vr-c2-standards-r2-reaffirm",
          "performer": "claude:bench-reviewer/vr-c2-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b7f16bc7bc9bb2827e92686bc09084c902334231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c2-standards-r2/reaffirm@f03e7fb9",
            "digest": "sha256:beb074339f5316bc580174d1fdf2a3d42054a929d2c7032c4f261de423e1dd9b",
            "excerpt": "Standards reaffirm: pass. Line 166 now gives 4 Standards findings with 1 binding, which matches the excerpt."
          },
          "axis": "Standards",
          "base": "ef2cd35f1052f7003c1f7647539d656be7585821",
          "tip": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "finding_ids": [],
          "supersedes": [
            "vr-c2-standards-r2"
          ]
        }
      ]
    },
    {
      "id": "VR-C3",
      "base": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
      "tip": "ab9d305b314690959e5f68d3c751ff90169b8f29",
      "plan_digest": "sha256:76c2080cabf7a682039d7570abbfba170c11c247f78c80a9f3ca6149334653bc",
      "source_digest": "9fa1b3df144d6b1bfc4425a4a3e65f769589d816",
      "acceptance_rows": [
        "VR28",
        "VR29",
        "VR30",
        "VR31",
        "VR32",
        "VR33",
        "VR34",
        "VR35"
      ],
      "verification": [
        {
          "id": "vr-c3-5-worktree-r1",
          "performer": "claude:bench-writer/vr-t5-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a32bf83de855b97f131505b44f6ef435412620c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t5-author-20261001/5-worktree@b5e755a9",
            "digest": "sha256:2943eec08cd65a785b9db8e7fb953327f75411af5701fdab7b00b48ebee3ac9d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,49383\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "5-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c3-6-worktree-r1",
          "performer": "claude:bench-writer/vr-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f7c6830f5ff7842e168993b1026c1ea0f40ebd08",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t6-author-20261001/6-worktree@a531f85b",
            "digest": "sha256:800597aee235d83931a6c68ad8a4eecf87f515f2ca6b3a1fb23075b8453c9767",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,52046\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "6-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c3-7-worktree-r1",
          "performer": "claude:bench-writer/vr-t7-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9fa1b3df144d6b1bfc4425a4a3e65f769589d816",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t7-author-20261001/7-worktree@c6252480",
            "digest": "sha256:8f3aba6cc1fa109db02771be8823eaf68b3de2b281332d3d34488e2d7eb6e673",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,62556\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "7-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "vr-c3-standards-r1",
          "performer": "claude:bench-reviewer/vr-c3-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9fa1b3df144d6b1bfc4425a4a3e65f769589d816",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c3-standards@ab9d305b",
            "digest": "sha256:d5223808608a4a64dfb19f84983f00c74f0ee124f0924c045c168172f67c9094",
            "excerpt": "Standards: 3 findings. Worst: the record called a weaker mustNoFingerprint check narrower; table-name literals have two test-side spellings."
          },
          "axis": "Standards",
          "base": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "tip": "ab9d305b314690959e5f68d3c751ff90169b8f29",
          "finding_ids": [
            "R17",
            "R18"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c3-spec-r1",
          "performer": "claude:bench-reviewer/vr-c3-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9fa1b3df144d6b1bfc4425a4a3e65f769589d816",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c3-spec@ab9d305b",
            "digest": "sha256:32bb37311c25e6d35a6debe95b6143e2af5d270a3c69c93d63dbd8fffd761978",
            "excerpt": "Spec: 2 findings. Worst: the record logged one looser no-fingerprint check where nine exist; the spec names mustFingerprint for an all-none read."
          },
          "axis": "Spec",
          "base": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "tip": "ab9d305b314690959e5f68d3c751ff90169b8f29",
          "finding_ids": [
            "R15",
            "R16"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c3-coverage-r1",
          "performer": "claude:bench-reviewer/vr-c3-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9fa1b3df144d6b1bfc4425a4a3e65f769589d816",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c3-coverage@ab9d305b",
            "digest": "sha256:f806b309c9e8ff78aece445b89dcf7fe655b1cfc392705045901084735a97b5b",
            "excerpt": "Coverage: 1 finding. Worst: TestCleanSetHostileOperand no longer pins the none placeholder; a probe that empties it is silent for clean tests."
          },
          "axis": "Coverage",
          "base": "f03e7fb9b46e5cdbe0ae55d94f171b209fe27447",
          "tip": "ab9d305b314690959e5f68d3c751ff90169b8f29",
          "finding_ids": [
            "R14"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  },
  "amendments": [
    {
      "from": "sha256:095153eee4b7547b074ceddeb588e1a6a77a4e4663f33dd70ed0b00e55144d39",
      "to": "sha256:fc14637963d16b0fca7a4058f4eac565d31a02cb563e2fcdc1331d0547ba5c47",
      "chunk_ids": {
        "VR-C1": [
          "VR-C1"
        ],
        "VR-C2": [
          "VR-C2"
        ],
        "VR-C3": [
          "VR-C3"
        ],
        "VR-C4": [
          "VR-C4"
        ],
        "VR-C5": [
          "VR-C5"
        ]
      }
    }
  ]
}
```
