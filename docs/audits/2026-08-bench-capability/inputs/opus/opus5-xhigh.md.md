# Bench — adversarial capability-realization audit

**Auditor:** Claude Opus 5, effort xhigh, Claude Code
**Target:** `gibbonmi/bench` @ `58d966e2`, branch `audit/opus`
**Date:** 2026-08-17
**Evidence appendix:** `audit/opus-xhigh/evidence/`

---

## A. Executive verdict

**Is Bench directionally correct?** Yes, and more so than I expected going in. Bench
is not ceremony. Its central bet — *replace model self-assessment with a
content-addressed external oracle, and mechanically enforce the guidance
invariants that a prose-only skill kit can only hope for* — is the right bet, and
large parts of it are built to a standard I could not fault after probing them.
`bench gate` runs in 64 s, caches on a `git write-tree` hash that correctly
includes untracked-unignored files, and I verified invalidation in both
directions. `bench coverage --check` caught 4 of 4 hostile spec-map inputs with
exact messages. `bench preflight` refused an empty ownership fence and then
caught an unauthorized path. The git guard blocks 20 destructive spellings
including `git restore --worktree` while correctly allowing `git restore --staged`.
The delegation guard denies an omitted `model` field — closing the silent-escalation
hole rather than just naming it. `bench setup` is plan-first and self-verifying.
Across 27,758 words of agent-facing prose there are **zero** instances of
*carefully / thoroughly / appropriately / robustly / properly*. This is disciplined work.

**Its strongest architectural idea** is the tree-hash-keyed verdict: one
content-addressed fact ("is *this exact tree* green?") in place of the claims-and-freshness
database the audit brief assumed I would find. Bench does not have a claim system,
and it is right not to. Second strongest is the **conformance suite** — 29 checks
that make doc/guidance invariants mechanically enforceable, which is the one thing
Bench genuinely adds over Pocock's prose.

**What must be preserved:** the gate-as-oracle chain (`bench gate` → `bench commit`
→ Stop hook), the conformance idea, `/bench-debug`'s repro-loop-first discipline
(which preserves every load-bearing constraint of Pocock's `diagnosing-bugs` and
adds three real ones), `craft-review`'s citation standard and "refute before you
report", `craft-delegate`'s independent-probe rule, and `bench setup`.

**The largest weakness is not complexity — it is that the oracle does not grade the
thing Bench is mostly made of.** `internal/conformance` is the repo's biggest
package (13,250 LOC, 29 registered checks). It reaches the live tree only through
`TestRootConformance`, which opens with `if os.Getenv("BENCH_CONFORMANCE_ROOT") == ""
{ capability.Environment(t, ...) }`. Nothing in `bench gate` sets that variable —
the only setter in the repo is `bench prep-release`. So `go test ./...` prints
`ok internal/conformance` while the live-root suite silently skips, and the gate
reports green. Run it by hand at HEAD and **10 checks are red**, including
`.agents/commands/bench-implement-spec.md missing Entry orientation` and
`missing Exit handoff` — a structural contract on a workflow phase file, broken by
commit `fa4e1f02` six days ago, which landed green through Bench's own pipeline.
The gate stopped grading the live root on 2026-07-05 (`72c037a1`), **43 days ago**.
Bench's own `craft-tdd` names the exact trap — *"a run whose only matched test
skipped still prints `ok` (Go does), so read the test output, not the summary"* —
and Bench's own `/bench-assess` run four days ago enumerated the gate's skips and
walked past the one that matters.

**Second weakness: there is no front door.** `bench status` in an un-adopted repo
prints `bench: clean — nothing pending`. A *valid staged spec* — the single most
action-relevant fact in the workflow — produces **no status row at all**; the
`specs` signal only fires for specs *awaiting retirement*. `/bench-what-next` is
labelled "Maintenance, not a workflow phase" and is hidden from the model. Pocock's
suite ships `ask-matt` ("*a router over the skills in this repo — you don't remember
every skill, so ask*"); Bench inherited the pipeline and **did not port the router**,
while growing to 11 phases.

**Third: a new repo's gate cannot go green.** `bench setup` scaffolds no
`.bench/gate-inputs.json`, so the gate runs its command with `PATH` only; the
scaffolded `gate.sh` then calls the repo-local wrapper, whose line 14
`${BENCH_HOME:-$HOME/.bench}` dies on unbound `$HOME` under `set -u`. Declare
`HOME` and it is *still* red, because the same scaffolded gate hard-fails on
`bench canary`, which exits 1 on the empty fixture inventory every new repo has.
I reproduced the whole chain in a throwaway repo.

**What should be removed:** `bench outline`'s bare invocation (24.6 KB to emit 200
of 5,335 symbols under no stated ranking); the 34 orphaned `specs/*/tickets/`
directories no command owns; the `bench structure` signal in its current form
(62 issues, unenforced, +12 in four days, with four rows of rot inside its own
suppression file); and one of `bench status`/`bench handoff`'s two verdict readers.

**What should be strengthened:** wire root conformance into the gate; fix the
verdict-staleness rule so a **red** record is drift-checked (it currently is not —
see below); add a `specs: staged` status signal; ship a router.

**Rewrite, consolidate, or improve incrementally?** **Incrementally, with one
structural consolidation.** A rewrite would destroy demonstrated behavior — the
oracle, the guards, the conformance checks — to fix problems that are, without
exception, small and local. The three highest-impact defects I found are a missing
env var, a mis-ordered `if`, and a missing status row.

**What to build next:** wire the 29 conformance checks into `bench gate` as a
first-class phase, and make an environment-class capability skip inside the oracle
a **red**, not a footnote. Everything else in this report is downstream of Bench
being able to see its own violations.

---

### The one-paragraph version

Bench's machinery is better than its reputation-by-inspection would suggest, and
its problem is not that it built too much — it is that the most valuable thing it
built is switched off, and it has no way to notice. The audit's governing question
("demonstrable improvement, or accumulated ceremony?") splits: the **CLI and
enforcement layer** is demonstrable, tested, and I broke very little of it. The
**28-artifact prose layer** is high-craft but has **zero controlled evidence**, and
the only benchmark that exists on this question (SWE-Skills-Bench, 565 instances)
found 80% of skills produce zero pass-rate gain, an average of +1.2%, up to 451%
token overhead, and three skills that *degraded* performance. Bench's prose layer
is a portfolio of ~28 untested hypotheses costing 11,306 words for one spec-backed
feature against Pocock's 3,115 for the same pipeline. That is the honest verdict:
**the enforcement half has earned its place; the prose half has not been asked to.**

---

## B. Audit environment

```text
AUDIT MODE:               in-repository dogfood, cold start, disposable worktree
HARNESS:                  Claude Code
MODEL:                    claude-opus-5
EFFORT:                   xhigh
BENCH REPOSITORY PATH:    /home/mgibs/workspace/bench-audit-opus (linked worktree)
                          main tree: /home/mgibs/workspace/bench
BENCH COMMIT SHA:         58d966e2f92f7f37eba07b6215e8eef45371b72d
BENCH BRANCH:             audit/opus
GIT STATUS (initial):     " M .claude/settings.json"  "?? audit/"
GIT STATUS (final):       " M .claude/settings.json"  "?? audit/"  (+ gitignored dist/, .logs/)
AUTO-LOADED INSTRUCTIONS: CLAUDE.md → @AGENTS.md + @.bench/BENCH.md — 2,342 words
AUTO-LOADED SKILLS:       16 bench-craft-* + prototype (symlinked .claude/skills/)
                          6 model-invocable phase commands; 5 hidden by
                          disable-model-invocation. Description block: 1,259 words.
AVAILABLE TOOLS:          Bash, Read, Write, Edit, Skill, WebFetch, Artifact,
                          Agent (deliberately unused — see limitations)
OUTPUT DIRECTORY:         audit/opus-xhigh/   (per operator instruction)
```

**Safety boundaries observed.** Read-only inspection by default; local builds and
local test runs; one throwaway git repo under the session scratchpad for the
adoption experiment (created and removed); temporary in-tree fixtures
(`specs/zz-audit-probe/`, two `internal/toon/toon.go` mutations, one
`zz_audit_probe.txt`) each `cp`-backed and restored. No commits, no pushes, no
network mutations, no secret reads, no release commands (`bench prep-release`,
`bench release`, `bench release-preflight` were **not** run — they are
release-adjacent). Four public web sources were read (arXiv abstracts and one
GitHub README); no data left this machine.

**Final tree state verified:** `git status --porcelain` matches the initial state
plus the audit directory. The gate cache was reconciled to green
(`bench gate --fresh` → `green @ cf6f619`) before close.

### Limitations, stated plainly

1. **No controlled model trials.** The audit's own most important question —
   does `/bench-debug` reproduce the practitioner-observed effect of upstream
   `diagnosing-bugs`, and do the 28 prose artifacts change outcomes at all —
   requires repeated independent trials across harnesses. I could not run them:
   self-benchmarking is invalid (I already hold the discipline), and my operating
   instructions barred spawning delegates. Section W specifies the design; no
   claim in section H is labelled OBSERVED for the effect itself.
2. **Codex was not exercised.** Every cross-harness claim is derived from
   configuration files (`agents/openai.yaml`, `.codex/hooks.json`, the conformance
   check that enforces them), not from a Codex session.
3. **Upstream is pinned to local checkouts.** `mattpocock/skills` at
   `d574778f` (2026-07-08) and the kunchenguid repos at their 2026-06 commits.
   Upstream may have moved; every upstream claim names the commit.
4. **The 2026 research sources** (workspace, J-CoT, SWE-Skills-Bench, SWE-Effi,
   ToM-SWE, the J-Space repos) sit at or past my training boundary. I fetched
   four of them; the rest are marked RESEARCH and are not load-bearing.
5. **`internal/worktree` (13,006 LOC) was probed, not audited.** I exercised
   `list`, `path`, `clean` (plan), and read the lifecycle contracts. Its
   concurrency and lease semantics deserve their own pass.

---

## C. Bench reconstructed

Derived from implementation first, then compared to prose. This is what the tree
says, not what the docs say.

### C.1 The actual object

```
86,674 lines of Go   (40,242 non-test)   +   28 agent-facing prose artifacts (27,758 words)
                                         +   5 shell hooks   +   1 gate script
```

Largest packages, and what they are for:

| package | LOC | test LOC | what it is |
|---|---:|---:|---|
| `conformance` | 13,250 | 12,838 | 29 checks that grade the kit's own guidance/docs |
| `worktree` | 13,006 | 7,651 | pooled isolated checkouts, leases, landing, cleanup |
| `gate` | 5,678 | 1,111 | the oracle: subject identity, phases, verdict stores |
| `adopt` | 4,594 | 933 | `setup`/`link`/`upgrade`/`doctor` |
| `roadmap` | 3,336 | 1,889 | board parse, context snapshot, drain accounting |

**26 % of the codebase (conformance + worktree) implements no user-facing
engineering capability.** Both are infrastructure for *trusting* the work rather
than doing it. That is a legitimate bet — but see F and S for whether each pays.

### C.2 The real path from intent to verified completion

```
                       ┌──────────────────────────────────────────┐
   user opens session  │ SessionStart hook → bench session-inspect│
                       │   prints: CLI location, resume result,   │
                       │   ambient dashboard, guard brief         │
                       └──────────────────┬───────────────────────┘
                                          │  (silent + exit 127 if dist/bench absent — E1)
                                          ▼
  user must already know which of 11 phases applies  ◄── NO ROUTER (D)
                                          │
        ┌────────────┬────────────┬───────┴────┬─────────────┬──────────────┐
        ▼            ▼            ▼            ▼             ▼              ▼
  /bench-setup  /bench-shape  /bench-write  /bench-debug  /bench-implement  /bench-what-next
     -repo         -idea         -spec                       -spec           (maintenance)
        │            │             │            │              │
        │      decisions/<t>.md    │      repro loop first     │
        │      (bench maps)   specs/<slug>/     │        specs/<slug>/tickets/
        │                      spec.md          │        one delegate each
        │                    Status: staged     │        isolated worktrees
        │                    coverage map       │        serial green commits
        │                    ownership fences   │        on ONE retained source
        │                          │            │              │
        │                          └────────────┴──────────────┤
        │                                                      ▼
        │                                     /bench-review-implementation
        │                                     3 parallel read-only axis delegates
        │                                     Standards · Spec · Coverage
        │                                     findings → reviews/<slug>.md
        │                                                      │
        │                                                      ▼
        │                                          bench worktree land
        │                                     composes + gates + publishes
        │                                     Status: implemented, releases source
        │                                                      │
        │                                                      ▼
        │                                            /bench-final-check
        │                                     reports retained evidence,
        │                                     writes capture/retros/<slug>.md
        │                                                      │
        └──────────────────────────────────────────────────────┤
                                                               ▼
                                                    /bench-what-next
                                    drains IDEAS + learnings + retros → ROADMAP.md
                                    rewrites capture/session-handoff.md
                                    ONE batch commit on green
```

**The enforcement layer, drawn separately, because it is the part that is real:**

| control | boundary | posture when core is missing | actually enforces |
|---|---|---|---|
| `bench gate` | invariant 1 | n/a (is the core) | gofmt·vet·test·race·system·shellcheck |
| `bench commit` | landing | n/a | gates, then commits named paths atomically |
| Stop hook | Claude+Codex | **fail-open**, no forged verdict | can't stop on red while `BENCH_SHIFT=1` |
| `block-dangerous-git` | Bash PreToolUse | fail-closed on substring `git` | 20 destructive spellings, 1-level wrappers |
| `check-agent-line` | Agent PreToolUse (**Claude only**) | **fail-open** | delegate model ∈ `.bench/lines.env` |
| `pre-push` | git | n/a | default-branch push; `.bench` drift when pinned |
| conformance (29) | — | — | **nothing — not wired to the gate (E6)** |

### C.3 Where information is introduced, transformed, and lost

| stage | authoritative artifact | mechanically checked? | can go stale? |
|---|---|---|---|
| idea | `capture/IDEAS.md` | drain count in `bench status` | no (inbox empties to zero) |
| shaped | `decisions/<topic>.md` | `bench maps` (**advisory only**) | **yes — 1 invalid map, 6 days**|
| specified | `specs/<slug>/spec.md` | `bench coverage --check` (**prose-invoked**) | yes |
| ticketed | `specs/<slug>/tickets/*.md` | `Blocked by:` parsed; `Writes:` **advisory** | yes |
| implementing | worktree lease + `intent` ledger | `bench worktree list`, status `intent` row | no |
| locally verified | `$GIT_COMMON_DIR/bench-gate-evidence/<sha>` | content-addressed on (tree, oracle) | **no — correct** |
| *reported* verdict | `$GIT_DIR/bench-last-gate` | **red records never drift-checked (E5)** | **yes** |
| integrated | `bench worktree land` | composes + gates + publishes | no |
| reviewed | `reviews/<slug>.md` | presence = "fix work remains" | yes |
| promoted | spec `Status:` line → deleted | `bench spec retire` refuses unmerged | yes |
| resumed | `capture/session-handoff.md` | pin block regenerated by `bench handoff`; **`## State` prose is not** | **yes** |

**The single most important structural fact:** every row marked "advisory only" or
"prose-invoked" above would be mechanically enforced if the conformance suite ran
in the gate. `coverage-map-validation` and `decision-map-integrity` are both
*registered conformance checks*. They are enforced in `bench prep-release` and
nowhere else.
---

## D. Entry-point verdict

