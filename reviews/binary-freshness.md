# Review pickup — binary-freshness

Frozen base `ee339a804cdd50d2bc9e7fd4c282594235fa1700`.
Reviewed tip `9250818f3cbd17fdd02d4ad26c03605bd019bfeb`.

Raw findings: Standards 10, Spec 9, Coverage 4. That is 23.
De-duplicated repair targets: 5. The remaining findings need a reviewer
decision and carry no repair target yet.

## Standards

10 findings: 3 hard violations and 7 judgment calls. The worst is S1.

- **S1 — the shim adds a third derivation of the envelope read, and it is
  wrong.** `ask-user`, and it blocks. AGENTS.md, "Code standard — one source
  per fact": "Two derivations of the same fact must collapse into one
  source." `.bench/hooks/block-bench-follow-on.sh:32` defines
  `envelope_command` with a greedy `sed`.
  `.bench/hooks/block-dangerous-git.sh:93` already defines the same function
  over the same field, character by character, scoped under `"tool_input"`.
  The coordinator traced the new reader on a real envelope: for
  `{"tool_name":"Bash","tool_input":{"command":"ls"},"cwd":"/home/u/bench"}`
  it extracts `ls"},"cwd":"/home/u/bench` and classifies a plain `ls` as a
  Bench call. Under a stale binary that refuses every `ls` in a checkout
  whose path holds `bench`, which is the failure story 22 forbids.
- **S2 — the new canary fixture pastes the live shared library.**
  `auto-fix`. Same rule.
  `tests/canary/guard-classifier-table/word-test-drops-xargs/files/dot-bench/lib/resolve-bench.sh`
  is the live 167-line file plus a header and one changed arm. The sibling
  fixture `guard-resolver-order-drift` uses a 15-line stub, and both guidance
  fixtures in this same diff use `BASE` with `MUTATE.json`.
- **S3 — the system-suite route is authored twice.** `ask-user`. `AGENTS.md`
  and `projects/benchkit.md` carry the same fact, and the second sentence is
  verbatim identical. The ticket asked for both files.
- **S4 — Middle Man.** `auto-fix`. `internal/adopt/doctor.go:302` keeps a
  one-line delegator to `KitSourceCheckout`, and its comment narrates the
  change, which `craft-comments` forbids.
- **S5 — a red record retained in two comments.** `auto-fix`.
  `internal/adopt/broker_test.go` and `cmd/bench/build_subject_mode_test.go`
  each keep a build-time red and cite "refusal three", an ordinal no code
  names.
- **S6 — a stale count in past-tense narration.** `auto-fix`.
  `internal/adopt/broker_test.go` says "Nine rows once reported healthy";
  `doctorRows` now holds eleven.
- **S7 — a comment describing the code that was.** `auto-fix`.
  `internal/preflight/source_tip_test.go:25` still says "exactly the five
  rows".
- **S8 — Mysterious Name.** `auto-fix`. `staleySealedCopy` in
  `internal/systemtest/owner_stale_seal_test.go`.
- **S9 — exported mutable global with no external reader.** `auto-fix`.
  `brokermanifest.Fields`; `rg` over the tree finds no consumer outside its
  own file.
- **S10 — three copies of the build-inputs fixture tree.** `no-op`. Seven
  packages already do this, and an unexported Go test helper cannot cross a
  package. Raised so the shared testing package can be ruled on once.

## Spec

9 findings. Every one of the 27 coverage rows resolves to a real test. Four
rows are partial. The worst is F1.

- **F1 — the build publishes the manifest where no route reads it.**
  `ask-user`, and it blocks. Spec line: "The manifest lands beside the
  resolved wrapper, as today." `internal/freshness/freshness_publish.go:29`
  derives the directory from the published executable, so an ordinary build
  writes `dist/bench-broker.manifest`. The complete reader enumeration is
  `bin/bench.sh:449` and `internal/adopt/broker.go:20`, and both name the
  wrapper's `bin/`, as `.gitignore:40` and `.bench/build-outputs.json` do.
  So stories 10 and 13 deliver nothing on the ordinary path, and
  `./dist/bench doctor` now reports a broker the landing will refuse.

  BF10 and BF14 miss it, because both grade the directory the writer chose.
