# Guard Bench commands from shell follow-ons

Status: staged

Roadmap: FT248

Decision source: named reviewed artifact `roadmap/FT248.md`, committed by the 2026-08-23 reviewed drain

Verification log: 1 iteration(s) to accept — Terra xhigh returned four blockers, and the author folded all findings after the capped review

## Problem

Agents append shell follow-ons to `bench` commands before they read the output.
The follow-on can hide evidence or change the command's non-interactive input.

The always-loaded rule forbids this behavior, but no Bash guard enforces it.
A repeated reflex can therefore bypass the rule in every harness session.

## Solution

A PreToolUse Bash guard refuses an outer shell follow-on on a `bench` invocation.
The refusal explains that each Bench response is bounded, complete, and self-contained.

The guard parses the complete Bash command.
It therefore sees an operator after `bench worktree exec <label> -- <command>`.
It does not treat a quoted child command as outer shell syntax.
It resolves a live executable alias without running its target.

The guard uses a shared shell-command token stream with the destructive-Git guard.
The hook stays a thin shim over one internal Bench command.
Its degraded rim warns and allows when the wrapper resolver or core cannot run.

The always-loaded rule states why a projection or shell follow-on is invalid.
The implementation adds no projection option.

## User stories

### Refuse an outer shell follow-on

Line: gpt-5.6-luna / medium.

The source fixes the behavior, and the existing hook process is the known seam.

1. As an agent, I want a bare `bench` command to run, so that I receive its complete response.
2. As a reviewer, I want an outer pipeline refused, so that the pipeline cannot hide Bench evidence.
3. As a reviewer, I want an outer redirection refused, so that the shell cannot replace Bench input or output.
4. As a reviewer, I want an outer `&&` refused, so that Bench does not become part of a command chain.
5. As a reviewer, I want an appended shell command refused, so that an extra command cannot follow Bench.
6. As an agent, I want the full `bench worktree exec` command inspected, so that syntax after its child command cannot escape inspection.
7. As an agent, I want the refusal to explain Bench output, so that I do not retry another projection.

### Preserve intended Bash calls

Line: gpt-5.6-luna / medium.

The hostile-input partitions require exact shell syntax and a narrow degraded posture.

8. As an agent, I want supported Bench flags to remain valid, so that the guard does not replace Bench argument validation.
9. As an agent, I want quoted child syntax to remain valid, so that `worktree exec` can run one explicit child command.
10. As an agent, I want non-Bench commands to remain valid, so that incidental Bench text does not trigger the guard.
11. As a reviewer, I want JSON-escaped operators classified, so that the tool envelope cannot hide a shell follow-on.
12. As a reviewer, I want routine command prefixes inspected, so that a prefix cannot hide the reflex.
13. As an agent, I want a broken guard to warn and allow, so that missing kit infrastructure does not brick Bash.
14. As a reviewer, I want both harness adapters to invoke the guard, so that the rule is portable.
15. As an agent, I want `bench guards` to advertise the refusal, so that I can learn the active deny surface.

### Recognize each Bench executable spelling

Line: gpt-5.6-luna / medium.

The real shell producer emits path spellings and live symlink aliases.

16. As a reviewer, I want a repository-relative Bench path recognized, so that the guard covers by-path kit calls.
17. As a reviewer, I want an absolute Bench path recognized, so that the guard covers calls outside the repository root.
18. As a reviewer, I want a live symlink alias recognized, so that an alias cannot hide a Bench follow-on.

### Teach the complete-response contract

Line: gpt-5.6-sol / high.

The always-loaded text steers future generation, so the leverage override applies.

19. As an agent, I want the rule to state that Bench responses are bounded and complete, so that I do not add insurance projections.
20. As a reviewer, I want no new projection option, so that the implementation does not concede the premise that caused the reflex.

## Implementation decisions

One shell-command module owns quote folding, heredoc-body removal, operator tokens, redirection tokens, and simple-command spans.
Its interface exposes those token classes and spans to both analyzers.
The destructive-Git analyzer consumes the same token stream and retains its current classifications.

A Bench-follow-on classifier consumes that token stream and an injected executable identity resolver.
It recognizes a direct `bench` token and a path whose basename is `bench` or `bench.sh`.
For another spelling, the identity resolver follows the executable path without running it.
A resolved target whose basename is `bench` or `bench.sh` identifies a live symlink alias.
The identity resolver uses the process PATH for a bare alias and the process working directory for a relative alias.
It recognizes the same routine command prefixes and one wrapper depth as the destructive-Git analyzer.

