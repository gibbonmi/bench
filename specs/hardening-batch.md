# Hardening batch

Status: implemented

## Problem

Four small, unrelated soft spots have accumulated in the kit, each cheap to fix and each a latent regression or friction if left. From the point of view of the people who touch the kit:

1. **A cold maintainer can't see that the model-id grammar's `$` anchor is load-bearing.** `modelid.SafeToken` rejects tokens with newlines because Go's default `$` is end-of-text (`\z`-style), not multiline — the property that stops a shell-newline injection in a bound model id. Nothing pins it, so a future refactor to `(?m)` multiline would silently accept `claude\nrm -rf` and no test would notice.
2. **`stophook.Run` — the completion oracle's I/O path — is untested in Go.** Its rc→verdict mapping (green allows, red blocks, the `rc==3` no-gate branch blocks *without* seeding a forged verdict) is exercised only end-to-end by the shell gate contracts. A change to the mapping or the no-forged-verdict guarantee could rot with the unit suite still green.
3. **A new adopter hits undocumented prerequisites.** The README never states that the from-clone install path needs Go or Node, that Windows is unsupported (WSL2 only), or the node-version-manager shim caveat — so a fresh install can wall silently.
4. **The README `.bench/` layout is stale, and three always-loaded skill descriptions carry dead weight.** The layout tree omits shipped/present pieces (`lib/`, `BENCH-reference.md`, `gate.sh`, `lines.env`), so the documented tree doesn't match reality; and the three longest craft-skill descriptions (synthesis, adr, seams) each end with a redundant trigger restatement that costs tokens in every session that loads it.

## Solution

Five factual, no-behavior-change slices, landing in one batch:

- **Corpus.** Add the newline-class rejects (trailing `\n`, embedded `\n`, CR) to the shared safe-token corpus **in a new slice consumed only by the direct `SafeToken` seam**, and iterate it from `modelid_test.go`. This pins the `$`-anchor grammar property without disturbing the conformance line-binding consumer, which parses each value through `lines.TierValue` first (see Implementation decisions for why that distinction is mandatory).
- **stophook.** Add a Go seam test for `Run` that drives a stub gate script in a temp git repo across the four verdict cases (green, red, `rc==3` no-gate, non-executable gate), asserting the return code, the BLOCKED-message posture, and the verdict-cache state.
- **README prerequisites.** State the three omitted prereqs.
- **README layout.** List the `.bench/` pieces the tree omits, with the two generated/consumer-created entries annotated so the tree documents real facts.
- **Skill descriptions.** Trim each of the three descriptions' redundant third trigger clause, preserving the concrete quoted trigger phrases and the `index:` frontmatter (so the generated skills index stays byte-equal and the gate stays green).

No enforcement, output, or grammar behavior changes anywhere: the new tests characterize behavior that already holds; the docs state facts that are already true; the descriptions shorten without dropping a trigger situation they uniquely own.

## User stories

1. As a kit maintainer, I want `SafeToken`'s rejection of newline-class tokens (trailing `\n`, embedded `\n`, CR) pinned at the grammar seam, so that a regression loosening the `$` anchor to multiline is caught before it accepts an injected shell newline in a bound model id. Line: claude-opus-4-8 / medium. The corpus is shared with the conformance line-binding seam, so placing these rows without breaking that gate check is oracle-adjacent correctness work, not a mechanical data append.

2. As a kit maintainer, I want `stophook.Run`'s gate-exec path covered by a Go seam test across green, red, `rc==3` no-gate, and non-executable-gate cases, so that the completion oracle's verdict mapping and no-forged-verdict guarantee can't rot with the unit suite still green. Line: claude-opus-4-8 / medium. The stop hook is a load-bearing oracle and the test authors real I/O — a stub gate script, a temp git repo, and verdict-cache assertions — so correctness matters more than speed.

3. As a new adopter, I want the README to state the install prerequisites it omits — Go or Node for the from-clone path, Windows-unsupported/WSL2, and the node-version-manager shim caveat — so that I don't hit an undocumented install wall. Line: claude-sonnet-5 / low. This is factual transcription of already-true prerequisites backstopped by the docs stale-command-reference scan, not compounding guidance prose, so it deviates from the profile's top-tier doc-authoring row deliberately.

