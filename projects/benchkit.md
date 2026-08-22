# Project: benchkit

Bench kit is the harness-agnostic agent-development workflow. The npm package
`redbench` ships it. Bench kit is not an application. It is shell, markdown, and
JSON that other repos consume. The deliverable is the `bench` CLI (`bin/bench.sh`),
the working agreement (`AGENTS.md`), the portable `.agents/` skills and commands, and
the harness adapters that call the shared `.bench/hooks/` scripts. The artifacts are
plain files, so the kit must work identically under Claude Code, Codex, and any other
AGENTS.md harness. This portability is the product.

## Working branch

`main`. (`/bench-final-check` names the commit-on-green policy as canonical. The
pre-push hook, not the commit step, guards the default branch, so `bench commit`
works on any branch. This line only states the binding.)

## Seams (test here; everything else is free to change)

- **The gate contract** (`.bench/gate.sh` / `bench gate`). This is the oracle surface.
  Every Bench operation routes through its exit code: 0 means shippable, and non-zero
  means not done. This is the highest seam: a weak gate weakens the whole system.
  Test it by feeding it a conformant tree (green) and a broken tree (red); never trust
  a reading of the diff. The gate package is the single deep owner of
  reusable-verdict authorization and durable execution; ADR 0002 records the
  accepted trust posture.
- **The `bench` CLI subcommands.** This is the operational shell surface. `bench help`,
  rendered from the Go `commandRegistry`, is the executable inventory, while
  `.bench/BENCH-reference.md` keeps category-level operational guidance. The contract is the
  stable command names and exit codes; the
  implementation behind each stays free to change. Keep the gate resolution
  (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) in one place.
  Missing-binary resolution stays network-silent by default and names the explicit
  `bench repair` action; automation opts into the same repair path with exactly
  `BENCH_REPAIR=1`, while `BENCH_OFFLINE=1` and `BENCH_NO_REPAIR` suppress it.
  Reviewed spec-backed builds keep serial green ticket commits on one retained
  integration source. Semantic review binds its frozen base and tip;
  `bench worktree land` composes and gates that pair on the destination, and its
  published commit owns the spec's `Status: implemented` flip.
  `bench worktree reclaim` is the one reader of the `$BENCH_HOME/worktrees` pool
  parent: bare, it plans reclaimable keys and removes nothing, and
  `--apply <fingerprint>` removes exactly what that plan named. It owns the single
  reclaimability predicate; `bench resume-clean` counts through that same predicate,
  reports the count and the verb, and never removes a pool key itself.
- **The AXI query surface** (`bench anchors`, `bench learnings`, `bench maps`, `bench guards`,
  `bench diff`, `bench coverage`, `bench roadmap`, and `bench worktree list`, and the
  shared flat-table TOON emitter behind them). This is the agent-facing read-only
  surface, and the AXI-conformant half of the hybrid output contract: TOON stdout,
  definitive empty states, structured errors on stdout, exit 0/1/2. The AXI contract
  fragments gate-test it. Each guard script's static leading-comment header supplies
  the guard manifest (Bench reads it as data and never executes it), and
  `bench guards --brief` is the surface the SessionStart hook injects.
  `bench diff` is the single coherent review snapshot: it reports the revision,
  aggregate, inventory, checkout, whitespace, and optional complete patch from
  one movement-checked read. Spec-backed review supplies the retained source's
  explicit frozen base; bare mode retains its recorded-base fallback for other
  work. `bench coverage --check` is
  the one parser for the acceptance-coverage-map convention; the gate's docs
  fragment consumes it instead of carrying its own.
- **The ambient dashboard** (`bench status`). This is the single deterministic renderer
  the SessionStart hook and the user both call: it ranks the signals that fire on a
  fixed severity ladder and leads with the next action. It reads gate state from the
  **gate cache** (`<git-dir>/bench-last-gate`, written durably by gate execution) —
  never from a cold gate run. The
  contract (gate-tested) shows only on signal, keeps a five-row budget, treats a
  stale-green as not a clean bill, and gives one combined capture-drain row (parked
  ideas plus open learnings) pointing at `/bench-drain`. A stale exact verdict always
  stays the strong stale row; no path-based reduced-scope softening applies. Its
  severity-1 git signal reports dirty paths from the named or current checkout while
  aggregating unpushed commits and unique local branches across the repository;
  severity-2 intent joins the shared common-directory ledger, compact by default and
  expanded by `--all`.
