# bench outline

Status: implemented

## Problem

An agent picking up a large or polyglot linked repo has to read whole files to
find where a seam lives — a function, a type, a heading. That is expensive and
noisy: a shell-function outline of this repo is ~1.6k tokens against ~67k of full
source, a ~40x compression. Today there is no cheap way to ask "where are the
declarations" without paying the full read.

## Solution

`bench outline [path]` — an on-demand repo seam map. It walks the tracked files
(optionally scoped to a path), runs a hand-rolled per-language pattern scan, and
emits an AXI-conformant TOON table of `file,line,kind,name` rows so an agent can
locate a seam by name and jump straight to `file:line`. It is regenerated on every
call; nothing is committed to the tree. It **locates** candidate seams; it does not
identify which are the project's blessed seams — `projects/<name>.md` owns that.

## User stories

1. As an agent, I want `bench outline` with no argument to emit every tracked
   declaration in the repo as a TOON `outline[N]{file,line,kind,name}:` table, so
   that I can find any seam without reading whole files. This story also carries the
   AXI posture: a definitive empty table when nothing matches, a structured
   not-in-a-repo error on stdout with exit 1, and usage on stdout with exit 2 for an
   unknown flag or a second positional argument — routed through `cmd/bench/main.go`,
   `bin/bench.sh`, and the `.bench/BENCH.md` inventory like every other query
   surface. Line: claude-sonnet-5 / low. This is mechanical wiring at the known
   AXI-command seam that `learnings`/`maps`/`structure` already established.

2. As an agent, I want each of the v1 languages recognized by its own declaration
   forms — shell function definitions, Go `func`/method/`type` declarations,
   Markdown headings, Python `def`/`class`, and JS/TS function/class/exported-const
   arrow declarations — so that the outline is useful across a polyglot repo. Line:
   claude-opus-4-8 / medium. The per-language regex forms are the genuinely
   uncertain part of this build and the place a wrong reading of "each parser" would
   silently under-deliver.

3. As an agent, I want an optional `path` argument to scope the outline to one file
   or one directory (default = repo root), resolved from my current working
   directory even when it is deeper than the root, so that I can index just the
   region I care about. Only tracked files are walked and binaries are skipped, with
   output paths always repo-root-relative. Line: claude-opus-4-8 / medium. The
   cwd-to-root pathspec resolution and the tracked-only/binary-skip walk are the
   correctness-sensitive edges here.

4. As an agent, I want one poisoned file — a control byte in a heading, or in a
   tracked path — to drop only its own row rather than fail the whole command, and I
   want a stable file-then-line ordering, so that the index is dependable on a
   hostile tree. Line: claude-opus-4-8 / medium. Control-byte safety at the TOON
   boundary is the subtle correctness point: the naive `toon.Table` call refuses the
   whole table on one bad byte.

5. As an agent reading the help and the output, I want the surface to say plainly
   that it *locates* seams (`file:line`) and does not identify which are the
   project's blessed seams, so that I do not mistake the index for the seam list in
   `projects/<name>.md`. This wording lands in the `-h`/help line, the `.bench/BENCH.md`
   CLI-inventory row (under "Oracle and diagnostics"), and the top-level bench.sh help
   block. Line: claude-sonnet-5 / low. It is a fixed one-line promise, decided in the
   handoff, transcribed into three surfaces.

## Implementation decisions

- **New package `internal/outline`**, one `Command([]string) (string, int)` entry,
  dispatched from the `commands` map in `cmd/bench/main.go` and routed by a
  `outline) route_porcelain "$@" ;;` case plus a help line in `bin/bench.sh`. A
  single `.bench/BENCH.md` CLI-inventory line under "Oracle and diagnostics". No new
  mechanism — this grows the existing router by a name, as every prior slice did.

- **One per-language pattern table, one source of the fact.** A single table maps a
  file extension to its language's ordered declaration patterns (each pattern a
  compiled regex plus the `kind` it emits and the capture group holding the `name`).
  Adding a language is adding a table entry — never a new code path. The enumerated
  v1 set and their `kind` values:
  - shell (`.sh`): `name() { … }` and `function name` → kind `function`.
  - Go (`.go`): `func Name`, `func (recv) Name` → kind `func`; `type Name` → kind `type`.
  - Markdown (`.md`): ATX headings `#`…`######` → kind `heading` (the prose anchors).
  - Python (`.py`): `def name` → kind `def`; `class Name` → kind `class`.
  - JS/TS (`.js`,`.jsx`,`.ts`,`.tsx`): `function name` → kind `function`; `class Name`
    → kind `class`; `export const name = (…) =>` → kind `const`.

