# roadmap-reconcile — surface roadmap rows whose work already shipped

Status: implemented

Source: `ASSESSMENT.md` backlog 4 (finding §3 high — the stale FT9 row, fixed
by hand in the 2026-07-06 spec-retire pass; this spec builds the mechanism so
the class can't recur silently).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

`ROADMAP.md`'s header promises "a row leaves when the work ships", but nothing
watches for the violation: the FT9 row still said "ready to spec" days after
FT9 was implemented and merged. The existing status signal fires on the merged
*spec file* awaiting retirement — but a row that outlives its spec (retire pass
missed the row) or a row whose named spec has shipped has no ambient signal,
so the drain cadence lags reality until someone happens to run
`/bench-what-next`.

## Solution

`bench status` cross-checks roadmap rows that name a spec path against the
tree: a row naming a `specs/<slug>.md` that is merged-implemented, or naming
one that no longer exists, fires a "roadmap row for shipped work" signal
pointing at the retire pass / `/bench-what-next`. Rows that name no spec stay
reviewer judgment, reconciled by `/bench-what-next` step 1 as today. The
roadmap convention — a row for spec'd work names its spec path — is recorded in
the ROADMAP header line so the detection surface is a stated convention, not an
inference.

## User stories

1. As a cold session, I want `bench status` to flag a roadmap row that names a
   `specs/<slug>.md` whose file is merged-implemented, so a shipped feature
   whose row survived is surfaced ambiently instead of waiting for a manual
   reconcile.
   Line: claude-opus-4-8 / medium. The dashboard's severity ladder and five-row
   budget are a gate-tested contract, so a new signal routes to the mid tier.

2. As a cold session, I want the same signal when a row names a spec path that
   no longer exists (the retire pass deleted the spec but missed the row), so a
   dangling row is caught by the tree itself.
   Line: claude-opus-4-8 / medium. Same contract surface as story 1.

3. As a reviewer, I want the ROADMAP header to state the convention that a row
   for spec'd work names its spec path, so the signal's detection surface is a
   documented contract and a row that omits the path is a visible choice to
   stay outside the ambient check.
   Line: claude-fable-5 / high. One header sentence, but it is workflow
   guidance every future drain inherits — the doc-authoring leverage override
   applies.

## Implementation decisions

- **Detection is path-literal, not semantic.** The scan looks for
  `specs/<slug>.md` tokens in roadmap lines and checks each against the tree:
  exists-and-merged-implemented → shipped-row signal; missing → dangling-row
  signal. No FT-number parsing, no title matching — a row that names no spec
  path is deliberately out of the signal's reach (what-next owns it). The
  merged-implemented test reuses the same one-source detector the existing
  retirement signal uses.
- **One signal, two details.** Both classes render under one `roadmap` signal
  name with a per-class detail ("row for merged work" / "row names a retired
  spec"), keeping the five-row budget cost to one row; the action points at the
  retire pass (`/bench-what-next` when the row needs judgment).
- **Severity ranks with housekeeping** — alongside the existing
  specs-awaiting-retirement row, below gate and git; the signal must never
  displace a red-gate or dirty-tree row.
- **A staged spec fires nothing.** Rows naming staged specs are the normal
  open-work state (including the nine rows this drain adds); only
  merged-implemented or missing spec files fire.
- **`bench status` only; no new subcommand.** The check is a read-only
  derivation inside the status renderer — the reconcile itself stays with
  `/bench-what-next` and the retire pass, which own the judgment.

## Testing decisions

- **What a good test is here:** drive the built binary's `bench status` in
  throwaway fixture repos with the three row states (names staged spec, names
  merged-implemented spec, names missing spec) and assert row text, action,
  and absence. Prior art: `internal/contract/runtime/runtime_status_test.go`
  (row rendering, budget, ladder).
- **Seam:** the `bench status` stdout renderer (runtime contract).
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: SessionStart hook or reviewer runs `bench status`
        │
        ▼
    ROADMAP.md rows       ──▶  [ status renderer                     ]  ──▶  ranked dashboard, incl.
    specs/*.md at HEAD    ──▶  [  scan rows for specs/<slug>.md      ]        "roadmap row for shipped
    (merged-implemented    ──▶ [  cross-check each against the tree  ]         work" + next action
     detector, one source)     [  rank on severity ladder            ]
                      ◀ tests attach here: runtime contract builds fixtures per row state
                        and asserts the signal fires for merged/missing, stays silent for
                        staged, and never displaces gate/git rows in the budget.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A row naming a merged-implemented spec fires the shipped-row signal with its action | status runtime contract | fixture ROADMAP row naming specs/x.md + merged implemented specs/x.md — no such row renders today | the signal is the feature; without the cross-check the dashboard stays silent and the assertion fails |
| 2 | A row naming a missing spec file fires the dangling-row detail | status runtime contract | fixture row naming specs/gone.md with no such file — silent today | catches the retire-pass-missed-the-row class that FT9's hand-fix worked around |
| 1 | A row naming a staged spec fires nothing | status runtime contract | already covered in spirit by today's silence, pinned as a guard: fixture with a staged spec asserts no roadmap signal | a naive "row names any spec" implementation would flag all nine staged rows this drain adds — this row makes that wrong implementation fail |
| 1 | The signal never displaces gate/git rows under the five-row budget | status runtime contract | fixture with a red gate + dirty tree + shipped row — asserting the ladder order fails if the new severity outranks them | pins the severity placement decision so the ambient contract survives the addition |
| 3 | The ROADMAP header states the name-your-spec-path convention | reviewer cold-read | not TDD-able — ROADMAP.md is per-repo content outside the kit's docs anchors; enforced at review and by the drain that writes rows | a header without the convention leaves the signal's coverage boundary undocumented; caught at review, stated openly |

### Edge inventory

- error path → **Won't handle**: an unreadable ROADMAP.md already renders
  status's existing no-roadmap posture; the scan adds no new failure mode
  (pure read).
- empty/absent input → covered: absent ROADMAP.md or no spec-path tokens →
  no signal (staged-spec guard row pins the no-false-fire side).
- boundary values → covered: multiple offending rows collapse into one signal
  row with a count — asserted in the story-1 fixture with two shipped rows.
- malformed input → row tokens that look like spec paths but carry stray
  formatting (backticks, bold) — the scan matches the path token inside
  markdown decoration; asserted in fixtures since the real roadmap wraps rows
  in bold/backticks.
- interrupted/partial state, re-run idempotency — **Won't handle**: read-only
  derivation.
- hostile environment (spec paths with spaces) — **Won't handle**: spec slugs
  are kebab-case by the template's own convention; a spaced path never enters
  through the sanctioned drain.

## Out of scope

- **Semantic staleness** (a row whose *wording* lies while naming nothing —
  "ready to spec" prose drift) — that is `/bench-what-next` step-1 judgment; a
  text-understanding check is a different capability and probably not a
  mechanical one. No estimate — it needs a shape decision first.
- **Auto-editing the roadmap from status** — status is a read-only reporter by
  contract; mutation stays with the drain phases. Estimate if ever wanted:
  ~5 edits, 3 gate runs.
