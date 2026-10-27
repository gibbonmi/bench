# FT451 background grade: research

Consumed by: the `decisions/background-grade.md` map, its decision tickets, and then `/bench-write-spec`.
Drift: a change to the grade key, the evidence store, the gate lock, a checkpoint subject, or the host limits.
Retire when: `/bench-write-spec` compiles the map into a spec.

## Recommendation, scope, and evidence status

Let a background grade compute the exact run of the next consumer from immutable inputs only.
Let it write the shared evidence store and hold no checkout lock.
Bound it machine-wide and below the foreground.

Scope: the gate reuse machinery at `main` commit `48f9241df` and the consumers that look up a green.
The scope also covers the run-binary selection, the process lifetime rules, and the host resource census.
This run changed no code and ran no gate.

Evidence status: one reader read every source in this session, by reviewer direction, so no paired read exists.
Each fact cites a path and a line at `48f9241df`, or a commit for a drained record.
No measurement ran.

## Facts

### Q1. How does a whole-gate verb reuse a green today?

- The grade key frames the graded tree with the grade identity (`internal/gate/engine.go:151-159`). ADR 0031 names the identity members (`docs/adr/0031-one-grade-key-names-each-whole-gate-green.md:7-19`).
- The evidence store keeps one file for each grade key in the git common directory (`internal/gate/engine.go:147-149`, `internal/gate/run_durability.go:105-117`). So every worktree of the repository shares it.
- The gate cache `bench-last-gate` and the gate lock `bench-gate.lock` live in the admin directory of each checkout (`internal/git/git.go:175`, `internal/gate/gate.go:380`, `internal/gate/verdict.go:151-157`).
- Before the lock, a run reuses a green only when the reuse writes nothing (`internal/gate/run_transaction.go:78-84`, `internal/gate/run_durability.go:68-81`). That holds when the checkout cache already describes the green, or when another run holds the lock.
- Each lookup refuses reuse while the checkout cache holds a pending record, even when the store holds a green under the key (`internal/gate/engine.go:219-226`, `internal/gate/verdict.go:177-190`). This is the GR44 edge.
- A run that cannot take the checkout lock refuses with `gate execution already in progress` and names the owner pid (`internal/gate/run_transaction.go:94-101`). No verb waits for the lock.
- A green expires 60 minutes after its source run, and a reuse keeps the source time (`internal/gate/verdict.go:24`, `internal/gate/engine.go:137-142`).
- A red run or a timed-out run removes the green under its key (`internal/gate/run_transaction.go:214-225`, `internal/gate/run_transaction.go:241-252`). The store holds whole greens only (`internal/gate/engine.go:165-167`).
- A lane pass writes only its lane record and no verdict (`internal/gate/lane.go:7-9`, `internal/gate/lane.go:190-197`).

### Q2. What subject and running executable does each consumer compute?

| consumer | graded tree | running executable | citation |
|---|---|---|---|
| `bench gate` in a worktree | the working tree | the executable that the wrapper resolves | `internal/gate/evaluation.go:39-56`, `bin/bench.sh:293-317` |
| chunk checkpoint | the HEAD tree without `reviews/<slug>.md` | the running executable, in its own process | `internal/gate/checkpoint.go:286-306`, `internal/reviewrecord/plan.go:209-218` |
| prepared chunk run | the tree of the chunk checkpoint | the running executable | `internal/gate/checkpoint.go:84-88`, `internal/gate/checkpoint.go:302-304` |
| complete checkpoint | the published tree of HEAD | the installed broker | `internal/gate/complete_checkpoint.go:16-31`, `internal/gate/complete_checkpoint.go:42-91` |
| landing | the published tree composed on the destination | the installed broker | `projects/benchkit.md:49-62`, ADR 0031 lines 29-31 |

- Each checkpoint, the prepared run, and the landing grade a private checkout through the prospective owner (`internal/gate/engine.go:54-67`).
- The published tree takes the spec flip or the tickets-only close, then the delivery closure (`internal/landing/published/published.go:19-38`).
- The chunk checkpoint frames no chunk id and no purpose (`internal/gate/checkpoint.go:383-387`). The complete checkpoint frames its purpose and its tip into the oracle digest only, which is not a key member (`internal/gate/checkpoint.go:388-397`).
- ADR 0031 keeps the chunk checkpoint in its own process and moves the complete checkpoint under the broker (ADR 0031 lines 31-37).
- A chunk checkpoint runs the delta set when its change since the reviewed tip is test-only and that tip has a whole green (`.bench/BENCH-reference.md:729-746`). The delta green goes to `<git-dir>/bench-last-delta` of the checkout, and no landing reads it.

