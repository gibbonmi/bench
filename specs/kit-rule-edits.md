# Kit rule edits

Status: implemented

## Problem

Two 2026-07-06 learnings drained into decisions but never landed in the artifacts that enforce them. First: Agent-tool worktrees were repeatedly cut several commits behind `main` — one delegate stalled, one was permission-denied its own fast-forward, one built against a stale base — and the mid-run fix that held was never written down, so the next delegation rediscovers it. Second: an explicit reviewer batch drain produced ten map-less specs under a defensible override, but `/bench-write-spec`'s entry contract says "refuse without a closed map" with no override path, so the reconciliation happened ad hoc and will again. Both fixes are decided; neither is in the owning file.

## Solution

Write the two decided rules into their owning artifacts, add-only, weakening nothing:

- `craft-delegate`'s charge section gains the stale-base opener as a standing requirement for worktree-isolated delegates, plus the orchestrator's reciprocal duty, plus an explicit read-only exemption; the good-charge example teaches the opener line.
- `/bench-write-spec`'s entry contract gains one sentence recording the reviewer-directed batch-drain override, so the override path is on the page without softening the default map gate.

Both are guidance prose that compounds through every session that loads it. No CLI code, no new conformance anchor, no skills-index change.

## User stories

1. As a session spawning a worktree-isolated delegate, I want `craft-delegate`'s charge section to require the delegate's charge to open with "run `git merge --ff-only main`, verify HEAD equals main, stop and report if the merge is denied or diverges", to state the orchestrator's reciprocal duty (a blocked worktree is fast-forwarded by the orchestrating session, which then resumes the same delegate), and to say read-only delegates are unaffected, so a delegate never silently builds against a stale base. Line: claude-fable-5 / high. This is kit guidance prose that every worktree-isolated delegation loads, so it takes the `craft-line` leverage override of top model and high effort.

2. As a session copying a charge from the template, I want `craft-delegate`'s good-charge example to demonstrate the opener line on a worktree-isolated write-delegation, so the opener is taught by example and not only by rule. Line: claude-fable-5 / high. The good-charge example is template prose that future charges copy verbatim, so it earns the same leverage override as the rule it teaches.

3. As a session running `/bench-write-spec` under a reviewer's batch drain, I want the entry-contract paragraph to record that an explicit reviewer-directed batch drain (an assessment or reviewed findings doc into specs) may substitute for per-spec maps with every defaulted decision flagged in-spec for post-hoc veto, and that absent that explicit instruction the map gate stands, so the override is legitimate without weakening the default. Line: claude-fable-5 / high. The entry contract governs every spec session, so its override clause takes the leverage tier despite being a one-sentence add.

## Implementation decisions

- **Two files, add-only markdown, no existing rule weakened.** No Go code, no new conformance check, and no skills-index change (frontmatter `description`/`index` unchanged, so `bench` menus and `skills-index.sh --check` are unaffected). This matches the map's Handoff: rule text only.
- **File 1 — `.agents/skills/bench-craft-delegate/SKILL.md`, `## The charge`.** Add the stale-base opener as a standing requirement for worktree-isolated delegates, the orchestrator's reciprocal fast-forward-and-resume duty, and the explicit read-only exemption. The good-charge example demonstrates the opener on a **worktree-isolated write-delegation**; the existing read-only review example stays read-only, because read-only delegates are unaffected and the opener would be false on a `Do not edit any file` charge.
- **File 2 — `.agents/commands/bench-write-spec.md`, `## Entry orientation`.** Add the batch-drain override sentence into the existing "This phase refuses to run without a complete map." paragraph — the map gate keyword and the override live in the same place so a reader sees both at once.
- **Authored inline by the invoking top-tier session; no delegate.** Doc authoring is the `craft-line` leverage override, and both `craft-delegate` and `/bench-write-spec` permit a same-session delegate only on an explicit reviewer ask, which is absent here.
- **Wording stays as terse as the surrounding text.** These are leverage artifacts loaded by every delegation and every spec session (Handoff domain watch-out); an add-only edit must not inflate the file.

## Testing decisions

- **What a good test is here.** The artifacts are prose; the only observable is text presence in the live file, read the same way the docs conformance scan reads it. Assert the decided substance appears — do not re-derive it, and do not pretend prose presence is behavioral TDD.
- **Which seams get tested, and prior art.** The existing text-presence seams are the docs conformance scan (`internal/conformance`, `TestRootConformance` → `checkWorkflowAnchors`, which already `strings.Contains`-checks anchors like `refuses to run without` and `a claim, not a result`) and the canary meta phase (`tests/canary/workflow-guidance-anchors/`, deliberately-broken copies asserted red by substring). Author-time `rg -F` greps mirror those checks. Per the Handoff, **no new anchor is added**, so the specific new sentences are verified by grep plus review; the gate's role is to stay green, catching only accidental collateral damage.
- **Gate command.** `.bench/gate.sh`.

### Seam diagram