### D.1 What exists today (OBSERVED)

There is no front door. The README's "Reviewer quick start" is a menu of seven
slash commands plus two maintenance ones; two more (`/bench-assess`,
`/bench-deepen`) bring the total to **11**. Nothing observes state and routes.

I tested the candidate routers:

| candidate | what it actually is | verdict |
|---|---|---|
| `/bench-what-next` | frontmatter: *"Roadmap maintenance … **Maintenance, not a workflow phase**"*, `disable-model-invocation: true` | **not a router.** It drains capture into `ROADMAP.md`. Naming it "what next" is the single most misleading identifier in the repo. |
| `bench status` | 11 signals with a `→ action` column | **closest thing.** But see D.2. |
| `bench roadmap` | emits `sequence[3]{rank,text,command}` with harness-native commands | a router *output* over a human-maintained section; the CLI does no judgment (documented and deliberate) |
| `/bench-shape-idea` | has a "Starting from the roadmap" section that reads `ROADMAP.md` and offers top items | a de-facto cold entry **only if** the work needs shaping |

### D.2 Three demonstrated entry-point defects

**(1) An un-adopted repo reports all-clear.** In a fresh `git init` repo with no
`.bench/` at all:

```
$ bench status
bench: clean — nothing pending          [rc=0]
```

Identical to a fully-adopted, fully-clean repo. AXI principle 5 (definitive empty
states) fails at the highest-traffic surface. There is no route to `bench setup`.

**(2) A staged spec produces no signal.** I created a valid staged spec with a
passing coverage map and an authorized ownership fence, then ran `bench status`:

```
▶ commit on green  (git)
  git        3 dirty paths ...   worktree ...   structure ...   decisions ...   gate ...
```

No `specs` row. Reading `internal/status/status.go`, the `specs` signal fires on
exactly one condition — `%d merged spec(s) awaiting retirement`. **The workflow's
single most action-relevant state — "an approved spec is staged; run
`/bench-implement-spec`" — is invisible to the ambient dashboard.**