### Q3. Which subject members read the live worktree during a run?

- A prospective evaluation captures the graded tree by its id and settles on the private checkout (`internal/gate/evaluation.go:58-67`). So an edit of the worktree does not reach the graded files.
- The same subject also reads the identity root, which is the worktree for each checkpoint and for the prepared run (`internal/gate/engine.go:65`, `internal/gate/checkpoint.go:305`).
- The baseline schedule enters the identity from that root (`internal/gate/subject.go:140`, `internal/gate/prospective.go:75-88`). The baseline runner digest enters the oracle digest from that root (`internal/gate/prospective.go:103-116`).
- The run compares the subject before and after the oracle, and a changed subject refuses with no verdict (`internal/gate/run_settle.go:16-39`).
- The baseline runner digest covers the Go files of `./cmd/bench` and its dependencies, `go.mod`, and three auxiliary inputs (`internal/freshness/freshness_buildinputs.go:50-90`, `scripts/go-build.inputs`).
- FT314 recorded that a plan commit during the prepared GR-C2 gate made that gate refuse with a changed subject (`git show 4645ff51c:capture/retros/gate-reuse.md`).
- `bench probe` mutates one file of the worktree, runs tests, and restores the file (`bench probe --help`).
- FT445 found that a reuse of the worktree's own executable breaks `bench probe` at its first mutant (`git show 61c5adec4:capture/retros/chunk-repair-cycle.md`, line 7).
- A prospective run builds its run binary from the private checkout when that checkout can build itself (`internal/gate/prospective.go:21-41`).

### Q4. Where do the grade identity and the run select the run binary?

- One helper answers the source root of the run binary (`internal/gate/gate.go:311-321`). It answers the graded tree for a prospective script that can build itself, and the wrapper kit otherwise.
- The run passes the resolution kind of its plan (`internal/gate/run_transaction.go:161-162`). An ordinary run in the kit resolves `gate-script`, so it builds from the wrapper kit.
- The identity passes the prospective kind whenever `.bench/gate-prospective.sh` exists in the graded root (`internal/gate/subject.go:93-100`, `internal/gate/subject.go:141-142`).
- So an ordinary run in the kit frames the build inputs of the graded tree while it builds from the wrapper kit.
- For a worktree of the kit's own repository, `KitForRoot` answers that worktree as its kit (`internal/gate/kit_root.go:5-20`). So the two answers differ only when `BENCH_KIT` names another checkout.
- The GR-C1 Coverage axis found this edge as finding C2, and the orchestrator flagged it for reviewer veto (`git show 4645ff51c^:reviews/gate-reuse.md`, line 265).

### Q5. What lifetime rules bind a process that outlives its verb?

- Every production child starts through one lifetime owner, and a new caller needs a census row (`docs/adr/0030-one-lifetime-owner-and-process-user-protection.md:7-13`).
- The owner applies the caller's normal-exit policy, which drains or preserves the process group (ADR 0030 lines 8-10). `bench worktree exec` preserves its child's group after a normal exit (`internal/worktree/exec.go:123-132`).
- Each resource that a child can still use carries process-user protection, so a removal waits until every user ends (ADR 0030 lines 15-29).
- The gate registers its oracle on the cache generation and on the run binary (`internal/gate/run_transaction.go:145-157`, `internal/gate/run_transaction.go:192-193`).
- FT372 records a nested landing gate that ran 920 s after its probe (`roadmap/FT372.md:10-11`, `roadmap/FT372.md:35`). It sealed a cache generation, so the next landing gate refused.
- The gate runs its phases serially, and ADR 0024 keeps phase overlap closed until a resource census (`docs/adr/0024-gate-phase-overlap-requires-a-resource-census.md:3-6`).

### Q6. What does the resource census show?

