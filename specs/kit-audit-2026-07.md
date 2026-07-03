# Kit audit — skills, docs, shell (2026-07-02)

Findings from a four-axis audit of the kit: skills prose/consistency, stale
references, shell-script bugs, and comparison against the reference skill corpus
(`reference-skill-repos/`, `gl-axi/skills/`, `regroup-app/.agents/skills/`).
Every finding below was verified against the files on disk; the guard bypass was
reproduced live. This doc is the pickup surface for a fresh session — it holds
the findings and the agreed packaging, not the fixes.

## State at audit time

Gate green. Learnings journal empty (fully integrated). Structure budgets pass.
All cross-wiring verified correct: every `bench` subcommand named in docs exists
in the dispatch, exit-handoff chains are consistent, Claude/Codex skill and hook
wiring matches the docs, `projects/benchkit.md` tier binding matches
`.bench/lines.env`.

## Package A — git-guard hardening (bug path, do first)

Route: `/bench-debug`. Red contract rows first, committed into the gate before
any shift (a red-gate rollback would otherwise destroy them), then the fix.
Guard-boundary code: run on the top tier despite the small diff.

### A1. Newline-separated commands bypass the guard (HIGH, reproduced)

`.bench/hooks/git-guard.py` treats `; && || | & ( )` as command boundaries but
not a literal newline — shlex consumes it as whitespace, so only the first
line's verb is classified. A multi-line command block (`git add` / `git commit`
/ `git push` on separate lines) is the most common way an agent batches git
operations, and the push sails through. The wrapper form
(`bash -c "cd x\ngit push"`) leaks the same way.

Reproduced 2026-07-02:

```
python3 .bench/hooks/git-guard.py "$(printf 'git add -A\ngit commit -m wip\ngit push origin main')"
→ empty verdict (allowed)
python3 .bench/hooks/git-guard.py 'git status && git push origin main'
→ "git push" (blocked)
```

The allow/block matrix in `.bench/gate-runtime-shift-contracts.sh` /
`gate-runtime-git-contracts.sh` covers `;`, `&&`, and wrappers but has no
newline row — add one red first.

Sketch of fix: in `tokenize()`, split the input on `\n` and join the per-line
token streams with a synthetic `;` so `scan()` resets its expect-command state.
The recursive wrapper scan re-tokenizes the inner string, so the wrapper case is
fixed by the same change. Tradeoff: a newline inside a quoted argument also
becomes a boundary — acceptable for an honest-mistake guard; it errs toward
blocking.

### A2. Guard fails open when `python3` is missing (MED, parked decision — decide here)

`.bench/hooks/block-dangerous-git.sh` extracts the command via
`python3 -c ... || true`; with no `python3` on PATH the extraction is empty and
the empty-command early-exit allows everything, including `git push`. This is
asymmetric with the same script's missing/empty-analyzer branch, which fails
closed with a dedicated contract. Same "cannot classify" condition, opposite
verdict. Already parked on the roadmap; A1 raises its priority — resolve the
fail-open/fail-closed decision as part of this package. If fail-closed: require
`command -v python3` after the analyzer-presence check, exit 2 with a BLOCKED
message, and contract it.

## Package B — skills and docs polish (one spec)

Route: `/bench-write-spec`, one batch spec. Merge with the already-parked
roadmap item "skill polish batch" — its items (single-source the 7-edge-class
list, craft-line token-cap sizing, craft-synthesis dogfood carve-out,
craft-seams structure-gate split, craft-design-system harness-neutral naming)
and the findings below are one work package. Kit edits: `craft-synthesis`
discipline applies; guidance prose is leverage per `craft-line`.

Ranked:

1. **`craft-line` decision table has no tie-break** (med-high).
   `bench-craft-line/SKILL.md:33-38` — the genuinely-uncertain row and the
   weak-gate row both match a stage that is both, but emit different lines, and
   "bump one tier" names no baseline. A weak reader under-provisions the
   uncertain seam — the error the skill itself calls expensive. Fix: rows read
   top-down, first match wins (or: genuine uncertainty dominates), and the bump
   row states its baseline ("one tier above whatever the other two signals
   selected").
2. **`craft-seams` and `craft-tdd` lack the contrastive good/bad pair**
   (med-high). `craft-skills` mandates one for any skill governing an output
   surface; `craft-seams` governs interface shape and ships zero examples. The
   reference corpus has directly reusable pairs (testable-vs-hardwired
   constructor injection; mocking GOOD/BAD at an SDK interface).
3. **README Layout drift** (med). `README.md:145-176` omits
   `check-agent-line.sh` (an active deny-capable hook), four of seven `bin/`
   scripts (`bench-query.sh`, `bench-diff.sh`, `bench-coverage.sh`,
   `bench-worktree.sh`), and `bench-craft-line` from a skills list that reads
   as exhaustive.
4. **Split `craft-seams` branch-only sections into `references/`** (med).
   "Design it twice" and the structure-gate/splitting treatment are reached by
   only some runs; the file is the longest craft skill. Already parked on the
   roadmap.
5. **Paste-ready sub-agent prompts in "Design it twice"** (med). Replace the
   four concept bullets with literal briefs the agent pastes (reference
   pattern: `Agent 1: "Minimize the interface — 1-3 entry points max…"`).
6. **Quoted user-phrasing triggers in frontmatter** (med). `craft-tdd`,
   `craft-adr`, `craft-seams` descriptions are conceptual; references quote the
   user's words (`mentions "red-green-refactor"`) and fire more reliably. One
   quoted trigger each.
7. **Trim `craft-cli` description** (med). ~40 words of operational-vs-AXI
   scoping nuance sit in the description (in context every turn) and duplicate
   the body verbatim. End the description after the "use whenever" clause.
8. **`craft-skills` hand-maintained roster** (low-med). Lines 16-19 enumerate
   the nine model-invoked skills — a second derivation of a
   frontmatter-derivable fact; drifts on the next add/remove. Replace with the
   rule: every `craft-*` skill is model-invoked, phase adapters are not.
9. **Bolded stop marker at `craft-tdd`'s over-fit point** (low-med). The
   central over-fit guard is buried in prose; borrow the reference pattern
   ("stop — this is the exact failure this skill prevents") at the
   pre-agreed-seams section.
10. **Cut `craft-design-system` repeated rationale** (low-med). The "seamless"
    thesis is stated three times around one instruction at lines 54-70; the
    two-bullet harness contrast carries it. Adjacent to the parked
    harness-neutral-naming item.
11. **Gate self-check 1d prefix match** (low). `.bench/gate.sh:56` —
    `\.bench/(gate|done)(\.sh)?` prefix-matches future `gate-*` references in
    `bin/bench.sh` and would flag them as extensionless refs. Harmless today;
    anchor the alternation.
12. **Prose sediment** (low). `craft-grill` states one-question-at-a-time three
    times (keep the Discipline bullet); `craft-seams` Pocock attribution is
    five lines for one actionable clause (keep the optional-upgrade clause).
13. **Dropped failure-mode names** (low). Reference skill-writing guidance
    names five diagnostic failure modes; `craft-skills` keeps three and drops
    Sprawl and No-op as named entries. Likely deliberate leanness — confirm or
    restore, reviewer's call.

## Stays parked (roadmap already holds these; not pulled into A or B)

- Single-source the `lines.env` tier parser duplicated across
  `check-agent-line.sh` and `_line-guard.sh`.
- Missing craft skills: craft-gate, craft-review, craft-delegate.
- Workflow-exit gaps (capped/unmet shift routing, superseded-spec retirement).

## Non-findings (checked, no action)

- Pre-rename paths in `specs/bench-naming-mechanical.md` are intentional and
  self-documented with the historical marker.
- `.bench/bin/` absent from this repo is correct — it is what `bench link`
  installs into consumer repos; both hooks probe both locations.
- `craft-cli`'s AXI-conformant set matches the gate's own enumeration; `bench
  diff`/`coverage` emitting TOON without formal conformance is consistent.

## Don't regress (kit strengths the reference corpus lacks)

- Completion deferred to the gate everywhere — references end on self-graded
  checklists.
- `craft-line`'s three-signal routing table (fix only the tie-break) and its
  enforcement wiring via the Agent-tool hook.
- Anti-skip reach clauses ("even on small changes") and cross-harness
  portability of the whole skill surface.
