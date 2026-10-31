## Outcome

The `exact-planning-readers` build (FT125) landed at `62e0b05f` on a green whole-project gate. Agents now read exact bytes through native readers: `bench read` for a file or a line range, `bench spec show` for a section or a story, `bench outline --symbol` for a Go declaration, and `bench worktree read` for an owned worktree. Every page carries a source identity, and a continuation refuses when the source changed.

`bench worktree list --spec` lists the worktrees that a spec binding relates, through the commitment scopes. An ordinary raw read of one file inside the repository stays allowed and gets native-form advice on stderr, and an owned read gets one census record. Census output heads now name `spec show` and `outline --symbol`. The spec closes 49 coverage rows; ER44 is a review-owned row by design.

## Cost

- landing: 62e0b05ffb51abfa62a315ade23d0cae75f1bd60, range 62e0b05ffb51abfa62a315ade23d0cae75f1bd60^1..62e0b05ffb51abfa62a315ade23d0cae75f1bd60^2
- wall: 1103 min
- commits: 116 total, 90 bookkeeping, 26 change
- graph: 10 tickets, chain 7, 5 chunks, 5 serial checkpoints
- dispatch: 10 authors, 16 repairs
- trace: commit 122 spans in 1154654 ms, lane 0 spans in 0 ms, landing 1 spans in 325116 ms
- wall per ticket: 110.3 min against 60 min, over budget
- bookkeeping per ticket: 9.0 against 7.1, over budget

The wall time includes an idle night between the two orchestrator sessions, from 00:26 to 11:20. Provider usage and charges are unknown.

## Gate-stage timings

- landing: commit 62e0b05ffb51abfa62a315ade23d0cae75f1bd60, trace c2514e3d9e8862933b2a3209c39de855
- gofmt: 188 ms
- vet: 1652 ms
- test: 223408 ms
- race: 9212 ms
- system: 46504 ms
- shellcheck: 541 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| none | No labeled probability claim was made in this build. | unknown | unknown | unknown | claude-opus-5-5 / high / orchestrator |

Brier mean unknown over 0 pairs, 0 abstained.

Every author met its ticket boundary. ER-C3 used five repair cycles, ER-C4 four, and ER-C5 three; each run past two cycles went under `--auto-approve`. Most late findings came from rules that a repair itself added, such as the related read label rule and the raw-read repository bound. The ER-C5 authors had the spec edge cases in their first charge, and the first ER-C5 review still found five defects.

## Coordinator catches

- The resumed handoff named the ER-C3 checkpoint, but the tree was at the ER-C4 repair cycle 3 assignment. The orchestrator resumed from the tree.
- The first resumed charge used the `bench worktree exec` form that `docs/known-issues.md` warns about. The reviewer asked for a known-issues check, and later charges used the `--in` form.
- Repair cycle 3 of ER-C4 found a second source of the `primary` keyword outside its fence. The orchestrator accepted it as R12 and joined it with the next confirming round's finding in one repair.
- The final `./...` verification found a red anchor needle from ticket 5c. The author's `--run` filter and every review axis had missed `TestAdmissionClasses`.
- The ordinary assessment record was not written before the landing. The orchestrator wrote it after the landing, with the earlier session folded into one attempt.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1.md | 2 | one-source, delegate-error |
| 2a.md | 0 | none |
| 2b.md | 2 | edge-found, one-source |
| 3a.md | 0 | none |
| 3b.md | 5 | edge-found, other, other, edge-found, delegate-error |
| 4a.md | 0 | none |
| 4b.md | 4 | edge-found, spec-row, edge-found, one-source |
| 5a.md | 0 | none |
| 5b.md | 0 | none |
| 5c.md | 3 | edge-found, one-source, check-gap |

## Agent-experience improvements

### Bench CLI

- Count the 27 `rg` calls of this landing apart from the other 21 raw heads, because sanctioned reads hide the other classes.
  Feeds: new
- Let `bench record review` take its excerpt on stdin, because the orchestrator wrote 22 scratch excerpt files with raw `cat` heredocs for a record that Bench owns.
  Feeds: new
- Add a review-prose append form to `bench record`, because the orchestrator appended each round's prose to the review record with raw `cat >>` calls.
  Feeds: new
- Let `bench assessment record` build a run from the completion plan's assignments, because the orchestrator wrote the record by hand with a `python3` script.
  Feeds: new
- Route a read of the review record tail to `bench read --lines`, because one raw `tail` call read state that the new reader serves.
  Feeds: none
- Let `bench record ticket` accept `--in`, because every author tried the `--in` form first and then used `bench worktree exec`.
  Feeds: new

### Skills

- Add to `craft-delegate` that a charge runs each touched invariant package, such as the anchors registry, whole and not through a `--run` filter.
  Feeds: new
- Add to the bounded repair policy that the orchestrator joins a repair's out-of-fence defect with the next confirming round's findings in one repair cycle.
  Feeds: none

### Process

- Read the bounded repair policy and `docs/known-issues.md` before the first charge of a resumed build, because this session proposed two cuts that the policy already owned.
  Feeds: none
- Commit a confirming round's record and the next repair assignment in one commit, and let the repair write its witness reds in its own record commit; this build used both cuts from ER-C5 on.
  Feeds: none
- Write the assessment record when the build starts, because a census file that the landing releases cannot feed a record that is written after the landing.
  Feeds: new