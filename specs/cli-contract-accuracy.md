# cli-contract-accuracy — the CLI says what it does and does what it says

Status: staged

Source: `ASSESSMENT.md` backlog 7 + 8 (findings §1 med, §4 med/low-med/low)
plus the two open code findings from the ft9 review (committed at
`reviews/ft9.md`, drained here).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

Several small advertisement/behavior mismatches, each cheap alone, together
erode the contract agents navigate by: both phase docs show `bench commit`
without its mandatory `-m` (an agent following either verbatim exits 2), and
the command's own synopsis marks `-m` optional; `bench coverage --check`
success prints nothing (indistinguishable from no-output failure);
`bench roadmap` outside a repo exits 0 with "no ROADMAP.md" while its sibling
in the same package exits 1 structurally; `bench diff --full bogus` blames
`--full`; post-resolution git failures in `bench diff` render as empty
sections at exit 0; a control-byte commit subject hard-fails all of
`bench diff --full` (posture undecided since the ft9 review); and the real
`maps --count` / `guards --brief` flags are documented nowhere an agent looks.

## Solution

Make each surface tell the truth: fix the two phase docs and the synopsis;
give `coverage --check` a definitive pass line and route its violation lines
through the canonical error renderer; give `roadmap` the standard not-in-repo
posture; make unknown-argument errors name the offending argument; surface
post-resolution git failures as structured errors; pin the control-byte
refusal posture as the documented contract; and document the two hidden flags
in their commands' help.

## User stories

