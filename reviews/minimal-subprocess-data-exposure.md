# Review findings — minimal-subprocess-data-exposure (FT88, commit c440010)

Advisory pickup state for `/bench-implement-spec`. The fix session deletes this
file in the same green commit that closes the findings.

## Standards

3 findings (2 minor, 1 nit). Worst: S1.

- **S1 (minor) — derivation check binds only one direction.**
  `internal/conformance/data_handling_test.go:689` asserts every
  `env.AgentPasslist` pattern appears in DATA_HANDLING.md, but not the reverse:
  the doc can advertise a pattern the constants dropped. One-source-per-fact
  (AGENTS.md) names exactly this enforcement/advertisement pair. The
  security-relevant direction (code wider than doc) is covered, hence minor.
  Fix: add the doc→constants direction so the sets are proven equal.
- **S2 (minor, judgment call) — single-escaper fingerprint is narrow (carried
  item 6).** `internal/conformance/data_handling_test.go:721` keys on the
  literal `\u%04x` verb. Reliably catches a pasted copy of the escaper — its
  stated job — but a second escaper via `strconv.QuoteRune`, `%U`, or split
  concatenation passes clean. Acceptable as shipped; recorded so the narrowness
  is a decided residual, not an oversight.
- **S3 (nit) — passlist token regex reads prose.**
  `internal/conformance/data_handling_test.go:669` captures every
  backtick-quoted name in the marked region, including Family-column prose
  (e.g. `XDG_CONFIG_HOME` at DATA_HANDLING.md:303), so an incidental mention
  could satisfy the check if that variable's table row were deleted.
  Superset-only today; tighten to the Pattern column when convenient.

Cleared: `internal/env` / `internal/sanitize` split vs split-or-grant
(justified, decision-map-planned deep modules); sanitizer single-sourcing
(every render site routes through `internal/sanitize`, `preview.go` deleted,
doc's 200-rune cap matches `internal/shift/validate.go:28`).

## Spec

1 finding (medium). Worst: P1. Coverage-map audit: **37/37 rows satisfied**, no
doubted rows.

- **P1 (medium, reopen recommended) — dashboard Roadmap/Sequence panels
  flatten to a single line (carried item 2).** `internal/dashboard/render.go:52-53`
  routes `RoadmapText` and `Sequence` through `sanitize.Controls`, which
  escapes `\n`/`\t`; both fields render inside `<pre>` (`render.go:245,248`),
  so the previously multi-line panels now show one line of literal `\n`
  tokens. Spec-literal per decision #14, but the `<pre>` consequence was
  under-priced: html/template already neutralizes markup, so the only real
  threat there was raw C0 bytes, which the deleted `dashboard.sanitize`
  stripped without flattening structure. Reviewer decision: reopen for the two
  `<pre>` fields (e.g. a newline-preserving sanitize variant) or accept the
  flattened rendering.

Carried-item verdicts: item 1 (DATA_HANDLING prose) accurate against the
shipped paths, no stub-beside-table gaps; item 3 (120-rune banner truncation)
acceptable — display-only, full objective persists in the 0600 scratch file
and commit subject; item 7 (loop marker sentinels + per-adapter byte equality)
satisfies stories 6/7 with no loop/adapter gap.

## Coverage

4 findings (2 medium, 2 low). Worst: C1.

- **C1 (medium) — rsync is an unguarded hard test dependency (carried item 4).**
  `internal/contract/surface/artifact/reproducibility_test.go:28` execs
  `rsync -a --delete` with no availability guard, unlike the sibling
  `SkipIfSubjectFileMissing`-gated tests. On a runner without rsync (tar is
  POSIX-mandated; rsync is not — e.g. minimal Alpine images) the contract
  phase reds for an environment reason, the exact class the 2026-07-20
  learning set out to kill. The overlay change itself STRENGTHENS the probe
  (tar could not mirror uncommitted deletions; `--delete` can, and nothing tar
  caught is lost). Fix: guard or fall back when the binary is absent.
- **C2 (medium) — env.allow `#` comment grammar untested.**
  `internal/env/allow.go:56` handles comment lines and DATA_HANDLING.md
  advertises them, but no case in `internal/env/env_test.go` writes one. A
  regression dropping the branch would fail-close valid files (`#` fails the
  name pattern). Add a row asserting a comment line is skipped.
- **C3 (low) — CRLF and BOM inputs unpinned.** CRLF works only via
  `TrimSpace` (`internal/env/allow.go:53-55`) and no test pins it; a UTF-8 BOM
  fails closed with a misleading "entry before any section header" error
  (`allow.go:67-68`). Both fail-closed; pin CRLF, and optionally strip or
  name the BOM in the error.
- **C4 (low) — duplicate and leading-whitespace entries unpinned**
  (`internal/env/allow.go:74-76`). Behavior is benign (harmless append;
  lenient accept) but unfixed by any test.

Carried-item verdicts: item 4 strengthened (see C1 for the one regression);
item 5 (line-routing probes moved to stdin transport) NOT weakened — model-flag
assertions unchanged, prompt-presence preserved via stdin-reading stubs, the
ledger objective-absence flip matches FT88's intent change. Opencode argv
exposure and `$(cat)` newline strip already pinned
(`internal/contract/runtime/runtime_prompt_transport_test.go:216-227`).