4. As a cold agent session, I want the README `.bench/` layout block to list the pieces it omits (`lib/`, `BENCH-reference.md`, `gate.sh`, `lines.env`) plus the link-installed `.bench/bin/` and the dev-build `dist/`, so that the documented tree matches what's on disk. Line: claude-sonnet-5 / low. Same factual-transcription character and same docs-scan backstop as story 3, so it shares the cheap routing.

5. As every agent session, I want the three longest craft-skill descriptions (synthesis, adr, seams) trimmed of their redundant third trigger clause while keeping every concrete quoted trigger phrase, so that always-loaded context costs fewer tokens without losing a trigger. Line: claude-fable-5 / high. Skill descriptions are always-loaded guidance prose whose every word compounds across every session that loads it, which is the exact case the profile's leverage override reserves for the top tier, so the invoking session authors this slice inline rather than delegating it.

## Implementation decisions

**Corpus — a separate slice, not an append (contestable; flagged for veto).** The newline-class rejects go into a **new exported slice** in `internal/modelid/modelidtest/corpus.go` (recommend `NewlineRejectedTokens`), consumed only by the direct `SafeToken` seam in `modelid_test.go` via an added loop. They must **not** be appended to `RejectedTokens`. Reason, verified against the live tree: `RejectedTokens` is iterated by two seams — the direct `SafeToken` unit test *and* the conformance test `TestLineBindingRejectsUnsafeModelTokens`, which writes each value into a temp `.bench/lines.env` as `BENCH_TIER_MID=<value>` and parses it through `lines.TierValue`. `TierValue` splits on `\n` and strips a trailing `\r`+whitespace, so a trailing-`\n` value truncates to a safe remainder and an embedded-`\n` value truncates at the first newline — in both cases `SafeToken` then sees a *valid* token, the "is not a safe model token" diagnostic never fires, and the conformance check goes red. (The CR row survives that consumer, but keeping all three newline-class rows together in one grammar-seam slice is cleaner than splitting them.) The map named "the corpus" generically; this separate-slice placement is the refinement forced by verifying the map against the current repo, and it is the call most worth a veto — on the slice name and on whether the CR row should instead ride the existing `RejectedTokens`.

**stophook — a new `TestRun` in the existing test file.** Add `TestRun` to `internal/stophook/stophook_test.go`, alongside the existing `Active`/`Tail`/`BlockMessage` table tests. It drives `Run(stdin, wrapper, armed, stderr)` with the process CWD set to a fresh `git init` temp repo, so `Run`'s internal `git.Root()` resolves there and `gate.Record` can write `<git-dir>/bench-last-gate`. The `wrapper` argument is a per-case executable stub script that, invoked as `<wrapper> gate`, exits with the case's chosen code. `git.TreeHash` succeeds in a commit-less repo (it falls back to the empty tree), so no seed commit is needed. Each case asserts the returned code, whether the BLOCKED message reached the stderr writer, and the verdict-cache state (present with `green`/`red`, or absent/untouched for the no-gate case). The non-executable-gate case points `wrapper` at a non-runnable path so `cmd.Run` returns a non-`*exec.ExitError`, exercising the `rc = 1` treat-as-red branch — the missing/non-executable-gate hostile input the map assigned to the stub.

**README — two factual edits, no restructuring.** (a) A prerequisites statement covering: the from-clone path needs Go (to build the compiled core) or Node (to run via `npx`); Windows is unsupported, use WSL2; a global install under a node version manager (nvm/asdf/fnm/volta) relies on the PATH shim that `bench doctor` reports and repairs. (b) The `.bench/` subtree in the Layout block gains `lib/`, `BENCH-reference.md`, `gate.sh`, and `lines.env`, plus the `.bench/bin/` local CLI set annotated as link-installed and the top-level `dist/` annotated as the gitignored dev build the gate rebuilds — so the tree documents real facts rather than phantom tracked paths. Whether `.bench/bin/` and `dist/` belong in the tree at all (versus prose-only) is contestable and flagged for veto, since neither is a tracked kit path. New prose must reference only real commands so the stale-command-reference scan stays green.

**Skill descriptions — `description:` only, `index:` untouched (authored inline, top tier).** Trim the `description:` frontmatter line of `bench-craft-synthesis`, `bench-craft-adr`, and `bench-craft-seams`, dropping each one's redundant closing trigger restatement:
- **synthesis** — drop "or any time you're deciding whether a kit change earns its place" (it restates "evaluating a change to the kit itself"); keep the `/bench-update-kit` and `/bench-what-next` triggers.
- **adr** — drop the redundant "any time you're about to document a change" trigger, but keep the substantive discipline "document the resulting state, not the change" and the quoted `"write an ADR"` / `"document this decision"` phrases.
- **seams** — drop the abstract "Use whenever" paraphrases the concrete quoted phrases already cover, keeping `"where should this test live"`, `"where does the interface go"`, `"is this the right boundary to test"`, and the "even on small changes" scope note.

