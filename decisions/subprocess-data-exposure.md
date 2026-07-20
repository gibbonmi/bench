# Minimal subprocess data exposure (FT88)

## Destination

Agents launch from a documented environment passlist and project gates keep
FT78's manifest-declared closed subject with a sentinel pinning it, prompt
text never travels through argv on a repository-controlled path, durable state
carries an objective identifier with sanitized rendering, and a data-handling
inventory plus sentinel contracts prove it. Sources: `RR:C-08`, `RC:H-01`.

## #1: What environment do agent and gate subprocesses receive?

Type: Grill

### Question
Today the shift adapter gets `os.Environ()+BENCH_SHIFT=1` and the gate gets
everything minus `BENCH_KIT`/`BENCH_WRAPPER`, so every credential in the parent
reaches repo-local scripts. Passlist vs denylist, where defaults live, and how a
repo opts in additions all change the build.

### Answer
Two static default passlists in the Go core, one per subprocess class.
Agent default: process basics (`PATH`, `HOME`, `USER`, `SHELL`, `TMPDIR`,
`TERM`, `COLORTERM`, `LANG`, `LC_*`, `XDG_*`) plus the `BENCH_*` variables the
loop itself owns plus a documented harness-auth section (the shipped adapters'
harness CLIs' own documented auth/config variables, enumerated in
`DATA_HANDLING.md`). Gate default: the same process basics plus common
build-tool families (`GO*`, `GIT_*`, `npm_config_*`, `NODE_*`, `CARGO_*`,
`RUSTUP_*`, `CI`) and the `BENCH_*` gate contract variables. Repo opt-ins live
in a committed `.bench/env.allow` with `[agent]` / `[gate]` sections, one name
or `PREFIX*` glob per line — explicit additions only; no wholesale
inherit-everything escape hatch. Malformed `env.allow` fails closed (refuse to
launch, name the line). Denylist posture rejected: unlisted-but-sensitive names
always slip a denylist.

## #2: Do the harness CLIs accept a stdin prompt?

Type: Research

### Question
Adapters currently pass the prompt as the harness CLI's positional argument; a
stdin transport only works end-to-end if each headless CLI reads a piped
prompt.

### Answer
Probed 2026-07-20. `claude -p`: yes — piped stdin prompt produced a normal
response (live probe, exit 0); `--input-format` is documented for `--print`.
`codex exec`: yes — help text states "If not provided as an argument (or if
`-` is used), instructions are read from stdin." `opencode run`: no documented
stdin contract — docs show only `opencode run [message..]` positional form
(https://opencode.ai/docs/cli/); not installed locally, so no live probe.

## #3: How does prompt text travel loop → adapter → harness?

Blocked by: #2
Type: Grill

### Question
The loop passes the full iteration prompt as the adapter's `$1`, and adapters
re-expose it as the harness CLI's argv — visible in process listings twice.

### Answer
Stdin end-to-end. The adapter contract changes from "prompt as `$1`" to
"prompt on stdin" (all three shipped adapters and the contract text in
`BENCH-reference.md` update together; no dual-mode transition — the kit ships
adapters and loop in lockstep). claude and codex adapters forward stdin to the
CLI's documented stdin path. opencode's CLI documents only positional prompts,
so its adapter keeps argv at the final hop with the residual exposure recorded
in `DATA_HANDLING.md` as a harness limitation to revisit when upstream
documents stdin. The mode-0600-file transport is the fallback design only if a
future harness can read neither stdin nor argv-free input; not built now.

## #4: What do durable records carry — objective text or an identifier?

Type: Grill

### Question
Objective text lands in commit subjects, intent records, stdout, and a 0644
`.bench-objective` worktree file. The row requires durable state to use an
objective identifier and sanitized rendering.

### Answer
Intake stays the choke point: `validateObjective` already rejects control
bytes; add a documented length cap. Shift commit subjects keep the sanitized
objective text — reviewer-authored, and readable Git history is a feature —
with the durability documented in `DATA_HANDLING.md`. Intent records carry the
existing entry key as the objective identifier and drop the free-text
objective field from durable storage; the full text lives only in the worktree
`.bench-objective` file, tightened to mode 0600 (recreated each shift, removed
with the worktree). Terminal summaries and structured output render through
one shared sanitizer that strips control sequences. No content-based secret
detection anywhere — sensitivity is handled by documenting which paths are
durable, not by guessing at content.

## #5: What do the sentinel contracts prove, and at which seam?

Type: Grill

### Question
The row demands proof that denied variables do not reach default subprocesses
and prompt content cannot leak into process listings, commits, or structured
output.

### Answer
Contract tests at the built-binary seam. Environment sentinel: export a marker
variable (e.g. `BENCH_TEST_SENTINEL=secret`) in the test parent, run a shift
with a stub adapter and a stub gate that each dump their environment; assert
the marker is absent from both and that passlisted/opted-in names survive.
Argv sentinel: the stub adapter reads `/proc/self/cmdline` and asserts the
prompt marker is absent, and stdin content matches the iteration prompt.
Durability sentinel: run a shift with a marker objective; assert intent
records store the key not the text, `.bench-objective` is mode 0600, and a
control-byte objective is rejected at intake before any record is written.

## #6: What does the data-handling inventory cover and where does it live?

Type: Grill

### Question
The row ships an inventory of every repository-controlled prompt, environment,
file, log, network, cache, and retention path.

### Answer
Root `DATA_HANDLING.md`, referenced from `SECURITY.md`, describing the current
decided state per invariant 3: each repository-controlled path (prompt
transport, both env passlists and the opt-in mechanism, worktree scratch
files, intent records, commit subjects, structured output, caches) with what
data reaches it and its retention. A conformance check asserts the passlist
constants and the inventory's variable listing derive from one source (the Go
core exports the lists; the doc or check reads them) so the advertisement
cannot drift from the enforcement. Whether consumers receive the file rides on
FT85's payload allowlist, not decided here.

