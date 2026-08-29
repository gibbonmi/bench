# Worktree exec comfort (FT254)

Status: shaping

## Destination

An agent does multi-step work in an assignment worktree through `bench worktree exec`
and its sibling verbs, with no raw command on the worktree path. Roadmap row FT254
holds the decided items; this map records the evidence behind each item and the
reviewer choices the row leaves open. The row's items are:

- The exec help shows stdin with a heredoc example.
- `bench worktree show <target> <rev>:<path>` prints the blob bytes unchanged.
- Every exec failure after target resolution prints the absolute worktree path. An
  earlier refusal names `bench worktree list`.
- `bench worktree resolve <target> <path>...` stages caller-edited paths without
  conflict markers and reports the remaining unmerged paths.
- The follow-on guard names the token that caused a refusal. A heredoc body with no
  `bench` word is not a follow-on. A non-bench follow-on step under
  `bench worktree exec` is accepted.
- `bench worktree exec --env KEY=VALUE` supplies worktree-local inputs.
- An assignment retained after release stays resolvable by `exec`, or release names
  the recovery verb.
- An exec-native gate form replaces `bash bin/bench.sh gate`. A declared environment
  and tool-closure form supports a focused package reproduction.
- `bench worktree create --from <target>` starts at a sibling's tip. The preflight
  `base-current` red names `bench worktree merge` as its `next=` remedy.
- A worktree's own grammar runs through `bench worktree exec <label> -- ./dist/bench`
  or a `bench worktree build <label>` form.
- A flag after `--` belongs to the child. The child's usage error, stderr, and exit
  code pass through unchanged.

## #1: What does `bench worktree exec` do today?

Blocked by: none
Type: Research

### Question

Which row claims already hold in `internal/worktree/exec.go` and `internal/usage`?
Cover stdin, the child environment, every error path, a flag after `--`, exit-code
pass-through, and the resolution of an assignment that release retained.

### Answer

Answered 2026-08-29 by a read-only research delegation; the asset in Sources holds
the citations, the probe transcript, and the test inventory.

Row claims against the code:

- stdin passes through: holds. `os.Stdin` reaches `cmd.Stdin` unwrapped. No test pins
  it, and the follow-on guard refuses the heredoc that would exercise it (see #2,
  case iv).
- Every failure after target resolution prints the absolute worktree path: does not
  hold. The `cmd.Start` failure prints only the `os/exec` error; the cancel path
  prints nothing and exits 130.
- An earlier refusal names `bench worktree list`: does not hold. `printTargetRefusal`
  prints `bench worktree exec: <reason>` with no next verb. The retained-release path
  already proves the `next=` shape, so exec can compose it.
- A flag after `--` belongs to the child: partly holds. Flag recognition stops at
  `--`, but the empty-positional guard in `usage.Parse` still fires past the
  terminator, so `rg -N '' README.md` is refused as `unknown argument: ""`. This is the
  `ft214-craft-spec-visit` occurrence, now reproduced. The row's intent covers every
  argument after `--`, not only a flag.
- The child's usage error, stderr, and exit code pass through unchanged: partly
  holds. Both pass through, but a child's exit 2 collides with the grammar refusal's
  exit 2, so only the stderr prefix `usage: bench worktree exec` separates them. This
  is now #11.
- An assignment retained after release stays resolvable: does not hold. `drain-batch`
  lists as `active` with `tree missing`, `list` advertises an exec row for it, and
  exec refuses with `owner marker does not match assignment <id>`, a component the
  operator cannot act on. Every release retain does print the absolute path and a
  `next=bench worktree release ...` line.

The child environment strips `BENCH_KIT`, `BENCH_WRAPPER`, `BENCH_RUN_BINARY`, and the
inherited `BENCH_HOME`. It sets `BENCH_HOME` from the verb. When the worktree's own
`bin/bench.sh` exists, it sets `BENCH_WRAPPER` to that file. A caller's leading
assignment (`FOO=bar bench worktree exec ...`) already reaches the child. So
`--env KEY=VALUE` adds a form the guard can see, not a capability.

## #2: How does the follow-on guard classify a call?

Blocked by: none
Type: Research

### Question

How does `internal/benchguard` split a call into segments? Which segment counts as a
Bench segment? Does the refusal path know the offending segment? Which verdict does
each row case get today? The row cases are a non-bench follow-on under `exec`, a
heredoc body with no `bench` word, and a heredoc body with a `bench` word.

### Answer

Answered 2026-08-29 by a read-only research delegation; the asset in Sources holds
the citations and the nine-case classification table.

