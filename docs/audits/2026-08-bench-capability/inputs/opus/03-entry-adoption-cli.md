# Evidence log 03 — entry point, adoption, CLI, controls

## E8 — un-adopted repo reports all-clear (OBSERVED)

Fresh `git init` repo, one file, one commit, no `.bench/`:

```
$ bench status
bench: clean — nothing pending          [rc=0]
```

Identical output to an adopted, fully-clean repo. No route to `bench setup`.
AXI principle 5 (definitive empty states) fails at the highest-traffic surface.

## E9 — `bench setup` is excellent; the gate it scaffolds cannot go green (OBSERVED)

```
$ bench setup --plan
bench setup: plan for .../newrepo
  git repository detected -> proceeding
  AGENTS.md absent -> will be created with the managed Bench block
  CLAUDE.md absent -> will be created with the managed @AGENTS.md / @.bench/BENCH.md import lines
  .bench/gate.sh: no build system detected -> a fail-closed stub will be written (red until configured)
  projects/newrepo.md absent -> will be scaffolded as a starting profile
  platform assets (.bench/, .agents/, .claude/, .codex/) -> installed or converged via the link lifecycle

$ bench setup --yes
  ok: AGENTS.md carries the Bench managed block
  ok: CLAUDE.md carries the Bench import lines
  red: .bench/gate.sh is still the unconfigured fail-closed stub
  ok: profile present at projects/newrepo.md
  ok: repo-local bench resolvable at .bench/bin/bench.sh
  ok: setup-emitted pointer present at .agents/commands/bench-setup-repo.md
next: configure .bench/gate.sh - replace the BENCH_SENTINEL sentinel with your project's real checks
```

Plan-first, converging, self-verifying, exact next action. Best surface in the repo.

**Then, following that exact instruction:**

```
$ grep -v BENCH_SENTINEL .bench/gate.sh > g && mv g .bench/gate.sh
$ bench gate
.../.bench/bin/bench.sh: line 14: HOME: unbound variable
gate: canary inventory validation failed
gate: red
```

Cause chain, each link verified:

1. `bench setup` scaffolds **no** `.bench/gate-inputs.json`.
2. `internal/gate/subject.go:65` starts the phase environment as
   `subject.Env = []string{"PATH=" + pathEnv}` and grows it **only** from a
   manifest's `Environment` list. `loadManifest` returns `nil` when the file is
   absent → the gate command runs with `PATH` alone.
3. The scaffolded `.bench/gate.sh:25` runs `"$bench" canary "$root"` where
   `$bench` is `.bench/bin/bench.sh`.
4. `bin/bench.sh:14` is `export BENCH_HOME="${BENCH_HOME:-$HOME/.bench}"` under
   `set -euo pipefail` → unbound `$HOME` is fatal.

Minimal repro: `env -u HOME bash .bench/bin/bench.sh version` →
`line 14: HOME: unbound variable`.

**Declaring `HOME` is not sufficient:**

```
$ printf '{"schema":1,"closure":"local","environment":["BENCH_HOME","HOME"],"paths":[],"tools":["bash","git"]}' > .bench/gate-inputs.json
$ bench gate --fresh
canary fixture inventory is empty
gate: canary inventory validation failed
gate: red
```

`bench canary` exits **1** on an empty inventory (verified:
`bench canary <newrepo>` → rc=1; `bench canary <kit>` → `canary inventory ok
(233 fixture bindings)`, rc=0), and the scaffolded gate hard-fails on it. Every
new repo has an empty fixture inventory.

Neither `bench setup`'s output nor `.agents/commands/bench-setup-repo.md`
mentions `canary` or `gate-inputs.json` (verified by grep). The phase's only
coverage is §3's *"if it errors for a reason other than real failing checks, fix
the wiring before declaring done."*

Scratch repo removed after the experiment.

## E10 — no signal for a staged spec (OBSERVED)

A temporary `specs/zz-audit-probe/spec.md` was created with `Status: staged`, a
valid coverage map, and a backticked ownership fence. `bench preflight build
zz-audit-probe` then returned a full check table (proving the spec was
well-formed and discoverable):

```
checks[5]{check,verdict,detail}:
  base-current,green,""
  paths-authorized,red,"not authorized by any ownership fence: .claude/settings.json"
  rows-owned,not-applicable,""
  rows-membership,not-applicable,""
  diff-nonempty,not-applicable,""
```

`bench status` in the same tree emitted **no `specs` row**. Reading
`internal/status/status.go:546`, the `specs` signal has exactly one form:
`"%d merged spec(s) awaiting retirement"`. The full signal set is
`gate · git · intent · worktree · guards · drain · structure · decisions ·
specs · reviews · roadmap` — with `specs` meaning *retirement*, never *staged*.

Probe spec removed; `specs/` back to 34 directories.

## E11 — `bench coverage --check` bites (OBSERVED)

Four hostile inputs against the probe spec, each caught with an exact message and
exit 1:

| input | result |
|---|---|
| story 2 declared, referenced by no row, no `Not covered:` line | rc=1 — *"coverage map leaves story 2 unreferenced; add a row or a `Not covered: story 2 — <reason>` line"* |
| + `Not covered:` line added | rc=0 — *"ok: coverage map valid — 1 row(s)"* |
| row references 5 stories | rc=1 — *"row 1 references story 3, which the spec does not declare"* |
| behavior states two predicates (`;`) | rc=1 — *"row 1 behavior states more than one predicate (';' outside backticks)"* |
| duplicate row id `P1` | rc=1 — *"row 2 has a duplicate row id 'P1' (first used at row 1)"* |

