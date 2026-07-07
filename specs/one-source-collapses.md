# one-source-collapses — the kit passes its own duplication audit

Status: staged

Source: `ASSESSMENT.md` backlog 6 (findings §4 med, §5 med, §6 med/low).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.
Two of the assessment's five instances resolve to decisions rather than edits
(recorded below), and one was found already gate-coupled during spec drafting.

## Problem

The kit's cardinal code standard — one source per fact — has live violations in
its own tree: the not-in-repo message exists in three phrasings
(`toon.NotInRepo`, a hand-duplicated literal in the commit command, a bare
`not in a git repo` in ~8 operational commands and the gate); the
acceptance-coverage-map schema is derived in both the write-spec command and
the `craft-tdd` skill; `bin/bench.sh`'s header comment carries a stale third
roster of the subcommand set; and the project profile's Seams section
re-enumerates a partial subcommand list that has already drifted (missing
`doctor`, `canary`, `version`, `gate pin`).

## Solution

Collapse each fact to one owner: every not-in-repo message renders
`toon.NotInRepo`'s phrase (AXI commands on stdout exit 1 as today; operational
commands the same sentence on stderr); the coverage-map schema lives only in
`/bench-write-spec` with `craft-tdd` pointing at it and keeping only the
TDD-specific row discipline; the `bin/bench.sh` header roster is deleted in
favor of the usage heredoc; and the profile's seam list points at the canonical
CLI Inventory instead of re-enumerating it.

## User stories

1. As an agent hitting any `bench` command outside a git repository, I want one
   phrase — `toon.NotInRepo`'s — on every surface (AXI commands: stdout, exit
   1; operational commands: stderr, unchanged exit codes), so the posture
   `git.go` already claims is actually true and no second phrasing can drift.
   Line: claude-sonnet-5 / medium. The change is mechanical routing through an
   existing function across known call sites, fully pinned by a contract probe,
   so the cheap tier at medium effort covers it.

2. As a spec author, I want `craft-tdd`'s Acceptance-rows section to take the
   five-field schema and red-signal definition from `/bench-write-spec` by
   pointer, keeping only the TDD-side discipline (rows go red one at a time,
   classification rules), so the schema has one home and the two most-used
   build surfaces can't diverge.
   Line: claude-fable-5 / high. Skill prose compounds through every session
   that loads it — the profile's skill-authoring leverage override applies.

