# skills-index hostile-input hardening

Status: staged

Decision source: `ROADMAP.md` FT208 at `20d52ca7f997235514539dbeb001da558d44cb25` (reviewer-approved drain)

Verification log: spec 3 + tickets 1 completed advisory round(s) — initial Sol/xhigh and independent Terra/high spec reviews blocked; one Terra/high repaired-spec re-review blocked, and the next re-review was canceled before verdict; one Terra/high ticket review blocked. All findings were repaired. The reviewer explicitly signed off the corrected spec and breakdown without a terminal advisory ACCEPT verdict.

## Problem

The generated skills index has one semantic owner, but the gate currently opens the
same producer files again through conformance and package-surface checks. Eleven
inherited hostile or degenerate inputs can therefore hang a later reader even after
`internal/skillsindex` is fixed, erase the generated block, silently omit a skill,
accept forged rendered lines or malformed structure, leave scratch files behind, or
report the wrong recovery action. `bench skills-index --write` makes destructive cases
operator-visible because a false empty source can replace tracked reference bytes.

## Solution

Keep skills-index parsing, rendering, diagnostics, and replacement in
`internal/skillsindex`, while extending the existing `internal/bounds` control-record
classifier with an explicit no-follow form. Every applicable gate consumer of
`SKILL.md`, `.bench/BENCH-reference.md`, and `.bench/consumer-payload.json` uses that
one bounded byte-classification owner before interpreting bytes. `internal/skillsindex`
then enforces literal directory enumeration, complete frontmatter, safe one-line
values, exactly one marker span, and canonical allowlist-row validation. Check and
write refuse attributed producer defects; write preserves the reference and cleans up
after rename failure or catchable interruption. Command diagnostics distinguish a
missing Git executable without changing `internal/git`.

## User stories

1. **Hostile producer topology cannot become authoritative input.** As the gate or an
   operator, I get a bounded, path-attributed refusal when a producer is a symlink or
   non-regular file; a root containing spaces or glob syntax is enumerated literally;
   a real orphan skill directory is diagnosed before command-adapter suppression; and
   empty is distinct from absent. Write leaves the reference unchanged. Line:
   `gpt-5.6-terra` / medium. This is oracle behavior at a known shared classifier seam.
2. **Only well-formed producer bytes reach their consumers.** As the gate or an
   operator, I get a fail-closed diagnostic for invalid UTF-8 or oversized producer
   bytes, incomplete leading frontmatter, a control rune in a rendered value, invalid
   marker cardinality, or a syntactically/semantically invalid allowlist. No later gate
   reader reopens the path unsafely. Line: `gpt-5.6-terra` / medium.
3. **Failed or interrupted replacement leaves honest recovery state.** As an operator,
   I find no `.bench/.skills-index-*` residue after rename failure or SIGINT, and the
   original reference remains authoritative unless replacement completes. Line:
   `gpt-5.6-terra` / high. Catchable interruption needs a deterministic fresh-process
   handshake through the production replacement path.
4. **Repository-discovery failures name the correct recovery action.** As an operator,
   I see required tool `git` named as missing or non-executable when it cannot launch,
   while executed Git outside a repository retains the existing not-in-repository
   diagnostic. Line: `gpt-5.6-terra` / medium. This is a known command seam and HI11.

Story partition: stories 1 and 2 share the bounded producer-read contract across
`internal/bounds`, `internal/skillsindex`, conformance, and package-surface consumers.
Story 3 is the higher process-lifecycle tracer. Story 4 is deliberately separate from
cleanup: it changes only command error classification and requires no `internal/git`
edit. Together they are one serial skills-index hardening build because every story
terminates at the same gate/operator command behavior.

Ticket partition: HI1-HI14 remain owned by behavior-bearing tickets. One final
user-facing completion ticket owns only the release-note entry and is blocked by every
behavior ticket, so it cannot advertise partial hardening and is not an implementation
dependency of any behavior.

## Implementation decisions

- `internal/bounds` remains the one owner of file-shape and bounded-byte
  classification. Add the smallest explicit no-follow composition beside `Classify`;
  do not change the existing follow-live-symlink contract used elsewhere and do not
  create a skills-index-local `os.FileMode` classifier. The new form performs `Lstat`,
  refuses every symlink before `Open`, accepts only regular files, reads at most
  `bounds.ControlRecordLimit` (2 MiB), and retains the existing UTF-8 and state
  vocabulary.
