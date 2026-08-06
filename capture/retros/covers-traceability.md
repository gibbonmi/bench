# Retro: covers-traceability

## Outcome

Promoted terminal on the first promotion attempt after two repair rounds:
squash `5ff1505`, candidate `29b9cbb`, spec `Status: implemented`. All 15
coverage-map rows landed with named test owners. Seven assignments (five
tickets plus two repair tickets) all integrated and released through the public
lifecycle; no synthesized plumbing. Three composed review rounds: round 1 found
one concrete defect (duplicated, divergent row-ID grammar) and one concrete
coverage gap (silent first-wins on multiple covers annotations); round 2's
Coverage axis found the prose-mention false-coverage gap, closed by the
reviewer-decided bracket-anchoring repair; round 3 was clean of concrete
findings. Six judgment findings left flagged on the review receipt for reviewer
veto. Mid-build the reviewer also directed a light-path landing: the
craft-tickets thin-by-default slicing flip (`1691017`), justified by this
build's own rows-per-ticket data.

## Gate-stage timings

Not instrumented per stage this run. Observed wall-clock: reduced-scope
(specs/-only) gated commits ~1–2 min; full fresh gates ~5–8 min warm; the one
fully cold gate (immediately after a `go clean -cache` on an 18G cache) was the
session's longest single wait. `internal/specbuild` full package test ~55–61 s
per run, paid once per checkpoint receipt and once per verification round —
the dominant repeated cost. Two full gates were spent on cache-wipe timing
(reviewer-requested wipes landed just before gate-heavy phases).

## Ticket-versus-spec-slice and delegate performance

Five tickets of 3–5 acceptance rows each, one fresh write-delegate per ticket,
all first-pass green at the charged tier (2× sonnet/low producers, 2×
opus/medium lifecycle policy, 1× fable/high guidance). No ladder escalations.
The fixture-and-seam inventory pasted into each charge (from one pre-build
read-only research delegate) is the visible reason first-pass quality held:
delegates wrote against `requirePromotionRefusal`, `ticketFixture`, and
`newCheckpointFixture` prior art instead of re-deriving harnesses. Repair
tickets (sonnet/low) also landed first-pass. The one delegate quality miss:
the grammar-repair delegate skipped the ticket's RG4/RG5 red demonstrations
("green from first run"), which the coordinator had to run itself.

## Coordinator catches

- Vacuous probe detected: the `local`-guard mutation on promote totality is an
  equivalent mutant (the ID grammar forbids `local` as a map ID), caught when
  the coordinator probe stayed green; replaced with a naming-completeness swap
  that bit. Round-3 Coverage independently confirmed the equivalent-mutant
  reading.
- The prose delegate's bite-test run was vacuous: `-run 'Bite\|ExampleAgreement'`
  matches a literal `|` in RE2, so zero tests ran and printed ok; the real bite
  red surfaced only in the coordinator's full gate.
- The grammar-repair delegate's skipped per-row reds (RG4/RG5) were run by the
  coordinator before acceptance; both bit.
- Review delegates' claims spot-checked against regexes and code already read;
  round-1 Standards' grammar-divergence citation verified before ticketing.

## Agent-experience improvements

### Bench CLI

- Checkpoint and review receipts are strict JSON with exact tree hashes,
  RFC3339Nano times, and exact changed-path sets, but nothing emits a receipt
  skeleton; the coordinator reverse-engineered the schema from
  `internal/specbuild/checkpoint.go`. A `bench spec build receipt --assignment
  <id>` scaffold (or a documented schema dump) would remove the failure-prone
  step.
- A `bench commit` run inside a generic worktree lands on the worktree's own
  branch and reports success; the coordinator must know to fast-forward main.
  Either the landing should name the destination branch loudly or porcelain
  should offer a main-landing route.
- Worktree gate verdicts live in the worktree's private gitdir, so a released
  worktree's green cannot authorize recomposition in the main checkout; the
  refusal names the remedy, but the round-trip cost one extra full gate.
- `bench gate --fresh` porcelain returns while `gate-run` continues detached;
  a chained follow-up command races it. A `--wait` flag or documented blocking
  behavior would prevent the race.

### Skills

- craft-delegate could name the RE2 pitfall in its charged-checks guidance:
  quoting `-run 'A\|B'` runs zero tests and prints ok — require delegates to
  `-list` the pattern first or use unquoted alternation.
- The thin-by-default slicing flip landed mid-build from this build's own data
  (3–5-row tickets, all first-pass green) — evidence the parked threshold was
  right to wait for.

### Process

- Cache wipes and gate-heavy phases interact badly: two wipes landed
  immediately before gate-heavy stretches, each converting the next gate to
  fully cold. Sequencing rule worth adopting: wipe caches only after the last
  planned full gate of a phase.
- The mid-build light-path landing forced an extra recomposition and an extra
  fresh-evidence gate. It worked, but landing order (prose before the composed
  review binds) had to be reasoned out manually; a lifecycle note on
  recomposition cost when interleaving light-path work would make the tradeoff
  explicit.
- Stale orphaned `gate.sh` process groups from earlier sessions (wrapper-kill
  leaks) surfaced and were swept; motivation recorded in the parked
  gate-invoke-to-Go idea.