The `index:` field is not touched. The generated skills index in `.bench/BENCH-reference.md` derives from `index:` (and optional `index-note:`) via `.bench/skills-index.sh`, never from `description:`, so the conformance skills-index-equality check stays byte-equal and green. No SKILL bodies change. This slice is authored inline by the invoking session at the top tier, not delegated.

## Testing decisions

- **What a good test is here.** Drive real behavior at the seam: `SafeToken` directly for the grammar property, and the real `Run` orchestration (stub gate, temp git repo, observed cache and stderr) for the stop hook — never a reading of the diff. The two doc slices have no runtime seam; their gate attachment is the docs conformance scan, and their acceptance is a maintainer read of stated-fact accuracy plus the skills-index-equality check.
- **Seams and prior art.**
  - `modelid.SafeToken` — the pure grammar seam, already table-driven from the shared corpus in `modelid_test.go` (`TestSafeToken`); the new slice extends the same loop. The conformance consumer `TestLineBindingRejectsUnsafeModelTokens` in `line_routing_checks_test.go` keeps iterating only `RejectedTokens`, unchanged.
  - `stophook.Run` — new `TestRun` in `stophook_test.go`. Prior art: the chdir-into-temp-repo pattern in `internal/shift`, `internal/worktree`, and `internal/coverage` tests; the verdict-cache convention in `internal/gate.Record` (writer) and `internal/status` `appendGate` (reader).
  - README and skill descriptions — no test seam; they ride the Go root conformance scans (stale command references for the README, generated skills-index equality for the descriptions).
- **Gate.** `.bench/gate.sh` (the project gate).

### Seam diagram

Seam 1 — the grammar (`SafeToken`, pure):

    trigger: TestSafeToken iterates NewlineRejectedTokens (new slice)
        │
        ▼
    "gpt\n"      ──▶  [ SafeToken → safeTokenRe.MatchString ]  ──▶  false (rejected)
    "gpt\n5"     ──▶  [   ^[A-Za-z0-9][...]*$  (Go $ = \z)  ]
    "gpt\r5"     ──▶  [                                     ]
                          ◀ tests attach here: assert SafeToken(value) == false for each row

Seam 2 — the stop-hook orchestration (`Run`, I/O):

    trigger: TestRun calls Run(stdin, stubGate, armed=true, &stderr) in a temp git repo
        │
        ▼
    stub `gate` exit code ──▶ [ stophook.Run: exec <wrapper> gate ]──▶ return 0 / 2
                              [  rc→verdict, rc!=3 → gate.Record  ]──▶ stderr BLOCKED?
                              [  git.Root() → <git-dir>/          ]──▶ bench-last-gate cache
                                  ◀ tests attach here: assert (return code, BLOCKED-on-stderr, cache line/absence)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `SafeToken` rejects a trailing-`\n` token | `SafeToken` via `NewlineRejectedTokens` | Characterization pin (holds today). Goes red under the mutation `$`→`(?m)$` or `\z`→`\Z` tolerance, which would accept `"gpt\n"`. | Go's default `$` is end-of-text; a multiline/`\Z`-tolerant anchor would let a token carrying a trailing shell newline through, and this row fails the moment it does. |
