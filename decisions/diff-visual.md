# Diff comprehension visual

## Destination

An opt-in Bench artifact that turns a finished change into an interactive visual:
how the changes flow together, where to start, what to focus on, and confidence
for production readiness. The agent writes a small structured report; kit
scaffolding fills in the rest and renders it. Principle: you can outsource your
thinking, but not your understanding. Show, don't tell.

## #1: Is this a Bench kit feature or a standalone tool?

Type: Grill

### Question
Placement decides who owns distribution, agent conventions, and build standards.

### Answer
Bench kit feature. The "generated when I review" flow is the kit's own workflow
seam, and the kit already owns agent-report conventions and distribution to
linked repos.

## #2: Where do the connection edges come from?

Type: Grill

### Question
Agent-declared vs static analysis vs hybrid — decides the report schema,
language posture, and build size.

### Answer
Agent-declared, over a CLI-emitted node skeleton. The CLI parses the diff and
pre-lists every node with a short stable ID; the agent only draws typed edges
between handed IDs and adds annotations. The validator rejects any edge naming
an unknown node ID, so invented nodes are structurally impossible. Token-cheap:
agent output is a few dozen structured lines referencing short IDs. Static
analysis rejected — see Out of scope.

Amended by #10: the agent may additionally declare untouched context nodes
(rendered grey) so the before/after chains can show where the change sits.
They are a distinct declared class, not diff nodes, and the validator admits
them as such.

## #3: What form does the rendered visual take?

Type: Grill

### Question
HTML file vs terminal vs served app.

### Answer
A single self-contained HTML file: inline JS/SVG, interactive graph (pan, zoom,
click-through), opened in the browser. No server, no external fetches.

## #4: What granularity is a node?

Type: Grill

### Question
File vs symbol altitude — fixes the skeleton schema and graph readability.

### Answer
One node per changed file. Symbols revealed by git's hunk headers (built-in
xfuncname patterns, no analysis dependency) are listed inside the node, shown
on click/hover; an edge may name the symbol it attaches to.

## #5: When is the report authored, and by whom?

Type: Grill

### Question
The authoring moment decides staleness risk and who scores confidence.

### Answer
Opt-in step of `/bench-final-check`, after the gate passes and the commit
lands: the diff is frozen and the report is stamped to that commit hash, so
staleness is detectable. This sits immediately upstream of the moments the
reviewer actually reads it — the push and PR decisions. The underlying CLI
subcommand also works standalone for ad-hoc runs. Reports carry an author-role
stamp because a final-check author is usually the implementer (self-graded
confidence stays visible as such).

## #6: What does the scoring component consist of?

Type: Grill

### Question
Graph-wired signals vs a Greptile-style dimension panel.

### Answer
Three signals, each with a visual job: an overall production-readiness
confidence with a mandatory one-line justification (the headline); a per-node
attention level (low / focus / hotspot) driving the graph's heat; a
tested/untested badge per node. Standing rule: a score exists only if it drives
a visual encoding or changes a reviewer action. Full dimension panels rejected —
see Out of scope.

Amended by #10: the coverage badge is three-valued — direct tests /
end-to-end only / test-infrastructure — and names the covering test; and the
attention level carries mandatory justification (hotspot: why + remediation +
stakes; focus: why), applying the standing rule to heat itself.

## #7: How is the edge vocabulary defined?

Type: Grill

### Question
Fixed set vs free-form — decides validation and whether edge styling means
anything.

### Answer
A small fixed set (~5 kinds) owned by the validator, each with distinct visual
styling, plus an optional one-line note per edge. The exact list is #11.

## #8: Which artifacts get committed?

Type: Grill

### Question
Report and/or rendered HTML — provenance vs derived-artifact bloat.

### Answer
The JSON report is tracked in the spec's folder (`specs/<slug>/`), joining the
spec's provenance and retiring with it (small, diffable, commit-stamped). The
HTML is derived deterministically from report + git and stays untracked,
regenerated on demand by the CLI. Landing path for a spec-less ad-hoc run is
open — see Not yet specified.

## #9: Renderer layout — vendor a library or hand-roll?

