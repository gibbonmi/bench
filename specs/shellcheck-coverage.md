# shellcheck-coverage — lint the enforcement shell, and skip loudly

Status: implemented

Source: `ASSESSMENT.md` backlog 9 (findings §2 med improvement, §2 low).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

The shellcheck phase lints `bin/bench.sh`, `.bench/hooks/*.sh`, and
`.bench/lib/*.sh` — but not the load-bearing enforcement shell it misses by
extension or location: the shift adapters (`.bench/adapters/claude|codex|
opencode`, extensionless), the gate entry script itself (`.bench/gate.sh`),
and the pre-push hook body (a Go string literal in the adopt package, never
lintable as shell). The phase is also Optional: where the binary is absent the
whole shell-lint defense silently skips, and a present-but-unexecutable
shellcheck (EACCES) is treated the same as absent.

## Solution

Extend the lint set to the adapters and the gate script; extract the pre-push
hook body from the Go literal into an embedded shell asset (one source the
compiler embeds and shellcheck lints); and make the Optional posture honest —
absent binary skips with a visible verdict line, while a present-but-broken
binary is a red, not a skip.

## User stories

1. As a kit maintainer, I want the shellcheck phase to lint the three shift
   adapters and `.bench/gate.sh`, so the shell that arms the loop and resolves
   the oracle is inside the lint defense, not beside it.
   Line: claude-sonnet-5 / medium. Extending a file-list builder at a known
   seam, with the profile's mid-effort floor for gate logic setting the
   effort.

2. As a kit maintainer, I want the pre-push hook body to live as an embedded
   shell asset (`go:embed`) that both `bench link` installs and the shellcheck
   phase lints, so the hook's one source is lintable instead of hiding inside
   a Go string.
   Line: claude-opus-4-8 / medium. Moving an enforcement artifact between
   representations touches the install path every consumer runs, so it
   warrants the mid tier.

3. As a reviewer reading gate output, I want an absent shellcheck to produce a
   visible skipped verdict (`phase shellcheck: skipped (not installed)`), so
   the defense's absence is a fact on the record, not silence.
   Line: claude-sonnet-5 / medium. Verdict-line plumbing in the phase runner,
   observable at the gate seam.

4. As a reviewer, I want a present-but-unexecutable shellcheck (EACCES or any
   non-not-found exec failure) to fail the phase rather than skip it, so a
   broken lint binary reads as breakage, not as a clean pass.
   Line: claude-opus-4-8 / medium. Exec-error classification inside the phase
   runner is oracle correctness — the worst class of bug in this kit — so mid
   tier.

## Implementation decisions

- **Optional stays, honestly.** The shipped decision (shellcheck is best-effort,
  not a hard dependency) is respected: absent binary → skip. What changes is
  the skip's visibility (story 3) and the boundary of "absent" — only
  exec-not-found qualifies; every other exec failure is red (story 4).
- **Adapters are named explicitly, not discovered by suffix.** The three
  adapter files are extensionless by contract; the argv builder adds them by
  name alongside the suffix-scanned directories, and `.bench/gate.sh` joins
  the same list. A missing adapter file stays a conformance concern (that
  check exists), not a shellcheck concern.
- **The pre-push asset move is a pure representation change**: the embedded
  file's bytes are what `bench link` installs — behavior identical, proven by
  the existing install/behavior contract tests staying green. The shellcheck
  phase lints the asset path in the kit tree. (Build note: the guards-wiring
  spec edits the same hook body for its `--describe` change; whichever builds
  second rebases trivially, and both specs carry this note.)
- **Shellcheck severity stays `-S warning`** — the extension must not ratchet
  strictness in the same change, or new files fail for style reasons and the
  diff conflates two decisions.
- **Bite-proof**: the canary family gains a needle proving the extended set is
  really linted (a fixture adapter with a deliberate lint error goes red) —
  without it, a path typo in the argv builder silently un-lints everything
  this spec adds.

## Testing decisions

