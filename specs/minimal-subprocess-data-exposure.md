# Minimal subprocess data exposure

Status: implemented

Compiled from `decisions/subprocess-data-exposure.md` (closed, all seven tickets
resolved; #7 is the reviewer-closed 2026-07-20 amendment dropping the gate
class). Roadmap row FT88; sources `RR:C-08`, `RC:H-01`.

## Problem

The shift adapter inherits the whole parent environment — `os.Environ()` plus
`BENCH_SHIFT=1` — and it is a repository-controlled script, so any credential
the reviewer happens to have exported — a cloud token, a registry password, an
unrelated project's API key — reaches code the reviewer did not write for the
purpose. The project gate is the half FT78 already closed: `bench gate`
launches the gate script with `PATH` plus only the names declared in
`.bench/gate-inputs.json`. But nothing gate-observable pins that closure, so a
future change widening the gate subject would go green silently.

Prompt text has the same shape of problem in a different channel. The loop hands
the full iteration prompt to the adapter as `argv[1]`, and each shipped adapter
re-exposes it as the harness CLI's own positional argument, so the prompt is
visible twice in any process listing on the machine. Objective text then lands
in durable state nobody scoped for it: a world-readable `.bench-objective` file,
a required free-text field in the intent ledger, and terminal renders with no
shared control-sequence policy behind them.

Nothing today proves any of this either way. There is no inventory of which
repository-controlled path carries what data, and no contract that would go red
if a future change widened the exposure again.

## Solution

The harness adapter launches from a documented passlist instead of the
inherited environment, with a committed `.bench/env.allow` as the only way to
add to it; the project gate keeps FT78's manifest-declared closed subject, and
a repo whose gate needs an extra variable declares it in
`.bench/gate-inputs.json`. The iteration prompt travels on stdin end-to-end, so it never appears
in argv on any hop Bench controls. Durable state references an objective by its
intent key; the full text lives only in a mode-0600 worktree file that dies with
the worktree. One sanitizer owns the control-sequence policy for every terminal
render. A root `DATA_HANDLING.md` inventories every repository-controlled path,
and a conformance check derives its variable listing from the same constants the
enforcement uses, so the advertisement cannot drift from the enforcement.

Sentinel contracts at the built-binary seam prove each claim: a marker variable
exported in the test parent reaches neither subprocess, a marker prompt appears
in the stub adapter's stdin but not its `/proc/self/cmdline`, and a marker
objective reaches the ledger only as a key.

## User stories

1. As a reviewer running a shift, I want the harness adapter to receive only a
   documented passlist of environment variables, so that an unrelated credential
   in my shell does not reach repository-controlled adapter code.
   Line: `gpt-5.6-luna` / medium. The passlist membership and the construction
   contract are both fixed by the decision map, so this is transcription at a
   known seam rather than a design question.

2. As a reviewer running the gate, I want a contract test proving a marker
   variable exported in my shell does not reach the gate subprocess while
   `PATH` and manifest-declared names survive, so that FT78's closed gate
   subject cannot be silently widened by a future change. (Amended per map #7:
   the original gate-class passlist re-solved a problem FT78 had closed.)
   Line: `gpt-5.6-terra` / medium. This is oracle authorship over
   already-shipped behavior, and the row's whole value hangs on a demonstrated
   mutation red rather than a TDD red.

3. As a reviewer whose repo's adapter needs one extra variable, I want to add
   it to a committed `.bench/env.allow` under an `[agent]` section, so that the
   addition is explicit and reviewable in the diff. A gate-side need has an
   existing home instead: declare the name in `.bench/gate-inputs.json`, which
   is both the opt-in and the verdict-identity declaration.
   Line: `gpt-5.6-luna` / medium. A small line-oriented parser with a fully
   enumerated grammar and a stated fail posture.

4. As a reviewer, I want a malformed or hostile `.bench/env.allow` to refuse to
   launch and name the offending line, so that a typo silently widens nothing
   and a wholesale wildcard cannot be smuggled in.
   Line: `gpt-5.6-luna` / medium. The fail-closed posture is decided; the work is
   enumerating the rejection cases the map and the profile checklist already name.

