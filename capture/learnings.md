
## 2026-08-06 — FT195 repair rounds: two ticket-authoring lessons [open]

What happened: the publication-transaction repair ticket charged prior-pair
restoration rows but no residue row, so the delegate's signal handler leaked
temp files the spec already forbade ("temporary staging is removed on handled
exit or signal"), costing a second repair round. Separately, the topology
repair enumerated the gate's attested publisher as an allowed production
caller after verifying it existed and was package-private — but not that any
production path reaches it; the fresh review found it dead, reversing the
two-publisher scope a round later.
Right behavior: (1) when a ticket's scope includes a signal or failure path,
walk the edge inventory for residue/cleanup classes and give them rows —
delegates build exactly to the charged rows; (2) admitting a caller into a
contract's allowlist requires proving production reachability, not source
presence.
Proposed rule change: none — both are covered by existing craft-spec
edge-inventory classes and the review axis; this is a charge-quality
reminder, not a rule gap.

## 2026-08-06 — Reviewer verdict: fence-breadth smell is too cheap [open]

What happened: the reviewer graded the transaction ticket's 4-directory fence
a sizing failure, not a paid-for exception — the one-line justification is
self-graded by the author at the moment they most want the slice whole, and
width multiplies the edge inventory on which rows get under-charged (the
observed chain: wide fence → unwalked residue/signal edge classes → two
repair rounds).
Right behavior: width buys obligations. Proposed craft-tickets change
("thin-by-requirement"): >2 directories requires the justification line PLUS
an edge-inventory account binding each added edge class (failure sites,
signal windows, residue/cleanup, concurrency) to an acceptance row or an
explicit review-rederivable "none applies"; an unaccounted class refuses at
review. >3 directories is a hard stop for reviewer sign-off before
assignment. Honest tracers keep the 2-3 range.
Proposed rule change: as above; land via craft-synthesis after FT195
promotes (reviewer-initiated, sign-off effectively given; exact wording
still theirs to approve).

## 2026-08-06 — Supersedes the thin-by-requirement proposal: simplify instead [open]

What happened: the reviewer rejected the width-buys-obligations design as
adding complexity to an already-large skill; decision is simplification with
thinness dominant. Retro test on the run's own wide ticket confirmed it:
"one observable outcome" justified a 4-directory fence, yet the work split
retroactively into two independently-green 2-file tickets — the observable
was one, the landings were two.
Right behavior: craft-tickets drops the fence-breadth and row-count numeric
smells and their justification-line escape hatches entirely. One rule: split
until splitting further would leave no independently-green landing; keeping a
group whole requires naming the specific red the thinner cut would strand — a
falsifiable claim review re-derives — never a description of the feature's
wholeness. Expand–migrate–contract stays. Tier economics stay in craft-line
(thin bounds the cost of a wrong cheap-tier bet; tier still follows seam
uncertainty, not fence size).
Proposed rule change: as above, reviewer-initiated; land via craft-synthesis
after FT195 promotes.

## 2026-08-06 — Promote red: unregistered tests in a classified family [open]

What happened: promote's prospective gate went red on
TestSubjectPackageTopology — the artifact fixture package's executable
registry of each surface package's top-level test names. The original FT195
build added fourteen posture tests without registering them, and two of this
session's repair tickets (stub relocation, seal-signal) added three more,
also unregistered; none of the tickets' Integration surfaces lines named the
classifier. It surfaced only at promote because the whole gate first runs on
a candidate there.
Right behavior: craft-tickets already states it — a test added to a
classified family carries its classifier as an integration surface; the
sibling-family search (start from an existing member, find every path that
classifies it) would have found the registry. Charge-time discipline, not a
rule gap.
Proposed rule change: none beyond the sizing change already queued.
