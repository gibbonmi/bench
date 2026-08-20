# Land spec amendments

Status: implemented

Decision source: reviewer-confirmed current conversation, 2026-08-20 (FT225). Five forks closed by the reviewer: an in-range spec amendment in the landing source is legitimate regardless of author and the landing composes it; the source's spec bytes win unconditionally over a destination-side amendment, with no divergence refusal; `paths-authorized` implicitly authorizes the active spec's own folder in both preflight modes; the reviewer typing `--base` and `--source-tip` at `bench worktree land` is the acceptance of every commit up to that tip, with no re-delivery step; one spec ships both behaviors rather than two.

Verification log: 2 iteration(s) to accept — one `opus`/medium round returned two blocking findings. First, the spec promised the source's bytes win unconditionally while the composition refuses an overlapping destination spec edit as a merge conflict before the overwrite runs; folded as the spec-path neutralization decision and row LS13. Second, the implicit own-folder authorization makes the ownership-fence document source-amendable in-range, an authorization-chain change the Bootstrap authority section had not named; folded as the typed-range control paragraph and its Won't handle. Folded partials: the implicit entry re-derived from the existing spec-path fact instead of a second Facts field, LS8's degenerate restated as the mode-keying guard, LS4's refusal reattributed to the transition's parse error, the self-fence Won't handle corrected to name this spec, ticket 3's sweep row narrowed to the fenced file, and the crafted `--spec` guard named in the edge inventory. The second iteration verified every fold against the tree and accepted.

## Problem

A reviewed landing refuses any source whose spec differs from the destination's
copy, because publication rewrites exactly those bytes to
`Status: implemented`. A review that amends a story, a build that writes its
verification log, or a coordinator that adds a coverage row therefore breaks
its own landing. Three landings in three days each paid the same detour: a
hand commit on the destination, then a rebase of the source, because an
amendment commit inside the reviewed range also reds `paths-authorized`
unless the spec fences its own folder. The flow refuses without naming a
remedy, and the remedy it forces runs outside the flow.

## Solution

The landing publishes the source's spec bytes. `LandReviewed` keeps proving
that the bytes it publishes are the source tip's committed spec, and stops
demanding that the destination carry the same bytes. The amendment sits inside
the reviewed range, so the reviewer saw it in the diff and accepted it by
typing the source tip at the landing.

`paths-authorized` gains one implicit authorization: the active spec's own
folder, `specs/<slug>/`, in build mode, in review mode, and in the landing's
final source authorization — one predicate shared by all three consumers.
Self-fence boilerplate becomes unnecessary.

A destination-side spec amendment made after the review base is overwritten by
the source's transition, as the reviewer decided. Nothing is lost to history:
the landing commit keeps the destination version reachable through its first
parent.

## User stories

### The landing composes in-source spec amendments

Line: opus / mid. Landing logic at an existing seam with existing harnesses,
the cached routing for gate-critical code.

1. As a reviewer landing a reviewed source whose review amended the spec, I want the landing to publish the source's spec bytes, so that an amendment needs no hand commit on the destination.
2. As a build session, I want a spec write I committed in-range, such as a verification log, to land through the same composition, so that authorship never changes legitimacy.
3. As a reviewer, I want the published spec to be exactly the implemented transition of the source-tip bytes, so that what the review saw is what ships.
4. As a landing operator, I want the landing to refuse when the supplied spec bytes are not the source tip's committed spec, so that bytes no reviewed commit carries cannot publish.
5. As a reviewer, I want a destination-side spec amendment made after the review base to be overwritten by the source's transition, so that the reviewed truth wins without a divergence refusal.
6. As a landing operator, I want a source-tip spec that does not parse as staged to keep refusing, so that a broken spec cannot flip to implemented.
7. As a landing operator, I want a resumed landing to keep verifying the published commit against the source's transitioned bytes, so that recovery evidence stays sound while the first-run check changes shape.

### The spec's own folder is implicitly authorized

