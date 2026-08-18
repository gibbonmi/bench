# bench-front-door

Status: staged

Decision source: docs/audits/2026-08-bench-capability/results-fable-high/action-items.yaml entry `A3` (named reviewed artifact from the 2026-08 capability audit; ledger L-03, L-08, L-15 carry its reproduced evidence; dependencies A1 and A2 landed at `a2914fd5`/`f1135fd6` and `2915fcc0`)

Introduces commands: /bench-drain, $bench-drain

Verification log: <pending> — filled at review close

## Problem

Nothing routes from observed state to the next thing to do. The wrapper's bare
`bench` prints a 44-line inventory; the binary's prints `no subcommand`; the two
disagree about `help`. `bench status` prints `clean — nothing pending` in a repo
that has never been adopted, and prints nothing at all for a valid staged spec, so
the middle of the workflow — spec staged, build next — is invisible. Six of the
board's eleven actions are prose (`fix before commit`, `split (craft-seams)`,
`resume interrupted work`), which nothing can invoke, and `bench handoff` already
has to skip past them to find a command. The one phase named like a router,
`/bench-what-next`, is roadmap maintenance. Each harness therefore starts a session
by guessing.

## Solution

One route, projected from the board that already exists. `bench status --route`
ranks the existing severity ladder and emits exactly one
`next[1]{state,why,command}` row plus a single runners-up line; it adds a `setup`
signal when `.bench/` is absent and a `specs: staged` state read from the spec
Status reader, and every board action becomes an invocable command or empty so the
route (and the handoff, which consumes the same selection) never has to skip
prose. Thin `/bench` (Claude) and `$bench` (Codex) adapters run the route and load
only the recommended phase. The wrapper's bare invocation prints the route and the
binary agrees with it, with the inventory on `bench help`. `/bench-what-next`
becomes `/bench-drain`, with a one-release alias. No new state file, no compiler,
no second derivation of any signal.

## User stories

### Group A — the route

Line: opus / high. Multi-seam Go work (a new flag, a shared selection owner that
`bench status` and `bench handoff` both consume, TOON output) with an exact
per-scenario contract; mid tier at high effort per `craft-line`.

1. As an agent starting a session, I want `bench status --route` to print exactly one `next[1]{state,why,command}` TOON row, so that I have one owner of "next" instead of a board to interpret.
2. As an agent, I want the route's lead to be the first row of the existing severity ladder that carries a command, so that the route never ranks differently from `bench status`.
3. As an agent, I want a second line `also:` naming the runners-up (`signal (detail) → command`, ladder order, empty commands omitted, `none` when there are none), so that I can see what the lead displaces without a second call.
4. As an agent in an adopted repo whose board is empty, I want the route to answer `clean` / `nothing pending` / `bench roadmap` when `ROADMAP.md` exists at the root and `/bench-drain` when it does not, so that a quiet tree still has a next command.
5. As an agent whose board has rows but none with a command (a locked-pending gate, an orphaned review pickup), I want the lead row rendered with an empty `command` and no fallback, so that I am not sent to the roadmap while a live gate owner holds the tree.
6. As a Codex session, I want `bench status --route --harness codex` to render `$bench-…` where Claude sees `/bench-…`, through the one prefix table `bench handoff --harness` already uses, so that the two harnesses cannot drift.
7. As the reviewer, I want `bench handoff`'s `## Next command` derived from the same selection the route prints — including the clean-board fallback and the empty-command state — so that the handoff and the route can never disagree about what is next.
8. As an agent, I want `--route --all`, `--harness` without `--route`, an unknown `--harness` value, and any trailing token to be usage errors (exit 2) with the grammar line, so that a malformed call never prints a route that looks authoritative.
9. As an agent outside a git repository, I want `bench status --route` to return the existing structured not-in-repo error (exit 1), so that the route never invents a state for a directory it cannot read.
10. As the reviewer, I want `bench status --route` in this repository at HEAD to emit one row, so that the front door demonstrably works on the kit itself.