**(3) The action column is half-invocable.** Of 11 signals, roughly five carry a
command a harness can run (`/bench-what-next`, `/bench-shape-idea`, `bench link`,
`bench spec retire <slug>`, `bench gate --fresh`); the rest carry prose
("commit on green", "split (craft-seams)", "resume interrupted work", "promote or
delete by hand", "fix before commit"). A router cannot be built on a mixed column.

### D.3 The finding that decides this section

Pocock's suite ships **`ask-matt`** — `disable-model-invocation: true`, described
as *"Ask which skill or flow fits your situation. A router over the skills in this
repo."* Its body is a decision tree: main flow, on-ramps, codebase health,
vocabulary layer, crossing sessions, standalone, precondition. It exists precisely
because *"you don't remember every skill, so ask."*

**Bench inherited Pocock's pipeline and did not port his router**, while growing
from 10 engineering skills to 11 phases + 16 craft skills. The README's own
provenance table lists nine Pocock-derived rows and `ask-matt` is not among them.
This is the clearest single case of an upstream capability lost in adaptation.

### D.4 Recommendation: `/bench` — a state-driven router, not an alias

Do **not** copy `ask-matt`. Bench can do strictly better, because Bench *observes
state* and Pocock's router only reads prose. But the router must be a projection
over facts Bench already computes — not a new subsystem.

```
/bench
 └─ bench status --route          (new flag on an existing command)
      ├─ reads the 11 existing signals + one new `specs: staged` signal
      ├─ ranks by the existing severity ordinals already in the row structs
      └─ emits ONE row:  next[1]{state,why,command}
 └─ the phase file loads only the recommended phase's guidance
```

Why this over the alternatives:

- **Clearer than a menu** because it answers with one command, not eleven.
- **Robust** because every input is already derived from the tree — no new source
  of truth, no second reader (the standard `.bench/BENCH.md` already imposes).
- **Cheap**: `bench status` already computes all of it; the missing pieces are one
  signal (`specs: staged`), one flag, and normalizing the action column.

**How it discovers state:** the existing `status` signal set —
gate · git · intent · worktree · guards · drain · structure · decisions · specs ·
reviews · roadmap — plus `specs: staged`. **How it handles ambiguity:** it emits
the highest-severity actionable row and *names the runners-up in one line*; it
never asks a question the tree answers. **How it resumes:** the `intent` signal
already reports interrupted work, and `capture/session-handoff.md` carries the
`## Next command` line — the router should read that line rather than re-deriving.

**Claude Code:** `/bench` as a `.agents/commands/bench.md` phase, and it should be
**model-invocable** so a cold agent can orient itself. **Codex:** `$bench` adapter;
per the existing conformance rule it will be `allow_implicit_invocation: false`,
which is the right default there. **Setup stays separate:** `bench setup` is an
explicitly different act, is already excellent, and the router's job in an
un-adopted repo is exactly one row — `setup · this repo is not adopted · bench setup`.

**Expert entry points remain directly invokable**, unchanged: `/bench-shape-idea`,
`/bench-write-spec`, `/bench-debug`, `/bench-implement-spec`,
`/bench-review-implementation`, `/bench-final-check`, and the maintenance four.

### D.5 The twelve scenarios

| # | scenario | today | with `/bench` |
|---|---|---|---|
| 1 | new user, Bench installed | reads a README menu of 11 | one row: `bench setup` |
| 2 | existing user, no active work | `clean — nothing pending`, no route | roadmap sequence rank 1 |
| 3 | unshaped idea | must know `/bench-shape-idea` | routed via `decisions` signal |
| 4 | shaped, no spec | must know `/bench-write-spec` | `bench maps` state `ready` → routed |
| 5 | spec, no tickets | `/bench-implement-spec` bounces back to write-spec | routed to `/bench-write-spec` |
| 6 | ready ticket | **no signal exists** | new `specs: staged` signal |
| 7 | interrupted implementation | `intent` row: "resume interrupted work" (prose) | invocable resume |
| 8 | failing tests / bug | `gate red` row → "fix before commit" (prose) | `/bench-debug` |
| 9 | unverified changes | `git N dirty paths` → "commit on green" (prose) | `bench commit` |
| 10 | local work awaiting review | **no signal** | `/bench-review-implementation` |
| 11 | expert knows the phase | works today | unchanged |
| 12 | Claude → Codex switch | `$bench-*` adapters exist; handoff file carries `## Next command` | unchanged; router reads the same file |

Scenarios **6 and 10 have no signal at all today** — the two states in the middle
of the workflow, which is exactly where a router earns its keep.

---

## E. Provenance map

Upstream pinned: `mattpocock/skills` @ `d574778f94cf620fcc8ce741584093bc650a61d3`
(2026-07-08); `kunchenguid/{axi,no-mistakes,treehouse,firstmate,lavish-axi,tasks-axi,programbench-bench}`
at the commits in evidence log 01. Bench maintains its own provenance table in
`README.md:423`; I verified it against the tree rather than trusting it, and it is
**substantially accurate** — with one omission (`ask-matt`) and one row I would
re-classify (`git-guardrails`, below).

### E.1 Pocock-derived

| Bench artifact | upstream | classification | what changed |
|---|---|---|---|
| `/bench-write-spec` + `craft-spec` | `to-spec` | **adapted, materially extended** | keeps: no-interview synthesis, seams sketched + confirmed, highest/fewest seams, glossary + ADRs, no file paths (they rot), the Problem/Solution/Stories/Implementation/Testing/Out-of-scope template. Adds: acceptance coverage map with a **machine-checked schema**, edge inventory, **Won't handle**, hostile-input checklist, ownership fences, `Decision source:` entry contract, `Status:` lifecycle, verification log, per-story `Line:`. Replaces the issue tracker with in-repo `specs/`. |
| `craft-tickets` | `to-tickets` | **adapted, faithful** | keeps tracer-bullet vertical slices, demoable-alone, one-context-window sizing, `Blocked by:` edges, the granularity/edges/merge-split quiz, **expand–contract for wide refactors** (near-verbatim inheritance), no file paths. Adds: one ticket = one fresh write-delegate charge; advisory `Writes:` disjointness. |
| `craft-tdd` | `tdd` | **adapted, materially extended** | keeps: behavior through public interfaces, seams, pre-agreed seams only, the three anti-patterns (implementation-coupled, tautological, horizontal slicing), red-before-green, one slice at a time, refactor-is-not-in-the-loop. **Inherits both reference files** (`tests.md`, `mocking.md`) — rewritten, not copied (77→50 and 59→45 lines). Adds: the vacuity check, the compile-error-red rule, edge-class walk, acceptance-row red-signal classification, and the `ok`-hides-a-skip warning. |
| `/bench-implement-spec` | `implement` (15 lines) | **superseded** | upstream is a stub ("use /tdd, run typechecks, /code-review, commit"). Bench's is a different artifact. |
| `craft-review` + `/bench-review-implementation` | `code-review` | **adapted, improved** | keeps: parallel sub-agents so axes don't pollute each other, Standards + Spec separation, "do not merge or rerank", per-axis count + worst issue, the **Fowler smell baseline** (adopted 2026-07-08 per `skills-assessment.md`), repo-overrides-baseline, smells are judgment calls. **Adds a third axis (Coverage)**, the citation standard, `no-op`/`auto-fix`/`ask-user` dispositions, "re-derive then compare", the universal-claim-needs-an-enumeration rule, and **"refute before you report"**. |
| `/bench-debug` | `diagnosing-bugs` | **adapted — see section H** | |
| `craft-seams` | `codebase-design` | **adapted** | deep modules, interface-is-the-test-surface, deletion test, leverage/locality. Adds design-it-twice, dependency categories, read-budget-before-exploring, structure-gate splitting. |
| `craft-grill` | `grilling` / `grill-with-docs` | **adapted, improved** | frontier rounds, recommendation-per-question, look up facts don't ask them. Adds: the facts-vs-decisions split and the **confirmation gate before enacting** (both adopted as *bug fixes* from upstream v1.1 per `skills-assessment.md`), "close on a predicate, not a label", scoping-is-a-decision. |
| `craft-domain` | `domain-modeling` | **adapted** | canonical terms + Avoid lists, CONTEXT.md glossary-only. Adds producer-derived equivalence partitions and code-versus-claim comparison. |
| `craft-skills` | `writing-great-skills` | **adapted, extended** | leading words, progressive disclosure, contrastive pairs. Adds negation/negative-space/sediment failure modes and "write for the weakest reader". |
| `/bench-deepen` | `improve-codebase-architecture` | **adapted** | survey → candidates → grill. Adds HTML report and mid-tier delegate. |
| `/bench-setup-repo` | `setup-matt-pocock-skills` | **adapted** | keeps the two-half split; explicitly detects and reuses existing Pocock structure (`CONTEXT.md`, `docs/adr/`, `docs/agents/`). Adds gate + profile + lines. |
| `prototype` | `prototype` | **adapted** | rewritten around "name the question first"; deliverable is the verdict, not the code. |
| **— none —** | **`ask-matt`** | **NOT PORTED** | the router. See D.3. |

### E.2 Kunchenguid-derived

| Bench artifact | upstream | classification |
|---|---|---|
| `bench worktree` | `treehouse` (pooled worktrees) | **adapted, heavily extended** — leases, assignments, landing, plan-before-apply cleanup, fingerprints |
| `bench shift` (gated loop) | `gnhf` + `no-mistakes` | **adapted** — the differentiator is *gate-on-green, never self-graded* |
| Stop hook + `.bench/gate.sh` | `no-mistakes` (external gate before push) | **adapted** — Bench moves the gate from push-time to iteration-time |
| TOON output, `help[]` rows, exit codes, empty states | AXI's 10 principles | **adapted with a declared exemption registry** (`cmd/bench/main.go`: 8 `axiApprovedRoot`, the rest `axiExempt` with a named reason) |
| isolated validation worktree | `no-mistakes` | **adapted** |
| structured findings | `no-mistakes` | **adapted** → `reviews/<slug>.md` + dispositions |

**Divergence worth naming:** `no-mistakes` has **one entry point** —
`git push no-mistakes`, plus `/no-mistakes` — with the whole pipeline behind it.
Bench took the pipeline and the isolation and left the single entry behind. That is
the same omission as `ask-matt`, arriving from the other parent.

### E.3 Bench-native (no upstream)

The gate verdict keyed to `git write-tree`; **the conformance suite** (29 checks
enforcing guidance invariants); the harness × tier `lines.env` binding and its
Agent-tool guard; the acceptance coverage map's machine-checked row schema
(`bench coverage --check`); ownership fences + `bench preflight`; `bench worktree land`;
the roadmap/capture drain loop (`/bench-what-next` + `bench roadmap --context`);
`craft-line`'s diff-owned/inherited/spec-predicted red classification; the
`git-guardrails` guard (README credits Pocock's `misc/git-guardrails-claude-code`;
the Bench implementation is a Go tokenizer with wrapper scanning and a documented
threat model — I would classify it **Bench-native, inspired by**).
---

## F. Component inventory

Format: purpose · enforcement · observed value · cost · problem · recommendation.
"Observed" means I ran it during this audit.

### F.1 Enforcement and oracle

| component | enforcement | observed | problem | rec |
|---|---|---|---|---|
| `bench gate` | **is** invariant 1 | green 64 s / 5,917 B; red on a planted gofmt violation; cache invalidates on tracked **and** untracked-unignored change; 0 s reuse | does not run the 29 conformance checks | **strengthen** — add a conformance phase |
| gate evidence store (`bench-gate-evidence/<sha256(tree,oracle)>`) | content-addressed | correct in both directions | none found | **keep** |
| `bench-last-gate` (per-worktree) | read by `status` + `handoff` | **contradicts the oracle** (E5) | red records are never drift-checked | **fix (P0)** |
| `bench commit` | gates then commits atomically | not exercised (would commit) | — | keep |
| Stop hook | can't stop on red while `BENCH_SHIFT=1` | fail-open verified with core removed | fail-open is moot (shift needs the core) — defensible | keep |
| `block-dangerous-git` | 20 destructive spellings; 1-level wrappers | 43 probes | enumerates spellings not effects; gaps: `git stash drop/clear`, `git filter-branch`, `git rm -rf`, `eval`; degraded rim blocks on substring `git` (`echo legitimate` → BLOCKED) | **simplify the degraded rim** (match a `git` *token*, not a substring); add the 3 gaps |
| `check-agent-line` | delegate model ∈ `lines.env` | 10 probes: denies unbound, denies **omitted/empty** model, exact-string match, honest fork warning | Claude-only by documented Codex limitation; fan-out count unenforced | **keep** — best-in-repo control |
| `pre-push` | default branch + `.bench` drift | reported "not installed" in this worktree | worktrees don't inherit the hook | note |
| **conformance (29 checks)** | **none — skipped in the gate** | **10 red at HEAD** | see E6 | **P0** |

### F.2 Query / context CLI

| command | AXI | bytes | observed | rec |
|---|---|---:|---|---|
| `bench status` | exempt (operational) | 394 | 11 signals, `→ action` column | **strengthen**: add `specs: staged`; normalize actions to invocable commands; distinguish un-adopted from clean |
| `bench roadmap` | approved | 2,221 | 73 rows, board totals, `sequence[3]` with commands, drain counts | keep |
| `bench roadmap --context` | approved | 19,155 | schema-4 snapshot, per-source byte counts | keep — the drain's single inventory |
| `bench maps` | approved | 797, **rc=1** | caught a genuinely invalid decision map | **wire to the gate** (it already is a conformance check that doesn't run) |
| `bench coverage --check` | approved | ~180 | **4/4 hostile inputs caught** with exact messages | **wire to the gate** (same) |
| `bench preflight` | exempt (operational) | ~300 | refused empty fence; then caught an unauthorized path | keep |
| `bench diff` | approved | 759 | frozen-base review input | keep |
| `bench guards` | approved | 482 | deny surface + wiring per harness | keep |
| `bench learnings` | approved | 43 | definitive empty state | keep |
| `bench anchors` | approved | 201 | rc=2 with usage on missing arg | keep |
| `bench worktree list` | approved (child) | 194 | 2 foreign worktrees, landed/ignored counts | keep |
| `bench handoff` | exempt | ~2.5 KB | **regenerates the pin block from the tree**, preserves `## State` prose | keep; inherits the E5 bug |
| `bench structure` | exempt | 5,661, **rc=1** | 62 issues (50 four days ago); 4 rows of rot in its own accept-list | **simplify or gate** — see S |
| `bench outline` | exempt | **24,660** | 200 of 5,335 symbols; `outline_meta` honestly discloses `truncated=true`; path-scoping works (`outline internal/toon` → 18 rows) | **change the bare default** — refuse or aggregate; no ranking rule is stated for the 200 |
| `bench models` | exempt | 396 | advisory discovery, never validation | keep |
| `bench canary` | exempt | 42 | rc=1 on empty inventory | **decouple from the scaffolded gate** (D/§ adoption bug) |
| `bench commands --brief` | exempt | 31 | 3-verb liveness probe | keep |

### F.3 Adoption

`bench setup` (plan-first, converges, self-verifies with ok/red rows, names an
exact `next:`) is the **best-designed surface in the repo**. Its one defect is
downstream: it scaffolds no `.bench/gate-inputs.json` and a `gate.sh` whose canary
line cannot pass in a new repo. Both are one-line fixes.

### F.4 Prose artifacts — 28 hypotheses

| artifact | words | enforcement of its own rules | evidence it changes outcomes |
|---|---:|---|---|
| `AGENTS.md` + `.bench/BENCH.md` | 2,342 | conformance (**skipped**) + `guidance-prose-budgets` (**skipped**) | none |
| 11 phase commands | 12,111 | `docs-currency-workflow` (**skipped**) | none |
| 16 craft skills | 13,305 | `skills-index-command-adapters` (**skipped**) | none |
| 9 reference files | ~350 lines | pointer-fired | none |

Every one of these is well-written; none has been A/B tested; and the checks that
would keep them honest are the checks that do not run. Section W is the fix.

### F.5 State artifacts

| artifact | who writes | who reads | staleness control |
|---|---|---|---|
| `ROADMAP.md` + `roadmap/FT<n>.md` | `/bench-what-next` | `bench roadmap`, `/bench-shape-idea` | `roadmap-detail-integrity` (skipped) |
| `capture/IDEAS.md` | `bench idea` | drain | must empty to zero each run — good |
| `capture/learnings.md` | any phase | drain | open-entry count in `status` |
| `capture/retros/<slug>.md` | `/bench-final-check` | drain | pending count in `status` |
| `capture/session-handoff.md` | `bench handoff` + phases | cold sessions | pin block regenerated; **`## State` prose is not** |
| `decisions/<topic>.md` | `/bench-shape-idea` | `/bench-write-spec` | `bench maps` (advisory) |
| `specs/<slug>/` | `/bench-write-spec` | everything | **34 orphan tickets-only dirs, 43 files, 15 days, no owner** |
| `reviews/<slug>.md` | review phase | implement phase | presence = work remains |

---

## G. Pocock integration assessment

**What was preserved.** Every load-bearing procedure. I checked clause by clause
against the pinned upstream and found no upstream *invariant* silently dropped in
`to-spec`, `to-tickets`, `tdd`, `code-review`, `codebase-design`, `domain-modeling`,
or `grilling`. The seam vocabulary, the tracer-bullet sizing rule, the two-axis
separation and its "do not rerank" rule, the three TDD anti-patterns, the smell
baseline, expand–contract — all present, most sharpened.

**What was improved, and how.** Three moves, all of the same kind: *Bench converted
a prose rule into a checkable predicate.*

1. **Coverage map.** `to-spec` says "write user stories, extensively". Bench adds a
   row schema and `bench coverage --check`, which I demonstrated refuses an
   unreferenced story, a >4-story row, a two-predicate behavior, and a duplicate
   row id. Upstream has no equivalent; this is the single biggest upgrade.
2. **Review citation standard.** `code-review` asks for findings; `craft-review`
   states what a finding must cite, adds "refute before you report", and adds the
   rule that *a universal claim without an enumeration is a sample*. That converts
   reviewer confidence into a checkable artifact property.
3. **Delegate verification.** Upstream has no delegation discipline at all.
   `craft-delegate` requires the coordinator's independent probe to differ from the
   delegate's **in kind and in site** — "a second probe at the same site is
   vacuous, and a vacuous probe is indistinguishable from a pass."

**What drifted.** `ask-matt` was not ported (D.3). Upstream's `implement` closes
with `/code-review` *before committing*; Bench routes review before landing but
after serial ticket commits — a deliberate, documented change, not drift.

**What is duplicated.** Very little, and the repo polices it hard. `.claude/commands`
is a **symlink** to `.agents/commands` (one tracked entry — I confirmed by inode and
`git ls-files`), so the 11 phases have one source. The one real duplication I found:
the eleven `.agents/skills/bench-*/SKILL.md` Codex adapters each carry
`disable-model-invocation: true` (a **Claude** key) *and* `agents/openai.yaml`'s
`allow_implicit_invocation: false` (the Codex key) — two spellings of one policy,
one of which is inert on the harness that reads the file.

**What should remain prompt-level:** seam choice, story breadth, hypothesis
ranking, scope cuts, finding severity, tier selection. All require judgment over
a specific tree.

**What should become deterministic (and mostly already is, but unwired):**
coverage-map validity, decision-map integrity, doc/command structural contracts,
skills-index/adapter parity, prose budgets, roadmap detail integrity. **Six of the
29 conformance checks already implement exactly this and none of them runs.**

**Do the consistency gains survive across harnesses?** Partly, and the asymmetry is
mechanically caused:

- The gate **enforces** `allow_implicit_invocation: false` on every Codex adapter
  (`checkCodexCommandAdapters`, `skills_index_checks_test.go:229`). There is **no
  corresponding check on the Claude side**, and 6 of 11 Claude commands are
  model-invocable. So Bench's stated doctrine — *"workflow phases are
  reviewer-chosen entry points, not background generation guidance"* — holds by
  construction on Codex and by accident on Claude.
- `check-agent-line` is Claude-only, for a documented Codex limitation
  (`SubagentStart` can't deny). The harness-independent backstop is the shift
  adapters' refusal — which only covers headless runs.
- `/bench-debug` is model-invocable on Claude and **explicit-only on Codex**. See H.5.

---

## H. `diagnosing-bugs` assessment

This section carries the audit's highest burden of proof, so it is explicit about
what I observed and what I could not.

### H.1 The causal mechanism — which constraints are load-bearing

Reading upstream (134 lines, pinned `d574778f`) against `/bench-debug` (141 lines),
the constraints that plausibly break repair loops are not equal. My ranking, from
the structure of the failure they prevent:

| # | constraint | what it prevents | load-bearing? |
|---|---|---|---|
| 1 | **No hypothesis before a red-capable loop exists** | the entire symptom-patching mode | **yes — this is the skill** |
| 2 | **Red-capable = asserts the user's *exact* symptom**, not "doesn't crash" | fixing a neighbouring bug | **yes** |
| 3 | **One command, already run once, invocation + output pasted** | claiming a loop that doesn't exist | **yes** — it is the only *checkable* completion criterion in the skill |
| 4 | Deterministic / raise the reproduction rate | debugging against noise | yes |
| 5 | Fast (2 s) | the loop being too expensive to use | supporting |
| 6 | Minimise one element at a time, every remaining piece load-bearing | a hypothesis space too large to search | yes |
| 7 | 3–5 **ranked, falsifiable** hypotheses with stated predictions | anchoring on the first plausible idea | yes |
| 8 | One variable at a time; debugger over logs; tagged `[DEBUG-a4f2]` | uninterpretable probes | supporting |
| 9 | Regression test at a **correct** seam, or "no seam is the finding" | false confidence | yes |
| 10 | **Re-run the original un-minimised loop after the fix** | a green unit test over a still-broken path | **yes** |
| 11 | Cleanup + name the correct hypothesis in the commit | instrumentation rot; no learning | supporting |

Items 1, 2, 3, and 10 are the ones I would refuse to weaken. They share a shape:
**each forces contact with reality at a point where the model would otherwise
substitute confidence.**

### H.2 What Bench preserved — clause by clause (OBSERVED)

All eleven are present in `/bench-debug`. Verbatim-equivalent on the four
load-bearing ones:

- *"Before reading code to form a theory, build **one command** … that you have
  **already run once** (paste the invocation and its output)"* — 1 and 3.
- *"asserts the user's *exact* symptom … Not 'runs without erroring.'"* — 2.
- *"repro → failing test → fix → green → re-run the full Phase 1 loop"* and
  *"Done means: the Phase 1 loop no longer reproduces"* — 10.
- *"If you catch yourself theorizing before this command exists, stop — jumping to
  a hypothesis without a red loop is the exact failure this prevents."* — 1.

### H.3 What Bench added (and these are real)

1. **Isolation before the first artifact.** *"create or select its isolated
   worktree before the first repro artifact exists … deciding isolation after
   Phase 1 has already produced artifacts is too late."* Upstream has no worktree
   model; this closes an attribution hole Bench's own gate would otherwise inherit.
2. **The lookalike-surface rule.** *"drives the real bug path through the accused
   command — the surface the claim names, invoked as reported, never a lookalike
   (a raw `git add` is not `bench commit`)."* This is a genuine sharpening of
   constraint 2 and it recurs in `/bench-what-next`'s defect-entry rule. It is
   Bench-native and it is good.
3. **Seam ownership.** The regression-test seam and the edit owner may differ:
   fix once at the narrowest shared owner, keep the test at the highest seam where
   the failure is observable, and record tangled ownership as the architecture
   finding rather than patching every caller.
4. **Repro survives shift rollback.** The quarantine-marker mechanism (commit the
   repro in the project's expected-failure form so a red iteration's rollback
   doesn't destroy it; the fix's green commit removes the marker). Upstream has no
   loop to survive. This is thoughtful.
5. **Retired-spec recovery.** `bench spec history <slug>` — because Bench deletes
   specs on merge, a shipped feature's spec is in git history. Without this the
   debugger would conclude the behavior was never specified.

### H.4 Dilution risks (OBSERVED — the honest list)

Three compressions genuinely weaken the discipline, and one is severe.

1. **The loop-construction menu shrank from 10 options to 5.** Upstream enumerates
   failing test · curl · CLI+snapshot diff · **headless browser** · replay a
   captured trace · throwaway harness · **property/fuzz loop** · **bisection
   harness** · **differential loop** · **HITL script**. Bench keeps five and drops
   the five in bold. Constraint 1 says "spend disproportionate effort here"; the
   menu is *how* an agent spends it. For a UI bug, an intermittent-output bug, or
   a regression between two known-good states, Bench's agent now has no pointer.
   **This is the most consequential loss** and the cheapest to reverse: it is a
   reference file (`references/loop-constructions.md`) behind a pointer, ~15 lines,
   costing nothing until Phase 1 fires.
2. **The two hard stop-gates were softened.** Upstream: *"**No red-capable command,
   no Phase 2.**"* and *"Do not proceed until you have reproduced **and**
   minimised."* Bench keeps the theorizing-stop but drops both explicit gates, and
   drops the checkbox form (`- [ ]`) that made the completion criterion a thing the
   model must *emit*. Per Bench's own `craft-skills` — *"a vague criterion invites
   premature completion"* — this is a self-inflicted regression.
3. **"Tighten the loop" as an explicit iteration is gone.** Upstream treats the
   loop as a product with its own improvement pass (faster / sharper / more
   deterministic). Bench folds it into three adjectives on the first attempt,
   losing the "once you have *a* loop, tighten it" move.

Not dilution, but worth naming: upstream's *"Be aggressive. Be creative. Refuse to
give up."* and *"Spend disproportionate effort here"* are dropped. Bench's
`craft-skills` would call these no-ops, and I think Bench is right — they are
motivational, not actionable.

### H.5 Cross-harness behavior (OBSERVED from configuration)

| | Claude Code | Codex |
|---|---|---|
| upstream `diagnosing-bugs` available? | **no** — `.claude/settings.json` sets `"diagnosing-bugs": "off"` | unknown (not configured by Bench) |
| Bench `/bench-debug` model-invocable? | **yes** (no `disable-model-invocation` in `.agents/commands/bench-debug.md`) | **no** — `agents/openai.yaml`: `allow_implicit_invocation: false`, gate-enforced |

**The behavior with the strongest practitioner evidence in the whole system cannot
be self-invoked on Codex.** On Claude an agent that notices "this is a bug" can
enter the discipline; on Codex the reviewer must type `$bench-debug`. Given that
the observed effect is *breaking an agent out of a repair loop* — a state the agent
is in precisely because it is not reasoning well about its own process — requiring
the human to notice and intervene removes the mechanism at the moment it is needed.

### H.6 Removal burden — not met, and I am not asking for it

I have no controlled evidence that `/bench-debug` ≡ upstream. I did not run trials.
Per the audit's own rule the burden is on any recommendation to remove, dilute, or
replace — so my recommendation is **preserve and repair**:

- **P1** restore the loop-construction menu as a pointer-fired reference file;
- **P1** restore the two hard gates (`No red-capable command, no Phase 2`;
  `Do not proceed until reproduced and minimised`) and the checkbox completion form;
- **P1** restore "tighten the loop" as a named step;
- **P2** make `$bench-debug` implicitly invocable on Codex, or add a Codex-side
  trigger that fires the discipline on repeated same-class failures (see I.4);
- **P2** benchmark it (section W, arms D/E/F) before any further compression.

**Do not compress this skill further without the benchmark.** It is the one artifact
in the repo with repeated practitioner evidence, and section W exists to protect it.
---

## I. Run the Dang Test Findings

Nine patterns where Bench permits or causes speculation in place of a cheap
available observation. Each uses the required format.

### I.1 The oracle grades itself by a summary line

- **Current behavior** — `bench gate`'s test phase runs `go test -count=1 ./...`,
  prints `[test] ok github.com/gibbonmi/bench/internal/conformance 7.087s`, and
  reports **green**. The live-root conformance suite skipped.
- **Available observation** — `BENCH_CONFORMANCE_ROOT="$PWD" go test -count=1
  ./internal/conformance -run '^TestRootConformance$'` → **10 failures at HEAD**,
  in 3.2 seconds.
- **Why observation is superior** — a summary line reports *that a binary ran*,
  not *what it graded*. Bench's own `craft-tdd` says so: *"a run whose only
  matched test skipped still prints `ok` (Go does), so read the test output, not
  the summary."* The doctrine is correct and the oracle ignores it.
- **Current trigger** — none. `bench prep-release` only.
- **Recommended trigger** — every `bench gate` run.
- **Recommended enforcement** — a `conformance` phase in the built-in table
  (`internal/gate/phases.go`) that sets `BENCH_CONFORMANCE_ROOT` and
  `BENCH_CONFORMANCE_TIER=dev`; **and** an `environment`-class capability skip
  inside the oracle becomes **red**, not a footnote. A capability skip
  (`fifo`, `privilege`) is a fact about the host; an environment skip is a fact
  about the *invocation*, and it is always the caller's bug.
- **Measurement** — count of conformance diagnostics at HEAD (10 → 0), and days
  between a violation landing and its detection (currently ≥ 6, was 43 for the
  wiring itself).

### I.2 The gate's skip disclosure expands the harmless class and aggregates the dangerous one

- **Current behavior** — `capability-skips: 7 (capability=6 environment=1)`, then
  `class=fifo: 3`, `class=privilege: 3`. The single `environment` skip — the
  entire 29-check suite — is never named, never expanded, never attributed.
- **Available observation** — the runner already knows the skipping test's name
  and reason string (`"BENCH_CONFORMANCE_ROOT not set"`).
- **Why superior** — Bench's own `craft-gate` requires *"one distinct message per
  failure mode"* and *"a red gate must name its cause without archaeology"*. The
  same standard should bind skips: a skip is a check that did not happen.
- **Recommended trigger** — every gate run that skips anything.
- **Recommended enforcement** — print `class=environment: 1 (TestRootConformance:
  BENCH_CONFORMANCE_ROOT not set)`; red for the environment class.
- **Measurement** — `/bench-assess` on 2026-08-13 enumerated the skips
  ("five skips, including one FIFO and three privilege skips") and did not
  identify the fifth. A correct disclosure would have made that impossible.

### I.3 Two verdict readers, one of which never checks drift

- **Current behavior** — same tree, same second, no writes between calls:
  ```
  bench gate    : gate: green (fresh verdict reused for this tree)
  bench status  : gate       red        → fix before commit
  bench handoff : Gate: red at `a353ae6` — current
  ```
- **Available observation** — the reproduction above (three commands, ~2 minutes).
- **Root cause, located** — `internal/gate/verdict.go:inspectSubjectAt` returns on
  `if rec.Status != "green"` **before** it evaluates `if rec.Tree != s.Tree`, so a
  red record is never drift-classified; and
  `internal/status/status.go:GateVerdict` defines
  `nonReusableGreen := ... && in.Status == "green" && !in.ReusableGreen;
  gi.Stale = nonReusableGreen` — structurally excluding red records from staleness.
- **Why superior** — the ambient dashboard is the cold-session orientation
  surface, and `bench handoff` is the cross-session resume artifact. Both now
  assert a red the oracle denies. The direction is fail-safe (never a false green)
  but it costs a wasted 64 s gate run at best and a phantom debugging session at
  worst. Note the *good* design that propagated it: `internal/handoff/facts.go`
  deliberately does not re-derive staleness — *"a second derivation of it is how
  the two surfaces come to disagree"* — so single-sourcing faithfully copied a
  wrong rule to every consumer. Single-sourcing amplifies correctness **and**
  defect.
- **Recommended enforcement** — evaluate tree/oracle drift before status in
  `inspectSubjectAt`; define `Stale` on the record, not on greenness; add the
  reproduction as a regression test.
- **Measurement** — the three-surface agreement assertion, run in the gate.

### I.4 Nothing detects a repair loop

- **Current behavior** — `craft-line`'s ladder counts *diff-owned reds* and
  escalates a tier after two at the same tier. That is a **spend** response. There
  is no trigger that says *stop patching and enter the diagnosis discipline*.
  `/bench-debug` is reviewer-invoked ("This phase is reviewer-invoked") and
  explicit-only on Codex.
- **Available observation** — the gate already records every verdict with its
  tree, oracle, status, and timestamp. Two consecutive reds whose failing phase
  and first failing check are identical is a computable predicate over records the
  system already writes.
- **Why superior** — the practitioner evidence for `diagnosing-bugs` is precisely
  that it *breaks agents out of repair loops*. Detecting the loop is the trigger
  the mechanism is missing; today it depends on a human noticing.
- **Recommended trigger** — second consecutive red with the same first-failing
  check and no new evidence.
- **Recommended enforcement** — `bench gate` emits one row:
  `repair_loop[1]{check,attempts,action}` → `/bench-debug`. Advisory first;
  measure before making it blocking.
- **Measurement** — repair rounds per ticket are **already collected**:
  `/bench-final-check`'s retro template mandates a repair-attribution table with a
  fixed cause vocabulary, and `/bench-what-next` tallies it. The last drain
  recorded *8 tickets, 6 one-shots, 2 with repairs (delegate-error 1, spec-row 1,
  other 2)*. This is a real measurement channel Bench already built and does not
  yet feed back.

### I.5 Advisory reds nobody is required to run

- **Current behavior** — `bench maps` rc=1 on
  `decisions/spec-build-review-gate-cadence.md`, which names
  `internal/specbuild/checkpoint.go`; `internal/specbuild` was deleted
  2026-08-11 (`dae240df`). `bench structure` rc=1 with 62 issues. The gate is
  green. The prior `/bench-assess` found the map defect on 2026-08-13 (M-O1) and
  it is still red.
- **Available observation** — both commands, ~1 second each.
- **Why superior** — a stale cross-reference is a *fact about the tree*, not a
  judgment. `decision-map-integrity` is already a registered conformance check.
- **Recommended enforcement** — wire conformance into the gate (I.1) and
  `decision-map-integrity` comes with it, free.
- **Measurement** — days-to-detection for a dangling `Sources:` path.

### I.6 Nothing routes anyone to the residue it deliberately refuses to clean

- **Current behavior** — 34 `specs/<slug>/` directories hold only `tickets/`
  (43 tracked files, oldest 2026-08-02). `bench spec retire <slug>` recognises the
  state precisely — *"incomplete retired spec folder … remove by hand after
  reviewing its residue; retire will not auto-clean it"* — a good, deliberate
  refusal. But `bench status`'s `specs` signal only counts `specs/*/spec.md`
  awaiting retirement, and `bench roadmap --context` reports `specs/,parsed,0`.
- **Available observation** — `find specs -type d -mindepth 1 -maxdepth 1` minus
  those with a `spec.md`. One command.
- **Why superior** — the refusal is correct; the *absence of a surface that ever
  states the debt* means the correct refusal is never acted on. 15 days of
  accumulation is the evidence.
- **Recommended enforcement** — a `specs: N tickets-only folder(s)` status row
  with action `review and remove by hand`.
- **Measurement** — count over time; it should not grow.

### I.7 The adoption path ships a gate that cannot go green

- **Current behavior** — reproduced end to end in a throwaway repo:
  `bench setup --yes` → `next: configure .bench/gate.sh - replace the
  BENCH_SENTINEL sentinel`. Remove the sentinel (the instruction) → gate still red
  with `.bench/bin/bench.sh: line 14: HOME: unbound variable` and
  `canary inventory validation failed`. Declare `HOME`/`BENCH_HOME` in a
  hand-written `gate-inputs.json` → still red: `canary fixture inventory is empty`.
- **Available observation** — `bench setup --yes && bench gate` in a temp repo.
  Ninety seconds.
- **Why superior** — the whole architecture rests on "the gate is the oracle". An
  adopter whose gate cannot go green has no working Bench, and neither
  `bench setup`'s output nor `/bench-setup-repo` names either cause.
  `/bench-setup-repo` §3 says *"if it errors for a reason other than real failing
  checks, fix the wiring"* — delegating a known, deterministic, kit-caused
  breakage to model judgment.
- **Recommended enforcement** — (a) `bench setup` scaffolds
  `.bench/gate-inputs.json` declaring `HOME`/`BENCH_HOME`; (b) `bin/bench.sh:14`
  becomes `${BENCH_HOME:-${HOME:?bench: HOME or BENCH_HOME must be set}/.bench}` so
  the failure names itself; (c) the scaffolded `gate.sh` emits the canary line
  **commented out** with a TODO — exactly what `/bench-setup-repo` §2 already tells
  humans to do with not-yet-real contracts; (d) an adoption smoke test in the gate:
  temp repo → `setup` → configure a trivial real check → gate green.
- **Measurement** — the smoke test itself.

### I.8 The degraded git guard blocks reading and permits deleting

- **Current behavior** — with `dist/bench` absent, the guard's fail-closed rim
  matches the literal substring `git` in the raw envelope. Observed:
  `cat .gitignore` → BLOCKED, `ls …/.git` → BLOCKED, `echo legitimate` → BLOCKED;
  `rm -rf .g""it` → ALLOWED.
- **Available observation** — the six probes above.
- **Why superior** — the rim's stated purpose is "cannot classify, so refuse the
  destructive surface". A substring is not the destructive surface. In this audit
  it cost three turns before I could read `git status`.
- **Recommended enforcement** — tokenize the command and match a bare `git` word,
  or reuse the existing Go tokenizer via a pure-shell fallback list. Keep
  fail-closed; narrow the predicate.
- **Measurement** — false-positive rate over a corpus of read-only commands.

### I.9 A missing core is silent at the one moment it matters

- **Current behavior** — at session open with no `dist/bench`, the SessionStart
  hook prints one line and exits 127; the ambient dashboard — the mechanism
  designed to orient a cold session — is absent, and nothing says why. This was
  the literal state of both this worktree **and the maintainer's main tree** at the
  start of this audit. Recovery is `go build -o dist/bench ./cmd/bench`, **1.4 s**,
  documented only inside the handoff file that isn't loaded.
- **Available observation** — running the hook with the binary moved aside.
- **Recommended enforcement** — the hook already resolves the wrapper; on a 127 it
  should print the one-line build command rather than exiting silently.
- **Measurement** — cold-start-to-first-useful-command turns.

### I.10 Proposed doctrine — short enough to stay active

Two sentences, for `.bench/BENCH.md` beside the four invariants. It states a
predicate, not an exhortation, because `craft-skills` is right that exhortations
are no-ops:

> **Empirical escape.** When a bounded, safe command can distinguish the
> hypotheses you are holding, run it before writing another sentence about them —
> and prefer the command that can come back *red*.
> A retry inherits the failed command, the observed failure, what you ruled out,
> and the one variable you changed; a retry that inherits none of these is not a
> retry.

Do not build a skill for this. Bench's problem is not that agents lack the idea —
`/bench-debug`, `craft-review`'s "refute before you report", `craft-gate`'s "prove
it bites", `craft-tdd`'s vacuity check, and `/bench-what-next`'s
repro-through-the-accused-command rule all already say it. The problem is that the
**oracle** does not follow it. Fix I.1 and I.3; the doctrine is the reminder, not
the mechanism.

### I.11 The auditor's own metrics, honestly

| metric | this audit |
|---|---|
| Observation Opportunity Delay | 0 steps in 8 of 9 investigations; **1 instance** of ~4 reasoning steps between "the two verdicts disagree" and running the tight three-surface experiment — I listed hypotheses about a second store before testing, and one of them (`dist/` exclusion) was falsified by a 10-second command I should have run first |
| Speculation After Discriminator | 1 (above) |
| Blank Retry Rate | 0/4 — each retry changed a named variable (`dist/` present→absent; substring→live core; sentinel→gate-inputs→canary) |
| Failure Inheritance Rate | 4/4 |
| Red-Loop Acquisition | conformance gap: 1 command. Verdict disagreement: 3 commands after 1 falsified hypothesis. Adoption chain: 3 commands. |
| Root-Cause Before Fix Rate | n/a — this audit proposes no fixes it applied |
| Corrections made | 2: I inferred `.claude/commands` were hardlinks (they are a symlinked directory) and I attributed early "stale" readings to a defect when they were self-inflicted tree churn. Both were caught by running a command instead of trusting the inference. |

---

## J. Prose audit

### J.1 The measurements, which mostly refute the priors

| probe | result |
|---|---|
| vague verbs (`carefully`, `thoroughly`, `appropriately`, `robustly`, `properly`) across 27,758 words | **0** |
| `relevant` / `sufficient` | 3 / 1 |
| `complete` / `clean` / `safe` | 22 / 14 / 6, nearly all with observable referents ("complete vertical path", "clean review", "safe model-id tokens") |
| negative-instruction density | 0.7 per 100 words in `AGENTS.md` and `.bench/BENCH.md`; 0.6–1.4 across skills; 185 total |
| always-loaded weight | 2,342 words + 1,259-word description block ≈ **4.8 k tokens** |
| `guidance-prose-budgets` conformance check | **passes** — Bench enforces per-file line budgets on its own guidance |
| duplication of the 11 phase files across harnesses | none — `.claude/commands` is one tracked symlink |

**The audit brief's hypotheses of vague-verb padding, negative-instruction
overload, and context dumping are not supported.** The always-loaded footprint is
lean and progressive disclosure genuinely works. This prose is better disciplined
than most production codebases' comments.

### J.2 What is actually wrong with the prose

**Volume, not quality.** One spec-backed feature loads:

```
always-loaded                                    3,601 words
/bench-write-spec + craft-{spec,seams,domain,tickets,tdd}   4,993
/bench-implement-spec + craft-{line,delegate}    2,662
/bench-review-implementation + craft-review      2,510
/bench-final-check                               1,141
                                                ───────
                                                14,907 words  ≈ 20 k tokens
Pocock's equivalent pipeline (to-spec, to-tickets, implement, tdd, code-review)
                                                 3,115 words  ≈  4 k tokens
```

**≈ 3.6× the guidance weight for the same pipeline shape.** SWE-Effi's token-snowball
finding and SWE-Skills-Bench's "modest savings to a 451% increase while pass rates
remain unchanged" both bear directly here, and Bench has no measurement of its own.
`roadmap/FT100` — *cut prose weight to demonstrated-delta clauses* — is already
rank 1 on the board. **It should not be executed by editorial judgment.** Every
clause deleted without a measurement is a coin flip in both directions. Run W first.

### J.3 Specific findings

| finding | class | recommendation |
|---|---|---|
| `/bench-what-next` is named like a router and is roadmap maintenance | **contradiction with user expectation** | **rename** to `/bench-drain` (or `/bench-roadmap`); free the name for the router |
| `.claude/README.md` says `.claude/skills/` and `.claude/commands/` are "adapter symlinks" — skills are, commands is a symlinked *directory* | stale detail | shorten |
| `.claude/README.md` justifies Codex-only phase skills because "a same-named skill duplicates the slash-menu entry" — in current Claude Code the 6 non-disabled commands appear as model-invocable skills anyway | stale rationale | rewrite the paragraph around invocation *policy*, not menu duplication |
| `.bench/BENCH-reference.md:106` carries the removed token `spec build` | stale | delete (conformance already catches it — see I.1) |
| `.agents/commands/bench-implement-spec.md` lacks `## Entry orientation` / `## Exit handoff` | **structural contract violation** | restore the headings |
| `.agents/commands/bench-final-check.md` "dropped whole-file implementation-retro replacement"; `.bench/BENCH.md` "dropped the implementation-retro drain owner" | contract violations | restore |
| eleven Codex adapters carry both `disable-model-invocation` (Claude key) and `allow_implicit_invocation: false` (Codex key) | duplicated policy, one inert | keep one; note the other in a comment |
| `disable-model-invocation` — a Claude-specific key — lives in the harness-neutral `.agents/commands/` | **model-specific contamination of the neutral layer** | acceptable in practice; document it as an adapter key that the neutral layer carries, and **add the missing Claude-side parity check** |
| `capture/session-handoff.md`'s `## State` says "`specs/` is empty" (34 tickets-only folders) and "about to land as one commit" (it landed) | **plan recorded as state** | the pin block is regenerated correctly by `bench handoff`; the prose is not. Constrain `## State` to assertions the tree can contradict, or shorten it to the three facts a resumer acts on |

**Hidden policy — behavior that exists only in prose:**

- **Effort** has no enforcement surface anywhere (`craft-line` says so plainly —
  an honest disclosure, and still hidden policy).
- **Fan-out count** is declared, never checked.
- **Iteration cap** is declared, never checked.
- **"Run `bench coverage --check`"** is a `/bench-write-spec` instruction; the
  mechanical check exists (`docs-currency-workflow`) and doesn't run.
- **The inline-edit allowance** ("exactly one source-line insertion, deletion, or
  replacement… spans the current reviewer request") is a precise, checkable
  predicate enforced by nothing.

---

## K. Work-state assessment

**Verdict: reject a new work-state layer; consolidate the five faces that exist.**

The J-Space report's `Goal / Core / Verified / Open / Next` ledger is a reasonable
shape, and its own authors concede the evidence is *"single runs without confidence
intervals … not unified-harness experiments … engineering interpretation of
black-box behavior."* Bench already has every field, distributed across artifacts
that are individually better than a YAML block would be:

| field | Bench's existing owner | can a fresh model use it? |
|---|---|---|
| **Goal** | spec `## Problem` / `## Solution`; or `capture/session-handoff.md` `## State` | yes for a spec; **prose-dependent** for the handoff |
| **Core** | the four invariants + `## Ownership fences` + `Blocked by:` | yes |
| **Verified** | the gate evidence record `(tree, oracle, status, recorded_at)` + red-then-green coverage rows | **yes, and it is genuinely evidence-referenced** — better than a self-reported list |
| **Open** | decision-map `Not yet specified` + open decision tickets, surfaced by `bench maps` | yes |
| **Next** | `bench status` action column + `bench roadmap` `sequence` + handoff `## Next command` | **three sources, unreconciled** — this is the gap |
| *attempts* | retro repair-attribution table (fixed cause vocabulary) | yes, but only post-hoc |
| *coverage* | `bench coverage <spec>` `rows[N]{story,behavior,seam}` | yes |
| *feedback loop* | `/bench-debug` Phase 1 command | **prose only — not stored anywhere** |

Answering the audit's questions directly:

- *Can a new model identify the objective without conversation history?* **Yes**
  with a live spec; **partly** otherwise — the handoff `## State` is prose that can
  drift, as it has at HEAD.
- *Are verified statements backed by current evidence references?* **Yes** — the
  tree-hash verdict is the best answer to this question I have seen in a harness.
- *Can open questions survive a reset without becoming accidental facts?* **Yes** —
  `bench maps` refuses readiness while any answer is a placeholder.
- *Is the next action executable and bounded?* **Not reliably** — six of eleven
  status actions are prose, and three surfaces answer the question independently.

**Recommendation: consolidate, do not adopt.** The single change worth making is
D.4's `bench status --route`: one `next[1]{state,why,command}` row derived from the
signals that already exist. That collapses three "Next" sources into one projection
without creating a second source of truth. The only genuinely missing field is
**feedback loop** — for a bug, the repro command should be stored (a `Repro:` line
in the handoff, or in the ticket) so a resuming session inherits the loop rather
than rebuilding it. That is P2, and it is small.

---

## L. Context compiler assessment

Bench already compiles context, and does it well: 4.8 k tokens always-loaded,
everything else pointer-fired. The `outline_meta` truncation disclosure, the
`bench roadmap --context --row <ids>` targeted fetch, and `craft-seams`'
declare-a-read-budget-before-exploring rule are all real progressive-disclosure
mechanisms.

**Smallest sufficient context for one implementation ticket**, derived from what
`/bench-implement-spec` + `craft-delegate` actually require a delegate to receive:

```
objective ...... the ticket's `What to build` (one paragraph)
constraints .... the spec's `## Ownership fences` (exact paths)
requirements ... this ticket's acceptance-coverage rows (story, behavior, seam,
                 why-it-catches) — NOT the whole coverage map
prior art ...... the exemplar file to mirror + the fixture/seam inventory
verification ... the focused check command + the named central-property mutation
line ........... resolved model id / effort / iteration cap
base ........... worktree root + expected tip + stale-base check
                                                    ────── ~600–900 words
```

It does **not** need: the whole spec, the roadmap, other tickets, `ASSESSMENT.md`,
historical reasoning, or the 16 craft skills. `craft-delegate` says exactly this —
*"Prefer compressed inputs — the named decision source, exact passages, coverage
rows, and the fence's fixture-and-seam inventory — so the delegate uses prior art
instead of re-deriving it."* The design is right.

**Where measurement is missing.** Bench has no instrumentation for total context
loaded, duplicate reads, or re-reads after compaction. It *does* have the
`/bench-final-check` retro's "Agent-experience improvements → Bench CLI / Skills /
Process" headings, which is a qualitative channel for the same thing. Section W
adds the quantitative one. This is the honest limit of section L: I can tell you
the design is sound and I cannot tell you what it costs in practice, because
nothing measures it.

**One concrete waste I did observe:** `bench outline` with no path emits 24,660
bytes (~7 k tokens) to deliver 200 of 5,335 symbols under no stated ranking, and
the rows are dominated by `tests/canary/**` fixtures and one `.webp`. Any agent
that runs the discoverable bare form pays 7 k tokens for a near-arbitrary sample.
Path-scoping works (`bench outline internal/toon` → 18 rows, ~600 bytes), so the
fix is to make the bare form refuse or aggregate.
---

## M. Claim/evidence assessment

**Bench has no claim system, and that is the right architecture.** There is no
`bench claim`, no claim store, no provenance labels, no replay command. A
repository-wide search for claim machinery returns nothing outside unrelated
prose. The audit brief's doctrine — *CLI establishes observation → agent supplies
interpretation → claim records epistemic state → gate enforces required evidence →
replay detects staleness* — describes a system Bench deliberately does not build.

What Bench built instead is **one content-addressed fact**:

```
subject  = (policy, identityRoot, TreeHash(root), gate resolution, PATH, manifest identity,
            + hashed contents of every declared path/tool/launcher in .bench/gate-inputs.json)
evidence = $GIT_COMMON_DIR/bench-gate-evidence/sha256(tree ‖ oracle)
           → {schema, state, status, tree, oracle, recorded_at}
```

I attacked it and it held:

| attack | result |
|---|---|
| stale source — modify a tracked file, re-run | full re-run (no reuse) ✓ |
| untracked file added (`audit/`, `zz_probe.txt`) | tree hash changes, full re-run ✓ |
| gitignored artifact added/removed (`dist/bench`) | tree hash unchanged ✓ (correct: `git add -A` respects `.gitignore`) |
| revert to a previously-green tree | green reused, 0 s ✓ |
| plant a real red, then revert | red recorded at its own tree; green correctly reused at the reverted tree ✓ |
| unsupported assertion / evidence without interpretation | not expressible — there is no assertion channel to abuse ✓ |
| **is the reported verdict the same as the authoritative one?** | **NO — E5/I.3** ✗ |
| claim-dependency invalidation | n/a — no dependency graph exists, and none is needed |
| direct file reads bypassing the CLI | possible and harmless; the verdict is not a file an agent would forge |
| freshness window | verdicts expire (`now.Sub(recorded) >= freshness`) ✓ |

**What behavior changes because each state exists** — the audit's "monitoring
without control" test:

| state | what it blocks or permits | control? |
|---|---|---|
| gate **green** for this exact tree | `bench commit` lands; the Stop hook allows stopping; `bench worktree land` composes | **real control** |
| gate **red** | commit refuses; Stop blocks while `BENCH_SHIFT=1` | **real control** |
| gate **stale** (green, drifted tree) | reuse refused, gate re-runs | **real control** |
| gate **red, drifted tree** | **nothing** — and it is reported as current | **defect (I.3)** |
| `bench structure` red (62) | **nothing** | monitoring without control |
| `bench maps` red (1) | **nothing** | monitoring without control |
| conformance red (10) | **nothing** | monitoring without control — and unobservable |
| `reviews/<slug>.md` present | means "fix work remains"; deleted by the fix commit | real, though prose-enforced |
| spec `Status: staged` | `bench spec retire` refuses; `bench worktree land` publishes the flip | real control |
| open decision ticket | `bench maps` refuses map readiness | **real control** |

Three of the four red-capable advisory surfaces control nothing. That is the
concrete meaning of "the gate is green while the kit is red."

---

## N. Gate / Review / Final Check / CI matrix

| | purpose | inputs | coverage | timing | owner | failure effect | bypass | independent? |
|---|---|---|---|---|---|---|---|---|
| **focused checks** | fast feedback inside a ticket | the delegate's own tests | its slice | during a ticket | delegate | none | trivially | no |
| **mutation probe** | is the test vacuous? | delegate's finished work | one central property | ticket close | delegate, then coordinator **at a different site and kind** | finding | prose-enforced | **yes — by construction** |
| **`bench gate`** | the oracle | whole tree | gofmt·vet·test·race·system·shellcheck | every commit, every shift iteration | Go core | commit refused | `--fresh` is an escape *upward*; no downward escape found | n/a (deterministic) |
| **conformance (29)** | guidance/doc invariants | whole tree | docs, skills index, maps, roadmap, coverage maps, prose budgets, routes, AXI registry | **`bench prep-release` only** | conformance pkg | none in dev | **skipped by default** | n/a |
| **`bench preflight`** | phase-entry agreement | spec + tree + explicit base | base currency, authorized paths, row ownership/membership, non-empty diff | before build; before review | Go core | phase stops | prose-invoked | partially |
| **`/bench-review-implementation`** | semantics the gate can't see | frozen base+tip diff, spec, standards docs, coverage rows, hostile-input list | 3 axes in parallel fresh contexts | after implementation | 3 read-only delegates | findings → repair tickets | advisory by design | **yes** — see below |
| **`bench worktree land`** | compose + gate + publish + release | reviewed base/tip pair | whole tree | integration | Go core | landing refused | — | n/a |
| **`/bench-final-check`** | report the applicable oracle | retained or fresh verdict | whole tree | last | reviewer | — | it *reports*, never re-grades a landed spec | no — deliberately |
| **CI** | release evidence | `.github/workflows/{release,native-runtime}.yml` | cross-platform, release | on release | GitHub | release blocked | — | yes |

**Useful redundancy vs bad redundancy.** Bench is unusually good here:

- Gate vs review defends **distinct failure modes** (regression/rule vs
  wrong-thing/wrong-way/uncovered-edge) — genuinely useful.
- The three review axes run in **parallel fresh contexts** specifically so one
  axis's derivation cannot seed another's, and `craft-review` requires each to
  **re-derive from its primary source before comparing the candidate** — that is
  real independence engineering, not "another model looks at it".
- The coordinator's mutation probe must differ from the delegate's **in kind and
  in site**, with the stated reason: *"a second probe at the same site is vacuous,
  and a vacuous probe is indistinguishable from a pass."* This is the best
  anti-bad-redundancy clause I found in the repo.

**Where independence is weaker than claimed.** The Spec axis receives the spec, and
the Standards axis receives the standards docs — but the diff arrives via
`bench diff --full`, which includes the commit log. Commit messages carry the
implementer's rationale, so the axes are partially anchored by the implementer's
own account of what they did. `craft-review` explicitly guards the *spec* case
(*"quote the applicable spec line rather than trusting what a ticket or commit
message claims was built"*) but nothing withholds the narrative. **A cheap
experiment (section W, arm G) would settle whether withholding the log changes
finding counts.**

**Adversarial cases, tested where possible:**

| case | result |
|---|---|
| passes tests, violates a requirement | possible — the gate has no requirement axis; that is *why* review exists |
| passes gate, fails conformance | **demonstrated** — 10 checks red at HEAD, gate green |
| passes gate, fails final-check | not reachable — final-check reports the gate |
| review finding the gate cannot detect | by design; dispositions route them |
| stale evidence | greens: correctly invalidated. reds: **not** (I.3) |
| partial coverage | `bench coverage --check` catches an unreferenced story ✓ (**not gated**) |
| local success claimed as global completion | `bench worktree land` gates the composed tree — good; but a ticket delegate's "focused checks green" is explicitly *not* accepted as done, which is correct |

---

## O. CLI / AXI assessment

Scored against AXI's 10 principles (`kunchenguid/axi` @ `6908df20`) by execution,
not doctrine.

| # | principle | Bench | evidence |
|---|---|---|---|
| 1 | token-efficient output (TOON) | **strong** | `internal/toon`; typed headers `rows[N]{cols}` on every query surface |
| 2 | minimal default schemas | **strong** | `learnings[0]{date,title}`, `guards[4]{...}`, `worktrees[2]{...}` — 2–8 fields |
| 3 | content truncation + `--full` | **mixed** | `bench diff --full` ✓; `bench test --full` ✓; `bench outline` truncates and *discloses* (`omitted_symbols=5135, truncated=true`) but offers **no** `--full` or paging — only path scoping |
| 4 | precomputed aggregates | **strong** | `board[1]{rows_shown,rows_total,sequence_trusted}`, `drain[1]{...}`, `outline_meta[1]{...}` |
| 5 | definitive empty states | **strong except one** | `learnings[0]{...}:` ✓, `help[0]{cmd,why}:` ✓ — but `bench status` in an **un-adopted repo** prints `clean — nothing pending`, indistinguishable from adopted-and-clean ✗ |
| 6 | structured errors + exit codes | **strong** | 0 success / 1 domain error / 2 usage, consistently across 24 probes; `toon.Errorf(problem, remedy)` pairs every error with a remedy |
| 7 | ambient context | **strong** | SessionStart hook prints the dashboard + guard brief (when the core exists) |
| 8 | content-first (no-arg shows data) | **strong except one** | `status`, `roadmap`, `maps`, `guards`, `learnings`, `diff` all show live data bare; `bench outline` shows 24 KB of near-arbitrary data — content-first taken past its useful limit |
| 9 | contextual disclosure (next steps) | **strong** | `help[N]{cmd,why}` rows; `bench maps` → `/bench-shape-idea`; `bench spec retire` → the exact remedy |
| 10 | consistent help | **strong** | `bench <verb> --help` uniform; `bench help` is the inventory; `bench nosuchverb` → rc=2 with the exact message |

**Bench goes beyond AXI in one respect worth naming:** it maintains a *declared
conformance registry* (`cmd/bench/main.go`) where 8 commands are `axiApprovedRoot`
and every other is `axiExempt` with a typed reason (`operational`, `mutation`,
`plumbing`, `release`). Honest scoping instead of blanket claims — and there is a
conformance check for it (`axi-query-registry`), which of course does not run.

**Measured agent cost of representative output:**

| question an agent asks | command | bytes | round trips |
|---|---|---:|---:|
| what needs attention? | `bench status` | 394 | 1 |
| what should I work on? | `bench roadmap` | 2,221 | 1 |
| what is the guard surface? | `bench guards` | 482 | 1 |
| what changed vs base? | `bench diff` | 759 | 1 |
| is this spec covered? | `bench coverage <spec>` | ~250 | 1 |
| where are the seams? | `bench outline` | **24,660** | 1 (+ archaeology) |
| what is my worktree state? | `bench worktree list` | 194 | 1 |
| is the tree green? | `bench gate` (cached) | ~120 | 1 |

Seven of eight are excellent. The eighth is the kill-list candidate.

| Current behavior | Agent cost | Failure mode | Proposed owner | Proposed change |
|---|---:|---|---|---|
| `bench outline` bare emits 200/5,335 symbols | ~7 k tok | near-arbitrary sample; help admits it "does not identify blessed seams" | CLI | bare form emits `outline_meta` + per-directory symbol counts only; require a path or `--full` for rows |
| `bench status` can't distinguish un-adopted from clean | 1 turn + a wrong mental model | false all-clear | CLI | `setup` signal when `.bench/` is absent |
| `bench status` has no staged-spec signal | the whole middle of the workflow is invisible | agent must guess the phase | CLI | add `specs: staged` |
| `bench status` action column half prose | router impossible | mixed | CLI | every action is an invocable command or empty |
| `bench structure` 62 issues, rc=1, ungated | permanent nag; +12 in 4 days | alarm fatigue | CLI + gate | gate it at a frozen budget, or drop it from `status` to an on-demand query |
| `bench-last-gate` read by 2 consumers, red never drift-checked | wasted 64 s gate runs; phantom reds | contradicts the oracle | Go core | I.3 |
| `.logs/gate-*.jsonl` accumulate (15 in one session) | disk only | none | Go core | prune on run |
| `bench-gate-evidence/` has no prune path | disk only (209 B each) | none | Go core | prune by age |

---

## P. Orchestration assessment

**The orchestrator's role is well-drawn and not bloated.** `craft-delegate` states
it precisely: *"The coordinator scopes, routes, and verifies work; a write-delegate
authors code."* The inline allowance is a hard number — one source-line insertion,
deletion, or replacement, spanning the whole reviewer request and not resetting
across tasks or slices. That is unusually disciplined; most harnesses leave this to
vibes.

**What each delegation boundary buys** — the audit's "weak answer / strong answer"
test:

| delegation | what it buys | verdict |
|---|---|---|
| write-delegate per ticket | isolated ownership (own worktree, own fence), bounded context, parallelism on disjoint `Writes:` | **strong** — isolated ownership |
| 3 review axes in parallel | **independent rediscovery** — fresh contexts so one axis can't seed another, each re-deriving from its primary source | **strong** — independent cold review |
| Research decision tickets | bounded research, off the critical path | **strong** |
| the `--reviewer <tier>` round in write-spec | separate verification of spec+tickets | **strong** |
| `/bench-deepen`'s mid-tier read-only survey | bounded research | fine |

I found **no ceremonial delegate**. Every boundary has a named non-"another model
looks at it" justification.

**Broadcast state — what is projected once vs rediscovered.** Bench gets this
right in principle:

- *Broadcast once:* the four invariants, the ownership fences, the coverage rows,
  the resolved line, the exemplar file, the fixture/seam inventory, the expected
  tip. All are facts.
- *Independently rediscovered:* implementation correctness, requirement
  satisfaction, absence of bugs, the review verdict. `craft-review`'s "re-derive,
  then compare" is the enforcing clause.
- *Never broadcast as fact:* a delegate's done-claim. `craft-delegate`: *"A
  delegate's done-claim is a claim, not a result"*, with a six-item verification
  list (run the gate, check every row red-then-green, `git status` the worktree,
  resolve every identifier in an absence claim, probe one behavior independently,
  spot-check citations).

**Where orchestration is unenforced.** Everything above is prose. Nothing checks
that the coordinator ran the six verifications, that fan-out matched the
declaration, that the probe site differed, or that a delegate stayed inside its
fence *while working* (`bench preflight` checks authorized paths, but at phase
entry, not continuously). `roadmap/FT213` — *"a read-only delegate reading a graded
tree gets its own worktree, and a delegate's claim about a gate signal gets an
oracle-verified probe"* — is exactly this gap, is rank 3 on the board, and was
*reproduced twice* in a retro. The system knows.

**Failure inheritance.** A blocked delegate's return shape is specified precisely —
*"repro command, red output digest, the failing surface it observed, and its
in-fence dirty paths"* — which is a genuinely good structured-retry state, better
than the "task + previous attempt failed" the audit warns about. It is prose-only,
and there is no mechanical carrier for it: it lives in the coordinator's context,
so it does not survive a compaction or a model swap. That is the P2 gap in K.

---

## Q. Model/harness boundaries

| layer | contents | boundary quality |
|---|---|---|
| **Bench core** | the Go binary; the gate, verdict, worktree, spec, coverage, roadmap, conformance packages; `.bench/gate.sh`; the four invariants; the craft skills' engineering semantics | **clean** — nothing model-specific found in `internal/` |
| **Claude Code adapter** | `.claude/settings.json` (hooks, `skillOverrides`), `.claude/skills/` symlinks, `.claude/commands` symlink, `check-agent-line --harness claude`, `TodoWrite` seeding from `bench coverage`, `BENCH_CLAUDE_{TOP,MID,CHEAP}` | **mostly clean**, one leak: `disable-model-invocation` is a Claude key living in the neutral `.agents/commands/*.md` frontmatter |
| **Codex adapter** | `.codex/hooks.json`, `.agents/skills/bench-*/agents/openai.yaml`, `BENCH_CODEX_*`, `.bench/adapters/codex` | **clean**, and the only one with a **gate-enforced** invocation policy |
| **Model-specific profile** | `.bench/lines.env` opaque tier tokens; `projects/<name>.md` `Lines` prose | **exemplary** — tier is the only shared identity, no family is canonical, ids are opaque strings compared exactly, and an unbound harness column *refuses to run* rather than borrowing |

The tier-binding design is the cleanest model-independence mechanism I have seen in
a harness: it refuses to guess. `bench models` is explicitly *advisory discovery,
never validation*, and the guard *"does not ask a provider whether an id exists."*

**The two real boundary defects:**

1. **Invocation-policy parity is enforced on one side only.** `checkCodexCommandAdapters`
   requires `allow_implicit_invocation: false` on all 11 Codex adapters. There is
   no Claude-side check, and 6 of 11 Claude commands are model-invocable. Same
   doctrine, one harness held to it. **Add the mirror check** — then decide
   deliberately which phases are model-invocable *on both*.
2. **`check-agent-line` is Claude-only** for a documented, re-checkable Codex
   limitation. The note even carries a re-check trigger (*"Re-check if the Codex
   changelog adds a spawn tool name or a deny-capable SubagentStart"*) — good
   practice, but nothing schedules the re-check.

**Can a fresh Codex session resume Claude's work, and vice versa?** Structurally
yes: `capture/session-handoff.md` carries a machine-regenerated pin block
(repository, path, branch, HEAD, dirty count, unpushed commits, spec status, gate
verdict) plus `## Next command`, and `bench handoff --harness codex` rewrites the
next command in that harness's form. The tree is authoritative where the two
disagree. **Two things degrade it:** the `## State` prose is model-authored and
drifts (it is wrong at HEAD), and the gate line it pins can be a phantom red (I.3).
Both are fixable; the mechanism is sound.
---

## R. Start-over decision

**Recommendation: (A) continue incrementally, with one bounded consolidation
inside it. Not (B), not (C).**

| | risk to demonstrated behavior | migration cost | duplicate truth | testability | harness compat | rollback | time to value | old-vs-new benchmarkable |
|---|---|---|---|---|---|---|---|---|
| **A. incremental** | **none** | ~2 weeks of the P0/P1 list | reduces (E5, verdict readers) | unchanged (86 k LOC of tests) | unchanged | per-commit | **days** | yes — same tree |
| **B. strangler spine** | moderate — the spine duplicates `bench status` + gate + handoff during migration | months | **increases during migration** | new spine untested | must re-derive both adapters | branch-level | months | hard — two systems |
| **C. rewrite** | **severe** — 29 conformance checks, 20 guard spellings, the tree-hash verdict, the lease machinery | quarters | n/a | starts at zero | rebuild | none | quarters | no |

The candidate spine the brief proposes —
`init → canonical entry → deterministic status → work state → context compilation
→ selected practice skill → executable feedback → evidence → gate → checkpoint` —
is a **description of what Bench already is**, with two gaps: "canonical entry"
(D.4) and "work state" (K, which I recommend rejecting as a new layer). Building a
spine to reach a design you already have is how you lose the parts you can't see.

**The one consolidation worth doing inside A:** collapse the two gate-verdict
readers to one (I.3) and make `bench status` the single projection of "what is
true and what is next" (D.4). That is ~200 lines of Go and one new status signal —
not a spine.

**What a rewrite would accidentally destroy**, concretely: the `git write-tree`
verdict key (which correctly includes untracked-unignored files and correctly
excludes gitignored ones — subtle and right); `git restore --source=` being in the
guard's deny surface while `--staged` is allowed; the `check-agent-line` denial of
an *omitted* model field; `bench spec retire`'s refusal to auto-clean an
incomplete folder; `bench worktree clean`'s plan-then-apply fingerprint;
the "coordinator's probe differs in kind **and site**" rule; the 29 conformance
checks. Every one of these is a scar from a real failure, and none is legible from
the outside.

**Do not rewrite because the system is hard to explain.** It is hard to explain
because it has 73 roadmap rows, not because it is incoherent. The architecture I
reconstructed cold in section C is coherent.

---

## S. Complexity kill list

| # | target | evidence | verdict |
|---|---|---|---|
| 1 | `bench outline` bare invocation | 24,660 B for 200/5,335 symbols, no ranking rule, help admits it doesn't find blessed seams | **simplify** — bare form emits meta + per-dir counts; rows need a path or `--full` |
| 2 | 34 orphan `specs/*/tickets/` dirs (43 files, 15 days) | invisible to `status` and `roadmap --context`; `spec retire` correctly refuses but nothing routes there | **delete the residue; add the signal** |
| 3 | `bench structure` as a permanent `status` row | 62 issues, ungated, +12 in 4 days, 4 rows of rot in its own accept-list | **move** off `status` to an on-demand query, **or** gate it at a frozen budget. A permanent unenforced 62 is alarm fatigue |
| 4 | the second gate-verdict reader (`bench-last-gate`) | contradicts the oracle (I.3) | **merge** into one reader |
| 5 | `disable-model-invocation` on the 11 Codex adapter SKILL.md files | inert on Codex; duplicates `allow_implicit_invocation` | **delete** one of the two |
| 6 | `/bench-deepen` | maintenance survey; overlaps `/bench-assess` (which also surfaces architecture candidates) and `craft-seams`' deletion test; produces an HTML report | **benchmark first** — I found no evidence either way, and it is one of two maintenance phases that duplicate a scoping step |
| 7 | `ASSESSMENT.md` + `skills-assessment.md` + `COMPLIANCE_ASSESSMENT.md` + `capture/architecture-review-*.html` (≈80 KB of assessment artifacts in-tree) | four overlapping self-assessments; the newest (4 days old) missed the biggest defect | **move to reference docs** or prune to one; they are not agent-facing and they age |
| 8 | `/bench-what-next`'s name | it is roadmap maintenance, explicitly "not a workflow phase" | **rename** — the name is the entry-point confusion |
| 9 | `capture/session-handoff.md` `## State` free prose | wrong on 3 facts at HEAD | **shorten** to what the tree can contradict |
| 10 | prose weight: 14,907 words per spec-backed feature vs Pocock's 3,115 | no measurement either way | **benchmark first (FT100)** — do not cut by editorial judgment |

**Agentless as a hostile baseline.** Agentless is `localize → repair → validate`.
Bench adds: an external oracle that is not the model; isolated ownership so
parallel work cannot collide; a coverage map that makes requirement loss
detectable; independent multi-axis review; a resumable handoff; and mechanical
guard rails on destructive actions. Each of those defends a failure mode
localize-repair-validate does not address — **for multi-session, multi-agent,
reviewer-supervised work**. For a single-shot bug fix on a well-tested repo,
Agentless would remove *almost everything Bench adds*, and Bench's own
lighter-path table already concedes the point (one independently-green ticket
crossing no declared seam skips the spec phase entirely). That table is Bench's
Agentless escape hatch, and it is the right shape.

---

## T. Missing capabilities (only after simplification), ranked by leverage

1. **Conformance in the gate** (P0) — turns 29 written invariants into enforced ones.
2. **A router** (P0/P1) — D.4. Without it every other improvement is behind a menu.
3. **Repair-loop detection** (P1) — I.4. The trigger the highest-evidence mechanism
   is missing, computable from records already written.
4. **A/B measurement harness** (P1) — W. Without it, FT100's prose cut and every
   future skill edit are coin flips.
5. **`specs: staged` status signal** (P1) — the invisible middle of the workflow.
6. **Adoption smoke test in the gate** (P1) — I.7; the onboarding path is broken
   and nothing would have told anyone.
7. **Claude-side invocation-policy parity check** (P2) — Q.
8. **Stored feedback-loop command for bugs** (P2) — K; so a resuming session
   inherits the repro instead of rebuilding it.

---

## U. Recommended architecture

The smallest coherent future Bench is the current Bench with **one addition, one
merge, and one wiring**:

```
                    ┌───────────────────────────────────────────────┐
   ADD  ────────────│  /bench   →  bench status --route             │
                    │  one row: next[1]{state,why,command}          │
                    │  projection over signals that already exist   │
                    └──────────────────┬────────────────────────────┘
                                       │
   ┌───────────────────────────────────┴───────────────────────────────────┐
   │  expert phases, unchanged and directly invokable                      │
   │  shape · write-spec · debug · implement · review · final-check        │
   │  maintenance: drain (was what-next) · assess · setup-repo · update-kit│
   └───────────────────────────────────┬───────────────────────────────────┘
                                       │
   ┌───────────────────────────────────┴───────────────────────────────────┐
   │  THE ORACLE                                                           │
   │  gofmt · vet · test · race · system · shellcheck                      │
   │  + conformance (29)   ◄── WIRE                                        │
   │  environment-class capability skip = RED                              │
   │  keyed on sha256(TreeHash ‖ oracle);  ONE reader  ◄── MERGE           │
   └───────────────────────────────────┬───────────────────────────────────┘
                                       │
   ┌───────────────────────────────────┴───────────────────────────────────┐
   │  guards: git · agent-line (both harnesses when Codex allows) · stop   │
   │  isolation: pooled worktrees, leases, fences, plan-then-apply cleanup │
   │  state: spec Status · coverage rows · decision maps · handoff pin     │
   └───────────────────────────────────────────────────────────────────────┘
```

Nothing else is added. No spine, no work-state layer, no new skill, no new agent
role.

---

## V. Prioritized roadmap

### P0 — the oracle must see the kit

**P0-1. Wire root conformance into `bench gate`; make an environment skip red.**
*Problem:* 29 checks, 13,250 LOC, unrun since 2026-07-05; 10 red at HEAD; a
doctrine commit violated a structural contract 6 days ago and landed green.
*Evidence:* E6 / I.1. *Change:* add a `conformance` phase to
`internal/gate/phases.go` setting `BENCH_CONFORMANCE_ROOT` + `BENCH_CONFORMANCE_TIER=dev`;
classify an `environment` capability skip inside the oracle as red; expand the skip
line to name the test and reason. *Benefit:* six advisory surfaces become enforced
at once. *Cost:* +3.2 s per gate run (measured); a one-time fix of 10 diagnostics.
*Dependency:* none. *Proof:* `bench gate` red at HEAD before the 10 fixes, green
after; a planted heading deletion in any `.agents/commands/*.md` reds the gate.

**P0-2. One verdict reader; drift-check reds.**
*Problem:* `bench gate` green while `bench status` says red and `bench handoff`
writes `Gate: red at <old-tree> — current` into the cold-start artifact.
*Evidence:* I.3, with the two single-line causes located. *Change:* in
`inspectSubjectAt`, evaluate tree/oracle drift **before** the status branch; in
`GateVerdict`, define `Stale` on the record rather than on greenness. *Benefit:*
the ambient dashboard and the resume artifact stop contradicting the oracle.
*Cost:* ~20 lines + a regression test. *Dependency:* none. *Proof:* the
three-surface agreement assertion from I.3 as a test.

### P1 — the front door and the feedback

**P1-1. `/bench` router + `bench status --route`.** *Change:* new `specs: staged`
signal; normalize the action column to invocable commands or empty; a `setup`
signal when `.bench/` is absent; `--route` emits one `next[1]{state,why,command}`;
a thin model-invocable `.agents/commands/bench.md` + `$bench` adapter; rename
`/bench-what-next` → `/bench-drain`. *Evidence:* D. *Proof:* the 12 scenarios in
D.5, each asserted as a fixture: given a tree state, `--route` emits exactly one
expected command.

**P1-2. Restore `/bench-debug`'s three compressions.** The loop-construction menu
as `references/loop-constructions.md` behind a pointer; the two hard gates
(`No red-capable command, no Phase 2`; `Do not proceed until reproduced and
minimised`) and the checkbox completion form; "tighten the loop" as a named step.
*Evidence:* H.4. *Cost:* ~20 lines in a reference file + 4 lines in the phase.
*Proof:* the diff against upstream shows each restored constraint; then arm F of W.

**P1-3. Adoption smoke test.** `bench setup` scaffolds `.bench/gate-inputs.json`;
`bin/bench.sh:14` names its own failure on unbound `HOME`; the scaffolded
`gate.sh` emits the canary line commented with a TODO; a gate fixture drives
temp-repo → setup → trivial real check → **green**. *Evidence:* I.7.
*Proof:* the fixture.

**P1-4. Repair-loop tripwire (advisory).** `bench gate` emits
`repair_loop[1]{check,attempts,action}` on the second consecutive red with the same
first-failing check. *Evidence:* I.4. *Proof:* fixture with two identical reds;
then measure against the retro repair-attribution tally already being collected.

**P1-5. The measurement harness.** Section W, arms A/B/E/F at minimum.
*Dependency:* blocks FT100.

### P2

- Claude-side invocation-policy parity check (Q).
- `specs: N tickets-only folder(s)` signal; delete the 34 orphans.
- Guard gaps: `git stash drop|clear`, `git filter-branch`, `git rm -rf`; tokenize
  the degraded rim so `echo legitimate` stops being blocked.
- `bench outline` bare form → meta + per-directory counts.
- Store the repro command in the handoff/ticket so a resume inherits the loop.
- Prune `.logs/` and `bench-gate-evidence/`.

### P3

- Codex `check-agent-line` equivalent when the platform allows; schedule the
  documented re-check.
- Consolidate the four in-tree assessment artifacts.

### EXPERIMENT

- W arms A/B/C (does Bench beat instructions-only?), D/E (does `/bench-debug`
  carry its weight?), F (upstream vs Bench-integrated), G (does withholding the
  commit log change review findings?).

### DELETE

- `disable-model-invocation` from the 11 Codex adapter SKILL.md files.
- The 34 orphan ticket directories (after promoting anything durable).
- `bench structure` from the permanent `status` row set (keep the command).

**One caution about ordering.** `roadmap/FT100` (cut prose weight) is currently
rank 1 on the board. **It should not run before P1-5.** Cutting 14,907 words of
guidance by editorial judgment, with no measurement, in a system whose own
`craft-skills` skill warns that *"every decision a skill declines to make is
silently delegated to the model's priors"*, is the highest-variance action
available. Measure, then cut.

---

## W. Benchmark plan

**Controls held fixed:** repository commit, task set, model, effort, harness, tool
permissions, wall-clock budget, token budget, starting context. **Trials:** ≥5 per
arm per task; report medians and the spread, never a single trajectory.

**Arms**

| arm | configuration |
|---|---|
| A | model + `AGENTS.md` only — no Bench phases, no craft skills, no gate |
| B | current Bench, full pipeline |
| C | proposed Bench (P0 + P1 applied) |
| D | current Bench **without** `/bench-debug` (bug tasks only) |
| E | current Bench **with** `/bench-debug` (bug tasks only) |
| F | upstream `diagnosing-bugs` vs Bench `/bench-debug`, same tasks |
| G | review with the commit log withheld vs supplied (independence probe, N.) |

**Tasks — use real ones, not toys.** Bench has an unusually good corpus available:
`git log --diff-filter=D -- specs/` lists every retired spec, and
`capture/retros/` (drained, but recoverable from history) records what each build
cost. Reconstruct 10–15 tasks by checking out the parent of each landing commit
and using the retired spec as the requirement. Bug tasks come from `roadmap/FT*.md`
`Occurrence:` lines, which already name reproduced defects with dates.

**Measures**

- *Outcome:* task success; requirement satisfaction against the retired spec's
  coverage rows; hidden defects found by a held-out reviewer; review findings; CI;
  first-fix success.
- *Epistemic:* unsupported assertions; stale evidence used; incorrect completion
  claims; claim/evidence mismatch.
- *Context:* total tokens; files read; duplicate reads; irrelevant reads; repeated
  searches; context loaded before it was relevant.
- *Execution:* tool calls; duplicate commands; **observation-opportunity delay**;
  speculation-after-discriminator; blank-retry rate; time-to-red-capable-loop;
  time-to-root-cause.
- *Recovery:* checkpoint completeness; cross-session resume; cross-model resume;
  repeated dead ends; failure-inheritance rate.
- *Orchestration:* subagent count; duplicated discovery; overlapping changes;
  integration conflicts; orchestrator tokens.
- *Complexity:* Bench operations required; latency; new failure modes; human
  interventions.

**Predefined success criteria (set before running):**

- **B beats A** on requirement satisfaction and epistemic quality by a margin
  larger than the trial spread, at **≤ 3× A's tokens**. If B wins on correctness
  but costs > 3×, that is a *finding*, not a pass — SWE-Effi's exact warning.
- **E beats D** on first-fix success and time-to-root-cause on bug tasks. If it
  does not, the practitioner evidence does not transfer to this repo/model and the
  compression debate reopens on evidence.
- **F**: if upstream ≥ Bench-integrated, the H.4 compressions cost something real;
  restore more.
- **C beats B** on tokens at equal or better outcome. That is the only honest
  justification for FT100.

**Per-skill ablation.** SWE-Skills-Bench (arXiv 2603.15401, 565 instances) found
39 of 49 skills produced **zero** pass-rate improvement, mean **+1.2 %**, token
overhead from modest savings to **+451 %**, and **three skills that degraded
performance by up to 10 %**. Bench has 28 prose artifacts and zero ablation data.
Run leave-one-out on the 16 craft skills over the same task set. Expect most to be
neutral; the value of the experiment is finding the two or three that are negative.

---

## X. Research reconciliation

| source | claim | Bench support | Bench contradiction | transfer limit | unproven |
|---|---|---|---|---|---|
| **Workspace / J-space** (transformer-circuits 2026; **fetched**) | verbalizable representations occupy 6–10 % of activation variance; information can be present yet inactive; ablating top J-space vectors impairs multi-hop reasoning while leaving recall intact | progressive disclosure (4.8 k always-loaded, rest pointer-fired) and `craft-delegate`'s "compressed inputs" align with *making the currently-relevant constraint small and explicit* | none | **large.** This is about neural representations, not text files. It does **not** license a work-state YAML, and the paper's mechanism (a limited internal workspace) is not the same thing as a prompt budget. I use it only as motivation, never as evidence | whether *any* text-level protocol changes J-space occupancy |
| **J-CoT** (2607.21981; not fetched) | carry selective intermediate state rather than all prior prose or a dense hidden state | Bench's ticket charges carry rows + fence + tip, not the transcript | none | Bench needs no latent reasoning; the transfer is analogical only | RESEARCH |
| **DeepSeek-V4 J-Space report** (**fetched**) | "capability realization loss" across reasoning mode / interface / tool schema / active representations / long-term state / verification; proposes Goal-Core-Verified-Open-Next, checkpointing, empirical escape, `fast/full/loop` gating | Bench independently has all five fields (K) and an empirical-escape norm | Bench's `Verified` is **stronger** — an evidence reference, not a self-report | **authors concede** single runs, no CIs, no unified harness, black-box interpretation | everything; treat as a design sketch |
| **J-Space Suite** | ledger + tiered gating | — | Bench already has tiers via `craft-line` | a second tier axis would **conflict** with `lines.env` | — |
| **ReAct** (2210.03629) | reason → act → observe → update | the gate is a forced observe step every iteration; `/bench-debug` Phase 1 forbids reasoning before observation | the four unenforced advisory surfaces permit reason→reason→conclude | none material | — |
| **Reflexion** (2303.11366) | feedback across attempts; episodic memory | retro repair-attribution table with a fixed cause vocabulary; `capture/learnings.md` reviewed drain | verbose reflective prose *is* preferred over structured state in `## State` | Bench's structured form (cause vocabulary) is arguably better than reflective text | does the drain change future outcomes? |
| **SWE-agent** (2405.15793) | the agent-computer interface is a first-class design object | Bench's CLI is designed for agents (TOON, typed schemas, `help[]` rows, exit codes, a declared AXI registry) | `bench outline` bare form is human-shaped; `bench status`' prose action column | — | measured turn/token savings |
| **Agentless** (2407.01489) | localize → repair → validate can match agentic scaffolds | Bench's lighter-path table is exactly this escape hatch | the full pipeline is 14,907 words of guidance for one feature | Agentless targets single-shot issue fixing, not supervised multi-session work | does B beat A? (W) |
| **Lost in the Middle** (2307.03172) | mid-context information is used unreliably | small always-loaded core; pointer-fired references; targeted `--row` fetches | none observed | — | Bench measures no context metrics |
| **SWE-Skills-Bench** (2603.15401; **fetched**) | 39/49 skills → zero pass-rate gain; mean +1.2 %; token overhead to +451 %; 3 skills degraded up to −10 %; only narrowly-tailored domain-specific skills gained (to +30 %) | Bench's *specialized* artifacts (`craft-gate` for this oracle, `craft-delegate`'s isolation, `/bench-debug`) are in the class that gained | **direct challenge to the other ~25** | benchmark used different skills, repos, models | **the central unproven question in this audit** |
| **SWE-Effi** (2509.09853; **fetched**) | effectiveness = accuracy ÷ resources; token snowball; expensive failures; scaffold quality matters less than scaffold-model integration | `craft-line`'s iteration cap and "stop and report, never grind" target expensive failures directly | 3.6× Pocock's guidance weight, unmeasured | — | Bench's cost-per-outcome |
| **ToM-SWE** (2510.21903; not fetched) | persistent user intent, stateful interaction | `projects/<name>.md`, `CONTEXT.md`, decision maps, `Lines` cache | — | **do not** import a second user-model agent | RESEARCH |
| **Coconut** (2412.06769) | latent reasoning | — | — | not applicable | — |

**The distinction the brief demands, stated plainly:**
`neural J-space (a measured 6–10 % of activation variance) ≠ J-CoT latent reasoning
(a training-time architecture) ≠ a text-level agent-control protocol (what Bench
is)`. Bench should adopt **none** of the J-Space vocabulary. Its useful residue is
one already-held principle: *the constraint that matters now should be small,
explicit, and evidence-backed* — which Bench implements better than the source
does, because its `Verified` field is a content-addressed hash rather than a
sentence.

---

## Y. Ten hardest truths

1. **Bench's largest package does not run in Bench's oracle.** 13,250 LOC,
   29 checks, skipped for 43 days on a missing environment variable; 10 red at HEAD.
2. **A Bench build violated a Bench contract and landed green six days ago.**
   `fa4e1f02` removed `## Entry orientation` and `## Exit handoff` from
   `bench-implement-spec.md`; the check for it exists and did not run.
3. **`bench gate` and `bench status` disagree about the same tree in the same
   second**, because a red verdict is never drift-checked. The ambient dashboard —
   the cold-session orientation surface — contradicts the oracle, and
   `bench handoff` writes the contradiction into the resume artifact.
4. **Bench's own adversarial self-assessment ran four days ago, enumerated the
   gate's skips, and walked past the one that disables 29 checks** — because the
   gate aggregates the dangerous skip class and expands the harmless ones.
5. **A newly adopted repo's gate cannot go green**, for two reasons the adopter did
   not cause and no Bench surface names. The whole architecture rests on the gate.
6. **Bench inherited Pocock's pipeline and dropped his router.** `ask-matt` exists
   upstream for exactly this reason; Bench grew to 11 phases without it, and
   `/bench-what-next` — which is named like the router — is roadmap maintenance
   and is hidden from the model.
7. **The workflow's middle is invisible.** A valid staged spec produces no
   `bench status` row; so does work awaiting review. The two states where a router
   would earn its keep are the two states with no signal.
8. **The prose is excellent and unmeasured.** Zero vague verbs in 27,758 words,
   0.7 negations per 100 words, enforced line budgets — and 3.6× Pocock's guidance
   weight for the same pipeline, with zero A/B evidence, against a benchmark that
   found 80 % of skills produce no gain and some produce harm.
9. **Single-sourcing propagated a defect perfectly.** `internal/handoff` refuses to
   re-derive staleness precisely so the surfaces cannot disagree — and therefore
   copies the wrong rule to every consumer. One source per fact makes correctness
   and incorrectness equally viral.
10. **The behavior with the strongest evidence has the weakest trigger.**
    `/bench-debug` preserves every load-bearing constraint of `diagnosing-bugs`,
    upstream is switched off in this repo, and on Codex the skill cannot be
    self-invoked — so the mechanism that breaks repair loops requires a human to
    notice the loop.

---

## Z. One next ticket

# Grade the kit with the kit: wire root conformance into `bench gate`

**Blocked by:** none.
**Writes:** `internal/gate/phases.go`, `internal/gate/capability_skips.go`,
`internal/conformance/gate_entry_test.go`, plus the 10 files the diagnostics name.

## What to build

`bench gate` runs the 29 registered conformance checks against the live root as a
first-class phase, and an `environment`-class capability skip inside the oracle is
**red** rather than a line in a footer. A guidance or documentation invariant that
the kit has written down becomes one the kit cannot land past.

## Why this one

Every other finding in this audit is either downstream of it or smaller than it.
Six advisory surfaces (`coverage-map-validation`, `decision-map-integrity`,
`roadmap-detail-integrity`, `docs-currency-workflow`, `skills-index-command-adapters`,
`guidance-prose-budgets`) become enforced by the same change. It is the difference
between Bench-as-prose-kit and Bench-as-enforced-platform, and it is the one thing
Bench has that Pocock's skills do not — currently switched off.

## Acceptance

- [ ] `bench gate` at `58d966e2` exits **non-zero**, naming all 10 current
      diagnostics with their existing messages.
- [ ] After the 10 fixes land, `bench gate` is green and the run reports
      `phase conformance: green`.
- [ ] Deleting `## Entry orientation` from any `.agents/commands/*.md` reds the
      gate with that check's own message; restoring it greens the gate.
- [ ] Unsetting `BENCH_CONFORMANCE_ROOT` in the phase reds the gate with
      `environment skip inside the oracle`, not a footer count.
- [ ] The skip line names the skipping test and its reason:
      `class=environment: 1 (TestRootConformance: BENCH_CONFORMANCE_ROOT not set)`.
- [ ] Added gate wall-clock is reported in the landing evidence
      (measured: +3.2 s against a 64 s baseline).
- [ ] `roadmap/FT133`'s two "matched but skipped / printed ok" occurrences are
      closed by this change, or the row states why they are not.

## Proof plan

1. **Red first:** wire the phase, run `bench gate` at HEAD, capture the 10
   diagnostics — the red is the evidence the phase is connected to the live root,
   not to fixtures.
2. **Mutation probe (omission kind):** delete the phase's env assignment; the gate
   must go red on the environment-skip rule, not silently pass. Restore.
3. **Mutation probe (different site and kind, per `craft-delegate`):** plant a
   heading deletion in a command file; the gate must red with
   `docs-currency-workflow`'s own message. Restore.
4. **Green:** fix the 10; gate green; `bench prep-release`'s ship-tier conformance
   must remain green and must not double-run the dev set.
5. **Cost:** report the wall-clock delta, since the gate taxes every shift iteration
   (`craft-gate`: *"keep checks cheap"*).
---

## Required adversarial questions

**1. What does Bench ask modern Claude or Codex to do that they already do
reliably without help?** Explore a repo before editing; use domain vocabulary;
write vertical slices; keep tests off private internals; separate a review's
concerns. `craft-tdd`'s three anti-patterns, `craft-seams`' "prefer an existing
seam", and `craft-review`'s two-axis rationale are largely restating current model
defaults. **Candidates for the leave-one-out ablation in W.**

**2. What does Bench make more consistent?** Mechanically: that "done" means an
external verdict on an exact tree; that a delegate runs on a bound model; that
destructive git is refused; that a shift never commits red; that a spec's stories
are all referenced or explicitly excepted; that a decision map cannot be declared
ready with an open ticket; that a phase cannot enter on a stale base or an
unauthorized path. Those are the ones I verified by execution.

**3. Which consistency gains are merely claimed, and which are observed?**
*Observed:* every item in Q2. *Claimed only:* effort adherence, iteration caps,
fan-out counts, the inline-edit allowance, the coordinator's six-item delegate
verification, review-axis independence, and the outcome effect of all 28 prose
artifacts.

**4. Which Pocock-derived behaviors are load-bearing?** Repro-loop-before-hypothesis
with exact-symptom assertion and an original-loop rerun; seams chosen and confirmed
before implementation; tracer-bullet vertical slicing with blocking edges;
axis-separated review with no reranking; the Fowler smell baseline; the frontier
grill with a recommendation per question; no file paths in durable artifacts.

**5. Where has Bench weakened an upstream skill through composition?** In
`/bench-debug`: the loop-construction menu cut from 10 options to 5 (dropping
headless-browser, property/fuzz, bisection, differential, HITL); the two hard
stop-gates and the checkbox completion form removed; "tighten the loop" removed as
a named step. And structurally: making `$bench-debug` explicit-only on Codex while
upstream `diagnosing-bugs` is model-invocable.

**6. Where has Bench improved an upstream skill through state or enforcement?**
The acceptance coverage map with `bench coverage --check` (4/4 hostile inputs
caught); the finding-citation standard plus "refute before you report" and
"a universal claim without an enumeration is a sample"; a third review axis; the
delegate-probe independence rule (different kind *and* site); ownership fences
checked by `bench preflight`; the quarantine-marker mechanism so a repro survives
shift rollback; `bench spec history` for recovering a deleted spec.

**7. Where does Bench ask a model to remember something software should enforce?**
Effort; iteration caps; fan-out; the one-line inline-edit allowance; running
`bench coverage --check`; the coordinator's six delegate verifications; the
probe-differs-in-site rule; that a phase file keeps its required headings (a check
exists — it does not run).

**8. Where does Bench ask software to decide something that requires judgment?**
Almost nowhere, and it says so out loud in the right places: `bench models` is
"advisory discovery, never validation"; `bench roadmap` "does no judgment" and
extracts a human-written sequence verbatim; `bench outline` "does not identify the
project's blessed seams"; `bench canary` is "inventory-only". The one place I would
push back is `bench structure`'s fixed 400-line file budget driving a permanent
`status` row — file length is a proxy for a judgment (`craft-seams` says "split
along responsibility, never line count") and the row nags on the proxy.

**9. Where does Bench confuse available context with active context?** Rarely.
The always-loaded core is 4.8 k tokens and everything else is pointer-fired. The
one real instance is `bench outline`'s bare form: 7 k tokens of near-arbitrary
symbols, mostly test fixtures, presented as if available equals useful.

**10. Where does Bench confuse a plan with current state?** `capture/session-handoff.md`'s
`## State`, at HEAD: it names a superseded HEAD, says `specs/` is empty (34
tickets-only folders), and describes a batch "about to land as one commit" that
landed. The pin block above it is regenerated correctly by `bench handoff`; the
prose below it is not, and the staleness detector measures write-recency rather
than content agreement, so it stays silent.

**11. Where does Bench confuse assertion with observation?** Structurally, almost
nowhere — the gate is an observation and there is no assertion channel to abuse.
The exception is the reporting layer: `bench status` and `bench handoff` *assert* a
red the oracle denies (I.3).

**12. Where does Bench confuse a passing test with requirement satisfaction?**
It does not, and it is explicit: the gate "cannot tell whether you built the right
thing the right way", which is why review exists and why the coverage map ties
stories to seams. The real gap is that the map's validity check does not run in
the gate.

**13. Where does Bench confuse repeated review with independent review?** It
mostly does not — three axes in parallel fresh contexts, each re-deriving from its
primary source, is real independence. The residual is that the diff arrives with
the commit log, so the axes see the implementer's rationale. Arm G of W settles it.

**14. Where does Bench confuse subagent count with useful parallelism?** It does
not. Every delegation boundary I found buys isolated ownership, independent
rediscovery, or bounded research. Fan-out is declared and unchecked, which is a
gap in enforcement, not in design.

**15. Where does Bench confuse detailed prose with control?** Everywhere the
prose is the only mechanism: effort, caps, fan-out, the inline allowance, delegate
verification, the review dispositions — and, most consequentially, the six
conformance checks that *are* mechanisms and are switched off, so that today they
function exactly as prose.

**16. Where does Bench allow blank retries?** `craft-line`'s ladder is
evidence-carrying by construction (classify reds, feed gate output back, name the
changed variable), and a blocked delegate's return shape is specified. Nothing
*checks* it, and nothing detects two consecutive identical reds — I.4.

**17. Where does Bench let an agent patch before reproducing?** Only outside
`/bench-debug`. Since the phase is reviewer-invoked (and explicit-only on Codex),
the default path for "the gate went red" is `/bench-implement-spec`'s repair route
or `craft-line`'s ladder, neither of which requires a red-capable repro first.
That is precisely the gap I.4 closes.

**18. Where can an agent keep thinking after a discriminating test is available?**
`bench maps`, `bench structure`, `bench coverage --check`, `bench preflight`, and
root conformance are all one-second discriminators that nothing requires anyone to
run. And the gate itself: `[test] ok internal/conformance` invites the agent to
stop, exactly as `craft-tdd` warns.

**19. Where can stale evidence authorize work?** Greens: nowhere I could break —
the tree-hash key held against five attacks. Reds: they authorize *inaction* —
`bench status` and `bench handoff` report a phantom red, which sends an agent to
fix nothing.

**20. Where can local verification masquerade as global completion?** Not at the
landing: `bench worktree land` gates the composed tree, and a delegate's focused
green is explicitly not accepted as done. The masquerade is one level up — the
whole-tree gate is green while the whole *kit* is red on 10 conformance checks.

**21. Which instructions disappear if the CLI improves?** "Run `bench coverage
--check`"; "read `bench status` and run the housekeeping rows"; "re-run the gate";
the handoff's "`dist/bench` must exist"; `/bench-setup-repo` §3's "fix the wiring";
and the whole class of "remember to check X" once X is a gate phase.

**22. Which CLI commands exist because the workflow is unnecessarily complicated?**
`bench outline` (a seam-finding aid whose own help disclaims finding seams);
`bench structure` (a proxy metric that generates a permanent unenforced row);
`bench maps --template` (a schema the phase file must not restate — correct
single-sourcing, but it exists because decision maps are a hand-authored format);
and arguably `bench dashboard` (an HTML snapshot of a board that `bench roadmap`
already projects).

**23. What would Agentless remove?** Shaping, specs, tickets, coverage maps,
fences, delegation, worktrees, review axes, retros, the roadmap — leaving
localize → repair → validate. Bench's own lighter-path table already permits
nearly that for a one-ticket change, which is the right concession. What Agentless
cannot remove without loss: the external oracle, the destructive-action guards, and
the resumable handoff.

**24. What would SWE-agent change about the interface?** Very little — Bench
already treats the CLI as an agent interface. It would flag `bench outline`'s bare
form and `bench status`' prose action column, and it would ask for measured turn
counts per question, which Bench does not collect.

**25. What would ReAct change about the cadence?** It would make the observation
mandatory rather than available: a gate red must be followed by an *observation*
(the failing check re-run, the repro built), not by another edit. That is I.4.

**26. What would Lost in the Middle change about context compilation?** It would
approve the current shape (small always-loaded core, pointer-fired references,
targeted `--row` fetches) and object to `bench outline`'s 200-row dump — 7 k tokens
whose middle nobody will use.

**27. What does J-space research suggest about active state without justifying
J-space terminology?** That presence in context is not activity in computation, so
the constraint that matters now should be small, explicit, and separately carried.
Bench already does this and does it better than the sketch it would be importing,
because its `Verified` is a content hash rather than a self-report. Adopt none of
the vocabulary.

**28. Which J-Space Suite ideas would be cargo cult if copied?** The J-space naming
itself; a `Goal/Core/Verified/Open/Next` YAML as a *new file* (a sixth source of
truth beside the five Bench has); the `fast/full/loop` tier axis (it would collide
with `craft-line`'s cheap/mid/top); and the benchmark table, which its own authors
mark as single-run, no CIs, not a unified harness.

**29. What state must survive a model swap?** The commit graph; the spec and its
coverage rows; the ownership fences; `.bench/lines.env`; the decision maps; the
gate verdict for the exact tree; the handoff pin block; the worktree lease and
assignment ledger; open review pickups; and — the one thing that currently does
not — the repro command for an in-flight bug.

**30. What state should deliberately not survive?** The conversation transcript;
the implementer's rationale (so review can re-derive); a delegate's internal
reasoning; a gate verdict for any other tree; and `capture/retros/*.md` after the
drain — Bench deletes them on purpose, which is right.

**31. Which facts should be broadcast?** The invariants; the resolved line; the
fences; the coverage rows for the ticket at hand; the expected tip and base; the
exemplar to mirror; the fixture/seam inventory. All are facts, all are cheap.

**32. Which conclusions require independent rediscovery?** Implementation
correctness; requirement satisfaction; absence of bugs; the review verdict; and any
delegate's done-claim. `craft-review`'s "re-derive, then compare" and
`craft-delegate`'s "a claim, not a result" both say this; neither is checked.

**33. Is `what-next` already the right router?** No. Its own frontmatter says
"Maintenance, not a workflow phase", it is hidden from the model, and it drains
capture into `ROADMAP.md`. It is well-designed for what it does and misnamed for
what people will expect.

**34. Is `/bench` a useful front door or only another alias?** Useful **only if**
it routes from observed state. As a menu or an alias it is another name to
remember. As `bench status --route` emitting one `next[1]{state,why,command}` row,
it is the missing projection over facts Bench already computes — and it beats
Pocock's `ask-matt`, which can only read prose.

**35. What should `bench init` own, if anything?** Nothing new. `bench setup`
already owns adoption and does it well (plan-first, converging, self-verifying,
with a named next action). `bench init` scaffolds only `.bench/gate.sh` and is a
subset. Keep it as plumbing; the router's row in an un-adopted repo should say
`bench setup`.

**36. Can a new user find the right path without reading the skill graph?** No.
Today: a README menu of eleven, a `bench status` that says "nothing pending" in an
un-adopted repo, and no signal for a staged spec. With P1-1: yes.

**37. Can a fresh Codex session resume Claude's work?** Structurally yes — the
handoff pin block is machine-regenerated, `bench handoff --harness codex` rewrites
the next command, `$bench-*` adapters exist, and the tree wins where they disagree.
Degraded by the drifting `## State` prose and the phantom-red gate line.

**38. Can a fresh Claude session resume Codex's work?** Same mechanism, same
answer, plus one asymmetry: Claude's model can self-invoke six phases, so it
recovers from a missing handoff more readily than Codex can.

**39. Can the `diagnosing-bugs` behavior be preserved with less prompt load?**
Probably — the four load-bearing constraints (H.1 items 1, 2, 3, 10) are ~150 words.
But "probably" is the wrong currency here: measure with W arms D/E/F first. Note
the direction of the current evidence is that Bench already compressed it and lost
three things (H.4).

**40. Does reducing it lose the mechanism that breaks repair loops?** The mechanism
is the *forced contact with reality before theorizing*, and that survives
compression. What did not survive is the **menu** that tells an agent how to build
the loop for a bug class it has no obvious loop for — and for those bugs, losing
the menu is losing the mechanism, because there is no loop to force.

**41. If Bench were rebuilt today, what 20 % preserves 80 % of demonstrated
value?** The tree-hash gate + `bench commit`; the destructive-git and agent-line
guards; the conformance suite (**wired**); `bench setup`; pooled worktrees with
ownership fences and `bench preflight`; the acceptance coverage map with
`--check`; `/bench-debug`; `craft-review`'s citation standard; the handoff pin
block; `.bench/lines.env`. Roughly the CLI plus five prose artifacts.

**42. What would a rewrite accidentally destroy?** Listed concretely in R: the
untracked-inclusive/gitignored-exclusive tree key; `git restore --source=` in the
deny surface with `--staged` allowed; denial of an *omitted* model field;
`spec retire`'s refusal to auto-clean a residue; `worktree clean`'s
plan-then-apply fingerprint; the probe-differs-in-kind-and-site rule; the 29
checks. All scars from real failures, none legible from outside.

**43. What is the smallest architectural spine that can absorb the proven
workflows?** The one in U — the current architecture plus a router, minus one
verdict reader, plus a conformance phase. The brief's candidate spine describes
Bench as it already exists; building it separately would be a re-derivation, which
is the defect Bench's own code standard forbids.

**44. What should be deleted before anything new is added?** The 34 orphan ticket
directories; the duplicate `disable-model-invocation` keys on the Codex adapters;
`bench structure` from the permanent status row set; `bench outline`'s bare-form
row dump; and three of the four in-tree assessment artifacts.

**45. What is the one next ticket with the highest leverage?** Section Z — wire
root conformance into `bench gate` and make an environment-class skip inside the
oracle red. It converts six written invariants into enforced ones in a single
change, and it is the only finding in this audit that would have caught the other
findings.

---

## Closing note

The question this audit was given was whether Bench improves the realization of
model capability or has accumulated ceremony. The answer is that Bench built an
unusually good enforcement layer and an unusually well-written prose layer, and
then — through a single missing environment variable — stopped pointing the first
at the second. It is not a system that mistook prose for control. It is a system
that built real control, wired it to the wrong entry point, and had no surface
that would say so.

Every defect in this report is small. That is the good news and it is also the
warning: a 64-second oracle, a mis-ordered `if`, an absent status row, and an
unset variable are individually trivial and collectively the difference between a
platform that grades itself and one that believes it does.