Line: opus / mid. Preflight decision-domain logic in the same cached routing.

8. As a build session, I want `paths-authorized` to authorize changes under my spec's own folder without a self-fence entry, so that spec and ticket writes need no boilerplate.
9. As a review session, I want that same implicit authorization in review mode, so that an in-range amendment never reds the preflight.
10. As a landing operator, I want the landing's final source authorization to accept amendment commits through that same one predicate, so that the preflight and the landing agree.
11. As a build session, I want a change under a different spec's folder to stay unauthorized, so that the implicit authorization covers only the active slug.
12. As a build session, I want a change under a sibling folder whose name merely extends my slug to stay unauthorized, so that the boundary is a path segment rather than a string prefix.
13. As a spec author, I want declared fence entries to keep their exact semantics, so that the implicit entry adds one authorization without changing any existing one.

### Guidance names the cadence

Line: fable / high. Guidance prose compounds through every session that loads
it — the leverage override's cached routing for doc authoring.

14. As a review session, I want the review-phase guidance to state that spec amendments commit to the landing source on the finding cadence, so that no session routes an amendment through a destination hand commit again.

### Reviewed exclusions

15. As a reviewer, I want no divergence detection against the destination, so that the landing never refuses on a destination-side spec change I already decided loses.
16. As a reviewer, I want no re-review trigger for post-review source movement, so that my typed `--base` and `--source-tip` remain the acceptance.
17. As a reviewer, I want no new flag or verb on the landing surface, so that the amendment path costs nothing to invoke.

## Implementation decisions

**One implicit-authorization source.** The paths-authorized check derives the
implicit entry from the spec path fact the gatherer already supplies — the
folder containing the resolved spec — and consults it beside the declared
fence entries, with the same segment-boundary rule. No second fact carries the
folder, so the printed spec path and the authorized folder cannot disagree.
Build preflight, review preflight, and the landing's `AuthorizeReviewedSource`
all route through that one gatherer-and-check pair, so the three consumers
cannot disagree either. The fence parser and the declared-entry semantics do
not change.

**The landing proves source provenance, not agreement.** The staged-spec check
keeps its source half — the bytes the landing will transition must be the
source tip's committed spec — and drops its destination half entirely. The
refusal message names the source-tip mismatch. Publication mechanics are
unchanged: the composed tree receives the implemented transition of those
bytes, and the prospective tree still gates before the commit exists.

**No destination comparison of any kind.** The landing reads the destination's
spec for no purpose beyond composition mechanics. A destination absence, a
stale copy, and a post-review amendment all land the same way.

**The composition neutralizes the spec path.** A destination-side spec change
that overlaps the source's amendment is a two-sided change, and a plain
merge-tree composition refuses it as a content conflict before the transition
ever runs. The decided behavior is unconditional, so the composition treats
the spec path as source-owned: the destination side is given the source's
bytes at that path before the merge, so the spec can never conflict, and the
transition then rewrites it. The exact mechanism is build discretion; the
behavior — an overlapping destination amendment still lands — is LS13's row.

**The resume path is already source-only.** It verifies that the published
commit carries the implemented transition of the source's staged bytes; that
comparison stays as it is.

**Acceptance is the landing invocation.** No re-review, re-delivery, or
freshness window guards post-review source movement. The existing tip checks —
the typed `--source-tip` must be the worktree head, the branch tip, and still
be both after the gate — remain the whole protection.

### Bootstrap authority

This change adds no executable hop. The landing's authentication chain — the
assignment ledger, the owner marker, the green-marker authorization, the
compare-and-swap on the destination ref — is untouched, and the composed
prospective tree still passes the full gate before any object becomes
reachable. The bytes being published gain a stronger provenance proof than
today's: they must be the committed content of the reviewed source tip, a
commit the reviewer named twice.

