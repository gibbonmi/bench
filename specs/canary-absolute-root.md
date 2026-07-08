# canary-absolute-root

Status: staged

Compiled from roadmap row FT40 (a reproduced defect, not a map — the reviewer
directed it into `specs/` as a build target). Reproduced today: `bench canary .`
from the repo root reports every canary "did not bite ... got exit 127". Class:
cwd/path assumption — profile checklist line "cwd deeper than the repo root when
the command assumes root".

## Problem

`bench canary <root>` accepts a root argument and joins the inner gate off it:
`gate := filepath.Join(root, ".bench", "gate.sh")`. When the argument is
relative (`bench canary .`), that gate path is relative too (`.bench/gate.sh`).
Each canary fixture then runs its inner gate from a fresh temp cwd
(`cmd.Dir = <materialized fixture temp dir>`), so the relative gate path resolves
against that temp cwd, not the repo root. The file isn't there, bash exits 127,
and every fixture reports `did not bite ... got exit 127` — the harness that
proves the gate's checks bite is itself silently broken whenever it's invoked
with a relative root.

Worse failure mode: a fixture whose `files/` tree carries its own
`.bench/gate.sh` materializes that file into the temp cwd. The relative gate
path then resolves to the *fixture's* gate, so the sweep runs the wrong gate
entirely and can report a false verdict instead of an obvious 127.

## Solution

Absolutize the root at the entry to the sweep, before any path is derived from
it, so the gate path is absolute and resolves against the repo root regardless
of the inner fixture's cwd. `filepath.Abs` also cleans the argument, so a
trailing slash collapses. Every downstream caller (the `bench canary` CLI entry
and any future caller of the sweep) inherits the fix from one place.

## User stories

1. As a kit maintainer running `bench canary .` (or any relative root) from the
   repo root or a subdirectory, I want the canary sweep to run each fixture's
   inner gate against the repo's real `.bench/gate.sh`, so that fixtures bite as
   designed instead of every one reporting `did not bite ... got exit 127`, and
   so that a fixture carrying its own `.bench/gate.sh` can never shadow the real
   gate.
   Line: claude-sonnet-5 / medium. This is a one-line absolutize at a known
   seam with a pre-agreed red-capable test, which is exactly the canary and
   conformance logic the profile routes to the cheap tier at medium effort.

## Implementation decisions

- **Seam: `Sweep(root, runner)` in `internal/canary/canary.go`.** Absolutize
  `root` at the top of `Sweep`, before `fixtures(...)` and before the
  `gate := filepath.Join(root, ".bench", "gate.sh")` derivation, via
  `filepath.Abs`. On the `filepath.Abs` error (only when the cwd is
  unresolvable), return that error rather than proceeding with a relative path.
  `Sweep` is chosen over the CLI `Run` entry because it is the single point
  where the gate path is derived from root and the reusable, unit-testable
  boundary where the canary tests already attach; `Run` (which resolves root
  from args or `git.Root`) inherits the fix through it, and `git.Root` is
  already absolute so absolutizing it again is idempotent.
- **No change to `defaultRunner` or the fixture materialization.** The bug is
  the relative gate path, not the inner cwd — the inner gate *must* run from the
  materialized fixture cwd; the gate path just has to be absolute so that cwd
  stops mattering.
- `filepath.Abs` cleans the path, so a trailing-slash root (`bench canary ./`)
  needs no separate handling.

## Testing decisions

- A good test here observes `Sweep`'s external behavior — the verdict it returns
  and the gate path it hands the runner — not internals. The seam is
  `internal/canary/canary_test.go` (package `canary`), which already drives
  `Sweep` with recording/stub runners; prior art:
  `TestSweepMaterializesFixtureAndRequiresTargetedBite` and
  `TestSweepUsesLiteralFixturePathWithSpacesAndGlobCharacters`.