The classifier refuses an outer pipeline, redirection, `&&`, `||`, semicolon, or newline on a Bash call that invokes Bench.
It examines the full outer command, including syntax after a `worktree exec` child argument.
It does not inspect operator characters that quote folding keeps inside one argument.

The classifier does not validate Bench subcommands, flags, or child arguments.
The Bench command grammar remains the one source for those facts.

The internal guard command reads `tool_input.command` from the PreToolUse envelope.
It owns the one refusal message and applies that message to every refused classifier verdict.
It prints one actionable refusal and exits 2 when the classifier refuses.
The refusal states that the Bench response is bounded, complete, and self-contained.

The hook shim uses the shared wrapper resolver.
A missing wrapper-resolver library, wrapper, binary, unreadable command field, or core error produces a warning and exit 0.
An empty command field produces the same warning and exit 0.
Only the reachable core returns an intentional exit 2.

Claude and Codex invoke the new script in their existing PreToolUse Bash group.
The script carries the static guard manifest that `bench guards` discovers.
The package and conformance inventories require the script and both wiring points.

The always-loaded CLI contract gives the reason before its prohibition.
The guidance anchor requires the reason and the valid-invalid example.
No public Bench command or flag joins the command registry.

### Bootstrap authority before execution

The harness loads the project hook configuration and invokes the tracked shim before the Bash command.
The shim locates the wrapper through the shared resolver and invokes the internal guard command.
The core then parses the envelope before the harness starts the requested Bash process.

No hop has an independent trust root outside the candidate tree.
The reviewer-visible trust assumption is an honest-mistake agent that does not alter or evade the guard.
The gate later verifies hook presence, harness wiring, package inclusion, and the planted refusal marker.
This design makes no authenticated-execution or security claim.

## Testing decisions

The highest seam is a real PreToolUse hook process with a marker in the refused shell follow-on.
The process test observes the exit code, diagnostic, and absent marker.

Classifier tests cover operator, redirection, quoting, prefix, wrapper, executable spelling, and command-position partitions.
The shared lexer tests retain every destructive-Git behavior after extraction.

The executable identity resolver is a local-substitutable filesystem dependency.
A temporary directory supplies the real target and the live symlink alias through the production resolver.

Conformance tests remove the script and each harness wire independently.
Each mutation must make the gate red with the named omission.
The guard query test observes the manifest-derived row.

The guidance canary removes the bounded-and-complete reason while it keeps the prohibition.
That mutation must make the guidance anchor red.