The authorization chain does change, and the change is named rather than
implied: the ownership-fence section lives in the spec the source may now
amend in-range, so the fence document becomes self-amendable — a source can
widen its own fences with an in-range commit. The record does not authenticate
itself; the control is the closed acceptance decision: the widening commit
sits inside the reviewed diff, and the reviewer accepts it by typing `--base`
and `--source-tip` at the landing. No mechanical pin replaces that reading.

## Testing decisions

- A good test exercises the landing's accept-or-refuse verdict and the exact
  published tree, or the preflight's per-row verdict — never the internal
  shape of either.
- Seams and prior art: `Decide` has a pure table-test harness over `Facts`
  (`internal/preflight`); `LandReviewed` has a repository-backed harness
  asserting refusals and published trees (`internal/landing`); the land
  command has an end-to-end harness with real assignments (`internal/worktree`).
- Gate seam: all three packages run in the ordinary `test` phase; the pure
  decision-domain tests stay repository-free per the profile's census.

### Seam diagram

    trigger: bench worktree land <path> --request <id> --base <b> --source-tip <t> --spec <slug> -m <msg>
        │
        ▼
    reviewed range  ──▶  [ AuthorizeReviewedSource ]  ──▶  authorized range fact
                              │ paths-authorized: declared fences
                              │ + implicit specs/<slug>/
                              ◀ tests attach here: Facts table rows
                              │ through Decide, repository-free
        ▼
    source tip spec  ──▶  [ LandReviewed ]  ──▶  two-parent landing commit,
                              │                  spec published as
                              │                  Implemented(source-tip bytes)
                              ◀ tests attach here: real repositories, assert
                                refusal strings and the published tree's bytes

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LS1 | 1, 3 | a landing whose source-tip spec differs from the destination's copy publishes the implemented transition of the source bytes | LandReviewed | a build that keeps the destination byte-equality demand refuses this fixture |
| LS2 | 4 | a landing whose supplied spec bytes are not the source tip's committed spec refuses | LandReviewed | a build that deletes the whole check publishes bytes no reviewed commit carries |
| LS3 | 5, 15 | a destination whose spec changed after the review base still lands, and the published spec is the source's transition | LandReviewed | a build that adds a divergence refusal or merges destination bytes goes red |
| LS4 | 6 | a source-tip spec that does not parse as staged refuses | LandReviewed | a rework that swallows the transition's parse error publishes a non-staged spec |
| LS5 | 1 | a destination that never carried the spec file still lands, and the published spec is the source's transition | LandReviewed | a build that reads the destination copy for any purpose fails on the absent file |
| LS6 | 7 | a resumed landing over a published commit carrying an amended source's transition verifies and completes | land resume command | a build that reworks the first-run check and breaks the resume comparison refuses a valid resume |
| LS7 | 8 | build preflight authorizes a changed path under the active spec's folder with no self-fence entry | Decide table | a build without the implicit entry reds paths-authorized on this row |
| LS8 | 9 | review preflight authorizes that same path the same way | Decide table | an implementation that keys the implicit entry on the mode field leaves one mode red while the other stays green |
| LS9 | 11 | a changed path under a different spec's folder stays unauthorized | Decide table | a build authorizing all of specs/ silences the fence for every foreign spec |
| LS10 | 12 | a changed path under a sibling folder whose name extends the slug stays unauthorized | Decide table | a string-prefix implementation authorizes the sibling folder |
| LS11 | 13 | a changed path authorized only by a declared fence entry remains authorized, and one outside every fence and the spec folder stays red | Decide table | a build that replaces declared entries with the implicit one flips one half of this row |
| LS12 | 2, 10 | `bench worktree land` over a source carrying an in-range spec-amendment commit, with no self-fence entry, lands green and publishes the amended transition | land command end-to-end | a build correct at each unit seam separately can still refuse composed |
| LS13 | 5 | a destination spec change that overlaps the source's amendment on the same lines still lands, and the published spec is the source's transition | LandReviewed | a plain merge-tree composition refuses this fixture as a content conflict |

