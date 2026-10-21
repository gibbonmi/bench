## Outcome

The FT419 candidate-binary-install build landed as `3b970d1c2d1686f05436dea57e5c845037723cb0` with a green whole-project gate. A spec can now declare persisted fields, and `bench worktree build --install` installs a declared build's candidate executable. The landing refresh holds while a declared build is open. A guarded write refuses on an uninstalled executable, and `bench worktree build --restore` clears declared ledger fields and rebuilds the primary executable. ADR 0016 and the guidance name the install exception. All 62 acceptance rows reconciled, and the change has 60 files, 7952 insertions, and 193 deletions.

The build ran six tickets in five chunks on Claude Code as a `--full --delegate --auto-approve` run with an author limit of 3. Fresh Opus/high authors and ten Opus/high repair sessions wrote the code and prose. Opus/medium ran every review axis, and one Opus/high and Sonnet/high consultant pair decided two behavioral questions by consensus. A Haiku staleness pass opened the build. The run took about 8.5 hours of wall time, from 11:35 to 20:09.

The chunks used 2, 1, 3, 2, and 1 repair cycles. CB-C3 took a third cycle under `--auto-approve`. The consultants added rows CB66 to CB69 and story 46, and each addition is flagged for reviewer veto.

## Gate-stage timings

- landing: commit 3b970d1c2d1686f05436dea57e5c845037723cb0, trace 7c43b5b1c6fcf6b87090b6bcdfa7eff0
- gofmt: 194 ms
- vet: 1441 ms
- test: 174666 ms
- race: 2267 ms
- system: 95491 ms
- shellcheck: 560 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 7 repair 1 | F1 one install instruction | verified | 8 | held | opus / high / repair |
| ticket 7 repair 1 | F2 split paragraph at 95 of 95 lines | verified | 8 | held | opus / high / repair |
| ticket 7 repair 1 | F3, S1, C6 version-cache refusal sentence | verified | 8 | held | opus / high / repair |
| ticket 7 repair 1 | F4 ADR Decision qualifier | verified | 7 | held | opus / high / repair |
| ticket 7 repair 1 | F5 ADR rewrap | claimed | 8 | held | opus / high / repair |
| ticket 7 repair 1 | S2 restore rebuilds from the primary sources | verified | 9 | held | opus / high / repair |
| ticket 7 repair 1 | C1 one declaration grammar | verified | 7 | held | opus / high / repair |
| ticket 7 repair 1 | C2 stop-short restore pointer | verified | 7 | held | opus / high / repair |
| ticket 7 repair 1 | C3 profile install exception | verified | 8 | held | opus / high / repair |
| ticket 7 repair 1 | C4 overlap route | verified | 8 | held | opus / high / repair |
| ticket 7 repair 1 | C5 recovery install removed | verified | 8 | held | opus / high / repair |

Brier mean 0.051 over 11 pairs, 0 abstained.

The context compaction dropped the claim tables of the earlier authors and repairs, so this table holds only the last repair. The confirming round closed every row above.

Each ticket stayed inside its chunk slice, but three tickets needed `Writes:` widenings. Ticket 6 grew `build.go` to 470 lines and took a new `build_restore.go`. Its cycle 2 grew `transaction_test.go` to 496 lines and took a new `remove_fields_test.go`. Review found five real defects that the authors' tests passed through. These were the lost missing-tree refusal, the linked-repository declaration read, the untested manifest location, two open declared builds, and the absent-ledger restore.

## Coordinator catches

- The first charge refused `--check-current`, because the orchestrator prepared it before the assignment commit.
- The first chunk base was wrong. The coordinator corrected it to the `main` tip `97a97926`.
- In CB-C2, the last author's record commit blocked ticket 2's entries. The coordinator reset the tree to `cd364bb6` with the record as a working change, recorded, and committed once.
- `bench record confirmation` refused the CB-C1 test-only repair, because the repair-assignment plan commit is in the delta. The chunk took a full confirming round.
- Two line-budget stops left dirty author work. The coordinator saved the diff, reset with the preserving envelope, widened the ticket, and resumed the author.
- The coordinator's own probe found the untested manifest location in CB-C3.
- The first landing used a new request token. The reauthorize route then lost the commitment binding, and the coordinator restored the original token.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-declare-persisted-fields.md | 2 | spec-row, delegate-error |
| 2-install-the-candidate.md | 1 | delegate-error |
| 3-hold-the-refresh.md | 1 | spec-row |
| 5-refuse-an-uninstalled-writer.md | 3 | spec-row, delegate-error, delegate-error |
| 6-restore-after-a-stopped-build.md | 2 | one-source, one-source |
| 7-document-the-install-route.md | 1 | one-source |

## Agent-experience improvements

### Bench CLI

- Census account: the landing counted 80 raw calls; rg 60 and ls 2 are sanctioned reads, and cat 4, cd 3, awk 1, and sed 1 are unknown delegate calls, so propose no Bench form for them.
  Feeds: none
- Census account, owned gap: python3 6 read the plan fence and the review record, so let `bench assessment record` derive its attempts from those two records.
  Feeds: new
- Census account, missed form: wc 3 counted guidance lines, so route a budget count to `bench test --check guidance-prose-budgets` in the author charge.
  Feeds: none
- Make the request-mismatch refusal of `bench worktree land` name the assignment's bound request token, so the operator does not replace it.
  Feeds: new
- Make `bench probe --check prose` copy the first prose diagnostic into its `first_failure` cell.
  Feeds: new
- Add a read-only verb that lists every anchor needle and canary string that pins a given guidance file.
  Feeds: new
- Make the test-only predicate of `bench record confirmation` ignore a repair-assignment record in the plan fence.
  Feeds: FT318
- Let each author of a multi-ticket chunk record its entries after the last author's record commit.
  Feeds: FT369

### Skills

- State in the implement phase that the orchestrator prepares each charge after its assignment commit.
  Feeds: none
- State in the author charge that a Go file within 30 lines of its cap takes a new sibling file before the edit.
  Feeds: new

### Process

- Send each behavioral review question to the two named consultants, because their consensus closed both questions in one pass.
  Feeds: none
- Keep three review axes on a confirming round, because the CB-C5 confirming axes each found a different advice item.
  Feeds: none