### Group B — two new signals (and one flagged extension)

Line: opus / high. New board states must fit the existing ladder without a second
ordering and read from the existing spec and map readers.

11. As an agent in a repo with no `.bench/` directory, I want a `setup` row (`no .bench/` → `bench setup`) to lead the board and the route, so that an un-adopted repo never reads as `clean`.
12. As an agent, I want the `specs` row family to gain a staged state — `N staged spec(s)` from the spec Status reader — so that a staged spec is visible on the board.
13. As an agent, I want the staged state's command to be `/bench-implement-spec specs/<slug>/spec.md` when exactly one spec is staged and bare `/bench-implement-spec` when several are, so that the route names the file when the choice is unambiguous and leaves the choice to the phase when it is not.
14. As the reviewer, I want the staged row to rank at severity 4 ahead of the drain row (append order breaks the tie under the stable sort), below intent and guards, so that a staged spec outranks maintenance but never a red gate, a dirty tree, or a hung shift.
15. As an agent whose tree carries both a staged spec and a merged spec awaiting retirement, I want two `specs` rows (staged at 4, retirement at 8), so that neither state hides the other.
16. As an agent whose staged spec has a paired `reviews/<slug>.md`, I want the command to stay `/bench-implement-spec …`, so that review findings return to the implementing phase on the same source as the review phase already instructs.
17. **[reviewer-flagged addition]** As an agent whose top-level `decisions/` holds only ready maps (no shaping or invalid map), I want the `decisions` row to read `N ready map(s)` with command `/bench-write-spec decisions/<map>.md` (one) or `/bench-write-spec` (several), so that the audit's own "shaped map ready" scenario yields a command; shaping and invalid maps keep winning the row.
18. As an agent, I want the setup and staged states to stay quiet when their inputs cannot be read (an unreadable `specs/`, a spec directory without `spec.md`, a `Status:` line inside a code fence, a CRLF Status line reading as `staged`), following the retirement state's existing posture, so that a housekeeping row can neither block the SessionStart hook nor fabricate a count.

### Group C — every action is a command or empty

Line: opus / high. A table-driven normalization across every producer, graded by
one predicate; the cheapest wrong implementation is one row left as prose.

19. As an agent, I want every board action to satisfy one predicate — opens `/bench-`, `bench `, or `git `, contains no step separator, carries no trailing prose — or to be the empty string, so that the route, the handoff, and the dashboard consume commands, never hints.
20. As an agent facing a red gate, I want the action `/bench-debug` (reviewer-contestable: the audit is silent; the alternative is `bench gate`), so that a failing oracle routes to the fix path.
21. As an agent facing a timeout, an unavailable verdict, or an exact-tip partial green, I want `bench gate --fresh`, and for interrupted-pending, invalid, or drifted records `bench gate`, so that every gate state names the run that resolves it.
22. As an agent facing a locked-pending gate, I want an empty action, so that nothing tells me to run over a live gate owner.
23. As an agent with a dirty tree, I want the `git` row's command to be `/bench-final-check`, and with a clean tree carrying unpushed commits or a unique branch `git push`, with the detail still listing every clause, so that a compound `a / b` action never reaches a reader.
24. As an agent with a live intent entry, I want `bench status --all`, so that the route names the surface that expands each interrupted entry with its kind, path, branch, and objective.
25. As an agent with a leased pool worktree or a typed worktree failure, I want `bench worktree list`, and for a porcelain failure `git worktree list`, and for an out-of-pool worktree `bench worktree clean <path>`, so that every worktree state names its inspection.
26. As an agent, I want `bench structure` for structural debt, `bench maps` for a failed decisions scan, `git status` for unavailable git state, and an empty action for an orphaned review pickup, so that no prose survives outside a detail cell.
27. As a dashboard reader, I want an empty action rendered as an empty cell, so that normalization cannot crash the HTML board or fabricate a hint.

### Group D — thin harness adapters

Line: sonnet / medium. Two markdown adapters and guide mentions, graded by the
existing conformance globs; mechanical once the route exists.

