# Bench capability-realization audit

Subject: `58d966e2f92f7f37eba07b6215e8eef45371b72d` (`audit/sol`)  
Primary evidence: [evidence.md](evidence.md)  
Required question responses: [questions.md](questions.md)

## A. Executive verdict

**Verdict: directionally correct, operationally uneven, and not yet demonstrated as a coherent capability-realization harness. Consolidate it through a strangler migration; do not rewrite it and do not continue adding phases.**

Bench's strongest architectural idea is that completion is an exact, replayable state transition rather than an agent's assertion. Its gate binds the work-tree identity and oracle identity, its prospective commit path gates the exact content to be committed, and its ownership/preflight checks reject moved subjects, unowned paths, uncited requirements, and phantom coverage. Those behaviors were exercised, not inferred. `[OBSERVED: E007–E009]`

The demonstrated behavior that must survive is the deterministic substrate: exact-tree gate evidence and invalidation, atomic landing, owned worktrees and recovery handles, preflight/coverage checks, compact read queries, hostile-input fixture ownership, and the reproduction-first debugging loop. The last item has repeated practitioner evidence but lacks a controlled Bench outcome study. `[OBSERVED: E007–E008, E013]` `[REPEATED PRACTITIONER EVIDENCE]`

The largest weakness is the gap between this substrate and the model-facing control plane. Bench has no normal-user router, no canonical active work record, no context compiler, no executable Review or Final Check, no general claim/evidence state, and no recorded proof that most of its 31,456 words of guidance improve task outcomes. Markdown describes transitions that software neither owns nor observes. The result can be simultaneously green, structurally noisy, semantically stale, and hard to enter. `[OBSERVED: E009–E019]`

Remove or retire fixed-count delegation, prompt-only `--full` orchestration, duplicate command/skill adapters, ticket-only spec sediment, invalid decision maps, stale handoff narrative, generic guidance that does not apply to the active profile, and public maintenance plumbing with no external consumer. Do this only with compatibility and benchmark evidence; do not erase the exact gate or debugging mechanism to make the system look smaller. `[INFERRED]`

Strengthen one spine: a logical **`bench`** router backed by deterministic status plus a small canonical work-state store and context compiler. `/bench` in Claude and `$bench` in Codex should be thin harness adapters to that same entry, not distinct workflows. The compiler should emit only goal, governing decisions, scope, fresh evidence, open uncertainty, next discriminator, failed-attempt diagnostics, and coverage needed for the selected action. `[INFERRED]`

Build exactly one tracer ticket next: the canonical work-state/router/context-compiler slice in section Z. It must read current status and handoff data, preserve backward compatibility, route representative states, expose why it chose the route, and beat current Bench on a preregistered A/B/C benchmark before more phases move behind it.

## B. Audit environment

| Field | Recorded value |
|---|---|
| Audit mode | In-repository dogfood plus external cold-observer reconstruction; not a paired live-harness audit |
| Harness | Codex |
| Model | Sol model line; the runtime did not expose a verifiable model build identifier |
| Effort / reasoning | xhigh, as assigned to this audit run |
| Repository | `/home/mgibs/workspace/bench-audit-sol` |
| Commit / branch | `58d966e2f92f7f37eba07b6215e8eef45371b72d` / `audit/sol` |
| Initial worktree | No staged or tracked changes; untracked audit charge only (`?? audit/`) |
| Final worktree | No staged or tracked changes; only `?? audit/`, containing the supplied charge and requested audit output. All temporary experiment/build/cache artifacts were removed. `[E023]` |
| Auto-loaded instructions | System/developer Codex policy and repository `AGENTS.md`; `.bench/BENCH.md` and `projects/benchkit.md` were then read as required project sources. `CONTEXT.md` and component guidance were inspected after implementation reconstruction. |
| Auto-available project skills | 17 craft/prototype skills were exposed to Codex. Eleven on-disk `bench-*` phase adapters were not exposed in this session. All 28 project skill files and 11 command documents were nevertheless inspected. `[OBSERVED: E015]` |
| Skills applied | `craft-synthesis`, `craft-cli`, `craft-skills`, `craft-seams`, `craft-spec`, `craft-tickets`, and `craft-adr`; they constrained evaluation of kit changes, CLI claims, boundaries, specs, tickets, and resulting-state documentation. |
| Tools | Read-only shell/search, exact-source local builds, local Go/shell tests, isolated local clones/fixtures, patch-based report writing, and read-only web/upstream retrieval |
| Output | `audit/sol-xhigh/` in this disposable audit worktree, as explicitly requested |
| Safety boundary | No commits, pushes, releases, deployments, secret reads, check weakening, or Bench implementation edits. Temporary artifacts were isolated and cleaned after observation. |

**Limitations.** `[OBSERVED]` This was one Codex trajectory, not a multi-model randomized study. Claude rendering was exercised and Claude adapter files were inspected, but no live Claude model performed the tasks. Review-with/without-narrative and four-way debugging comparisons require repeated model runs and were designed, not fabricated. Ship-tier services and publication were not invoked because they cross the safe boundary. Research establishes plausible mechanisms and counterevidence, not Bench-specific causality.

## C. Bench reconstructed

The implementation-first reconstruction is:

```text
installation/link
  -> wrapper + pinned platform binary + harness adapters/hooks
  -> static root inventory OR exact binary subcommand dispatch

user/model chooses a phase (mostly Markdown)
  -> reads AGENTS + BENCH + profile + selected craft/phase guidance
  -> may shape -> spec -> tickets -> implement/debug -> review -> final-check
  -> phase state is distributed through Git, specs, tickets, capture, and prose

deterministic substrate
  -> status/read queries
  -> worktree ownership / intent ledger
  -> preflight + coverage
  -> test / gate(tree identity + oracle identity)
  -> commit / release state machines

resume
  -> handoff pins + arbitrary reviewer-authored State body
  -> first syntactically invocable suggested command
  -> next harness repeats context loading and phase selection
```

The actual states and transitions are not one state machine. They are five overlapping systems:

1. **Git is durable truth:** content, refs, staged prospective tree, commits, worktree metadata.
2. **The gate is a deterministic verifier:** current tree plus oracle identity becomes green/red/timeout evidence; drift prevents reuse. `[REPRODUCED: E007–E008]`
3. **Ownership and acceptance checks are deterministic but phase-entered:** intent/worktree records, fences, acceptance rows, and preflight verdicts reject known classes once called. `[REPRODUCED: E007]`
4. **Workflow phases are prompt protocols:** Shape, Build Spec, Implement, Diagnose, Review, and Final Check rely on a model to perform and remember transitions. `review`, `final-check`, and `claims` are not executable subcommands. `[DIVERGED: E010]`
5. **Resume state is a partially structured document:** repository/branch/commit pins are renderable; `State` is arbitrary prose with commit-age freshness only. `[DIVERGED: E011–E012]`

The happy path advertised by the prose is reasonable—idea → decision → spec → vertical tickets → fresh implementation → independent semantic review → gate/landing → handoff—but there is no single owner of the path. A user can enter directly at a CLI command or skill, skip semantic phases, write unsupported state, and still obtain a green development gate. This is not a gate bug: it is an absent control-plane contract.

Actual transition inventory:

| State | Entry; required artifact/evidence | Legal next / regression or repair | Owner; persistence; resume / cross-model | Mechanical representation and bypass |
|---|---|---|---|---|
| Idea | `bench idea` or capture edit; idea row, no evidence | Shape, defer, delete; clarify ambiguity | User/maintainer; `capture/IDEAS.md`; visible as file | Parking is mechanical; selection/shaping is prompt and skippable |
| Shaped decision | Shape phase; chart/decision ticket, reviewer choices | Spec, prototype, deepen, park; reopen decision | Shaper + reviewer; decision files/maps | Prompt state inferred from artifacts; can jump directly to code |
| Specified | Write Spec; `spec.md` with status/coverage rows | Ticket, implement bounded slice, revise/retire | Spec author/reviewer; spec file; cross-model readable | Schema/coverage partially mechanical; semantic completeness is review judgment |
| Ticketed | Ticket slicing; ticket files with scope/deps/rows | Assign/implement, block, reslice | Orchestrator; Git files | Preflight can check references/fences; about 40 ticket-only files currently have no live spec owner |
| Assigned | Worktree create/intent; assignment/base/token/fence | Implement, release, reauthorize, clean | Orchestrator/worktree owner; Git metadata + ledger | Strong mechanical identity; direct unmanaged work can bypass phase ownership |
| Implementing | Edit in owned subject; diff/attempts/tests | Locally verify, debug, checkpoint, abandon | Writer/shift agent; Git/worktree, limited shift evidence | Content durable; general objective/attempt state is not |
| Locally verified | Focused tests/preflight; command results and coverage | Review or gate; regress to implement/debug on failure | Builder; mostly logs/conversation | Tests executable; no general durable claim links result to requirement |
| Reviewed | Review prompt over fixed diff/spec; finding/pickup prose | Repair, gate/land when no findings, escalate | Independent-axis roles in prose; review files | Advisory only; can be skipped, independence unrecorded, persistence route contradictory |
| Development-complete / integrated | Fresh exact gate then prospective commit/merge | Checkpoint/final/CI; regress on tree/oracle drift or CI failure | Gate oracle + integrator; Git and gate record | Strong mechanical transition for declared oracle; does not require semantic review |
| Promoted / released | Prep/preflight then release state transition | Submit, promote, rollback, status | Human operator + release machine; release records/registry | Strong CLI exists, but current workflow's raw publish is a bypass |
| Checkpointed/resumable | `bench handoff`; pins plus State and next | Resume any routed phase; rewrite checkpoint | Closing agent; committed Markdown | Pins derived; State arbitrary, commit-age only, cross-model prefixes partial |

The requirement trace is mechanically strong only at islands: a named acceptance row can survive spec → ticket → preflight/coverage, and a planted uncited/phantom omission turns those checks red. It can still disappear before preflight, live only in a ticket-only directory that coverage cannot open, or be semantically “satisfied” by prose because Review is advisory. `[OBSERVED: E007, E010, E016]`

Classification summary:

| Area | Result | Reason |
|---|---|---|
| Gate, commit, ownership, preflight | **REPRODUCED** | Implementation, prose, and adversarial tests agree. |
| Compact read queries | **REPRODUCED** | Content-first TOON-style output and structured failure generally match AXI guidance. |
| Root entry | **DIVERGED** | Wrapper and exact binary disagree; neither routes intent. `[E004, E019]` |
| Review / Final Check | **UNREPRODUCIBLE** | Prompt procedures exist; no executable phase/control evidence exists. |
| General claims | **UNSOURCED** | No claim schema/command; prose artifacts accept assertions. |
| Handoff freshness | **DIVERGED** | Commit-distance freshness passes while State contradicts the live tree. |
| Codex/Claude parity | **DIVERGED** | Common prefix rendering exists, but adapter availability and hooks differ. |
| Spec-build lifecycle | **PARTIAL** | Coverage/preflight are real; current ticket-only artifacts are invisible to live spec status. |

## D. Entry-point verdict

The one canonical normal-user entry point should be **`bench`**.

