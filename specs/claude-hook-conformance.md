# claude-hook-conformance — gate the primary harness's hook wiring

Status: implemented

Source: `ASSESSMENT.md` backlog 2 (findings §2 high ×2).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

Deleting the Stop, PreToolUse:Bash, or SessionStart wiring from
`.claude/settings.json` leaves the gate green: conformance checks only the Agent
matcher on Claude, while the Codex equivalents (Stop → stop.sh,
Bash → block-dangerous-git.sh) are both gated. The primary harness has strictly
less conformance-gated hook wiring than the secondary — the largest silent
de-enforcement hole the assessment found.

## Solution

Extend the conformance suite so `.claude/settings.json` is held to at least the
Codex standard: Stop must run `stop.sh`, PreToolUse:Bash must run
`block-dangerous-git.sh`, SessionStart must run `session-start.sh`, alongside
the existing Agent → `check-agent-line.sh` check. Prove the new checks bite with
a canary fixture, mirroring the existing needle discipline.

## User stories

1. As a reviewer, I want a conformance diagnostic when `.claude/settings.json`
   is present but its Stop event does not run `.bench/hooks/stop.sh`, so gutting
   the completion oracle's interactive layer on the primary harness turns the
   gate red instead of passing silently.
   Line: claude-sonnet-5 / medium. The check is a structural clone of the gated
   Codex pattern at a known conformance seam, and the profile's mid-effort floor
   for oracle logic sets the effort.

2. As a reviewer, I want the same diagnostic for a PreToolUse:Bash group that
   does not run `.bench/hooks/block-dangerous-git.sh`, so the destructive-git
   guard's wiring is enforcement, not advertisement.
   Line: claude-sonnet-5 / medium. Same clone, same seam, same floor as story 1.

3. As a reviewer, I want the same diagnostic for a SessionStart event that does
   not run `.bench/hooks/session-start.sh`, so the ambient dashboard's injection
   point can't be dropped unnoticed.
   Line: claude-sonnet-5 / medium. Same clone, same seam, same floor as story 1.

4. As a reviewer, I want a canary fixture proving at least one of the new Claude
   wiring checks bites (a fixture settings.json with the Stop wiring removed goes
   red with the targeted substring), so the new checks can't rot into
   always-pass.
   Line: claude-opus-4-8 / medium. Canary fixtures are the gate guarding the
   gate, and wiring a new needle into the fixture families is oracle work that
   warrants the mid tier.

## Implementation decisions

- **Mirror `checkCodexHooks`, don't generalize prematurely.** The new checks live
  beside the existing Claude Agent-matcher check and read the same parsed
  structure; one parse of settings.json feeds all four wiring assertions. If the
  build finds the Claude and Codex checkers collapsing naturally into one
  helper, that is welcome (one source per fact), but the contract is the
  diagnostics, not the helper shape.
- **Absent file skips; present file is checked.** `checkCodexHooks` returns no
  diagnostics when `.codex/hooks.json` is missing, and the Claude checks keep
  that posture for `.claude/settings.json` — the kit repo always has the file,
  and canary fixtures hide dot-dirs, so fail-closed-on-absent would misfire
  across every fixture family. (Default call, flagged: this means *deleting the
  whole file* still evades the check; see Out of scope.)
- **Malformed JSON stays owned by the JSON-validity check** — the wiring checks
  return nothing on unmarshal failure exactly as both existing checkers do,
  because a second answerer for "is this valid JSON" would duplicate a fact the
  validity family owns.
- **Match by command substring** (`.bench/hooks/<name>.sh`), the established
  convention — Claude wiring uses `$CLAUDE_PROJECT_DIR` prefixes and Codex uses
  `git rev-parse` prefixes, so the hook-script path is the stable token.
- **The needle lands in the canary family that owns the wiring checks**, as a
  broken-fixture directory with the targeted-substring EXPECT, following the
  existing per-family needle pattern in `tests/canary/`.

## Testing decisions

