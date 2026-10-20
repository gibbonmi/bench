# Candidate binary install for a persisted-field build

Status: staged

Roadmap: FT419

Decision source: named reviewed artifact `roadmap/FT419.md` (drain d-f9b840c3b468, reviewer approval 2026-10-13), with the reviewer-confirmed conversation answers of 2026-10-20 that select the install route and the held refresh

Verification log: pending

## Problem

A kit build can add a field to a persisted record. The two records in scope are the intent ledger and the review record. Each Bench executable decodes both records strictly, so an older executable refuses a record that holds the new field.

The intent ledger is one file under the Git common directory, and every checkout shares it (`internal/intent/intent.go`, `Address`). On 2026-10-12, a chunk checkpoint ran the build worktree's own `dist/bench` and wrote `commitment.gaps` into that ledger. The `main` executable then refused the whole ledger until the reviewer approved a hand removal of the key.

The review record is a tracked file in each tree. The landing runs only the installed promotion broker (`bin/bench.sh`, `land_route`), and the broker decodes the review record in its own process (`internal/gate/checkpoint.go`, `applyCheckpoint`). On 2026-10-12, the installed broker refused the new `folds` field. The build landed on its candidate binary only by a reviewer decision.

No guidance tells an orchestrator which executable runs a checkpoint or a landing of such a build. No check refuses the mismatch before the first bad write.

## Solution

A spec declares each persisted field that its build adds, on one `Persisted fields:` header line. The reviewer signs off that line with the spec. The strict decoders stay unchanged, and no record gains a version marker.

Before the first ticket dispatch, the orchestrator runs `bench worktree build --install <target>` from the primary checkout. The verb builds the candidate and publishes it as the primary checkout's executable. Then the verb runs the existing `doctor --fix` route, so the broker manifest binds the candidate. After the install, the hooks, `bench status`, the checkpoints, the record verbs, and the landing all decode with the candidate's code.

Three guarded entries of a declared build refuse before any ledger read or write when the running executable is not the installed candidate. The chunk checkpoint and the gap close of `bench record assignment` write the shared ledger. The complete checkpoint decodes the review record. Each refusal names the install command.

While the declared build is open, the post-landing refresh of any other landing holds the installed candidate and does not rebuild. After the build's own landing, the refresh returns the primary checkout to a normal `main` build. When a build stops before its landing, the operator runs `bench worktree build --restore <slug>`. That verb clears the declared ledger fields under the installed candidate and then reinstalls the `main` executable.

## User stories

Line: opus / high.
Implementation-line reason: CB-C2 is the hardest chunk, because it moves landing authority to candidate code through a trust boundary that ADR 0016 guards. The seams are known, and system tests can drive two real executables, but a wrong executable selection passes every unit test.
Harder chunks: CB-C2.

### Declare the persisted fields

1. As a spec author, I want to declare each persisted field that my build adds on one header line, so that the install obligation is machine-visible.
2. As a reviewer, I want the declaration to name the record and the field path, so that my sign-off covers an exact field.
3. As a spec author, I want a spec with no declaration to keep today's behavior, so that most builds pay nothing.
4. As an orchestrator, I want a malformed declaration to refuse with a named face, so that a typing error does not disable the guard.
5. As an orchestrator, I want an unknown record name to refuse, so that only the two strictly decoded shared records take the obligation.
6. As a maintainer, I want each record owner to answer whether its own decoder knows a field path, so that the answer has one source.
7. As a maintainer, I want that answer to come from the record's own struct tags, so that a new field needs no second registration.
8. As a maintainer, I want a path through a list field to resolve through the list element, so that a nested field such as `commitment.gaps.chunk` resolves.

### Install the candidate

9. As an orchestrator, I want one verb to install my build's candidate as the repository executable and broker, so that every reader decodes the new field.
10. As a reviewer, I want the install to read its authority from the default-branch copy of the bound spec, so that candidate code cannot authorize its own install.
11. As a reviewer, I want the install to refuse a target whose bound spec on the default branch declares no persisted field, so that the exception stays narrow.
12. As a reviewer, I want the install to refuse an unbound target, so that only a committed outcome can install a candidate.
13. As a maintainer, I want the install to refuse in a repository that declares no Bench build inputs, so that a linked repository is never changed.
14. As an orchestrator, I want the worktree copy and the primary copy of the candidate to be the same bytes, so that a writer in either location matches the broker manifest.
15. As a reviewer, I want the broker manifest to bind the primary `dist/bench` path through the existing `doctor --fix` route, so that the wrapper's land route stays unchanged.
16. As an orchestrator, I want to run the install again after a later Go commit, so that the installed candidate holds every field that the build adds.
17. As a reviewer, I want the install to refuse unless the running executable is the primary checkout's `dist/bench`, so that an uninstalled candidate never runs the install.
18. As an orchestrator, I want the install to refuse when a package or version-cache executable comes before the primary `dist/bench` in the wrapper's order, so that an install never leaves a second executable in charge.

### Refuse a ledger write on the wrong executable

19. As an orchestrator, I want a declared build's writers to pass on the installed candidate before any ticket adds the declared field, so that the guard never blocks its own build.
20. As an orchestrator, I want a declared build's chunk checkpoint to refuse when the running executable is not the manifest-bound broker, so that no uninstalled or rebuilt candidate writes the shared ledger.
21. As an orchestrator, I want a declared build's complete checkpoint to refuse on the same digest mismatch, so that the completion path cannot bypass the guard.
22. As an orchestrator, I want the gap close of `bench record assignment` for a declared build to refuse on the same digest mismatch, so that the first dispatch cannot write the new key.
23. As an orchestrator, I want each refusal to come before any ledger read or write, so that the named refusal wins over a decode error.
24. As an orchestrator, I want each refusal to name `bench worktree build --install` with my target label, so that I have one exact next command.
25. As an orchestrator, I want the refusal to print the slot form of the install command when the ledger is unreadable, so that the route stays exact without an owner lookup.
26. As an orchestrator, I want a version-cache executable to refuse the same way, so that a cache that appears after the install cannot bypass the guard.

