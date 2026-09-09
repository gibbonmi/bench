# Review pickup — FT311 delegate preparation

Frozen base `aafeb009534b9050d2902a2d69f68d306230321a`.
Reviewed tip `65b3ca8e8e92a58a148bd4a2d75a749f43ff5255`.

Three native same-family axes ran in isolated read-only worktrees. One
cross-family Codex pass ran the standing falsification, because the diff changes
kit guidance. One finding came from a peer session, and the coordinator verified
it.

Raw finding count: 20. Standards 9, Spec 4, Coverage 4, falsification 3.
De-duplicated repair targets after collapse: 16. One further observation is
parked, and it needs no repair here.

## Standards

Finding count: 9. Worst issue: the charge packet renders from two sources, and
the two copies have drifted.

### ST1. The charge packet renders twice

Disposition: `auto-fix`.

`AGENTS.md` states that two derivations of one fact must collapse into one
source. `renderCharge` at `internal/preflight/charge.go:48-115` and
`renderReviewPacket` at `internal/preflight/review.go:128-206` each author the
same tail. The drift is visible to a caller. `charge.go:97` names the evidence
table column `path`. `review.go:185` names the same column `source`.

### ST2. The anchor bite harness is pasted three times

Disposition: `auto-fix`.

`internal/conformance/docs_workflow_helpers_test.go:431-497` owns the pattern.
`internal/conformance/ft311_preparation_test.go:12-50` and `:52-113` repeat it.
The copy dropped the original uniqueness guard. One test therefore registers two
subtests with the same name.

### ST3. A live-tree test escapes its classification

Disposition: `auto-fix`.

`internal/conformance/tier_test.go:422-426` states that `classifiedLiveTreeTests`
is the sole classification of tests that read the live tree.
`internal/conformance/ft311_preparation_test.go:119` resolves the root by hand.
The detector cannot see that route. The new test is therefore absent from the
map.

### ST4. A move deleted a why-comment

Disposition: `auto-fix`.

The rationale that bound `readCandidate` to the hostile-input rule is gone from
the `internal/spec/` production code.

### ST5. Both guidance files sit at exactly their budget

Disposition: `ask-user`.

No rule forbids the long lines or the abutting headings. The base already
carried both patterns. The axis asks whether zero growth by compression meets
the budget intent.

### ST6. The build widened its own fence

Disposition: `ask-user`.

Commits `eb9f9785`, `cd97a787`, and `6f7c7b93` added ownership-fence rows. They
also rewrote the `Writes:` lines of three tickets. `6f7c7b93` rewrote the answer
text of two answered decision tickets. The spec is still staged. The tree
neither shows nor refutes an approval.

### ST7. Positional source slices couple two renderers to list order

Disposition: `auto-fix`.

`charge.go:70-71` and `review.go:145-146` map an index to a meaning. A reorder
reassigns the columns with no compile error.

### ST8. A new anchor registry file carries no doc comment

Disposition: `auto-fix`.

Every sibling file in `internal/anchors` carries one.

### ST9. Operator diagnostics carry a roadmap identifier

Disposition: `ask-user`.

Every new diagnostic starts with an FT number. The family it copies uses a
subject-descriptive prefix.

## Spec

Finding count: 4. Worst issue: the review phase never consumes the prepared
shared evidence that ticket 3 built.

### SP1. The review phase does not consume the prepared evidence

Disposition: `auto-fix`.

Ticket 5 requires the phase to consume the ticket 3 shared evidence. No
`--charge` invocation exists in `.agents/commands/bench-review-implementation.md`.
Lines 62 to 85 still collect the diff and the consumers separately. A peer
session raised this point. The coordinator confirmed it by an enumeration over
that file.

### SP2. The charge output contract deviates

Disposition: `auto-fix`.

The spec pins ten labels at `specs/ft311-delegate-preparation/spec.md:111`. The
build form emits twelve. The `fence` cell and the `evidence` cell hold the
identical string. The review form makes the `fence` cell and the `ticket` cell
identical. The coordinator observed both results in real runs.

### SP3. The triage input contract landed nowhere

Disposition: `auto-fix`.

Spec lines 213 and 214 name the inputs that triage must receive. The list is the
failed command, the selected checks, the execution evidence, the findings, and
the frozen source. A repository sweep finds that requirement in no guidance file.

### SP4. Two acceptance rows name tests that do not exist

Disposition: `ask-user`.

Rows DP7 and DP8 name `TestChargeProjection` and `TestChargeFullRetrieval`. The
build shipped one merged test. Both behaviors are covered, so this is a false
seam name. A seam-column amendment needs spec-change authority.

## Coverage

Finding count: 4. Worst issue: seven promised charge cells can be emptied, and
the owning suite stays green.

### CV1. The charge row cells are graded only by co-located tables

Disposition: `auto-fix`.

Four probes at the reviewed tip emptied required cells. Each one returned
`silent` against `./internal/preflight`, and each one restored. The coordinator
confirmed the baseline as pass, 18104 ms, zero skips. The assertions use a
substring match over the whole packet. The co-located tables therefore satisfy
them.

### CV2. Six review-owned rows have no stored record

Disposition: `ask-user`.

Rows DP19, DP21, DP22, DP23, DP24, and DP25 name a live exercise. No artifact in
the tree holds a dispatch record or a return. The finding discipline states that
evidence which lives only in a delegate return is not citable at the landing.
Row DP20 is the exception, because this three-axis run is its exercise.

### CV3. A cross-package contract has no test

Disposition: `auto-fix`.

`internal/preflight/review.go:83-105` hard-codes the meta shape that
`internal/consumers/command.go:99` declares. A rename there refuses every review
charge. The consumers suite stays green.

### CV4. One hostile-input class is open

Disposition: `auto-fix`.

The project profile lists a working directory deeper than the repository root.
No charge test and no proposal test runs from a nested directory. The spec edge
inventory disposes of no such class.

## Cross-harness falsification

The Codex pass ran at mid tier and high effort with stdin closed. It returned
three findings. The coordinator dispositioned each one.

### FA1. The handoff names no destination harness

Disposition: `no-op`.

The repo refutes the concern.
`.agents/skills/bench-craft-delegate/SKILL.md:19-24` already carries the same
pattern for a write delegation, and it predates this spec. The destination is a
runtime fact about the operator environment.

### FA2. The live review exercise is open

Accepted. This finding merges into CV2.

### FA3. The anchor matcher is a substring match

Disposition: `no-op` for this fence, and parked.

Codex enumerated no surviving contradictory guidance, and its own sweep found
none. The finding describes a limit of the anchor mechanism across the whole
kit. A semantic route checker is a new seam, and it needs its own spec.

The Codex pass also confirmed the removal half of row DP26 independently. It
found no surviving same-family CLI route and no surviving inline-axis route.