- The host has 16 logical cores and 23 GiB of memory on WSL2, and `go env GOFLAGS` gives `-p=4 -parallel=4` (observed in this session on 2026-10-27).
- The profile caps concurrent test delegates at two (`projects/benchkit.md:646-649`).
- The gate gives its children a closed environment (`internal/gate/subject.go:69-71`, `internal/gate/subject.go:101-111`, `internal/gate/subject.go:133-138`). It holds `PATH`, the declared `BENCH_HOME` and `HOME` (`.bench/gate-inputs.json`), and the build cache.
- The run adds the run binary, `BENCH_KIT`, the baseline entry, and the run-log entry (`internal/gate/run_transaction.go:185-191`). It also adds the span entry (`internal/gate/engine.go:239`).
- A delta run adds `BENCH_GATE_DELTA` (`internal/gate/delta_checkpoint.go:84`). That entry is a Bench input outside the grade key, so it is a precedent for a background width input.
- No Bench setting selects the `-p` width of a gate phase (`internal/gate/phases.go:220-236`). FT255 owns that shared setting and its default of two (`roadmap/FT255.md:8-13`).
- The reviewer decided on 2026-10-25 that a gate measurement runs at `-p=2` only (`roadmap/FT255.md:19-21`). The FT451 bound criterion repeats it (`.bench/commitment.json:1337-1339`).
- These workloads run beside each other on this host:
  - a whole gate for a landing, a checkpoint, or a fold
  - a fast lane at each worktree commit, which takes seconds (ADR 0017 line 13)
  - an author's own `bench test` runs and probes
  - up to two test-running delegates
- The FT255 record holds the contention evidence (`roadmap/FT255.md:27-37`):
  - seven focused race runs beside a gate raised package walls from about 70 s to 175 s
  - four occurrences of load-coupled reds that passed alone
  - one run took all capacity and stopped WSL on 2026-08-25
  - outer phase overlap extended the posture span from 50.9 s to 85.4 s
- The FT445 landing took 207 s in its test phase and 48 s in its system phase (`git show 61c5adec4:capture/retros/chunk-repair-cycle.md`, lines 11-21). The FT314 landing took 334 s and 53 s (`git show 4645ff51c:capture/retros/gate-reuse.md`, lines 7-14).
- So a whole landing gate took about 4.3 to 6.5 minutes.
- A whole `--fresh` checkpoint took 5 min 40 s. The FT445 test-only repair cycle took 40 min 23 s, and the author's own runs took about 26 min of it.

### Q7. Which work makes lane-passed tips, and which tips reach a consumer?

- `bench commit` runs the fast lane and publishes on a lane pass (`internal/commit/commit.go:150-199`). `bench worktree merge` grades its composed tree with the target's fast lane (`projects/benchkit.md:75-80`).
- A build's serial ticket authors commit on one integration source (`.agents/skills/bench-craft-delegate/SKILL.md:99`). The chunk checkpoint runs at the HEAD of that source.
- A build also makes plan commits and record commits on that source (`.agents/commands/bench-implement-spec.md:58-64`).
- A record commit changes only the review record, so the chunk key stays. A plan commit changes the spec or a ticket, so the chunk key moves.
- The prepared run is opt-in at two guidance points (`.agents/commands/bench-implement-spec.md:64-66`).
- A repository with no declared lane keeps the full-gate commit (`docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md:32-34`). So only the kit and a linked repo with a declared lane make lane-passed tips.

## Inferences

- I1. Assume that a background grade holds the checkout lock and writes the checkout cache (Q1). Then every foreground gate verb in that worktree refuses until the grade ends.
- I2. Assume that the subject of a background grade reads the live worktree (Q3). Then an author edit or a `bench probe` mutant of a Go file during the run loses the verdict.
- I2a. A subject that reads only the graded commit and machine inputs keeps its verdict.
- I3. A background grade answers a consumer only when it reproduces that consumer's tree, its running executable, and its schedule decision (Q2). The chunk consumer and the landing consumer need different trees and different executables.
- I4. A record commit does not move the chunk key (Q7). So a grade for each tip repeats a run with no new subject, unless the grades deduplicate by key.
- I5. A grade of a superseded key answers no consumer (Q6). On this host, concurrency already made reds and one crash, so that grade is pure cost.
- I6. On a runner with two to four cores, one background whole gate at the default width competes with every foreground run (Q6).
- I7. A green that ages past 60 minutes answers no consumer, so a long review or repair round can waste a background green (Q1).
- I8. The closed gate environment means that a background width needs a Bench-owned input. Whether that input enters the grade key decides whether a background green and a foreground green share a key (Q6).
- I9. A green of a worktree tip answers the landing only when the destination did not move after that tip (Q2). The landing composes on the current destination.
- I10. The consumer payload manifest does not name `.bench/gate-prospective.sh`, so a linked repo probably resolves the ordinary kind in both selections of Q4.