- **The capture inbox and working roadmap** (`bench idea` → `capture/IDEAS.md`;
  `bench roadmap` → the `ROADMAP.md` index and its `roadmap/FT<n>.md` detail
  owners, one per row, each holding that row's body, `Occurrence:` ledger, and
  `Sources:` line). This is capture-and-forget: park an out-of-scope idea and
  commit to nothing; ideas graduate only through a `/bench-drain` drain into the
  working roadmap. The contract (gate-tested in a throwaway repo): `idea` appends one
  dated line and creates the inbox; a no-arg `idea` errors without appending;
  `roadmap` prints the working document plus drain status, or its
  `## Recommended sequence` when nothing needs draining. `capture/IDEAS.md`,
  `ROADMAP.md`, and `roadmap/` stay per-consumer content — never in the kit's
  `package.json` `files[]`.
- **The kit content surface** (`.agents/skills/*/SKILL.md`, `.agents/commands/*.md`).
  This is portable harness-facing content. The contract is structural: every skill
  carries YAML frontmatter (name and description) and follows progressive
  disclosure; every command is a phase the index names. The `.bench/BENCH.md` skills
  index is generated
  from each skill's `index:` frontmatter (`bench skills-index --write`);
  craft skills use `craft-*` visible names so `$bench` menus show
  only human-run phase adapters. Codex derives command-adapter skills from
  `.agents/commands/` and documents them in `.bench/BENCH.md`. The gate's conformance
  layer enforces those contracts so disk, docs, and adapters do not drift.
  `.claude/` paths are adapters, not a second source of truth.
- **The safe managed-asset lifecycle** (`bench link` and `bench unlink`). This is the
  adoption surface: it preflights and stages the complete write set, syncs durable
  content, and atomically promotes it or rolls the repository back. Relink reconciles
  old and new manifests: clean removed assets leave, while modified or project-owned
  collisions remain in place and produce a machine-readable partial result. Unlink
  removes only clean manifest-owned assets and reports any residuals with the same
  partial posture.
  The lifecycle preserves project-owned `AGENTS.md` text and installs the `.bench/bin/`
  local CLI set the shared hooks use when no global `bench` command is on PATH.
- **The distributable artifact contract** (wrapper and native package tarballs).
  The exact tarballs are the acceptance subject: one private-staging builder consumes
  the canonical asset manifest and platform matrix, artifact tests inspect and install
  its output, and the native workflow executes the same host smoke used locally. The
  installed shim keeps maintenance on the installed kit and routes operations through
  the linked repository's tracked launcher; staging details stay free to change.
- **AGENTS.md** — the canonical working agreement for project-owned content. `CLAUDE.md`
  imports it (and `.bench/BENCH.md`); never duplicate content there. The four invariants
  and the communication rules stay canonical in `.bench/BENCH.md`; the craft skill and
  command indexes live here. The gate checks those indexes, checks command-adapter skills
  against `.bench/BENCH.md`, and checks that no shared rule's literal marker phrase
  reappears in AGENTS.md — a substring check, so it reds a verbatim copy; a
  paraphrased restatement is review's to catch, not the gate's.

No UI. There is **no design source** for this repo.

## Hostile-input checklist (shell CLI)

`/bench-write-spec`'s edge inventory walks these edge classes for this domain —
the hostile inputs shell CLIs actually meet. Walk every class before locking a
coverage map; a class skipped here returns as a regression.

- paths and directory names containing spaces or glob characters
- control bytes (ESC, BEL) in git-sourced text — commit subjects, branch names,
  paths — which `toon.Table` refuses rather than renders