- **What a good test is here:** observe the gate phase from outside — its
  verdict lines and reaction to planted lint errors and broken binaries —
  never the argv slice. Prior art: the gate runtime contract tests
  (`internal/contract/runtime/runtime_gate_test.go`) and the canary fixture
  families under `tests/canary/`.
- **Seams:** the gate phase runner's verdict output (runtime contract, with a
  PATH-controlled shellcheck), and the canary meta-gate for the bite-proof.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: `bench gate` (shellcheck phase), driven by runtime contract with a
             controlled PATH
        │
        ▼
    bin/bench.sh + hooks + lib     ──▶  [ shellcheck phase              ]  ──▶  phase verdict:
    + adapters + .bench/gate.sh    ──▶  [  argv over the extended set   ]        green / red /
    + embedded pre-push asset      ──▶  [  exec classification:         ]        skipped (not installed)
    PATH: real | absent | broken   ──▶  [   not-found→skip, else→red    ]
                      ◀ tests attach here: runtime contract runs the gate with shellcheck
                        absent (skip line), with a chmod-000 stub (red), and the canary
                        needle plants a lint error in a fixture adapter (red, targeted).

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A lint error planted in a fixture adapter turns the gate red with the targeted substring | canary meta-gate | the needle fixture doesn't exist and adapters aren't linted — planting the error today leaves the gate green | proves the extended set is actually linted; an argv typo that drops the adapters fails here, not in review |
| 1 | `.bench/gate.sh` is in the linted set | shellcheck phase (runtime contract) | already covered once story 1's needle pattern exists — a second needle or the same fixture family covers it; red today by the same green-on-planted-error probe | same silent-un-lint failure mode as the adapters |
| 2 | `bench link` installs a pre-push hook byte-identical to the embedded asset, and behavior postures are unchanged | surface contract (link + prepush tests) | already covered — the existing install and hook-behavior tests pin the bytes' effect; they stay green through the representation move, stated openly as regression cover rather than new TDD | the risk is divergence during extraction; the existing suites are exactly the net for it |
| 2 | The embedded asset path is in the linted set | shellcheck phase (runtime contract) | plant a lint error in the asset in a fixture — green today (the body is a Go string no linter reads) | the whole point of the extraction; red proves shellcheck now sees the hook |
| 3 | With shellcheck absent from PATH, the phase reports a visible skipped verdict | gate output (runtime contract) | run the gate with a PATH omitting shellcheck — today the skip is silent, so asserting the skip line fails | absence becomes a fact on the record; a regression to silence fails the assertion |
| 4 | With a present-but-unexecutable shellcheck on PATH, the phase goes red naming the exec failure | gate output (runtime contract) | chmod-000 stub on PATH — today this is treated as skip (gate green), so the red assertion fails | the EACCES-masking defect is directly asserted away |

### Edge inventory

- error path → rows: broken binary (4), planted lint errors (1, 2).
- empty/absent input → row: absent binary (3); an empty adapters directory is
  a conformance diagnostic already, **Won't handle** here.
- boundary values → covered: the set now spans suffix-scanned dirs, named
  extensionless files, and an embedded asset — each class has a row.
- malformed input → **Won't handle**: shellcheck's own parsing owns malformed
  shell; the phase only routes files to it.
- interrupted/partial state, re-run idempotency — **Won't handle**: the phase
  is a single exec with no state.
- hostile environment (tool missing from PATH) → rows 3 and 4 are exactly the
  profile's missing-tool class, split into its two honest cases.

## Out of scope

- **Ratcheting shellcheck severity** (`-S style` or error-on-info) — a
  separate lint-policy decision over the whole existing set. Estimate:
  ~2 edits, 3 gate runs (the extra runs for triaging newly-flagged lines).
- **Linting consumer-repo hooks at link time** — a separate adoption-surface
  capability (bench link running lint on install); the kit gate only defends
  the kit tree. Estimate: ~3 edits, 2 gate runs.
