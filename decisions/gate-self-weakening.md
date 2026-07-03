# Gate self-weakening tripwire

Closes the HIGH-sev quality-assessment finding: both `run_gate` (bin/bench.sh:69)
and the Stop hook (.bench/hooks/stop.sh:79) honor the **working-tree** `.bench/gate.sh`,
so an agent that weakens that file commits red work green. Benchkit defends via its
`tests/canary/` harness (weaken a check → its fixture stops biting → gate red), but
`bench init` scaffolds a canary-less gate and the runner silently skips when
`tests/canary/` is absent — so consumer repos have no automated defense, only the
human reviewer. This is the one finding that touches invariant #1 directly.

The build is three deliverables: the **init tripwire** (mechanism, tickets #2–#5),
an **ADR** (decided-state tradeoff, so future scans stop re-flagging the seam; the
repo's first — `docs/adr/` does not exist yet — written path-free per craft-adr, and
paired with a one-line decided-state comment at the seam so a code-reading scan finds
the marker inline), and a **craft-gate note** (guidance for the agent authoring or
editing a gate). All ticket decisions below are resolved; the map is ready for a spec.

## #1: What threat are we defending against?

Blocked by: —
Type: Grill

### Question
Working-tree honoring means the gate file is editable by the agent it grades. Do we
defend against a lazy agent taking a shortcut, or a determined adversary?

### Answer
The **accidental/lazy** agent — keep it minimal. Nothing in the working tree can
stop a determined adversary: it deletes the tripwire in the same edit that weakens
the gate. What defense buys in a consumer repo is making a weakening **loud and
visible** — a red gate plus a reviewable diff — which is the posture benchkit's
canaries already provide (they surface deletion, they don't prevent it). The
adversarial case (gate pinned outside the writable tree, hash-verified in pre-push)
is explicitly **out of scope** here and would be a separate map.

## #2: What is the tripwire mechanism?

Blocked by: #1
Type: Grill

### Question
What does `bench init` plant, and what does it assert?

### Answer
**One real seed canary**, planted by `bench init`. The scaffolded gate ships the
canary **harness active** plus one seed fixture that bites a scaffolded **example
check** (see #5 for the exact check and why it is one of *two* scaffolded checks).
This is the only option that defends from day one without pretending a universal
cross-stack check exists: the scaffolded gate starts with no real stack checks, so
the seed canary bites a scaffolded example check to prove the harness actually
executes, and the harness itself is self-defending (see #4). As the consumer adds
real checks, the seed fixture is a copy-paste worked example for canary-per-check.
Rejected: a "gate names its own checks" heuristic (vacuously satisfiable) and a
dormant harness (defends nothing until someone opts in, which the lazy agent won't).

## #3: Where does the canary-runner live?

Blocked by: #2
Type: Grill

### Question
The runner is ~30 lines, currently inline in benchkit's own `.bench/gate.sh`. If
the init scaffold pastes it into the consumer stub, that is a second copy — a
one-source-per-fact violation. How is it placed?

### Answer
Extract the runner into a shipped lib, `.bench/lib/canary-run.sh`. Benchkit's own
gate sources it (replacing its inline block); the scaffolded consumer gate sources
it too, since `bench link` already copies `.bench/lib/` into the consumer
(bin/bench-link.sh:238). **One source, both consume** — and benchkit's own canaries
then guard the shared runner. This strictly improves current state: the inline block
becomes a sourced, testable lib. The extraction is a clean refactor with its own
before/after.

## #4: Absent-harness behavior?

Blocked by: #3
Type: Grill

### Question
Today the runner **skips** when `tests/canary/` is absent. Should the shared lib go
red instead?

### Answer
**err-if-absent, in the shared lib, no opt-out.** Absence or an empty harness → red
gate. That *is* the tripwire: the lazy escape (`rm -rf tests/canary/`) becomes a red
gate and a visible diff instead of a silent pass. Benchkit is unaffected (it always
has canaries) and behaves identically to consumers. There is no deliberate opt-out —
the only real cost of the harness is gate runtime, answered by keeping the seed
harness minimal (one fixture), not by allowing it to be turned off. Someone can
still edit the err-if-absent line away, but that is the residual inherent to a
working-tree-honored gate, covered by the accepted threat model in #1 and recorded
in the ADR.

## #5: How does the scaffold avoid the fail-until-configured / vacuous-baseline deadlock?

Blocked by: #2, #4
Type: Grill

### Question
A seed canary needs a scaffolded check to bite, but two existing mechanics
constrain it. (1) The current scaffold `exit 3`s until edited (bin/bench.sh:267) —
a real safety property: an unconfigured gate is red, so a consumer cannot commit
work against an empty gate. (2) The runner's attribution baseline (gate.sh:216-228)
runs the gate against an empty `git init` repo and errs if any fixture's EXPECT
substring appears in that empty output. A single "configure me" check that the seed
canary bites satisfies (1) but fires on the empty baseline too, so its EXPECT is
flagged vacuous — a fresh `bench init` gate goes red for the wrong reason. How is
this resolved?

### Answer
The scaffold ships **two distinct checks**, not one:

- An **un-canaried "configure me" sentinel** — fires unconditionally, so a fresh
  gate is red until the consumer configures (preserves the fail-until-configured
  property). It is deliberately not canaried: it exists to be deleted, so its rot
  does not matter.
- A **canaried forbidden-token example check** — greps tracked files for a planted
  marker (e.g. `DO-NOT-COMMIT`) and `err`s if found. It passes on a normal repo and
  on the empty baseline (neither contains the marker), so its seed canary is not
  vacuous; it bites only its own fixture (one file containing the marker). The
  "run a command, `err` on failure" shape is the copy-paste template a consumer
  follows for real checks.

Why two: the two jobs are contradictory postures. "Keep a fresh gate red" must fire
everywhere (including the empty baseline); "prove the harness is live" must fire only
on a purpose-built fixture. One check cannot do both without reintroducing the
vacuous deadlock or dropping the safety property.

**Transition** (prescribed in the scaffold's own comments): (1) delete the sentinel
line once a first real check exists — that deletion *is* "configured"; (2) treat the
forbidden-token check + its seed canary as a living template — for each real check
added, add a canary that bites it; (3) the example check may be replaced by a real
one, but deleting the example requires deleting its canary in the same edit (an
orphaned canary references a gone check → "canary did not bite" → red), and not all
canaries can be deleted (err-if-absent, #4). Net: the harness always has at least one
biting fixture, and adding-a-check-adds-a-canary is the path of least resistance.

**Noted residual** (recorded in the ADR, not a defect): after the consumer deletes
the sentinel but before adding real checks, the gate is green on the example check
alone. That is the consumer's deliberate, visible configure step — covered by the
lazy-not-adversarial threat model in #1 — not a hole to engineer against.