3. As a maintainer reading `bin/bench.sh`, I want the header comment's stale
   subcommand roster deleted (the header keeps its one-paragraph purpose
   statement; the usage heredoc stays the file's roster), so the third copy of
   the inventory stops existing.
   Line: claude-sonnet-5 / low. A comment deletion with no behavior surface —
   cheap tier, low effort.

4. As a cold session reading the profile, I want the Seams section's CLI bullet
   to point at the canonical inventory in `.bench/BENCH.md` and state the seam
   contract (stable names + exit codes) without re-enumerating subcommands, so
   the partial second copy stops drifting.
   Line: claude-opus-4-8 / low. A pointer rewrite in project-profile prose —
   mechanical, but profile wording steers cold sessions, so it rides mid
   rather than cheap.

5. As a reviewer, I want the profile's tier-binding paragraph to name the
   conformance check that already keeps it honest (instead of the manual
   "keep it in sync" note), so the advertised coupling matches the enforcement
   that exists.
   Line: claude-opus-4-8 / low. One-sentence profile correction, same class as
   story 4.

## Implementation decisions

- **Not-in-repo: one phrase, two surfaces — the split is the operational/AXI
  contract, not the wording.** AXI query commands keep `toon.NotInRepo` on
  stdout exit 1 (unchanged). Operational commands (`commit`, gate plumbing,
  worktree/shift paths that today say `not in a git repo`) print the same
  rendered sentence to stderr with their existing exit codes. `commit.go`'s
  hand-copied literal is replaced by a call. `gate`'s exit-3 posture is
  unchanged — only the phrase unifies. (Default call, flagged: unifying exit
  codes across operational commands is *not* attempted; that would be a
  behavior change beyond the duplication fix.)
- **Coverage-map schema owner is the command, not a new skill.** The assessment
  suggested a `craft-spec` skill as one option; creating a new
  model-invocable skill is a separate capability (parked on the roadmap), and
  the command file already hosts the canonical schema today. `craft-tdd` keeps
  its row-classification rules — that content is TDD's own fact, not a copy.
- **Verified during drafting: the tier binding is already gate-coupled.** The
  assessment's §6 med claim ("no gate check couples them") is wrong —
  `checkLineBinding` cross-checks the profile prose against `.bench/lines.env`
  model ids and alias declarations. No collapse is needed; story 5 only fixes
  the profile's advertisement of a manual coupling that is actually enforced.
- **The `--full` help/doc pair is accepted as honest repetition** (the ft9
  review's "decide, and either collapse or accept"): the runtime help constant
  and the Go package doc serve different readers through different channels,
  and a shared source would contort both. Decision recorded here; no edit.
- **The parked FT6 dispatch-reconcile meta-check stays parked.** A roster
  conformance check coupling the usage heredoc to `.bench/BENCH.md` was
  considered and rejected — the reviewer parked exactly that class of
  meta-check pending evidence of real rot, and this spec respects the closed
  decision.

## Testing decisions

- **What a good test is here:** run each affected command outside any git
  repository and assert the single phrase and per-surface channel/exit — the
  only black-box-observable fact this spec changes. Prose collapses (stories
  2–5) are anchor- or reviewer-enforced and classified honestly below. Prior
  art: `internal/contract/axi/axi_fail_closed_test.go` (not-in-repo postures)
  and the conformance docs anchors.
- **Seams:** the built binary's stdout/stderr outside a repo (contract), and
  the conformance docs-anchor layer for the `craft-tdd` pointer.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: contract test runs each `bench <cmd>` from a non-repo cwd
        │
        ▼
    cwd outside any repo  ──▶  [ every bench command          ]  ──▶  AXI cmds: stdout
                               [  git.Root() fails            ]        `error: not in a git repository — …`, exit 1
                               [  → one rendered phrase        ]  ──▶  operational cmds: same sentence on
                               [    (toon.NotInRepo)           ]        stderr, existing exit codes
                      ◀ tests attach here: a table-driven probe sweeps the subcommand set
                        outside a repo and asserts every emission contains the one phrase
                        and honors its surface's channel + exit contract.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Every command's not-in-repo emission contains the one rendered phrase on its contract channel | built binary outside a repo (contract) | table-driven sweep asserting the phrase per command — red today for `gate` (`gate: not in a git repo`) and the other bare-phrase operational commands | any surviving second phrasing fails its table row; the sweep is the enforcement the "uniform posture" comment never had |
| 1 | AXI commands keep stdout/exit-1; operational commands keep their existing exit codes | built binary outside a repo (contract) | already covered for AXI (fail-closed suite pins them); the operational rows extend the same probe and are red only where the phrase changes channel | guards against "unify the phrase" accidentally becoming "unify the posture" — the exit codes must not move |
| 2 | `craft-tdd` carries the schema by pointer; the five-field enumeration appears only in the write-spec command | conformance docs anchor | anchor asserting the command file carries the schema and craft-tdd carries the pointer phrase — red today (both carry the enumeration) | if a later edit re-inlines the schema into the skill, the anchor goes red instead of the copies silently diverging |
| 3 | `bin/bench.sh` header no longer carries a subcommand roster | reviewer cold-read | not TDD-able — a comment has no behavior surface; enforced at review (the shellcheck phase only lints, and pinning comment absence would over-fit the gate) | stated openly: this is a review-checked deletion, not TDD coverage |
| 4 | The profile seam bullet points at the BENCH.md inventory without enumerating subcommands | reviewer cold-read | not TDD-able — profile prose outside the anchor families; review-enforced | same honest disclosure as story 3 |
| 5 | The tier-binding paragraph names the conformance check as the coupling | reviewer cold-read | not TDD-able — one profile sentence; the coupling itself is already gate-tested by checkLineBinding | the enforcement exists; only its advertisement changes, so review is the right venue |

### Edge inventory

- error path → covered: the not-in-repo state *is* the error path; every
  command's row asserts it.
- empty/absent input → **Won't handle** beyond the sweep: commands' per-arg
  usage postures are unchanged and already pinned by their own suites.
- boundary values → covered: the sweep includes both surface classes (AXI +
  operational) so the two-channel contract is exercised at its boundary.
- malformed input → **Won't handle**: no parsing changes in this spec.
- interrupted/partial state, re-run idempotency → **Won't handle**: message
  rendering has no state.
- hostile environment (cwd deleted under the process, `GIT_DIR` pointing at
  garbage) → **Won't handle**: `git.Root()`'s failure detection is upstream of
  this change and unchanged; the sweep only asserts what renders after it
  fails.

## Out of scope

- **A model-invocable spec-authoring skill** (`craft-spec`) owning the
  coverage-map discipline — a separate capability the assessment's §5 low
  finding motivates (the structural root of the duplication); parked as a
  ROADMAP row for a `craft-synthesis`-disciplined build. Estimate: ~6 edits,
  2 gate runs.
- **Unifying operational-command exit codes** for the not-in-repo state — a
  behavior contract change across the operational surface, distinct from the
  phrasing collapse. Estimate: ~5 edits, 3 gate runs.