- **Pure parser seam.** `Symbols(path string, content []byte) []Symbol` (Symbol =
  `{Line int; Kind, Name string}`) selects the table entry by extension, scans line
  by line, and returns rows in ascending line order. A path whose extension has no
  table entry yields no rows. A final line lacking a trailing newline is still a
  line and is still scanned. This is a line-regex indexer: it does **not** track
  comments, string literals, or Markdown code fences, so a commented-out or fenced
  declaration is indexed as a benign candidate — consistent with "locates, not
  identifies," where the agent confirms by reading the line.

- **Tracked-file walk via git, NUL-framed.** The file list comes from
  `git ls-files -z` (reusing `internal/git`), split on NUL so a path with spaces or
  an embedded newline survives whole and git never C-quotes it. A `path` argument is
  resolved from the process cwd to an absolute path, then to a root-relative literal
  pathspec (`:(literal)…`) passed after `--`, so a glob character in the path is
  matched literally rather than expanded, and a directory argument scopes to the
  files beneath it. Output paths are repo-root-relative regardless of cwd. A file is
  read and skipped as binary when its content contains a NUL byte; a file with no
  table entry contributes no rows.

- **Columns are `file,line,kind,name` — no signature column in v1.** The name plus
  `file:line` is the anchor; a signature adds table width and is the most likely
  carrier of a control byte the TOON boundary would reject. Minimal schema per AXI.