The classifier reads `.tool_input.command`. It strips heredoc bodies but keeps the
`<<` operator, folds quoted strings into one word, and splits segments on control
operators. At the first Bench-headed segment it scans the whole token stream for any
redirection or control operator. So an operator before the Bench segment also
refuses.

The refusal string is a constant, and `Classify` returns a bare bool. The operator
token is in scope at the refusal point and is discarded. To name it is a return-type
widening plus one call site, not a redesign.

Row claims against the code:

- "The guard names the token" does not hold today.
- "A heredoc body with no `bench` word is not a follow-on" holds today and is pinned
  by `TestClassifyFollowOns`.
- "A non-bench follow-on under `exec` is accepted" does not hold. A quoted follow-on
  inside `sh -c '...'` passes; an unquoted `;` or `&&` after the exec child refuses.

Census records nothing for a `bench`-headed call, so an exec child's verb head never
counts as a raw call (`TestRecordSkipsABenchCall`). Two reviewer questions come out of
this ticket, now recorded as #9 and #10.

## #3: How does a worktree run the gate, its own build, and a focused reproduction?

Blocked by: none
Type: Research

### Question

Why does the worktree usage list `bash bin/bench.sh gate --fresh`? Which checkout's
binary does a bare `bench` serve inside an exec child? Which environment does a
hand-run focused package test need, and which existing verb already builds or tests?

### Answer

Answered 2026-08-29 by a read-only research delegation; the asset in Sources holds
the citations and the seam table.

The usage trailer `bash bin/bench.sh gate --fresh` is a hand-added string in
`internal/usage/worktree.go`, not a verb, and the same line is a public help row on
`gate`. Its reason is the wrapper-owned run binary: `.bench/gate.sh` refuses a
wrapper-less entry and names that invocation. No ADR or map explains the worktree
placement.

Row claims against the code:

- An exec-native gate form exists: holds in effect, not as a verb. `bench worktree
  exec <label> -- bench gate` grades the worktree, because the child's re-pointed
  `BENCH_WRAPPER` makes the gate build its own binary from the graded tree. The
  trailer still advertises the raw path.
- A declared environment and tool-closure form for a focused reproduction: partly
  holds. `bench test` materializes `BENCH_RUN_BINARY`, `BENCH_KIT`, `GOCACHE`, and the
  conformance pair for ordinary packages and named checks. It has no build-tag form,
  so `-tags=system ./internal/systemtest` has only a hand-built environment. This is
  now #12.
- A worktree's own grammar through `exec <label> -- ./dist/bench`: does not hold.
  Nothing builds a worktree-local `dist/bench`; `runbinary.Build` is exported with no
  CLI surface.
- The wrapper on PATH serves the main checkout's build: holds. `bench_binary_path`
  re-anchors `dist/bench` at the main tree when the worktree has none, by design. So a
  bare `bench` in an exec child runs the main checkout's binary for any new grammar.
  `bench gate` and `bench test` in that child own a private build of the graded tree.

`scripts/go-build.sh` reads no `BENCH_HOME`; the occurrence's `BENCH_HOME` need came
from `runbinary`'s build environment, which refuses a non-absolute home. `GOFLAGS`
appears nowhere in the tree; the `-p=4` cap is a host `go env -w` setting.

## #4: What do merge conflicts, create, and preflight leave behind today?

Blocked by: none
Type: Research

### Question

What state does `bench worktree merge` leave on a conflict, and what does it name as
recovery? How does Git expose an in-progress merge, cherry-pick, or revert in a
linked worktree? What base does `bench worktree create` use? Can `--from <target>`
compose the merge verb's target resolution? What does the preflight
`base-current` red print? Which cuts did the retired worktree-merge spec decline?

### Answer

Answered 2026-08-29 by a read-only research delegation; the asset in Sources holds
the citations, the Git doc references, and the verbatim retired-spec passages.

`bench worktree merge` composes through `git merge-tree --write-tree` in the root
repository. It never touches the worktree index. A conflict is a refusal with a
`refusal_paths` table and an empty `next` field. It leaves no `MERGE_HEAD` and no
change. The landing's conflict refusal names the hand route verbatim:
`git -C '<path>' merge '<destination>'`, then `bench commit`, review, and land.

So a conflicted index only ever comes from a raw Git command. The worktree-merge
spec's "Won't handle" list placed that command outside the command boundary.
`resolve` therefore stages into a state that a raw command started, in every case.

Git exposes the in-progress state per worktree through
`git -C <worktree> rev-parse --git-path MERGE_HEAD`. The same call resolves
`CHERRY_PICK_HEAD`, `REVERT_HEAD`, and `sequencer`. The tree reads none of them.
`--git-path` serves only `index`, `hooks`, and the lease.

