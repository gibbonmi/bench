# FT303 CLI assessment

Subject: `730695ce` (the `main` tip on 2026-09-06)

This is the ticket #7 research object for `decisions/software-factory.md`. It
records the application inventory, the attributed log sources, the grouped
calls, the current verbs, and the ranked candidates. The reviewer selects the
improvement in ticket #10. This asset selects nothing.

The window starts at the factory decision, 2026-09-04 00:00 UTC, and ends at
the extraction, 2026-09-06 08:45 UTC. The git-ignored machine-local inventory
at `research/ft303-cli-assessment/` holds the extraction scripts and the raw
counts.

## Application inventory

| application | state | evidence |
| --- | --- | --- |
| Claude Code CLI 2.1.263 | installed | `~/.local/bin/claude` |
| Codex CLI 0.153.4 | installed | `~/.local/bin/codex` |
| OpenCode | configuration only | `~/.config/opencode` exists; no binary on `PATH` |
| ChatGPT desktop | absent | no Linux configuration; no Windows AppData entry for the host user |
| Claude desktop | absent | no Linux configuration; no Windows AppData entry for the host user |

An installed application proves no command requirement. The inventory bounds
which log sources can exist on this machine.

## Log sources

| source | path | window content | state |
| --- | --- | --- | --- |
| Claude Code transcripts | `~/.claude/projects/-home-mgibs-workspace-bench/` | 19 sessions, 109 delegate transcripts under `<session>/subagents/`, 5681 Bash calls | read |
| Codex rollouts | `~/.codex/sessions/2026/09/` | 17 sessions with the repository or a pool worktree as `cwd`, 692 exec calls | read |
| Bench census | `~/.bench/census/bench-2826441890/` | no record; each release deletes its record | absent |
| Bench seam record | `~/.bench/otel/bench-2826441890/traces.jsonl` | 72,320 spans | counted, not mined |
| Gate logs | `.logs/gate-*.jsonl` | 20 runs from 2026-09-05T19:34Z | counted, not mined |
| Retired retros | the parents of the four drain commits | four census heads and their asks | read |
| Claude prompt history | `~/.claude/history.jsonl` | prompts, no calls | not mined |
| Shell history | `~/.bash_history` | the user's own calls | not mined |

The census is the one FT243 source that is absent for the window. Four retro
heads survive it: structural-refactor-pass with 7 at the landing and 31 across
ten worktrees, git-admin-readers with 34, ledger-settle-policy with 23, and
pin-removal with 3. The transcripts replace the census for the window. They
reach more calls, because the census counts only a call that names the pool
path.

## Method

The extraction reads each Bash `tool_use` block from the Claude transcripts
and each `exec` tool call from the Codex rollouts. A Codex call carries one or
more `cmd` strings inside a script, and each string counts as one call. The
verb head obeys the glossary rule. Prefix resolution steps over `env`,
`timeout`, `xargs`, and a leading assignment, and `git` and `bench` keep their
subcommand. A census-defined raw call names the pool path with a head other
than `bench`.

The counts are lower bounds. The extraction does not reconstruct a hook-owned
call, a child process, or a variable-held executable. A Read, Edit, or Write
tool call is never a call. A raw-call total decides nothing. The ranking below
rests on a repeated pattern, a confirmed verb gap, and a cost per instance.

## Current verbs

The window holds 6373 calls in 36 sessions: 4642 with a Bench head and 1731
with another head. Delegates made 3414 calls, and 2120 calls chained two or
more segments.

| verb | direct segments |
| --- | --- |
| `bench worktree exec` | 3748 |
| `bench worktree create` | 97 |
| `bench worktree merge` | 59 |
| `bench worktree land` | 48 |
| `bench worktree path` | 45 |
| `bench handoff` | 40 |
| `bench outline` | 37 |
| `bench worktree clean` | 31 |
| `bench status` | 28 |
| `bench structure` | 23 |
| `bench worktree list` | 22 |
| `bench gate-prose` | 22 |
| `bench roadmap` | 20 |
| `bench consumers` | 19 |
| `bench commit` | 15 |
| `bench test` | 4 |
| `bench gate` | 2 |

`bench worktree exec` carries 59 percent of all calls. Its child heads lead
with `rg` 606, `go test` 428, `bash` 347, `cp` 208, `sh` 189, and
`git status` 167. Then come `bench commit` 164, `sed` 150,
`bench gate-prose` 135, `git rev-parse` 129, `cmp` 122, and `git diff` 98. The Bench verbs that run as an exec child are
`bench commit` 164, `bench gate-prose` 135, `bench preflight` 60,
`bench structure` 55, `bench coverage` 50, and `bench test` 42. The gate group
below holds 239 calls, almost all as an exec child.

