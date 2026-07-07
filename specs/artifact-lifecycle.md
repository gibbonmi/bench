# artifact-lifecycle — own the backward path (spec retirement + review pickup)

Status: staged

Source: `ASSESSMENT.md` backlog 1 + 11 (findings §1 high/med, §3 med/low).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

The pipeline's forward path (idea → map → spec → build → review → gate → commit)
is well-owned; the backward path depends on someone remembering. Live instances
found by the 2026-07-06 assessment: a merged spec awaiting a manual
promote-then-delete that only prose describes; `reviews/<slug>.md` whose deletion
is promised by the review phase but never read, resolved, or deleted by
`/bench-implement-spec`; a review pickup that orphans when the reviewer accepts
residual risk or when its spec retires first; and a review artifact born
untracked, invisible to git-based checks and flipping the gate verdict stale.

## Solution

Give the backward path owners. A new `bench spec retire <slug>` performs the
mechanical half of retirement (validated deletion of the spec and its review
pickup); the phase docs close the review-pickup loop (implement-spec reads and
deletes it; review-implementation commits it at birth; final-check deletes it on
an accepted-residual-risk route); and `bench status` flags an orphaned review
pickup so no artifact state is unowned.

## User stories

1. As a reviewer, I want `bench spec retire <slug>` to delete `specs/<slug>.md`
   and, when present, `reviews/<slug>.md`, only after validating the spec is
   merged-implemented (the same detector the status retirement signal uses), and
   to print what it removed plus the judgment duties that remain (promote durable
   content, remove the ROADMAP row, commit as `spec-retire: <slug>`), so the
   mechanical half of promote-then-delete is owned by the CLI instead of prose
   memory.
   Line: claude-opus-4-8 / medium. This is a new operational subcommand whose
   refusal postures guard a destructive deletion, so it routes to the mid tier
   even though the seam itself is known.

2. As a reviewer, I want `bench spec retire` to refuse loudly on every unsafe
   input — a staged (unmerged) spec, an unknown slug, a slug whose spec is
   missing but whose review pickup exists, extra or unknown arguments — so a typo
   or a premature retire is non-destructive.
   Line: claude-sonnet-5 / medium. This is edge behavior at the same
   contract-observable seam as story 1, so the cheap tier at medium effort
   covers it.

3. As a build agent, I want `/bench-implement-spec` to say that when
   `reviews/<spec-slug>.md` exists for the spec being built, its findings are
   part of the build target and the green fix commit that closes them names and
   deletes the file, so the review phase's deletion promise has an owner on the
   implement side.
   Line: claude-fable-5 / high. The command file is a leverage artifact loaded
   by every implement-spec session, so the profile's doc-authoring override
   routes it to the top tier.

4. As a reviewer, I want `/bench-review-implementation` step 5 to commit
   `reviews/<spec-slug>.md` in the same session that writes it, and
   `/bench-final-check` to delete the pickup file (naming it in the landing
   commit) when I accept residual risk and skip the fix pass, so the artifact is
   never untracked drift and never outlives the decision it captured.
   Line: claude-fable-5 / high. Same leverage class — two workflow command files
   whose wording every future session inherits.

5. As a cold session, I want `bench status` to flag a `reviews/<slug>.md` that
   has no matching `specs/<slug>.md` as an orphaned review pickup with a
   clean-up action, so a pickup file that escaped its lifecycle is surfaced
   instead of silently reading as pending work.
   Line: claude-opus-4-8 / medium. The ambient dashboard's severity ladder and
   row budget are a gate-tested contract, so a new signal routes to the mid
   tier.

## Implementation decisions

- **`bench spec retire` is a sibling of `bench spec implemented`** in the spec
  package: same slug/path resolution, same stdout error style (`toon.Usage` exit
  2, `toon.Errorf` exit 1), plain success lines. It deletes files and prints the
  remaining duties; it never commits and never runs the gate — `bench commit`
  owns commit discipline.