- The no-follow disposition is exact: an absent path is `absent`; a present zero-byte
  regular file is `empty`; bytes through exactly 2 MiB are classified normally; 2 MiB
  plus one is `unreadable` with `ReadOversized`; invalid UTF-8 is `malformed`; a live
  symlink and a dangling symlink are both `wrong-type` without following; directories,
  FIFOs, devices, sockets, and other non-regular objects are `wrong-type` without
  opening. These are the only new `bounds` promises.
- Inventory of applicable gate readers to migrate, preserving each consumer's own
  parser and diagnostic policy:
  - `SKILL.md`: `internal/skillsindex` plus conformance's skill-frontmatter,
    craft-name, Codex-adapter, AXI-query-registry, workflow-anchor, guidance-token,
    and prose-budget readers in `skills_index_checks_test.go`, `checks_test.go`,
    `axi_query_registry_test.go`, `docs_workflow_helpers_test.go`,
    `guidance_token_sweep_test.go`, and `prose_budget_test.go`; workflow anchors also
    cross `internal/anchors.EvaluateGroup`, whose registry targets include skill files.
  - `BENCH-reference.md`: `internal/skillsindex` plus the command-guide/adapter and
    docs-workflow readers in `skills_index_checks_test.go`,
    `docs_workflow_checks_test.go`, and `docs_workflow_helpers_test.go`, including
    `internal/anchors.EvaluateGroup` when it evaluates registered reference targets.
  - `consumer-payload.json`: `internal/skillsindex`, package-core's packed-asset
    guard, package-shipped-surface, and `internal/packagesurface` contract-document
    inventory.
  A generic helper may delegate to the new classifier only for exactly these three
  producer classes: `.agents/skills/*/SKILL.md`, `.bench/BENCH-reference.md`, and
  `.bench/consumer-payload.json`. Its behavior for every other path remains unchanged;
  producer-specific absence/empty diagnostics stay at their caller. No parser
  predicates move into `bounds`, and FT208 does not silently widen to other files read
  by the same generic helper.
- Preserve `internal/skillsindex` as the one owner of skills-index discovery, entry
  rendering, marker parsing, diagnostics, and replacement. Conformance invokes
  `skillsindex.Check`; it must not copy marker text, line shape, or allowlist policy.
  Composition tests through the actual conformance/package checks are required because
  package-local tests alone cannot detect a later unsafe reopen.
- Enumerate `<root>/.agents/skills` literally. Classify every real child directory and
  its expected `SKILL.md` before asking whether a same-named
  `.agents/commands/<name>.md` suppresses the adapter from the rendered index. Thus
  `.agents/skills/bench-example/` without `SKILL.md` is diagnosed even when
  `.agents/commands/bench-example.md` exists; the paired valid command-named skill is
  classified successfully and then suppressed. Roots containing spaces, `[`, `]`,
  `*`, or `?` remain literal.
- Producer absence semantics remain caller-owned. Missing `SKILL.md` in a real skill
  directory is an orphan defect. Missing `BENCH-reference.md` is unverifiable. Missing
  optional `consumer-payload.json` continues to mean no kit-only prefixes. Empty skill
  bytes produce the distinct missing-`index:` diagnostic; an empty reference is
  `empty`, not `missing`, and blocks write; a present empty allowlist is invalid JSON
  and blocks check and write.
- Frontmatter starts at byte zero with a `---` line and requires a closing `---` line.
  Text before the opener and an unclosed opener cannot authorize `index:` or
  `index-note:`. Preserve first exact-key value and no-trailing-newline behavior inside
  a valid fence. Returned values refuse every Unicode control rune, including tab, CR,
  ESC, BEL, NUL, and DEL; graphic Unicode remains valid.
- Marker parsing succeeds only for exactly one start marker followed by exactly one
  end marker in the whole reference. Missing, reversed, duplicated, or unclosed spans
  block check and write without changing bytes.
