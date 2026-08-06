
## 2026-08-06 — FT195 repair rounds: two ticket-authoring lessons

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
