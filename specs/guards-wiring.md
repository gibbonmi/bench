# guards-wiring — advertise enforcement by wiring, not file presence

Status: staged

Source: `ASSESSMENT.md` backlog 5 (findings §2 med ×2).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

`bench guards` discovers hooks by file presence in `.bench/hooks/`, so it
advertises a deny surface that may not fire: on a Codex-only session it lists
`check-agent-line` even though `.codex/hooks.json` wires no Agent matcher. And
the managed pre-push hook's `--describe` unconditionally advertises
"`.bench` drift from bench gate pin" while the hook itself skips drift
enforcement (warn-only) when no pin file exists. Both are the kit's own defect
class — advertisement stronger than enforcement — inside the enforcement layer.

## Solution

`bench guards` gains a `wired` field per row, derived by scanning the harness
hook configs (`.claude/settings.json`, `.codex/hooks.json`) for each hook
script, with pre-push reporting its installed state as today. The pre-push
`--describe` manifest becomes state-aware: its denies line includes the drift
clause only when a pin actually exists, so the manifest is generated from the
rules currently in force.

## User stories

1. As an agent orienting in a repo, I want each `bench guards` hook row to carry
   a `wired` field naming the harness configs that actually reference the script
   (`claude`, `codex`, both, or `none`), so the deny surface I read matches the
   hooks that can actually fire here.
   Line: claude-opus-4-8 / medium. The guards table is the enforcement layer's
   self-report and the wiring scan is a new derivation across two config
   formats, so it warrants the mid tier.

2. As an agent, I want the pre-push row to keep its real install check
   (managed / unmanaged / not installed) unchanged, so the one guard that
   already measures wiring stays the pattern, not a regression.
   Line: claude-sonnet-5 / low. This is a regression guard on behavior the
   contract suite already pins, so the cheap tier at low effort is right.

3. As a reviewer, I want the pre-push `--describe` denies line to include
   "`.bench` drift from bench gate pin" only when a pin file exists — and to say
   the drift check is disarmed when it doesn't — so `bench guards` output
   collapses advertisement to enforcement instead of overstating it.
   Line: claude-opus-4-8 / medium. The manifest is generated inside enforcement
   shell that runs at push time, so a defect here mis-describes or breaks the
   push guard; the mid tier fits.

4. As a kit maintainer, I want the pre-push hook body edited at its one source
   (the Go template that `bench link` installs) with the installed kit-repo hook
   refreshed in the same change, so the manifest fix cannot fork the template
   from the installed reality.
   Line: claude-sonnet-5 / medium. Mechanical propagation at a known seam with
   contract coverage, so the cheap tier at medium effort covers it.

## Implementation decisions

