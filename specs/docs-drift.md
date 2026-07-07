# docs-drift — the platform docs describe the current decided state

Status: staged

Source: `ASSESSMENT.md` backlog 10 + 12 (findings §6 med/low, §5 low-med/low).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

Four platform documents have drifted behind shipped decisions, and two skill
pointers misdirect their weakest readers. ADR 0001 still frames the tripwire as
chosen "rather than write protection" though push-time pin protection has since
shipped (and the ADR embeds file-path tokens invariant 3 forbids);
BENCH-reference's Hook Layers omits the shipped pin verification; CONTEXT.md
defines "skill" by its `.claude/` mirror path instead of the canonical
`.agents/skills/` source; nine hook/adapter plumbing subcommands sit in the
always-loaded BENCH.md inventory costing tokens every session; `craft-cli`'s
trigger only fires for CLIs that already declare AXI conformance — excluding
the moment the conformance decision is being made; and `craft-review` /
`craft-tdd` call the edge-inventory step "the canonical list" without naming
the domain checklist that completes it.

## Solution

One documentation pass, each edit at the fact's canonical home: rewrite ADR
0001 to the current tripwire-and-pin state; add the pin layer to Hook Layers;
repoint CONTEXT.md's skill definition; demote the plumbing subcommands to
BENCH-reference; widen `craft-cli`'s trigger to the conformance decision
point; and make the two edge-inventory pointers name both sources (generic
classes in the write-spec step, domain checklist in the profile).

## User stories

1. As a teammate who just walked in, I want ADR 0001 to describe the current
   decided state — the working-tree gate defended by the canary tripwire *and*
   push-time pin verification — with the title reframed and file-path tokens
   removed, so the record reads as what is, not the history of what changed.
   Line: claude-fable-5 / high. ADR authoring is the profile's doc-authoring
   leverage override, and `craft-adr` governs the rewrite.

2. As a session consulting Hook Layers, I want the pre-push entry to name both
   denials — direct pushes to the default branch and `.bench` drift from the
   gate pin — so the reference matches the shipped hook. (Build note: the
   guards-wiring spec makes the hook's own manifest state-aware; this line
   describes the layer, not the live pin state.)
   Line: claude-fable-5 / high. Reference prose every harness consults; same
   override.

3. As a cold session reading CONTEXT.md, I want "skill" defined by
   `.agents/skills/*/SKILL.md` with `.claude/` named as the mirror, so the
   ubiquitous language points at the canonical source.
   Line: claude-opus-4-8 / low. A one-line vocabulary correction — mechanical,
   but in the language file cold sessions trust, so mid rather than cheap.

4. As every future session, I want the nine hook/adapter plumbing subcommands
   moved from BENCH.md's always-loaded CLI Inventory to a plumbing list in
   BENCH-reference, with the inventory noting where plumbing lives, so the
   per-session token cost buys only what sessions actually invoke.
   Line: claude-fable-5 / high. BENCH.md is the highest-leverage always-loaded
   file; restructuring its inventory is exactly what the override exists for.

5. As a session about to build a new agent-facing CLI, I want `craft-cli`'s
   trigger to fire when *deciding whether* a CLI should conform to AXI — not
   only after the project already declares conformance — so the skill loads at
   the decision point it governs.
   Line: claude-fable-5 / high. Skill frontmatter steers every trigger match;
   the skill-authoring override applies.

6. As the weakest reader of `craft-review` or `craft-tdd`, I want their
   edge-inventory pointers to name both real sources — the generic edge-class
   list in `/bench-write-spec`'s edge-inventory step and the domain
   hostile-input checklist in the project profile — so the pointer teaches
   where the classes actually live.
   Line: claude-fable-5 / high. Skill prose, same override.

## Implementation decisions

- **ADR 0001 keeps its filename.** The file name is the record's identity and
  is referenced from the profile; only the title line, the state description,
  and the Considered-options framing change. The rewrite absorbs the current
  "Considered options" paragraph (which already describes the pin) into the
  decided state, per invariant 3: current state, no decision history.
- **File-path tokens leave the ADR body**; the mechanism is described by role
  (the managed pre-push hook, the recorded pin), not by path.
- **The BENCH.md demotion keeps cold-pickup guarantees**: the inventory remains
  canonical in BENCH.md for the commands sessions drive, and gains one line
  saying plumbing subcommands are enumerated in BENCH-reference — a pointer,
  not a copy. The BENCH-reference plumbing list becomes the one enumeration of
  the nine. The heading-duplication conformance check constrains the section
  naming; the build honors it.
