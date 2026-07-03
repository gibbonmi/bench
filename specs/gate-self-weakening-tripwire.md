# Gate self-weakening tripwire

Decision map: `decisions/gate-self-weakening.md` (closed, 5 tickets resolved).

## Problem
Both `run_gate` and the Stop hook execute the **working-tree** `.bench/gate.sh`, so
an agent that weakens that file commits red work green. Benchkit itself is protected —
its `tests/canary/` harness makes a weakened check stop biting, which turns the gate
red — but that protection does not ship. `bench init` scaffolds a canary-less gate,
and the runner silently skips when `tests/canary/` is absent. A consumer repo
therefore has no automated defense against a lazy agent gutting its own oracle; only
the human reviewer stands between a self-weakened gate and a green commit. This is the
one quality-assessment finding that touches invariant #1 directly.

## Solution
Ship the tripwire the kit already uses on itself, and make its absence loud.

1. **Extract** the canary runner from `.bench/gate.sh` into a shared, shipped lib
   (`.bench/lib/canary-run.sh`) that benchkit's own gate and every consumer's
   scaffolded gate both source — one source, both consume.
2. **err-if-absent**: the shared runner goes *red* when `tests/canary/` is absent or
   empty, so the lazy escape (`rm -rf tests/canary/`) becomes a red gate and a visible
   diff instead of a silent pass.
3. **Scaffold a live harness**: `bench init` plants a working canary harness plus two
   scaffolded checks — an un-canaried *sentinel* (keeps a fresh gate red until
   configured) and a canaried *forbidden-token example* (proves the harness executes
   and is the copy-paste template for real canaries).
