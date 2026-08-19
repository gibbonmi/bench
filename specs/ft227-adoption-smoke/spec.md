# ft227-adoption-smoke

Status: staged

Decision source: `roadmap/FT227.md` — the board row the reviewer opened from the 2026-08 capability audit (ledger L-14) and reconciled onto the board in the drain committed as `eefc96f4` on 2026-08-19; its one occurrence is the audit's end-to-end reproduction, re-run in a scratch repository while this spec was authored.

Verification log: 1 iteration to accept — the round (`opus` / high, read-only) returned ten partials and no blocking finding; all ten were folded here: SM3 made `--fresh`, SD5 moved to a zero-signal fixture with an honest why-clause, the declared-tool and ungraded-binary edges written as Won't handle, the `BENCH_HOME` mutation probe recorded, the exact `gate: green` line assertion, the private `BENCH_HOME` bound on the setup launch too, the `paths` justification restated, `bench link` named beside `bench init`, and the round's own answer (the guard alone makes the fresh green) written into Problem and Solution.

## Problem

A repository adopted by `bench setup` cannot run its own gate green, even after
the operator does the one thing the scaffold asks for.

`bench setup` on a repository with no detected build system writes the
fail-closed stub `scaffoldGate()` produces. That stub ends by validating the
canary inventory through the installed wrapper:
`"$bench" canary "$root" || err "canary inventory validation failed"`. Two
things make that line red on every fresh repository:

