# Status Valve (FT30)

## #1: How does a reader see past the five-row budget?

Blocked by: —
Type: Grill

### Question
Six signals fire today; the five-row budget truncates to `+1 more` and
`bench status` has no flag to expand — the hidden row is exactly the
roadmap-reconcile signal built yesterday. A budget needs an overflow valve.

### Answer
`bench status --all` prints every row, budget off; the truncation line
becomes `+N more (bench status --all)` so the remedy is on the page that
truncates. The default five-row budget and severity ordering are unchanged —
the budget is right for ambient SessionStart injection; the valve is for the
reader who wants the rest. Rejected: raising the budget (defeats its
noise-control purpose); auto-expanding when overflowing (ambient surface must
stay bounded).

## #2: What form must action strings take?

Blocked by: —
Type: Grill

### Question
Status action strings recommend non-invocable surfaces — the specs row says
`promote-then-delete (spec-retire)`, the decisions row names a skill — against
the platform's own "recommend in the form the harness can invoke" rule.

### Answer
Every action string is something the reader can run or invoke verbatim: a
`bench` subcommand or a canonical `/bench-*` phase (the kit-prose convention;
harness adapters translate phase names at their own layer). Concretely: specs
row → `bench spec retire <slug>`, decisions row → `/bench-shape-idea`,
structure row → keep the split recommendation but name it via the phase that
loads the skill, orphaned worktree branches → `bench worktree clean` (once
[FT28] makes that true — land FT28 first or together), reviews row keeps
"promote or delete by hand" (genuinely manual). The spec enumerates the final
string per row so the contract tests pin exact output.

## Handoff

1. **Module boundaries.** `internal/status` owns rendering, the flag, and the
   strings; `bin/bench.sh` routes the new flag; no other package changes.
2. **Contracts.** `bench status` (no args): unchanged except new truncation
   line text and revised action strings. `bench status --all`: all rows, same
   format, no truncation line. Unknown args stay exit 2.
3. **Deep vs thin.** All thin — rendering and copy changes on an existing
   surface.
4. **Black-box assertables.** Row-count behavior at 6+ signals with and
   without `--all`; exact truncation line; exact action strings per row.
5. **Gate attachment.** The runtime status contract tests already pin row
   output — extend them; canary fixtures embedding status output text need the
   same refresh.
6. **Hostile-input owners.** status owns flag parsing (`--all` plus junk →
   exit 2).
7. **Uncertainty flags.** n/a — surface is settled.
8. **Rejected alternatives.** Bigger budget; auto-expansion; per-row
   verbosity flags.
9. **Domain watch-outs.** The SessionStart hook injects `bench status` output
   into every session — string changes ripple into any conformance fixture
   that pins that text; sweep fixtures in the same diff.

Dependency order: build with or after FT28 (worktree action string names its
sweep).
