# Spec sourcing — every spec is compiled from a Handoff, never synthesized inline

Found via `/bench-debug` 2026-07-03: `/bench-shape-idea` advertised two grill
bypasses ("already clear / no real fog → skip straight to `/bench-write-spec`")
and `/bench-write-spec` accepted a no-map source ("sketch the seams from
scratch"), so the top model started authoring a spec directly from conversation.
Nothing in the gate bit. All tickets below were resolved in the discovery grill.

## #1: How does work reach a spec?

Type: Grill

### Question
What is the only permitted path from an idea (roadmap, conversation, anywhere)
to `specs/<feature>.md`?

### Answer
Grill → `decisions/<topic>.md` with a placeholder-free `## Handoff` → spec
compiled from that Handoff. Applies to **every** spec, however simple the idea —
a simple idea yields a short grill and a zero-open-ticket map in one sitting.
Both bypass sentences leave `/bench-shape-idea`; "don't manufacture tickets"
stays. `/bench-write-spec` refuses to run without a named map whose Handoff is
complete (`bench maps` shows no row for it). Bugs are exempt: they route to
`/bench-debug`, which needs no spec.

## #2: Who writes the spec, and where?

Type: Grill

### Question
Which model and venue compiles Handoff → spec, given the goals: the top model
does not author specs, and context windows stay small?

### Answer
Mid tier authors every spec. Default venue: a **fresh mid-tier session** — the
grill session ends at the closed map; the new session starts cold reading only
the map and the repo. Same-session mid **delegate** is the exception, taken only
on the reviewer's explicit ask (e.g. the spec must exist before the grill
session ends); the delegate is read-only and returns the spec text for the
invoking session to write, or runs worktree-isolated. Rationale: a fresh session
ends the expensive big-context session instead of keeping it alive as an idle
orchestrator, and sign-off is interactive anyway.

## #3: What reviews a spec draft before reviewer sign-off?

Type: Grill

### Question
Is there a model review pass between the mid-tier draft and the human sign-off?

### Answer
A **conditional** top-tier reviewer sub-agent, spawned by the spec-writing
session in a fresh small context (Handoff + draft only), firing only when the
Handoff carries uncertainty flags (item 7) or the draft deviates from the map.
It returns findings and a recommend/block verdict — **advisory**; sign-off stays
the reviewer's. No standing top-tier pass: the profile's
no-standing-top-tier-opt-out rule holds, and a complete Handoff makes the
mechanical case conformance work the gate already covers.

## #4: How does the invariant stay enforced?

Type: Grill

### Question
What stops the bypass prose from drifting back in?

### Answer
Gate anchors in the kit-conformance layer assert the Handoff-sourcing prose in
both command files, each with a red-by-construction canary fixture proving the
check bites. The `/bench-debug` repro loop (two `rg` assertions, observed red
2026-07-03) graduates into that canary.

## Handoff

1. **Module boundaries.** `/bench-shape-idea` prose (owns grill → map → Handoff;
   its entry and exit lose the bypass sentences); `/bench-write-spec` prose (owns
   the Handoff-required entry contract, the venue rule, and the conditional
   review rule); the gate's kit-conformance layer (owns the anchors); the canary
   suite (owns bite-proof); the project profile's Lines section (owns the
   spec-authoring routing note).
2. **Contracts.** `/bench-write-spec` input: a `decisions/<topic>.md` whose
   `## Handoff` is placeholder-free; on a missing or open map it stops and names
   the map to close, it does not draft. Spec-writing delegate: reads map + repo,
   writes nothing, returns spec text. Reviewer sub-agent: map + draft in,
   findings + advisory verdict out. Gate: non-zero with a per-anchor error
   substring.
3. **Deep vs thin.** The command prose is the deep unit (judgment lives there);
   the gate anchors are deliberately thin greps — they pin sentences, not
   semantics. No new shell modules.
4. **Black-box assertables.** Anchor greps' exit codes on the two command files;
   canary run goes red with the targeted substring; `bench maps` already asserts
   Handoff completeness (no new parser).
5. **Gate attachment.** Kit-conformance layer for the anchors; canary layer
   proves each new check bites. Venue/review conduct (which session, when the
   reviewer fires) is prose the gate cannot see — flag for review-phase
   attention, no TDD seam.
6. **Hostile-input owners.** Anchor wording must tolerate hard-wrapped prose
   (single-line fixed-string greps missed one bypass in the discovery repro) —
   owned by the gate anchor implementation. Others from the shell checklist:
   n/a — prose + grep change, no new CLI surface.
7. **Uncertainty flags.** None — all four tickets closed in the grill.
8. **Rejected alternatives.** Standing top-tier approval pass (moves merge
   authority to a model; violates no-standing-top-tier); same-session top-model
   spec authorship (the defect itself); roadmap-graduations-only scope (leaves
   conversation-born specs unguarded); a separate approval sub-agent named
   "approve" (sign-off is the reviewer's, structurally).
9. **Domain watch-outs.** Command/skill edits load at session start, so this
   change steers the *next* session, not the one that lands it. The workflow's
   right-size valve ("ask before deviating") is reviewer-gated and survives —
   it permits skipping the pipeline with explicit OK, never silently.

Dependency order: n/a — single spec.