1. **The gate environment is closed and nothing declares what the wrapper
   needs.** `bench gate` launches the gate script with `PATH` plus only the
   names declared under `environment` in `.bench/gate-inputs.json`. `bench
   setup` writes no such file, so the script runs with `PATH` alone. The
   installed wrapper's first statement, `export
   BENCH_HOME="${BENCH_HOME:-$HOME/.bench}"` under `set -u`, dies with
   `HOME: unbound variable` before it can route anything. The kit's own
   `.bench/gate-inputs.json` declares exactly this (`BENCH_HOME`, `HOME`) for
   exactly this reason; a linked repository receives no equivalent.
2. **An empty inventory is an error, and every new repository has one.** `bench
   canary` reports `canary fixture inventory is empty` and exits 1 when
   `tests/canary` does not exist. That is correct for the kit and for
   `prep-release`, which rely on it. The scaffolded stub, however, treats that
   exit as a gate failure, so a repository that has not yet written a fixture is
   red for a reason the operator did nothing to cause.

The two reds are sequential only because the stub is unguarded: the guard alone
makes a fresh repository green, because a repository with no `tests/canary` never
invokes the wrapper. The manifest becomes load-bearing from the moment the
adopter creates a fixture directory, which is the first moment the wrapper runs
inside the gate.

The audit reproduced this end to end: `bench setup --yes`, remove the sentinel
line as the stub instructs, `bench gate` — red on `HOME`, then, after a
hand-written manifest, red on the empty inventory. The same sequence was re-run
in a scratch repository on 2026-08-19 with the current tree and produced the same
two reds, and the same scratch repository went green once a manifest declaring
`BENCH_HOME` and `HOME` was present and the canary call was guarded on the
inventory directory's existence.

Nothing in the kit's gate observes adoption from the adopter's side. The system
package installs the kit into a disposable repository and runs the wrapper, but
never runs an adopted repository's scaffolded gate. So the failure is invisible
to the oracle and was found by an audit, not by a red.

## Solution

Three small changes and one piece of evidence.

- `bench setup` seeds `.bench/gate-inputs.json` — schema 1, local closure,
  `environment` declaring `BENCH_HOME` and `HOME`, no paths, the wrapper's own
  tool set — through the same transaction that writes the gate and the profile.
  It is a seed: written when absent, never touched when present, never recorded
  in the link manifest, so an operator's later edit can never read back as a
  modified-managed conflict. The plan preview names it like every other write.
- The scaffolded stub guards its inventory call on the existence of
  `$root/tests/canary`. A repository with no inventory directory skips
  validation; a repository with one still validates it, and a present-but-empty
  directory still reds as an empty inventory. `bench canary` itself does not
  change. `bench init` and `bench setup`'s zero-signal branch keep sharing the
  one `scaffoldGate()` source.
- The guard alone turns the fresh repository green; the manifest is what keeps it
  green once the adopter writes a fixture. Both land here so that the first
  fixture does not reintroduce the audit's red.
- The sentinel stays. A fresh stub is still red until the operator removes the
  sentinel line; this spec makes the gate green *after* that documented step,
  not before it.
- One system journey, under the kit gate's `system` phase, adopts one of the
  owner's disposable repositories with `bench setup --yes` and runs its
  scaffolded gate through the installed wrapper: red with the sentinel, green
  once the sentinel line is gone, green again with one project fixture in
  `tests/canary` (which proves the wrapper ran under the closed environment and
  resolved the installed binary), red naming `HOME` once the seeded manifest is
  removed, and red naming the empty inventory when `tests/canary` exists but
  holds nothing.

## User stories

### Seeding the gate input declaration

Line: `opus` / medium. One seed entry in an existing transaction plus a preview
line; the adopt package's conventions for plan entries are established and the
content is a short constant.

1. As an adopter, I want `bench setup` to write `.bench/gate-inputs.json`
   declaring `BENCH_HOME` and `HOME`, so that the scaffolded gate's wrapper call
   runs under the closed gate environment instead of dying on an unbound
   variable.
2. As an adopter, I want the seeded manifest to be one the gate accepts as valid
   — schema 1, `local` closure, all five fields present — so that the
   declaration is load-bearing rather than read as "manifest invalid" and
   silently dropped.
3. As an adopter, I want the plan preview to name the manifest it will seed, or
   say it is leaving mine alone, so that nothing lands silently.
4. As an adopter, I want a `.bench/gate-inputs.json` I already own to be left
   byte-identical and never recorded in the link manifest, so that a later hand
   edit never reads back as a modified-managed conflict.
5. As an adopter, I want the manifest seeded whether or not a build system was
   detected, so that the declaration does not depend on the shape of the
   proposed gate command.
6. As an adopter, I want a second `bench setup` run to leave the seeded manifest
   byte-identical, so that re-running adoption is idempotent.
7. As a teammate who just walked in, I want the platform reference to say that
   `bench setup` seeds this file, so that I know where the declaration came from
   before I extend it.

### Guarding the scaffolded inventory call

Line: `opus` / medium. A three-line change inside one generated shell script and
the string test that pins it; the edge is the present-but-empty directory.

8. As an adopter, I want the scaffolded gate to skip inventory validation when
   `tests/canary` does not exist, so that a new repository's gate does not fail
   on an inventory it has not created.
9. As an adopter, I want the scaffolded gate to still validate the inventory
   once `tests/canary` exists, so that a broken inventory stays red.
10. As an adopter, I want a present-but-empty `tests/canary` to red the gate
    naming the empty inventory, so that a directory I created for fixtures is
    not silently accepted as nothing.
11. As a maintainer, I want `bench init`'s scaffold and `bench setup`'s
    zero-signal stub to stay one script, so that the guard cannot drift between
    the two writers.
12. As a maintainer, I want `bench canary` unchanged — an empty inventory
    remains an error — so that the kit's own conformance and `prep-release`
    semantics are untouched.
13. As an adopter, I want the sentinel to keep a fresh stub red until I remove
    it, so that this fix does not fabricate a green gate.
14. As an adopter, I want the guard to resolve `tests/canary` against the
    repository root, so that the gate gives the same answer from any working
    directory.

### Proving adoption goes green

Line: `opus` / medium. A system journey in an established tagged package with
existing process helpers; the care goes into the legs' ordering and into keeping
the operator's home untouched.

15. As a maintainer, I want a system journey that adopts a disposable repository
    with `bench setup --yes`, so that adoption is exercised end to end by the
    kit's own gate.
16. As a maintainer, I want that journey to see the fresh stub red with the
    sentinel remedy before any edit, so that the fail-closed posture is observed,
    not assumed.
17. As a maintainer, I want the journey to remove exactly the sentinel line —
    the documented operator step — and then see `bench gate` green through the
    installed wrapper with no `tests/canary` present, so that the guard is proved
    at the real seam.
18. As a maintainer, I want the journey to add one project fixture under
    `tests/canary` and see the gate green with the inventory reported as one
    binding, so that the wrapper is proved to run under the closed environment
    and to resolve the installed `.bench/dist/bench`.
19. As a maintainer, I want the journey to remove the seeded manifest and see the
    gate red naming `HOME`, so that the declaration is proved load-bearing rather
    than decorative.
20. As a maintainer, I want the journey to empty `tests/canary` and see the gate
    red naming the empty inventory, so that the present-but-empty edge is pinned.
21. As a maintainer, I want the journey to bind a private `BENCH_HOME` and leave
    it empty, so that a kit test never writes under the operator's home.
22. As a maintainer, I want the journey to use one of the owner's three
    disposable repositories and the selected executable only, so that the
    system package's repository budget and identity ledger hold.
23. As a maintainer, I want the journey to run under the kit gate's `system`
    phase, so that the evidence is oracle-owned rather than a script someone
    remembers to run.
24. As a teammate who just walked in, I want the project profile to name the
    adoption journey beside the stripped-distribution one, so that the system
    package's advertised shape matches the tree.

## Implementation decisions

- **The manifest is a seed, not a managed asset.** It uses the existing `seed`
  plan-entry kind, exactly as the starter profile does: written through the one
  FT84 transaction, skipped when a file already exists at the path, never
  recorded in `link-manifest.tsv`. An operator extends this file by hand (more
  environment names, declared tools, paths); a managed kind would turn every such
  edit into a conflict on the next `bench setup` or relink.
- **One Go source for the seeded bytes.** A `scaffoldGateInputs()` function
  beside `scaffoldGate()` and `scaffoldProfile()` returns the exact content. Its
  value: `{"schema": 1, "closure": "local", "environment": ["BENCH_HOME",
  "HOME"], "paths": [], "tools": ["bash", "basename", "dirname", "git",
  "readlink", "uname"]}`, indented, trailing newline. `environment` is what the
  wrapper needs: `HOME` for its first statement, `BENCH_HOME` so an operator's
  override reaches the wrapper the gate runs instead of being silently replaced
  by `$HOME/.bench`. `tools` is the wrapper's own tool set, which is the kit's
  list minus the kit-build-only `go` and `node`. `paths` is empty: the gate
  script is hashed by resolution and the installed wrapper by the tree hash,
  while `.bench/dist/bench` is gitignored by design and stays outside the
  subject either way (see the edge inventory). Nothing outside this function
  spells the content.
- **The preview names the seed.** `renderSetupPreview` gains one line in the
  same voice as the profile line: absent → will be seeded declaring `BENCH_HOME`
  and `HOME`; present → left as-is. `inspectRepo` gains the one fact.
- **The guard is the stub's, not `bench canary`'s.** `scaffoldGate()` wraps its
  inventory call: `if [ -d "$root/tests/canary" ]; then … fi`. Directory
  existence is the whole predicate — a present-but-empty directory falls through
  to `bench canary`, which reports the empty inventory and reds the gate. The
  path is rooted at `$root`, never relative, so the answer is the same from any
  working directory. `bench canary`, its empty-inventory message, and every
  caller of `Inventory` are untouched; the kit's own gate does not route through
  this stub.
- **The sentinel marker becomes exported.** The journey removes the sentinel
  line by the same marker the doctor row and the setup remedy print. The
  constant is exported from `internal/adopt` and the journey reads it there, so
  the marker stays one source.
- **Journey placement and budget.** The journey is a new tagged file in
  `internal/systemtest` (the owner file is near its line budget). It adopts
  `owner.repos[1]`: the only journeys that touch that repository run `canary
  <kit>` and an unknown command from it and never read its tree, so the
  repository budget stays at three. Setup runs on the selected executable with
  `BENCH_KIT` at the real kit; every gate run goes through the installed
  `.bench/bin/bench.sh` with `BENCH_RUN_BINARY` pointed at the selected
  executable. Every launch the journey makes — setup and each gate
  leg — binds one private `BENCH_HOME` under the test's temporary directory,
  through the owner's `observeSelected` plus `runAt` pattern (the
  reauthorize journey's precedent), because `runSelected` and `runWrapper`
  carry fixed environments with no `BENCH_HOME` override. The stub is not a
  phase-table gate, so `gate-run` injects no `BENCH_RUN_BINARY` into it: inside
  the gate the environment is `PATH` plus the two declared names, and the stub's
  wrapper resolves the `.bench/dist/bench` copy `bench setup` installed — the
  production path, the same bytes as the selected executable.
- **Leg order and freshness.** The legs run in the order the stories list them.
  The first gate run is plain `bench gate`; every later leg passes `--fresh`, so
  a leg whose subject happens to hash identically to an earlier one (an empty
  directory does not change the tree hash) cannot be answered from a reusable
  green. A green leg asserts the exact `gate: green` line, never the substring:
  a reused verdict prints `gate: green (fresh verdict reused for this tree)`,
  and the substring would accept it.
- **`bench setup` exits 3 on a zero-signal repository by design** — the doctor's
  gate row is red while the sentinel is present. The journey asserts that exit
  and the red row, not a zero.
- **The reference doc gains one clause.** `.bench/BENCH-reference.md`'s sentence
  on `.bench/gate-inputs.json` says `bench setup` seeds it with the names the
  installed wrapper needs. The profile's system-package paragraph names the
  adoption journey beside the stripped-distribution one.

## Testing decisions

- The external behavior a good test exercises: an empty repository, adopted by
  `bench setup --yes` and edited only as the stub instructs, runs `bench gate`
  green through its installed wrapper; and each of the two scaffold changes is
  load-bearing, because removing the manifest or emptying the inventory
  directory turns that same gate red for the named reason.
- Seams. The adopt package's existing ordinary tests already run `setup` against
  a temporary git repository with `BENCH_KIT` set (`setup_prompt_test.go`'s
  fixture); the seeded bytes, the preview line, the preserved operator file, and
  idempotency attach there. The stub's text attaches to the existing
  `scaffoldGate()` string test, whose canary-line expectation changes to the
  guarded form. The end-to-end behavior attaches to the tagged system package,
  whose owner already provides process helpers, the selected executable, and
  disposable repositories.
- Gate seam: the ordinary `test` phase for the adopt package; the `system` phase
  for the journey. No shell or gate edit in the kit.
- `bench canary`'s empty-inventory behavior is pinned by its existing tests; this
  spec adds no row for it beyond the journey's present-but-empty leg.
- Mutation probes recorded in the verification log, not retained: (a) write the
  manifest as `inline` rather than `seed` → the preserved-operator-file row reds;
  (b) guard on `[ -d tests/canary ]` without `$root` → no ordinary row reds, the
  string row is the catch; (c) drop `HOME` from the seeded environment → the
  journey's fixture leg reds with the wrapper's own error; (d) drop `BENCH_HOME`
  from the seeded environment → only SD1 reds, which is the omission that
  licenses SD1's independently authored spelling of the bytes.

### Seam diagram

    adopter: bench setup --yes                    adopter: bench gate (via .bench/bin/bench.sh)
        │                                             │
        ▼                                             ▼
    inspectRepo ──▶ renderSetupPreview             gate-run: PATH + declared env ──▶ [ .bench/gate.sh stub ]
        │              ◀ preview test                                                   │
        ▼                                                                               ├─ sentinel? ──▶ red
    convergeSetup ──▶ one transaction ──▶ .bench/gate.sh        (scaffoldGate)          ├─ [ -d $root/tests/canary ] ──▶ no ──▶ green
                                      ──▶ .bench/gate-inputs.json (scaffoldGateInputs, seed)  └─ yes ──▶ .bench/bin/bench.sh canary
                                      ──▶ projects/<name>.md      (seed)                           │  (needs HOME; resolves .bench/dist/bench)
                       ◀ adopt tests: bytes, preserved file, idempotency                           ▼
                                                                                             bench canary ──▶ ok / "inventory is empty"
                       ◀ system journey: the whole path, five legs, private BENCH_HOME

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SD1 | 1, 2 | after `setup --yes` in a temporary repository, `.bench/gate-inputs.json` exists with exactly the seeded bytes: schema 1, `local`, environment `BENCH_HOME` and `HOME`, empty paths, the wrapper's tool list | adopt ordinary test over the setup fixture | a setup that writes nothing, or writes a shape the gate's strict loader rejects, leaves the wrapper call without `HOME` |
| SD2 | 3 | the plan preview names `.bench/gate-inputs.json` as about to be seeded when absent, and as left as-is when present | adopt ordinary test over `renderSetupPreview` | a silent write violates the preview's "nothing is acted on silently" contract and is invisible to an operator reading the plan |
| SD3 | 4 | a pre-existing `.bench/gate-inputs.json` with operator content survives `setup --yes` byte-identical and does not appear in `link-manifest.tsv` | adopt ordinary test | a managed kind overwrites the operator's declaration or records it, so the next setup reports it as a conflict |
| SD4 | 6 | a second `setup --yes` leaves the seeded file byte-identical | adopt ordinary test | a writer that rewrites or appends turns re-adoption into a diff |
| SD5 | 5 | `setup --yes` in a zero-signal repository (no `go.mod`, `Makefile`, `package.json`, or `Cargo.toml`) also writes the seeded bytes | adopt ordinary test over a second, zero-signal fixture | the existing fixture carries `go.mod`, so a seed attached only to the detected-ecosystem branch passes SD1 and leaves the stub — the one gate that calls the wrapper — without `HOME`; the declaration must not depend on the inferred gate command |
| SG1 | 8, 11, 14 | `scaffoldGate()` contains the inventory call wrapped in `if [ -d "$root/tests/canary" ]`, rooted at `$root`, and `bench init` writes that same text | adopt string test over `scaffoldGate()` and the init round trip | an unguarded call reds every fresh repository; a relative path answers differently by working directory; a second writer drifts |
| SG2 | 13 | `scaffoldGate()` still carries the sentinel line and marker | adopt string test | a stub without the sentinel is a fabricated green gate |
| SG3 | 12 | `bench canary` on an absent and on an empty `tests/canary` still reports `canary fixture inventory is empty` and exits 1 | existing canary inventory test, unchanged | a fix that softened the command instead of the stub would change `prep-release` and the kit's own inventory semantics |
| SM1 | 15, 22 | `bench setup --yes` in `owner.repos[1]` exits 3, prints the red sentinel doctor row, and leaves `.bench/gate.sh`, `.bench/gate-inputs.json`, `.bench/bin/bench.sh`, and `.bench/dist/bench` on disk | system journey | adoption that does not install the wrapper or the binary copy would make every later leg fail for the wrong reason |
| SM2 | 16 | before any edit, `bench gate` through the installed wrapper exits 1 and its stderr names the sentinel remedy | system journey | a stub that is green before configuration is the fabricated-gate failure |
| SM3 | 17, 8 | after deleting the sentinel line, `bench gate --fresh` through the installed wrapper exits 0 and prints the exact line `gate: green` with no `tests/canary` present | system journey | the original red: an unguarded inventory call, or an unbound `HOME` reached through it, reds here |
| SM4 | 18, 9, 2 | with `tests/canary/<family>/<fixture>/` holding one file, `bench gate --fresh` exits 0, prints `canary inventory ok (1 fixture bindings)`, and prints the exact line `gate: green` | system journey | proves the wrapper ran under the closed environment with `HOME` declared and resolved the installed binary; a manifest the loader rejects, or a wrapper that cannot find its binary, reds here |
| SM5 | 19 | with `.bench/gate-inputs.json` removed and the fixture still present, `bench gate --fresh` exits 1 and its stderr contains `HOME: unbound variable` | system journey | a gate that is green without the manifest would mean the declaration is decorative and the closure is leaking |
| SM6 | 20, 10 | with the manifest restored and `tests/canary` present but empty, `bench gate --fresh` exits 1 and its stderr contains `canary fixture inventory is empty` | system journey | a guard on emptiness rather than existence passes here and silently accepts an inventory the operator started |
| SM7 | 21 | every launch the journey makes binds the private `BENCH_HOME`, and after every leg that directory is still empty | system journey | a leg that writes under `BENCH_HOME` would, in production, write under the operator's home |

Not covered: story 7 — one clause of reference prose; review reads it.
Not covered: story 23 — placement in the tagged package, which the kit gate's `system` phase runs whole; no assertion adds information.
Not covered: story 24 — profile prose; review reads it.

### Edge inventory

- **Absent vs present-but-empty:** `tests/canary` absent is green (SM3);
  present-but-empty is red naming the empty inventory (SM6). The manifest absent
  is red naming `HOME` (SM5); the manifest present is green (SM4).
- **Paths with spaces or glob characters:** the owner's disposable repositories
  are named `repository [journey]-…`, so every journey leg runs under a root with
  a space and brackets; the guard quotes `"$root/tests/canary"` (SG1).
- **Cwd deeper than the root:** the stub `cd`s to `$root` and the guard is rooted
  there (SG1, story 14).
- **Re-run idempotency:** second `setup` (SD4); a second `bench gate` after a
  green is answered from the record, which is why the later legs pass `--fresh`.
- **Required tool missing from PATH:** the journey's gate environment carries no
  `BENCH_RUN_BINARY` and no global `bench`; the stub resolves the repo-local
  wrapper and the wrapper resolves `.bench/dist/bench` (SM4). A *declared* tool
  missing from `PATH` is the next line.
- **Won't handle:** a declared tool (`readlink`, `uname`, …) absent from the
  host's `PATH` — the gate opens the subject as "declared tool unavailable" and
  still runs and returns the script's exit code; no green is retained or
  reusable, and nothing reds. The cost is verdict caching, not correctness, and
  the missing coreutil fails the wrapper on its own.
- **Won't handle:** `.bench/dist/bench` is outside the gate subject — link
  gitignores it, so it is in neither the tree hash nor the seeded `paths`, and
  swapping that binary leaves a green reusable. The kit's own gate does not
  route through the stub, and declaring the binary would tie every adopter's
  verdict to an arch-specific untracked file; an adopter who wants it graded
  adds it to `paths` by hand.
- **Invocation through every shipped surface:** setup through the selected
  executable, the gate through the installed wrapper, the inventory through the
  stub's wrapper call (SM1–SM4).
- **Hand-edited file without a trailing newline:** the seed is written with one;
  an operator's own file is never read by setup, only left alone (SD3).
- **A live symlink where a directory is expected:** `tests/canary` as a symlink
  to a real directory is followed by `-d` and validated, the same as git would
  materialize it.
- **Won't handle:** a dangling symlink or a special file at `tests/canary` —
  `-d` answers false and the stub skips validation, as for an absent directory;
  the adopter who planted it sees it in `git status`, and `bench canary` run by
  hand still reports the inventory as empty.
- **Won't handle:** `HOME` genuinely unset in the operator's own environment —
  the declaration passes through what exists; the wrapper's `HOME: unbound
  variable` is the remaining diagnostic, and naming it better is a separate
  wrapper change (Out of scope).
- **Won't handle:** `BENCH_HOME` unset when `bench gate` is invoked on the
  binary directly rather than through the wrapper — the subject opens as
  "declared environment unavailable" and the gate still runs; the stub's inner
  wrapper then defaults `BENCH_HOME` from `HOME`. The profile's cold-session
  notes already record this.
- **Won't handle:** a hand-emptied or malformed `.bench/gate-inputs.json` —
  the gate's strict loader already reports it as invalid and opens the subject;
  the observable failure is the same `HOME` red as an absent manifest, and the
  operator's remedy is the same file.
- **Won't handle:** repositories linked before this change — they receive the
  seed on their next `bench setup` run, not on `bench upgrade` (Out of scope).

## Ownership fences

Tickets are serial on one retained integration source. Reviewer disposition:
approve, merge, or split at sign-off.

- `internal/adopt/setup.go`
- `internal/adopt/init.go`
- `internal/adopt/doctor_rows.go`
- `internal/adopt/setup_report.go`
- `internal/adopt/adopt_test.go`
- `internal/adopt/setup_prompt_test.go`
- `internal/adopt/setup_test.go`
- `internal/systemtest/`
- `.bench/BENCH-reference.md`
- `projects/benchkit.md`
- `specs/ft227-adoption-smoke/`
- `capture/session-handoff.md`

## Out of scope

- **`bench init` and `bench link` seeding the manifest** the way `bench setup`
  does: 3 edits (the seed write in `Init`, the same in `Link`'s plan, their
  tests), 1 gate run. `init` is the older gate-only scaffold and `link`
  installs assets without a gate; the row names `setup`.
- **The wrapper naming its own failure on an unbound `HOME`** (`: "${HOME:?…}"`
  ahead of the `BENCH_HOME` export): 2 edits (`bin/bench.sh`, a wrapper test),
  1 gate run. A diagnostic improvement, not part of making adoption green.
- **A detected-ecosystem leg in the journey** (a `Makefile` with a `test:`
  target, or a `go.mod`, adopted and gated green): 2 edits (one leg, one
  fixture), 1 gate run. The row's evidence is the empty repository.
- **Seeding the manifest into already-linked repositories on `bench upgrade`:**
  3 edits (an upgrade-path seed, its preview, a test), 1 gate run. Upgrade
  converges managed assets; a seed there is a new policy.
- **A doctor row for an absent `.bench/gate-inputs.json`:** 2 edits (the row,
  its test), 1 gate run.

## Further notes

The scratch reproduction that grounded this spec: `bench setup --yes` in an
empty repository with a private `BENCH_HOME`, then the sentinel line deleted,
then `bench gate` through the installed wrapper — `HOME: unbound variable`,
red. With a hand-written manifest declaring `BENCH_HOME` and `HOME` — `canary
fixture inventory is empty`, red. With the guard added to the stub — green with
no inventory, red naming the empty inventory with an empty `tests/canary`, green
reporting one binding with one project fixture, and red naming `HOME` again once
the manifest was removed. The private `BENCH_HOME` stayed empty throughout.