`git add <path>` collapses the stages from the working-tree bytes. It does not look
for conflict markers. A `--continue` needs the state file plus a resolved index.
`git diff --name-only --diff-filter=U` lists unmerged paths once each.

`bench worktree create` starts at the default branch tip, with `HEAD` as the
fallback. It passes the start through `requestedStart` of `createAt`, which
`--refresh` already exercises. The merge verb's `mergeSiblingTip` resolves a target
to a verified tip. So `create --from <target>` composes two existing seams. The
ledger records only `Start`. No field names a source assignment.

The preflight grammar has no bare form. The row's "bare" means
`bench preflight <mode> <slug>` with no `--base`. The red text is
`default branch tip is not an ancestor of HEAD`. `CheckResult` has no remedy field,
and the table is `check,verdict,detail`. So a `next=` needs a schema change, not a
string edit. This is now #13, and FT162 is the sibling row on false `base-current`
reds.

The worktree-merge spec retired at `049d12b0`. Its Out of scope names the
preflight `next=` remedy and `create --from` as declined cuts, in the words FT254
uses. The census spec retired at `bf4b1f9e`. Its decision map recorded four
exec-comfort decisions, quoted in the asset. None of these shipped, and `--env` has
no retired provenance. FT263 (`bench commit` on a pending `MERGE_HEAD`) is open with
`Next: decide`.

The four retired census decisions are:

- #10: one stdin help line and one acceptance row.
- #11: `resolve` stages, reports, and continues the operation; the continue step is
  flagged for veto.
- #12: `show` prints the blob bytes unchanged and returns Git's exit code.
- #13: a `worktree: <absolute path>` stderr line before every post-resolution
  failure, and `bench worktree list` as the next action before resolution.

`bench worktree list` shows `tree missing` when `os.Stat` fails on the path. For an
assignment row it still prints `path` and `exec` actions and no recovery. Only a
foreign row gets the orphan action.

## #5: Does `resolve` continue an operation Bench did not start?

Blocked by: #4
Type: Grill

### Question

`bench worktree resolve <target> <path>...` stages caller-edited paths and reports the
remaining unmerged paths. Does it also continue a merge, cherry-pick, or revert that a
raw Git command started, or does it refuse and name the raw continuation?

The retired census map answered "continue with the message unchanged and print the
commit". It flagged that step for veto, because it makes Bench own a Git state it did
not create. Research #4 adds the other side. If `resolve` stops at staging, the next
step is `bench commit` on a pending `MERGE_HEAD`. FT263 records that this path
authored a one-parent commit. So the no-continue answer is safe only once FT263
lands, and the continue answer overlaps FT263's contract.

Recommendation: no continue. `resolve` stages and reports. When no unmerged path
remains, it names the exact `--continue` form for the live operation. FT263 then
owns the commit contract in one place.

### Answer

— (open)

## #6: Which recovery does a retained assignment get?

Blocked by: #1
Type: Grill

### Question

`bench worktree release` can retain an assignment. Does `exec` stay resolvable
against the retained tree, or does release name the reachable recovery verb? Research
#1 shows that release already names `next=bench worktree release ...` for every retain
reason. The open edge is an assignment whose tree is gone. `bench worktree list`
shows it as `active` and offers an exec row. Exec then refuses on the owner marker.
Does `list` stop advertising exec for a missing tree? Does the exec refusal name the
missing tree and the recovery verb?

### Answer

— (open)

## #7: Which form runs a worktree's own grammar?

Blocked by: #3
Type: Grill

### Question

Choose the convention `bench worktree exec <label> -- ./dist/bench <verb>`, or a
`bench worktree build <label>` verb that builds the worktree's tree and serves it.
Research #3 narrows the choice. The `./dist/bench` convention needs a producer that
does not exist. A `build` verb is a new CLI surface over the exported
`runbinary.Build`. Both are needed only when a session must run new grammar as a
bare `bench`. `bench gate` and `bench test` under exec already grade the worktree.

### Answer

— (open)

## #8: One spec, or independent slices?

Blocked by: #1, #2, #3, #4
Type: Grill

### Question

The row holds several independently useful behaviors. They are the exec surface, the
guard, the gate and build forms, `show`, `resolve`, and `create --from` with the
preflight remedy. Ship one bundled spec, or split into ordered slices?

Research narrows the costs. Two items fit the light path now. `bench worktree list`
stops offering `path` and `exec` on a `tree missing` row. The merge conflict refusal
carries the `next=` that the landing refusal already composes.

