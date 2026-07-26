# Review — FT86 fail-closed control records

Reviewed range: `a223629..f9cc907` (seven build commits, already landed on `main`;
`bench diff` reports empty because the work is on the default branch).
Spec: `specs/ft86-fail-closed-control-records.md`.

## Standards

11 findings. Worst: the failed-read predicate is hand-derived at nine sites.

1. **Failed-read predicate duplicated nine times.** `state == StateUnreadable ||
   state == StateWrongType` is re-derived at `internal/roadmap/roadmap.go:218,232`,
   `internal/maps/maps.go:268,302,307`, and `internal/status/status.go:446,447,481,532,615`,
   while `internal/bounds` — the package that owns the vocabulary — exports no
   predicate for it. `internal/roadmap/roadmap.go:120` and
   `internal/roadmap/context_parse.go:659` (`degradedState`) are two further,
   differently-scoped variants. Adding a state to `FileState` requires editing
   every one of them to stay correct. AGENTS.md: "two derivations of the same
   fact … must collapse to one source."

2. **Three per-package copies of one error grammar.** `learnings.journalError`
   (`internal/learnings/learnings.go:230`), `maps.mapsError`
   (`internal/maps/maps.go:409`), and `roadmap.roadmapReadError`
   (`internal/roadmap/roadmap.go:800`) each render the same
   `toon.Errorf("<path> is <state>", reason)` line. `roadmapReadError`'s own
   comment advertises it as sharing "one grammar with every other migrated
   surface" — the enforcement and its advertisement in separate places.
   `internal/status/status.go:449,455,475,533` adds four hand-built copies of the
   sibling `unknown (<path> is <state>)` cell grammar. Neither grammar lives in
   `toon`, which owns every other AXI line shape.

3. **Story-reference grammar encoded twice.** `storyRefRe`
   (`internal/coverage/coverage.go:43`) and `storyPartRe` (`:45`) both encode the
   number / optional en-dash-or-hyphen range grammar. `:213` dereferences
   `storyPartRe.FindStringSubmatch(...)[1]` with no nil guard. Not a live panic —
   `storyRefRe` currently guarantees every trimmed comma-part matches — but the
   absence of a guard makes the two regexes a lockstep obligation, where drift
   turns a validation miss into a crash.

4. **`IDEAS.md` has two readers again.** `roadmap.lineCount`
   (`internal/roadmap/roadmap.go:246-253`) re-implements `ideaLines`' body
   (`:861-872`) instead of composing it, and `ParkedIdeas`' doc comment at `:166`
   still claims the count and the listing "both go through `ideaLines`, one
   source" — now false.

5. **Fixture harness pasted three times.** The chmod-0o000 + `capability.Privilege`
   block appears at `internal/contract/axi/axi_maps_fail_closed_test.go:31-38,63-70,95-102`,
   while the same diff extracts `chmodUnreadable`
   (`internal/contract/axi/axi_status_fail_closed_test.go:15`) for exactly that
   block and documents it as matching the maps file's pattern. The standard names
   this case verbatim.

6. **Control-record paths hardcoded beside their owners.**
   `internal/status/status.go:449,455` hardcodes `"IDEAS.md"` and
   `".bench/learnings.md"` while `roadmap.ideasFile`
   (`internal/roadmap/roadmap.go:38`) and `learnings.journalPath`
   (`internal/learnings/learnings.go:176`) own them;
   `internal/roadmap/roadmap.go:218` adds a third derivation via
   `filepath.Join(root, ".bench", "learnings.md")`.

7. **Comment narrates the change.** `internal/roadmap/context_types.go:721-723`
   describes "the same tally `readDirSource` used to compute inline" — citing a
   function this diff deletes. `craft-comments`: no narration of the change.