- **F2 — the spec contradicts itself on the publish arity.** `auto-fix`. The
  Implementation decisions line says two arguments; the recorded amendment
  and the tree say one.
- **F3 — story 17 is conditional.** `ask-user`. `bin/bench.sh:466` reaches
  the version refusal only when `package_version` answers, so an install root
  with no readable `package.json` continues past a still-`dev` manifest. The
  recovery still cannot loop.
- **F4 — BF26 and BF27 grade `Decide` only.** `auto-fix`.
  `internal/preflight/gather.go:201` `binarySealFacts` spells the path and
  calls `freshness.Verify`, and no test reaches it.
- **F5 — a third copy of the rebuild sentence.** `auto-fix`.
  `.bench/hooks/block-bench-follow-on.sh:14` repeats `shell_quote` and
  `rebuild_action` from `.bench/hooks/session-start.sh:38`. A system test
  pins it against the Go source, so it cannot drift silently.
- **F6 — `requireSourceRoot` is behavior nobody asked for.** `no-op`. It is
  fail-closed and every caller names a root.
- **F7 — `commands --brief` is not the decided signature.** `no-op`.
  Behavior-equivalent; the decision line describes a shape the tree lacks.
- **F8 — the seal row's applicability predicate differs by consumer.**
  `no-op`. No producer of the divergent state exists today.
- **F9 — the residual `kitSourceCheckout` alias.** `no-op`. Standards owns
  it as S4.

## Coverage

4 findings. The worst is C1.

- **C1 — the shell word test disagrees with `benchguard.InvokesBench` on a
  quoted head, and the shared table pins none of that class.** `ask-user`.
  Eight of nine measured forms disagree: `"bench" gate`, `'bench' gate`,
  `\bench gate`, `( bench gate )`, `$(bench gate)`, `be"nch" gate`,
  `ls; "bench" gate`, `env "X=1" bench help`. The Go half runs a real lexer
  that strips quotes; the shell half runs `set -- $1` under `set -f`, which
  never unquotes. So BF25 passes vacuously over the class, and a Bench verb
  runs stale.

  The spec's Won't-handle cuts name symlinks and `bash -c` only, so quoting
  is undecided. The sample is nine forms of an open class.
- **C2 — the shim's envelope reader does not decode `\uXXXX`.** `auto-fix`,
  and it blocks. `{"tool_input":{"command":"ls && bench gate"}}`
  passes at exit 0 while the literal spelling refuses at exit 2. The profile's
  hostile-input checklist names this class directly. It shares one repair with
  S1.
- **C3 — preflight reports not applicable for a dangling `dist/bench`
  symlink, where doctor reds.** `auto-fix`, and it blocks.
  `internal/preflight/gather.go:202` uses `os.Stat`, which follows the link;
  `internal/adopt/doctor_rows.go:238` uses `os.Lstat`. The `--full` build
  preflight is the one that fails open. BF27 says "no `dist/bench`"; a
  dangling symlink is not that.
- **C4 — the land route's post-rebuild re-read dies bare when the manifest is
  gone.** `ask-user`. Under `set -euo pipefail`, `done < "$manifest"` on an
  absent file aborts with no `bench:` diagnostic and no 127. The first read is
  guarded; the second is not. Reachability is unproven, because
  `freshness.Publish` restores the manifest on rollback.

## The five repair targets

1. Collapse the envelope read to one source in `.bench/lib/resolve-bench.sh`
   and source it from both hooks. Closes S1 and C2.
2. Publish the build's manifest where the routes read it. Closes F1.
3. Grade `binarySealFacts` and use `os.Lstat`. Closes C3 and F4.
4. Replace the pasted fixture library with the `BASE` and `MUTATE.json`
   shape. Closes S2.
5. Fold the comment, naming, and stale-decision fixes. Closes S4 to S9, F2,
   and F7.