28. As a Claude Code user, I want `/bench` (`.agents/commands/bench.md`, model-invocable, mirrored through `.claude/commands`) to run `bench status --route`, follow the one row — reading `.agents/commands/<phase>.md` when the command opens `/bench-`, running the command exactly otherwise — and load nothing else, so that a session starts at the front door and not with a guess.
29. As a Codex user, I want `$bench` (`.agents/skills/bench/SKILL.md` + `agents/openai.yaml` with `allow_implicit_invocation: false`) to run `bench status --route --harness codex` and behave identically, so that both harnesses enter through the same door.
30. As the reviewer, I want `/bench` and `$bench` named in the shared guide (`.bench/BENCH.md` commands bullet and workflow), so that A1's conformance checks (`checkCommandGuideReferences`, `checkCodexCommandAdapters`, `checkClaudeSkillMirror`) grade the new pair green and every linked repo receives the front door.
31. As the reviewer, I want one table-driven test proving the same board renders `/bench-…` for claude and `$bench-…` for codex through the shared route owner (the existing `TestFirstInvocable` extended or moved), so that the audit's prefix acceptance is a red-capable test.

### Group E — wrapper root and binary alignment

Line: opus / high. Shell dispatch plus a Go verb graded by a fail-closed dispatch
registry; small but two surfaces must agree byte-for-byte.

32. As a user typing bare `bench`, I want the route (not the inventory), and the inventory on `bench help`, `--help`, `-h`, so that the public entry point routes.
33. As a user invoking the binary directly with no arguments, I want the same route, and on `help`/`--help`/`-h` the same inventory as the wrapper, so that the wrapper and the binary agree.
34. As the reviewer, I want the inventory text held once — in the binary's `help` verb, with the wrapper's help arm routing to it — so that the two entries cannot drift; the wrapper still owns wrapper-only verbs and the inventory stays text, not a registry projection.
35. As the reviewer, I want the new `help` dispatch name registered in the fail-closed routing registry (`subcommandRouting`, exempt: takes no arguments), so that the conformance suite stays green rather than red on the added verb.

### Group F — rename `/bench-what-next` → `/bench-drain`

Line: sonnet / medium. A mechanical rename fan-out over code strings, anchors,
conformance paths, canary fixtures, and docs, graded by A1's conformance and
`bench canary`; the alias keeps the stale-command sweep green.

36. As a user, I want `/bench-drain` (`.agents/commands/bench-drain.md`, `disable-model-invocation: true`) and `$bench-drain` to be the maintenance phase, so that the drain is named for what it does.
37. As a user on a handoff written before the rename, I want `/bench-what-next` and `$bench-what-next` to keep working for one release as thin alias files that say "renamed; read and follow `.agents/commands/bench-drain.md`", so that `## Next command` in an old handoff still resolves.
38. As the reviewer, I want every production string that names the phase — the `drain` and `roadmap` board rows, `bench roadmap`'s help row, the learnings `HarnessPhase`, the harness-phase whitelist (which accepts both names during the alias window), the dashboard's empty-roadmap prose, and `bench init`'s scaffold prose — to say `/bench-drain`, so that no live surface advertises the old name.
39. As the reviewer, I want every anchor row that pins `bench-what-next.md` re-pointed to `bench-drain.md` (needles and diagnostics included), and the two anchors whose needle is the phrase `/bench-what-next` to carry the new phrase, so that the guidance contracts follow the file.
40. As the reviewer, I want the conformance tests that read the phase file by path (`recurrence_maintenance_contract_test`, `docs_workflow_helpers_test`, the `registry_test` fixture table) to read `bench-drain.md`, and the six `what-next-*` canary fixture directories renamed `drain-*` with their `bench-what-next.md` fixture files renamed and EXPECT/BASE/MUTATE.json content following, so that `bench canary` and the conformance suite are green after the rename.
41. As a reader, I want the live prose that names the phase — README, `.bench/BENCH.md`, `.bench/BENCH-reference.md`, CONTEXT.md, `projects/benchkit.md`, the sibling command files and `craft-synthesis`, ROADMAP.md, and the two `docs/*.md` guides — to say `/bench-drain`, and dated occurrence lines, `docs/audits/`, and `capture/` left untouched, so that current guidance is right and history is not rewritten.
42. As the reviewer, I want the stale-command sweep to stay green during the alias window (it derives valid `/bench-*` tokens from `.agents/commands/*.md`) and to go red the release the alias files are deleted, so that leftover references cannot outlive the alias.