5. As a reviewer, I want a contract test that exports a marker variable and
   proves it reaches neither the adapter nor the gate while passlisted and
   opted-in names survive, so that a future change widening the environment turns
   the gate red.
   Line: `gpt-5.6-terra` / medium. This is oracle authorship — the profile routes
   gate and conformance logic to mid effort because a check that fails to bite is
   worse than no check.

6. As a reviewer, I want the loop to hand the iteration prompt to the adapter on
   stdin rather than as an argument, so that the prompt does not appear in the
   machine's process listing.
   Line: `gpt-5.6-luna` / medium. One call site with a decided transport and an
   existing stdin plumbing path on the command.

7. As a reviewer using any shipped adapter, I want it to forward stdin to its
   harness CLI's documented stdin path, so that the prompt stays out of argv on
   the final hop too.
   Line: `gpt-5.6-luna` / medium. Three near-identical POSIX shell edits against
   CLI contracts probed and recorded in the map.

8. As an agent reading the adapter contract, I want `.bench/BENCH-reference.md`
   to describe the stdin transport as the only contract, so that a custom
   adapter written against the docs is not written against the removed argv form.
   Line: `gpt-5.6-sol` / high. This is the leverage override in `craft-line`: the
   reference prose is loaded by sessions in every linked repo, so a wrong
   sentence here compounds far past the cost of writing it.

9. As a reviewer, I want a contract test asserting the stub adapter's
   `/proc/self/cmdline` carries no prompt marker while its stdin matches the
   iteration prompt exactly, so that a regression back to argv is caught by the
   gate rather than by a process listing.
   Line: `gpt-5.6-terra` / medium. Oracle authorship again, and the assertion has
   to distinguish a genuinely absent marker from a test that never ran the adapter.

10. As a reviewer, I want an over-long objective rejected at intake with a usage
    error naming the cap, so that unbounded text cannot flow into a commit
    subject, a scratch file, and the ledger at once.
    Line: `gpt-5.6-luna` / low. One added branch beside the existing empty and
    control-byte checks in the same function.

11. As a reviewer, I want the worktree `.bench-objective` file created mode 0600,
    so that the one place the full objective text persists is readable only by
    the user who started the shift.
    Line: `gpt-5.6-luna` / low. A single file-mode constant with an existing
    contract-test fixture to assert against.

12. As a reviewer, I want intent records to reference an objective by their
    existing entry key instead of storing its text, so that the durable ledger
    carries an identifier rather than free-form content.
    Line: `gpt-5.6-luna` / medium. The field removal is mechanical, but it crosses
    the ledger's validity rule and the status renderer, so the call sites need
    reading before editing.

13. As a reviewer reading `bench status`, I want an interrupted shift's row to
    still show me which objective it was, so that dropping the ledger field does
    not cost me the ambient signal I use to decide whether to resume.
    Line: `gpt-5.6-luna` / medium. Reading the objective back from the recorded
    worktree path is a small change, but it has to degrade cleanly when the file
    or the worktree is gone.

14. As a reviewer, I want every terminal render of operator-influenced text to go
    through one sanitizer, so that the control-sequence policy has a single source
    instead of three near-copies that drift.
    Line: `gpt-5.6-terra` / medium. This is a collapse across three packages with
    byte-pinned tests on one of them, so the risk is in the call sites rather than
    in the policy.

15. As a reviewer, I want a contract test proving a marker objective reaches the
    ledger only as a key, that `.bench-objective` is mode 0600, and that a
    control-byte objective is rejected before any record is written, so that all
    three durability claims are gate-observable.
    Line: `gpt-5.6-terra` / medium. Oracle authorship, and the ordering assertion
    — rejected *before* any write — is the part a careless test would miss.

