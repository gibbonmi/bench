# Enforcement verification

Status: staged

## Problem

Two of Bench's harness-independent backstops don't verify themselves, so a repo can quietly lose the protection they promise:

1. **The agent-line guard fails open on an omitted model.** When a delegation envelope carries no `resolvedModel`/`model`, the guard warns and allows. But an omitted model inherits the invoking session's model — the exact silent-escalation path invariant #2 exists to stop. The counter-rule lives only in `craft-delegate` prose, not in enforcement.
2. **Nothing verifies the pre-push hook after a clone.** `bench doctor` checks the PATH shim but never that `.git/hooks/pre-push` exists, is Bench-managed, or isn't diverted by `core.hooksPath`. git does not clone hooks, so a fresh clone silently drops the default-branch backstop, and no ambient signal surfaces the gap.

## Solution

- Tighten the one degraded branch that is also the attack path the guard exists for: in a **routed** repo (a parseable `.bench/lines.env` with a complete tier binding), an Agent call with a missing, empty, or whitespace-only model is **denied** (exit 2) with a message naming the fix — pass a bound alias from `lines.env`. Every other degraded branch (malformed envelope JSON, absent envelope, unrouted repo, incomplete binding) keeps today's fail-open rim, because those have no binding to enforce and a broken guard must never brick delegation.
- `bench doctor` gains a **read-only** pre-push row: red when the hook is absent, present-but-foreign (no managed marker), or diverted by a `core.hooksPath` that doesn't carry a managed hook; the remedy is `bench link` and the row names the resolved install path. The same row fires on the kit repo itself (FT10).
- `bench status` gains a low-noise guards signal that fires once — on the primary checkout of a routed repo — when the pre-push is missing or unmanaged, so the gap is ambient rather than doctor-only.

## User stories

1. As a reviewer relying on the line guard, I want a delegation that omits the model **denied** in a routed, completely-bound repo, so a silent inherit of the session's model can't slip past invariant #2. Line: claude-opus-4-8 / medium. This flips a documented fail-open rim on the oracle-adjacent enforcement path, where a wrong verdict is the worst class of bug in a kit whose premise is that the guard bites.

2. As a delegating agent, I want the missing-model deny message to name the fix — pass one of the bound aliases from `lines.env` — so I can re-delegate correctly instead of guessing. Line: claude-opus-4-8 / medium. The message wording is load-bearing enforcement text that the gate grades only as a substring, so it earns the same tier as the verdict it accompanies.

3. As the owner of an unrouted repo, I want the missing-model, malformed-JSON, absent-envelope, and incomplete-binding branches to keep failing open, so a repo with no binding to enforce is never bricked by the guard. Line: claude-opus-4-8 / medium. These are the residual rims the decision deliberately preserves, and pinning them as regression guards belongs with the verdict change that could break them.

4. As a maintainer, I want `bench doctor` to show a red pre-push row when the hook is absent, present-but-foreign, or diverted by `core.hooksPath`, so a fresh clone's lost backstop is caught the next time doctor runs. Line: claude-sonnet-5 / low. It is a thin composition of the already-embedded template, managed marker, and hooks-dir resolver at a known seam.

5. As a kit maintainer, I want the same doctor row to fire on the kit repo itself and name the pre-push install path, so the kit's own backstop is verified the same way a consumer repo's is (FT10). Line: claude-sonnet-5 / low. It is the same row under a different repo, adding only that the detail names the resolved path.

6. As a maintainer, I want doctor to only report a foreign pre-push and never overwrite or relabel it, so a user-authored hook is never clobbered by a diagnostic. Line: claude-sonnet-5 / low. The pre-push row is read-only by construction, so the story is a guard that doctor's write path stays confined to the shim.

7. As a session reading `bench status`, I want a low-noise guards row — remedy `bench link`, ranked just below the worktree rows — when the pre-push is missing or unmanaged in a routed repo, so the gap is ambient without crowding the gate and git rows. Line: claude-sonnet-5 / low. It is a new signal composed from the shared pre-push classifier, slotted into the existing severity ladder.