## Implementation decisions

- **One route owner.** The selection that turns the severity ladder into "the next
  command" — the invocable predicate, the harness prefix table and its translation,
  the first-invocable walk, and the clean-board fallback — moves into the status
  package (which already owns the ladder) or a package both status and handoff import;
  the handoff package consumes it and keeps no copy. `bench handoff`'s `Action`,
  `Signal`, and no-invocable state are read from that owner. Recording it there rather
  than in handoff avoids the import cycle (handoff already imports status).
- **Route output.** `next[1]{state,why,command}` via the TOON table writer, then one
  `also:` line. `state` is the signal name, `why` is that row's detail, `command` is
  its action; the fallback row is `clean` / `nothing pending` / `bench roadmap` or
  `/bench-drain` (ROADMAP.md present or absent). A board with rows but no command
  renders the lead row with an empty `command`; the fallback applies only to an empty
  board. Runners-up are every later row with a non-empty command; `also: none` otherwise.
- **Grammar.** `bench status [--all] [--route [--harness claude|codex]]`; `--route`
  and `--all` are exclusive; `--harness` requires `--route`; harness names come from
  the shared prefix table exactly as `bench handoff` validates them.
- **Signals.** `setup`: severity 0, appended before the gate row, fires when the git
  root has no `.bench/` directory; detail `no .bench/`, command `bench setup`.
  `specs` staged: severity 4, appended before the drain row; count from the spec
  Status reader (`Status: staged` on a parsed `spec.md`), unreadable or malformed specs
  count zero as the retirement state already does. **Flagged:** `decisions` ready
  state, severity 6 on the existing row, only when no shaping/invalid map is active.
- **Normalization table** (state → command; unchanged commands not listed): gate
  locked-pending → empty; interrupted-pending, invalid, drifted → `bench gate`;
  unavailable, timeout, exact-tip partial → `bench gate --fresh`; red → `/bench-debug`
  (contestable). git dirty → `/bench-final-check`; else unpushed or unique branch →
  `git push`; state unavailable → `git status`. intent (both) → `bench status --all`.
  worktree typed failure and leased → `bench worktree list`; porcelain failure →
  `git worktree list`; out-of-pool → `bench worktree clean <path>`. structure →
  `bench structure`. decisions scan failure → `bench maps`. reviews orphan → empty.
  drain and roadmap → `/bench-drain`. The invocable predicate gains the `git ` prefix
  and rejects trailing prose after a command; placeholders `<path>`/`<slug>`
  remain allowed. Typed worktree failure text moves into the detail cell.
- **Adapters.** `/bench` is model-invocable (no `disable-model-invocation` key);
  `$bench` follows the Codex adapter shape (`name: bench`, references the command file,
  `allow_implicit_invocation: false`). Both load only the phase the route names.
- **Root and help.** Wrapper bare invocation and binary no-arg both run
  `status --route` and return its exit code; the inventory text lives in the binary's
  `help` verb (exempt in the routing registry) and the wrapper's `help|--help|-h` arm
  routes to it. Trade the reviewer decides: `bench help` then needs the binary
  (`route_binary` exits 127 with its install message when absent); the alternative is
  leaving the heredoc in the wrapper and having the binary's `help` print a one-line
  pointer to `bench help`.
- **Rename.** `bench-drain.md` is the phase; `bench-what-next.md` and its Codex pair
  remain as thin aliases (same `disable-model-invocation` posture) for one release; the
  harness-phase whitelist accepts both during the window; the anchors registry, the
  path-hardcoded conformance tests, and the canary fixture tree follow the file.