## Tested results

None. This run executed no gate and no measurement.

## Proposals

- P1. Let a background grade compute exactly the subject and the schedule of the next consumer of its worktree.
- P2. Read every subject member from the graded commit or from machine inputs. Let the identity read the run's own run-binary selection, so one derivation answers both.
- P3. Write the shared evidence store only, and hold a claim for each key in place of the checkout lock. Let a consumer that meets an in-flight run of its own key wait for that run.
- P4. Deduplicate by grade key and stop the grade of a superseded key. Run at most one background grade at once on the machine, at `-p=2`, below the foreground.
- P5. Retire `--prepare` and its two guidance points when the chunk consumer is in scope.

## Option table: where a background grade keeps its state

| option | consequence | evidence |
|---|---|---|
| the checkout cache and the checkout lock, as the prepared run does | every foreground gate verb in the worktree refuses while the grade runs | `internal/gate/run_transaction.go:94-101`, `internal/gate/engine.go:219-226` |
| the shared store and a claim for each key | a foreground verb meets the grade only on the same key, where it can wait | `internal/gate/engine.go:147-149` |

## Diagram

```
 lane pass (bench commit, bench worktree merge)
        | starts
        v
 background grade  --writes green-->  evidence store
 (private checkout of the tip)        (git common dir, one file per grade key)
                                              ^
                                              | looks up the grade key
 consumer: chunk checkpoint or landing -------+
 (checkout cache and checkout lock in the admin dir of its own checkout)
```

## Contradictions

- C1. The roadmap row says that the grade identity and the ordinary run select the run binary from different sources. The comment at `internal/gate/gate.go:311-313` says that both read one helper. Both statements hold, because both call one helper with different kind arguments (Q4).

## Unknowns

- U1. This run did not find which subject member moved when a plan commit refused the prepared GR-C2 gate. A plan commit changes no build input, and the drained record names only the refusal.
- U2. This run did not verify that Go reads its `-p` width from the configuration file under `HOME` inside a gate run. The open learning at `capture/learnings.md:7-10` has no repro, and this run did not open the Go documentation.
- U3. This run did not measure the cost of one background gate beside foreground work at `-p=2`.
- U4. No record counts how often the destination moves between the last commit of a worktree and its landing.
- U5. No record counts how many spec-less or tickets-only kit landings take the whole gate and not the reader set. With U4, this count sets the hit rate of the landing consumer.
- U6. This run did not verify that the lowest I/O priority takes effect on WSL2.
- U7. This run did not verify whether a whole green and a red with a landing context are the only extra writes of a run. The map review found the reader table and the reader hold (`internal/gate/run_transaction.go:249`, `internal/gate/run_transaction.go:256-258`).

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, inferences, tested results, and proposals stay in separate sections.
- [x] Each question has one section.
- [x] The option table records each consequence with a citation.
- [x] The contradiction and the unknowns stay.
- [x] A diagram shows the store and lock relation.
- [x] The unknowns state what this run could not verify.
- [x] Each load-bearing claim cites a path and a line, a commit, or an observation of this session.
- [ ] Paired independent reads: not done, by reviewer direction.
- [x] An independent map review on Opus at high effort spot-checked the Q1 to Q4 citations against the code, and this asset carries its corrections.

## Validation plan

- The spec writer reproduces U1 with one prepared run and one plan commit before the subject rows lock.
- The first build ticket measures U3 at `-p=2` under the bound that the map decides.
- U2 graduates through the repro that `capture/learnings.md:9` names.
- The census ticket records U6, because its log must show that the priority took effect.
- The spec writer counts U5 from the landing history before the landing chunk.
- The spec writer lists every write of a whole run before the write rows lock, to close U7.
