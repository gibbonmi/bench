# landing-refusal-diagnostics

Status: implemented

Decision source: named reviewed artifact — `roadmap/FT233.md` (FT233, graduated through reviewed drains; nine dated occurrences 2026-08-19..20; Group E's occurrence evidence sits in `roadmap/FT169.md`, 2026-08-18, while the diagnostic itself is FT233's fifth)

Verification log: 2 iteration(s) to accept — the round returned nine blocking findings; the author folded the design completions and mechanical fixes, the re-verify accepted them, and findings 1 and 5 stand as flagged reviewer decisions in the approval table

## Problem

Five independent diagnostics in the landing path each cost a session of guesswork.
An abbreviated `--source-tip` is refused with a drift-shaped message
("worktree source tip mismatch") that names no identity, while an abbreviated
`--base` silently resolves — the asymmetry that made the FT226/FT227 refusals
unreadable. A landing that published its commit and then stopped at the marker,
reconcile, or release step exits 1, so the one irreversible success in the flow
reads as a failure. The gate's own `.logs/` progress records count as ignored
residue in the source worktree, so running the gate blocks the release that
follows it. Several refusals name a predicate but none of the values it
compared, so diagnosis needs follow-up commands. And a request token stored
only as a SHA-256 digest cannot be recovered by a fresh session, and no refusal
names the sanctioned replacement path.

## Solution

Batched under one gate because each fix is small: landing-family refusals say
what they saw, what they wanted, and what to do next. An abbreviated commit
identity is named as an abbreviation with the full identity printed (and still
refused — the exact-identity posture does not move, and `--base` joins it as a
flagged reviewer decision). A published-but-incomplete landing exits with its
own code and prints the resume invocation. The gate's runtime log root joins
the one ignored-residue allowance everywhere in the lifecycle. The evidenced
refusals carry the offending paths or identities. A lost request token gets a
refusal that names the assignment and the `bench worktree reauthorize`
recovery.

Vocabulary, used consistently below: a **refusal** exits 1 and publishes
nothing; an **incomplete** landing has published its commit and stopped at a
follow-up step; a **released** landing exits 0; the **runtime root** is the
gate's ignored `.logs/` directory; an **abbreviated identity** is a hex prefix
(length 4–39) of the full commit identity a comparison wanted; a **resume
invocation** is a pasteable command with every value filled except the caller
token; a **continuation template** is a command with the values the refusing
command holds filled and placeholders for the rest.

## User stories

### Group A — an abbreviated identity is named, never misread as drift

Line: opus / medium. Literal comparisons at one known seam; mid model because
the surface guards landing authority.

1. As a landing operator, I want a short `--source-tip` refused with the
   abbreviation named and the full worktree identity printed, so that I fix the
   flag instead of chasing phantom drift.
2. As a landing operator, I want an abbreviated `--base` refused with the
   `--base` flag named as the abbreviated value, so that landing identities are
   uniformly exact. (Flagged: the tree accepts a short `--base` today, so this
   is a new refusal precondition — a reviewer decision, not a default.)
3. As a resuming operator, I want a short `--resume` value refused with the
   resolved full identity printed, so that the retry can pin the exact
   published commit.
4. As a landing operator, I want genuine drift — a full requested tip that
   differs from the worktree HEAD — reported with both identities and no
   abbreviation claim, so that drift and a typo read differently.
5. As a reviewer, I want no abbreviated identity ever accepted, so that
   landing authority still pins exact commits.

### Group B — a published landing never reads as a bare failure

Line: opus / high. The exit contract changes at an irreversible boundary; a
wrong code misleads every automation reader.

6. As a landing operator, I want a landing that published and then stopped
   (marker, reconcile, or release) to exit with a distinct incomplete code, so
   that the outcome cannot read as a failed landing.
7. As a landing operator, I want the incomplete output to print the resume
   invocation — published commit, identity flags, spec, and path filled, the
   request token as a placeholder — so that a fresh session resumes without
   archaeology.
8. As a resuming operator, I want a resume that stops at a follow-up step to
   exit with the same distinct code, so that one code means one thing on both
   paths.
9. As a script author, I want refusals to keep exit 1 and usage errors to keep
   exit 2, so that existing readers keep their meanings.
10. As a landing operator, I want a fully released landing to keep exit 0 and
    `worktree=released`, so that the green path is unchanged.
11. As a cold-session agent, I want the landing exit semantics stated in
    `.bench/BENCH-reference.md`, so that an incomplete outcome is interpretable
    without reading Go.

