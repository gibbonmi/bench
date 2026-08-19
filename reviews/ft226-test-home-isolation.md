# Review pickup — ft226-test-home-isolation

Base `b2d3d625` · reviewed tip `2275f7ae` · three read-only axis delegates, `opus` / high.

Raw findings: 12 (Standards 5, Spec 3, Coverage 4).
De-duplicated repair targets: 8. Coverage C1 and Spec 3 name one fix.

## Standards

Count 5. Worst: S1 — comment prose that argues the diff's own correctness, the
one register violation `craft-comments` names explicitly.

- **S1 · auto-fix.** `internal/worktree/main_test.go:12-13` and `:22-25` argue
  correctness to a diff reviewer ("what makes the record honest", "so neither can
  mask the others") rather than stating a constraint the next reader needs.
  Cites `.agents/skills/bench-craft-comments/SKILL.md`, "The register". The
  initialization-order sentence is a genuine why and stays.
- **S2 · no-op once S1 lands.** `main_test.go` carries 21 comment lines in 167
  against a package test aggregate near 4.6% (`reauthorize_test.go` 0/243,
  `ownership_test.go` 3/397, `resume_test.go` 16/467, `worktree_test.go` 33/488).
  Verified by the coordinator. The excess is mostly S1's prose.
- **S3 · ask-user.** `main_test.go:163-167` `withinDir` re-derives
  `internal/worktree/worktree.go:144-147` `insidePool`, differing only in
  treating `rel == "."` as inside; `internal/freshness/freshness.go:500` is a
  third copy. Cites AGENTS.md "one source per fact". Verified by the coordinator
  against the source. The Spec axis independently asserted no duplicate exists
  and was wrong. Judgment: independence from production code is arguably wanted
  for the `Pool` assertion, and the self-inclusive case genuinely differs.
- **S4 · auto-fix.** `main_test.go:60-63` and `:76-79` append every top-level
  entry to both `entries` and `lines`, so the count verdict and the printed
  report are two derivations of one list. One slice of `{path, origins}` with a
  render method single-sources it.
- **S5 · auto-fix.** `capture/session-handoff.md` restates the spec's Build
  verification log (sweep totals, the whole probe (b) correction). The file's own
  Shape section forbids the second copy. A fresh session needs the pointer and
  the veto flag, not the numbers.

Clean elsewhere: the `t.Setenv` at `reauthorize_test.go:230` matches the sibling
fixture convention; `requireTest`/`mustWrite`/`mustMkdirAll` are reused;
`TestMain` follows the `internal/systemtest/owner_test.go` prior art.

## Spec

Count 3. Worst: P1 — ticket 02 ships a checked acceptance box stating a command
behavior the build's own log records as false.

- **P1 · ask-user.** `specs/ft226-test-home-isolation/tickets/02-...md:35-37` is
  checked for `TMPDIR=/nonexistent go test ...`, while the spec's verification
  log records that this command never reaches `TestMain`. The coordinator ticked
  every box with a blanket replace and introduced this. DT4 itself is genuinely
  satisfied — `main_test.go:29-33` exits 1 on the `MkdirTemp` error before
  `m.Run`, evidenced two ways. The reviewer decides whether the spec's Testing
  decisions probe (b) and this box take the `GOTMPDIR=/tmp` form.
- **P2 · no-op.** EV1's "before" evidence is paraphrased rather than a
  transcribed `find`, and the sentinel path is a placeholder. Credible and
  consistent with the Problem section's measurement, but asserted.
- **P3 · auto-fix, same target as C1.** A walk error suppresses the per-entry
  residue listing DT1 requires.

Per-row verdict: IS1, IS2, IS3, DT1, DT2, DT4, EX1 delivered; DT3 delivered for
the `m.Run` and residue legs, its removal-failure leg untested by the row's own
design; EV1, SW1, SW2, SW3 asserted-only through the verification log, with SW1's
four hostile shapes proven in a scratch pool first. No production file is in the
diff, every path is inside the declared fence, and nothing from Out of scope was
built. Sweep arithmetic is self-consistent: 1,710 + 9 = 1,719; 9 survivors + 10
keys leaked by `main`'s concurrent gate runs = the 19 SW3 compares.

## Coverage

Count 4. Worst: C1 — the residue report goes silent precisely when a leak is
accompanied by an unreadable directory, defeating DT2's offender naming.

- **C1 · auto-fix.** `main_test.go:57-62` returns `residueReport{}` with the
  error, and `main_test.go:84` propagates the `WalkDir` callback error. A leaked
  worktree under a mode-`0o000` directory — a shape this package already creates
  at `orphan_sweep_test.go:179` and `classifier_shape_test.go:58` — makes
  `TestMain:41-44` print only `permission denied`, suppressing every residue path.
  Accumulate the origins error onto the report instead of replacing it.
- **C2 · auto-fix.** `TestHomeResidueListsLeakedWorktreesWithOrigin:129-160`
  plants only a well-formed pointer. A `.git` that is a directory and a `.git`
  with no `gitdir:` line are the spec's Malformed-input edge and are unexercised.
- **C3 · auto-fix.** `main_test.go:119-120` guards IS2's "not under the inherited
  value" half on `inheritedBenchHome != ""`, so it drops silently when a
  developer runs with `BENCH_HOME` unset; and `withinDir(p, p)` is never
  exercised.
- **C4 · no-op.** `main_test.go:35-38` exits without removing the home it just
  created. Reachable only on an invalid environment-variable name.

Refuted, not reported: the SIGINT helper child exits normally so its cleanup
runs; an `m.Run` panic is the spec's declared interrupted-state edge; `HOME`
isolation is a Won't handle; `worktree_test.go:93` restores correctly, so EX1
holds.
