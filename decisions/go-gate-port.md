# go-gate-port — slice 8 of the Go rewrite

Child of `decisions/go-rewrite.md` (#6, slice order): the gate content itself, the
last slice — after it the strangler window closes. Bootstrap evidence: the gate is
`.bench/gate.sh` (231 lines of inline checks) + 12 contract fragments (~3.4k lines)
+ `.bench/lib/canary-run.sh` (63 lines, **shipped** — `bench init`/`link` install
it and consumer-scaffolded gates source it) + ~60 fixtures in `tests/canary/`. The
check population splits two ways: **behavior contracts** (runtime, shift, link,
doctor, axi×3, go routing, package — self-provision temp repos, exec the CLI, and
guard themselves out of canary fixtures) and **root-grading conformance checks**
(the inline gate.sh checks plus the docs/line/coverage/anchor fragments — greps
over the repo at `$root`; these are what the fixtures target). Parent constraints
carried: `.bench/gate.sh` stays the stable entry and exit-code contract; the
canary layer is this slice's net; every ported check must still fail its fixture
with its targeted error before its shell twin retires.

## #1: What invokes the ported checks?

Type: Grill

### Answer
**Pure `go test` — no kit-conformance subcommand in the binary.** Behavior
contracts become Go integration tests (exec the built binary in `t.TempDir`
fixtures — the `contract` harness's provision/report/cleanup becomes a test
helper, carrying `--space-path`/`--no-repo` capability). Root-grading conformance
checks become root-parameterized tests (graded root via env, default the kit
repo) so the canary sweep can point them at a fixture tree. `.bench/gate.sh`
shrinks to a few-line entry: resolve root, build the dev binary, `go vet`,
`go test -count=1` (the gate never trusts the test cache — contract tests exec
git/bash/binaries, which escape Go's input tracking), canary sweep. Rejected: a
`bench selfcheck` subcommand (kit-conformance logic is dead weight in the shipped
binary) and a hybrid split (two runners for one population).

## #2: Does the canary runner port too?

Type: Grill

### Answer
**Yes — `bench canary` in the binary.** Reviewer's call, weighing the tradeoff
table: one language for the whole gate stack and deterministic Go testing for
anything coded outweigh keeping the 63-line runner on the consumer-readable text
surface. The subcommand stays generic — copy each `tests/canary/<fixture>/files`
tree to a throwaway repo (restoring `dot-` prefixes), exec **the repo's own gate**
(`.bench/gate.sh`) against it under `BENCH_CANARY_INNER=1`, assert red plus the
EXPECT substring, plus the absent-or-empty-harness red and the vacuous-EXPECT
baseline — no kit-specific logic inside, so kit and consumers run the same sweep.
Rejected: staying shell (forgoes the testability and one-language wins the
reviewer prioritized).

## #3: How do existing consumer gates survive the runner moving?

Type: Grill

### Answer
**`.bench/lib/canary-run.sh` becomes a one-glance shim** delegating to
`bench canary` and folding the exit into `fail` — consistent with the permanent
exec rim (launcher, hook shims). Existing consumer gates keep sourcing it and keep
working after a re-link; new scaffolds (`internal/adopt/init.go`) call
`bench canary` directly. Rejected: removing the lib (every linked repo's gate goes
red until a human edits it) and keeping sourcing as the permanent API (permanent
indirection for new consumers).

## #4: How does the slice cut into specs?

Type: Grill

### Answer
**Three specs, a→b→c.** (a) Canary machinery — `bench canary`, the lib shim, the
scaffold update, gate.sh's inner-run mode, and the root-parameterized conformance
harness: the retire-rule machinery the rest depends on. (b) Conformance checks —
ported check-by-check, each fixture proving the Go check bites before its shell
twin retires. (c) Behavior contracts — the bulk, family-by-family as stories in
one spec. Each is a session-sized diff with the gate green throughout. Rejected:
a+b combined (bigger first diff) and one staged spec (multi-session, riskiest
diff of the migration in one review).

## Handoff

1. **Module boundaries.** The canary sweep is a new `internal/` package behind a
   `bench canary` subcommand (naming is the spec's call, `gitguard` precedent).
   Conformance checks and behavior contracts live as `go test` code in the kit
   module — never in the shipped binary. The `contract` harness becomes a shared
   test helper package. Deleted by the end of (c): every `gate-*.sh` fragment and
   `gate-contract-runner.sh`; `canary-run.sh` shrinks to the shim (#3);
   `.bench/gate.sh` shrinks to the few-line entry (#1). Permanently text: the
   exec rim (npm launcher, postinstall, generated pre-push, hook entry shims,
   doctor PATH shim, planted `.bench/bin` shims, consumer gate/done files).
2. **Contracts.** Gate exit 0/non-zero at `.bench/gate.sh` (exit 3 outside a git
   repo) and its `gate: green`/`gate: red` lines carry. Fixture format carries
   unchanged: `files/` tree, `EXPECT` substring file, `dot-` prefix restore,
   absent-or-empty harness is red, vacuous-EXPECT baseline. `BENCH_CANARY_INNER`
   carries: the sweep never recurses, and the inner gate runs only the
   root-grading subset. New: `bench canary` exits non-zero when any fixture fails
   to bite, attributing per fixture by name; conformance tests accept a graded
   root via env; the gate runs `go test -count=1`.
3. **Deep vs thin.** The binary + the kit's test code are the deep units (sweep
   logic, all check logic, the fixture-provision helper). `.bench/gate.sh` and
   the `canary-run.sh` shim are thin pass-throughs — small enough to read in one
   glance.
4. **Black-box assertables.** Gate exit/stdout at the entry, before and after
   each flip. The ~60 fixtures are the bite-proof for every conformance check —
   the retire rule is per check: fixture red with the targeted substring under
   the Go check, then the shell twin deletes in the same diff. `bench canary`
   asserts standalone: a repo with a biting fixture set exits 0, a rotted or
   absent one exits non-zero naming the fixture.
5. **Gate attachment.** For specs (a) and (b) the not-yet-ported shell gate
   remains the oracle while its twins port. Gate-blind: spec (c)'s behavior
   contracts have no canary fixtures — the flip is same-diff and review-graded
   (assertion-for-assertion parity against the deleted fragment), flagged for
   the review phase rather than the gate.
6. **Hostile-input owners.** Paths with spaces → the test helper carries
   `--space-path` provisioning. Absent vs empty → the harness-absent-or-empty
   red and the empty-EXPECT/vacuous baseline carry. Missing tool → the go.mod
   hard-red toolchain check carries in the entry. No-trailing-newline
   hand-edited files → EXPECT reading and the grep-based conformance parsers.
   Symlink/dot-dir fixtures → `dot-` restore in the sweep. Re-run idempotency →
   sweep temp dirs created and removed per fixture. SIGINT mid-sweep → throwaway
   dirs only, no repo state touched.
7. **Uncertainty flags.** Per-fixture inner-run cost: ~60 inner gate runs each
   invoking `go test` machinery — spec (a) must measure and keep the sweep at or
   under today's wall-clock (options: `-run` filter, a prebuilt `go test -c`
   conformance binary), escalating per `craft-line` if the mechanism is unclear.
   Whether any consumer of the `gate: green`/`gate: red` stdout lines exists
   beyond humans (the stop hook reads the gate cache, not stdout) — spec (a)
   verifies before treating the lines as droppable surface. Cache posture for
   *pure* unit tests (module `go test` split from `-count=1` contract runs) is
   the spec's call within the #1 constraint.
8. **Rejected alternatives.** `bench selfcheck` conformance subcommand and the
   hybrid runner (#1); canary runner staying shell (#2); removing the lib or
   keeping sourcing as the permanent API (#3); two-spec and one-spec cuts (#4);
   trusting the go test cache in the gate (#1).
9. **Domain watch-outs.** EXPECT substrings couple fixtures to check message
   text: a ported check keeps its targeted substring or updates the fixture in
   the same diff — the vacuous-EXPECT baseline guards against emptying one. A
   check exists in shell or Go, never both: the shell twin retires in the same
   diff its fixture bites the Go check. The sweep execs the repo's own gate —
   kit and consumer follow one path; kit-specific knowledge never enters the
   subcommand. The gate must not trust `go test`'s cache for anything that
   execs a subprocess.

Dependency order: three specs (a) canary machinery → (b) conformance checks →
(c) behavior contracts; (c) closes the parent map's slice order and the
strangler window.