- Export one root-package parser for payload bytes: JSON decoding followed by the
  existing canonical row validator. Filesystem-backed skills-index, package-core,
  package-shipped-surface, and packagesurface consumers classify before passing bytes
  to it; embedded `PayloadRows` passes its `go:embed` bytes directly to the same parser.
  Unknown audience, empty/unsafe source, and duplicate source are one
  enumerated semantic-invalid partition. Check refusing `Entries`' allowlist error is
  new behavior: current `Check` discards that error; write already refuses it.
- Replacement owns cleanup from successful `CreateTemp` until successful rename. A
  deferred removal covers write, close, chmod, cancellation, and rename errors;
  successful rename disarms it. The write/command path accepts a context and checks it
  immediately before replacement.
- The SIGINT test uses a production-reached pre-replacement operation seam. After the
  real production path creates and writes the sibling temp, the child-provided barrier
  publishes one exact marker on an inherited pipe and blocks on the same cancellation
  context. The parent waits for that marker, sends SIGINT, and asserts nonzero exit,
  original bytes, and no residue. There is no directory polling, sleep, or temp-name
  appearance oracle.
- `git.Root` already returns the underlying `exec.Cmd.Run` error through `Output`, so
  no `internal/git` edit is required. `skillsindex.Command` classifies an
  executable-not-found cause as missing/non-executable Git; only an executed Git probe
  failure retains `toon.NotInRepo`.
- The operator-visible refusal/cleanup change receives one concise typed
  `CHANGELOG.md` entry under the standing `craft-synthesis` policy. A final
  user-facing completion ticket owns only that entry and names every HI1-HI14 behavior
  ticket as a blocker; it records the complete shipped outcome after behavior lands
  rather than participating in implementation.
- Bootstrap authority is not applicable: FT208 authenticates no executable and makes
  no trusted-execution claim; its process test exercises the already-built test binary.

## Testing decisions

- `internal/bounds` table tests the complete no-follow partition: absent, empty,
  exact-limit, oversized, invalid UTF-8, live symlink, dangling symlink, directory,
  FIFO, device, socket, and ordinary regular file. Real filesystem fixtures own all
  constructible rows; existing capability helpers cover device/socket availability.
- Package tests drive `Entries`, `Check`, and `Write` against real producer objects and
  assert diagnostics, bounded completion, unchanged reference bytes, and no residue.
  The classifier table replaces the proposed local synthetic-mode seam; only injected
  rename failure and pre-replacement context barrier remain justified lower seams.
- Composition rows invoke the registered conformance functions, not replicas:
  `load-validity-metadata`, `skills-index-command-adapters`, `docs-currency-workflow`,
  `line-routing`, `guidance-prose-budgets`, and `axi-query-registry` for their
  inventoried `SKILL.md` reads; the skills-index and docs checks, including the
  `anchors.EvaluateGroup` route, for `BENCH-reference.md`; and both package-core and
  package-shipped-surface paths for `consumer-payload.json`. Each uses
  hostile FIFO/symlink/oversized/invalid-UTF-8 fixtures as applicable and must complete
  with an attributed refusal. This is what prevents a later `os.ReadFile` from hanging
  after `internal/skillsindex` itself has gone green.
- Payload tests drive the root parser over invalid JSON and every semantic-invalid row,
  then drive `Check`, `Write`, package-core, package-shipped-surface, and packagesurface
  with the same bytes. This proves canonical row predicates remain single-sourced.
- The fresh-process SIGINT handshake is exactly the marker/barrier protocol above. The
  missing-Git partition also crosses the real `cmd/bench` registry dispatch so the
  operator surface, stdout, and exit code—not merely a helper return—are asserted.
- `bench skills-index` is AXI-exempt because it mutates; generic AXI envelope tests do
  not observe this behavior. The exact command-registry and conformance rows below are
  its surface/gate coverage.

