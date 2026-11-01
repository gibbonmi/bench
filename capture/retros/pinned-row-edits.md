## Outcome

FT437 landed at `46e80c1af` on a green whole-project gate. Plan and approve now read a changed pinned row at the working file of the checkout. The commit and the landing read that row at the candidate tree. One root-relative link walk in `internal/bounds` serves the retro write and the working read, and it wraps `bounds.ErrLinkRefused`. All 29 coverage rows pass, and the landing closed the FT437 row.

## Cost

- landing: 46e80c1af0cc6f428c811a3a22095da0106cd501, range 46e80c1af0cc6f428c811a3a22095da0106cd501^1..46e80c1af0cc6f428c811a3a22095da0106cd501^2
- wall: 57 min
- commits: 9 total, 6 bookkeeping, 3 change
- graph: 1 tickets, chain 1, 1 chunks, 1 serial checkpoints
- dispatch: 1 authors, 1 repairs
- trace: commit 9 spans in 79133 ms, lane 0 spans in 0 ms, landing 1 spans in 328117 ms
- wall per ticket: 57.9 min against 60 min
- bookkeeping per ticket: 6.0 against 7.1
- delegate tokens: about 810k over 14 sessions; the author took 177k, the repair 83k, and the six completed axes 308k. Three stopped Fable axes returned no count. Provider charges remain unknown.

## Gate-stage timings

- landing: commit 46e80c1af0cc6f428c811a3a22095da0106cd501, trace e7d2ccc9758b5a513a435e98162be38a
- gofmt: 184 ms
- vet: 1678 ms
- test: 229787 ms
- race: 9997 ms
- system: 46741 ms
- shellcheck: 532 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| PE-C1 review | S1: a pinned-row fixture is pasted twice | confirmed | 7 | 1 | opus / medium / reviewer |
| PE-C1 review | S2: the changed-from expectation is built three times | confirmed | 6 | 1 | opus / medium / reviewer |
| PE-C1 review | P1: the PE21, PE22, and PE25 catch clauses are false | confirmed | 9 | 1 | opus / medium / reviewer |
| PE-C1 review | P2: the PE25 landing test lacks the face guard | confirmed | 7 | 1 | opus / medium / reviewer |
| PE-C1 review | P3: the protectedCandidate sentence is false | confirmed | 6 | 1 | opus / medium / reviewer |
| PE-C1 review | C1: each half of the ChangedSource match survives a mutation | confirmed by probe | 4 | 1 | opus / medium / reviewer |
| PE-C1 review | S3, S4, C2, C3 | advice | 5, 3, 7, 6 | abstained | opus / medium / reviewer |
| PE-C1 author | the planned commit probe stays silent alone | confirmed by the Spec axis | unknown | abstained | opus / medium / implementer |

Brier mean: 0.145 over 6 pairs; 5 abstained.

The author implemented the ticket inside its fence on the first pass. It kept `sourceAt` unchanged, so two `Writes:` entries stayed untouched. The repair closed four test findings in one cycle, and the orchestrator closed the two prose findings in the plan commit.

## Coordinator catches

- The staleness pass found the new `bounds.ErrLinkRefused` contract. The orchestrator decided that the moved walk wraps it, so that one sentinel marks each link refusal.
- The review preflight refused `--source-tip` at the frozen code tip, because it derives the tip from the record commit.
- `bench record confirmation` refused the test-only route for `commitmenttest/repo.go`, so the chunk took a confirming round.
- The refreeze of a test-only repair owed all five plan probes again.
- The Spec axis cited line numbers of the diff, not of the files, so the orchestrator checked each citation at the tip.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-read-changed-pinned-rows.md | 1 | one-source |

## Agent-experience improvements

### Bench CLI

- Add a `bench record pickup` form that writes the per-axis pickup sections of a review round, because the orchestrator appended them with two raw heredoc calls (census class: owned gap).
  Feeds: new
- Name `bench spec show <slug> --section coverage` and `bench worktree read --section` in the review and author charges, so that a session does not search spec rows with `rg` (census class: missed form, 2 calls).
  Feeds: new
- Count `rg` code searches apart in the census, because 33 sanctioned `rg` reads hide the other classes in the head totals (census class: sanctioned read).
  Feeds: new
- Two raw calls, one `git log` and one `sed`, have no visible purpose in the census rows, so this close proposes no change for them (census class: unknown).
  Feeds: none
- Measure the test-only predicate of `bench record ticket` on the repair delta, so that a test-only repair does not owe every plan probe again.
  Feeds: new
- Let `bench record confirmation` accept a test-support package that only test files import, so that a fixture move closes without a confirming round.
  Feeds: new
- Make `bench preflight review` accept the frozen code tip of the last chunk entry as `--source-tip`, or name the derived tip in its refusal.
  Feeds: new

### Skills

- Send one staleness delegate for a spec that has at most one ticket, because two delegates read the same seven drift paths.
  Feeds: new
- State in `docs/known-issues.md` row 4 that its workaround applies at a refreeze, because the first freeze refuses `bench record amendment`.
  Feeds: none

### Process

- Declare the expected repair rounds and their confidence in the line declaration, as the Claude scorecard requires, because this build omitted it.
  Feeds: none