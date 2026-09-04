# Repair the review findings

Blocked by: wire-the-guard-facts-in-the-subcommand.md
Writes: internal/gitguard/scan.go, internal/gitguard/verdict.go, internal/gitguard/gitguard.go, internal/gitguard/verdict_test.go, internal/gitguard/gitguard_test.go, internal/gitguard/checker_junction_test.go, internal/git/push_destination_test.go, internal/gittest/gittest.go, cmd/bench/main_test.go, specs/agent-push-guard/spec.md, tests/canary/canonical-path-owner/second-derivation, tests/canary/injected-ports/unregistered-port, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: PG28, PG42, PG43, PG44, PG45, PG46

## What to build

This ticket records the repairs the review round of 2026-09-04 accepted. The
review pickup `reviews/agent-push-guard.md` holds the findings and their
dispositions.

A push that carries a global `-C`, `--git-dir`, or `--work-tree` denies with
the unresolved label. The push facts describe the process working directory,
not the redirected repository. A push under an `xargs` prefix
denies with the same label, because its destination arrives on stdin. The
destination rule reads `@` as `HEAD` and strips a `heads/` prefix the way it
strips `refs/heads/`. The unresolved advice gains a first clause that names
the redirect fix, and the message still ends with the sentence row PG28 pins.

One shared builder, `gittest.TopicTrackingRepo`, owns the topic-tracking
fixture that three test packages built by hand. Row PG28 names the test that
holds its assertion.

## Acceptance

- [ ] `git -C /other push origin topic`, `git --git-dir=/x push origin topic`, and `git --work-tree=/x push origin topic` each return the unresolved label.
- [ ] `git -C /tmp reset --hard` keeps the reset label.
- [ ] `xargs git push` returns the unresolved label.
- [ ] `git push origin @` returns the default-branch label on a checked-out `main` and the allow verdict on a checked-out `topic`.
- [ ] `git push origin heads/main` returns the default-branch label, and `git push origin heads/topic` returns the allow verdict.
- [ ] The unresolved block message ends with `Name the remote and the branch: git push <remote> <branch>.`
- [ ] No test package builds the topic-tracking repository by hand.
- [ ] Self-probe: omit the redirect check, and report the three redirect rows red.
