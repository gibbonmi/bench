# open-learnings — shape the remaining learnings into buildable specs

Source: the two open entries in `.bench/learnings.md` after the 2026-07-04
learnings integration pass. Both entries are generalizable but still need product
decisions before `/bench-write-spec` can compile them into implementation work.

Current frontier: resolve #3. #1 split the learnings into two specs; #2 settled
the review-findings artifact lifecycle.

## #1: Are the two remaining learnings one spec or two?

Blocked by: none
Type: Grill

### Question
The remaining learnings point at two different product surfaces: the ambient
dashboard's stale-gate wording and the semantic-review phase's findings storage.
Do they belong in one spec because they both improve cold-session pickup, or in
separate specs because they touch different seams and contracts?

### Answer
Two separate specs. The review-findings learning changes the
`/bench-review-implementation` artifact workflow; the stale-gate learning changes
`bench status` behavior and severity wording. Keeping them separate gives each a
clean Handoff and tighter acceptance coverage.

## #2: Where do review findings persist, and when are they retired?

Blocked by: none
Type: Grill

### Question
`/bench-review-implementation` can surface findings the reviewer intends to fix,
but today those findings live only in chat. The build needs a durable pickup
surface and a lifecycle: where the artifact lives, when it is written, who edits
it, and when it is deleted or folded away.

### Answer
Persist actionable review findings in `reviews/<spec-slug>.md`. The review
phase writes the file only when findings need a later fix; clean reviews and
accepted residual risks do not create an artifact. The file records findings by
axis with citations and the worst issue so a later session can resume without
chat history. The fix session deletes the file in the green fix commit after the
findings are resolved.

## #3: How should `bench status` classify stale gate drift?

Blocked by: none
Type: Grill

### Question
The ambient dashboard currently reports any gate-cache tree mismatch as stale.
Some stale states are benign capture drift, while others mean committed code has
moved past the last green gate. The build needs a precise classification rule:
which paths count as benign, how mixed drift is worded, and what action should
lead.

### Answer
— (open: choose the benign/real drift rule and dashboard wording)