**Enforcement:** prose-invoked by `/bench-write-spec`. The mechanical enforcement
exists — `internal/conformance/registry/registry.go:155` maps the
`coverage-map-validation` fixture family to the `docs-currency-workflow` check,
which calls `coverage.Command(["--check", path])` over staged specs
(`docs_workflow_checks_test.go:264`) — and that check is one of the 29 that do not
run in the gate (evidence log 02).

## E12 — `check-agent-line` deny surface (OBSERVED, 10 probes)

`.bench/lines.env` binds `BENCH_CLAUDE_{TOP,MID,CHEAP} = fable, opus, sonnet`.

| envelope `model` | rc | message |
|---|---|---|
| `opus` / `sonnet` / `fable` | 0 | allowed |
| `haiku` | 2 | *DENIED: delegation model 'haiku' is not a bound tier; harness claude binds top=fable mid=opus…* |
| `claude-opus-4-1` | 2 | DENIED — unbound |
| *(field absent)* | 2 | *DENIED: the delegation envelope has a missing or empty model field — an omitted model silently…* |
| `""` | 2 | DENIED — same |
| `OPUS` | 2 | DENIED — exact-string comparison, deliberately case-sensitive |
| `" opus "` | 2 | DENIED — no trimming, deliberately |
| `subagent_type: fork` | 0 | *WARNING: a fork delegation inherits this session's model, which no hook can…* |

Denying the **omitted** field is the notable one: it closes the silent-escalation
hole rather than documenting it.

## E13 — git guard, 43 probes (OBSERVED)

**Blocked (correct):** `git reset --hard`, `git clean -f*`, `git checkout -- .`,
`git checkout HEAD -- .`, `git restore .`, `git restore --worktree .`,
`git restore --source=HEAD .`, `git push` (all forms incl. `--delete`),
`git rebase -i`, `git branch -D`, `git update-ref -d`, `git reflog expire`,
`git commit --amend`, `git worktree remove --force`; and the same inside
`bash -c`, `sh -c`, `xargs`, and behind a `GIT_DIR=` prefix.

**Allowed (correct):** `git status`, `git log`, `git commit -m`, `git merge`,
`git switch`, `git checkout <branch>`, `git restore --staged .`,
`git worktree prune`.

**Allowed (gaps):** `git stash drop`, `git stash clear`, `git filter-branch --all`,
`git rm -rf .`, `git gc --prune=now`, `git reset --keep`, `git reset --merge`,
`git am --abort`, `git apply -R`, `git symbolic-ref`, `git revert`,
`eval 'git reset --hard'` (documented: wrapper scan is one level), `rm -rf .git`
(documented as out of scope).

The guard enumerates spellings, not effects — but the `git restore` family being
covered, with `--staged` correctly exempt, is careful work I did not expect.

## E14 — `bench worktree clean` is plan-first and refuses a dirty tree (OBSERVED)

```
$ bench worktree clean "$PWD"
worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:
  /home/mgibs/workspace/bench-audit-opus,retain,dirty,count=16 bytes=15690858 ...,
  none,"09d7ceb8…","ignored residuals require --discard-ignored"
ignored_paths[16]{path}:
  .logs/gate-*.jsonl  (15)
  dist/bench
```

Correct refusal; nothing removed. Note the plan would list `dist/bench` as an
ignored residual — and `capture/session-handoff.md:45` warns that `dist/bench`
must exist for CLI resolution. `roadmap/FT177` already records this
(*"landing guard should warn before removing a load-bearing `dist/bench`"*).

## E15 — CLI exit-code and size sweep (OBSERVED, 24 commands)

Consistent contract: **0** success, **1** domain error/red, **2** usage error.

```
bench preflight build <slug>   rc=1  164 B
bench coverage <slug>          rc=1  180 B
bench spec retire <orphan>     rc=1  182 B  "incomplete retired spec folder … retire will not auto-clean it"
bench maps                     rc=1  797 B
bench structure                rc=1  5,661 B (62 issues)
bench status                   rc=0  394 B
bench roadmap                  rc=0  2,221 B
bench learnings                rc=0  43 B    "learnings[0]{date,title}:"  — definitive empty state
bench guards                   rc=0  482 B
bench diff                     rc=0  759 B
bench worktree list            rc=0  194 B
bench canary                   rc=0  42 B
bench outline                  rc=0  24,660 B  ← the outlier
bench anchors                  rc=2  49 B    usage
bench coverage                 rc=2  63 B    usage
bench idea                     rc=2  26 B    usage
bench nosuchverb               rc=2  39 B    "bench: unknown subcommand"
```

`bench outline` honestly discloses its truncation
(`outline_meta[1]{...}: "1380","1311","69","5335","200","5135","true"` plus
`outline_skips[69]`), but offers no `--full` or paging — only path scoping, which
works well (`bench outline internal/toon` → 18 rows).

## E16 — drift measured against the prior self-assessment

`ASSESSMENT.md` (2026-08-13, four days before this audit) is a `/bench-assess`
product: six read-only mid-tier area sweeps, top-tier adversarial synthesis.

| metric | 2026-08-13 | 2026-08-17 |
|---|---|---|
| `bench structure` issues | 50 (L-O1) | **62** |
| gate capability skips | 5 ("one FIFO and three privilege skips") | **7** (`capability=6 environment=1`) |
| invalid decision maps | 1 (M-O1, `spec-build-review-gate-cadence`) | **1, same one** |
| root conformance in the gate | **not mentioned** | absent since 2026-07-05 |

The assessment enumerated the gate's skips and accounted for four of five. The
fifth — the environment-class skip that disables all 29 conformance checks — was
not identified. The gate's own output makes this easy to miss: it expands
`class=fifo` and `class=privilege` and prints the environment class as a bare
count.