- **Merged-implemented validation reuses the status detector.** The definition of
  "merged spec awaiting retirement" (spec content at HEAD carries the
  `Status: implemented` marker) already lives behind the status specs row; retire
  calls that one source rather than growing a second parser. (Default call,
  flagged: retire refuses on a spec that is implemented in the working tree but
  not yet at HEAD — that state means the finishing commit hasn't landed.)
- **Deletion order: review pickup first, then spec.** An interrupt between the
  two leaves a valid spec that a re-run retires cleanly, never an orphaned
  review file.
- **Orphaned pickup is refused, not auto-cleaned.** `bench spec retire <slug>`
  with no spec file errors and names the orphan and the manual fix; deleting a
  review file without its spec is a reviewer judgment the status row (story 5)
  surfaces. (Default call, flagged for veto.)
- **The review artifact becomes tracked state at birth** (story 4). Committing
  it keeps `bench tree-hash` semantics untouched — the tripwire still counts
  untracked files, and the assessment's tree-hash question resolves as "commit
  the artifact", not "exempt transient state". (Default call, flagged: the
  alternative — excluding `reviews/` from tree-hash — weakens the stale
  tripwire for every future untracked file class.)
- **Doc anchors follow the existing pattern**: new load-bearing phrases in the
  three edited command files are pinned by tightening
  `internal/conformance/docs_workflow_helpers_test.go` anchors, the same move
  the ft9 spec used for its call-site repoint.
- **Status row severity**: the orphan signal ranks with the housekeeping rows
  (near the existing specs-awaiting-retirement row), inside the five-row budget;
  it must not displace gate/git rows.

## Testing decisions

- **What a good test is here:** drive the built `bench` binary in a throwaway
  git fixture and assert stdout/stderr and exit codes — never internals. Prior
  art: `internal/contract/runtime/runtime_spec_test.go` (spec flip postures),
  `runtime_status_test.go` (dashboard rows), `runtime_commit_test.go` (named-set
  staging).
- **Seams:** the `bench spec` command surface (runtime contract), the
  `bench status` renderer (runtime contract), and the conformance docs-anchor
  layer for the three command-file edits.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

Seam 1 — the retire operation:

    trigger: reviewer (or a session doing the retire pass) runs `bench spec retire <slug>`
        │
        ▼
    slug/path arg  ──▶  [ spec.Command "retire"                 ]  ──▶  stdout: retired: specs/<slug>.md
    tree at HEAD   ──▶  [  resolve → validate merged-implemented ]         retired: reviews/<slug>.md (when present)
    reviews/<slug> ──▶  [  → delete pickup, delete spec          ]         next: promote / ROADMAP row / spec-retire commit
                        [                                        ]  ──▶  refusals: toon errors, exit 1/2
                      ◀ tests attach here: runtime contract drives the built binary in a
                        fixture repo (staged spec, implemented spec, orphan pickup) and
                        asserts output lines, exit codes, and which files remain on disk.