### Group C — the gate's own logs never block the lifecycle

Line: opus / high. The allowance feeds destructive eligibility; a wrong
widening loses user data.

12. As a releasing operator, I want a source worktree whose only ignored
    residue is the runtime root to release without flags or hand cleanup, so
    that running the gate cannot block the release that follows it.
13. As a releasing operator, I want ignored residue outside declared build
    outputs and the runtime root still retained fail-closed, so that the
    allowance stays narrow.
14. As an operator, I want the destination check, the release eligibility, and
    the clean eligibility to give the same verdict on the same residue, so that
    residue cannot pass one side and block another.
15. As an operator, I want the runtime-root allowance to hold when
    `.bench/build-outputs.json` is absent, so that a repo with no declaration
    still releases over gate logs.

### Group D — an evidenced refusal names what it saw and what to do

Line: opus / medium. The values are computed at each site or one parser switch
away with in-tree prior art; the work is carrying them to one grammar.

16. As a landing operator, I want "landing destination is not clean" to name
    the offending paths with a bounded listing and the true total, so that the
    first refusal is diagnosable without a follow-up command.
17. As an operator, I want ignored-residue refusals to name the residue paths
    with a bounded listing and the true total, so that I know what to clear or
    declare.
18. As a reauthorizing operator, I want a `--base` that fails the ancestry
    proof, when the addressed assignment's identity proofs otherwise hold, to
    name the recorded start as the wanted value, so that a moved destination
    HEAD does not send me source-reading.
19. As an operator, I want every enriched refusal to share one field grammar
    from one renderer — each command keeping its existing output channel — so
    that per-site drift cannot return.
20. As an operator with hostile path names (spaces, glob characters, control
    bytes, newlines), I want refusal output to carry no raw control byte and
    keep its record and listing structure, so that a hostile path cannot forge
    or split a record.

### Group E — a lost request token has a named recovery

Line: opus / medium. One refusal branch plus a hint; the authority design must
not move.

21. As a fresh session, I want a landing or release refusal — when the request
    digest matches no assignment but the target path holds exactly one active
    assignment — to name that assignment id and a `bench worktree reauthorize`
    continuation template, so that recovery never depends on the dead session's
    memory.
22. As a security reviewer, I want the stored request digest never accepted as
    the caller token, so that a ledger record cannot authenticate itself.
23. As an operator, I want no recovery hint when the target path holds no
    assignment, so that the hint never points at another session's worktree.

## Implementation decisions

- **One field grammar, one renderer, unchanged channels.** A typed refusal
  value in the worktree package carries `detail` plus optional `observed`,
  `wanted`, and `next`. One renderer produces the field text for every
  enriched site. `bench worktree land` keeps its stdout
  `refused{detail=…,observed=…,wanted=…,next=…}` record (absent fields
  omitted); `bench worktree release` and `bench worktree reauthorize` keep
  their existing `bench worktree <verb>: ` stderr prefix and exit 1, and
  append the same field grammar after their detail text. Plain errors render
  detail-only, unchanged. No field ever carries a raw control byte.
- **Bounded path listings.** Offending-path and residue listings render as a
  TOON path table after the record or message line — prior art: the cleanup
  plan's ignored-paths preview — under the classifier's existing entry limit,
  with the true total always stated. Path bytes never enter the `k=v` grammar.
- **The `next=` contract.** `next=` values follow the orphan-line precedent:
  shell-quoted, and a value that is not line-safe is replaced by its pointer
  form rather than escaped into a command that cannot work. A caller token is
  never echoed or persisted: the resume invocation fills the published commit,
  `--base`, `--source-tip`, `--spec`, and the path exactly and prints a
  placeholder for `--request`; the Group E continuation template fills the
  assignment id and path and prints placeholders for the new token and any
  identity the refusing command does not hold, and its placeholder text names
  the full-identity requirement (land refuses abbreviations even though
  reauthorize resolves them).
- **Exit contract.** 0 = released; 1 = refusal (nothing published); 2 = usage;
  3 = published-but-incomplete, on both the first run and the resume.
  `craft-cli` fixes 0/1/2 for AXI-approved query surfaces; `bench worktree
  land` is not one (only `worktree list` is AXI-approved), and the kit already
  carries a third non-verdict code — the gate's no-gate exit 3. The incomplete
  record gains `next=` carrying the resume invocation. `landedIncomplete` and
  `resumeIncomplete` collapse to one renderer — one source for the incomplete
  record shape.
