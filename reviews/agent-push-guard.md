# Review pickup: agent-push-guard

Frozen base `5044b03dedd6a598fea6e6dc885f761e02def193`, reviewed tip
`063223038c07e8f72652a1277db0c8ccf790a3af`. Three axes ran at `opus` low.
Raw findings: 11. Repair targets after collapse: 7.

## Standards

Count: 3. Worst: two hand-rolled repository fixtures duplicate the scaffold
that `internal/gittest` owns.

- `internal/gitguard/checker_junction_test.go` (`gcTopicRepo`) and `internal/git/push_destination_test.go` (`newTopicRepo`) build the same topic-tracking repository. AGENTS.md names a fixture harness pasted N times as the defect. Disposition: auto-fix. Move one builder into `internal/gittest` and call it from both packages.
- `cmd/bench/main_test.go` (`guardPushRepo`) repeats the init and identity calls that `gittest.RepoOnBranch` owns. Disposition: auto-fix. Compose `RepoOnBranch`, then keep only the commit and the topic checkout.
- `internal/gitguard/checker_junction_test.go` (`realChecker`) derives the guard wiring a second time beside `guardGit`. A collapse needs an exported constructor and a new import edge from `internal/gitguard` to `internal/git`, which the spec excludes. Disposition: ask-user. Reported, not repaired.

## Spec

Count: 4. Worst: a push with a global `-C`, `--git-dir`, or `--work-tree` is graded against the current directory's repository.

- `cmd/bench/main.go` (`guardProbeRoot`) and `internal/gitguard/scan.go` (`FindSubcommand`): the scanner strips the redirect, and the facts read the current directory. `git -C /other push origin main` allows when the current default differs. Disposition: auto-fix under the fail-closed rule. A redirected push denies as unresolved, and the advice sentence gains a first clause. Recorded for reviewer veto.
- `specs/agent-push-guard/spec.md` row PG28 names `TestBlockMessageNamesLabel`; the assertion lives in `TestBlockMessageCarriesUnresolvedAdvice`. Disposition: auto-fix. Amend the seam column.
- `--force-if-includes` is allowed by fall-through and pinned by a test row, but no spec line decides it. Disposition: no-op in code. Add one Won't-handle line.
- The `git push main` row pins an edge the spec lists as Won't handle. Disposition: no-op in code. Move the edge into the inventory as handled.

## Coverage

Count: 4. Worst: `xargs git push` reaches the bare rule with no args and allows a push whose destination arrives on stdin.

- `xargs git push` with a checked-out topic returns allow. `checkoutVerdict` and `restoreVerdict` deny on `viaXargs`; `pushVerdict` never reads the flag. Disposition: auto-fix. Deny as unresolved, with a row beside the existing xargs rows.
- `git push origin @` on a checked-out `main` returns allow. Git reads `@` as `HEAD`. Disposition: auto-fix. Read `@` as `HEAD` in `pushDestination`, with a row.
- `git push origin heads/main` returns allow. Git reads the `heads/` prefix as `refs/heads/`. Disposition: auto-fix. Strip `heads/` as well, with a row.
- `git -C /other push origin topic` is the Spec finding above. Disposition: collapsed into that repair target.
