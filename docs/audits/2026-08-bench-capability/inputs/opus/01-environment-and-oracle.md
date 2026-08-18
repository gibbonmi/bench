# Evidence log 01 — environment, hooks, oracle

## Environment (Phase 0)

```
AUDIT MODE:               in-repository dogfood (disposable worktree), cold start
HARNESS:                  Claude Code
MODEL:                    claude-opus-5 (Opus 5)
EFFORT:                   xhigh
BENCH REPOSITORY PATH:    /home/mgibs/workspace/bench-audit-opus
                          (linked worktree; main tree /home/mgibs/workspace/bench)
BENCH COMMIT SHA:         58d966e2f92f7f37eba07b6215e8eef45371b72d
BENCH BRANCH:             audit/opus
WORKTREE AT START:        " M .claude/settings.json" + "?? audit/"
WORKTREE AT END:          same + "?? audit/opus-xhigh/" (report), "?? dist/" (gitignored build)
REMOTE:                   https://github.com/gibbonmi/bench.git
SIBLING WORKTREES:        /home/mgibs/workspace/bench (main), bench-audit-sol (audit/sol)
AUTO-LOADED INSTRUCTIONS: CLAUDE.md -> @AGENTS.md + @.bench/BENCH.md (2,342 words total)
AUTO-LOADED SKILLS:       16 bench-craft-* + prototype (.claude/skills symlinks)
                          + 6 model-invocable phase commands
                          (bench-debug, bench-final-check, bench-implement-spec,
                           bench-review-implementation, bench-shape-idea, bench-write-spec)
                          5 phase commands hidden by disable-model-invocation:
                          (bench-assess, bench-deepen, bench-setup-repo,
                           bench-update-kit, bench-what-next)
                          skill description block: 1,259 words
AVAILABLE TOOLS:          Bash, Read, Write, Edit, Skill, Agent (unused), Artifact,
                          deferred: WebFetch/WebSearch (unused — all upstream sources
                          were available as local pinned checkouts)
OUTPUT DIRECTORY:         audit/opus-xhigh/ (per operator instruction; overrides the
                          prompt's "outside the repo" suggestion)
```

Upstream sources inspected as **local pinned checkouts** (no network used):

| repo | commit | date | origin |
|---|---|---|---|
| mattpocock/skills | `d574778f94cf620fcc8ce741584093bc650a61d3` | 2026-07-08 | github.com/mattpocock/skills |
| kunchenguid/axi | `6908df208f6e8b3c3c2c5bf009a081f99191c2de` | 2026-06-26 | github.com/kunchenguid/axi |
| kunchenguid/no-mistakes | `87b2abf78888d8af738903415f5f4b58e61e2396` | 2026-06-27 | github.com/kunchenguid/no-mistakes |
| kunchenguid/treehouse | `68fa3d2556542add76bf80255787b8625a5041a6` | 2026-06-24 | github.com/kunchenguid/treehouse |
| kunchenguid/firstmate | `08691310f92d99b178d058023913b772397709ef` | 2026-06-27 | github.com/kunchenguid/firstmate |
| kunchenguid/lavish-axi | `e4556721bfc52473f509f6670ff9067939b19139` | 2026-06-26 | github.com/kunchenguid/lavish-axi |
| kunchenguid/tasks-axi | `6d200480f592bfbb050c25fe5dee41f30ddc3853` | 2026-06-23 | github.com/kunchenguid/tasks-axi |
| kunchenguid/programbench-bench | `2d1528dce843b501ee8f582d8a463b96ccef4d0c` | 2026-06-24 | github.com/kunchenguid/programbench-bench |

These are the revisions on disk; they may lag current upstream HEAD. Every upstream
claim in this report is pinned to the commit above, not to "current upstream".

## E1 — cold start: the CLI does not exist (OBSERVED)

At session open, `dist/bench`, `node_modules/@redbench/*`, and `$BENCH_HOME/cache`
were all absent in **both** the audit worktree and the main tree. The global
`/home/mgibs/.local/bin/bench` shim points at `/home/mgibs/workspace/bench/bin/bench.sh`,
which resolves the same missing binary. Every `bench` verb exited 127.

Recovery: `go build -o dist/bench ./cmd/bench` — **1.4 s**. Nothing told me to do it;
`capture/session-handoff.md` line 45 mentions it, but the handoff is not auto-loaded.

## E2 — degraded-core hook posture is inverted (OBSERVED)

With `dist/bench` moved aside, each hook was fed a representative envelope:

| hook | boundary | degraded behavior | exit |
|---|---|---|---|
| `session-start.sh` | SessionStart | dies after one line; ambient dashboard absent | 127 |
| `stop.sh` | Stop | **fail-OPEN** — "allowing this stop without a gate verdict" | 0 |
| `check-agent-line.sh` | PreToolUse:Agent | **fail-OPEN** — "allowing delegation" | 0 |
| `block-dangerous-git.sh` | PreToolUse:Bash | **fail-CLOSED** on substring `git` | 2 |

The fail-open rims are documented and reasoned (a missing oracle must not forge a
green). The Stop rim is largely moot in practice: `BENCH_SHIFT=1` is only set by
`bench shift`, which itself needs the core. The `check-agent-line` rim is not moot.

