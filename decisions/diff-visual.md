# Diff comprehension visual

Status: shaping

## Destination

An opt-in Bench artifact turns a finished change into an interactive visual. It
shows how the changes flow together, where to start, and what to focus on, and
it states confidence for production readiness. The agent writes a small
structured report; kit
scaffolding fills in the rest and renders it. Principle: you can outsource your
thinking, but not your understanding. Show, don't tell. Context nodes are agent
claims that the CLI does not verify.

## #1: Is this a Bench kit feature or a standalone tool?

Blocked by: none
Type: Grill

### Question
Placement decides who owns distribution, agent conventions, and build standards.

### Answer
Bench kit feature. The "generated when I review" flow is the kit's own workflow
seam, and the kit already owns agent-report conventions and distribution to
linked repos.

## #2: Where do the connection edges come from?

Blocked by: none
Type: Grill

### Question
Agent-declared vs static analysis vs hybrid — decides the report schema,
language posture, and build size.

### Answer
Agent-declared, over a CLI-emitted node skeleton. The CLI parses the diff and
pre-lists every node with a short stable ID. The agent only draws typed edges
between handed IDs and adds annotations. The validator rejects any edge naming
an unknown node ID, so invented nodes are structurally impossible. Token-cheap:
agent output is a few dozen structured lines referencing short IDs. Static
analysis rejected — see Out of scope.

Amended by #10: the agent may additionally declare untouched context nodes
(rendered grey) so the before/after chains can show where the change sits.
They are a distinct declared class, not diff nodes, and the validator admits
them as such.

## #3: What form does the rendered visual take?

Blocked by: none
Type: Grill

### Question
HTML file vs terminal vs served app.

### Answer
A single self-contained HTML file: inline JS/SVG, interactive graph (pan, zoom,
click-through), opened in the browser. No server, no external fetches.

## #4: What granularity is a node?

Blocked by: none
Type: Grill

### Question
File vs symbol altitude — fixes the skeleton schema and graph readability.

### Answer
One node per changed file. Symbols revealed by git's hunk headers (built-in
xfuncname patterns, no analysis dependency) are listed inside the node, and
shown on click/hover. An edge may name the symbol it attaches to.

## #5: When is the report authored, and by whom?

Blocked by: none
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

Blocked by: none
Type: Grill

### Question
Graph-wired signals vs a Greptile-style dimension panel.

### Answer
Three signals, each with a visual job. An overall production-readiness
confidence carries a mandatory one-line justification (the headline). A
per-node attention level (low / focus / hotspot) drives the graph's heat, and
a tested/untested badge marks each node. Standing rule: a score exists only if
it drives a visual encoding or changes a reviewer action. Full dimension
panels rejected — see Out of scope.

Amended by #10: the coverage badge is three-valued — direct tests /
end-to-end only / test-infrastructure — and names the covering test; and the
attention level carries mandatory justification (hotspot: why + remediation +
stakes; focus: why), applying the standing rule to heat itself.

## #7: How is the edge vocabulary defined?

Blocked by: none
Type: Grill

### Question
Fixed set vs free-form — decides validation and whether edge styling means
anything.

### Answer
A small fixed set (~5 kinds) owned by the validator, each with distinct visual
styling, plus an optional one-line note per edge. The exact list is #11.

## #8: Which artifacts get committed?

Blocked by: none
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

Blocked by: none
Type: Grill

### Question
Automatic graph layout is the hard part of the visual; dependency posture is a
reviewer decision.

### Answer
Superseded by #12: no layout library is vendored. The dependency-precedent
reasoning here still stands for any future need, but the approved visual's
layouts are computable without one.

## #10: What does the visual actually look like?

Blocked by: none
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
3. **CALL CHAINS · BEFORE → AFTER** — one lane per trust path. The before
   chain is all grey, with its operational pain named in red. The after
   chain marks untouched context grey, modified amber with diffstat, new
   green-dashed, and artifacts dotted; decision forks show the refusal
   branch. This comparison is the centerpiece.
4. **CHANGE MAP** — hub-style node graph, nodes carry coverage badge and heat
   border, click-through to reading order.
5. **OBSERVABILITY** — what an operator will see in operation, each entry
   owned by a file.
6. **RETRO** — from `capture/retros/<slug>.md` when present, omitted when
   absent; findings phrased "what was wrong. Now what."
7. **READING ORDER** — numbered rows (heat chip, one-line guidance, coverage
   badge, diffstat, NEW tag). Expansion opens with a boxed blue `+ view diff`
   toggle (GitHub-style two-gutter table, embedded), then role ("what this
   file is"), what changed, and why you should care. It also shows coverage
   naming the covering test, observability, heat justification, symbols, and
   connections. Tests and config/docs ripple collapse into two expandable
   group rows.

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
lines. The FT131 prototype showed drawn test edges turn a 31-file map into
spaghetti.

## #12: Which layout engine renders the map and chains?

Blocked by: none
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
- Embedded-diff size policy for very large changes: FT131's 22KB was fine, but
  a 5MB diff is not thought through.

## Spec-writer discretion

## Out of scope

- Static-analysis edge computation — syntactic edges are the weakest review
  signal and carry a per-language maintenance tail.
- A live-served app or any runtime service (chartr-style cockpit).
- Full Greptile/CodeRabbit dimension panels — panel scores sit beside the code
  rather than driving the visual.
- A vendored layout library — #12 superseded that route because the approved
  layouts are computable without one.
- Drawing all five edge kinds as graph lines — `tests` and `ripple` remain group
  rows and badges to avoid spaghetti.
- A floating detail panel — inline row expansion won.
- An intro animation enabled by default — it remains opt-in behind one flag and
  interactions fast-forward it.

## Sources

- Path: `decisions/diff-visual-prototype.html`
  Supports: the approved FT131 prototype behind #10's layout, reading order, and interaction decisions.
  Drift: local throwaway prototype; compare it with the current decision map before relying on it.
- Path: `decisions/diff-visual-ft131-report.html`
  Supports: the shareable offline FT131 sample report, including the 31-file range and visual encodings cited by #10.
  Drift: local derived sample; re-read if the underlying report format changes.
- URL: `https://github.com/rengwu/chartr`
  Supports: #3's local interactive “star-map” graph expectations, from the README fetched 2026-07-30.
  Drift: mutable upstream repository; re-verify before citing onward.
- URL: `https://greptile.com/greptile-vs-coderabbit`
  Supports: the reviewer-supplied comparison of scoring panels used to grill #6; it was not fetched during shaping.
  Drift: reviewer-supplied external reference; fetch and verify before citing onward.
- URL: `https://www.greptile.com/`
  Supports: the reviewer-supplied CLI “full-fidelity” marketing image that informed #10's confidence box, diagram section, and issue rows, fetched 2026-07-30.
  Drift: mutable marketing site; re-find the cited image before citing onward.
- URL: `https://github.com/gibbonmi/bench/compare/5417f0b...a3acb3c`
  Supports: the historical FT131 diff, spec, and tickets that supplied #10's sample data, coverage-badge refinement, and retro reconstruction.
  Drift: historical repository evidence; re-open the named range because its spec folder is retired from the current tree.