One coherent exec-surface slice holds the exec help line, the `worktree:` path line,
the `bench worktree list` next action, the `--` empty-argument fix, and the guard
token. The retired decisions are its provenance. `show` is small and independent.
`resolve` waits on #5 and FT263. `create --from` composes existing seams. The
preflight `next=` needs a schema change (#13).

Recommendation: three slices. Slice 1 is exec surface plus guard plus `show`. Slice 2
is `create --from`, the gate usage trailer, and #7's build form. Slice 3 is `resolve`,
sequenced after FT263. The two light-path items land ahead of any slice.

### Answer

— (open)

## #9: Which token does a refusal name when the operator precedes the Bench segment?

Blocked by: #2
Type: Grill

### Question

In `cat roadmap/FT254.md; echo ---; bench maps`, the first operator in the stream is
the `;` after `cat`. The operator adjacent to the Bench segment is the `;` after
`echo`. Which one does the refusal cite: the first in the stream, or the one adjacent
to the Bench segment?

### Answer

— (open)

## #10: Does a non-bench segment before an exec call pass?

Blocked by: #2
Type: Grill

### Question

The row accepts a non-bench follow-on step after an exec child. So
`bench worktree exec L -- cp a b; cp b a` passes. Does
`cp a b && bench worktree exec L -- true` also pass? Or does a segment before the
Bench segment still refuse?

### Answer

— (open)

## #11: Which exit code does an exec grammar refusal use?

Blocked by: #1
Type: Grill

### Question

The child's exit code passes through, and a child's usage error exits 2. The exec
grammar refusal also exits 2, so the code alone cannot separate the two. Keep exit 2
for both and rely on the `usage: bench worktree exec` stderr prefix, or move the
grammar refusal to another code? The second option changes an established exit
contract.

### Answer

— (open)

## #12: Does `bench test` gain a system-suite form?

Blocked by: #3
Type: Grill

### Question

`bench test` owns the declared environment for ordinary packages and named
conformance checks. The tagged system suite (`-tags=system ./internal/systemtest`)
runs only inside the gate or with a hand-built environment. Give `bench test` a
system form (a `--check`-style name or a tag flag), or leave the system suite
gate-only?

### Answer

— (open)

## #13: Does the preflight `next=` remedy stay in scope?

Blocked by: #4
Type: Grill

### Question

The `base-current` red names no remedy, and the preflight check table has no remedy
field. To name `bench worktree merge` needs a `next` column on every check row, or
a special case in the `base-current` detail string. Widen the check schema, special-
case the one string, or decline the cut again as the worktree-merge spec did?

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

- How `Classify` carries the offending token to the refusal printer, given that the
  token is already in scope at the refusal point.
- Which existing helper `create --from <target>` composes to resolve a sibling tip,
  and how the self-reference check adapts to a target that does not exist yet.
- Whether the `worktree: <absolute path>` line composes the `next` field that the
  retained-release refusal already proves.

## Out of scope

- A fold of an uncommitted sibling snapshot through the merge verb, an interactive
  conflict-resolution mode, and a tagged system journey for the merge verb. The
  worktree-merge spec declined these three cuts, and FT254 did not fold them.
- A `bench worktree pick` verb or any verb that starts a merge, cherry-pick, or
  revert: the retired census map placed it out of scope.
- The `bench commit` contract on a pending `MERGE_HEAD`: FT263 owns it.
- False `base-current` reds: FT162 owns them.
- Conflict resolution inside `bench worktree merge`: the verb refuses and names the
  paths; the hand resolution stays raw Git.

## Sources

- Path: `decisions/assets/ft254-exec-comfort-research.md`
  Supports: #1 through #4 and the factual premises of #5 through #13. Four read-only research delegations ran 2026-08-29 at `a712f84f`. The coordinator spot-checked two or three citations per delegation against the tree. Its #4 section quotes the retired worktree-merge spec (`049d12b0^:specs/worktree-merge/spec.md`) and the retired census decision map (`bf4b1f9e^:specs/exec-census/decisions/exec-census.md`) verbatim. Those quotes settle #8's declined cuts and the veto flag behind #5.
  Drift: re-resolve the line citations before `/bench-write-spec` reads this map if `internal/worktree`, `internal/benchguard`, `internal/usage`, `internal/preflight`, or `bin/bench.sh` change. The quoted retired passages are commit-pinned and do not drift.
- Path: `roadmap/FT263.md`
  Supports: #5's dependency on the `bench commit` `MERGE_HEAD` contract.
  Drift: re-read when FT263 leaves `Next: decide`.