- **`wired` is derived, never declared.** The guards package scans the two
  harness configs for each hook script's path token — the same substring
  convention conformance uses — so there is no registry to drift. A harness
  config that is absent contributes nothing (a repo without `.codex/` simply
  can't wire Codex).
- **Schema change is additive**: `guards[N]{guard,boundary,denies}` grows a
  `wired` column; `--brief` keeps its one-line-per-guard shape but appends the
  wired harnesses so the SessionStart injection stays honest. Downstream
  consumers are the AXI contract tests and the session-start hook — both
  updated in the same change.
- **`session-start.sh` stays out of the table** — informational hooks
  (denies: nothing) are already filtered, and wiring reporting doesn't change
  that filter.
- **Pre-push `--describe` reads pin state at describe time**
  (`git rev-parse --git-path bench-gate-pin`), so the manifest always states the
  rules currently in force: pinned → both denies; unpinned → push-to-main only,
  plus an explicit "drift check disarmed — run `bench gate pin`" note in the
  `why`. (Default call, flagged: the alternative — fail-closed pushes when
  unpinned — changes adoption ergonomics for every linked repo and is the
  reviewer's decision, recorded in Out of scope.)
- **One source for the hook body.** The pre-push script lives as a template in
  the adopt package (`link_hook.go`); the edit happens there, and the kit repo's
  installed hook is refreshed through the sanctioned install path, not by
  hand-editing `.git/hooks/pre-push`.

## Testing decisions

- **What a good test is here:** drive `bench guards` (and the hook's
  `--describe` directly) in throwaway fixture repos with different wiring and
  pin states, asserting the TOON rows — never the scan internals. Prior art:
  `internal/contract/axi/axi_guards_test.go` (aggregation behavior),
  `internal/contract/surface/prepush_test.go` (installed hook behavior).
- **Seams:** the `bench guards` stdout table (AXI contract) and the installed
  pre-push hook's `--describe` output (surface contract).
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: agent or SessionStart hook runs `bench guards [--brief]`
        │
        ▼
    .bench/hooks/*.sh        ──▶  [ guards.Rows                     ]  ──▶  guards[N]{guard,boundary,denies,wired}
    .claude/settings.json    ──▶  [  discover scripts               ]        pre-push row: managed/unmanaged/
    .codex/hooks.json        ──▶  [  scan configs → wired field     ]        not installed (unchanged)
    installed pre-push       ──▶  [  pre-push install check (as-is) ]
    + pin file state         ──▶  [  hook --describe (state-aware)  ]
                      ◀ tests attach here: AXI contract fixtures wire/unwire each config and
                        assert the wired cell; surface contract toggles the pin file and
                        asserts the --describe denies line changes with it.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A hook script wired in claude but not codex renders `wired` naming exactly the wiring configs | `bench guards` stdout (AXI contract) | fixture with the kit's own asymmetric wiring — today the table has no `wired` header, so the header/cell assertion fails | a presence-only implementation cannot produce the asymmetric cell; the assertion demands the scan |
| 1 | A hook script referenced by neither config renders `wired: none` (definitive, not blank) | `bench guards` stdout (AXI contract) | fixture with an extra orphan hook script — no such cell exists today | a blank cell would violate the definitive-empty-state rule this row pins |
| 1 | `--brief` carries the wired information per guard | `bench guards --brief` stdout | today's brief lines carry no wiring — the assertion fails | the SessionStart injection is the highest-leverage reader; an un-updated brief path would keep overstating there |
| 2 | pre-push row still reports managed / unmanaged / not installed by marker | `bench guards` stdout (AXI contract) | already covered — existing aggregation tests pin these three postures; re-run after the schema change | the additive column must not regress the one wiring-true row that exists today |
| 3 | `--describe` with a pin present includes the drift clause; with no pin it drops the clause and notes the disarm | installed pre-push `--describe` (surface contract) | fixture repo, delete the pin file, run the hook with `--describe` — today the drift clause renders unconditionally, so the no-pin assertion fails | the manifest is the advertisement; this row is literally the collapse of advertisement to enforcement |
| 4 | The template source and a freshly installed hook agree byte-for-byte on the describe logic | surface contract (link installs into a fixture repo) | already covered structurally — link tests install from the template; assert the new describe behavior through a fresh install, red until the template (not the local hook) carries it | editing only the kit's installed hook would leave every future consumer with the old manifest; installing fresh in the test forces the template path |

### Edge inventory

- error path → covered: unreadable/absent harness configs contribute no wiring
  (story-1 fixtures include an absent `.codex/`); unmanaged and missing pre-push
  postures already pinned (story 2).
- empty/absent input → rows: `wired: none` definitive cell; absent pin file
  (story 3).
- malformed input → **Won't handle** in guards: an unparseable harness config
  scans as not-wired rather than erroring — guards is a read-only reporter and
  the JSON-validity conformance check owns malformedness; the decision is
  recorded here rather than silent.
- boundary values → covered: multiple hook groups per event and multiple
  harnesses per script are the kit repo's own shape (story 1).
- re-run idempotency / interrupted state — **Won't handle**: read-only
  derivation, nothing to resume.
- hostile environment (paths with spaces in config commands) → the substring
  match tolerates quoting differences by matching the script's relative path
  token; asserted implicitly by the kit-shape fixture, **Won't handle** beyond
  that (config command strings are kit-authored, not user input).

## Out of scope

- **Fail-closed pre-push when unpinned** (block pushes until `bench gate pin`
  has run) — a separate enforcement-posture capability that changes adoption
  ergonomics for every linked repo; the reviewer owns that call. Estimate:
  ~2 edits, 2 gate runs.
- **Wiring-awareness for arbitrary third-party harnesses** (beyond claude +
  codex) — a separate discovery capability; today's two configs are the shipped
  adapters. Estimate: ~4 edits, 2 gate runs per added harness format.
