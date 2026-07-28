# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-07-28 — a write-charge named the family binding but not the registration seam  [open]

- **What happened:** The FT152 build charged a delegate to add four canary fixtures
  and named the family→scope binding in `internal/conformance/registry/registry.go`
  as the gate layer. It did not name `canaryFixtureRegistry` in
  `registry_test.go`, which requires a per-fixture classification. The delegate's
  work was correct against its charge; the full gate went red on four
  `is unclassified` errors. The repair cost one extra delegate round.
- **Right behavior:** A charge that adds a member to a registered family names
  every registration seam that member must appear in, not just the family's own
  binding. Charge-time is when the coordinator can still enumerate them cheaply —
  `rg` for one existing sibling's name across the package finds all of them.
- **Proposed rule change:** In `craft-delegate`'s "The charge", extend the
  gate-layer clause: when the delegate adds a member to an enumerated or
  registered family, the charge names each registry it must be entered in, found
  by grepping an existing sibling rather than recalled.

## 2026-07-28 — a universal claim shipped from a truncated read, inside the diff that forbids it  [open]

- **What happened:** While repairing FT152 I told a delegate "every existing
  `workflow-guidance-anchors` fixture is listed there", having read only the
  visible head of a long Go slice. The delegate checked and found
  `terminal-repair-bound-anchor` registered but absent from every bite list. The
  false claim was in the same session that shipped story 3's quantifier clause —
  a claim over a whole set is verified by enumerating the set, not by extending
  one measured member.
- **Right behavior:** Enumerate before quantifying, or state the claim as a
  sample. The cost here was zero because the delegate verified rather than
  complied; that is the only reason it was caught.
- **Proposed rule change:** None. Story 3's hook now sits in
  `/bench-implement-spec` at the point of use, which is exactly the repair this
  incident argues for. Recording it as evidence the hook earns its place, not as
  a new rule.

## 2026-07-28 — the falsification pass caught what three fresh-context axes cleared, again  [open]

- **What happened:** On FT152's build, a three-axis semantic review returned 29
  findings and cleared the `--full` section's scope fence. A cross-harness Codex
  falsification pass at the top binding, charged to refute rather than grade,
  found that the scope fence re-derived the fix-don't-park boundary as
  path-width-based where `.bench/BENCH.md` sets it decision-based — so the two
  rules give opposite answers for a small out-of-story bug. This is the second
  such result after FT91's draft, where the same pairing returned block after
  this repo's own axes had cleared it.
- **Right behavior:** Unclear, and that is why this is captured rather than
  acted on. FT152 ships the pass as an offer gated behind a judgment-sized
  trigger; two-for-two suggests the trigger may be set too high, but two runs is
  not a measurement.
- **Proposed rule change:** None yet — a third data point, or a run where the
  pass returns nothing, should decide it. Candidate if the pattern holds: make
  the falsification pass standing for kit-guidance diffs specifically, where a
  defect compounds through every session that loads the prose.
