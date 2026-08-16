# Worktree landed-assignment visibility and bulk retirement (FT210)

Status: ready

## Destination

A landed-but-unreleased assignment stops reading as healthy, and retiring a build's
worktrees is one sanctioned plan/apply command instead of a hand-written loop.
Today `state=active` means only "never released": the writer signal (`lease`) and
landedness (`landed`) are separate `bench worktree list` columns, the session-start
summary folds all of it into `retained active=N`, `release` needs the creating request,
`land` releases only the source it lands, and `clean` retires one path per call. Scope
is one spec holding three pieces: (A) a derived `landed` classification projected by
`list` and the resume summary, (B) `bench worktree clean --landed` plan/apply,
(C) workflow prose that makes release-at-acceptance and a final-check sweep mechanical.

Every ticket below was closed by the reviewer in the 2026-08-16 shaping session
(four grill rounds); each Answer is the exact predicate fixed.

## #1: How does a landed-but-unreleased assignment become visible?

Blocked by: none
Type: Grill

### Question

New ledger state, derived classification, or summary-line-only?

### Answer

Derived classification, no new ledger state. Canonical predicate: `state=active ∧
landed=true ∧ lease≠live`. Dirty or ignored residue affects the plan action, not the
classification. `list` and the resume summary project it; `retained active=N` splits
into `active` (everything else still active) and `landed`.

## #2: Who guarantees release for per-ticket worktrees?

Blocked by: none
Type: Grill

### Question

Coordinator at ticket acceptance, closing-phase sweep, or both?

### Answer

Both: the coordinator runs `bench worktree release` for each accepted independent
slice, and `/bench-final-check` runs the bulk sweep as the backstop.

## #3: What is the sanctioned bulk form?

Blocked by: none
Type: Grill

### Question

A selector on `clean`, a fold into FT199's `bench branches retire`, or a bulk release?

### Answer

`bench worktree clean --landed`: a selector on the existing plan/apply verb reusing
clean's ownership-free retirement (FT148 #1) and refusal contract. FT199 stays at ref
level. Bulk release is rejected: release's request digest is the ownership proof.

## #4: Scope edges — automatic removal and recovered rows

Blocked by: none
Type: Grill

### Question

Does the session-start automatic path ever remove a landed tree, and are `recovered`
rows in the bulk sweep?

### Answer

No and no. The automatic path keeps retaining active assignments (FT148 design);
retirement is visibility plus the explicit bulk verb. `recovered` rows hold unpublished
work and remain per-row reviewer decisions (FT98/FT199).

## #5: One spec or a split?

Blocked by: #1, #2, #3
Type: Grill

### Question

Pieces A, B, and C ship as one spec or separately?

### Answer

One spec: B's plan rows are A's classifier projected, and C is the prose that names B.

## #6: Dirty or unknown-lease rows in the bulk plan

Blocked by: #3
Type: Grill

### Question

Retain, or preserve under recovery refs like per-path clean?

### Answer

Retain and list. Bulk apply removes only rows whose planned action is remove; a row
with uncommitted work or an unknown lease is retained with its reason, and the plan
advertises `bench worktree clean <path>` for it. Bulk never authors recovery refs.

## #7: Closing-sweep scope

Blocked by: #2, #3
Type: Grill

### Question

Repo-wide `clean --landed`, or only the trees the build created?

### Answer

Repo-wide. The predicate is safe independent of who cut the tree; no build-scoped
bookkeeping is kept.

## #8: Ledger record after bulk apply

Blocked by: #3
Type: Grill

### Question

Settled in the apply, or left for the next reconcile?

### Answer

Settled in the same apply: the record is dropped or completed so `list` shows no
missing-tree active rows and the resume summary's open count is right immediately.

## #9: Flag composition

Blocked by: #3
Type: Grill

### Question

How does `--landed` compose with clean's grammar?

### Answer

`--landed` with a `<path>` operand is an invocation error (exit 2, the existing clean
usage line). `--discard-ignored`, `--full`, and `--discard-branch` apply to every
planned row. Branch disposition matches per-path clean exactly: a branch the tool's
own landedness proof licenses is deleted with its tree, and `--discard-branch` is the
operator supplying that proof for a branch the derivation cannot show. (Reopened and
re-closed by the reviewer 2026-08-16 during spec authoring: the first answer's
"landed branches stay behind for FT199" rested on the author's false premise that
per-path clean never deletes without the flag.)

## #10: Plan/apply contract

Blocked by: #3
Type: Grill

### Question

One fingerprint over the set, per-row fingerprints, or best-effort apply?

### Answer

One fingerprint binding the exact row set and each row's planned action and OID.
`--apply <fp>` refuses the whole apply if any row drifted (new landed row, tree
changed, lease went live).

## #11: Empty selection

Blocked by: #3
Type: Grill

### Question

`clean --landed` with zero landed assignments?

### Answer

Exit 0, empty plan, no fingerprint. `--apply` against an empty set is an invocation
error because no fingerprint was issued.

## #12: Enforcement mechanism

Blocked by: #2
Type: Grill

### Question

What makes release-at-acceptance and the closing sweep mechanical?

### Answer

A phase step plus anchored prose: `/bench-final-check` runs the bare `clean --landed`
plan and applies it as a named step before its landing report, which carries the sweep
result; `craft-delegate` names release-at-acceptance; both sentences ride the anchor
registry. No new gate check on ledger state.

## #13: Canonical term

Blocked by: #1
Type: Grill

### Question

Name of the classification in `retained <reason>=N` and list actions?

### Answer

`landed` — one word for the fact the `landed` column already carries per row. Avoid:
orphan, stale, idle, abandoned, unreleased.

## #14: Advertised remedy and orphan interplay

Blocked by: #1, #3
Type: Grill

### Question

What do `list` and the resume summary advertise when landed > 0, and how does it
interact with the age-based orphan listing?

### Answer

Both emit the bare `bench worktree clean --landed` plan invocation as the AXI action.
A row that is both aged and landed is counted and listed as `landed`; the age-based
orphan candidate list keeps only non-landed actives.

## #15: Unknown landedness

Blocked by: #1
Type: Grill

### Question

Row whose `landed` reads `unknown` (or no resolvable default branch)?

### Answer

Plain `active`, never `landed`: skipped by the bulk plan, counted under `active`.

## Not yet specified

## Spec-writer discretion

- Whether the `landed` reason is a `CleanupReason` value in the existing retained-reason
  ordering or a sibling projection, provided the summary token is `landed=N`.
- Where the shared classifier lives, provided `list`, the resume summary, and the
  `--landed` planner derive the predicate from one source.
- Row layout of the bulk plan output within the existing `renderCleanup` shape.

## Out of scope

- A request-derivation override or any bulk `release` (FT148 #1 stands).
- Automatic removal of landed trees at session start.
- `recovered` rows and their recovery refs (FT98, FT199).
- Branch inventory and classification across refs (FT199).
- A gate check on ledger state.

## Sources