- **One residue allowance.** `ignoredWithinBuildOutputs` — the variant whose
  *additional* predicate is always false — is deleted; every declared-outputs
  allowance composes `landing.RuntimeIgnoredPath`, which stays the one source
  of the runtime-root fact. Release and clean eligibility therefore treat
  runtime-root records as the destination check already does.
  **Flagged consequence (reviewer veto):** with the allowance in force,
  `bench worktree clean` plans `discard-remove` over runtime-root records
  without `--discard-ignored`; the plan/apply fingerprint gate still stands,
  and the records are the gate's own pruned-at-twenty run logs.
- **Abbreviation diagnosis.** At a landing-family identity comparison, a
  case-insensitive hex prefix of length 4–39 of the compared identity
  classifies as abbreviated: the refusal names the flag, states the
  abbreviation, and prints the full identity as `wanted`. The value is never
  accepted. A prefix shorter than 4 or non-hex input takes the drift shape
  with both identities and no abbreviation claim.
  **Flagged (reviewer veto):** for `--base` this adds a refusal the tree does
  not have today — a short `--base` currently resolves and lands — so story 2
  and LR2 are a posture decision, the one deliberate exception to the
  no-new-precondition rule below. The FT227 occurrence's tip-mismatch string
  cannot be reproduced from a short `--base` alone (that string has one
  producer, the worktree-HEAD-versus-`--source-tip` comparison), so the
  occurrence is treated as imprecise about mechanism and the row pins the
  chosen posture instead.
- **Bootstrap authority (recovery hint).** The hint discloses only the
  assignment id — already public as the first column of
  `bench worktree list` — and routes through `bench worktree reauthorize`,
  whose tree-derived exact-identity proofs are the trust root: worktree path,
  owner marker and lock reason, checked-out branch, worktree HEAD equal to the
  resolved tip, assignment branch tip equal to the resolved tip, and
  recorded-start ancestry. (The request CAS is a concurrency guard, not a
  proof.) Trace: operator → land/release (refusal, grants no authority) →
  reauthorize (re-proves identity from the tree, swaps in a new token) → land.
  The stored digest never authenticates: a record cannot authenticate itself.
- **Reauthorize wanted-value.** The base-ancestry refusal names the addressed
  assignment's recorded start; the lookup order inside the command is free,
  the behavior is pinned by LR18.
- **FT199 boundary.** No refusal precondition is added or removed, with one
  flagged exception: the `--base` posture decision above. Group C narrows one
  precondition's input: gate-owned runtime records stop counting as residue.
  Which refusals survive long-term stays FT199's coordinator's call.
- **FT169 boundary (flagged for reviewer veto).** FT169 reserves the
  interrupted-landing recovery *model* as reviewer-undecided, and its
  2026-08-18 occurrence is the evidence behind FT233's fifth diagnostic.
  Group E decides no recovery model: it only names the already-shipped
  reauthorize path inside a refusal and grants nothing.
- **FT224 boundary.** The property-wide refusal split (component naming on
  every site, the wrong `--discard-ignored` hint text, the light-path
  decision) stays on FT224. This spec enriches only the sites its occurrence
  ledger evidences, and the Group E branch is the one slice of the four-way
  string it touches.

## Testing decisions

- A good test drives a lifecycle command against a fixture repository and
  asserts the output record or message, the named fields, and the exit code —
  never a reading of internal state alone. Post-publication step failures are
  driven through the existing stub seams (`advanceLandingMarker`,
  `reconcileLanding`, `releaseLandingAssignment`).
- Seams with prior art: the landing command seam (`land_test.go` in-process
  calls plus built-binary journeys via `publicLandingFixture`), the
  eligibility-facts seam (`eligibility_test.go`), the allowance predicate unit
  seam (`build_outputs.go` tests), and the one system-journey refusal
  assertion (`internal/systemtest/owner_test.go`).
- The gate observes through its `test`, `race`, and `system` phases
  (`go test -count=1 ./...` and friends); no new gate phase.