## Testing decisions

- A good test drives a real repository fixture into one scenario and observes the
  route's single row (and the wrapper/binary bytes), never the internals of ranking.
- Seams and prior art: the status package's in-memory board tests (`status_test.go`
  fixtures via `gittest`); the handoff package's `TestFirstInvocable`; the dashboard
  render test; the shell/binary agreement through the cmd/bench main test, which already drives
  `bin/bench.sh` with bash and the in-process `Command{}.Run`; the conformance suite and `bench canary`
  for the adapters and the rename; the routing registry test for the `help` verb.
- The gate seam: `bench gate` runs the Go tests, the conformance suite (A1), the canary
  inventory, `bench coverage --check`, and `bench skills-index --check`.

### Seam diagram

    trigger: `/bench` · `$bench` · bare `bench` · session start · `bench handoff`
        │
        ▼
    repo state ──▶ [ status.Signals (ladder) ] ──▶ [ route owner: first command / fallback / harness prefix ]
                     ◀ tests: gittest fixtures per scenario, assert one next row + also line
                                                              │
                                        ┌─────────────────────┼──────────────────────┐
                                        ▼                     ▼                      ▼
                              bench status --route      bench handoff Next     adapters (.md) → phase file
                              ◀ TOON bytes             ◀ render test          ◀ conformance globs

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R1 | 1 | `bench status --route` on a fixture with a red gate prints exactly one `next[1]{state,why,command}` table with the row `gate,red,/bench-debug` | status route test | a build that prints the board, or two rows, fails the exact byte match |
| R2 | 2 | a fixture whose ladder is [gate red, git dirty, drain] routes `gate` and never `git` or `drain` | status route test | a build that re-ranks (e.g. prefers a phase) picks `drain` and fails |
| R3 | 3 | the same fixture prints `also: git (1 dirty path) → /bench-final-check; drain (1 idea(s), 0 open learning(s), 0 pending retro(s)) → /bench-drain` as its second line | status route test | a build that omits runners-up or includes an empty-command row fails the exact line |
| R4 | 3 | a fixture with one invocable row prints `also: none` | status route test | an omitted `also:` line fails |
| R5 | 4 | an adopted, empty-board fixture with `ROADMAP.md` routes `clean,nothing pending,bench roadmap` | status route test | a build printing `clean — nothing pending` with no command fails |
| R6 | 4 | the same fixture without `ROADMAP.md` routes `clean,nothing pending,/bench-drain` | status route test | a build that always says `bench roadmap` fails |
| R7 | 5 | a fixture whose only row is gate `locked-pending` routes `gate,locked-pending,` with an empty command and `also: none` | status route test | a build applying the fallback prints `bench roadmap` and fails |
| R8 | 6 | `--harness codex` on the drain fixture prints `$bench-drain` in the row and in `also:`, and `--harness claude` prints `/bench-drain` | status route test | a build translating only the lead row fails on `also:` |
| R9 | 7 | `bench handoff` on the drain fixture writes `## Next command` naming `/bench-drain`, and on the empty-board fixture with `ROADMAP.md` names `bench roadmap` | handoff render test | a handoff keeping its own selection prints no command on the clean board and fails |
| R10 | 7 | the handoff's no-invocable state is set exactly when the route's lead has an empty command on a non-empty board | handoff facts test | a second derivation of "no command" drifts on the locked-pending fixture |
| R11 | 8 | `--route --all`, `--harness codex` alone, `--route --harness opencode`, and `--route extra` each exit 2 with the grammar line | status usage test | a build accepting any of the four prints a route and fails |
| R12 | 9 | `bench status --route` outside a repository returns the structured not-in-repo error and exit 1 | status command test | a build that panics or prints `clean` fails |
| R13 | 10 | `bench status --route` run at this repository's HEAD emits exactly one `next[1]` row | cmd/bench main test driving `bin/bench.sh` | the audit's own acceptance; a route that emits zero or two rows on the kit fails |
| R14 | 11 | a git repo without `.bench/` renders the board lead `▶ bench setup  (setup)` and routes `setup,no .bench/,bench setup` | status board + route test | a build still printing `clean — nothing pending` fails |
| R15 | 11 | an adopted fixture emits no `setup` row | status board test | a build keying on something other than the `.bench/` directory fires falsely |
| R16 | 12, 13 | a fixture with one `specs/<s>/spec.md` carrying `Status: staged` renders `specs  1 staged spec(s)  → /bench-implement-spec specs/<s>/spec.md` | status board test | a build with no staged state prints nothing for the spec |
| R17 | 13 | a fixture with two staged specs renders `2 staged spec(s)` with the bare command `/bench-implement-spec` | status board test | a build that always names the first path fails |
| R18 | 14 | a fixture with a staged spec and a drain source orders `specs` before `drain` and after `guards` | status board test | a build appending staged after drain, or at severity 8, fails the order |
| R19 | 15 | a fixture with one staged and one implemented spec renders two `specs` rows, staged first | status board test | a build collapsing the family to one row fails |
| R20 | 16 | a fixture with a staged spec and `reviews/<s>.md` still routes `/bench-implement-spec specs/<s>/spec.md` | status route test | a build routing to review or land fails |
| R21 | 17 | a fixture whose `decisions/` holds one `Status: ready` map and nothing else routes `decisions,1 ready map(s),/bench-write-spec decisions/<m>.md` | status route test | without the flagged state the board is empty and the fallback answers `bench roadmap` |
| R22 | 17 | a fixture with one ready and one shaping map keeps `1 unresolved map(s)` → `/bench-shape-idea` | status board test | a build letting ready win fails |
| R23 | 18 | a spec directory without `spec.md`, a `Status: staged` inside a code fence, and an unreadable `specs/` FIFO each count zero staged and the board returns | status board test | a build counting the fenced line, or blocking on the FIFO, fails |
| R24 | 19 | a table-driven test over every producible board row asserts each action is invocable-or-empty under the one predicate | status action test | leaving one prose action (the cheapest wrong build) turns exactly that row red |
| R25 | 19 | the invocable predicate rejects `bench gate --fresh for a whole-tree verdict`, `/bench-final-check / push`, and `benchmark the loop`, and accepts `git push`, `bench spec retire <slug>`, `/bench-drain` | route owner unit test | a prefix-only predicate admits the first two |
| R26 | 20, 21, 22 | gate fixtures red, timeout, invalid, unavailable, drifted, exact-tip partial, locked-pending render commands `/bench-debug`, `bench gate --fresh`, `bench gate`, `bench gate --fresh`, `bench gate`, `bench gate --fresh`, and empty respectively | status board test | any state keeping its old prose fails its cell |
| R27 | 23 | a dirty tree renders `git` command `/bench-final-check` with detail listing dirty and unpushed clauses together | status board test | a build joining `commit on green / /bench-final-check` fails |
| R28 | 23 | a clean tree with unpushed commits renders `git push`, and a unique branch alone renders `git push` | status board test | a build keeping `push` prose fails |
| R29 | 24 | a live intent entry renders command `bench status --all` | status board test | `resume interrupted work` surviving fails |
| R30 | 25 | leased-pool, typed-failure, porcelain-failure, and out-of-pool worktree fixtures render `bench worktree list`, `bench worktree list`, `git worktree list`, `bench worktree clean <path>` | status board test | any old worktree prose fails its cell |
| R31 | 26 | structural debt, a failed decisions scan, unavailable git state, and an orphaned pickup render `bench structure`, `bench maps`, `git status`, and empty | status board test | any of the four old prose strings fails |
| R32 | 27 | the HTML dashboard on a board containing an empty-action row renders `<td></td>` for that cell and exits 0 | dashboard render test | a build that panics on the empty string or substitutes text fails |
| R33 | 28 | `.agents/commands/bench.md` exists with `description:` and no `disable-model-invocation` key, and its body names `bench status --route` and the two follow-on rules | conformance docs check | a missing or user-only-invocable file fails |
| R34 | 29 | `.agents/skills/bench/SKILL.md` names `bench` and references `.agents/commands/bench.md`, and `agents/openai.yaml` carries `allow_implicit_invocation: false` | conformance `checkCodexCommandAdapters` | the existing glob turns red on the missing pair |
| R35 | 30 | `/bench` and `$bench` appear in `.bench/BENCH.md` or `.bench/BENCH-reference.md`, and no `.claude/skills/bench` entry exists | conformance guide-reference and skill-mirror checks | the existing checks turn red on omission or a duplicate mirror |
| R36 | 31 | one table-driven test feeds a board carrying `/bench-drain` to the route owner and asserts `/bench-drain` for claude and `$bench-drain` for codex | route owner unit test | a translation applied only in handoff leaves the route untested |
| R37 | 32 | bare `bin/bench.sh` prints the `next[1]` table and `bin/bench.sh help` prints the inventory beginning `bench — Pocock pipeline` | cmd/bench main test driving `bin/bench.sh` | a wrapper still defaulting to the heredoc fails |
| R38 | 33 | the built binary with no arguments prints the same `next[1]` table as the wrapper, and `dist/bench help` prints byte-identical inventory to `bin/bench.sh help` | cmd/bench main test (`Command{}.Run` + `bin/bench.sh`) | `no subcommand` exit 2, or two inventories, fails |
| R39 | 34 | the inventory string appears once in the tree — in the Go `help` verb — and `bin/bench.sh` contains no `Pocock pipeline` heredoc | cmd/bench main test reading the tree | a second copy of the inventory turns red |
| R40 | 35 | `subcommandRouting` carries a `help` row and the routing check is green | conformance routing test | the fail-closed registry is red for the unregistered verb |
| R41 | 36 | `.agents/commands/bench-drain.md` carries the drain phase body and `disable-model-invocation: true`, and `.agents/skills/bench-drain/` carries the Codex pair | conformance adapter checks | the rename without the pair turns the glob red |
| R42 | 37 | `.agents/commands/bench-what-next.md` and `.agents/skills/bench-what-next/` exist as thin aliases naming `bench-drain.md`, and `bench handoff --next /bench-what-next` on a fixture still validates | conformance stale-command sweep + handoff test | deleting the alias makes old handoffs stale; the sweep goes red on any old reference |
| R43 | 38 | the drain and roadmap board rows, `bench roadmap` help, the learnings phase action, `bench init` scaffold, and the dashboard empty-roadmap prose all say `/bench-drain`, and `validHarnessPhase` accepts both names | status/roadmap/learnings/adopt/dashboard/axi tests | one string left behind fails its package test |
| R44 | 39 | `bench anchors .agents/commands/bench-drain.md` lists every row formerly pinned to `bench-what-next.md`, and the two phrase anchors need `/bench-drain` | anchors registry test | an anchor left on the old path pins a file that no longer carries the body |
| R45 | 40 | the six renamed `drain-*` canary fixtures validate under `bench canary` and the conformance registry table names them | canary + `registry_test` | a stale fixture name is red in the registry table |
| R46 | 41 | the stale-command sweep, `checkCommandGuideReferences`, and `bench skills-index --check` are green after the prose rename | conformance + skills-index | a missed `$bench-what-next →` mapping row or a dangling reference is red |
| R47 | 42 | with the alias files removed in a throwaway copy, the stale-command sweep reports every remaining `/bench-what-next` reference | conformance sweep mutation | proves the alias window is the only thing keeping old references green |