### Seam diagram

    Claude or Codex PreToolUse:Bash
                  │
                  ▼
        [ thin hook shim ] ──▶ shared wrapper resolver
                  │
                  ▼
        [ internal guard command ]
                  │
                  ▼
    command text ──▶ [ shared shell token stream ] ──▶ [ Bench-follow-on classifier ] ──▶ allow or refuse
                                                                    ▲
                                         executable identity resolver

    classifier tests attach through the token-stream and identity-resolver interfaces

    process tests attach at the hook exit, diagnostic, and marker

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FOG01 | 1 | A bare `bench gate --fresh` call exits through Bench without a guard refusal. | PreToolUse hook process | A classifier that rejects all Bench calls blocks the valid form. |
| FOG02 | 2 | `bench help \| touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A parser that ignores the pipeline creates the marker. |
| FOG03 | 2 | `bench help \|& touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | An operator spelling list misses the combined token and creates the marker. |
| FOG04 | 3 | `bench gate 2>&1` exits 2 before Bench runs. | PreToolUse hook process | A tokenizer that drops redirections allows the call. |
| FOG05 | 3 | `bench gate </dev/null` exits 2 before Bench runs. | PreToolUse hook process | An output-only check misses changed non-interactive input. |
| FOG06 | 4 | `bench help && touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A pipeline-only classifier allows the command chain and creates the marker. |
| FOG07 | 5 | `bench help; touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A classifier that ignores semicolons executes the extra command. |
| FOG08 | 5 | A newline between `bench help` and `touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A whitespace-only lexer loses the command separator. |
| FOG09 | 6 | An outer redirection after `bench worktree exec <label> -- bench gate --fresh` exits 2. | Bench-follow-on classifier | A child-only scan omits syntax that the outer shell applies. |
| FOG10 | 6 | An outer pipeline after `bench worktree exec <label> -- bench gate --fresh` exits 2. | Bench-follow-on classifier | A scan that stops at `--` omits the outer pipeline. |
| FOG11 | 7 | The internal guard routes a refused verdict through its one bounded, complete, and self-contained message owner. | internal guard command | An operator-specific denial can omit the reason when it bypasses the message owner. |
| FOG12 | 8 | `bench gate --fresh` remains allowed as one simple command. | Bench-follow-on classifier | Argument validation inside the guard rejects a supported flag. |
| FOG13 | 9 | `bench worktree exec <label> -- bash -lc 'go test && go vet'` remains allowed. | Bench-follow-on classifier | A quote-blind scan treats child syntax as an outer operator. |
| FOG14 | 10 | `rg bench AGENTS.md` remains allowed. | Bench-follow-on classifier | A token-presence check mistakes an argument for an invocation. |
| FOG15 | 10 | `printf hi > bench` remains allowed. | Bench-follow-on classifier | A basename search mistakes a redirection target for a verb. |
| FOG16 | 11 | A JSON-escaped `&&` after a Bench invocation exits 2. | internal guard command | An opaque escape placeholder welds the two commands together. |
| FOG17 | 11 | A JSON-escaped redirection after a Bench invocation exits 2. | internal guard command | A raw-text search misses the encoded operator. |
| FOG18 | 12 | One `bash -lc` wrapper around a Bench pipeline exits 2. | Bench-follow-on classifier | A top-level-only scan allows the routine wrapper. |
| FOG19 | 13 | A missing wrapper-resolver library or wrapper prints a warning and exits 0. | degraded hook process | A fail-closed shim bricks unrelated Bash recovery. |
| FOG20 | 13 | A core error prints a warning and exits 0. | degraded hook process | A nonzero passthrough turns infrastructure failure into denial. |
| FOG21 | 14 | Claude's PreToolUse Bash group invokes the follow-on guard. | conformance wiring check | A Claude-only omission leaves one harness unguarded. |
| FOG22 | 14 | Codex's PreToolUse Bash group invokes the follow-on guard. | conformance wiring check | A Codex-only omission leaves one harness unguarded. |
| FOG23 | 15 | `bench guards` reports the guard name, Bash seam, denial, and both wired harnesses. | guard query process | A missing manifest or wire disappears from the advertised surface. |
| FOG24 | 19 | The always-loaded CLI contract says each Bench response is bounded and complete evidence. | guidance conformance anchor | A prohibition without its reason survives a prose-only edit. |
| FOG25 | 20 | `bench gate --brief` remains a usage error. | public Bench command | A projection flag concedes the source's explicit exclusion. |
| FOG26 | 5 | `bench help extra \|\| touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A classifier that checks only successful chains allows the failure chain and creates the marker. |
| FOG27 | 12 | `env X=1 bench help \| touch <marker>` exits 2 and leaves the marker absent. | Bench-follow-on classifier | A bare-verb detector misses the invocation behind the environment prefix. |
| FOG28 | 13 | Malformed envelope JSON prints a warning and exits 0. | internal guard command | A silent allow passes without the required degraded diagnostic. |
| FOG29 | 13 | An absent command field prints a warning and exits 0. | internal guard command | A missing-field allow passes without the required degraded diagnostic. |
| FOG30 | 13 | A non-string command field prints a warning and exits 0. | internal guard command | A type-blind decoder treats the field as an ordinary empty command. |
| FOG31 | 13 | An empty command field prints a warning and exits 0. | internal guard command | An empty-string fast path allows without the required degraded diagnostic. |
| FOG32 | 13 | A command field with an escaped NUL prints a warning and exits 0. | internal guard command | A control-byte truncation can hide the operator without a degraded warning. |
| FOG33 | 13 | A missing platform binary prints a warning and exits 0. | degraded hook process | A wrapper exit 127 can become a denial or a silent allow. |
| FOG34 | 16 | From `<root>/internal/systemtest`, `../../bin/bench.sh help \| touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | An exact executable-token match or root-only assumption misses the repository-relative path. |
| FOG35 | 17 | `<absolute-root>/bin/bench.sh help \| touch <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A relative-path allowlist misses the worktree-specific absolute path. |
| FOG36 | 18 | `./kit-command help \| touch <marker>` exits 2 and leaves the marker absent when `kit-command` links to `bin/bench.sh`. | PreToolUse hook process | A basename-only classifier treats the alias name as a non-Bench command. |
| FOG37 | 3 | `bench help > <marker>` exits 2 and leaves the marker absent. | PreToolUse hook process | A redirection classifier that handles only input and descriptor duplication creates the stdout marker. |
| FOG38 | 3 | `bench gate <<'EOF'<newline>input<newline>EOF` exits 2 before Bench runs. | PreToolUse hook process | A tokenizer that removes the heredoc operator with its body allows the input redirection. |
| FOG39 | 10 | `cat <<'EOF'<newline>bench gate \| tail -20<newline>EOF` remains allowed. | Bench-follow-on classifier | A tokenizer that scans heredoc body text mistakes data for a Bench invocation. |
| FOG40 | 12 | `X=1 bench help \| touch <marker>` exits 2 and leaves the marker absent. | Bench-follow-on classifier | A prefix resolver that handles only the `env` utility misses a shell assignment. |