### Seam diagram

    trigger: /bench-final-check (or the operator) runs
             `bench worktree land` / `release` / `reauthorize`
        │
        ▼
    flags + repo state ──▶ [ worktree lifecycle commands:      ] ──▶ output records
                           [ proofs → publication → follow-ups ]     land: stdout refused{…} | landed{…,next=…}
                                                                     release/reauthorize: stderr message + fields
                                                                     exit 0 / 1 / 2 / 3
                      ◀ tests attach here: in-process command calls and built-binary
                        journeys drive fixture repositories and stubbed follow-up
                        steps; assertions read records, fields, and exit codes
    (the eligibility facts and the allowance predicate are unit-tested
     below the command seam)

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LR1 | 1 | a land with `--source-tip` shortened to 12 hex and every other value correct exits 1 with a refusal naming the abbreviation and printing the full worktree HEAD | command | today the refusal is the bare drift string with no identity, so the row reds on the missing full identity |
| LR2 | 2 | a land with `--base` abbreviated and a full correct `--source-tip` exits 1 with a refusal that names `--base` as abbreviated and contains no tip-mismatch detail | command | the tree accepts a short `--base` today, so this row locks the flagged posture decision and reds on silent acceptance |
| LR3 | 3 | a resume with a short `--resume` value exits 1 with the resolved full identity in the refusal | command | the current message names no identity, so a retry has nothing to paste |
| LR4 | 4 | a land with a full-length `--source-tip` that differs from the worktree HEAD exits 1 with both identities in the refusal and no abbreviation claim | command | drift and typo currently produce one identical string, and the no-claim clause reds an implementation that calls every mismatch an abbreviation |
| LR5 | 5 | after a land with an abbreviated `--source-tip`, destination HEAD, source branch tip, and green marker are unchanged | command | reds an implementation that resolves and accepts the short identity |
| LR6 | 6 | a first-run landing whose release step fails after publication exits 3 with `worktree=incomplete:release` | command | today the exit is 1, indistinguishable from a refusal |
| LR7 | 6, 8 | a marker failure and a reconcile failure after publication each exit 3 | command | every post-publication step carries the distinct code, not only release |
| LR8 | 7 | the incomplete output contains the resume invocation with the published commit, both identity flags, the spec, and the path filled exactly and a placeholder for the request token | command | a fresh session must paste it; a missing value, a wrong value, or an echoed caller token reds |
| LR9 | 8 | a resume that stops at a follow-up step exits 3, not 1 | command | the resume path shares one meaning per code |
| LR10 | 9 | a pre-publication refusal still exits 1 and a grammar error still exits 2 | command | reds if the new code leaks into the refusal or usage paths |
| LR11 | 10 | a complete landing still exits 0 with `worktree=released` | command | freezes the green path against the contract change |
| LR12 | 12 | a release of a source whose only ignored residue is runtime-root gate records completes and removes the worktree | command | today the release retains with the residuals refusal, which is the FT233 hand-cleanup loop |
| LR13 | 13 | a release of a source with one ignored path outside the runtime root and outside declared outputs retains fail-closed | eligibility | reds if the allowance widens beyond the runtime root |
| LR14 | 14 | the residue set that passes the destination check also passes release eligibility and clean eligibility | eligibility | reds when one side allows what another blocks, the split that produced FT233's occurrences |
| LR15 | 15 | with `.bench/build-outputs.json` absent, a runtime-root-only source still releases | eligibility | the absent-declaration edge must not disable the runtime-root allowance |
| LR16 | 16 | the destination-not-clean refusal output names offending paths in the bounded table with the true total | command | today the refusal names zero paths |
| LR17 | 17 | the ignored-residue refusal output names residue paths in the bounded table with the true total | command | today the refusal names zero paths |
| LR18 | 18 | reauthorize with a non-ancestor `--base` and an otherwise-proven assignment names the recorded start in its stderr refusal fields | command | pins FT226, where the wanted value was absent |
| LR19 | 19, 20 | a refusal whose observed path carries a newline, ESC, and a comma emits no raw control byte and keeps the record line and the listing table structurally intact | command | the bytes that split a record are the ones asserted; raw pass-through or a corrupted table reds, and one renderer is what makes this assertable once |
| LR20 | 21 | a land with an unknown request against a path holding one active assignment refuses and names the assignment id and the reauthorize continuation template | command | today the four-way mismatch string names no recovery |
| LR21 | 21 | a release with an unknown request against the same shape names the same recovery in its stderr refusal fields | command | both lifecycle verbs carry the hint, not only land |
| LR22 | 22 | a land passing the stored digest itself as `--request` refuses and publishes nothing | command | reds any literal-digest acceptance, the cheapest wrong recovery |
| LR23 | 23 | a land with an unknown request against a path with no assignment refuses without a reauthorize hint | command | an unconditional hint would point at another session's worktree |

Not covered: story 11 — prose documentation in `.bench/BENCH-reference.md`; the
review round verifies it, no mechanical oracle asserts doc semantics.