- **Control-byte safety at the row, not the table.** Before emission, a row whose
  `file` or `name` cell contains a byte spec-TOON cannot represent (a control
  character other than tab/newline/return) is dropped, so one poisoned heading or
  tracked path removes only itself and the command still exits 0 with the rest of the
  index. Rows are sorted by file (git's ls-files order) then ascending line, giving a
  deterministic table. Emission is `toon.Table("outline", []string{"file","line","kind","name"}, rows)`.

- **On-demand, stateless.** Every call regenerates; nothing is written to the tree,
  no scratch state, no lease. Re-running on an unchanged tree yields byte-identical
  output.

## Testing decisions

- **What a good test is here:** exercise external behavior at a seam — the parser's
  `(line,kind,name)` rows for a fixture buffer, and the command's TOON/exit contract
  in a throwaway repo — never the internals of the regex table.
- **Two seams, matching the repo's precedent:** a fast in-package unit test over the
  pure `Symbols` parser (exhaustive across the five languages, where the pattern
  table's correctness lives), and the AXI black-box contract driving the built
  `bench outline` in a throwaway fixture repo (where the walk, path scoping, TOON
  shape, empty state, exit codes, and hostile inputs live). Prior art:
  `internal/learnings` + `internal/maps` (package unit tests) and
  `internal/contract/axi/*` (built-binary fixture contracts).
- **Gate command:** `.bench/gate.sh` (the project gate). New source adds a
  `TestRootConformance` package to build/vet, and the AXI fragment adds an outline
  contract file.

### Seam diagram

Seam A — the pure per-language parser (unit):

    trigger: outline_test.go drives Symbols() with a fixture buffer per language
        │
        ▼
    (path ext, file bytes) ──▶ [ langTable[ext] → per-line regex scan ] ──▶ []Symbol{line,kind,name}
    (final line, no \n)    ──▶ [                                       ]
                                   ◀ tests attach here: feed one buffer per language,
                                     assert the exact (line,kind,name) rows

Seam B — the command / AXI contract (built binary, black box):

    trigger: `bench outline [path]` in a throwaway fixture repo
        │
        ▼
    argv[path] ──▶ [ git ls-files -z → dispatch by ext → Symbols → row-filter → toon.Table ] ──▶ TOON stdout + exit
    (no arg)   ──▶ [                                                                         ]
                       ◀ tests attach here: f.Bench("outline", path); assert
                         header / rows / empty state / exit code

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | no-arg outline emits an `outline[N]{file,line,kind,name}:` header and one row per tracked declaration | B: `f.Bench("outline")` in a fixture repo | contract test asserting the header and a known row fails against an empty-stub command | an always-empty command emits `outline[0]{…}` and the known-row assertion goes red |
| 1 | nothing to index emits the definitive empty table `outline[0]{file,line,kind,name}:` at exit 0 | B: fixture repo with only unparseable files | assertion on the exact empty header fails if the command errors or prints nothing | pins the AXI empty-state posture against a silent or error exit |
| 1 | outside a repo prints a structured not-in-repo error on stdout, exit 1 | B: `f.Bench("outline")` run outside any repo | exit-code and `error:`-on-stdout assertion fails | catches an error routed to stderr or a wrong exit code, invisible to an agent |
| 1 | unknown flag or a second positional argument prints usage on stdout, exit 2 | B: `f.Bench("outline","--bogus")` and `f.Bench("outline","a","b")` | assertion on exit 2 and `usage` substring fails | a command that ignores bad args exits 0 and the assertion goes red |
| 2 | a Go `func`, method, and `type` are indexed with kinds func/func/type | A: `Symbols("x.go", …)` | table assertion on the three Go rows fails against a parser with no Go entry | the degenerate empty parser returns no rows; the Go rows go red |
| 2 | a shell `name()` and `function name` are indexed with kind function | A: `Symbols("x.sh", …)` | shell-row assertion fails against a parser missing the shell entry | enumerates the shell language so a one-language port cannot pass |
| 2 | a Markdown ATX heading is indexed with kind heading | A: `Symbols("x.md", …)` | heading-row assertion fails against a parser missing the Markdown entry | enumerates the prose-anchor language explicitly |
| 2 | a Python `def` and `class` are indexed with kinds def/class | A: `Symbols("x.py", …)` | Python-row assertion fails against a parser missing the Python entry | enumerates the Python language explicitly |
| 2 | a JS/TS function, class, and exported const arrow are indexed with kinds function/class/const | A: `Symbols("x.ts", …)` | JS/TS-row assertion fails against a parser missing the JS/TS entry | enumerates the JS/TS language and its three forms explicitly |
| 3 | a file-path argument scopes the outline to just that file's rows | B: `f.Bench("outline","pkg/one.go")` in a two-file repo | assertion that only `one.go` rows appear fails against a cwd-blind whole-repo scan | a command that ignores the path argument emits the other file's rows and goes red |
| 3 | a directory-path argument scopes the outline to files beneath it | B: `f.Bench("outline","pkg")` | assertion that only `pkg/**` rows appear fails against an unscoped scan | pins directory scoping distinct from whole-repo |
| 3 | invocation from a subdirectory deeper than the root indexes the whole repo with root-relative paths | B: `runBenchInDir(sub, "outline")` | assertion on root-relative rows fails against a command that assumes cwd == root | catches the cwd-deeper-than-root hostile class |
| 3 | an untracked file is not indexed; a binary file is skipped | B: fixture repo with an untracked source file and a NUL-containing tracked file | assertion that neither contributes a row fails if the walk reads the worktree or does not NUL-test | catches a walk that ignores git tracking or emits binary noise |
| 4 | a control byte in a heading (or in a tracked path) drops only that row; the command still exits 0 with the rest | B: fixture repo with one ESC-bearing heading plus clean declarations | assertion `exit 0 && clean rows present && poisoned row absent` fails against a naive `toon.Table` call | a naive emit refuses the whole table on one bad byte, so the whole command exits 1 and the assertion goes red |
| 4 | rows are ordered by file then ascending line | B: fixture repo with out-of-order declarations across two files | exact-order row assertion fails against an unsorted emit | pins deterministic ordering so agents and the contract can match rows |
| 5 | the help text and the output-surface promise say the tool locates seams and does not identify blessed seams | B: `f.Bench("outline","-h")` asserts the expectation-setting clause | substring assertion on the "locates … not … blessed" wording fails against generic help | pins the decided expectation-setting so the index is not mistaken for the seam list |
| edge (of 3) | a `path` argument containing spaces or a glob character is matched literally | B: fixture with a spaced/globbed dir, `f.Bench("outline","<that path>")` | assertion that the intended files are scoped fails if the pathspec expands the glob or splits on the space | covers the spaces/glob hostile class end-to-end through the shell router |
| edge (of 2) | a declaration on a final line lacking a trailing newline is still indexed | A: `Symbols` on a buffer whose last line has no `\n` | assertion the last symbol is present fails against a parser that drops the newline-less tail | covers the no-trailing-newline hostile class |
| edge (of 3) | a present-but-empty tracked file yields zero rows (distinct from an absent path) | B: fixture with an empty tracked `.go` file | assertion of zero rows for that file, exit 0, fails if the empty read errors | covers the absent-vs-empty hostile class |
| edge (of 1) | git unavailable / `ls-files` failure prints a loud structured error, exit 1 | B: fixture with a broken index (corrupt `.git/index`) | assertion on the git-failure error line and exit 1 fails against a silent empty table | covers the required-tool-missing / git-failure hostile class |
| edge (of 1) | re-running on an unchanged tree yields byte-identical output | B: two `f.Bench("outline")` runs, compared | byte-equality assertion fails against a non-deterministic emit | covers the re-run idempotency class |

### Edge inventory

The shell-CLI hostile-input checklist from `projects/benchkit.md`, each landed as a
row above or a **Won't handle** line here:

- paths/dirs with spaces or glob characters — **row** (edge of 3).
- control bytes (ESC/BEL) in git-sourced text (paths) and in file content near
  matches — **row** (story 4): dropped at the row, command stays green.
- hand-edited file whose last line lacks a trailing newline — **row** (edge of 2).
- absent file vs present-but-empty file — **row** (edge of 3): empty tracked file →
  zero rows; an untracked/absent path → definitive empty table, exit 0.
- unquoted multi-word arguments (`$*` vs `$1`) — **row** (edge of 3): `route_porcelain`
  forwards `"$@"`, so a quoted spaced path reaches argv whole; the spaces row proves it
  end-to-end.
- required tool missing / git failure — **row** (edge of 1): loud structured error, exit 1.
- invocation through a symlink rather than the real path — **Won't handle**: outline
  resolves the root via `git rev-parse` and reads no `argv[0]`; the shared binary-routing
  contract already covers symlinked invocation, so re-testing it here duplicates coverage.
- invocation through every shipped surface (real CLI, by-path CLI, hooks, adapters)
  — **Won't handle**: outline is a plain query command; the multi-surface routing
  contract that guarantees each surface reaches the same routed implementation is a
  generic property already gate-tested, not re-proven per command. The one obligation
  here — being *in* the router — is story 1.
- interrupt (SIGINT) mid-run — **Won't handle**: a single-pass read-only command that
  writes no scratch state, lease, or worktree; an interrupt leaves nothing to clean up.
- re-run idempotency — **row** (edge of 1): stateless, byte-identical on an unchanged tree.
- cwd deeper than the repo root — **row** (story 3): resolved cwd-to-root, output
  root-relative.

Parser false positives (a declaration inside a comment, string, or Markdown code
fence) — **Won't handle** in v1: this is a line-regex indexer, not a language parser;
a false-positive row still points at a real line the agent reads to confirm, which is
exactly the "locates, not identifies" contract. Comment/string/fence awareness is a
separate future capability (Out of scope).

## Out of scope

- **A signature/preview column** — a separate presentation capability layered on the
  same rows, deferred for its control-byte width cost; the name plus `file:line` is
  the v1 anchor. ~4 edits, 2 gate runs.
- **More languages** (Rust, Java, Ruby, C/C++, …) — each is one table entry plus its
  parser rows, a separate additive capability; add on demonstrated demand. ~2 edits,
  1 gate run per language.
- **`--json` (or any second output format)** — a distinct output-contract capability;
  AXI TOON is the one v1 surface. ~3 edits, 1 gate run.
- **A committed/cached outline artifact** — deliberately excluded, not deferred: a
  committed index duplicates knowledge and drifts; on-demand is a closed decision.
- **Seam *identification*** (marking which indexed symbols are the project's blessed
  seams) — a separate capability owned by `projects/<name>.md`; outline only locates.
- **Comment/string/fence awareness** in the parser — a real language-parsing
  capability, a separate future spec; v1 accepts benign false positives. ~5 edits, 2
  gate runs.
