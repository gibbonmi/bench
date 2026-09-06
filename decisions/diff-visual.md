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

## Notes

## Decisions so far

- [Is this a Bench kit feature or a standalone tool?](diff-visual/tickets/1.md): Bench kit feature.
- [Where do the connection edges come from?](diff-visual/tickets/2.md): Agent-declared, over a CLI-emitted node skeleton.
- [What form does the rendered visual take?](diff-visual/tickets/3.md): A single self-contained HTML file: inline JS/SVG, interactive graph (pan, zoom, click-through), opened in the browser.
- [What granularity is a node?](diff-visual/tickets/4.md): One node per changed file.
- [When is the report authored, and by whom?](diff-visual/tickets/5.md): Opt-in step of `/bench-final-check`, after the gate passes and the commit lands: the diff is frozen.
- [What does the scoring component consist of?](diff-visual/tickets/6.md): Three signals, each with a visual job.
- [How is the edge vocabulary defined?](diff-visual/tickets/7.md): A small fixed set (~5 kinds) owned by the validator, each with distinct visual styling.
- [Which artifacts get committed?](diff-visual/tickets/8.md): The JSON report is tracked in the spec's folder (`specs/<slug>/`), joining the spec's provenance and retiring with it (small, diffable, commit-stamped).
- [Renderer layout — vendor a library or hand-roll?](diff-visual/tickets/9.md): Superseded by #12: no layout library is vendored.
- [What does the visual actually look like?](diff-visual/tickets/10.md): Prototyped live against the real FT131 change (5417f0b..a3acb3c, 31 files) and approved.
- [What is the exact edge-kind list?](diff-visual/tickets/11.md): Five kinds in the schema, three drawn.
- [Which layout engine renders the map and chains?](diff-visual/tickets/12.md): No dependency.

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
