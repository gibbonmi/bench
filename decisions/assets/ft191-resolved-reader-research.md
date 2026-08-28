# FT191 research: a resolved read surface for Bench

Consumed by: the FT191 `/bench-write-spec` run, or an FT191 shaping map if the reviewer routes one first.
Drift: re-verify after a change to `internal/outline/outline.go`, the FT173 AXI contract, `bench-craft-seams`, `bench-craft-review`, or `go.mod`.
Retire when: FT191's spec stages, or the reviewer retires FT191.

External citations carry the retrieval date 2026-08-28. Re-fetch a cited page
before any engine decision contested after 2026-Q4.

## Question and scope

FT191 (`roadmap/FT191.md`) asks for a reader that emits helpers, doubles, and
fixtures with `file:line`. It also asks for "every executable consumer from
registries and the call graph". The open decision is projection versus
existing-outline mode. This asset answers the seven charge decisions: engine,
surface shape, consumers, enforcement, determinism, measurement, and kill
criteria. It also tests the "make code research deterministic" framing. Prior
art is Benzi; every Benzi number below is vendor-run.

## Prior art: what Benzi shows

All Benzi facts pin to `github.com/oooscoos/Benzi` at commit `78781f91`
(retrieved 2026-08-28). The engine is closed and proprietary, and the license
forbids reverse engineering. Every architecture claim therefore rests on
vendor prose. Its `swebench/metadata.yaml` states `oss: false` and
`verified: false`. Every number is vendor-run; the vendor's own
`benchmark/README.md` says so itself.

Four findings carry weight for FT191:

1. **The used surface is small.** The vendor advertises "35+ tools". Its 500
   published SWE-bench trajectories use the deep resolution tools zero times
   (`backflow`, `forwardflow`, `trace_path`, `call_tree`, `external_calls`).
   The index tools that carry load are simple: `get_definition` (344 calls),
   `search_symbols` (240), and `get_callers` (127). Meanwhile `shell` (7,822),
   `read_source` (4,641), and `benzi_grep` (3,127) dominate. Lesson: build a
   small resolver surface — definition, callers, references, symbol search —
   and nothing speculative.
2. **The symptom map is deterministic, with published weights.** Traceback
   frames weigh 3, quoted code weighs 2, and named symbols weigh 1. Sites
   rank by convergence (`swebench/SWE_BENCH_REPORT.md` §2). The model still
   decides everything. The layer is benchmark configuration, not the shipped
   product. The weights are a usable precedent for `bench-debug` orientation.
3. **The headline numbers do not replay.** The cross-harness lines-read
   figure (9,125 versus 20,704) cannot be recomputed from the published
   artifacts. The control-arm transcripts are gitignored by design. A replay
   of the vendor's own `lines_read.py` over the published `runs.jsonl` gives
   10,005, not 9,125. Two live vendor pages disagree on the SWE-bench cost
   ($37.33 versus $28.92). Treat every Benzi figure as directional only.
4. **A model layer confounds the SWE-bench score.** A second-pass verifier
   requested revision on 494 of 500 instances. Every instance is therefore a
   two-draft run. The 78.2% resolve rate is not attributable to the index
   alone.

Benzi publishes no claim that its index is byte-reproducible at a commit. It
states honest depth limits: runtime tracing is Python-only.

## §1 Engine, under the build contract

The build contract is fixed. `scripts/go-build.sh` pins `CGO_ENABLED=0`,
`-trimpath`, and the toolchain. `go.mod` holds one dependency, and the
dependency standard prices every new one.

### Tree-sitter: rejected

- Both official-lineage Go bindings require cgo. `tree-sitter/go-tree-sitter`
  opens with a cgo preamble; `smacker/go-tree-sitter` vendors the C sources
  and is dormant since 2024-08-27.
- The cgo-free paths are immature. `malivvan/tree-sitter` (wazero WASM route)
  is self-described pre-release, with 5 stars. `odvcencio/gotreesitter`
  (pure-Go runtime, 206 grammars, 2 modules) is six months old. It runs
  ~4.8× slower than the C runtime at full parse, and it has withdrawn one
  benchmark headline.
- The deeper problem is categorical: tree-sitter parses, and it does not
  resolve names. The resolution layer built on it, `github/stack-graphs`, was
  archived 2025-09-09 with four languages covered. GitHub unshipped precise
  code navigation, and the paper's author says the project "withered on the
  vine". A tree-sitter choice means per-language name resolution written from
  scratch, on an unproven runtime. That is Benzi's closed engine rebuilt by
  hand.

