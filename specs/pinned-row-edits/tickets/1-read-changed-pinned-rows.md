# Read a changed pinned row at its publishing location

Blocked by: none
Writes: internal/commitment/delivery.go, internal/commitment/repository/sources.go, internal/commitment/repository/candidate.go, internal/commitment/repository/repository.go, internal/commitment/repository/closure.go, internal/bounds/classify.go, internal/roadmap/retro.go, internal/commitment/pinned_row_test.go (new), internal/commitment/pinned_row_hostile_test.go (new), internal/commitment/head_binding_test.go, internal/commit/commitment_test.go, internal/commitment/repository/head_binding_test.go, internal/conformance/testdata/fixture-skip-classes.json, tests/canary/package-core-guard/bounds-classify-limit-restated, tests/canary/package-core-guard/bounds-read-limit-restated, CHANGELOG.md, .bench/commitment.json
Covers: PE01, PE02, PE03, PE04, PE05, PE06, PE07, PE08, PE09, PE10, PE11, PE12, PE13, PE14, PE15, PE16, PE17, PE18, PE19, PE20, PE21, PE22, PE23, PE24, PE25, PE26, PE27, PE28, PE29

## What to build

An approved plan and one commit record a pinned-row decision in its row file.
Read these spec sections before you start: Implementation decisions, Testing decisions, the acceptance coverage map, and the edge inventory.

Make these changes:

- Add one predicate to the commitment package beside `NewSource`. It is true for an outcome source of the proposal when a current outcome source binds the same id and path, and no current outcome source binds the exact same identity. Compose it on `namesSource`.
- Give `readTree` the new answer of the read-location table in the spec. At plan and approve, a changed outcome source reads the working file under the store root. At commit and landing, it reads the candidate tree that each caller already passes. Keep each other answer of the table. The new answer reaches the three `readTree` callers in `sources.go` and `candidate.go`.
- Move `refuseSymlinkComponents` from `internal/roadmap/retro.go` into `internal/bounds/classify.go`, beside `RefuseLinks`, as one exported walk. Give its refusal a neutral text that suits a read and a write and keeps the words `symbolic link`, which `TestRetroRefusesSymlinkedDestinationComponents` requires. Make that refusal wrap `bounds.ErrLinkRefused`, so that one sentinel marks each link refusal of the package. That test needs no edit. Make the retro write call the moved walk, and leave no copy in `internal/roadmap`.
- Read the working file with the moved walk from the store root, and then with `bounds.ClassifyNoFollow`. Accept only a parsed classification.
- Factor the normalization of `sourceAt` once, so that the working read and the tree read compute one identity for the same bytes. If the signature of `sourceAt` changes, update each of its callers in `closure.go`, `repository.go`, and `sources.go`. No test calls `sourceAt` directly.
- Keep each refusal text and each refusal face that the spec names. Add no error literal that holds a next command, because the route-bypass check scans these packages.
- Rewrite the `changed-row` subtest of `TestHeadBindingKeepsTheDefaultTipForEveryOtherSource` as the PE17 witness. Update the comments above `readTree` and above that test, so that they state the new rule.
- Add the planned tests of the coverage map. Put the hostile-file table in `pinned_row_hostile_test.go` and the other plan and approve tests in `pinned_row_test.go`.
- A FIFO test follows the precedent of `internal/commit/lane_test.go`, which calls `capability.Unavailable`. Add its row to `internal/conformance/testdata/fixture-skip-classes.json`.
- The two `bounds-*-limit-restated` canary units read `internal/bounds/classify.go`. Keep their BASE, MUTATE.json, EXPECT, CHECK, and assertions unchanged, and confirm that both still bite.
- Add one Unreleased entry to `CHANGELOG.md` that names the new read locations.

The ticket changes no field of the intent ledger or the review record.
`.bench/commitment.json` is on the `Writes:` line for the spec-stage binding commit only. The author does not edit it.

## Acceptance

- [ ] Plan and approve accept the new identity of a pinned row while the edit stays only in the working file, for one row, for two rows, and for a store root below a linked directory (PE01, PE02, PE06, PE27).
- [ ] Plan refuses a new identity that the working file does not hold, including a HEAD-only edit, and the refusal writes nothing (PE03, PE04, PE05).
- [ ] Approve refuses a plan whose working row changed after the plan (PE07).
- [ ] Plan refuses each hostile working path of the edge inventory (PE08, PE09, PE10, PE11, PE12, PE13, PE14).
- [ ] The working identity drops the `Occurrences:` ledger line, and a malformed ledger refuses (PE15, PE16).
- [ ] A kept identity, a new unpinned source, and a new deliverable keep their reads (PE17, PE18, PE19).
- [ ] `bench commit` records the approved policy, the board, and the row edit in one commit, and it refuses the policy-only and the rewritten-row commits (PE20, PE21, PE22).
- [ ] A pinned-row edit with no policy change still refuses at commit (PE23).
- [ ] The landing admits the composed approved edit and refuses a composed row that differs from the approval (PE24, PE25).
- [ ] The changelog names the new read locations (PE26).
- [ ] The retro write keeps its link refusal through the moved walk, and no copy of the walk survives (PE28, PE29).
