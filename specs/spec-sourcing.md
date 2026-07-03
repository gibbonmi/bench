# Spec sourcing — every spec compiled from a closed map's Handoff

Status: implemented

Source map: `decisions/spec-sourcing.md` (all four tickets resolved in the
discovery grill; Handoff item 7 flags no uncertainty). This spec is the map's
only slice and carries every decision forward.

## Problem

`/bench-shape-idea` advertised two skip-to-spec bypasses ("already clear / no
real fog → skip straight to `/bench-write-spec`") and `/bench-write-spec`
accepted a no-map source ("sketch the seams from scratch"). Together they let the
top model author a spec directly from conversation — bypassing the grill → map →
Handoff path that keeps seams pre-agreed and context small. Nothing in the gate
bit; the drift was pure prose.

## Solution

Close the bypass: every spec is compiled from a `decisions/<topic>.md` whose
`## Handoff` is complete, however simple the idea — a simple idea yields a short
grill and a zero-open-ticket map in one sitting. `/bench-shape-idea` loses both
bypass sentences (keeping "don't manufacture tickets"). `/bench-write-spec`
refuses to run without a named map whose Handoff `bench maps` shows closed, is
authored by the **mid tier in a fresh session** (grill session ends at the map),
and spawns a **conditional** top-tier reviewer sub-agent only when the Handoff
carries uncertainty flags or the draft deviates from the map. Two thin gate
anchors pin the Handoff-sourcing prose in both command files, each with a red
canary, so the bypass cannot silently drift back.

## User stories

1. As a reviewer shaping an idea, I want `/bench-shape-idea` stripped of both
   skip-to-spec bypass sentences — the entry-orientation "if the idea is already
   clear, skip straight to `/bench-write-spec`" and the exit "recommend skipping
   straight to `/bench-write-spec`" — replaced by the rule that even a clear idea
   yields a (short, zero-open-ticket) map with a `## Handoff`, while "don't
   manufacture tickets" stays, so no idea reaches a spec without a map behind it.
   Line: claude-opus-4-8 / high. The map fully settled this prose but it is
   gate-blind beyond a single anchor, so it gets the mid model at high effort to
   get the wording right rather than cheap transcription.

2. As a spec-writer, I want `/bench-write-spec`'s entry contract to **refuse to
   run** without a named `decisions/<topic>.md` whose `## Handoff` is complete —
   `bench maps` showing no row for it — replacing the "where no map produced this
   spec, sketch the seams from scratch" path with "name the map to close and
   stop; never draft from conversation alone," so the no-map source is closed at
   the consumer.
   Line: claude-opus-4-8 / high. This is the load-bearing entry contract of the
   feature and its exact wording is what the anchor pins and review reads, so it
   stays on mid at high effort.

3. As the transition between phases, I want `/bench-write-spec` to state the
   venue rule — the mid tier authors every spec, the default venue is a fresh
   mid-tier session that starts cold on the map plus the repo, and a same-session
   read-only (or worktree-isolated) mid delegate is the exception taken only on
   the reviewer's explicit ask — and `/bench-shape-idea`'s exit to recommend that
   fresh mid-tier session for `/bench-write-spec`, so the top model never authors
   a spec and big-context sessions do not linger as idle orchestrators.
   Line: claude-opus-4-8 / high. The venue rule is gate-blind prose the map
   decided, and its wording defines who authors specs, so it stays on mid at high
   effort.

4. As a spec-writing session, I want `/bench-write-spec` to describe the
   **conditional** top-tier reviewer sub-agent — spawned in a fresh small context
   (Handoff plus draft only), firing only when the Handoff carries uncertainty
   flags (item 7) or the draft deviates from the map, returning findings and an
   advisory recommend/block verdict, with no standing top-tier pass and sign-off
   staying the reviewer's — so a genuinely uncertain draft gets a model check
   without moving merge authority to a model.
   Line: claude-opus-4-8 / high. This is gate-blind guidance prose the map
   settled, and its wording governs when a top-tier model spends, so it stays on
   mid at high effort.

5. As a cold session, I want `projects/benchkit.md`'s Lines section to carry the
   spec-authoring routing note — spec authoring runs mid tier in a fresh session,
   distinct from the top-tier command/doc-authoring leverage override — so the
   routing is discoverable where a cold session already reads its lines.
   Line: claude-opus-4-8 / medium. It is a small routing note the map decided,
   gate-blind, so it runs on mid at medium effort.

6. As the gate, I want two thin kit-conformance anchors — `/bench-shape-idea`
   must not contain the skip-to-spec bypass fragment (negative), and
   `/bench-write-spec` must contain the map-required refuse contract (positive),
   each anchor matching against whitespace-normalized text so a hard wrap that
   splits the fragment across a newline cannot defeat the grep — plus a red canary
   fixture per anchor proving it bites (and a third, wrapped, fixture proving the
   negative anchor's wrap-tolerance bites), so the Handoff-sourcing prose cannot
   silently drop out of either command.
   Line: claude-opus-4-8 / medium. Anchor checks and canary fixtures are oracle
   content where a wrong anchor is an always-pass, and the profile routes gate
   correctness to mid effort.

## Implementation decisions

- **Two command files change prose; one profile file gains a note.**
  `/bench-shape-idea` loses two bypass sentences and gains the map-is-mandatory
  rule plus a fresh-mid-session exit recommendation. `/bench-write-spec` gains the
  refuse-without-complete-map entry contract (replacing the no-map path), the
  venue rule, and the conditional-reviewer rule. `projects/benchkit.md` Lines
  gains the spec-authoring routing note. No new shell modules, no schema change.

- **The refuse contract is backed by the existing `bench maps` parser.**
  "Handoff complete" means `bench maps` emits no row for the map (its
  refuse-to-close loop already keeps a row until the Handoff is present and
  placeholder-free). The prose tells the agent to run `bench maps`; there is no
  new parser and no new subcommand.

- **Two anchors, opposite polarity, both thin greps that pin sentences, not
  semantics.** `/bench-shape-idea` gets a **negative** anchor (the bypass
  fragment ``straight to `/bench-write-spec` `` must be absent — it is the
  invariant part of any skip bypass and does not appear in the legitimate exit
  recommendation). `/bench-write-spec` gets a **positive** anchor (the phrase
  `refuses to run without` must be present). Each is an inline `grep`/`err` in the
  kit-conformance docs fragment, guarded on the file existing, matching the
  existing inline-anchor precedent (the roadmap-promotion seam) rather than
  overloading the `require_anchor` helper whose message reads "acceptance coverage
  anchor."

- **Anchors match against whitespace-normalized text (the discovery repro's
  lesson).** A line-oriented fixed-string grep silently misses a needle that spans
  a hard wrap; the discovery repro hit exactly this. Keeping our own needle on one
  line covers the positive anchor (we control write-spec's wrapping) but not the
  negative one — a future editor reintroducing the bypass controls its own
  wrapping. So both anchors pipe the file through `tr -s '[:space:]' ' '` (collapse
  every whitespace run, newlines included, to one space) before the fixed-string
  grep, so a wrap that splits the fragment across a newline cannot defeat the
  match. A third canary (`shape-idea-bypass-wrapped`) plants the wrapped
  reintroduction and proves the wrap-tolerance bites.

- **Kit prose routes mid + high, a declared step down from the leverage
  override.** `craft-line`'s leverage override would route command-phase edits
  top + high, but the map settled every decision, the Handoff-sourcing prose is
  anchored, and the venue/reviewer/profile prose lands on review — so this is
  high-effort transcription of decided content, not guidance design. This mirrors
  `specs/spec-handoff-lifecycle`, which the reviewer approved; the reviewer may
  bump any story to top.

- **Venue and review conduct are gate-blind by design (Handoff item 5).** Which
  session authors the spec, and when the conditional reviewer fires, are prose the
  gate cannot see; they get no anchor and land on the review axis and reviewer
  veto. Only the two Handoff-sourcing anchors are gate-observable.

## Testing decisions

- A good test here runs the real gate against a minimal fixture tree in a
  throwaway repo and asserts it goes red with the anchor's targeted substring —
  never a reading of the prose. Prior art: the handoff-anchor canaries
  (`tests/canary/shape-idea-handoff-anchor`,
  `tests/canary/write-spec-handoff-anchor`) and the whole `tests/canary/` suite.
- One seam is tested: the kit-content conformance layer (gate layer 3 anchors,
  proven to bite by layer 7 canaries). Stories 1–5's prose semantics beyond the
  two anchored sentences are gate-blind and land on review, not a new seam.
- Gate command: `.bench/gate.sh` (the project gate), which already runs the
  conformance anchors and the canary suite.

### Seam diagram — kit-content conformance (Handoff-sourcing anchors)

    trigger: every gate run (layer 3 anchors) / canary run (layer 7)
        │
        ▼
    .agents/commands/         ──▶  [ Handoff-sourcing anchors:      ]  ──▶  err() lines,
      bench-shape-idea.md          [  shape-idea: bypass fragment    ]       gate red/green
      bench-write-spec.md     ──▶  [    ABSENT (negative grep)       ]
                                   [  write-spec: refuse contract     ]
                                   [    PRESENT (positive grep)       ]
                                        ◀ tests attach here: three canary fixtures —
                                          one reintroduces the bypass, one the same
                                          bypass hard-wrapped, one drops the
                                          contract; the gate must go red with each
                                          targeted substring (and stay green on disk)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `/bench-shape-idea` on disk contains no ``straight to `/bench-write-spec` `` bypass fragment | conformance anchor (layer 3) | `rg 'straight to \`/bench-write-spec\`' .agents/commands/bench-shape-idea.md` exits 0 today (two hits, observed red 2026-07-03) — the discovery repro | the negative anchor greps the fragment and reds the gate whenever it is present; removal makes it green |
| 6 (guards 1) | gate errs when the bypass fragment is reintroduced into shape-idea | canary (layer 7) | new `shape-idea-bypass` fixture run with the anchor not yet added — the gate does not err, so the canary's "did not bite" assertion is red | proves the negative anchor bites; an always-pass anchor fails its canary |
| 6 (guards 1b) | gate errs when the bypass fragment is reintroduced **hard-wrapped** across a newline | canary (layer 7) | `shape-idea-bypass-wrapped` fixture: a line-oriented grep MISSES it (`grep -F` exits 1), so a non-normalized anchor does not err and the canary is red | proves the whitespace-normalized negative anchor catches a wrapped reintroduction the future editor controls; the discovery repro's own failure class |
| 2 | `/bench-write-spec` on disk contains the `refuses to run without` map-required contract | conformance anchor (layer 3) | new positive anchor run before the contract sentence is added — gate red with "dropped the map-required entry contract" | the positive anchor greps the contract phrase and reds the gate when it is absent |
| 6 (guards 2) | gate errs when the refuse contract is dropped from write-spec | canary (layer 7) | new `write-spec-map-required` fixture (contract sentence removed) run with the anchor not yet added — gate does not err, canary "did not bite" red | proves the positive anchor bites; a rotted anchor fails its canary |
| 3 | write-spec states the venue rule; shape-idea exit recommends the fresh mid-tier session | review axis | not TDD-able — gate-blind prose (Handoff item 5); verified by the review axis and reviewer veto | no anchor greps semantics; a wrong or missing venue rule is a review finding, not a gate red |
| 4 | write-spec describes the conditional top-tier reviewer sub-agent | review axis | not TDD-able — gate-blind prose (Handoff item 5); review axis + reviewer veto | same as story 3: the gate cannot grade when the reviewer fires |
| 5 | `projects/benchkit.md` Lines carries the spec-authoring routing note | review axis | not TDD-able — gate-blind prose; verified by review that the note reads correctly | a routing note is prose the gate does not parse |

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist, per mapped
behavior (the two anchors):

- **Hard-wrapped bypass prose** — covered (story 6): the Handoff's named domain
  edge, and the class the discovery repro caught. Both anchors match against
  whitespace-normalized text (`tr -s '[:space:]' ' '`) so a wrap that splits the
  fragment across a newline cannot defeat the match — load-bearing for the
  *negative* anchor, whose reintroduction a future editor wraps as it pleases. The
  `shape-idea-bypass-wrapped` canary plants the wrapped form and proves it bites.
- **Command file absent vs present-but-empty** — covered: each anchor is guarded
  on `[ -f "$file" ]`, so an absent file skips the check (and the canary's
  empty-fixture vacuous baseline stays clean); a present file with the fragment
  missing/added fires the targeted err.
- **Fixture is not otherwise conformant** — covered by the runner's
  attribution-by-substring: a minimal single-file fixture over-fails on unrelated
  checks, and only the targeted substring is asserted.
- **Missing trailing newline / CRLF in the command file** — covered: the needles
  are mid-line fragments, so line-terminator variation at the end of the line does
  not affect the match.
- **cwd deeper than repo root; re-run idempotency** — covered: the gate `cd`s to
  the repo root and is read-only, so both are already handled by the gate harness.
- **Won't handle:** a differently-worded future bypass — the negative anchor pins
  the known fragment, not every paraphrase; a novel bypass phrasing is a review
  finding, consistent with "pin sentences, not semantics" (Handoff item 3).
- **Won't handle:** the negative anchor is fence-blind — a file quoting the bypass
  fragment inside a code fence would trip it; safe because only
  `bench-shape-idea.md` is scanned and its post-fix prose does not contain the
  fragment (verified: the legitimate exit recommendation reads "`/bench-write-spec`
  on a fresh mid-tier session", not "straight to").
- **Won't handle:** semantic quality of the venue, reviewer, and profile prose —
  gate-blind by design (Handoff item 5); the review axis owns it.

## Out of scope

- **Adversarial gate pinning** (hash-verify `.bench/gate.sh` outside the writable
  tree in `pre-push`) — a separate threat model (a determined agent weakening the
  gate itself), distinct from this lazy-prose-drift guard, already parked on the
  roadmap — ~6 edits, ~4 gate runs.