16. As a reviewer or an auditor, I want a root `DATA_HANDLING.md` describing every
    repository-controlled prompt, environment, file, log, network, cache, and
    retention path with what data reaches it, so that the exposure surface is
    documented rather than reconstructed from code.
    Line: `gpt-5.6-terra` / high. This deliberately does **not** take the
    profile's top-tier doc-authoring row: that override is justified by prose
    loaded into every session, and this file is repo-only reference material an
    auditor reads on demand. High effort, mid model.

17. As a reviewer, I want a conformance check asserting the inventory's variable
    listing derives from the passlist constants, so that adding a variable to the
    code without documenting it turns the gate red.
    Line: `gpt-5.6-terra` / medium. Gate and conformance logic per the profile's
    cached routing.

18. As a reviewer, I want a canary fixture proving that derivation check actually
    bites, so that the check cannot rot into an always-pass without the meta layer
    noticing.
    Line: `gpt-5.6-terra` / medium. The conformance-family canary rule already
    requires it, and the fixture has to fail for its own targeted reason.

## Implementation decisions

**New package `internal/env` — the deep module.** Policy, parsing, and
construction sit behind one call: `(repoRoot) → []string` or a fail-closed
error naming the malformed `.bench/env.allow` line. The package serves the
agent class only — the project gate's environment is FT78's manifest-declared
subject, owned by `internal/gate`, which this feature pins rather than
rebuilds. The package exports the default passlist as data so the conformance
check in story 17 can read the same values the enforcement uses — this is the
single-source requirement, not a convenience.

**Default passlist.** A process-basics set: `PATH`, `HOME`,
`USER`, `LOGNAME`, `SHELL`, `TMPDIR`, `TERM`, `COLORTERM`, `LANG`, `LC_*`,
`XDG_*`. `HOME` and `XDG_CONFIG_HOME` are load-bearing for git's own config
resolution, and the `LC_*` family must be a glob because exact-name matching
breaks locales on real systems.

