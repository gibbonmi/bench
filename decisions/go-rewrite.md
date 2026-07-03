# Go rewrite — move the kit's executable logic from shell to Go

Graduated from the roadmap 2026-07-03. The parked line weighed two paths —
swap grep→rg in the scripts, or rewrite the shell in Go. The bootstrap grill
rejected the rg swap (#1) and made this map the Go-rewrite decision. Evidence
gathered at bootstrap: 30 shell files / ~5.7k LOC plus a 418-line Python hook
analyzer; a full gate run costs 82s wall / 49s CPU, paid per shift iteration.

What Go would buy (assessed at bootstrap, packaging excluded): a unit-test
seam for the parsers/emitters the canary layer currently compensates for; the
shell-footgun classes of the hostile-input checklist dissolve structurally;
one-source-per-fact becomes enforceable across the CLI/hook boundary (the
gate_tree_hash and hooks_dir mirrors are live defects); one language absorbs
git-guard.py; gate loop ~2–4× faster (fork/exec churn vanishes, git subprocess
time remains); cheap runway for the five parked parser candidates. The main
non-packaging counterweight: agents currently read and hot-patch the kit as
plain text — Go inserts a compile step and makes the core opaque.

## #1: Is the grep→rg swap worth doing?

Type: Grill

### Question
The parked idea's first path: replace grep with ripgrep in the kit scripts
for speed and consistency.

### Answer
**Rejected.** Every one of the ~100+ grep call sites is a single-file
assertion (`grep -qF needle file`) or a pipe filter — zero recursive tree
searches, the only place ripgrep is faster, so the speed win is nil. And rg
would be the kit's first hard non-POSIX dependency, shipped to every linked
repo (the gate deliberately keeps even shellcheck best-effort to avoid this).
Do not reopen for consistency's sake; reopen only if recursive search
appears in a hot path.

## #2: What is candidate for rewrite, and what stays plain text regardless?

Type: Grill

### Question
Scope boundary of the language question itself.

### Answer
Candidate: executable logic only — `bin/*.sh`, the gate fragments, `.bench`
lib and hooks, and `git-guard.py` (~6.1k LOC across two languages). All
markdown content (skills, commands, BENCH.md, docs, profiles) and the JSON
adapters stay plain text under every outcome — the content surface is the
product's portable half and is never in question.

## #3: Does the kit accept a compiled core?

Type: Grill

### Question
The product-identity call, and the one that can kill the whole idea: today
"plain files are the product" — agents read the scripts as source, patch them,
and the edit takes effect immediately; the gate's parse layer (`bash -n`) and
canaries assume text. A Go core makes the executable half opaque to agents,
inserts compile between edit and effect, and changes what "reading the kit"
means for every future session. Is that identity change acceptable in
principle? If no, the map closes as rejected and #4–#6 die.

### Answer
**Accepted in principle, conditional on #4.** The identity change is
acceptable because the legibility loss is smaller than the slogan implies:
Go source is at least as agent-legible as bash, and `go run` keeps
edit→effect near-live in the kit repo. The condition is consumer-side and
hard: #4 must find a distribution shape with **no toolchain requirement on
consumer machines** and **an auditable surface** (e.g. thin shell shims
exec'ing a versioned binary — consumers can still read what their hooks
invoke). If #4 cannot deliver both, this map closes as rejected.

## #4: How would consumers get the binary, and what does the toolchain cost?

Type: Research

### Question
The packaging reality deferred at bootstrap. npm currently ships shell that
`bench link` copies into consumer repos and hooks exec by path. For a Go
core: prebuilt per-platform binaries vs `go install` vs source build;
whether the Go toolchain becomes a consumer or CI dependency; what the
pre-push hook execs on a machine with no toolchain; what `bench link` copies
and how binary/asset version skew is detected; what the gate's parse layer
becomes (go build/vet/test in place of `bash -n`). Acceptance bar (set by
#3): no toolchain requirement on consumer machines, and an auditable
consumer surface. A shape that misses either bar closes the map as
rejected. Output: a short summary asset with a recommended distribution
shape and its hard dependencies.

### Answer
**Bar met — recommended shape: npm platform packages (esbuild pattern).**
Prebuilt Go binaries as os/cpu-filtered `optionalDependencies`
(`@benchkit/<os>-<arch>`, four targets), `bin/bench` a thin launcher that
execs the matching one. No consumer toolchain — node ≥ 22 is already
required today. Auditable surface holds: the pre-push hook is already
self-contained inline shell (no binary dependency, survives unchanged);
`bench link` plants thin exec shims in `.bench/bin/` instead of full CLI
source; harness hook entries stay `.sh` shims calling binary subcommands —
every line executing as text in a consumer repo stays readable. Version
skew: kit-version stamp in the link manifest, checked against
`bench --version` by session-start/doctor. Kit repo: Go toolchain is
dev/CI-only; the gate's parse layer gains `go build`/`vet`/`test` beside
`bash -n` for the remaining shell. Rejected shapes (postinstall download,
`go install`, committed binaries) and residual risks for #5 (npm
optional-deps lockfile edge, `--ignore-scripts`, 4-target release matrix):
`decisions/assets/go-distribution.md`.

## #5: Go or no-go?

Blocked by: #3, #4
Type: Grill

### Question
The decision this map exists for: do the bootstrap-assessed benefits
(testability, footgun elimination, unification, gate speed) justify the
identity change (#3) at the distribution cost (#4)? A ~6k-LOC rewrite with
the black-box gate as the only regression net is the migration risk to
price in.

### Answer
**Go — strangler-only; big-bang is off the table.** The rewrite is
committed because the kit's recurring defect classes (mirrored facts across
script boundaries, pipefail-dependent fallbacks, the Python outgrowth) are
language-caused, the canary layer is a standing tax on shell's
untestability, and the kit's growth trend (parser candidates) is shell's
worst fit — while the migration risk is bounded by the black-box gate: the
AXI contracts and gate checks assert stdout and exit codes, so they keep
biting across a port. Typed, tested Go is also more delegable to the cheap
tier than bash. Accepted costs, eyes open: ~10–15 spec-sized sessions of
port work, and the permanent 5-package release matrix from #4. Constraint
binding on #6: every port lands seam-by-seam behind the existing CLI
dispatch with the old gate green throughout — no flag-day cutover.
Rejected: no-go with a shell unit harness (half the testability win,
nothing for footguns/duplication/speed); pilot-then-re-decide (pays the
release-matrix cost for a fraction of the value while deferring the
decision this map exists to make).

## #6: Rewrite scope, migration order, and testing shape

Blocked by: #5
Type: Grill

### Question
Only if #5 is go: big-bang vs strangler (seam-by-seam behind the existing
CLI dispatch); which seam first (the AXI query surface and parsers are the
highest-value, lowest-risk candidates; the gate itself migrates last, under
the old gate's watch); what the unit layer covers vs what stays black-box
gate contract; what the canary layer becomes.

### Answer
**Strangler in eight slices, value-first and risk-ascending; the gate ports
too, last.** Order: (1) walking skeleton — Go module, CI cross-compile
matrix, launcher, `bench version` ported to prove #4's pipeline before real
logic moves; (2) AXI query surface + TOON emitter; (3) `bench status`
renderer; (4) `git-guard.py` absorbed as a Go package (`--describe` is the
contract); (5) hook logic → binary subcommands behind `.sh` shims; (6)
`doctor` + `link` — the highest-stakes mutators, ported once Go idioms are
settled; (7) worktree + shift loop; (8) gate fragments → `go test`, with
`.bench/gate.sh` remaining the stable entry point and exit-code contract.
Risk posture for slice 8 (reviewer's call, cost flagged): the unchanged
shell gate is the regression net for slices 1–7 only; the gate port's own
net is the canary layer — red-by-construction fixtures are
language-agnostic, and every ported check must still fail its fixture with
its targeted error before the shell check retires. Per-slice rule: a seam
ports only when black-box contracts cover it first — the shift/worktree
loop needs contract backfill before slice 7. Done per slice = gate green +
`go build`/`vet`/`test` green. Permanently shell (from #4): postinstall,
generated pre-push, hook entry shims, the doctor PATH shim, the launcher,
and the planted `.bench/bin` shims.

## Handoff

1. **Module boundaries.** One Go module, one binary: `cmd/bench` +
   `internal/<seam>` packages mirroring the slice list in #6 (query/TOON,
   status, guard-analyzer, hooks, doctor, link, worktree/shift, gate
   checks). Outside the module, permanently shell: the npm launcher,
   postinstall, generated pre-push, hook entry shims, doctor PATH shim,
   planted `.bench/bin` shims, and `.bench/gate.sh` as the gate's stable
   entry. Distribution: `benchkit` wrapper + four `@benchkit/<os>-<arch>`
   platform packages (#4).
2. **Contracts.** All existing observable contracts carry over unchanged —
   that is the strangler's premise: CLI subcommand names and exit codes;
   the AXI hybrid contract (TOON stdout, definitive empty states,
   structured errors on stdout, exit 0/1/2); each guard's `--describe`
   self-manifest; link-manifest semantics and safe-link refusal behaviors;
   the generated pre-push text; gate exit 0/non-zero at `.bench/gate.sh`.
   New: the launcher execs the platform binary or fails loudly naming the
   missing package; kit-version stamp in the link manifest checked by
   session-start/doctor.
3. **Deep vs thin.** The binary is the deep unit — all logic, parsing, and
   platform variance live behind its subcommands. Launcher, shims,
   postinstall, and the gate entry are thin pass-throughs with no seam of
   their own; they must stay small enough to read in one glance.
4. **Black-box assertables.** Per slice: identical stdout/exit assertions
   before and after the port — the existing gate fragments (AXI, link,
   runtime, status contracts) run unchanged against the ported binary.
   New unit layer: `go test` table tests on parsers/emitters/pure logic.
   Slice 8: every ported gate check still turns its canary fixture red
   with the targeted error substring.
5. **Gate attachment.** The unchanged shell gate is the oracle for slices
   1–7; it gains `go build`/`vet`/`test` in its parse layer once the module
   exists. Slice 8's net is the canary layer. Gate-blind spots: a real
   multi-package `npm i -g` (manual smoke per release), and the
   shift/worktree loop until its contract backfill lands (must precede
   slice 7).
6. **Hostile-input owners.** Most checklist classes dissolve in Go; what
   remains maps to the shell rim and the port seams — paths with
   spaces/globs → shim quoting and Go exec argv (never a shell string);
   absent vs empty file → parser package table tests; required tool
   missing → the launcher's missing-platform-package error; symlink
   invocation → launcher target resolution; SIGINT mid-loop → slice 7
   (worktree/shift) scratch-state cleanup; re-run idempotency → slice 6
   (relink) and slice 1 (repeat install); cwd deeper than root → root
   resolution in the binary, asserted per subcommand; no-trailing-newline
   hand-edited files → the maps/learnings/roadmap parsers.
7. **Uncertainty flags.** Shift/worktree contract coverage is thin —
   spec-writer scopes the backfill before slice 7, escalating per
   `craft-line` if the loop's observable surface is unclear. npm
   optional-deps lockfile edge (asset, #4) — verify launcher repair
   behavior against current npm during slice 1. Go toolchain version
   pinning policy (go.mod directive vs CI matrix) — settle in slice 1.
8. **Rejected alternatives.** grep→rg swap (#1); no-go with a shell unit
   harness and pilot-then-re-decide (#5); big-bang cutover (#5); keeping
   the gate shell permanently (#6 — reviewer chose to port it,
   canary-netted); postinstall-download as primary, `go install`, and
   committed binaries (#4/asset).
9. **Domain watch-outs.** During the strangler window shell and Go coexist
   — the CLI dispatch routes per-subcommand and a half-ported seam must
   never have two live implementations. npm installs only the matching
   platform package; a lockfile written on one OS can omit another's
   optional dep (repair fallback exists; `--ignore-scripts` consumers lose
   it). An orphaned or stale binary after version moves is detected by the
   manifest version stamp, not by the binary itself. Go builds are
   reproducible only with `-trimpath` and pinned toolchains.

Dependency order: slices 1→8 as numbered in #6; contract backfill for the
shift/worktree loop lands before slice 7; slices 2–3 are independent of
each other after slice 1 and may swap if a spec needs to.