1. As a build agent following `/bench-implement-spec` or `/bench-final-check`,
   I want every `bench commit` call site in both docs to show the mandatory
   `-m "<msg>"` (and the commit package's synopsis line to agree), so
   following the doc verbatim commits instead of exiting 2.
   Line: claude-opus-4-8 / low. A mechanical call-site correction in two
   command files — the same downward deviation from the doc-authoring override
   that the ft9 spec recorded for its call-site swap.

2. As an agent validating a spec, I want `bench coverage --check` on a valid
   map to print a definitive pass line (map counted, exit 0), so success is
   distinguishable from silence per the AXI empty-state rule.
   Line: claude-sonnet-5 / low. One output line at a contract-pinned seam.

3. As an agent, I want `coverage --check` violation lines rendered through the
   canonical `toon.Errorf` (`error: <kind> — <hint>`), so the error shape has
   one renderer as the AXI contract intends.
   Line: claude-sonnet-5 / low. Internal single-sourcing with byte-visible
   output pinned by contract assertions.

4. As an agent outside a repo, I want `bench roadmap` to exit 1 with the
   structured not-in-repo error (and keep exit 0 "no ROADMAP.md" for a repo
   without a roadmap), so absence-of-repo and absence-of-roadmap stop
   conflating.
   Line: claude-sonnet-5 / low. A two-branch posture fix pinned by the
   fail-closed contract suite. Build note: lands cleanly after (or with) the
   one-source-collapses not-in-repo story; the phrase source is shared.

5. As an agent mistyping a flag, I want `bench diff --full bogus` (and any
   valid-flag-then-unknown-arg order) to report `unknown argument: bogus`, so
   the usage error names the actual offender.
   Line: claude-sonnet-5 / low. Arg-attribution fix at the diff seam. Build
   note: the review-after-merge spec rewrites the same parser into a loop —
   build these together or that spec first; both specs carry this note.

6. As a review agent, I want a post-resolution git failure in `bench diff`
   (files listing, log, or `--full` body) to surface as a structured error at
   exit 1 instead of rendering empty sections at exit 0, so emptiness always
   means "no changes", never "git broke".
   Line: claude-sonnet-5 / medium. Error-path plumbing at a known seam; the
   PATH-shim probe pattern already exists in the test helpers.

7. As a review agent, I want the control-byte-subject posture pinned and
   documented: a commit subject carrying a control byte makes `bench diff
   --full` exit 1 with the structured unrepresentable-cell error (matching the
   existing posture for control-byte paths), and the `--full` help names the
   refusal, so the ft9 review's open finding becomes a decided, tested
   contract.
   Line: claude-sonnet-5 / low. A pinning test plus one help sentence; the
   behavior already exists. (Default call, flagged: graceful degradation —
   sanitizing subjects — was the alternative; consistency with the existing
   path posture won. Parked in Out of scope.)

8. As an agent discovering surfaces, I want `bench maps -h` to document
   `--count` and `bench guards -h` to document `--brief`, so real flags stop
   being findable only in Go source. The main usage heredoc stays plumbing-free
   by design.
   Line: claude-sonnet-5 / low. Two help strings at their command seams.

## Implementation decisions

- **Doc fixes pin by anchor, not by prose review**: the two `bench commit`
  call sites get their `-m` requirement pinned by tightening the existing
  docs-anchor family, the ft9 story-5 move.
- **The pass line is one definitive sentence** (`ok: coverage map valid — N
  row(s)` shape, exact wording at build) on stdout exit 0; multi-map behavior
  is unchanged (`--check` takes one spec today and keeps doing so).
- **`roadmap`'s repo-present behavior is untouched** — missing/zero-byte
  roadmap remains the maintenance prompt at exit 0; only the no-repo branch
  changes, adopting `toon.NotInRepo` exit 1 like its sibling `idea`.
- **Unknown-argument attribution scans for the first unrecognized argument**
  rather than blaming `args[0]`; the fix lands in the diff parser (the cited
  instance) and the same attribution rule guides the parser rewrite the
  review-after-merge spec performs.
- **False-empty fix keeps the loud-base guarantee layered**: base resolution
  failures already error loudly; this change makes the three post-resolution
  calls (`changedFiles`, `commitLog`, body) propagate their errors through the
  structured renderer instead of returning nil. The ft9 review's "emptiness as
  success" tradeoff is thereby decided in favor of honest errors.
- **Control bytes: refusal is the contract.** `toon.Table`'s refusal already
  fires; the work is a pinning contract test (ESC in a fixture commit
  subject), a help-text sentence, and the profile checklist line (promoted at
  the ft9 retire). No sanitization layer.
- **`reviews/ft9.md` is deleted by the green commit that lands stories 6–7**
  (its two open code findings drain into them; the third finding is decided in
  the one-source-collapses spec), closing the pickup file's lifecycle through
  the sanctioned route.

## Testing decisions

- **What a good test is here:** black-box probes of the built binary in
  fixtures — exit codes, stdout channel, exact phrases. Prior art:
  `internal/contract/axi/axi_wave2_test.go` (diff postures),
  `axi_fail_closed_test.go` (not-in-repo), the coverage contract probes, and
  the PATH-stub helpers used by conformance for subprocess shimming.
- **Seams:** the affected command surfaces (`coverage`, `roadmap`, `diff`,
  `maps`, `guards` help) via the AXI contract; the two phase docs via the
  conformance anchor family.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: contract suite drives each corrected surface
        │
        ▼
    fixture repos          ──▶  [ coverage --check | roadmap |      ]  ──▶  definitive pass line /
    non-repo cwd           ──▶  [ diff [--full] | maps -h |         ]        structured errors /
    PATH-shimmed git       ──▶  [ guards -h                          ]        honest exit codes
    ESC-subject commit     ──▶  [                                    ]
                      ◀ tests attach here: each story is one probe family asserting
                        stdout text + exit code; the git-failure family runs the binary
                        with a shim that fails a chosen git subcommand.

    trigger: `bench gate` (conformance docs anchors)
        │
        ▼
    bench-implement-spec.md / bench-final-check.md  ──▶  [ docs anchor family ]  ──▶  red until
                                                                                       `-m` appears at both call sites

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Both phase docs show `bench commit -m` at their call sites; the synopsis marks `-m` required | conformance docs anchor | tighten the anchor to require the `-m` token at both call sites — red today (both omit it) | an agent-misleading call site can't return without turning the anchor red; the synopsis rides the same diff, review-checked |
| 2 | `--check` on a valid map prints the pass line at exit 0 | coverage stdout (AXI contract) | run `--check` against a valid spec fixture — stdout is empty today, so the `ok:` assertion fails | pins the definitive-empty-state rule; silent success can't come back |
| 3 | `--check` violations render as `error: <kind> — <hint>` via the canonical renderer | coverage stdout (AXI contract) | assert the exact canonical shape against a broken-map fixture — the current hand-rolled line already matches loosely, so the row pins shape + single renderer; red only if wording shifts during the change | guards the one-renderer rule while the line is touched; divergence from the canonical shape fails |
| 4 | `bench roadmap` outside a repo exits 1 with the structured not-in-repo error; inside a repo without ROADMAP.md it keeps exit 0 | roadmap stdout (AXI contract) | run from a non-repo cwd — today exits 0 with `no ROADMAP.md`, failing the exit-1 assertion | the conflation is the defect; the paired in-repo probe stops an overcorrection |
| 5 | `bench diff --full bogus` reports `unknown argument: bogus` | diff stdout (AXI contract) | the probe asserts `bogus` in the usage line — today it renders `--full`, failing | misattribution is directly asserted away |
| 6 | A post-resolution git failure yields a structured error at exit 1, never empty sections at exit 0 | diff stdout with PATH-shimmed git (contract) | shim git to fail `log` after base resolution — today `--full` exits 0 with `log[0]`, failing the exit-1 assertion | the false-empty class is only reachable through a failing subprocess; the shim makes it a first-class red |
| 7 | An ESC-byte commit subject makes `--full` exit 1 with the unrepresentable-cell error, and help names the refusal | diff stdout (AXI contract) | already covered by behavior, un-pinned — the probe passes on day one and is the regression guard; the help-text assertion is red until the sentence lands | the posture was undecided; the test turns it into contract, and the help assertion is the TDD-able half |
| 8 | `maps -h` names `--count`; `guards -h` names `--brief` | maps/guards stdout (AXI contract) | assert the flag tokens in each help output — absent today | undocumented real flags are the defect; the help is the agent-visible surface |

### Edge inventory

- error path → rows: not-in-repo (4), unknown argument (5), git failure (6),
  unrepresentable cell (7), plus broken-map violations (3).
- empty/absent input → rows: valid-map pass line (2) and roadmap's
  absent-file-vs-absent-repo split (4).
- boundary values → covered: the shim family (6) exercises each of the three
  post-resolution call sites (files, log, body) as separate probes.
- malformed input → row 7 (control bytes); comma/quote subjects stay covered
  by the existing escape-once contract test.
- interrupted/partial state, re-run idempotency — **Won't handle**: all
  touched commands are read-only or single-emission; no state machine changes.
- hostile environment → covered by the PATH-shim family (6): a broken git on
  PATH is the canonical hostile environment for these commands; the profile's
  missing-tool class is the same probe with the shim removed from PATH —
  **Won't handle** separately, since command startup already errors loudly
  without git (existing behavior, not touched here).

## Out of scope

- **Graceful degradation for control-byte subjects** (sanitize instead of
  refuse) — a separate rendering-policy capability that would also have to
  cover paths for consistency; the refusal posture is decided above. Estimate:
  ~4 edits, 2 gate runs.
- **A general unknown-argument attribution helper shared by every command
  parser** — the attribution rule is fixed here at its cited instance and
  inherited by the parser the review-after-merge spec rewrites; sweeping the
  other single-arg parsers is a separate mechanical pass. Estimate: ~6 edits,
  2 gate runs.