4. **Document the decided state** so future audits stop re-flagging the working-tree
   seam as an open hole: an ADR (the repo's first) and a `craft-gate` note.

We defend against the **lazy** agent, not a determined adversary — nothing in the
writable tree can stop an adversary who deletes the tripwire in the same edit. The
defense is loudness: a red gate plus a reviewable diff.

## User stories

1. As a kit dev, I want the canary runner extracted verbatim into
   `.bench/lib/canary-run.sh` and sourced by benchkit's gate, so its logic has one
   home instead of being duplicated when a consumer scaffold needs it.
   Line: claude-opus-4-8 / medium. A wrong source order can silently disarm the
   oracle without turning the gate red, so the move earns the mid model even though it
   is mechanical.

2. As a consumer, I want the shared runner to turn the gate red when `tests/canary/`
   is absent, so deleting the harness to escape a failing check is a red gate and a
   visible diff, not a silent pass.
   Line: claude-opus-4-8 / medium. This is oracle behavior whose correctness is the
   whole point of the feature, and it must not regress benchkit's own inner runs.

3. As a consumer, I want the shared runner to turn the gate red when `tests/canary/`
   exists but holds no fixture, so an emptied-out harness is caught the same as a
   deleted one.
   Line: claude-opus-4-8 / medium. Same oracle-correctness surface as story 2; the two
   absent/empty cases are one code path with two asserted inputs.

4. As a consumer, I want `bench init` to scaffold a gate that is red until I configure
   it, via an un-canaried sentinel check, so a freshly-initialised gate can never
   commit real work green before I have wired real checks.
   Line: claude-opus-4-8 / medium. Every consumer inherits this scaffold; the profile
   routes gate/conformance logic to mid effort because a wrong oracle is the worst bug
   class in the kit.

5. As a consumer, I want `bench init` to scaffold a live canary harness — a
   forbidden-token example check plus its seed fixture under
   `tests/canary/example/` — so the tripwire defends from day one and I have a
   worked template to copy for each real check I add.
   Line: claude-opus-4-8 / medium. The seed fixture must be non-vacuous against the
   scaffolded gate's own empty-repo baseline, a correctness constraint the gate grades.

6. As a consumer, I want the scaffolded harness to stay non-vacuous and biting —
   removing the sentinel yields a green gate whose seed canary still passes — so the
   scaffold I ship is a correct example, not a broken one.
   Line: claude-opus-4-8 / medium. Vacuous-EXPECT and did-not-bite are the two ways the
   seed can be silently wrong; both are gate-observable and worth mid effort.

7. As a kit dev, I want a second `bench init` to leave an existing gate and harness
   untouched, so re-running setup never clobbers a configured gate.
   Line: claude-sonnet-4-6 / low. Idempotence is a guard already present in `init`
   (`[ ! -e ]`); this only extends it to the new scaffolded paths, fully gate-observable.

8. As a reviewer, I want the working-tree-honored gate recorded as a decided tradeoff
   in an ADR — accepted limit, lazy-not-adversarial threat model, and the noted
   green-window residual — so a future audit reads it as accepted rather than
   re-raising it as an open hole.
   Line: claude-fable-5 / high. Doc authoring is the profile's leverage override:
   decided-state prose that governs every future scan is high-value per token spent.

9. As an agent authoring or editing a gate, I want a `craft-gate` note explaining why
   the tripwire exists and that deleting it is a loud, reviewable act, so the guidance
   that shapes gate edits carries the intent, not just the mechanism.
   Line: claude-fable-5 / high. Skill prose compounds through every gate-authoring
   session that loads it; the leverage override applies.

## Implementation decisions

- **New shipped lib `.bench/lib/canary-run.sh`.** Holds the canary block currently
  inline in `.bench/gate.sh` (the `BENCH_CANARY_INNER` guard, the vacuous-EXPECT
  attribution baseline, and the per-fixture biting loop). Benchkit's gate sources it at
  the #7 position, replacing the inline block, exactly as it already sources the
  `gate-*.sh` fragments from `$gate_dir`. **Lib contract:** the caller must have
  defined `root`, `err`, and `fail` in scope before sourcing (benchkit's gate and the
  new scaffold both do); the lib assumes the gate under test lives at
  `$root/.bench/gate.sh`. `.bench/lib/` is already copied into consumers by
  `bench link` and linted by shellcheck.

- **err-if-absent lives inside the `BENCH_CANARY_INNER != 1` guard.** The current
  `&& [ -d tests/canary ]` skip becomes an explicit branch: outer run + absent-or-empty
  harness → `err "canary harness absent — tests/canary/ has no fixtures; the gate
  cannot prove its own checks bite"`. Placement inside the inner-guard is load-bearing:
  inner fixture runs (`BENCH_CANARY_INNER=1`) still skip the whole block, so err-if-absent
  never fires recursively and never trips benchkit's own inner runs.

- **`init()` scaffold grows from a bare stub into a minimal real gate.** It gains the
  `root`/`err`/`fail` harness and `exit "$fail"` tail, then in order: the **sentinel**
  (`err "configure .bench/gate.sh — replace this sentinel with real checks"`, fires
  unconditionally), the **forbidden-token example** (greps tracked files for a planted
  marker and `err`s on a hit), and `. "$gate_dir/lib/canary-run.sh"`. It also scaffolds
  `tests/canary/example/EXPECT` (the example check's error substring) and
  `tests/canary/example/files/` containing one file carrying the marker. The marker is a
  literal chosen not to occur in a clean repo. Existing `[ ! -e ]` idempotence guards
  extend to the new paths.

- **Sentinel is deliberately un-canaried.** It fires on every run including the
  attribution baseline's empty-repo run, so canarying it would be flagged vacuous by
  construction. It exists to be deleted; its rot does not matter. Only the example
  check is canaried.

- **err-if-absent is a behavioral contract, not a canary.** Because it lives inside the
  inner-guard, no inner fixture run can reach it, so it cannot be a canary fixture. It
  is tested like the existing link/status contracts: a throwaway repo, an outer gate
  run, assert red + substring. This is consistent with how benchkit already tests its
  behavioral (non-conformance) surfaces.

- **ADR is the repo's first** (`docs/adr/` does not exist yet). Written path-free per
  `craft-adr`; paired with a one-line decided-state comment at the gate seam pointing to
  it, so a code-reading scan finds the marker inline.

- **Out of scope by construction:** pinning the gate outside the writable tree
  (hash-verified in pre-push) — the adversarial defense — is a separate capability, not
  built here (see Out of scope).

## Testing decisions

- **A good test here** drives a real `bench` subprocess against a tree and reads its
  exit code and stderr — never a reading of the diff. Both seams already exist in
  `gate-link-contracts.sh`; no new seam is introduced.
- **Seams (two, both existing prior art):**
  - `bench gate` invoked as a subprocess against a fixture tree — fronts the canary
    runner, err-if-absent, and the scaffolded checks. Prior art: the canary block
    (`specs/canary.md`) and every `gate-*.sh` fragment.
  - `bench init` invoked in a throwaway repo — fronts the scaffolder. Prior art:
    `gate-link-contracts.sh:5-8` (the learnings/gate scaffold assertion).
- **Gate command:** the project gate, `bench gate` (runs all fragments, the canary
  meta-check, and shellcheck on the new lib).

### Seam diagram

Seam A — the gate as oracle over a tree (canary runner + err-if-absent):

    trigger: `bench gate` (outer) / the canary loop re-invoking the gate (inner)
        │
        ▼
    tree with tests/canary/*     ──▶  [ .bench/gate.sh                 ]  ──▶  exit 0/≠0
    tree with NO tests/canary    ──▶  [   sources lib/canary-run.sh    ]  ──▶  stderr "gate: …"
    BENCH_CANARY_INNER=1 set     ──▶  [   (inner run skips the block)  ]
                                        ◀ tests attach here: run `bench gate` in a
                                          throwaway repo, assert exit ≠ 0 and the
                                          targeted substring (or exit 0 when expected)

Seam B — the init scaffolder:

    trigger: `bench init` in a fresh `git init` repo
        │
        ▼
    empty repo   ──▶  [ bench.sh init() ]  ──▶  .bench/gate.sh (sentinel + example + source)
                 ──▶  [                 ]  ──▶  tests/canary/example/{EXPECT,files/…}
                                              ◀ tests attach here: after init, assert the
                                                scaffolded files exist and that running the
                                                scaffolded gate is red via the sentinel;
                                                with the sentinel removed, green + canary passes

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | after extraction, benchkit's own canary suite still runs and bites | Seam A (benchkit gate) | `bench gate` on benchkit goes red if a known fixture stops biting after the move; run before wiring the source line | if the lib is mis-sourced the fixtures are never checked and a deliberately-broken fixture in the suite stops turning the gate red |
| 2 | absent `tests/canary/` on an outer run → red + "canary harness absent" | Seam A (throwaway) | new contract fragment: throwaway repo, no `tests/canary/`, `bench gate`; fails first because the pre-change runner skips silently (exit 0) | if err-if-absent is missing the run exits 0 and the assertion on red+substring fails |
| 3 | present-but-empty `tests/canary/` → red + same substring | Seam A (throwaway) | same fragment, second input: `mkdir tests/canary` with no fixture, assert red | an empty dir passing `[ -d ]` would slip through a dir-only check; asserts the emptiness branch |
| 4 | freshly-`init`ed gate is red via the sentinel | Seam B | throwaway `bench init`, run scaffolded gate, assert exit ≠ 0 + sentinel substring; fails first because today's stub is a bare `exit 3` with no sentinel string | a scaffold that came up green (or red for the wrong reason) would fail the substring assertion |
| 5 | `init` scaffolds `tests/canary/example/{EXPECT,files/}` | Seam B | extend the `gate-link-contracts.sh:7` post-init assertion to require the example fixture; fails first because init does not create it today | if init skips the harness the file-exists assertions fail — the same shape as the existing learnings-scaffold check |
| 6 | seed canary is non-vacuous and bites: sentinel removed → gate green and canary passes | Seam B → Seam A | in the throwaway, delete the sentinel line, run `bench gate`, assert exit 0 (proves example passes clean, canary bit its fixture, EXPECT not vacuous) | a vacuous EXPECT or a non-biting fixture makes the scaffolded gate red for a canary reason, failing the exit-0 assertion |
| 7 | second `bench init` leaves an existing gate and `tests/canary/` untouched | Seam B | throwaway: init, edit the gate, init again, assert the edit survives; fails first if the new scaffold paths lack the `[ ! -e ]` guard | re-run idempotency — a clobbering second init would overwrite the user's configured gate |
| 8 | `BENCH_CANARY_INNER=1` still skips the whole block, incl. err-if-absent | Seam A | fragment: set `BENCH_CANARY_INNER=1`, no `tests/canary/`, assert the absent-harness err does **not** fire | if err-if-absent is placed outside the inner-guard, every inner fixture run spuriously errors and benchkit's own gate breaks |
| 8 (ADR) | working-tree tradeoff recorded, path-free, with a seam comment pointing to it | — | not TDD-able (prose) — verified in the Standards review axis against `craft-adr` (no paths, decided-state only) | a scan re-flags the seam if the decided state is undocumented; the review axis, not the gate, is the check |
| 9 | `craft-gate` note carries the tripwire's intent | — | not TDD-able (prose) — Standards axis against `craft-skills`; the docs-currency gate check still asserts the skill's anchors are intact | guidance drift is a review-axis concern; the gate only guards the skill's structural anchors |

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist:

- **absent vs present-but-empty `tests/canary/`** — both asserted, rows 2 and 3.
- **re-run idempotency (second `init`)** — row 7.
- **interrupted/partial state (recursion via inner run)** — row 8.
- **boundary: sentinel-removed green window** — row 6 asserts the intended green; the
  residual (a consumer who removes the sentinel and adds no real check) is a decided,
  documented tradeoff (story 8 ADR), not a defect.
- **paths/filenames with spaces or glob chars** — *Won't handle*: the scaffolded
  example check greps a fixed-content seed fixture with no exotic names; the check is
  written quote-safe, but hardening a consumer's arbitrary tracked-path names is their
  gate's concern, not the scaffold's.
- **EXPECT file with no trailing newline** — *Won't handle*: inherits the existing
  canary harness's `cat`/`grep -qF` handling unchanged (`specs/canary.md`); this feature
  adds no new EXPECT-parsing path.
- **invocation through a symlink / cwd deeper than repo root** — *Won't handle here*:
  the scaffolded gate resolves `root` via `git rev-parse --show-toplevel`, the same
  handling benchkit's gate already proves (`gate-runtime-contracts.sh:44`); no new code
  path.
- **required tool missing from PATH** — *Won't handle*: the scaffold uses only `grep`,
  `git`, and `mktemp`, already assumed present by every existing fragment.

## Out of scope

- **Adversarial gate pinning** — hash-verifying `.bench/gate.sh` outside the writable
  tree in the `pre-push` hook, so a determined agent can't weaken the gate at all. A
  separate capability (a new enforcement surface, not the rest of this one) and
  explicitly a different threat model than the lazy-agent defense built here. Estimate
  to build later: ~6 edits, ~4 gate runs (pre-push hook change, a manifest of the
  pinned hash, a contract fragment, a canary). Parked on the roadmap so it survives this
  spec.
- **A meta-check requiring one canary per `err()`** — the roadmap's "canary-per-check"
  enforcement. Distinct from this feature (it grades the whole gate's coverage, not the
  tripwire) and already parked. Estimate: ~4 edits, ~3 gate runs.
