# Learnings — usage journal

## 2026-08-08 — Debug receipts bound repair scope, not ticket count  [open]

The GP1 bootstrap debug receipt named one cause and the union of paths required
to repair it. Reading `bench-implement-spec`'s singular “takes a repair ticket”
literally produced one wide ticket that combined an exact-green bootstrap
publisher with its installed-only consumer. The required ticket-breakdown
review showed both halves could land green independently and resliced them into
an ordered producer→consumer chain.

The right behavior is to treat a debug receipt as the repair cause, maximum
ownership fence, and proceed condition. Run `craft-tickets` inside that bound:
it may derive one repair ticket or several reciprocal, independently-green
tickets, but their union may not escape the receipt and the blocked assignment
refreshes only after the terminal repair ticket lands.

Proposed rule change: make the debug-repair section plural-aware — “a receipt
takes an ownership-fenced repair sequence, never a small spec” — and state that
the sequence is resliced by `craft-tickets`, with one ticket remaining the
common case.

## 2026-08-08 — Ticket-shaping prose must dogfood its own decomposition  [open]

The repair-reslicing hotfix made its target prose locally clearer, but its own
ticketing demonstrated a worse pattern for later changes. It claimed that a
command-only or skill-only split would strand a project-gate red without
reproducing one, kept independently-green producer and consumer changes in one
wide ticket, and then accumulated duplicated fixture inventories and overlapping
prose anchors while trying to close that ticket's mutation rows. The guidance
change therefore risked teaching future ticket work to justify a wide slice in
prose and to add enforcement sediment instead of re-deriving the split.

The right behavior is to dogfood ticket-shaping guidance on the candidate that
changes it. Attempt every proposed thinner landing, require the claimed stranded
red to come from the real project gate, and run integration-surface and
single-source discovery again after coverage is added. A locally clearer command
or skill is not an improvement when applying its rule to the same diff produces
wider tickets, duplicate inventories, or overlapping enforcement.

Proposed rule change: extend `craft-synthesis`'s consistency and dogfood loops
for ticket-shaping prose. Apply `craft-tickets` — and the external `to-tickets`
tracer-bullet check when it is the source precedent — to the candidate itself;
compare the resulting slices and blocking edges with the pre-change behavior;
and reject or reslice a candidate that cannot reproduce its keep-together red,
makes unrelated tickets wider, or introduces a second source for fixture or
anchor knowledge.

## 2026-08-09 — Phase-close capture must not repay the authoritative gate  [open]

A spec promotion owns the exact composed gate and landing, but required phase-close
retro or learning capture happens after that terminal outcome. Writing the tracked
capture then makes the tree look like new ungated work and routes the session toward
an immediate second whole-project gate, even though the only new content records the
outcome of the gate that just completed.

The right behavior is to keep phase-close capture durable and reviewable without
automatically paying the whole-project oracle again. The capture must not weaken or
misstate project-green authority, and its eventual reviewed drain still needs an
ordinary authoritative landing path.

Proposed rule change: define a capture-only accounting path for post-promotion retros
and learnings. It should preserve the sole promotion gate, avoid an immediate
follow-up gate solely because capture was recorded, and make the deferred capture
visible to the next reviewed drain rather than silently treating it as shipped code.
