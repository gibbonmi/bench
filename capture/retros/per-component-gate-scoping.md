# Retro: per-component-gate-scoping

## Outcome

Promoted at `f341800` (candidate `2e1a61c`, run `a74f334f…`), spec marked
implemented by promotion and retired at `627be9c`. 17 planned tickets plus two
review-round repair tickets, all ported or authored by write-delegates and
landed through the full `bench spec build` lifecycle — the first complete
assign→checkpoint→integrate→review→promote run in this repo. All 21 acceptance
coverage rows green. Five findings flagged for reviewer veto (S1 reduced-run
reachability comment, S2 provenance comments, S3 dead ReusableGreen branch,
Sp1 derivation-source check does not bind an entry to its named derivation,
Sp5 stale superseded repair ticket); one accepted residual recorded in the
profile (field-set slice drift is unobservable to the record-class registry).

## Gate-stage timings

Full `--fresh` runs during the build: ~10–14 min wall each, paid four times to
re-green recomposition after capture-only commits (reduced verdicts are
deliberately non-reusable as whole-tree green — the exact cost this spec
removes). Reduced runs for capture-only commits: under 1 min. Promotion's own
composed gate: green in one run. First post-promotion gate (retire commit):
3m28s wall — the run that authored the initial component slots; later
capture-only runs inherit per-component evidence.

## Ticket-versus-spec-slice and delegate performance

Ticket-sized charges with an exact reference state ("port `git show <sha>:<path>`,
byte-verify, adapt only if compilation demands") outperformed every alternative:
12 of 17 ports came back byte-identical with zero adaptation, each in one
iteration. The two synthesis points both came from charge defects, not delegate
error: the identity charge mis-stated the branch's commit order (the delegate
detected it and synthesized the correct intermediate state), and the round-1
repair charge pinned a property that was already covered (caught by the round-2
composed review). Sonnet handled every mechanical port cleanly; opus earned its
tier only on the decision-function, identity, slot, attest, and carry tickets,
matching the spec's routing note.

## Coordinator catches

- Delegate probes and coordinator probes must differ in kind and in *site*:
  three coordinator swap probes were vacuous on the first try (gitignored
  `dist/bench` in a snapshot-based fixture, a wrong `-run` package, a mutation
  the check's charter doesn't cover) and needed re-aiming before they bit.
  A vacuous probe looks identical to a passing one — only the deliberate red
  distinguishes them.
- Receipt rows must be re-derived from the ticket file at checkpoint time, not
  from the coverage map: two checkpoints refused because tickets carried more
  bracketed rows (PS31/32, PS34/35, PS36) than the map suggested.
- The round-2 composed review overturned a repair the round-1 review had
  requested and the repair delegate had "completed" — fresh-context review of
  the exact composition catches what the authoring context cannot.

## Agent-experience improvements

### Bench CLI

- `bench spec build checkpoint` receipts are hand-assembled JSON against an
  undocumented schema (row outcome vocabulary `passed|already-covered|not-tdd-able`
  was discovered by reading `receiptRows`). A `bench spec build receipt
  <assignment> --probe <cmd>`-style generator would remove the whole class of
  invalid-receipt refusals, including the row-set mismatches.
- Recomposition after a capture-only tip move demands full whole-tree green,
  which before this build meant a 10–14 min `--fresh` run per interleaved
  commit; with per-component scoping live this cost should collapse — verify on
  the next mixed-session day.
- `error: invalid spec build receipt` names no failing condition; the diagnosis
  needed an in-package harness. Naming the first failed check would have saved
  a debug cycle.

### Skills

- craft-tickets' acceptance-row template does not match `resolveTicket`'s
  parser (bracketed row ids, single-line comma-separated backticked fences);
  the 17 staged tickets needed mechanical normalization before assignment. One
  side must own the format (journaled 2026-08-02).
- craft-delegate could name the "same fence, different site" rule for probe
  pairs explicitly; the charge template's probe-kind vocabulary (omission/swap)
  was carried by the coordinator, not the skill.

### Process

- Two sessions sharing one checkout serialize on every lifecycle mutation
  (clean-checkout preconditions) and every landing (tip moves wedge or
  recompose the run). The zero-checkpoint wedge cost three abandon/restart
  cycles before the first checkpoint existed. Rule that held afterward:
  never commit to main while a spec-build run has un-checkpointed assignments,
  and get one checkpoint in as early as possible.
- Delegate scratchpads collided: a stale `mine/` backup directory from one
  delegate clobbered four out-of-fence files in a later delegate's restore
  glob (caught by its own `git status` check). Backup paths inside a delegate
  worktree should be worktree-local and uniquely named, never shared scratchpad
  globs.
