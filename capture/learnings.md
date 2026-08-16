# Learnings — usage journal

## 2026-08-16 — FT210 spec: two loops, both blocked twice under a two-round cap [open]

**What happened.** The FT210 spec's falsification loop (Sol/high, `codex exec`) blocked
on 10 findings; five were my own false current-code claims (per-path clean already
deletes proven-landed branches, per-path fingerprints bind the whole ledger so cannot be
frozen across a set, `ClassifyPathShape` is not on the planner path, "already covered"
tests that assert something else) and one reopened a map answer built on one of those
premises. The ticket loop blocked on granularity: I kept groups whole without naming
the stranded red, and consumer tickets said "the shared classifier" instead of
restating the predicate. Round 2 in each loop was fix-verification only and still left
partials, folded after the cap.

**Right behavior.** Before writing a "today X" or "already covered" cell, run or read the
exact function/test named — the code-versus-claim check is per cell, not per section.
When a grill recommendation states a premise about current code, verify it before
asking; a wrong premise closes a decision the reviewer then has to reopen. Slice tickets
to one coverage row by default and only merge with a named stranded red.

**Proposed rule change.** `craft-spec`: an "already covered" red-signal cell names the
test function and the assertion line, not the file. `craft-grill`: a recommendation's
one-clause why may not assert current-code behavior the author has not read this
session.

## 2026-08-16 — craft-tickets and to-tickets slice the same spec 17 vs 5 [open]

**What happened.** On the FT210 spec (20 coverage rows), a review applying
`craft-tickets` literally ("split until no independently-green landing remains") produced
17 tickets, twelve of them single-row test-only leaves. Pocock's `to-tickets` (complete
vertical slice, demoable alone, one fresh context window, prefactor first) produced 5
with the same dependency spine. The reviewer approved the 5.

**Right behavior.** A ticket is a demoable vertical slice sized to a fresh context window;
a coverage row that only adds a test to a seam its parent slice already opened is that
slice's acceptance row, not a ticket.

**Proposed rule change.** Remake `craft-tickets` on `to-tickets`' rules (and its quiz —
granularity, edges, merge/split — as the approval round), the same way `craft-spec` is
being remade on `to-spec`; retire the "smallest independently-green" splitting sentence
and its anchors/canaries in the same change.
