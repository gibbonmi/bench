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

## 2026-08-01 — light-path ticket receipts have no terminal lifecycle  [open]

- **What happened:** Five shipped light-path changes still have ticket-only
  folders under `specs/`, but no `spec.md`. Their implementation commits landed,
  while `bench spec history` returns no record for their slugs and `bench spec
  retire` cannot target them. The folders therefore persist as committed receipts
  without being classified as active specs, retained history, or retireable
  lifecycle state.
- **Right behavior:** A light-path close gives its ticket receipt one explicit
  terminal disposition. If receipts remain under `specs/`, a canonical reader
  identifies why they remain and distinguishes them from active specs; otherwise
  the close removes them after promoting any durable content. Sessions should not
  have to infer lifecycle state from the absence of `spec.md`.
- **Proposed rule change:** Make the light-path final-check path either retain and
  index ticket-only receipts through one existing spec reader, or remove their
  folders at close. Choose one owner and one policy rather than adding a second
  ad-hoc archive convention.

## 2026-08-01 — a write-delegate's charge should name the mutation it must run against its own work  [open]

- **What happened:** Across the per-component-gate-scoping build, three delegates
  returned rows they believed green whose mutations did not actually red: two because
  the kit-shaped fixture made the right and wrong answers identical (a literal build-input
  list happened to equal the derivation's answer; a corrupted `go.mod` short-circuited at
  the first component so later derivations were never asked), and one because the ticket's
  literal operation sequence could not isolate the property (a slot copied to another
  component's path was refused by the identity comparison, so dropping the component
  comparison changed nothing). Each of those three was caught by the delegate itself,
  because the charge told it to apply the mutation and observe the red rather than to
  reason that the code was correct. The one that was NOT caught had no such instruction:
  the component-slot class's central property — that a slot stops answering once its
  component's identity moves — was pinned by nothing, and a coordinator probe that made
  slots non-content-addressed left the entire package green. The rows I wrote covered
  byte-identity, four refusal classes, and per-component isolation, and still missed the
  one property the class exists for.
- **Right behavior:** Two distinct probes, and the charge must carry both. First, the
  delegate's SELF-PROBE: the charge names a specific mutation — the one that breaks the
  feature's central property, not a peripheral one — and requires the delegate to apply it
  to its own finished work, report the observed result, and add the missing row if the
  package still passes. Naming it in the charge is what makes it a check rather than a
  reasoning exercise; "test it thoroughly" produces the same blind spot the rows already
  have. Second, the coordinator's INDEPENDENT probe, which must be a different KIND from
  any the delegate ran — a registry mis-binding, a transposition, a widening, a
  determinism break — because a second instance of the same mutation kind is correlated
  evidence, not independent evidence. `craft-delegate` already requires the coordinator
  probe; what is missing is the self-probe on the charge side, and the explicit
  different-kind requirement. The generalization: rows are written from what the author
  believes the feature does, so they inherit the author's blind spot; a named mutation
  attacks the property instead of the belief.
- **Proposed rule change:** Fold into **FT164**, which owns ticket and repair charge
  quality — the self-probe is a charge-side duty and belongs beside the enumerated-family
  and blast-radius clauses already queued there. The coordinator half (independent probe
  must differ in kind from the delegate's) sharpens the existing "probe at least one
  accepted behavior independently" line in
  `.agents/skills/bench-craft-delegate/SKILL.md`, which today says the probe must be
  independent of the delegate's tests but not that it must differ in kind. Kit edit under
  `craft-synthesis`.

- 2026-08-02  Two format gaps between authored tickets and the specbuild parser: acceptance
  rows must read `- [ ] [ID] text` (bracketed id) and `Ownership fence:` must be one line of
  comma-separated backticked paths — the craft-tickets template and the pcgs tickets as staged
  used `ID —` rows and wrapped prose fences, which parse to zero rows / garbage fences and
  refuse assignment. I normalized the 17 pcgs tickets mechanically (landed 3972744). Proposed
  rule: either teach craft-tickets the parseable shape or make resolveTicket accept the
  template's shape; one of the two must own the format.
- 2026-08-02  Committing to main while a spec-build run holds un-checkpointed assignments
  wedges the run: the tip moves, every operation demands recomposition, and recomposePromotion
  cannot replay an empty candidate ("No valid patches in input"), so checkpoint, start, and
  abandon all refuse — abandon's own recomposition exemption is FT176's
  exempt-abandon-from-the-recomposition-refusal. Right behavior: sequence capture/main commits
  strictly outside run windows (before start or after integrate), and empty-candidate
  recomposition should fast-forward the base rather than replay a patch.