- **What a good test is here:** run the conformance check against a fixture tree
  and assert the diagnostic strings — the same black-box style the existing
  wiring checks use — plus the canary run proving red-for-the-planted-reason.
  Prior art: `internal/conformance/validity_checks_test.go` (checkCodexHooks),
  `internal/conformance/line_routing_checks_test.go` (Claude Agent matcher),
  `tests/canary/` fixture families.
- **Seams:** the conformance diagnostic surface (root conformance run against a
  tree), and the canary meta-gate.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: `bench gate` (conformance phase; canary phase re-runs it on fixtures)
        │
        ▼
    .claude/settings.json  ──▶  [ Claude wiring checks              ]  ──▶  diagnostics:
    (tree under grade,      ──▶ [  Stop → stop.sh                   ]        "claude settings.json <event>
     BENCH_CONFORMANCE_ROOT)    [  PreToolUse Bash → block-danger…  ]         missing or does not run …"
                                [  SessionStart → session-start.sh  ]
                                [  Agent → check-agent-line.sh (existing) ]
                      ◀ tests attach here: point the check at fixture trees with each
                        wiring present/removed and assert the diagnostic set; the canary
                        fixture asserts the full gate goes red on the planted removal.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Removing the Stop → stop.sh wiring from settings.json yields a conformance diagnostic | conformance diagnostics | a fixture settings.json without the Stop hook passes conformance today — asserting the diagnostic fails | the diagnostic is the enforcement; if the check is absent or matches too loosely, the assertion fails |
| 2 | Removing the Bash → block-dangerous-git.sh wiring yields a diagnostic | conformance diagnostics | same probe for the Bash matcher — passes conformance today | same as story 1, for the git guard |
| 3 | Removing the SessionStart → session-start.sh wiring yields a diagnostic | conformance diagnostics | same probe for SessionStart — passes conformance today | same as story 1, for the dashboard injection |
| 1 | An intact kit-repo settings.json yields no new diagnostics | conformance diagnostics | already covered — the kit tree runs root conformance in every gate; a false-positive check turns the gate red immediately | guards against a check that fires on the healthy wiring shape ($CLAUDE_PROJECT_DIR prefix) |
| 4 | The gate run against a fixture with Stop wiring removed goes red with the targeted substring | canary meta-gate | the canary fixture does not exist today; adding EXPECT before the check exists would fail the canary run | proves the new checks bite; a check rotted into always-pass fails the canary, not just review |
| 1–3 | Absent `.claude/settings.json` yields no wiring diagnostics (posture parity with Codex) | conformance diagnostics | already covered by posture, un-pinned — add the probe with the file removed; it passes on day one and is a regression guard, stated openly | pins the skip-on-absent decision so a later fail-closed change is a deliberate edit, not drift |

### Edge inventory

- error path → covered: each removed wiring is a coverage row.
- empty/absent input → rows: absent settings.json (skip posture pinned);
  empty hooks object behaves as removed wiring (story 1–3 probes).
- malformed input → **Won't handle** here: invalid JSON is owned by the existing
  JSON-validity check; the wiring checks return nothing on unmarshal failure by
  design (one owner per fact).
- boundary values → covered: multiple hook groups per event (the healthy kit
  file has two PreToolUse groups) are exercised by the intact-tree row.
- interrupted/partial state, re-run idempotency — **Won't handle**: the checks
  are pure reads over a tree; nothing to resume or make idempotent.
- hostile environment (fixture dot-dir hiding) → covered by story 4: the canary
  family's `dot-claude` prefix convention is part of the fixture.

## Out of scope

- **Fail-closed on a deleted `.claude/settings.json`** — a separate enforcement
  decision that would have to distinguish "kit repo" from "fixture/consumer
  tree" (canary fixtures and consumer repos legitimately lack the file);
  needs its own posture design. Estimate: ~3 edits, 2 gate runs.
- **A Codex Agent-matcher line guard** (`check-agent-line` parity on the
  secondary harness, assessment §2 med) — a new enforcement capability gated on
  whether Codex hooks support an Agent matcher at all; parked as a ROADMAP row
  pending that research. Estimate: ~4 edits, 2 gate runs once feasibility is
  known.