- **`craft-cli` trigger wording**: the description adds the decision moment
  ("or when deciding whether a new agent-facing CLI should conform") while
  keeping the existing conforming-project trigger — widening, not moving.
- **Pointer fixes don't move content**: the generic edge-class list stays
  canonical in the write-spec step (it hosts and declares it); the two skills
  simply stop calling it "the" canonical list without naming the profile
  checklist half. One sentence each.
- **Anchor coverage**: BENCH.md and BENCH-reference are inside the
  docs-currency and token-diet conformance families; the build re-runs those
  checks and adds/updates anchors only where an existing family already pins
  the edited phrase — no new anchor classes for prose whose absence the
  reviewer catches (recorded honestly in the map).

## Testing decisions

- **What a good test is here:** almost none to write — this is prose whose
  correctness is semantic. The gate guards regressions structurally (workflow
  anchors, token-diet placement, heading single-sourcing, stale-reference
  sweeps, skills-index generation), and the reviewer enforces the content at
  cold-read. Rows below classify honestly rather than dressing prose review
  as TDD — the ft4 spec set the precedent for this posture.
- **Seams:** the conformance docs families (anchors, token-diet,
  skills-index) and reviewer cold-read.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: `bench gate` (conformance docs families) + reviewer cold-read
        │
        ▼
    ADR 0001, BENCH.md,        ──▶  [ docs conformance families        ]  ──▶  diagnostics on anchor loss,
    BENCH-reference, CONTEXT,  ──▶  [  anchors · token-diet · heading  ]        heading duplication,
    craft-cli / craft-review / ──▶  [  single-source · skills-index    ]        stale refs, index drift
    craft-tdd SKILL.md         ──▶  [  · stale-ref sweep               ]
                      ◀ tests attach here: run `bench gate` after each edit; the families
                        catch structural rot. Content correctness attaches at reviewer
                        cold-read — stated openly, not disguised as TDD.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | ADR 0001 describes tripwire + pin as current state, no file paths, title reframed | reviewer cold-read | not TDD-able — ADR content is semantic; the stale-reference sweep and profile reference to the ADR file stay green structurally | a wrong or history-shaped rewrite passes every check; the reviewer is the enforcement, stated openly |
| 2 | Hook Layers names both pre-push denials | reviewer cold-read | not TDD-able — reference prose; no anchor family pins this sentence today and minting one over-fits the gate | same honest posture as story 1 |
| 3 | CONTEXT.md defines skill by the `.agents/skills/` source | reviewer cold-read | not TDD-able — vocabulary prose | same |
| 4 | The nine plumbing subcommands appear once, in BENCH-reference; BENCH.md carries the pointer | conformance docs families + cold-read | the token-diet and heading-single-source checks run today and must stay green through the move — red only if the move violates them; the once-only property is review-checked | the families catch the structural failure modes of the move (duplicated headings, diet violations); the count of enumerations is reviewer-checked |
| 5 | `craft-cli`'s description fires at the conformance decision point | conformance skills-index + cold-read | the skills index is generated from frontmatter and its equality is gate-checked — the regenerated index goes red if the build edits the description without regenerating | the index check catches the mechanical half (frontmatter/index drift); trigger wording quality is reviewer-checked |
| 6 | Both skills' edge-inventory pointers name the write-spec step and the profile checklist | reviewer cold-read | not TDD-able — one sentence per skill | same honest posture as story 1 |

### Edge inventory

This spec touches no code path; the code-path edge classes are n/a (the ft4
precedent). The documentation-specific edges walked:

- **anchor inside a rewritten sentence** → covered by the story-4 row: the
  docs families re-run after every edit; any load-bearing phrase that moves
  turns its family red.
- **stale cross-references after the demotion** → covered: the stale-reference
  sweep walks the edited files; a reference to a moved list line is caught
  structurally where tokens are involved, reviewer-checked where prose.
- **skills-index regeneration forgotten** → covered by the story-5 row (the
  generated-index equality check).
- **shipped-file dogfood referents** (BENCH-reference must stay
  consumer-generic while gaining the plumbing list) → the existing
  dogfood-referent check runs on both files; the build keeps the new text
  generic. Covered by the existing family, no new row.
- **re-run idempotency, interrupted state, hostile input** — **Won't handle**:
  no runtime surface.

## Out of scope

- **A conformance family for ADR currency** (detecting an ADR whose framing
  lags shipped mechanisms) — semantic drift detection is a separate capability
  and probably a review-time duty, not a gate check. No estimate — needs a
  shape decision first.
- **Auditing the other ADRs and reference sections for similar drift** — a
  separate sweep with its own findings; this spec fixes the four instances the
  assessment verified. Estimate: ~1 session of read-only sweep, then its own
  drain.
