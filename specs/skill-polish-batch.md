# Skill polish batch — skills and docs sharpening (Package B of the 2026-07 kit audit)

## Problem

The kit's guidance prose has accumulated drift and dullness that a fresh agent
pays for on every load: an ambiguous routing table in `craft-line`, mandated
contrastive examples missing from the two skills that govern output surfaces,
a README layout that under-reports the kit, triplicated edge-class lists that
will drift, and assorted sediment (repeated rationale, attribution bloat,
hand-maintained rosters). Each defect is small; together they lower the
probability that an agent generates the intended response — and guidance prose
multiplies through every session that loads it.

## Solution

One batch pass over the skill and doc surface: fix the routing ambiguity, add
the mandated contrastive pairs, split branch-only content into `references/`,
single-source the duplicated lists, sharpen frontmatter triggers, sync the
README, and anchor one latent gate regex. Merges the audit's 13 polish
findings with the parked "skill polish batch" roadmap items into a single
coherent work package; the roadmap line retires when this lands.

Source findings: `specs/kit-audit-2026-07.md` Package B (this spec supersedes
its ranked list as the work surface; the audit doc stays as evidence).

## User stories

Actor key: *loading agent* = any agent that loads the skill in a future
session; *reviewer* = the human who owns the merge.

1. As a loading agent, I want the `craft-line` decision table to resolve
   unambiguously when a stage is both genuinely uncertain and weakly gated —
   rows read top-down, first match wins, and the bump row names its baseline
   ("one tier above what the other signals selected") — so that I never
   under-provision the uncertain seam.
   Line: claude-fable-5 / high. This is guidance prose that routes every
   future line declaration, so the leverage override applies.
2. As a loading agent, I want `craft-line` to give a derivation method for the
   ~token cap (sized from expected iterations times per-iteration cost, with
   the review-delegate 60k row as the worked example), so that declared caps
   are estimates rather than vibes.
   Line: claude-fable-5 / high. The cap method steers every declared line, so
   the leverage override applies.
3. As a loading agent, I want `craft-seams` to carry one contrastive good/bad
   pair (testable constructor-injection seam vs hardwired dependency), so that
   the skill meets the standard `craft-skills` mandates for output-surface
   skills and I can pattern-match instead of interpret.
   Line: claude-fable-5 / high. A canonical example shapes every interface an
   agent designs, so the leverage override applies.
4. As a loading agent, I want `craft-seams`' branch-only sections ("Design it
   twice" and the structure-gate/splitting treatment) moved to `references/`
   files with one-line pointers from `SKILL.md`, so that the always-loaded
   body stays short and the deep content is reached exactly when its branch
   applies.
   Line: claude-fable-5 / high. Restructuring the most-loaded craft skill is
   leverage work, so the override applies.
5. As a loading agent, I want the "Design it twice" guidance to give literal
   paste-ready sub-agent briefs instead of four concept bullets, so that
   running the technique is a copy, not a composition exercise.
   Line: claude-fable-5 / high. The briefs are executed verbatim by future
   agents, so their wording is the whole value and routes top.
6. As a loading agent, I want the Pocock attribution in `craft-seams` cut to
   its one actionable clause, so that five lines of provenance stop diluting
   the instruction density.
   Line: claude-fable-5 / high. Same file and same leverage as stories 3-5,
   so it rides the same routing.
7. As a loading agent, I want `craft-tdd` to carry one contrastive good/bad
   pair (mocking at an owned SDK interface vs mocking internals), so that the
   over-mocking failure the skill warns about has a concrete picture.
   Line: claude-fable-5 / high. A canonical test example steers every TDD
   pass, so the leverage override applies.
