# Second-wave check parsers: `bench diff` and `bench coverage`

## Problem
Two review-phase signals are still hand-derived, and one of them is derived
wrong. `/bench-review-implementation` pins its diff against the default branch,
but a shift branch's true base is the pre-shift HEAD — a shift stacked on
unmerged work reviews the wrong diff today. And the acceptance-coverage audit
hand-parses spec markdown tables that the gate already parses separately with
its own embedded validator — two derivations of the same signal, the exact bug
class wave 1 (`decisions/state-surface.md`) was built to end.

## Solution
Two AXI-conformant query subcommands beside the wave-1 trio, sharing the same
TOON emitter, error helpers, and gate conformance pattern. `bench diff` is the
single source of review-base truth: the shift loop records the pre-shift HEAD
in branch-scoped git config at branch creation, and `diff` resolves recorded
base first, merge-base with the default branch as fallback, then lists the
changed files. `bench coverage <spec>` parses the existing gate-enforced
acceptance-coverage-map convention and reports mechanical state plus extracted
rows; the gate's embedded validator is replaced by the same parser, leaving one
derivation. `/bench-review-implementation` re-points to both. Decisions:
`decisions/second-wave-parsers.md` (closed).

## User stories
1. As a reviewer, I want `bench shift` to record the pre-shift HEAD in
   `branch.<name>.benchBase` at branch creation, so the true review base
   survives the loop instead of dying with a local variable.
   Line: claude-sonnet-4-6 / low. One config write at a known point in the
   shift loop, fully observable by the gate's stub-adapter shift test.
2. As a review agent, I want `bench diff` to resolve the review base as the
   recorded `benchBase` key when it names a reachable commit, falling back to
   merge-base with the default branch when the key is absent or its commit is
   gone, so stacked shifts review their own work and plain branches keep
   working. Line: claude-sonnet-4-6 / medium. Small resolution ladder in shell
   at the wave-1 seam, but the fallback ordering carries the bug fix.
3. As a review agent, I want `bench diff` to print scalar preamble lines
   (branch, base sha, resolution method) followed by a definitive
   `files[N]{status,path}` TOON table of changed files, so the review phase
   consumes one structured surface instead of assembling git plumbing.
   Line: claude-sonnet-4-6 / medium. Plumbing over `git diff --name-status`
   through the existing TOON emitter, all gate-observable.
4. As an agent in a broken context, I want `bench diff` to answer with a
   structured stdout error and exit 1 outside a git repository or when no base
   is resolvable, and usage plus exit 2 on unknown arguments, so failures are
   parseable instead of silent. Line: claude-sonnet-4-6 / low. The exact AXI
   error posture wave 1 already established, copied to one more command.
5. As a review agent, I want `bench coverage <spec>` to parse the canonical
   acceptance-coverage-map convention and emit the spec's mechanical state
   (mapped, historical, or no-map) plus one TOON row per data row, so the
   coverage audit is grounded in extraction instead of ad-hoc markdown
   parsing. Line: claude-sonnet-4-6 / medium. A table parser at a known seam
   against a convention that is already pinned by fixtures.
6. As an agent, I want `bench coverage` to answer a missing argument or unknown
   flag with usage and exit 2, and a nonexistent spec path with a structured
   error and exit 1, so misuse is distinguishable from an unmapped spec.
   Line: claude-sonnet-4-6 / low. Same AXI posture as story 4.