## #7: Does the gate class survive FT78's closed gate subject?

Type: Grill (reviewer-closed 2026-07-20, in-session, after an implementation
stop-short)

### Question
Decision #1's gate premise was stale when the build reached it: FT78's verdict
identity already launches a project gate from a closed subject — `PATH` plus
only the names declared in `.bench/gate-inputs.json` — so the specced gate
passlist is strictly *wider* than what ships. The only live `gateEnv` caller is
the kit's own four-phase runner, and wiring the gate passlist there turned
Bench's own gate red (silent archive corruption in the conformance phase's
release-evidence probe, reproduced twice against a green baseline).

### Answer
Drop the gate class; `internal/env` is agent-class-only. The gate side of the
feature becomes a sentinel contract pinning FT78's closed subject: a marker
variable exported in the parent must not reach the gate subprocess, while
`PATH` and manifest-declared names survive. `.bench/env.allow` keeps `[agent]`
as its only known section — an unknown section, including a stale `[gate]`,
is rejected fail-closed — and the gate's opt-in mechanism is the existing
manifest declaration, documented in `DATA_HANDLING.md`. The kit's own
four-phase runner keeps its inherited environment, grouped with the other
subprocess classes FT87 owns. Rewiring the verdict path onto a passlist was
rejected: it reaches into verdict-identity hashing and ADR 0002's closed trust
posture for no added closure.

## Not yet specified

- Exact membership of the two default passlists beyond the families above —
  fixed at spec time against each harness CLI's and build tool's documentation.
- Objective length-cap value.

## Out of scope

- Env minimization for other subprocess classes (git fetch, model discovery,
  hooks, canary) — FT87 owns bounded subprocess behavior; revisit after both.
- Content-based secret detection/redaction — documented-durability policy only.
- Consumer payload membership for `DATA_HANDLING.md` — FT85 owns the allowlist.
- Host-level controls (OS sandboxing, IAM, SIEM) — outside the repo roadmap.

## Handoff

1. **Module boundaries.** `internal/env` (new): the agent passlist constants,
   `env.allow` parsing, env construction for the agent class. `internal/shift`:
   adapter stdin transport, objective file perms, intent-key identifier.
   `internal/gate`: unchanged — the project-gate environment is FT78's
   manifest-declared subject, pinned by a sentinel rather than rebuilt.
   `.bench/adapters/*`: stdin forwarding. `internal/sanitize` (new or folded
   into an existing owner): the one shared control-sequence sanitizer. Root
   `DATA_HANDLING.md`.
2. **Contracts.** `internal/env`: (repo root) → ordered agent env slice or a
   fail-closed error naming the malformed `env.allow` line. Adapter contract:
   prompt on stdin, `BENCH_SHIFT=1`, exit code passthrough. Intent entries:
   key-only objective reference. `validateObjective`: empty/control-byte/
   over-cap → exit-2 usage error.
3. **Deep vs thin.** `internal/env` is the deep module (policy + parsing +
   construction behind one call). Adapters stay thin pass-throughs — their
   assertable is the harness argv/stdin shape, not logic. The sanitizer is
   deep enough to own every rendering call site.
4. **Black-box assertables.** Stub adapter env dump (marker absent, passlist
   present); stub gate env dump (marker absent, `PATH` and manifest-declared
   names present — the FT78 closed-subject pin); `/proc/self/cmdline`
   prompt-marker absence; stdin content equality; `.bench-objective` mode
   bits; intent JSON field shape; exit codes and stderr for malformed
   `env.allow` and rejected objectives.
5. **Gate attachment.** All sentinel contracts run as contract tests at the
   built-binary seam inside the existing gate contract phase; the
   inventory-vs-constants check runs in the conformance phase. No seam is
   gate-invisible; no manual verify.
6. **Hostile-input owners.** Control-byte/oversized objectives →
   `validateObjective`. Malformed/hostile `env.allow` (bad glob, traversal,
   wrong section) → `internal/env` fail-closed parse. Control sequences in
   render paths → shared sanitizer. Hostile env values (huge, multiline) pass
   through untouched by design — passlists filter names, never values.
7. **Uncertainty flags.** opencode stdin support is undocumented — if upstream
   adds it before build, drop the argv residual; the spec should cite the doc
   state at spec time. (The gate-passlist compatibility flag is resolved by
   #7: the gate keeps FT78's subject, and a repo whose gate needs a variable
   declares it in `.bench/gate-inputs.json` — `DATA_HANDLING.md` states that
   remedy.)
8. **Rejected alternatives.** Denylist posture; wholesale inherit-all escape
   hatch; dual-mode argv+stdin adapter transition; objective-ID-only commit
   subjects; content-based secret scanning; building the 0600-file transport
   now.
9. **Domain watch-outs.** The passlist filters names only — a passlisted
   variable's value may still be sensitive; that is a documented-durability
   fact, not a defect. `LC_*`-style glob families are load-bearing on real
   systems; exact-name-only matching breaks locales.

Dependency order: `internal/env` passlists + sentinel env contracts →
stdin transport + adapter contract change → objective identifier/perms/
sanitizer → `DATA_HANDLING.md` + conformance derivation check.
