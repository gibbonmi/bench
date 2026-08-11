# Retro: bench-preflight

Built 2026-08-11 via `/bench-implement-spec --full`, single session, all 25
coverage rows landed across 7 implementation tickets plus 4 review-repair
tickets (12 gated commits, every landing green first try).

## What worked

- **Front-loading the repair loop into the breakdown.** The Codex
  (`gpt-5.6-sol`/high) ticket review returned 12 blocking findings plus 2
  verify-round residuals before any code existed; the build then landed with
  zero mid-ticket reds. Reslices at planning time were sentences; the same
  defects post-code would each have been a worktree + gate cycle.
- **Ticket precision carried sonnet.** Reviewer-directed sonnet delegates
  handled every implementation ticket (fable only for the kit-prose leverage
  ticket) because charges carried exact producers, signatures, exemplars, and
  per-token mutation tables. The story-level opus ceilings were never needed.
- **Independent coordinator probes earned their cost.** Two delegate rounds
  were rejected and re-charged on probe evidence (census entry gated nothing;
  worktree-change semantics had no red), and two coordinator probes were
  themselves vacuous on first attempt and redone — the probe-the-probe
  discipline caught defects on both sides.
- **Dogfooding mid-build.** The freshly landed command red-flagged its own
  build twice: a phantom `PF99` token in a repair ticket (rows-membership) and
  out-of-fence session-capture files at final landing (paths-authorized).
  Both were real drift.

## What to improve

- **Coordinator cwd discipline.** Two compound Bash commands ran their
  post-commit steps inside a worktree cwd, once cutting a nested worktree from
  the wrong root. Rule adopted mid-session: landing steps (ff, release,
  create) run as separate main-checkout commands.
- **Probe backup hygiene.** A stale round-1 scratchpad backup clobbered a
  delegate's round-2 diff (recovered from the delegate's own backup). Logged
  to `capture/learnings.md` with a proposed per-round-unique-name rule.
- **Semantic review remains the load-bearing layer.** All 10 composed-review
  findings sat behind correctly-cited coverage rows and a green gate: roughly
  seven were structurally invisible to the mutation machinery (undeclared
  predicates, standards violations), two slipped probe sampling, one was an
  over-claimed spec row. The fresh-context review delegate is not optional.
- **Mutation tables under-specify hostile classes.** P1 (cross-line
  parenthetical) had an honest green mutation against a single-line fixture —
  the declared red was real but the hostile class wasn't in any row. Edge
  inventories should name the multi-line variant explicitly when a grammar is
  line-oriented.

## Open at close (reviewer disposition pending)

- Post-repair review's six minors: provenance comment in `gather_test.go`,
  stale `RenderError` doc line in `internal/toon`, `fenceTokensInLine`
  parameter-list smell, two missing edge fixtures (odd backtick across lines —
  verified fail-closed; nested FIFO), RG2 wording (fixed in the landing).
- Held vetoes: S4 (ticket-ID comment density in tests), V2 (mixed-tag map
  policy, parked to `capture/IDEAS.md`).
- C16's shell-route mutation cannot bite: `bin/bench.sh`'s `*) route_binary`
  fallback dispatches unknown tokens through the registry (pre-existing).
