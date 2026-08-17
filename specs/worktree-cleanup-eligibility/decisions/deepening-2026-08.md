# Deepening survey — FT216, FT217, and FT218

Status: ready

## Destination

Deliver the three remaining source-verified deepening opportunities from
`capture/architecture-review-20260817T104714.html`: one worktree eligibility
verdict (FT216), one adopt lifecycle decision (FT217), and named Git readers for
recurring repository facts (FT218). FT216 and FT217 enter the normal spec path;
FT218 takes one light-path ticket per qualifying reader. The gate transaction,
verdict registry, shift objective, green marker, and skills-index reader are
already deep and are outside this map.

## #1: Route each deepening and preserve the refactor exit proof

Blocked by: none
Type: Grill

### Question

Which workflow route owns each candidate, and what proves a behavior-preserving
deepening did not become an unreviewed feature change?

### Answer

FT216 and FT217 each receive their own `/bench-write-spec` pass from this map.
The FT216 pass compiles the map into
`specs/worktree-cleanup-eligibility/decisions/`; FT217 and FT218 cite that one
copy until the remaining decisions are re-homed before FT216 retires.
FT218 uses one light-path ticket per reader that recurs at three or more
production sites. Every route carries the same exit proof: the pre-existing
suite passes with test logic unchanged; mechanical test renames are permitted,
but changing an assertion or expected outcome stops the refactor and routes the
behavior delta through a separately reviewed feature or defect path.

## #2: FT216 — one worktree eligibility verdict

Blocked by: #1
Type: Grill

### Question

What does the eligibility module decide, which consumers read it, and may this
deepening repair precedence anomalies it exposes?

### Answer

Characterize the current effective refusal precedence first, one reachable
Action/ReasonCode outcome at a time, and preserve it exactly. One deep
eligibility module then owns the ordered ownership, assignment, lock, lease,
landedness, recovery, tracked-state, nested-state, and ignored-residue evidence
and returns one decided verdict with that evidence. Explicit cleanup projects
the operator-facing plan; automatic cleanup is a stricter reading of the same
verdict and does not infer eligibility from formatted strings.
`ApplyExplicit`, `PlanAutomatic`, `LandCommand`, and `ReleaseCommand` read
that decision. `DiscardBranch` remains derived after landedness and cannot be
reached by automatic cleanup. Any precedence anomaly is a separate
`/bench-debug` repair, never bundled into this deepening. ADR 0005 is rewritten
with the build to name the eligibility verdict while retaining its conjunctive
ownership and preservation posture.

## #3: FT217 — one adopt lifecycle decision

Blocked by: #1
Type: Grill

### Question

Which adopt verbs consume the lifecycle decision, what stays local to each verb,
and does the deepening add a dry-run surface?

### Answer

`link`, `setup`, `upgrade`, and `unlink` all execute one immutable
inventory comparison that decides the ordered add, change, remove, and preserve
operations. Each verb retains only its distinct validation, safety checks,
reporting, and atomic filesystem implementation. Tests reach the public command
seam for composition behavior and the decision interface for exhaustive
inventory partitions. This deepening adds no dry-run surface and no other
observable behavior.

## #4: FT218 — named Git readers for recurring repository facts

Blocked by: #1
Type: Grill

### Question

Which Git incantations become named readers after FT189, and does the deepening
replace the generic process plumbing?

### Answer

Promote a repository fact only when the same incantation occurs at three or more
production sites, and give each qualifying fact one light-path ticket. The
current first ticket is private worktree administration-directory resolution:
six production readers repeat the same absolute `--git-dir` query.
`CommonDir` already owns shared administration-directory resolution and is
consumed rather than recreated. Any further peeled-commit or porcelain-status
reader requires a fresh current-tree census before its ticket is cut.
`Output`, `Raw`, and `OK` remain generic plumbing; the deepening makes
callers ask for named facts without deleting the escape hatch.

## Not yet specified

## Spec-writer discretion

- Names and file placement inside the modules, provided the public command seams
  and decided responsibilities above do not move.
- Ticket slicing within each spec, following `craft-tickets` expand–migrate–
  contract order.
- Exact coverage-row wording and fixture placement, while the exit proof and
  behavior predicates above remain intact.

## Out of scope

- Correcting worktree eligibility precedence during FT216.
- Adding an adopt dry-run surface during FT217.
- Deleting `git.Output`, `git.Raw`, or `git.OK` during FT218.
- Reopening the already-landed gate transaction, verdict registry, shift
  objective, green marker, or skills-index reader deepenings.
- Building FT108's refactor skill as a prerequisite.

## Sources

- Path: `capture/architecture-review-20260817T104714.html`
  Supports: the current source verification, before/after relationships, deletion tests, and recommendation for #2 through #4.
  Drift: re-run the source survey if any named module moves before its spec or light-path ticket is authored.
- Path: `docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md`
  Supports: #2's conjunctive ownership, preservation, and fail-closed constraints.
  Drift: re-read if ADR 0005 is revised before FT216 begins.
- Path: `.bench/BENCH.md`
  Supports: #1's spec-versus-light-path routing and green, independently shippable cadence.
  Drift: re-read if the right-size table or workflow changes before implementation.
- Path: `roadmap/FT216.md`
  Supports: #2's current worktree evidence and tracked destination.
  Drift: re-read if FT207 or another worktree cleanup change lands before FT216 is authored.
- Path: `roadmap/FT217.md`
  Supports: #3's current adopt call-path evidence and tracked destination.
  Drift: re-read if an adopt lifecycle verb changes before FT217 is authored.
- Path: `roadmap/FT218.md`
  Supports: #4's current named-reader scope and light-path destination.
  Drift: re-run the recurrence census after any internal Git or worktree reader change.