### Land through the installed candidate

27. As a reviewer, I want the installed candidate to land its own build through the unchanged land route, so that the manifest digest stays the landing's authentication.
28. As a maintainer, I want the post-landing refresh to leave the broker manifest bound to the primary executable, so that the next landing authenticates.

### Restore after a stopped build

29. As an operator, I want the restore to preserve a copy of the ledger before it writes, so that a wrong clear can be undone.
30. As an operator, I want the restore to clear each ledger field that the stopped build's spec declares, so that the `main` executable can read the ledger again.
31. As an operator, I want the restore to rebuild the primary executable from `main` and rerun `doctor --fix` under it, so that `main` code authors the broker manifest, the shim, and the pre-push hook again.
32. As an operator, I want the restore to refuse when the running executable cannot read the ledger, so that the restore never writes a ledger it cannot read.
33. As an operator, I want the restore to refuse while the ledger holds an open build, so that it never clears the fields of a running build between chunks.

### Document the exception

34. As a teammate who just walked in, I want ADR 0016 to state the one install exception, so that the broker rule and the install agree.
35. As a spec author, I want `craft-spec` to require each declared spec to name its pre-landing writers and its landing binary, so that the edge inventory carries the decision.
36. As an orchestrator, I want `/bench-implement-spec` to name the install step before the first ticket dispatch, so that the build runs it without a reviewer prompt.
37. As an orchestrator, I want `.bench/BENCH-reference.md` to describe the install and the restore, so that the command notes match the verbs.

### Reviewed exclusions

38. As a maintainer, I want an older binary's tolerant purge read to stay with its existing owner, so that this build does not absorb FT140's residual.
39. As a maintainer, I want another worktree's own older build to stay out of scope, so that the existing one-version-at-a-time rule governs it.

### Keep the candidate installed

40. As an orchestrator, I want an unrelated landing during a chunk gap of my declared build to leave the installed candidate in place, so that the shared ledger stays readable.
41. As a maintainer, I want the refresh to report `held` only while a declared build is open and the manifest binds the primary executable, so that a normal stale primary still rebuilds.
42. As an operator, I want the install to find the bound assignment when the ledger holds a field that the running executable does not know, so that the install can recover a lost candidate.
43. As an orchestrator, I want the install to publish the executable, the seal, and the manifest as one transaction, so that an interrupt leaves the earlier set in place.
44. As an orchestrator, I want each refusal route to name the bound integration target of the open declared build, so that a refusal in an author worktree names the right target.
45. As an orchestrator, I want `/bench-implement-spec` to rerun the install after each Go commit of a declared build, so that a rebuilt copy never meets a stale manifest.

## Implementation decisions

### Terms

- **Persisted field:** a JSON field path inside the intent ledger or the review record that a build adds.
- **Declaration:** the spec header line `Persisted fields:`. Its value is a comma-separated list of `<record>:<path>` items.
- **Record names:** `intent-ledger` and `review-record`. No other name is valid.
- **Field path:** dot-separated JSON tag names from the record root, for example `intent-ledger:commitment.gaps` or `review-record:folds`.
- **Primary checkout:** the main working tree of the repository, which holds the shim target `bin/bench.sh` and the broker manifest.
- **Primary executable:** the canonical physical path of the primary checkout's `dist/bench`.
- **Installed candidate:** a candidate executable at the primary executable path that the broker manifest binds by path and digest.
- **Open build:** a claimed outcome with an active bound assignment whose deliverable is a staged spec, as the landing-window derivation computes it. Its chunk-gap state does not matter, and the restore refuses on an open build of any outcome.
- **Declared build:** an open build whose default-branch bound spec declares a persisted field.
- **Integration target:** the first bound assignment of the declared build in the open-build derivation.
- **Held install:** a declared build exists, and the primary broker manifest digest equals the primary executable's digest.

### Declaration and field knowledge

The `internal/spec` header scan owns the declaration, beside `Status:` and `Roadmap:`. It returns the parsed items or a parse error. A spec with no line returns an empty list.

The parser refuses an empty item, an unknown record name, an empty path segment, and a duplicate item. The refusal face is `persisted-fields-invalid`, and each consumer prints it.

Each record owner exports one predicate that reports whether a field path resolves through its root type. `internal/intent/ledger` answers for `intent-ledger`, and `internal/reviewrecord` answers for `review-record`. The predicate walks the `json` struct tags by reflection, descends through pointers and slice elements, and refuses a path through a non-struct. The struct tags are the one source.

### The candidate check

One check owner serves the three ledger writers of a declared build. It takes the declaration, the running executable's path, and the primary broker manifest as arguments. It applies only where the repository declares Bench build inputs. It refuses with the face `candidate-not-installed` when the running executable's SHA-256 digest differs from the digest in the primary broker manifest.

The check does not test whether the running executable knows a declared field. Before the ticket that adds a field commits, the installed candidate cannot know that field, and it also cannot write it. The digest test alone stops an uninstalled or rebuilt executable from writing the shared ledger. After the field's ticket commits, the orchestrator reinstalls, because a rebuilt worktree copy no longer matches the manifest digest.

The face's route is `bench worktree build --install <target>` from the primary checkout, then a rerun of the refused command. When the ledger reads, the target label is the integration target of the declared build, whatever worktree the caller runs in. When the ledger does not read, the route prints the literal slot `<target>`.

Before the install, the `main` executable is the manifest-bound broker, so the check passes on it. That pass is safe, because `main` code cannot write a field that it does not know. The guidance step owns the install before the first dispatch.