`bench` with no subcommand should accept optional plain intent, combine it with deterministic repository/work-state facts, and return one route, the reason, missing prerequisites, and one copy-paste next command. `/bench` (Claude Code) and `$bench` (Codex) should only adapt harness syntax to the same logical entry. They are not additional canonical entries.

Representative routing contract:

| State / intent | Canonical route |
|---|---|
| New install | `bench init` transaction, then `bench` |
| No active work | show definitive empty state and accept an idea or expert command |
| Unshaped ambiguous idea | Shape/grill practice |
| Shaped wide change without spec | Spec practice |
| Spec without bounded work | Ticket slicing |
| Ready bounded ticket | Implement |
| Exact failure or reported bug | Debug |
| Interrupted work | Resume from canonical checkpoint and next discriminator |
| Dirty/unverified change | Focused verification, then gate |
| Locally complete | Independent review, then landing readiness |
| Expert names a workflow | Directly enter it after prerequisites are checked |
| Harness switch | Render the same state with the target harness prefix |

`bench init` should own only a transactional minimal install: platform files, selected gate scaffold, adapters/hooks, one project-profile pointer, and initialization of the work-state store. It should not interview the architecture or author policy. Expert entries remain direct: `bench gate`, `bench status`, `bench preflight`, `bench coverage`, `bench worktree`, `bench test`, `bench commit`, `bench release`, and named practices such as debug or review.

Today the wrapper's no-argument behavior prints more than forty commands and exits 0; the exact binary returns `no subcommand` and exits 2. `help` also differs. `commands --brief` is intentionally a three-command liveness probe. None distinguishes user states, and `what-next` is roadmap/capture maintenance rather than a front door. `[OBSERVED: E001, E004, E018–E019]`

Claude currently has phase-command symlinks and a Claude-only agent-line hook; this Codex session exposed craft skills but not the on-disk phase adapters. Prefix translation for derived handoff actions is single-sourced, but manual `--next` strings remain deliberately opaque. Therefore the same repository is not the same experience across harnesses. `[OBSERVED: E012, E015]`

## E. Provenance map

