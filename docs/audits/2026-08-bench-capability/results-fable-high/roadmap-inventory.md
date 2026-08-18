# Roadmap inventory — what exists, verified against `58d966e2`

## Roadmap-like artifacts discovered

| Artifact | What it is | Count / state at HEAD | Owner command |
|---|---|---|---|
| `ROADMAP.md` | index: one bold line per row, sections, dependency tables, recommended sequence | 73 rows; sequence = FT100, FT207, FT213; release status **NO-GO** | `/bench-what-next` (drain), `bench roadmap` (read) |
| `roadmap/FT<n>.md` | row bodies (problem, occurrences, sources) | 73 files, 12,626 words; heaviest FT162 (708), FT89 (662), FT164 (630) | same |
| `specs/*/` | staged specs + tickets | **0 `spec.md`**, 34 dirs holding only `tickets/` (43 files); 28 are `light-path-*` — the standing light-path route writes one ticket and has no close/retire step (e.g. `light-path-duplicate-acceptance-ids` landed 2026-08-03 as `b6b77897`, ticket dir still present) | `bench spec …`, `bench coverage`, `bench preflight` (all report "spec not found" for these) |
| `decisions/*.md` + `bench maps` | decision maps / tickets | 7 unresolved maps: 5 fog "Not yet specified" (bounded-network-resource-cli, diff-visual, gate-concurrency, gate-critical-path, gate-pipeline), 1 frontier task (cost-follows-project-size), 1 **invalid** (spec-build-review-gate-cadence → deleted `internal/specbuild/*`) | `/bench-shape-idea`, `bench maps` |
| `capture/IDEAS.md`, `capture/learnings.md`, `capture/retros/` | inbox, journal, retros | all drained to 0 at HEAD (`bench roadmap` drain row 0/0/0) | `bench idea`, `/bench-what-next` |
| `capture/session-handoff.md` | resume artifact | pin block regenerated; `## State` prose stale at HEAD (see ledger L-07) | `bench handoff` |
| `ASSESSMENT.md` (2026-08-13), `skills-assessment.md`, `COMPLIANCE_ASSESSMENT.md`, `capture/architecture-review-20260817T104714.html` | prior self-assessments (~80 KB) | ASSESSMENT.md §6 already records the raw-`npm publish` bypass; none records the conformance-in-gate loss | `/bench-assess` |
| `docs/release-runbook.md`, `.github/workflows/release.yml` | release procedure vs implementation | runbook requires `bench release submit/promote`; workflow uses raw `npm publish` (ledger L-13) | — |
| `.bench/structure-accept` | structure suppression list | 4 stale accept rows (deleted paths) | `bench structure` |
| `roadmap/FT133.md`, `FT120.md`, `FT213.md` occurrences | in-row records of the `TestRootConformance` skip trap | 4 occurrences 2026-07-26 → 2026-08-16; none names that the dev phase was removed 2026-08-09 | — |

Deferred review findings live inside rows (FT140, FT142, FT208); known-bug lists are the
"False greens" and "Reds the diff doesn't own" sections; unimplemented acceptance criteria
survive only as tickets-only folders (above); proposed commands: `bench branches retire`
(FT199), `bench capture commit` (FT166), `bench transcript` (FT204), `bench spec show`/
`--symbol` (FT125), scoped gate (FT215).

## Findings from the two audits that have no roadmap row today

| Finding | Ledger | Row? |
|---|---|---|
| Live-root conformance not graded in the dev gate; 10 red | L-01 | none (only occurrences of the skip in FT133/FT120/FT213) |
| Red verdicts never drift-checked; status/handoff contradict gate | L-02 | none |
| No router / staged-spec signal / un-adopted status | L-03/15/16 | none (`what-next` rename not on board) |
| Adoption gate cannot go green (gate-inputs.json, HOME, canary) | L-14 | none |
| Test-isolation leak into `~/.bench/worktrees` (759 orphans) | L-06 | none |
| `/bench-debug` compressions vs upstream; Codex trigger | L-11/L-10 | FT89 and FT112 touch `/bench-debug` prose; nothing on the compressions |
| Release workflow raw `npm publish` | L-13 | none as a row (release-readiness item 4 only) |
| Guard degraded-rim substring; deny gaps | L-20 | none |
| `bench outline` bare 24.6 KB | L-19 | FT125/FT191 adjacent, not this |
| Codex adapters carry inert Claude key; no Claude parity check | L-30 | none |

## Row-by-row verification

Legend — *exists?* = the problem the row states is still observable at HEAD (fresh check
noted where one was cheap); *partial* = part of the row's work has landed; *occ* =
`Occurrence:` lines. Row summaries were gathered by three read-only delegates and
spot-verified; anything marked ✔ was re-checked directly in this run.

