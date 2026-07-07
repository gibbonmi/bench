# First-Install Fixes (FT26)

## #1: How do the profile templates reach a fresh install?

Blocked by: —
Type: Grill

### Question
`bench init` and `/bench-setup-repo` point at `projects/<name>.md` example
profiles, but `projects/` is not in `package.json` `files[]` — on the
advertised npm/npx path the templates don't exist and the agent improvises.

### Answer
Add `projects/` to `files[]` so the three example profiles ship as-is.
Rejected: embedding template text in the Go binary — duplicates content that
already lives as files and drifts from the shipped examples.

## #2: Which stale kit docs stop shipping?

Blocked by: —
Type: Grill

### Question
npm ships the kit's own `HANDOFF.md`, `CLAUDE.md`, and `AGENTS.md`; adoption
never reads them (`link` generates consumer copies from constants), and the
stale handoff is exported to every consumer.

### Answer
Drop all three from `files[]`. They are the kit repo's working files, not
package content; `link` generates the consumer-facing equivalents. (HANDOFF.md
itself is deleted from the repo by the after-merge build ([FT25]).)

## #3: How does the committed-binary hazard get closed?

Blocked by: —
Type: Grill

### Question
`bench link` copies the running arch-specific binary to `.bench/dist/bench`
with no ignore entry — a consumer who commits it hands a different-arch
teammate a broken CLI, and the Stop-hook oracle then fails open silently.

### Answer
Link writes `.bench/dist/.gitignore` (ignoring the binary) alongside the copy,
manifest-tracked like other managed files. Rejected: arch-checking at hook
runtime — treats the symptom, adds a moving part to every hook invocation.

## #4: How does the scaffolded gate find the CLI?

Blocked by: —
Type: Grill

### Question
The scaffolded `.bench/gate.sh` calls bare `bench canary`; a machine relying
only on the `.bench/bin` local CLI (the documented "global bench optional"
story) gets a gate that cannot go green.

### Answer
The scaffold resolves the repo-local CLI first — `"$(dirname "$0")/bin/bench"`
when present and executable — falling back to `bench` on PATH. One resolver
line in the scaffold template; POSIX sh, kit shell rules apply.

## #5: What proves the fresh-install path stays green?

Blocked by: #1, #3, #4
Type: Grill

### Question
Every one of these defects is invisible in the kit repo because it never
installs itself cold. What coverage catches the class, not the instances?

### Answer
Two layers: (a) extend the package-surface contract test to pin `files[]`
contents — `projects/` present, `HANDOFF.md`/`CLAUDE.md`/`AGENTS.md` absent;
(b) a fresh-install contract test that creates a throwaway git repo in a temp
dir, runs the built binary's `link` and `init` with PATH stripped to
essentials (no global `bench`), and asserts the scaffolded gate resolves and
the profile templates are reachable. Uses the built `dist/bench` per the
runtime-contract convention. Rejected: a full `npm pack`+install smoke in the
gate — adds npm and network variance to a 30s gate for coverage layer (a)
already provides.

## Handoff

1. **Module boundaries.** `package.json` owns the surface list;
   `internal/adopt` (link.go, init.go) owns the gitignore write and scaffold
   template; `internal/contract/surface` owns the files[] pin;
   `internal/contract/runtime` owns the fresh-install smoke.
2. **Contracts.** `files[]` ships `projects/` and no kit working docs; link
   writes `.bench/dist/.gitignore`; scaffolded gate.sh prefers
   `.bench/bin/bench` over PATH; both contract tests red on regression.
3. **Deep vs thin.** The fresh-install smoke is the deep piece (temp repo,
   stripped PATH, built binary); everything else is one-line edits.
4. **Black-box assertables.** files[] JSON contents; presence of the
   gitignore in a linked temp repo; gate scaffold exit code in a PATH-stripped
   env; template file presence after init.
5. **Gate attachment.** All four seams land in gate-observed contract tests;
   nothing is manual-verify.
6. **Hostile-input owners.** The smoke test owns the no-PATH environment; link
   owns pre-existing `.bench/dist/.gitignore` (idempotent re-link must not
   duplicate or clobber a user's wider ignore).
7. **Uncertainty flags.** Whether the manifest format needs a version bump for
   the new managed file — the implementer checks how link fingerprints
   additions on re-link.
8. **Rejected alternatives.** Binary-embedded templates; hook-time arch check;
   npm-pack smoke in the gate.
9. **Domain watch-outs.** postinstall must never fail an install — the smoke
   runs in the kit's gate, not in the consumer's install path.

Dependency order: n/a — single spec.
