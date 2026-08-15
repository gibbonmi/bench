# Review — skills-index-reader

Frozen base `5b41322a` · reviewed tip `2ca77cb2` · three axes run read-only through
`codex exec -m gpt-5.6-sol` at high effort, one fresh context each.

Raw findings: 23 (Standards 5, Spec 6, Coverage 12). De-duplicated repair targets: 17.

Each finding is classified **diff-owned** (introduced by this diff) or **inherited**
(behavior byte-identical to the pre-collapse implementation, faithfully moved). This is a
behavior-preserving refactor, so inherited edges are pre-existing gaps the move exposed,
not regressions it caused.

## Standards

5 findings, all hard violations. Worst: an invocation advertised as invalid mutates a
tracked file.

1. **`--check --write` is accepted and writes** — diff-owned. `usage.Parse` treats the two
   as independent flags and `Command` then selects `--write`; both help surfaces advertise
   them as exclusive. AGENTS.md requires an enforcement and its advertisement to share one
   source. `internal/skillsindex/command.go:15,37`, `bin/bench.sh:343`. Disposition:
   `ask-user` — exclusivity may need the shared `usage.Grammar` seam widened.
2. **Marker boundaries parsed twice** — diff-owned. `markerBlock` and `replaceMarkerBlock`
   each scan for the marker pair independently and diverge on a duplicated block.
   `skillsindex.go:274,290`. Disposition: `auto-fix` — derive the span once.
3. **Index-line format encoded twice** — diff-owned. The renderer builds the line at
   `skillsindex.go:63`; `byName` re-parses rendered lines with `indexLineRe` at `:262` to
   recover names already held in `Entry`. The module's own SI6 guard bans this literal in
   the conformance files. Disposition: `auto-fix`.
4. **Six comments narrate history or cite provenance** — diff-owned. "used to be derived
   twice", "today's posture", "now that the retired shell … is gone". `craft-comments`
   requires timeless present tense. `skillsindex.go:7,73,246`, `command.go:21`,
   `command_test.go:10`, `skillsindex_test.go:197`. Disposition: `auto-fix`.
5. **Three conformance comments restate forwarding code or argue their own correctness** —
   diff-owned. `checks_test.go:301`, `skills_index_checks_test.go:183,259`. Disposition:
   `auto-fix`.

## Spec

6 findings. Rows SI1, SI2, SI5, SI6, SI8, SI9 closed; SI3, SI4, SI7 partial; SI10 open.
Worst: the same `--check --write` mutation. All seven diagnostic strings verified
byte-exact against the spec, and the SI8 sweep residue matches its permitted list.