Occurrence dispositions beyond the rows: FT233 occurrence 9 (an
untracked-but-kept path became ignored residue and stopped the release) stays
a blocking residue by design — story 13/LR13 keeps it fail-closed, LR17 makes
the refusal name it, and Group B's exit 3 stops it from reading as a failed
landing.

Cheapest wrong implementation per group, and the row that reds it: claim
abbreviation for every mismatched value → LR4; resolve and accept the short
identity → LR5; flip the exit code globally → LR10 and LR11; print `next=`
without the published commit or with the caller token echoed → LR8; allow all
ignored residue → LR13; allow only when a declaration file exists → LR15;
append raw paths into the record → LR19; enrich nothing and keep the bare
strings → LR16 and LR17; accept the ledger digest as the token → LR22; hint
unconditionally → LR23.

### Edge inventory

- An abbreviated value shorter than 4 hex, or non-hex, takes the drift shape
  with both identities (LR4's shape), never an abbreviation claim.
- An uppercase abbreviated identity is classified case-insensitively (folds
  under LR1's behavior).
- `.bench/build-outputs.json` absent versus present-but-empty: both leave the
  runtime-root allowance in force (LR15 drives absent; the build asserts the
  empty-declaration variant beside it).
- A genuine caller token that happens to be 64 hex digests as before; LR22
  only forbids accepting the *stored* digest.
- More than one active assignment on the target path: no recovery hint (the
  hint requires exactly one match).
- The Group E hint feeds the Group A refusal: reauthorize resolves a short
  tip, land refuses one, so the continuation template's placeholder text
  names the full-identity requirement (pinned in the `next=` decision).
- `bench worktree clean` over runtime-root records: plans `discard-remove`
  without `--discard-ignored` (the flagged Group C consequence; plan/apply
  still gates the removal).
- **Won't handle:** a `.logs` symlink — removal deletes the link object and
  never follows it; git worktree removal already guarantees non-traversal.
- **Won't handle:** movement between a refusal's comparison and its print —
  the refusal names the snapshot it compared; the movement-checked chokepoint
  is FT200's.
- **Won't handle:** the stale-executable freshness refusal — its remedy text
  is owned by the freshness verifier and passes through unchanged.
- **Won't handle:** the property-wide four-way refusal split and the
  `--discard-ignored` hint text — FT224 owns that property.
- **Won't handle:** adding or removing refusal preconditions beyond the
  flagged `--base` posture decision — FT199's coordinator settles which
  refusals survive.
- **Won't handle:** changing reauthorize's acceptance of short `--source-tip`
  values — the authority half of FT169; the hint's placeholder text is the
  bridge.
- **Won't handle:** a `request` column in `bench worktree list` — recovery
  routes through reauthorize; displaying the digest adds no recovery power.

## Ownership fences

- `specs/landing-refusal-diagnostics/`
- `internal/worktree/land.go`
- `internal/worktree/land_test.go`
- `internal/worktree/reauthorize.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/build_outputs.go`
- `internal/worktree/subshell.go`
- `internal/worktree/eligibility.go`
- `internal/worktree/eligibility_test.go`
- `internal/worktree/classifier.go`
- `internal/worktree/classifier_shape_test.go`
- `internal/worktree/ownership.go`
- `internal/worktree/ownership_test.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `internal/worktree/lifecycle_test.go`
- `internal/landing/landing.go`
- `internal/landing/landing_test.go`
- `internal/systemtest/owner_test.go`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `capture/session-handoff.md`

## Out of scope

- FT224's property-wide refusal naming (every site names its component and its
  fixing command, testable per component) — ~30 edits, ~6 gate runs.
- FT169's authority decisions and the interrupted-landing recovery model —
  shape-idea first; ~20 edits, ~5 gate runs after shaping.
- FT199's branch-retirement coordinator over a repository-wide ref inventory —
  ~40 edits, ~8 gate runs.
- A `request` column in `bench worktree list` — 4 edits, 1 gate run.
- Serving light-path (tickets-only) specs from `bench worktree land` — FT224's
  reviewer decision; ~15 edits, 3 gate runs once decided.

## Further notes

The FT227/FT226 occurrence ledger reports an abbreviated `--base` surfacing as
the tip-mismatch string; the review round traced that string to its one
producer (the worktree-HEAD-versus-`--source-tip` comparison) and confirmed a
short `--base` resolves and lands today. The occurrence is therefore treated
as imprecise about mechanism, and LR2 locks the flagged posture decision
instead of a reproduction.