Seam 2 — the orphan signal:

    trigger: SessionStart hook or reviewer runs `bench status`
        │
        ▼
    reviews/*.md   ──▶  [ status renderer                   ]  ──▶  ranked dashboard rows,
    specs/*.md     ──▶  [  pair pickups with spec files      ]        incl. "orphaned review pickup"
                        [  rank on severity ladder           ]
                      ◀ tests attach here: runtime contract builds a fixture with
                        reviews/x.md and no specs/x.md and asserts the row text
                        and its action; a paired fixture asserts no row fires.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench spec retire <slug>` on a merged-implemented spec deletes the spec and its review pickup and prints the remaining duties at exit 0 | spec command runtime contract | `bench spec retire x` in a fixture with merged specs/x.md — today exits 2 `unknown argument: retire`, so RequireExit(0) and the `retired:` assertions fail | a retire that doesn't exist, deletes the wrong files, or skips the pickup fails the output and on-disk assertions |
| 1 | retire with no review pickup present retires the spec alone and says so | spec command runtime contract | same fixture without reviews/x.md — exits 2 today | an implementation that errors on the missing pickup, or silently claims to have deleted one, fails the row |
| 2 | retire on a staged spec refuses at exit 1 without deleting anything | spec command runtime contract | fixture with staged specs/x.md — `bench spec retire x` must exit 1 and leave the file; exits 2 (unknown arg) today, and the refusal text does not exist | a retire that deletes unmerged work is the destructive failure this row pins |
| 2 | retire on an unknown slug exits 1 naming the tried paths; extra/unknown args exit 2 | spec command runtime contract | `bench spec retire nope` and `bench spec retire x y` — both exit 2 (unknown arg) today with the wrong message | wrong exit codes or a silent success on a typo would fail the exact posture assertions |
| 2 | retire on an orphaned pickup (review file, no spec) refuses and names the orphan | spec command runtime contract | fixture with reviews/x.md only — no such refusal text exists today | an auto-clean or a generic not-found would fail the named-orphan assertion |
| 3 | `/bench-implement-spec` names the review-pickup duty (read findings, delete in the green fix commit) | conformance docs anchor | tighten `docs_workflow_helpers_test.go` to require a `reviews/` anchor in bench-implement-spec.md — red today (`rg reviews/` hits only the review command) | if the duty is dropped or reworded away, the anchor goes red instead of the lifecycle silently reopening |
| 4 | `/bench-review-implementation` step 5 commits the pickup at birth; `/bench-final-check` deletes it on accepted residual risk | conformance docs anchor | anchors for the commit-at-birth phrase and the final-check deletion phrase — both absent today | the two docs are the only owners of these duties; losing the phrases loses the owner |
| 5 | `bench status` flags an orphaned review pickup with a clean-up action | status runtime contract | fixture with reviews/x.md and no specs/x.md — no such row renders today | if the signal is missing or mispaired, the row assertion fails; a paired-fixture probe guards against false fires |

### Edge inventory

Walked per behavior; each resolved as a coverage row above or a **Won't handle**
line here.

- error path → rows: staged refusal, unknown slug, orphaned pickup, unknown args.
- empty/absent input → rows: missing pickup (retire proceeds), missing spec
  (refuse); `bench spec retire` with no slug exits 2 (usage) — covered with the
  unknown-arg row.
- interrupted/partial state → decision recorded: pickup deleted before spec, so
  an interrupt leaves a re-runnable state; the re-run row below covers recovery.
- re-run idempotency → covered by the unknown-slug row: a second retire finds no
  spec, exits 1, deletes nothing.
- malformed input (slug with spaces, path arg) → resolution rides
  `spec.Resolve`, already exercised by the flip contract tests; asserted once in
  the story-1 fixture by retiring via a path argument.
- hostile environment (cwd deeper than root) → slug fallback anchors at the repo
  root exactly as `bench spec implemented` does — same code path, already
  contract-tested for flip; **Won't handle** a separate retire probe.
- control bytes in slugs — **Won't handle**: slugs are repo file basenames the
  reviewer created; the profile's control-byte class applies to git-sourced
  text rendered through TOON tables, and retire renders none.
- symlinked invocation / shipped surfaces — **Won't handle**: retire adds no new
  dispatch; routing is covered by the existing surface contract.

## Out of scope

- **`bench spec history <slug>`** (fold the duplicated
  `git log --grep=spec-retire` recovery incantation from two command docs into
  the CLI) — a separate read-only query capability with its own AXI contract;
  parked as a ROADMAP row. Estimate: ~4 edits, 2 gate runs.
- **Auto-removing the retired spec's ROADMAP row from the CLI** — row wording is
  reviewer judgment reconciled by `/bench-what-next`; mechanizing markdown row
  surgery is a distinct (and riskier) capability. Estimate if ever wanted:
  ~6 edits, 3 gate runs.