### Edge inventory

- The outer operator is `|`, `|&`, `&&`, `||`, `;`, or a newline.
- The redirection changes stdin, stdout, stderr, an explicit file descriptor, or two descriptors.
- The envelope spells `<`, `>`, or `&` literally or as a JSON Unicode escape.
- Quotes contain operator characters in a child argument.
- A heredoc body contains Bench text, while its outer heredoc operator remains visible.
- An environment assignment or routine command prefix precedes the Bench executable.
- The caller uses direct `bench`, repository-relative `./bin/bench.sh`, absolute `<root>/bin/bench.sh`, or live symlink alias `./kit-command`.
- The command starts in the repository root or a deeper directory.
- The command field is absent, empty, non-string, malformed, or contains a control byte.
- The wrapper-resolver library, wrapper, platform binary, or core command is unavailable.
- The core exits with an intentional refusal, a normal allow, or an analyzer error.
- Claude and Codex load the script from their shipped Bash hook groups.

**Won't handle** a Bench spelling assembled by shell expansion — a directly spelled Bench invocation remains the surviving in-scope caller.

**Won't handle** a dangling or cyclic symlink alias — a live symlink alias remains the surviving in-scope caller.

**Won't handle** wrappers deeper than one `-c` hop — the direct and one-wrapper callers remain in scope.

**Won't handle** deliberate edits to the hook, configuration, or core — the gate remains the surviving integrity authority.

**Won't handle** non-Bash tool calls — each harness's Bash PreToolUse call remains in scope.

## Ownership fences

**01-extract-shared-shell-token-stream.md writer**

- `internal/shellcommand/`
- `internal/gitguard/`

**02-refuse-bench-shell-follow-ons.md writer**

- `.bench/hooks/block-bench-follow-on.sh`
- `.claude/settings.json`
- `.codex/hooks.json`
- `bin/bench.sh`
- `cmd/bench/`
- `internal/benchguard/`
- `internal/systemtest/`
- `internal/conformance/`
- `internal/guards/`
- `internal/packagesurface/`
- `tests/canary/load-validity-metadata/`
- `tests/canary/line-routing/`
- `tests/canary/package-core-guard/`
- `projects/benchkit.md`
- `CHANGELOG.md`

**03-teach-the-complete-bench-response.md writer**

- `.bench/BENCH.md`
- `internal/anchors/registry_data.go`
- `tests/canary/workflow-guidance-anchors/`

The extraction writer publishes the token-stream contract before the enforcement writer consumes it.
The guidance writer has a disjoint fence and can ship on its own gate.

## Out of scope

- Evasion-resistant shell parsing is a separate security capability: at least 20 edits and 4 gate runs.
- Full Bench argument validation in the guard duplicates the command grammar: at least 15 edits and 3 gate runs.
- Follow-on guards for non-Bench programs are a separate policy capability: at least 8 edits and 2 gate runs.
- A new output projection is explicitly excluded by the decision source: at least 10 edits and 2 gate runs.

## Further notes

The reviewed artifact calls for the same fail-open rim as the kit's thin delegation guard.
The destructive-Git guard has a different fail-closed degraded rim.
This spec follows the explicit FT248 fail-open decision.