| FT | Title (short) | Subsystem | Kind | Exists? | Partial | Occ | Overlaps / notes |
|---|---|---|---|---|---|---|---|
| FT6 | parked bundle (~7 unrelated parks) | other | parked | n/a | no | 0 | per-item graduation triggers; not one row |
| FT8 | Sonnet 5 mid-tier revisit | delegation | parked | time box 2026-09-01 not reached | no | 0 | binding lives in `lines.env`/profile |
| FT24 | Codex agent-line guard parity | hooks | parked | yes — `.codex/hooks.json` has no Agent guard ✔ | no | 0 | duplicates re-check note in BENCH-reference |
| FT38 | dashboard visual identity | CLI | decision | yes | v1 shipped | 0 | reviewer taste decision |
| FT58 | hardened pool roots | worktree | bug | yes (`_ = chmodPool`) | no | 0 | part decision (propagate vs best-effort) |
| FT71 | versioned local shift evidence (bank) | gate | feature | yes | no | 1 | bank track; after FT169 |
| FT89 | guidance coherence / current-state docs | prose | guidance | partly (some occurrences already fixed) | yes | 8 | includes a real build: generate root help/inventory from `commandRegistry` (→ router/root work); overlaps FT100/FT106/FT113 |
| FT92 | attributed subject drift + consumer input hygiene | gate | bug | yes (bare drift message) | no | 0 | two unrelated halves |
| FT94 | single-sourced resume summary golden | CLI | bug | yes | ? | 1 | standards-debt batch |
| FT98 | preserve-then-discard primitive, four faces | worktree | feature | yes | yes (3 commits shipped) | 9 | accreted; duplicate 2026-07-30 occurrence |
| FT99 | spec problem-premise verification | spec | guidance | partly | ? | 8 | overlaps FT106 |
| FT100 | prose-weight pass | prose | guidance | yes | yes | 1 | rank 1 on board; row itself says anchor consolidation is prerequisite; sequenced after FT89 and FT170's benchmark decision |
| FT101 | per-context scope for monorepos | other | feature | yes | no | 0 | contains a gate-scoping sub-decision |
| FT102 | escalation-policy cross-check in synthesis | prose | guidance | partly | ? | 1 | 6 lines |
| FT103 | existence-checked absence evidence (gate half) | gate | bug | yes (per-row `Source` never stat'd) | yes | 1 | narrow |
| FT104 | load-induced commit refusals | gate | guidance | partly | ? | 3 | reds-not-owned section |
| FT106 | doc claims re-verified against tree | capture | feature | partly | no | 0 | FT172 wants its probe |
| FT108 | refactor lane with mechanical exit test | prose | feature | yes (no craft-refactor) | no | 0 | after FT89 |
| FT111 | provenance tags outlive specs | prose | guidance | yes (live `FT76 story 3` tags) | no | 1 | land with FT179 |
| FT112 | approximation-stays-green is not a cleared bug | prose | guidance | partly (`/bench-debug` prose) | ? | 1 | same file as FT89's debug items |
| FT113 | `bench commit --spec` residuals | CLI | bug | yes (help wording; no misuse guard) | yes | 4 | 2026-08-17 occurrence |
| FT115 | load-robust deadlines from bounds | gate | bug | partly | ? | 2 | near-duplicate FIFO half with FT120 |
| FT117 | FT87 parser-surface follow-ups | CLI | bug | partly (one named leaf already routed) | yes | 1 | body stale in ≥1 leaf |
| FT120 | gate/canary/contract harness defects | gate | bug | partly | ? | 7 | includes 2026-08-14 `TestRootConformance` skip occurrence (→ L-01) |
| FT125 | reader surfaces return the slice | CLI | feature | yes | no | 1 | needs its own validation |
| FT130 | capture write mid-lifecycle | gate | bug | partly (detection shipped) | partly | 2 | policy choice open |
| FT133 | coverage --check verifies red-signal citations | gate | feature | **premise stale** — `red signal` column retired from schema (`coverage_test.go:28-33`) | partly | 8 | 2 of 8 occurrences are the `TestRootConformance` skip (→ L-01) |
| FT138 | instrument build economics | other | decision | partly (JSONL `ElapsedMS` exists) | partly | 2 | prerequisite for any measured prose cut |
| FT140 | review residuals wanting a verdict | capture | decision | prose | ? | 5 | five unrelated verdicts |
| FT141 | `gate pin` records red verdicts | gate | feature | yes | no | 5 | **name collision**: `bench gate pin` already exists as the pre-push pin porcelain |
| FT142 | FT91 review residuals | release | decision | partly | ? | 2 | citations outside tree |
| FT144 | kit specs have two audiences | prose | guidance | yes | no | 5 | three asks accreted |
| FT158 | cross-harness falsification standing | delegation | guidance | yes | partly | 4 | absorbed six disciplines |
| FT162 | full-run/phase-close single subject & handoff | CLI | feature | partly (handoff sections shipped) | partly | 11 | most accreted row; overlaps handoff extension (L-04) |
| FT164 | repair-lane charges; done-claim resolves owners | delegation | guidance | partly | yes (core shipped 2026-08-03) | 16 | |
| FT165 | fold domain-modeling into shape-idea | prose | guidance | **no** — `bench-craft-domain` exists and `bench-shape-idea.md:14` charges it ✔ | done | 1 | satisfied |
| FT166 | `bench capture commit` porcelain | CLI | feature | yes | no | 2 | two deliverables |
| FT168 | fixture-selecting canary | gate | feature | yes | partly | 2 | narrow |
| FT169 | worktree landing command owns stale-base dance | worktree | feature | shipped; residual refusal-message defects | yes | 11 | rewrite to residuals |
| FT170 | behavioral red/green evaluation for skill guidance | prose | decision | yes | no | 0 | parked, yet gates FT100 (ROADMAP:273) — this *is* the measurement harness |
| FT172 | roadmap parser/context snapshot completeness | capture | feature | partly (snapshot shipped) | yes | 6 | after FT106 |
| FT173 | AXI residual: active assignment, deleted tree | worktree | decision | yes (`list.go:98-113`) | no | 1 | smallest well-defined item |
| FT174 | ticket grammar parser | spec | feature | yes (no `ParseTicket` in landed tree) | partly | 12 | two tickets-only folders carve pieces off it |
| FT177 | stale `dist/bench` → silent no-op probes; landing guard deletes binary | gate | bug | yes (no mod-time check); "documented nowhere" occurrence stale | partly | 5 | live foot-gun half relates to cold-start (L-29) |
| FT178 | bare `bench worktree` traps automation | worktree | bug | yes (`bin/bench.sh:290,367`) | no | 1 | do not probe by running |
| FT179 | comment quality / reviewer register | prose | guidance | partly | ? | 2 | with FT111; cites non-existent `ParseTicket` |
| FT180 | spec-optional route at shape-idea exit | spec | decision | partly (map homing shipped) | partly | 2 | the light-path residue (34 dirs) is this route's missing close step |
| FT182 | Planned receipt over absent target wedges abandon | worktree | bug | yes (`resume.go:171` branch) | no | 1 | precisely scoped |
| FT185 | gate results join TOON contract | CLI | feature | yes (`fmt.Fprintln` in runner.go) | no | 1 | spec-ready |
| FT190 | injected interfaces need real-producer test or exemption | gate | feature | **no** — `injected_ports_test.go` enforces exactly this ✔ | done | 1 | satisfied |
| FT191 | fixture-and-seam inventory for charges | CLI | feature | yes | no | 2 | two asks |
| FT192 | one source per fact in spec/ticket prose | prose | guidance | partly (parser retired) | yes | 3 | craft-spec visit with FT209/FT214 |
| FT197 | Go core owns gate invocation/process lifetime | gate | feature | partly (process groups shipped; shell hop remains) | yes | 0 | body stale |
| FT199 | branch-retirement coordinator | worktree | feature | yes (no `branches` verb) | no | 6 | bundles 6 landing-refusal defects; overlaps FT206 |
| FT200 | preflight mechanical at landing chokepoint | gate | decision | yes (preflight not in gate.sh/prepush) | yes (command shipped) | 0 | Sol's "automatic preflight" ask |
| FT201 | cancel-signal registrations one source | gate | feature | yes (`skillsindex/command.go:52` drifts) | partly | 1 | live drift instance |
| FT202 | standing test-support fence; census scope | spec | decision | yes | no | 2 | |
| FT204 | bounded transcript/session query | CLI | decision | yes | no | 0 | new surface |
| FT205 | craft-delegate end-of-life pair | prose | guidance | yes | partly | 5 | same file as FT213/FT221 |
| FT206 | exact-candidate review sees destination metadata | worktree | feature | yes (`landing.go:385`) | no | 2 | duplicate surface with FT199 |
| FT207 | worktree-mutating paths share malformed-admin refusal | worktree | decision | partly | no | 1 | rank 2 on board; unblocked by FT189 |
| FT208 | skills-index hardening residuals (3 decisions) | CLI | decision | yes | no | 1 | three independent decisions |
| FT209 | refactor proves itself by differential; grouping cardinality | prose | guidance | yes | no | 1 | craft-spec/craft-domain edits |
| FT212 | `worktree clean --landed` invalid invocation | worktree | bug | **no** — `--landed` accepted, rc 0, usage lists it ✔ | done (FT216) | 1 | stale row |
| FT213 | read-only delegate gets own worktree; oracle-verified probe | delegation | guidance | yes (`SKILL.md:99`) | no | 2 | rank 3; one code guard inside |
| FT214 | build may not edit its own spec rows/fences | prose | guidance | yes | no | 1 | craft-spec visit |
| FT215 | changed-package-scoped gate | gate | feature | yes | no | 1 | perf; measured timings in row |
| FT217 | one adopt-lifecycle decision | adopt | feature | yes (only `upgrade` uses `PlanLifecycle`) | partly | 1 | with FT218 |
| FT218 | named git readers | git | feature | yes (6 sites) | no | 1 | with FT217 |
| FT219 | deepen refreshes ready map frontier | capture | guidance | yes | no | 1 | single insertion point |
| FT220 | write-spec censuses shared readers | spec | guidance | yes | no | 1 | single insertion point |
| FT221 | craft-delegate cp-aside checklist step | prose | guidance | yes (`SKILL.md:94`) | no | 1 | same file as FT205/FT213 |
| FT222 | per-repair-class delegate-tier preference | delegation | decision | yes | no | 1 | profile table addition |