Not covered: story 14 — guidance prose; the review round and the diff observe it, no gate seam grades sentence content.
Not covered: story 16 — an exclusion satisfied by writing nothing; review observes the absence in the diff.
Not covered: story 17 — same; the landing grammar is untouched and any new flag would be visible in the diff.

### Edge inventory

Walked against the profile's hostile-input checklist for shell CLIs.

- Spec absent at the source tip versus present-but-empty: both refuse today through spec resolution and the staged parse, and neither path changes; LS4 asserts the parse-failure class.
- Spec absent at the destination: LS5 — the landing reads the destination copy for no purpose, so absence composes.
- A command whose own write changes a fact it reports: the landing rewrites the spec it publishes; the existing contract already states the post-write truth, and LS1 asserts the published bytes exactly.
- Prefix boundary of the implicit entry: LS10; the segment-boundary rule is the declared-entry rule, shared, not re-derived.
- Slug spelled with a trailing slash or oddly cased: the implicit entry derives from the resolved spec path the gatherer already holds. A crafted `--spec` argument cannot widen it: the typed-status lookup matches only folder-enumerated slugs, so a path-shaped argument fails before authorization — that lookup is the guard, and it stays.
- A destination spec change overlapping the source's amendment: LS13 — the composition neutralizes the spec path, so the overlap never becomes a merge conflict.
- Hand-edited spec whose last line lacks a trailing newline: bytes pass through the existing transition untouched apart from the status line; no new parser is introduced.
- Re-run idempotency: LS6 drives the resume path over a published amended landing in a fresh process.
- State serialized by one process and reloaded by a fresh one: the resume comparison is exactly that boundary, and LS6 drives the second process rather than reusing the first's structures.

**Won't handle:** destination-side divergence detection — the reviewer decided the source's bytes win unconditionally; the surviving in-scope caller is any reviewer who amends the destination mid-build, whose version stays reachable through the landing commit's first parent.

**Won't handle:** re-review of post-review source movement — the typed `--base` and `--source-tip` are the acceptance; the surviving in-scope caller is every landing invocation, which still refuses when the typed tip is not the worktree head and branch tip.

**Won't handle:** implicit authorization anywhere outside the active spec's folder — declared fences keep governing every other path; the surviving in-scope caller is every build path outside `specs/<slug>/`.

**Won't handle:** migrating old self-fence entries — a self-fence stays valid and becomes redundant; the one staged spec today is this one, whose explicit self-fence is deliberate bootstrap, not a repair target.

**Won't handle:** an in-range commit that widens the spec's own ownership fences — the fence document is now source-amendable by design; the widening commit sits inside the reviewed diff, and the reviewer's typed `--base` and `--source-tip` are the control. The surviving in-scope caller is every landing, which still refuses a tip the reviewer did not type.

## Ownership fences

Ticket `authorize-the-spec-folder-implicitly`:

- `internal/preflight/`

Ticket `publish-the-source-spec-bytes`:

- `internal/landing/`
- `internal/worktree/`

Ticket `name-the-amendment-cadence-in-review-guidance`:

- `.agents/commands/bench-review-implementation.md`

All tickets additionally:

- `specs/land-spec-amendments/`

The spec-folder entry is explicit here on purpose: the implicit authorization
this spec ships does not exist until this build lands, so this build is the
last one that needs the self-fence.

## Out of scope

- **A divergence advisory on the landing report** — an informational line when
  the destination's spec moved after the review base, without a refusal. A
  separate capability over the same seam. 4 edits, 1 gate run.
- **Restructuring the landing-refusal surface** — FT169, FT224, and FT233 all
  edit that surface and are flagged as a joint restructure candidate. Roughly
  10 edits, 3 gate runs, owned by those rows.

## Further notes

The three FT225 occurrences (2026-08-18 bench-front-door, 2026-08-19 FT226,
2026-08-20 worktree-exec-run-binary) are the producer-derived partitions behind
stories 1, 2, and 12: a review amendment, a build-authored spec write, and a
coordinator amendment landing only through a rebase.