## Grouped calls

| group | calls | by delegates | sessions | gap |
| --- | --- | --- | --- | --- |
| search: `rg`, `grep`, `find` | 1623 | 922 | 35 | none |
| read: `cat`, `sed -n`, `head`, `tail`, `ls`, `wc` | 1490 | 745 | 29 | none; 226 exec children read a file that the Read tool can open |
| git read: `log`, `status`, `diff`, `show`, `rev-parse` | 881 | 307 | 34 | none; FT304 owns the projection |
| focused test: `go test`, `bench test` | 710 | 570 | 17 | `go test` 619 against `bench test` 92; feeds FT290 |
| scripting: `python3`, `awk` | 376 | 141 | 26 | none; 155 of 290 python calls edit or scan Markdown |
| mutation probe: copy aside, mutate, focused test, `cmp`, restore | 334 | 189 | 15 | no verb; FT168 and FT98 own the ask |
| build: `go build`, `go vet`, `gofmt` | 329 | 291 | 11 | none; the lane runs them |
| gate | 239 | 75 | 21 | none |
| `sed -i` edits | 64 | 10 | 9 | none; the Edit tool opens the worktree path |
| git write | 59 | 27 | 11 | none; `bench commit` and the landing own the writes |

The 198 census-defined raw calls spread over 18 sessions. Their heads are `cd`
27, `rg` 26, `python3` 17, `cp` 16, `ls` 11, `sed` 8, and `grep` 8.

## Candidates

### 1. A mutation-probe verb

A copy-aside probe ran 102 complete sequences in 10 sessions, on 56 distinct
Go files, with 120 copy-aside calls. A sequence costs a median of 3 calls and
a 90th percentile of 5. The three days hold 57, 62, and 1 sequences. The census
counts a probe that names the pool path as raw shell, so the signal cannot
tell a probe from a leak.

The current grammar has no probe verb. FT168 owns
`bench probe` at `Next: ticket`. FT98 owns the preserve-and-restore face and
states that it needs a new `bench worktree` subcommand, a declared seam.

### 2. A production-or-test projection on `bench consumers` and `bench outline`

`bench consumers` ran 42 times in 9 sessions, 24 of them in the FT302 survey
session. Eight calls used `--full`, and eight used `--changed`. `bench outline`
ran 38 times in 3 sessions, 29 with `--full`.

Seven calls within two of a
consumers or outline call grepped test names. The FT302 survey delegate
reported one 91.7 KB overflow at 60 percent test symbols. The CLI contract
forbids a pipe, so the projection must be CLI-owned. Only the FT303
occurrences own this ask.

### 3. A `bench structure --path <prefix>` filter

`bench structure` ran 97 times: 44 with `--growth` from the lane, 8 with
`--since`, and 49 bare. A bare call prints 111 rows. Three bare calls piped
the answer into `rg`, against the CLI contract. Only the FT303 occurrence owns
this ask.

### 4. Guidance repairs on the worktree surface

`bench worktree path` ran 59 times in 12 sessions, and the 198 raw calls above
follow it. A one-line note that the path serves the file tools only is a
guidance repair with no verb. Two further asks are served already. The shell
`time` prefix wrapped 99 exec calls in 4 sessions, and a heredoc fed 362 exec
calls. Sixteen exec children referenced `PWD`, which is the one unserved ask.
FT254 owns exec comfort.

### 5. An inbox-emptying verb for the drain

One drain of three in the window met the harness classifier's refusal of a
shell truncation. The evidence is one event.

## Findings outside the selection

- `~/.bench/otel/` holds 25,015 record directories and 309 MB that test
  binaries wrote, 1,000 to 5,700 per day since 2026-08-30. The service name is
  `authorization.test`, and `benchhome` falls back to `~/.bench` when
  `BENCH_HOME` is unset. The idea is parked in `capture/IDEAS.md` on
  2026-09-06.
- The census cannot serve a window after a landing, because the release
  deletes the record. The transcripts are the durable source, and FT243 can
  name them.

## Ranking

1. The probe verb: the largest raw group with a verb gap, 3 calls per instance, and a census signal that mislabels it.
2. The production filter: a weekly pattern at the deepening cadence, one overflow, and seven follow-up greps.
3. The structure path filter: 49 whole-census reads and three contract deviations.
4. The path note: one guidance line.
5. The inbox verb: one event.

The assessment recommends candidate 1 as the one improvement. Candidate 4 is a
guidance repair. The reviewer decides in #10 whether it rides on the same
landing.