### Edge inventory

- Several staged specs → bare `/bench-implement-spec` (R17). Both staged and merged → two rows (R19).
- Spec dir without `spec.md`, fenced `Status:`, FIFO at `specs/` → count zero (R23). CRLF Status lines already parse through the shared reader; covered by the same reader's tests.
- Empty board with/without ROADMAP.md (R5, R6); rows-but-no-command (R7).
- `--harness` unknown, `--route --all`, trailing token, `--harness` alone (R11); outside a repo (R12).
- Un-adopted repo that is also dirty: `setup` leads (severity 0, appended first) — R14's fixture may be dirty.
- Non-primary checkout (pool worktree): the guards row stays quiet as today; `setup` reads the checkout's own root.
- Detached HEAD / no upstream: the git row's clauses degrade as today; the command precedence still yields one command (R27, R28).
- **Won't handle:** a route across worktrees (reading another checkout's spec or review state) — the route reads the tree it runs in plus the shared ledger; the destination-side "reviewed, awaiting land" scenario answers `bench worktree list` through the leased row (R30) because landing needs a request token only the assignment holds and a clean review leaves no artifact; `bench worktree list` remains the caller.
- **Won't handle:** an empty `ROADMAP.md` (zero bytes) still counts as present for the fallback — `bench roadmap` is then the surface that reports its emptiness.
- **Won't handle:** auto-detecting the harness — `--harness` is explicit as in `bench handoff`; the adapters pass it.
- Bootstrap authority: no trusted-execution or refusal-before-execution claim is introduced; the route names commands, it does not execute them.

