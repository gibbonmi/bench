# go-guard-port — slice 4 of the Go rewrite

Status: implemented

Map: `decisions/go-guard-port.md` (closed; child of `decisions/go-rewrite.md`, #6
slice order). This spec builds slice 4 of the 8-slice strangler: the whole
destructive-git guard absorbed into the Go core. The 418-line Python analyzer
(`.bench/hooks/git-guard.py`) and the envelope parse that `block-dangerous-git.sh`
does in Python both move behind a `bench guard-git` plumbing subcommand; the hook
becomes a thin fail-closed shim; `git-guard.py` and the hooks `__pycache__` are
deleted; and with the last `.py` gone the gate's python-parse check (gate.sh 1c)
and the `analyzer-parse-broken` canary retire. This removes the kit's last
standalone Python program and the guard's own python3 dependency; two hooks
(`stop.sh`, `check-agent-line.sh`) still parse their envelopes with inline
`python3 -c` and shed it in slice 5, so python3-the-runtime is not gone kit-wide
until then.

## Problem

The destructive-git guard is the one enforcement that makes invariant #4 real: the
agent has no push/reset/rebase authority. Today it straddles two runtimes.
`block-dangerous-git.sh` needs python3 twice per Bash tool call — once to parse the
PreToolUse JSON envelope, once to run the analyzer — and the analyzer itself is 418
lines of subtle shell tokenizing (newline-as-separator, malformed-quote fallback,
redirection stripping, one-level wrapper scan, env/xargs/timeout prefixes) that
lives only in Python. Porting the analyzer alone (the map's original slice-4 scope)
would leave the hook still needing python3 for the envelope and double its
degradation branches. So the grill widened the slice to the whole hook: the shim's
Python dies with the analyzer, and the guard no longer depends on python3. (Two
other hooks, `stop.sh` and `check-agent-line.sh`, still use inline `python3 -c`
for envelope parsing; they are slice-5 scope, so python3 leaves the kit entirely
only then, not here.)

## Solution

A new `internal/gitguard` package owns the analyzer — tokenizer, prefix/wrapper
resolution, scan, the per-verb verdicts, the deny table, the PreToolUse envelope
parse, the block message, and the describe-classes enumeration — the single source
that both classifies and advertises, exactly as `git-guard.py` was. `internal/git`
gains the two ref/branch existence checks the analyzer needs. `cmd/bench` gains a
`guard-git` plumbing subcommand (a direct case in `run()`, like `version` — it needs
stdin, a stderr verdict, and a distinct error exit the `commands` map can't express)
that reads the envelope on stdin and yields the verdict, plus a `--describe-classes`
sub-mode feeding the shim's manifest.

`block-dangerous-git.sh` becomes a one-glance shim with exactly two rims of its own:
resolve the bench wrapper (reusing the resolver `stop.sh` already carries), pipe the
envelope to `bench guard-git`, and pass its verdict through — with **binary
unresolvable → refuse anything git-shaped** and **binary errored → refuse** as the
fail-closed rims. The block message and the classification move into Go; the shim
keeps only the two degradation postures and `--describe`.

`git-guard.py`, the hooks `__pycache__`, gate.sh's 1c py-compile check, and the
`analyzer-parse-broken` canary are deleted; the three fail-closed contracts that
today blank/copy `git-guard.py` or strip python3 re-point to the binary-missing
shape; README and any link/package reference to the analyzer update — all in the
same change so the stale-reference sweep stays green.

The parity net is strong and unchanged: the 79-case allow/block matrix in
`gate-runtime-git-contracts.sh` drives the hook end-to-end through the JSON envelope
and runs verbatim against the ported shim. As insurance during the port, a
throwaway build-time differential pass runs a generated command corpus through both
the Python and the Go analyzer; the verdict diff must be empty before `git-guard.py`
is deleted. That harness is never committed — it is the map's #2 decision, not a
permanent artifact.

## User stories

1. As the destructive-git guard, I want the analyzer ported to `internal/gitguard`
   — the newline-aware tokenizer (newline-as-separator, operator-token collapsing,
   malformed-quote fallback split, redirection stripping, quoted newlines surviving
   into wrapper strings), the env/command/nohup/timeout/xargs prefix resolution, the
   one-level `sh|bash|zsh -c` wrapper scan, the fifteen per-verb verdict functions
   with their carve-outs, the deny table, and the PreToolUse envelope parse — so
   that every classification the hook makes reads from one Go source with no python3.
   Line: claude-opus-4-8 / high. This is the subtlest code in the kit and a silent
   tokenizer or verdict divergence disables the guard without any red signal, so it
   sets the analyzer idiom at the mid tier with high effort exactly where the map
   flagged the parity net thinnest.

2. As the analyzer, I want `internal/git` to gain `RefResolves` and `BranchExists`
   — the 2-second-bounded `git rev-parse --verify --quiet` checks the checkout and
   forced-creation verdicts call, preserving the fail-safe defaults (a ref that
   cannot be resolved reads as unresolvable so checkout of an unknown target blocks;
   a branch whose existence cannot be determined reads as existing so forced
   creation blocks) — so that ref-truth has one home and a hung git cannot stall a
   Bash tool call. Line: claude-opus-4-8 / medium. Small surface, but the
   fail-*toward-blocking* defaults are load-bearing safety details a wrong port
   would silently invert.

3. As the hook, I want a `bench guard-git` plumbing subcommand — a direct case in
   `cmd/bench` `run()` that reads the envelope on stdin, classifies through
   `internal/gitguard`, writes the full `BLOCKED: …` message to stderr and exits 2
   on a block, exits 0 on allow, exits 3 on a genuine failure to run (with a
   top-level recover mapping any panic to 3 so exit 2 means *only* an intentional
   block), plus a `--describe-classes` sub-mode that prints the comma-joined deny
   labels to stdout and reads no stdin — so that the deep unit owns the verdict, the
   message, and the advertisement. Line: claude-opus-4-8 / medium. The exit-code
   contract is new and safety-critical: a block that leaks as exit 0, or a panic
   that masquerades as a block, is the failure this story's discipline prevents.

4. As `block-dangerous-git.sh`, I want to become a thin fail-closed shim — resolve
   the bench wrapper (the ~8-line search inlined this slice), pipe the envelope to `bench
   guard-git` and pass its exit code and stderr through, with **binary unresolvable
   or missing (rc 127) → refuse anything git-shaped and leave the rest of the shell
   usable** and **binary ran but errored (rc ∉ {0,2}) → refuse** as the two rims,
   and answer `--describe` by composing name/boundary/why with the denies clause
   from `bench guard-git --describe-classes` — so that the guard fails safe when the
   core is gone and the shim carries no classification logic of its own. Line:
   claude-opus-4-8 / high. The fail-closed posture is the entire guarantee; a shim
   that grants authority when the binary is absent is worse than a crash, so it
   takes the mid tier and high effort.

5. As the gate, I want the three fail-closed contracts re-pointed from the Python
   world to the binary world — `gate-axi-contracts.sh`'s analyzer-missing and
   python3-missing cases become binary-unresolvable and binary-missing
   (git-shaped-only) cases, the redundant empty-`.py` case retires, a new
   binary-errored case asserts fail-closed on rc 3, and
   `gate-axi-guards-contracts.sh`'s python3-missing `--describe` manifest becomes the
   analyzer-missing manifest — so that the "cannot classify → refuse" property stays
   gate-locked while its trigger changes, and the new binary-absent posture is
   enforced rather than trusted. Line: claude-opus-4-8 / medium. Touching gate
   contracts is the worst defect class in this kit (`craft-gate`); adapting the
   triggers without weakening any assertion takes the mid tier.

6. As the porter, I want a throwaway build-time differential pass — a generated
   command corpus (the matrix cases plus fuzzed wrapper/redirect/prefix/multiline
   variants) run through both `git-guard.py` and `bench guard-git`, asserting an
   empty verdict diff — that must pass before `git-guard.py` is deleted, and is
   never committed — so that the port's parity is proven where the tokenizer is
   subtlest, not just where the 79-case matrix reaches. Line: claude-opus-4-8 /
   medium. It is the safety net for an irreversible deletion; a harness that samples
   too thin or diffs the wrong field gives false confidence, so it takes the mid
   tier even though it is thrown away.

7. As a kit maintainer, I want `git-guard.py` and the hooks `__pycache__` deleted,
   gate.sh's 1c py-compile check and the `analyzer-parse-broken` canary retired, and
   every dangling reference fixed in the same commit — README's CLI-tree analyzer
   line and guard-row prose, and any `bin/bench-link.sh` or package/link-contract
   mention of the analyzer — so that no command has two implementations, the kit
   advertises no python3 dependency it no longer has, and the stale-reference sweep
   stays green. Line: claude-sonnet-5 / medium. Mechanical deletion and reference
   edits, fully observed by `bash -n`, the gate load, and the docs stale-reference
   sweep.

## Implementation decisions

- **Package layout.** New `internal/gitguard`: the analyzer, split as the tokenizer,
  the prefix/wrapper resolution, the scan, the per-verb verdict functions, the deny
  table, the envelope parse, the block-message composition, and `DescribeClasses`.
  Each stays under the file/dir budgets the kit enforces (the 418-line Python file
  splits along its natural seams; this is not the rejected "split to meet the
  budget" — the file is being rewritten in Go, and Go packages are multi-file by
  convention). `internal/git` gains `RefResolves` and `BranchExists`. `cmd/bench`
  gains the `guard-git` **direct case** in `run()` — not a `commands`-map entry,
  because the map's `func([]string) (string, int)` signature writes its return to
  stdout and cannot express stdin, a stderr verdict, or the distinct error exit;
  `version` is the precedent for a direct case that the map can't hold.

- **Ref/branch checks stay CWD-relative and fail toward blocking.** `RefResolves`
  and `BranchExists` run `git rev-parse --verify --quiet` in the process CWD (the
  agent's working directory, where the Bash call runs — the hook does not pass a
  `-C root`), under a 2-second `context` bound. On any error or timeout,
  `RefResolves` returns false (unknown ref → checkout blocks) and `BranchExists`
  returns true (unknown branch → forced creation blocks). This reproduces the
  Python `ref_resolves`/`branch_exists` defaults exactly, including their asymmetry.

- **The verdict contract — binary owns the message, shim owns the rims.** The
  binary emits its verdict through its exit code and stderr; the shim never
  re-derives a verdict. The mapping is the load-bearing decision, stated as a table
  because prose blurs it:

  | actor | condition | action |
  |---|---|---|
  | `bench guard-git` | allow | exit 0, no output |
  | `bench guard-git` | block | full `BLOCKED: \`<label>\` — …` to **stderr**, exit 2 |
  | `bench guard-git` | genuine run failure / panic | exit 3 (recover maps panic → 3) |
  | shim | wrapper unresolvable **or** `guard-git` rc 127 | read stdin; substring `git` → `BLOCKED: guard degraded (analyzer missing) …` exit 2; else exit 0 |
  | shim | `guard-git` rc 0 or 2 | pass through (message already on stderr) |
  | shim | `guard-git` rc ∉ {0,2,127} | `BLOCKED: guard analyzer error …` exit 2 (fail closed) |

  This preserves the two distinct degradation postures the Python shim had:
  *cannot-run-at-all* (was python3-missing, now binary-missing) fails closed on
  git-shaped only so the shell stays usable; *ran-but-errored* fails closed
  unconditionally. The panic→3 recover is load-bearing, not belt-and-braces: an
  unrecovered Go panic exits with status **2**, so without the recover a crash would
  masquerade as an intentional block. With it, rc 2 *from the binary* means only
  "block". The claim is scoped below the wrapper — a bash syntax error in the
  wrapper itself also exits 2 and passes through as a block without a `BLOCKED:`
  prefix, which is still fail-closed (over-block), just not a clean verdict.

- **`--describe-classes` feeds the manifest; no second `.sh`.** The shim's
  `--describe` prints the fixed name/boundary/why and generates the denies clause by
  running `bench guard-git --describe-classes` (comma-joined deny labels from the
  same table `classify` reads — one source for enforce and advertise, as today). If
  the binary is unresolvable, the denies line degrades to `manifest unavailable
  (analyzer missing)`, the exact text the analyzer-missing contract already asserts.
  Per the map's domain watch-out, the analyzer's Go home must **not** grow a
  hooks-dir `.sh` that answers `--describe` a second time — `bench guards` aggregates
  `*.sh` manifests, and a duplicate would collide.

- **Wrapper resolution is inlined this slice, extracted in slice 5.** The shim must
  locate the bench wrapper (`.bench/bin/bench.sh`, `bin/bench.sh`, then `bench` on
  PATH) to call `guard-git`; `stop.sh` already carries this exact search as its
  `bench_cmd()`. Extracting it to a shared `.bench/lib/` snippet now would honor
  one-source-per-fact, but it buys a **new fail-open path**: a shim that `source`s a
  lib gains a failure mode the inline version cannot have — lib missing or broken →
  the shim errors before its rims run, and a PreToolUse hook exiting non-2 is a
  *non-blocking* error, so the guard silently grants. It would also add `.bench/lib/`
  surface `bench link` must ship, which story 7 does not touch. So the shim inlines
  the ~8-line search (self-contained, no `source`, cannot fail open on a missing
  lib). This is honest, incidental repetition for exactly one slice: slice 5 rewrites
  `stop.sh` into its own shim, and the resolver collapses to one source *then*, when
  both consumers are being touched anyway. The code standard's one-source rule bites
  on *durable* duplication; this is a deliberate one-slice window with a named
  collapse point, not drift.

- **The analyzer shares the binary but not the authority.** Item 9's asymmetry —
  bench's own controlled rollback inside `bench shift` runs in-process, outside the
  agent's shell, while the agent's Bash calls route through this hook — stays
  structural: classification runs *only* on the `guard-git` path fed by the shim, and
  `internal/git`'s own callers (rollback, tree-hash, the query commands) never route
  through `internal/gitguard`. The guard intercepts the agent's tool calls; it does
  not police bench's own git. Keep it that way as the packages grow.

- **Deletion is part of this change, not a follow-up.** `git-guard.py` and
  `.bench/hooks/__pycache__/` are deleted; gate.sh's 1c block (the `*.py`
  py-compile loop and its two comments) is removed with the last `.py`; the
  `tests/canary/analyzer-parse-broken/` canary is removed (its guarded surface, the
  py-compile check, is gone — verify no registry names it by string). README's
  CLI-tree analyzer line and the guard-row prose update to the Go layout; any
  `bin/bench-link.sh` or package/link-contract mention of `git-guard.py` drops. The
  `guard-describe-boundary-dropped` canary and gate.sh's 2b Codex check reference the
  surviving `block-dangerous-git.sh` and are unchanged.

- **Item 7 (Codex envelope) — settled, no escalation.** Both harnesses wire the
  same shim reading the same stdin: `.claude/settings.json` and `.codex/hooks.json`
  both run `block-dangerous-git.sh` on `PreToolUse:Bash`, and gate.sh's 2b check
  already verifies Codex fires that event and shells the command. The envelope
  carries `tool_input.command` in the shape the current Python parse assumes, so the
  Go parse reads the one shape for both harnesses. The uncertainty flag resolves by
  reading the adapters, as the Handoff allowed.

## Testing decisions

- **What a good test is here.** Acceptance drives the hook end-to-end — pipe a JSON
  envelope to `block-dangerous-git.sh` and assert exit code plus the `BLOCKED:`
  message — never Go internals. The 79-case matrix already does exactly this and is
  the primary regression net; it runs verbatim. Go table tests are additional, at
  the pure-function seam in `internal/gitguard` (tokenizer, prefixes, wrapper scan,
  each verdict, envelope parse, describe-classes), where the shell-untestability tax
  on the tokenizer finally retires.
- **Seams.** Three: the **hook enforcement seam** (the matrix + fail-closed
  contracts driving the shim → binary — prior art: `gate-runtime-git-contracts.sh`
  passes today against the Python shim); the **`go test` unit seam** (table tests
  beside `internal/gitguard` and `internal/git` — prior art: slice 2/3's
  `internal/maps`, `internal/structure` tests); and the **degradation/manifest seam**
  (the `gate-axi-contracts.sh` fail-closed and `gate-axi-guards-contracts.sh`
  `--describe` contracts — prior art: they exist today against python3-missing and
  re-point to binary-missing).
- **Gate command:** `.bench/gate.sh` (the project gate). Its Go layer already runs
  `gofmt -l`, `go vet`, `go test ./...`, `go build`, and the cross-compile targets,
  so the new packages and tests are graded with no new gate wiring; the shell
  fragments run the ported hook unchanged.
- **The differential parity pass is not gate-observable, by design.** It needs
  `git-guard.py` alive to diff against, and it is deleted before commit. It is a
  build-time step the porter runs and discards — insurance for the `.py` deletion,
  not a coverage row. Recorded as a Won't-handle line below so its absence from the
  committed gate is a decision on the page.

### Seam diagram — hook enforcement seam (the 79-case matrix)

    trigger: agent Bash tool call → PreToolUse hook (Claude/Codex) runs block-dangerous-git.sh
        │
        ▼
    JSON envelope ──▶ [ block-dangerous-git.sh: resolve wrapper → pipe stdin ]
    on stdin         [ bench guard-git → internal/gitguard.Classify         ]
                     [   ParseEnvelope → Tokenize → resolvePrefixes → scan  ] ──▶ exit 0 (allow), no output
    (git subprocess  [   → per-verb verdict → deny label → block message    ] ──▶ exit 2 + `BLOCKED: …` on stderr
     via internal/git.RefResolves / BranchExists for checkout/branch checks)
        ◀ tests attach here: gate-runtime-git-contracts.sh pipes {"tool_input":{"command":…}} to the
          shim and asserts exit 0/2 + `BLOCKED:` — the 79-case matrix, run VERBATIM. This is a parity
          net: it passes today against Python and must stay green across the port; a verdict regression
          is what turns it red. Not red-first by design (it is the parity guarantee, not new coverage).

### Seam diagram — `go test` unit seam (analyzer internals)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    command string / ──▶ [ internal/gitguard.Tokenize / resolvePrefixes / scan ] ──▶ deny label or ""
    envelope bytes       [ internal/gitguard.<verb>Verdict (×15)               ] ──▶ block? bool
    (in-memory           [ internal/gitguard.ParseEnvelope                     ] ──▶ command or ""
     fixtures)           [ internal/gitguard.DescribeClasses                   ] ──▶ comma-joined labels
                         [ internal/git.RefResolves / BranchExists             ] ──▶ bool (2s-bounded)
        ◀ tests attach here: table tests pin the tokenizer edges, prefix/wrapper resolution, each
          verdict function's option/carve-out matrix, envelope parse, and DescribeClasses == the deny
          table. Red before the package exists → does not compile.

### Seam diagram — degradation / manifest seam (shim rims + --describe)

    trigger: `bench guards` runs `bash block-dangerous-git.sh --describe`; or a Bash call with the core gone
        │
        ▼
    --describe (no stdin) ──▶ [ shim: name/boundary/why + denies from `bench guard-git --describe-classes` ]
                              [   wrapper/binary unresolvable → denies: manifest unavailable (analyzer missing) ] ──▶ exit 0
    envelope, binary gone  ─▶ [ shim rim 1: substring `git` → BLOCKED (fail closed) ; else allow ]            ──▶ git exit 2 / rest exit 0
    envelope, binary rc 3  ─▶ [ shim rim 2: rc ∉ {0,2,127} → BLOCKED (analyzer error) ]                       ──▶ exit 2
        ◀ tests attach here: gate-axi-contracts.sh (analyzer-missing / binary-missing / binary-errored
          fail-closed) and gate-axi-guards-contracts.sh (--describe manifest) — re-pointed from
          python3-missing to binary-missing. Red until the shim emits the new posture.

### Acceptance coverage map

Per-item granularity is stated where a behavior quantifies over a set (each
tokenizer edge, each verdict function, each degradation rim).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the full 79-case allow/block matrix passes through the ported shim → binary | hook enforcement | already covered (parity net): the matrix passes today against Python; it stays green across the port and turns red on any verdict divergence — not red-first by design | any tokenizer or verdict-function drift flips at least one of the 79 asserted verdicts |
| 1 | tokenizer edges — newline-as-separator, operator-token collapsing, malformed-quote fallback split, redirection stripping, quoted-newline survival into wrapper strings — each edge | go test unit | `internal/gitguard` tokenizer table test before `Tokenize` exists → does not compile | pins the subtlest code below the matrix, where the black-box cases can't reach every internal boundary |
| 1 | prefix resolution (env / command / nohup / timeout / xargs) and the one-level `sh\|bash\|zsh -c` wrapper scan — each prefix and the nested scan | go test unit | `internal/gitguard` prefix/wrapper table test before the functions exist → does not compile | a mis-ported prefix skip or an over-deep wrapper scan changes which token is classified as the verb |
| 1 | each of the 15 verdict functions' option/carve-out matrix (push, reset --hard, clean -f, branch -D/-f + worktree- carve-out, checkout, switch, restore, rebase, stash drop/clear, commit --amend, update-ref -d, tag -d, reflog expire, worktree remove --force + .claude/worktrees carve-out) — per verdict | go test unit | `internal/gitguard` verdict table tests before each function → does not compile | pins every carve-out and option boundary the matrix samples but does not exhaust |
| 1 | non-JSON stdin or an envelope without `tool_input.command` → allow (exit 0) | hook enforcement + go test unit | `ParseEnvelope` table test before it exists → no compile; a shim contract piping `{}` and non-JSON asserting exit 0 | preserves the current allow-on-unparseable behavior; catches a port that blocks or crashes on a malformed envelope |
| 3 | `bench guard-git --describe-classes` prints the comma-joined deny labels equal to the deny table | go test unit + degradation/manifest | `DescribeClasses` table test before it exists → no compile; the guards-aggregation contract asserts the denies clause via the shim | one source for enforce and advertise; a divergence between classify and describe surfaces here |
| 3 | a block exits 2 with `BLOCKED:` on stderr; an allow exits 0; a panic maps to exit 3 (never 2) | hook enforcement + go test unit | the matrix asserts exit 2 + `BLOCKED:`; a `internal/gitguard` table test drives the recover wrapper with a panicking classify stub and asserts the mapped exit is 3, not 2 (a recovered binary has no reliable crash *input*, so the wrapper is the honest seam) | catches a block that leaks as allow, or a panic that the shim would pass through as a disguised block |
| 4 | binary unresolvable → `--describe` reports `manifest unavailable (analyzer missing)`, exit 0 | degradation/manifest | the re-pointed analyzer-missing contract (`gate-axi-contracts.sh`): shim copied with no reachable bench → denies text asserted (text unchanged; the missing thing changes from `git-guard.py` to the binary) | catches a `--describe` that errors or advertises a manifest when the core is gone |
| 4 | binary missing → git-shaped input BLOCKED exit 2, non-git input allowed exit 0 | degradation/manifest | the re-pointed python3-missing → **binary-missing** contract (Handoff item 5, gate-visible): PATH/fixture with no bench → `git push` exit 2 + `BLOCKED:`, `ls -la` exit 0 | catches a shim that grants destructive authority (or breaks the whole shell) when the binary is absent — starts red until rim 1 exists |
| 4 | binary ran but errored (rc ∉ {0,2,127}) → BLOCKED exit 2 (fail closed) | degradation/manifest | **new** contract: a stub wrapper on PATH that exits 3 → assert `BLOCKED:` exit 2 | catches a shim that treats an analyzer crash as allow — starts red until rim 2 exists |
| 5 | the three fail-closed contracts and the `--describe` manifest contract assert the binary-world posture, no assertion weakened; the empty-`.py` case retires | degradation/manifest | the adapted contracts fail until the shim emits the binary-missing / binary-errored / analyzer-missing behavior above | proves the "cannot classify → refuse" property stays gate-locked while its trigger moves from Python to the binary |
| 7 | `git-guard.py`, `__pycache__`, gate 1c, and the `analyzer-parse-broken` canary are gone; no dangling reference in README, gate.sh, link/package contracts | gate load + docs stale-reference sweep | the docs sweep against a tree still naming `git-guard.py` → red; a gate load with a leftover 1c `.py` reference or an orphaned canary registry entry → red | the conformance layer fails when a deleted file is still referenced or the retired check still names a `.py` |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist and the
map's item-6 owners; each resolved as a coverage row above or a **Won't handle** line
here.

- **multi-line command blocks** — covered: tokenizer newline handling (matrix
  `:90–93` + table test); a destructive verb on any line blocks.
- **malformed quoting** — covered: the fallback-split path (table test) and the
  matrix's `bash -c` / wrapper cases.
- **worktree carve-out paths with spaces or `..` traversal** — covered: matrix
  `:74–79` plus path-normalization table tests; traversal out of `.claude/worktrees/`
  still blocks.
- **required tool / binary missing** — covered: story 4's rim-1 binary-missing
  contract (git-shaped fail-closed, rest usable).
- **non-JSON stdin / absent `tool_input.command`** — covered: the allow-on-
  unparseable row (envelope parse table test + shim contract).
- **hostile command text (wrappers, redirects, xargs, env prefixes)** — covered: the
  analyzer's own job, matrix `:100–111`; the threat model stays honest-mistake, one
  wrapper level deep, by design.
- **re-run idempotency / interrupted state** — safe by construction: the hook is
  stateless, classifying per call with no persisted state, so a re-run is identical;
  no test needed.
- **Won't handle: deliberate multi-level wrapper nesting** (`sh -c 'sh -c "git
  push"'`) — the wrapper scan goes exactly one level deep by design; the backstops
  for a misaligned agent are the git pre-push hook and pooled-worktree isolation, not
  this hook. Parity with today.
- **Won't handle: a ref/branch `git rev-parse` that hangs** — bounded at 2 seconds,
  then the fail-toward-blocking default applies; the bound is not made configurable.
- **Won't handle: a repository path containing a literal newline** — parity with the
  Python (which already misreads it); no in-scope caller produces one.
- **Won't handle: `#` as a shell comment** — a deliberate divergence from the Python
  tokenizer (shlex dropped `#`-to-end-of-line; the Go `lex` treats `#` as an ordinary
  word char). The divergence is fail-safe in both directions: a `#`-commented
  destructive verb over-blocks rather than slipping through, and the Python's own comment
  handling had the opposite hole (the comment ate the newline separator, so
  `git status # x`⏎`git push` slipped a real push past the guard). Honoring comments is
  unnecessary for an honest-mistake guard and would reintroduce that hole. Pinned by the
  `hash is a word` tokenizer table test.
- **Won't handle: the throwaway differential parity pass under the committed gate** —
  it requires `git-guard.py` alive to diff against and is deleted before commit; it
  is the map's #2 build-time insurance, run by the porter and discarded, not a
  permanent coverage row.

## Out of scope

- **Slice 5 — the remaining hooks** (`session-start.sh`, `stop.sh`,
  `check-agent-line.sh`, and the adapters' `_line-guard.sh` ported to Go behind
  their shims). A distinct capability with its own spec by the parent map's slice
  order; this slice only applies the shim pattern to the one hook whose Python it
  removes. Estimate to build later: per the parent map, one spec-sized session.

- **Evasion resistance / multi-level wrapper expansion** — an explicit non-goal of
  the guard's threat model (honest-mistake, one level deep), not deferred scope.
  There is no future spec for it; the pre-push hook and worktree isolation are the
  deliberate backstops.

- **Collapsing the shim into a pure in-process guard (no `.sh`)** — not possible:
  the harness invokes a hook by file path (`PreToolUse:Bash` → a script), so a thin
  shim is structural, not a deferrable cut.