The three entries call the check before any ledger read or write. The chunk checkpoint and the gap close write the shared ledger, and the complete checkpoint decodes the review record:

- `bench gate --checkpoint <spec> --chunk <id>` reads the declaration from the spec in the checkpoint source tree. It calls the check before the owner lookup, the review-record read, and the chunk-gap write.
- `bench gate --checkpoint <spec> --complete` calls the check at the same point.
- `bench record <spec> assignment` reads the declaration from the spec at the calling worktree's HEAD commit. It calls the check before its record read and its gap-close transaction.

### The install verb

`bench worktree build --install <target>` is a new form of the existing build verb. It keeps the build verb's one positional argument.

The verb refuses before any write in this order:

1. The repository declares no Bench build inputs.
2. The running executable's canonical physical path is not the primary executable. The route names the primary checkout as the directory to run from.
3. The wrapper's order selects a package or version-cache executable before the primary `dist/bench`. The refusal names that executable's path.
4. The target has no bound assignment in the ledger. The verb reads the assignments and the commitment bindings through the existing tolerant read mode and writes no ledger byte. Thus an unknown field elsewhere in the ledger does not stop the lookup, and the install can recover a lost candidate.
5. The default-branch tip holds no copy of the bound spec, or that copy declares no persisted field.

Then the verb builds the target's tree into the target's `dist/bench` through the existing build owner. The worktree copy stays in place. The verb publishes a copy of that executable, its seal, and the broker manifest at the primary checkout through the existing freshness publication owner. That owner locks the directory, keeps backups, and restores the earlier set on a failure or a signal. Last, the verb runs the primary checkout's `bin/bench.sh doctor --fix` as a child process. The child runs with the primary checkout as its working directory. Its environment unsets `BENCH_KIT` and `BENCH_WRAPPER` and sets `BENCH_RUN_BINARY` to the primary executable. `doctor --fix` writes the PATH shim and the pre-push hook and binds that executable in the broker manifest again. A failed child leaves the published set in place, and the verb reports the child's refusal; a rerun of the install repeats both steps.

After each Go commit of a declared build, the orchestrator reruns the install before the next dispatch or checkpoint. A plain `bench worktree build` of the integration worktree changes its digest, so the next guarded entry refuses and names the install.

The verb prints one table with the worktree, the executable digest, and the manifest path. Its `next` line names the first ticket dispatch.

### Bootstrap authority

The trace has four hops:

1. The operator's shell runs the PATH shim from the primary checkout. The shim runs the primary checkout's tracked `bin/bench.sh`, which is `main` code. The wrapper resolves the primary `dist/bench`. When the wrapper resolves a worktree executable or a cache executable instead, the install's second and third refusals stop the install before any write.
2. That trusted executable reads the authority: the bound assignment in the ledger and the `Persisted fields:` line of the default-branch copy of the bound spec. A reviewer-approved spec-stage landing wrote that copy. The candidate tree's spec copy is never the authority. On a reinstall, the running executable is the earlier installed candidate, which the same declaration already authorized.
3. The trusted executable runs the candidate tree's build script and then runs the candidate executable as the `doctor --fix` child. That child writes the broker manifest, the PATH shim, and the pre-push hook. No validator authenticates the candidate bytes. This is the reviewer-visible trust assumption: the reviewer's sign-off of the declaration authorizes candidate code to become the broker.
4. At the landing, the wrapper's unchanged `land_route` authenticates the broker by the manifest path, digest, version, and platform before it runs it.

The seal and the manifest record the installed bytes. Candidate code writes both, so neither one authenticates the candidate's behavior. The trust root of the exception is the reviewer's sign-off of the declaration, and the manifest digest binds the exact bytes that the install published.

### The landing

The land route, the manifest authentication, and the landing owner stay unchanged. With the candidate installed, the broker is the candidate, and it decodes every declared field. A landing on an uninstalled broker keeps today's strict-decode refusal; this spec adds no landing face.

### The held refresh

The post-landing refresh verifies the primary seal against the primary sources. An installed candidate's seal records the candidate tree's sources, so it never verifies there. Without a hold, any landing during a chunk gap of a declared build would rebuild `main` over the candidate. The ledger would then hold a field that the new primary executable refuses.

The refresh therefore reads the held-install predicate first. When it holds, the refresh writes nothing, and its effect result is `held`. A `held` result exits 0. When the predicate does not hold, the refresh runs as today.

The build's own completion landing reconciles the delivered outcome before its effects run. At its refresh, the build is no longer open, so the refresh rebuilds `main` and rebinds the manifest.

### The restore

`bench worktree build --restore <slug>` runs from the primary checkout, and its positional argument is the stopped build's spec slug. The declared ledger fields are the `intent-ledger` items of the default-branch copy of `specs/<slug>/spec.md`.

The verb refuses before any write in this order:

1. The repository declares no Bench build inputs.
2. The default-branch copy of the spec declares no persisted field.
3. The running executable's strict ledger read fails. This refusal comes before the open-build test, because that test reads the ledger too.
4. The ledger holds an open build of any outcome, because the clear is ledger-wide. The refusal names `bench commitment block` and `bench worktree release` as the routes that end the build.

Then the verb copies the ledger to `$BENCH_HOME/restore/<repo-key>/` and refuses when the copy fails. One strict ledger transaction then removes each declared field that the ledger holds. A declared field that the ledger does not hold needs no change, because the build stopped before any writer wrote it. For a nested path such as `commitment.gaps.chunk`, the transaction removes the last segment from each element along the path. Then the verb rebuilds the primary `dist/bench` from the primary sources through the existing rebuild owner. Last, it runs `doctor --fix` under the rebuilt executable with the install's child environment, so `main` code authors the broker manifest, the shim, and the pre-push hook again.