Type: Grill

### Question
Automatic graph layout is the hard part of the visual; dependency posture is a
reviewer decision.

### Answer
Superseded by #12: no layout library is vendored. The dependency-precedent
reasoning here still stands for any future need, but the approved visual's
layouts are computable without one.

## #10: What does the visual actually look like?

Type: Prototype

### Question
Layout, the reading-order tour mechanic, how attention heat and tested badges
are encoded, what makes it fun — none of this can be decided in prose. Charge
the `prototype` skill: build a throwaway render of a realistic sample report
and iterate with the reviewer live.

### Answer
Prototyped live against the real FT131 change (5417f0b..a3acb3c, 31 files) and
approved. Assets beside this map: `diff-visual-prototype.html` (working copy),
`diff-visual-ft131-report.html` (shareable offline sample). The report's job
statement, per the reviewer: answer "what did the AI just build", not "is this
quality code" — the gate and review phase own quality.

The visual: a single-file terminal-styled report. Sections in order:
1. **THE PROBLEM / THE CHANGE** — two-three plain sentences each, no em
   dashes, written for a reader with marginal codebase context.
2. **CONFIDENCE SCORE** — pct + verdict + one-line justification.
3. **CALL CHAINS · BEFORE → AFTER** — one lane per trust path; before chain
   all grey with its operational pain named in red; after chain marks
   untouched context grey, modified amber with diffstat, new green-dashed,
   artifacts dotted; decision forks show the refusal branch. This comparison
   is the centerpiece.
4. **CHANGE MAP** — hub-style node graph, nodes carry coverage badge and heat
   border, click-through to reading order.
5. **OBSERVABILITY** — what an operator will see in operation, each entry
   owned by a file.
6. **RETRO** — from `.bench/retros/<slug>.md` when present, omitted when
   absent; findings phrased "what was wrong. Now what."