8. As a loading agent, I want a bolded stop marker at `craft-tdd`'s
   pre-agreed-seams section ("stop — this is the exact failure this skill
   prevents"), so that the central over-fit guard interrupts generation
   instead of hiding in prose.
   Line: claude-fable-5 / high. The marker exists to alter agent behavior at
   the highest-risk moment, so it routes top.
9. As a loading agent, I want the frontmatter descriptions of `craft-tdd`,
   `craft-adr`, and `craft-seams` to each quote one user phrasing (for
   example: mentions "red-green-refactor", "write an ADR", "where should this
   test live"), so that the skills fire on the words users actually type.
   Line: claude-fable-5 / high. Descriptions are the firing mechanism for the
   whole skill, so the leverage override applies.
10. As a loading agent, I want `craft-cli`'s description ended after its "use
    whenever" clause, with the ~40 words of operational-vs-AXI scoping nuance
    living only in the body, so that every-turn context spends fewer tokens
    and the fact has one source.
    Line: claude-fable-5 / high. Trimming a trigger without dulling it is
    precision prose work, so it routes top.
11. As a loading agent, I want `craft-skills`' hand-maintained nine-name
    roster replaced with the rule that generates it (every `craft-*` skill is
    model-invoked; phase adapters are not), so that the fact has one source
    and cannot drift on the next skill add.
    Line: claude-fable-5 / high. The rule governs future skill authoring, so
    the leverage override applies.
12. As a loading agent, I want the two dropped diagnostic failure modes
    (Sprawl, No-op) restored to `craft-skills` as one-line named entries,
    completing the five-mode vocabulary the reference guidance uses — unless
    the reviewer confirms the three-mode leanness was deliberate.
    Line: claude-fable-5 / high. Failure-mode names are diagnostic vocabulary
    for future skill reviews, so the leverage override applies.
13. As a loading agent, I want `craft-grill` to state one-question-at-a-time
    once (the Discipline bullet stays; the other two statements go), so that
    the skill's real content isn't padded by repetition.
    Line: claude-fable-5 / high. Same leverage class as the other skill-prose
    stories, so it rides the same routing.
14. As a loading agent, I want `craft-design-system` to state its "seamless"
    thesis once (the two-bullet harness contrast carries it), name tools
    harness-neutrally, and give a stop route for headless runs, so that the
    skill works identically on every harness the kit claims to support.
    Line: claude-fable-5 / high. Cross-harness skill prose is the product's
    portability claim, so the leverage override applies.
15. As a loading agent, I want `craft-synthesis` to carve out proportionality
    for prose-only kit edits (the full dogfood discipline scales down when no
    behavior changes), so that small wording fixes stop inheriting the full
    ceremony.
    Line: claude-fable-5 / high. The carve-out changes when a whole discipline
    applies, so it routes top.
16. As a maintainer of the kit, I want the seven-edge-class list to have one
    canonical enumeration (in `/bench-write-spec`, the phase that walks it)
    with `craft-tdd` and `/bench-review-implementation` keeping their
    gate-required anchors as pointers to it, so that the list cannot drift
    into three divergent copies.
    Line: claude-fable-5 / high. The edge-class list defines coverage breadth
    for every feature build, so the leverage override applies.
17. As a reviewer, I want `README.md`'s Layout section to match disk — adding
    `check-agent-line.sh`, the four missing `bin/` scripts (`bench-query.sh`,
    `bench-diff.sh`, `bench-coverage.sh`, `bench-worktree.sh`), and
    `bench-craft-line` — so that a list that reads as exhaustive is
    exhaustive.
    Line: claude-opus-4-8 / medium. Mechanical enumeration against disk with
    no judgment content, but the gate cannot grade completeness, so it takes
    the mid row rather than cheap.
18. As a maintainer of the kit, I want the gate self-check 1d alternation
    anchored so `\.bench/(gate|done)(\.sh)?` stops prefix-matching future
    `gate-*` references in `bin/bench.sh`, so that the check cannot produce a
    false positive when a legitimately-named script appears.
    Line: claude-opus-4-8 / medium. Gate code takes the profile's cached
    gate/conformance routing of mid effort.

## Implementation decisions

- **Discipline**: every edit follows `craft-synthesis`; guidance prose is
  leverage per `craft-line` (hence the top-tier lines above). Kit content
  edits regenerate the skills index (`.bench/skills-index.sh --write`)
  whenever `index:` frontmatter changes — the gate diffs committed vs
  generated.
- **Progressive-disclosure rule for splits (story 4/5)**: `SKILL.md` keeps the
  trigger sentence and a one-line pointer; the moved body lands in
  `references/<topic>.md`. A split that leaves no pointer is a defect (dead
  content), not a trim.
- **Single-source location (story 16)**: the canonical enumeration lives in
  the `/bench-write-spec` command (its step 5 is the walking site). The other
  two files keep the anchor wording the gate's conformance layer requires,
  reworded as pointers — the gate must stay green through the change, which is
  the check that the anchors survived.
- **Descriptions vs body (stories 9, 10)**: trigger phrasing lives in
  frontmatter descriptions; scoping nuance lives in the body. One fact, one
  source.
- **Roadmap retirement**: the parked "skill polish batch" line in `ROADMAP.md`
  is deleted in the same change that lands this spec's implementation — the
  idea has graduated.
- **No behavior changes** outside story 18's one-character-class regex anchor;
  everything else is markdown. Story 18 touches `.bench/gate.sh` only.
- Story 12 defaults to restoring the two names; the reviewer can strike it at
  approval if the leanness was deliberate.

## Testing decisions

What a good test is here: the gate is the only oracle with authority, and for
prose it grades **structure, wiring, and anchors** — not semantics. So the
honest testing posture is: (a) every structural consequence of an edit
(frontmatter validity, skills-index equality, anchor presence, npm files[]
resolution) is gate-covered and must stay green through every story; (b) the
one behavior change (story 18) gets a direct red demonstration; (c) prose
semantics are graded by `/bench-review-implementation` against
`specs/kit-audit-2026-07.md`'s finding descriptions — declared here as review
coverage, not pretended into TDD coverage.

Prior art: the gate's kit-conformance layer (layer 3) and canary layer
(layer 7) already test this surface; `tests/canary/` fixtures are the
established mechanism for proving a gate check bites.

Gate command: `bench gate` (the project gate, `.bench/gate.sh`).

### Seam diagram

Seam 1 — kit content surface, fronted by the gate's conformance layer:

    trigger: bench gate (run after every story's edit)
        │
        ▼
    skill/command markdown edits  ──▶  [ gate layers 2+3:            ]  ──▶  green / red + per-file attribution
    .bench/BENCH.md index block   ──▶  [ frontmatter, skills-index    ]
                                       [ equality, anchors, files[]   ]
                      ◀ tests attach here: run `bench gate`; red names the drifted file

Seam 2 — gate self-check 1d, fronted by the gate run itself (story 18):

    trigger: bench gate (layer 1 self-checks)
        │
        ▼
    bin/bench.sh source lines  ──▶  [ 1d extensionless-ref regex ]  ──▶  flag / pass
                      ◀ tests attach here: grep -E the old and new pattern against a
                        fixture line containing a `gate-*` name; old flags it, new passes it

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2 | craft-line table resolves top-down with named baseline; cap has a method | seam 1 | not TDD-able — prose semantics; gate has no grader for table logic. Review coverage: reviewer reads the two rows against audit finding 1 | a wrong fix leaves the ambiguity the finding describes; only a reader catches that |
| 3, 7 | contrastive pair present in craft-seams and craft-tdd | seam 1 | not TDD-able as semantics; structural half already covered (frontmatter/structure checks stay green). Review coverage: pair exists and contrasts | gate proves the file still conforms; only review proves the example teaches |
| 4, 5 | branch-only content moved to references/ with live pointers | seam 1 | already covered structurally — gate layer 2/3 stays green through the move; a broken skills-index or lost anchor goes red | the move's failure mode (anchor or index drift) is exactly what layer 3 flags per file |
| 6, 8, 13, 14, 15 | sediment cut; stop marker added; carve-out stated once | seam 1 | not TDD-able — pure prose. Review coverage against audit findings 9, 10, 12 | no observable seam exists for wording; declaring review coverage is the honest row |
| 9, 10 | descriptions carry quoted triggers; craft-cli description ends at "use whenever" | seam 1 | not TDD-able for firing behavior (trigger matching happens in the harness, not the gate); structural: gate green on frontmatter | a dulled trigger is invisible to every automated check we own; flagged as Won't-handle below |
| 11, 12 | roster replaced by rule; failure modes restored | seam 1 | not TDD-able — prose. Review coverage: no enumerated skill names remain in craft-skills body | the defect is a second derivation; a reader greps for the nine names and finds none |
| 16 | edge-class list enumerated once; other files point at it | seam 1 | already covered — gate layer 3's edge-coverage anchor check goes red if an anchor is lost in the rewrite (verified bite: canary layer asserts layer-3 checks fail on broken fixtures) | losing an anchor is the drift this story risks, and it is precisely what the existing check flags |
| 17 | README Layout matches disk | seam 1 | not TDD-able — no gate check grades README completeness. Review coverage: diff the section against `ls` of hooks/, bin/, skills/ | future drift has no oracle either; see Won't handle and the parked `bench doctor` idea |
| 18 | anchored 1d regex passes `gate-*` refs, still flags true extensionless refs | seam 2 | direct red demonstration at implement time: `grep -E` the current pattern against a `gate-runtime-x.sh`-style line — it matches (wrongly) today; the anchored pattern must not, while the existing canary fixture for 1d stays red | the pair (old matches / new doesn't, canary still bites) is exactly the false-positive-fixed-without-false-negative contract |

### Edge inventory

Walked per behavior against the profile's checklist (shell CLI classes apply
only to story 18; the rest is a markdown surface with its own edge classes —
wiring, anchors, index equality, pointer liveness):

- **Frontmatter edit without index regeneration** (stories 9, 10) — covered:
  gate layer 3 index-equality goes red; that red is the regression signal.
- **Split leaves a dead references/ file or a pointerless SKILL.md** (stories
  4, 5) — coverage row 4/5 plus the implementation decision making a
  pointerless split a defect by definition; review checks pointer liveness.
- **Anchor lost while single-sourcing** (story 16) — covered: layer-3 anchor
  check, bite proven by canary.
- **npm files[] breakage from new references/ files** (story 4) — covered:
  gate layer 2 resolves every packaged path; new files under `.agents/` are
  already in the packaged tree.
- **Regex over-anchor** (story 18) — the fix must not stop flagging genuine
  extensionless `\.bench/gate` refs; covered by the existing 1d canary
  fixture staying red.
- **Won't handle: trigger-firing regression** (stories 9, 10) — no owned
  seam observes harness trigger matching; mitigated by keeping each "use
  whenever" clause intact, accepted as unobservable.
- **Won't handle: future README/layout drift** (story 17) — building a
  drift check is the parked `bench doctor` idea, a separate capability; this
  story fixes the instance, not the class.
- **Won't handle: semantic drift between the canonical edge-class list and
  its pointers' paraphrases** (story 16) — pointers will name the list, not
  paraphrase it, which removes the surface; a wording-equality gate check is
  not worth its brittleness.

## Out of scope

- **Single-source the `lines.env` tier parser** (`check-agent-line.sh` /
  `_line-guard.sh`) — a shared-lib extraction across the hook/adapter
  boundary, its own capability with its own contract rows. Estimate: ~6
  edits, ~4 gate runs. Stays parked.
- **Missing craft skills (craft-gate, craft-review, craft-delegate)** — each
  is a new artifact with its own shaping, not polish of an existing one.
  Estimate: ~5 edits, ~3 gate runs per skill. Stays parked.
- **Workflow-exit gaps** (capped/unmet shift routing, superseded-spec
  retirement) — command-phase feature work, not prose polish. Estimate: ~8
  edits, ~4 gate runs. Stays parked.
- **Package A (git-guard newline bypass + python3 fail-open)** — the bug
  path, routed to `/bench-debug` per `specs/kit-audit-2026-07.md`; deliberately
  not in this spec. Estimate: ~6 edits, ~5 gate runs.
