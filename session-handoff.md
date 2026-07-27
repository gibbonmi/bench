# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` @ `29ea46c` — 16 unpushed commits ahead of `origin/main`; `IDEAS.md` dirty with 3 parked lines
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented, deliberately unretired)
Gate: green at `08482c6`; every commit since is doc-only

## State

- **The worktree-pool investigation is delivered and its decisions are closed
  (reviewer sign-off 2026-07-27).** Root cause: `bench worktree release` matches
  only the exact plaintext request string that created the assignment
  (`internal/worktree/ownership.go:344`); the ledger stores a one-way digest and
  the harness hook derives the request from the Claude session id, so a dead
  session's worktrees are structurally unreleasable. Confirmed in code and
  reproduced through the accused command. Decided: orphans route to
  `bench worktree clean` by design — a request-derivation override for `release`
  is rejected as voiding the ownership model.
- **All 20 pool branches are disposable.** Twelve are landed and clean; the three
  FT87s3 branches and one agent branch are drafts of work main shipped reshaped
  (bench's `landed=false` is the FT98 patch-id gap, seen in the wild); the FT91
  concurrency-budget arm is abandoned per the retarget — sign-off covers
  deleting its branch. The four `fix/ft86-review-*` worktrees need their assign
  branch recreated before `clean` accepts them. One unregistered empty directory
  remains for plain `rm`.
- **The cleanup sequence below has not been run.** It is reviewer-run: branch
  deletion and `rm` are guard-refused from agents, and `clean --apply` was left
  to the same pass. `clean` writes a recovery ref before destroying anything
  dirty and deletes branches itself only with landed proof.
- **Ledger residue outlives the pool drain.** Assignment `72b9811f`
  (ft91-gate-phase-split) is `active` with its tree already gone — nothing today
  compacts an active record with a missing tree; FT147(b) covers it. The ~21
  recovery refs (the "preserved" wall at session start) retire per-ref after the
  drain: `bench worktree recovery <ref>` to inspect, `--apply <fingerprint>` to
  retire when the plan allows.
- **FT147 is signed off and awaits the drain.** The leak: kit prose orders
  worktree creation 12× for every retirement instruction (`release` is named in
  no guidance file), assignments have no timestamp/lease/reaper and the resume
  classifier hard-retains `active`, and FT98's landed proof misses reshaped
  commits. Fix shape as signed off: (a) prose — `craft-delegate` close-out duty
  (the coordinator releases each worktree it cut, at done-claim acceptance),
  `bench-implement-spec` stop-short names a retirement owner, BENCH.md inventory
  names the subcommands; (b) code — created-at timestamp plus an `orphaned`
  classifier verdict surfaced by resume as ready-to-run clean commands;
  (c) rides on FT98. Prose edits still go through `craft-synthesis` as a build.
  `/bench-what-next` should convert the parked release-refusal idea line into
  this row rather than promote it unverified.
- **Half of the FT146 row still needs your removal verdict** — the artifact
  contract-test half was a mis-attribution; nothing is left to build there.
- **`ft91-gate-phase-split` stays unretired on purpose** — retiring it destroys
  your veto surface on stories 4, 5, and 9.
- **Push needs `bench gate pin` first** — interactive TTY, so it is yours.
- Drain pending: 3 parked ideas, 1 open learning.

## Cleanup sequence (reviewer-run, in order)

```bash
POOL=~/.bench/worktrees/bench-2826441890

# 1. The 16 assign-branch worktrees: plan, then apply with the printed fingerprint.
for w in 2e244b92*-71e2c598* 3262b7b3*-ef31c1dc* 3be84fa2*-b7fc1bae* 5096abf2*-e48a7258* \
         7c14860e*-1d9a75cf* 877c99aa*-2cd80392* 8ddd3feb*-b0ba1a29* 9ff77f4a*-cab958cc* \
         cc14983b*-a4de1307* fc05fad7*-9566bd64* ebd6824d*-43613d5d* b6192cf8*-0507f12a* \
         88203ee7*-5d7a0f1b* ccb9f04c*-f632161c* ad1d12a2*-d73215b5* e8e880c7*-324fd4dc*; do
  p=$(echo $POOL/$w)
  fp=$(bench worktree clean "$p" | awk -F, 'NR==2{print $6}')
  [ ${#fp} -eq 64 ] && bench worktree clean "$p" --apply "$fp" || echo "SKIPPED (retain?): $p"
done

# 2. The four fix/ft86-review-* worktrees: restore the assign branch so the
#    registration matches the ledger, re-plan (the retain should flip), apply.
declare -A FIX=(
  [0431f7e9aa9247fc30c6479b0b0b0f0e-48e207f910fb3750ea69c997e9cc9471]=bench/assign/0431f7e9aa9247fc30c6479b0b0b0f0e/48e207f910fb3750ea69c997e9cc9471
  [30a9b2979dacaef03a1e0f30b419a0da-b8c88a40709f436a455638a5c8f59da0]=bench/assign/30a9b2979dacaef03a1e0f30b419a0da/b8c88a40709f436a455638a5c8f59da0
  [478c3ac90c1c2d3628a111a7710d7b72-3f0111d86e5be84a13ea97b49cad20ee]=bench/assign/478c3ac90c1c2d3628a111a7710d7b72/3f0111d86e5be84a13ea97b49cad20ee
  [cfa2d2942c5089c406ada0811dce4735-d72ea6f38ca718dba5a6f3103bf305ff]=bench/assign/cfa2d2942c5089c406ada0811dce4735/d72ea6f38ca718dba5a6f3103bf305ff
)
for w in "${!FIX[@]}"; do
  git -C "$POOL/$w" checkout -b "${FIX[$w]}"
  fp=$(bench worktree clean "$POOL/$w" | awk -F, 'NR==2{print $6}')
  [ ${#fp} -eq 64 ] && bench worktree clean "$POOL/$w" --apply "$fp" || echo "SKIPPED (retain?): $POOL/$w"
done

# 3. The unregistered empty directory (outside bench's authority).
rm -r "$POOL/eb04fb610b43d73caa1f0d9e4a4e5f4d"

# 4. Branches clean leaves behind (it deletes only proof-landed ones itself).
#    -d proves its own safety; the -D set is superseded drafts + the abandoned
#    FT91 arm, all signed off 2026-07-27.
git branch -d fix/ft86-review-coverage fix/ft86-review-singlesource \
              fix/ft86-review-failclosed fix/ft86-review-git
git branch -D \
  bench/assign/88203ee7d603436167ef591631a6d1da/5d7a0f1b9b784de9b8a7717a04217ad7 \
  bench/assign/ccb9f04ce873be7000146c9486df1584/f632161c5b91ce924d471c7842aa3a4f \
  bench/assign/ad1d12a248a7e4e4ce87f8839970ac77/d73215b551ad2003b9c623b113be1068 \
  bench/assign/e8e880c740767330755ca5f6f22d7fd3/324fd4dc5335ae6b08cf65bd3b2e52eb

# 5. Retire recovery refs once the pool is empty: inspect each, apply when the
#    plan says it can retire. (Per-ref by hand — payload verdicts differ.)
git for-each-ref 'refs/bench/recovery/**' --format='%(refname)'
# then per ref:  bench worktree recovery <ref>   →   bench worktree recovery <ref> --apply <fingerprint>
```

## Next command

`/bench-what-next`

(after running the cleanup sequence above)

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