To that it adds `BENCH_*` and the shipped adapters' documented harness
variables: `ANTHROPIC_*`, `CLAUDE_CODE_*`, `CLAUDE_CONFIG_DIR`, `API_TIMEOUT_MS`
(Claude Code); `CODEX_*`, `RUST_LOG`, `SSL_CERT_FILE` (Codex); `OPENCODE_*` and
`OPENAI_API_KEY` (opencode's documented provider substitution); and the cloud
credential chains Claude Code documents for Bedrock and Vertex routing — `AWS_*`,
`GOOGLE_*`, `GCLOUD_PROJECT`, `CLOUD_ML_REGION`, `VERTEX_LOCATION`. Every name
here is cited to an official page in `DATA_HANDLING.md`; provider keys for
harnesses Bench does not ship an adapter for (Groq, Gemini, Azure, Cloudflare)
are deliberately absent — they are a one-line `env.allow` addition, which is what
the opt-in mechanism exists for.

**A default glob must not straddle two families.** Every glob above covers a
namespace a single owner controls; a name whose prefix is shared with a foreign
family is enumerated exactly instead. This is the rule that kept `GO*` out of
the retired gate-class draft — it would also have matched
`GOOGLE_APPLICATION_CREDENTIALS` — and it binds every future addition to the
default list. A unit-test edge row checks each default glob against a fixture
of foreign names, so an exact set cannot later be "simplified" into a wider
glob.

**The project gate keeps FT78's closed subject, and this feature pins it.**
`bench gate` launches the gate script with `PATH` plus only the names declared
under `environment` in `.bench/gate-inputs.json` — already strictly tighter
than any default passlist, and part of the verdict identity. What is missing is
the regression guard: a sentinel contract proving a marker exported in the
parent never reaches the gate script, while `PATH` and a manifest-declared name
survive. The sentinel is green on the shipped tree, so its red capability is
proven once by a targeted mutation appending the marker to the subject
environment, with that demonstrated red recorded in the test. A repo whose gate
script needs an extra variable declares it in `.bench/gate-inputs.json`;
`DATA_HANDLING.md` states that remedy alongside the observable failure — the
gate red with the project's own error naming the missing variable. The kit's
own four-phase runner keeps its inherited environment: its phases are Bench's
own repository-controlled commands, and minimizing them is grouped with the
other subprocess classes FT87 owns (see Out of scope).

**`.bench/env.allow` grammar.** Optional, committed, absent means defaults only —
never an error. Line-oriented: `#` comments, blank lines, an `[agent]` section
header — the only known section — and one entry per line that is either an
exact name or a `PREFIX*` glob. Rejected with a fail-closed error naming the
line number and the reason: an entry before any section header, an unknown
section name (including `[gate]` — the gate opt-in is the manifest, not this
file), a bare `*`
or an entry whose glob is not a single trailing `*`, an entry containing a `/`
or `=`, and any character outside the portable environment-name set. A rejected
file refuses the launch; it never degrades to defaults, because a silently-ignored
opt-in is indistinguishable from a working one.

**Prompt transport.** The loop writes the iteration prompt to the adapter's
stdin and passes no positional argument. The adapter contract becomes: prompt on
stdin, `BENCH_SHIFT=1` armed, exit code passed through. Adapters change in
lockstep with the loop — no dual-mode transition, because a dual-mode adapter
would keep the argv path alive and testable-green while leaking. `claude -p` and
`codex exec` both read a piped prompt (probed 2026-07-20). `opencode run`
documents only the positional form as of 2026-07-20, so its adapter reads stdin
and passes it to the CLI positionally; the residual final-hop exposure is
recorded in `DATA_HANDLING.md` as a harness limitation to drop when upstream
documents stdin.

**Objective handling.** `validateObjective` gains a length cap of **200 runes**,
rejecting with the existing exit-2 usage posture and a message naming the actual
and maximum length. The cap is counted in runes, not bytes, so a multibyte
objective is not cut at a third of its apparent length. `.bench-objective` is
created mode 0600. `intent.Entry` drops its `Objective` field and the ledger's
validity rule stops requiring it; the entry key is the objective identifier.
Shift commit subjects keep the sanitized objective text — it is reviewer-authored
and a readable git history is a feature — with that durability documented.

**`bench status` keeps its objective detail without storing it.** For shift
entries, the renderer reads the objective back from `<worktree>/.bench-objective`,
which exists for exactly as long as the intent entry is live, and degrades to the
key when the file or the worktree is gone. Claude-agent entries have no worktree
and render the key alone. *This is a spec-time refinement the map did not name —
it follows the map's rule that the full text lives only in the worktree file, but
the consequence for the ambient dashboard was not priced there. Flagged for veto.*

**One sanitizer, by collapse rather than addition.** `intent.Preview` is already
a correct escaping sanitizer in the wrong package. It moves to a new
`internal/sanitize` as two functions: `Controls` (escape every control rune,
no cap) and `Preview` (`Controls` plus the existing 120-rune cap and byte-count
suffix). `dashboard.sanitize` is deleted and its call sites use `Controls`; the
shift loop's objective banner and result detail render through `Preview`;
`internal/intent` and `internal/status` call the new package. `internal/toon`
keeps its *refusal* semantics — refusing a control-bearing cell rather than
silently rendering a stripped one is a closed AXI-contract decision recorded in
the profile, and a distinct policy from sanitizing a human-facing string.

**No content-based secret detection anywhere.** Sensitivity is handled by
documenting which paths are durable, not by guessing at content. A passlisted
variable's value may still be sensitive; that is a documented fact, not a defect —
the passlist filters names, never values.

**`DATA_HANDLING.md` is repo-only for now.** It is not added to `package.json`
`files[]`; whether consumers receive it rides on FT85's payload allowlist.
`SECURITY.md` gains a pointer to it.

## Testing decisions

A good test here drives the real binary and observes what a subprocess or a
durable artifact actually contains — never what a Go-level unit believes it
passed. Every sentinel is a marker value planted in the parent and looked for in
the child, so a test that silently failed to launch anything cannot pass by
finding nothing.

Prior art: `internal/contract/runtime/runtime_gate_test.go` drives stub gates
through `BENCH_GATE`; `internal/contract/runtime/runtime_gate_proof_helpers_test.go`
writes an executable `.bench/gate.sh` stub that appends to a file; the shared
`contract.NewFixture` / `Fixture.Bench` / `Fixture.WriteExecutable` harness in
`internal/contract/helper.go` and `command.go` builds throwaway repos with an
isolated environment. `internal/conformance/line_routing_static_test.go`
(`checkLineBinding`) is the model for the doc-versus-constants derivation check.

The gate command is the project gate: `.bench/gate.sh`.

### Seam diagram

**Seam A — adapter environment construction (`internal/env`, observed at the built binary).**

    trigger: `bench shift`
        │
        ▼
    parent os.Environ()  ──▶  [ internal/env.Build(root) ]  ──▶  ordered env slice
    .bench/env.allow     ──▶  [                          ]  ──▶  fail-closed error
                                   ◀ tests attach here: a stub adapter dumps its
                                     environment to a file outside the fixture repo; the
                                     test asserts on the dumped contents and on the CLI's
                                     exit code and stderr

**Seam A2 — gate subject environment (FT78's manifest-declared subject, pinned at the built binary).**

    trigger: `bench gate` in a fixture repo with a stub gate
        │
        ▼
    parent os.Environ()      ──▶  [ gate verdict subject: PATH +   ]  ──▶  gate script env
    .bench/gate-inputs.json  ──▶  [ manifest-declared names (FT78) ]
                                   ◀ tests attach here: the stub gate dumps its
                                     environment; the test asserts the parent marker is
                                     absent while PATH and a manifest-declared name
                                     survive

**Seam B — prompt transport (loop → adapter, observed at the built binary).**

    trigger: each iteration of `bench shift`
        │
        ▼
    iteration prompt  ──▶  [ shift loop → adapter process ]  ──▶  adapter stdin
                           [                              ]  ──▶  adapter argv
                                   ◀ tests attach here: the stub adapter records its own
                                     /proc/self/cmdline and its full stdin to two files;
                                     the test asserts the marker's absence from one and
                                     byte equality against the prompt in the other

**Seam C — durable objective state (observed at the built binary).**

    trigger: `bench shift <objective>`
        │
        ▼
    objective text  ──▶  [ validateObjective → worktree file  ]  ──▶  .bench-objective (0600)
                         [ → intent ledger → commit → stdout  ]  ──▶  bench-intent.json (key)
                                                                 ──▶  commit subject (text)
                                   ◀ tests attach here: the test reads the fixture's ledger
                                     JSON, stats the scratch file's mode bits, and checks
                                     exit code plus ledger absence for a rejected objective

**Seam D — inventory-versus-constants derivation (`internal/conformance`).**

    trigger: the gate's conformance phase; the canary phase re-runs it on a broken fixture
        │
        ▼
    internal/env constants  ──▶  [ checkDataHandlingDerivation ]  ──▶  diagnostics, or none
    DATA_HANDLING.md        ──▶  [                             ]
                                   ◀ tests attach here: TestRootConformance over the real
                                     tree (green) and a canary fixture that adds a constant
                                     without documenting it (red, targeted substring)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A marker variable exported in the test parent is absent from the adapter's dumped environment | A | new contract test; red before the passlist exists because today's adapter inherits `os.Environ()` verbatim | The marker can only be absent if construction filters by name; an unfiltered inherit passes it through unchanged |
| 1 | Passlisted names present in the parent survive into the adapter's environment | A | new contract test; red against a degenerate empty-environment implementation | Pins the passlist as a filter rather than a wipe, which the marker-absence row alone would accept |
| 2 | A marker variable exported in the test parent is absent from the gate's dumped environment | A2 | green on the shipped tree — FT78 already closes the subject; red capability proven once by a targeted mutation appending the marker to the subject environment, with the demonstrated red recorded in the test | Pins the closed subject so a future change widening the gate environment turns the gate red instead of green-by-silence |
| 2 | `PATH` and a name declared in the fixture's `.bench/gate-inputs.json` survive into the gate's environment | A2 | new contract test; red against a subject that ignores the manifest declaration or wipes the environment | Pins the subject as PATH-plus-declared rather than empty, which the marker-absence row alone would accept |
| 3 | A name added under `[agent]` in `.bench/env.allow` reaches the adapter | A | new contract test; red before parsing exists because the file is ignored entirely | The opt-in is the only sanctioned widening; a build that ignored the file would green every default-passlist row while breaking the one mechanism consumers are told to use |
| 3 | An entry under an unknown section — including a stale `[gate]` — refuses the launch and names its line | A | new contract test; red against a parser that silently skips unknown sections | A silently-ignored section is indistinguishable from a working opt-in, and `[gate]` is the form the retired draft grammar would most plausibly leave behind |
| 3 | A `PREFIX*` entry admits every matching name in its section | A | new contract test; red against an exact-match-only parser | The map names glob support as load-bearing; exact-only matching is the likeliest cheap implementation |
| 4 | A malformed `env.allow` refuses the launch, exits non-zero, and names the offending line number | A | new contract test; red before the parser exists because the file is ignored and the shift proceeds | Fail-closed is a posture, not a detail: a parser that skipped bad lines would still launch and still leak |
| 5 | The environment sentinel contract runs inside the gate's contract phase | A | `.bench/gate.sh` on a tree where the new contract file is present but the passlist is not wired | Confirms the proof is gate-attached rather than a test only a human runs |
| 6, 7 | The stub adapter's `/proc/self/cmdline` contains no prompt marker | B | new contract test; red today because the loop passes the prompt as `argv[1]` | Reads the kernel's own view of the process's arguments, which is the exact surface a process listing exposes |
| 6, 7 | The stub adapter's stdin equals the iteration prompt byte for byte | B | new contract test; red against an implementation that removes argv without wiring stdin | Without it, deleting the argument alone — the cheapest way to green the cmdline row — would pass |
| 7 | Each shipped adapter, run directly with a prompt on stdin and a stub harness CLI on `PATH`, delivers that prompt to the CLI | B | new contract test, once per adapter; red against adapters left reading `$1` | Without it the cheapest wrong implementation — drop the loop's argument and leave the adapters unchanged — greens every stub-adapter row while every real shift launches with an empty prompt |
| 7 | The stub harness CLI's own `/proc/self/cmdline` carries no prompt marker for the claude and codex adapters | B | new contract test; red today because both adapters pass the prompt positionally after `--` | The final hop is where the map's exposure claim actually lands; the loop-side assertion says nothing about it |
| 7 | The opencode adapter is asserted to pass the prompt positionally, matching its documented limitation | B (edge of 7) | new contract test; red against an adapter changed to a stdin form the CLI does not document | Pins the known residual as a deliberate, documented state rather than letting it drift silently either way |
| 8 | `.bench/BENCH-reference.md` describes stdin as the adapter contract and no longer describes the `$1` form | D | `.bench/gate.sh` conformance phase; red while the reference still documents a transport the loop no longer uses | The kit's docs-currency and stale-reference checks are the gate's existing owner of prose that contradicts code |
| 9 | The argv sentinel contract runs inside the gate's contract phase | B | `.bench/gate.sh` on a tree with the contract file present and the loop still passing argv | Same gate-attachment proof as story 5, for the transport half |
| 10 | An objective longer than 200 runes exits 2 with a message naming the actual and maximum length | C | new contract test; red before the cap exists because an over-long objective is accepted and a shift begins | Pins both the rejection and its usage-error posture, which a silent truncation would fail |
| 10 | A 200-rune objective is accepted and a 201-rune objective is rejected | C (edge of 10) | new contract test; red against an off-by-one boundary | Boundary values are where a cap implemented against the wrong comparison operator shows |
| 11 | `.bench-objective` in a live shift worktree has mode 0600 | C | new contract test; red today because the file is written 0644 | Stats the real file, so a change to the constant that missed the call site cannot pass |
| 12 | The intent ledger JSON for a shift contains no field carrying the objective text, and a marker objective's text appears nowhere in the file | C | new contract test; red today because `Entry.Objective` stores it verbatim | Searching the whole serialized file, not one named field, survives the field being renamed rather than removed |
| 12 | An entry with no objective text is accepted by the ledger's validity rule | C | existing `internal/intent` unit tests; red once the field is dropped while `validEntry` still requires it | The validity rule is the gate that would reject every new entry, breaking the ledger rather than the field |
| 13 | `bench status` shows the objective for a live interrupted shift after the ledger field is gone | C | new contract test; red against an implementation that drops the field without reading the worktree file | Preserves the ambient contract the profile names; a key-only render would otherwise pass every other row |
| 13 | `bench status` renders the entry without error when the worktree or its `.bench-objective` is absent | C (edge of 13) | new contract test; red against an implementation that propagates the read error | A missing scratch file is the normal end state, so a strict read would break status for every completed shift |
| 14 | A control sequence in operator-influenced terminal text renders escaped rather than raw, through the shared sanitizer | C | existing `internal/dashboard` control-byte test plus a new assertion on the shift banner; red on the shift banner today, which applies no sanitizer of its own | The banner is the one render path with no sanitizer behind it, so it is where the collapse is actually load-bearing |
| 14 | `internal/toon` continues to refuse a control-bearing cell rather than sanitizing it | C | existing byte-pinned `internal/toon` tests; red if the collapse rewires TOON onto the escaping policy | Guards the closed AXI decision against being folded into the collapse by mistake |
| 14 | Exactly one non-test package outside `internal/toon` implements control-rune escaping | D | new conformance check; red on the current tree, which has the routine in both `internal/intent` and `internal/dashboard` | The story's claim is single-sourcing, not rendering; without this a fourth copy added for the banner would pass every behavioral row |
| 15 | A control-byte objective is rejected with no intent entry and no worktree created | C | new contract test; red against a validator that runs after the ledger write | Ordering is the claim; asserting rejection alone would pass an implementation that wrote first and cleaned up after |
| 16, 17 | Adding a name to a passlist constant without documenting it in `DATA_HANDLING.md` turns the conformance phase red with a targeted diagnostic | D | new canary fixture; red by construction | This is the derivation claim itself — the check exists so the doc cannot drift from the constants |
| 17 | The conformance check passes on the real tree | D | `.bench/gate.sh` conformance phase | Proves the shipped inventory's variable listing is actually complete rather than the check being satisfied by an empty doc |
| 16 | The inventory's non-variable sections — prompt, file, log, network, cache, and retention paths — describe the paths that exist | D | not TDD-able: no machine-checkable derivation exists for prose describing a path's retention, so this is a review-owned residual carried into `/bench-review-implementation` | Stated so the gap is a decision on the page; the derivation check covers only the variable listing, and a doc that is a variable table plus stubs would otherwise pass |
| 18 | The canary fixture fails for its own targeted substring and not for an unrelated reason | D | `.bench/gate.sh` canary phase | The canary layer's own contract: a fixture that goes red for the wrong reason proves nothing about the check |
| edge of 3, 4 | An `env.allow` whose last line lacks a trailing newline parses that line normally | A | new unit test in `internal/env`; red against a naive split that drops the final entry | Named in the profile's hostile-input checklist; a dropped final entry silently narrows the passlist |
| edge of 3, 4 | An `env.allow` that is present but empty is accepted and yields the defaults | A | new unit test; red against a parser that treats zero entries as malformed | The profile checklist requires absent and present-but-empty to be distinct asserted behaviors |
| edge of 4 | A bare `*` entry and a mid-name glob are both rejected by line number | A | new unit test; red against a parser that admits any `*` | A bare `*` is the wholesale inherit-everything hatch the map explicitly rejected |
| edge of 1 | No default glob admits a name outside the family it cites | A | new unit test enumerating the glob entries against a fixture of foreign names | The straddle rule binds every future addition; without the row an exact set could be "simplified" into a wider glob that admits a name nobody documented |
| edge of 6, 7 | A prompt containing a newline and a leading dash arrives on stdin unaltered | B | new contract test; red against an implementation that line-buffers or strips the prompt | Iteration prompts are multi-line by construction, and a dash-leading prompt is the case the old `--` guard existed for |
| edge of 1 | A passlisted variable holding a multi-line or very large value passes through unaltered | A | new unit test; red against an implementation that truncates or splits on newlines | The map states values pass through untouched by design; a value-mangling filter would corrupt legitimate configuration |

### Edge inventory

The canonical classes plus the project profile's hostile-input checklist, walked
per behavior. Everything not landed as a row above is a decision here.

- **Error path, malformed input, boundary values, empty/absent input, hostile
  environment** — landed as rows above.
- **Re-run idempotency** — landed implicitly: `.bench-objective` is recreated
  each shift and the ledger is keyed per entry, both already covered by existing
  worktree-reuse contracts.
- Paths or directory names containing spaces or glob characters — **Won't
  handle** as a new row: `contract.NewFixture`'s `WithSpacePath()` already runs
  the shift surface under a spaced path, and nothing in this feature parses a
  path.
- Special files (FIFOs, devices, sockets) at `.bench/env.allow` — **Won't
  handle**: the file is read with an ordinary bounded read that fails closed on
  any read error, which is the same posture a special-file rejection would reach.
- Invocation through a symlink, and through every shipped surface — **Won't
  handle**: this feature adds no new entry point, and the existing surface-parity
  contract already covers `bench shift` and `bench gate` across the real kit CLI,
  the linked by-path CLI, hooks, and adapters.
- Interrupt (SIGINT) mid-loop — **Won't handle** as a new row: the scratch-file
  and lease behavior under interrupt is unchanged by this feature, and the
  existing preserve-and-recover contracts own it. The one interaction — a
  retained worktree keeping a 0600 `.bench-objective` alive — is documented in
  `DATA_HANDLING.md` under retention rather than asserted.
- cwd deeper than the repo root — **Won't handle**: `internal/env` takes the
  resolved repo root as a parameter rather than discovering it, so there is no
  cwd assumption to break.
- Missing required tool on PATH — **Won't handle** as a new row: `PATH` is on
  the passlist (and in FT78's gate subject) precisely so this does not change,
  and the existing adapter-not-executable preflight already asserts the
  failure mode.
- Compatibility probe: the two harness stdin contracts were probed live and by
  documentation on 2026-07-20 and recorded in the map; opencode's absence of a
  documented stdin contract is carried as a stated limitation rather than a
  promise, so there is no unproven compatibility claim in this spec.

## Out of scope

- **Environment minimization for the other subprocess classes** — git fetch,
  model discovery, hooks, the canary runner, and the kit's own four-phase gate
  runner each launch with their own inherited environment. Separate capability:
  FT87 owns bounded subprocess behavior generally, and the two features would
  collide at the same call sites; the four-phase runner joins this group
  because a passlist there corrupted the kit's own release-evidence probe when
  tried (map #7). Estimate to fold in later once `internal/env` exists: 8
  edits, 4 gate runs.
- **Content-based secret detection or redaction** — scanning values for
  credential-shaped strings is a different discipline with its own false-positive
  policy, explicitly rejected by the map in favor of documented durability.
  Estimate: 10 edits, 5 gate runs.
- **Consumer delivery of `DATA_HANDLING.md`** — whether the file ships in the
  npm payload is FT85's allowlist decision, not this feature's.
  Estimate: 2 edits, 1 gate run, once FT85 lands.
- **The mode-0600 file prompt transport** — the fallback design for a future
  harness that can read neither stdin nor argv-free input. Not built now because
  no shipped harness needs it and an unused transport rots.
  Estimate: 5 edits, 2 gate runs.