- control bytes a sink *permits* but cannot survive: `toon.Representable` allows
  tab, newline, and return because the encoder escapes them, so a line-structured
  markdown or single-line sink borrowing that predicate accepts a value that
  splits its own field. Assert the permitted bytes, not only the refused ones —
  a test that exercises only what the predicate rejects proves nothing about the
  half that reaches the document
- a whitespace predicate applied to hand-edited text, where the class boundary
  is what a reader can see. `unicode.IsSpace` carries the whole `White_Space`
  property, so every `Zs` separator — U+00A0, U+2000 through U+200A, U+3000 —
  is a separator, while U+200B and U+FEFF are not. An ASCII-only predicate
  borrowed from a wire-format class silently drops the first group; a predicate
  widened to zero-width characters makes the rule unreadable to the person
  repairing the line. Name the exact predicate and assert both sides of it — a
  test that only exercises the accepted side cannot tell the two apart
- a command whose own write changes a fact it reports: an artifact-rewriting
  command that also states tree cleanliness, staleness, or a derived next step
  falsifies its own output the moment it lands. Assert repeated application in
  the *tracked* configuration, and decide per field whether it excludes its own
  write or states the post-write truth — an untracked fixture holds the tree
  still and passes either way
- a path read out of a file, which is not evidence about anything until it is
  resolved the way its writer resolves it. A git worktree's `.git` pointer may
  carry a relative `gitdir:`, which git resolves against the pointer's own
  directory; statting it against the process working directory answers a
  question about a different path, and a command that treats the resulting
  absence as proof will act on live data. Trimming is the same trap — a path may
  legitimately end in a space. Decide per field whether a value is absolute, and
  treat unresolvable as unknown rather than as absent
- hand-edited files whose last line lacks a trailing newline
- absent file vs present-but-empty file (distinct behaviors, both asserted)
- special files in any discovered path — script discovery, control-record reads,
  the spec and decision sweeps (FIFOs, devices, sockets) — must be rejected
  before reading so neither static inspection nor an ambient command can block
- a dangling symlink where a file is expected: a plain read reports it as
  not-found, so a reader that does not stat first classifies a broken link as an
  authoritative empty state
- a *live* symlink where a generator's input file is expected: following it makes
  bytes outside the graded tree authoritative, and a generator that then rewrites
  that path writes through the link to a target the tree never named. Refusing a
  broken link is not enough — the working link is the destructive half
- unquoted multi-word arguments (`$*` vs `$1`)
- lifecycle guidance that names every sanctioned operation but routes one step
  through raw Git anyway; swap the route while preserving all command tokens, so
  a token-presence check cannot pass synthesized commit, ref, replay, or worktree
  plumbing
- an operator or a command name arriving JSON-escaped, where a hook or guard reads a
  tool envelope rather than a shell: Go's `encoding/json` HTML-escapes `&`, `<`, and `>`
  by default, so `&&` reaches the reader as `\u0026\u0026` at least as often as literal.
  A decoder that maps every escape to one opaque placeholder preserves word structure and
  destroys separator structure, which welds two commands into one token and takes the
  second verb out of command position. Assert the escaped spelling of a separator, not
  only the literal one
- an operator *run* read against a list of spellings: a lexer emits `|&`, `;;`, and `;&`
  as one token each, and a membership test over the common operators matches none of
  them, so the word after the run stays out of command position. Decide the class by
  shape — a token made only of operator characters — and name the exception (redirection
  opens no command position) rather than enumerating the inclusions
- a flag's value read as a positional: a parser resolving a subcommand or
  positional as "first token not `-`-prefixed" skips flags but not their values
  or anything after `--`, so `git stash -m list` and `git stash -- list` resolve
  to `list` and an allow verdict lands on the mutation the guard exists to
  refuse. The assertion must supply the sought token only as a flag value or
  pathspec — a case that also spells the positional explicitly passes for the
  wrong reason