8. As someone running pool or linked worktrees, I want the status signal to fire only on the primary checkout, so worktrees that share the main `.git` don't double-report the same hook. Line: claude-sonnet-5 / low. The dedup mirrors the existing worktree-classifier's primary-checkout detection, so it reuses a known seam.

## Implementation decisions

**Modules touched.** `internal/lines` owns the verdict change; `internal/adopt` owns both the pre-push classifier and doctor's row; `internal/status` owns the ambient signal. `.bench/hooks/check-agent-line.sh` and the adapter wiring are unchanged — the hook already shells into `bench check-agent-line`, which calls `lines.AgentLineVerdict`, so the tightened verdict propagates to the hook, the conformance suite, and the canary without touching shell.

**Verdict change (deep seam), minimal diff, no signature change.** `AgentLineVerdict(stdin, linesEnvExists, linesEnvContent) → (exitCode, stderr)` is unchanged in shape. The only new logic is in the `model == ""` branch: when the repo is routed **and** the binding is complete (all three `BENCH_TIER_*` set), an empty model denies (exit 2) with the alias-fix message; otherwise it falls through to today's fail-open warning. The present-model matching path (alias/tier match → allow; unbound → deny) is untouched. Branch order stays as today so the present-model messages don't move.

**Whitespace and empty models fold into the missing branch.** `ModelFromEnvelope` treats an all-whitespace field as blank for its non-empty test, so a whitespace `resolvedModel` correctly falls back to `model`, and an all-blank envelope yields `""`. This closes the "whitespace model value" hostile input at the parse boundary and makes the omitted/empty/whitespace cases one uniform branch in the verdict.

**Deny message content.** The missing-model deny is a distinct message from the unbound-model deny: it states the model field is missing/empty and instructs re-delegating on a bound alias, listing the bound tiers and aliases via the existing binding-describe formatting. Both remain enforcement text the gate matches by substring.

**Pre-push classifier (shared, one source).** A single exported classifier in `internal/adopt` reuses the embedded `prePushTemplate`, the `bench:managed-pre-push` marker, and the hooks-dir resolver that already live in `link_hook.go` — same package, so Handoff uncertainty flag #7 resolves to "no import cycle, reuse directly." It resolves the effective hooks directory git will use for the repo, honoring `core.hooksPath`, and classifies the pre-push there as **absent**, **foreign** (present, no managed marker), **diverted** (`core.hooksPath` set to a dir with no managed hook), or **managed**. The fingerprint is the managed marker line the template carries — not byte-identity, which would false-red across default-branch token substitution and benign template evolution.

**Doctor row is read-only.** Doctor composes the classifier into a pre-push row rendered whenever doctor runs inside a git worktree; the row is red on absent/foreign/diverted, names the resolved install path, and gives `bench link` as the remedy. Doctor never writes the hook — self-heal is rejected (least surprise), and a foreign hook is reported, never overwritten or relabeled. Doctor's exit follows its existing red posture: a red pre-push row makes doctor exit 1 even when the shim is healthy.

**Status signal.** `status` imports the `adopt` classifier (cycle-free: `adopt` imports only `git`). A new signal fires only on the **primary checkout** — detected by comparing the git dir against the git-common-dir, the same test the worktree classifier's `canonicalRoot` uses — so pool and linked worktrees sharing the `.git` don't double-report, and only in a **routed** repo (`.bench/lines.env` present). The signal sorts immediately after the worktree signals and before the drain signal; the housekeeping severities below it are renumbered as needed to keep every severity unique (the renderer relies on uniqueness) and preserve the ladder. Remedy is `bench link`.

## Testing decisions