7. **READING ORDER** — numbered rows (heat chip, one-line guidance, coverage
   badge, diffstat, NEW tag); expansion opens with a boxed blue `+ view diff`
   toggle (GitHub-style two-gutter table, embedded), then role ("what this
   file is"), what changed, why you should care, coverage naming the covering
   test, observability, heat justification, symbols, connections. Tests and
   config/docs ripple collapse into two expandable group rows.

Cross-cutting rules settled here: files tier into flow / tests / ripple, and
only flow earns chain-actor and narrated-step slots; heat must argue its case
(hotspot: why + what would cool it + does it matter; focus: why; low:
nothing); all prose minimizes jargon; the intro animation is off by default
behind one flag, and any interaction fast-forwards it.

## #11: What is the exact edge-kind list?

Blocked by: #10
Type: Grill

### Question
Which ~5 relationship kinds (candidates: calls, data-flow, tests, configures,
ripple) earn a slot and a styling — the prototype shows which stylings actually
differentiate.

### Answer
Five kinds in the schema, three drawn. `calls`, `data-flow`, and `configures`
render as styled edges in the chains and map. `tests` and `ripple` are
member-to-owner links that render as group rows and badges, never as drawn
lines — the FT131 prototype showed drawn test edges turn a 31-file map into
spaghetti.

## #12: Which layout engine renders the map and chains?

Type: Grill

### Question
#9 chose vendored d3-force, but the approved visual is structured, not
force-directed: lane-based before/after chains plus a hub map. Decide whether
d3-force survives, a layered/dagre-style pass replaces it, or the hub uses a
simple radial heuristic with no dependency at all.

### Answer
No dependency. Chains are arithmetic (linear lanes plus a fork). The hub map
uses a deterministic radial heuristic: highest-degree node center, producers
left, consumers right. Pathological shapes are bounded by the flow-tier cap
from #10's tiering. Supersedes #9.

## Not yet specified

- Opt-in mechanism detail (reviewer ask at final-check vs profile standing rule).
- How the authoring agent reconstructs the before-chain (reads the old revision
  itself vs the CLI hands it prior state).
- Subcommand names, report/skeleton schema field names; landing path for a
  report from a spec-less ad-hoc run.
- Heuristics for suggesting generation on big diffs.
- Sharing the HTML beyond a local open (PR-host embedding).

## Out of scope

- Static-analysis edge computation — syntactic edges are the weakest review
  signal and carry a per-language maintenance tail.
- A live-served app or any runtime service (chartr-style cockpit).
- Full Greptile/CodeRabbit dimension panels — panel scores sit beside the code
  rather than driving the visual.

## Handoff

1. **Module boundaries.** Skeleton emitter (Go): parses a commit range into the
   node skeleton — ids, paths, status, diffstat, hunk-header symbols — and
   extracts per-file diffs. Report validator (Go): admits the agent's report
   against the skeleton. Renderer (Go): report + git → one self-contained HTML.
   Workflow step: opt-in authoring pass in `/bench-final-check`. CLI routing:
   plumbing subcommand(s) in the Go core, routed by `bin/bench.sh`.
2. **Contracts.** Emitter: range in, skeleton out, nonzero on unparseable range.
   Validator: report in, exit zero only when every edge names a known id or
   declared context node, every hotspot/focus carries its justification fields,
   confidence carries its one-liner, and the commit stamp matches; each failure
   names the offending field. Renderer: deterministic — same report + repo
   yields byte-identical HTML — and the output performs zero external loads.
   Exact field names and wire format are deliberately unspecified (fog).
3. **Deep vs thin.** Emitter, validator, and renderer are deep modules owning
   diff parsing, schema policy, and layout respectively. `bin/bench.sh` routing
   and the final-check hook are thin pass-throughs with no seams of their own.
4. **Black-box assertables.** Skeleton output for a fixture repo range; validator
   exit codes and diagnostics per rejection class; rendered HTML contains the
   #10 sections and omits RETRO when no retro file exists; no-external-load
   check on emitted HTML; two identical runs are byte-identical.
5. **Gate attachment.** Go tests at all three module seams plus a contract test
   that renders a fixture repo end to end. The visual's *quality* is not
   gate-observable — approving look-and-feel stays a manual verify against the
   prototype assets.
6. **Hostile-input owners.** Paths with spaces/globs/symlinks: emitter and
   renderer escaping. Hostile file contents reaching HTML (script injection via
   diff text or agent prose): renderer escaping, validator length caps. Malformed
   or fabricated agent report: validator. Stale report vs tree: commit-stamp
   refusal (#5). Oversized diffs: renderer embed policy — flagged below.
7. **Uncertainty flags.** Embedded-diff size policy for very large changes
   (FT131's 22KB was fine; a 5MB diff is not thought through). Before-chain
   provenance (fog). Opt-in mechanism and subcommand naming (fog). These go to
   the spec-writer, not silently defaulted.
8. **Rejected alternatives.** Static-analysis edges; served app; dimension
   panels; vendored layout library (#9 superseded); all five edge kinds drawn
   (#11); floating detail panel (inline row expansion won); intro animation on
   by default (off behind one flag, interaction fast-forwards).
9. **Domain watch-outs.** The JSON report is provenance and is tracked; the HTML
   is derived and stays untracked, regenerated on demand. Confidence is
   author-stamped and usually self-graded — the report renders that stamp
   visibly. Context nodes are agent claims the CLI does not verify.

Dependency order: schema core (emitter + validator) → renderer → final-check
integration. Slicing stays the reviewer's call.

## Sources

- https://github.com/rengwu/chartr (README, fetched 2026-07-30) — fed #3's
  interactivity expectations (local interactive "star-map" graph). Drift-prone:
  re-verify before citing onward.
- greptile.com/greptile-vs-coderabbit — reviewer-supplied reference for scoring
  panels; not fetched. #6 was grilled live against the reviewer's description.
- greptile.com CLI "full-fidelity" marketing image (fetched 2026-07-30) —
  reviewer-supplied; fed #10's terminal-report form (confidence box, diagram
  section, issue rows).
- FT131 diff, spec, and tickets at 5417f0b..a3acb3c — fed #10's sample data,
  the coverage-badge refinement, and the retro reconstruction. Drift-prone:
  the spec folder is retired; read it from git history.
- All other decisions grilled in-session.