## Ownership fences

- `internal/status/`
- `internal/handoff/`
- `internal/maps/`
- `internal/axi/action.go`, `internal/axi/action_test.go`
- `internal/dashboard/`
- `internal/roadmap/roadmap.go`, `internal/roadmap/roadmap_test.go`
- `internal/learnings/learnings.go`, `internal/learnings/learnings_test.go`
- `internal/adopt/init.go`
- `internal/anchors/registry_data.go`
- `internal/conformance/`
- `cmd/bench/`
- `bin/bench.sh`
- `tests/canary/`
- `.agents/commands/`
- `.agents/skills/bench/`, `.agents/skills/bench-drain/`, `.agents/skills/bench-what-next/`, `.agents/skills/bench-craft-synthesis/SKILL.md`
- `.bench/BENCH.md`, `.bench/BENCH-reference.md`
- `README.md`, `CONTEXT.md`, `ROADMAP.md`, `projects/benchkit.md`
- `docs/reporesident-distillation.md`, `docs/greenfield-build-sequence.md`
- `capture/session-handoff.md`
- `specs/bench-front-door/`

## Out of scope

- Generating root help, the `.bench/BENCH.md` inventory, and `commands --brief` from `commandRegistry` (FT89 remainder) — ~6 edits, 2 gate runs.
- Deleting the `what-next` alias files after one release and sweeping the residue — 2 deletions plus the references the sweep names, 1 gate run.
- A8 (constrained `## State`, `Repro:` line, Next from the router in the handoff body) — depends on this spec; own spec.
- A9 hygiene batch (session-start hint, log pruning, dist/bench landing warning) — own spec.
- L-16 prose contradictions beyond the lines the rename touches — ~5 edits, 1 gate run.
- Any change to the expert phases' behavior (A3 non-goal).
- Renumbering the severity ladder to give the staged state its own ordinal — ~4 edits (ladder comment, two tests, dashboard test), 1 gate run; the reviewer may swap it in for the tie rule.

## Further notes

- No `Roadmap:` line: the audit's board has not been drained into `ROADMAP.md`; the reviewer override records the audit order there. The drain that folds A3 in also folds FT89 (root/inventory generation, partial) and FT180 (route decision surface).
- Reviewer-visible calls, in one place: (1) `decisions: ready` state (story 17) is an addition beyond the source's two named signals, demanded by its own fixture list; (2) gate red → `/bench-debug`; (3) inventory relocation into the binary and the 127 trade; (4) staged severity 4 tie rule vs renumbering; (5) the "reviewed, awaiting land" scenario answering `bench worktree list`.
- Migration: an old handoff's `## Next command` naming `/bench-what-next` resolves through the alias for one release.