Delegate skill text:

    trigger: docs conformance scan (TestRootConformance) + author rg + every delegating session
        │
        ▼
    edit  ──▶  [ .agents/skills/bench-craft-delegate/SKILL.md ]  ──▶  charge-section rule text
    edit  ──▶  [   (## The charge, good-charge example)        ]  ──▶  taught opener example
                      ◀ tests attach here: `rg -F 'git merge --ff-only main' <file>`;
                        gate reads the live file, canary meta re-runs against fixtures

Write-spec command text:

    trigger: docs conformance scan (checkWorkflowAnchors) + author rg + every /bench-write-spec session
        │
        ▼
    edit  ──▶  [ .agents/commands/bench-write-spec.md      ]  ──▶  entry-contract paragraph
    edit  ──▶  [   (## Entry orientation, map-gate para)   ]      (map gate + batch-drain override)
                      ◀ tests attach here: `rg -F 'batch drain' <file>`;
                        existing write-spec-map-required canary pins the surrounding contract

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Charge section requires the stale-base opener for worktree-isolated delegates, states the orchestrator's fast-forward-and-resume duty, and exempts read-only delegates | `bench-craft-delegate/SKILL.md` text (`## The charge`) | `rg -F 'git merge --ff-only main' .agents/skills/bench-craft-delegate/SKILL.md` exits 1 today (verified absent); author-time grep, **not a gate anchor** — Handoff item 1 forbids new conformance code | The substring is absent while the opener sentence is missing; the read-only exemption is confirmed by reading, since no gate check pins prose semantics |
| 2 | Good-charge example demonstrates the opener line on a worktree-isolated write-delegation | `bench-craft-delegate/SKILL.md` text (good-charge example block) | grep for the opener inside the example block; not gate-enforced, verified by read | An absent or read-only example returns no match and reads as untaught; `craft-skills` review confirms the example is a write-delegation, not the read-only one |
| 3 | Entry contract carries the batch-drain override sentence in the map-gate paragraph, with the default map gate intact | `bench-write-spec.md` text (`## Entry orientation`) | `rg -F 'batch drain' .agents/commands/bench-write-spec.md` exits 1 today (verified absent); the existing `write-spec-map-required` canary still pins the surrounding `refuses to run without` contract | The substring is absent while the override sentence is missing; the canary independently proves the default map gate was not dropped |
| 1,2,3 | The edits keep the gate green: no pinned anchor dropped, every canary EXPECT still matches, and no embedded fixture is left stale | `.bench/gate.sh` (docs conformance scan + canary meta phase) | Run `.bench/gate.sh`; it exits 0 today and must stay 0. Adding prose cannot turn it red by itself — this row guards against an edit that accidentally deletes `refuses to run without a complete map` or `a claim, not a result`, or that changes a substring a canary fixture pins | The docs conformance scan reads the live files and the canary re-runs the gate against the broken fixtures; damaging a pinned anchor or a fixture EXPECT flips the exit code |

### Edge inventory

- **Pinned-anchor collision** → gate-green coverage row. The edits must not drop `refuses to run without a complete map` or `a claim, not a result`; both survive because the adds sit in the same paragraphs without removing those phrases.
- **Canary-fixture drift** → gate-green coverage row. Verified finding: no canary fixture embeds the changed text — the `write-spec-*` fixtures pin other anchors (map-required removal, `map's Handoff`, the ROADMAP-row sentence) and the delegate SKILL.md is embedded in no fixture — so **no fixture refresh is needed**. Conditional safety net: refresh a fixture only if the author's final wording happens to alter a substring that fixture pins (it should not).
- **Read-only exemption** → story 1 behavior. The opener must not bind read-only delegates; the example that teaches it is worktree-isolated, not the read-only review charge.
- **Verbosity of the added prose** — Won't handle: terseness is `craft-skills` and review judgment, not gate-observable.
- **Skills-index drift** — Won't handle: frontmatter descriptions are unchanged, so `skills-index.sh --check` is unaffected by design.
- **Duplicate insertion on a re-run of the author pass** — Won't handle: a single add-only edit with no runtime surface; review catches a double insert.
- **Gate-pinning the new sentences with a dedicated anchor** — Won't handle here: the Handoff scopes this pass to rule text only; see Out of scope.

## Out of scope

- **A permission-layer change letting worktree-isolated delegates fast-forward themselves** — a separate capability that widens the delegate write surface for a narrow need, and the map's rejected alternative; build later as ~3-4 edits (permission policy + adapter/hook wiring + a bite canary), ~3 gate runs.
- **A dedicated conformance anchor plus canary fixture pinning each new sentence so a future refactor can't silently drop it** — a separate gate-hardening capability, excluded here because the Handoff scopes this pass to rule text with no new code; ~2 edits (a `require()` line + a broken-copy fixture with its EXPECT) and 1 gate run per pinned sentence. Worth parking on the roadmap if leverage-prose erosion recurs.
