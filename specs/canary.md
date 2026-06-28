# Gate canary — prove the gate still bites

## Problem
benchkit is self-hosted: the gate (`.bench/gate.sh`) is both the oracle and code
under active development. It runs its conformance checks against the repo and exits
0/1, but nothing verifies the checks still *catch* anything. A check that silently
always-passes — edited wrong, or rotted by a refactor — would let the gate go green
while the oracle has lost its teeth. The gate `bash -n`s the CLI and hooks but never
tests its own checks against known-bad input. (This is the gap behind the real
regressions in `cea2f42` and the `learnings.md` bug: structurally valid, behaviorally
broken, caught only after the fact.)

## Solution
A **fixture canary** inside the gate. Keep a small set of deliberately-broken kit
fixtures, each targeting one failure mode a real check claims to catch. The gate runs
itself against each fixture and asserts the run goes RED *and* emits that fixture's
expected error substring. If a fixture stops biting, the corresponding check has
rotted and the gate fails — the oracle proves itself on every run.

## User stories
1. As a kit dev, I want the gate to fail if any covered conformance check stops
   catching its target, so a rotted oracle can't pass silently.
2. As a kit dev, I want one fixture per real failure mode — bad frontmatter, a
   dangling AGENTS.md↔skill index reference, an unscaffolded `.bench/*` contract file,
   an extensionless `.bench/gate` reference, invalid JSON, a missing `files[]` path —
   so the checks that caught the historical regressions stay proven.
3. As a kit dev, I want each fixture to assert its *specific* failure substring, so a
   passing fixture means *that* check bit (attribution), not that some unrelated check
   happened to fire.
4. As a kit dev, I want the canary to run as part of `bench gate`, so the guard runs
   whenever the oracle does — no fourth check surface to remember.
5. As a kit dev, I want the canary not to recurse, so the inner gate run against a
   fixture doesn't re-trigger the canary.
6. As a kit dev, I want fixtures minimal (the broken file plus an expected substring),
   so they don't drift or need maintenance when unrelated parts of the kit change.
7. As a kit dev, I want a clear convention that adding a new conformance check comes
   with adding its fixture, so coverage and proof grow together.

## Implementation decisions
- **Seam — one.** The gate invoked as a subprocess against a fixture, observed via
  exit code + stderr. This exercises the real check logic (the actual gate code runs),
  not a reimplementation. It generalizes the existing check-1d pattern (run the real
  path in a throwaway dir, assert the outcome).
- **Canary lives in `.bench/gate.sh`** as a new check (the gate is benchkit-specific
  and not shipped, so this guards benchkit's own checks only). It runs near-last so
  per-check failures it depends on are already defined.
- **Attribution by substring, not by isolation.** Because the gate accumulates
  failures (`err` sets `fail=1` and continues), a minimal fixture over-fails on
  unrelated checks — that is fine. The canary asserts RED **and** that the targeted
  substring is present; it never asserts the *absence* of other errors. This is what
  lets fixtures stay minimal and drift-free: new checks only add noise to inner runs,
  never false-green the canary. *(This is the one call to veto — the alternative,
  near-complete valid fixtures broken in one spot, reads cleaner but rots with every
  kit change. Chosen against it for maintenance.)*
- **Fixture layout:** `tests/canary/<name>/` holds `EXPECT` (the substring) and
  `files/` (the broken kit tree). The canary copies `files/` into a throwaway
  `git init` dir, runs the gate there with a recursion guard, and asserts.
- **Recursion guard:** the inner run is invoked with `BENCH_CANARY_INNER=1`; the
  canary check skips itself when that var is set. No fixture needs `git add` — `git
  init` alone satisfies the gate's `git rev-parse`, and working-tree globs see
  uncommitted files.
- **Not shipped:** `tests/` stays out of `package.json` `files[]` (dev-only, like
  `.bench/` and `decisions/`).

## Testing decisions
- **What a good test is here:** run the real gate against a known-broken fixture and
  assert the external verdict — red exit + the expected message. Behavior at the seam,
  never gate internals.
- **Seam + prior art:** subprocess-gate-in-tmpdir, observed by exit/stderr. Prior art
  is check 1d (`bench init` in a throwaway repo, assert a file appears).
- **No clean red — disclosed.** The canary is a meta-check over checks that already
  pass, so it goes green on first write rather than red-then-green. Its red-*capability*
  is validated once during the build: temporarily neutralize one real check (e.g. make
  the frontmatter check always-pass), confirm the matching fixture's substring vanishes
  and the canary fails, then revert. That demonstration is the "red" and is recorded,
  not committed.
- **Gate command:** `bash .bench/gate.sh` (the project gate; the canary becomes part
  of it). "Done" = gate green with the canary included, across all six fixtures.

## Out of scope
- **Behavioral skill-firing tests** — whether a skill actually triggers and produces
  sane output. A separate capability: needs a running-agent harness, not a
  deterministic gate (decided in `decisions/dogfooding.md` #3). ~3–4h of agent time.
- **CLI-loop canary** — proving `bench shift`/worktree behavior end-to-end (commit on
  green, block on red). Distinct seam (the loop, not the gate); its own future spec.
  ~1–2h.
- **Gate self `bash -n`/shellcheck** — your call: you chose "fixture canary" over
  "both layered," so syntax-self-check is deliberately excluded (a broken gate fails to
  parse loudly anyway). ~5 min if reopened.
