# One change, one full grade

Status: ready

## Destination

A landed change pays the whole-project gate once, at the landing. A worktree
commit records fast-lane evidence and never a whole-project verdict. Scope: `bench commit` in a worktree and `bench worktree land`. The floor of
that one run is the second map, `specs/worktree-test-floor/decisions/worktree-test-floor.md`.
This map supersedes FT215's headline, a changed-package gate path, which prior
rulings close.

## #1: What makes the gate run twice for one change?

Blocked by: none
Type: Grill

### Question

The last build paid ten full runs for four landings. Which step makes the second
run unavoidable today?

### Answer

Resolved 2026-08-25: the landing's authored transition. The landing rewrites the
spec's `Status: staged` line and removes the tickets-only folder after the source
was graded, so the composed tree is new. A second cause rides with it. The landing keys
evidence to `main`'s identity root and baseline runner (ADR 0016). The worktree
commit keys to its own, so even a byte-identical tree does not reuse across
checkouts. See `specs/one-change-one-grade/decisions/assets/ft215-prior-rulings.md`.

## #2: Where does the one full grade live?

Blocked by: #1
Type: Grill

### Question

Option A: the landing reuses the source's green and applies its own transition
without a run. This needs ADR 0016 reopened, because the source's evidence keys
to the worktree's runner. Option D: the worktree commit stops running the
whole-project gate, and the landing is the one full grade. This obeys ADR 0016
and `spec-build-review-gate-cadence` #2 and #5. Option K: keep both runs and cut
only the floor. Recommendation: D.

### Answer

Resolved 2026-08-25: option D. The worktree `bench commit` stops running the
whole-project gate. The landing is the one full grade, and it grades the tree
with the flip already applied. ADR 0016 stays closed.

## #3: What evidence does a worktree `bench commit` require?

Blocked by: #2
Type: Grill

### Question

Under D, what does `bench commit` run before it publishes onto the worktree
branch? Recommendation: a fast lane of gofmt (shipped), the prose check on the
changed Markdown, `go vet`, and `go build ./...`. The lane records its own record
class and never green.

### Answer

Resolved 2026-08-25: the fast lane. `bench commit` runs gofmt on the changed Go
files, the prose check on the changed Markdown, `go vet`, and `go build ./...`
before it publishes onto the worktree branch. The lane records its own record
class and never a green verdict.

## #4: What iteration evidence does the agent run before the landing?

Blocked by: #2
Type: Grill

### Question

Under D, does the agent still run `bench gate` in the worktree before it lands?
If yes, the second run remains. Recommendation: focused `go test` on the touched
packages as iteration evidence, the fast lane at commit, and the landing as the
one full gate. Invariant 1 stays as written, because the landing runs `bench gate`.

### Answer

Resolved 2026-08-25: the agent does not run the whole-project gate in the
worktree as a step of the workflow. Iteration evidence is focused `go test` on
the touched packages, then the fast lane at commit. The landing runs the one
`bench gate`, so invariant 1 keeps its wording.

## #6: One spec or two?

Blocked by: #2
Type: Grill

### Question