### Seam diagram

    gate checks · cmd/bench registry
                 │
                 ▼
    producer path ─▶ bounds no-follow classification (2 MiB + UTF-8)
                 │                  │
                 ├──▶ skillsindex parser/check/write ─▶ diagnostics/replacement
                 ├──▶ conformance skill/docs consumers
                 └──▶ payload parser ─▶ package-core · shipped-surface · packagesurface
                                      │
                         temp created/written → barrier marker → SIGINT/cancel

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| HI1 | 1 | The no-follow classifier enumerates absent, empty, exact 2 MiB, 2 MiB+1 oversized, invalid UTF-8, live/dangling symlink, FIFO, device, socket, directory, and regular-file dispositions exactly as specified. | `internal/bounds` public classifier | Not yet observed: the new no-follow API/table do not exist; the current `Classify` follows a live symlink. | It prevents a second mode classifier and catches follow-before-open, unbounded-read, and byte-validation omissions at their shared owner. |
| HI2 | 1 | A root containing spaces and `[ ] * ?` with one indexed skill is enumerated literally; `Check` sees it and `Write` does not erase it. | skillsindex module | Not yet observed: current `filepath.Glob` returns zero for a hostile literal root. | It fails the cheapest surviving glob implementation. |
| HI3 | 1 | A real orphan directory is diagnosed before adapter suppression; the paired command-named directory with a valid `SKILL.md` is classified and then suppressed; an empty regular skill stays distinct. | skillsindex module | Not yet observed: current glob omits both missing-file directories before suppression. | The pair fixes ordering and prevents suppression from hiding a missing producer. |
| HI4 | 1 | Missing, empty, oversized, invalid-UTF-8, and wrong-type reference files retain distinct attributed refusals and all block `Write` without changing bytes. | skillsindex module | Not yet observed: current read helper collapses failures and empty bytes to `missing`. | It proves the shared states survive caller policy instead of collapsing back to false absence. |
| HI5 | 2 | Text before `---` or an unclosed opener cannot supply a field; a complete leading fence returns the first exact-key value with or without trailing newline. | skillsindex module | Not yet observed: current parser accepts late and unclosed fences. | It pins the complete leading-fence predicate without full YAML scope. |
| HI6 | 2 | `index:` and `index-note:` reject tab, CR, ESC, BEL, NUL, DEL, and every other Unicode control rune; graphic Unicode renders on one line. | skillsindex module | Not yet observed: current renderer accepts the control partition. | The refused/permitted pair catches sink-line injection and overbroad ASCII-only fixes. |
| HI7 | 2 | Zero, reversed, duplicated, or unclosed marker pairs are malformed; exactly one ordered pair succeeds; malformed writes preserve bytes. | skillsindex module | Not yet observed: current parser accepts the first plausible pair. | It catches stop-at-first-pair implementations. |
| HI8 | 2 | Payload bytes reject invalid JSON plus unknown audience, empty/unsafe source, and duplicate source through the canonical parser; absent remains optional; empty is invalid; all present invalid states make `Check` and `Write` refuse. | root payload parser → skillsindex | Not yet observed: current `Check` discards `Entries`' error and direct readers bypass canonical validation. | It catches JSON-only, copied-partial-validator, and check-error-discard degenerates. |
| HI9 | 3 | Rename failure after temp creation returns the error, preserves original bytes, and leaves no sibling temp. | replacement operation seam | Not yet observed: current rename return leaves the temp. | It proves cleanup authority spans the last failure path. |
| HI10 | 3 | Child publishes the exact pre-replacement marker after real temp creation/write and blocks on context; parent waits, SIGINTs, then observes nonzero exit, original bytes, and no residue. | fresh process + production-reached barrier | Not yet observed: current command has no cancellation/barrier path and leaves residue when interrupted in the vulnerable interval. | The handshake makes signal timing deterministic and cannot pass via polling or sleep. |
| HI11 | 4 | Empty `PATH` yields exit 1 naming required tool `git`, never not-in-repository; Git present outside a repo retains the established diagnostic. | `cmd/bench` registry dispatch | Not yet observed: current `Command` collapses both errors to `toon.NotInRepo`. | The paired command-surface partition proves underlying error classification without changing Git discovery. |
| HI12 | 1 | The registered `load-validity-metadata`, `skills-index-command-adapters`, `docs-currency-workflow` (including `anchors.EvaluateGroup`), `line-routing`, `guidance-prose-budgets`, and `axi-query-registry` paths complete over hostile `SKILL.md` with the shared attributed refusal. | actual conformance checks | Not yet observed: direct conformance, AXI guidance, and anchor-registry readers bypass a hardened skillsindex path. | Composition degenerate—skillsindex goes green while any later SKILL reader hangs—goes red here. |
| HI13 | 2 | Registered `skills-index-command-adapters` and `docs-currency-workflow`, including its `anchors.EvaluateGroup` route, complete over hostile `BENCH-reference.md` with the shared attributed refusal. | actual conformance checks | Not yet observed: docs, adapter, and anchor-registry readers reopen the reference directly. | Composition degenerate—module check refuses but a later guide or anchor reader hangs or erases state—goes red here. |
| HI14 | 2 | Registered `package-core-guard`, `package-shipped-surface`, and their `packagesurface.ContractDocumentInputs` consumer all refuse the same hostile or invalid `consumer-payload.json` through shared classification and canonical parsing. | actual package gate consumers | Not yet observed: all three currently open/parse the root file independently. | Composition degenerate—skillsindex refuses while a package reader hangs or accepts invalid rows—goes red here. |