### The ADR 0016 exception

ADR 0016 gains this decided text: "One exception applies. A kit build whose reviewer-approved spec declares a persisted-record field installs its own candidate as the repository executable and broker before its first ticket dispatch. The installed candidate then owns that build's landing, and the post-landing refresh holds it while that build is open. The reviewer's approval of the declaration is the trust root of this exception. The install is the only route by which candidate code gains landing authority."

### Reader sweep and stale claims

The reader sweep found these readers of the broker rule: ADR 0016, `bin/bench.sh` (`land_route` comment), `internal/worktree/land.go` (`LandCommand` doc), `internal/adopt/broker.go` (`WriteBrokerManifest` doc), and `projects/benchkit.md` (Working-branch broker sentence). The install changes none of these statements except ADR 0016, because the landing still runs only the manifest-bound broker and never consults a repository executable.

Three readers verify the primary seal against the primary sources: the `commands --brief` probe (`cmd/bench/main.go`), the doctor binary-seal row (`internal/adopt/doctor_rows.go`), and the preflight seal row (`internal/preflight/binary_seal.go`). After an install, each one reports the primary executable as stale. A **Won't handle** line below disposes of them. A fourth reader, the post-landing refresh (`internal/worktree/land_effects.go`, `refreshBroker`), acts on the same verdict, and the held refresh above governs it.

`projects/benchkit.md` claims that `bench worktree land` refuses a dev executable that was not built from the current sources. The current landing has no such refusal (`internal/worktree/land.go`, `LandCommand`; `internal/worktree/land_freshness_test.go`, `TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal`). This spec flags the claim for reviewer veto and does not change it.

## Implementation chunks

Every chunk follows its predecessor's green checkpoint and its independent semantic review. The order puts each producer before each consumer that names it.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CB-C1 / 1-declare-persisted-fields.md | A spec declaration resolves against the running executable's decoders | CB1–CB12 | TestPersistedFieldsHeader, TestLedgerKnownField, TestReviewRecordKnownField | no |
| CB-C2 / 2-install-the-candidate.md, 3-hold-the-refresh.md, 4-land-through-the-installed-candidate.md | One verb installs the candidate as executable and broker, the refresh holds it, and the candidate lands its build | CB13–CB25, CB50–CB56 | TestBuildInstall, TestRefreshHoldsAnInstalledCandidate, TestCandidateInstallLandsANewLedgerField, TestCandidateInstallSurvivesAGapLanding | yes |
| CB-C3 / 5-refuse-an-uninstalled-writer.md | Each guarded entry refuses on an executable that is not the installed candidate | CB26–CB36, CB57 | TestCandidateCheck, TestCheckpointPersistedFields, TestAssignmentPersistedFields, TestCandidateNotInstalledRoute | no |
| CB-C4 / 6-restore-after-a-stopped-build.md | One verb restores `main` after a stopped build | CB37–CB43, CB49 | TestBuildRestore | no |
| CB-C5 / 7-document-the-install-route.md | The ADR, the spec discipline, and the build phase name the install route | CB44–CB48, CB58, CB59 | TestCandidateInstallAnchors | no |