| 1 | `SafeToken` rejects an embedded-`\n` token | same | Characterization pin. Goes red under a multiline `(?m)` switch, which would match `"gpt\n5"` line-by-line. | An embedded-newline token is the injection the grammar exists to block; only an anchor that spans the whole string keeps it rejected. |
| 1 | `SafeToken` rejects a CR (`\r`) token | same | Characterization pin. Goes red if `\r` were added to the allowed character class. | CR is not in the allowed class; the row pins that the class stays CR-free. |
| 2 | Green gate → allow | `Run` (green stub) | Characterization pin. Goes red if the `rc==0 → return 0` branch is inverted or removed. | A stop hook that blocks on green would brick every clean shift; this row fixes exit 0 + a `green` cache line on a zero-exit gate. |
| 2 | Red gate → block with BLOCKED message | `Run` (exit-2 stub) | Characterization pin. Goes red if non-zero rc were treated as allow (return 0) or the BLOCKED write dropped. | A red gate that allows the stop is the worst-class oracle bug; this row fixes exit 2, a `red` cache line, and the BLOCKED header on stderr. |
| 2 | `rc==3` no-gate → block **without** seeding a verdict | `Run` (exit-3 stub) | Characterization pin. Goes red if the `rc != 3` cache guard is dropped (cache would gain a spurious `red` line) or if `rc==3` were mapped to allow. | The no-forged-verdict guarantee means a resolved-to-None gate must block the stop yet never key a verdict to the tree; this row fixes exit 2 with the cache absent/untouched. |
| 2 | Non-executable/missing gate → treat as red and block | `Run` (non-runnable `wrapper`) | Characterization pin. Goes red if a `cmd.Run` start failure were mapped to allow instead of `rc = 1`. | A gate that can't be executed must fail toward enforcement, not silently allow the stop; this row fixes exit 2 + a `red` cache on the start-failure branch. |
| 3 | README states the three install prerequisites | docs conformance scan + maintainer read | Not TDD-able: factual doc content. Verified by the stale-command-reference scan staying green and a maintainer confirming each prereq is true. | Stated honestly as a non-gate-observable content addition; the scan only guards that no dead command reference is introduced. |
| 4 | README `.bench/` tree lists the omitted pieces | docs conformance scan + maintainer read | Not TDD-able: layout is prose. Verified by a maintainer diffing the tree against the live `.bench/` contents. | Same honest exception as story 3; there is no parser for layout-tree accuracy. |
| 5 | Each description drops its redundant clause, keeps its quoted triggers, leaves `index:` intact | conformance skills-index equality + maintainer read | The gate runs the skills-index equality check; it stays **green** because the index reads `index:`, not `description:`. It goes red only if `index:` is touched. | Directly pins the Handoff assertable "skills-index check stays green after description edits"; trigger-phrase preservation is semantic and checked by a maintainer read, stated as the honest non-gate-observable part. |

Cheapest-wrong-implementation check: an always-green stub of `SafeToken` (returns `true`) reddens every story-1 row; a `Run` that always returns 0 reddens the red/no-gate/non-exec rows; a `Run` that drops the `rc != 3` guard reddens the no-gate row specifically. The map is not passed by any of those degenerate implementations.

### Edge inventory

Edge classes walked per behavior, each resolved as a coverage row above or a **Won't handle** line here:

- `SafeToken` empty/absent input — already covered by the `empty` row in `RejectedTokens`; no new row.
- `SafeToken` Unicode line/paragraph separators (U+2028/U+2029) — **Won't handle**: they're outside the allowed character class, so already rejected by the class, not the `$` anchor; a dedicated row would test the class, not the property this slice pins.
- `Run` not armed (`armed=false`) — resolved as a small assertion in `TestRun`: returns 0, the stub is never invoked, the cache is untouched.
- `Run` already active (`stop_hook_active` true) — resolved as a small assertion in `TestRun`: returns 0 without invoking the gate; `Active`'s own detection is already covered by `TestActive`.
- `Run` empty gate output on a red verdict — **Won't handle**: `Tail("")` behavior is already pinned by `TestTail`/`TestBlockMessage`; re-asserting it through `Run` adds nothing.
- `Run` non-git CWD (`git.Root` fails) — **Won't handle**: `Record` is a documented no-op when the root can't resolve; the seam test always runs in a temp git repo, and the shell gate contracts exercise the non-git path.
- `Run` SIGINT mid-gate — **Won't handle**: process-signal handling is not this seam's concern and no behavior here changes it.
- Skill-description over-trim dropping a concrete trigger phrase — guarded by story 5's acceptance criterion that the quoted phrases remain; a maintainer read confirms, since no gate check grades description content.

## Out of scope

- **Enforcement semantics (FT27).** Changing what the guard or the stop hook *decides* — a separate behavior change with its own decision map. This batch only characterizes the current decisions.
- **Status output (FT30).** Any change to what `bench status` renders — its own ticket and seam. Out of this batch entirely.
- **Further README restructuring.** Reorganizing sections, rewriting the quick-start, or reflowing the layout beyond the two factual additions in stories 3–4 — a separate documentation pass, roughly 4 edits, 2 gate runs, deferred because it introduces prose-design decisions this fact-only batch shouldn't carry.
