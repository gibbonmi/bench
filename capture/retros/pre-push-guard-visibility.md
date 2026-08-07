# Retro — pre-push-guard-visibility

## Outcome

Promoted terminal: run `9975d21b`, candidate `a57e2334`, promotion commit
`e1c24b83` (tree `f4fd683b`, retained evidence
`v1:4cc9bdc64527f0974df789c09084cb8612b5648d24cc2eb162596d56d072ea73`),
`Status: implemented` authored by promote. Nineteen assignments integrated and
released across the original thirteen tickets plus six review-driven repair
tickets. Three fresh composed review rounds; every accepted finding closed with
citations before promotion. The final round also fixed a spec-authored defect:
the clean-skip predicate keyed on the old manifest hash would have made
`bench upgrade` a content no-op; the reviewer amended the predicate to
prospective bytes-and-mode match and the amendment landed in the closing
commit after promotion (the run pins staged spec content).

## Gate-stage timings

Observed on the full-gate `bench commit` runs this session (five runs, stable
across them): test ≈228s dominated by `internal/gate` (227s), specbuild
(145–165s), contract/runtime (136s), posture (109–112s), worktree (84–86s),
offline/prepared (75–81s); conformance-suite ≈75–77s (package-core-guard 16–19s,
line-routing 7s); race ≈8–9s; canary, build, gofmt, vet, shellcheck marginal.
Wall clock per gate ≈10 minutes. Checkpoints and integrations paid zero gates;
promote paid one.

## Ticket-versus-spec-slice and delegate performance

All charges were ticket-sized with two-file ownership fences; no delegate ever
received a spec slice. Six write delegates (four opus, two sonnet) all returned
in-fence, uncommitted, with TDD-ordered evidence and their own mutation probes;
none required a repair re-charge for delegate error. The one blocked-adjacent
event was handled correctly: the clean-skip delegate flagged the upgrade
propagation question as a spec-level call instead of absorbing it, and the
fresh review round converted it into the accepted F1 finding. Fresh-context
review delegates (opus) out-performed expectations: round two demonstrated F1
by compiling and running both trees rather than arguing from the diff.

## Coordinator catches

- The delegate claim that `internal/guards` no longer reads `PrePushMarker`
  was wrong for the main tree but right for the candidate composition — caught
  by probing the candidate rather than trusting the enumeration.
- A coordinator probe sed pattern that failed to match produced a meaningless
  green (`doctorEnv` vs `DoctorEnv`, then a NUL-byte escaping mismatch on
  restore); both were caught by checking that the mutation was actually present
  before trusting the red/green, which is the probe discipline working.
- Two stale-binary refusals: `dist/bench` predating the FT135 fix made the
  repaired recomposition appear still broken, and the post-promotion binary
  made green surface tests fail; both diagnosed by timestamp-versus-commit
  comparison before touching any state.
- Checkpoint receipt refused on `outcome: "green"` — the schema wants
  `passed`; diagnosed with a throwaway in-package test rather than guessing.

## Repair attribution

| ticket | repair rounds | causes |
|---|---|---|
| narrow-pre-push-readme-claim | 0 | none |
| repair-upgrade-hook-planning | 0 | none |
| render-guards-hook-health | 0 | none |
| expose-hook-health-record | 0 | none |
| render-doctor-hook-health | 0 | none |
| repair-stale-hook-with-doctor | 0 | none |
| converge-symlinked-hook-lifecycle | 0 | none |
| align-upgrade-hook-plan | 0 | none |
| contract-hook-health-entry-point | 0 | none |
| restrict-upgrade-hook-absence | 0 | none |
| signal-stale-hook-status | 0 | none |
| converge-live-hook-installation | 0 | none |
| render-hook-refresh-plan | 0 | none |
| repair-recompose-bootstrap-marker | 0 | none |
| restore-doctor-hook-execute-mode | 0 | none |
| correct-hook-marker-comment | 0 | none |
| repair-link-symlink-conflict-order | 1 | spec-row |
| repair-hook-symlink-fifo-gate | 0 | none |
| repair-doctor-surface-defects | 0 | none |
| repair-clean-skip-propagation | 0 | none |

The one repair round: the symlink-conflict-order ticket implemented the spec's
clean-skip row faithfully and the row itself was defective (old-manifest-hash
predicate); the follow-up ticket repaired the behavior under the amended row.

## Agent-experience improvements

### Bench CLI

- Lifecycle bookkeeping (ticket files, and any reviewer-approved spec
  amendment) pays a full gate per repair round because it must land on main
  before `assign`; a staged-ticket route or docs-scoped commit tier would make
  repair rounds gate-free until promote. Captured in `capture/learnings.md`.
- The checkpoint receipt's `outcome` vocabulary (`passed`,
  `already-covered`, `not-tdd-able`) is discoverable only by reading
  `receiptRows`; a one-line refusal message naming the accepted values would
  have saved a diagnostic round trip.
- `bench spec build status --full` does not print the run id; it had to be
  read from the state file to author receipts.

### Skills

- The craft-delegate charge template's mutation-probe requirement carried its
  weight: every delegate red was demonstrated, and the two coordinator probe
  mishaps were caught by the verify-the-mutation-exists discipline the skill
  prescribes.

### Process

- Fresh-context review delegates that run the code (not just read the diff)
  caught the one defect all prior layers missed; keep charging reviews with
  permission to build and execute the candidate.
- Parking a reviewer-approved spec amendment until after terminal promotion is
  the working pattern under the current SpecTip pin; the closing commit is its
  landing slot.