1. **`--check --write` writes instead of exiting 2** — diff-owned. Spec: "`bench
   skills-index [--check|--write]`", "usage exits 2" (spec.md:93). Disposition: `auto-fix`.
2. **Allowlist read bypasses the canonical row validator** — inherited. Spec: "reads
   … through the root package's `PayloadRow`/`PayloadKitOnlyPrefixes` … so the module adds
   no JSON parsing of its own" (spec.md:78). `PayloadRows()` validates but reads the
   *embedded* allowlist, so it cannot serve a temp root; the pre-collapse
   `kitOnlySkillSources` had the identical bypass. `skillsindex.go:235` vs
   `consumer_payload.go:42`. Disposition: `ask-user` — needs a root-scoped API.
3. **SI7 is not tested at the seam it names** — diff-owned. Spec: "CLI seam (`cmd/bench`
   dispatch test)" (spec.md:171); the test calls `skillsindex.Command` in-process and
   `cmd/bench/*_test.go` has no skills-index case. Disposition: `auto-fix`.
4. **SI3's exact three-line assertion is no longer pinned** — diff-owned. Spec: "check
   returns exactly the three attributed lines in that (alphabetical) order" (spec.md:167).
   The collision repair extended the mandated fixture instead of adding a case beside it,
   so it now expects four. `skillsindex_test.go:140,145`. Disposition: `auto-fix` — keep
   both fixtures. Cause: the coordinator's repair instruction, not the delegate.
5. **SI4 is falsely classified as already covered** — spec-authoring defect. The four named
   canaries' EXPECT files cover missing-entry, dangling-entry, missing-index and attributed
   drift; none covers the non-attributable `block drifted from generated form` case the row
   claims. Disposition: `ask-user`.
6. **SI10's permitted-edit list was exceeded** — diff-owned, non-behavioral.
   `kitOnlySkillSources` was deleted rather than collapsed-and-comment-repointed, and the
   unlisted `skillNameFromIndexLine` was deleted too. No pre-existing assertion text
   changed. Disposition: `ask-user` — reviewer veto on widening the list.

## Coverage

12 findings. Worst: no producer path is classified before use, so an untrusted filesystem
shape can hang the gate indefinitely or drive a destructive rewrite from a falsely empty
state.

Inherited (behavior identical pre-collapse; verified against `5b41322a`):

1. **Glob metacharacters in the repo root** — `filepath.Glob` treats `/tmp/repo[1]` as
   pattern syntax; the ignored error yields zero skills, so `--write` erases the index.
   `skillsindex.go:80`. Disposition: `auto-fix`.
2. **No path classification before read** — a FIFO at `SKILL.md`,
   `consumer-payload.json` or `BENCH-reference.md` blocks indefinitely; devices and
   sockets collapse into misleading empty states; symlinks import bytes from outside
   `.agents`; dangling links read as authoritative empty. The profile's hostile-input
   checklist names every one of these. `skillsindex.go:78,101,230,313`. Disposition:
   `auto-fix`.
3. **Orphan skill directory with no `SKILL.md` is silently omitted** while the same
   directory with an empty `SKILL.md` is diagnosed. `skillsindex.go:80`. Disposition:
   `auto-fix`.
4. **Present-but-empty reference file reports as missing** (read failure and zero bytes
   both become `""`), and a block holding one blank line compares equal to an empty block,
   so `--check` is clean while `--write` changes tracked bytes. `skillsindex.go:142,150`.
   Disposition: `auto-fix`.
5. **The "leading fence" rule is not enforced** — prose, then `---`, `index: injected`,
   `---` is read as frontmatter, as is an unclosed opening fence. `skillsindex.go:106`.
   Disposition: `auto-fix`.
6. **Control bytes pass into a line-structured sink** — `index: safe\r- forged` renders a
   second line for any consumer treating CR as a break. The profile names this class
   explicitly. `skillsindex.go:61,113`. Disposition: `ask-user` — permitted byte set undecided.
7. **Marker cardinality unchecked** — two complete blocks means `Check` reads the first and
   passes while a stale second survives; a duplicated start marker lets `Write` succeed and
   leave a tree that immediately fails `Check`. `skillsindex.go:274,290`. Disposition: `auto-fix`.
9. **Semantically invalid allowlist bypasses validation** — a misspelled `audience` parses,
   marks nothing, and lets `--write` strip `(kit-only)` markers. `skillsindex.go:235`.
   Disposition: `ask-user`. Same repair target as Spec 2.
10. **Temp-file residue** — a rename failure or SIGINT after `CreateTemp` leaves
    `.bench/.skills-index-*`. The spec's only won't-handle is a crash mid-rename.
    `skillsindex.go:204`. Disposition: `ask-user`.

Diff-owned:

8. **`--check --write` conflicting pair untested** — SI7 covers each flag and `--bogus`,
   not the pair. `command_test.go:47`. Disposition: `auto-fix`. Same repair target as
   Standards 1 / Spec 1.
11. **`git` absent from PATH misclassified** as "not in a git repository".
    `command.go:33`, `internal/git/git.go:89`. Disposition: `ask-user`.
12. **No process-boundary proof for the verb** — `skills-index` is absent from the
    dispatcher's retained-route tests and the AXI suite excludes mutation commands, so
    wrapper routing, deep cwd and fresh-process reload can regress silently.
    `cmd/bench/command_registry_test.go:428`. Disposition: `auto-fix`. Same repair target
    as Spec 3.