The red signals are honest for this documentation-only repair: no implementation test
was authorized. Each behavior ticket must make its mapped test the first edit,
record the current red, and stop if the red cannot be reproduced.

### Edge inventory

- Error paths: HI1, HI4, HI7-HI14.
- Empty/absent: HI1, HI3, HI4, and HI8 enumerate producer-specific dispositions; the
  existing zero-skill-root case remains the zero-cardinality control.
- Boundary values: HI1 fixes exact limit/limit+1; HI2 uses one skill; HI7 fixes marker
  cardinality zero/one/many.
- Malformed: invalid UTF-8 in HI1/HI4 and syntax/semantic partitions HI5-HI8.
- Interrupted/partial state: HI9 and deterministic HI10.
- Re-run idempotency: existing second-write byte equality plus HI9/HI10 no-residue.
- Process lifecycle: HI10 fresh process; HI11 actual command dispatch environments.
- Hostile environment: HI1, HI2, HI6, HI10, HI11, and composition HI12-HI14.
- Project checklist: spaces/globs HI2; sink controls HI6; missing trailing newline HI5;
  absent/empty HI1/HI3/HI4/HI8; special files and links HI1/HI12-HI14; missing tool
  HI11; interrupt HI10; rerun HI9/HI10; fresh process HI10/HI11.

Won't handle: argument splitting, positional grammar, prompting stdin, worktree
destruction, and host-backed fsync pressure do not reach the changed reader contract.
Wrapper/symlinked executable invocation and deep-cwd discovery are unchanged
dispatcher/repository concerns; the surviving in-scope caller is the real `cmd/bench`
registry dispatch in HI11. Generic AXI coverage is explicitly not credited.

## Ownership fences

- `internal/bounds/classify.go`
- `internal/bounds/classify_test.go`
- `internal/skillsindex/`
- `consumer_payload.go`
- `consumer_payload_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/skills_index_checks_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/guidance_token_sweep_test.go`
- `internal/conformance/prose_budget_test.go`
- `internal/conformance/package_core_checks_test.go`
- `internal/conformance/package_shipped_surface_test.go`
- `internal/anchors/registry.go`
- `internal/anchors/registry_test.go`
- `internal/packagesurface/contract_documents.go`
- `internal/packagesurface/contract_documents_test.go`
- `cmd/bench/command_registry_test.go`
- `CHANGELOG.md`

Reviewer disposition requested: approve only these exact files/prefix. No gate registry,
fixture family, operating-guide document, `internal/git`, wrapper, or release code is
required. Cross-fence contract: `bounds` owns shape/bytes; root package owns payload-row
validity; skillsindex owns index semantics; each gate consumer owns only its diagnostic
and interpretation.

## Out of scope

- Full YAML parsing/schema validation for all skill frontmatter (FT89): separate
  compatibility/dependency capability, about 6 edits and 4 gate runs.
- Changing index line format, skill-name grammar, or indexed population: separate
  product compatibility, about 5 edits and 3 gate runs after reviewer decision.
- Hostile hardening of files other than the three inventoried producer classes: a
  repository-wide capability, at least 12 edits and 6 gate runs. Every current gate
  reader of these three classes is in scope; this cut cannot hide a later reopen.
- Recovery after SIGKILL, power loss, or host-filesystem durability failure: separate
  durable-transaction/fsync contract, about 7 edits and 5 gate runs. Catchable SIGINT
  cleanup remains in scope as HI10.