CB-C2 precedes CB-C3 and CB-C4, because the refusal route and the restore name the install verb and reuse its child environment.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"CB-C1","tickets":["1-declare-persisted-fields.md"],"verification":[{"id":"spec","command":"bench test --package ./internal/spec"},{"id":"ledger","command":"bench test --package ./internal/intent/ledger"},{"id":"review","command":"bench test --package ./internal/reviewrecord"}]},{"id":"CB-C2","tickets":["2-install-the-candidate.md","3-hold-the-refresh.md","4-land-through-the-installed-candidate.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"system","command":"bench test --check system --run TestCandidateInstall"}]},{"id":"CB-C3","tickets":["5-refuse-an-uninstalled-writer.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"record","command":"bench test --package ./internal/reviewrecord/..."},{"id":"routes","command":"bench test --package ./internal/refusalroute"}]},{"id":"CB-C4","tickets":["6-restore-after-a-stopped-build.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]},{"id":"CB-C5","tickets":["7-document-the-install-route.md"],"verification":[{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"prose","command":"bench test --check prose"}]}],"final_verification":[{"id":"system","command":"bench test --check system"},{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/candidate-binary-install/spec.md"}]}
```

## Testing decisions

A good test drives a real verb or a real exported predicate and observes its printed face, its exit code, and the bytes of the ledger, the manifest, or the primary executable. No test asserts on a log line.

The engineering seams are:

- the `internal/spec` header scan, with `Facts` and `metadata` as prior art;
- the two record-owner predicates, with the strict decoders in `internal/intent/intent.go` and `internal/reviewrecord/parse.go` as prior art;
- the one candidate check, which takes the executable path and the manifest as arguments;
- the build verb in `internal/worktree/build.go`, with `TestExecRunsTheWorktreeDistBench` as prior art for a real worktree executable;
- the checkpoint entry in `internal/gate/checkpoint.go` and the assignment entry in `internal/reviewrecord/recordcmd/assignment.go`, with the face registry in `internal/refusalroute`;
- one system test that builds two real executables from two trees, with `TestSystemOwnerRefusesAStaleSelectedSeal` and `TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory` as prior art.

The system test is the gate seam that observes both incidents end to end. Its fixture tree adds one new ledger field and one new review-record field to a copy of the kit source.

### Seam diagram

    trigger: orchestrator runs bench worktree build --install <target> from the primary checkout
        │
        ▼
    default-branch spec declaration ──▶ [ install verb in the primary executable ] ──▶ primary dist/bench = candidate, manifest path and digest = candidate
                                            ◀ tests attach here: system test reads the manifest and the primary executable digest

    trigger: bench gate --checkpoint (--chunk or --complete), or bench record <spec> assignment
        │
        ▼
    spec declaration + running executable + manifest ──▶ [ one candidate check ] ──▶ candidate-not-installed, or the existing write
                                            ◀ tests attach here: planted spec, planted manifest, and injected executable path drive each entry

    trigger: bench worktree build --restore <slug> from the primary checkout
        │
        ▼
    default-branch spec declaration + ledger ──▶ [ restore verb ] ──▶ ledger copy, cleared fields, primary dist/bench = main
                                            ◀ tests attach here: ledger bytes, copy path, and primary seal verification

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CB1 | 1 | The header scan returns `intent-ledger:commitment.gaps` and `review-record:folds` from one `Persisted fields:` line | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A scan that ignores the line returns an empty list |
| CB2 | 3 | The header scan returns an empty list for a spec with no `Persisted fields:` line | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A scan that defaults to a non-empty list fails |
| CB3 | 2 | The header scan ignores a `Persisted fields:` line inside a code fence | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A scan without the fence state reads a template example |
| CB4 | 4 | The parser refuses an empty item with `persisted-fields-invalid` | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A parser that skips empty items accepts `intent-ledger:a,,` |
| CB5 | 5 | The parser refuses the record name `gate-evidence` with `persisted-fields-invalid` | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A parser without the closed name set accepts any record |
| CB6 | 4 | The parser refuses `intent-ledger:commitment..gaps` with `persisted-fields-invalid` | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A parser that splits without a segment check accepts an empty segment |
| CB7 | 4 | The parser refuses a duplicate item with `persisted-fields-invalid` | planned TestPersistedFieldsHeader in internal/spec/persisted_test.go | A parser that keeps duplicates passes this row |
| CB8 | 6 | The ledger predicate resolves `commitment.gaps` | planned TestLedgerKnownField in internal/intent/ledger/fields_test.go | A predicate with a hand list omits the field |
| CB9 | 8 | The ledger predicate resolves `commitment.gaps.chunk` through the slice element | planned TestLedgerKnownField in internal/intent/ledger/fields_test.go | A walker that stops at a slice refuses the nested path |
| CB10 | 7 | The ledger predicate refuses `commitment.unknown` | planned TestLedgerKnownField in internal/intent/ledger/fields_test.go | A predicate that answers true for every path passes nothing here |
| CB11 | 7 | The ledger predicate refuses `schema.extra` because `schema` is not a struct | planned TestLedgerKnownField in internal/intent/ledger/fields_test.go | A walker that ignores the field kind accepts the path |
| CB12 | 6 | The review-record predicate resolves `folds` and refuses `foldz` | planned TestReviewRecordKnownField in internal/reviewrecord/fields_test.go | A predicate wired to the ledger type answers the wrong record |
| CB13 | 9 | The install publishes the target's sealed executable bytes at the primary executable path | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb that only builds the worktree leaves the primary bytes unchanged |
| CB14 | 14 | The install leaves the worktree `dist/bench` and the primary `dist/bench` with one SHA-256 digest | planned TestBuildInstall in internal/worktree/build_install_test.go | Two separate builds can differ and break the digest match |
| CB15 | 15 | After an install whose `doctor --fix` child starts beside a built worktree `dist/bench`, the manifest `path` is the primary executable and its digest is that executable's digest | planned TestCandidateInstallLandsANewLedgerField in internal/systemtest/candidate_install_test.go | A child that resolves the worktree copy binds a worktree path with the same digest |
| CB16 | 13 | The install refuses in a repository with no Bench build inputs and writes no file | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb without the build-input guard rewrites a linked repository |
| CB17 | 12 | The install refuses a target with no bound assignment and writes no file | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb that skips the binding lookup installs any worktree |
| CB18 | 11 | The install refuses when the default-branch copy of the bound spec has no `Persisted fields:` line | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb without the declaration guard installs every build |
| CB19 | 10 | The install refuses when only the target tree's copy of the spec declares a field | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb that reads the candidate's spec copy lets candidate code authorize itself |
| CB20 | 16 | A second install after a new Go commit replaces both copies with the new digest | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb that skips an existing install keeps the stale candidate |
| CB21 | 17 | An install whose running executable is a worktree `dist/bench` refuses before its build and writes no file | planned TestCandidateInstallLandsANewLedgerField in internal/systemtest/candidate_install_test.go | A verb without the path check lets an uninstalled candidate read its own authority |
| CB22 | 18 | An install refuses and names the cache path when `$BENCH_HOME/cache/bin/<version>` holds an executable for this platform | planned TestBuildInstall in internal/worktree/build_install_test.go | A verb without the order check installs a candidate that the cache then hides |
| CB23 | 9 | After the install, `bench status` from the primary checkout exits 0 with a ledger that holds a `gaps` key | planned TestCandidateInstallLandsANewLedgerField in internal/systemtest/candidate_install_test.go | Without the install, the older primary executable refuses the ledger as in incident 1 |
| CB24 | 27 | The installed candidate lands a build whose published tree holds a review record with a `folds` key | planned TestCandidateInstallLandsANewLedgerField in internal/systemtest/candidate_install_test.go | Without the install, the installed broker refuses the review record as in incident 2 |
| CB25 | 28 | After that landing, the primary broker manifest digest equals the primary `dist/bench` digest | planned TestCandidateInstallLandsANewLedgerField in internal/systemtest/candidate_install_test.go | A refresh that rebuilds without the manifest leaves the next landing unauthenticated |
| CB26 | 19 | The candidate check passes when the running digest equals the manifest digest and the running predicate does not resolve a declared path | planned TestCandidateCheck in internal/gate/candidate_check_test.go | A check that also tests field knowledge refuses every first dispatch of a declared build |
| CB27 | 20 | The candidate check refuses with `candidate-not-installed` when the running digest differs from the manifest digest | planned TestCandidateCheck in internal/gate/candidate_check_test.go | A check without the digest test lets an uninstalled candidate write the shared ledger |
| CB28 | 3 | The candidate check passes a spec with no declaration and an unrelated manifest digest | planned TestCandidateCheck in internal/gate/candidate_check_test.go | A check that always compares digests breaks every undeclared build |
| CB29 | 20 | A declared chunk checkpoint on a mismatched digest refuses with `candidate-not-installed` and leaves the ledger bytes unchanged | planned TestCheckpointPersistedFields in internal/gate/checkpoint_persisted_test.go | A check placed after the gap write still writes the new key |
| CB30 | 21 | A declared `--complete` checkpoint on a mismatched digest refuses with `candidate-not-installed` and leaves the ledger bytes unchanged | planned TestCheckpointPersistedFields in internal/gate/checkpoint_persisted_test.go | A check wired only into the chunk path lets the complete path decode with the wrong executable |
| CB31 | 22 | A declared `bench record <spec> assignment` on a mismatched digest refuses with `candidate-not-installed` and leaves the ledger bytes unchanged | planned TestAssignmentPersistedFields in internal/reviewrecord/recordcmd/assignment_persisted_test.go | A gap close without the check writes the new key at the first dispatch |
| CB32 | 23 | With an unreadable ledger, the declared chunk checkpoint prints `candidate-not-installed` and not the decode error | planned TestCheckpointPersistedFields in internal/gate/checkpoint_persisted_test.go | A check placed after the owner lookup prints the decode error first |
| CB33 | 24 | With a readable ledger, the `candidate-not-installed` route prints `bench worktree build --install` followed by the integration target label | planned TestCandidateNotInstalledRoute in internal/refusalroute/route_test.go | A route that names only a family gives no exact next command |
| CB34 | 25 | With an unreadable ledger, the `candidate-not-installed` route prints `bench worktree build --install <target>` with the literal slot | planned TestCandidateNotInstalledRoute in internal/refusalroute/route_test.go | A route that reads the label from a failed lookup prints an empty or wrong target |
| CB35 | 24 | The `candidate-not-installed` route prints the rerun of the refused command after the install command | planned TestCandidateNotInstalledRoute in internal/refusalroute/route_test.go | A route without the rerun leaves the refused write undone |
| CB36 | 26 | A declared checkpoint that runs a version-cache executable after an install refuses with `candidate-not-installed` | planned TestCheckpointPersistedFields in internal/gate/checkpoint_persisted_test.go | A check that trusts any executable with the field lets the cache bypass the install |
| CB37 | 29 | The restore writes the ledger copy under `$BENCH_HOME/restore/<repo-key>/` before it changes the ledger | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore without the copy loses the ledger on a wrong clear |
| CB38 | 30 | The restore leaves the ledger with no `gaps` key when the named spec declares `intent-ledger:commitment.gaps` | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore that only rebuilds leaves a ledger that `main` refuses |
| CB39 | 30 | The restore removes `chunk` from each `gaps` element when the named spec declares `intent-ledger:commitment.gaps.chunk` | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore that clears only top-level keys leaves the nested field |
| CB40 | 31 | After the restore, the primary `dist/bench` verifies against the primary sources and the manifest digest equals its digest | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore that only clears fields leaves the candidate installed |
| CB41 | 32 | The restore refuses before any write when the running executable's strict ledger read fails | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore that runs the open-build test first prints a raw decode error and no named refusal |
| CB42 | 33 | With one open build whose chunk gap is open, the restore refuses before any write | planned TestBuildRestore in internal/worktree/build_restore_test.go | A guard keyed to the landing-window refusal lets the restore clear a running build between chunks |
| CB43 | 30 | The restore refuses when the default-branch copy of the named spec declares no persisted field | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore without the declaration guard clears fields that no spec names |
| CB44 | 34 | ADR 0016 holds the anchored sentence "The install is the only route by which candidate code gains landing authority." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the exception sentence reds the anchor |
| CB45 | 35 | `craft-spec` holds the anchored sentence "A spec that declares a persisted field names each pre-landing writer and the landing binary in its edge inventory." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the edge-inventory rule reds the anchor |
| CB46 | 36 | `/bench-implement-spec` holds the anchored sentence "Before the first ticket dispatch of a declared build, run `bench worktree build --install <target>` from the primary checkout." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the step reds the anchor |
| CB47 | 37 | `.bench/BENCH-reference.md` holds the anchored sentence "`bench worktree build --restore <slug>` clears the declared ledger fields of a stopped build and reinstalls the `main` executable." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the note reds the anchor |
| CB48 | 34 | The edited ADR and guidance pass the prose check | `bench test --check prose` | A sentence over the STE bound reds the prose check |
| CB49 | 30 | The restore of a spec whose declared ledger field the ledger does not hold completes and leaves the primary `dist/bench` verified against the primary sources | planned TestBuildRestore in internal/worktree/build_restore_test.go | A restore that requires each declared field to exist cannot restore a build that stopped before its field ticket |
| CB50 | 40 | After an unrelated landing during an open chunk gap of a declared build, the primary `dist/bench` digest is unchanged and `bench status` from the primary checkout exits 0 | planned TestCandidateInstallSurvivesAGapLanding in internal/systemtest/candidate_install_test.go | A refresh without the hold rebuilds `main` over the candidate, and `main` refuses the ledger |
| CB51 | 41 | With a declared build open and the manifest digest equal to the primary digest, the refresh reports `held` and writes no file | planned TestRefreshHoldsAnInstalledCandidate in internal/worktree/land_effects_test.go | A refresh without the predicate rebuilds the candidate away |
| CB52 | 41 | With a declared build open and the manifest digest different from the primary digest, the refresh rebuilds and reports `complete` | planned TestRefreshHoldsAnInstalledCandidate in internal/worktree/land_effects_test.go | A hold keyed only to the open build keeps a stale primary forever |
| CB53 | 41 | With no declared build open and an unverified primary seal, the refresh rebuilds and reports `complete` | planned TestRefreshHoldsAnInstalledCandidate in internal/worktree/land_effects_test.go | A hold that ignores the open build stops every normal refresh |
| CB54 | 42 | With a ledger that holds an unknown commitment field, an install run by the `main` executable finds the bound assignment and publishes the candidate | planned TestBuildInstall in internal/worktree/build_install_test.go | A strict ledger read in the install leaves a lost candidate with no recovery |
| CB55 | 42 | The install leaves the ledger bytes unchanged | planned TestBuildInstall in internal/worktree/build_install_test.go | A tolerant read that also writes would drop the unknown field |
| CB56 | 43 | When the publication fails after the executable write, the primary executable, seal, and manifest keep their earlier bytes | planned TestBuildInstall in internal/worktree/build_install_test.go | A rename sequence without the publication owner leaves a candidate beside the old manifest |
| CB57 | 44 | A refusal from an author worktree prints `bench worktree build --install` followed by the integration target label and not the author label | planned TestCandidateNotInstalledRoute in internal/refusalroute/route_test.go | A route that takes the caller's own owner lookup names a target that the install refuses |
| CB58 | 45 | `/bench-implement-spec` holds the anchored sentence "After each Go commit of a declared build, rerun `bench worktree build --install <target>` before the next dispatch." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the reinstall step reds the anchor |
| CB59 | 41 | `.bench/BENCH-reference.md` holds the anchored sentence "A `held` refresh keeps an installed candidate while its declared build is open." | planned TestCandidateInstallAnchors in internal/anchors/registry_candidate_install_test.go | Removal of the effect note reds the anchor |

Not covered: story 38 — review-owned exclusion; FT140 owns the tolerant purge read.
Not covered: story 39 — review-owned exclusion; `docs/known-issues.md` states the one-version-at-a-time rule.

### Edge inventory

The pre-landing writers and the landing binary of the declared fields of the 2026-10-12 build are as follows:

- `intent-ledger:commitment.gaps`: the pre-landing writers are the gap close of `bench record assignment` at the first dispatch and the chunk checkpoint's gap open. Both run in the caller's own process.
- `review-record:folds`: the pre-landing writer is `bench record fold`, which writes the tree-local review record in the build worktree.
- The landing binary is the installed candidate, which the manifest binds through the unchanged land route.

This spec adds no persisted field of its own, so its own build declares none.

The audience is this repository only. A linked repository declares no Bench build inputs, so CB16 refuses the install there, and the declaration stays inert.

Hostile inputs of the declaration attach to CB4–CB7: an empty item, an unknown record name, an empty segment, and a duplicate item. A path through a non-struct attaches to CB11.

The edges and their dispositions are:

- A worktree working directory selects the worktree `dist/bench`. CB21 refuses that executable, and CB15 proves the manifest path.
- The version cache precedes the primary `dist/bench` in the wrapper. CB22 refuses the install, and CB36 refuses a later checkpoint.
- A Go commit after the install leaves the installed candidate stale, as FT406's known issue states. A stale candidate cannot write a field that it does not know. After a rebuild of the worktree copy, CB27 refuses the new digest, and CB20 reinstalls.
- The owner lookup fails open to "no owner" on a strict-decode refusal. CB32 places the check before that lookup, and CB34 pins the route without it.
- An open build between chunks does not refuse other landings. CB42 keys the restore guard to the open build, not to the window refusal.
- The post-landing refresh rebuilds the primary executable when it does not verify. An installed candidate never verifies against the primary sources. CB50–CB53 hold the candidate during the build, and CB25 proves the rebuild after the build's own landing.
- An unrelated landing or a session-start rebuild advice can still leave a `main` executable beside a ledger with the new field. CB54 lets the install recover from that state.
- Before the install, the `main` executable passes the digest check. It cannot write a field that it does not know, and CB46 owns the install step.
- A refusal in an author worktree names the integration target. CB57 pins that label.

**Won't handle** lines:

- An older executable's tolerant purge read drops an unknown commitment field on its rewrite — FT140 owns the tolerant read's schema guard, and the installed candidate's own purge keeps the field.
- Another worktree's own older `dist/bench` refuses the ledger after the install — the one-version-at-a-time rule in `docs/known-issues.md` governs it, and the installed candidate serves the primary checkout and the build worktree.
- A strict-decode refusal in the owner lookup of `bench commit`, preflight, the handoff, or the tree target under an older executable — the install removes the older executable from the build's own path, and CB32 covers the checkpoint.
- Shared gate evidence under the common directory with a new field — its reader fails soft to an invalid cache entry, and the gate then runs again.
- A candidate that changes the package version — the land route's version check refuses it, and a version change belongs to the release preparation, not to a build.
- A candidate tree that removes its own `Persisted fields:` line — the removal only disables the writer guard, the install still reads the default-branch copy, and review grades the spec edit.
- A landing of a declared build on an uninstalled broker — it keeps today's strict-decode refusal, and CB30 refuses the complete checkpoint that comes before the landing.
- The `commands --brief` probe, the doctor seal row, and the preflight seal row report the installed candidate as stale — a rebuild that obeys them uninstalls the candidate, and CB27 then refuses the next declared write.
- A restore after the operator already rebuilt `main` by hand — CB41 refuses it, and an install (CB54) followed by the restore is the route.
- A reinstall while a chunk checkpoint runs — the orchestrator runs the install and the checkpoint in sequence in one session, and the check runs at the checkpoint's entry only.
- The preflight `binary-seal` row names a plain `bench worktree build` as its remedy — after that build, the next guarded entry refuses and names the install, and FT344 owns the remedy wording.
- A malformed `Persisted fields:` line passes the spec stage — the install refuses it with `persisted-fields-invalid`, and a preflight row for the line is out of scope.

## Out of scope

- FT354 strict-read ownership of the review-record decoder: 0 edits, 0 gate runs here; FT354 owns its staged spec.
- FT406 rebuild on commit: 0 edits, 0 gate runs here; CB20 covers a manual reinstall until FT406 lands.
- FT414 precondition projection of the install obligation in the record and landing verbs: about 6 edits, 2 gate runs as its own capability.
- FT424 chunk-gap landing of waiting outcomes: 0 edits here.
- A `bench preflight build` row for the declaration: about 5 edits, 2 gate runs as its own capability, because the row-list assertions need their own sweep.
- A `bench status` row that names a pending restore: about 4 edits, 1 gate run as its own capability, because it needs an exact installed-candidate predicate.
- A schema guard on the tolerant purge read: FT140 owns it; about 3 edits, 1 gate run.

## Further notes

### Pre-review proof checklist

- Cited symbols: `Address`, `applyCheckpoint`, `land_route`, `LandCommand`, `WriteBrokerManifest`, `refreshBroker`, `brokerChangeNotice`, `CloseGap`, `openBuilds`, `metadata`, and `Facts` resolve in the files named above.
- Import edges: `internal/spec` imports neither record package; each record owner adds only `reflect`. The gate, record, and worktree packages already import `internal/spec`, `internal/intent`, and `internal/reviewrecord` (uncertain until `go list` runs at slicing).
- Source-row clauses and occurrences: the source-to-row table below covers each clause.
- Promised field labels: `Persisted fields:`, `intent-ledger`, `review-record`, `candidate-not-installed`, and `persisted-fields-invalid`.
- Changed-function callers: `metadata` serves `Facts` and `AwaitsRetirement`; `buildWith` serves `BuildCommand`; `applyCheckpoint` serves the checkpoint and the prospective gate; the assignment `BeforeWrite` serves `bench record assignment`. The slicer runs `bench consumers` for each.
- Copy survival: none.
- Rendered-shape readers: the build verb's grammar gains `--install` and `--restore`; the slicer names every help-inventory fixture and usage test of `bench worktree build`. The landing's `effects` row gains the result `held`; the slicer names every reader of the effect result set, including `.bench/BENCH-reference.md`.
- New result labels: `held`.
- Pin operators: digest rows compare exact SHA-256 strings; path rows compare canonical physical paths; refusal rows compare the exact face name.
- Entry reads: each writer entry reads the primary manifest and the running executable; the check takes both values as injected arguments behind its entry.
- Derived expectations: the field predicates answer from the struct tags, and the tests assert literal paths that the test author writes.
- Consolidated rules: one candidate check serves the three writer entries.
- Quantified obligations: "each declared field that the ledger holds" applies at CB38, CB39, and CB49.
- Workflow-step writes: the install writes the primary executable, its seal, the manifest, the shim, and the pre-push hook; the restore writes the ledger copy, the ledger, and the same five paths. No checkpoint digest reads these paths.

### Source-to-row table

| source clause | occurrence | rows or disposition |
| --- | --- | --- |
| A build that adds a persisted ledger or review-record field writes and lands through a binary that reads that field | `roadmap/FT419.md:1` | CB13–CB15, CB24, CB26–CB31 |
| The strict decoder stays | `roadmap/FT419.md:8` | No decoder change; CB8–CB12 read the same struct tags |
| The checkpoints and the landing of a kit build run through the build's own candidate binary | `roadmap/FT419.md:8-9` | CB15, CB23, CB24, CB29, CB30 |
| The spec names the pre-landing writer and the landing binary in its edge inventory | `roadmap/FT419.md:10-11` | Edge inventory above; CB45 |
| BW-D checkpoint wrote `commitment.gaps`, and the `main` binary refused the ledger | `roadmap/FT419.md:16` | CB23, CB29, CB31, CB32 |
| The installed binary refused the review-record `folds` field | `roadmap/FT419.md:17` | CB24, CB30 |
| Install route through the authenticated `doctor --fix` route; the wrapper stays unchanged | conversation 2026-10-20 | CB15, CB21, CB24, CB25 |
| ADR 0016 gets one narrow exception | conversation 2026-10-20 | CB44 |
| The post-landing refresh holds an installed candidate while its declared build is open | conversation 2026-10-20 | CB50–CB53, CB59 |

### Flagged additions and residuals

These items go beyond the literal decision source, and each one serves the decided route:

- The `Persisted fields:` header line makes the obligation machine-visible.
- The `candidate-not-installed` face turns a late decode error into an early named refusal at each ledger writer.
- The restore verb handles a build that stops before its landing.

Flag for reviewer veto: the 2026-10-20 answer placed the install before the build's first checkpoint. The gap close of `bench record assignment` writes the ledger at the first dispatch, before any checkpoint. This spec therefore places the install before the first ticket dispatch, which follows the source's "writes through a binary that reads that field". The route is unchanged.

Flag for reviewer veto: `projects/benchkit.md` claims that the landing refuses a stale dev executable, and the current code has no such refusal. This spec does not change the claim.

Flag for reviewer veto: the install runs candidate code as the `doctor --fix` child before any landing. That child also writes the PATH shim and the pre-push hook. While the build is open, the installed candidate also serves every session's hooks, including the fail-closed git guard, and it brokers every unrelated landing. This is the trust assumption that the install route accepts, and the reviewer's declaration sign-off is its root.

This planning pass ran `bench coverage --check` only. It supplies no implementation, runtime, mutation, or gate evidence.