The fail-closed rim matched on the literal substring `git` anywhere in the raw
PreToolUse envelope. Demonstrated false positives that cost this audit real turns:

```
cat .gitignore                  -> BLOCKED
ls /path/to/repo/.git           -> BLOCKED
echo legitimate                 -> BLOCKED   (le-GIT-imate)
```

and a demonstrated false negative in the same degraded state:

```
rm -rf .g""it                   -> ALLOWED
```

## E3 — live git guard deny surface (OBSERVED)

With the core present, 43 probes through the real `bench guard-git`:

BLOCKED (correct): `git reset --hard`, `git clean -f*`, `git checkout -- .`,
`git checkout HEAD -- .`, `git restore .`, `git restore --worktree .`,
`git restore --source=HEAD .`, `git push` (all forms), `git rebase -i`,
`git branch -D`, `git update-ref -d`, `git reflog expire`, `git commit --amend`,
`git worktree remove --force`, and the same inside `bash -c`, `sh -c`,
`xargs`, and with a `GIT_DIR=` prefix.

ALLOWED (correct): `git status`, `git log`, `git commit -m`, `git merge`,
`git switch`, `git checkout <branch>`, `git restore --staged .`, `git worktree prune`.

ALLOWED (gaps — all plausible honest mistakes):
- `git stash drop` / `git stash clear` — destroys stashed work
- `git filter-branch --all` — history rewrite, while `rebase -i` is blocked
- `git rm -rf .` — deletes tracked files
- `git gc --prune=now`, `git reset --keep`, `git reset --merge`
- `eval 'git reset --hard'` — documented: wrapper scan is one level, `eval` not covered
- `rm -rf .git` — documented as out of scope

The guard enumerates **spellings**, not **effects**. It is nonetheless materially
better than I assumed before probing: `git restore` (the modern synonym for the
blocked `git checkout --`) is covered, including `--source=`, and `--staged` is
correctly allowed.

## E4 — the gate is real, fast, and correctly tree-keyed (OBSERVED)

```
bench gate               -> green, 64 s, 5,917 bytes of output, 95 lines
bench gate (re-run)      -> green, 0 s ("fresh verdict reused for this tree")
phases: gofmt, vet, test, race, system, shellcheck
slowest: internal/worktree 50.2 s of the 64 s
```

Cache invalidation verified in both directions:

| probe | tree hash | gate |
|---|---|---|
| baseline | `cb77222` | green (reused) |
| + untracked `zz_audit_probe.txt` | changes | full re-run |
| + comment appended to tracked `internal/toon/toon.go` | changes | full re-run |
| both reverted | `cb77222` | green (reused) |
| + `dist/bench` removed | `cb77222` (unchanged — gitignored) | n/a |

`git.TreeHash` seeds a throwaway index from HEAD, `git add -A` (respects
`.gitignore`), `write-tree`. Untracked-unignored files are inside the identity.
This is a genuinely good design: one content-addressed fact, no claims database.

Red detection verified: appending `func   auditProbeBadFmt( ) {}` to
`internal/toon/toon.go` produced `gate: red`.

## E5 — two verdict stores disagree; `bench status` contradicts the oracle (OBSERVED, defect)

`bench gate` reuses a green from a **content-addressed per-subject evidence store**
(`$GIT_COMMON_DIR/bench-gate-evidence/<sha256(tree,oracle)>`, read by
`inspectEvidence` -> `evidencePath`). It does **not** refresh the **per-worktree
last-executed record** (`$GIT_DIR/bench-last-gate`, written by `durableReplace`).
`bench status` reads only the latter.

Reproduction (each line is a real observation from one run):

```
state0  tree=cb77222   bench-last-gate: tree=cb77222 status=green
state1  append gofmt violation -> tree=626eef9
        bench gate            -> "gate: red"
        bench-last-gate: tree=626eef9 status=red
state2  revert                 -> tree=cb77222
        bench gate            -> "gate: green (fresh verdict reused for this tree)"
        bench-last-gate: tree=626eef9 status=red        <-- not refreshed
```

Consumer-visible consequence, same tree, same second:

```
$ bench gate
gate: green (fresh verdict reused for this tree)

$ bench status
▶ fix before commit  (gate)
  gate       red                            → fix before commit
```

The **ambient dashboard — the SessionStart surface designed to orient a cold
session — reports RED as its headline while the oracle reports GREEN.**

A second, milder manifestation was the first symptom I hit: with `bench-last-gate`
holding an older executed green (`7f0d995`) and the current tree green-by-reuse
(`cb77222`), `bench status` printed
`gate stale (gated tree 7f0d995, work tree cb77222) → re-run the gate`,
prescribing a 64 s re-run of a gate that had just certified the exact tree.

Failure direction: `status` over-reports problems; it did **not** produce a false
green (it compares cached vs current tree and degrades to "stale"). So this is a
false-alarm / wasted-work / trust-erosion defect, not a safety hole. It is still a
direct violation of the repo's own "one source per fact" standard: two readers of
the gate verdict, disagreeing.

`bench gate --fresh` reconciles both stores.

Minor: `$GIT_COMMON_DIR/bench-gate-evidence/` has no pruning path (3 files, 209 bytes
each, at audit time).