8. **Comments argue their own correctness.** `internal/coverage/coverage.go:308-311`
   ("so the two branches below are exhaustive") and
   `internal/roadmap/context_types.go:680` ("It reuses `bounds.OutlineFileLimit`
   rather than declaring a second constant") both address the diff's reviewer
   rather than the next reader.

9. **Doc comment wider than the enforcement.**
   `internal/coverage/coverage.go:224` promises membership validated "against the
   exact declared story set", but the range branch (`:216-222`) member-checks only
   `start` and `end`. See Spec finding 5 — same defect, both axes.

10. **Duplicated failing call.** `internal/worktree/lifecycle.go:1064` — when
    `RemoteDefaultRef` returns `""`, the expression reduces to
    `!worktreeAdd(root, cand, "") && !worktreeAdd(root, cand, "")`, running the
    identical failing add twice.

11. **Two `"main"` constants.** `git.mainCandidate`
    (`internal/git/default_branch.go:9`) and `adopt.fallbackProtectedBranch`
    (`internal/adopt/link_hook.go:14`). Defensible as different facts (an
    object-DB probe vs. a guard fail-safe); flagged so the split is confirmed
    rather than assumed. See Spec finding 2.

## Spec

7 findings, 32 of 34 coverage rows verified sound. Worst: two enumerated call
sites took a posture the spec's table forbids.

1. **Story 18 — two call sites reversed their assigned posture.** The
   implementation-decisions table assigns `internal/adopt/link_hook.go:95` and
   `internal/adopt/link_stage.go:97` "skip the candidate". The build instead
   introduced `fallbackProtectedBranch = "main"`
   (`internal/adopt/link_hook.go:96`) and routes both through `protectedBranch`,
   so an unresolvable default still bakes the literal `main` — against story 18's
   own text, "no code path can fabricate `main`". The in-code reasoning (a guard
   protecting nothing is worse than a guard on a guessed branch) may well be the
   right call; it is an undeclared deviation from an enumerated posture and is
   the reviewer's to accept or reject. The other nine call sites match.

2. **Story 16 — both coverage rows cite a test that cannot run.** The rows name
   `go test ./internal/contract/axi -run TestAXIOutlineSymlinkSkipped`. No such
   function exists; the real test is the lowercase subtest
   `testAXIOutlineSymlinkSkipped` (`internal/contract/axi/axi_outline_test.go:140`),
   dispatched from `TestAXIOutlineContracts` via `contract.RunParallel`. Verified:
   the cited command exits 0 with `[no tests to run]`. Outline was genuinely left
   unmigrated, so the non-change held — but the guard the spec calls "the thing
   that forbids migrating outline" cannot go red as written.

3. **Story 11 — `bench maps --count` still fabricates a zero.**
   `internal/maps/maps.go:302-303` returns `"0\n", 0` when `decisions/` is
   unreadable or wrong-type, on the very surface the status adapter consumes,
   while `bench maps` on the same fixture exits 1 with the error row. Only the
   Go-level `UnresolvedCount` carries state. `TestAXIMapsCountMatchesListing`
   exercises a failing file *inside* the directory, never the whole-directory
   failure, so no row catches it.

4. **Story 17 — one of three migrated signals is unexercised.** `appendMaps`
   builds the `unknown (decisions is …)` row (`internal/status/status.go:478`),
   but `TestAXIStatusUnknownRow` and `TestAXIStatusUnknownNotSuppressed` assert
   only the drain and roadmap rows.

5. **Story 15 — ranges validate endpoints, not membership.**
   `internal/coverage/coverage.go:216-222` calls `checkStoryMember` on `start` and
   `end` only, so `2-4` passes against a declared set of {1, 2, 4} despite 3 being
   undeclared. "Membership validated exactly" is not what the range arm does.

6. **Built but not specified.** `isTemplatePlaceholder`
   (`internal/learnings/learnings.go:117`) exempts the scaffold's `## <date>`
   example from malformed classification — no story asks for it.
   `maps.UnresolvedCount` now counts unsupported-schema and unreadable files as
   unresolved; the spec never states that equivalence, and its consequence is
   that a routine `decisions/README.md` turns `bench maps` into exit 1
   (`internal/contract/axi/axi_test.go:284-299` re-baselines exit 0 → 1 for
   exactly that case).

7. **Story 12 — present-but-empty `ROADMAP.md` still reads as absent.** `bench
   roadmap` prints `no ROADMAP.md` for an empty file, against the Solution's
   "anything but absent" rule. Already parked in `IDEAS.md` by the builder and
   pinned by a green runtime contract; listed here for disposition, not as news.

**Verified sound:** the other 32 rows' tests exist at the named seam and assert
the named behavior. Both canary fixtures carry correct EXPECT files
(`tests/canary/coverage-map-validation/no-map-not-historical/EXPECT`,
`tests/canary/package-core-guard/default-branch-refabricated/EXPECT`) and are
wired into `internal/conformance/fixture_bite_test.go:134,183`; the second also
ships a negative case proving the sweep stays silent on the surviving struct
field.

## Coverage

8 findings. Worst: a FIFO at `specs/*.md` hangs `bench status` with no timeout.

1. **FIFO in `specs/` hangs `bench status` and `bench dashboard` forever.**
   Verified in a throwaway repo: `mkfifo specs/hang.md` → `bench status` hits
   `timeout 5` (exit 124). Reached through `internal/status/status.go:567`
   (`retirementCount`) and `:627` (`roadmapReconcileCounts`), plus
   `internal/spec/spec.go:128` (`readCandidate`) — all raw `os.ReadFile`.
   `TestClassifyFIFOWithoutOpen` (`internal/bounds/classify_test.go:161`) proves
   `bounds` type-checks before opening; these callers never reach it. The
   profile's hostile-input checklist names this class verbatim. Note the scope
   caveat: the spec explicitly holds `internal/spec`'s housekeeping counters out
   of the migration, so this is partly pre-existing — but `bench status` is what
   the SessionStart hook runs, and it now hangs. (I could not reproduce the
   delegate's companion claim that `bench coverage --check` hangs; it exited 2.)

2. **FIFO at `IDEAS.md` hangs `bench dashboard`.** Verified: exit 124.
   `roadmap.ideaLines` (`internal/roadmap/roadmap.go:247`) and `RoadmapText`
   (`:158`) stayed on `os.ReadFile` while `lineCount` beside them migrated, so
   two readers of one file now disagree about whether it is safe to open. This is
   migration residue inside the spec's own scope.

3. **Non-UTF8 control record makes `bench status` render a clean board.**
   `printf '\xff\xfe' > IDEAS.md` → `bench status --all` prints no drain row and
   no roadmap row at all, exit 0, while `bench learnings` and `bench roadmap` on
   the same bytes exit 1 with `is malformed — invalid UTF-8`.
   `internal/roadmap/roadmap.go:215-225,229-240` and
   `internal/status/status.go:612-619` switch on `StateUnreadable`/`StateWrongType`
   only, so `StateMalformed` falls through the default branch and parses the
   garbage. `axi_status_fail_closed_test.go` drives only chmod-0000; a malformed
   row belongs beside the unreadable one in `statusUnknownFixture`.

4. **`bench maps --count` on an unenumerable `decisions/`.** `chmod 000
   decisions` → `--count` prints `0` at exit 0 while `bench maps` exits 1. Same
   defect as Spec finding 3; recorded on both axes because the Coverage gap is
   that `TestAXIMapsCountMatchesListing` only makes a file *inside* the directory
   unreadable.

5. **A repo with no commits kills the whole `--context` snapshot.** `git init`
   with zero commits → `bench roadmap --context` returns
   `error: roadmap context failed — exit status 128` from the unguarded
   `rev-parse --abbrev-ref HEAD` at `internal/git/git.go:284`, even though FT86
   taught `Facts` to degrade an unresolvable default. `bench diff` on the same
   repo degrades correctly. Related: `ResolvedDefault`'s **zero**-local-branch
   path (`internal/git/default_branch.go:23-26`) has no test —
   `internal/git/git_test.go` covers one branch and two.

6. **Read-limit divergence at one record.** `learnings.Command` reads with
   `ModelReadLimit` (5 MiB, `internal/learnings/learnings.go:177`);
   `roadmap.learningCount` reads the same file with `controlRecordLimit` (2 MiB,
   `internal/roadmap/context_types.go:11`). A 3 MiB journal renders rows in
   `bench learnings` and reads `unknown (… is unreadable)` in `bench status`. No
   oversized fixture exists at any consuming parser — only at the `bounds` seam
   (`internal/bounds/classify_test.go:59-66`).

7. **`StateWrongType` for a file is asserted only behind capability guards.**
   `internal/bounds/classify_test.go:71-78` reaches it via `requireSocket`
   (filed under `capability.Fifo`, so a socket-less host is tallied as a FIFO
   gap) and `:179-191` via `/dev/null` plus a socket. On a host without unix
   sockets nothing asserts `wrong-type` for a file, and no AXI test drives
   wrong-type at any surface — though a directory at the record path needs no
   capability at all and works. `BENCH_REQUIRE_CAPABILITIES=1` is satisfiable on
   this host (zero skips) but unsatisfiable as root, since `requireUnreadable`
   (`internal/bounds/classify_test.go:276`) and `chmodUnreadable`
   (`internal/contract/axi/axi_status_fail_closed_test.go:16`) must skip there.

8. **`coverage --check` on a spec with no `## User stories`.** A spec with a
   valid map and no story section emits
   `references story 1, which the spec does not declare (has: )` — a degenerate
   empty-set clause at `internal/coverage/coverage.go:241`
   (`declaredStoriesList`). Every AXI fixture includes `axiCoverageStories`, so
   nothing covers a missing or differently-headed story section.