- **What a good test is here.** Drive real behavior at the seam: the pure verdict for the lines change, and the real `bench doctor` / `bench status` CLI for the two rows — never a reading of the diff. The three seams already have test families; extend them.
- **Seams and prior art.**
  - `lines.AgentLineVerdict` — pure unit table `TestAgentLineVerdict` (plus `TestAgentLineVerdictDenyMessage`) in `internal/lines`, and black-box through the real hook via the conformance check `checkAgentHookBehavior` in `line_routing_checks_test.go`.
  - Doctor pre-push row — black-box via `bench doctor` in the `TestDoctorReportFixContracts` family (`internal/contract/surface`), alongside `testDoctorReport`.
  - Status guards signal — black-box via `bench status` in the `TestRuntimeStatusContracts` family (`internal/contract/runtime`), alongside `testRuntimeStatusWarmPool` (worktree dedup) and `testRuntimeStatusRoadmapReconcile` (ladder + primary/branch gating).
- **Built-binary note.** The contract/runtime seams exercise the compiled CLI; rebuild `dist/bench` before hand-running them.
- **Ratified posture flip (invariant #1).** The existing `checkAgentHookBehavior` case "does not fail open on a missing model field" asserts exit 0 in the routed fixture — it must be **updated to assert exit 2** with the alias-fix message, and a **new** unrouted missing-model case added asserting exit 0. This is a behavior change ratified by the closed decision map, not test-weakening; call it out in the diff so review reads it as intended.
- **Gate.** `.bench/gate.sh` (the project gate).

### Seam diagram

Seam 1 — the verdict (pure, plus the real hook):

    trigger: PreToolUse Agent hook → bench check-agent-line
        │
        ▼
    stdin envelope ──▶ [ lines.AgentLineVerdict ] ──▶ exitCode (0 allow / 2 deny)
    linesEnvExists ──▶ [  (+ ModelFromEnvelope) ] ──▶ stderr (warn or DENIED msg)
    linesEnvContent ─▶ [                        ]
                          ◀ tests attach here: unit table drives the three args;
                            checkAgentHookBehavior drives the real hook in temp routed/
                            unrouted repos and reads exit code + stderr

Seam 2 — doctor row:

    trigger: bench doctor (inside a git worktree)
        │
        ▼
    repo hooks dir ──▶ [ adopt pre-push classifier ] ──▶ absent│foreign│diverted│managed
    core.hooksPath ─▶ [   + doctorReport row       ] ──▶ stdout row + exit (1 on red)
                          ◀ tests attach here: BenchEnv(..., "doctor") over temp repos
                            with absent / foreign / diverted / managed pre-push; assert
                            row text, install path, and exit

Seam 3 — status signal:

    trigger: bench status (SessionStart hook or user)
        │
        ▼
    primary-checkout? ─▶ [ appendPrePush → render ] ──▶ board row ranked worktree<row<drain
    routed?          ──▶ [   (adopt classifier)   ] ──▶ (or omitted)
    hook state       ──▶ [                         ]
                          ◀ tests attach here: Bench("status") from primary vs pool
                            worktree; assert row present once, action, and ladder order

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | routed + complete binding + omitted model → exit 2 | `TestAgentLineVerdict` unit + `checkAgentHookBehavior` | new unit case `AgentLineVerdict(no-model envelope, exists=true, fullBinding)` expects exit 2 — fails today (returns 0); the conformance routed case flips 0→2 | the cheapest wrong impl leaves the fail-open branch, which returns exit 0 and trips both signals |
| 1 | routed + complete binding + whitespace/empty model → exit 2 | `TestAgentLineVerdict` unit | new unit case with `resolvedModel:"   "` expects exit 2 — fails today | a build that only trims omission, not whitespace, returns 0 and goes red |
| 2 | missing-model deny names the alias fix | `TestAgentLineVerdict` deny-message unit | new assertion the exit-2 stderr names the missing field, "bound alias", and the tier list — fails today | a deny that reuses the unbound wording (no alias-fix guidance) fails the substring check |
| 3 | unrouted missing-model stays exit 0 | `TestAgentLineVerdict` unit + `checkAgentHookBehavior` | new unrouted missing-model case expects exit 0 — already green (regression guard, not TDD-able red-first) | a build that denies whenever the model is empty, ignoring routing, breaks this and goes red |
| 3 | malformed JSON, absent envelope, incomplete binding stay exit 0 | `TestAgentLineVerdict` unit | existing/added cases expect exit 0 — already green (regression guards) | a routing-blind tightening that denies these degraded branches trips them |
| 4 | absent pre-push → doctor red, names path, exit 1 | `TestDoctorReportFixContracts` | new contract case: temp repo, no hook → `doctor` exit 1, stdout names "pre-push" + install path — fails today (no such row) | the cheapest wrong impl (no row) leaves doctor exit 0 with no pre-push text |
| 4 | foreign pre-push (no marker) → doctor red "not bench-managed" | `TestDoctorReportFixContracts` | new contract case: user hook lacking the marker → red row — fails today | a marker-blind check that calls any present hook healthy stays green and goes red here |
| 4 | `core.hooksPath` diverted to a dir with no managed hook → doctor red | `TestDoctorReportFixContracts` | new contract case sets `core.hooksPath` to an empty dir → red "diverted" — fails today | an impl reading only `.git/hooks` misses the divert and reports green |
| 5 | kit-repo-shaped repo → same row fires and names install path | `TestDoctorReportFixContracts` | new contract case (bench repo, no managed hook) asserts the row names the resolved path — fails today | a row gated to link-manifest-only repos would skip the kit repo and go red |
| 6 | foreign pre-push left byte-identical after doctor | `TestDoctorReportFixContracts` | new contract case reads the foreign file before/after `doctor` and asserts unchanged bytes plus a red row | a build that installs on detection (like the link path) changes the bytes and trips the assertion |
| 7 | routed repo, pre-push missing → status guards row, action `bench link`, ranked worktree<row<drain | `TestRuntimeStatusContracts` | new contract case removes the hook in a routed repo → row present with action and ladder order — fails today | no-row impl omits the text; a mis-severity impl fails the ordering assertion |
| 8 | signal fires only on the primary checkout | `TestRuntimeStatusContracts` | new contract case adds a pool/linked worktree and runs `status` from it → no pre-push row; from primary → exactly one | an impl that reads the shared `.git` from every checkout double-reports and fails the from-pool assertion |

### Edge inventory

Edge classes walked per behavior; each resolves to a coverage row above or a **Won't handle** line here.

- Empty/absent input — omitted model (rows 1, story 1), empty stdin / absent envelope → fail-open (row, story 3).
- Malformed input — non-JSON envelope → fail-open (row, story 3); foreign pre-push with no marker (row, story 4).
- Boundary — whitespace-only model folded into the missing branch (row, story 1); `core.hooksPath` set to a dir that *does* carry a managed hook → **green**, asserted as the not-diverted boundary in the diverted contract case.
- Hostile environment — pool/linked worktree sharing `.git` (row, story 8); spaces in the hooks path reuse the existing spaced-path doctor fixture convention, so path handling is not re-derived.
- Interrupted/partial state — n/a: all three surfaces are pure reads with no write or lease.
- Re-run idempotency — doctor and status are read-only and idempotent by construction; no separate row.
- Trailing-newline / hand-edited hook — the marker check is a newline-agnostic substring, so a hook whose last line lacks a newline still classifies correctly; noted, no separate row.

**Won't handle:**

- Evasion-resistant model enforcement — the guard is the honest-mistake backstop its own threat model declares; a determined agent can still bypass it, and hardening past that is a different capability, not this one.
- Doctor auto-installing or repairing the pre-push — rejected by the decision (least surprise); the remedy is `bench link`, and the row stays read-only.
- A hard gate check for pre-push presence — rejected: the gate must run green in worktrees and CI clones where hooks are legitimately absent (Handoff item 8).
- Distinguishing a non-Bench git repo from a Bench one for the doctor row — doctor assumes a Bench context; running it in an unrelated repo shows the row harmlessly, and gating it further is not worth the surface.

## Out of scope

- **FT31 — enforcement-postures ADR.** One decision record capturing the accepted honor-system residuals, including the fail-open rims this spec deliberately preserves (unrouted, incomplete binding, malformed/absent envelope). It is a separate documentation capability with its own roadmap row, not the rest of this feature — this spec records the residuals in prose but does not author the ADR. Est. ~2 edits, 1 gate run.