The cadence slice (#2–#4) and the floor cut are independently useful.
Recommendation: two specs, cadence first.

### Answer

Resolved 2026-08-25: two specs. This map is the cadence slice, and it is the
first spec. The floor cut has its own map, `specs/worktree-test-floor/decisions/worktree-test-floor.md`,
and it is the second spec.

## #7: Multi-ticket specs under D

Blocked by: #2
Type: Grill

### Question

Under D, ticket commits stack on fast-lane evidence only, and the landing gate is
the single proof. Accept that a red surfaces at the landing, not at the ticket
commit? Recommendation: accept; one run costs about 100 s.

### Answer

Resolved 2026-08-25: accepted. Ticket commits stack on fast-lane evidence, and
the landing gate is the single proof. A red that surfaces at the landing routes
back to the worktree and costs one run.

## #8: Does `bench shift` keep the whole-project gate as its loop oracle?

Blocked by: #2
Type: Grill

### Question

A shift is unattended work. Its loop and the stop hook (`bench stop-verdict`)
run `bench gate` and key the verdict to the tree. Under D, a shift that commits
and lands pays two runs. Recommendation: unchanged; the shift keeps the
whole-project gate as its loop oracle, and its second run is accepted because no
reviewer sits between its iterations.

### Answer

Resolved 2026-08-25: unchanged. The shift keeps the whole-project gate as its
loop oracle and its stop rule. A shift that commits and lands pays a second run,
and that cost is accepted for unattended work.

## #9: How does invariant 4's "commit on green" read under D?

Blocked by: #2
Type: Grill

### Question

`.bench/BENCH.md` says "Commit on green, never on red" and "Done means
`bench gate` exits zero". Under D a worktree commit publishes on fast-lane
evidence. Recommendation: "green" in invariant 4 means the landing's gate, and a
worktree commit requires fast-lane green. The spec-writer drafts the exact
sentence for veto at spec sign-off.

### Answer

Resolved 2026-08-25: "green" in invariant 4 means the landing's gate. A
worktree commit requires fast-lane green. The spec-writer drafts the exact
sentence in `.bench/BENCH.md`, and the reviewer vetoes it at spec sign-off.

## #10: What does `bench commit` run in a linked repo?

Blocked by: #3
Type: Grill

### Question

`bench commit` is one code path, and the kit's fast lane is Go-specific. A
linked repo owns its gate and may have no Go. Recommendation: a project declares
its fast lane, and the kit ships its own declaration. A linked repo with no
declaration keeps today's full-gate commit, so D applies only where a lane is
declared.

### Answer

Resolved 2026-08-25: a project declares its fast lane, and the kit ships its own
declaration. A linked repo with no declaration keeps today's full-gate commit.
D applies only where a lane is declared.

## #11: What does `bench commit --dry-run` run?

Blocked by: #3
Type: Grill

### Question

Today `--dry-run` gates the composed snapshot and commits nothing.
Recommendation: `--dry-run` runs exactly what the commit runs, the fast lane;
`bench gate` remains the whole-project check on demand.

### Answer

Resolved 2026-08-25: `--dry-run` runs exactly what the commit runs, the fast
lane, and commits nothing. `bench gate` remains the whole-project check on
demand.

## Not yet specified


## Spec-writer discretion

- The seam that runs the prose check on the changed Markdown outside the `test`
  phase. The name and the store of the fast lane's record class.
- The exact sentence for invariant 4 in `.bench/BENCH.md` (#9), for veto at spec
  sign-off.
- The shape of the fast-lane declaration (#10), as long as the kit ships its own
  and an undeclared linked repo keeps the full-gate commit.

## Out of scope

- A changed-package or diff-derived gate path. Four maps rule it unsound, and the
  profile pins no per-package loop.
- A Markdown-only or capture-only lane. FT183 #1 retired it as unreachable, and
  `prose-mechanics` reads Markdown inside the gate.
- The flip's authorship. The landing verb stays the one flip and close author (FT113).
- The ticket-local evidence machinery of `spec-build-review-gate-cadence`. It waits
  on FT173 and FT130.
- Reopening ADR 0016's evidence key. #2 chose option D.
- The `internal/worktree` test floor. `specs/worktree-test-floor/decisions/worktree-test-floor.md` owns it.

## Sources

- Path: `specs/one-change-one-grade/decisions/assets/ft215-prior-rulings.md`
  Supports: #1's answer and the option framing of #2. It records the four closed rulings against diff-scoped gating, the retired Markdown lane, ADR 0016's evidence key, and the current gate shape with code citations. Produced 2026-08-25 by one read-only research delegate and checked against the tree.
  Drift: re-verify if `internal/gate/subject.go`'s key framing, `internal/landing/landing.go`'s transition, or ADR 0016 changes before the spec reads this map.