7. As a kit dev, I want the gate's embedded node coverage validator replaced by
   a `bench coverage --check` mode invoked per spec — resolved relative to the
   gate script's own tree so canary fixtures keep tripping — preserving every
   current error condition (canonical header, five cells, non-empty cells,
   story references within the spec's numbering, historical opt-out), so the
   convention has exactly one parser. Line: claude-opus-4-8 / medium. This is
   oracle surgery: the profile routes gate and conformance logic to the mid
   tier because a wrong gate is the worst bug class in this repo.
8. As a kit dev, I want the existing coverage canaries to keep biting through
   the re-pointed check, and new canary fixtures proving the `diff` and
   `coverage` AXI contract checks bite, so the oracle's teeth are proven, not
   assumed. Line: claude-opus-4-8 / medium. Canary attribution is gate logic
   and routes mid for the same reason as story 7.
9. As a reviewer, I want `/bench-review-implementation`'s pin-the-diff step to
   direct agents to `bench diff` for base and file list (git remaining for the
   content diff), so the phase consumes the fixed base instead of re-deriving
   the wrong one. Line: claude-fable-5 / high. Command prose is guidance that
   compounds through every future session — the profile's leverage override
   routes it top; this spec's approval is the escalation ask.
10. As a cold agent, I want the dispatcher, `bench` usage text, `.bench/BENCH.md`
    command list, and the benchkit profile's AXI seam list to name `diff` and
    `coverage`, so the surface is discoverable and the docs-currency gate
    checks stay green. Line: claude-sonnet-4-6 / low. Mechanical list edits
    the gate already polices.

## Implementation decisions
- Both parsers live beside the wave-1 trio in the query module, sharing
  `toon_table`, `axi_error`, and `axi_usage`. No new files, no new dependencies;
  the CLI stays plain shell (portability is the product — node exists in this
  repo but not in every linked repo, which is why the validator port goes to
  awk/bash rather than the CLI shelling to node).
- `bench diff` takes no positional arguments (current branch only; `-h` and
  usage-error handling per AXI). Resolution ladder: `branch.<name>.benchBase`
  config if set and `git cat-file -e <sha>^{commit}` — else merge-base with
  `default_branch()` — else structured error. The preamble names which method
  resolved (`recorded`, `merge-base`, `merge-base (recorded sha unreachable)`)
  so a review agent can see when a fallback happened. Detached HEAD has no
  branch config by construction and takes the merge-base path.
- File rows come from `git diff --name-status --no-renames <base>...HEAD`
  (three-dot, matching the review phase's convention). Rename detection is
  disabled so every row is a flat `{status,path}` pair — statuses A/M/D and
  friends, never the two-path R rows.
- The shift loop writes the config key immediately after creating the shift
  branch, before iteration 1 — the value is the `base` it already computed.
  Orphaned keys after branch deletion are harmless residue the fallback
  ignores (closed in the decision map).
- `bench coverage <spec>` output: scalar `spec:` and `state:` preamble
  (mapped | historical | no-map), then `rows[N]{story,seam,red_signal}` for
  mapped specs (behavior and why-cells stay in the file; the three emitted
  fields are what the review audit keys on). `--check` emits nothing on a
  valid map and structured errors + exit 1 on violations — the gate's mode.
  Judgment statuses (missing / partial / falsely-classified) remain review's;
  the parser only extracts.
- The gate's docs-contract fragment loops `specs/*.md` invoking
  `<gate-script-dir>/../bin/bench.sh coverage --check` — resolved from the
  gate script's own location, never the working tree's copy, because canary
  inner runs execute the real gate against minimal fixture trees that carry
  no CLI. Error message substrings keep the current phrasing so existing
  canary EXPECT files keep matching.

## Testing decisions
- Good tests here exercise the CLI as a subprocess in throwaway fixture repos
  and assert stdout shape + exit codes — the established AXI-contract pattern.
  No unit-testing of internal functions.
- Seams: (1) the CLI query surface, tested by the gate's AXI contract
  fragment; (2) the gate contract itself, tested by canary fixtures; (3) the
  shift loop's recording, tested by the existing stub-adapter shift pattern.
  All three exist today; no new seam is introduced.
- Gate command: `bench gate` (`.bench/gate.sh`).

### Seam diagram — CLI query surface (`bench diff`, `bench coverage`)

    trigger: review agent / gate AXI fragment runs the subcommand
        │
        ▼
    git state (branch, config,      ┌──────────────────────────────┐
    commits) ─────────────────────▶ │ query module: resolve base /  │ ──▶ TOON on stdout
    specs/<f>.md ────────────────▶  │ parse coverage map            │ ──▶ exit 0/1/2
                                    └──────────────────────────────┘
                      ◀ tests attach here: fixture repo + subprocess run,
                        assert preamble lines, table header/rows, exit code

### Seam diagram — gate coverage check (single parser)

    trigger: `bench gate` (and canary inner runs against fixture trees)
        │
        ▼
    specs/*.md in cwd ──▶ [ gate fragment → bench coverage --check ] ──▶ red + substring
                                                                          or silence
                      ◀ tests attach here: canary fixtures (broken map trips,
                        conformant tree passes), EXPECT substrings unchanged

### Seam diagram — shift base recording

    trigger: `bench shift "<objective>"` with a stub BENCH_AGENT
        │
        ▼
    pre-shift HEAD ──▶ [ shift_loop: create branch, write benchBase ] ──▶ branch + config key
                      ◀ tests attach here: after the stub shift, assert
                        `git config branch.<name>.benchBase` == recorded HEAD
                        and `bench diff` resolves method `recorded`

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | after a stub-adapter shift, `branch.<name>.benchBase` equals the pre-shift HEAD | stub shift in fixture repo | `git config` lookup on the shift branch currently returns empty (key never written) | fails if the loop stops recording or records the wrong sha |
| 2 | on a branch stacked on unmerged work, `bench diff` base equals the recorded sha, not merge-base with default | CLI in fixture repo | `bench diff` today prints the general usage banner and exits 0 (no subcommand) | this is the review-base bug; the row is red until the ladder prefers the key |
| 2 | with the key absent, base falls back to merge-base with the default branch; method says `merge-base` | CLI in fixture repo | same red as above | fails if fallback breaks plain (non-shift) branches |
| 2 (edge) | key present but sha unreachable → merge-base fallback, method names the fallback | CLI in fixture repo | same red as above | a pruned or fixture-copied key must degrade loudly, not error or lie |
| 3 | changed files emit as `files[N]{status,path}`; N matches; a path with spaces round-trips quoted | CLI in fixture repo | same red as above | fails if rows drop, misalign, or the emitter stops escaping |
| 3 (edge) | no changes since base → definitive `files[0]{status,path}:` and exit 0 | CLI in fixture repo | same red as above | an empty diff must be a definitive answer, not an error or silence |
| 4 | outside a git repo: structured error on stdout, exit 1; unknown argument: usage, exit 2 | CLI outside repo | same red as above | fails if error posture drifts from the AXI contract |
| 4 (edge) | unborn default branch / no merge-base and no key → structured error naming the cause, exit 1 | CLI in fixture repo | same red as above | an unresolvable base must be a parseable failure, not raw git noise |
| 5 | a mapped spec emits `state: mapped` + one row per table data row with story/seam/red-signal cells | CLI in this repo's fixture | `bench coverage <spec>` prints the usage banner and exits 0 today | fails if extraction misparses the canonical table |
| 5 (edge) | escaped `\|` inside a cell stays one cell; story ranges (`2–4`) and `edge of N` rows parse | CLI in fixture repo | same red as above | the cell grammar the node validator handles must survive the port |
| 5 (edge) | `historical` marker → `state: historical`, zero rows, exit 0; spec without the heading → `state: no-map`, exit 0 | CLI in fixture repo | same red as above | both are definitive states, distinct from errors |
| 6 | missing spec argument or unknown flag → usage, exit 2; nonexistent path → structured error, exit 1 | CLI in fixture repo | same red as above | misuse must be distinguishable from an unmapped spec |
| 7 | `--check` on a malformed map exits 1 with the current error phrasings (header, cell count, empty cell, story out of range) | CLI in fixture repo | node validator currently emits these from inside the gate; `bench coverage --check` does not exist | fails if any validation rule is lost in the port |
| 7 | the gate goes red on a broken map via the re-pointed check and green on this repo | `bench gate` | already red pre-port via the node block — parity must hold after | proves the swap didn't weaken the oracle |
| 8 | `broken-coverage-map` and `coverage-axis-anchor` canaries still bite; new canaries trip the `diff`/`coverage` AXI checks | canary layer | new fixtures' EXPECT substrings unmatched until checks exist | a check without a biting fixture rots into an always-pass |
| 9 | pin-the-diff step names `bench diff` as the base/file-list source | docs anchor check | anchor grep for the command in the phase file fails today | fails if the phase keeps re-deriving the wrong base |
| 10 | dispatcher routes both; usage text, `.bench/BENCH.md` list, and profile seam list name them | `bench gate` docs-currency checks | `bench diff` falls through to the usage banner today | fails if the surface ships undiscoverable |

### Edge inventory
Walked per the profile's hostile-input checklist (shell CLI):
- paths with spaces/globs in changed files → coverage row (story 3).
- absent vs empty: missing spec file vs no-map vs historical vs empty diff →
  coverage rows (stories 3, 5, 6).
- unquoted multi-word arguments → `diff` takes none; `coverage` takes one path;
  usage-error rows cover strays (stories 4, 6).
- interrupted mid-loop → **Won't handle:** the config write is a single atomic
  git command; a shift killed before it behaves exactly as today (fallback).
- re-run idempotency → **Won't handle:** both commands are read-only; each
  shift creates a fresh timestamped branch, so keys never collide, and
  orphaned keys are decided-harmless in the map.
- hostile environment (no global `bench`, symlink invocation) → **Won't
  handle:** unchanged from wave 1; the commands add no new tool dependency.
- cwd deeper than repo root → both resolve root via `git rev-parse` like wave
  1; exercised implicitly by the AXI fragment's fixture layout (subdir run
  included in story 3's assertions).
- CRLF / missing trailing newline in a hand-edited spec table → coverage row
  (story 5 edge, the cell grammar assertions include a CRLF line).

## Out of scope
- **`bench refs` / `bench doctor` / `bench detect`** — separate capabilities
  parked on the roadmap pending learnings-funnel evidence (decision map #1).
  Estimate if later admitted: ~8–12 edits, ~6 gate runs each.
- **A `--base <ref>` override on `bench diff`** — a distinct
  review-arbitrary-ranges capability nobody consumes yet; the phase reviews
  the current branch. Estimate: 3 edits, 2 gate runs.
- **Emitting `behavior`/`why` cells from `bench coverage`** — a wider-schema
  capability; the audit keys on story/seam/red-signal and the file stays the
  source for prose cells. Estimate: 2 edits, 2 gate runs.