- a grammar token quoted in surrounding prose: a documented example of an inline
  annotation (a backticked `(covers AB1)` in a row's own text) parses as live
  syntax unless the grammar anchors the annotation to its position
- non-ASCII whitespace in hand-edited markdown: an NBSP where a
  position-anchored grammar permits only space and tab silently unanchors the
  token, so the diagnostic must stay fail-closed rather than false-positive
- required tool missing from PATH (no global `bench`, no `readlink -f`)
- invocation through a symlink rather than the real path
- invocation through every shipped surface: real kit CLI, linked-repo by-path
  CLI, hooks, and adapters must all reach the same routed implementation
- destructive worktree state: foreign or identity-mismatched registrations, reused
  paths, the primary checkout, ignored residuals, dirty nested repositories, and
  plan/apply drift all fail closed without losing recovery state
- a two-sided merge of a path the tool itself rewrites before composing: when
  the path exists at the merge base and only one side changes an attribute —
  file mode is the live case — the merge silently auto-resolves to the changed
  side, so a wrong value on the rewritten side raises no conflict and an
  edit-shaped fixture stays green with the plumbing broken. Only the shape where
  both sides add a path absent from base and destination turns a differing
  attribute into a real conflict, so a red-capable assertion drives that add/add
  shape — the edit shape proves nothing about what the rewrite carried
- state serialized by one process and reloaded by a fresh one: the writer's
  in-memory value and the reader's re-parse agree at unit level and diverge
  across the boundary, so the assertion drives a second process rather than
  reusing the first's structures. Recomposition and recovery suites that stop
  at the first success prove one path and leave every other recomposition
  route unwalked
- cwd deeper than the repo root when the command assumes root
- non-TTY stdin on a prompting command must fail closed naming its
  non-interactive flags; `/dev/null` stdin reads as a character device, so
  TTY-detection contracts must drive a pipe, not the default null device
- host-backed filesystems under host-side I/O pressure: on WSL2, ext4 lives
  behind a VHDX whose file and directory `fsync` calls can stall for seconds
  even when guest-side CPU, memory, and `fsync` stress stay green
- an unterminated delimiter in authored Markdown: a fence, an HTML comment, or
  a frontmatter block that never closes. A parser that strips or skips such a
  region must decide where the region ends. Name that decision and assert it. A
  parser that swallows the rest of the file grades nothing after the opening
  marker

Known residual risk: `bench setup`'s real-TTY confirm wiring is one untested
constructor line binding stdin. Testing it needs a pty dependency, which is a
reviewer decision the FT76 spec deliberately left open.

## Gate (`.bench/gate.sh`)

The dev gate answers one question for the exact current subject: does the complete kit
work from this tree? A non-reused run owns one selected Bench executable built through
`scripts/go-build.sh`, passes its cleaned absolute path as `BENCH_RUN_BINARY`, and
removes the private build after every terminal outcome. Ordinary commands route through
the production Go `Command` registry; wrapper, executable identity, freshness, and
process behavior remain at the bounded system seam.

The kit phase table is exactly:

| phase | authoritative argv |
|---|---|
| `gofmt` | `bench gate-go gofmt <root>` |
| `vet` | `go -C <root> vet ./...` |
| `test` | `go test -count=1 ./...` |
| `race` | one `go test -race -count=1 -v` invocation derived from `internal/racetests.Tests` |
| `system` | `go test -count=1 -tags=system ./internal/systemtest` |
| `shellcheck` | the stable shell-file inventory, optional when shellcheck is absent |

Go owns package scheduling inside the one ordinary test driver, and that driver grades
the live tree: the `test` phase carries the graded root and the dev tier to the
conformance entry point, so the registry's checks run inside the oracle rather than only
under `prep-release`. The phase materializes that environment on the same terms as race
and system — kit-only, and only where the entry test is declared — so a linked repo stays
unaffected. There is no separate
contract or conformance dev driver, per-package loop, nested Go test, fixture-executing
canary phase, component partition, or stripped-subject phase schedule. An
environment-class skip observed by the oracle is red and names the test that emitted it:
a check the gate failed to stage has no verdict, so it cannot count as green. The race runner
verifies every registry sentinel executed. The tagged system package has one
`TestMain` owner, at most three disposable repositories, one selected executable
identity ledger, teardown on green/red/interrupt/timeout, and exactly one
stripped-distribution journey beside one adoption journey, which adopts a
disposable repository with `bench setup --yes` and drives its scaffolded gate
through the installed wrapper under a private `BENCH_HOME`.

The five command decision domains — gate, adopt, preflight, canary,
and freshness — consume immutable values in process. Their ordinary tests create no
repositories and start no operating-system processes. `internal/git` owns the one
ordinary repository adapter; `internal/gate` owns the one ordinary controlled process
group adapter.

Canary fixtures are immutable inputs to registered conformance checks. Each retained kit
fixture has exactly one check owner; its ordinary mutation test calls that owner directly,
requires the fixture-specific red, restores the fixture subject, and requires that red to
disappear. The top-level `bench canary` command validates and aggregates the complete
canary inventory without invoking owners or starting a successor process. Linked repos
receive that inventory validation and own planted-reason proof in their native tests. One
tagged system journey proves the selected executable reaches the production inventory path.

The workflow-guidance family pins the spec-to-ticket handoff from identified rows and
approved ownership fences through ticket evidence, ledger review, and fence-drift
repair. Its auto-discovered mutations keep each clause in
the section where a fresh agent acts on it.

The conformance registry remains the single source for check order, subject, input
derivation, implementation, and canary family ownership. This table is the profile's
current-state advertisement of its non-meta input bindings:

| conformance check | input source |
|---|---|
| `kit-compliance` | `catch-all` |
| `canary-fixture-compliance` | `catch-all` |
| `load-validity-metadata` | `catch-all` |
| `skills-index-command-adapters` | `catch-all` |
| `docs-currency-workflow` | `catch-all` |
| `gate-entry-contract` | `gate-entry` |
| `ordinary-build-census` | `catch-all` |
| `offline-smoke-proof` | `offline-smoke` |
| `handoff-shape-single-source` | `catch-all` |
| `harness-prefix-single-source` | `go-source` |
| `package-shipped-surface` | `catch-all` |
| `line-routing` | `catch-all` |
| `package-core-guard` | `catch-all` |
| `release-evidence-probe` | `catch-all` |
| `bench-sh-routes` | `bench-routes` |
| `default-branch-single-source` | `go-source` |
| `data-handling-derivation` | `go-source+data-handling` |
| `single-control-escaper` | `go-source` |
| `bounds-policy` | `catch-all` |
| `marker-wait-deadlines` | `go-source` |
| `subcommand-routing` | `go-source` |
| `axi-query-registry` | `catch-all` |
| `skip-ownership` | `go-source` |
| `decision-map-integrity` | `decision-documents` |
| `injected-port-registry` | `go-source` |
| `guidance-prose-budgets` | `benchkit-profile` |
| `roadmap-detail-integrity` | `roadmap-board` |
| `structure-accept-currency` | `catch-all` |
| `retro-improvement-markers` | `capture-retros` |
| `row-next-grammar` | `catch-all` |
| `prose-mechanics` | `catch-all` |

A green verdict records the exact whole subject and oracle. Reuse is allowed only for a
current exact green; partial/component and reduced-scope records are legacy input classes
that fail closed and are never authored. Prospective execution uses the same complete
phase architecture. A stale exact verdict stays stale rather than being softened by path
classification.

`BENCH_REQUIRE_CAPABILITIES=1` makes capability skips fatal. Without it, capability
rows stay informational because a developer host may legitimately lack optional
facilities. The release workflows enable strict capability posture.

The ship tier remains `bench prep-release`, run once per release. It requires a current
exact dev-green verdict and owns release-evidence verification, cross-platform artifacts,
reproducibility, stress/cross-compile coverage, and publication/preflight rehearsal. Dev
green proves the complete branch-native source architecture; it does not grant publish
authority or claim release artifacts were reproduced.

### Guidance prose budgets

Guidance prose that outgrows a session's attention stops being read. This table is the
one source for how long each subject may be, and the `guidance-prose-budgets` check
parses it rather than repeating its numbers. An exact row beats the glob row a subject
also matches, so raising or lowering a budget is an edit here and nowhere else.

| subject | limit |
|---|---|
| `.bench/BENCH.md` | 180 |
| `.agents/commands/bench-implement-spec.md` | 75 |
| `.agents/commands/bench-write-spec.md` | 73 |
| `.agents/commands/bench-debug.md` | 170 |
| `.agents/skills/bench-craft-tickets/SKILL.md` | 100 |
| `.agents/skills/bench-craft-spec/SKILL.md` | 150 |
| `.agents/skills/*/SKILL.md` | 120 |

The glob row classifies a newly added skill, so a skill arrives budgeted without
anybody editing the checker. Every other `.agents/commands/*.md` file stays outside the
reviewed universe, and the `.claude/skills/*` adapter symlinks are distribution surfaces
rather than subjects — a symbolic link or special file found where a subject belongs is
refused unread.

### Prose mechanics

The `prose-mechanics` check grades the two ASD-STE100 rules a program can measure. It
reds a sentence longer than 25 words and a paragraph with more than six sentences. The
check reads every authored Markdown file under the graded root. `.bench/prose-exclusions`
names the paths the check does not grade. Each prose exclusion row gives one path and a
one-clause reason, and the reviewer owns that file. The semantic STE rules stay with
top-tier review.

## Lines (model + effort routing)

The routing rubric — the three-signal decision table and the escalation ladder —
lives in the `craft-line` skill. This section holds what is project-specific: the
binding, the cached routings, and the escalation policy.

**Harness × tier binding** (advisory candidates from `bench models`; set
2026-07-10): tier is the only shared identity, and each harness holds its own
column — no family is canonical.

| tier | codex | claude | opencode |
|---|---|---|---|
| top | `gpt-5.6-sol` | `fable` | unbound |
| mid | `gpt-5.6-terra` | `opus` | unbound |
| cheap | `gpt-5.6-luna` | `sonnet` | unbound |

These opaque safe tokens are this repo's current choices, not a namespace rule;
the token grammar and discovery posture live in `craft-line`. OpenCode stays
unadopted here, so its column stays unbound and its adapter refuses to launch
rather than borrowing another harness's ids. Machine-readable source:
`.bench/lines.env`, read by the Agent-tool hook and the shift adapters, each
naming its own harness; the `checkLineBinding` conformance check cross-checks
this table against it, so drift between the two turns the gate red. Code-author
venue follows `craft-delegate`; line choice does not override its threshold. A
headless shift declares the tier — `BENCH_MODEL=cheap` — and the adapter's own
column supplies the id.

**Escalation policy:** no standing top-tier opt-out — any bump to the table's top
row pauses and asks the reviewer, whichever harness column you run in
(the ladder lives in `craft-line`). Tier moves still get declared — no silent
escalation.

- **Skill / command / doc authoring** → **top model, high effort**. This is the
  leverage override in `craft-line`: guidance prose compounds through every
  session that loads it while the edit costs few tokens. The `craft-skills` and
  `craft-adr` skills apply. Spend here.
- **Spec and ticket authoring** → **the session holding the decision source, at
  whatever tier it runs**.
  `/bench-write-spec` accepts exactly one of three sources: a ready compiled
  map, a reviewer-confirmed current conversation, or a named reviewed artifact,
  and authors the spec and tickets from that source and the current tree. Top +
  high remains a reviewer-approved escalation. After ticket approval, a fresh
  mid-tier session starts the build. Distinct from the doc-authoring leverage
  override above: that spends the top tier on the kit's guidance prose.
- **`bench` CLI shell plumbing** → cheap model, low–medium effort at the known seam.
  This stays mechanical once the gate-resolution and worktree-pool shapes exist.
- **Gate / conformance logic** → mid effort. Correctness of the oracle matters more
  than speed — a wrong gate is the worst class of bug in a kit whose whole premise is
  "the gate is the oracle."
- **Spec-and-tickets review round** (`/bench-write-spec`) → **mid model, high
  effort**, read-only and same-family through the harness's native agent surface.
  It reviews the spec and its ticket breakdown together against `craft-tickets`
  after `/bench-write-spec` slices it; `/bench-write-spec` owns the round's
  operating protocol.
- **Review-axis delegate** (`/bench-review-implementation`, one per axis) → mid
  model, medium effort, **~1 iteration each** (three axes can run in parallel).
  This stays read-heavy: each delegate takes the full diff plus standards docs and
  runs verification commands.

## Notes for cold sessions

- Read `AGENTS.md` first — the working agreement. The four invariants and the
  communication rules stay canonical in `.bench/BENCH.md` (AGENTS.md points there); read
  that too. `CLAUDE.md` imports both; edit `AGENTS.md` or `.bench/BENCH.md`, not it.
- `CONTEXT.md` pins the ubiquitous language (gate, oracle, shift, seam, line, …). Use
  those terms exactly; do not invent synonyms.
- The kit's portability across harnesses is a closed decision. Claude and Codex hook
  adapters are interactive layers on top of shared `.bench/hooks/` scripts and the
  harness-independent substrate (the `bench shift` loop and the git `pre-push` hook) —
  never the only thing enforcing an invariant.
- The `projects/gl-axi.md` example profile is a shipped template, not a live
  project. This file (`benchkit.md`) is the profile for this repo.
- A symbolic link inside an allowlisted kit payload tree gets refused, not followed.
  Closed decision (2026-07-23): following the link would ship bytes the allowlist
  never named, so the allowlist would stop being the complete statement of what a
  consumer receives. Do not reopen it as a link/upgrade ergonomics fix.
- Never build `dist/bench` with plain `go build`; use
  `bash scripts/go-build.sh <root> <out>` so the binary carries the package
  version the version and upgrade contracts require. `bench worktree land`
  now refuses a dev executable that was not built from the current sources, and
  names that command as the remedy. The proof is self-attestation, so it cannot
  catch an executable that predates it or one patched to skip its own check, and
  the seal cannot authenticate its own writer: a hand `go build` plus a hand
  `freshness-publish` produces a seal the check accepts. The independent roots
  that survive both are the gate's private exact-source build and the operator's
  sanctioned rebuild.
- To exercise or measure a durable worktree artifact directly, invoke that
  worktree's own `./dist/bench`. `bench` on PATH resolves to the main checkout's
  wrapper and may belong to a different source tree. Gate and `bench test` runs
  do not reuse that artifact: their owner builds one private exact-source binary,
  propagates it through every ordinary child, and removes it after teardown.
- A direct `./dist/bench` invocation needs `BENCH_HOME` exported (the wrapper
  exports it; the gate-inputs closure declares it). Without it the prospective
  subject stays open, so a landing's gate can run fully green and still refuse
  as `prospective authorization refused: infrastructure`.
- Never mutate the repository while a gate is running. The gate binds its
  verdict to the starting subject and rejects a run whose subject changes.
- When SessionStart diagnoses incomplete environment closure, prepend the
  recovered tool directory to the ambient PATH. Do not replace the harness
  toolchain that is already present.
- Canary mutation tests are ordinary in-process checks. Do not add a gate, wrapper,
  `go test`, or `go run` constructor to a fixture owner; the architecture census treats
  that as a regression.
- Never stop a gate by killing only its shell wrapper. Signal `gate-run`, which owns
  teardown of the gate script's process group and descendants.
- Nothing under `.claude/` is a copy. `.claude/commands` is a git-tracked
  symlink to `../.agents/commands`; `.claude/skills` is a real directory whose
  every entry is a symlink to `../../.agents/skills/<name>`. So editing
  `.agents/` is the whole edit — there is no mirror to sync and no content drift
  for a check to guard — but *adding* a skill still needs its new
  `.claude/skills/` symlink entry created. Read "mirror" in any artifact as
  "symlink"; FT152's spec and map both assumed copied trees and specified
  mirror work that did not exist.