### Headless LSP: rejected as the reader, kept as prior art

- All six major servers run headless over stdio, and the wrapper pattern is
  proven (`microsoft/multilspy`; `oraios/serena`, 28k stars).
- The protocol is nondeterministic by design. LSP 3.17 defines
  `ContentModified` for answers invalidated mid-flight. Answers depend on
  background-index progress. Pyright returns references only for files sent
  through `didOpen` (`microsoft/pyright#10086`) — silently wrong, not an
  error.
- The cost is GB-class memory per language server. A citation cannot replay
  against a moving index, so a reader a hook re-runs cannot sit on this
  substrate.

### SCIP: viable as a format, deferred as an engine

- SCIP is alive and openly governed since 2026 (`github.com/scip-code/scip`,
  Apache-2.0, protobuf). The Go reader is a separate module with four runtime
  dependencies, protobuf included, and no cgo.
- The per-language indexers are uneven. `scip-go` and `scip-java` are active.
  `scip-typescript` and `scip-python` have had no default-branch commit for
  ~11 months, and `scip-clang` is beta. LSIF is dead: lsif.dev is archived,
  and the `scip` CLI dropped LSIF conversion entirely.
- A SCIP index is not byte-reproducible as produced. The format embeds an
  absolute `project_root` and the indexer argv, and it guarantees no field
  order. Semantic stability needs a pinned indexer version plus a
  consumer-side pass (`CanonicalizeDocument`, `SortDocuments`).
- `scip-go` is itself built on `golang.org/x/tools` and shells out to `go`.
  For Go, SCIP serializes the same engine the recommendation below uses
  directly.

### Recommendation: Go-native `golang.org/x/tools`, in process

| Fact | Evidence |
|---|---|
| Dependency delta: 3 modules (`golang.org/x/tools`, plus `x/mod` and `x/sync` indirect) | probe `go mod tidy`; re-run by the coordinator on this machine, 2026-08-28 |
| `CGO_ENABLED=0 -trimpath` build passes (3.75 MB probe binary) | same probe, same machine |
| Load, type-check, SSA, and CHA over this repo: ~1.5–2.7 s (~434k parsed lines with deps) | delegate measurement on this machine, warm and cold `GOCACHE` |
| Node and edge counts stable across runs (21,514 / 360,553) | delegate measurement, 4 runs |
| Raw ordering nondeterministic; a ~170-line naming tail varies through type aliases | delegate measurement; see §5 |
| Runtime needs a `go` binary (`go/packages` shells out to `go list`) | `x/tools/internal/gocommand/invoke.go`; failure reproduced with `go` off PATH |

The license is BSD-3 (MIT-compatible), the org is official, and the footprint
is build-time-only. That is the dependency standard's exact shape. The query
mapping:

- References: enumerate `go/types.Info.Uses` across the loaded packages.
  This yields exact edges for static identifier uses. The gopls enumerations
  are internal and not importable, so the walk is written here; it is modest.
- Implements: `types.Implements`, pairwise over the package scopes.
- Call graph: `x/tools/go/callgraph` CHA for a library tree — sound, and
  over-approximate on interface calls. RTA needs a `main`, and VTA is marked
  experimental; defer both.
- No algorithm sees reflection-mediated calls. The output says so.

At ~2 s per full pass on this repo, the resolver regenerates per call, like
outline. There is no persisted index, no staleness state, and no freshness
contract. Determinism then reduces to canonical ordering plus stable symbol
identity (§5).

**What a linked non-Go repo gets: nothing resolved.** It keeps
`bench outline`, and it gains the citation convention (§3), which is
engine-free. The SCIP reader stays the named additive path if a linked repo
ever produces indexes. Do not build that consumer now: one adapter is
hypothetical (`bench-craft-seams`).

## §2 Projection or sibling (FT191's open question)

**Recommendation: a sibling reader, not an outline mode.** Split FT191's two
row families at spec time:

- The inventory family (helpers, doubles, prior-art fixtures) is
  candidate-class. It can extend `bench outline` through the existing pattern
  table (`internal/outline/outline.go:71`: a new form "is adding a table
  entry"). The rows keep outline's promise, and the agent confirms each
  candidate by reading the line.
- The consumer family (callers, references, registry consumers) is
  verified-class. It needs the resolver sibling. Its value is exactly that
  the agent does not confirm each row by hand.

Four facts force the split:

1. **The trust levels differ.** Outline "LOCATES candidate seams; it does not
   IDENTIFY the project's blessed seams" (`internal/outline/outline.go:14`).
   A resolved callers row is a verified edge. One table with two trust levels
   either downgrades verified rows or silently upgrades regex rows. Both
   readings break the promise.
2. **The machinery is disjoint.** Outline shares only `internal/toon`,
   `internal/git`, and `internal/usage` with a resolver. A `--callers` flag
   on outline merges two fact sources under one grammar. Every outline change
   would then carry a resolver question. The one-source standard cuts against
   that coupling.
3. **The empty state lies for unsupported languages.** Resolution starts
   Go-only. An empty callers table for a `.ts` symbol reads as "no callers",
   a universal claim; the truth is "cannot resolve this language". A sibling
   can refuse, or mark the language unsupported, per row. `craft-review`
   pins the stake: "a universal claim without an enumeration is a sample".
4. **The seams-skill tests agree.** Deletion test: without the resolver, each
   consumer hand-rolls its own `rg` pipeline, so complexity reappears. The
   one-adapter rule: one engine exists today, so build the Go resolver
   directly behind the command. Do not build a multi-engine interface before
   a second engine is real.

**LOCATE/IDENTIFY verdict.** A callers query has crossed from locating into
identifying: it identifies reference edges. It still does not identify
blessed seams; `projects/<name>.md` keeps that. The promise gains a third
clause rather than an exception. Outline locates candidates; the resolver
identifies edges for supported languages; the project profile identifies
blessed seams. Restate the resolver's promise in its help text, the way
`internal/outline/outline.go:49` does.

The sibling is a new AXI query surface: TOON tables, definitive empty states,
structured stdout errors, exit 0/1/2, and bounded output with `--full`. The
FT173 contract leaves AXI widening as a reviewer decision
(`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`, Out of
scope). FT191's spec is where that decision lands.

## §3 Consumers by phase

| Consumer | What it gets | Row shape | Advisory or enforceable |
|---|---|---|---|
| `craft-seams` reroute (write-spec, implement) | callers and references instead of declarations; the read budget buys confirmation reads, not discovery reads | `consumers[N]{symbol,file,line,via}`, `via` ∈ call, reference, implements, registry | advisory; the seam choice stays model judgment |
| `craft-research` citations | a replayable citation for every enumeration-class claim | citation fields: repo SHA, exact argv, bench version, output hash | replay enforceable; claim relevance is not (§4) |
| `bench-review-implementation` | the blast list: touched symbols' callers and holders over the frozen base..tip | `blast[N]{changed_symbol,consumer_file,consumer_line,via}` | presence and freshness enforceable; consideration is not |
| `bench-debug` orientation | quoted identifiers and traceback frames resolved to sites before Phase 1 | `sites[N]{identifier,file,line,kind,ambiguity}` | advisory; never ranks hypotheses |
| `bench structure` / stop hook | nothing new, initially | — | no new gate check (§4) |

The review consumer is the strongest. `craft-review` already demands the
enumeration: "The citation is the grep over all of it, or the per-member run"
(`.agents/skills/bench-craft-review/SKILL.md`). The resolver is an
enumeration engine whose enumerations replay. FT191's own ledger records the
motivating miss: in FT210, a planner-local counter stayed green while two
commands reached target Git through older consumers (`roadmap/FT191.md`). A
blast list over the diff surfaces exactly that class.

The debug consumer needs a false-positive rule. `bench-debug` forbids theory
before the repro loop exists (`.agents/commands/bench-debug.md`, Phase 1). A
ranked symptom map on a vague report anchors the model before any red signal
exists. So the orientation output serves loop construction only: which
surface to invoke. It emits exact-match sites with an explicit `ambiguity`
count, and it refuses to rank a name that resolves to many sites. Benzi's
weight scheme is the precedent; surface the weights as columns, never as one
merged rank.

## §4 Enforcement versus shaping

The test: a use is enforcement only when a hook or the gate can pass or fail
it without any reading of model output.

| Proposed use | Class | Why |
|---|---|---|
| citation replay: recompute the cited invocation at the recorded SHA; compare the output hash | enforcement, of citation integrity | mechanical given SHA, argv, and bench version; it proves the evidence is real, never that the claim follows |
| "a callers query ran before the write" | looks like enforcement; is "the model was told to" | binding the query argument to the diff's touched symbols stops a trivial query; nothing stops an ignored query — do not build it |
| the review blast list attached to the review artifact | enforcement of presence and freshness only | a phase hook recomputes and byte-compares; consideration stays judgment; keep it at the review hook, never in the gate |
| symptom map before `bench-debug` Phase 1 | shaping | orientation input only |
| seams reroute returning consumers | shaping | same |
| a gate red for "touched symbol with unenumerated consumers" | rejected for the gate | it grades an advisory artifact, not the tree; pasting the list unread satisfies it; the ceiling is a `bench status` advisory count, the `bench structure` pattern |

The gameable case has a direct answer. Requiring a callers query before a
write is satisfiable by a trivial query, and a scoped variant is satisfiable
by an ignored one. So the discipline does not gate the write. It puts the
enumeration where ignoring is visible to a person: the review blast list.

## §5 Determinism table

| Component | Class | Note |
|---|---|---|
| resolver graph build at a SHA | deterministic, with discipline | pinned toolchain (`go.mod`), declared build tags, canonical sort |
| symbol naming for generic instantiations | measured hazard | a ~170-line edge-name tail varied on 2 of 4 runs (`os.DirEntry` versus `io/fs.DirEntry` alias spelling); canonicalize identity through the type-checker's object, not the printed name |
| query output ordering | deterministic | sort rows before the TOON emitter; outline models this |
| symptom-map site extraction | deterministic given the extractor | the regex extractor is fixed; which report text is fed stays user input |
| symptom-map ranking | rejected | §3 replaces ranking with ambiguity counts |
| citation replay | deterministic | given SHA, argv, and bench version; record all three plus the output hash |
| review blast enumeration | deterministic | base..tip are frozen by the review contract (`projects/benchkit.md`) |
| index build (SCIP path) | deterministic only after normalization | absolute root and argv embedded; no field-order guarantee |
| claim selection, synthesis, seam choice, hypothesis ranking | stochastic (model judgment) | never a gate input |
| Benzi symptom map and scope cards | vendor-claimed deterministic; unverifiable | engine closed |

The alias tail is the one open determinism hazard. Whether canonical symbol
identity closes it fully is a runnable-probe question, not a reading
question. The FT191 spec names that probe: N resolver runs over this repo,
byte-equal after canonicalization, as a red-capable check.

## §6 Measurement — the smallest isolating experiment

Bench has no benchmark. The kit itself is the ~100k-class subject: 126,407 Go
lines across 1,784 tracked files (counted 2026-08-28). The design is a paired
A/B on identical tasks: same harness, same model and tier, same charges
except the read surface. Arm A gets `bench outline` as-is; arm B adds the
resolver, named in the seams reroute. Hide `.git` during runs; Benzi's
`_hidden_git_history` exists because a control run read the fix commit out of
history.

Two task families, both with tree-derived ground truth:

1. **Seam location** (~15 tasks): re-derive consumer lists for symbols whose
   true consumer sets are known from landed specs. Metrics: files read under
   the declared seam budget, reroute frequency, wrong-seam rate.
2. **Planted-regression review** (~25 plants): take landed commits and plant
   regressions in callers outside the diff's context window — the FT210
   class. Metric: review catch rate per arm.

Noise floor: below 20 planted regressions per arm, a catch-rate delta under
~30 points is noise (paired McNemar, 80% power). Below 10 paired seam tasks,
the files-read comparison detects only very large effects. Repetition (3 runs
per task) controls model nondeterminism; the harness exposes no temperature.
Expect a modest exploration effect: Benzi's trajectories ran index queries at
roughly one-fifteenth the volume of shell, read, and grep. The review catch
rate, not files read, is where the resolver should show first.

Adapt Benzi's harness design, not its code. The code is Windows-only, holds
its task list as code, and gives each CLI a bespoke arm. The portable parts
are the grading shape (baseline, then grade, fail-to-pass plus pass-to-pass),
the `.git` hiding, and the checkout-fix-then-revert-sources worktree recipe.
Two recorded traps: an asymmetric turn cap on one arm voids the comparison,
and per-task variance of 14–16% makes single runs meaningless.

## §7 Kill criteria, ranked by likelihood

1. **Non-Go linked repos gain nothing beyond outline — most likely, and
   partially true by construction.** The resolver is Go-only, so the kit is
   the first and main beneficiary. Two softeners: the citation convention is
   language-free, and the SCIP reader stays a later additive path. Reviewer
   input needed: the language mix of the first external validation target.
2. **The host harness already covers it — bites the navigation framing,
   spares the evidence framing.** Claude Code ships a native LSP tool today:
   definition, references, implementations, call hierarchy, eleven official
   server plugins (`code.claude.com/docs/en/tools-reference`). Codex CLI
   ships nothing native (`openai/codex#8633` open, unanswered), and OpenCode
   gates its tool behind an experimental flag. No harness gives determinism
   at a SHA, replayable citations, or a frozen blast list. Aim the reader at
   evidence, not navigation.
3. **Citation replay is too heavy for `craft-research` — medium; design
   answers it.** The reader emits its own citation line (SHA, argv, version,
   output hash), so adoption is copy-paste. Replay runs on demand, not per
   claim.
4. **Symptom-map false positives anchor the model — low-medium; posture
   answers it.** Emit ambiguity counts, refuse to rank wide matches, and
   build the layer last.
5. **No engine fits the build contract — does not bite.** Measured: three
   modules, pure Go, CGO off, ~2 s per pass on this repo.

## Verdict on the "deterministic research" framing

The stated goal, "make code research deterministic", is wrong as stated. The
working hypothesis survives with one sharpening. The reader is a pure
function of the tree at a SHA; research stays model judgment. The achievable
target holds: every load-bearing enumeration-class claim in a research
asset, a spec seam section, or a review finding can cite a reader invocation
that a hook or reviewer re-runs.

The sharpening: **citation replay does not need the resolver.** A replayable
citation is SHA, argv, tool version, and output hash — and `rg` and
`bench outline` satisfy that today. The resolver changes what a citation can
say for Go: a verified edge set instead of a textual sample. It does not
change whether citations replay. The spec can therefore ship the citation
convention and the resolver as separable deliverables.

Benzi's own record supports the "research stays judgment" half. With a full
resolved index mounted, its 500 trajectories still ran 7,822 `shell`, 4,641
`read_source`, and 3,127 `benzi_grep` calls against ~824 index queries. The
index informed the model; it did not replace reading. So aim the reader at
citations, blast lists, and review enumeration — the evidence layer — and
leave navigation to the harness.

Claim classes the target does not cover; the spec says so explicitly:

- **Runtime and behavioral claims** (concurrency, performance, wire
  compatibility). These need runnable probes; the craft-research contract
  already routes them there (`decisions/craft-research.md`, #7).
- **Dynamic edges**: reflection, `go:linkname`, `exec` of external binaries,
  and shell-to-Go CLI invocation edges. Registry rows
  (`cmd/bench/main.go:77`) and `rg` sweeps cover these as candidate-class
  citations only.
- **Absence claims in unresolved languages.** For a non-Go file the honest
  citation stays a textual sweep, and the claim stays a sample.
- **Quality and intent claims.** No reader warrants "this seam is right".

## Conflicts

- Benzi's two live surfaces disagree on the SWE-bench cost ($37.33 versus
  $28.92); unresolved.
- A replay of Benzi's `lines_read.py` gives 10,005 against the published
  9,125; probably data drift, unresolved.
- `gotreesitter` withdrew one performance headline, and its benchmark framing
  is unusual; its claims stay unverified.

## Residual unknowns

- No SCIP indexer was executed; every SCIP timing and reproducibility
  statement rests on documentation and source reading.
- The gopls initial-workspace-load series on perf.golang.org was
  unretrievable (404), so absolute gopls load times stay unpublished.
- Benzi's full tool list is unpublished; 30 of the "35+" names were
  reconstructed from trajectories.
- Whether canonical symbol identity fully closes the alias-naming tail is
  unproven; the FT191 spec names it as a probe (see §5).
- The language mix of the first external validation target is reviewer input
  for kill criterion 1.

## Verification record

- The coordinator read every cited local file directly in this session.
- The coordinator re-ran the recommendation-deciding probe (dependency
  delta, `CGO_ENABLED=0` build) on this machine, 2026-08-28.
- The delegate call-graph measurements (timings, counts, alias tail) were
  not re-run; they are single-source, from this machine.
- Web claims carry URLs from three read-only delegate returns, retrieved
  2026-08-28. The coordinator cross-checked the joins between returns (for
  example, `scip-go` resting on `x/tools` in both surveys) without
  re-fetching pages.
- Benzi repo quotes come from raw files at a pinned commit; benzi.fly.dev
  quotes passed through a summarizing fetch and sit one step from source.

Sources not opened: the Benzi engine and its transcripts (withheld), the
perf.golang.org series (404), and `gotreesitter` internals.
