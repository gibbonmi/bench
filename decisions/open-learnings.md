# open-learnings — shape the remaining learnings into buildable specs

Source: the two open entries in `.bench/learnings.md` after the 2026-07-04
learnings integration pass. Both entries are generalizable but still need product
decisions before `/bench-write-spec` can compile them into implementation work.

Current frontier: closed. This map yields two separate specs: review-findings
persistence and benign stale-gate status wording.

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
Classify stale gate drift by a fixed capture-only allowlist: `ROADMAP.md` and
`.bench-notes.md`. If every changed path since the gated tree is in that
allowlist, `bench status` treats the stale verdict as benign capture drift. Any
code, docs, config, or mixed drift stays the current stronger stale-gate signal
and leads with `re-run the gate`.

## #4: What exact dashboard row should benign stale-gate drift render?

Blocked by: #3
Type: Grill

### Question
The stale-gate spec needs exact observable text for the benign case. What should
the signal detail and action say when all drift is limited to the capture-only
allowlist?

### Answer
Render the benign case as signal `gate`, detail `stale (capture-only drift)`,
and action `re-run when convenient`. The row remains a gate signal, but its
wording is deliberately softer than code/config drift.

## Handoff

1. **Module boundaries.** Review findings: `/bench-review-implementation` owns
   when the artifact is written and retired; `reviews/<spec-slug>.md` is a
   transient project-local pickup file, not a second spec. Stale gate:
   `internal/status` owns dashboard classification and wording; `internal/git`
   may own a tree-to-tree changed-path helper if the implementation needs one.
2. **Contracts.** Review findings: actionable findings that need a later fix
   create `reviews/<spec-slug>.md`; clean reviews and accepted residual risks do
   not. The file records findings by Standards/Spec/Coverage axis, includes
   citations and the worst issue, and is deleted in the green fix commit after
   resolution. Stale gate: all changed paths since the gated tree inside
   `ROADMAP.md` and `.bench-notes.md` render `gate | stale (capture-only drift)
   | re-run when convenient`; any other path, invalid cached tree, or mixed drift
   keeps the existing stale row and action `re-run the gate`.
3. **Deep vs thin.** Review findings is guidance-only phase behavior; the command
   file stays the deep source for the workflow, and the artifact is just the
   durable handoff surface. Stale-gate classification stays inside the status
   renderer path; callers still invoke only `bench status`.
4. **Black-box assertables.** Review findings: conformance anchors can pin that
   `/bench-review-implementation` names `reviews/<spec-slug>.md`, writes only
   actionable findings, and deletes the file after a green fix. Stale gate:
   status tests or runtime contracts cover capture-only ROADMAP drift,
   capture-only `.bench-notes.md` drift, real code/docs drift, and mixed drift
   with exact row text.
5. **Gate attachment.** Review findings is partly gate-blind because agents write
   the artifact during the review phase; the gate can pin the command prose but
   review/user verification owns whether a real review writes a useful file.
   Stale gate is gate-observable through `internal/status` tests and the status
   runtime contract.
6. **Hostile-input owners.** Stale-gate path comparison is root-relative and exact:
   `ROADMAP.md` and `.bench-notes.md` only, not nested or similarly named paths.
   Added, modified, or deleted allowlisted paths are benign when they are the
   entire diff. Missing or malformed gate cache data falls back to the existing
   stronger stale-gate row. Review artifacts use the spec slug from the reviewed
   spec path; no-spec reviews keep the current chat-only behavior unless the
   reviewer supplies a slug.
7. **Uncertainty flags.** None.
8. **Rejected alternatives.** One combined spec; appending findings to the spec;
   writing review artifacts for clean reviews; expanding benign stale-gate drift
   to all markdown or docs; changing the gate-cache key to ignore capture files.
9. **Domain watch-outs.** `ROADMAP.md` affects the status footer and
   `.bench/learnings.md` affects the learnings row, so neither "capture" nor
   "markdown" is enough as a class. Keep the allowlist deliberately small until
   real use proves another path is safe.

Dependency order: two independent specs; recommended order is review-findings
persistence first, then benign stale-gate status wording.