- **Pre-agreed test double — a "resolving runner"** that reproduces
  `defaultRunner`'s cwd semantics without spawning bash: it resolves the gate as
  bash would (`if !filepath.IsAbs(call.Gate) { call.Gate = filepath.Join(call.Cwd, call.Gate) }`),
  then returns:
  - the fixture's targeted diagnostic at exit 1 when the resolved path is the
    real repo gate (the fixture bites),
  - a distinct "shadow gate ran" output at exit 0 when it resolves to some other
    existing file (the fixture's own materialized gate), and
  - exit 127 with a "No such file or directory" line when the resolved path is
    absent.

  This double turns the exact field symptom (`got exit 127`) into a deterministic
  unit assertion and lets one test cover both failure modes.
- **Run the sweep from a cwd that is not the repo root, with a relative root
  argument that is not literally `.`** (e.g. `t.Chdir` into the parent and pass
  the root's basename). This keeps the degenerate fix — special-casing the
  string `"."` — red, because an arbitrary relative root must also be
  absolutized.
- Gate command: `bench gate` (the Go tests run under it).
- Go 1.25: use `t.Chdir` for the cwd change.

### Seam diagram

    trigger: bench canary <relative-root>  (CLI Run → Sweep),
             and the canary_test.go resolving-runner double
        │
        ▼
    root arg (relative, e.g. ".")  ──▶  [ Sweep: filepath.Abs(root),  ]  ──▶  RunCall{Gate: <absolute>/.bench/gate.sh, Cwd: <fixture temp dir>}
    fixture temp cwd               ──▶  [ then join gate off abs root, ]  ──▶  inner gate resolves against repo root, not the temp cwd
                                        [ run each fixture's gate      ]
                      ◀ tests attach here: canary_test.go drives Sweep with a
                        relative root from a different cwd and a resolving runner;
                        observes the returned verdict and the recorded Gate path

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | relative root run from a different cwd → gate handed to the runner is absolute and resolves to the repo gate; fixture bites | `internal/canary/canary_test.go` — `Sweep` with resolving runner, `t.Chdir` off the repo root, relative root arg ≠ `.` | pre-fix `Sweep` returns `canary 'relroot' did not bite (want red + "targeted diagnostic"; got exit 127)`; post-fix returns nil, and `filepath.IsAbs(recordedGate)` is true | with a relative gate path the inner cwd resolves it to a nonexistent file → 127; asserting IsAbs on an arbitrary relative root also kills the `=="."`-only degenerate fix |
| 1 (edge) | fixture carrying its own `.bench/gate.sh` cannot shadow the real gate | same seam — fixture `files/` tree materializes a `dot-bench/gate.sh`; resolving runner tags which gate it resolved | pre-fix `Sweep` returns `canary 'shadow-gate' did not bite (want red + "targeted diagnostic"; got exit 0)` (the fixture's own gate ran, emitting no diagnostic); post-fix returns nil and the recorded gate equals the repo gate, not the fixture's | a relative gate path resolves to the fixture's materialized gate in the temp cwd, running the wrong gate and reporting a false verdict; the absolute path pins it to the repo gate |
| 1 (edge) | trailing-slash root is cleaned | same seam — relative root passed with a trailing separator | covered by the same assertions: post-fix `recordedGate` is absolute with no trailing/doubled separator (`filepath.Abs` cleans) | `filepath.Abs` applies `Clean`; a doubled or trailing separator in the recorded gate path would show the clean was skipped |

Degenerate-implementation check: the cheapest wrong fix rewrites only the exact
string `"."`; the row-1 assertion runs with a relative root that is not `.`, so
that stub stays red.

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist and the canonical
classes:

- **cwd deeper than / different from repo root when the command assumes root**
  (the profile's named class, the defect itself) — coverage row 1.
- **Relative root argument** (`.`, `sub/dir`, basename from parent) — coverage
  row 1.
- **Root with trailing slash** — coverage row 3 (`filepath.Abs` cleans).
- **Fixture carrying its own `.bench/gate.sh`** — coverage row 2 (shadow gate
  cannot win once the path is absolute).
- **Already-absolute root** (the `git.Root` path from bare `bench canary`) —
  **Won't handle** as a new case: `filepath.Abs` is idempotent on an absolute
  path, and the existing `Sweep` tests (all use `t.TempDir`, which is absolute)
  already assert this path stays green.
- **`filepath.Abs` error** (unresolvable cwd) — **Won't handle** with a test:
  it fires only when the process cwd has been removed, is not reproducible in a
  standard test, and is handled by returning the error rather than proceeding.
- **Paths with spaces / glob characters** — already covered by
  `TestSweepUsesLiteralFixturePathWithSpacesAndGlobCharacters`; absolutizing
  changes nothing there (`filepath.Abs` preserves literal bytes).
- **Re-run idempotency** — n/a: `Sweep` is stateless per call; no persisted
  state changes.

## Out of scope

- **Absolutizing root in other subcommands that take a root argument** (e.g.
  `bench structure`, `bench guards`) — a separate audit of each command's own
  cwd assumptions, only warranted if a defect is reproduced there; this spec
  fixes the one reproduced surface — ~1 edit + ~1 test each, ~2 gate runs, per
  command, if a repro appears.
- **A CLI-level integration test that runs `bench canary .` through a real
  built binary against real fixtures** — a separate, heavier runtime-contract
  capability; the `Sweep`-seam unit test pins the behavior deterministically and
  cheaply, and the real end-to-end path is exercised by the gate's own canary
  sweep on every run — ~2 edits, ~2 gate runs.