| Source | Exact revision | Bench realization | Assessment |
|---|---|---|---|
| Matt Pocock skills, adopted v1.1 | `mattpocock/skills@d574778f94cf620fcc8ce741584093bc650a61d3` (`v1.1.0`), adopted by Bench `5e3c0ba98500a904be47673be975f8770cb33d0d` | Grilling, domain language, ADRs, specs, tickets, seams, TDD, debugging, review, comments, delegation, design-system use, skill craft, prototyping | Deep, recognizable integration; Bench adds state/fences/gate. |
| Matt Pocock later comparison | `84fdeffd12f2ee307994d1eb6feb48173b6e0502` | Roadmap comparison source | Pin exists, but later deltas are not a systematic merge ledger. |
| Matt Pocock current audit snapshot | [`9c9f36ccd3995266cd675468af71639c8dde1ec5`](https://github.com/mattpocock/skills/tree/9c9f36ccd3995266cd675468af71639c8dde1ec5) | Current upstream comparator for 12 corresponding practices and Ask Matt routing | Current upstream has a clearer router and more explicit context boundaries. |
| Kunchenguid no-mistakes current snapshot | [`6859d1e827f5ab2592a4703d3bab8734a38c9aa5`](https://github.com/kunchenguid/no-mistakes/tree/6859d1e827f5ab2592a4703d3bab8734a38c9aa5) | Disposable-worktree review/testing pipeline, safe-mechanical versus judgment boundary, isolated review turns | Conceptual influence is named, but exact historical import revision is absent. `[E022]` |
| Kunchenguid AXI current snapshot | [`408a6536625e5b05e5c56e6c4a04fe83e1f510a5`](https://github.com/kunchenguid/axi/tree/408a6536625e5b05e5c56e6c4a04fe83e1f510a5) | TOON tables, minimal query schemas, structured errors, ambient context, contextual disclosure | Strong influence on read commands; root inventory violates the same economy. Historical import pin absent. |
| Bench-native | Subject commit `58d966e2f92f7f37eba07b6215e8eef45371b72d` | Exact-tree/oracle gate, prospective commit, worktree ownership, shift, preflight/coverage, intent ledger, status/handoff, decision maps, release state | The most demonstrated value; retain behind a smaller interface. |
| External research | Pinned papers/repos in section X | Evaluation lenses for selective workspace, action cadence, interface design, recovery, efficiency, and skill effectiveness | Research constrains claims; it does not prove Bench outcomes. |

The repository can reconstruct the Pocock adoption point exactly. It cannot reconstruct which no-mistakes or AXI source revision produced each integrated behavior. That provenance gap should be repaired prospectively with source revision plus a behavioral delta, not by guessing history. `[OBSERVED: E021–E022]`

## F. Component inventory

The compact cells below cover the required fields: **purpose; input → output → consumer; enforcement/provenance; observed value; cost/problems; recommendation.** “Prompt” means the model is the enforcement boundary; “check” means executable enforcement exists.

### Skills

| Component | Purpose; input → output → consumer | Enforcement / provenance | Observed value; cost / problems | Recommendation |
|---|---|---|---|---|
| `bench-craft-adr` | Record durable decision; resolved choice → resulting-state ADR → future agents | Prompt; Pocock-derived | Clear memory boundary; overlaps maps/docs if used indiscriminately | **Keep**, load on demand |
| `bench-craft-cli` | Design agent-facing CLI; interface task → AXI-shaped contract → CLI authors | Prompt + downstream tests; AXI/Bench | Bench queries demonstrate value; root violates it | **Keep** as reference and dogfood it |
| `bench-craft-comments` | Govern valuable comments; code context → sparse rationale → maintainers | Prompt; Pocock-derived | Sensible prose discipline; modern models often do baseline behavior | **Simplify** |
| `bench-craft-delegate` | Scope/verify delegates; task graph → charges/results → orchestrator | Prompt; Pocock/Bench | Ownership language useful; mandatory delegation and fixed ceremony unmeasured | **Benchmark first**, make delegation optional |
| `bench-craft-design-system` | Consume UI system; UI task → compliant component choices → UI implementer | Prompt; Pocock-derived | Useful only in UI profiles; cold cost elsewhere | **Move** to profile-selected guidance |
| `bench-craft-domain` | Stabilize vocabulary/edges; fuzzy concepts → glossary/scenarios → spec/review | Prompt; Pocock/Bench | Distinguishes names from decisions well | **Keep**, progressive disclosure |
| `bench-craft-gate` | Author the oracle; gate change → biting check → repository | Executable proof required; Bench-native | Load-bearing one-source and fail-posture discipline | **Keep and strengthen** |
| `bench-craft-grill` | Resolve decision frontier; idea/facts → closed choices → shaper/spec | Prompt; Pocock-derived | Strong anti-assumption mechanism; can over-interview small work | **Keep**, route only ambiguity |
| `bench-craft-line` | Select model/effort; work risk → tier choice → orchestrator | Mostly prompt; Bench-native | Makes cost explicit; Codex hook parity absent and outcome delta unmeasured | **Benchmark**, then simplify or enforce |
| `bench-craft-review` | Review Standards/Spec/Coverage; fixed diff → cited findings → implementer | Prompt; Pocock/Bench | Three lenses are useful; independence and completion are assertions | **Keep slim judgment core**, execute isolation/state |
| `bench-craft-seams` | Place interfaces/tests; behavior boundary → deep seam → design/spec/test | Prompt; Pocock-derived | High-value judgment not reducible to file size | **Keep** |
| `bench-craft-skills` | Author reliable skills; skill change → lean trigger/content → kit maintainer | Prompt + conformance anchors; Bench-native | Useful for kit work; self-referential active-load cost | **Move** to maintainer profile |
| `bench-craft-spec` | Specify wide change; decisions/edges → stories + coverage rows → builders/review | Prompt plus coverage schema; Pocock/Bench | Coverage rows are an enforceable improvement; overkill for bounded work | **Keep for wide work** |
| `bench-craft-synthesis` | Fold kit candidates safely; proposed change → assessed proposal → reviewer | Prompt; Bench-native | Protects closed decisions and dogfood loop | **Keep** in kit-maintainer profile |
| `bench-craft-tdd` | Apply red/green/refactor at chosen seam; behavior → regression + code → builder | Prompt plus tests/gate; Pocock-derived | Useful when seam chosen; can overfit or become ceremony | **Keep slim**, trigger selectively |
| `bench-craft-tickets` | Slice vertical tracer bullets; spec → bounded tickets → fresh writers | Prompt + ticket schema/preflight; Pocock/Bench | Vertical slices valuable; mandatory fresh writer per ticket is unmeasured | **Keep for wide builds**, relax delegate mandate |
| `prototype` | Answer one named decision with disposable code; hypothesis → running evidence → reviewer | Prompt, cleanup obligation; Pocock-derived | Direct empirical escape | **Keep** |
| `bench-assess` | Assess repository across six areas; repo → findings/roadmap → maintainer | Phase adapter; Bench-native | Broad coverage; fixed six-agent fan-out and output load unmeasured | **Simplify / benchmark** |
| `bench-debug` | Break repair loops; symptom → reproduced root cause + regression/fix → maintainer | Prompt plus tests/gate; Pocock `diagnosing-bugs` + Bench | Mechanism preserved and valuable; isolation/shift prose dilutes it | **Preserve and strengthen**; see H |
| `bench-deepen` | Explore alternate designs; chosen question → alternatives/decision → reviewer | Phase adapter; Bench-native/Pocock | Useful for consequential ambiguity; delegate count is ceremony | **Benchmark first** |
| `bench-final-check` | Report landing/readiness and capture learnings; landed work → closure summary → user | Prompt only; Bench-native | Some housekeeping has executable owners; no differential control | **Merge** into executable landing/status |
| `bench-implement-spec` | Build tickets; staged spec → gated slices → reviewer | Prompt orchestration + CLI substrate; Bench/Pocock | Uses strong preflight/gate; `--full` is prose, not orchestration | **Simplify**; delete claimed mode until executable |
| `bench-review-implementation` | Run three semantic axes; diff/spec → findings/pickups → implementer | Prompt only; Pocock/Bench | Useful rubric; persistence, commit, and independence contradictions | **Move control to runner**, keep rubric |
| `bench-setup-repo` | Install/profile Bench; repository → linked kit/profile → users | Prompt + `link/init`; Bench-native | Setup capabilities exist; mixes transaction with architecture interview | **Split**: deterministic init, optional shaping |
| `bench-shape-idea` | Resolve an idea into a decision ticket; idea → chart/decision → spec/prototype | Prompt; Pocock Ask Matt / Bench | Useful for genuinely ambiguous work; can stop early or over-shape | **Keep only for ambiguous/wide work** |
| `bench-update-kit` | Evaluate upstream/candidate kit changes; source diff → proposal → maintainer | Prompt; Bench-native synthesis | Appropriate for this kit; irrelevant to linked projects | **Move** to maintainer profile |
| `bench-what-next` | Drain roadmap/capture and write handoff; board → maintenance action/state → maintainer | Prompt + status/handoff helpers; Bench-native | Useful maintenance ritual; falsely resembles router and produced stale state | **Rename/split**; never canonical entry |
| `bench-write-spec` | Turn closed decisions into build contract; decision → spec/tickets → builders | Prompt + coverage/preflight; Pocock/Bench | Strong for wide changes; duplicates shape/ticket narration | **Consolidate** with wide-work compiler |

The eleven `.agents/commands/bench-*.md` files are harness entry documents for the eleven `bench-*` phase skills above. Their purpose/input/output/consumer is identical to the matching phase skill; enforcement is prompt discovery, provenance is Bench-native adapter code, observed value is Claude discoverability, cost is duplicated inventory/anchors, and the current Codex session did not surface them. **Recommendation: generate thin adapters from one registry and test both harness inventories.** `[OBSERVED: E013, E015]`

### CLI subsystems and gates/checks

| Component | Purpose; input → output → consumer | Enforcement / provenance | Observed value; cost / problems | Recommendation |
|---|---|---|---|---|
| Wrapper/install (`link`, `unlink`, `upgrade`, `doctor`, `repair`, `version`) | Install/resolve pinned binary; repo/platform → runnable Bench → users/hooks | Shell + manifests; Bench-native | Safe lifecycle intent; cold source checkout lacked pinned binary and wrapper semantics differ from binary | **Keep**, unify dispatch/help and add source fallback diagnostics |
| Root registry/help | Discover commands; no args/help → inventory → user/agent | Static wrapper prose | Inventory is exhaustive but >40 lines and not a router; binary disagrees | **Replace normal root** with router; retain `commands` inventory |
| Query plane (`status`, `roadmap`, `learnings`, `maps`, `guards`, `diff`, `outline`, `anchors`, `commands`) | Inspect ambient state; repo → compact TOON → agent | Executable; AXI/Bench | Generally content-first, bounded, structured | **Keep**, reduce public maintenance queries |
| Test/canary/structure | Verify packages, fixture ownership, structural budgets; tree → evidence → gate/maintainer | Executable | 233 planted reasons bite; `bench test` polluted subject under documented `BENCH_HOME` | **Keep**, fix isolation; do not equate anchor tests with behavior |
| Coverage/preflight | Validate requirement membership and fences; spec/diff → verdict rows → build/review | Executable; Bench-native | Catches uncited/phantom/out-of-fence cases | **Keep and invoke automatically** |
| Gate/cache/pin | Run oracle and reuse exact evidence; tree+oracle → green/red/timeout record → commit/hook | Executable; Bench-native | Strongest proven control, including stale and moved-subject refusal | **Preserve unchanged behind new spine** |
| Commit/spec lifecycle | Gate prospective tree, commit named paths, flip/retire/history; paths/spec → Git transition → developer | Executable; Bench-native | Atomic landing is high value | **Keep** |
| Worktree/intent/shift | Own isolated work, execute/release/recover, loop to green; objective/base → assignment/evidence → agents | Executable + prompts; no-mistakes/Bench | Identity and failure-split tests pass; substantial state/command surface | **Keep deep module**, simplify public happy path |
| Handoff | Render pins/next action; repo + arbitrary State → Markdown → next harness | Executable renderer + prompt-authored body | Pins and prefixes recover; semantic freshness is false-positive prone | **Replace body** with schema-backed checkpoint |
| Dashboard | Render board; repository state → HTML → human | Executable | Useful observer surface, not workflow control | **Keep optional** |
| Models/line | List IDs and validate declared line; environment → advisory choice/check → orchestrator | Partial executable; Bench-native | Some conformance; cross-harness enforcement asymmetric | **Benchmark then complete or demote** |
| Release (`prep-release`, preflight, prepare/submit/promote/rollback/status) | Verify and govern publication; version/artifacts → evidence/state/publication → operator | Executable state machine; Bench-native | Deep, resumable control; workflow currently bypasses it with raw publish | **P0 wire as only publish path** |
| Skills index / anchors | Detect prose registry/anchor drift; files → conformance results → kit maintainer | Executable; Bench-native | Prevents silent deletion; creates public plumbing and rewards literal-text survival | **Internalize**, retain only proven contract anchors |

### Hooks

| Hook | Purpose; input → output → consumer | Enforcement / provenance | Observed value; cost / problems | Recommendation |
|---|---|---|---|---|
| Session start/resume | Reconcile worktrees and print status; session → dashboard → agent | Harness hook; Bench-native | Ambient recovery useful; large warning set does not route | **Keep**, feed router/context capsule |
| `block-dangerous-git` | Prevent unmanaged destructive Git; shell input → allow/deny → agent | Claude/Codex hook; Bench-native | Fail-closed when core absent; blocked read-only command containing `git` during cold audit | **Keep**, narrow parser and guarantee bootstrap path |
| `check-agent-line` | Enforce delegated model tier; agent call → allow/deny → orchestrator | Claude only | Claimed governance has no Codex parity | **Benchmark**, then wire symmetrically or remove |
| Stop | Refuse stopping armed shift on red gate; session state → deny/advice → agent | Claude/Codex | Valuable terminal-condition control | **Keep** |
| Pre-push | Require pinned gate evidence before push; ref/tree → allow/deny → Git user | Not installed in audit worktree | Strong design but inactive here | **Install transactionally or report definitive absence** |

### Workflows, state artifacts, and roles

| Component | Purpose; input → output → consumer | Enforcement / provenance | Observed value; cost / problems | Recommendation |
|---|---|---|---|---|
| Shape → spec → tickets | Convert ambiguity to bounded build; idea/decisions → coverage contract/slices → builders | Mostly prompt; Pocock/Bench | Strong artifacts when wide; no live specs, ticket sediment remains | **Route proportionally**, retire sediment |
| Implement/debug | Change behavior; ticket/symptom → tested diff → review | Prompt plus preflight/test/gate | Debug has causal loop; generic implementation can edit before reproduction | **Keep debug trigger explicit**, add empirical tripwire |
| Review → pickup → final | Semantic certification and closure; fixed diff/spec → findings/readiness → user | Prompt | Rubrics exist; no independent runner or executable Final differential | **Consolidate around review record + landing state** |
| Release | Rehearse/authorize/publish/recover; version → release record → operator | Mixed executable and workflow prose | State machine strong; raw-publish bypass is critical | **Make machine path exclusive** |
| Git/tree/gate evidence | Source and verification identity; content → immutable refs/records → all roles | Executable | Fresh, replayable substrate | **Keep canonical** |
| Intent ledger/objective scratch | Assignment ownership/objective; objective/base → handles → worktree/shift | Executable but narrow/ephemeral | Good concurrency identity; not general work state | **Reuse as inputs**, do not stretch into goal store |
| Specs/tickets/reviews/maps | Requirements, slices, pickups, decisions; prose/schema → files → phase agents | Mixed checks/prompt | Coverage rows useful; current orphan/ticket-only state and invalid map | **Migrate live facts**, archive/retire rest |
| Capture/roadmap/retros/scorecards | Backlog and learning; observations → prose rows → maintainer | Prompt plus query parsers | Institutional memory; high cold/context and staleness cost | **Keep off active path**, summarize by pointer |
| Session handoff | Resume metadata; tree + prose → handoff → fresh session | Commit-age check only | Pin recovery works; arbitrary state can lie | **Replace active semantics**, retain generated view |
| Human reviewer/operator | Makes judgment/authorizes risky transitions | Policy boundary | Necessary for semantic and publication decisions | **Keep explicit** |
| Orchestrator | Selects line, state, delegates, integration | Prompt role | Needed for wide work; lacks canonical broadcast state | **Back with work-state store** |
| Writer delegate | Implements one bounded slice | Prompt role + owned worktree | Isolation can help; mandatory fresh writer is unmeasured | **Use when dependency graph merits it** |
| Review-axis delegates | Rediscover Standards/Spec/Coverage findings | Prompt role | Axes useful; independence unverifiable | **Run in isolated contexts with retained metadata** |
| Shift agent | Iterates objective in owned worktree until green | Executable loop + model | Failure split observed | **Keep specialized** |
| Gate oracle | Decides development completion for configured subject | Executable | Exact and replayable within declared capability set | **Keep authority bounded to its tier** |
| Release operator | Authorizes ship/publication transitions | Mixed | Necessary external boundary | **Force governed release commands** |

## G. Pocock integration assessment

Bench did not merely rename Pocock's skills. It preserved recognizable causal practices and added a deterministic substrate.

**Preserved.** `[OBSERVED]` Frontier-round grilling, facts-versus-decisions separation, stable domain language, resulting-state ADRs, acceptance-led specs, vertical tracer tickets, seam selection, selective TDD, reproduce/minimize/hypothesize debugging, fixed-base review, sparse comments, and disposable prototypes remain legible in the current skills. The upstream Ask Matt flow also supports one context through grill/spec/tickets and a fresh implementation context, which resembles Bench's intended phase boundaries.

**Improved.** `[OBSERVED]` Bench adds acceptance-coverage rows, build/review fences, exact-base diffs, worktree ownership, preflight, gate evidence, prospective commits, recovery handles, and a release state machine. These additions turn several upstream recommendations into checkable preconditions or transitions. The passing adversarial tests in E007 are the strongest evidence of a real consistency gain.

**Drifted or weakened.** `[OBSERVED]`

- `diagnosing-bugs` is embedded in a larger worktree/shift phase. Its text says review should check a regression “it didn't author,” although the same phase authors it; that weakens the independence claim.
- Review adds a useful Coverage axis, but its pickup persistence route conflicts with its “no gate” rule, and “the gate and you decide done” conflicts with gate-only completion authority. `[E017]`
- Shape can produce a chart and stop where upstream frontier elicitation would continue to a resolved decision.
- Current upstream Ask Matt is an explicit router; Bench distributes routing among static inventory, status, handoff, and What Next.
- Current no-mistakes makes review turns explicitly isolated and treats implementer/fix narratives as claims. Bench says “fresh” without recording isolation.

**Duplicated.** `[OBSERVED]` The same phase identities appear in 11 skill adapters and 11 command documents, while cross-cutting rules recur across BENCH, profile, craft skill, phase skill, and conformance anchors. Some repetition is adapter glue, but command inventories, delegation rules, gate reminders, and phase-close narration carry duplicated knowledge that can drift. The 31,456-word estate is not all active at once, yet the mandated cold set is already about 6,000 words. `[E014–E015]`

**Boundary recommendation.** Judgment stays prompt-level: interrogating ambiguity, naming domain edges, selecting seams, evaluating semantics, ranking hypotheses, and choosing a proportionate slice. Deterministic software should own: prerequisite discovery, phase/state transitions, ownership, exact refs, coverage membership, freshness, attempt inheritance, context selection, reviewer isolation metadata, landing, and publication.

The observed deterministic gains survive any harness that invokes the same binary and repository checks. The prompt-level gains are not shown to survive. `[OBSERVED]` Claude has phase adapters and an Agent-line hook that Codex lacked in this run; therefore cross-harness consistency is partial, not established. `[E012, E015]`

## H. `diagnosing-bugs` assessment

### Causal mechanism

The load-bearing mechanism is not the prose length or the skill name. It is this forced sequence:

1. acquire an executable loop for the exact reported symptom before repair;
2. minimize the reproducer and map the execution path;
3. state a small set of falsifiable hypotheses;
4. run the cheapest safe discriminator, changing one variable;
5. carry the observed failure into the next attempt;
6. identify root cause before changing production behavior;
7. add a regression at the seam that can independently bite;
8. rerun both the new regression and the original signal;
9. clean instrumentation and fixtures.

`[REPEATED PRACTITIONER EVIDENCE]` The audit charge records repeated practitioner success breaking repair loops with the upstream practice. `[RESEARCH]` ReAct supports interleaving reasoning with observations, and Reflexion supports retaining explicit feedback across attempts. Neither proves that this particular wording is optimal.

### Audit evidence

`[OBSERVED]` Current upstream `diagnosing-bugs` is a focused 138-line reproduction/minimization/hypothesis loop. Bench's roughly 140-line debug command retains that sequence and adds isolated entry, preflight/shift, gate, and handoff integration. Focused gate and shift-failure tests passed. There was no live multi-trial model comparison, so the audit demonstrates textual and substrate preservation—not fewer repair attempts, higher root-cause accuracy, or better first-fix success. `[E007; controlled experiment matrix]`

### Dilution risks

- Worktree, line, delegation, gate, and handoff rules compete with the causal loop for attention.
- Saying a check “it didn't author” does not create independence inside the same phase.
- A generic implementation route can patch before reproducing unless the router recognizes failure-shaped intent.
- Conformance anchors can preserve the words while the model skips the behavior.
- A compressed replacement can destroy the mechanism if it removes exact-symptom reproduction or the original-signal rerun.

### Cross-harness behavior

`[OBSERVED]` The repository file is harness-neutral and Claude exposes a debug adapter; this Codex session did not expose the on-disk `bench-debug` phase adapter. A user can still ask for debugging or load the file, but automatic discovery parity failed. The deterministic tests/gate are portable; the trigger and prompt compliance are not demonstrated across models.

### Focused benchmark

Run at least 20 seeded bug instances per condition, stratified across deterministic unit failure, integration failure, flaky/concurrency failure, and underspecified report. Randomize model/harness and repository order; use identical budgets and hidden root-cause labels. Compare:

- no specialized skill;
- current Bench default route;
- current upstream `diagnosing-bugs` alone;
- Bench-integrated debug;
- proposed compact causal card plus executable attempt ledger.

Measure time/tool calls/tokens to red-capable loop, speculation after a discriminator becomes available, blank retry rate, failure inheritance, root-cause-before-fix rate, root-cause accuracy, first-fix success, regression mutation bite, original-signal rerun, cleanup, and total cost. Require confidence intervals and report failures, not just aggregate pass rate.

**Recommendation: PRESERVE the causal loop, CHANGE its delivery, DO NOT REPLACE it yet.** Put the trigger, attempt record, and original-signal requirement in the router/work-state substrate; keep hypothesis formation and seam choice as a compact skill; move examples and worktree mechanics behind contextual disclosure. Delete or rewrite only after the focused benchmark shows non-inferiority on root-cause and original-signal metrics with lower context cost.

## I. Run the Dang Test Findings

| Current behavior | Available observation | Why observation is superior | Current trigger | Recommended trigger | Recommended enforcement | Measurement |
|---|---|---|---|---|---|---|
| Root help/status are reasoned about as entry points | Invoke wrapper and exact binary with no args/help against representative states | Reveals actual route, exit, and schema | User already knows `help`, `status`, or What Next | Any normal `bench` invocation or ambiguous intent | Router contract tests over state fixtures | Route accuracy; clarification count; tokens to first useful action |
| Prompt phase is treated as a control | Invoke its subcommand and inspect retained state owner | Distinguishes narrative from executable transition | Skill/command prose is loaded | Claimed Review/Final/claim enforcement | Phase registry must declare `prompt`, `advisory`, or executable owner | Unsupported-control claim count |
| Gate freshness is discussed abstractly | Change tree/oracle/subject and attempt reuse | Directly tests stale authorization | Gate reuse path | Any reuse of verification evidence | Digest-bound record, fail closed on drift | Stale-evidence acceptance rate |
| Handoff is called fresh because HEAD matches | Compare every semantic field with live sources | Commit age cannot detect same-commit false prose | Commit-distance status | Resume or harness swap | Schema-derived fields with source digests; reject contradictions | Resume fact accuracy; stale-field rate |
| Requirement preservation is assumed from files | Delete/uncite/phantom one acceptance row and run preflight/coverage | Tests the omission that matters | Model remembers to call preflight | Spec handoff, build start, review start | Automatic preflight at phase transition | Requirement-loss detection timing |
| Reviewer independence is asserted | Run clean-context and narrative-seeded reviews over the same fixed diff | Measures anchoring and unique valid findings | Prompt requests “fresh” delegates | Semantic certification | Isolated context, retained inputs/identity, overlap report | Valid unique findings; false positives; overlap |
| Agent retries after failure | Compare the next action with and without prior diagnostics | Shows whether failure changed state | Shift-only recovery or model memory | Same failure class twice or patch after failed test | Attempt ledger and blank-retry tripwire | Blank retry rate; failure inheritance rate |
| Agent proposes a patch for a bug | Run the exact symptom or build the cheapest red-capable loop | Converts speculation to causal evidence | User/model explicitly selects Debug | Failure language, failing command, or second repair | Debug route blocks production edits until a red observation or explicit bounded exception | Red-loop acquisition; root-cause-before-fix |
| Local green is described as completion | Query capability skips and ship/publication state | Exposes untested boundaries | Final prose is expected to qualify it | Completion/publication claim | Typed tier claim with dependencies | Local-as-global misclaim rate |
| Skill text is assumed useful | Randomized task trials with and without skill | Measures outcomes and overhead | Skill addition plus conformance anchor | Add/retain/change any active skill | Benchmark gate for active-load growth | Pass delta, token delta, time, regressions |

**Empirical-escape doctrine:** When two active hypotheses predict different safe executable observations, run the cheapest discriminating observation now. Record subject, command, output identity, and what it falsified. A retry is allowed only with new evidence or a materially different hypothesis. For bugs, no production repair precedes an exact-symptom red loop unless the record names why reproduction is unsafe or impossible; completion requires rerunning the original signal.

`[OBSERVED]` Gate, preflight, coverage, canary, and shift already embody parts of this doctrine. Ordinary phase prompts do not. Review even prohibits tests during the review pass, so an immediately available discriminator can be deferred. The right boundary is to preserve reviewer non-mutation while allowing read-only/focused observation or a separately recorded probe.

## J. Prose audit

| Defect | Current evidence | Disposition |
|---|---|---|
| Duplication | Eleven phase skills plus eleven command documents; gate/delegation/phase-close rules repeated through profile, skills, and anchors | Generate adapters from one registry; compile active policy |
| Contradiction | Review says persist/commit pickup state, says no gate in phase, while public `bench commit` gates; it also shares done authority with the reviewer contrary to gate-only policy. `[E017]` | One state transition owner; make review advisory and gate authority explicit |
| Negative-instruction overload | Thousands of words prescribe what agents must not do across multiple layers | Turn safety-critical negatives into guards; delete baseline-model reminders after benchmark |
| Vague verbs | “deepen,” “assess,” “synthesize,” “fresh,” “independent,” and “done” lack uniform observable postconditions | Define typed outputs or keep them explicitly advisory |
| Cargo-cult reasoning | Fixed agent counts, tier ladders, “fresh writer” rules, and dense checklists can look rigorous without an outcome study | Benchmark or remove; never use count as a proxy for independence |
| Context dumping | About 6,000 words are mandated before task-specific guidance; 31,456 words exist across the guidance estate. `[E014]` | Compile a small active capsule with pointers |
| Hidden policy | `--full`, review, final check, next routing, and line governance are primarily Markdown semantics | Expose capabilities/state in the CLI or label them prompt-only |
| Stale prose | Handoff `State` contradicts live repository pins at the same current commit; ticket-only specs survive outside status. `[E011, E016]` | Derive state; archive historical artifacts |
| Model-specific contamination | Claude has phase symlinks and Agent-line enforcement; Codex exposure differs. `[E015]` | Core registry + generated adapters + parity tests |
| Output bloat | Root emits a 40-plus command catalog; status emits many findings but chooses by syntactic invocability | Root route first; disclose inventories on request |

Deletion candidates are the duplicate adapter bodies, prompt-only `--full` promise, fixed delegation counts, generic design-system content from non-UI active context, obsolete ticket-only spec directories, invalid maps, and stale handoff narrative. `[INFERRED]` Exact removals require the migration ticket and conformance updates; this audit did not alter them.

## K. Work-state assessment

Bench should **adopt and strengthen** a small canonical work state, but not canonize a decorative `Goal / Core / Verified / Open / Next` Markdown template.

Recommended schema:

```text
work_id, goal, status, owner
subject: repository, base, tip/worktree, scope fences
decisions: ids + source pointers
requirements: ids + coverage state
verified: claim ids with subject/tool/result/freshness
open: uncertainty + competing hypotheses
next: one action + why + prerequisites
attempts: action, observation, diagnostics, falsified hypotheses
feedback_loop: original signal + exact rerun
```

`Goal` and `Next` are singular routing facts. `Core` should be renamed to explicit accepted decisions/constraints so it cannot become a prose junk drawer. `Verified` must contain evidence references, not conclusions. `Open` should distinguish missing facts from unresolved judgment. `Attempts`, `Coverage`, and `Feedback loop` are mandatory for recovery but were absent from the five-field proposal.

The store should survive model and harness swaps but not preserve chain-of-thought, obsolete plans, persona, tool chatter, or falsified speculation. Git remains content truth; the work record indexes live sources rather than copying them. Handoff becomes a generated view. `[OBSERVED]` Current handoff recovers pins but preserves semantically stale arbitrary State, and the intent ledger owns assignments rather than the active objective. `[E011–E012]`

## L. Context compiler assessment

Bench has a file loader, not a context compiler. `[OBSERVED]` It makes broad policy available and relies on the model to select it; the cold mandate alone is roughly 6,000 words. A compiler must select context from current state and task intent, order it by immediate use, and provide pointers for zoom-in.

Smallest sufficient context by representative work:

| Work | Always include | Add on demand | Exclude initially |
|---|---|---|---|
| Bounded code change | Goal, base/tip, fence, named acceptance rows, relevant project standard, last fresh tests, next action | Seam/TDD guidance if boundary is uncertain | Roadmap, release, skill-authoring, unrelated craft skills |
| Reported bug | Exact symptom/original command, observed output identity, scope, attempts/diagnostics, next discriminator, debug causal card | Worktree recovery and seam reference when needed | Spec ceremony, design-system, broad assessment |
| Wide feature | Closed decisions, glossary edges, spec rows, dependency graph, current slice, ownership, fresh evidence | Tickets and specialist practice for current slice | Completed tickets' prose and full profiles |
| Independent review | Immutable base/tip, spec/decision sources, standards pointers, changed files, fresh machine evidence | One axis rubric per isolated reviewer | Implementer conclusion/fix summary; other reviewers' findings |
| Release | Version, artifact digests, tier/capability evidence, authorization, release-state record, next transition | Recovery runbook for active failure | Shape/spec/debug guidance |

Target a 200–400-token routing header and usually no more than 1,500 task-specific tokens before the first observation, with source pointers rather than copied prose. This is a benchmark target, not a universal hard limit. The compiler should emit deterministic sections (`goal`, `subject`, `fresh_evidence`, `open`, `next`, `sources`) and explain every included source. `[INFERRED]`

## M. Claim/evidence assessment

The gate record is the good model: a claim is about an exact subject, depends on an exact oracle, records an outcome, can be replayed, and becomes stale on drift. `[OBSERVED: E007–E008]` Extend the pattern without pretending all semantic claims are executable.

| Dimension | Current state | Required state |
|---|---|---|
| Provenance | Strong for Git/gate; weak for handoff prose, reviews, and Kunchenguid imports | Source IDs/revisions and author/tool for every durable claim |
| Freshness | Tree/oracle gate freshness strong; handoff only commit-distance | Dependency digests and semantic source derivation |
| Replay | Gate commands/logs replayable; prompt judgments not reproducibly isolated | Retain inputs, fixed refs, rubric/version, reviewer context metadata |
| Dependencies | Gate explicit; other artifacts imply relationships by path/name | Typed dependency graph and invalidation |
| Control transitions | Gate/commit/release commands have states; phases mostly narrative | Explicit allowed transitions with authority and failure posture |
| Bypasses | Raw handoff assertions and raw `npm publish` workflow path | Unsupported claims remain advisory; publication path denies raw bypass |
| Output | Query TOON is compact; root/help and prose artifacts are large | Content-first claim rows, aggregates, definitive empty, pointers |
| Failure | Gate fails closed; prompt phase can omit or retry blankly | Structured error with preserved diagnostics and next discriminator |

A minimal durable claim needs: `id`, `type`, `subject`, `outcome`, `source/provenance`, `observed_at`, `tool/rubric version`, `dependencies`, `freshness rule`, `replay recipe`, and `authority`. Semantic review findings can be evidence-bearing advisory claims without becoming “machine truth.” Completion claims must declare their tier; seven local capability skips cannot silently authorize a global/ship claim. `[OBSERVED: E008, E020]`

## N. Gate / Review / Final Check / CI matrix

| Control | Responsibility | Coverage | Timing / subject | Independence and authority | Overlap / gap |
|---|---|---|---|---|---|
| Focused tests / `bench test` | Fast behavioral feedback | Selected packages; capability skips reported | During implementation; current files/environment | Executable, but subject isolation bug observed under audit environment | Feeds gate; not completion |
| Preflight / coverage | Ownership fences and acceptance membership | Paths, citations, rows, spec lifecycle | Build/review entry at base/tip | Executable; only bites when invoked | Should be automatic transition prerequisite |
| Development Gate | Decide configured repository completion | gofmt, vet, test, race, system, shellcheck here; 7 capability skips | Exact work tree + exact oracle; before landing | Deterministic authority for declared dev tier | Does not certify semantic satisfaction, clean board, native/ship tiers |
| Review | Find semantic Standards/Spec/Coverage defects | Fixed diff plus sources, in prose | After implementation, before fixes/landing | Advisory; “fresh” axes not recorded or enforced | Needed judgment; persistence and commit route contradictory |
| Final Check | Report landing/readiness and housekeeping | Status, capture, retrospective, handoff | After work/review/gate | Prompt-only; no executable subcommand or distinct evidence | Mostly overlaps status/gate/landing; no reproduced differential |
| Native CI | Re-run repository checks in hosted capabilities | CI workflow matrix, environment-specific capabilities | Pushed commit/ref | External executable evidence; not exercised here | Pre-push absent; local skips may differ |
| Ship/release checks | Artifacts, cross-build, authorization, publication | Prep/release preflight and state machine | Version/artifact/registry | Deterministic transitions plus human authorization | Raw `npm publish` workflow bypasses governed submit/promote path |

`[OBSERVED]` A fresh gate passed in 59.8 seconds and reused its exact verdict; the same tree still had 62 structure issues, seven unresolved maps, one dirty audit path, and seven capability skips. That is acceptable only if outputs say precisely “development gate green for this oracle,” not “all requirements globally complete.” `[E008–E009]`

## O. CLI / AXI assessment

Bench's read plane is substantially AXI-shaped; its root and workflow plane are not.

| AXI property | Actual behavior | Verdict |
|---|---|---|
| Token-efficient defaults | `status`, `diff`, `maps`, `guards`, `coverage`, and `preflight` use compact records | **Good** |
| Minimal schema | Selected queries avoid prose, but root prints a full 40-plus command catalog | **Mixed** |
| Truncation / aggregates | Status and roadmap summarize counts/top rows; full detail is optional in several commands | **Good** |
| Definitive empty | Query commands generally state absence; missing binary exits 127 clearly | **Mostly good** |
| Structured errors / exits | Exact binary rejects unknown commands with exit 2; preflight/gate have structured outcomes | **Good**, except prompt phases look like capabilities but are unknown |
| Ambient context | Status detects branch/tree/worktree/gate state | **Good**, but does not own user intent or active work |
| Content first | Query plane yes; wrapper help and instruction loading no | **Mixed** |
| Contextual disclosure | `--full`/help exist in parts; public surface still exposes maintenance internals | **Mixed** |
| Consistent help | Wrapper accepts `help`; exact binary rejects it; no-arg behavior differs | **Fail** `[E004, E019]` |
| Side-effect clarity | Most queries are read-only; `bench test` unexpectedly created ten worktree fixtures under the subject with documented environment variables | **Fail in tested case** `[E006]` |

Every registered command appears in exactly one sweep row below. “Bounded” means default output is compact/deterministic; command-specific `--help` remains the supported grammar discovery path.

| Commands | Discovery/default/help/output/exit | Safety, agent cost, failure mode | Owner / change |
|---|---|---|---|
| no args, `help` | Wrapper: full inventory/0; binary: no subcommand/2; `help` also disagrees | High discovery tokens; behavior depends on entry executable | Core registry: route on no args, consistent help |
| `commands --brief`, `version` | Three-command liveness probe and version string; bounded | Name can be mistaken for inventory | Keep explicit probe wording |
| `link`, `unlink`, `upgrade` | Manifest-driven install lifecycle; help exists | Mutating but planned; broad footprint risk | Installer: keep transactional/dry-run behavior |
| `init` | Gate scaffold | Name suggests broader product initialization than behavior | P0 root/init contract: own minimal transaction |
| `doctor`, `repair` | Diagnose/fix PATH/pinned binary; explicit mutation flags | Cold checkout had no pinned binary; round trip before real commands | Installer: improve bootstrap error/repair route |
| `models` | Advisory IDs; bounded | Model freshness and line benefit unstable | Profile owner: benchmark, version |
| `structure` | Size/crowding findings | 62 advisory issues; size can masquerade as architecture judgment | Query owner: keep advisory, profile thresholds |
| `skills-index`, `anchors` | Prose registry and literal anchors | Public maintenance surface; encourages text-presence confidence | Move to internal conformance namespace |
| `idea`, `roadmap`, `learnings`, `maps` | Capture/backlog queries; bounded top/TOON views | Historical state can crowd active work; one invalid map | Maintenance plane, outside normal route |
| `status` | Ambient board + next action | Strong facts, but first-invocable heuristic skips higher prose severities | Feed P0 router; stop claiming semantic next action |
| `handoff` | Rewrites pins and preserves State; harness option | Same-commit stale semantics, opaque manual prefix | Generate from canonical state |
| `dashboard` | HTML/stdout board | Optional human cost; no authority | Keep view-only |
| `canary` | Fixture registry validation | 233 bindings prove planted mutation ownership, not behavior | Keep internal/gate evidence wording precise |
| `guards` | Deny-surface inventory | Useful ambient safety; warning is not installed control | Keep, add wired-state source |
| `diff` | Frozen review base and optional full diff | Good bounded default; full output intentional | Keep |
| `coverage` | Acceptance rows and `--check` | Strong for live spec; ticket-only slug fails “spec not found” | Keep; migration must retire/orphan old tickets |
| `preflight review`, `preflight build` | One verdict row per phase-entry check | Strong but model must invoke it | Make transition prerequisite |
| `test` | Fresh package tests, failures/skips | Passed but created ten fixture dirs in subject under documented environment | Fix isolation; retain structured skips |
| `outline` | Candidate seams by file/line | Correctly refuses semantic seam judgment | Keep observation-only |
| `gate`, `gate --fresh`, `gate pin` | Exact oracle, cache, pre-push pin | Strong; “green” can be overgeneralized beyond tier | Preserve; type capability/tier in claim |
| `prep-release`, `release-preflight` | Artifact/cross-build/auth rehearsal | External capabilities may skip; bounded safe verify modes | Keep ship evidence explicit |
| `release prepare`, `submit`, `promote`, `rollback`, `status` | Governed resumable npm state | Strong but bypassed by raw workflow publication | P0 exclusive publish path |
| `worktree` create/release/clean/refresh/help | Owned isolated lifecycle | Deep grammar and coordination cost | Keep happy path; disclose recovery on demand |
| `worktree list`, `path`, `exec`, `reauthorize` | Structured lookup/execution/recovery | Valuable exact ownership; expert-only complexity | Keep deep module, contextual help |
| `shift` | Gated pooled loop | Preserves one failure class; can become autonomous long loop | Keep specialized with attempt state |
| `commit` | Prospective exact gate then named-path commit | High-value, mutating, requires explicit paths | Preserve |
| `spec implemented`, `retire`, `history` | Validated lifecycle/history | Current live spec set empty; ticket sediment outside model | Keep lifecycle; add orphan migration |

| Current behavior | Agent cost | Failure mode | Proposed owner | Proposed change |
|---|---:|---|---|---|
| Wrapper/binary root disagreement | Extra probe and branch in reasoning | Wrong grammar/exit assumed | Core capability registry | Generate both dispatch/help paths from one contract |
| Full inventory as no-arg default | 40-plus lines before intent is known | User picks plausible but nonexistent phase command | Logical router | Return one route; inventory on explicit request |
| Compact TOON read queries | Low | Some advisory rows mistaken for controls | Query plane | Keep schema; tag authority/freshness |
| Prompt-only phase names | Filesystem/skill archaeology | Formal prose mistaken for executable capability | Phase registry | Declare capability kind and direct adapter |
| Manual preflight invocation | One remembered round trip | Ownership/requirement issue detected late | Transition engine | Invoke automatically at build/review entry |
| `bench test` fixture leakage | Cleanup round trip and dirty tree | Verification mutates its subject | Test runner | Isolate fixture root and assert clean-subject postcondition |

The most costly ergonomic flaw is that an agent must know whether it is talking to the shell wrapper, exact binary, prompt adapter, or prose-only phase. `bench review` looking plausible but returning “unknown subcommand” is a schema failure at the product boundary. The root should return a small route object such as:

```text
state: interrupted
reason: active work record has a failed observation and no fresh successor
next: bench resume <work-id>
needs: none
experts: status, debug, gate
```

Keep public commands whose semantics users need: init/link, route/status, test, preflight/coverage, gate/commit, worktree/shift, release, and a discoverable expert inventory. Internalize `skills-index`, literal anchors, cache plumbing, and lifecycle repair details unless an external consumer is demonstrated. `[INFERRED]`

## P. Orchestration assessment

**Orchestrator role.** `[INFERRED]` A wide task needs one owner for decomposition, immutable bases, dependency order, broadcast facts, acceptance integration, and escalation. It does not need a narrator that restates every skill. Back the role with canonical work state and make every assignment reference a work ID and requirement slice.

**Delegates and ownership.** `[OBSERVED]` Worktree assignments, fences, base identities, execution, release, and recovery are deep, test-backed modules. Keep them. Fixed numbers—six assessment agents, three review axes, a required writer for each ticket, and sometimes three alternative designers—are not evidence of useful parallelism. Delegate only independent graph nodes whose latency or context isolation justifies coordination cost.

**Broadcast state.** Broadcast canonical names, goal, accepted decisions, immutable refs, scope, requirements, fresh evidence, active failure diagnostics, and next action by ID/pointer. Do not broadcast implementer conclusions to an independent reviewer or full prior conversation to a writer. `[INFERRED]`

**Independent review.** A fresh model call is not enough. Record fixed base/tip, sources, rubric version, provided context hash, model/harness identity when available, and absence/presence of implementer narrative. Run axes in clean contexts; merge cited findings afterward. The reviewer remains advisory and does not mutate the subject during discovery.

**Failure inheritance.** `[OBSERVED]` Shift preserves progress and staging-fault recovery evidence; ordinary phases do not have an attempt schema. Every failed action should append observation identity, diagnostics, and affected hypothesis. A same-class retry with no new evidence should be rejected or explicitly waived. `[E007]`

**Integration.** The orchestrator may combine green bounded slices only after revalidating dependency edges and the integrated subject. The exact gate, not delegate reports, authorizes the configured dev transition. Publication requires the release state machine and human authority. The present guidance has the right conceptual roles but asks prose and fixed agent counts to stand in for state and isolation.

**Ticket/spec-build audit.** `[OBSERVED]` Ticket templates carry scope, ownership, dependencies, assumptions, requirement IDs, acceptance, tests, validation, escalation, and bounded output guidance. Preflight detects out-of-fence paths and missing/phantom citations. Yet no live `specs/*/spec.md` exists, roughly forty ticket files remain, ticket-only directories are ignored by spec status, and coverage reports `spec not found`. Partial/provisional work can persist in Git/worktrees, blocked work in prose, and shift failures in specialized evidence; no general start → assign → checkpoint → integrate → review → promote machine exists. Whole-project acceptance is therefore not owned at exactly one semantic transition: coverage is mechanical membership, Review is advisory satisfaction, and Gate is configured code verification. `[E007, E010, E016]`

**Dogfood results.**

| Intended experience | What was actually exercised | Result |
|---|---|---|
| Normal front door | Wrapper/binary no args/help/status | Inventory or error, no intent route |
| Shape/spec/work units | Inspected all phase contracts; queried current maps/spec/coverage; ran preflight/coverage adversarial owners | Prompt procedures plus real checks; no live spec to traverse |
| Controlled implementation | Ran repository-owned mutation/ownership/gate fixtures rather than editing Bench, per audit prohibition | Deterministic enforcement reproduced without altering implementation |
| Diagnose bug | Compared current upstream/integrated mechanisms; exercised shift/gate failure owners | Mechanism visible; no repeated model outcome trial |
| Gate/evidence replay | Fresh full gate, reuse, red/timeout/moved-subject tests | Reproduced strongly |
| Review/Final/claims | Invoked each plausible subcommand; inspected prompt/state owners | Unknown commands; controls unreproduced/absent |
| Checkpoint/resume/model swap | Rewrote Claude/Codex handoffs in isolated clones and compared live state | Pins/prefixes partial; stale State preserved |
| Release | Inspected commands/workflow and prior assessment; no external publication | Governed machine exists; workflow bypass confirmed statically |

The dogfood friction was itself diagnostic: the Git guard needed a locally built core before even read-only Git text was accepted, the wrapper and binary disagreed, the documented `bench test` environment polluted the subject with fixtures, and audit state could not be expressed except as untracked prose. `[E001–E006]`

## Q. Model/harness boundaries

| Layer | Should contain | Current leakage / finding | Recommendation |
|---|---|---|---|
| **Bench core** | State schemas, router decision inputs, ownership, preflight, coverage, queries, gate/cache, landing, release, adapter registry | Root dispatch differs by wrapper/binary; semantic phases masquerade as capabilities | One executable registry and typed capability metadata |
| **Claude Code adapter** | `/bench` and direct phase shims, hook syntax, prefix rendering | Phase symlinks and Agent-line hook exist; behavior is richer than Codex | Generate from core registry; parity-test only intended differences |
| **Codex adapter** | `$bench`/skill metadata, Codex hook syntax, prefix rendering | Project phase adapters existed on disk but were absent from this session; no Agent-line parity | Make installation/discovery observable and testable |
| **Model-specific profile** | Proven model limits, tier IDs, effort/cost parameters, benchmark results | Generic workflow policy and fixed line rules can leak into model choice; runtime model build ID unavailable | Keep optional, versioned, and data-backed; no domain truth |

State that must cross boundaries: goal, refs, fences, decisions, requirements, coverage, evidence/freshness, open uncertainties, attempts/diagnostics, next action, and recovery handles. State that should not cross: private reasoning, persona, conversational narration, obsolete plans, falsified speculation, implementer conclusions presented as review facts, and full inactive guidance.

`[OBSERVED]` Handoff rendering translates derived phase prefixes between Claude and Codex, but an explicit `--next /bench-debug` is left untouched and arbitrary State remains the same. That is correct for opaque user input but proves the document is not a semantic migration mechanism. `[E012]`

## R. Start-over decision

| Option | Benefit | Risk | Judgment |
|---|---|---|---|
| A. Continue incrementally | Lowest immediate migration cost | Adds another layer to 31k words and five overlapping state systems; preserves ambiguous entry and bypasses | **Reject as default** |
| B. Architectural consolidation / strangler | Preserves proven deep modules while replacing root, state, and context one slice at a time | Temporary adapters and dual reads; requires disciplined deletion | **Recommend** |
| C. Full rewrite | Clean conceptual surface | High chance of losing exact-tree/oracle binding, prospective commits, ownership/recovery, hostile fixtures, release evidence, and debug causal loop | **Reject** |

| Decision criterion | Incremental | Strangler | Rewrite |
|---|---|---|---|
| Risk to demonstrated consistency | Low per edit, cumulative drift high | Low if old deep modules remain oracle | High |
| Migration complexity | Superficially low, never closes dual concepts | Medium and explicit compatibility window | Very high |
| Sources of truth | Preserves all current duplication | Can make new state canonical and old views derived | Starts one source but needs risky data migration |
| Testability | Existing local tests; architecture stays hard to compare | Old/new routes can run side by side on fixed fixtures | New tests may miss historical hostile cases |
| Harness compatibility | Existing asymmetry persists | Registry can generate both while old adapters remain rollback | Must rebuild adapters at once |
| Rollback | Per-edit only | Route flag/adapter rollback to existing workflow | Expensive, likely branch-level |
| Time to user value | Small fixes quickly, core confusion remains | One tracer delivers entry/resume value early | Long |
| Benchmarkability | Diffuse changes confound comparison | Direct A/B/C per migrated slice | Old/new comparable only after large build |

The seam is clear enough for a strangler: put a new logical router, canonical work record, context compiler, and typed evidence view in front of existing status/preflight/gate/worktree/release modules. Read old handoff/spec/intent data during migration; emit the new state as canonical; turn old documents into generated views; delete adapters and sediment only after equivalence and outcome benchmarks pass. `[INFERRED]`

## S. Complexity kill list

Complexity budget (ordinal, based on observed mechanism and maintenance/context surface):

| Major capability | Demonstrated benefit | Complexity | Benefit / complexity | Class |
|---|---:|---:|---:|---|
| Exact gate + prospective commit | 5 | 4 | High | Keep |
| Ownership/worktree/recovery | 5 | 5 | High for parallel/risky work | Keep; simplify surface |
| Preflight/coverage | 4 | 3 | High | Strengthen |
| Compact query plane | 4 | 2 | High | Keep |
| Debug causal loop | 4 (practitioner + mechanism) | 2 | High, pending causal benchmark | Keep / benchmark delivery |
| Release state machine | 4 | 5 | Medium-high at ship boundary | Keep; remove bypass |
| Specs/tickets | 3 for wide work | 4 | Medium, low for bounded work | Route proportionally |
| Review axes | 3 plausible, 1 demonstrated control | 3 | Unknown | Benchmark; execute isolation |
| Status/handoff/What Next combination | 2 | 4 | Low current ratio | Merge/replace active-state path |
| Fixed delegation/line ceremony | 1 unmeasured | 4 | Low/unknown | Benchmark then simplify/delete |
| Literal skill anchors and duplicate adapters | 1 conformance | 3 | Low | Internalize/generate |
| Prompt-only Final Check | 1 unreproduced | 2 | Low | Merge |

| Action | Candidate | Why |
|---|---|---|
| **DELETE** | Prompt-only `bench-implement-spec --full` claim | It is narration, not a supported orchestration capability. |
| **DELETE** | Retired ticket-only spec directories after migration/archive | They are invisible to spec status and coverage and misrepresent current work. `[E016]` |
| **DELETE** | Invalid/dead decision maps after owner review | They create warning load without a resolvable source. `[E009]` |
| **DELETE** | Stale handoff `State` as an authored active source | It is semantically false while freshness reports current. `[E011]` |
| **MERGE** | Eleven command documents and matching phase adapters | Generate thin harness adapters from one registry. |
| **MERGE** | Final Check into landing/status/checkpoint | No independent executable differential was reproduced. |
| **MERGE** | What Next routing hints into root router; leave roadmap drain separately named | Current name invites the wrong entry behavior. |
| **SIMPLIFY** | Root command inventory | Route first; inventory on request. |
| **SIMPLIFY** | Review rubric/pickup protocol | Keep semantic axes; one unambiguous advisory-record transition. |
| **SIMPLIFY** | Debug delivery | Protect causal loop, disclose worktree/gate mechanics only when needed. |
| **SIMPLIFY** | Setup Repo | Deterministic init first; optional architecture interview separately. |
| **MOVE** | Design-system guidance | Load only for UI-profile tasks. |
| **MOVE** | Skill-authoring, synthesis, and kit-update guidance | Kit-maintainer profile, not linked-project cold context. |
| **MOVE** | `skills-index`, anchors, and registry repair commands | Internal maintainer namespace unless outside use is proven. |
| **BENCHMARK FIRST** | Mandatory delegation, fixed agent counts, fresh writer per ticket | Coordination cost and outcome benefit are unmeasured. |
| **BENCHMARK FIRST** | Model-line/effort ladder | Enforcement is asymmetric and model/scaffold interactions change. |
| **BENCHMARK FIRST** | Any rewrite or large skill compression | Could erase the gate/debug mechanisms that actually bite. |

## T. Missing capabilities

Only after the kill-list consolidation, add these in leverage order:

Capability-realization loss map:

| Stage | Potential capability loss | Current Bench mitigation | Evidence it works | Cost | Remaining gap |
|---|---|---|---|---:|---|
| Entry | User/model cannot select a useful path | Static root inventory, status, handoff, skills | Commands execute, but routing gap confirmed | High discovery/context | No intent router |
| Intent | Objective/scope become inferred or stale | Idea capture, intent ledger, handoff | Assignment identity works | Multiple stores | No canonical active goal |
| Shaping | Assumptions harden into design | Grill/domain/shape practices | Mechanism inspection only | Prompt turns/tokens | No marginal outcome evidence; can over-shape |
| Specification | Edges/requirements omitted | Spec schema and acceptance rows | Missing/phantom row tests bite | Authoring/maintenance | Semantic adequacy advisory; no live specs |
| Ticketing | Global requirement drifts in local slice | Row citations, fences, preflight | Adversarial tests bite | Duplicate prose/artifacts | Ticket-only sediment; mandatory writer cost |
| Active context | Relevant fact buried or omitted | Mandatory files and skill triggers | Availability observed | ~6k cold words | No selection/compiler |
| Implementation | Scope drift or premature edit | Worktree ownership, preflight, TDD advice | Ownership/fence checks pass | Deep lifecycle | Checks are late-bound; generic path can patch before repro |
| Debugging | Repair loop/speculation | Integrated diagnosing-bugs | Causal constraints preserved; practitioner evidence | Low-medium prompt load | Trigger/parity/outcome delta unmeasured |
| Tools | Agent misuses or cannot discover interface | TOON queries and guards | Query behavior reproduced | Broad command surface | Root/wrapper/binary inconsistency; test side effect |
| Evidence | Assertion mistaken for observation | Gate records/logs | Freshness/reuse tests strong | Cache/schema | No general typed semantic claim/attempt state |
| Gate | Stale/local result authorizes wrong subject | Tree+oracle identities, fresh option | Strongly reproduced | ~60s full run here | Capability skips/tier can be overgeneralized |
| Review | Model misses semantic defect or is anchored | Three axes/fixed diff guidance | Rubric only | Multiple agents/tokens | Independence and resolution bypass unverified |
| Integration | Green slices conflict or incomplete whole | Prospective commit, gate, worktree bases | Exact landing tests pass | Coordination | Semantic whole-work coverage ambiguous |
| Resume | New model repeats work or trusts stale state | Handoff pins/State | Pins refresh | Full prose reload | Semantic freshness/attempts missing |
| Model swap | Adapter difference changes behavior | Harness prefixes and shared repository files | Prefix tests partial | Duplicate adapters/hooks | Phase discovery and line-hook parity fail |

Highest-leverage losses are entry/intent/active-context/resume because one canonical state-and-router spine addresses all four; evidence/attempt state is next because it converts observations into safe transitions. Review becomes worth automating only after those inputs exist.

1. **Canonical work state plus intent router and context compiler.** This closes entry, recovery, active-context, failure-inheritance, and harness-swap gaps together.
2. **Typed claim/evidence/freshness graph.** Reuse gate evidence; distinguish advisory semantic findings from transition-authorizing proof.
3. **Executable independent-review runner.** Isolate inputs, retain provenance, merge cited findings, and forbid implementer narrative by default.
4. **Empirical-escape/attempt tripwire.** Detect blank retries and require a discriminator/original-signal loop for bug-shaped work.
5. **Exclusive governed publication path.** The capability exists; workflow and hooks must make bypass unavailable.
6. **Harness parity contract.** Generate adapters, inspect installed discovery, and test intended Claude/Codex differences.

More skills, more dashboards, a second planning ontology, latent-state branding, and automatic multi-agent expansion are not missing capabilities.

## U. Recommended architecture

```text
Claude /bench       Codex $bench       shell bench
       \                 |                 /
        +-------- logical Bench entry ----+
                         |
               intent + deterministic status
                         |
                  canonical work state
          (goal, subject, scope, requirements,
           evidence, open, attempts, next)
                         |
                   context compiler
             /           |             \
       judgment      executable       direct expert
       practices      transitions       commands
  grill/spec/debug   preflight/gate    status/test/
  seams/review       commit/release    worktree/etc.
             \           |             /
                typed evidence graph
                         |
             generated handoff/dashboard
```

The architecture has one root, one active-state owner, one source registry, and explicit boundaries between judgment, observation, and authority. Existing gate, commit, worktree, coverage, and release packages remain deep modules. Markdown skills become small decision aids selected by the compiler. The generated handoff is a view, not a second database. `[INFERRED]`

## V. Prioritized roadmap

| Priority | Problem | Evidence | Proposed change | Benefit | Cost | Dependency | Proof of success |
|---|---|---|---|---|---|---|---|
| **P0** | No canonical entry, active state, or selective context; resume can preserve false prose | E011–E019 | Build the single tracer in Z: logical root router + minimal work record + compiler + generated handoff, reading old sources compatibly | One entry; reliable resume; less context; basis for later controls | High design/migration risk, medium implementation | Existing status, handoff, intent, spec and gate readers | Six representative state fixtures route correctly; stale handoff cannot authorize; Claude/Codex render same semantics; A/B/C meets W thresholds |
| **P0** | Release workflow bypasses governed publication with raw `npm publish` | E020 and prior repository assessment | Make `bench release submit/promote` the only documented/automated publication path; guard direct path in release workflow | Prevents unsupported ship transition; preserves resumability | Medium; external CI/registry validation needed | Existing release state machine and authorization profile | Dry-run/fixture proves direct publish denied, resumed submit/promote succeeds, rollback/status evidence retained |
| **P1** | Review independence and persistence are prose assertions | E010, E017; no-mistakes comparator | Add an isolated review runner and typed advisory finding record over fixed refs; keep three judgment rubrics | Verifiable isolation and replay without making review the oracle | Medium-high; model invocation adapters | P0 work/evidence IDs and adapter registry | Seeded defects: isolated reviewers' inputs are auditable; no narrative leakage; valid unique-find rate non-inferior; pickup transition unambiguous |
| **P1** | General claims, attempts, and failure inheritance have no owner | E010–E012; shift is only partial exception | Add typed claim/evidence/attempt records using gate record semantics; reject blank retry without new observation/waiver | Freshness, replay, empirical escape, reliable model swap | Medium | P0 schema and current gate evidence | Tree/source mutation invalidates dependent claims; controlled repeat inherits diagnostics; unsupported prose remains advisory |
| **P2** | Wrapper/binary/harness registries drift and Codex adapters are absent | E004, E015, E019 | Generate help and thin Claude/Codex adapters from one capability registry; parity-test intended differences | Predictable discovery, fewer duplicated anchors | Medium | P0 logical route and install lifecycle | Wrapper/binary help contract identical; installed Claude/Codex inventory contains same logical phases; conformance deletion bites one owner |
| **P3** | Board/dashboard shows signals but not causal active state | E003, E009, E018 | Render work/evidence graph in dashboard after schema stabilizes | Human observability without adding authority | Low-medium | P0/P1 | Dashboard is a pure generated view; empty/stale/blocked states render definitively |
| **EXPERIMENT** | Skill, line, delegation, and debug marginal value is unknown | E013–E015 and research in X | Run W, including focused debug and context-compiler arms, before changing active guidance | Evidence-based retention and model/profile tuning | High evaluation compute/time | Frozen task set, instrumentation, multiple model/harness runs | Confidence intervals for outcome and resource metrics; preregistered keep/change/delete decisions |
| **DELETE** | Duplicate adapters, ticket sediment, invalid maps, stale handoff body, fixed-count rules, and irrelevant cold guidance consume attention | E009, E011, E014–E017 | After P0/P2 equivalence, archive or delete exact owner-reviewed set and remove anchors from the same source registry | Smaller surface and one source per fact | Medium migration/review; risk of deleting hidden behavior | P0 generated views, P2 registry, EXPERIMENT for behavioral rules | Cold context and guidance size materially fall; gate/canary green; benchmark non-inferior; no live source points to removed artifacts |

## W. Benchmark plan

### Main A/B/C comparison

- **A — baseline:** model + repository instructions, no Bench workflow.
- **B — current:** current Bench at the audited commit, including current root/discovery and skills.
- **C — proposed:** the architecture in U with the Z tracer, compact selected practices, and existing deterministic substrate.

Use a frozen, hidden-answer suite of at least 60 repository tasks: 15 bounded repairs, 15 reported bugs, 10 wide changes with requirements, 10 interrupted/resume tasks, and 10 review/release-state tasks. Include small and large repositories, hostile filenames/input, ambiguous reports, stale evidence, ownership violations, and requirement omissions. Every condition receives the same repository commit, budget, tools, and safe external boundary. Randomize task order and condition; run at least three seeds across at least two model families and both Claude/Codex adapters where available. Separate model adaptation runs from evaluation runs.

### Focused debugging comparison

- **D:** current Bench with no specialized debugging skill/route.
- **E:** current Bench-integrated debug.
- **F:** current upstream `diagnosing-bugs` alone.
- **G:** proposed compact causal card plus executable work/attempt state.

Use at least 20 instances per arm and include unit, integration, concurrency/flaky, configuration, and misleading-symptom bugs. Seed a known discriminating test and known root cause without exposing them. Mutate or remove each authored regression to prove it bites.

### Measures

| Dimension | Measures |
|---|---|
| Outcome | Task success; exact requirement satisfaction; hidden-test pass; regression quality; semantic defect escape; first-fix success |
| Epistemic quality | Root-cause accuracy; root-cause-before-fix; observation opportunity delay; speculation after discriminator; unsupported claim rate; stale-evidence use |
| Context efficiency | Input/output tokens; active guidance tokens; useful-source retrieval; lost-requirement rate; token overhead at equal outcome |
| Execution efficiency | Wall time; tool calls; model turns; repair attempts; gate runs; cost; time to first red-capable loop |
| Recovery | Resume fact accuracy; time to next useful action; failure inheritance; blank retry rate; original-signal rerun; cleanup rate |
| Orchestration | Delegate utilization; wait/coordination time; conflicts; valid unique review findings; reviewer anchoring; integrated failures |
| Complexity | Commands/concepts needed; install failures; context size; state owners; bypass count; maintenance fixture count |

### Analysis and decision rules

Pre-register task success and exact-requirement satisfaction as primary outcomes. Report per-task paired differences, confidence intervals, model/harness interactions, capability skips, and full cost—not only means. Treat a skill as retained active guidance only if it improves its target outcome or materially reduces severe failure without an unacceptable resource regression. For C, require non-inferior task success (lower confidence bound no worse than 3 percentage points), at least 30% lower active guidance tokens, at least 50% fewer stale/unsupported state fields on resume, and no regression in gate/ownership/coverage enforcement. For compact debug G, require non-inferior root-cause accuracy and original-signal rerun, lower observation delay, and no higher blank-retry rate. These thresholds are proposals to preregister, not observed results.

## X. Research reconciliation

| Source | Claim | Bench support | Bench contradiction | Transfer limit | Unproven question |
|---|---|---|---|---|---|
| [Verbalizable Representations Form a Global Workspace in Language Models](https://transformer-circuits.pub/2026/workspace/index.html) / [paper](https://arxiv.org/abs/2607.15495) | `[RESEARCH]` Selective, reportable, deliberately maintained representations can act as a global workspace and broadcast to multiple processes | Gate/status expose compact shared facts; proposed state/compiler matches selectivity | Bench treats thousands of available words as usable active context | Neural representation findings do not establish a textual workflow schema | Which compiled Bench facts improve model decisions, and by how much? |
| [J-CoT](https://arxiv.org/abs/2607.21981) | `[RESEARCH]` Selectively propagate vocabulary-indexed intermediate information rather than all prior language or the whole hidden state | Small evidence/next-action records are analogous at an interface level | Handoffs and profiles carry broad prose | Preprint latent-computation mechanism is not a mandate for latent reasoning or Goal/Core Markdown | What textual selection policy is best for repository work? |
| [DeepSeek V4 × J-Space Capability Realization Report](https://github.com/Tiger3807861189/DeepSeek-V4-J-Space-Capability-Realization-Report) | `[RESEARCH]` Community report hypothesizes strong interface/path/context sensitivity | Wrapper/binary and harness-adapter differences provide concrete sensitivity cases | Bench has no paired outcome evidence for its J-space-like claims | Single-run/community evidence without robust intervals is not universal model evidence | Do its reported gains replicate on Bench tasks and current models? |
| [J-Space Cognition Suite](https://github.com/Tiger3807861189/J-Space-Cognition-Suite-V3.6) | `[RESEARCH]` Proposes persistent Goal/Core/Verified/Open/Next, checkpoints, empirical escape, failure inheritance, and coverage | Gate records, shift recovery, and the proposed work schema realize some mechanisms | Current handoff is arbitrary prose; attempts and claim graph absent | Personas, dense rails, cognitive branding, and fixed labels may be cargo cult | Which individual mechanism survives ablation and cross-model trials? |
| [ReAct](https://arxiv.org/abs/2210.03629) | `[RESEARCH]` Interleaving reasoning, action, observation, and updates improves grounded task solving | Debug, gate, preflight, and shift create observation loops | Shape/review/maintenance can reason at length without a discriminator | Benchmark environments differ from full software repositories; not every judgment has an executable action | Does an empirical tripwire reduce repair cost without premature testing? |
| [Reflexion](https://arxiv.org/abs/2303.11366) | `[RESEARCH]` Verbal feedback and episodic memory can improve later attempts | Shift preserves some failure/recovery evidence | Ordinary phases lack structured attempt inheritance and allow blank retries | Benefit of reflective prose does not prove long self-authored narratives are the best state | Is a small diagnostic record better than free-form reflection for repair? |
| [SWE-agent](https://arxiv.org/abs/2405.15793) | `[RESEARCH]` The agent-computer interface materially affects repository-agent performance | TOON queries, structured exits, worktree/gate commands are agent-oriented | Root inventory, wrapper/binary mismatch, and prose-only commands burden discovery | Results depend on model, tools, and task distribution | Does the proposed root/schema improve success at lower tokens? |
| [SWE-bench](https://arxiv.org/abs/2310.06770) | `[RESEARCH]` Real issue resolution requires repository-level localization, modification, and validation | Bench owns validation and repository state deeply | Current internal conformance is not an end-task benchmark | Benchmark contamination, environment reproducibility, and issue distribution limit direct generalization | How does Bench compare on a controlled, uncontaminated repository suite? |
| [Agentless](https://arxiv.org/abs/2407.01489) | `[RESEARCH]` A simple localization → repair → validation pipeline is a strong baseline without elaborate autonomous planning | Bench's direct queries/test/gate can implement this baseline | Mandatory shape/spec/ticket/delegate ceremony exceeds it for bounded work | Wide, policy-heavy, multi-owner work may need state/coordination Agentless does not target | Which Bench layers beat the simple baseline after full cost accounting? |
| [Lost in the Middle](https://arxiv.org/abs/2307.03172) | `[RESEARCH]` Long-context use varies with information position; presence does not ensure effective use | Compact query outputs and progressive references help | Roughly 6k cold words and 31k guidance estate rely on availability | Retrieval QA findings do not determine exact coding-agent context limits | What size/order maximizes Bench task outcomes by model? |
| [Coconut](https://arxiv.org/abs/2412.06769) | `[RESEARCH]` Continuous latent reasoning can carry useful intermediate state without verbalizing every step | Reinforces not persisting full model narration | Current handoff/capture favor textual narration | Training-time latent architecture cannot be implemented by a repository harness | No Bench architecture decision depends on Coconut; keep it only as a conceptual caution |
| [ToM-SWE](https://arxiv.org/abs/2510.21903) | `[RESEARCH]` Persistent user intent and interaction state can aid long-horizon software work | Work IDs, goal, decisions, and resume state align | Current handoff State is stale/unvalidated; intent ledger is assignment-only | A separate user-model agent adds complexity and privacy/state risks not justified here | Does explicit user-intent state reduce clarification and wrong-scope edits? |
| [SWE-Effi](https://arxiv.org/abs/2509.09853) | `[RESEARCH]` Evaluate correctness with tokens, time, expensive failures, and scaffold-model fit | AXI outputs and gate reuse reduce some cost | Guidance/context and fixed orchestration can create token snowball and expensive failure | Efficiency rankings change by model and infrastructure | Does Bench's reliability benefit remain after total resource cost? |
| [SWE-Skills-Bench](https://arxiv.org/abs/2603.15401) | `[RESEARCH]` In its study, 39/49 skills did not improve pass rate; average improvement was about 1.2%, some added up to 451% tokens at unchanged pass rate, and three degraded results through mismatch | Bench has conformance/ownership for skill text | No Bench skill has a paired marginal-outcome benchmark; 233 canaries prove anchors, not behavior | Its skills/tasks/models are not identical to Bench | Which Bench skills improve their named target after overhead and version mismatch? |
| [Matt Pocock skills](https://github.com/mattpocock/skills/tree/9c9f36ccd3995266cd675468af71639c8dde1ec5) | `[RESEARCH + OBSERVED]` Routes engineering work through focused practices with explicit context boundaries; debugging requires reproduction and discriminators | Bench preserves most causal practices and adds deterministic state/fences | Composition dilutes debug/review and lacks upstream's clear router | Practitioner design is not randomized causal evidence | Which integrated changes improve outcomes over the pinned upstream alone? |
| [no-mistakes](https://github.com/kunchenguid/no-mistakes/tree/6859d1e827f5ab2592a4703d3bab8734a38c9aa5) | `[RESEARCH + OBSERVED]` Uses disposable worktrees, isolated review stages, structured findings, safe-mechanical boundaries, validation, PR and CI | Bench worktrees/gate/fences are deep realizations | Review isolation is prompt-only and release workflow bypasses governed publication | Different host/CI/product assumptions; current source is not Bench's historical import | Can Bench retain isolation with less orchestration and exact historical provenance? |
| [AXI](https://github.com/kunchenguid/axi/tree/408a6536625e5b05e5c56e6c4a04fe83e1f510a5) | `[RESEARCH + OBSERVED]` Agent CLIs should be content-first, minimal, bounded, structured, ambient, and consistent | Most read queries align | Root help is large/inconsistent; prompt phases are undiscoverable capabilities; test side effect violates expectation | AXI principles do not choose Bench's workflow or semantic authority | Does the routed root measurably reduce discovery time and tokens? |

The three “J” concepts remain separate: neural global-workspace evidence is not J-CoT's latent mechanism, and neither is a text-level Bench control protocol. The transferable engineering hypothesis is selective, fresh, shared state; its concrete schema still needs the benchmark in W.

## Y. Ten hardest truths

1. Bench's strongest results come from ordinary deterministic software—identity, ownership, checks, and state transitions—not from its largest prompts.
2. A 233-fixture conformance suite can prove every sentence is pinned while proving nothing about whether a model follows it.
3. There is no current normal-user entry point; a catalog, status board, handoff heuristic, and roadmap ritual do not add up to a router.
4. The handoff can be semantically wrong at the exact commit that declares it fresh.
5. Review and Final Check are advisory prompt performances, not controls, and their apparent formality obscures that boundary.
6. The green development gate is trustworthy for its declared oracle but cannot establish semantic, native, ship, or publication completion.
7. Mandatory agents and fresh contexts do not create independence unless inputs, identities, and information boundaries are observable.
8. Most skill value is unmeasured, and current research makes “the prose sounds good” an indefensible retention test.
9. A rewrite would likely destroy more demonstrated value than it creates; indefinite incrementalism would preserve the incoherent control plane.
10. The next unit of work must be one state-and-routing tracer with a benchmark, not another phase, skill, dashboard, or ontology.

## Z. One next ticket

### CR-001 — Route one active work record through one canonical `bench` entry

**Outcome.** Implement a vertical tracer in which `bench` (no subcommand, with optional intent) reads deterministic repository facts plus one minimal canonical work record, selects one route, explains why, compiles the smallest action context, and renders equivalent Claude and Codex continuations. Existing gate, status, preflight, coverage, worktree, and release commands remain unchanged behind it.

**Acceptance contract.** `[INFERRED]`

- Define one versioned record with goal, subject/base/tip, scope, decision/requirement pointers, evidence IDs and freshness, open uncertainty, attempts/diagnostics, original feedback signal, and one next action.
- Read current intent/spec/handoff/status sources for backward compatibility; do not treat handoff prose as authoritative. Emit handoff/dashboard as generated views.
- Route fixture states for: new idea, ready ticket, failing bug, interrupted attempt, dirty unverified change, and locally complete work. Expert subcommands bypass routing without behavior change.
- Emit one minimal structured result: state, route, reason, prerequisites, next command, selected source pointers, and stale/contradictory inputs.
- Generate `/bench` and `$bench` adapters from the same registry and prove semantic equality apart from prefix syntax.
- Preserve exact gate evidence, ownership fences, requirement-loss detection, and governed release semantics; do not reimplement them in the router.

**Concrete proof plan.** Start with failing contract fixtures for all six routes, same-commit stale handoff rejection, wrapper/binary root parity, Claude/Codex rendering, blank-retry inheritance, and definitive empty state. Demonstrate that current root fails those fixtures. Implement the tracer, then run focused owners, all 233 canary owners, the full project gate, and the A/B/C pilot on at least 12 representative tasks. Accept only if routing is correct in every deterministic fixture, no stale handoff field authorizes a transition, existing enforcement remains green, active guidance tokens fall at least 30%, and pilot success is non-inferior under the preregistered W rule. If the outcome threshold fails, retain the new state reader as experimental and do not migrate another phase.

**Out of scope.** No Review runner, skill rewrite, new dashboard, broad schema migration, release rewrite, or deletion of legacy artifacts belongs in this ticket.

The 45 required adversarial questions are answered individually in [questions.md](questions.md), with the same evidence labels and no additional implementation recommendation.
